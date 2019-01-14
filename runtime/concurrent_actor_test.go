package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/internal/tests"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	runtimemocks "github.com/thewizardplusplus/tick-tock/runtime/mocks"
	waitermocks "github.com/thewizardplusplus/tick-tock/runtime/waiter/mocks"
)

func TestConcurrentActor(test *testing.T) {
	type args struct {
		makeStates   func(context context.Context, log *commandLog) StateGroup
		initialState string
		inboxSize    int
		messages     []string
	}

	for _, testData := range []struct {
		name     string
		args     args
		errCount int
		wantLog  []int
	}{
		{
			name: "success with messages (with an unbuffered inbox)",
			args: args{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_2": {withCalls()},
						"message_3": {withCalls()},
					})
				},
				initialState: "state_1",
				messages:     []string{"message_2", "message_3"},
			},
			wantLog: []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		},
		{
			name: "success with messages (with a buffered inbox)",
			args: args{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_2": {withCalls()},
						"message_3": {withCalls()},
					})
				},
				initialState: "state_1",
				inboxSize:    tests.BufferedInbox,
				messages:     []string{"message_2", "message_3"},
			},
			wantLog: []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		},
		{
			name: "success without messages",
			args: args{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), nil)
				},
				initialState: "state_1",
			},
		},
		{
			name: "error",
			args: args{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_2": {withErrOn(2)},
						"message_3": {withErrOn(2)},
					})
				},
				initialState: "state_1",
				messages:     []string{"message_2", "message_3"},
			},
			errCount: 2,
			wantLog:  []int{10, 11, 12, 15, 16, 17},
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			actor := &Actor{nil, testData.args.initialState}
			contextCopy := new(contextmocks.Context)
			contextOriginal := new(contextmocks.Context)
			if len(testData.args.messages) != 0 {
				contextCopy.On("SetStateHolder", actor).Return()
				contextOriginal.On("Copy").Return(contextCopy)
			}

			var log commandLog
			actor.states = testData.args.makeStates(contextCopy, &log)

			waiter := new(waitermocks.Waiter)
			if messageCount := len(testData.args.messages); messageCount != 0 {
				waiter.On("Add", 1).Times(messageCount)
				waiter.On("Done").Times(messageCount)
			}

			errorHandler := new(runtimemocks.ErrorHandler)
			if testData.errCount != 0 {
				errorHandler.
					On("HandleError", mock.MatchedBy(func(error) bool { return true })).
					Times(testData.errCount)
			}

			synchronousWaiter := tests.NewSynchronousWaiter(waiter)
			dependencies := Dependencies{synchronousWaiter, errorHandler}
			concurrentActor := NewConcurrentActor(actor, testData.args.inboxSize, dependencies)
			concurrentActor.Start(contextOriginal)

			for _, message := range testData.args.messages {
				concurrentActor.SendMessage(message)
			}
			synchronousWaiter.Wait()

			mock.AssertExpectationsForObjects(test, contextCopy, contextOriginal, waiter, errorHandler)
			assert.ElementsMatch(test, testData.wantLog, log.commands)
			checkStates(test, actor.states)
		})
	}
}

func TestConcurrentActorGroup(test *testing.T) {
	type args struct {
		makeStates   func(context context.Context, log *commandLog) StateGroup
		initialState string
	}

	for _, testData := range []struct {
		name     string
		args     []args
		messages []string
		wantLog  []int
	}{
		{
			name: "success with actors",
			args: []args{
				{
					makeStates: func(context context.Context, log *commandLog) StateGroup {
						return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
							"message_2": {withCalls()},
							"message_3": {withCalls()},
						})
					},
					initialState: "state_1",
				},
				{
					makeStates: func(context context.Context, log *commandLog) StateGroup {
						return newLoggableStates(context, log, 2, 2, group(5, 20), loggableCommandOptions{
							"message_2": {withCalls()},
							"message_3": {withCalls()},
						})
					},
					initialState: "state_1",
				},
			},
			messages: []string{"message_2", "message_3"},
			wantLog:  []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39},
		},
		{
			name:     "success without actors",
			messages: []string{"message_2", "message_3"},
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			waiter := new(waitermocks.Waiter)
			if messageCount := len(testData.args) * len(testData.messages); messageCount != 0 {
				waiter.On("Add", 1).Times(messageCount)
				waiter.On("Done").Times(messageCount)
			}

			contextFirstCopy := new(contextmocks.Context)
			contextSecondCopy := new(contextmocks.Context)
			contextOriginal := new(contextmocks.Context)
			contextOriginal.On("Copy").Return(contextFirstCopy)
			if len(testData.args) != 0 {
				contextFirstCopy.On("Copy").Return(contextSecondCopy)
			}

			var log commandLog
			var concurrentActors ConcurrentActorGroup
			synchronousWaiter := tests.NewSynchronousWaiter(waiter)
			errorHandler := new(runtimemocks.ErrorHandler)
			dependencies := Dependencies{synchronousWaiter, errorHandler}
			for _, args := range testData.args {
				states := args.makeStates(contextSecondCopy, &log)
				defer checkStates(test, states)

				actor := &Actor{states, args.initialState}
				contextSecondCopy.On("SetStateHolder", actor).Return()

				concurrentActor := NewConcurrentActor(actor, tests.UnbufferedInbox, dependencies)
				concurrentActors = append(concurrentActors, concurrentActor)
			}

			contextFirstCopy.On("SetMessageSender", concurrentActors).Return()
			concurrentActors.Start(contextOriginal)
			for _, message := range testData.messages {
				concurrentActors.SendMessage(message)
			}
			synchronousWaiter.Wait()

			mock.AssertExpectationsForObjects(
				test,
				contextFirstCopy,
				contextSecondCopy,
				contextOriginal,
				waiter,
				errorHandler,
			)
			assert.ElementsMatch(test, testData.wantLog, log.commands)
		})
	}
}
