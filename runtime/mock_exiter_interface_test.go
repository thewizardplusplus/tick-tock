// Code generated by mockery v1.0.0. DO NOT EDIT.

package runtime

import mock "github.com/stretchr/testify/mock"

// MockExiterInterface is an autogenerated mock type for the ExiterInterface type
type MockExiterInterface struct {
	mock.Mock
}

// Exit provides a mock function with given fields: code
func (_m *MockExiterInterface) Exit(code int) {
	_m.Called(code)
}
