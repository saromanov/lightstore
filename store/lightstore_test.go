package store

import (
	"testing"
)

func TestOpenLightstore(t *testing.T) {
	light := Open(nil)
	defer light.Close()
	if !light.IsCreated() {
		t.Fatalf("unable to create db")
	}
}

func TestSetData(t *testing.T) {
	light := Open(nil)
	defer light.Close()
	err := light.Write(func(txn *Txn) error {
		err := txn.Set([]byte("foo"), []byte("bar"))
		if err != nil {
			return err
		}
		return txn.Commit()
	})
	if err != nil {
		t.Fatalf("unable to write data: %v", err)
	}
}
