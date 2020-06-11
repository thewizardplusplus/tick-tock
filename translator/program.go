package translator

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/commands"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
)

// Options ...
type Options struct {
	InboxSize    int
	InitialState string
}

// Translate ...
func Translate(
	actors []*parser.Actor,
	declaredIdentifiers mapset.Set,
	options Options,
	dependencies runtime.Dependencies,
) (
	translatedActors runtime.ConcurrentActorGroup,
	err error,
) {
	for index, actor := range actors {
		translatedStates, err := translateStates(actor.States, declaredIdentifiers)
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
			dependencies,
		))
	}

	return translatedActors, nil
}

func translateStates(states []*parser.State, declaredIdentifiers mapset.Set) (
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

		translatedMessages, settedStates, err := translateMessages(state.Messages, declaredIdentifiers)
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

func translateMessages(messages []*parser.Message, declaredIdentifiers mapset.Set) (
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

		translatedCommands, settedState, err := translateCommands(message.Commands, declaredIdentifiers)
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

func translateCommands(commands []*parser.Command, declaredIdentifiers mapset.Set) (
	translatedCommands runtime.CommandGroup,
	settedState string,
	err error,
) {
	localDeclaredIdentifiers := declaredIdentifiers.Clone()
	for index, command := range commands {
		translatedCommand, newSettedState, _, didReturn, err :=
			translateCommand(command, localDeclaredIdentifiers)
		if err != nil {
			return nil, "", errors.Wrapf(err, "unable to translate the command #%d", index)
		}
		if didReturn && index != len(commands)-1 {
			return nil, "", errors.Errorf("unreachable commands after the command #%d", index)
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

func translateCommand(command *parser.Command, declaredIdentifiers mapset.Set) (
	translatedCommand runtime.Command,
	topLevelSettedState string,
	settedStates mapset.Set,
	didReturn bool,
	err error,
) {
	switch {
	case command.Let != nil:
		var expression expressions.Expression
		expression, settedStates, err = translateExpression(command.Let.Expression, declaredIdentifiers)
		if err != nil {
			return nil, "", nil, false, errors.Wrap(err, "unable to translate the let command")
		}

		translatedCommand = commands.NewLetCommand(command.Let.Identifier, expression)
		declaredIdentifiers.Add(command.Let.Identifier)
	case command.Send != nil:
		translatedCommand = commands.NewSendCommand(*command.Send)
	case command.Set != nil:
		translatedCommand = commands.NewSetCommand(*command.Set)
		topLevelSettedState = *command.Set
		settedStates = mapset.NewSet(*command.Set)
	case command.Return:
		translatedCommand = commands.ReturnCommand{}
		didReturn = true
	case command.Expression != nil:
		var expression expressions.Expression
		expression, settedStates, err = translateExpression(command.Expression, declaredIdentifiers)
		if err != nil {
			return nil, "", nil, false, errors.Wrap(err, "unable to translate the expression command")
		}

		translatedCommand = commands.NewExpressionCommand(expression)
	}
	if settedStates == nil {
		settedStates = mapset.NewSet()
	}

	return translatedCommand, topLevelSettedState, settedStates, didReturn, nil
}
