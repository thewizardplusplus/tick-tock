package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolean(test *testing.T) {
	for _, data := range []struct {
		name  string
		value interface{}
	}{
		{
			name:  "true",
			value: True,
		},
		{
			name:  "false",
			value: False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			assert.IsType(test, Boolean(0), data.value)
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
		wantResult Boolean
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "success/nil",
			args:       args{Nil{}},
			wantResult: False,
			wantErr:    assert.NoError,
		},
		{
			name:       "success/float64/greater than zero",
			args:       args{23.0},
			wantResult: True,
			wantErr:    assert.NoError,
		},
		{
			name:       "success/float64/less than zero",
			args:       args{-23.0},
			wantResult: True,
			wantErr:    assert.NoError,
		},
		{
			name:       "success/float64/equal to zero",
			args:       args{0.0},
			wantResult: False,
			wantErr:    assert.NoError,
		},
		{
			name:       "success/*Pair/nonempty",
			args:       args{&Pair{"one", &Pair{"two", nil}}},
			wantResult: True,
			wantErr:    assert.NoError,
		},
		{
			name:       "success/*Pair/empty",
			args:       args{(*Pair)(nil)},
			wantResult: False,
			wantErr:    assert.NoError,
		},
		{
			name:       "error",
			args:       args{func() {}},
			wantResult: False,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotResult, gotErr := NewBoolean(data.args.value)

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
		want Boolean
	}{
		{
			name: "true",
			args: args{true},
			want: True,
		},
		{
			name: "false",
			args: args{false},
			want: False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewBooleanFromGoBool(data.args.value)

			assert.Equal(test, data.want, got)
		})
	}
}
