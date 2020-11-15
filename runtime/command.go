package runtime

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

//go:generate mockery -name=Command -inpkg -case=underscore -testonly

// Command ...
type Command interface {
	Run(context context.Context) (result interface{}, err error)
}

// CommandGroup ...
type CommandGroup []Command

// Run ...
func (commands CommandGroup) Run(context context.Context) (result interface{}, err error) {
	result = types.Nil{}
	for index, command := range commands {
		if result, err = command.Run(context); err != nil {
			return nil, errors.Wrapf(err, "unable to run the command #%d", index)
		}
	}

	return result, nil
}

// ParameterizedCommandGroup ...
type ParameterizedCommandGroup struct {
	parameters []string
	commands   CommandGroup
}

// NewParameterizedCommandGroup ...
func NewParameterizedCommandGroup(
	parameters []string,
	commands CommandGroup,
) ParameterizedCommandGroup {
	return ParameterizedCommandGroup{parameters, commands}
}

// ParameterizedRun ...
func (parameterizedCommands ParameterizedCommandGroup) ParameterizedRun(
	ctx context.Context,
	arguments []interface{},
) (result interface{}, err error) {
	values := context.ZipValues(parameterizedCommands.parameters, arguments)
	context.SetValues(ctx, values)

	result, err = parameterizedCommands.commands.Run(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "unable to run parameterized commands")
	}

	return result, nil
}
