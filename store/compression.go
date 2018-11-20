package store

import (
	"bytes"
	"compress/zlib"
	"io"
)

// Compress provides compression of data
func Compress(data []byte) []byte {
	var b bytes.Buffer

	w := zlib.NewWriter(&b)
	w.Write(data)
	w.Close()
	return b.Bytes()
}

// Decompress provides decompression of data
func Decompress(data []byte) []byte {
	var buf bytes.Buffer
	b := bytes.NewReader(data)
	r, err := zlib.NewReader(b)
	if err != nil {
		panic(err)
	}
	io.Copy(buf, r)
	return buf.Bytes()
}
