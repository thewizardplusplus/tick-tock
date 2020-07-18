package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
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
			name: "success/not equal/same types",
			args: args{
				leftValue:  23.0,
				rightValue: 42.0,
			},
			wantResult: assert.False,
			wantErr:    assert.NoError,
		},
		{
			name: "success/not equal/different types/Nil",
			args: args{
				leftValue:  types.Nil{},
				rightValue: 23.0,
			},
			wantResult: assert.False,
			wantErr:    assert.NoError,
		},
		{
			name: "success/not equal/different types/float64",
			args: args{
				leftValue:  23.0,
				rightValue: &types.Pair{Head: 12.0, Tail: &types.Pair{Head: 23.0, Tail: nil}},
			},
			wantResult: assert.False,
			wantErr:    assert.NoError,
		},
		{
			name: "success/not equal/different types/*Pair",
			args: args{
				leftValue:  &types.Pair{Head: 12.0, Tail: &types.Pair{Head: 23.0, Tail: nil}},
				rightValue: types.Nil{},
			},
			wantResult: assert.False,
			wantErr:    assert.NoError,
		},
		{
			name: "error/unsupported type",
			args: args{
				leftValue:  func() {},
				rightValue: types.Nil{},
			},
			wantResult: assert.False,
			wantErr:    assert.Error,
		},
		{
			name: "error/unable to compare",
			args: args{
				leftValue:  &types.Pair{Head: 12.0, Tail: &types.Pair{Head: 23.0, Tail: nil}},
				rightValue: &types.Pair{Head: 12.0, Tail: &types.Pair{Head: types.Nil{}, Tail: nil}},
			},
			wantResult: assert.False,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotResult, gotErr := types.Equals(data.args.leftValue, data.args.rightValue)

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
		wantResult types.ComparisonResult
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "Nil/success",
			args: args{
				leftValue:  types.Nil{},
				rightValue: types.Nil{},
			},
			wantResult: types.Equal,
			wantErr:    assert.NoError,
		},
		{
			name: "Nil/error",
			args: args{
				leftValue:  types.Nil{},
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
			wantResult: types.Less,
			wantErr:    assert.NoError,
		},
		{
			name: "float64/success/equal",
			args: args{
				leftValue:  23.0,
				rightValue: 23.0,
			},
			wantResult: types.Equal,
			wantErr:    assert.NoError,
		},
		{
			name: "float64/success/greater",
			args: args{
				leftValue:  42.0,
				rightValue: 23.0,
			},
			wantResult: types.Greater,
			wantErr:    assert.NoError,
		},
		{
			name: "float64/error",
			args: args{
				leftValue:  23.0,
				rightValue: &types.Pair{Head: 12.0, Tail: &types.Pair{Head: 23.0, Tail: nil}},
			},
			wantResult: 0,
			wantErr:    assert.Error,
		},
		{
			name: "*Pair/success/less",
			args: args{
				leftValue:  &types.Pair{Head: 12.0, Tail: &types.Pair{Head: 23.0, Tail: nil}},
				rightValue: &types.Pair{Head: 12.0, Tail: &types.Pair{Head: 42.0, Tail: nil}},
			},
			wantResult: types.Less,
			wantErr:    assert.NoError,
		},
		{
			name: "*Pair/success/equal",
			args: args{
				leftValue:  &types.Pair{Head: 12.0, Tail: &types.Pair{Head: 23.0, Tail: nil}},
				rightValue: &types.Pair{Head: 12.0, Tail: &types.Pair{Head: 23.0, Tail: nil}},
			},
			wantResult: types.Equal,
			wantErr:    assert.NoError,
		},
		{
			name: "*Pair/success/greater",
			args: args{
				leftValue:  &types.Pair{Head: 12.0, Tail: &types.Pair{Head: 42.0, Tail: nil}},
				rightValue: &types.Pair{Head: 12.0, Tail: &types.Pair{Head: 23.0, Tail: nil}},
			},
			wantResult: types.Greater,
			wantErr:    assert.NoError,
		},
		{
			name: "*Pair/error/incorrect type",
			args: args{
				leftValue:  &types.Pair{Head: 12.0, Tail: &types.Pair{Head: 23.0, Tail: nil}},
				rightValue: types.Nil{},
			},
			wantResult: 0,
			wantErr:    assert.Error,
		},
		{
			name: "*Pair/error/unable to compare",
			args: args{
				leftValue:  &types.Pair{Head: 12.0, Tail: &types.Pair{Head: 23.0, Tail: nil}},
				rightValue: &types.Pair{Head: 12.0, Tail: &types.Pair{Head: types.Nil{}, Tail: nil}},
			},
			wantResult: 0,
			wantErr:    assert.Error,
		},
		{
			name: "unsupported type",
			args: args{
				leftValue:  func() {},
				rightValue: types.Nil{},
			},
			wantResult: 0,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotResult, gotErr := types.Compare(data.args.leftValue, data.args.rightValue)

			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}
