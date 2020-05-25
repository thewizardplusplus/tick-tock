package commands

import (
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// SendCommand ...
type SendCommand struct {
	message string
}

// NewSendCommand ...
func NewSendCommand(message string) SendCommand {
	return SendCommand{message}
}

// Run ...
func (command SendCommand) Run(context context.Context) (result interface{}, err error) {
	context.SendMessage(command.message)
	return types.Nil{}, nil
}
