package backend

import "github.com/google/btree"

// BTree provides implementation of btree
// over google btree package
type BTree struct {
	repair *Repair
	engine *btree.BTree
}

// InitBTree provides initialization of the B-Tree
func InitBTree(degree int) *BTree {
	bt := &BTree{
		repair: NewRepair(),
		engine: btree.New(degree),
	}
	return bt
}

// Put provides store of the key and value
func (d *BTree) Put(key []byte, value interface{}, op ItemOptions) error {
	if op.Immutable {
		return nil
	}
	return nil
}

// Get provides getting of value by the key
func (d *BTree) Get(key []byte) (interface{}, error) {
	/*value, ok := d.Value.Get(key)
	if !ok {
		return nil, errors.New("unable to find element")
	}
	return value, nil*/
	return nil, nil
}

// Exist provides implementation for checking of key is exist
func (d *BTree) Exist(key []byte) bool {
	return false
}

// FIrst retruns min element from storage
func (d *BTree) First() interface{} {
	return d.engine.Min()
}

// Remove provides removing of the record
func (d *BTree) Delete(key []byte) error {
	return nil
}

// Next provides iteration to the next element of storage
func (d *BTree) Next(i int) interface{} {
	return nil
}
