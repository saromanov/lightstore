package store

import (
	"bytes"
	"testing"
)

func TestTxnStart(t *testing.T) {
	st := newStore(nil)
	txn := st.NewTransaction(true)
	err := txn.Set([]byte("foo"), []byte("bar"))
	if err != nil {
		t.Fatal("unable to insert data")
	}
	result := txn.Get([]byte("foo"))
	if bytes.Compare(result, []byte("bar")) != 0 {
		t.Fatal("unable to get result")
	}
}

func TestIteratorWithNoOptions(t *testing.T) {
	st := newStore(nil)
	txn := st.NewTransaction(true)
	err := txn.Set([]byte("foo"), []byte("bar"))
	if err != nil {
		t.Fatal("unable to insert data")
	}
	it := txn.NewIterator(IteratorOptions{})
	for it.First(); it.Valid(); it.Next() {
		itm := it.Item()
		err := itm.Value(func(v []byte) error {
			return nil
		})
	}
}

func TestIteratorWithSize(t *testing.T) {
	st := newStore(nil)
	txn := st.NewTransaction(true)
	err := txn.Set([]byte("foo"), []byte("barfoo"))
	if err != nil {
		t.Fatal("unable to insert data")
	}
	it := txn.NewIterator(IteratorOptions{
		Size: 4,
	})
	for it.First(); it.Valid(); it.Next() {
		itm := it.Item()
	}
}
