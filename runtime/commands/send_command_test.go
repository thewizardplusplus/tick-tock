package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

func TestSendCommand(test *testing.T) {
	context := new(mocks.Context)
	context.On("SendMessage", "test").Return()

	gotResult, gotErr := NewSendCommand("test").Run(context)

	mock.AssertExpectationsForObjects(test, context)
	assert.Equal(test, types.Nil{}, gotResult)
	assert.NoError(test, gotErr)
}
