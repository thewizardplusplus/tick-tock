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
		initialState context.State
	}

	for _, testData := range []struct {
		name    string
		args    args
		want    *Actor
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				states: StateGroup{
					"state_0": ParameterizedMessageGroup{},
					"state_1": ParameterizedMessageGroup{},
				},
				initialState: context.State{Name: "state_0"},
			},
			want: &Actor{
				states: StateGroup{
					"state_0": ParameterizedMessageGroup{},
					"state_1": ParameterizedMessageGroup{},
				},
				currentState: context.State{Name: "state_0"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				states: StateGroup{
					"state_0": ParameterizedMessageGroup{},
					"state_1": ParameterizedMessageGroup{},
				},
				initialState: context.State{Name: "state_unknown"},
			},
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
			currentState context.State
		}
		args struct {
			state context.State
		}
	)

	for _, testData := range []struct {
		name             string
		fields           fields
		args             args
		wantCurrentState context.State
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "success with a different state",
			fields: fields{
				states: StateGroup{
					"state_0": ParameterizedMessageGroup{},
					"state_1": ParameterizedMessageGroup{},
				},
				currentState: context.State{Name: "state_0"},
			},
			args: args{
				state: context.State{Name: "state_0"},
			},
			wantCurrentState: context.State{Name: "state_0"},
			wantErr:          assert.NoError,
		},
		{
			name: "success with a same state",
			fields: fields{
				states: StateGroup{
					"state_0": ParameterizedMessageGroup{},
					"state_1": ParameterizedMessageGroup{},
				},
				currentState: context.State{Name: "state_0"},
			},
			args: args{
				state: context.State{Name: "state_0"},
			},
			wantCurrentState: context.State{Name: "state_0"},
			wantErr:          assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				states: StateGroup{
					"state_0": ParameterizedMessageGroup{},
					"state_1": ParameterizedMessageGroup{},
				},
				currentState: context.State{Name: "state_0"},
			},
			args: args{
				state: context.State{Name: "state_unknown"},
			},
			wantCurrentState: context.State{Name: "state_0"},
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
			currentState context.State
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
					return newLoggableStates(context, log, 2, group(2), group(5), loggableCommandOptions{
						"message_3": {withCalls()},
					})
				},
				currentState: context.State{Name: "state_1"},
			},
			args: args{
				contextCopy: new(mocks.Context),
				message:     context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "success with state arguments",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					messageConfig := parameterizedGroup(2, "one", "two")
					return newLoggableStates(context, log, 2, messageConfig, group(5), loggableCommandOptions{
						"message_3": {withCalls()},
					})
				},
				currentState: context.State{
					Name:      "state_1",
					Arguments: []interface{}{5, 12},
				},
			},
			args: args{
				contextCopy: func() context.Context {
					context := new(mocks.Context)
					context.On("SetValue", "one", 5).Return()
					context.On("SetValue", "two", 12).Return()

					return context
				}(),
				message: context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "success with message arguments",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) StateGroup {
					messageConfig := parameterizedGroup(2, "one", "two")
					return newLoggableStates(context, log, 2, messageConfig, group(5), loggableCommandOptions{
						"message_3": {withParameters([]string{"two", "three"}), withCalls()},
					})
				},
				currentState: context.State{
					Name:      "state_1",
					Arguments: []interface{}{5, 12},
				},
			},
			args: args{
				contextCopy: func() context.Context {
					context := new(mocks.Context)
					context.On("SetValue", "one", 5).Return()
					context.On("SetValue", "two", 12).Return()
					context.On("SetValue", "two", 23).Return()
					context.On("SetValue", "three", 42).Return()

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
					return newLoggableStates(context, log, 2, group(2), group(5), loggableCommandOptions{
						"message_3": {withErrOn(2)},
					})
				},
				currentState: context.State{Name: "state_1"},
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
