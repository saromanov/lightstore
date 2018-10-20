package store

//https://docs.oracle.com/cd/E17275_01/html/api_reference/C/txn.html

// write defines writing to store
type write struct {
	key []byte
}

// Txn represents transaction
type Txn struct {
	writes []write
	count  int64
	id     int64
	reads  []int64
}

// Entry defines new key value pair
type Entry struct {
	key   []byte
	value []byte
}

// NewTransaction creates a new transaction
func (s *Store) NewTransaction() *Txn {
	txn := &Txn{}
	return txn
}
