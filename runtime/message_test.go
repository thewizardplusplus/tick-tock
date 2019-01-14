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
		makeMessages func(log *commandLog) MessageGroup
		args         args
		wantLog      []int
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name:         "success with an empty group",
			makeMessages: func(log *commandLog) MessageGroup { return nil },
			args:         args{"two"},
			wantErr:      assert.NoError,
		},
		{
			name: "success with an unknown message",
			makeMessages: func(log *commandLog) MessageGroup {
				return MessageGroup{
					"one": newLoggableCommands(log, 5),
					"two": newLoggableCommands(log, 5, withIDFrom(5)),
				}
			},
			args:    args{"unknown"},
			wantErr: assert.NoError,
		},
		{
			name: "success with a known message",
			makeMessages: func(log *commandLog) MessageGroup {
				return MessageGroup{
					"one": newLoggableCommands(log, 5),
					"two": newLoggableCommands(log, 5, withIDFrom(5), withCalls()),
				}
			},
			args:    args{"two"},
			wantLog: []int{5, 6, 7, 8, 9},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			makeMessages: func(log *commandLog) MessageGroup {
				return MessageGroup{
					"one": newLoggableCommands(log, 5),
					"two": newLoggableCommands(log, 5, withIDFrom(5), withErrOn(2)),
				}
			},
			args:    args{"two"},
			wantLog: []int{5, 6, 7},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log commandLog
			messages := testData.makeMessages(&log)
			err := messages.ProcessMessage(testData.args.message)

			assert.Equal(test, testData.wantLog, log.commands)
			checkMessages(test, messages)
			testData.wantErr(test, err)
		})
	}
}
