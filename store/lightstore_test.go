package store

import (
	"testing"
)

func TestOpenLightstore(t *testing.T) {
	light := Open(nil)
	if !light.IsCreated() {
		t.Fatalf("unable to create db")
	}
}
