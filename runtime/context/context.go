package context

// MessageSender ...
//go:generate mockery -name=MessageSender -case=underscore
type MessageSender interface {
	SendMessage(message string)
}

// StateHolder ...
//go:generate mockery -name=StateHolder -case=underscore
type StateHolder interface {
	SetState(state string) error
}

// Context ...
//go:generate mockery -name=Context -case=underscore
type Context interface {
	MessageSender
	StateHolder

	SetMessageSender(sender MessageSender)
	SetStateHolder(holder StateHolder)
}

// DefaultContext ...
type DefaultContext struct {
	MessageSender
	StateHolder
}

// SetMessageSender ...
func (context *DefaultContext) SetMessageSender(sender MessageSender) {
	context.MessageSender = sender
}

// SetStateHolder ...
func (context *DefaultContext) SetStateHolder(holder StateHolder) {
	context.StateHolder = holder
}
