package datastructures

// Storage provides basic abstraction over
// store of key values
type Storage interface {
	Get([]byte) ([]byte, error)
	Put([]byte, []byte) error
	Delete([]byte) error
	Exist([]byte) bool
}
