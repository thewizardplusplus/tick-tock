package context

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
		want ValueGroup
	}{
		{
			name: "lengths of both lists are equal",
			args: args{
				parameters: []string{"one", "two"},
				arguments:  []interface{}{23, 42},
			},
			want: ValueGroup{"one": 23, "two": 42},
		},
		{
			name: "lengths of both lists are equal/both lists are empty",
			args: args{
				parameters: nil,
				arguments:  nil,
			},
			want: ValueGroup{},
		},
		{
			name: "argument list is longer",
			args: args{
				parameters: []string{"one", "two"},
				arguments:  []interface{}{12, 23, 42},
			},
			want: ValueGroup{"one": 12, "two": 23},
		},
		{
			name: "argument list is longer/without parameters",
			args: args{
				parameters: nil,
				arguments:  []interface{}{23, 42},
			},
			want: ValueGroup{},
		},
		{
			name: "parameter list is longer",
			args: args{
				parameters: []string{"one", "two", "three"},
				arguments:  []interface{}{23, 42},
			},
			want: ValueGroup{"one": 23, "two": 42, "three": types.Nil{}},
		},
		{
			name: "parameter list is longer/without arguments",
			args: args{
				parameters: []string{"one", "two"},
				arguments:  nil,
			},
			want: ValueGroup{"one": types.Nil{}, "two": types.Nil{}},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := ZipValues(data.args.parameters, data.args.arguments)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestSetValues(test *testing.T) {
	type args struct {
		holder ValueHolder
		values ValueGroup
	}

	for _, data := range []struct {
		name string
		args args
	}{
		{
			name: "with values",
			args: args{
				holder: func() ValueHolder {
					holder := new(MockValueHolder)
					holder.On("SetValue", "one", 1)
					holder.On("SetValue", "two", 2)

					return holder
				}(),
				values: ValueGroup{"one": 1, "two": 2},
			},
		},
		{
			name: "without values",
			args: args{
				holder: new(MockValueHolder),
				values: nil,
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			SetValues(data.args.holder, data.args.values)

			mock.AssertExpectationsForObjects(test, data.args.holder)
		})
	}
}
