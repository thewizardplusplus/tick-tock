package commands

import (
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// ReturnCommand ...
type ReturnCommand struct{}

// Run ...
func (command ReturnCommand) Run(context context.Context) (result interface{}, err error) {
	return nil, runtime.ErrReturn
}
