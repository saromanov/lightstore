package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompression(t *testing.T) {
	res := compress([]byte("abcd"))
	assert.Equal(t, res, []byte{0x78, 0x9c, 0x4a, 0x4c, 0x4a, 0x4e, 0x1, 0x4, 0x0, 0x0, 0xff, 0xff, 0x3, 0xd8, 0x1, 0x8b}, "should be equal")
	assert.Equal(t, decompress(res), []byte("abcd"), "should be equal")
}
