// Code generated by mockery v1.0.0. DO NOT EDIT.
package runtime

import mock "github.com/stretchr/testify/mock"

// MockWaiter is an autogenerated mock type for the Waiter type
type MockWaiter struct {
	mock.Mock
}

// Add provides a mock function with given fields: delta
func (_m *MockWaiter) Add(delta int) {
	_m.Called(delta)
}

// Done provides a mock function with given fields:
func (_m *MockWaiter) Done() {
	_m.Called()
}
