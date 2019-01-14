package commands

import (
	"io"

	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// OutCommand ...
type OutCommand struct {
	message string
	writer  io.Writer
}

// NewOutCommand ...
func NewOutCommand(message string, writer io.Writer) OutCommand {
	return OutCommand{message, writer}
}

// Run ...
func (command OutCommand) Run(context context.Context) error {
	_, err := command.writer.Write([]byte(command.message))
	return err
}
