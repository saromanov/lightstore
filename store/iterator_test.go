package store

import (
	"testing"
)

func TestIteratorBasic(t *testing.T) {
	s := newStore(nil)
	if !s.IsCreated() {
		t.Fatalf("unable to create db")
	}
	key := []byte("key")
	valueFirst := []byte("value")
	stored := s.Set(key, valueFirst)
	if stored != nil {
		t.Fatalf("unable to store data")
	}
}
