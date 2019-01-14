package translator

import (
	"io"

	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/commands"
)

// TranslateStates ...
// TODO: disable the empty state group.
// TODO: disable few initial states.
// TODO: disable same states names.
// TODO: check setted states for existence.
// TODO: return an initial state.
func TranslateStates(writer io.Writer, states []*parser.State) runtime.StateGroup {
	translatedStates := make(runtime.StateGroup)
	for _, state := range states {
		translatedStates[state.Name] = TranslateMessages(writer, state.Messages)
	}

	return translatedStates
}

// TranslateMessages ...
// TODO: disable same messages names.
// TODO: return a map of messages names to setted states (with nonempty states only).
func TranslateMessages(writer io.Writer, messages []*parser.Message) runtime.MessageGroup {
	translatedMessages := make(runtime.MessageGroup)
	for _, message := range messages {
		translatedMessages[message.Name] = TranslateCommands(writer, message.Commands)
	}

	return translatedMessages
}

// TranslateCommands ...
// TODO: disable few set commands.
// TODO: return the setted state (by analogy with the TranslateCommand() function).
func TranslateCommands(writer io.Writer, commands []*parser.Command) runtime.CommandGroup {
	var translatedCommands runtime.CommandGroup
	for _, command := range commands {
		translatedCommand, _ := TranslateCommand(writer, command)
		translatedCommands = append(translatedCommands, translatedCommand)
	}

	return translatedCommands
}

// TranslateCommand ...
func TranslateCommand(writer io.Writer, command *parser.Command) (
	translatedCommand runtime.Command,
	settedState string,
) {
	if command.Send != nil {
		translatedCommand = commands.NewSendCommand(*command.Send)
	}
	if command.Set != nil {
		translatedCommand = commands.NewSetCommand(*command.Set)
		settedState = *command.Set
	}
	if command.Out != nil {
		translatedCommand = commands.NewOutCommand(writer, *command.Out)
	}
	if command.Exit {
		translatedCommand = commands.ExitCommand{}
	}

	return translatedCommand, settedState
}
