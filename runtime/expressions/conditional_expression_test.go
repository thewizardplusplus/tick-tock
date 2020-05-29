package expressions

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/mocks"
)

type SignedCommand struct {
	*mocks.Command

	Sign string
}

func NewSignedCommand(sign string) SignedCommand {
	return SignedCommand{new(mocks.Command), sign}
}

func TestNewConditionalExpression(test *testing.T) {
	conditionalCases := []ConditionalCase{
		{NewSignedExpression("one-condition"), NewSignedCommand("one-command")},
		{NewSignedExpression("two-condition"), NewSignedCommand("two-command")},
	}
	got := NewConditionalExpression(conditionalCases)

	checkConditionalCases(test, conditionalCases)
	assert.Equal(test, conditionalCases, got.conditionalCases)
}

func TestConditionalExpression_Evaluate(test *testing.T) {
	type fields struct {
		conditionalCases []ConditionalCase
	}
	type args struct {
		context context.Context
	}

	for _, data := range []struct {
		name       string
		fields     fields
		args       args
		wantResult interface{}
		wantErr    assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	} {
		test.Run(data.name, func(test *testing.T) {
			expression := ConditionalExpression{
				conditionalCases: data.fields.conditionalCases,
			}
			gotResult, gotErr := expression.Evaluate(data.args.context)

			checkConditionalCases(test, data.fields.conditionalCases)
			mock.AssertExpectationsForObjects(test, data.args.context)
			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}

func checkConditionalCases(test *testing.T, conditionalCases []ConditionalCase) {
	for _, conditionalCase := range conditionalCases {
		mock.AssertExpectationsForObjects(test, conditionalCase.Condition, conditionalCase.Command)
	}
}
