// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// StateHolder is an autogenerated mock type for the StateHolder type
type StateHolder struct {
	mock.Mock
}

// SetState provides a mock function with given fields: state
func (_m *StateHolder) SetState(state string) error {
	ret := _m.Called(state)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(state)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}