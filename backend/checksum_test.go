package backend

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	assert.Equal(t, Hash("123"),"202cb962ac59075b964b07152d234b70", "hashes is not equal")
}

func TestHashEmpty(t *testing.T) {
	assert.Equal(t, Hash(""),"", "hashes is not equal")
	recover()
}