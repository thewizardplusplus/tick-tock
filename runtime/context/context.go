package context

// MessageSender ...
type MessageSender interface {
	SendMessage(message string)
}
