package context_test

import (
	"testing"
	"unsafe"

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

func TestDefaultContext_SetStateHolder(test *testing.T) {
	holder := new(mocks.StateHolder)
	defaultContext := context.DefaultContext{}
	defaultContext.SetStateHolder(holder)

	mock.AssertExpectationsForObjects(test, holder)
	assert.Equal(test, context.DefaultContext{StateHolder: holder}, defaultContext)
}

func TestDefaultContext_Copy(test *testing.T) {
	sender := new(mocks.MessageSender)
	holder := new(mocks.StateHolder)
	defaultContext := &context.DefaultContext{MessageSender: sender, StateHolder: holder}
	defaultContextCopy := defaultContext.Copy()

	mock.AssertExpectationsForObjects(test, sender, holder)
	assert.Equal(test, defaultContext, defaultContextCopy)
	if assert.IsType(test, &context.DefaultContext{}, defaultContextCopy) {
		assert.NotEqual(
			test,
			unsafe.Pointer(defaultContext), unsafe.Pointer(defaultContextCopy.(*context.DefaultContext)),
		)
	}
}
