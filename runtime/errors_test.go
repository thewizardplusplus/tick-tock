package runtime

import (
	"testing"
	"testing/iotest"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/tests/mocks"
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
		name        string
		args        args
		wantMessage string
		wantCode    int
	}{
		{
			name:        "success with a common error",
			args:        args{iotest.ErrTimeout},
			wantMessage: "error: timeout",
			wantCode:    1,
		},
		{
			name: "success with an user exit error (direct)",
			args: args{ErrUserExit},
		},
		{
			name: "success with an user exit error (wrapped)",
			args: args{errors.Wrap(errors.Wrap(ErrUserExit, "level #1"), "level #2")},
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			writer := new(mocks.Writer)
			if len(testData.wantMessage) != 0 {
				writer.On("Write", []byte(testData.wantMessage)).Return(len(testData.wantMessage), nil)
			}

			exiter := new(mocks.Exiter)
			exiter.On("Exit", testData.wantCode).Return()

			NewDefaultErrorHandler(writer, exiter.Exit).HandleError(testData.args.err)

			writer.AssertExpectations(test)
			exiter.AssertExpectations(test)
		})
	}
}
