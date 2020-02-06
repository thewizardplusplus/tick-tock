// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "github.com/thewizardplusplus/tick-tock/runtime/context"

import mock "github.com/stretchr/testify/mock"

// Expression is an autogenerated mock type for the Expression type
type Expression struct {
	mock.Mock
}

// Evaluate provides a mock function with given fields: _a0
func (_m *Expression) Evaluate(_a0 context.Context) (interface{}, error) {
	ret := _m.Called(_a0)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(context.Context) interface{}); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
