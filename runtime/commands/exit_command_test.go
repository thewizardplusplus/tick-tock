package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestExitCommand(test *testing.T) {
	context := new(mocks.Context)
	err := ExitCommand{}.Run(context)

	mock.AssertExpectationsForObjects(test, context)
	assert.Equal(test, runtime.ErrUserExit, err)
}
