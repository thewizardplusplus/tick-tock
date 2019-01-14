package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestExitCommand(test *testing.T) {
	context := new(mocks.Context)
	assert.Panics(test, func() { ExitCommand{}.Run(context) })
	context.AssertExpectations(test)
}
