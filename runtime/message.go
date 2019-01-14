package runtime

import "github.com/thewizardplusplus/tick-tock/runtime/context"

// MessageGroup ...
type MessageGroup map[string]CommandGroup

// ProcessMessage ...
// TODO: wrap the error with the message name.
func (messages MessageGroup) ProcessMessage(context context.Context, message string) error {
	return messages[message].Run(context)
}
