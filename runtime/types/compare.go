package types

import (
	"github.com/pkg/errors"
)

// ComparisonResult ...
type ComparisonResult int

// ...
const (
	Less ComparisonResult = iota
	Equal
	Greater
)

// Equals ...
func Equals(leftValue interface{}, rightValue interface{}) (bool, error) {
	// if operands have different types, they aren't equal
	switch leftValue.(type) {
	case Nil:
		if _, ok := rightValue.(Nil); !ok {
			return false, nil
		}
	case float64:
		if _, ok := rightValue.(float64); !ok {
			return false, nil
		}
	case *Pair:
		if _, ok := rightValue.(*Pair); !ok {
			return false, nil
		}
	default:
		return false, errors.Errorf(
			"unsupported type %T of the left value for comparison for equality",
			leftValue,
		)
	}

	result, err := Compare(leftValue, rightValue)
	if err != nil {
		return false, errors.Wrap(err, "unable to compare values for equality")
	}

	return result == Equal, nil
}

// Compare ...
func Compare(leftValue interface{}, rightValue interface{}) (ComparisonResult, error) {
	var result ComparisonResult
	switch typedLeftValue := leftValue.(type) {
	case Nil:
		if _, ok := rightValue.(Nil); !ok {
			return 0, errors.Errorf(
				"incorrect type of the right value for comparison (%T instead %T)",
				rightValue,
				leftValue,
			)
		}

		result = Equal
	case float64:
		typedRightValue, ok := rightValue.(float64)
		if !ok {
			return 0, errors.Errorf(
				"incorrect type of the right value for comparison (%T instead %T)",
				rightValue,
				leftValue,
			)
		}

		switch {
		case typedLeftValue < typedRightValue:
			result = Less
		case typedLeftValue == typedRightValue:
			result = Equal
		case typedLeftValue > typedRightValue:
			result = Greater
		}
	case *Pair:
		typedRightValue, ok := rightValue.(*Pair)
		if !ok {
			return 0, errors.Errorf(
				"incorrect type of the right value for comparison (%T instead %T)",
				rightValue,
				leftValue,
			)
		}

		var err error
		if result, err = typedLeftValue.Compare(typedRightValue); err != nil {
			return 0, errors.Wrap(err, "unable to compare pairs")
		}
	default:
		return 0, errors.Errorf("unsupported type %T of the left value for comparison", leftValue)
	}

	return result, nil
}
