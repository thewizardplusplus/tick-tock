package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnknownStateError_Error(test *testing.T) {
	for _, testData := range []struct {
		name string
		err  unknownStateError
		want string
	}{
		{
			name: "nonempty",
			err:  "test",
			want: "unknown state test",
		},
		{
			name: "empty",
			want: "unknown state <empty>",
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			got := testData.err.Error()
			assert.Equal(test, testData.want, got)
		})
	}
}
