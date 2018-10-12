package store

import "testing"

func TestOpen(t *testing.T) {
	s := Open(nil)
	if !s.IsCreated() {
		t.Fatalf("unable to create db")
	}
}
