package translator

import (
	"io"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/commands"
)

// TranslateStates ...
func TranslateStates(writer io.Writer, states []*parser.State) (
	translatedStates runtime.StateGroup,
	initialState string,
	err error,
) {
	if len(states) == 0 {
		return nil, "", errors.New("no states")
	}

	translatedStates = make(runtime.StateGroup)
	messagesWithSettings := make(map[string][]string)
	for _, state := range states {
		if _, ok := translatedStates[state.Name]; ok {
			return nil, "", errors.Errorf("duplicate state %s", state.Name)
		}

		if state.Initial {
			if len(initialState) != 0 {
				err := errors.Errorf("second initial state %s (first was %s)", state.Name, initialState)
				return nil, "", err
			}

			initialState = state.Name
		}

		translatedMessages, settedStates, err := TranslateMessages(writer, state.Messages)
		if err != nil {
			return nil, "", errors.Wrapf(err, "unable to translate the state %s", state.Name)
		}

		translatedStates[state.Name] = translatedMessages
		for message, state := range settedStates {
			messagesWithSettings[state] = append(messagesWithSettings[state], message)
		}
	}

	for state, messages := range messagesWithSettings {
		if _, ok := translatedStates[state]; !ok {
			return nil, "", errors.Errorf("unknown state %s in messages %v", state, messages)
		}
	}

	if len(initialState) == 0 {
		initialState = states[0].Name
	}

	return translatedStates, initialState, nil
}

// SettedStateGroup ...
type SettedStateGroup map[string]string

// TranslateMessages ...
func TranslateMessages(writer io.Writer, messages []*parser.Message) (
	translatedMessages runtime.MessageGroup,
	settedStates SettedStateGroup,
	err error,
) {
	translatedMessages = make(runtime.MessageGroup)
	settedStates = make(SettedStateGroup)
	for _, message := range messages {
		if _, ok := translatedMessages[message.Name]; ok {
			return nil, nil, errors.Errorf("duplicate message %s", message.Name)
		}

		translatedCommands, settedState, err := TranslateCommands(writer, message.Commands)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to translate the message %s", message.Name)
		}

		translatedMessages[message.Name] = translatedCommands
		if len(settedState) != 0 {
			settedStates[message.Name] = settedState
		}
	}

	return translatedMessages, settedStates, nil
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
			err := errors.Errorf("second set command %s (first was %s)", newSettedState, settedState)
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
