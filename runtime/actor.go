package runtime

import "github.com/thewizardplusplus/tick-tock/runtime/context"

// Actor ...
type Actor struct {
	states       StateGroup
	currentState context.State
}

// NewActor ...
func NewActor(states StateGroup, initialState context.State) (*Actor, error) {
	if _, ok := states[initialState.Name]; !ok {
		return nil, newUnknownStateError(initialState.Name)
	}

	return &Actor{states, initialState}, nil
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
	context = context.Copy()
	context.SetStateHolder(actor)

	return actor.states.ProcessMessage(context, actor.currentState, message)
}
