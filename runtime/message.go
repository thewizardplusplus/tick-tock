package runtime

// MessageGroup represents a map of message names to command lists.
type MessageGroup map[string]CommandGroup

// ProcessMessage executes a command list corresponding to a certain message.
// It supports empty groups and unknown messages, in both cases nothing happens.
func (messages MessageGroup) ProcessMessage(message string) error {
	return messages[message].Run()
}
