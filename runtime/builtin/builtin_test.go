package builtin

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			name:       "inf",
			expression: expressions.NewIdentifier("inf"),
			wantResult: math.Inf(+1),
			wantErr:    assert.NoError,
		},
		{
			name:       "pi",
			expression: expressions.NewIdentifier("pi"),
			wantResult: math.Pi,
			wantErr:    assert.NoError,
		},
		{
			name:       "e",
			expression: expressions.NewIdentifier("e"),
			wantResult: math.E,
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
		{
			name: "type/success/float64",
			expression: expressions.NewFunctionCall("type", []expressions.Expression{
				expressions.NewNumber(23),
			}),
			wantResult: types.NewPairFromText("num"),
			wantErr:    assert.NoError,
		},
		{
			name: "type/success/*types.Pair",
			expression: expressions.NewFunctionCall("type", []expressions.Expression{
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
			}),
			wantResult: types.NewPairFromText("list"),
			wantErr:    assert.NoError,
		},
		{
			name: "type/error",
			expression: expressions.NewFunctionCall("type", []expressions.Expression{
				expressions.NewIdentifier("type"),
			}),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "size",
			expression: expressions.NewFunctionCall("size", []expressions.Expression{
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
			}),
			wantResult: 3.0,
			wantErr:    assert.NoError,
		},
		{
			name: "floor",
			expression: expressions.NewFunctionCall("floor", []expressions.Expression{
				expressions.NewNumber(2.5),
			}),
			wantResult: 2.0,
			wantErr:    assert.NoError,
		},
		{
			name: "ceil",
			expression: expressions.NewFunctionCall("ceil", []expressions.Expression{
				expressions.NewNumber(2.5),
			}),
			wantResult: 3.0,
			wantErr:    assert.NoError,
		},
		{
			name: "trunc",
			expression: expressions.NewFunctionCall("trunc", []expressions.Expression{
				expressions.NewNumber(2.5),
			}),
			wantResult: 2.0,
			wantErr:    assert.NoError,
		},
		{
			name: "round",
			expression: expressions.NewFunctionCall("round", []expressions.Expression{
				expressions.NewNumber(2.5),
			}),
			wantResult: 3.0,
			wantErr:    assert.NoError,
		},
		{
			name: "pow",
			expression: expressions.NewFunctionCall("pow", []expressions.Expression{
				expressions.NewNumber(2),
				expressions.NewNumber(3),
			}),
			wantResult: 8.0,
			wantErr:    assert.NoError,
		},
		{
			name: "sqrt",
			expression: expressions.NewFunctionCall("sqrt", []expressions.Expression{
				expressions.NewNumber(4),
			}),
			wantResult: 2.0,
			wantErr:    assert.NoError,
		},
		{
			name: "abs",
			expression: expressions.NewFunctionCall("abs", []expressions.Expression{
				expressions.NewNumber(-23),
			}),
			wantResult: 23.0,
			wantErr:    assert.NoError,
		},
		{
			name: "head/success",
			expression: expressions.NewFunctionCall("head", []expressions.Expression{
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
			}),
			wantResult: 12.0,
			wantErr:    assert.NoError,
		},
		{
			name: "head/error",
			expression: expressions.NewFunctionCall("head", []expressions.Expression{
				expressions.NewIdentifier(translator.EmptyListConstantName),
			}),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "tail/success/nonempty tail",
			expression: expressions.NewFunctionCall("tail", []expressions.Expression{
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
			}),
			wantResult: &types.Pair{
				Head: 23.0,
				Tail: &types.Pair{
					Head: 42.0,
					Tail: nil,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "tail/success/empty tail",
			expression: expressions.NewFunctionCall("tail", []expressions.Expression{
				expressions.NewFunctionCall(
					translator.ListConstructionFunctionName,
					[]expressions.Expression{
						expressions.NewNumber(23),
						expressions.NewIdentifier(translator.EmptyListConstantName),
					},
				),
			}),
			wantResult: (*types.Pair)(nil),
			wantErr:    assert.NoError,
		},
		{
			name: "tail/error",
			expression: expressions.NewFunctionCall("tail", []expressions.Expression{
				expressions.NewIdentifier(translator.EmptyListConstantName),
			}),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "num/success",
			expression: expressions.NewFunctionCall("num", []expressions.Expression{
				expressions.NewString("23"),
			}),
			wantResult: 23.0,
			wantErr:    assert.NoError,
		},
		{
			name: "num/error/list conversion",
			expression: expressions.NewFunctionCall("num", []expressions.Expression{
				expressions.NewFunctionCall(
					translator.ListConstructionFunctionName,
					[]expressions.Expression{
						expressions.NewNumber(float64('t')),
						expressions.NewFunctionCall(
							translator.ListConstructionFunctionName,
							[]expressions.Expression{
								expressions.NewFunctionCall(
									translator.ListConstructionFunctionName,
									[]expressions.Expression{
										expressions.NewNumber(float64('h')),
										expressions.NewFunctionCall(
											translator.ListConstructionFunctionName,
											[]expressions.Expression{
												expressions.NewNumber(float64('i')),
												expressions.NewIdentifier(translator.EmptyListConstantName),
											},
										),
									},
								),
								expressions.NewFunctionCall(
									translator.ListConstructionFunctionName,
									[]expressions.Expression{
										expressions.NewNumber(float64('s')),
										expressions.NewFunctionCall(
											translator.ListConstructionFunctionName,
											[]expressions.Expression{
												expressions.NewNumber(float64('t')),
												expressions.NewIdentifier(translator.EmptyListConstantName),
											},
										),
									},
								),
							},
						),
					},
				),
			}),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "num/error/string conversion",
			expression: expressions.NewFunctionCall("num", []expressions.Expression{
				expressions.NewString("test"),
			}),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "str/success/float64",
			expression: expressions.NewFunctionCall("str", []expressions.Expression{
				expressions.NewNumber(23),
			}),
			wantResult: types.NewPairFromText("23"),
			wantErr:    assert.NoError,
		},
		{
			name: "str/success/*types.Pair/tree in the head",
			expression: expressions.NewFunctionCall("str", []expressions.Expression{
				expressions.NewFunctionCall(
					translator.ListConstructionFunctionName,
					[]expressions.Expression{
						expressions.NewFunctionCall(
							translator.ListConstructionFunctionName,
							[]expressions.Expression{
								expressions.NewNumber(float64('h')),
								expressions.NewFunctionCall(
									translator.ListConstructionFunctionName,
									[]expressions.Expression{
										expressions.NewNumber(float64('i')),
										expressions.NewIdentifier(translator.EmptyListConstantName),
									},
								),
							},
						),
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
			}),
			wantResult: types.NewPairFromText("[[104,105],23,42]"),
			wantErr:    assert.NoError,
		},
		{
			name: "str/success/*types.Pair/tree in the tail",
			expression: expressions.NewFunctionCall("str", []expressions.Expression{
				expressions.NewFunctionCall(
					translator.ListConstructionFunctionName,
					[]expressions.Expression{
						expressions.NewNumber(12),
						expressions.NewFunctionCall(
							translator.ListConstructionFunctionName,
							[]expressions.Expression{
								expressions.NewFunctionCall(
									translator.ListConstructionFunctionName,
									[]expressions.Expression{
										expressions.NewNumber(float64('h')),
										expressions.NewFunctionCall(
											translator.ListConstructionFunctionName,
											[]expressions.Expression{
												expressions.NewNumber(float64('i')),
												expressions.NewIdentifier(translator.EmptyListConstantName),
											},
										),
									},
								),
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
			}),
			wantResult: types.NewPairFromText("[12,[104,105],42]"),
			wantErr:    assert.NoError,
		},
		{
			name: "str/error/JSON marshalling",
			expression: expressions.NewFunctionCall("str", []expressions.Expression{
				expressions.NewFunctionCall(
					translator.ListConstructionFunctionName,
					[]expressions.Expression{
						expressions.NewNumber(12),
						expressions.NewFunctionCall(
							translator.ListConstructionFunctionName,
							[]expressions.Expression{
								expressions.NewIdentifier("str"),
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
			}),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "str/error/unsupported type",
			expression: expressions.NewFunctionCall("str", []expressions.Expression{
				expressions.NewIdentifier("str"),
			}),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "strs/success",
			expression: expressions.NewFunctionCall("strs", []expressions.Expression{
				expressions.NewString(`"test"`),
			}),
			wantResult: types.NewPairFromText(`"\"test\""`),
			wantErr:    assert.NoError,
		},
		{
			name: "strs/error",
			expression: expressions.NewFunctionCall("strs", []expressions.Expression{
				expressions.NewFunctionCall(
					translator.ListConstructionFunctionName,
					[]expressions.Expression{
						expressions.NewNumber(float64('t')),
						expressions.NewFunctionCall(
							translator.ListConstructionFunctionName,
							[]expressions.Expression{
								expressions.NewFunctionCall(
									translator.ListConstructionFunctionName,
									[]expressions.Expression{
										expressions.NewNumber(float64('h')),
										expressions.NewFunctionCall(
											translator.ListConstructionFunctionName,
											[]expressions.Expression{
												expressions.NewNumber(float64('i')),
												expressions.NewIdentifier(translator.EmptyListConstantName),
											},
										),
									},
								),
								expressions.NewFunctionCall(
									translator.ListConstructionFunctionName,
									[]expressions.Expression{
										expressions.NewNumber(float64('s')),
										expressions.NewFunctionCall(
											translator.ListConstructionFunctionName,
											[]expressions.Expression{
												expressions.NewNumber(float64('t')),
												expressions.NewIdentifier(translator.EmptyListConstantName),
											},
										),
									},
								),
							},
						),
					},
				),
			}),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "strl/success",
			expression: expressions.NewFunctionCall("strl", []expressions.Expression{
				expressions.NewFunctionCall(
					translator.ListConstructionFunctionName,
					[]expressions.Expression{
						expressions.NewFunctionCall(
							translator.ListConstructionFunctionName,
							[]expressions.Expression{
								expressions.NewNumber(float64('"')),
								expressions.NewFunctionCall(
									translator.ListConstructionFunctionName,
									[]expressions.Expression{
										expressions.NewNumber(float64('o')),
										expressions.NewFunctionCall(
											translator.ListConstructionFunctionName,
											[]expressions.Expression{
												expressions.NewNumber(float64('n')),
												expressions.NewFunctionCall(
													translator.ListConstructionFunctionName,
													[]expressions.Expression{
														expressions.NewNumber(float64('e')),
														expressions.NewFunctionCall(
															translator.ListConstructionFunctionName,
															[]expressions.Expression{
																expressions.NewNumber(float64('"')),
																expressions.NewIdentifier(translator.EmptyListConstantName),
															},
														),
													},
												),
											},
										),
									},
								),
							},
						),
						expressions.NewFunctionCall(
							translator.ListConstructionFunctionName,
							[]expressions.Expression{
								expressions.NewFunctionCall(
									translator.ListConstructionFunctionName,
									[]expressions.Expression{
										expressions.NewNumber(float64('"')),
										expressions.NewFunctionCall(
											translator.ListConstructionFunctionName,
											[]expressions.Expression{
												expressions.NewNumber(float64('t')),
												expressions.NewFunctionCall(
													translator.ListConstructionFunctionName,
													[]expressions.Expression{
														expressions.NewNumber(float64('w')),
														expressions.NewFunctionCall(
															translator.ListConstructionFunctionName,
															[]expressions.Expression{
																expressions.NewNumber(float64('o')),
																expressions.NewFunctionCall(
																	translator.ListConstructionFunctionName,
																	[]expressions.Expression{
																		expressions.NewNumber(float64('"')),
																		expressions.NewIdentifier(translator.EmptyListConstantName),
																	},
																),
															},
														),
													},
												),
											},
										),
									},
								),
								expressions.NewIdentifier(translator.EmptyListConstantName),
							},
						),
					},
				),
			}),
			wantResult: types.NewPairFromText(`["\"one\"","\"two\""]`),
			wantErr:    assert.NoError,
		},
		{
			name: "strl/error/incorrect type",
			expression: expressions.NewFunctionCall("strl", []expressions.Expression{
				expressions.NewFunctionCall(
					translator.ListConstructionFunctionName,
					[]expressions.Expression{
						expressions.NewFunctionCall(
							translator.ListConstructionFunctionName,
							[]expressions.Expression{
								expressions.NewNumber(float64('"')),
								expressions.NewFunctionCall(
									translator.ListConstructionFunctionName,
									[]expressions.Expression{
										expressions.NewNumber(float64('o')),
										expressions.NewFunctionCall(
											translator.ListConstructionFunctionName,
											[]expressions.Expression{
												expressions.NewNumber(float64('n')),
												expressions.NewFunctionCall(
													translator.ListConstructionFunctionName,
													[]expressions.Expression{
														expressions.NewNumber(float64('e')),
														expressions.NewFunctionCall(
															translator.ListConstructionFunctionName,
															[]expressions.Expression{
																expressions.NewNumber(float64('"')),
																expressions.NewIdentifier(translator.EmptyListConstantName),
															},
														),
													},
												),
											},
										),
									},
								),
							},
						),
						expressions.NewFunctionCall(
							translator.ListConstructionFunctionName,
							[]expressions.Expression{
								expressions.NewNumber(23),
								expressions.NewIdentifier(translator.EmptyListConstantName),
							},
						),
					},
				),
			}),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "strl/error/string conversion",
			expression: expressions.NewFunctionCall("strl", []expressions.Expression{
				expressions.NewFunctionCall(
					translator.ListConstructionFunctionName,
					[]expressions.Expression{
						expressions.NewFunctionCall(
							translator.ListConstructionFunctionName,
							[]expressions.Expression{
								expressions.NewNumber(float64('"')),
								expressions.NewFunctionCall(
									translator.ListConstructionFunctionName,
									[]expressions.Expression{
										expressions.NewNumber(float64('o')),
										expressions.NewFunctionCall(
											translator.ListConstructionFunctionName,
											[]expressions.Expression{
												expressions.NewNumber(float64('n')),
												expressions.NewFunctionCall(
													translator.ListConstructionFunctionName,
													[]expressions.Expression{
														expressions.NewNumber(float64('e')),
														expressions.NewFunctionCall(
															translator.ListConstructionFunctionName,
															[]expressions.Expression{
																expressions.NewNumber(float64('"')),
																expressions.NewIdentifier(translator.EmptyListConstantName),
															},
														),
													},
												),
											},
										),
									},
								),
							},
						),
						expressions.NewFunctionCall(
							translator.ListConstructionFunctionName,
							[]expressions.Expression{
								expressions.NewFunctionCall(
									translator.ListConstructionFunctionName,
									[]expressions.Expression{
										expressions.NewNumber(float64('"')),
										expressions.NewFunctionCall(
											translator.ListConstructionFunctionName,
											[]expressions.Expression{
												expressions.NewFunctionCall(
													translator.ListConstructionFunctionName,
													[]expressions.Expression{
														expressions.NewNumber(float64('h')),
														expressions.NewFunctionCall(
															translator.ListConstructionFunctionName,
															[]expressions.Expression{
																expressions.NewNumber(float64('i')),
																expressions.NewIdentifier(translator.EmptyListConstantName),
															},
														),
													},
												),
												expressions.NewFunctionCall(
													translator.ListConstructionFunctionName,
													[]expressions.Expression{
														expressions.NewNumber(float64('w')),
														expressions.NewFunctionCall(
															translator.ListConstructionFunctionName,
															[]expressions.Expression{
																expressions.NewNumber(float64('o')),
																expressions.NewFunctionCall(
																	translator.ListConstructionFunctionName,
																	[]expressions.Expression{
																		expressions.NewNumber(float64('"')),
																		expressions.NewIdentifier(translator.EmptyListConstantName),
																	},
																),
															},
														),
													},
												),
											},
										),
									},
								),
								expressions.NewIdentifier(translator.EmptyListConstantName),
							},
						),
					},
				),
			}),
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

func TestValues_nan(test *testing.T) {
	ctx := context.NewDefaultContext()
	context.SetValues(ctx, Values)

	expression := expressions.NewIdentifier("nan")
	result, err := expression.Evaluate(ctx)

	if assert.NoError(test, err) {
		require.IsType(test, float64(0), result)
		assert.True(test, math.IsNaN(result.(float64)))
	}
}

func TestValues_inDelta(test *testing.T) {
	for _, data := range []struct {
		name       string
		expression expressions.Expression
		want       float64
	}{
		{
			name: "sin",
			expression: expressions.NewFunctionCall("sin", []expressions.Expression{
				expressions.NewNumber(23),
			}),
			want: -0.846220,
		},
		{
			name: "cos",
			expression: expressions.NewFunctionCall("cos", []expressions.Expression{
				expressions.NewNumber(23),
			}),
			want: -0.532833,
		},
		{
			name: "tn",
			expression: expressions.NewFunctionCall("tn", []expressions.Expression{
				expressions.NewNumber(23),
			}),
			want: 1.588153,
		},
		{
			name: "arcsin",
			expression: expressions.NewFunctionCall("arcsin", []expressions.Expression{
				expressions.NewNumber(0.5),
			}),
			want: 0.523598,
		},
		{
			name: "arccos",
			expression: expressions.NewFunctionCall("arccos", []expressions.Expression{
				expressions.NewNumber(0.5),
			}),
			want: 1.047197,
		},
		{
			name: "arctn",
			expression: expressions.NewFunctionCall("arctn", []expressions.Expression{
				expressions.NewNumber(0.5),
			}),
			want: 0.463647,
		},
		{
			name: "angle",
			expression: expressions.NewFunctionCall("angle", []expressions.Expression{
				expressions.NewNumber(2),
				expressions.NewNumber(3),
			}),
			want: 0.982793,
		},
		{
			name: "exp",
			expression: expressions.NewFunctionCall("exp", []expressions.Expression{
				expressions.NewNumber(2.3),
			}),
			want: 9.974182,
		},
		{
			name: "ln",
			expression: expressions.NewFunctionCall("ln", []expressions.Expression{
				expressions.NewNumber(23),
			}),
			want: 3.135494,
		},
		{
			name: "lg",
			expression: expressions.NewFunctionCall("lg", []expressions.Expression{
				expressions.NewNumber(23),
			}),
			want: 1.361727,
		},
	} {
		ctx := context.NewDefaultContext()
		context.SetValues(ctx, Values)

		got, err := data.expression.Evaluate(ctx)

		if assert.NoError(test, err) {
			require.IsType(test, float64(0), got)
			assert.InDelta(test, data.want, got.(float64), 1e-6)
		}
	}
}
