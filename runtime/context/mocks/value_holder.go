// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ValueHolder is an autogenerated mock type for the ValueHolder type
type ValueHolder struct {
	mock.Mock
}

// SetValue provides a mock function with given fields: name, value
func (_m *ValueHolder) SetValue(name string, value interface{}) {
	_m.Called(name, value)
}