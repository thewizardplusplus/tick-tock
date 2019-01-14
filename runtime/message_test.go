package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageGroup_ProcessMessage(test *testing.T) {
	type args struct {
		message string
	}

	for _, testData := range []struct {
		name         string
		makeMessages func(log *[]int) MessageGroup
		args         args
		wantLog      []int
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name:         "success with an empty group",
			makeMessages: func(log *[]int) MessageGroup { return nil },
			args:         args{"two"},
			wantErr:      assert.NoError,
		},
		{
			name: "success with an unknown message",
			makeMessages: func(log *[]int) MessageGroup {
				return MessageGroup{
					"one": newLoggableCommands(log, 5, 0),
					"two": newLoggableCommands(log, 5, 5),
				}
			},
			args:    args{"unknown"},
			wantErr: assert.NoError,
		},
		{
			name: "success with a known message",
			makeMessages: func(log *[]int) MessageGroup {
				return MessageGroup{
					"one": newLoggableCommands(log, 5, 0),
					"two": newCalledLoggableCommands(log, 5, 5, -1),
				}
			},
			args:    args{"two"},
			wantLog: []int{5, 6, 7, 8, 9},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			makeMessages: func(log *[]int) MessageGroup {
				return MessageGroup{
					"one": newLoggableCommands(log, 5, 0),
					"two": newCalledLoggableCommands(log, 5, 5, 2),
				}
			},
			args:    args{"two"},
			wantLog: []int{5, 6, 7},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log []int
			messages := testData.makeMessages(&log)
			err := messages.ProcessMessage(testData.args.message)

			assert.Equal(test, testData.wantLog, log)
			for _, commands := range messages {
				checkLoggableCommands(test, commands)
			}
			testData.wantErr(test, err)
		})
	}
}
