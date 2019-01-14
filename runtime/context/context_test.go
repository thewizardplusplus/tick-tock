package context_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestDefaultContext_SetMessageSender(test *testing.T) {
	sender := new(mocks.MessageSender)
	defaultContext := context.DefaultContext{}
	defaultContext.SetMessageSender(sender)

	mock.AssertExpectationsForObjects(test, sender)
	assert.Equal(test, context.DefaultContext{MessageSender: sender}, defaultContext)
}
