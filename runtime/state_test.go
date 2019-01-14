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
		makeStates func(log *[]int) StateGroup
		args       args
		wantLog    []int
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			makeStates: func(log *[]int) StateGroup {
				return StateGroup{
					"state_one": MessageGroup{
						"message_one": newLoggableCommands(log, 5, 0),
						"message_two": newLoggableCommands(log, 5, 5),
					},
					"state_two": MessageGroup{
						"message_three": newLoggableCommands(log, 5, 10),
						"message_four":  newCalledLoggableCommands(log, 5, 15, -1),
					},
				}
			},
			args:    args{"state_two", "message_four"},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name:       "error with an empty group",
			makeStates: func(log *[]int) StateGroup { return nil },
			args:       args{"state_unknown", "message_unknown"},
			wantErr:    assert.Error,
		},
		{
			name: "error with an unknown state",
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
			args:    args{"state_unknown", "message_unknown"},
			wantErr: assert.Error,
		},
		{
			name: "error on command execution",
			makeStates: func(log *[]int) StateGroup {
				return StateGroup{
					"state_one": MessageGroup{
						"message_one": newLoggableCommands(log, 5, 0),
						"message_two": newLoggableCommands(log, 5, 5),
					},
					"state_two": MessageGroup{
						"message_three": newLoggableCommands(log, 5, 10),
						"message_four":  newCalledLoggableCommands(log, 5, 15, 2),
					},
				}
			},
			args:    args{"state_two", "message_four"},
			wantLog: []int{15, 16, 17},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log []int
			states := testData.makeStates(&log)
			err := states.ProcessMessage(testData.args.state, testData.args.message)

			assert.Equal(test, testData.wantLog, log)
			for _, messages := range states {
				checkMessages(test, messages)
			}
			testData.wantErr(test, err)
		})
	}
}