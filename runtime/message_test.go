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
		context context.Context
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
				context: new(mocks.Context),
				message: context.Message{Name: "unknown"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with an unknown message",
			makeMessages: func(context context.Context, log *commandLog) MessageGroup {
				return newLoggableMessages(context, log, group(2), group(5), nil)
			},
			args: args{
				context: new(mocks.Context),
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
				context: new(mocks.Context),
				message: context.Message{Name: "message_1"},
			},
			wantLog: []int{5, 6, 7, 8, 9},
			wantErr: assert.NoError,
		},
		{
			name: "success with message arguments",
			makeMessages: func(context context.Context, log *commandLog) MessageGroup {
				return newLoggableMessages(context, log, group(2), group(5), loggableCommandOptions{
					"message_1": {withParameters([]string{"one", "two"}), withCalls()},
				})
			},
			args: args{
				context: func() context.Context {
					context := new(mocks.Context)
					context.On("SetValue", "one", 23).Return()
					context.On("SetValue", "two", 42).Return()

					return context
				}(),
				message: context.Message{
					Name:      "message_1",
					Arguments: []interface{}{23, 42},
				},
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
				context: new(mocks.Context),
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
				context: new(mocks.Context),
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
				context: new(mocks.Context),
				message: context.Message{Name: "message_1"},
			},
			wantLog: []int{5, 6, 7},
			wantErr: assert.NoError,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log commandLog
			messages := testData.makeMessages(testData.args.context, &log)
			err := messages.ProcessMessage(testData.args.context, testData.args.message)

			mock.AssertExpectationsForObjects(test, testData.args.context)
			checkMessages(test, messages)
			assert.Equal(test, testData.wantLog, log.commands)
			testData.wantErr(test, err)
		})
	}
}

func TestParameterizedMessageGroup(test *testing.T) {
	type fields struct {
		parameters   []string
		makeMessages func(context context.Context, log *commandLog) MessageGroup
	}
	type args struct {
		context   context.Context
		arguments []interface{}
		message   context.Message
	}

	for _, testData := range []struct {
		name    string
		fields  fields
		args    args
		wantLog []int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success with state arguments",
			fields: fields{
				parameters: []string{"one", "two"},
				makeMessages: func(context context.Context, log *commandLog) MessageGroup {
					return newLoggableMessages(context, log, group(2), group(5), loggableCommandOptions{
						"message_1": {withCalls()},
					})
				},
			},
			args: args{
				context: func() context.Context {
					context := new(mocks.Context)
					context.On("SetValue", "one", 23).Return()
					context.On("SetValue", "two", 42).Return()

					return context
				}(),
				arguments: []interface{}{23, 42},
				message:   context.Message{Name: "message_1"},
			},
			wantLog: []int{5, 6, 7, 8, 9},
			wantErr: assert.NoError,
		},
		{
			name: "success with message arguments",
			fields: fields{
				parameters: []string{"one", "two"},
				makeMessages: func(context context.Context, log *commandLog) MessageGroup {
					return newLoggableMessages(context, log, group(2), group(5), loggableCommandOptions{
						"message_1": {withParameters([]string{"two", "three"}), withCalls()},
					})
				},
			},
			args: args{
				context: func() context.Context {
					context := new(mocks.Context)
					context.On("SetValue", "one", 5).Return()
					context.On("SetValue", "two", 12).Return()
					context.On("SetValue", "two", 23).Return()
					context.On("SetValue", "three", 42).Return()

					return context
				}(),
				arguments: []interface{}{5, 12},
				message: context.Message{
					Name:      "message_1",
					Arguments: []interface{}{23, 42},
				},
			},
			wantLog: []int{5, 6, 7, 8, 9},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				parameters: []string{"one", "two"},
				makeMessages: func(context context.Context, log *commandLog) MessageGroup {
					return newLoggableMessages(context, log, group(2), group(5), loggableCommandOptions{
						"message_1": {withErrOn(2)},
					})
				},
			},
			args: args{
				context: func() context.Context {
					context := new(mocks.Context)
					context.On("SetValue", "one", 23).Return()
					context.On("SetValue", "two", 42).Return()

					return context
				}(),
				arguments: []interface{}{23, 42},
				message:   context.Message{Name: "message_1"},
			},
			wantLog: []int{5, 6, 7},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log commandLog
			messages := testData.fields.makeMessages(testData.args.context, &log)
			parameterizedMessages := NewParameterizedMessageGroup(testData.fields.parameters, messages)
			err := parameterizedMessages.ParameterizedProcessMessage(
				testData.args.context,
				testData.args.arguments,
				testData.args.message,
			)

			mock.AssertExpectationsForObjects(test, testData.args.context)
			checkMessages(test, messages)
			assert.Equal(test, testData.wantLog, log.commands)
			testData.wantErr(test, err)
		})
	}
}
