package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	testutils "github.com/thewizardplusplus/tick-tock/internal/test-utils"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	runtimemocks "github.com/thewizardplusplus/tick-tock/runtime/mocks"
	waitermocks "github.com/thewizardplusplus/tick-tock/runtime/waiter/mocks"
)

func TestConcurrentActor(test *testing.T) {
	type args struct {
		contextCopy  context.Context
		makeStates   func(context context.Context, log *commandLog) StateGroup
		initialState string
		inboxSize    int
		messages     []context.Message
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
				contextCopy: new(contextmocks.Context),
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_2": {withCalls()},
						"message_3": {withCalls()},
					})
				},
				initialState: "state_1",
				messages: []context.Message{
					{Name: "message_2"},
					{Name: "message_3"},
				},
			},
			wantLog: []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		},
		{
			name: "success with messages (with a buffered inbox)",
			args: args{
				contextCopy: new(contextmocks.Context),
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_2": {withCalls()},
						"message_3": {withCalls()},
					})
				},
				initialState: "state_1",
				inboxSize:    testutils.BufferedInbox,
				messages: []context.Message{
					{Name: "message_2"},
					{Name: "message_3"},
				},
			},
			wantLog: []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		},
		{
			name: "success with message arguments",
			args: args{
				contextCopy: func() context.Context {
					context := new(contextmocks.Context)
					context.On("SetValue", "one", 23).Return()
					context.On("SetValue", "two", 42).Return()

					return context
				}(),
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_2": {withParameters([]string{"one", "two"}), withCalls()},
					})
				},
				initialState: "state_1",
				inboxSize:    testutils.BufferedInbox,
				messages: []context.Message{
					{Name: "message_2", Arguments: []interface{}{23, 42}},
				},
			},
			wantLog: []int{10, 11, 12, 13, 14},
		},
		{
			name: "success without messages",
			args: args{
				contextCopy: new(contextmocks.Context),
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), nil)
				},
				initialState: "state_1",
			},
		},
		{
			name: "error",
			args: args{
				contextCopy: new(contextmocks.Context),
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_2": {withErrOn(2)},
						"message_3": {withErrOn(2)},
					})
				},
				initialState: "state_1",
				messages: []context.Message{
					{Name: "message_2"},
					{Name: "message_3"},
				},
			},
			errCount: 2,
			wantLog:  []int{10, 11, 12, 15, 16, 17},
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			actor := &Actor{nil, testData.args.initialState}
			contextOriginal := new(contextmocks.Context)
			if len(testData.args.messages) != 0 {
				testData.args.contextCopy.(*contextmocks.Context).On("SetStateHolder", actor).Return()
				contextOriginal.On("Copy").Return(testData.args.contextCopy)
			}

			var log commandLog
			actor.states = testData.args.makeStates(testData.args.contextCopy, &log)

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

			synchronousWaiter := testutils.NewSynchronousWaiter(waiter)
			dependencies := Dependencies{synchronousWaiter, errorHandler}
			concurrentActor := NewConcurrentActor(actor, testData.args.inboxSize, dependencies)
			concurrentActor.Start(contextOriginal)

			for _, message := range testData.args.messages {
				concurrentActor.SendMessage(message)
			}
			synchronousWaiter.Wait()

			mock.AssertExpectationsForObjects(
				test,
				testData.args.contextCopy,
				contextOriginal,
				waiter,
				errorHandler,
			)
			checkStates(test, actor.states)
			assert.ElementsMatch(test, testData.wantLog, log.commands)
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
		messages []context.Message
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
			messages: []context.Message{
				{Name: "message_2"},
				{Name: "message_3"},
			},
			wantLog: []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39},
		},
		{
			name: "success without actors",
			messages: []context.Message{
				{Name: "message_2"},
				{Name: "message_3"},
			},
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
			synchronousWaiter := testutils.NewSynchronousWaiter(waiter)
			errorHandler := new(runtimemocks.ErrorHandler)
			dependencies := Dependencies{synchronousWaiter, errorHandler}
			for _, args := range testData.args {
				states := args.makeStates(contextSecondCopy, &log)
				defer checkStates(test, states)

				actor := &Actor{states, args.initialState}
				contextSecondCopy.On("SetStateHolder", actor).Return()

				concurrentActor := NewConcurrentActor(actor, testutils.UnbufferedInbox, dependencies)
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

func TestConcurrentActorGroup_withMessageArguments(test *testing.T) {
	contextSecondCopy := new(contextmocks.Context)
	contextSecondCopy.On("SetValue", "one", 23).Return()
	contextSecondCopy.On("SetValue", "two", 42).Return()

	contextFirstCopy := new(contextmocks.Context)
	contextFirstCopy.On("Copy").Return(contextSecondCopy)

	contextOriginal := new(contextmocks.Context)
	contextOriginal.On("Copy").Return(contextFirstCopy)

	var log commandLog
	states := newLoggableStates(contextSecondCopy, &log, 2, 2, group(5), loggableCommandOptions{
		"message_2": {withParameters([]string{"one", "two"}), withCalls()},
	})
	actor := &Actor{states, "state_1"}
	contextSecondCopy.On("SetStateHolder", actor).Return()

	waiter := new(waitermocks.Waiter)
	waiter.On("Add", 1).Times(1)
	waiter.On("Done").Times(1)

	synchronousWaiter := testutils.NewSynchronousWaiter(waiter)
	errorHandler := new(runtimemocks.ErrorHandler)
	concurrentActor := NewConcurrentActor(actor, testutils.UnbufferedInbox, Dependencies{
		Waiter:       synchronousWaiter,
		ErrorHandler: errorHandler,
	})
	concurrentActors := ConcurrentActorGroup{concurrentActor}
	contextFirstCopy.On("SetMessageSender", concurrentActors).Return()

	concurrentActors.Start(contextOriginal)
	concurrentActors.SendMessage(context.Message{
		Name:      "message_2",
		Arguments: []interface{}{23, 42},
	})

	synchronousWaiter.Wait()

	mock.AssertExpectationsForObjects(
		test,
		contextFirstCopy,
		contextSecondCopy,
		contextOriginal,
		waiter,
		errorHandler,
	)
	checkStates(test, states)
	assert.ElementsMatch(test, []int{10, 11, 12, 13, 14}, log.commands)
}
