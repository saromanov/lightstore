package store

import "fmt"

// Item defines struct for item on iterator
type Item struct {
	key   []byte
	value []byte
	err   error
	next  *Item
	txn   *Txn
}

// String retruns string representation of Item
func (i *Item) String() string {
	return fmt.Sprintf("%q", i.Key())
}

// Iterator provides iterating over the KV pairs
type Iterator struct {
	txn     *Txn
	lastKey []byte
	closed  bool
}

// Key returns key of the item
func (i *Item) Key() []byte {
	return i.key
}

// Copykey provides copy of the current key
func (i *Item) CopyKey(src []byte) []byte {
	return copy(i.key, src)
}

func copy(f, s []byte) []byte {
	return append(f[:0], s...)
}
