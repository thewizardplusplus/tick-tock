package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendCommand(test *testing.T) {
	context := new(MockContext)
	context.On("SendMessage", "test").Return()

	err := NewSendCommand("test").Run(context)

	context.AssertExpectations(test)
	assert.NoError(test, err)
}
