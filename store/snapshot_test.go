package store

import (
	"testing"
)

func TestSnapshot(t *testing.T) {
	light, err := Open(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer light.Close()
	err = light.Write(func(txn *Txn) error {
		err := txn.Set([]byte("foo"), []byte("bar"))
		if err != nil {
			return err
		}
		return txn.Commit()
	})
	if err != nil {
		t.Fatalf("unable to write data: %v", err)
	}

	snap := NewSnapshot(light.getStore(), "./snapshot1")
	err = snap.Write(nil)
	if err != nil {
		t.Fatalf("unable to write snapshot: %v", err)
	}
}
