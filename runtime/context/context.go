package context

// MessageSender ...
type MessageSender interface {
	SendMessage(message string)
}

// StateHolder ...
type StateHolder interface {
	SetState(state string) error
}

// Context ...
//go:generate mockery -name=Context -case=underscore
type Context interface {
	SendMessage(message string)
	SetState(state string) error
	SetMessageSender(sender MessageSender)
	SetStateHolder(holder StateHolder)
}
