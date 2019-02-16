package store

// index provides representation of index
type index struct {
	name  string
	data  string
	fnc   func(a, b []byte) bool
	store *Store
}
