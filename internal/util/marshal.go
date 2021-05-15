package util

import (
	"encoding/json"
	"reflect"
)

// CanMarshal returns if an object implements the json.Marshaler interface and thus can be marshaled
func CanMarshal(obj interface{}) bool {
	rv := reflect.ValueOf(obj)
	_, ok := rv.Interface().(json.Marshaler)
	return ok
}
