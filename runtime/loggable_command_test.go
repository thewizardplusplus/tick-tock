package runtime

import (
	"testing"
	"testing/iotest"
)

type loggableCommand struct {
	MockCommand

	log *[]int
	id  int
}

func newLoggableCommand(log *[]int, id int) *loggableCommand {
	return &loggableCommand{MockCommand{}, log, id}
}

func (command *loggableCommand) Run() error {
	*command.log = append(*command.log, command.id)
	return command.MockCommand.Run()
}

type loggableCommandMode int

const (
	loggableCommandQuiet loggableCommandMode = iota
	loggableCommandCalls
	loggableCommandErr
)

type loggableCommandConfig struct {
	mode     loggableCommandMode
	idOffset int
	errIndex int
}

func (config loggableCommandConfig) moddedErrIndex() int {
	if config.mode != loggableCommandErr {
		return -1
	}

	return config.errIndex
}

type loggableCommandOption func(*loggableCommandConfig)

func withIDFrom(offset int) loggableCommandOption {
	return func(config *loggableCommandConfig) { config.idOffset = offset }
}

func withCalls() loggableCommandOption {
	return func(config *loggableCommandConfig) { config.mode = loggableCommandCalls }
}

func withErrOn(index int) loggableCommandOption {
	return func(config *loggableCommandConfig) {
		config.mode = loggableCommandErr
		config.errIndex = index
	}
}

func newLoggableCommands(log *[]int, count int, options ...loggableCommandOption) CommandGroup {
	var config loggableCommandConfig
	for _, option := range options {
		option(&config)
	}

	var commands CommandGroup
	for i := 0; i < count; i++ {
		command := newLoggableCommand(log, i+config.idOffset)
		if config.mode == loggableCommandCalls || i <= config.moddedErrIndex() {
			var err error
			if i == config.moddedErrIndex() {
				err = iotest.ErrTimeout
			}

			command.On("Run").Return(err)
		}

		commands = append(commands, command)
	}

	return commands
}

func checkCommands(test *testing.T, commands CommandGroup) {
	for _, command := range commands {
		command.(*loggableCommand).AssertExpectations(test)
	}
}

func checkMessages(test *testing.T, messages MessageGroup) {
	for _, commands := range messages {
		checkCommands(test, commands)
	}
}

func checkStates(test *testing.T, states StateGroup) {
	for _, messages := range states {
		checkMessages(test, messages)
	}
}
