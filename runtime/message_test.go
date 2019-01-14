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
					"one": CommandGroup{newLoggableCommand(log, 1), newLoggableCommand(log, 2)},
					"two": CommandGroup{newLoggableCommand(log, 3), newLoggableCommand(log, 4)},
				}
			},
			args:    args{"unknown"},
			wantErr: assert.NoError,
		},
		{
			name: "success with a known message",
			makeMessages: func(log *[]int) MessageGroup {
				messages := MessageGroup{
					"one": CommandGroup{newLoggableCommand(log, 1), newLoggableCommand(log, 2)},
					"two": CommandGroup{newLoggableCommand(log, 3), newLoggableCommand(log, 4)},
				}
				messages["two"][0].(*loggableCommand).On("Run").Return(nil)
				messages["two"][1].(*loggableCommand).On("Run").Return(nil)

				return messages
			},
			args:    args{"two"},
			wantLog: []int{3, 4},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			makeMessages: func(log *[]int) MessageGroup {
				messages := MessageGroup{
					"one": CommandGroup{newLoggableCommand(log, 1), newLoggableCommand(log, 2)},
					"two": CommandGroup{newLoggableCommand(log, 3), newLoggableCommand(log, 4)},
				}
				messages["two"][0].(*loggableCommand).On("Run").Return(iotest.ErrTimeout)

				return messages
			},
			args:    args{"two"},
			wantLog: []int{3},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log []int
			messages := testData.makeMessages(&log)
			err := messages.ProcessMessage(testData.args.message)

			assert.Equal(test, testData.wantLog, log)
			for _, commands := range messages {
				for _, command := range commands {
					command.(*loggableCommand).AssertExpectations(test)
				}
			}
			testData.wantErr(test, err)
		})
	}
}
