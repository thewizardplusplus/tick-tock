package runtime

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestMessageGroup(test *testing.T) {
	type args struct {
		message context.Message
	}

	for _, testData := range []struct {
		name         string
		makeMessages func(context context.Context, log *commandLog) MessageGroup
		args         args
		wantLog      []int
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name:         "success with an empty group",
			makeMessages: func(context context.Context, log *commandLog) MessageGroup { return nil },
			args: args{
				message: context.Message{Name: "two"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with an unknown message",
			makeMessages: func(context context.Context, log *commandLog) MessageGroup {
				return newLoggableMessages(context, log, group(2), group(5), nil)
			},
			args: args{
				message: context.Message{Name: "unknown"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with a known message",
			makeMessages: func(context context.Context, log *commandLog) MessageGroup {
				return newLoggableMessages(context, log, group(2), group(5), loggableCommandOptions{
					"message_1": {withCalls()},
				})
			},
			args: args{
				message: context.Message{Name: "message_1"},
			},
			wantLog: []int{5, 6, 7, 8, 9},
			wantErr: assert.NoError,
		},
		{
			name: "common error",
			makeMessages: func(context context.Context, log *commandLog) MessageGroup {
				return newLoggableMessages(context, log, group(2), group(5), loggableCommandOptions{
					"message_1": {withErrOn(2)},
				})
			},
			args: args{
				message: context.Message{Name: "message_1"},
			},
			wantLog: []int{5, 6, 7},
			wantErr: assert.Error,
		},
		{
			name: "direct return error",
			makeMessages: func(context context.Context, log *commandLog) MessageGroup {
				return newLoggableMessages(context, log, group(2), group(5), loggableCommandOptions{
					"message_1": {withCustomErrOn(ErrReturn, 2)},
				})
			},
			args: args{
				message: context.Message{Name: "message_1"},
			},
			wantLog: []int{5, 6, 7},
			wantErr: assert.NoError,
		},
		{
			name: "wrapped return error",
			makeMessages: func(context context.Context, log *commandLog) MessageGroup {
				return newLoggableMessages(context, log, group(2), group(5), loggableCommandOptions{
					"message_1": {withCustomErrOn(errors.Wrap(errors.Wrap(ErrReturn, "level #1"), "level #2"), 2)},
				})
			},
			args: args{
				message: context.Message{Name: "message_1"},
			},
			wantLog: []int{5, 6, 7},
			wantErr: assert.NoError,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			context := new(mocks.Context)
			var log commandLog
			messages := testData.makeMessages(context, &log)
			err := messages.ProcessMessage(context, testData.args.message)

			mock.AssertExpectationsForObjects(test, context)
			checkMessages(test, messages)
			assert.Equal(test, testData.wantLog, log.commands)
			testData.wantErr(test, err)
		})
	}
}
