package commands

import (
	"io"

	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// OutCommand ...
type OutCommand struct {
	writer  io.Writer
	message string
}

// NewOutCommand ...
func NewOutCommand(writer io.Writer, message string) OutCommand {
	return OutCommand{writer, message}
}

// Run ...
// TODO: wrap the Write() error with the method name.
func (command OutCommand) Run(context context.Context) error {
	_, err := command.writer.Write([]byte(command.message))
	return err
}
