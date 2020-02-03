package context_test

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestNewDefaultContext(test *testing.T) {
	got := context.NewDefaultContext()

	assert.Nil(test, got.MessageSender)
	assert.Nil(test, got.StateHolder)
	assert.Equal(test, context.DefaultValueStore{}, got.CopyableValueStore)
}

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

func TestDefaultContext_SetValueStore(test *testing.T) {
	store := new(mocks.CopyableValueStore)
	defaultContext := context.DefaultContext{}
	defaultContext.SetValueStore(store)

	mock.AssertExpectationsForObjects(test, store)
	assert.Equal(test, context.DefaultContext{CopyableValueStore: store}, defaultContext)
}

func TestDefaultContext_Copy(test *testing.T) {
	store := new(mocks.CopyableValueStore)
	store.On("Copy").Return(store)

	sender := new(mocks.MessageSender)
	holder := new(mocks.StateHolder)
	defaultContext := &context.DefaultContext{
		MessageSender:      sender,
		StateHolder:        holder,
		CopyableValueStore: store,
	}
	defaultContextCopy := defaultContext.Copy()

	mock.AssertExpectationsForObjects(test, sender, holder, store)
	assert.Equal(test, defaultContext, defaultContextCopy)
	if assert.IsType(test, &context.DefaultContext{}, defaultContextCopy) {
		assert.NotEqual(
			test,
			unsafe.Pointer(defaultContext), unsafe.Pointer(defaultContextCopy.(*context.DefaultContext)),
		)
	}
}
