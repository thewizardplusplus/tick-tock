package tests

import (
	"io"
	"sync"

	"github.com/spf13/afero"
	"github.com/thewizardplusplus/tick-tock/runtime/waiter"
)

// ...
const (
	UnbufferedInbox = iota
	BufferedInbox
)

// Reader ...
//go:generate mockery -name=Reader -case=underscore
type Reader interface {
	io.Reader
}

// Writer ...
//go:generate mockery -name=Writer -case=underscore
type Writer interface {
	io.Writer
}

// FileSystem ...
//go:generate mockery -name=FileSystem -case=underscore
type FileSystem interface {
	afero.Fs
}

// File ...
//go:generate mockery -name=File -case=underscore
type File interface {
	afero.File
}

// Exiter ...
//go:generate mockery -name=Exiter -case=underscore
type Exiter interface {
	Exit(code int)
}

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

// GetAddress ...
func GetAddress(s string) *string {
	return &s
}
