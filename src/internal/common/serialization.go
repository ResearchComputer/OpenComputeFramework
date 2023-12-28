package common

import "encoding/json"

func DictionaryToBytes(dict map[string]interface{}) ([]byte, error) {
	return json.Marshal(dict)
}
