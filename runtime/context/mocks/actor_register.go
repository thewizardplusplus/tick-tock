// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import context "github.com/thewizardplusplus/tick-tock/runtime/context"
import mock "github.com/stretchr/testify/mock"

// ActorRegister is an autogenerated mock type for the ActorRegister type
type ActorRegister struct {
	mock.Mock
}

// RegisterActor provides a mock function with given fields: actor
func (_m *ActorRegister) RegisterActor(actor context.Actor) {
	_m.Called(actor)
}
