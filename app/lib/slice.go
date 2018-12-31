package lib

import (
	"fmt"
	"reflect"
)

// InterfaceToSlice cast interface{} to slice
func InterfaceToSlice(slice interface{}) []string {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return nil
	}
	ret := make([]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		ret[i] = fmt.Sprintf("%v", s.Index(i).Interface())
	}
	return ret
}
