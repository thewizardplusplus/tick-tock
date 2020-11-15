package expressions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

func TestNewIdentifier(test *testing.T) {
	got := NewIdentifier("test")

	assert.Equal(test, "test", got.name)
}

func TestIdentifier_Evaluate(test *testing.T) {
	type fields struct {
		name string
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
			name:   "success",
			fields: fields{"test"},
			args: args{
				context: func() context.Context {
					context := new(MockContext)
					context.On("Value", "test").Return(2.3, true)

					return context
				}(),
			},
			wantResult: 2.3,
			wantErr:    assert.NoError,
		},
		{
			name:   "error",
			fields: fields{"test"},
			args: args{
				context: func() context.Context {
					context := new(MockContext)
					context.On("Value", "test").Return(nil, false)

					return context
				}(),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			identifier := Identifier{
				name: data.fields.name,
			}
			gotResult, gotErr := identifier.Evaluate(data.args.context)

			mock.AssertExpectationsForObjects(test, data.args.context)
			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}
