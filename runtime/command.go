package runtime

import "github.com/pkg/errors"

// Context ...
//go:generate mockery -name=Context -inpkg -case=underscore -testonly
type Context interface {
	SendMessage(message string)
	SetState(state string) error
}

// Command ...
//go:generate mockery -name=Command -inpkg -case=underscore -testonly
type Command interface {
	Run() error
}

// CommandGroup ...
type CommandGroup []Command

// Run ...
func (commands CommandGroup) Run() error {
	for index, command := range commands {
		if err := command.Run(); err != nil {
			return errors.Wrapf(err, "unable to run the command #%d", index)
		}
	}

	return nil
}
