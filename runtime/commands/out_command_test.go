package commands

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	testsmocks "github.com/thewizardplusplus/tick-tock/tests/mocks"
)

func TestOutCommand(test *testing.T) {
	type fields struct {
		message string
	}

	for _, testData := range []struct {
		name       string
		fields     fields
		writingErr error
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:    "success with a nonempty message",
			fields:  fields{"test"},
			wantErr: assert.NoError,
		},
		{
			name:    "success with an empty message",
			wantErr: assert.NoError,
		},
		{
			name:       "error",
			fields:     fields{"test"},
			writingErr: iotest.ErrTimeout,
			wantErr:    assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			writer := new(testsmocks.Writer)
			writeCall := writer.On("Write", []byte(testData.fields.message))
			if testData.writingErr == nil {
				writeCall.Return(len(testData.fields.message), nil)
			} else {
				writeCall.Return(0, testData.writingErr)
			}

			context := new(contextmocks.Context)
			err := NewOutCommand(writer, testData.fields.message).Run(context)

			writer.AssertExpectations(test)
			context.AssertExpectations(test)
			testData.wantErr(test, err)
		})
	}
}
