package runtime

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// StateGroup ...
type StateGroup map[string]ParameterizedMessageGroup

// Contains ...
func (states StateGroup) Contains(state context.State) bool {
	_, ok := states[state.Name]
	return ok
}

// ProcessMessage ...
func (states StateGroup) ProcessMessage(
	context context.Context,
	state context.State,
	message context.Message,
) error {
	if !states.Contains(state) {
		return newUnknownStateError(state)
	}

	err := states[state.Name].ParameterizedProcessMessage(context, state.Arguments, message)
	if err != nil {
		return errors.Wrapf(err, "unable to process the state %s", state.Name)
	}

	return nil
}

// ParameterizedStateGroup ...
type ParameterizedStateGroup struct {
	StateGroup

	parameters []string
}

// NewParameterizedStateGroup ...
func NewParameterizedStateGroup(parameters []string, states StateGroup) ParameterizedStateGroup {
	return ParameterizedStateGroup{states, parameters}
}

// ParameterizedProcessMessage ...
func (parameterizedStates ParameterizedStateGroup) ParameterizedProcessMessage(
	ctx context.Context,
	arguments []interface{},
	state context.State,
	message context.Message,
) error {
	if !parameterizedStates.Contains(state) {
		return newUnknownStateError(state)
	}

	values := context.ZipValues(parameterizedStates.parameters, arguments)
	context.SetValues(ctx, values)

	if err := parameterizedStates.StateGroup.ProcessMessage(ctx, state, message); err != nil {
		return errors.Wrap(err, "unable to process parameterized states")
	}

	return nil
}
