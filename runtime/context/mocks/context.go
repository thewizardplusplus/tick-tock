// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "github.com/thewizardplusplus/tick-tock/runtime/context"
import mock "github.com/stretchr/testify/mock"

// Context is an autogenerated mock type for the Context type
type Context struct {
	mock.Mock
}

// Copy provides a mock function with given fields:
func (_m *Context) Copy() context.Context {
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

// SendMessage provides a mock function with given fields: message
func (_m *Context) SendMessage(message string) {
	_m.Called(message)
}

// SetMessageSender provides a mock function with given fields: sender
func (_m *Context) SetMessageSender(sender context.MessageSender) {
	_m.Called(sender)
}

// SetState provides a mock function with given fields: state
func (_m *Context) SetState(state string) error {
	ret := _m.Called(state)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(state)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetStateHolder provides a mock function with given fields: holder
func (_m *Context) SetStateHolder(holder context.StateHolder) {
	_m.Called(holder)
}
