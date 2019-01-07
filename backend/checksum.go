package backend

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/saromanov/gohasha"
)

// Checksum provides getting of checksum on string
func Checksum(data string) string {
	check := sha256.New()
	check.Write([]byte(data))
	return hex.EncodeToString(check.Sum(nil))
}

// Hash provides hashing of data
func Hash(data string) string {
	if data == "" {
		return ""
	}
	hashstr, err := gohasha.GoHasha(&gohasha.GohashaOptions{Data: data})
	if err != nil {
		panic(err)
	}
	return hashstr
}
