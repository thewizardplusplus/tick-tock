package runtime

// StateGroup represents a map of states names to message maps.
type StateGroup map[string]MessageGroup

// ProcessMessage executes a message map corresponding to a certain state.
// If the state is unknown, it'll cause an error.
func (states StateGroup) ProcessMessage(state string, message string) error {
	messages, ok := states[state]
	if !ok {
		return newUnknownStateError(state)
	}

	return messages.ProcessMessage(message)
}
