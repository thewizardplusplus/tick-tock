package runtime

import "github.com/thewizardplusplus/tick-tock/runtime/context"

// Actor ...
type Actor struct {
	states       StateGroup
	currentState string
}

// NewActor ...
func NewActor(states StateGroup, initialState string) (*Actor, error) {
	if _, ok := states[initialState]; !ok {
		return nil, newUnknownStateError(initialState)
	}

	return &Actor{states, initialState}, nil
}

// SetState ...
func (actor *Actor) SetState(state string) error {
	if _, ok := actor.states[state]; !ok {
		return newUnknownStateError(state)
	}

	actor.currentState = state
	return nil
}

// ProcessMessage ...
func (actor *Actor) ProcessMessage(context context.Context, message context.Message) error {
	context = context.Copy()
	context.SetStateHolder(actor)

	return actor.states.ProcessMessage(context, actor.currentState, message)
}
