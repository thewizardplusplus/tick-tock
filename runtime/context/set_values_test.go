package context_test

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

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
