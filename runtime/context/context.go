package context

// Message ...
type Message struct {
	Name      string
	Arguments []interface{}
}

//go:generate mockery -name=MessageSender -inpkg -case=underscore -testonly

// MessageSender ...
type MessageSender interface {
	SendMessage(message Message)
}

// State ...
type State struct {
	Name      string
	Arguments []interface{}
}

//go:generate mockery -name=StateHolder -inpkg -case=underscore -testonly

// StateHolder ...
type StateHolder interface {
	SetState(state State) error
}

// Actor ...
type Actor interface {
	MessageSender

	Start(context Context, arguments []interface{})
}

// ActorRegister ...
//go:generate mockery -name=ActorRegister -case=underscore
type ActorRegister interface {
	RegisterActor(actor Actor, arguments []interface{})
}

// Context ...
//go:generate mockery -name=Context -case=underscore
type Context interface {
	MessageSender
	StateHolder
	ActorRegister
	ValueStore

	SetMessageSender(sender MessageSender)
	SetStateHolder(holder StateHolder)
	SetActorRegister(register ActorRegister)
	SetValueStore(store CopyableValueStore)
	Copy() Context
}

// DefaultContext ...
type DefaultContext struct {
	MessageSender
	StateHolder
	ActorRegister
	CopyableValueStore
}

// NewDefaultContext ...
func NewDefaultContext() *DefaultContext {
	valueStore := make(DefaultValueStore)
	return &DefaultContext{CopyableValueStore: valueStore}
}

// SetMessageSender ...
func (context *DefaultContext) SetMessageSender(sender MessageSender) {
	context.MessageSender = sender
}

// SetStateHolder ...
func (context *DefaultContext) SetStateHolder(holder StateHolder) {
	context.StateHolder = holder
}

// SetActorRegister ...
func (context *DefaultContext) SetActorRegister(register ActorRegister) {
	context.ActorRegister = register
}

// SetValueStore ...
func (context *DefaultContext) SetValueStore(store CopyableValueStore) {
	context.CopyableValueStore = store
}

// Copy ...
func (context DefaultContext) Copy() Context {
	return &DefaultContext{
		MessageSender:      context.MessageSender,
		StateHolder:        context.StateHolder,
		ActorRegister:      context.ActorRegister,
		CopyableValueStore: context.CopyableValueStore.Copy(),
	}
}
