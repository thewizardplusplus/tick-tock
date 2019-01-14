package translator

import (
	"io"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/commands"
)

// Dependencies ...
type Dependencies struct {
	runtime.Dependencies

	OutWriter io.Writer
}

// Translate ...
func Translate(inboxSize int, actors []*parser.Actor, dependencies Dependencies) (
	translatedActors runtime.ConcurrentActorGroup,
	err error,
) {
	for index, actor := range actors {
		translatedStates, initialState, err := translateStates(actor.States, dependencies.OutWriter)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to translate the actor #%d", index)
		}

		translatedActor, _ := runtime.NewActor(translatedStates, initialState) // nolint: gosec
		translatedActors = append(translatedActors, runtime.NewConcurrentActor(
			inboxSize,
			translatedActor,
			dependencies.Dependencies,
		))
	}

	return translatedActors, nil
}

func translateStates(states []*parser.State, outWriter io.Writer) (
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

		translatedMessages, settedStates, err := translateMessages(state.Messages, outWriter)
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

type settedStateGroup map[string]string

func translateMessages(messages []*parser.Message, outWriter io.Writer) (
	translatedMessages runtime.MessageGroup,
	settedStates settedStateGroup,
	err error,
) {
	translatedMessages = make(runtime.MessageGroup)
	settedStates = make(settedStateGroup)
	for _, message := range messages {
		if _, ok := translatedMessages[message.Name]; ok {
			return nil, nil, errors.Errorf("duplicate message %s", message.Name)
		}

		translatedCommands, settedState, err := translateCommands(message.Commands, outWriter)
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

func translateCommands(commands []*parser.Command, outWriter io.Writer) (
	translatedCommands runtime.CommandGroup,
	settedState string,
	err error,
) {
	for _, command := range commands {
		translatedCommand, newSettedState := translateCommand(command, outWriter)
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

func translateCommand(command *parser.Command, outWriter io.Writer) (
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
		translatedCommand = commands.NewOutCommand(*command.Out, outWriter)
	}
	if command.Exit {
		translatedCommand = commands.ExitCommand{}
	}

	return translatedCommand, settedState
}
