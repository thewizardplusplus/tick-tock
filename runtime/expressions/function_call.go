package expressions

import (
	"reflect"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// NewFunctionCall ...
func NewFunctionCall(name string, arguments []Expression) FunctionApplying {
	return NewFunctionApplying(
		arguments,
		func(context context.Context, arguments []interface{}) (interface{}, error) {
			function, ok := context.Value(name)
			if !ok {
				return nil, errors.Errorf("unknown function %s", name)
			}

			var reflectedArguments []reflect.Value
			for _, argument := range arguments {
				reflectedArguments = append(reflectedArguments, reflect.ValueOf(argument))
			}

			var results []reflect.Value
			var err error
			func() {
				defer func() {
					if err2 := recover(); err2 != nil {
						err = errors.Errorf("%v", err2)
					}
				}()

				results = reflect.ValueOf(function).Call(reflectedArguments)
			}()
			if err != nil {
				return nil, errors.Wrapf(err, "unable to call the function %s", name)
			}
			if resultCount := len(results); resultCount != 1 {
				return nil, errors.Errorf(
					"incorrect number of results (%d) for the function %s",
					resultCount,
					name,
				)
			}

			return results[0].Interface(), nil
		},
	)
}
