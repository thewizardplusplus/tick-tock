package runtime

import (
	"testing"
	"testing/iotest"

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
			makeCommands: func(log *[]int) CommandGroup { return makeCommands(log, 5, 5) },
			wantLog:      []int{0, 1, 2, 3, 4},
			wantErr:      assert.NoError,
		},
		{
			name:         "error",
			makeCommands: func(log *[]int) CommandGroup { return makeCommands(log, 5, 2) },
			wantLog:      []int{0, 1, 2},
			wantErr:      assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log []int
			commands := testData.makeCommands(&log)
			err := commands.Run()

			assert.Equal(test, testData.wantLog, log)
			checkLoggableCommands(test, commands)
			testData.wantErr(test, err)
		})
	}
}

func makeCommands(log *[]int, count int, errIndex int) CommandGroup {
	commands := newLoggableCommands(log, count, 0)
	// expect execution of all commands from first to failed one, inclusive
	for i := 0; i <= errIndex; i++ {
		var err error
		// return an error from a failed command
		if i == errIndex {
			err = iotest.ErrTimeout
		}

		if i < len(commands) {
			commands[i].(*loggableCommand).On("Run").Return(err)
		}
	}

	return commands
}
