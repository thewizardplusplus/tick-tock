package runtime

// Command is implemented by any command. Real type doesn't have to be pure and can store a state.
type Command interface {
	Run() error
}
