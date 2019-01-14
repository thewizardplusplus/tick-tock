package runtime

// Actor ...
type Actor struct {
	currentState string
	states       StateGroup
}

// NewActor ...
func NewActor(initialState string, states StateGroup) (*Actor, error) {
	if _, ok := states[initialState]; !ok {
		return nil, newUnknownStateError(initialState)
	}

	return &Actor{initialState, states}, nil
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
func (actor *Actor) ProcessMessage(context Context, message string) error {
	context.SetActor(actor)
	return actor.states.ProcessMessage(context, actor.currentState, message)
}
