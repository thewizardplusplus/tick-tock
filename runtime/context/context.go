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
	ValueStore

	SetMessageSender(sender MessageSender)
	SetStateHolder(holder StateHolder)
	SetValueStore(store CopyableValueStore)
	Copy() Context
}

// DefaultContext ...
type DefaultContext struct {
	MessageSender
	StateHolder
	CopyableValueStore
}

// SetMessageSender ...
func (context *DefaultContext) SetMessageSender(sender MessageSender) {
	context.MessageSender = sender
}

// SetStateHolder ...
func (context *DefaultContext) SetStateHolder(holder StateHolder) {
	context.StateHolder = holder
}

// SetValueStore ...
func (context *DefaultContext) SetValueStore(store CopyableValueStore) {
	context.CopyableValueStore = store
}

// Copy ...
func (context DefaultContext) Copy() Context {
	valueStoreCopy := context.CopyableValueStore.Copy()
	return &DefaultContext{context.MessageSender, context.StateHolder, valueStoreCopy}
}
