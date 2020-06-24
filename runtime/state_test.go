package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

func TestStateGroup(test *testing.T) {
	type args struct {
		context context.Context
		state   string
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
				context: new(mocks.Context),
				state:   "state_1",
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
					context := new(mocks.Context)
					context.On("SetValue", "one", types.Nil{}).Return()
					context.On("SetValue", "two", types.Nil{}).Return()

					return context
				}(),
				state:   "state_1",
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
					context := new(mocks.Context)
					context.On("SetValue", "one", types.Nil{}).Return()
					context.On("SetValue", "two", types.Nil{}).Return()
					context.On("SetValue", "two", 23).Return()
					context.On("SetValue", "three", 42).Return()

					return context
				}(),
				state: "state_1",
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
				context: new(mocks.Context),
				state:   "state_unknown",
				message: context.Message{Name: "message_unknown"},
			},
			wantErr: assert.Error,
		},
		{
			name: "error with an unknown state",
			makeStates: func(context context.Context, log *commandLog) StateGroup {
				return newLoggableStates(context, log, 2, group(2), group(5), nil)
			},
			args: args{
				context: new(mocks.Context),
				state:   "state_unknown",
				message: context.Message{Name: "message_unknown"},
			},
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
				context: new(mocks.Context),
				state:   "state_1",
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
