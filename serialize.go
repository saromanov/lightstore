package lightstore

import
(
	"encoding/json"
)
//This module provides serialization

type KeySerialize struct {
	Key string
	Value string
}

//JsonSerialization...
func JsonSerialization(data interface{}) string{
	result, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return string(result)
}
