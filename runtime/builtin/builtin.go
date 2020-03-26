package builtin

import (
	"math"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
	"github.com/thewizardplusplus/tick-tock/translator"
)

// ...
// nolint: gochecknoglobals
var (
	Values = context.ValueGroup{
		translator.EmptyListConstantName: (*types.Pair)(nil),

		translator.ListConstructionFunctionName: func(
			head interface{},
			tail *types.Pair,
		) (*types.Pair, error) {
			return &types.Pair{Head: head, Tail: tail}, nil
		},
		translator.AdditionFunctionName: func(a interface{}, b interface{}) (interface{}, error) {
			switch typedA := a.(type) {
			case float64:
				if typedB, ok := b.(float64); ok {
					return typedA + typedB, nil
				}
			case *types.Pair:
				if typedB, ok := b.(*types.Pair); ok {
					return typedA.Append(typedB), nil
				}
			default:
				return nil, errors.Errorf(
					"unsupported type %T of the argument #0 for the function %s",
					a,
					translator.AdditionFunctionName,
				)
			}

			return nil, errors.Errorf(
				"incorrect type of the argument #1 for the function %s (%T instead %T)",
				translator.AdditionFunctionName,
				b,
				a,
			)
		},
		translator.SubtractionFunctionName: func(a float64, b float64) (float64, error) {
			return a - b, nil
		},
		translator.MultiplicationFunctionName: func(a float64, b float64) (float64, error) {
			return a * b, nil
		},
		translator.DivisionFunctionName: func(a float64, b float64) (float64, error) {
			return a / b, nil
		},
		translator.ModuloFunctionName: func(a float64, b float64) (float64, error) {
			return math.Mod(a, b), nil
		},
		translator.NegationFunctionName: func(a float64) (float64, error) {
			return -a, nil
		},
		translator.KeyAccessorFunctionName: func(pair *types.Pair, index float64) (interface{}, error) {
			integerIndex := int(index)
			item, ok := pair.Item(integerIndex)
			if !ok {
				return nil, errors.Errorf("index %d is out of the range", integerIndex)
			}

			return item, nil
		},
	}
)
