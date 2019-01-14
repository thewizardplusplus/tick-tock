package runtime

// StateGroup ...
type StateGroup map[string]MessageGroup

// ProcessMessage ...
func (states StateGroup) ProcessMessage(state string, message string) error {
	messages, ok := states[state]
	if !ok {
		return newUnknownStateError(state)
	}

	return messages.ProcessMessage(nil, message)
}
