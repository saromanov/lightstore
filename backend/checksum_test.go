package backend

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHash(t *testing.T) {
	assert.Equal(t, Hash("123"), "202cb962ac59075b964b07152d234b70", "hashes is not equal")
}

func TestHashEmpty(t *testing.T) {
	assert.Equal(t, Hash(""), "", "hashes is not equal")
}

func TestChecksum(t *testing.T) {
	assert.Equal(t, Checksum("123"), "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3", "hashes is not equal")
}

func TestChecksumEmpty(t *testing.T) {
	assert.Equal(t, Checksum(""), "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855", "hashes is not equal")
}
