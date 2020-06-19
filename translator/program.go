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
	messagesWithSettingsByStates := make(map[string][]string)
	for _, state := range states {
		if _, ok := translatedStates[state.Name]; ok {
			return nil, errors.Errorf("duplicate state %s", state.Name)
		}

		translatedMessages, settedStatesByMessages, err :=
			translateMessages(state.Messages, declaredIdentifiers)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to translate the state %s", state.Name)
		}

		translatedStates[state.Name] = translatedMessages
		for message, settedStates := range settedStatesByMessages {
			for _, state := range settedStates.ToSlice() {
				messagesWithSettingsByStates[state.(string)] = append(
					messagesWithSettingsByStates[state.(string)],
					message,
				)
			}
		}
	}

	for state, messages := range messagesWithSettingsByStates {
		if _, ok := translatedStates[state]; !ok {
			return nil, errors.Errorf("unknown state %s in messages %v", state, messages)
		}
	}

	return translatedStates, nil
}

type settedStateGroup map[string]mapset.Set

func translateMessages(messages []*parser.Message, declaredIdentifiers mapset.Set) (
	translatedMessages runtime.MessageGroup,
	settedStatesByMessages settedStateGroup,
	err error,
) {
	translatedMessages = make(runtime.MessageGroup)
	settedStatesByMessages = make(settedStateGroup)
	for _, message := range messages {
		if _, ok := translatedMessages[message.Name]; ok {
			return nil, nil, errors.Errorf("duplicate message %s", message.Name)
		}

		translatedCommands, settedStates, err := translateCommands(message.Commands, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to translate the message %s", message.Name)
		}

		translatedMessages[message.Name] =
			runtime.NewParameterizedCommandGroup(message.Parameters, translatedCommands)
		settedStatesByMessages[message.Name] = settedStates
	}

	return translatedMessages, settedStatesByMessages, nil
}

func translateCommands(commands []*parser.Command, declaredIdentifiers mapset.Set) (
	translatedCommands runtime.CommandGroup,
	settedStates mapset.Set,
	err error,
) {
	localDeclaredIdentifiers := declaredIdentifiers.Clone()
	settedStates = mapset.NewSet()
	var topLevelSettedState string
	for index, command := range commands {
		translatedCommand, topLevelSettedState2, settedStates2, didReturn, err :=
			translateCommand(command, localDeclaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to translate the command #%d", index)
		}
		if didReturn && index != len(commands)-1 {
			return nil, nil, errors.Errorf("unreachable commands after the command #%d", index)
		}

		translatedCommands = append(translatedCommands, translatedCommand)
		settedStates = settedStates.Union(settedStates2)

		if len(topLevelSettedState2) == 0 {
			continue
		}
		if len(topLevelSettedState) != 0 {
			return nil, nil, errors.Errorf(
				"second set command %s (first was %s)",
				topLevelSettedState2,
				topLevelSettedState,
			)
		}
		topLevelSettedState = topLevelSettedState2
	}

	return translatedCommands, settedStates, nil
}

func translateCommand(command *parser.Command, declaredIdentifiers mapset.Set) (
	translatedCommand runtime.Command,
	topLevelSettedState string,
	settedStates mapset.Set,
	didReturn bool,
	err error,
) {
	settedStates = mapset.NewSet()
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
		translatedCommand = commands.NewSendCommand(command.Send.Name)
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

	return translatedCommand, topLevelSettedState, settedStates, didReturn, nil
}
