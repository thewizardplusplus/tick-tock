package commands

import (
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

func TestSetCommand(test *testing.T) {
	for _, testData := range []struct {
		name       string
		settingErr error
		wantResult interface{}
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "success",
			settingErr: nil,
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name:       "error",
			settingErr: iotest.ErrTimeout,
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			context := new(mocks.Context)
			context.On("SetState", "test").Return(testData.settingErr)

			gotResult, gotErr := NewSetCommand("test").Run(context)

			mock.AssertExpectationsForObjects(test, context)
			assert.Equal(test, testData.wantResult, gotResult)
			testData.wantErr(test, gotErr)
		})
	}
}
