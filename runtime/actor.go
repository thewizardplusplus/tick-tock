package runtime

import (
	"fmt"

	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// Actor ...
type Actor struct {
	states       StateGroup
	currentState context.State
}

// SetState ...
func (actor *Actor) SetState(state context.State) error {
	if !actor.states.Contains(state) {
		return newUnknownStateError(state)
	}

	actor.currentState = state
	return nil
}

// ProcessMessage ...
func (actor Actor) ProcessMessage(context context.Context, message context.Message) error {
	return actor.states.ProcessMessage(context, actor.currentState, message)
}

// ActorFactory ...
type ActorFactory struct {
	name         string
	states       StateGroup
	initialState context.State
}

// NewActorFactory ...
func NewActorFactory(
	name string,
	states StateGroup,
	initialState context.State,
) (ActorFactory, error) {
	if !states.Contains(initialState) {
		return ActorFactory{}, newUnknownStateError(initialState)
	}

	return ActorFactory{name, states, initialState}, nil
}

// Name ...
func (factory ActorFactory) Name() string {
	return factory.name
}

// Name ...
func (factory ActorFactory) String() string {
	return fmt.Sprintf("<class %s>", factory.name)
}

// MarshalText ...
func (factory ActorFactory) MarshalText() (text []byte, err error) {
	return []byte(factory.String()), nil
}

// CreateActor ...
func (factory ActorFactory) CreateActor() *Actor {
	return &Actor{factory.states, factory.initialState}
}
