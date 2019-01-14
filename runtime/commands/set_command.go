package commands

import "github.com/thewizardplusplus/tick-tock/runtime/context"

// SetCommand ...
type SetCommand struct {
	state string
}

// NewSetCommand ...
func NewSetCommand(state string) SetCommand {
	return SetCommand{state}
}

// Run ...
// TODO: wrap the SetState() error with the method name.
func (command SetCommand) Run(context context.Context) error {
	return context.SetState(command.state)
}
