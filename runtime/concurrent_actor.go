package runtime

import (
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/waiter"
)

// Dependencies ...
type Dependencies struct {
	waiter.Waiter
	ErrorHandler
}

// ConcurrentActor ...
type ConcurrentActor struct {
	innerActor   *Actor
	inbox        chan string
	dependencies Dependencies
}

// NewConcurrentActor ...
func NewConcurrentActor(actor *Actor, inboxSize int, dependencies Dependencies) ConcurrentActor {
	return ConcurrentActor{actor, make(chan string, inboxSize), dependencies}
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
