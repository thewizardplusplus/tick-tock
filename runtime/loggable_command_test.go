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

func newLoggableCommands(log *[]int, count int, idOffset int) CommandGroup {
	var commands CommandGroup
	for i := 0; i < count; i++ {
		commands = append(commands, newLoggableCommand(log, i+idOffset))
	}

	return commands
}

func newCalledLoggableCommands(log *[]int, count int, idOffset int, errIndex int) CommandGroup {
	commands := newLoggableCommands(log, count, idOffset)
	for index, command := range commands {
		// expect execution of all commands from first to failed one, inclusive;
		// error index -1 means a failed command is missing
		if errIndex != -1 && index > errIndex {
			break
		}

		var err error
		// return an error from a failed command
		if index == errIndex {
			err = iotest.ErrTimeout
		}

		command.(*loggableCommand).On("Run").Return(err)
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

func (command *loggableCommand) Run() error {
	*command.log = append(*command.log, command.id)
	return command.MockCommand.Run()
}
