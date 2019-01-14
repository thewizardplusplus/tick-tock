package commands

import (
	"io"

	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// Writer ...
//go:generate mockery -name=Writer -case=underscore
type Writer interface {
	io.Writer
}

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
func (command OutCommand) Run(context context.Context) error {
	_, err := command.writer.Write([]byte(command.message))
	return err
}
