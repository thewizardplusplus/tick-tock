// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package expressions

import (
	mapset "github.com/deckarep/golang-set"
	context "github.com/thewizardplusplus/tick-tock/runtime/context"

	mock "github.com/stretchr/testify/mock"
)

// MockContext is an autogenerated mock type for the Context type
type MockContext struct {
	mock.Mock
}

// Copy provides a mock function with given fields:
func (_m *MockContext) Copy() context.Context {
	ret := _m.Called()

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// RegisterActor provides a mock function with given fields: actor, arguments
func (_m *MockContext) RegisterActor(actor context.Actor, arguments []interface{}) {
	_m.Called(actor, arguments)
}

// SendMessage provides a mock function with given fields: message
func (_m *MockContext) SendMessage(message context.Message) {
	_m.Called(message)
}

// SetActorRegister provides a mock function with given fields: register
func (_m *MockContext) SetActorRegister(register context.ActorRegister) {
	_m.Called(register)
}

// SetMessageSender provides a mock function with given fields: sender
func (_m *MockContext) SetMessageSender(sender context.MessageSender) {
	_m.Called(sender)
}

// SetState provides a mock function with given fields: state
func (_m *MockContext) SetState(state context.State) error {
	ret := _m.Called(state)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.State) error); ok {
		r0 = rf(state)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetStateHolder provides a mock function with given fields: holder
func (_m *MockContext) SetStateHolder(holder context.StateHolder) {
	_m.Called(holder)
}

// SetValue provides a mock function with given fields: name, value
func (_m *MockContext) SetValue(name string, value interface{}) {
	_m.Called(name, value)
}

// SetValueStore provides a mock function with given fields: store
func (_m *MockContext) SetValueStore(store context.CopyableValueStore) {
	_m.Called(store)
}

// Value provides a mock function with given fields: name
func (_m *MockContext) Value(name string) (interface{}, bool) {
	ret := _m.Called(name)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// ValuesNames provides a mock function with given fields:
func (_m *MockContext) ValuesNames() mapset.Set {
	ret := _m.Called()

	var r0 mapset.Set
	if rf, ok := ret.Get(0).(func() mapset.Set); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(mapset.Set)
		}
	}

	return r0
}
