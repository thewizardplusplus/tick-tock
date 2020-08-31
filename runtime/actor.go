package runtime

import (
	"fmt"

	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// Actor ...
type Actor struct {
	states       ParameterizedStateGroup
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
func (actor Actor) ProcessMessage(
	context context.Context,
	arguments []interface{},
	message context.Message,
) error {
	return actor.states.ParameterizedProcessMessage(context, arguments, actor.currentState, message)
}

// ActorFactory ...
type ActorFactory struct {
	name         string
	states       ParameterizedStateGroup
	initialState context.State
}

// NewActorFactory ...
func NewActorFactory(
	name string,
	states ParameterizedStateGroup,
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

// String ...
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
