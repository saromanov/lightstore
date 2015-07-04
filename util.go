package lightstore

import
(
	"strconv"
)

func Bool(str string) bool {
	res, err := strconv.ParseBool(str)
	if err != nil {
		return false
	}

	return res
}