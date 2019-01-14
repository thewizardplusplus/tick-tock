package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateGroup_ProcessMessage(test *testing.T) {
	type args struct {
		state   string
		message string
	}

	for _, testData := range []struct {
		name       string
		makeStates func(context Context, log *commandLog) StateGroup
		args       args
		wantLog    []int
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			makeStates: func(context Context, log *commandLog) StateGroup {
				return StateGroup{
					"state_one": MessageGroup{
						"message_one": newLoggableCommands(context, log, 5),
						"message_two": newLoggableCommands(context, log, 5, withIDFrom(5)),
					},
					"state_two": MessageGroup{
						"message_three": newLoggableCommands(context, log, 5, withIDFrom(10)),
						"message_four":  newLoggableCommands(context, log, 5, withIDFrom(15), withCalls()),
					},
				}
			},
			args:    args{"state_two", "message_four"},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name:       "error with an empty group",
			makeStates: func(context Context, log *commandLog) StateGroup { return nil },
			args:       args{"state_unknown", "message_unknown"},
			wantErr:    assert.Error,
		},
		{
			name: "error with an unknown state",
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
			args:    args{"state_unknown", "message_unknown"},
			wantErr: assert.Error,
		},
		{
			name: "error on command execution",
			makeStates: func(context Context, log *commandLog) StateGroup {
				return StateGroup{
					"state_one": MessageGroup{
						"message_one": newLoggableCommands(context, log, 5),
						"message_two": newLoggableCommands(context, log, 5, withIDFrom(5)),
					},
					"state_two": MessageGroup{
						"message_three": newLoggableCommands(context, log, 5, withIDFrom(10)),
						"message_four":  newLoggableCommands(context, log, 5, withIDFrom(15), withErrOn(2)),
					},
				}
			},
			args:    args{"state_two", "message_four"},
			wantLog: []int{15, 16, 17},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			context := new(MockContext)
			var log commandLog
			states := testData.makeStates(context, &log)
			err := states.ProcessMessage(context, testData.args.state, testData.args.message)

			context.AssertExpectations(test)
			assert.Equal(test, testData.wantLog, log.commands)
			checkStates(test, states)
			testData.wantErr(test, err)
		})
	}
}
