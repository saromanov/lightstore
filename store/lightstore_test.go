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
