package datastructures

import (
	"errors"

	"github.com/google/btree"
)

// BTree provides implementation of btree
// over google btree package
type BTree struct {
	repair *Repair
	Value  *btree.BTree
}

func InitBTree(degree int) *BTree {
	bt := &BTree{
		repair: NewRepair(),
		tree:   btree.New(degree),
	}
	return bt
}

func (d *BTree) Put(key []byte, value interface{}, op ItemOptions) error {
	if op.Immutable {
		return nil
	}
	d.Value.Set(key, value)
	return nil
}

// Get provides getting of value by the key
func (d *BTree) Get(key []byte) (interface{}, error) {
	value, ok := d.Value.Get(key)
	if !ok {
		return nil, errors.New("unable to find element")
	}
	return value, nil
}

func (d *BTree) GetFromRepair(key string) (*RepairItem, error) {
	return d.repair.GetFromRepair(key, "")
}

// Exist provides implementation for checking of key is exist
func (d *BTree) Exist(key []byte) bool {
	_, ok := d.Value.Get(key)
	if !ok {
		return false
	}
	return true
}

// Remove provides removing of the record
func (d *SkipList) Delete(key []byte) error {
	_, ok := d.Value.Delete(key)
	if !ok {
		return errors.New("unable to delete element")
	}
	return nil
}
