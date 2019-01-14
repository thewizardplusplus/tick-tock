package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandGroup_Run(test *testing.T) {
	for _, testData := range []struct {
		name         string
		makeCommands func(log *[]int) CommandGroup
		wantLog      []int
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name:         "success without commands",
			makeCommands: func(log *[]int) CommandGroup { return nil },
			wantErr:      assert.NoError,
		},
		{
			name:         "success with commands",
			makeCommands: func(log *[]int) CommandGroup { return newCalledLoggableCommands(log, 5, 0, -1) },
			wantLog:      []int{0, 1, 2, 3, 4},
			wantErr:      assert.NoError,
		},
		{
			name:         "error",
			makeCommands: func(log *[]int) CommandGroup { return newCalledLoggableCommands(log, 5, 0, 2) },
			wantLog:      []int{0, 1, 2},
			wantErr:      assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log []int
			commands := testData.makeCommands(&log)
			err := commands.Run()

			assert.Equal(test, testData.wantLog, log)
			checkCommands(test, commands)
			testData.wantErr(test, err)
		})
	}
}
