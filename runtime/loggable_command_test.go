package runtime

import "testing"

type loggableCommand struct {
	MockCommand

	log *[]int
	id  int
}

func newLoggableCommand(log *[]int, id int) *loggableCommand {
	return &loggableCommand{MockCommand{}, log, id}
}

func checkLoggableCommands(test *testing.T, commands CommandGroup) {
	for _, command := range commands {
		command.(*loggableCommand).AssertExpectations(test)
	}
}

func (command *loggableCommand) Run() error {
	*command.log = append(*command.log, command.id)
	return command.MockCommand.Run()
}
