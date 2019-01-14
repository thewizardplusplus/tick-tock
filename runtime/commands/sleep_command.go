package commands

import (
	"time"

	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// Randomizer ...
type Randomizer func() float64

// Sleeper ...
type Sleeper func(duration time.Duration)

// SleepDependencies ...
type SleepDependencies struct {
	Randomizer
	Sleeper
}

// SleepCommand ...
type SleepCommand struct {
	minimum      float64
	maximum      float64
	dependencies SleepDependencies
}

// NewSleepCommand ...
func NewSleepCommand(minimum float64, maximum float64, dependencies SleepDependencies) (
	SleepCommand,
	error,
) {
	if minimum < 0 {
		return SleepCommand{}, errors.New("negative minimum")
	}
	if maximum < 0 {
		return SleepCommand{}, errors.New("negative maximum")
	}
	if maximum < minimum {
		return SleepCommand{}, errors.New("maximum less minimum")
	}

	return SleepCommand{minimum, maximum, dependencies}, nil
}

// Run ...
func (command SleepCommand) Run(context context.Context) error {
	delay := command.dependencies.Randomizer()*(command.maximum-command.minimum) + command.minimum
	command.dependencies.Sleeper(time.Duration(delay * float64(time.Second)))

	return nil
}
