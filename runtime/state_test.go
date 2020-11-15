package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

func TestStateGroup(test *testing.T) {
	type args struct {
		context context.Context
		state   context.State
		message context.Message
	}

	for _, testData := range []struct {
		name       string
		makeStates func(context context.Context, log *commandLog) StateGroup
		args       args
		wantLog    []int
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			makeStates: func(context context.Context, log *commandLog) StateGroup {
				return newLoggableStates(context, log, 2, group(2), group(5), loggableCommandOptions{
					"message_3": {withCalls()},
				})
			},
			args: args{
				context: new(MockContext),
				state:   context.State{Name: "state_1"},
				message: context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "success with state arguments",
			makeStates: func(context context.Context, log *commandLog) StateGroup {
				messageConfig := parameterizedGroup(2, "one", "two")
				return newLoggableStates(context, log, 2, messageConfig, group(5), loggableCommandOptions{
					"message_3": {withCalls()},
				})
			},
			args: args{
				context: func() context.Context {
					context := new(MockContext)
					context.On("SetValue", "one", 5).Return()
					context.On("SetValue", "two", 12).Return()

					return context
				}(),
				state: context.State{
					Name:      "state_1",
					Arguments: []interface{}{5, 12},
				},
				message: context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "success with message arguments",
			makeStates: func(context context.Context, log *commandLog) StateGroup {
				messageConfig := parameterizedGroup(2, "one", "two")
				return newLoggableStates(context, log, 2, messageConfig, group(5), loggableCommandOptions{
					"message_3": {withParameters([]string{"two", "three"}), withCalls()},
				})
			},
			args: args{
				context: func() context.Context {
					context := new(MockContext)
					context.On("SetValue", "one", 5).Return()
					context.On("SetValue", "two", 12).Return()
					context.On("SetValue", "two", 23).Return()
					context.On("SetValue", "three", 42).Return()

					return context
				}(),
				state: context.State{
					Name:      "state_1",
					Arguments: []interface{}{5, 12},
				},
				message: context.Message{
					Name:      "message_3",
					Arguments: []interface{}{23, 42},
				},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name:       "error with an empty group",
			makeStates: func(context context.Context, log *commandLog) StateGroup { return nil },
			args: args{
				context: new(MockContext),
				state:   context.State{Name: "state_unknown"},
				message: context.Message{Name: "message_unknown"},
			},
			wantLog: nil,
			wantErr: assert.Error,
		},
		{
			name: "error with an unknown state",
			makeStates: func(context context.Context, log *commandLog) StateGroup {
				return newLoggableStates(context, log, 2, group(2), group(5), nil)
			},
			args: args{
				context: new(MockContext),
				state:   context.State{Name: "state_unknown"},
				message: context.Message{Name: "message_unknown"},
			},
			wantLog: nil,
			wantErr: assert.Error,
		},
		{
			name: "error on command execution",
			makeStates: func(context context.Context, log *commandLog) StateGroup {
				return newLoggableStates(context, log, 2, group(2), group(5), loggableCommandOptions{
					"message_3": {withErrOn(2)},
				})
			},
			args: args{
				context: new(MockContext),
				state:   context.State{Name: "state_1"},
				message: context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log commandLog
			states := testData.makeStates(testData.args.context, &log)
			err := states.ProcessMessage(testData.args.context, testData.args.state, testData.args.message)

			mock.AssertExpectationsForObjects(test, testData.args.context)
			checkStates(test, states)
			assert.Equal(test, testData.wantLog, log.commands)
			testData.wantErr(test, err)
		})
	}
}

func TestParameterizedStateGroup(test *testing.T) {
	type fields struct {
		parameters []string
		makeStates func(context context.Context, log *commandLog) StateGroup
	}
	type args struct {
		context   context.Context
		arguments []interface{}
		state     context.State
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
			name: "success with actor arguments",
			fields: fields{
				parameters: []string{"one", "two"},
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, group(2), group(5), loggableCommandOptions{
						"message_3": {withCalls()},
					})
				},
			},
			args: args{
				context: func() context.Context {
					context := new(MockContext)
					context.On("SetValue", "one", 5).Return()
					context.On("SetValue", "two", 12).Return()

					return context
				}(),
				arguments: []interface{}{5, 12},
				state:     context.State{Name: "state_1"},
				message:   context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "success with state arguments",
			fields: fields{
				parameters: []string{"one", "two"},
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					messageConfig := parameterizedGroup(2, "two", "three")
					return newLoggableStates(context, log, 2, messageConfig, group(5), loggableCommandOptions{
						"message_3": {withCalls()},
					})
				},
			},
			args: args{
				context: func() context.Context {
					context := new(MockContext)
					context.On("SetValue", "one", 5).Return()
					context.On("SetValue", "two", 12).Return()
					context.On("SetValue", "two", 23).Return()
					context.On("SetValue", "three", 42).Return()

					return context
				}(),
				arguments: []interface{}{5, 12},
				state: context.State{
					Name:      "state_1",
					Arguments: []interface{}{23, 42},
				},
				message: context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "success with message arguments",
			fields: fields{
				parameters: []string{"one", "two"},
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					messageConfig := parameterizedGroup(2, "two", "three")
					return newLoggableStates(context, log, 2, messageConfig, group(5), loggableCommandOptions{
						"message_3": {withParameters([]string{"three", "four"}), withCalls()},
					})
				},
			},
			args: args{
				context: func() context.Context {
					context := new(MockContext)
					context.On("SetValue", "one", 5).Return()
					context.On("SetValue", "two", 12).Return()
					context.On("SetValue", "two", 23).Return()
					context.On("SetValue", "three", 42).Return()
					context.On("SetValue", "three", 100).Return()
					context.On("SetValue", "four", 1000).Return()

					return context
				}(),
				arguments: []interface{}{5, 12},
				state: context.State{
					Name:      "state_1",
					Arguments: []interface{}{23, 42},
				},
				message: context.Message{
					Name:      "message_3",
					Arguments: []interface{}{100, 1000},
				},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "error with an unknown state",
			fields: fields{
				parameters: []string{"one", "two"},
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, group(2), group(5), nil)
				},
			},
			args: args{
				context:   new(MockContext),
				arguments: []interface{}{5, 12},
				state:     context.State{Name: "state_unknown"},
				message:   context.Message{Name: "message_unknown"},
			},
			wantLog: nil,
			wantErr: assert.Error,
		},
		{
			name: "error on command execution",
			fields: fields{
				parameters: []string{"one", "two"},
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, group(2), group(5), loggableCommandOptions{
						"message_3": {withErrOn(2)},
					})
				},
			},
			args: args{
				context: func() context.Context {
					context := new(MockContext)
					context.On("SetValue", "one", 5).Return()
					context.On("SetValue", "two", 12).Return()

					return context
				}(),
				arguments: []interface{}{5, 12},
				state:     context.State{Name: "state_1"},
				message:   context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log commandLog
			states := testData.fields.makeStates(testData.args.context, &log)
			parameterizedStates := NewParameterizedStateGroup(testData.fields.parameters, states)
			err := parameterizedStates.ParameterizedProcessMessage(
				testData.args.context,
				testData.args.arguments,
				testData.args.state,
				testData.args.message,
			)

			mock.AssertExpectationsForObjects(test, testData.args.context)
			checkStates(test, states)
			assert.Equal(test, testData.wantLog, log.commands)
			testData.wantErr(test, err)
		})
	}
}
