package expressions

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
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

	var arguments []reflect.Value
	for index, argument := range expression.arguments {
		result, err2 := argument.Evaluate(context)
		if err2 != nil {
			return nil, errors.Wrapf(
				err2,
				"unable to evaluate the argument %d for the function %s",
				index,
				expression.name,
			)
		}

		arguments = append(arguments, reflect.ValueOf(result))
	}

	var results []reflect.Value
	func() {
		defer func() {
			if err2 := recover(); err2 != nil {
				err = errors.Errorf("%v", err2)
			}
		}()

		results = reflect.ValueOf(function).Call(arguments)
	}()
	if err != nil {
		return nil, errors.Wrapf(err, "unable to call the function %s", expression.name)
	}
	if resultCount := len(results); resultCount != 1 {
		return nil, errors.Errorf(
			"incorrect number of results (%d) for the function %s",
			resultCount,
			expression.name,
		)
	}

	return results[0].Interface(), nil
}
