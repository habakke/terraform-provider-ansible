package util

import (
	"encoding/json"
	"reflect"
)

func CanMarshal(obj interface{}) bool {
	rv := reflect.ValueOf(obj)
	_, ok := rv.Interface().(json.Marshaler)
	return ok
}
