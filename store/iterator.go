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

// Value returns current value of Item
func (i *Item) Value(f func([]byte) error) error {
	return nil
}

// IteratorOptions defines options
// for iterator
type IteratorOptions struct {
	Prefix []byte
	Size   uint
	Limit  uint
}

// Iterator provides iterating over the KV pairs
type Iterator struct {
	txn     *Txn
	lastKey []byte
	closed  bool
	item    *Item
	opt     IteratorOptions
	element int
	limit   uint
}

// Item retruns current item from iterator
func (it *Iterator) Item() *Item {
	return it.item
}

// First defines start point for iteration
func (it *Iterator) First() {
	it.lastKey = it.lastKey[:0]
	it.item = it.makeItem()
}

func (it *Iterator) makeItem() *Item {
	store := it.txn.store
	key := store.first()
	return &Item{
		key:   key,
		value: store.Get(key),
	}
}

// Valid returns false if current item is invalid
func (it *Iterator) Valid() bool {
	return it.item != nil && it.item.key != nil
}

// Next provides getting of the next element
// on iterator
func (it *Iterator) Next() *Item {
	it.element++
	if it.limit > 0 && it.element > int(it.limit) {
		return nil
	}
	if it.txn == nil {
		return nil
	}
	key := it.txn.store.next(it.element)
	it.item = &Item{
		key:   key,
		value: it.txn.store.Get(key),
	}
	return it.item
}

// Close provides closing of iterator
func (it *Iterator) Close() {
	it.txn = nil
	it.item = nil
}

// Key returns key of the item
func (i *Item) Key() []byte {
	return i.key
}

// CopyKey provides copy of the current key
func (i *Item) CopyKey(src []byte) []byte {
	return copy(i.key, src)
}

func copy(f, s []byte) []byte {
	return append(f[:0], s...)
}
