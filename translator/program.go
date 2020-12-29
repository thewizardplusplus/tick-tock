package translator

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/commands"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
)

// Options ...
type Options struct {
	InboxSize    int
	InitialState context.State
}

// TranslateProgram ...
func TranslateProgram(
	program *parser.Program,
	declaredIdentifiers mapset.Set,
	options Options,
	dependencies runtime.Dependencies,
) (
	definitions context.ValueGroup,
	translatedActors []runtime.ConcurrentActorFactory,
	err error,
) {
	definitions = make(context.ValueGroup)
	localDeclaredIdentifiers := declaredIdentifiers.Clone()
	for index, definition := range program.Definitions {
		translatedActorClass, wasActor, err :=
			translateDefinition(definition, localDeclaredIdentifiers, options, dependencies)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to translate the definition #%d", index)
		}

		definitionName := translatedActorClass.Name()
		if _, ok := definitions[definitionName]; ok {
			return nil, nil, errors.Errorf("duplicate definition %s", definitionName)
		}

		definitions[definitionName] = translatedActorClass
		if wasActor {
			translatedActors = append(translatedActors, translatedActorClass)
		}
	}

	return definitions, translatedActors, nil
}

func translateDefinition(
	definition *parser.Definition,
	declaredIdentifiers mapset.Set,
	options Options,
	dependencies runtime.Dependencies,
) (
	translatedActorClass runtime.ConcurrentActorFactory,
	wasActor bool,
	err error,
) {
	switch {
	case definition.Actor != nil:
		actorClass := (*parser.ActorClass)(definition.Actor)
		translatedActorClass, err =
			translateActorClass(actorClass, declaredIdentifiers, options, dependencies)
		if err != nil {
			return runtime.ConcurrentActorFactory{}, false, errors.Wrapf(
				err,
				"unable to translate the actor %s",
				definition.Actor.Name,
			)
		}

		declaredIdentifiers.Add(definition.Actor.Name)
		wasActor = true
	case definition.ActorClass != nil:
		translatedActorClass, err =
			translateActorClass(definition.ActorClass, declaredIdentifiers, options, dependencies)
		if err != nil {
			return runtime.ConcurrentActorFactory{}, false, errors.Wrapf(
				err,
				"unable to translate the actor class %s",
				definition.ActorClass.Name,
			)
		}

		declaredIdentifiers.Add(definition.ActorClass.Name)
	}

	return translatedActorClass, wasActor, nil
}

func translateActorClass(
	actorClass *parser.ActorClass,
	declaredIdentifiers mapset.Set,
	options Options,
	dependencies runtime.Dependencies,
) (
	translatedActorClass runtime.ConcurrentActorFactory,
	err error,
) {
	localDeclaredIdentifiers := declaredIdentifiers.Clone()
	for _, parameter := range actorClass.Parameters {
		localDeclaredIdentifiers.Add(parameter)
	}

	states, err := translateStates(actorClass.States, localDeclaredIdentifiers)
	if err != nil {
		return runtime.ConcurrentActorFactory{}, errors.Wrap(err, "unable to translate states")
	}

	parameterizedStates := runtime.NewParameterizedStateGroup(actorClass.Parameters, states)
	actorFactory, err :=
		runtime.NewActorFactory(actorClass.Name, parameterizedStates, options.InitialState)
	if err != nil {
		return runtime.ConcurrentActorFactory{}, errors.Wrap(err, "unable to construct the factory")
	}

	concurrentActorFactory :=
		runtime.NewConcurrentActorFactory(actorFactory, options.InboxSize, dependencies)
	return concurrentActorFactory, nil
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

		localDeclaredIdentifiers := declaredIdentifiers.Clone()
		for _, parameter := range state.Parameters {
			localDeclaredIdentifiers.Add(parameter)
		}

		translatedMessages, settedStatesByMessages, err :=
			translateMessages(state.Messages, localDeclaredIdentifiers)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to translate the state %s", state.Name)
		}

		translatedStates[state.Name] =
			runtime.NewParameterizedMessageGroup(state.Parameters, translatedMessages)
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

		localDeclaredIdentifiers := declaredIdentifiers.Clone()
		for _, parameter := range message.Parameters {
			localDeclaredIdentifiers.Add(parameter)
		}

		translatedCommands, settedStates, err :=
			translateCommands(message.Commands, localDeclaredIdentifiers)
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
		expression, settedStates, err = TranslateExpression(command.Let.Expression, declaredIdentifiers)
		if err != nil {
			return nil, "", nil, false, errors.Wrap(err, "unable to translate the let command")
		}

		translatedCommand = commands.NewLetCommand(command.Let.Identifier, expression)
		declaredIdentifiers.Add(command.Let.Identifier)
	case command.Start != nil:
		translatedCommand, settedStates, err = translateStartCommand(command.Start, declaredIdentifiers)
		if err != nil {
			return nil, "", nil, false, errors.Wrap(err, "unable to translate the start command")
		}
	case command.Send != nil:
		translatedCommand, settedStates, err = translateSendCommand(command.Send, declaredIdentifiers)
		if err != nil {
			return nil, "", nil, false, errors.Wrap(err, "unable to translate the send command")
		}
	case command.Set != nil:
		translatedCommand, settedStates, err = translateSetCommand(command.Set, declaredIdentifiers)
		if err != nil {
			return nil, "", nil, false, errors.Wrap(err, "unable to translate the set command")
		}

		topLevelSettedState = command.Set.Name
		settedStates.Add(command.Set.Name)
	case command.Return:
		translatedCommand = commands.ReturnCommand{}
		didReturn = true
	case command.Expression != nil:
		var expression expressions.Expression
		expression, settedStates, err = TranslateExpression(command.Expression, declaredIdentifiers)
		if err != nil {
			return nil, "", nil, false, errors.Wrap(err, "unable to translate the expression command")
		}

		translatedCommand = commands.NewExpressionCommand(expression)
	}

	return translatedCommand, topLevelSettedState, settedStates, didReturn, nil
}

func translateStartCommand(startCommand *parser.StartCommand, declaredIdentifiers mapset.Set) (
	translatedCommand runtime.Command,
	settedStates mapset.Set,
	err error,
) {
	var actorFactory expressions.Expression
	switch {
	case startCommand.Name != nil:
		identifier := *startCommand.Name
		if !declaredIdentifiers.Contains(identifier) {
			return nil, nil, errors.Errorf("unknown identifier %s", identifier)
		}

		actorFactory = expressions.NewIdentifier(identifier)
		settedStates = mapset.NewSet()
	case startCommand.Expression != nil:
		actorFactory, settedStates, err =
			TranslateExpression(startCommand.Expression, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to translate the actor class for the start command")
		}
	}

	var arguments []expressions.Expression
	for index, argument := range startCommand.Arguments {
		result, settedStates2, err := TranslateExpression(argument, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(
				err,
				"unable to translate the argument #%d for the start command",
				index,
			)
		}

		arguments = append(arguments, result)
		settedStates = settedStates.Union(settedStates2)
	}

	translatedCommand = commands.NewStartCommand(actorFactory, arguments)
	return translatedCommand, settedStates, nil
}

func translateSendCommand(sendCommand *parser.SendCommand, declaredIdentifiers mapset.Set) (
	translatedCommand runtime.Command,
	settedStates mapset.Set,
	err error,
) {
	var arguments []expressions.Expression
	settedStates = mapset.NewSet()
	for index, argument := range sendCommand.Arguments {
		result, settedStates2, err := TranslateExpression(argument, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(
				err,
				"unable to translate the argument #%d for the send command",
				index,
			)
		}

		arguments = append(arguments, result)
		settedStates = settedStates.Union(settedStates2)
	}

	translatedCommand = commands.NewSendCommand(sendCommand.Name, arguments)
	return translatedCommand, settedStates, nil
}

func translateSetCommand(setCommand *parser.SetCommand, declaredIdentifiers mapset.Set) (
	translatedCommand runtime.Command,
	settedStates mapset.Set,
	err error,
) {
	var arguments []expressions.Expression
	settedStates = mapset.NewSet()
	for index, argument := range setCommand.Arguments {
		result, settedStates2, err := TranslateExpression(argument, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(
				err,
				"unable to translate the argument #%d for the set command",
				index,
			)
		}

		arguments = append(arguments, result)
		settedStates = settedStates.Union(settedStates2)
	}

	translatedCommand = commands.NewSetCommand(setCommand.Name, arguments)
	return translatedCommand, settedStates, nil
}
