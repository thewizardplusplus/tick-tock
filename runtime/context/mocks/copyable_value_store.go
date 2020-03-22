// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "github.com/thewizardplusplus/tick-tock/runtime/context"
import mock "github.com/stretchr/testify/mock"

// CopyableValueStore is an autogenerated mock type for the CopyableValueStore type
type CopyableValueStore struct {
	mock.Mock
}

// Copy provides a mock function with given fields:
func (_m *CopyableValueStore) Copy() context.CopyableValueStore {
	ret := _m.Called()

	var r0 context.CopyableValueStore
	if rf, ok := ret.Get(0).(func() context.CopyableValueStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.CopyableValueStore)
		}
	}

	return r0
}

// SetValue provides a mock function with given fields: name, value
func (_m *CopyableValueStore) SetValue(name string, value interface{}) {
	_m.Called(name, value)
}

// Value provides a mock function with given fields: name
func (_m *CopyableValueStore) Value(name string) (interface{}, bool) {
	ret := _m.Called(name)

	var r0 interface{}
	if rf, ok := ret.Get(0).(func(string) interface{}); ok {
		r0 = rf(name)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	var r1 bool
	if rf, ok := ret.Get(1).(func(string) bool); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Get(1).(bool)
	}

	return r0, r1
}

// ValuesNames provides a mock function with given fields:
func (_m *CopyableValueStore) ValuesNames() context.ValueNameGroup {
	ret := _m.Called()

	var r0 context.ValueNameGroup
	if rf, ok := ret.Get(0).(func() context.ValueNameGroup); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.ValueNameGroup)
		}
	}

	return r0
}
