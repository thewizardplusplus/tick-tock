package runtime

// StateGroup ...
type StateGroup map[string]MessageGroup

// ProcessMessage ...
func (states StateGroup) ProcessMessage(context Context, state string, message string) error {
	messages, ok := states[state]
	if !ok {
		return newUnknownStateError(state)
	}

	return messages.ProcessMessage(context, message)
}
