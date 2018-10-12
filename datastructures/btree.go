package datastructures

import "github.com/google/btree"

// BTree provides implementation of btree
// over google btree package
type BTree struct {
	repair *Repair
	tree   *btree.BTree
}

func InitBTree(degree int) *BTree {
	bt := &BTree{
		repair: NewRepair(),
		tree:   btree.New(degree),
	}
	return bt
}
