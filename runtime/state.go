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
	state string,
	message context.Message,
) error {
	messages, ok := states[state]
	if !ok {
		return newUnknownStateError(state)
	}

	if err := messages.ParameterizedProcessMessage(context, nil, message); err != nil {
		return errors.Wrapf(err, "unable to process the state %s", state)
	}

	return nil
}
