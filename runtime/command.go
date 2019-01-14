package runtime

import "github.com/pkg/errors"

// Command is implemented by any command. Real type doesn't have to be pure and can store a state.
//go:generate mockery -name=Command -inpkg -case=underscore -testonly
type Command interface {
	Run() error
}

// CommandGroup represents a list of commands in a certain order.
type CommandGroup []Command

// Run executes commands sequentially one by one. If any command fails, execution stops.
func (commands CommandGroup) Run() error {
	for index, command := range commands {
		if err := command.Run(); err != nil {
			return errors.Wrapf(err, "unable to run the command #%d", index)
		}
	}

	return nil
}
