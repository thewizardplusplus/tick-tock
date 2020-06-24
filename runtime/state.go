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
		return newUnknownStateError(state.Name)
	}

	if err := messages.ParameterizedProcessMessage(context, state.Arguments, message); err != nil {
		return errors.Wrapf(err, "unable to process the state %s", state.Name)
	}

	return nil
}
