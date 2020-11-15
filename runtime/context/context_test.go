package context

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewDefaultContext(test *testing.T) {
	got := NewDefaultContext()

	assert.Nil(test, got.MessageSender)
	assert.Nil(test, got.StateHolder)
	assert.Equal(test, DefaultValueStore{}, got.CopyableValueStore)
}

func TestDefaultContext_SetMessageSender(test *testing.T) {
	sender := new(MockMessageSender)
	defaultContext := DefaultContext{}
	defaultContext.SetMessageSender(sender)

	mock.AssertExpectationsForObjects(test, sender)
	assert.Equal(test, DefaultContext{MessageSender: sender}, defaultContext)
}

func TestDefaultContext_SetStateHolder(test *testing.T) {
	holder := new(MockStateHolder)
	defaultContext := DefaultContext{}
	defaultContext.SetStateHolder(holder)

	mock.AssertExpectationsForObjects(test, holder)
	assert.Equal(test, DefaultContext{StateHolder: holder}, defaultContext)
}

func TestDefaultContext_SetActorRegister(test *testing.T) {
	register := new(MockActorRegister)
	defaultContext := DefaultContext{}
	defaultContext.SetActorRegister(register)

	mock.AssertExpectationsForObjects(test, register)
	assert.Equal(test, DefaultContext{ActorRegister: register}, defaultContext)
}

func TestDefaultContext_SetValueStore(test *testing.T) {
	store := new(MockCopyableValueStore)
	defaultContext := DefaultContext{}
	defaultContext.SetValueStore(store)

	mock.AssertExpectationsForObjects(test, store)
	assert.Equal(test, DefaultContext{CopyableValueStore: store}, defaultContext)
}

func TestDefaultContext_Copy(test *testing.T) {
	store := new(MockCopyableValueStore)
	store.On("Copy").Return(store)

	sender := new(MockMessageSender)
	holder := new(MockStateHolder)
	register := new(MockActorRegister)
	defaultContext := &DefaultContext{
		MessageSender:      sender,
		StateHolder:        holder,
		ActorRegister:      register,
		CopyableValueStore: store,
	}
	defaultContextCopy := defaultContext.Copy()

	mock.AssertExpectationsForObjects(test, sender, holder, store)
	assert.Equal(test, defaultContext, defaultContextCopy)
	if assert.IsType(test, &DefaultContext{}, defaultContextCopy) {
		assert.NotEqual(
			test,
			unsafe.Pointer(defaultContext),
			unsafe.Pointer(defaultContextCopy.(*DefaultContext)),
		)
	}
}
