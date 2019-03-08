package store

import "errors"

var (
	errNoWrites          = errors.New("unable to write on read-only mode")
	errEmptyKey          = errors.New("key is empty")
	errLargeKeySize      = errors.New("key size is larger then limit")
	errNoStorage         = errors.New("storage is not defined or transaction was closed")
	errMaxKeySize        = errors.New("key size is greather then max")
	errMaxValueSize      = errors.New("value size is greather then max")
	errNoIndexName       = errors.New("index name is not defined")
	errNoIndexNameOrData = errors.New("index name of data is not defined")
	errIndexExist        = errors.New("index with this name already exist")
	errFnIsNotDefined    = errors.New("function is not defined")
)
