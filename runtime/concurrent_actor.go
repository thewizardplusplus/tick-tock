package runtime

import "github.com/pkg/errors"

//go:generate mockery -name=ErrorHandler -inpkg -case=underscore -testonly
type ErrorHandler interface {
	HandleError(err error)
}

type ConcurrentActor struct {
	inbox        chan string
	innerActor   *Actor
	errorHandler ErrorHandler
}

func NewConcurrentActor(actor *Actor, errorHandler ErrorHandler) ConcurrentActor {
	return ConcurrentActor{make(chan string), actor, errorHandler}
}

func (actor ConcurrentActor) Start() {
	go func() {
		for message := range actor.inbox {
			if err := actor.innerActor.ProcessMessage(message); err != nil {
				err = errors.Wrapf(err, "unable to process the message %s", message)
				actor.errorHandler.HandleError(err)
			}
		}
	}()
}

func (actor ConcurrentActor) SendMessage(message string) {
	go func() { actor.inbox <- message }()
}
