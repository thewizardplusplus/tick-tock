package runtime

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUnknownStateError(test *testing.T) {
	got := newUnknownStateError("test")
	assert.Equal(test, "unknown state test", got.Error())
}
