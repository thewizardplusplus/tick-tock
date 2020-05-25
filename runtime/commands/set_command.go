package commands

import (
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// SetCommand ...
type SetCommand struct {
	state string
}

// NewSetCommand ...
func NewSetCommand(state string) SetCommand {
	return SetCommand{state}
}

// Run ...
func (command SetCommand) Run(context context.Context) (result interface{}, err error) {
	if err := context.SetState(command.state); err != nil {
		return nil, err
	}

	return types.Nil{}, nil
}
