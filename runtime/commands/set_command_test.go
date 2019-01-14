package commands

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestSetCommand(test *testing.T) {
	for _, testData := range []struct {
		name       string
		settingErr error
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			wantErr: assert.NoError,
		},
		{
			name:       "error",
			settingErr: iotest.ErrTimeout,
			wantErr:    assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			context := new(mocks.Context)
			context.On("SetState", "test").Return(testData.settingErr)

			err := NewSetCommand("test").Run(context)

			context.AssertExpectations(test)
			testData.wantErr(test, err)
		})
	}
}
