package backend

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/saromanov/gohasha"
)

func Checksum(data string) string {
	check := sha256.New()
	check.Write([]byte(data))
	return hex.EncodeToString(check.Sum(nil))
}

func Hash(data string) string {
	hashstr, err := gohasha.GoHasha(&gohasha.GohashaOptions{Data: data})
	if err != nil {
		panic(err)
	}
	return hashstr
}
