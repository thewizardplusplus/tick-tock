package runtime

import "github.com/thewizardplusplus/tick-tock/runtime/context"

// StateGroup ...
type StateGroup map[string]MessageGroup

// ProcessMessage ...
// TODO: wrap the unknownStateError error with the method name.
// TODO: wrap the ProcessMessage() error with the state name.
func (states StateGroup) ProcessMessage(
	context context.Context,
	state string,
	message string,
) error {
	messages, ok := states[state]
	if !ok {
		return newUnknownStateError(state)
	}

	return messages.ProcessMessage(context, message)
}
