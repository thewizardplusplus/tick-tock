package expressions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

func TestNewString(test *testing.T) {
	got := NewString("hi")

	want := &types.Pair{Head: float64('h'), Tail: &types.Pair{Head: float64('i'), Tail: nil}}
	assert.Equal(test, want, got.value)
}

func TestString_Evaluate(test *testing.T) {
	context := new(mocks.Context)
	string := String{
		value: &types.Pair{Head: float64('h'), Tail: &types.Pair{Head: float64('i'), Tail: nil}},
	}
	gotResult, gotErr := string.Evaluate(context)

	mock.AssertExpectationsForObjects(test, context)
	wantResult := &types.Pair{Head: float64('h'), Tail: &types.Pair{Head: float64('i'), Tail: nil}}
	assert.Equal(test, wantResult, gotResult)
	assert.NoError(test, gotErr)
}
