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
			name: "success with messages",
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
				waiter.On("Add", 1).Times(messageCount)
				waiter.On("Done").Times(messageCount)
			}

			errorHandler := new(MockErrorHandler)
			if testData.errCount != 0 {
				errorHandler.
					On("HandleError", mock.MatchedBy(func(error) bool { return true })).
					Times(testData.errCount)
			}

			concurrentActor := NewConcurrentActor(actor, waiter, errorHandler)
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
							},
						}
					},
				},
				{
					initialState: "state_two",
					makeStates: func(log *[]int) StateGroup {
						return StateGroup{
							"state_one": MessageGroup{
								"message_one": newLoggableCommands(log, 5, 20),
								"message_two": newLoggableCommands(log, 5, 25),
							},
							"state_two": MessageGroup{
								"message_three": newCalledLoggableCommands(log, 5, 30, -1),
								"message_four":  newCalledLoggableCommands(log, 5, 35, -1),
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

			var log []int
			var concurrentActors ConcurrentActorGroup
			errorHandler := new(MockErrorHandler)
			for _, args := range testData.args {
				states := args.makeStates(&log)
				defer checkStates(test, states)

				actor, err := NewActor(args.initialState, states)
				require.NoError(test, err)

				concurrentActor := NewConcurrentActor(actor, waiter, errorHandler)
				concurrentActors = append(concurrentActors, concurrentActor)
			}

			concurrentActors.Start()
			for _, message := range testData.messages {
				concurrentActors.SendMessage(message)
			}
			waiter.Wait()

			assert.ElementsMatch(test, testData.wantLog, log)
			waiter.AssertExpectations(test)
			errorHandler.AssertExpectations(test)
		})
	}
}
