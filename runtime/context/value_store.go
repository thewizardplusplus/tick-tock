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
