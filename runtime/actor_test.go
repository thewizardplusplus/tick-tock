package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
