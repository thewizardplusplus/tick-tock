package builtin

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
	"github.com/thewizardplusplus/tick-tock/translator"
)

func TestValues(test *testing.T) {
	for _, data := range []struct {
		name       string
		expression expressions.Expression
		wantResult interface{}
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "empty list",
			expression: expressions.NewIdentifier(translator.EmptyListConstantName),
			wantResult: (*types.Pair)(nil),
			wantErr:    assert.NoError,
		},
		{
			name: "list construction",
			expression: expressions.NewFunctionCall(
				translator.ListConstructionFunctionName,
				[]expressions.Expression{
					expressions.NewNumber(12),
					expressions.NewFunctionCall(
						translator.ListConstructionFunctionName,
						[]expressions.Expression{
							expressions.NewNumber(23),
							expressions.NewFunctionCall(
								translator.ListConstructionFunctionName,
								[]expressions.Expression{
									expressions.NewNumber(42),
									expressions.NewIdentifier(translator.EmptyListConstantName),
								},
							),
						},
					),
				},
			),
			wantResult: types.NewPairFromSlice([]interface{}{12.0, 23.0, 42.0}),
			wantErr:    assert.NoError,
		},
		{
			name: "addition/success/float64",
			expression: expressions.NewFunctionCall(
				translator.AdditionFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewNumber(3)},
			),
			wantResult: 5.0,
			wantErr:    assert.NoError,
		},
		{
			name: "addition/success/*types.Pair",
			expression: expressions.NewFunctionCall(
				translator.AdditionFunctionName,
				[]expressions.Expression{expressions.NewString("te"), expressions.NewString("st")},
			),
			wantResult: types.NewPairFromText("test"),
			wantErr:    assert.NoError,
		},
		{
			name: "addition/error/argument #0",
			expression: expressions.NewFunctionCall(
				translator.AdditionFunctionName,
				[]expressions.Expression{
					expressions.NewIdentifier(translator.AdditionFunctionName),
					expressions.NewIdentifier(translator.EmptyListConstantName),
				},
			),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "addition/error/argument #1",
			expression: expressions.NewFunctionCall(
				translator.AdditionFunctionName,
				[]expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewIdentifier(translator.EmptyListConstantName),
				},
			),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "subtraction",
			expression: expressions.NewFunctionCall(
				translator.SubtractionFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewNumber(3)},
			),
			wantResult: -1.0,
			wantErr:    assert.NoError,
		},
		{
			name: "multiplication",
			expression: expressions.NewFunctionCall(
				translator.MultiplicationFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewNumber(3)},
			),
			wantResult: 6.0,
			wantErr:    assert.NoError,
		},
		{
			name: "division",
			expression: expressions.NewFunctionCall(
				translator.DivisionFunctionName,
				[]expressions.Expression{expressions.NewNumber(10), expressions.NewNumber(2)},
			),
			wantResult: 5.0,
			wantErr:    assert.NoError,
		},
		{
			name: "modulo",
			expression: expressions.NewFunctionCall(translator.ModuloFunctionName, []expressions.Expression{
				expressions.NewNumber(10),
				expressions.NewNumber(3),
			}),
			wantResult: 1.0,
			wantErr:    assert.NoError,
		},
		{
			name: "negation",
			expression: expressions.NewFunctionCall(
				translator.NegationFunctionName,
				[]expressions.Expression{expressions.NewNumber(23)},
			),
			wantResult: -23.0,
			wantErr:    assert.NoError,
		},
		{
			name: "key accessor/success",
			expression: expressions.NewFunctionCall(
				translator.KeyAccessorFunctionName,
				[]expressions.Expression{
					expressions.NewFunctionCall(
						translator.ListConstructionFunctionName,
						[]expressions.Expression{
							expressions.NewNumber(12),
							expressions.NewFunctionCall(
								translator.ListConstructionFunctionName,
								[]expressions.Expression{
									expressions.NewNumber(23),
									expressions.NewFunctionCall(
										translator.ListConstructionFunctionName,
										[]expressions.Expression{
											expressions.NewNumber(42),
											expressions.NewIdentifier(translator.EmptyListConstantName),
										},
									),
								},
							),
						},
					),
					expressions.NewNumber(1),
				},
			),
			wantResult: 23.0,
			wantErr:    assert.NoError,
		},
		{
			name: "key accessor/error",
			expression: expressions.NewFunctionCall(
				translator.KeyAccessorFunctionName,
				[]expressions.Expression{
					expressions.NewFunctionCall(
						translator.ListConstructionFunctionName,
						[]expressions.Expression{
							expressions.NewNumber(12),
							expressions.NewFunctionCall(
								translator.ListConstructionFunctionName,
								[]expressions.Expression{
									expressions.NewNumber(23),
									expressions.NewFunctionCall(
										translator.ListConstructionFunctionName,
										[]expressions.Expression{
											expressions.NewNumber(42),
											expressions.NewIdentifier(translator.EmptyListConstantName),
										},
									),
								},
							),
						},
					),
					expressions.NewNumber(23),
				},
			),
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		ctx := context.NewDefaultContext()
		context.SetValues(ctx, Values)

		gotResult, gotErr := data.expression.Evaluate(ctx)
		assert.Equal(test, data.wantResult, gotResult)
		data.wantErr(test, gotErr)
	}
}
