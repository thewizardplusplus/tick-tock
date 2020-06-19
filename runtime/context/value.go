package context

import (
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// ValueGroup ...
type ValueGroup map[string]interface{}

// ValueHolder ...
//go:generate mockery -name=ValueHolder -case=underscore
type ValueHolder interface {
	SetValue(name string, value interface{})
}

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
