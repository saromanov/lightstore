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
	err = txn.Commit()
	if err != nil {
		t.Fatal("unable to apply commit")
	}
	result := txn.Get([]byte("foo"))
	if bytes.Compare(result, []byte("bar")) != 0 {
		t.Fatal("unable to get result")
	}
}
