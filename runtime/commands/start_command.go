package commands

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// StartCommand ...
type StartCommand struct {
	actorFactory expressions.Expression
}

// NewStartCommand ...
func NewStartCommand(actorFactory expressions.Expression) StartCommand {
	return StartCommand{actorFactory}
}

// Run ...
func (command StartCommand) Run(context context.Context) (result interface{}, err error) {
	actorFactory, err := command.actorFactory.Evaluate(context)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to evaluate the actor class for the start command")
	}

	typedActorFactory, ok := actorFactory.(runtime.ConcurrentActorFactory)
	if !ok {
		return nil, errors.Errorf(
			"unsupported type %T of the actor class for the start command",
			actorFactory,
		)
	}

	actor := typedActorFactory.CreateActor()
	context.RegisterActor(actor, nil)

	return types.Nil{}, nil
}
