package translator

import (
	"testing"

	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/commands"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

func TestTranslateExpression(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "Expression/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Expression/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Expression/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			expression := new(parser.Expression)
			err := parser.ParseToAST(data.args.code, expression)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				TranslateExpression(expression, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateListConstruction(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "ListConstruction/nonempty/success",
			args: args{
				code:                "12 : test",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				ListConstructionFunctionName,
				[]expressions.Expression{expressions.NewNumber(12), expressions.NewIdentifier("test")},
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "ListConstruction/nonempty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
					: when
						=> 24
							set two()
						=> 43
							set three()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				ListConstructionFunctionName,
				[]expressions.Expression{
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(23),
							Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
						},
						{
							Condition: expressions.NewNumber(42),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
					}),
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(24),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
						{
							Condition: expressions.NewNumber(43),
							Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
						},
					}),
				},
			),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "ListConstruction/nonempty/error",
			args: args{
				code:                "12 : unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "ListConstruction/empty/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "ListConstruction/empty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "ListConstruction/empty/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			listConstruction := new(parser.ListConstruction)
			err := parser.ParseToAST(data.args.code, listConstruction)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateListConstruction(listConstruction, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateNilCoalescing(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "NilCoalescing/nonempty/success",
			args: args{
				code:                "12 ?? 23 ?? 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewNilCoalescingOperator(
				expressions.NewNumber(12),
				expressions.NewNilCoalescingOperator(
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				),
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "NilCoalescing/nonempty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
					?? when
						=> 24
							set two()
						=> 43
							set three()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewNilCoalescingOperator(
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(23),
						Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
					},
					{
						Condition: expressions.NewNumber(42),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
				}),
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(24),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
					{
						Condition: expressions.NewNumber(43),
						Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
					},
				}),
			),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "NilCoalescing/nonempty/error",
			args: args{
				code:                "12 ?? 23 ?? unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "NilCoalescing/empty/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "NilCoalescing/empty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "NilCoalescing/empty/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			nilCoalescing := new(parser.NilCoalescing)
			err := parser.ParseToAST(data.args.code, nilCoalescing)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateNilCoalescing(nilCoalescing, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateDisjunction(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "Disjunction/nonempty/success",
			args: args{
				code:                "12 || 23 || 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewBooleanOperator(
				expressions.NewNumber(12),
				expressions.NewBooleanOperator(
					expressions.NewNumber(23),
					expressions.NewNumber(42),
					types.True,
				),
				types.True,
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Disjunction/nonempty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
					|| when
						=> 24
							set two()
						=> 43
							set three()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewBooleanOperator(
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(23),
						Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
					},
					{
						Condition: expressions.NewNumber(42),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
				}),
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(24),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
					{
						Condition: expressions.NewNumber(43),
						Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
					},
				}),
				types.True,
			),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "Disjunction/nonempty/error",
			args: args{
				code:                "12 || 23 || unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Disjunction/empty/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Disjunction/empty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Disjunction/empty/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			disjunction := new(parser.Disjunction)
			err := parser.ParseToAST(data.args.code, disjunction)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateDisjunction(disjunction, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateConjunction(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "Conjunction/nonempty/success",
			args: args{
				code:                "12 && 23 && 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewBooleanOperator(
				expressions.NewNumber(12),
				expressions.NewBooleanOperator(
					expressions.NewNumber(23),
					expressions.NewNumber(42),
					types.False,
				),
				types.False,
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Conjunction/nonempty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
					&& when
						=> 24
							set two()
						=> 43
							set three()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewBooleanOperator(
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(23),
						Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
					},
					{
						Condition: expressions.NewNumber(42),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
				}),
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(24),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
					{
						Condition: expressions.NewNumber(43),
						Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
					},
				}),
				types.False,
			),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "Conjunction/nonempty/error",
			args: args{
				code:                "12 && 23 && unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Conjunction/empty/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Conjunction/empty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Conjunction/empty/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			conjunction := new(parser.Conjunction)
			err := parser.ParseToAST(data.args.code, conjunction)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateConjunction(conjunction, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateBinaryOperation_equality(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "Equality/nonempty/success",
			args: args{
				code:                "12 == 23 != 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(EqualFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(NotEqualFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Equality/nonempty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
					== when
						=> 24
							set two()
						=> 43
							set three()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(EqualFunctionName, []expressions.Expression{
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(23),
						Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
					},
					{
						Condition: expressions.NewNumber(42),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
				}),
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(24),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
					{
						Condition: expressions.NewNumber(43),
						Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
					},
				}),
			}),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "Equality/nonempty/error",
			args: args{
				code:                "12 == 23 != unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Equality/empty/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Equality/empty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Equality/empty/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			equality := new(parser.Equality)
			err := parser.ParseToAST(data.args.code, equality)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateBinaryOperation(equality, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateBinaryOperation_comparison(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "Comparison/nonempty/success/less",
			args: args{
				code:                "12 < 23 < 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(LessFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(LessFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Comparison/nonempty/success/less or equal",
			args: args{
				code:                "12 <= 23 <= 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(LessOrEqualFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(LessOrEqualFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Comparison/nonempty/success/great",
			args: args{
				code:                "12 > 23 > 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(GreaterFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(GreaterFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Comparison/nonempty/success/great or equal",
			args: args{
				code:                "12 >= 23 >= 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(GreaterOrEqualFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(GreaterOrEqualFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Comparison/nonempty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
					< when
						=> 24
							set two()
						=> 43
							set three()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(LessFunctionName, []expressions.Expression{
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(23),
						Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
					},
					{
						Condition: expressions.NewNumber(42),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
				}),
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(24),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
					{
						Condition: expressions.NewNumber(43),
						Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
					},
				}),
			}),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "Comparison/nonempty/error",
			args: args{
				code:                "12 < 23 < unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Comparison/empty/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Comparison/empty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Comparison/empty/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			comparison := new(parser.Comparison)
			err := parser.ParseToAST(data.args.code, comparison)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateBinaryOperation(comparison, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateBinaryOperation_bitwiseDisjunction(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "BitwiseDisjunction/nonempty/success",
			args: args{
				code:                "12 | 23 | 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				BitwiseDisjunctionFunctionName,
				[]expressions.Expression{
					expressions.NewNumber(12),
					expressions.NewFunctionCall(BitwiseDisjunctionFunctionName, []expressions.Expression{
						expressions.NewNumber(23),
						expressions.NewNumber(42),
					}),
				},
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseDisjunction/nonempty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
					| when
						=> 24
							set two()
						=> 43
							set three()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				BitwiseDisjunctionFunctionName,
				[]expressions.Expression{
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(23),
							Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
						},
						{
							Condition: expressions.NewNumber(42),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
					}),
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(24),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
						{
							Condition: expressions.NewNumber(43),
							Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
						},
					}),
				},
			),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseDisjunction/nonempty/error",
			args: args{
				code:                "12 | 23 | unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "BitwiseDisjunction/empty/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseDisjunction/empty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseDisjunction/empty/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			bitwiseDisjunction := new(parser.BitwiseDisjunction)
			err := parser.ParseToAST(data.args.code, bitwiseDisjunction)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateBinaryOperation(bitwiseDisjunction, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateBinaryOperation_bitwiseExclusiveDisjunction(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "BitwiseExclusiveDisjunction/nonempty/success",
			args: args{
				code:                "12 ^ 23 ^ 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				BitwiseExclusiveDisjunctionFunctionName,
				[]expressions.Expression{
					expressions.NewNumber(12),
					expressions.NewFunctionCall(BitwiseExclusiveDisjunctionFunctionName, []expressions.Expression{
						expressions.NewNumber(23),
						expressions.NewNumber(42),
					}),
				},
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseExclusiveDisjunction/nonempty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
					^ when
						=> 24
							set two()
						=> 43
							set three()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				BitwiseExclusiveDisjunctionFunctionName,
				[]expressions.Expression{
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(23),
							Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
						},
						{
							Condition: expressions.NewNumber(42),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
					}),
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(24),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
						{
							Condition: expressions.NewNumber(43),
							Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
						},
					}),
				},
			),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseExclusiveDisjunction/nonempty/error",
			args: args{
				code:                "12 ^ 23 ^ unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "BitwiseExclusiveDisjunction/empty/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseExclusiveDisjunction/empty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseExclusiveDisjunction/empty/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			bitwiseExclusiveDisjunction := new(parser.BitwiseExclusiveDisjunction)
			err := parser.ParseToAST(data.args.code, bitwiseExclusiveDisjunction)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateBinaryOperation(bitwiseExclusiveDisjunction, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateBinaryOperation_bitwiseConjunction(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "BitwiseConjunction/nonempty/success",
			args: args{
				code:                "12 & 23 & 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				BitwiseConjunctionFunctionName,
				[]expressions.Expression{
					expressions.NewNumber(12),
					expressions.NewFunctionCall(BitwiseConjunctionFunctionName, []expressions.Expression{
						expressions.NewNumber(23),
						expressions.NewNumber(42),
					}),
				},
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseConjunction/nonempty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
					& when
						=> 24
							set two()
						=> 43
							set three()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				BitwiseConjunctionFunctionName,
				[]expressions.Expression{
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(23),
							Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
						},
						{
							Condition: expressions.NewNumber(42),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
					}),
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(24),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
						{
							Condition: expressions.NewNumber(43),
							Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
						},
					}),
				},
			),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseConjunction/nonempty/error",
			args: args{
				code:                "12 & 23 & unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "BitwiseConjunction/empty/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseConjunction/empty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseConjunction/empty/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			bitwiseConjunction := new(parser.BitwiseConjunction)
			err := parser.ParseToAST(data.args.code, bitwiseConjunction)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateBinaryOperation(bitwiseConjunction, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateBinaryOperation_shift(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "Shift/nonempty/success",
			args: args{
				code:                "5 << 12 >> 23 >>> 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				BitwiseLeftShiftFunctionName,
				[]expressions.Expression{
					expressions.NewNumber(5),
					expressions.NewFunctionCall(BitwiseRightShiftFunctionName, []expressions.Expression{
						expressions.NewNumber(12),
						expressions.NewFunctionCall(BitwiseUnsignedRightShiftFunctionName, []expressions.Expression{
							expressions.NewNumber(23),
							expressions.NewNumber(42),
						}),
					}),
				},
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Shift/nonempty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
					<< when
						=> 24
							set two()
						=> 43
							set three()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				BitwiseLeftShiftFunctionName,
				[]expressions.Expression{
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(23),
							Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
						},
						{
							Condition: expressions.NewNumber(42),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
					}),
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(24),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
						{
							Condition: expressions.NewNumber(43),
							Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
						},
					}),
				},
			),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "Shift/nonempty/error",
			args: args{
				code:                "5 << 12 >> 23 >>> unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "Shift/empty/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Shift/empty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Shift/empty/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			shift := new(parser.Shift)
			err := parser.ParseToAST(data.args.code, shift)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateBinaryOperation(shift, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateBinaryOperation_addition(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "Addition/nonempty/success/addition",
			args: args{
				code:                "12 + 23 + 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(AdditionFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(AdditionFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Addition/nonempty/success/subtraction",
			args: args{
				code:                "12 - 23 - 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(SubtractionFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(SubtractionFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Addition/nonempty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
					+ when
						=> 24
							set two()
						=> 43
							set three()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(AdditionFunctionName, []expressions.Expression{
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(23),
						Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
					},
					{
						Condition: expressions.NewNumber(42),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
				}),
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(24),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
					{
						Condition: expressions.NewNumber(43),
						Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
					},
				}),
			}),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "Addition/nonempty/error",
			args: args{
				code:                "12 + 23 + unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Addition/empty/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Addition/empty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Addition/empty/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			addition := new(parser.Addition)
			err := parser.ParseToAST(data.args.code, addition)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateBinaryOperation(addition, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateBinaryOperation_multiplication(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "Multiplication/nonempty/success/multiplication",
			args: args{
				code:                "12 * 23 * 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(MultiplicationFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(MultiplicationFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Multiplication/nonempty/success/division",
			args: args{
				code:                "12 / 23 / 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(DivisionFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(DivisionFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Multiplication/nonempty/success/modulo",
			args: args{
				code:                "12 % 23 % 42",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(ModuloFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(ModuloFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Multiplication/nonempty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
					* when
						=> 24
							set two()
						=> 43
							set three()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(MultiplicationFunctionName, []expressions.Expression{
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(23),
						Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
					},
					{
						Condition: expressions.NewNumber(42),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
				}),
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(24),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
					{
						Condition: expressions.NewNumber(43),
						Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
					},
				}),
			}),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "Multiplication/nonempty/error",
			args: args{
				code:                "12 * 23 * unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Multiplication/empty/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Multiplication/empty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Multiplication/empty/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			multiplication := new(parser.Multiplication)
			err := parser.ParseToAST(data.args.code, multiplication)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateBinaryOperation(multiplication, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateUnary(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "Unary/nonempty/success",
			args: args{
				code:                "-~!23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				ArithmeticNegationFunctionName,
				[]expressions.Expression{
					expressions.NewFunctionCall(BitwiseNegationFunctionName, []expressions.Expression{
						expressions.NewFunctionCall(LogicalNegationFunctionName, []expressions.Expression{
							expressions.NewNumber(23),
						}),
					}),
				},
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Unary/nonempty/success/with setted states",
			args: args{
				code: `
					-~!when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				ArithmeticNegationFunctionName,
				[]expressions.Expression{
					expressions.NewFunctionCall(BitwiseNegationFunctionName, []expressions.Expression{
						expressions.NewFunctionCall(LogicalNegationFunctionName, []expressions.Expression{
							expressions.NewConditionalExpression([]expressions.ConditionalCase{
								{
									Condition: expressions.NewNumber(23),
									Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
								},
								{
									Condition: expressions.NewNumber(42),
									Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
								},
							}),
						}),
					}),
				},
			),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Unary/nonempty/error",
			args: args{
				code:                "-~!unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Unary/empty/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Unary/empty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Unary/empty/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			unary := new(parser.Unary)
			err := parser.ParseToAST(data.args.code, unary)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateUnary(unary, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateAccessor(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "Accessor/nonempty/success/expressions",
			args: args{
				code:                "test[12][23]",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(KeyAccessorFunctionName, []expressions.Expression{
				expressions.NewFunctionCall(KeyAccessorFunctionName, []expressions.Expression{
					expressions.NewIdentifier("test"),
					expressions.NewNumber(12),
				}),
				expressions.NewNumber(23),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Accessor/nonempty/success/names",
			args: args{
				code:                "test.one.two",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(KeyAccessorFunctionName, []expressions.Expression{
				expressions.NewFunctionCall(KeyAccessorFunctionName, []expressions.Expression{
					expressions.NewIdentifier("test"),
					expressions.NewString("one"),
				}),
				expressions.NewString("two"),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Accessor/nonempty/success/names and expressions",
			args: args{
				code:                "test.one[12].two[23]",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(KeyAccessorFunctionName, []expressions.Expression{
				expressions.NewFunctionCall(KeyAccessorFunctionName, []expressions.Expression{
					expressions.NewFunctionCall(KeyAccessorFunctionName, []expressions.Expression{
						expressions.NewFunctionCall(KeyAccessorFunctionName, []expressions.Expression{
							expressions.NewIdentifier("test"),
							expressions.NewString("one"),
						}),
						expressions.NewNumber(12),
					}),
					expressions.NewString("two"),
				}),
				expressions.NewNumber(23),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Accessor/nonempty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;[
						when
							=> 24
								set two()
							=> 43
								set three()
						;
					]
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(KeyAccessorFunctionName, []expressions.Expression{
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(23),
						Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
					},
					{
						Condition: expressions.NewNumber(42),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
				}),
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(24),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
					{
						Condition: expressions.NewNumber(43),
						Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
					},
				}),
			}),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "Accessor/nonempty/error/atom translating",
			args: args{
				code:                "unknown[12][23]",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Accessor/nonempty/error/key translating",
			args: args{
				code:                "test[12][unknown]",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Accessor/empty/success",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Accessor/empty/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Accessor/empty/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			accessor := new(parser.Accessor)
			err := parser.ParseToAST(data.args.code, accessor)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateAccessor(accessor, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateAtom(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "Atom/number/integer",
			args: args{
				code:                "23",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/number/floating-point",
			args: args{
				code:                "2.3",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(2.3),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/symbol/latin1",
			args: args{
				code:                "'t'",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(116),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/symbol/not latin1",
			args: args{
				code:                "'т'",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(1090),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/string",
			args: args{
				code:                `"test"`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewString("test"),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/list definition/success",
			args: args{
				code:                "[12, 23, 42]",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				ListConstructionFunctionName,
				[]expressions.Expression{
					expressions.NewNumber(12),
					expressions.NewFunctionCall(ListConstructionFunctionName, []expressions.Expression{
						expressions.NewNumber(23),
						expressions.NewFunctionCall(ListConstructionFunctionName, []expressions.Expression{
							expressions.NewNumber(42),
							expressions.NewIdentifier(EmptyListConstantName),
						}),
					}),
				},
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/list definition/success/with setted states",
			args: args{
				code: `[
					when
						=> 23
							set one()
						=> 42
							set two()
					;,
					when
						=> 24
							set two()
						=> 43
							set three()
					;,
				]`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				ListConstructionFunctionName,
				[]expressions.Expression{
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(23),
							Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
						},
						{
							Condition: expressions.NewNumber(42),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
					}),
					expressions.NewFunctionCall(
						ListConstructionFunctionName,
						[]expressions.Expression{
							expressions.NewConditionalExpression([]expressions.ConditionalCase{
								{
									Condition: expressions.NewNumber(24),
									Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
								},
								{
									Condition: expressions.NewNumber(43),
									Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
								},
							}),
							expressions.NewIdentifier(EmptyListConstantName),
						},
					),
				},
			),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/list definition/error",
			args: args{
				code:                "[12, 23, unknown]",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Atom/hash table definition/success",
			args: args{
				code:                "{x: 12, y: 23, z: 42}",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				HashTableConstructionFunctionName,
				[]expressions.Expression{
					expressions.NewFunctionCall(HashTableConstructionFunctionName, []expressions.Expression{
						expressions.NewFunctionCall(HashTableConstructionFunctionName, []expressions.Expression{
							expressions.NewIdentifier(EmptyHashTableConstantName),
							expressions.NewString("x"),
							expressions.NewNumber(12),
						}),
						expressions.NewString("y"),
						expressions.NewNumber(23),
					}),
					expressions.NewString("z"),
					expressions.NewNumber(42),
				},
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/hash table definition/success/with setted states",
			args: args{
				code: `{
					x: when
						=> 23
							set one()
						=> 42
							set two()
					;,
					y: when
						=> 24
							set two()
						=> 43
							set three()
					;,
				}`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				HashTableConstructionFunctionName,
				[]expressions.Expression{
					expressions.NewFunctionCall(HashTableConstructionFunctionName, []expressions.Expression{
						expressions.NewIdentifier(EmptyHashTableConstantName),
						expressions.NewString("x"),
						expressions.NewConditionalExpression([]expressions.ConditionalCase{
							{
								Condition: expressions.NewNumber(23),
								Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
							},
							{
								Condition: expressions.NewNumber(42),
								Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
							},
						}),
					}),
					expressions.NewString("y"),
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(24),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
						{
							Condition: expressions.NewNumber(43),
							Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
						},
					}),
				},
			),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/hash table definition/error",
			args: args{
				code:                "{x: 12, y: 23, z: unknown}",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "Atom/function call/success",
			args: args{
				code:                "test(12, 23, 42)",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall("test", []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewNumber(23),
				expressions.NewNumber(42),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/function call/success/with setted states",
			args: args{
				code: `test(
					when
						=> 23
							set one()
						=> 42
							set two()
					;,
					when
						=> 24
							set two()
						=> 43
							set three()
					;,
				)`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall("test", []expressions.Expression{
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(23),
						Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
					},
					{
						Condition: expressions.NewNumber(42),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
				}),
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(24),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
					{
						Condition: expressions.NewNumber(43),
						Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
					},
				}),
			}),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/function call/error",
			args: args{
				code:                "test(12, 23, unknown)",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Atom/conditional expression/success",
			args: args{
				code: `
					when
						=> 12
							23
							42
						=> 13
							24
							43
						=> 14
							25
							44
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(12),
					Command: runtime.CommandGroup{
						commands.NewExpressionCommand(expressions.NewNumber(23)),
						commands.NewExpressionCommand(expressions.NewNumber(42)),
					},
				},
				{
					Condition: expressions.NewNumber(13),
					Command: runtime.CommandGroup{
						commands.NewExpressionCommand(expressions.NewNumber(24)),
						commands.NewExpressionCommand(expressions.NewNumber(43)),
					},
				},
				{
					Condition: expressions.NewNumber(14),
					Command: runtime.CommandGroup{
						commands.NewExpressionCommand(expressions.NewNumber(25)),
						commands.NewExpressionCommand(expressions.NewNumber(44)),
					},
				},
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/conditional expression/success/with setted states",
			args: args{
				code: `
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/conditional expression/error",
			args: args{
				code: `
					when
						=> 12
							23
							42
						=> 13
							24
							43
						=> 14
							25
							unknown
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Atom/identifier/success",
			args: args{
				code:                "test",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewIdentifier("test"),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/identifier/error",
			args: args{
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Atom/expression/success",
			args: args{
				code:                "(23)",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/expression/success/with setted states",
			args: args{
				code: `
					(when
						=> 23
							set one()
						=> 42
							set two()
					;)
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(23),
					Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
				},
				{
					Condition: expressions.NewNumber(42),
					Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/expression/error",
			args: args{
				code:                "(unknown)",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			atom := new(parser.Atom)
			err := parser.ParseToAST(data.args.code, atom)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateAtom(atom, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateListDefinition(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "ListDefinition/success/few items",
			args: args{
				code:                "[12, 23, 42]",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				ListConstructionFunctionName,
				[]expressions.Expression{
					expressions.NewNumber(12),
					expressions.NewFunctionCall(ListConstructionFunctionName, []expressions.Expression{
						expressions.NewNumber(23),
						expressions.NewFunctionCall(ListConstructionFunctionName, []expressions.Expression{
							expressions.NewNumber(42),
							expressions.NewIdentifier(EmptyListConstantName),
						}),
					}),
				},
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "ListDefinition/success/few items/with setted states",
			args: args{
				code: `[
					when
						=> 23
							set one()
						=> 42
							set two()
					;,
					when
						=> 24
							set two()
						=> 43
							set three()
					;,
				]`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				ListConstructionFunctionName,
				[]expressions.Expression{
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(23),
							Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
						},
						{
							Condition: expressions.NewNumber(42),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
					}),
					expressions.NewFunctionCall(ListConstructionFunctionName, []expressions.Expression{
						expressions.NewConditionalExpression([]expressions.ConditionalCase{
							{
								Condition: expressions.NewNumber(24),
								Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
							},
							{
								Condition: expressions.NewNumber(43),
								Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
							},
						}),
						expressions.NewIdentifier(EmptyListConstantName),
					}),
				},
			),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "ListDefinition/success/no items",
			args: args{
				code:                "[]",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewIdentifier(EmptyListConstantName),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "ListDefinition/error",
			args: args{
				code:                "[12, 23, unknown]",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			listDefinition := new(parser.ListDefinition)
			err := parser.ParseToAST(data.args.code, listDefinition)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateListDefinition(listDefinition, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateHashTableDefinition(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "HashTableDefinition/success/name/few entries/unknown key identifiers",
			args: args{
				code:                "{x: 12, y: 23, z: 42}",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				HashTableConstructionFunctionName,
				[]expressions.Expression{
					expressions.NewFunctionCall(HashTableConstructionFunctionName, []expressions.Expression{
						expressions.NewFunctionCall(HashTableConstructionFunctionName, []expressions.Expression{
							expressions.NewIdentifier(EmptyHashTableConstantName),
							expressions.NewString("x"),
							expressions.NewNumber(12),
						}),
						expressions.NewString("y"),
						expressions.NewNumber(23),
					}),
					expressions.NewString("z"),
					expressions.NewNumber(42),
				},
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "HashTableDefinition/success/name/few entries/known key identifiers",
			args: args{
				code:                "{x: 12, y: 23, z: 42}",
				declaredIdentifiers: mapset.NewSet("x", "y", "z"),
			},
			wantExpression: expressions.NewFunctionCall(
				HashTableConstructionFunctionName,
				[]expressions.Expression{
					expressions.NewFunctionCall(HashTableConstructionFunctionName, []expressions.Expression{
						expressions.NewFunctionCall(HashTableConstructionFunctionName, []expressions.Expression{
							expressions.NewIdentifier(EmptyHashTableConstantName),
							expressions.NewString("x"),
							expressions.NewNumber(12),
						}),
						expressions.NewString("y"),
						expressions.NewNumber(23),
					}),
					expressions.NewString("z"),
					expressions.NewNumber(42),
				},
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "HashTableDefinition/success/name/few entries/with setted states",
			args: args{
				code: `{
					x: when
						=> 23
							set one()
						=> 42
							set two()
					;,
					y: when
						=> 24
							set two()
						=> 43
							set three()
					;,
				}`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				HashTableConstructionFunctionName,
				[]expressions.Expression{
					expressions.NewFunctionCall(HashTableConstructionFunctionName, []expressions.Expression{
						expressions.NewIdentifier(EmptyHashTableConstantName),
						expressions.NewString("x"),
						expressions.NewConditionalExpression([]expressions.ConditionalCase{
							{
								Condition: expressions.NewNumber(23),
								Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
							},
							{
								Condition: expressions.NewNumber(42),
								Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
							},
						}),
					}),
					expressions.NewString("y"),
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(24),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
						{
							Condition: expressions.NewNumber(43),
							Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
						},
					}),
				},
			),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "HashTableDefinition/success/expression",
			args: args{
				code:                "{[test]: 23}",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				HashTableConstructionFunctionName,
				[]expressions.Expression{
					expressions.NewIdentifier(EmptyHashTableConstantName),
					expressions.NewIdentifier("test"),
					expressions.NewNumber(23),
				},
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "HashTableDefinition/success/expression/with setted states",
			args: args{
				code: `{
					[when
						=> 23
							set one()
						=> 42
							set two()
					;]: 23,
				}`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				HashTableConstructionFunctionName,
				[]expressions.Expression{
					expressions.NewIdentifier(EmptyHashTableConstantName),
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(23),
							Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
						},
						{
							Condition: expressions.NewNumber(42),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
					}),
					expressions.NewNumber(23),
				},
			),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "HashTableDefinition/success/no entries",
			args: args{
				code:                "{}",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewIdentifier(EmptyHashTableConstantName),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "HashTableDefinition/error/unknown identifier in the expression",
			args: args{
				code:                "{[unknown]: 23}",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "HashTableDefinition/error/unknown identifier in the value",
			args: args{
				code:                "{x: 12, y: 23, z: unknown}",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			hashTableDefinition := new(parser.HashTableDefinition)
			err := parser.ParseToAST(data.args.code, hashTableDefinition)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateHashTableDefinition(hashTableDefinition, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateFunctionCall(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "FunctionCall/success/few arguments",
			args: args{
				code:                "test(12, 23, 42)",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall("test", []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewNumber(23),
				expressions.NewNumber(42),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "FunctionCall/success/few arguments/with setted states",
			args: args{
				code: `test(
					when
						=> 23
							set one()
						=> 42
							set two()
					;,
					when
						=> 24
							set two()
						=> 43
							set three()
					;,
				)`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall("test", []expressions.Expression{
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(23),
						Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
					},
					{
						Condition: expressions.NewNumber(42),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
				}),
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(24),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
					},
					{
						Condition: expressions.NewNumber(43),
						Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
					},
				}),
			}),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "FunctionCall/success/no arguments",
			args: args{
				code:                "test()",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewFunctionCall("test", nil),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "FunctionCall/error/unknown function",
			args: args{
				code:                "unknown(12, 23, 42)",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "FunctionCall/error/argument translating",
			args: args{
				code:                "test(12, 23, unknown)",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			functionCall := new(parser.FunctionCall)
			err := parser.ParseToAST(data.args.code, functionCall)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateFunctionCall(functionCall, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateConditionalExpression(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, data := range []struct {
		name             string
		args             args
		wantExpression   expressions.Expression
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "ConditionalExpression/success/single conditional case/nonempty",
			args: args{
				code: `
					when
						=> 12
							23
							42
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(12),
					Command: runtime.CommandGroup{
						commands.NewExpressionCommand(expressions.NewNumber(23)),
						commands.NewExpressionCommand(expressions.NewNumber(42)),
					},
				},
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "ConditionalExpression/success/single conditional case/empty",
			args: args{
				code: `
					when
						=> 12
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(12),
					Command:   runtime.CommandGroup(nil),
				},
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "ConditionalExpression/success/few conditional cases/nonempty",
			args: args{
				code: `
					when
						=> 12
							23
							42
						=> 13
							24
							43
						=> 14
							25
							44
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(12),
					Command: runtime.CommandGroup{
						commands.NewExpressionCommand(expressions.NewNumber(23)),
						commands.NewExpressionCommand(expressions.NewNumber(42)),
					},
				},
				{
					Condition: expressions.NewNumber(13),
					Command: runtime.CommandGroup{
						commands.NewExpressionCommand(expressions.NewNumber(24)),
						commands.NewExpressionCommand(expressions.NewNumber(43)),
					},
				},
				{
					Condition: expressions.NewNumber(14),
					Command: runtime.CommandGroup{
						commands.NewExpressionCommand(expressions.NewNumber(25)),
						commands.NewExpressionCommand(expressions.NewNumber(44)),
					},
				},
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "ConditionalExpression/success/few conditional cases/empty",
			args: args{
				code: `
					when
						=> 12
						=> 13
						=> 14
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewNumber(12),
					Command:   runtime.CommandGroup(nil),
				},
				{
					Condition: expressions.NewNumber(13),
					Command:   runtime.CommandGroup(nil),
				},
				{
					Condition: expressions.NewNumber(14),
					Command:   runtime.CommandGroup(nil),
				},
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "ConditionalExpression/success/without conditional cases",
			args: args{
				code:                "when;",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewConditionalExpression(nil),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "ConditionalExpression/success/nonempty/with setted states",
			args: args{
				code: `
					when
						=> when
							=> 23
								set one()
							=> 42
								set two()
						;
							when
								=> 24
									set two()
								=> 43
									set three()
							;
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(23),
							Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
						},
						{
							Condition: expressions.NewNumber(42),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
					}),
					Command: runtime.CommandGroup{
						commands.NewExpressionCommand(
							expressions.NewConditionalExpression([]expressions.ConditionalCase{
								{
									Condition: expressions.NewNumber(24),
									Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
								},
								{
									Condition: expressions.NewNumber(43),
									Command:   runtime.CommandGroup{commands.NewSetCommand("three", nil)},
								},
							}),
						),
					},
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "ConditionalExpression/success/empty/with setted states",
			args: args{
				code: `
					when
						=> when
							=> 23
								set one()
							=> 42
								set two()
						;
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewConditionalExpression([]expressions.ConditionalCase{
				{
					Condition: expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(23),
							Command:   runtime.CommandGroup{commands.NewSetCommand("one", nil)},
						},
						{
							Condition: expressions.NewNumber(42),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two", nil)},
						},
					}),
					Command: runtime.CommandGroup(nil),
				},
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "ConditionalExpression/error/condition translating",
			args: args{
				code: `
					when
						=> 12
							23
							42
						=> 13
							24
							43
						=> unknown
							25
							44
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "ConditionalExpression/error/command translating",
			args: args{
				code: `
					when
						=> 12
							23
							42
						=> 13
							24
							43
						=> 14
							25
							unknown
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			conditionalExpression := new(parser.ConditionalExpression)
			err := parser.ParseToAST(data.args.code, conditionalExpression)
			require.NoError(test, err)

			gotExpression, gotSettedStates, gotErr :=
				translateConditionalExpression(conditionalExpression, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}
