package store

import (
	"bytes"
	"fmt"
	"testing"
)

func TestTxnStart(t *testing.T) {
	st := newStore(nil)
	txn := st.NewTransaction(true)
	err := txn.Set([]byte("foo"), []byte("bar"))
	if err != nil {
		t.Fatalf("unable to insert data: %v", err)
	}
	txn.Commit()
	result := txn.Get([]byte("foo"))
	if result == nil {
		t.Fatal("result value is empty")
	}
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
		if err != nil {
			t.Fatalf("unable to get value: %v", err)
		}
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
		fmt.Println(itm)
	}
}
