package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestNewActor(test *testing.T) {
	type args struct {
		initialState string
		states       StateGroup
	}

	for _, testData := range []struct {
		name    string
		args    args
		want    *Actor
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			args:    args{"state_0", StateGroup{"state_0": nil, "state_1": nil}},
			want:    &Actor{"state_0", StateGroup{"state_0": nil, "state_1": nil}},
			wantErr: assert.NoError,
		},
		{
			name:    "error",
			args:    args{"state_unknown", StateGroup{"state_0": nil, "state_1": nil}},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			got, err := NewActor(testData.args.initialState, testData.args.states)
			assert.Equal(test, testData.want, got)
			testData.wantErr(test, err)
		})
	}
}

func TestActor_SetState(test *testing.T) {
	type (
		fields struct {
			currentState string
			states       StateGroup
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
			fields:           fields{"state_0", StateGroup{"state_0": nil, "state_1": nil}},
			args:             args{"state_1"},
			wantCurrentState: "state_1",
			wantErr:          assert.NoError,
		},
		{
			name:             "success with a same state",
			fields:           fields{"state_0", StateGroup{"state_0": nil, "state_1": nil}},
			args:             args{"state_0"},
			wantCurrentState: "state_0",
			wantErr:          assert.NoError,
		},
		{
			name:             "error",
			fields:           fields{"state_0", StateGroup{"state_0": nil, "state_1": nil}},
			args:             args{"state_unknown"},
			wantCurrentState: "state_0",
			wantErr:          assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			actor := Actor{testData.fields.currentState, testData.fields.states}
			err := actor.SetState(testData.args.state)
			assert.Equal(test, testData.wantCurrentState, actor.currentState)
			testData.wantErr(test, err)
		})
	}
}

func TestActor_ProcessMessage(test *testing.T) {
	type (
		fields struct {
			currentState string
			makeStates   func(context context.Context, log *commandLog) StateGroup
		}
		args struct {
			message string
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
				currentState: "state_1",
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_3": {withCalls()},
					})
				},
			},
			args:    args{"message_3"},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				currentState: "state_1",
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					return newLoggableStates(context, log, 2, 2, group(5), loggableCommandOptions{
						"message_3": {withErrOn(2)},
					})
				},
			},
			args:    args{"message_3"},
			wantLog: []int{15, 16, 17},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			actor := Actor{testData.fields.currentState, nil}
			context := new(mocks.Context)
			context.On("SetStateHolder", &actor).Return()

			var log commandLog
			actor.states = testData.fields.makeStates(context, &log)

			err := actor.ProcessMessage(context, testData.args.message)

			context.AssertExpectations(test)
			assert.Equal(test, testData.wantLog, log.commands)
			checkStates(test, actor.states)
			testData.wantErr(test, err)
		})
	}
}
