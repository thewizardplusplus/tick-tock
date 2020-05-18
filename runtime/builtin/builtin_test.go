package builtin

import (
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"testing"
	"time"

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
			name:       "nil",
			expression: expressions.NewIdentifier("nil"),
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name:       "false",
			expression: expressions.NewIdentifier("false"),
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "true",
			expression: expressions.NewIdentifier("true"),
			wantResult: types.True,
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
			name: "equal/success/false/same types",
			expression: expressions.NewFunctionCall(
				translator.EqualFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewNumber(3)},
			),
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name: "equal/success/false/different types",
			expression: expressions.NewFunctionCall(
				translator.EqualFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewIdentifier("nil")},
			),
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name: "equal/success/true",
			expression: expressions.NewFunctionCall(
				translator.EqualFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewNumber(2)},
			),
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name: "equal/error",
			expression: expressions.NewFunctionCall(
				translator.EqualFunctionName,
				[]expressions.Expression{
					expressions.NewIdentifier(translator.EqualFunctionName),
					expressions.NewIdentifier("nil"),
				},
			),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "not equal/success/false/same types",
			expression: expressions.NewFunctionCall(
				translator.NotEqualFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewNumber(3)},
			),
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name: "not equal/success/false/different types",
			expression: expressions.NewFunctionCall(
				translator.NotEqualFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewIdentifier("nil")},
			),
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name: "not equal/success/true",
			expression: expressions.NewFunctionCall(
				translator.NotEqualFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewNumber(2)},
			),
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name: "not equal/error",
			expression: expressions.NewFunctionCall(
				translator.NotEqualFunctionName,
				[]expressions.Expression{
					expressions.NewIdentifier(translator.NotEqualFunctionName),
					expressions.NewIdentifier("nil"),
				},
			),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "less/success/false",
			expression: expressions.NewFunctionCall(
				translator.LessFunctionName,
				[]expressions.Expression{expressions.NewNumber(4), expressions.NewNumber(2)},
			),
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name: "less/success/true",
			expression: expressions.NewFunctionCall(
				translator.LessFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewNumber(3)},
			),
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name: "less/error",
			expression: expressions.NewFunctionCall(
				translator.LessFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewIdentifier("nil")},
			),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "less or equal/success/false",
			expression: expressions.NewFunctionCall(
				translator.LessOrEqualFunctionName,
				[]expressions.Expression{expressions.NewNumber(4), expressions.NewNumber(2)},
			),
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name: "less or equal/success/true/less",
			expression: expressions.NewFunctionCall(
				translator.LessOrEqualFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewNumber(3)},
			),
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name: "less or equal/success/true/equal",
			expression: expressions.NewFunctionCall(
				translator.LessOrEqualFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewNumber(2)},
			),
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name: "less or equal/error",
			expression: expressions.NewFunctionCall(
				translator.LessOrEqualFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewIdentifier("nil")},
			),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "greater/success/false",
			expression: expressions.NewFunctionCall(
				translator.GreaterFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewNumber(3)},
			),
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name: "greater/success/true",
			expression: expressions.NewFunctionCall(
				translator.GreaterFunctionName,
				[]expressions.Expression{expressions.NewNumber(4), expressions.NewNumber(2)},
			),
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name: "greater/error",
			expression: expressions.NewFunctionCall(
				translator.GreaterFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewIdentifier("nil")},
			),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "greater or equal/success/false",
			expression: expressions.NewFunctionCall(
				translator.GreaterOrEqualFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewNumber(3)},
			),
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name: "greater or equal/success/true/greater",
			expression: expressions.NewFunctionCall(
				translator.GreaterOrEqualFunctionName,
				[]expressions.Expression{expressions.NewNumber(4), expressions.NewNumber(2)},
			),
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name: "greater or equal/success/true/equal",
			expression: expressions.NewFunctionCall(
				translator.GreaterOrEqualFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewNumber(2)},
			),
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name: "greater or equal/error",
			expression: expressions.NewFunctionCall(
				translator.GreaterOrEqualFunctionName,
				[]expressions.Expression{expressions.NewNumber(2), expressions.NewIdentifier("nil")},
			),
			wantResult: nil,
			wantErr:    assert.Error,
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
			name: "arithmetic negation",
			expression: expressions.NewFunctionCall(
				translator.ArithmeticNegationFunctionName,
				[]expressions.Expression{expressions.NewNumber(23)},
			),
			wantResult: -23.0,
			wantErr:    assert.NoError,
		},
		{
			name: "logical negation/success/false",
			expression: expressions.NewFunctionCall(
				translator.LogicalNegationFunctionName,
				[]expressions.Expression{expressions.NewIdentifier("false")},
			),
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name: "logical negation/success/true",
			expression: expressions.NewFunctionCall(
				translator.LogicalNegationFunctionName,
				[]expressions.Expression{expressions.NewIdentifier("true")},
			),
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name: "logical negation/error",
			expression: expressions.NewFunctionCall(
				translator.LogicalNegationFunctionName,
				[]expressions.Expression{expressions.NewIdentifier(translator.LogicalNegationFunctionName)},
			),
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "key accessor/index in range",
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
			name: "key accessor/index out of range",
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
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name: "type/success/nil",
			expression: expressions.NewFunctionCall("type", []expressions.Expression{
				expressions.NewIdentifier("nil"),
			}),
			wantResult: types.NewPairFromText("nil"),
			wantErr:    assert.NoError,
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
			name: "bool/success/false",
			expression: expressions.NewFunctionCall("bool", []expressions.Expression{
				expressions.NewString(""),
			}),
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name: "bool/success/true",
			expression: expressions.NewFunctionCall("bool", []expressions.Expression{
				expressions.NewString("test"),
			}),
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name: "bool/error",
			expression: expressions.NewFunctionCall("bool", []expressions.Expression{
				expressions.NewIdentifier("bool"),
			}),
			wantResult: nil,
			wantErr:    assert.Error,
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
			name: "is_nan/false",
			expression: expressions.NewFunctionCall("is_nan", []expressions.Expression{
				expressions.NewNumber(23),
			}),
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name: "is_nan/true",
			expression: expressions.NewFunctionCall("is_nan", []expressions.Expression{
				expressions.NewIdentifier("nan"),
			}),
			wantResult: types.True,
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
			name: "num/success/correct number",
			expression: expressions.NewFunctionCall("num", []expressions.Expression{
				expressions.NewString("23"),
			}),
			wantResult: 23.0,
			wantErr:    assert.NoError,
		},
		{
			name: "num/success/incorrect number",
			expression: expressions.NewFunctionCall("num", []expressions.Expression{
				expressions.NewString("test"),
			}),
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name: "num/error",
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
			name: "str/success/nil",
			expression: expressions.NewFunctionCall("str", []expressions.Expression{
				expressions.NewIdentifier("nil"),
			}),
			wantResult: types.NewPairFromText("nil"),
			wantErr:    assert.NoError,
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
			name: "str/success/*types.Pair/with the nil type",
			expression: expressions.NewFunctionCall("str", []expressions.Expression{
				expressions.NewFunctionCall(
					translator.ListConstructionFunctionName,
					[]expressions.Expression{
						expressions.NewNumber(12),
						expressions.NewFunctionCall(
							translator.ListConstructionFunctionName,
							[]expressions.Expression{
								expressions.NewIdentifier("nil"),
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
			wantResult: types.NewPairFromText("[12,null,42]"),
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
			name: "strb/success/false",
			expression: expressions.NewFunctionCall("strb", []expressions.Expression{
				expressions.NewString(""),
			}),
			wantResult: types.NewPairFromText("false"),
			wantErr:    assert.NoError,
		},
		{
			name: "strb/success/true",
			expression: expressions.NewFunctionCall("strb", []expressions.Expression{
				expressions.NewString("test"),
			}),
			wantResult: types.NewPairFromText("true"),
			wantErr:    assert.NoError,
		},
		{
			name: "strb/error",
			expression: expressions.NewFunctionCall("strb", []expressions.Expression{
				expressions.NewIdentifier("strb"),
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
		test.Run(data.name, func(test *testing.T) {
			ctx := context.NewDefaultContext()
			context.SetValues(ctx, Values)

			gotResult, gotErr := data.expression.Evaluate(ctx)

			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
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
		test.Run(data.name, func(test *testing.T) {
			ctx := context.NewDefaultContext()
			context.SetValues(ctx, Values)

			got, err := data.expression.Evaluate(ctx)

			if assert.NoError(test, err) {
				require.IsType(test, float64(0), got)
				assert.InDelta(test, data.want, got.(float64), 1e-6)
			}
		})
	}
}

func TestValues_random(test *testing.T) {
	const numberCount = 10

	ctx := context.NewDefaultContext()
	context.SetValues(ctx, Values)

	expression := expressions.NewFunctionCall("seed", []expressions.Expression{
		expressions.NewNumber(23),
	})
	got, err := expression.Evaluate(ctx)

	assert.Equal(test, types.Nil{}, got)
	assert.NoError(test, err)

	var numbers []float64
	for i := 0; i < numberCount; i++ {
		expression := expressions.NewFunctionCall("random", nil)
		result, err := expression.Evaluate(ctx)

		assert.IsType(test, float64(0), result)
		assert.NoError(test, err)

		if number, ok := result.(float64); ok {
			numbers = append(numbers, number)
		}
	}

	rand.Seed(23)

	var wantNumbers []float64
	for i := 0; i < numberCount; i++ {
		wantNumber := rand.Float64()
		wantNumbers = append(wantNumbers, wantNumber)
	}

	assert.InDeltaSlice(test, wantNumbers, numbers, 1e-6)
}

func TestValues_env(test *testing.T) {
	type args struct {
		name expressions.Expression
	}

	const envName = "TEST"
	for _, data := range []struct {
		name       string
		prepare    func(test *testing.T)
		args       args
		wantResult interface{}
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success/existing variable/nonempty value",
			prepare: func(test *testing.T) {
				err := os.Setenv(envName, "test")
				require.NoError(test, err)
			},
			args: args{
				name: expressions.NewString(envName),
			},
			wantResult: types.NewPairFromText("test"),
			wantErr:    assert.NoError,
		},
		{
			name: "success/existing variable/empty value",
			prepare: func(test *testing.T) {
				err := os.Setenv(envName, "")
				require.NoError(test, err)
			},
			args: args{
				name: expressions.NewString(envName),
			},
			wantResult: (*types.Pair)(nil),
			wantErr:    assert.NoError,
		},
		{
			name: "success/nonexistent variable",
			prepare: func(test *testing.T) {
				err := os.Unsetenv(envName)
				require.NoError(test, err)
			},
			args: args{
				name: expressions.NewString(envName),
			},
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name:    "error",
			prepare: func(test *testing.T) {},
			args: args{
				name: expressions.NewFunctionCall(
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
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			previousValue, wasSet := os.LookupEnv(envName)
			defer func() {
				if wasSet {
					err := os.Setenv(envName, previousValue)
					require.NoError(test, err)
				}
			}()
			data.prepare(test)

			ctx := context.NewDefaultContext()
			context.SetValues(ctx, Values)

			expression := expressions.NewFunctionCall("env", []expressions.Expression{data.args.name})
			gotResult, gotErr := expression.Evaluate(ctx)

			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}

func TestValues_time(test *testing.T) {
	ctx := context.NewDefaultContext()
	context.SetValues(ctx, Values)

	expression := expressions.NewFunctionCall("time", nil)
	result, err := expression.Evaluate(ctx)

	if assert.NoError(test, err) {
		require.IsType(test, float64(0), result)

		resultTime := time.Unix(0, int64(result.(float64)*1e9))
		assert.WithinDuration(test, time.Now(), resultTime, time.Minute)
	}
}

func TestValues_sleep(test *testing.T) {
	startTime := time.Now()

	ctx := context.NewDefaultContext()
	context.SetValues(ctx, Values)

	expression := expressions.NewFunctionCall("sleep", []expressions.Expression{
		expressions.NewNumber(2.3),
	})
	result, err := expression.Evaluate(ctx)

	elapsedTime := int64(time.Since(startTime))
	assert.GreaterOrEqual(test, elapsedTime, int64(2300*time.Millisecond))
	assert.Less(test, elapsedTime, int64(time.Minute))
	assert.Equal(test, types.Nil{}, result)
	assert.NoError(test, err)
}

// based on https://talks.golang.org/2014/testing.slide#23 by Andrew Gerrand
func TestValues_exit(test *testing.T) {
	if os.Getenv("EXIT_TEST") == "TRUE" {
		ctx := context.NewDefaultContext()
		context.SetValues(ctx, Values)

		expression := expressions.NewFunctionCall("exit", []expressions.Expression{
			expressions.NewNumber(23),
		})
		expression.Evaluate(ctx) // nolint: errcheck

		return
	}

	command := exec.Command(os.Args[0], "-test.run=TestValues_exit")
	command.Env = append(os.Environ(), "EXIT_TEST=TRUE")

	err := command.Run()

	assert.IsType(test, (*exec.ExitError)(nil), err)
	assert.EqualError(test, err, "exit status 23")
}

func TestValues_input(test *testing.T) {
	for _, data := range []struct {
		name       string
		prepare    func(test *testing.T, tempFile *os.File)
		expression expressions.Expression
		wantResult interface{}
	}{
		{
			name: "in/part of symbols/success/part of symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			expression: expressions.NewFunctionCall("in", []expressions.Expression{
				expressions.NewNumber(2),
			}),
			wantResult: types.NewPairFromText("te"),
		},
		{
			name: "in/part of symbols/success/all symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			expression: expressions.NewFunctionCall("in", []expressions.Expression{
				expressions.NewNumber(4),
			}),
			wantResult: types.NewPairFromText("test"),
		},
		{
			name:    "in/part of symbols/error/without symbols",
			prepare: func(test *testing.T, tempFile *os.File) {},
			expression: expressions.NewFunctionCall("in", []expressions.Expression{
				expressions.NewNumber(2),
			}),
			wantResult: types.Nil{},
		},
		{
			name: "in/part of symbols/error/with lack of symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			expression: expressions.NewFunctionCall("in", []expressions.Expression{
				expressions.NewNumber(5),
			}),
			wantResult: types.Nil{},
		},
		{
			name: "in/all symbols/with symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			expression: expressions.NewFunctionCall("in", []expressions.Expression{
				expressions.NewNumber(-1),
			}),
			wantResult: types.NewPairFromText("test"),
		},
		{
			name:    "in/all symbols/without symbols",
			prepare: func(test *testing.T, tempFile *os.File) {},
			expression: expressions.NewFunctionCall("in", []expressions.Expression{
				expressions.NewNumber(-1),
			}),
			wantResult: (*types.Pair)(nil),
		},
		{
			name: "inln/part of symbols/success/part of symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			expression: expressions.NewFunctionCall("inln", []expressions.Expression{
				expressions.NewNumber(2),
			}),
			wantResult: types.NewPairFromText("te"),
		},
		{
			name: "inln/part of symbols/success/all symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			expression: expressions.NewFunctionCall("inln", []expressions.Expression{
				expressions.NewNumber(4),
			}),
			wantResult: types.NewPairFromText("test"),
		},
		{
			name:    "inln/part of symbols/error/without symbols",
			prepare: func(test *testing.T, tempFile *os.File) {},
			expression: expressions.NewFunctionCall("inln", []expressions.Expression{
				expressions.NewNumber(2),
			}),
			wantResult: types.Nil{},
		},
		{
			name: "inln/part of symbols/error/with lack of symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			expression: expressions.NewFunctionCall("inln", []expressions.Expression{
				expressions.NewNumber(5),
			}),
			wantResult: types.Nil{},
		},
		{
			name: "inln/all symbols/success/with symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test #1\ntest #2\n")
				require.NoError(test, err)
			},
			expression: expressions.NewFunctionCall("inln", []expressions.Expression{
				expressions.NewNumber(-1),
			}),
			wantResult: types.NewPairFromText("test #1\n"),
		},
		{
			name: "inln/all symbols/success/without symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("\ntest #2\n")
				require.NoError(test, err)
			},
			expression: expressions.NewFunctionCall("inln", []expressions.Expression{
				expressions.NewNumber(-1),
			}),
			wantResult: types.NewPairFromText("\n"),
		},
		{
			name: "inln/all symbols/error/with symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			expression: expressions.NewFunctionCall("inln", []expressions.Expression{
				expressions.NewNumber(-1),
			}),
			wantResult: types.Nil{},
		},
		{
			name:    "inln/all symbols/error/without symbols",
			prepare: func(test *testing.T, tempFile *os.File) {},
			expression: expressions.NewFunctionCall("inln", []expressions.Expression{
				expressions.NewNumber(-1),
			}),
			wantResult: types.Nil{},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			previousStdin := os.Stdin
			defer func() { os.Stdin = previousStdin }()

			tempFile, err := ioutil.TempFile("", "test.*")
			require.NoError(test, err)
			defer os.Remove(tempFile.Name()) // nolint: errcheck
			defer tempFile.Close()           // nolint: errcheck

			data.prepare(test, tempFile)
			err = tempFile.Close()
			require.NoError(test, err)

			tempFile, err = os.Open(tempFile.Name())
			require.NoError(test, err)
			os.Stdin = tempFile

			ctx := context.NewDefaultContext()
			context.SetValues(ctx, Values)

			gotResult, gotErr := data.expression.Evaluate(ctx)

			assert.Equal(test, data.wantResult, gotResult)
			assert.NoError(test, gotErr)
		})
	}
}

func TestValues_output(test *testing.T) {
	for _, data := range []struct {
		name       string
		prepare    func(test *testing.T, tempFile *os.File)
		expression expressions.Expression
		wantResult interface{}
		wantOutput string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:    "out/success",
			prepare: func(test *testing.T, tempFile *os.File) { os.Stdout = tempFile },
			expression: expressions.NewFunctionCall("out", []expressions.Expression{
				expressions.NewString("test"),
			}),
			wantResult: types.Nil{},
			wantOutput: "test",
			wantErr:    assert.NoError,
		},
		{
			name:    "out/error",
			prepare: func(test *testing.T, tempFile *os.File) { os.Stdout = tempFile },
			expression: expressions.NewFunctionCall("out", []expressions.Expression{
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
			wantOutput: "",
			wantErr:    assert.Error,
		},
		{
			name:    "outln/success",
			prepare: func(test *testing.T, tempFile *os.File) { os.Stdout = tempFile },
			expression: expressions.NewFunctionCall("outln", []expressions.Expression{
				expressions.NewString("test"),
			}),
			wantResult: types.Nil{},
			wantOutput: "test\n",
			wantErr:    assert.NoError,
		},
		{
			name:    "outln/error",
			prepare: func(test *testing.T, tempFile *os.File) { os.Stdout = tempFile },
			expression: expressions.NewFunctionCall("outln", []expressions.Expression{
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
			wantOutput: "",
			wantErr:    assert.Error,
		},
		{
			name:    "err/success",
			prepare: func(test *testing.T, tempFile *os.File) { os.Stderr = tempFile },
			expression: expressions.NewFunctionCall("err", []expressions.Expression{
				expressions.NewString("test"),
			}),
			wantResult: types.Nil{},
			wantOutput: "test",
			wantErr:    assert.NoError,
		},
		{
			name:    "err/error",
			prepare: func(test *testing.T, tempFile *os.File) { os.Stderr = tempFile },
			expression: expressions.NewFunctionCall("err", []expressions.Expression{
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
			wantOutput: "",
			wantErr:    assert.Error,
		},
		{
			name:    "errln/success",
			prepare: func(test *testing.T, tempFile *os.File) { os.Stderr = tempFile },
			expression: expressions.NewFunctionCall("errln", []expressions.Expression{
				expressions.NewString("test"),
			}),
			wantResult: types.Nil{},
			wantOutput: "test\n",
			wantErr:    assert.NoError,
		},
		{
			name:    "errln/error",
			prepare: func(test *testing.T, tempFile *os.File) { os.Stderr = tempFile },
			expression: expressions.NewFunctionCall("errln", []expressions.Expression{
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
			wantOutput: "",
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			previousStdout, previousStderr := os.Stdout, os.Stderr
			defer func() { os.Stdout, os.Stderr = previousStdout, previousStderr }()

			tempFile, err := ioutil.TempFile("", "test.*")
			require.NoError(test, err)
			defer os.Remove(tempFile.Name()) // nolint: errcheck
			defer tempFile.Close()           // nolint: errcheck
			data.prepare(test, tempFile)

			ctx := context.NewDefaultContext()
			context.SetValues(ctx, Values)

			gotResult, gotErr := data.expression.Evaluate(ctx)

			err = tempFile.Close()
			require.NoError(test, err)

			gotOutputBytes, err := ioutil.ReadFile(tempFile.Name())
			require.NoError(test, err)

			assert.Equal(test, data.wantResult, gotResult)
			assert.Equal(test, data.wantOutput, string(gotOutputBytes))
			data.wantErr(test, gotErr)
		})
	}
}
