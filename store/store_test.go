package store

import (
	"testing"
)

func TestOpen(t *testing.T) {
	s := Open(nil)
	if !s.IsCreated() {
		t.Fatalf("unable to create db")
	}
	key := []byte("key")
	valueFirst := []byte("value")
	stored := s.Set(key, valueFirst)
	if !stored {
		t.Fatalf("unable to store data")
	}
	data := s.Get(key)
	value := data.([]byte)
	if string(value) != string(valueFirst) {
		t.Fatalf("unable to get data")
	}
}
