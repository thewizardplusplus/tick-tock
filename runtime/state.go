package runtime

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// StateGroup ...
type StateGroup map[string]ParameterizedMessageGroup

// ProcessMessage ...
func (states StateGroup) ProcessMessage(
	context context.Context,
	state context.State,
	message context.Message,
) error {
	messages, ok := states[state.Name]
	if !ok {
		return newUnknownStateError(state)
	}

	if err := messages.ParameterizedProcessMessage(context, state.Arguments, message); err != nil {
		return errors.Wrapf(err, "unable to process the state %s", state.Name)
	}

	return nil
}

// ParameterizedStateGroup ...
type ParameterizedStateGroup struct {
	parameters []string
	states     StateGroup
}

// NewParameterizedStateGroup ...
func NewParameterizedStateGroup(parameters []string, states StateGroup) ParameterizedStateGroup {
	return ParameterizedStateGroup{parameters, states}
}

// ParameterizedProcessMessage ...
func (parameterizedStates ParameterizedStateGroup) ParameterizedProcessMessage(
	ctx context.Context,
	arguments []interface{},
	state context.State,
	message context.Message,
) error {
	if _, ok := parameterizedStates.states[state.Name]; !ok {
		return newUnknownStateError(state)
	}

	values := context.ZipValues(parameterizedStates.parameters, arguments)
	context.SetValues(ctx, values)

	if err := parameterizedStates.states.ProcessMessage(ctx, state, message); err != nil {
		return errors.Wrap(err, "unable to process parameterized states")
	}

	return nil
}
