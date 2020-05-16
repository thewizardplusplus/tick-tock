package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEquals(test *testing.T) {
	type args struct {
		leftValue  interface{}
		rightValue interface{}
	}

	for _, data := range []struct {
		name       string
		args       args
		wantResult assert.BoolAssertionFunc
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success/equal",
			args: args{
				leftValue:  23.0,
				rightValue: 23.0,
			},
			wantResult: assert.True,
			wantErr:    assert.NoError,
		},
		{
			name: "success/not equal",
			args: args{
				leftValue:  23.0,
				rightValue: 42.0,
			},
			wantResult: assert.False,
			wantErr:    assert.NoError,
		},
		{
			name: "error",
			args: args{
				leftValue:  23.0,
				rightValue: Nil{},
			},
			wantResult: assert.False,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotResult, gotErr := Equals(data.args.leftValue, data.args.rightValue)

			data.wantResult(test, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}

func TestCompare(test *testing.T) {
	type args struct {
		leftValue  interface{}
		rightValue interface{}
	}

	for _, data := range []struct {
		name       string
		args       args
		wantResult ComparisonResult
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "Nil/success",
			args: args{
				leftValue:  Nil{},
				rightValue: Nil{},
			},
			wantResult: Equal,
			wantErr:    assert.NoError,
		},
		{
			name: "Nil/error",
			args: args{
				leftValue:  Nil{},
				rightValue: 23.0,
			},
			wantResult: 0,
			wantErr:    assert.Error,
		},
		{
			name: "float64/success/less",
			args: args{
				leftValue:  23.0,
				rightValue: 42.0,
			},
			wantResult: Less,
			wantErr:    assert.NoError,
		},
		{
			name: "float64/success/equal",
			args: args{
				leftValue:  23.0,
				rightValue: 23.0,
			},
			wantResult: Equal,
			wantErr:    assert.NoError,
		},
		{
			name: "float64/success/greater",
			args: args{
				leftValue:  42.0,
				rightValue: 23.0,
			},
			wantResult: Greater,
			wantErr:    assert.NoError,
		},
		{
			name: "float64/error",
			args: args{
				leftValue:  23.0,
				rightValue: &Pair{12.0, &Pair{23.0, nil}},
			},
			wantResult: 0,
			wantErr:    assert.Error,
		},
		{
			name: "*Pair/success/less",
			args: args{
				leftValue:  &Pair{12.0, &Pair{23.0, nil}},
				rightValue: &Pair{12.0, &Pair{42.0, nil}},
			},
			wantResult: Less,
			wantErr:    assert.NoError,
		},
		{
			name: "*Pair/success/equal",
			args: args{
				leftValue:  &Pair{12.0, &Pair{23.0, nil}},
				rightValue: &Pair{12.0, &Pair{23.0, nil}},
			},
			wantResult: Equal,
			wantErr:    assert.NoError,
		},
		{
			name: "*Pair/success/greater",
			args: args{
				leftValue:  &Pair{12.0, &Pair{42.0, nil}},
				rightValue: &Pair{12.0, &Pair{23.0, nil}},
			},
			wantResult: Greater,
			wantErr:    assert.NoError,
		},
		{
			name: "*Pair/error/incorrect type",
			args: args{
				leftValue:  &Pair{12.0, &Pair{23.0, nil}},
				rightValue: Nil{},
			},
			wantResult: 0,
			wantErr:    assert.Error,
		},
		{
			name: "*Pair/error/unable to compare",
			args: args{
				leftValue:  &Pair{12.0, &Pair{23.0, nil}},
				rightValue: &Pair{12.0, &Pair{Nil{}, nil}},
			},
			wantResult: 0,
			wantErr:    assert.Error,
		},
		{
			name: "unsupported type",
			args: args{
				leftValue:  func() {},
				rightValue: Nil{},
			},
			wantResult: 0,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotResult, gotErr := Compare(data.args.leftValue, data.args.rightValue)

			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}
