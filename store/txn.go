package store

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/satori/go.uuid"
)

var errTransactionCompleted = errors.New("transaction was completed")

//https://docs.oracle.com/cd/E17275_01/html/api_reference/C/txn.html

// Txn represents transaction
type Txn struct {
	writes    EntrySlice
	count     int64
	id        string
	reads     EntrySlice
	write     bool
	store     *Store
	timestamp int64
	complete  bool
}

// Entry defines new key value pair
type Entry struct {
	key       []byte
	value     []byte
	expire    uint64
	timestamp int64
	isDelete  bool
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
		reads:     []*Entry{},
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
	t.close()
}

// Delete provides removing value by the key
func (t *Txn) Delete(key []byte) error {
	entry := &Entry{
		key:       key,
		timestamp: time.Now().Unix(),
		isDelete:  true,
	}
	if err := t.beforeSet(entry); err != nil {
		return err
	}
	t.writes = append(t.writes, entry)
	return nil
}

// Commit applies a new commit after modification
func (t *Txn) Commit() error {
	if t.store == nil {
		return errNoStorage
	}
	if len(t.writes) == 0 {
		return nil
	}
	sort.Sort(t.writes)
	t.handleCommit(t.writes, func(e *Entry) {
		if e.isDelete {
			t.store.Remove(e.key)
		}
		fmt.Println("EEE: ", t.store.Set(e.key, e.value))
	})
	t.Close()
	return nil
}

func (t *Txn) handleCommit(writes EntrySlice, action func(*Entry)) {
	for _, w := range t.writes {
		action(w)
	}
}

// close provides closing of transaction
func (t *Txn) close() {
	t.store = nil
	t.writes = t.writes[:0]
	t.reads = t.reads[:0]
	t.complete = true
}

// Rollback provides rolling back current
// pending transactions
func (t *Txn) Rollback() {
	t.rollback()
}

func (t *Txn) rollback() {
	t.writes = t.writes[:0]
}

// NewIterator creates a new iterator under transaction
func (t *Txn) NewIterator(opt IteratorOptions) (*Iterator, error) {
	if t.complete {
		return nil, errTransactionCompleted
	}
	return &Iterator{
		txn:   t,
		opt:   opt,
		limit: opt.Limit,
		item: &Item{
			key: t.store.first(),
		},
	}, nil
}

// Get returns value by the key
func (t *Txn) Get(key []byte) []byte {
	entry := &Entry{
		key: key,
	}
	t.reads = append(t.reads, entry)
	return t.store.Get(key)
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

// SetWithTTL provides setting of key with Time To Live(TTL)
func (t *Txn) SetWithTTL(key, value []byte, duration time.Time) error {
	entry := &Entry{
		key:       key,
		value:     value,
		timestamp: time.Now().Unix(),
		expire:    uint64(duration.Unix()),
	}
	return t.set(entry)
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
	if uint(len(entry.key)) > t.store.config.MaxKeySize {
		return errLargeKeySize
	}
	return nil
}

// beforeGet provides providing of operations before output
// of key to get
func (t *Txn) beforeGet(entry *Entry) error {
	return nil
}
