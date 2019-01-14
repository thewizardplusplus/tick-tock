package tests

import (
	"sync"

	"github.com/thewizardplusplus/tick-tock/runtime/waiter"
)

// SynchronousWaiter ...
type SynchronousWaiter struct {
	mock   waiter.Waiter
	syncer *sync.WaitGroup
}

// NewSynchronousWaiter ...
func NewSynchronousWaiter(waiter waiter.Waiter) SynchronousWaiter {
	return SynchronousWaiter{waiter, new(sync.WaitGroup)}
}

// Add ...
func (waiter SynchronousWaiter) Add(delta int) {
	waiter.mock.Add(delta)
	waiter.syncer.Add(delta)
}

// Done ...
func (waiter SynchronousWaiter) Done() {
	waiter.mock.Done()
	waiter.syncer.Done()
}

// Wait ...
func (waiter SynchronousWaiter) Wait() {
	waiter.syncer.Wait()
}
