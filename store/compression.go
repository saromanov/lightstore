package store

import (
	"bufio"
	"bytes"
	"compress/zlib"
	"io"
)

// compress provides compression of data
func compress(data []byte) []byte {
	var b bytes.Buffer

	w := zlib.NewWriter(&b)
	w.Write(data)
	w.Close()
	return b.Bytes()
}

// decompress provides decompression of data
func decompress(data []byte) []byte {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	b := bytes.NewReader(data)
	r, err := zlib.NewReader(b)
	if err != nil {
		panic(err)
	}
	io.Copy(writer, r)
	writer.Flush()
	return buf.Bytes()
}
