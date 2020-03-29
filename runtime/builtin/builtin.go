package builtin

import (
	"encoding/json"
	"math"
	"strconv"

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
		"nan":                            math.NaN(),
		"inf":                            math.Inf(+1),
		"pi":                             math.Pi,
		"e":                              math.E,

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
			item, ok := pair.Item(index)
			if !ok {
				return nil, errors.Errorf("index %g is out of the range", index)
			}

			return item, nil
		},
		"type": func(value interface{}) (*types.Pair, error) {
			var name string
			switch value.(type) {
			case float64:
				name = "num"
			case *types.Pair:
				name = "list"
			default:
				return nil, errors.Errorf("unsupported type %T of the argument #0 for the function type", value)
			}

			return types.NewPairFromText(name), nil
		},
		"size": func(pair *types.Pair) (float64, error) {
			return float64(pair.Size()), nil
		},
		"floor": func(a float64) (float64, error) {
			return math.Floor(a), nil
		},
		"ceil": func(a float64) (float64, error) {
			return math.Ceil(a), nil
		},
		"trunc": func(a float64) (float64, error) {
			return math.Trunc(a), nil
		},
		"round": func(a float64) (float64, error) {
			return math.Round(a), nil
		},
		"sin": func(a float64) (float64, error) {
			return math.Sin(a), nil
		},
		"cos": func(a float64) (float64, error) {
			return math.Cos(a), nil
		},
		"tn": func(a float64) (float64, error) {
			return math.Tan(a), nil
		},
		"arcsin": func(a float64) (float64, error) {
			return math.Asin(a), nil
		},
		"arccos": func(a float64) (float64, error) {
			return math.Acos(a), nil
		},
		"arctn": func(a float64) (float64, error) {
			return math.Atan(a), nil
		},
		"angle": func(x float64, y float64) (float64, error) {
			return math.Atan2(y, x), nil
		},
		"pow": func(base float64, exponent float64) (float64, error) {
			return math.Pow(base, exponent), nil
		},
		"sqrt": func(a float64) (float64, error) {
			return math.Sqrt(a), nil
		},
		"exp": func(a float64) (float64, error) {
			return math.Exp(a), nil
		},
		"ln": func(a float64) (float64, error) {
			return math.Log(a), nil
		},
		"lg": func(a float64) (float64, error) {
			return math.Log10(a), nil
		},
		"abs": func(a float64) (float64, error) {
			return math.Abs(a), nil
		},
		"head": func(pair *types.Pair) (interface{}, error) {
			if pair == nil {
				return nil, errors.New("head of an empty list")
			}

			return pair.Head, nil
		},
		"tail": func(pair *types.Pair) (*types.Pair, error) {
			if pair == nil {
				return nil, errors.New("tail of an empty list")
			}

			return pair.Tail, nil
		},
		"num": func(pair *types.Pair) (float64, error) {
			text, err := pair.Text()
			if err != nil {
				return 0, errors.Wrap(err, "unable to convert the list to a string")
			}

			number, err := strconv.ParseFloat(text, 64)
			if err != nil {
				return 0, errors.Wrap(err, "unable to convert the string to a number")
			}

			return number, nil
		},
		"strs": func(pair *types.Pair) (*types.Pair, error) {
			text, err := pair.Text()
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert the list to a string")
			}

			text = strconv.Quote(text)
			return types.NewPairFromText(text), nil
		},
		"strl": func(pair *types.Pair) (*types.Pair, error) {
			var items []string
			for index, item := range pair.Slice() {
				itemPair, ok := item.(*types.Pair)
				if !ok {
					return nil, errors.Errorf(
						"incorrect type of the item #%d for conversion to a string list (%T instead *types.Pair)",
						index,
						item,
					)
				}

				itemText, err := itemPair.Text()
				if err != nil {
					return nil, errors.Wrapf(err, "unable to convert the item #%d to a string", index)
				}

				items = append(items, itemText)
			}

			textBytes, _ := json.Marshal(items) // nolint: gosec
			return types.NewPairFromText(string(textBytes)), nil
		},
	}
)
