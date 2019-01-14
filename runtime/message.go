package runtime

// MessageGroup ...
type MessageGroup map[string]CommandGroup

// ProcessMessage ...
func (messages MessageGroup) ProcessMessage(context Context, message string) error {
	return messages[message].Run(context)
}
