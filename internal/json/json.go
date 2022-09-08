package json

import (
	jsoniter "github.com/json-iterator/go"
)

var jsoniterJson = jsoniter.ConfigCompatibleWithStandardLibrary

//func init() {
//	extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)
//}

func Marshal(v interface{}) ([]byte, error) {
	return jsoniterJson.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
	return jsoniterJson.Unmarshal(data, v)
}
