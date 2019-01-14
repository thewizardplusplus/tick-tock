package runtime

import "github.com/pkg/errors"

//go:generate mockery -name=Waiter -inpkg -case=underscore -testonly
type Waiter interface {
	Add(delta int)
	Done()
}

//go:generate mockery -name=ErrorHandler -inpkg -case=underscore -testonly
type ErrorHandler interface {
	HandleError(err error)
}

type ConcurrentActor struct {
	inbox        chan string
	innerActor   *Actor
	errorHandler ErrorHandler
}

func NewConcurrentActor(inboxSize int, actor *Actor, errorHandler ErrorHandler) ConcurrentActor {
	return ConcurrentActor{make(chan string, inboxSize), actor, errorHandler}
}

func (actor ConcurrentActor) Start(waiter Waiter) {
	go func() {
		for message := range actor.inbox {
			if err := actor.innerActor.ProcessMessage(message); err != nil {
				err = errors.Wrapf(err, "unable to process the message %s", message)
				actor.errorHandler.HandleError(err)
			}

			waiter.Done()
		}
	}()
}

func (actor ConcurrentActor) SendMessage(message string) {
	go func() { actor.inbox <- message }()
}

type ConcurrentActorGroup struct {
	actors []ConcurrentActor
	waiter Waiter
}

func NewConcurrentActorGroup(actors []ConcurrentActor, waiter Waiter) ConcurrentActorGroup {
	return ConcurrentActorGroup{actors, waiter}
}

func (group ConcurrentActorGroup) Start() {
	for _, actor := range group.actors {
		actor.Start(group.waiter)
	}
}

func (group ConcurrentActorGroup) SendMessage(message string) {
	group.waiter.Add(len(group.actors))

	for _, actor := range group.actors {
		actor.SendMessage(message)
	}
}
