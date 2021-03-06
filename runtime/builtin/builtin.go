package builtin

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
	"github.com/thewizardplusplus/tick-tock/translator"
)

type lineBreakMode int

const (
	withoutLineBreak lineBreakMode = iota
	withLineBreak
)

// nolint: gochecknoglobals
var (
	bufferedStdin = bufio.NewReader(os.Stdin)
)

// ...
// nolint: gochecknoglobals
var (
	Values = context.ValueGroup{
		translator.EmptyListConstantName:      (*types.Pair)(nil),
		translator.EmptyHashTableConstantName: (types.HashTable)(nil),
		"nil":                                 types.Nil{},
		"false":                               types.False,
		"true":                                types.True,
		"nan":                                 math.NaN(),
		"inf":                                 math.Inf(+1),
		"pi":                                  math.Pi,
		"e":                                   math.E,

		translator.ListConstructionFunctionName: func(
			head interface{},
			tail *types.Pair,
		) (*types.Pair, error) {
			return &types.Pair{Head: head, Tail: tail}, nil
		},
		translator.HashTableConstructionFunctionName: types.HashTable.With,
		translator.EqualFunctionName: func(a interface{}, b interface{}) (types.Boolean, error) {
			isEqual, err := types.Equals(a, b)
			if err != nil {
				return 0, errors.Wrap(err, "unable to compare values for equality")
			}

			return types.NewBooleanFromGoBool(isEqual), nil
		},
		translator.NotEqualFunctionName: func(a interface{}, b interface{}) (types.Boolean, error) {
			isEqual, err := types.Equals(a, b)
			if err != nil {
				return 0, errors.Wrap(err, "unable to compare values for equality")
			}

			return types.NewBooleanFromGoBool(!isEqual), nil
		},
		translator.LessFunctionName: func(a interface{}, b interface{}) (types.Boolean, error) {
			compareResult, err := types.Compare(a, b)
			if err != nil {
				return 0, errors.Wrap(err, "unable to compare values")
			}

			booleanResult := compareResult == types.Less
			return types.NewBooleanFromGoBool(booleanResult), nil
		},
		translator.LessOrEqualFunctionName: func(a interface{}, b interface{}) (types.Boolean, error) {
			compareResult, err := types.Compare(a, b)
			if err != nil {
				return 0, errors.Wrap(err, "unable to compare values")
			}

			booleanResult := compareResult == types.Less || compareResult == types.Equal
			return types.NewBooleanFromGoBool(booleanResult), nil
		},
		translator.GreaterFunctionName: func(a interface{}, b interface{}) (types.Boolean, error) {
			compareResult, err := types.Compare(a, b)
			if err != nil {
				return 0, errors.Wrap(err, "unable to compare values")
			}

			booleanResult := compareResult == types.Greater
			return types.NewBooleanFromGoBool(booleanResult), nil
		},
		translator.GreaterOrEqualFunctionName: func(a interface{}, b interface{}) (types.Boolean, error) {
			compareResult, err := types.Compare(a, b)
			if err != nil {
				return 0, errors.Wrap(err, "unable to compare values")
			}

			booleanResult := compareResult == types.Greater || compareResult == types.Equal
			return types.NewBooleanFromGoBool(booleanResult), nil
		},
		translator.BitwiseDisjunctionFunctionName: func(a float64, b float64) (float64, error) {
			return float64(int64(a) | int64(b)), nil
		},
		translator.BitwiseExclusiveDisjunctionFunctionName: func(a float64, b float64) (float64, error) {
			return float64(int64(a) ^ int64(b)), nil
		},
		translator.BitwiseConjunctionFunctionName: func(a float64, b float64) (float64, error) {
			return float64(int64(a) & int64(b)), nil
		},
		translator.BitwiseLeftShiftFunctionName: func(a float64, b float64) (float64, error) {
			return float64(int64(a) << uint64(b)), nil
		},
		translator.BitwiseRightShiftFunctionName: func(a float64, b float64) (float64, error) {
			return float64(int64(a) >> uint64(b)), nil
		},
		translator.BitwiseUnsignedRightShiftFunctionName: func(a float64, b float64) (float64, error) {
			if a < 0 {
				a = a + (1 << 32)
			}
			return float64(int64(a) >> uint64(b)), nil
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
			case types.HashTable:
				if typedB, ok := b.(types.HashTable); ok {
					return typedA.Merge(typedB), nil
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
		translator.ArithmeticNegationFunctionName: func(a float64) (float64, error) {
			return -a, nil
		},
		translator.BitwiseNegationFunctionName: func(a float64) (float64, error) {
			return float64(^int64(a)), nil
		},
		translator.LogicalNegationFunctionName: func(value interface{}) (types.Boolean, error) {
			boolean, err := types.NewBoolean(value)
			if err != nil {
				return 0, errors.Wrap(err, "unable to convert the value to a boolean")
			}

			return types.NegateBoolean(boolean), nil
		},
		translator.KeyAccessorFunctionName: func(
			value interface{},
			key interface{},
		) (interface{}, error) {
			var item interface{}
			switch typedValue := value.(type) {
			case *types.Pair:
				typedKey, ok := key.(float64)
				if !ok {
					return nil, errors.Errorf(
						"incorrect type of the argument #1 for the function %s (%T instead float64)",
						translator.KeyAccessorFunctionName,
						key,
					)
				}

				item, ok = typedValue.Item(typedKey)
				if !ok {
					return types.Nil{}, nil
				}
			case types.HashTable:
				var err error
				item, err = typedValue.Item(key)
				switch err {
				case nil:
				case types.ErrNotFound:
					return types.Nil{}, nil
				default:
					return nil, errors.Wrap(err, "unable to get an item from the hash table")
				}
			default:
				return nil, errors.Errorf(
					"unsupported type %T of the argument #0 for the function %s",
					value,
					translator.KeyAccessorFunctionName,
				)
			}

			return item, nil
		},
		"type": func(value interface{}) (*types.Pair, error) {
			var name string
			switch value.(type) {
			case types.Nil:
				name = "nil"
			case float64:
				name = "num"
			case *types.Pair:
				name = "list"
			case types.HashTable:
				name = "hash"
			case runtime.ConcurrentActorFactory:
				name = "class"
			default:
				return nil, errors.Errorf("unsupported type %T of the argument #0 for the function type", value)
			}

			return types.NewPairFromText(name), nil
		},
		"name": func(factory runtime.ConcurrentActorFactory) (*types.Pair, error) {
			name := factory.Name()
			return types.NewPairFromText(name), nil
		},
		"size": func(value interface{}) (float64, error) {
			typedValue, ok := value.(interface{ Size() int })
			if !ok {
				return 0, errors.Errorf("unsupported type %T of the argument #0 for the function size", value)
			}

			size := float64(typedValue.Size())
			return size, nil
		},
		"bool": types.NewBoolean,
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
		"is_nan": func(a float64) (float64, error) {
			isNaN := math.IsNaN(a)
			return types.NewBooleanFromGoBool(isNaN), nil
		},
		"seed": func(seed float64) (types.Nil, error) {
			rand.Seed(int64(seed))
			return types.Nil{}, nil
		},
		"random": func() (float64, error) {
			return rand.Float64(), nil // nolint: gosec
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
		"num": func(text *types.Pair) (interface{}, error) {
			textAsString, err := text.Text()
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert the list to a string")
			}

			number, err := strconv.ParseFloat(textAsString, 64)
			if err != nil {
				return types.Nil{}, nil
			}

			return number, nil
		},
		"str": func(value interface{}) (*types.Pair, error) {
			var text string
			switch typedValue := value.(type) {
			case float64:
				text = strconv.FormatFloat(typedValue, 'g', -1, 64)
			case *types.Pair, types.HashTable:
				var err error
				text, err = marshalToJSON(value)
				if err != nil {
					return nil, err
				}
			case fmt.Stringer:
				text = typedValue.String()
			default:
				return nil, errors.Errorf(
					"unsupported type %T of the argument #0 for the function str",
					typedValue,
				)
			}

			return types.NewPairFromText(text), nil
		},
		"strb": func(value interface{}) (*types.Pair, error) {
			boolean, err := types.NewBoolean(value)
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert the value to a boolean")
			}

			var text string
			switch boolean {
			case types.False:
				text = "false"
			case types.True:
				text = "true"
			}

			return types.NewPairFromText(text), nil
		},
		"strs": func(text *types.Pair) (*types.Pair, error) {
			textAsString, err := text.Text()
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert the list to a string")
			}

			textAsString = strconv.Quote(textAsString)
			return types.NewPairFromText(textAsString), nil
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

			text, _ := marshalToJSON(items) // nolint: gosec
			return types.NewPairFromText(text), nil
		},
		"strh": func(table types.HashTable) (*types.Pair, error) {
			text, err := marshalToJSON(table)
			if err != nil {
				return nil, err
			}

			return types.NewPairFromText(text), nil
		},
		"strhh": func(table types.HashTable) (*types.Pair, error) {
			pairs := make(map[string]string)
			for key, value := range table {
				keyText, ok := key.(string)
				if !ok {
					return nil, errors.Errorf(
						"incorrect type of the key for conversion to a string (%T instead *types.Pair)",
						key,
					)
				}

				valuePair, ok := value.(*types.Pair)
				if !ok {
					return nil, errors.Errorf(
						"incorrect type of the value for conversion to a string (%T instead *types.Pair)",
						value,
					)
				}

				valueText, err := valuePair.Text()
				if err != nil {
					return nil, errors.Wrap(err, "unable to convert the value to a string")
				}

				pairs[keyText] = valueText
			}

			text, _ := marshalToJSON(pairs) // nolint: gosec
			return types.NewPairFromText(text), nil
		},
		"with": types.HashTable.With,
		"keys": func(table types.HashTable) (*types.Pair, error) {
			keys := table.Keys()
			return types.NewPairFromSlice(keys), nil
		},
		"env": func(name *types.Pair) (interface{}, error) {
			nameText, err := name.Text()
			if err != nil {
				return nil, errors.Wrap(err, "unable to convert the list to a string")
			}

			value, ok := os.LookupEnv(nameText)
			if !ok {
				return types.Nil{}, nil
			}

			return types.NewPairFromText(value), nil
		},
		"time": func() (float64, error) {
			timestamp := time.Now().UnixNano()
			return float64(timestamp) / 1e9, nil
		},
		"sleep": func(duration float64) (types.Nil, error) {
			time.Sleep(time.Duration(duration * 1e9))
			return types.Nil{}, nil
		},
		"exit": func(exitCode float64) (types.Nil, error) {
			os.Exit(int(exitCode))
			return types.Nil{}, nil
		},
		"in": func(count float64) (interface{}, error) {
			if count >= 0 {
				return readChunk(count)
			}

			textBytes, err := ioutil.ReadAll(bufferedStdin)
			if err != nil {
				return types.Nil{}, nil
			}

			return types.NewPairFromText(string(textBytes)), nil
		},
		"inln": func(count float64) (interface{}, error) {
			if count >= 0 {
				return readChunk(count)
			}

			textBytes, err := bufferedStdin.ReadBytes('\n')
			if err != nil && err != io.EOF {
				return types.Nil{}, nil
			}

			// remove the line break in case of success
			if err == nil {
				textBytes = textBytes[:len(textBytes)-1]
			}

			return types.NewPairFromText(string(textBytes)), nil
		},
		"out": func(text *types.Pair) (types.Nil, error) {
			return print(os.Stdout, text, withoutLineBreak)
		},
		"outln": func(text *types.Pair) (types.Nil, error) {
			return print(os.Stdout, text, withLineBreak)
		},
		"err": func(text *types.Pair) (types.Nil, error) {
			return print(os.Stderr, text, withoutLineBreak)
		},
		"errln": func(text *types.Pair) (types.Nil, error) {
			return print(os.Stderr, text, withLineBreak)
		},
	}
)

func marshalToJSON(value interface{}) (string, error) {
	var err error
	value, err = types.GetDeepValue(value)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(value); err != nil {
		return "", errors.Wrap(err, "unable to marshal the value to JSON")
	}

	textBytes := buffer.Bytes()
	textBytes = textBytes[:buffer.Len()-1] // remove the line break

	return string(textBytes), nil
}

func readChunk(size float64) (interface{}, error) {
	chunkBytes := make([]byte, int(size))
	readSize, err := io.ReadFull(bufferedStdin, chunkBytes)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return types.Nil{}, nil
	}

	chunkBytes = chunkBytes[:readSize]
	return types.NewPairFromText(string(chunkBytes)), nil
}

func print(writer io.Writer, text *types.Pair, mode lineBreakMode) (types.Nil, error) {
	textAsString, err := text.Text()
	if err != nil {
		return types.Nil{}, errors.Wrap(err, "unable to convert the list to a string")
	}

	fmt.Fprint(writer, textAsString) // nolint: errcheck
	if mode == withLineBreak {
		fmt.Fprintln(writer) // nolint: errcheck
	}

	return types.Nil{}, nil
}
