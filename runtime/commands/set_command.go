package commands

import "github.com/thewizardplusplus/tick-tock/runtime"

// SetCommand ...
type SetCommand struct {
	state string
}

// NewSetCommand ...
func NewSetCommand(state string) SetCommand {
	return SetCommand{state}
}

// Run ...
func (command SetCommand) Run(context runtime.Context) error {
	return context.SetState(command.state)
}
