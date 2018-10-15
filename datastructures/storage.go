package datastructures

// Storage provides basic abstraction over
// store of key values
type Storage interface {
	Get([]byte) (interface{}, error)
	Put([]byte, interface{}, ItemOptions) error
	Delete([]byte) error
	Exist([]byte) bool
}
