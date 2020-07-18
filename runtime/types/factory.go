package types

import (
	"reflect"
)

func isActorClass(value interface{}) bool {
	// can't use this type directly because it occurs an import cycle
	return reflect.TypeOf(value).Name() == "ConcurrentActorFactory"
}

func getActorClassName(actorClass interface{}) string {
	// can't use this type directly because it occurs an import cycle
	results := reflect.ValueOf(actorClass).MethodByName("Name").Call(nil)
	return results[0].Interface().(string)
}
