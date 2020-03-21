package translator

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/commands"
)

// DeclaredIdentifierGroup ...
type DeclaredIdentifierGroup map[string]struct{}

// Options ...
type Options struct {
	InboxSize    int
	InitialState string
}

// Dependencies ...
type Dependencies struct {
	Commands commands.Dependencies
	Runtime  runtime.Dependencies
}

// Translate ...
func Translate(
	actors []*parser.Actor,
	declaredIdentifiers DeclaredIdentifierGroup,
	options Options,
	dependencies Dependencies,
) (
	translatedActors runtime.ConcurrentActorGroup,
	err error,
) {
	for index, actor := range actors {
		translatedStates, err := translateStates(actor.States, declaredIdentifiers, dependencies.Commands)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to translate the actor #%d", index)
		}

		translatedActor, err := runtime.NewActor(translatedStates, options.InitialState)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to construct the actor #%d", index)
		}

		translatedActors = append(translatedActors, runtime.NewConcurrentActor(
			translatedActor,
			options.InboxSize,
			dependencies.Runtime,
		))
	}

	return translatedActors, nil
}

func translateStates(
	states []*parser.State,
	declaredIdentifiers DeclaredIdentifierGroup,
	dependencies commands.Dependencies,
) (
	translatedStates runtime.StateGroup,
	err error,
) {
	if len(states) == 0 {
		return nil, errors.New("no states")
	}

	translatedStates = make(runtime.StateGroup)
	messagesWithSettings := make(map[string][]string)
	for _, state := range states {
		if _, ok := translatedStates[state.Name]; ok {
			return nil, errors.Errorf("duplicate state %s", state.Name)
		}

		translatedMessages, settedStates, err :=
			translateMessages(state.Messages, declaredIdentifiers, dependencies)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to translate the state %s", state.Name)
		}

		translatedStates[state.Name] = translatedMessages
		for message, state := range settedStates {
			messagesWithSettings[state] = append(messagesWithSettings[state], message)
		}
	}

	for state, messages := range messagesWithSettings {
		if _, ok := translatedStates[state]; !ok {
			return nil, errors.Errorf("unknown state %s in messages %v", state, messages)
		}
	}

	return translatedStates, nil
}

type settedStateGroup map[string]string

func translateMessages(
	messages []*parser.Message,
	declaredIdentifiers DeclaredIdentifierGroup,
	dependencies commands.Dependencies,
) (
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

		translatedCommands, settedState, err :=
			translateCommands(message.Commands, declaredIdentifiers, dependencies)
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

func translateCommands(
	commands []*parser.Command,
	declaredIdentifiers DeclaredIdentifierGroup,
	dependencies commands.Dependencies,
) (
	translatedCommands runtime.CommandGroup,
	settedState string,
	err error,
) {
	localDeclaredIdentifiers := make(DeclaredIdentifierGroup)
	for identifier := range declaredIdentifiers {
		localDeclaredIdentifiers[identifier] = struct{}{}
	}

	for index, command := range commands {
		translatedCommand, newSettedState, err :=
			translateCommand(command, localDeclaredIdentifiers, dependencies)
		if err != nil {
			return nil, "", errors.Wrapf(err, "unable to translate the command #%d", index)
		}

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

func translateCommand(
	command *parser.Command,
	declaredIdentifiers DeclaredIdentifierGroup,
	dependencies commands.Dependencies,
) (
	translatedCommand runtime.Command,
	settedState string,
	err error,
) {
	switch {
	case command.Let != nil:
		expression, err2 := translateExpression(command.Let.Expression, declaredIdentifiers)
		if err2 != nil {
			return nil, "", errors.Wrap(err2, "unable to translate the let command")
		}

		declaredIdentifiers[command.Let.Identifier] = struct{}{}
		translatedCommand = commands.NewLetCommand(command.Let.Identifier, expression)
	case command.Send != nil:
		translatedCommand = commands.NewSendCommand(*command.Send)
	case command.Set != nil:
		translatedCommand = commands.NewSetCommand(*command.Set)
		settedState = *command.Set
	case command.Out != nil:
		translatedCommand = commands.NewOutCommand(*command.Out, dependencies.OutWriter)
	case command.Sleep != nil:
		if translatedCommand, err = commands.NewSleepCommand(
			*command.Sleep.Minimum,
			*command.Sleep.Maximum,
			dependencies.Sleep,
		); err != nil {
			return nil, "", err
		}
	case command.Exit:
		translatedCommand = commands.ExitCommand{}
	case command.Expression != nil:
		expression, err2 := translateExpression(command.Expression, declaredIdentifiers)
		if err2 != nil {
			return nil, "", errors.Wrap(err2, "unable to translate the expression command")
		}

		translatedCommand = commands.NewExpressionCommand(expression)
	}

	return translatedCommand, settedState, nil
}
