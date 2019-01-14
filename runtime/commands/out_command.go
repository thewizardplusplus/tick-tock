package commands

import (
	"io"

	"github.com/thewizardplusplus/tick-tock/runtime"
)

// Writer ...
//go:generate mockery -name=Writer -inpkg -case=underscore -testonly
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
func (command OutCommand) Run(context runtime.Context) error {
	_, err := command.writer.Write([]byte(command.message))
	return err
}
