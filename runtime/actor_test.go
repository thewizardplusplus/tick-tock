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
			makeStates   func(log *[]int) StateGroup
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
				makeStates: func(log *[]int) StateGroup {
					return StateGroup{
						"state_one": MessageGroup{
							"message_one": newLoggableCommands(log, 5, 0),
							"message_two": newLoggableCommands(log, 5, 5),
						},
						"state_two": MessageGroup{
							"message_three": newLoggableCommands(log, 5, 10),
							"message_four":  newCalledLoggableCommands(log, 5, 15, -1),
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
				makeStates: func(log *[]int) StateGroup {
					return StateGroup{
						"state_one": MessageGroup{
							"message_one": newLoggableCommands(log, 5, 0),
							"message_two": newLoggableCommands(log, 5, 5),
						},
						"state_two": MessageGroup{
							"message_three": newLoggableCommands(log, 5, 10),
							"message_four":  newCalledLoggableCommands(log, 5, 15, 2),
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
			var log []int
			states := testData.fields.makeStates(&log)
			err := Actor{testData.fields.currentState, states}.ProcessMessage(testData.args.message)

			assert.Equal(test, testData.wantLog, log)
			checkStates(test, states)
			testData.wantErr(test, err)
		})
	}
}
