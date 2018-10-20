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
	key   []byte
	value []byte
}

// NewTransaction creates a new transaction
func (s *Store) NewTransaction() *Txn {
	txn := &Txn{
		store: s,
	}
	return txn
}

// Commit applies a new commit after modification
func (t *Txn) Commit() error {
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
