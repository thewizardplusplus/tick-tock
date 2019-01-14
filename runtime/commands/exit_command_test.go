package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExitCommand(test *testing.T) {
	context := new(MockContext)
	assert.Panics(test, func() { ExitCommand{}.Run(context) })
	context.AssertExpectations(test)
}
