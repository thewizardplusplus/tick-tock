package types

import (
	"github.com/pkg/errors"
)

// Boolean ...
//
// This type should be an alias.
type Boolean = float64

// ...
const (
	False Boolean = iota
	True
)

// NewBoolean ...
func NewBoolean(value interface{}) (Boolean, error) {
	if isActorClass(value) {
		return True, nil
	}

	var result Boolean
	switch typedValue := value.(type) {
	case Nil:
		result = False
	case float64:
		result = NewBooleanFromGoBool(typedValue != 0)
	case *Pair:
		result = NewBooleanFromGoBool(typedValue != nil)
	default:
		return False, errors.Errorf("unsupported type %T for conversion to boolean", value)
	}

	return result, nil
}

// NewBooleanFromGoBool ...
func NewBooleanFromGoBool(value bool) Boolean {
	if value {
		return True
	}

	return False
}

// NegateBoolean ...
func NegateBoolean(value Boolean) Boolean {
	if value == False {
		return True
	}

	return False
}
