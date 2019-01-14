package tests

import (
	"io"
	"sync"

	"github.com/spf13/afero"
	"github.com/thewizardplusplus/tick-tock/runtime/mocks"
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
	*mocks.Waiter
	*sync.WaitGroup
}

// NewSynchronousWaiter ...
func NewSynchronousWaiter() SynchronousWaiter {
	return SynchronousWaiter{new(mocks.Waiter), new(sync.WaitGroup)}
}

// Add ...
func (waiter SynchronousWaiter) Add(delta int) {
	waiter.Waiter.Add(delta)
	waiter.WaitGroup.Add(delta)
}

// Done ...
func (waiter SynchronousWaiter) Done() {
	waiter.Waiter.Done()
	waiter.WaitGroup.Done()
}

// GetAddress ...
func GetAddress(s string) *string {
	return &s
}
