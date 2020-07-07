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
		inbox        inbox
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
				inbox:        make(inbox, testutils.UnbufferedInbox),
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
				inbox:        make(inbox, testutils.BufferedInbox),
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
				inbox: make(inbox, testutils.UnbufferedInbox),
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
				inbox: make(inbox, testutils.UnbufferedInbox),
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
				inbox:        make(inbox, testutils.UnbufferedInbox),
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
				inbox:        make(inbox, testutils.UnbufferedInbox),
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
			concurrentActor := ConcurrentActor{
				innerActor: actor,
				inbox:      testData.fields.inbox,
				dependencies: Dependencies{
					Waiter:       synchronousWaiter,
					ErrorHandler: errorHandler,
				},
			}
			go concurrentActor.Start(contextOriginal)
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

func TestNewConcurrentActorGroup(test *testing.T) {
	contextFirstCopy := new(contextmocks.Context)
	contextFirstCopy.
		On("SetMessageSender", mock.AnythingOfType("*runtime.ConcurrentActorGroup")).
		Return()
	contextFirstCopy.
		On("SetActorRegister", mock.AnythingOfType("*runtime.ConcurrentActorGroup")).
		Return()

	contextOriginal := new(contextmocks.Context)
	contextOriginal.On("Copy").Return(contextFirstCopy)

	got := NewConcurrentActorGroup(contextOriginal)

	mock.AssertExpectationsForObjects(test, contextOriginal, contextFirstCopy)
	assert.Equal(test, got.context, contextFirstCopy)
	assert.Nil(test, got.actors)
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

			contextFirstCopy := new(contextmocks.Context)
			contextOriginal := new(contextmocks.Context)
			if len(testData.fields) != 0 {
				contextOriginal.On("Copy").Return(contextFirstCopy)
			}

			var log commandLog
			synchronousWaiter := testutils.NewSynchronousWaiter(waiter)
			errorHandler := new(runtimemocks.ErrorHandler)
			concurrentActors := &ConcurrentActorGroup{context: contextOriginal}
			for _, args := range testData.fields {
				states := args.makeStates(contextFirstCopy, &log)
				defer checkStates(test, states)

				actor := &Actor{states, args.currentState}
				contextFirstCopy.On("SetStateHolder", actor).Return()

				concurrentActors.RegisterActor(ConcurrentActor{
					innerActor: actor,
					inbox:      make(inbox, testutils.UnbufferedInbox),
					dependencies: Dependencies{
						Waiter:       synchronousWaiter,
						ErrorHandler: errorHandler,
					},
				})
			}

			for _, message := range testData.messages {
				concurrentActors.SendMessage(message)
			}
			synchronousWaiter.Wait()

			mock.AssertExpectationsForObjects(
				test,
				contextOriginal,
				contextFirstCopy,
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
		contextFirstCopy context.Context
		message          context.Message
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
				contextFirstCopy: func() context.Context {
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
				contextFirstCopy: func() context.Context {
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

			contextOriginal := new(contextmocks.Context)
			contextOriginal.On("Copy").Return(testData.args.contextFirstCopy)

			var log commandLog
			states := testData.fields.makeStates(testData.args.contextFirstCopy, &log)
			actor := &Actor{states, testData.fields.currentState}
			testData.args.contextFirstCopy.(*contextmocks.Context).On("SetStateHolder", actor).Return()

			synchronousWaiter := testutils.NewSynchronousWaiter(waiter)
			errorHandler := new(runtimemocks.ErrorHandler)
			concurrentActors := &ConcurrentActorGroup{context: contextOriginal}
			concurrentActors.RegisterActor(ConcurrentActor{
				innerActor: actor,
				inbox:      make(inbox, testutils.UnbufferedInbox),
				dependencies: Dependencies{
					Waiter:       synchronousWaiter,
					ErrorHandler: errorHandler,
				},
			})
			concurrentActors.SendMessage(testData.args.message)
			synchronousWaiter.Wait()

			mock.AssertExpectationsForObjects(
				test,
				contextOriginal,
				testData.args.contextFirstCopy,
				waiter,
				errorHandler,
			)
			checkStates(test, states)
			assert.ElementsMatch(test, testData.wantLog, log.commands)
		})
	}
}
