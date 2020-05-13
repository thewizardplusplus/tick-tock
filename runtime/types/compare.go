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
			return 0, errors.Wrapf(err, "unable to compare pairs")
		}
	default:
		return 0, errors.Errorf("unsupported type %T of the left value for comparison", leftValue)
	}

	return result, nil
}
