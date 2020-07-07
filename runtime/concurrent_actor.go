package runtime

import (
	"sync"

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

// Start ...
func (actor ConcurrentActor) Start(context context.Context) {
	for message := range actor.inbox {
		if err := actor.innerActor.ProcessMessage(context, message); err != nil {
			actor.dependencies.ErrorHandler.HandleError(err)
		}

		actor.dependencies.Waiter.Done()
	}
}

// SendMessage ...
func (actor ConcurrentActor) SendMessage(message context.Message) {
	// waiter increment should call synchronously
	// otherwise the program may end before all messages are processed
	actor.dependencies.Waiter.Add(1)

	// simulate an unbounded channel
	go func() { actor.inbox <- message }()
}

// ConcurrentActorFactory ...
type ConcurrentActorFactory struct {
	actorFactory ActorFactory
	inboxSize    int
	dependencies Dependencies
}

// NewConcurrentActorFactory ...
func NewConcurrentActorFactory(
	actorFactory ActorFactory,
	inboxSize int,
	dependencies Dependencies,
) ConcurrentActorFactory {
	return ConcurrentActorFactory{actorFactory, inboxSize, dependencies}
}

// CreateActor ...
func (factory ConcurrentActorFactory) CreateActor() ConcurrentActor {
	actor := factory.actorFactory.CreateActor()
	inbox := make(inbox, factory.inboxSize) // nolint: vetshadow
	return ConcurrentActor{actor, inbox, factory.dependencies}
}

// ConcurrentActorGroup ...
type ConcurrentActorGroup struct {
	context context.Context
	locker  sync.RWMutex
	actors  []context.Actor
}

// NewConcurrentActorGroup ...
func NewConcurrentActorGroup(context context.Context) *ConcurrentActorGroup {
	contextCopy := context.Copy()
	group := &ConcurrentActorGroup{context: contextCopy}
	group.context.SetMessageSender(group)
	group.context.SetActorRegister(group)

	return group
}

// RegisterActor ...
func (group *ConcurrentActorGroup) RegisterActor(actor context.Actor) {
	group.locker.Lock()
	defer group.locker.Unlock()

	group.actors = append(group.actors, actor)
	go actor.Start(group.context)
}

// SendMessage ...
func (group *ConcurrentActorGroup) SendMessage(message context.Message) {
	group.locker.RLock()
	defer group.locker.RUnlock()

	for _, actor := range group.actors {
		actor.SendMessage(message)
	}
}
