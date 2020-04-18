package tests

import (
	"reflect"
)

// GetReflectionAddress ...
func GetReflectionAddress(v interface{}) uintptr {
	return reflect.ValueOf(v).Pointer()
}
