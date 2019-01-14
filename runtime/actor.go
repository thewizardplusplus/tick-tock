package runtime

// Actor keeps a state map of an actor and its current one.
type Actor struct {
	currentState string
	states       StateGroup
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
