package runtime

import (
	"fmt"
	"sync"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/mocks"
)

type commandLog struct {
	sync.Mutex

	commands []int
}

func (log *commandLog) registerCommand(command int) {
	log.Lock()
	defer log.Unlock()

	log.commands = append(log.commands, command)
}

type loggableCommand struct {
	mock *mocks.Command
	log  *commandLog
	id   int
}

func newLoggableCommand(log *commandLog, id int) loggableCommand {
	return loggableCommand{new(mocks.Command), log, id}
}

func (command loggableCommand) Run(context context.Context) (result interface{}, err error) {
	command.log.registerCommand(command.id)
	return command.mock.Run(context)
}

type groupConfig struct {
	size       int
	idOffset   int
	parameters []string
}

func group(size int, idOffset ...int) groupConfig {
	config := groupConfig{size: size}
	if len(idOffset) != 0 {
		config.idOffset = idOffset[0]
	}

	return config
}

type loggableCommandMode int

const (
	loggableCommandCalls loggableCommandMode = iota + 1
	loggableCommandErr
)

type loggableCommandConfig struct {
	groupConfig

	mode     loggableCommandMode
	err      error
	errIndex int
}

func (config loggableCommandConfig) moddedErrIndex() int {
	if config.mode != loggableCommandErr {
		return -1
	}

	return config.errIndex
}

type loggableCommandOption func(*loggableCommandConfig)

func withParameters(parameters []string) loggableCommandOption {
	return func(config *loggableCommandConfig) { config.parameters = parameters }
}

func withCalls() loggableCommandOption {
	return func(config *loggableCommandConfig) { config.mode = loggableCommandCalls }
}

func withErrOn(index int) loggableCommandOption {
	return func(config *loggableCommandConfig) {
		config.mode = loggableCommandErr
		config.err = iotest.ErrTimeout
		config.errIndex = index
	}
}

func withCustomErrOn(err error, index int) loggableCommandOption {
	return func(config *loggableCommandConfig) {
		config.mode = loggableCommandErr
		config.err = err
		config.errIndex = index
	}
}

func newLoggableCommands(
	context context.Context,
	log *commandLog,
	config groupConfig,
	options ...loggableCommandOption,
) CommandGroup {
	commandConfig := loggableCommandConfig{groupConfig: config}
	for _, option := range options {
		option(&commandConfig)
	}

	var commands CommandGroup
	for i := 0; i < commandConfig.size; i++ {
		command := newLoggableCommand(log, i+commandConfig.idOffset)
		if commandConfig.mode == loggableCommandCalls || i <= commandConfig.moddedErrIndex() {
			var err error
			if i == commandConfig.moddedErrIndex() {
				err = commandConfig.err
			}

			command.mock.On("Run", context).Return(command.id, err)
		}

		commands = append(commands, command)
	}

	return commands
}

func newLoggableParameterizedCommands(
	context context.Context,
	log *commandLog,
	config groupConfig,
	options ...loggableCommandOption,
) ParameterizedCommandGroup {
	commandConfig := loggableCommandConfig{groupConfig: config}
	for _, option := range options {
		option(&commandConfig)
	}

	commands := newLoggableCommands(context, log, config, options...)
	return NewParameterizedCommandGroup(commandConfig.parameters, commands)
}

type loggableCommandOptions map[string][]loggableCommandOption

func newLoggableMessages(
	context context.Context,
	log *commandLog,
	messageConfig groupConfig,
	commandConfig groupConfig,
	options loggableCommandOptions,
) MessageGroup {
	messages := make(MessageGroup)
	for i := messageConfig.idOffset; i < messageConfig.idOffset+messageConfig.size; i++ {
		message := fmt.Sprintf("message_%d", i)
		config := group(commandConfig.size, i*commandConfig.size+commandConfig.idOffset)
		messages[message] = newLoggableParameterizedCommands(context, log, config, options[message]...)
	}

	return messages
}

func newLoggableParameterizedMessages(
	context context.Context,
	log *commandLog,
	messageConfig groupConfig,
	commandConfig groupConfig,
	options loggableCommandOptions,
) ParameterizedMessageGroup {
	messages := newLoggableMessages(context, log, messageConfig, commandConfig, options)
	return NewParameterizedMessageGroup(messageConfig.parameters, messages)
}

func newLoggableStates(
	context context.Context,
	log *commandLog,
	stateCount int,
	messageCount int,
	commandConfig groupConfig,
	options loggableCommandOptions,
) StateGroup {
	states := make(StateGroup)
	for i := 0; i < stateCount; i++ {
		state := fmt.Sprintf("state_%d", i)
		config := group(messageCount, i*messageCount)
		states[state] = newLoggableParameterizedMessages(context, log, config, commandConfig, options)
	}

	return states
}

func checkCommands(test *testing.T, commands CommandGroup) {
	for _, command := range commands {
		mock.AssertExpectationsForObjects(test, command.(loggableCommand).mock)
	}
}

func checkMessages(test *testing.T, messages MessageGroup) {
	for _, parameterizedCommands := range messages {
		checkCommands(test, parameterizedCommands.commands)
	}
}

func checkStates(test *testing.T, states StateGroup) {
	for _, parameterizedMessages := range states {
		checkMessages(test, parameterizedMessages.messages)
	}
}
