package runtime

import (
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

// Waiter ...
//go:generate mockery -name=Waiter -case=underscore
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
func (actor ConcurrentActor) Start(context context.Context) {
	go func() {
		for message := range actor.inbox {
			if err := actor.innerActor.ProcessMessage(context, message); err != nil {
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
func (actors ConcurrentActorGroup) Start(context context.Context) {
	context.SetMessageSender(actors)

	for _, actor := range actors {
		actor.Start(context)
	}
}

// SendMessage ...
func (actors ConcurrentActorGroup) SendMessage(message string) {
	for _, actor := range actors {
		actor.SendMessage(message)
	}
}
