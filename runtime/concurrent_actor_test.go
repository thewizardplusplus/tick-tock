package runtime

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		makeStates   func(context Context, log *commandLog) StateGroup
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
				initialState: "state_two",
				makeStates: func(context Context, log *commandLog) StateGroup {
					return StateGroup{
						"state_one": MessageGroup{
							"message_one": newLoggableCommands(context, log, 5),
							"message_two": newLoggableCommands(context, log, 5, withIDFrom(5)),
						},
						"state_two": MessageGroup{
							"message_three": newLoggableCommands(context, log, 5, withIDFrom(10), withCalls()),
							"message_four":  newLoggableCommands(context, log, 5, withIDFrom(15), withCalls()),
						},
					}
				},
				messages: []string{"message_three", "message_four"},
			},
			wantLog: []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		},
		{
			name: "success with messages (with a buffered inbox)",
			args: args{
				inboxSize:    1,
				initialState: "state_two",
				makeStates: func(context Context, log *commandLog) StateGroup {
					return StateGroup{
						"state_one": MessageGroup{
							"message_one": newLoggableCommands(context, log, 5),
							"message_two": newLoggableCommands(context, log, 5, withIDFrom(5)),
						},
						"state_two": MessageGroup{
							"message_three": newLoggableCommands(context, log, 5, withIDFrom(10), withCalls()),
							"message_four":  newLoggableCommands(context, log, 5, withIDFrom(15), withCalls()),
						},
					}
				},
				messages: []string{"message_three", "message_four"},
			},
			wantLog: []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
		},
		{
			name: "success without messages",
			args: args{
				initialState: "state_two",
				makeStates: func(context Context, log *commandLog) StateGroup {
					return StateGroup{
						"state_one": MessageGroup{
							"message_one": newLoggableCommands(context, log, 5),
							"message_two": newLoggableCommands(context, log, 5, withIDFrom(5)),
						},
						"state_two": MessageGroup{
							"message_three": newLoggableCommands(context, log, 5, withIDFrom(10)),
							"message_four":  newLoggableCommands(context, log, 5, withIDFrom(15)),
						},
					}
				},
			},
		},
		{
			name: "error",
			args: args{
				initialState: "state_two",
				makeStates: func(context Context, log *commandLog) StateGroup {
					return StateGroup{
						"state_one": MessageGroup{
							"message_one": newLoggableCommands(context, log, 5),
							"message_two": newLoggableCommands(context, log, 5, withIDFrom(5)),
						},
						"state_two": MessageGroup{
							"message_three": newLoggableCommands(context, log, 5, withIDFrom(10), withErrOn(2)),
							"message_four":  newLoggableCommands(context, log, 5, withIDFrom(15), withErrOn(2)),
						},
					}
				},
				messages: []string{"message_three", "message_four"},
			},
			errCount: 2,
			wantLog:  []int{10, 11, 12, 15, 16, 17},
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			actor := &Actor{testData.args.initialState, nil}
			context := new(MockContext)
			if len(testData.args.messages) != 0 {
				context.On("SetActor", actor).Return()
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
		makeStates   func(context Context, log *commandLog) StateGroup
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
					initialState: "state_two",
					makeStates: func(context Context, log *commandLog) StateGroup {
						return StateGroup{
							"state_one": MessageGroup{
								"message_one": newLoggableCommands(context, log, 5),
								"message_two": newLoggableCommands(context, log, 5, withIDFrom(5)),
							},
							"state_two": MessageGroup{
								"message_three": newLoggableCommands(context, log, 5, withIDFrom(10), withCalls()),
								"message_four":  newLoggableCommands(context, log, 5, withIDFrom(15), withCalls()),
							},
						}
					},
				},
				{
					initialState: "state_two",
					makeStates: func(context Context, log *commandLog) StateGroup {
						return StateGroup{
							"state_one": MessageGroup{
								"message_one": newLoggableCommands(context, log, 5, withIDFrom(20)),
								"message_two": newLoggableCommands(context, log, 5, withIDFrom(25)),
							},
							"state_two": MessageGroup{
								"message_three": newLoggableCommands(context, log, 5, withIDFrom(30), withCalls()),
								"message_four":  newLoggableCommands(context, log, 5, withIDFrom(35), withCalls()),
							},
						}
					},
				},
			},
			messages: []string{"message_three", "message_four"},
			wantLog:  []int{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39},
		},
		{
			name:     "success without actors",
			messages: []string{"message_three", "message_four"},
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			waiter := synchronousWaiter{new(MockWaiter), new(sync.WaitGroup)}
			if messageCount := len(testData.args) * len(testData.messages); messageCount != 0 {
				waiter.On("Add", 1).Times(messageCount)
				waiter.On("Done").Times(messageCount)
			}

			context := new(MockContext)
			var log commandLog
			var concurrentActors ConcurrentActorGroup
			errorHandler := new(MockErrorHandler)
			for _, args := range testData.args {
				states := args.makeStates(context, &log)
				defer checkStates(test, states)

				actor := &Actor{args.initialState, states}
				context.On("SetActor", actor).Return()

				concurrentActor := NewConcurrentActor(0, actor, Dependencies{waiter, errorHandler})
				concurrentActors = append(concurrentActors, concurrentActor)
			}

			context.On("SetActors", concurrentActors).Return()
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
