// Code generated by mockery v1.0.0. DO NOT EDIT.

package commands

import mock "github.com/stretchr/testify/mock"
import runtime "github.com/thewizardplusplus/tick-tock/runtime"

// MockContext is an autogenerated mock type for the Context type
type MockContext struct {
	mock.Mock
}

// SendMessage provides a mock function with given fields: message
func (_m *MockContext) SendMessage(message string) {
	_m.Called(message)
}

// SetActor provides a mock function with given fields: actor
func (_m *MockContext) SetActor(actor *runtime.Actor) {
	_m.Called(actor)
}

// SetActors provides a mock function with given fields: actors
func (_m *MockContext) SetActors(actors runtime.ConcurrentActorGroup) {
	_m.Called(actors)
}

// SetState provides a mock function with given fields: state
func (_m *MockContext) SetState(state string) error {
	ret := _m.Called(state)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(state)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}