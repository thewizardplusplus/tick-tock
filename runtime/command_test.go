package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

func TestCommandGroup(test *testing.T) {
	for _, testData := range []struct {
		name         string
		makeCommands func(context context.Context, log *commandLog) CommandGroup
		wantLog      []int
		wantResult   interface{}
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name:         "success without commands",
			makeCommands: func(context context.Context, log *commandLog) CommandGroup { return nil },
			wantLog:      nil,
			wantResult:   types.Nil{},
			wantErr:      assert.NoError,
		},
		{
			name: "success with commands",
			makeCommands: func(context context.Context, log *commandLog) CommandGroup {
				return newLoggableCommands(context, log, group(5), withCalls())
			},
			wantLog:    []int{0, 1, 2, 3, 4},
			wantResult: 4,
			wantErr:    assert.NoError,
		},
		{
			name: "error",
			makeCommands: func(context context.Context, log *commandLog) CommandGroup {
				return newLoggableCommands(context, log, group(5), withErrOn(2))
			},
			wantLog:    []int{0, 1, 2},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			context := new(mocks.Context)
			var log commandLog
			commands := testData.makeCommands(context, &log)
			gotResult, gotErr := commands.Run(context)

			mock.AssertExpectationsForObjects(test, context)
			checkCommands(test, commands)
			assert.Equal(test, testData.wantLog, log.commands)
			assert.Equal(test, testData.wantResult, gotResult)
			testData.wantErr(test, gotErr)
		})
	}
}

func TestParameterizedCommandGroup(test *testing.T) {
	type fields struct {
		parameters   []string
		makeCommands func(context context.Context, log *commandLog) CommandGroup
	}
	type args struct {
		context   context.Context
		arguments []interface{}
	}

	for _, testData := range []struct {
		name       string
		fields     fields
		args       args
		wantLog    []int
		wantResult interface{}
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				parameters: []string{"one", "two"},
				makeCommands: func(context context.Context, log *commandLog) CommandGroup {
					return newLoggableCommands(context, log, group(5), withCalls())
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
			},
			wantLog:    []int{0, 1, 2, 3, 4},
			wantResult: 4,
			wantErr:    assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				parameters: []string{"one", "two"},
				makeCommands: func(context context.Context, log *commandLog) CommandGroup {
					return newLoggableCommands(context, log, group(5), withErrOn(2))
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
			},
			wantLog:    []int{0, 1, 2},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log commandLog
			commands := testData.fields.makeCommands(testData.args.context, &log)
			parameterizedCommands := NewParameterizedCommandGroup(testData.fields.parameters, commands)
			gotResult, gotErr :=
				parameterizedCommands.ParameterizedRun(testData.args.context, testData.args.arguments)

			mock.AssertExpectationsForObjects(test, testData.args.context)
			checkCommands(test, commands)
			assert.Equal(test, testData.wantLog, log.commands)
			assert.Equal(test, testData.wantResult, gotResult)
			testData.wantErr(test, gotErr)
		})
	}
}
