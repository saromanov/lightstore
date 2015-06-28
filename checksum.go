package lightstore

import
(
	"crypto/sha256"
	"encoding/hex"
)

func Checksum(data string) string {
	check := sha256.New()
	check.Write([]byte(data))
	return hex.EncodeToString(check.Sum(nil))
}