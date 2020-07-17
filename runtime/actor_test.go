package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestActor_SetState(test *testing.T) {
	type fields struct {
		states       StateGroup
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
				states:       StateGroup{"state_0": {}, "state_1": {}},
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
				states:       StateGroup{"state_0": {}, "state_1": {}},
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
				states:       StateGroup{"state_0": {}, "state_1": {}},
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
		makeStates   func(context context.Context, log *commandLog) StateGroup
		currentState context.State
	}
	type args struct {
		contextCopy context.Context
		message     context.Message
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
			var log commandLog
			states := testData.fields.makeStates(testData.args.contextCopy, &log)
			actor := Actor{states, testData.fields.currentState}
			testData.args.contextCopy.(*mocks.Context).On("SetStateHolder", &actor).Return()

			contextOriginal := new(mocks.Context)
			contextOriginal.On("Copy").Return(testData.args.contextCopy)

			err := actor.ProcessMessage(contextOriginal, testData.args.message)

			mock.AssertExpectationsForObjects(test, contextOriginal, testData.args.contextCopy)
			checkStates(test, states)
			assert.Equal(test, testData.wantLog, log.commands)
			testData.wantErr(test, err)
		})
	}
}

func TestNewActorFactory(test *testing.T) {
	type args struct {
		name         string
		states       StateGroup
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
				states:       StateGroup{"state_0": {}, "state_1": {}},
				initialState: context.State{Name: "state_0"},
			},
			wantActorFactory: ActorFactory{
				name:         "Test",
				states:       StateGroup{"state_0": {}, "state_1": {}},
				initialState: context.State{Name: "state_0"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				name:         "Test",
				states:       StateGroup{"state_0": {}, "state_1": {}},
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
		states:       StateGroup{"state_0": {}, "state_1": {}},
		initialState: context.State{Name: "state_0"},
	}
	got := factory.Name()

	assert.Equal(test, "Test", got)
}

func TestActorFactory_String(test *testing.T) {
	factory := ActorFactory{
		name:         "Test",
		states:       StateGroup{"state_0": {}, "state_1": {}},
		initialState: context.State{Name: "state_0"},
	}
	got := factory.String()

	assert.Equal(test, "<class Test>", got)
}

func TestActorFactory_CreateActor(test *testing.T) {
	factory := ActorFactory{
		name:         "Test",
		states:       StateGroup{"state_0": {}, "state_1": {}},
		initialState: context.State{Name: "state_0"},
	}
	got := factory.CreateActor()

	want := &Actor{
		states:       StateGroup{"state_0": {}, "state_1": {}},
		currentState: context.State{Name: "state_0"},
	}
	assert.Equal(test, want, got)
}
