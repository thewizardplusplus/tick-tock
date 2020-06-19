package context_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

func TestZipValues(test *testing.T) {
	type args struct {
		parameters []string
		arguments  []interface{}
	}

	for _, data := range []struct {
		name string
		args args
		want context.ValueGroup
	}{
		{
			name: "lengths of both lists are equal",
			args: args{
				parameters: []string{"one", "two"},
				arguments:  []interface{}{23, 42},
			},
			want: context.ValueGroup{"one": 23, "two": 42},
		},
		{
			name: "lengths of both lists are equal/both lists are empty",
			args: args{
				parameters: nil,
				arguments:  nil,
			},
			want: context.ValueGroup{},
		},
		{
			name: "argument list is longer",
			args: args{
				parameters: []string{"one", "two"},
				arguments:  []interface{}{12, 23, 42},
			},
			want: context.ValueGroup{"one": 12, "two": 23},
		},
		{
			name: "argument list is longer/without parameters",
			args: args{
				parameters: nil,
				arguments:  []interface{}{23, 42},
			},
			want: context.ValueGroup{},
		},
		{
			name: "parameter list is longer",
			args: args{
				parameters: []string{"one", "two", "three"},
				arguments:  []interface{}{23, 42},
			},
			want: context.ValueGroup{"one": 23, "two": 42, "three": types.Nil{}},
		},
		{
			name: "parameter list is longer/without arguments",
			args: args{
				parameters: []string{"one", "two"},
				arguments:  nil,
			},
			want: context.ValueGroup{"one": types.Nil{}, "two": types.Nil{}},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := context.ZipValues(data.args.parameters, data.args.arguments)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestSetValues(test *testing.T) {
	type args struct {
		holder context.ValueHolder
		values context.ValueGroup
	}

	for _, data := range []struct {
		name string
		args args
	}{
		{
			name: "with values",
			args: args{
				holder: func() context.ValueHolder {
					holder := new(mocks.ValueHolder)
					holder.On("SetValue", "one", 1)
					holder.On("SetValue", "two", 2)

					return holder
				}(),
				values: context.ValueGroup{"one": 1, "two": 2},
			},
		},
		{
			name: "without values",
			args: args{
				holder: new(mocks.ValueHolder),
				values: nil,
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			context.SetValues(data.args.holder, data.args.values)

			mock.AssertExpectationsForObjects(test, data.args.holder)
		})
	}
}
