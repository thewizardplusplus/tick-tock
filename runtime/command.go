package runtime

import "github.com/pkg/errors"

// Context ...
//go:generate mockery -name=Context -inpkg -case=underscore -testonly
type Context interface {
	SendMessage(message string)
	SetState(state string) error
	SetActor(actor *Actor)
	SetActors(actors ConcurrentActorGroup)
}

// Command ...
//go:generate mockery -name=Command -inpkg -case=underscore -testonly
type Command interface {
	Run(context Context) error
}

// CommandGroup ...
type CommandGroup []Command

// Run ...
func (commands CommandGroup) Run(context Context) error {
	for index, command := range commands {
		if err := command.Run(context); err != nil {
			return errors.Wrapf(err, "unable to run the command #%d", index)
		}
	}

	return nil
}
