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
	type fields struct {
		makeStates   func(context context.Context, log *commandLog) StateGroup
		currentState context.State
		inboxSize    int
	}
	type args struct {
		contextCopy context.Context
		messages    []context.Message
	}

	for _, testData := range []struct {
		name     string
		fields   fields
		args     args
		errCount int
		wantLog  []int
	}{
		{
			name: "success with messages (with an unbuffered inbox)",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, group(2), group(5), loggableCommandOptions{
						"message_2": {withCalls()},
						"message_3": {withCalls()},
					})
				},
				currentState: context.State{Name: "state_1"},
				inboxSize:    testutils.UnbufferedInbox,
			},
			args: args{
				contextCopy: new(contextmocks.Context),
				messages:    []context.Message{{Name: "message_2"}, {Name: "message_3"}},
			},
			errCount: 0,
			wantLog:  []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		},
		{
			name: "success with messages (with a buffered inbox)",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, group(2), group(5), loggableCommandOptions{
						"message_2": {withCalls()},
						"message_3": {withCalls()},
					})
				},
				currentState: context.State{Name: "state_1"},
				inboxSize:    testutils.BufferedInbox,
			},
			args: args{
				contextCopy: new(contextmocks.Context),
				messages:    []context.Message{{Name: "message_2"}, {Name: "message_3"}},
			},
			errCount: 0,
			wantLog:  []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		},
		{
			name: "success with state arguments",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					messageConfig := parameterizedGroup(2, "one", "two")
					return newLoggableStates(context, log, 2, messageConfig, group(5), loggableCommandOptions{
						"message_2": {withCalls()},
					})
				},
				currentState: context.State{
					Name:      "state_1",
					Arguments: []interface{}{5, 12},
				},
				inboxSize: testutils.UnbufferedInbox,
			},
			args: args{
				contextCopy: func() context.Context {
					context := new(contextmocks.Context)
					context.On("SetValue", "one", 5).Return()
					context.On("SetValue", "two", 12).Return()

					return context
				}(),
				messages: []context.Message{{Name: "message_2"}},
			},
			errCount: 0,
			wantLog:  []int{10, 11, 12, 13, 14},
		},
		{
			name: "success with message arguments",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					messageConfig := parameterizedGroup(2, "one", "two")
					return newLoggableStates(context, log, 2, messageConfig, group(5), loggableCommandOptions{
						"message_2": {withParameters([]string{"two", "three"}), withCalls()},
					})
				},
				currentState: context.State{
					Name:      "state_1",
					Arguments: []interface{}{5, 12},
				},
				inboxSize: testutils.UnbufferedInbox,
			},
			args: args{
				contextCopy: func() context.Context {
					context := new(contextmocks.Context)
					context.On("SetValue", "one", 5).Return()
					context.On("SetValue", "two", 12).Return()
					context.On("SetValue", "two", 23).Return()
					context.On("SetValue", "three", 42).Return()

					return context
				}(),
				messages: []context.Message{
					{
						Name:      "message_2",
						Arguments: []interface{}{23, 42},
					},
				},
			},
			errCount: 0,
			wantLog:  []int{10, 11, 12, 13, 14},
		},
		{
			name: "success without messages",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, group(2), group(5), nil)
				},
				currentState: context.State{Name: "state_1"},
				inboxSize:    testutils.UnbufferedInbox,
			},
			args: args{
				contextCopy: new(contextmocks.Context),
				messages:    nil,
			},
			errCount: 0,
			wantLog:  nil,
		},
		{
			name: "error",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, group(2), group(5), loggableCommandOptions{
						"message_2": {withErrOn(2)},
						"message_3": {withErrOn(2)},
					})
				},
				currentState: context.State{Name: "state_1"},
				inboxSize:    testutils.UnbufferedInbox,
			},
			args: args{
				contextCopy: new(contextmocks.Context),
				messages:    []context.Message{{Name: "message_2"}, {Name: "message_3"}},
			},
			errCount: 2,
			wantLog:  []int{10, 11, 12, 15, 16, 17},
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log commandLog
			states := testData.fields.makeStates(testData.args.contextCopy, &log)
			actor := &Actor{states, testData.fields.currentState}
			if len(testData.args.messages) != 0 {
				testData.args.contextCopy.(*contextmocks.Context).On("SetStateHolder", actor).Return()
			}

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

			contextOriginal := new(contextmocks.Context)
			if len(testData.args.messages) != 0 {
				contextOriginal.On("Copy").Return(testData.args.contextCopy)
			}

			synchronousWaiter := testutils.NewSynchronousWaiter(waiter)
			concurrentActor := NewConcurrentActor(actor, testData.fields.inboxSize, Dependencies{
				Waiter:       synchronousWaiter,
				ErrorHandler: errorHandler,
			})
			concurrentActor.Start(contextOriginal)
			for _, message := range testData.args.messages {
				concurrentActor.SendMessage(message)
			}
			synchronousWaiter.Wait()

			mock.AssertExpectationsForObjects(
				test,
				contextOriginal,
				testData.args.contextCopy,
				waiter,
				errorHandler,
			)
			checkStates(test, states)
			assert.ElementsMatch(test, testData.wantLog, log.commands)
		})
	}
}

func TestConcurrentActorFactory(test *testing.T) {
	actorFactory := ActorFactory{
		states: StateGroup{
			"state_0": ParameterizedMessageGroup{},
			"state_1": ParameterizedMessageGroup{},
		},
		initialState: context.State{Name: "state_0"},
	}
	dependencies := Dependencies{
		Waiter:       new(waitermocks.Waiter),
		ErrorHandler: new(runtimemocks.ErrorHandler),
	}
	factory := NewConcurrentActorFactory(actorFactory, 23, dependencies)
	got := factory.CreateActor()

	mock.AssertExpectationsForObjects(test, dependencies.Waiter, dependencies.ErrorHandler)

	assert.Equal(test, 23, cap(got.inbox))
	got.inbox = nil

	want := ConcurrentActor{
		innerActor: &Actor{
			states: StateGroup{
				"state_0": ParameterizedMessageGroup{},
				"state_1": ParameterizedMessageGroup{},
			},
			currentState: context.State{Name: "state_0"},
		},
		dependencies: dependencies,
	}
	assert.Equal(test, want, got)
}

func TestConcurrentActorGroup(test *testing.T) {
	type fields struct {
		makeStates   func(context context.Context, log *commandLog) StateGroup
		currentState context.State
	}

	for _, testData := range []struct {
		name     string
		fields   []fields
		messages []context.Message
		wantLog  []int
	}{
		{
			name: "success with actors",
			fields: []fields{
				{
					makeStates: func(context context.Context, log *commandLog) StateGroup {
						return newLoggableStates(context, log, 2, group(2), group(5), loggableCommandOptions{
							"message_2": {withCalls()},
							"message_3": {withCalls()},
						})
					},
					currentState: context.State{Name: "state_1"},
				},
				{
					makeStates: func(context context.Context, log *commandLog) StateGroup {
						return newLoggableStates(context, log, 2, group(2), group(5, 20), loggableCommandOptions{
							"message_2": {withCalls()},
							"message_3": {withCalls()},
						})
					},
					currentState: context.State{Name: "state_1"},
				},
			},
			messages: []context.Message{{Name: "message_2"}, {Name: "message_3"}},
			wantLog:  []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39},
		},
		{
			name:     "success without actors",
			fields:   nil,
			messages: []context.Message{{Name: "message_2"}, {Name: "message_3"}},
			wantLog:  nil,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			waiter := new(waitermocks.Waiter)
			if messageCount := len(testData.fields) * len(testData.messages); messageCount != 0 {
				waiter.On("Add", 1).Times(messageCount)
				waiter.On("Done").Times(messageCount)
			}

			contextSecondCopy := new(contextmocks.Context)
			contextFirstCopy := new(contextmocks.Context)
			if len(testData.fields) != 0 {
				contextFirstCopy.On("Copy").Return(contextSecondCopy)
			}

			contextOriginal := new(contextmocks.Context)
			contextOriginal.On("Copy").Return(contextFirstCopy)

			var log commandLog
			var concurrentActors ConcurrentActorGroup
			synchronousWaiter := testutils.NewSynchronousWaiter(waiter)
			errorHandler := new(runtimemocks.ErrorHandler)
			for _, args := range testData.fields {
				states := args.makeStates(contextSecondCopy, &log)
				defer checkStates(test, states)

				actor := &Actor{states, args.currentState}
				contextSecondCopy.On("SetStateHolder", actor).Return()

				concurrentActor := NewConcurrentActor(actor, testutils.UnbufferedInbox, Dependencies{
					Waiter:       synchronousWaiter,
					ErrorHandler: errorHandler,
				})
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
				contextOriginal,
				contextFirstCopy,
				contextSecondCopy,
				waiter,
				errorHandler,
			)
			assert.ElementsMatch(test, testData.wantLog, log.commands)
		})
	}
}

func TestConcurrentActorGroup_withArguments(test *testing.T) {
	type fields struct {
		makeStates   func(context context.Context, log *commandLog) StateGroup
		currentState context.State
	}
	type args struct {
		contextSecondCopy context.Context
		message           context.Message
	}

	for _, testData := range []struct {
		name    string
		fields  fields
		args    args
		wantLog []int
	}{
		{
			name: "success with state arguments",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					messageConfig := parameterizedGroup(2, "one", "two")
					return newLoggableStates(context, log, 2, messageConfig, group(5), loggableCommandOptions{
						"message_2": {withCalls()},
					})
				},
				currentState: context.State{
					Name:      "state_1",
					Arguments: []interface{}{5, 12},
				},
			},
			args: args{
				contextSecondCopy: func() context.Context {
					context := new(contextmocks.Context)
					context.On("SetValue", "one", 5).Return()
					context.On("SetValue", "two", 12).Return()

					return context
				}(),
				message: context.Message{Name: "message_2"},
			},
			wantLog: []int{10, 11, 12, 13, 14},
		},
		{
			name: "success with message arguments",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					messageConfig := parameterizedGroup(2, "one", "two")
					return newLoggableStates(context, log, 2, messageConfig, group(5), loggableCommandOptions{
						"message_2": {withParameters([]string{"two", "three"}), withCalls()},
					})
				},
				currentState: context.State{
					Name:      "state_1",
					Arguments: []interface{}{5, 12},
				},
			},
			args: args{
				contextSecondCopy: func() context.Context {
					context := new(contextmocks.Context)
					context.On("SetValue", "one", 5).Return()
					context.On("SetValue", "two", 12).Return()
					context.On("SetValue", "two", 23).Return()
					context.On("SetValue", "three", 42).Return()

					return context
				}(),
				message: context.Message{
					Name:      "message_2",
					Arguments: []interface{}{23, 42},
				},
			},
			wantLog: []int{10, 11, 12, 13, 14},
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			waiter := new(waitermocks.Waiter)
			waiter.On("Add", 1).Times(1)
			waiter.On("Done").Times(1)

			contextFirstCopy := new(contextmocks.Context)
			contextFirstCopy.On("Copy").Return(testData.args.contextSecondCopy)

			contextOriginal := new(contextmocks.Context)
			contextOriginal.On("Copy").Return(contextFirstCopy)

			var log commandLog
			states := testData.fields.makeStates(testData.args.contextSecondCopy, &log)
			actor := &Actor{states, testData.fields.currentState}
			testData.args.contextSecondCopy.(*contextmocks.Context).On("SetStateHolder", actor).Return()

			synchronousWaiter := testutils.NewSynchronousWaiter(waiter)
			errorHandler := new(runtimemocks.ErrorHandler)
			concurrentActor := NewConcurrentActor(actor, testutils.UnbufferedInbox, Dependencies{
				Waiter:       synchronousWaiter,
				ErrorHandler: errorHandler,
			})
			concurrentActors := ConcurrentActorGroup{concurrentActor}
			contextFirstCopy.On("SetMessageSender", concurrentActors).Return()

			concurrentActors.Start(contextOriginal)
			concurrentActors.SendMessage(testData.args.message)
			synchronousWaiter.Wait()

			mock.AssertExpectationsForObjects(
				test,
				contextOriginal,
				contextFirstCopy,
				testData.args.contextSecondCopy,
				waiter,
				errorHandler,
			)
			checkStates(test, states)
			assert.ElementsMatch(test, testData.wantLog, log.commands)
		})
	}
}
