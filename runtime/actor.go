package runtime

// Actor keeps a state map of an actor and its current one.
type Actor struct {
	currentState string
	states       StateGroup
}

// NewActor creates a new actor with a certain state map and sets which one is initial.
// If the latter isn't contained in a certain state map, it'll cause an error.
func NewActor(initialState string, states StateGroup) (*Actor, error) {
	if _, ok := states[initialState]; !ok {
		return nil, newUnknownStateError(initialState)
	}

	return &Actor{initialState, states}, nil
}

// SetState changes a current state of an actor.
// If a new state isn't contained in a state map of the actor, it'll cause an error.
func (actor *Actor) SetState(state string) error {
	if _, ok := actor.states[state]; !ok {
		return newUnknownStateError(state)
	}

	actor.currentState = state
	return nil
}

// ProcessMessage executes a message map corresponding to a current state of an actor.
func (actor Actor) ProcessMessage(message string) error {
	return actor.states.ProcessMessage(actor.currentState, message)
}
