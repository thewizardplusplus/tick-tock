package expressions

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// Identifier ...
type Identifier struct {
	name string
}

// NewIdentifier ...
func NewIdentifier(name string) Identifier {
	return Identifier{name}
}

// Evaluate ...
func (expression Identifier) Evaluate(context context.Context) (result interface{}, err error) {
	value, ok := context.Value(expression.name)
	if !ok {
		return nil, errors.Errorf("unknown identifier %s", expression.name)
	}

	return value, nil
}
