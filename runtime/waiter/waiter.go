package waiter

// Waiter ...
//go:generate mockery -name=Waiter -case=underscore
type Waiter interface {
	Add(delta int)
	Done()
}
