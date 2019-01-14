package runtime

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestNewUnknownStateError(test *testing.T) {
	got := newUnknownStateError("test")
	assert.Equal(test, "unknown state test", got.Error())
}

func TestDefaultErrorHandler(test *testing.T) {
	type args struct {
		err error
	}

	for _, testData := range []struct {
		name      string
		args      args
		wantPanic func(assert.TestingT, assert.PanicTestFunc, ...interface{}) bool
	}{
		{
			name:      "nil error",
			wantPanic: assert.NotPanics,
		},
		{
			name:      "not nil error",
			args:      args{iotest.ErrTimeout},
			wantPanic: assert.Panics,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			testData.wantPanic(test, func() { DefaultErrorHandler{}.HandleError(testData.args.err) })
		})
	}
}
