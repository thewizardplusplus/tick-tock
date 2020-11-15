package runtime

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

func TestActor_SetState(test *testing.T) {
	type fields struct {
		states       ParameterizedStateGroup
		currentState context.State
	}
	type args struct {
		state context.State
	}

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
				states:       ParameterizedStateGroup{StateGroup: StateGroup{"state_0": {}, "state_1": {}}},
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
				states:       ParameterizedStateGroup{StateGroup: StateGroup{"state_0": {}, "state_1": {}}},
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
				states:       ParameterizedStateGroup{StateGroup: StateGroup{"state_0": {}, "state_1": {}}},
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
	type fields struct {
		makeStates   func(context context.Context, log *commandLog) ParameterizedStateGroup
		currentState context.State
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
			name: "success",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) ParameterizedStateGroup {
					options := loggableCommandOptions{"message_3": {withCalls()}}
					return newLoggableParameterizedStates(context, log, group(2), group(2), group(5), options)
				},
				currentState: context.State{Name: "state_1"},
			},
			args: args{
				context:   new(MockContext),
				arguments: nil,
				message:   context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "success with actor arguments",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) ParameterizedStateGroup {
					stateConfig := parameterizedGroup(2, "one", "two")
					options := loggableCommandOptions{"message_3": {withCalls()}}
					return newLoggableParameterizedStates(context, log, stateConfig, group(2), group(5), options)
				},
				currentState: context.State{Name: "state_1"},
			},
			args: args{
				context: func() context.Context {
					context := new(MockContext)
					context.On("SetValue", "one", 5).Return()
					context.On("SetValue", "two", 12).Return()

					return context
				}(),
				arguments: []interface{}{5, 12},
				message:   context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "success with state arguments",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) ParameterizedStateGroup {
					return newLoggableParameterizedStates(
						context,
						log,
						parameterizedGroup(2, "one", "two"),
						parameterizedGroup(2, "two", "three"),
						group(5),
						loggableCommandOptions{"message_3": {withCalls()}},
					)
				},
				currentState: context.State{
					Name:      "state_1",
					Arguments: []interface{}{23, 42},
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
				message:   context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "success with message arguments",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) ParameterizedStateGroup {
					return newLoggableParameterizedStates(
						context,
						log,
						parameterizedGroup(2, "one", "two"),
						parameterizedGroup(2, "two", "three"),
						group(5),
						loggableCommandOptions{
							"message_3": {withParameters([]string{"three", "four"}), withCalls()},
						},
					)
				},
				currentState: context.State{
					Name:      "state_1",
					Arguments: []interface{}{23, 42},
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
				message: context.Message{
					Name:      "message_3",
					Arguments: []interface{}{100, 1000},
				},
			},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				makeStates: func(context context.Context, log *commandLog) ParameterizedStateGroup {
					options := loggableCommandOptions{"message_3": {withErrOn(2)}}
					return newLoggableParameterizedStates(context, log, group(2), group(2), group(5), options)
				},
				currentState: context.State{Name: "state_1"},
			},
			args: args{
				context:   new(MockContext),
				arguments: nil,
				message:   context.Message{Name: "message_3"},
			},
			wantLog: []int{15, 16, 17},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log commandLog
			states := testData.fields.makeStates(testData.args.context, &log)
			actor := Actor{states, testData.fields.currentState}

			err := actor.ProcessMessage(
				testData.args.context,
				testData.args.arguments,
				testData.args.message,
			)

			mock.AssertExpectationsForObjects(test, testData.args.context)
			checkStates(test, states.StateGroup)
			assert.Equal(test, testData.wantLog, log.commands)
			testData.wantErr(test, err)
		})
	}
}

func TestNewActorFactory(test *testing.T) {
	type args struct {
		name         string
		states       ParameterizedStateGroup
		initialState context.State
	}

	for _, testData := range []struct {
		name             string
		args             args
		wantActorFactory ActorFactory
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				name:         "Test",
				states:       ParameterizedStateGroup{StateGroup: StateGroup{"state_0": {}, "state_1": {}}},
				initialState: context.State{Name: "state_0"},
			},
			wantActorFactory: ActorFactory{
				name:         "Test",
				states:       ParameterizedStateGroup{StateGroup: StateGroup{"state_0": {}, "state_1": {}}},
				initialState: context.State{Name: "state_0"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				name:         "Test",
				states:       ParameterizedStateGroup{StateGroup: StateGroup{"state_0": {}, "state_1": {}}},
				initialState: context.State{Name: "state_unknown"},
			},
			wantActorFactory: ActorFactory{},
			wantErr:          assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			gotActorFactory, err :=
				NewActorFactory(testData.args.name, testData.args.states, testData.args.initialState)

			assert.Equal(test, testData.wantActorFactory, gotActorFactory)
			testData.wantErr(test, err)
		})
	}
}

func TestActorFactory_Name(test *testing.T) {
	factory := ActorFactory{
		name:         "Test",
		states:       ParameterizedStateGroup{StateGroup: StateGroup{"state_0": {}, "state_1": {}}},
		initialState: context.State{Name: "state_0"},
	}
	got := factory.Name()

	assert.Equal(test, "Test", got)
}

func TestActorFactory_String(test *testing.T) {
	factory := ActorFactory{
		name:         "Test",
		states:       ParameterizedStateGroup{StateGroup: StateGroup{"state_0": {}, "state_1": {}}},
		initialState: context.State{Name: "state_0"},
	}
	got := factory.String()

	assert.Equal(test, "<class Test>", got)
}

func TestActorFactory_MarshalText(test *testing.T) {
	factory := ActorFactory{
		name:         "Test",
		states:       ParameterizedStateGroup{StateGroup: StateGroup{"state_0": {}, "state_1": {}}},
		initialState: context.State{Name: "state_0"},
	}
	// it's an example of an implicit call of the runtime.ActorFactory.MarshalText() method;
	// you also can use json.Encoder with its method SetEscapeHTML() to avoid HTML escaping
	gotBytes, gotErr := json.Marshal(factory)

	assert.Equal(test, []byte(`"\u003cclass Test\u003e"`), gotBytes)
	assert.NoError(test, gotErr)
}

func TestActorFactory_CreateActor(test *testing.T) {
	factory := ActorFactory{
		name:         "Test",
		states:       ParameterizedStateGroup{StateGroup: StateGroup{"state_0": {}, "state_1": {}}},
		initialState: context.State{Name: "state_0"},
	}
	got := factory.CreateActor()

	want := &Actor{
		states:       ParameterizedStateGroup{StateGroup: StateGroup{"state_0": {}, "state_1": {}}},
		currentState: context.State{Name: "state_0"},
	}
	assert.Equal(test, want, got)
}
