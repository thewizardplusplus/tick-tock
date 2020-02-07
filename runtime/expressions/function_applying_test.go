package expressions

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/internal/tests"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestNewFunctionApplying(test *testing.T) {
	arguments := []Expression{NewSignedExpression("one"), NewSignedExpression("two")}
	handler := func([]interface{}) (interface{}, error) { panic("not implemented") }
	got := NewFunctionApplying(arguments, handler)

	for _, argument := range arguments {
		mock.AssertExpectationsForObjects(test, argument)
	}
	assert.Equal(test, arguments, got.arguments)
	assert.Equal(test, tests.GetReflectionAddress(handler), tests.GetReflectionAddress(got.handler))
}

func TestFunctionApplying_Evaluate(test *testing.T) {
	type fields struct {
		arguments []Expression
		handler   FunctionHandler
	}
	type args struct {
		context context.Context
	}

	for _, data := range []struct {
		name       string
		fields     fields
		args       args
		wantResult interface{}
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success with arguments",
			fields: fields{
				arguments: []Expression{
					func() Expression {
						expression := NewSignedExpression("one")
						expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.3, nil)

						return expression
					}(),
					func() Expression {
						expression := NewSignedExpression("two")
						expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(4.2, nil)

						return expression
					}(),
				},
				handler: func(arguments []interface{}) (interface{}, error) {
					var result float64
					for _, argument := range arguments {
						result += argument.(float64)
					}

					return result, nil
				},
			},
			args:       args{new(contextmocks.Context)},
			wantResult: 6.5,
			wantErr:    assert.NoError,
		},
		{
			name: "success without arguments",
			fields: fields{
				arguments: nil,
				handler:   func([]interface{}) (interface{}, error) { return 2.3, nil },
			},
			args:       args{new(contextmocks.Context)},
			wantResult: 2.3,
			wantErr:    assert.NoError,
		},
		{
			name: "error on argument evaluation",
			fields: fields{
				arguments: []Expression{
					func() Expression {
						expression := NewSignedExpression("one")
						expression.
							On("Evaluate", mock.AnythingOfType("*mocks.Context")).
							Return(nil, iotest.ErrTimeout)

						return expression
					}(),
					NewSignedExpression("two"),
				},
				handler: func([]interface{}) (interface{}, error) { panic("not implemented") },
			},
			args:       args{new(contextmocks.Context)},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error on handler calling",
			fields: fields{
				arguments: []Expression{
					func() Expression {
						expression := NewSignedExpression("one")
						expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.3, nil)

						return expression
					}(),
					func() Expression {
						expression := NewSignedExpression("two")
						expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(4.2, nil)

						return expression
					}(),
				},
				handler: func([]interface{}) (interface{}, error) { return nil, iotest.ErrTimeout },
			},
			args:       args{new(contextmocks.Context)},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			expression := FunctionApplying{
				arguments: data.fields.arguments,
				handler:   data.fields.handler,
			}
			gotResult, gotErr := expression.Evaluate(data.args.context)

			for _, argument := range data.fields.arguments {
				mock.AssertExpectationsForObjects(test, argument)
			}
			mock.AssertExpectationsForObjects(test, data.args.context)
			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}
