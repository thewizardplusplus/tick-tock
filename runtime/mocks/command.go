// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "github.com/thewizardplusplus/tick-tock/runtime/context"
import mock "github.com/stretchr/testify/mock"

// Command is an autogenerated mock type for the Command type
type Command struct {
	mock.Mock
}

// Run provides a mock function with given fields: _a0
func (_m *Command) Run(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}