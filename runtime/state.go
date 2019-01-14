package runtime

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// StateGroup ...
type StateGroup map[string]MessageGroup

// ProcessMessage ...
func (states StateGroup) ProcessMessage(
	context context.Context,
	state string,
	message string,
) error {
	messages, ok := states[state]
	if !ok {
		return newUnknownStateError(state)
	}

	if err := messages.ProcessMessage(context, message); err != nil {
		return errors.Wrapf(err, "unable to process the state %s", state)
	}

	return nil
}
