package runtime

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// MessageGroup ...
type MessageGroup map[string]CommandGroup

// ProcessMessage ...
func (messages MessageGroup) ProcessMessage(context context.Context, message string) error {
	if _, err := messages[message].Run(context); err != nil && errors.Cause(err) != ErrReturn {
		return errors.Wrapf(err, "unable to process the message %s", message)
	}

	return nil
}
