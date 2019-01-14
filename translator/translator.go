package translator

import (
	"io"

	"github.com/pkg/errors"
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
		translatedMessages[message.Name], _, _ = TranslateCommands(writer, message.Commands)
	}

	return translatedMessages
}

// TranslateCommands ...
func TranslateCommands(writer io.Writer, commands []*parser.Command) (
	translatedCommands runtime.CommandGroup,
	settedState string,
	err error,
) {
	for _, command := range commands {
		translatedCommand, newSettedState := TranslateCommand(writer, command)
		translatedCommands = append(translatedCommands, translatedCommand)
		if len(newSettedState) == 0 {
			continue
		}
		if len(settedState) != 0 {
			err = errors.Errorf("second set command %s (first was %s)", newSettedState, settedState)
			return nil, "", err
		}

		settedState = newSettedState
	}

	return translatedCommands, settedState, nil
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
