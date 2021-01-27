package context

import (
	mapset "github.com/deckarep/golang-set"
)

// ValueStore ...
type ValueStore interface {
	ValueHolder

	ValuesNames() mapset.Set
	Value(name string) (value interface{}, ok bool)
}

//go:generate mockery --name=CopyableValueStore --inpackage --case=underscore --testonly

// CopyableValueStore ...
type CopyableValueStore interface {
	ValueStore

	Copy() CopyableValueStore
}

// DefaultValueStore ...
type DefaultValueStore map[string]interface{}

// ValuesNames ...
func (store DefaultValueStore) ValuesNames() mapset.Set {
	valuesNames := mapset.NewSet()
	for valueName := range store {
		valuesNames.Add(valueName)
	}

	return valuesNames
}

// Value ...
func (store DefaultValueStore) Value(name string) (value interface{}, ok bool) {
	value, ok = store[name]
	return value, ok
}

// SetValue ...
func (store DefaultValueStore) SetValue(name string, value interface{}) {
	store[name] = value
}

// Copy ...
func (store DefaultValueStore) Copy() CopyableValueStore {
	copy := make(DefaultValueStore)
	for name, value := range store {
		copy[name] = value
	}

	return copy
}
