package runtime

import (
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/waiter"
)

type inbox chan context.Message

// Dependencies ...
type Dependencies struct {
	waiter.Waiter
	ErrorHandler
}

// ConcurrentActor ...
type ConcurrentActor struct {
	innerActor   *Actor
	inbox        inbox
	dependencies Dependencies
}

// NewConcurrentActor ...
func NewConcurrentActor(actor *Actor, inboxSize int, dependencies Dependencies) ConcurrentActor {
	inbox := make(inbox, inboxSize) // nolint: vetshadow
	return ConcurrentActor{actor, inbox, dependencies}
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
func (actor ConcurrentActor) SendMessage(message context.Message) {
	actor.dependencies.Waiter.Add(1)
	go func() { actor.inbox <- message }()
}

// ConcurrentActorGroup ...
type ConcurrentActorGroup []ConcurrentActor

// Start ...
func (actors ConcurrentActorGroup) Start(context context.Context) {
	context = context.Copy()
	context.SetMessageSender(actors)

	for _, actor := range actors {
		actor.Start(context)
	}
}

// SendMessage ...
func (actors ConcurrentActorGroup) SendMessage(message context.Message) {
	for _, actor := range actors {
		actor.SendMessage(message)
	}
}
