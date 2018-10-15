package store

import (
	"testing"
)

func TestOpen(t *testing.T) {
	s := Open(nil)
	if !s.IsCreated() {
		t.Fatalf("unable to create db")
	}
	stored := s.Set([]byte("key"), []byte("value"))
	if !stored {
		t.Fatalf("unable to store data")
	}

	key := []byte("key")
	data := s.Get(key)
	value := data.([]byte)
	if string(key) != string(value) {
		t.Fatalf("unable to get data")
	}
}
