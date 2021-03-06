package store

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTxnStart(t *testing.T) {
	st, err := newStore(nil)
	if err != nil {
		t.Fatalf("unable to init store: %v", err)
	}
	defer st.Close()
	txn := st.NewTransaction(true)
	err = txn.Set([]byte("foo"), []byte("bar"))
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
	if !bytes.Equal(result, []byte("bar")) {
		t.Fatal("unable to get result")
	}
}

func TestIteratorWithNoOptions(t *testing.T) {
	st, err := newStore(nil)
	if err != nil {
		t.Fatalf("unable to init store: %v", err)
	}
	txn := st.NewTransaction(true)
	for i := 0; i < 10; i++ {
		err := txn.Set([]byte(fmt.Sprintf("%d", i)), []byte(fmt.Sprintf("bar+%d", i)))
		if err != nil {
			t.Fatal(err)
		}
	}
	txn.Commit()
	txn2 := st.NewTransaction(false)
	it, _ := txn2.NewIterator(IteratorOptions{})
	for it.First(); it.Valid(); it.Next() {
		itm := it.Item()
		val := itm.Value()
		if val == nil {
			t.Fatal("unable to get value")
		}
	}
}

func TestIteratorWithClosedTransaction(t *testing.T) {
	st, err := newStore(nil)
	if err != nil {
		t.Fatalf("unable to init store: %v", err)
	}
	txn := st.NewTransaction(true)
	for i := 0; i < 10; i++ {
		err := txn.Set([]byte(fmt.Sprintf("%d", i)), []byte(fmt.Sprintf("bar+%d", i)))
		if err != nil {
			t.Fatal(err)
		}
	}
	txn.Commit()
	_, err = txn.NewIterator(IteratorOptions{})
	assert.EqualError(t, err, errTransactionCompleted.Error(), "error response is not equal")
}

func TestIteratorWithNoSize(t *testing.T) {
	st, err := newStore(nil)
	if err != nil {
		t.Fatalf("unable to init store: %v", err)
	}
	txn := st.NewTransaction(true)
	err = txn.Set([]byte("foo"), []byte("bar"))
	if err != nil {
		t.Fatal("unable to insert data")
	}
	txn.Commit()
	txn2 := st.NewTransaction(true)
	it, _ := txn2.NewIterator(IteratorOptions{})
	count := 0
	for it.First(); it.Valid(); it.Next() {
		count++
	}
	assert.Equal(t, count, 1, "Count of return elements is not equal")
}

func TestIteratorWithLimit(t *testing.T) {
	st, err := newStore(nil)
	if err != nil {
		t.Fatalf("unable to init store: %v", err)
	}
	txn := st.NewTransaction(true)
	for i := 0; i < 10; i++ {
		err := txn.Set([]byte(fmt.Sprintf("%d", i)), []byte(fmt.Sprintf("bar+%d", i)))
		if err != nil {
			t.Fatal(err)
		}
	}
	txn.Commit()
	txn2 := st.NewTransaction(true)
	it, _ := txn2.NewIterator(IteratorOptions{
		Limit: 4,
	})
	count := 0
	for it.First(); it.Valid(); it.Next() {
		count++
	}
	assert.Equal(t, count, 4, "Count of return elements is not equal")
}
