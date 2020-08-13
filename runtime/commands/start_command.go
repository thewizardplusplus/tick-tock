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
	arguments    []expressions.Expression
}

// NewStartCommand ...
func NewStartCommand(
	actorFactory expressions.Expression,
	arguments []expressions.Expression,
) StartCommand {
	return StartCommand{actorFactory, arguments}
}

// Run ...
func (command StartCommand) Run(context context.Context) (result interface{}, err error) {
	actorFactory, err := command.actorFactory.Evaluate(context)
	if err != nil {
		return nil, errors.Wrap(err, "unable to evaluate the actor class for the start command")
	}

	typedActorFactory, ok := actorFactory.(runtime.ConcurrentActorFactory)
	if !ok {
		return nil, errors.Errorf(
			"unsupported type %T of the actor class for the start command",
			actorFactory,
		)
	}

	var arguments []interface{}
	for index, argument := range command.arguments {
		result, err := argument.Evaluate(context)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"unable to evaluate the argument #%d for the start command",
				index,
			)
		}

		arguments = append(arguments, result)
	}

	actor := typedActorFactory.CreateActor()
	context.RegisterActor(actor, arguments)

	return types.Nil{}, nil
}
