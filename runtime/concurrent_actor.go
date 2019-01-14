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
	waiter       Waiter
	errorHandler ErrorHandler
}

func NewConcurrentActor(actor *Actor, waiter Waiter, errorHandler ErrorHandler) ConcurrentActor {
	return ConcurrentActor{make(chan string), actor, waiter, errorHandler}
}

func (actor ConcurrentActor) Start() {
	go func() {
		for message := range actor.inbox {
			if err := actor.innerActor.ProcessMessage(message); err != nil {
				err = errors.Wrapf(err, "unable to process the message %s", message)
				actor.errorHandler.HandleError(err)
			}

			actor.waiter.Done()
		}
	}()
}

func (actor ConcurrentActor) SendMessage(message string) {
	actor.waiter.Add(1)
	go func() { actor.inbox <- message }()
}

type ConcurrentActorGroup []ConcurrentActor

func (actors ConcurrentActorGroup) Start() {
	for _, actor := range actors {
		actor.Start()
	}
}

func (actors ConcurrentActorGroup) SendMessage(message string) {
	for _, actor := range actors {
		actor.SendMessage(message)
	}
}
