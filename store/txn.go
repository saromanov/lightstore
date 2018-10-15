package store

// write defines writing to store
type write struct {
	key []byte
}

// Txn represents transaction
type Txn struct {
	writes []write
	count  int64
	id     int64
}
