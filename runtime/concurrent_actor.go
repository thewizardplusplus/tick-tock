package runtime

import (
	"sync"

	syncutils "github.com/thewizardplusplus/go-sync-utils"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

type inbox chan context.Message

// Dependencies ...
type Dependencies struct {
	syncutils.WaitGroup
	ErrorHandler
}

// ConcurrentActor ...
type ConcurrentActor struct {
	innerActor   *Actor
	inbox        inbox
	dependencies Dependencies
}

// Start ...
func (actor ConcurrentActor) Start(context context.Context, arguments []interface{}) {
	context = context.Copy()
	context.SetStateHolder(actor.innerActor)

	for message := range actor.inbox {
		if err := actor.innerActor.ProcessMessage(context.Copy(), arguments, message); err != nil {
			actor.dependencies.ErrorHandler.HandleError(err)
		}

		actor.dependencies.WaitGroup.Done()
	}
}

// SendMessage ...
func (actor ConcurrentActor) SendMessage(message context.Message) {
	// waiter increment should call synchronously
	// otherwise the program may end before all messages are processed
	actor.dependencies.WaitGroup.Add(1)

	// use unbounded sending to avoid a deadlock
	syncutils.UnboundedSend(actor.inbox, message)
}

// ConcurrentActorFactory ...
type ConcurrentActorFactory struct {
	ActorFactory

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
	actor := factory.ActorFactory.CreateActor()
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
	group := &ConcurrentActorGroup{context: context.Copy()}
	group.context.SetMessageSender(group)
	group.context.SetActorRegister(group)

	return group
}

// RegisterActor ...
func (group *ConcurrentActorGroup) RegisterActor(actor context.Actor, arguments []interface{}) {
	group.locker.Lock()
	defer group.locker.Unlock()

	group.actors = append(group.actors, actor)
	go actor.Start(group.context, arguments)
}

// SendMessage ...
func (group *ConcurrentActorGroup) SendMessage(message context.Message) {
	group.locker.RLock()
	defer group.locker.RUnlock()

	for _, actor := range group.actors {
		actor.SendMessage(message)
	}
}
