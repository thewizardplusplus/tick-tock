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
		makeMessages func(context Context, log *commandLog) MessageGroup
		args         args
		wantLog      []int
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name:         "success with an empty group",
			makeMessages: func(context Context, log *commandLog) MessageGroup { return nil },
			args:         args{"two"},
			wantErr:      assert.NoError,
		},
		{
			name: "success with an unknown message",
			makeMessages: func(context Context, log *commandLog) MessageGroup {
				return newLoggableMessages(context, log, 2, 0, 5, 0, nil)
			},
			args:    args{"unknown"},
			wantErr: assert.NoError,
		},
		{
			name: "success with a known message",
			makeMessages: func(context Context, log *commandLog) MessageGroup {
				return newLoggableMessages(context, log, 2, 0, 5, 0, loggableCommandOptions{
					"message_1": {withCalls()},
				})
			},
			args:    args{"message_1"},
			wantLog: []int{5, 6, 7, 8, 9},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			makeMessages: func(context Context, log *commandLog) MessageGroup {
				return newLoggableMessages(context, log, 2, 0, 5, 0, loggableCommandOptions{
					"message_1": {withErrOn(2)},
				})
			},
			args:    args{"message_1"},
			wantLog: []int{5, 6, 7},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			context := new(MockContext)
			var log commandLog
			messages := testData.makeMessages(context, &log)
			err := messages.ProcessMessage(context, testData.args.message)

			context.AssertExpectations(test)
			assert.Equal(test, testData.wantLog, log.commands)
			checkMessages(test, messages)
			testData.wantErr(test, err)
		})
	}
}
