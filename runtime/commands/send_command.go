package commands

import "github.com/thewizardplusplus/tick-tock/runtime/context"

// SendCommand ...
type SendCommand struct {
	message string
}

// NewSendCommand ...
func NewSendCommand(message string) SendCommand {
	return SendCommand{message}
}

// Run ...
func (command SendCommand) Run(context context.Context) error {
	context.SendMessage(command.message)
	return nil
}
