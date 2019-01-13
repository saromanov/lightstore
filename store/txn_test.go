package store

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxnStart(t *testing.T) {
	st := newStore(nil)
	defer st.Close()
	txn := st.NewTransaction(true)
	err := txn.Set([]byte("foo"), []byte("bar"))
	if err != nil {
		t.Fatalf("unable to insert data: %v", err)
	}
	err = txn.Commit()
	if err != nil {
		t.Fatalf("unable to commit: %v", err)
	}
	txn2 := st.NewTransaction(false)
	result := txn2.Get([]byte("foo"))
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
	for i := 0; i < 10; i++ {
		err := txn.Set([]byte(fmt.Sprintf("%d", i)), []byte(fmt.Sprintf("bar+%d", i)))
		if err != nil {
			t.Fatal(err)
		}
	}
	txn.Commit()
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

func TestIteratorWithNoSize(t *testing.T) {
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
	for i := 0; i < 10; i++ {
		err := txn.Set([]byte(fmt.Sprintf("%d", i)), []byte(fmt.Sprintf("bar+%d", i)))
		if err != nil {
			t.Fatal(err)
		}
	}
	txn.Commit()
	it := txn.NewIterator(IteratorOptions{
		Size: 4,
	})
	count := 0
	for it.First(); it.Valid(); it.Next() {
		itm := it.Item()
		err := itm.Value(func(v []byte) error {
			return nil
		})
		if err != nil {
			t.Fatalf("unable to get value: %v", err)
		}
		count++
	}
	assert.Equal(t, count, 4, "Count of return elements is not equal")
}
