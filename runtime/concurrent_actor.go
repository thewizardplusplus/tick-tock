package runtime

import "github.com/pkg/errors"

// Waiter ...
//go:generate mockery -name=Waiter -inpkg -case=underscore -testonly
type Waiter interface {
	Add(delta int)
	Done()
}

// Dependencies ...
type Dependencies struct {
	Waiter       Waiter
	ErrorHandler ErrorHandler
}

// ConcurrentActor ...
type ConcurrentActor struct {
	inbox        chan string
	innerActor   *Actor
	dependencies Dependencies
}

// NewConcurrentActor ...
func NewConcurrentActor(inboxSize int, actor *Actor, dependencies Dependencies) ConcurrentActor {
	return ConcurrentActor{make(chan string, inboxSize), actor, dependencies}
}

// Start ...
func (actor ConcurrentActor) Start() {
	go func() {
		for message := range actor.inbox {
			if err := actor.innerActor.ProcessMessage(message); err != nil {
				err = errors.Wrapf(err, "unable to process the message %s", message)
				actor.dependencies.ErrorHandler.HandleError(err)
			}

			actor.dependencies.Waiter.Done()
		}
	}()
}

// SendMessage ...
func (actor ConcurrentActor) SendMessage(message string) {
	actor.dependencies.Waiter.Add(1)
	go func() { actor.inbox <- message }()
}

// ConcurrentActorGroup ...
type ConcurrentActorGroup []ConcurrentActor

// Start ...
func (actors ConcurrentActorGroup) Start() {
	for _, actor := range actors {
		actor.Start()
	}
}

// SendMessage ...
func (actors ConcurrentActorGroup) SendMessage(message string) {
	for _, actor := range actors {
		actor.SendMessage(message)
	}
}
