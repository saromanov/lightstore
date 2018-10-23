package store

import "errors"

var (
	errNoWrites     = errors.New("unable to write on read-only mode")
	errEmptyKey     = errors.New("key is empty")
	errLargeKeySize = errors.New("key size is larger then limit")
)

const maxKeySize = 16384

//https://docs.oracle.com/cd/E17275_01/html/api_reference/C/txn.html

// Txn represents transaction
type Txn struct {
	writes []*Entry
	count  int64
	id     int64
	reads  []int64
	write  bool
	store  *Store
}

// Entry defines new key value pair
type Entry struct {
	key    []byte
	value  []byte
	expire uint64
}

// pendingWritesIterator provides iteration over
// pednding writes
type pendingWritesIterator struct {
	entries []*Entry
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
		store:  s,
		writes: []*Entry{},
		reads:  []int64{},
		write:  write,
	}
	return txn
}

// DB retruns reference to store
func (t *Txn) DB() *Store {
	return t.store
}

// Close defines ending of transaction
func (t *Txn) Close() {
	t.store = nil
}

// Delete provides removing value by the key
func (t *Txn) Delete(key []byte) error {
	entry := &Entry{
		key: key,
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
	for _, w := range t.writes {
		t.store.Set(w.key, w.value)
	}
	return nil
}

// Set writes a new key value pair to the pending writes
// It'll be applying after transaction
func (t *Txn) Set(key, value []byte) error {
	entry := &Entry{
		key:   key,
		value: value,
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
	if len(entry.key) > maxKeySize {
		return errLargeKeySize
	}
	return nil
}
