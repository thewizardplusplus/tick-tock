package translator

import (
	"io"

	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/commands"
)

// TranslateCommand ...
func TranslateCommand(writer io.Writer, command *parser.Command) runtime.Command {
	var translatedCommand runtime.Command
	if command.Send != nil {
		translatedCommand = commands.NewSendCommand(*command.Send)
	}
	if command.Set != nil {
		translatedCommand = commands.NewSetCommand(*command.Set)
	}
	if command.Out != nil {
		translatedCommand = commands.NewOutCommand(writer, *command.Out)
	}
	if command.Exit {
		translatedCommand = commands.ExitCommand{}
	}

	return translatedCommand
}
