package lightstore

import
(
	"encoding/json"
	"log"
)
//This module provides serialization

//JsonSerialization...
func JsonSerialization(data interface{}) string{
	result, err := json.Marshal(data)
	if err != nil {
		log.Printf(err)
	}

	return string(result)
}
