package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestNewActor(test *testing.T) {
	type args struct {
		states       StateGroup
		initialState string
	}

	for _, testData := range []struct {
		name    string
		args    args
		want    *Actor
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			args:    args{StateGroup{"state_0": nil, "state_1": nil}, "state_0"},
			want:    &Actor{StateGroup{"state_0": nil, "state_1": nil}, "state_0"},
			wantErr: assert.NoError,
		},
		{
			name:    "error",
			args:    args{StateGroup{"state_0": nil, "state_1": nil}, "state_unknown"},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			got, err := NewActor(testData.args.states, testData.args.initialState)
			assert.Equal(test, testData.want, got)
			testData.wantErr(test, err)
		})
	}
}

func TestActor_SetState(test *testing.T) {
	type (
		fields struct {
			states       StateGroup
			currentState string
		}
		args struct {
			state string
		}
	)

	for _, testData := range []struct {
		name             string
		fields           fields
		args             args
		wantCurrentState string
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name:             "success with a different state",
			fields:           fields{StateGroup{"state_0": nil, "state_1": nil}, "state_0"},
			args:             args{"state_1"},
			wantCurrentState: "state_1",
			wantErr:          assert.NoError,
		},
		{
			name:             "success with a same state",
			fields:           fields{StateGroup{"state_0": nil, "state_1": nil}, "state_0"},
			args:             args{"state_0"},
			wantCurrentState: "state_0",
			wantErr:          assert.NoError,
		},
		{
			name:             "error",
			fields:           fields{StateGroup{"state_0": nil, "state_1": nil}, "state_0"},
			args:             args{"state_unknown"},
			wantCurrentState: "state_0",
			wantErr:          assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			actor := Actor{testData.fields.states, testData.fields.currentState}
			err := actor.SetState(testData.args.state)
			assert.Equal(test, testData.wantCurrentState, actor.currentState)
			testData.wantErr(test, err)
		})
	}
}

func TestActor_ProcessMessage(test *testing.T) {
	type (
		fields struct {
			makeStates   func(context context.Context, log *commandLog) StateGroup
			currentState string
		}
		args struct {
			contextCopy context.Context
			message     context.Message
		}
	)

	for _, testData := range []struct {
		name    string
		fields  fields
		args    args
		wantLog []int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_3": {withCalls()},
					})
				},
				currentState: "state_1",
			},
			args: args{
				contextCopy: new(mocks.Context),
				message:     context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "success with message arguments",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_3": {withParameters([]string{"one", "two"}), withCalls()},
					})
				},
				currentState: "state_1",
			},
			args: args{
				contextCopy: func() context.Context {
					context := new(mocks.Context)
					context.On("SetValue", "one", 23).Return()
					context.On("SetValue", "two", 42).Return()

					return context
				}(),
				message: context.Message{
					Name:      "message_3",
					Arguments: []interface{}{23, 42},
				},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_3": {withErrOn(2)},
					})
				},
				currentState: "state_1",
			},
			args: args{
				contextCopy: new(mocks.Context),
				message:     context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			contextOriginal := new(mocks.Context)
			contextOriginal.On("Copy").Return(testData.args.contextCopy)

			actor := Actor{nil, testData.fields.currentState}
			testData.args.contextCopy.(*mocks.Context).On("SetStateHolder", &actor).Return()

			var log commandLog
			actor.states = testData.fields.makeStates(testData.args.contextCopy, &log)

			err := actor.ProcessMessage(contextOriginal, testData.args.message)

			mock.AssertExpectationsForObjects(test, testData.args.contextCopy, contextOriginal)
			checkStates(test, actor.states)
			assert.Equal(test, testData.wantLog, log.commands)
			testData.wantErr(test, err)
		})
	}
}
