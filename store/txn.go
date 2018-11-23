package store

import (
	"errors"
	"sort"
	"time"

	"github.com/satori/go.uuid"
)

var (
	errNoWrites     = errors.New("unable to write on read-only mode")
	errEmptyKey     = errors.New("key is empty")
	errLargeKeySize = errors.New("key size is larger then limit")
	errNoStorage    = errors.New("storage is not defined")
)

const maxKeySize = 16384

//https://docs.oracle.com/cd/E17275_01/html/api_reference/C/txn.html

// Txn represents transaction
type Txn struct {
	writes    EntrySlice
	count     int64
	id        string
	reads     []int64
	write     bool
	store     *Store
	timestamp int64
}

// Entry defines new key value pair
type Entry struct {
	key       []byte
	value     []byte
	expire    uint64
	timestamp int64
}

type EntrySlice []*Entry

func (a EntrySlice) Len() int           { return len(a) }
func (a EntrySlice) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a EntrySlice) Less(i, j int) bool { return a[i].timestamp < a[j].timestamp }

// pendingWritesIterator provides iteration over
// pednding writes
type pendingWritesIterator struct {
	entries EntrySlice
	nextIdx int
	readTs  uint64
}

// Start is moving iteration to the start
func (p *pendingWritesIterator) Start() {
	p.nextIdx = 0
}

// Next moves to the second element of entries
func (p *pendingWritesIterator) Next() {
	p.nextIdx++
}

// GetKey returns current key on iterator
func (p *pendingWritesIterator) GetKey() []byte {
	return p.entries[p.nextIdx].key
}

// NewTransaction creates a new transaction
func (s *Store) NewTransaction(write bool) *Txn {
	txn := &Txn{
		id:        uuid.Must(uuid.NewV4()).String(),
		store:     s,
		writes:    []*Entry{},
		reads:     []int64{},
		write:     write,
		timestamp: time.Now().UnixNano(),
	}
	return txn
}

// Gettimestamp returns timestamp of current transaction
func (t *Txn) GetTimestamp() int64 {
	return t.timestamp
}

// ID returns id of transaction
func (t *Txn) ID() string {
	return t.id
}

// DB retruns reference to store
func (t *Txn) DB() *Store {
	return t.store
}

// Close defines ending of transaction
func (t *Txn) Close() {
	t.store = nil
	t.writes = []*Entry{}
	t.reads = []int64{}
}

// Delete provides removing value by the key
func (t *Txn) Delete(key []byte) error {
	entry := &Entry{
		key:       key,
		timestamp: time.Now().Unix(),
	}
	if err := t.beforeSet(entry); err != nil {
		return err
	}
	t.writes = append(t.writes, entry)
	return nil
}

// Commit applies a new commit after modification
func (t *Txn) Commit() error {
	if len(t.writes) == 0 {
		return nil
	}
	if t.store == nil {
		return errNoStorage
	}
	sort.Sort(t.writes)
	for _, w := range t.writes {
		t.store.Set(w.key, w.value)
	}
	return nil
}

// Set writes a new key value pair to the pending writes
// It'll be applying after transaction
func (t *Txn) Set(key, value []byte) error {

	entry := &Entry{
		key:       key,
		value:     value,
		timestamp: time.Now().Unix(),
	}
	return t.set(entry)
}

// Get returns value by the key
func (t *Txn) Get(key []byte) []byte {
	return t.store.Get(key)
}

func (t *Txn) set(entry *Entry) error {
	if err := t.beforeSet(entry); err != nil {
		return err
	}
	t.writes = append(t.writes, entry)
	return nil
}

// beforeSet checks properties before set of the entry
func (t *Txn) beforeSet(entry *Entry) error {
	if !t.write {
		return errNoWrites
	}
	if len(entry.key) == 0 {
		return errEmptyKey
	}
	if len(entry.key) > maxKeySize {
		return errLargeKeySize
	}
	return nil
}
