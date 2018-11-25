package store

import (
	"testing"
)

func TestTxnStart(t *testing.T) {
	txn := NewTransaction(true)
	txn.Set([]byte("foo"), []byte("bar"))
}
