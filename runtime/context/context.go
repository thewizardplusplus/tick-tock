package context

// MessageSender ...
type MessageSender interface {
	SendMessage(message string)
}

// StateHolder ...
type StateHolder interface {
	SetState(state string) error
}
