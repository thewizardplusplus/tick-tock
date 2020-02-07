package tests

import (
	"reflect"
)

// GetNumberAddress ...
func GetNumberAddress(f float64) *float64 {
	return &f
}

// GetStringAddress ...
func GetStringAddress(s string) *string {
	return &s
}

// GetReflectionAddress ...
func GetReflectionAddress(v interface{}) uintptr {
	return reflect.ValueOf(v).Pointer()
}
