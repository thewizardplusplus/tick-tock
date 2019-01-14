package runtime

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

type synchronousWaiter struct {
	*MockWaiter
	*sync.WaitGroup
}

func (waiter synchronousWaiter) Add(delta int) {
	waiter.MockWaiter.Add(delta)
	waiter.WaitGroup.Add(delta)
}

func (waiter synchronousWaiter) Done() {
	waiter.MockWaiter.Done()
	waiter.WaitGroup.Done()
}

func TestConcurrentActor(test *testing.T) {
	type args struct {
		inboxSize    int
		initialState string
		makeStates   func(context context.Context, log *commandLog) StateGroup
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
				initialState: "state_1",
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_2": {withCalls()},
						"message_3": {withCalls()},
					})
				},
				messages: []string{"message_2", "message_3"},
			},
			wantLog: []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		},
		{
			name: "success with messages (with a buffered inbox)",
			args: args{
				inboxSize:    1,
				initialState: "state_1",
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_2": {withCalls()},
						"message_3": {withCalls()},
					})
				},
				messages: []string{"message_2", "message_3"},
			},
			wantLog: []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		},
		{
			name: "success without messages",
			args: args{
				initialState: "state_1",
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), nil)
				},
			},
		},
		{
			name: "error",
			args: args{
				initialState: "state_1",
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_2": {withErrOn(2)},
						"message_3": {withErrOn(2)},
					})
				},
				messages: []string{"message_2", "message_3"},
			},
			errCount: 2,
			wantLog:  []int{10, 11, 12, 15, 16, 17},
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			actor := &Actor{testData.args.initialState, nil}
			context := new(mocks.Context)
			if len(testData.args.messages) != 0 {
				context.On("SetStateHolder", actor).Return()
			}

			var log commandLog
			actor.states = testData.args.makeStates(context, &log)

			waiter := synchronousWaiter{new(MockWaiter), new(sync.WaitGroup)}
			if messageCount := len(testData.args.messages); messageCount != 0 {
				waiter.On("Add", 1).Times(messageCount)
				waiter.On("Done").Times(messageCount)
			}

			errorHandler := new(MockErrorHandler)
			if testData.errCount != 0 {
				errorHandler.
					On("HandleError", mock.MatchedBy(func(error) bool { return true })).
					Times(testData.errCount)
			}

			dependencies := Dependencies{waiter, errorHandler}
			concurrentActor := NewConcurrentActor(testData.args.inboxSize, actor, dependencies)
			concurrentActor.Start(context)

			for _, message := range testData.args.messages {
				concurrentActor.SendMessage(message)
			}
			waiter.Wait()

			context.AssertExpectations(test)
			assert.ElementsMatch(test, testData.wantLog, log.commands)
			checkStates(test, actor.states)
			waiter.AssertExpectations(test)
			errorHandler.AssertExpectations(test)
		})
	}
}

func TestConcurrentActorGroup(test *testing.T) {
	type args struct {
		initialState string
		makeStates   func(context context.Context, log *commandLog) StateGroup
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
					initialState: "state_1",
					makeStates: func(context context.Context, log *commandLog) StateGroup {
						return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
							"message_2": {withCalls()},
							"message_3": {withCalls()},
						})
					},
				},
				{
					initialState: "state_1",
					makeStates: func(context context.Context, log *commandLog) StateGroup {
						return newLoggableStates(context, log, 2, 2, group(5, 20), loggableCommandOptions{
							"message_2": {withCalls()},
							"message_3": {withCalls()},
						})
					},
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
			waiter := synchronousWaiter{new(MockWaiter), new(sync.WaitGroup)}
			if messageCount := len(testData.args) * len(testData.messages); messageCount != 0 {
				waiter.On("Add", 1).Times(messageCount)
				waiter.On("Done").Times(messageCount)
			}

			context := new(mocks.Context)
			var log commandLog
			var concurrentActors ConcurrentActorGroup
			errorHandler := new(MockErrorHandler)
			for _, args := range testData.args {
				states := args.makeStates(context, &log)
				defer checkStates(test, states)

				actor := &Actor{args.initialState, states}
				context.On("SetStateHolder", actor).Return()

				concurrentActor := NewConcurrentActor(0, actor, Dependencies{waiter, errorHandler})
				concurrentActors = append(concurrentActors, concurrentActor)
			}

			context.On("SetMessageSender", concurrentActors).Return()
			concurrentActors.Start(context)
			for _, message := range testData.messages {
				concurrentActors.SendMessage(message)
			}
			waiter.Wait()

			context.AssertExpectations(test)
			assert.ElementsMatch(test, testData.wantLog, log.commands)
			waiter.AssertExpectations(test)
			errorHandler.AssertExpectations(test)
		})
	}
}
