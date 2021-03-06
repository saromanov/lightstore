package store

import (
	"testing"
)

func TestOpen(t *testing.T) {
	s, err := Open(nil)
	if err != nil {
		t.Fatal(err)
	}
	if !s.IsCreated() {
		t.Fatalf("unable to create db")
	}
}

func TestWrite(t *testing.T) {
	s, err := newStore(nil)
	if err != nil {
		t.Fatal(err)
	}
	if !s.IsCreated() {
		t.Fatalf("unable to create db")
	}
	key := []byte("key")
	valueFirst := []byte("value")
	stored := s.Set(key, valueFirst)
	if stored != nil {
		t.Fatalf("unable to store data")
	}
	data := s.Get(key)
	if string(data) != string(valueFirst) {
		t.Fatalf("unable to get data")
	}
}
