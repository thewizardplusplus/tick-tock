package runtime

import (
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// Actor ...
type Actor struct {
	states       StateGroup
	currentState context.State
}

// SetState ...
func (actor *Actor) SetState(state context.State) error {
	if _, ok := actor.states[state.Name]; !ok {
		return newUnknownStateError(state.Name)
	}

	actor.currentState = state
	return nil
}

// ProcessMessage ...
func (actor *Actor) ProcessMessage(context context.Context, message context.Message) error {
	contextCopy := context.Copy()
	contextCopy.SetStateHolder(actor)

	return actor.states.ProcessMessage(contextCopy, actor.currentState, message)
}

// ActorFactory ...
type ActorFactory struct {
	states       StateGroup
	initialState context.State
}

// NewActorFactory ...
func NewActorFactory(states StateGroup, initialState context.State) (ActorFactory, error) {
	if _, ok := states[initialState.Name]; !ok {
		return ActorFactory{}, newUnknownStateError(initialState.Name)
	}

	return ActorFactory{states, initialState}, nil
}

// CreateActor ...
func (factory ActorFactory) CreateActor() *Actor {
	return &Actor{factory.states, factory.initialState}
}
