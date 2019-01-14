package commands

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	testsmocks "github.com/thewizardplusplus/tick-tock/internal/tests/mocks"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
)

func TestNewSleepCommand(test *testing.T) {
	type args struct {
		minimum float64
		maximum float64
	}

	for _, testData := range []struct {
		name    string
		args    args
		want    SleepCommand
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			args:    args{1.2, 3.4},
			want:    SleepCommand{1.2, 3.4, SleepDependencies{}},
			wantErr: assert.NoError,
		},
		{
			name:    "error with a negative minimum",
			args:    args{-1.2, 3.4},
			want:    SleepCommand{},
			wantErr: assert.Error,
		},
		{
			name:    "error with a negative maximum",
			args:    args{1.2, -3.4},
			want:    SleepCommand{},
			wantErr: assert.Error,
		},
		{
			name:    "error with a maximum less a minimum",
			args:    args{3.4, 1.2},
			want:    SleepCommand{},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			randomizer := new(testsmocks.Randomizer)
			sleeper := new(testsmocks.Sleeper)
			dependencies := SleepDependencies{randomizer.Randomize, sleeper.Sleep}

			got, err := NewSleepCommand(testData.args.minimum, testData.args.maximum, dependencies)
			got.dependencies = SleepDependencies{}

			mock.AssertExpectationsForObjects(test, randomizer, sleeper)
			assert.Equal(test, testData.want, got)
			testData.wantErr(test, err)
		})
	}
}

func TestSleepCommand_Run(test *testing.T) {
	type fields struct {
		minimum float64
		maximum float64
	}

	for _, testData := range []struct {
		name                   string
		fields                 fields
		initializeDependencies func(randomizer *testsmocks.Randomizer, sleeper *testsmocks.Sleeper)
	}{
		{
			name:   "success with a maximum greater a minimum",
			fields: fields{1.2, 3.4},
			initializeDependencies: func(randomizer *testsmocks.Randomizer, sleeper *testsmocks.Sleeper) {
				randomizer.On("Randomize").Return(0.25)
				sleeper.On("Sleep", 1750*time.Millisecond).Return()
			},
		},
		{
			name:   "success with a maximum equal a minimum",
			fields: fields{1.2, 1.2},
			initializeDependencies: func(randomizer *testsmocks.Randomizer, sleeper *testsmocks.Sleeper) {
				randomizer.On("Randomize").Return(0.25)
				sleeper.On("Sleep", 1200*time.Millisecond).Return()
			},
		},
		{
			name: "success with a zero maximum and a zero minimum",
			initializeDependencies: func(randomizer *testsmocks.Randomizer, sleeper *testsmocks.Sleeper) {
				randomizer.On("Randomize").Return(0.25)
				sleeper.On("Sleep", time.Duration(0)).Return()
			},
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			randomizer := new(testsmocks.Randomizer)
			sleeper := new(testsmocks.Sleeper)
			testData.initializeDependencies(randomizer, sleeper)

			dependencies := SleepDependencies{randomizer.Randomize, sleeper.Sleep}
			context := new(contextmocks.Context)
			err := SleepCommand{testData.fields.minimum, testData.fields.maximum, dependencies}.Run(context)

			mock.AssertExpectationsForObjects(test, randomizer, sleeper, context)
			assert.NoError(test, err)
		})
	}
}
