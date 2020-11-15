package runtime

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/internal/test-utils/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
)

func TestDefaultErrorHandler(test *testing.T) {
	const message = "error: timeout\n"
	writer := new(MockWriter)
	writer.On("Write", []byte(message)).Return(len(message), nil)

	exiter := new(mocks.Exiter)
	exiter.On("Exit", 1).Return()

	NewDefaultErrorHandler(writer, exiter.Exit).HandleError(iotest.ErrTimeout)

	mock.AssertExpectationsForObjects(test, writer, exiter)
}

func TestNewUnknownStateError(test *testing.T) {
	got := newUnknownStateError(context.State{Name: "test"})
	assert.Equal(test, "unknown state test", got.Error())
}
