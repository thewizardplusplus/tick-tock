package runtime

import (
	"sync"
	"testing"
	"testing/iotest"
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
	MockCommand

	log *commandLog
	id  int
}

func newLoggableCommand(log *commandLog, id int) *loggableCommand {
	return &loggableCommand{MockCommand{}, log, id}
}

func (command *loggableCommand) Run(context Context) error {
	command.log.registerCommand(command.id)
	return command.MockCommand.Run(context)
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

func newLoggableCommands(
	context Context,
	log *commandLog,
	count int,
	options ...loggableCommandOption,
) CommandGroup {
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

			command.On("Run", context).Return(err)
		}

		commands = append(commands, command)
	}

	return commands
}

type loggableCommandOptions map[string][]loggableCommandOption

func newLoggableStates(
	context Context,
	log *commandLog,
	groupSize int,
	idOffset int,
	options loggableCommandOptions,
) StateGroup {
	messages := make(MessageGroup)
	for index, name := range []string{"one", "two", "three", "four"} {
		completedName := "message_" + name
		comletedOptions := append(options[completedName], withIDFrom(index*groupSize+idOffset))
		messages[completedName] = newLoggableCommands(context, log, groupSize, comletedOptions...)
	}

	return StateGroup{
		"state_one": MessageGroup{
			"message_one": messages["message_one"],
			"message_two": messages["message_two"],
		},
		"state_two": MessageGroup{
			"message_three": messages["message_three"],
			"message_four":  messages["message_four"],
		},
	}
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
