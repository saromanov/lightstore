package store

import (
	"compress/zlib"
	"bytes"
)
// Compress provides compression of data
func Compress(data []byte)[]byte {
	var b bytes.Buffer

	w := zlib.NewWriter(&b)
	w.Write(data)
	w.Close()
	return b.Bytes()
}