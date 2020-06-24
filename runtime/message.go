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

// ParameterizedMessageGroup ...
type ParameterizedMessageGroup struct {
	parameters []string
	messages   MessageGroup
}

// NewParameterizedMessageGroup ...
func NewParameterizedMessageGroup(
	parameters []string,
	messages MessageGroup,
) ParameterizedMessageGroup {
	return ParameterizedMessageGroup{parameters, messages}
}

// ParameterizedProcessMessage ...
func (parameterizedMessages ParameterizedMessageGroup) ParameterizedProcessMessage(
	ctx context.Context,
	arguments []interface{},
	message context.Message,
) error {
	values := context.ZipValues(parameterizedMessages.parameters, arguments)
	context.SetValues(ctx, values)

	if err := parameterizedMessages.messages.ProcessMessage(ctx, message); err != nil {
		return errors.Wrap(err, "unable to process parameterized messages")
	}

	return nil
}
