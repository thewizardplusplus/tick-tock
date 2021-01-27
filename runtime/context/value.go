package context

import (
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

//go:generate mockery --name=ValueHolder --inpackage --case=underscore --testonly

// ValueHolder ...
type ValueHolder interface {
	SetValue(name string, value interface{})
}

// ValueGroup ...
type ValueGroup map[string]interface{}

// ZipValues ...
func ZipValues(parameters []string, arguments []interface{}) ValueGroup {
	values := make(ValueGroup)
	for index, name := range parameters {
		var value interface{}
		if index < len(arguments) {
			value = arguments[index]
		} else {
			value = types.Nil{}
		}

		values[name] = value
	}

	return values
}

// SetValues ...
func SetValues(holder ValueHolder, values ValueGroup) {
	for name, value := range values {
		holder.SetValue(name, value)
	}
}
