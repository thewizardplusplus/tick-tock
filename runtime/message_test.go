package runtime

import (
	"testing"
	"testing/iotest"

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
					"one": newLoggableCommands(log, 2, 0),
					"two": newLoggableCommands(log, 2, 2),
				}
			},
			args:    args{"unknown"},
			wantErr: assert.NoError,
		},
		{
			name: "success with a known message",
			makeMessages: func(log *[]int) MessageGroup {
				messages := MessageGroup{
					"one": newLoggableCommands(log, 2, 0),
					"two": newLoggableCommands(log, 2, 2),
				}
				messages["two"][0].(*loggableCommand).On("Run").Return(nil)
				messages["two"][1].(*loggableCommand).On("Run").Return(nil)

				return messages
			},
			args:    args{"two"},
			wantLog: []int{2, 3},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			makeMessages: func(log *[]int) MessageGroup {
				messages := MessageGroup{
					"one": newLoggableCommands(log, 2, 0),
					"two": newLoggableCommands(log, 2, 2),
				}
				messages["two"][0].(*loggableCommand).On("Run").Return(iotest.ErrTimeout)

				return messages
			},
			args:    args{"two"},
			wantLog: []int{2},
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
