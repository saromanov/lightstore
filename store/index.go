package store

// index provides representation of index
type index struct {
	name    string
	pattern string
	F       func(a, b []byte) bool
}
