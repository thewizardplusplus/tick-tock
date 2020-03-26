package expressions

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

const (
	expectedResultCount = 2
)

const (
	dataResultIndex = iota
	errorResultIndex
)

var (
	errorType = reflect.TypeOf((*error)(nil)).Elem()
)

// FunctionCall ...
type FunctionCall struct {
	name      string
	arguments []Expression
}

// NewFunctionCall ...
func NewFunctionCall(name string, arguments []Expression) FunctionCall {
	return FunctionCall{name, arguments}
}

// Evaluate ...
func (expression FunctionCall) Evaluate(context context.Context) (result interface{}, err error) {
	function, ok := context.Value(expression.name)
	if !ok {
		return nil, errors.Errorf("unknown function %s", expression.name)
	}

	functionType := reflect.TypeOf(function)
	if functionType.Kind() != reflect.Func {
		return nil, errors.Errorf("%s isn't function, it's %T", expression.name, function)
	}
	if functionType.NumIn() != len(expression.arguments) {
		return nil, errors.Errorf(
			"incorrect count of %s function arguments (%d instead %d)",
			expression.name,
			len(expression.arguments),
			functionType.NumIn(),
		)
	}
	if functionType.NumOut() != expectedResultCount {
		return nil, errors.Errorf(
			"incorrect count of %s function results (%d instead %d)",
			expression.name,
			functionType.NumOut(),
			expectedResultCount,
		)
	}
	if !functionType.Out(errorResultIndex).Implements(errorType) {
		return nil, errors.Errorf(
			"incorrect type of the result #%d of the function %s (%s instead %s)",
			errorResultIndex,
			expression.name,
			functionType.Out(errorResultIndex),
			errorType,
		)
	}

	var arguments []reflect.Value
	for index, argument := range expression.arguments {
		result, err2 := argument.Evaluate(context)
		if err2 != nil {
			return nil, errors.Wrapf(
				err2,
				"unable to evaluate the argument #%d for the function %s",
				index,
				expression.name,
			)
		}
		if !reflect.TypeOf(result).AssignableTo(functionType.In(index)) {
			return nil, errors.Errorf(
				"incorrect type of the argument #%d for the function %s (%T instead %s)",
				index,
				expression.name,
				result,
				functionType.In(index),
			)
		}

		arguments = append(arguments, reflect.ValueOf(result))
	}

	results := reflect.ValueOf(function).Call(arguments)
	if results[errorResultIndex].Interface() != nil {
		return nil, errors.Wrapf(
			results[errorResultIndex].Interface().(error),
			"unable to call the function %s",
			expression.name,
		)
	}

	return results[dataResultIndex].Interface(), nil
}
