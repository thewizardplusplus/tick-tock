package context

// ValueStore ...
type ValueStore interface {
	Value(name string) (value interface{}, ok bool)
	SetValue(name string, value interface{})
}

// CopyableValueStore ...
//go:generate mockery -name=CopyableValueStore -case=underscore
type CopyableValueStore interface {
	ValueStore

	Copy() CopyableValueStore
}

// DefaultValueStore ...
type DefaultValueStore map[string]interface{}

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
		copy.SetValue(name, value)
	}

	return copy
}
