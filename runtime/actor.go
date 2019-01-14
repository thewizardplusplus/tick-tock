package runtime

// Actor keeps a state map of an actor and its current one.
type Actor struct {
	currentState string
	states       StateGroup
}
