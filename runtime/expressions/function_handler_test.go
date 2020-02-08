package expressions

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestNewArithmeticFunctionHandler(test *testing.T) {
	type fields struct {
		handler ArithmeticFunctionHandler
	}
	type args struct {
		context   context.Context
		arguments []interface{}
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
				handler: func(context context.Context, arguments []float64) (float64, error) {
					var result float64
					for _, argument := range arguments {
						result += argument
					}

					return result, nil
				},
			},
			args: args{
				context:   new(contextmocks.Context),
				arguments: []interface{}{2.3, 4.2},
			},
			wantResult: 6.5,
			wantErr:    assert.NoError,
		},
		{
			name: "success without arguments",
			fields: fields{
				handler: func(context.Context, []float64) (float64, error) { return 2.3, nil },
			},
			args: args{
				context:   new(contextmocks.Context),
				arguments: nil,
			},
			wantResult: 2.3,
			wantErr:    assert.NoError,
		},
		{
			name: "error on argument typecasting",
			fields: fields{
				handler: func(context.Context, []float64) (float64, error) { panic("not implemented") },
			},
			args: args{
				context:   new(contextmocks.Context),
				arguments: []interface{}{2, 4.2},
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error on handler calling",
			fields: fields{
				handler: func(context.Context, []float64) (float64, error) { return 0, iotest.ErrTimeout },
			},
			args: args{
				context:   new(contextmocks.Context),
				arguments: []interface{}{2.3, 4.2},
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := NewArithmeticFunctionHandler(data.fields.handler)
			gotResult, gotErr := handler(data.args.context, data.args.arguments)

			mock.AssertExpectationsForObjects(test, data.args.context)
			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}
