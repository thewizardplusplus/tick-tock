// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ErrorHandler is an autogenerated mock type for the ErrorHandler type
type ErrorHandler struct {
	mock.Mock
}

// HandleError provides a mock function with given fields: err
func (_m *ErrorHandler) HandleError(err error) {
	_m.Called(err)
}