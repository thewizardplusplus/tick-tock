package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCommandGroup_Run(test *testing.T) {
	for _, testData := range []struct {
		name         string
		makeCommands func(context Context, log *commandLog) CommandGroup
		wantLog      []int
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name:         "success without commands",
			makeCommands: func(context Context, log *commandLog) CommandGroup { return nil },
			wantErr:      assert.NoError,
		},
		{
			name: "success with commands",
			makeCommands: func(context Context, log *commandLog) CommandGroup {
				return newLoggableCommands(context, log, 5, withCalls())
			},
			wantLog: []int{0, 1, 2, 3, 4},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			makeCommands: func(context Context, log *commandLog) CommandGroup {
				return newLoggableCommands(context, log, 5, withErrOn(2))
			},
			wantLog: []int{0, 1, 2},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log commandLog
			commands := testData.makeCommands(nil, &log)
			err := commands.Run()

			assert.Equal(test, testData.wantLog, log.commands)
			checkCommands(test, commands)
			testData.wantErr(test, err)
		})
	}
}
