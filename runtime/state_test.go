package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestStateGroup(test *testing.T) {
	type args struct {
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
				return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
					"message_3": {withCalls()},
				})
			},
			args: args{
				state:   "state_1",
				message: context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name:       "error with an empty group",
			makeStates: func(context context.Context, log *commandLog) StateGroup { return nil },
			args: args{
				state:   "state_unknown",
				message: context.Message{Name: "message_unknown"},
			},
			wantErr: assert.Error,
		},
		{
			name: "error with an unknown state",
			makeStates: func(context context.Context, log *commandLog) StateGroup {
				return newLoggableStates(context, log, 2, 2, group(5), nil)
			},
			args: args{
				state:   "state_unknown",
				message: context.Message{Name: "message_unknown"},
			},
			wantErr: assert.Error,
		},
		{
			name: "error on command execution",
			makeStates: func(context context.Context, log *commandLog) StateGroup {
				return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
					"message_3": {withErrOn(2)},
				})
			},
			args: args{
				state:   "state_1",
				message: context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			context := new(mocks.Context)
			var log commandLog
			states := testData.makeStates(context, &log)
			err := states.ProcessMessage(context, testData.args.state, testData.args.message)

			mock.AssertExpectationsForObjects(test, context)
			checkStates(test, states)
			assert.Equal(test, testData.wantLog, log.commands)
			testData.wantErr(test, err)
		})
	}
}
