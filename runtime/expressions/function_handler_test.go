package expressions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewArithmeticFunctionHandler(test *testing.T) {
	type fields struct {
		handler ArithmeticFunctionHandler
	}
	type args struct {
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
				handler: func(arguments []float64) (float64, error) {
					var result float64
					for _, argument := range arguments {
						result += argument
					}

					return result, nil
				},
			},
			args:       args{[]interface{}{2.3, 4.2}},
			wantResult: 6.5,
			wantErr:    assert.NoError,
		},
		{
			name: "success without arguments",
			fields: fields{
				handler: func([]float64) (float64, error) { return 2.3, nil },
			},
			args:       args{nil},
			wantResult: 2.3,
			wantErr:    assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				handler: func([]float64) (float64, error) { panic("not implemented") },
			},
			args:       args{[]interface{}{2, 3}},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			handler := NewArithmeticFunctionHandler(data.fields.handler)
			gotResult, gotErr := handler(data.args.arguments)

			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}
