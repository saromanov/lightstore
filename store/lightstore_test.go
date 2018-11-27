package store

import (
	"fmt"
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

func TestTwoCommits(t *testing.T) {
	light := Open(nil)
	defer light.Close()
	err := light.Write(func(txn *Txn) error {
		err := txn.Set([]byte("foo"), []byte("bar"))
		if err != nil {
			return err
		}

		if err := txn.Commit(); err != nil {
			return err
		}
		if err := txn.Commit(); err == nil {
			return fmt.Errorf("should return error on second commit")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unable to write data: %v", err)
	}
}

func TestWriteOnReadOnly(t *testing.T) {
	light := Open(nil)
	defer light.Close()
	err := light.View(func(txn *Txn) error {
		err := txn.Set([]byte("foo"), []byte("bar"))
		if err == nil {
			return fmt.Errorf("unable to write on read-only transaction")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unable to write data: %v", err)
	}
}

func TestGetData(t *testing.T) {
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

	err = light.View(func(txn *Txn) error {
		data := txn.Get([]byte("foo"))
		if string(data) != "bar" {
			return fmt.Errorf("unable to get data")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unable to get: %v", err)
	}
}

func TestGetNotFoundData(t *testing.T) {
	light := Open(nil)
	defer light.Close()

	err := light.View(func(txn *Txn) error {
		data := txn.Get([]byte("foo"))
		if len(data) != 0 {
			return fmt.Errorf("expecting empty result")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unable to get: %v", err)
	}
}

func TestViewData(t *testing.T) {
	light := Open(nil)
	defer light.Close()

	err := light.View(func(txn *Txn) error {
		data := txn.Get([]byte("foo"))
		if len(data) != 0 {
			return fmt.Errorf("expecting empty result")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unable to get: %v", err)
	}
}
