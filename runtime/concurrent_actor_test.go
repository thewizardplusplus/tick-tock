package runtime

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
		makeStates   func(log *[]int) StateGroup
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
				makeStates: func(log *[]int) StateGroup {
					return StateGroup{
						"state_one": MessageGroup{
							"message_one": newLoggableCommands(log, 5, 0),
							"message_two": newLoggableCommands(log, 5, 5),
						},
						"state_two": MessageGroup{
							"message_three": newCalledLoggableCommands(log, 5, 10, -1),
							"message_four":  newCalledLoggableCommands(log, 5, 15, -1),
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
				makeStates: func(log *[]int) StateGroup {
					return StateGroup{
						"state_one": MessageGroup{
							"message_one": newLoggableCommands(log, 5, 0),
							"message_two": newLoggableCommands(log, 5, 5),
						},
						"state_two": MessageGroup{
							"message_three": newCalledLoggableCommands(log, 5, 10, -1),
							"message_four":  newCalledLoggableCommands(log, 5, 15, -1),
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
				makeStates: func(log *[]int) StateGroup {
					return StateGroup{
						"state_one": MessageGroup{
							"message_one": newLoggableCommands(log, 5, 0),
							"message_two": newLoggableCommands(log, 5, 5),
						},
						"state_two": MessageGroup{
							"message_three": newLoggableCommands(log, 5, 10),
							"message_four":  newLoggableCommands(log, 5, 15),
						},
					}
				},
			},
		},
		{
			name: "error",
			args: args{
				initialState: "state_two",
				makeStates: func(log *[]int) StateGroup {
					return StateGroup{
						"state_one": MessageGroup{
							"message_one": newLoggableCommands(log, 5, 0),
							"message_two": newLoggableCommands(log, 5, 5),
						},
						"state_two": MessageGroup{
							"message_three": newCalledLoggableCommands(log, 5, 10, 2),
							"message_four":  newCalledLoggableCommands(log, 5, 15, 2),
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
			var log []int
			states := testData.args.makeStates(&log)
			actor, err := NewActor(testData.args.initialState, states)
			require.NoError(test, err)

			waiter := synchronousWaiter{new(MockWaiter), new(sync.WaitGroup)}
			if messageCount := len(testData.args.messages); messageCount != 0 {
				waiter.WaitGroup.Add(messageCount)
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
			concurrentActor.Start()

			for _, message := range testData.args.messages {
				concurrentActor.SendMessage(message)
			}
			waiter.Wait()

			assert.ElementsMatch(test, testData.wantLog, log)
			checkStates(test, states)
			waiter.AssertExpectations(test)
			errorHandler.AssertExpectations(test)
		})
	}
}

func TestConcurrentActorGroup(test *testing.T) {
	type args struct {
		initialState string
		makeStates   func(log *[]int) StateGroup
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
					makeStates: func(log *[]int) StateGroup {
						return StateGroup{
							"state_one": MessageGroup{
								"message_one": newLoggableCommands(log, 5, 0),
								"message_two": newLoggableCommands(log, 5, 5),
							},
							"state_two": MessageGroup{
								"message_three": newCalledLoggableCommands(log, 5, 10, -1),
								"message_four":  newCalledLoggableCommands(log, 5, 15, -1),
								"message_five":  newCalledLoggableCommands(log, 5, 20, -1),
							},
						}
					},
				},
				{
					initialState: "state_two",
					makeStates: func(log *[]int) StateGroup {
						return StateGroup{
							"state_one": MessageGroup{
								"message_one": newLoggableCommands(log, 5, 25),
								"message_two": newLoggableCommands(log, 5, 30),
							},
							"state_two": MessageGroup{
								"message_three": newCalledLoggableCommands(log, 5, 35, -1),
								"message_four":  newCalledLoggableCommands(log, 5, 40, -1),
								"message_five":  newCalledLoggableCommands(log, 5, 45, -1),
							},
						}
					},
				},
			},
			messages: []string{"message_three", "message_four", "message_five"},
			wantLog: []int{
				10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24,
				35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49,
			},
		},
		{
			name:     "success without actors",
			messages: []string{"message_three", "message_four", "message_five"},
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			waiter := synchronousWaiter{new(MockWaiter), new(sync.WaitGroup)}
			actorCount, messageCount := len(testData.args), len(testData.messages)
			waiter.On("Add", actorCount).Times(messageCount)
			if actorCount != 0 {
				waiter.On("Done").Times(actorCount * messageCount)
			}

			var log []int
			var concurrentActors []ConcurrentActor
			errorHandler := new(MockErrorHandler)
			for _, args := range testData.args {
				states := args.makeStates(&log)
				defer checkStates(test, states)

				actor, err := NewActor(args.initialState, states)
				require.NoError(test, err)

				concurrentActor := NewConcurrentActor(0, actor, Dependencies{waiter, errorHandler})
				concurrentActors = append(concurrentActors, concurrentActor)
			}

			group := NewConcurrentActorGroup(concurrentActors, waiter)
			group.Start()
			for _, message := range testData.messages {
				group.SendMessage(message)
			}
			waiter.Wait()

			assert.ElementsMatch(test, testData.wantLog, log)
			waiter.AssertExpectations(test)
			errorHandler.AssertExpectations(test)
		})
	}
}
