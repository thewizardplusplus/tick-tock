package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// it checks that the types.Boolean type is an alias of the float64 type
func TestBoolean(test *testing.T) {
	for _, data := range []struct {
		name  string
		value interface{}
	}{
		{
			name:  "true",
			value: types.True,
		},
		{
			name:  "false",
			value: types.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			assert.IsType(test, types.Boolean(0), data.value)
			assert.IsType(test, float64(0), data.value)
		})
	}
}

func TestNewBoolean(test *testing.T) {
	type args struct {
		value interface{}
	}

	for _, data := range []struct {
		name       string
		args       args
		wantResult types.Boolean
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "success/nil",
			args:       args{types.Nil{}},
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "success/float64/greater than zero",
			args:       args{23.0},
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "success/float64/less than zero",
			args:       args{-23.0},
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "success/float64/equal to zero",
			args:       args{0.0},
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "success/*Pair/nonempty",
			args:       args{&types.Pair{Head: "one", Tail: &types.Pair{Head: "two", Tail: nil}}},
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "success/*Pair/empty",
			args:       args{(*types.Pair)(nil)},
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "success/HashTable/nonempty",
			args:       args{types.HashTable{"one": "two", "three": "four"}},
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "success/HashTable/empty",
			args:       args{(types.HashTable)(nil)},
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name: "success/actor class",
			args: args{
				value: func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test",
						runtime.ParameterizedStateGroup{StateGroup: runtime.StateGroup{"state_0": {}}},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 0, runtime.Dependencies{})
				}(),
			},
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "error",
			args:       args{func() {}},
			wantResult: types.False,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotResult, gotErr := types.NewBoolean(data.args.value)

			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}

func TestNewBooleanFromGoBool(test *testing.T) {
	type args struct {
		value bool
	}

	for _, data := range []struct {
		name string
		args args
		want types.Boolean
	}{
		{
			name: "true",
			args: args{true},
			want: types.True,
		},
		{
			name: "false",
			args: args{false},
			want: types.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := types.NewBooleanFromGoBool(data.args.value)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestNegateBoolean(test *testing.T) {
	type args struct {
		value types.Boolean
	}

	for _, data := range []struct {
		name string
		args args
		want types.Boolean
	}{
		{
			name: "true",
			args: args{types.True},
			want: types.False,
		},
		{
			name: "false",
			args: args{types.False},
			want: types.True,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := types.NegateBoolean(data.args.value)

			assert.Equal(test, data.want, got)
		})
	}
}
