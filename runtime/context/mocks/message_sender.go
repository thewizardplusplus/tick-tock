// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// MessageSender is an autogenerated mock type for the MessageSender type
type MessageSender struct {
	mock.Mock
}

// SendMessage provides a mock function with given fields: message
func (_m *MessageSender) SendMessage(message string) {
	_m.Called(message)
}