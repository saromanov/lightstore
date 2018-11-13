package store

// Item defines struct for item on iterator
type Item struct {
	key   []byte
	value []byte
	err   error
	next  *Item
	txn   *Txn
}

// Key returns key of the item
func (i *Item) Key() []byte {
	return i.key
}

func copy(f, s []byte) []byte {
	return append(f[:0], s...)
}
