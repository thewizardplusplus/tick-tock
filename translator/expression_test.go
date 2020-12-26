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
				translateExpression(expression, data.args.declaredIdentifiers)

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

func TestTranslateEquality(test *testing.T) {
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
				translateEquality(equality, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateComparison(test *testing.T) {
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
				translateComparison(comparison, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateBitwiseDisjunction(test *testing.T) {
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
				translateBitwiseDisjunction(bitwiseDisjunction, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateBitwiseExclusiveDisjunction(test *testing.T) {
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

			gotExpression, gotSettedStates, gotErr := translateBitwiseExclusiveDisjunction(
				bitwiseExclusiveDisjunction,
				data.args.declaredIdentifiers,
			)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateBitwiseConjunction(test *testing.T) {
	type args struct {
		bitwiseConjunction  *parser.BitwiseConjunction
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
				bitwiseConjunction: func() *parser.BitwiseConjunction {
					bitwiseConjunction := new(parser.BitwiseConjunction)
					err := parser.ParseToAST("12 & 23 & 42", bitwiseConjunction)
					require.NoError(test, err)

					return bitwiseConjunction
				}(),
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
				bitwiseConjunction: func() *parser.BitwiseConjunction {
					const code = `
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
					`

					bitwiseConjunction := new(parser.BitwiseConjunction)
					err := parser.ParseToAST(code, bitwiseConjunction)
					require.NoError(test, err)

					return bitwiseConjunction
				}(),
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
				bitwiseConjunction: func() *parser.BitwiseConjunction {
					bitwiseConjunction := new(parser.BitwiseConjunction)
					err := parser.ParseToAST("12 & 23 & unknown", bitwiseConjunction)
					require.NoError(test, err)

					return bitwiseConjunction
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "BitwiseConjunction/empty/success",
			args: args{
				bitwiseConjunction: func() *parser.BitwiseConjunction {
					bitwiseConjunction := new(parser.BitwiseConjunction)
					err := parser.ParseToAST("23", bitwiseConjunction)
					require.NoError(test, err)

					return bitwiseConjunction
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseConjunction/empty/success/with setted states",
			args: args{
				bitwiseConjunction: func() *parser.BitwiseConjunction {
					const code = `
						when
							=> 23
								set one()
							=> 42
								set two()
						;
					`

					bitwiseConjunction := new(parser.BitwiseConjunction)
					err := parser.ParseToAST(code, bitwiseConjunction)
					require.NoError(test, err)

					return bitwiseConjunction
				}(),
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
				bitwiseConjunction: func() *parser.BitwiseConjunction {
					bitwiseConjunction := new(parser.BitwiseConjunction)
					err := parser.ParseToAST("unknown", bitwiseConjunction)
					require.NoError(test, err)

					return bitwiseConjunction
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateBitwiseConjunction(data.args.bitwiseConjunction, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateShift(test *testing.T) {
	type args struct {
		shift               *parser.Shift
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
				shift: func() *parser.Shift {
					shift := new(parser.Shift)
					err := parser.ParseToAST("5 << 12 >> 23 >>> 42", shift)
					require.NoError(test, err)

					return shift
				}(),
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
				shift: func() *parser.Shift {
					const code = `
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
					`

					shift := new(parser.Shift)
					err := parser.ParseToAST(code, shift)
					require.NoError(test, err)

					return shift
				}(),
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
				shift: func() *parser.Shift {
					shift := new(parser.Shift)
					err := parser.ParseToAST("5 << 12 >> 23 >>> unknown", shift)
					require.NoError(test, err)

					return shift
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "Shift/empty/success",
			args: args{
				shift: func() *parser.Shift {
					shift := new(parser.Shift)
					err := parser.ParseToAST("23", shift)
					require.NoError(test, err)

					return shift
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Shift/empty/success/with setted states",
			args: args{
				shift: func() *parser.Shift {
					const code = `
						when
							=> 23
								set one()
							=> 42
								set two()
						;
					`

					shift := new(parser.Shift)
					err := parser.ParseToAST(code, shift)
					require.NoError(test, err)

					return shift
				}(),
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
				shift: func() *parser.Shift {
					shift := new(parser.Shift)
					err := parser.ParseToAST("unknown", shift)
					require.NoError(test, err)

					return shift
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateShift(data.args.shift, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateAddition(test *testing.T) {
	type args struct {
		addition            *parser.Addition
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
				addition: func() *parser.Addition {
					addition := new(parser.Addition)
					err := parser.ParseToAST("12 + 23 + 42", addition)
					require.NoError(test, err)

					return addition
				}(),
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
				addition: func() *parser.Addition {
					addition := new(parser.Addition)
					err := parser.ParseToAST("12 - 23 - 42", addition)
					require.NoError(test, err)

					return addition
				}(),
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
				addition: func() *parser.Addition {
					const code = `
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
					`

					addition := new(parser.Addition)
					err := parser.ParseToAST(code, addition)
					require.NoError(test, err)

					return addition
				}(),
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
				addition: func() *parser.Addition {
					addition := new(parser.Addition)
					err := parser.ParseToAST("12 + 23 + unknown", addition)
					require.NoError(test, err)

					return addition
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Addition/empty/success",
			args: args{
				addition: func() *parser.Addition {
					addition := new(parser.Addition)
					err := parser.ParseToAST("23", addition)
					require.NoError(test, err)

					return addition
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Addition/empty/success/with setted states",
			args: args{
				addition: func() *parser.Addition {
					const code = `
						when
							=> 23
								set one()
							=> 42
								set two()
						;
					`

					addition := new(parser.Addition)
					err := parser.ParseToAST(code, addition)
					require.NoError(test, err)

					return addition
				}(),
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
				addition: func() *parser.Addition {
					addition := new(parser.Addition)
					err := parser.ParseToAST("unknown", addition)
					require.NoError(test, err)

					return addition
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateAddition(data.args.addition, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateMultiplication(test *testing.T) {
	type args struct {
		multiplication      *parser.Multiplication
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
				multiplication: func() *parser.Multiplication {
					multiplication := new(parser.Multiplication)
					err := parser.ParseToAST("12 * 23 * 42", multiplication)
					require.NoError(test, err)

					return multiplication
				}(),
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
				multiplication: func() *parser.Multiplication {
					multiplication := new(parser.Multiplication)
					err := parser.ParseToAST("12 / 23 / 42", multiplication)
					require.NoError(test, err)

					return multiplication
				}(),
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
				multiplication: func() *parser.Multiplication {
					multiplication := new(parser.Multiplication)
					err := parser.ParseToAST("12 % 23 % 42", multiplication)
					require.NoError(test, err)

					return multiplication
				}(),
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
				multiplication: func() *parser.Multiplication {
					const code = `
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
					`

					multiplication := new(parser.Multiplication)
					err := parser.ParseToAST(code, multiplication)
					require.NoError(test, err)

					return multiplication
				}(),
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
				multiplication: func() *parser.Multiplication {
					multiplication := new(parser.Multiplication)
					err := parser.ParseToAST("12 * 23 * unknown", multiplication)
					require.NoError(test, err)

					return multiplication
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Multiplication/empty/success",
			args: args{
				multiplication: func() *parser.Multiplication {
					multiplication := new(parser.Multiplication)
					err := parser.ParseToAST("23", multiplication)
					require.NoError(test, err)

					return multiplication
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Multiplication/empty/success/with setted states",
			args: args{
				multiplication: func() *parser.Multiplication {
					const code = `
						when
							=> 23
								set one()
							=> 42
								set two()
						;
					`

					multiplication := new(parser.Multiplication)
					err := parser.ParseToAST(code, multiplication)
					require.NoError(test, err)

					return multiplication
				}(),
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
				multiplication: func() *parser.Multiplication {
					multiplication := new(parser.Multiplication)
					err := parser.ParseToAST("unknown", multiplication)
					require.NoError(test, err)

					return multiplication
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateMultiplication(data.args.multiplication, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateUnary(test *testing.T) {
	type args struct {
		unary               *parser.Unary
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
				unary: func() *parser.Unary {
					unary := new(parser.Unary)
					err := parser.ParseToAST("-~!23", unary)
					require.NoError(test, err)

					return unary
				}(),
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
				unary: func() *parser.Unary {
					const code = `
						-~!when
							=> 23
								set one()
							=> 42
								set two()
						;
					`

					unary := new(parser.Unary)
					err := parser.ParseToAST(code, unary)
					require.NoError(test, err)

					return unary
				}(),
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
				unary: func() *parser.Unary {
					unary := new(parser.Unary)
					err := parser.ParseToAST("-~!unknown", unary)
					require.NoError(test, err)

					return unary
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Unary/empty/success",
			args: args{
				unary: func() *parser.Unary {
					unary := new(parser.Unary)
					err := parser.ParseToAST("23", unary)
					require.NoError(test, err)

					return unary
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Unary/empty/success/with setted states",
			args: args{
				unary: func() *parser.Unary {
					const code = `
						when
							=> 23
								set one()
							=> 42
								set two()
						;
					`

					unary := new(parser.Unary)
					err := parser.ParseToAST(code, unary)
					require.NoError(test, err)

					return unary
				}(),
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
				unary: func() *parser.Unary {
					unary := new(parser.Unary)
					err := parser.ParseToAST("unknown", unary)
					require.NoError(test, err)

					return unary
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateUnary(data.args.unary, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateAccessor(test *testing.T) {
	type args struct {
		accessor            *parser.Accessor
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
				accessor: func() *parser.Accessor {
					accessor := new(parser.Accessor)
					err := parser.ParseToAST("test[12][23]", accessor)
					require.NoError(test, err)

					return accessor
				}(),
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
				accessor: func() *parser.Accessor {
					accessor := new(parser.Accessor)
					err := parser.ParseToAST("test.one.two", accessor)
					require.NoError(test, err)

					return accessor
				}(),
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
				accessor: func() *parser.Accessor {
					accessor := new(parser.Accessor)
					err := parser.ParseToAST("test.one[12].two[23]", accessor)
					require.NoError(test, err)

					return accessor
				}(),
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
				accessor: func() *parser.Accessor {
					const code = `
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
					`

					accessor := new(parser.Accessor)
					err := parser.ParseToAST(code, accessor)
					require.NoError(test, err)

					return accessor
				}(),
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
				accessor: func() *parser.Accessor {
					accessor := new(parser.Accessor)
					err := parser.ParseToAST("unknown[12][23]", accessor)
					require.NoError(test, err)

					return accessor
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Accessor/nonempty/error/key translating",
			args: args{
				accessor: func() *parser.Accessor {
					accessor := new(parser.Accessor)
					err := parser.ParseToAST("test[12][unknown]", accessor)
					require.NoError(test, err)

					return accessor
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Accessor/empty/success",
			args: args{
				accessor: func() *parser.Accessor {
					accessor := new(parser.Accessor)
					err := parser.ParseToAST("23", accessor)
					require.NoError(test, err)

					return accessor
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Accessor/empty/success/with setted states",
			args: args{
				accessor: func() *parser.Accessor {
					const code = `
						when
							=> 23
								set one()
							=> 42
								set two()
						;
					`

					accessor := new(parser.Accessor)
					err := parser.ParseToAST(code, accessor)
					require.NoError(test, err)

					return accessor
				}(),
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
				accessor: func() *parser.Accessor {
					accessor := new(parser.Accessor)
					err := parser.ParseToAST("unknown", accessor)
					require.NoError(test, err)

					return accessor
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateAccessor(data.args.accessor, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateAtom(test *testing.T) {
	type args struct {
		atom                *parser.Atom
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
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST("23", atom)
					require.NoError(test, err)

					return atom
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/number/floating-point",
			args: args{
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST("2.3", atom)
					require.NoError(test, err)

					return atom
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(2.3),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/symbol/latin1",
			args: args{
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST("'t'", atom)
					require.NoError(test, err)

					return atom
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(116),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/symbol/not latin1",
			args: args{
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST("'т'", atom)
					require.NoError(test, err)

					return atom
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(1090),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/string",
			args: args{
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST(`"test"`, atom)
					require.NoError(test, err)

					return atom
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewString("test"),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/list definition/success",
			args: args{
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST("[12, 23, 42]", atom)
					require.NoError(test, err)

					return atom
				}(),
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
				atom: func() *parser.Atom {
					const code = `[
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
					]`

					atom := new(parser.Atom)
					err := parser.ParseToAST(code, atom)
					require.NoError(test, err)

					return atom
				}(),
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
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST("[12, 23, unknown]", atom)
					require.NoError(test, err)

					return atom
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Atom/hash table definition/success",
			args: args{
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST("{x: 12, y: 23, z: 42}", atom)
					require.NoError(test, err)

					return atom
				}(),
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
				atom: func() *parser.Atom {
					const code = `{
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
					}`

					atom := new(parser.Atom)
					err := parser.ParseToAST(code, atom)
					require.NoError(test, err)

					return atom
				}(),
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
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST("{x: 12, y: 23, z: unknown}", atom)
					require.NoError(test, err)

					return atom
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "Atom/function call/success",
			args: args{
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST("test(12, 23, 42)", atom)
					require.NoError(test, err)

					return atom
				}(),
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
				atom: func() *parser.Atom {
					const code = `test(
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
					)`

					atom := new(parser.Atom)
					err := parser.ParseToAST(code, atom)
					require.NoError(test, err)

					return atom
				}(),
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
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST("test(12, 23, unknown)", atom)
					require.NoError(test, err)

					return atom
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Atom/conditional expression/success",
			args: args{
				atom: func() *parser.Atom {
					const code = `
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
					`

					atom := new(parser.Atom)
					err := parser.ParseToAST(code, atom)
					require.NoError(test, err)

					return atom
				}(),
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
				atom: func() *parser.Atom {
					const code = `
						when
							=> 23
								set one()
							=> 42
								set two()
						;
					`

					atom := new(parser.Atom)
					err := parser.ParseToAST(code, atom)
					require.NoError(test, err)

					return atom
				}(),
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
				atom: func() *parser.Atom {
					const code = `
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
					`

					atom := new(parser.Atom)
					err := parser.ParseToAST(code, atom)
					require.NoError(test, err)

					return atom
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Atom/identifier/success",
			args: args{
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST("test", atom)
					require.NoError(test, err)

					return atom
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewIdentifier("test"),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/identifier/error",
			args: args{
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST("unknown", atom)
					require.NoError(test, err)

					return atom
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Atom/expression/success",
			args: args{
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST("(23)", atom)
					require.NoError(test, err)

					return atom
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/expression/success/with setted states",
			args: args{
				atom: func() *parser.Atom {
					const code = `
						(when
							=> 23
								set one()
							=> 42
								set two()
						;)
					`

					atom := new(parser.Atom)
					err := parser.ParseToAST(code, atom)
					require.NoError(test, err)

					return atom
				}(),
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
				atom: func() *parser.Atom {
					atom := new(parser.Atom)
					err := parser.ParseToAST("(unknown)", atom)
					require.NoError(test, err)

					return atom
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateAtom(data.args.atom, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateListDefinition(test *testing.T) {
	type args struct {
		listDefinition      *parser.ListDefinition
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
				listDefinition: func() *parser.ListDefinition {
					listDefinition := new(parser.ListDefinition)
					err := parser.ParseToAST("[12, 23, 42]", listDefinition)
					require.NoError(test, err)

					return listDefinition
				}(),
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
				listDefinition: func() *parser.ListDefinition {
					const code = `[
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
					]`

					listDefinition := new(parser.ListDefinition)
					err := parser.ParseToAST(code, listDefinition)
					require.NoError(test, err)

					return listDefinition
				}(),
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
				listDefinition: func() *parser.ListDefinition {
					listDefinition := new(parser.ListDefinition)
					err := parser.ParseToAST("[]", listDefinition)
					require.NoError(test, err)

					return listDefinition
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewIdentifier(EmptyListConstantName),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "ListDefinition/error",
			args: args{
				listDefinition: func() *parser.ListDefinition {
					listDefinition := new(parser.ListDefinition)
					err := parser.ParseToAST("[12, 23, unknown]", listDefinition)
					require.NoError(test, err)

					return listDefinition
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateListDefinition(data.args.listDefinition, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateHashTableDefinition(test *testing.T) {
	type args struct {
		hashTableDefinition *parser.HashTableDefinition
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
				hashTableDefinition: func() *parser.HashTableDefinition {
					hashTableDefinition := new(parser.HashTableDefinition)
					err := parser.ParseToAST("{x: 12, y: 23, z: 42}", hashTableDefinition)
					require.NoError(test, err)

					return hashTableDefinition
				}(),
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
				hashTableDefinition: func() *parser.HashTableDefinition {
					hashTableDefinition := new(parser.HashTableDefinition)
					err := parser.ParseToAST("{x: 12, y: 23, z: 42}", hashTableDefinition)
					require.NoError(test, err)

					return hashTableDefinition
				}(),
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
				hashTableDefinition: func() *parser.HashTableDefinition {
					const code = `{
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
					}`

					hashTableDefinition := new(parser.HashTableDefinition)
					err := parser.ParseToAST(code, hashTableDefinition)
					require.NoError(test, err)

					return hashTableDefinition
				}(),
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
				hashTableDefinition: func() *parser.HashTableDefinition {
					hashTableDefinition := new(parser.HashTableDefinition)
					err := parser.ParseToAST("{[test]: 23}", hashTableDefinition)
					require.NoError(test, err)

					return hashTableDefinition
				}(),
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
				hashTableDefinition: func() *parser.HashTableDefinition {
					const code = `{
						[when
							=> 23
								set one()
							=> 42
								set two()
						;]: 23,
					}`

					hashTableDefinition := new(parser.HashTableDefinition)
					err := parser.ParseToAST(code, hashTableDefinition)
					require.NoError(test, err)

					return hashTableDefinition
				}(),
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
				hashTableDefinition: func() *parser.HashTableDefinition {
					hashTableDefinition := new(parser.HashTableDefinition)
					err := parser.ParseToAST("{}", hashTableDefinition)
					require.NoError(test, err)

					return hashTableDefinition
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewIdentifier(EmptyHashTableConstantName),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "HashTableDefinition/error/unknown identifier in the expression",
			args: args{
				hashTableDefinition: func() *parser.HashTableDefinition {
					hashTableDefinition := new(parser.HashTableDefinition)
					err := parser.ParseToAST("{[unknown]: 23}", hashTableDefinition)
					require.NoError(test, err)

					return hashTableDefinition
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "HashTableDefinition/error/unknown identifier in the value",
			args: args{
				hashTableDefinition: func() *parser.HashTableDefinition {
					hashTableDefinition := new(parser.HashTableDefinition)
					err := parser.ParseToAST("{x: 12, y: 23, z: unknown}", hashTableDefinition)
					require.NoError(test, err)

					return hashTableDefinition
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateHashTableDefinition(data.args.hashTableDefinition, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateFunctionCall(test *testing.T) {
	type args struct {
		functionCall        *parser.FunctionCall
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
				functionCall: func() *parser.FunctionCall {
					functionCall := new(parser.FunctionCall)
					err := parser.ParseToAST("test(12, 23, 42)", functionCall)
					require.NoError(test, err)

					return functionCall
				}(),
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
				functionCall: func() *parser.FunctionCall {
					const code = `test(
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
					)`

					functionCall := new(parser.FunctionCall)
					err := parser.ParseToAST(code, functionCall)
					require.NoError(test, err)

					return functionCall
				}(),
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
				functionCall: func() *parser.FunctionCall {
					functionCall := new(parser.FunctionCall)
					err := parser.ParseToAST("test()", functionCall)
					require.NoError(test, err)

					return functionCall
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewFunctionCall("test", nil),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "FunctionCall/error/unknown function",
			args: args{
				functionCall: func() *parser.FunctionCall {
					functionCall := new(parser.FunctionCall)
					err := parser.ParseToAST("unknown(12, 23, 42)", functionCall)
					require.NoError(test, err)

					return functionCall
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "FunctionCall/error/argument translating",
			args: args{
				functionCall: func() *parser.FunctionCall {
					functionCall := new(parser.FunctionCall)
					err := parser.ParseToAST("test(12, 23, unknown)", functionCall)
					require.NoError(test, err)

					return functionCall
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateFunctionCall(data.args.functionCall, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateConditionalExpression(test *testing.T) {
	type args struct {
		conditionalExpression *parser.ConditionalExpression
		declaredIdentifiers   mapset.Set
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
				conditionalExpression: func() *parser.ConditionalExpression {
					const code = `
						when
							=> 12
								23
								42
						;
					`

					conditionalExpression := new(parser.ConditionalExpression)
					err := parser.ParseToAST(code, conditionalExpression)
					require.NoError(test, err)

					return conditionalExpression
				}(),
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
				conditionalExpression: func() *parser.ConditionalExpression {
					const code = `
						when
							=> 12
						;
					`

					conditionalExpression := new(parser.ConditionalExpression)
					err := parser.ParseToAST(code, conditionalExpression)
					require.NoError(test, err)

					return conditionalExpression
				}(),
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
				conditionalExpression: func() *parser.ConditionalExpression {
					const code = `
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
					`

					conditionalExpression := new(parser.ConditionalExpression)
					err := parser.ParseToAST(code, conditionalExpression)
					require.NoError(test, err)

					return conditionalExpression
				}(),
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
				conditionalExpression: func() *parser.ConditionalExpression {
					const code = `
						when
							=> 12
							=> 13
							=> 14
						;
					`

					conditionalExpression := new(parser.ConditionalExpression)
					err := parser.ParseToAST(code, conditionalExpression)
					require.NoError(test, err)

					return conditionalExpression
				}(),
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
				conditionalExpression: func() *parser.ConditionalExpression {
					conditionalExpression := new(parser.ConditionalExpression)
					err := parser.ParseToAST("when;", conditionalExpression)
					require.NoError(test, err)

					return conditionalExpression
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewConditionalExpression(nil),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "ConditionalExpression/success/nonempty/with setted states",
			args: args{
				conditionalExpression: func() *parser.ConditionalExpression {
					const code = `
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
					`

					conditionalExpression := new(parser.ConditionalExpression)
					err := parser.ParseToAST(code, conditionalExpression)
					require.NoError(test, err)

					return conditionalExpression
				}(),
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
				conditionalExpression: func() *parser.ConditionalExpression {
					const code = `
						when
							=> when
								=> 23
									set one()
								=> 42
									set two()
							;
						;
					`

					conditionalExpression := new(parser.ConditionalExpression)
					err := parser.ParseToAST(code, conditionalExpression)
					require.NoError(test, err)

					return conditionalExpression
				}(),
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
				conditionalExpression: func() *parser.ConditionalExpression {
					const code = `
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
					`

					conditionalExpression := new(parser.ConditionalExpression)
					err := parser.ParseToAST(code, conditionalExpression)
					require.NoError(test, err)

					return conditionalExpression
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "ConditionalExpression/error/command translating",
			args: args{
				conditionalExpression: func() *parser.ConditionalExpression {
					const code = `
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
					`

					conditionalExpression := new(parser.ConditionalExpression)
					err := parser.ParseToAST(code, conditionalExpression)
					require.NoError(test, err)

					return conditionalExpression
				}(),
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateConditionalExpression(data.args.conditionalExpression, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}
