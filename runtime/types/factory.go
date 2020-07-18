package types

import (
	"reflect"
)

func isActorClass(value interface{}) bool {
	// can't use this type directly because it occurs an import cycle
	return reflect.TypeOf(value).Name() == "ConcurrentActorFactory"
}
