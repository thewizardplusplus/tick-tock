package runtime

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// MessageGroup ...
type MessageGroup map[string]ParameterizedCommandGroup

// ProcessMessage ...
func (messages MessageGroup) ProcessMessage(
	context context.Context,
	message context.Message,
) error {
	_, err := messages[message.Name].ParameterizedRun(context, message.Arguments)
	if err != nil && errors.Cause(err) != ErrReturn {
		return errors.Wrapf(err, "unable to process the message %s", message)
	}

	return nil
}
