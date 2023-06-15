package json

import (
	jsoniter "github.com/json-iterator/go"
)

func mapper() jsoniter.API {
	return jsoniter.ConfigCompatibleWithStandardLibrary
}

func Stringify(data interface{}) (string, error) {
	return mapper().MarshalToString(data)
}

func FromJSON(json []byte, result interface{}) error {
	return mapper().Unmarshal(json, result)
}
