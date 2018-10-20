package store

import "errors"

var errNoWrites = errors.New("unable to write on read-only mode")

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

// NewTransaction creates a new transaction
func (s *Store) NewTransaction() *Txn {
	txn := &Txn{
		store:  s,
		writes: []*Entry{},
		reads:  []int64{},
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
	if !t.write {
		return errNoWrites
	}
	entry := &Entry{
		key:   key,
		value: value,
	}
	return t.set(entry)
}

func (t *Txn) set(entry *Entry) error {
	t.writes = append(t.writes, entry)
	return nil
}

// beforeSet checks properties before set of the entry
func (t *Txn) beforeSet(entry *Entry) error {
	if !t.write {
		return errNoWrites
	}
	return nil
}
