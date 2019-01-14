package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewActor(test *testing.T) {
	type args struct {
		initialState string
		states       StateGroup
	}

	for _, testData := range []struct {
		name    string
		args    args
		want    *Actor
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			args:    args{"state_one", StateGroup{"state_one": nil, "state_two": nil}},
			want:    &Actor{"state_one", StateGroup{"state_one": nil, "state_two": nil}},
			wantErr: assert.NoError,
		},
		{
			name:    "error",
			args:    args{"state_unknown", StateGroup{"state_one": nil, "state_two": nil}},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			got, err := NewActor(testData.args.initialState, testData.args.states)
			assert.Equal(test, testData.want, got)
			testData.wantErr(test, err)
		})
	}
}

func TestActor_SetState(test *testing.T) {
	type (
		fields struct {
			currentState string
			states       StateGroup
		}
		args struct {
			state string
		}
	)

	for _, testData := range []struct {
		name             string
		fields           fields
		args             args
		wantCurrentState string
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name:             "success with a different state",
			fields:           fields{"state_one", StateGroup{"state_one": nil, "state_two": nil}},
			args:             args{"state_two"},
			wantCurrentState: "state_two",
			wantErr:          assert.NoError,
		},
		{
			name:             "success with a same state",
			fields:           fields{"state_one", StateGroup{"state_one": nil, "state_two": nil}},
			args:             args{"state_one"},
			wantCurrentState: "state_one",
			wantErr:          assert.NoError,
		},
		{
			name:             "error",
			fields:           fields{"state_one", StateGroup{"state_one": nil, "state_two": nil}},
			args:             args{"state_unknown"},
			wantCurrentState: "state_one",
			wantErr:          assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			actor := Actor{testData.fields.currentState, testData.fields.states}
			err := actor.SetState(testData.args.state)
			assert.Equal(test, testData.wantCurrentState, actor.currentState)
			testData.wantErr(test, err)
		})
	}
}

func TestActor_ProcessMessage(test *testing.T) {
	type (
		fields struct {
			currentState string
			makeStates   func(context Context, log *commandLog) StateGroup
		}
		args struct {
			message string
		}
	)

	for _, testData := range []struct {
		name    string
		fields  fields
		args    args
		wantLog []int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				currentState: "state_two",
				makeStates: func(context Context, log *commandLog) StateGroup {
					return StateGroup{
						"state_one": MessageGroup{
							"message_one": newLoggableCommands(context, log, 5),
							"message_two": newLoggableCommands(context, log, 5, withIDFrom(5)),
						},
						"state_two": MessageGroup{
							"message_three": newLoggableCommands(context, log, 5, withIDFrom(10)),
							"message_four":  newLoggableCommands(context, log, 5, withIDFrom(15), withCalls()),
						},
					}
				},
			},
			args:    args{"message_four"},
			wantLog: []int{15, 16, 17, 18, 19},
			wantErr: assert.NoError,
		},
		{
			name: "error",
			fields: fields{
				currentState: "state_two",
				makeStates: func(context Context, log *commandLog) StateGroup {
					return StateGroup{
						"state_one": MessageGroup{
							"message_one": newLoggableCommands(context, log, 5),
							"message_two": newLoggableCommands(context, log, 5, withIDFrom(5)),
						},
						"state_two": MessageGroup{
							"message_three": newLoggableCommands(context, log, 5, withIDFrom(10)),
							"message_four":  newLoggableCommands(context, log, 5, withIDFrom(15), withErrOn(2)),
						},
					}
				},
			},
			args:    args{"message_four"},
			wantLog: []int{15, 16, 17},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var log commandLog
			states := testData.fields.makeStates(nil, &log)
			err := Actor{testData.fields.currentState, states}.ProcessMessage(testData.args.message)

			assert.Equal(test, testData.wantLog, log.commands)
			checkStates(test, states)
			testData.wantErr(test, err)
		})
	}
}
