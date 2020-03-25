package translator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/internal/tests"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
)

func TestTranslateExpression(test *testing.T) {
	type args struct {
		expression          *parser.Expression
		declaredIdentifiers context.ValueNameGroup
	}

	for _, data := range []struct {
		name           string
		args           args
		wantExpression expressions.Expression
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "Expression/success",
			args: args{
				expression: &parser.Expression{
					ListConstruction: &parser.ListConstruction{
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewNumber(23),
			wantErr:        assert.NoError,
		},
		{
			name: "Expression/error",
			args: args{
				expression: &parser.Expression{
					ListConstruction: &parser.ListConstruction{
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{
										Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotErr := translateExpression(data.args.expression, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateListConstruction(test *testing.T) {
	type args struct {
		listConstruction    *parser.ListConstruction
		declaredIdentifiers context.ValueNameGroup
	}

	for _, data := range []struct {
		name           string
		args           args
		wantExpression expressions.Expression
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "ListConstruction/nonempty/success",
			args: args{
				listConstruction: &parser.ListConstruction{
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
							},
						},
					},
					ListConstruction: &parser.ListConstruction{
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: tests.GetStringAddress("test")}},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewFunctionCall(
				ListConstructionFunctionName,
				[]expressions.Expression{expressions.NewNumber(12), expressions.NewIdentifier("test")},
			),
			wantErr: assert.NoError,
		},
		{
			name: "ListConstruction/nonempty/error",
			args: args{
				listConstruction: &parser.ListConstruction{
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
							},
						},
					},
					ListConstruction: &parser.ListConstruction{
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{
										Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "ListConstruction/empty/success",
			args: args{
				listConstruction: &parser.ListConstruction{
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewNumber(23),
			wantErr:        assert.NoError,
		},
		{
			name: "ListConstruction/empty/error",
			args: args{
				listConstruction: &parser.ListConstruction{
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{
									Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotErr :=
				translateListConstruction(data.args.listConstruction, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateAddition(test *testing.T) {
	type args struct {
		addition            *parser.Addition
		declaredIdentifiers context.ValueNameGroup
	}

	for _, data := range []struct {
		name           string
		args           args
		wantExpression expressions.Expression
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "Addition/nonempty/success/addition",
			args: args{
				addition: &parser.Addition{
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
						},
					},
					Operation: "+",
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
							},
						},
						Operation: "+",
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(42)}},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewFunctionCall(AdditionFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(AdditionFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantErr: assert.NoError,
		},
		{
			name: "Addition/nonempty/success/subtraction",
			args: args{
				addition: &parser.Addition{
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
						},
					},
					Operation: "-",
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
							},
						},
						Operation: "-",
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(42)}},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewFunctionCall(SubtractionFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(SubtractionFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantErr: assert.NoError,
		},
		{
			name: "Addition/nonempty/error",
			args: args{
				addition: &parser.Addition{
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
						},
					},
					Operation: "+",
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
							},
						},
						Operation: "+",
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{
										Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Addition/empty/success",
			args: args{
				addition: &parser.Addition{
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewNumber(23),
			wantErr:        assert.NoError,
		},
		{
			name: "Addition/empty/error",
			args: args{
				addition: &parser.Addition{
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{
								Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotErr := translateAddition(data.args.addition, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateMultiplication(test *testing.T) {
	type args struct {
		multiplication      *parser.Multiplication
		declaredIdentifiers context.ValueNameGroup
	}

	for _, data := range []struct {
		name           string
		args           args
		wantExpression expressions.Expression
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "Multiplication/nonempty/success/multiplication",
			args: args{
				multiplication: &parser.Multiplication{
					Unary: &parser.Unary{
						Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
					},
					Operation: "*",
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
						},
						Operation: "*",
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(42)}},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewFunctionCall(MultiplicationFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(MultiplicationFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantErr: assert.NoError,
		},
		{
			name: "Multiplication/nonempty/success/division",
			args: args{
				multiplication: &parser.Multiplication{
					Unary: &parser.Unary{
						Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
					},
					Operation: "/",
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
						},
						Operation: "/",
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(42)}},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewFunctionCall(DivisionFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(DivisionFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantErr: assert.NoError,
		},
		{
			name: "Multiplication/nonempty/success/modulo",
			args: args{
				multiplication: &parser.Multiplication{
					Unary: &parser.Unary{
						Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
					},
					Operation: "%",
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
						},
						Operation: "%",
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(42)}},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewFunctionCall(ModuloFunctionName, []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewFunctionCall(ModuloFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				}),
			}),
			wantErr: assert.NoError,
		},
		{
			name: "Multiplication/nonempty/error",
			args: args{
				multiplication: &parser.Multiplication{
					Unary: &parser.Unary{
						Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
					},
					Operation: "*",
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
						},
						Operation: "*",
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{
									Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Multiplication/empty/success",
			args: args{
				multiplication: &parser.Multiplication{
					Unary: &parser.Unary{
						Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewNumber(23),
			wantErr:        assert.NoError,
		},
		{
			name: "Multiplication/empty/error",
			args: args{
				multiplication: &parser.Multiplication{
					Unary: &parser.Unary{
						Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")}},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotErr :=
				translateMultiplication(data.args.multiplication, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateUnary(test *testing.T) {
	type args struct {
		unary               *parser.Unary
		declaredIdentifiers context.ValueNameGroup
	}

	for _, data := range []struct {
		name           string
		args           args
		wantExpression expressions.Expression
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "Unary/nonempty/success",
			args: args{
				unary: &parser.Unary{
					Operation: "-",
					Unary: &parser.Unary{
						Operation: "-",
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewFunctionCall(NegationFunctionName, []expressions.Expression{
				expressions.NewFunctionCall(NegationFunctionName, []expressions.Expression{
					expressions.NewNumber(23),
				}),
			}),
			wantErr: assert.NoError,
		},
		{
			name: "Unary/nonempty/error",
			args: args{
				unary: &parser.Unary{
					Operation: "-",
					Unary: &parser.Unary{
						Operation: "-",
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{
								Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Unary/empty/success",
			args: args{
				unary: &parser.Unary{
					Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewNumber(23),
			wantErr:        assert.NoError,
		},
		{
			name: "Unary/empty/error",
			args: args{
				unary: &parser.Unary{
					Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")}},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotErr := translateUnary(data.args.unary, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateAccessor(test *testing.T) {
	type args struct {
		accessor            *parser.Accessor
		declaredIdentifiers context.ValueNameGroup
	}

	for _, data := range []struct {
		name           string
		args           args
		wantExpression expressions.Expression
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "Accessor/nonempty/success",
			args: args{
				accessor: &parser.Accessor{
					Atom: &parser.Atom{Identifier: tests.GetStringAddress("test")},
					Keys: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
										},
									},
								},
							},
						},
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewFunctionCall(KeyAccessorFunctionName, []expressions.Expression{
				expressions.NewFunctionCall(KeyAccessorFunctionName, []expressions.Expression{
					expressions.NewIdentifier("test"),
					expressions.NewNumber(12),
				}),
				expressions.NewNumber(23),
			}),
			wantErr: assert.NoError,
		},
		{
			name: "Accessor/nonempty/error/atom translating",
			args: args{
				accessor: &parser.Accessor{
					Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
					Keys: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
										},
									},
								},
							},
						},
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Accessor/nonempty/error/key translating",
			args: args{
				accessor: &parser.Accessor{
					Atom: &parser.Atom{Identifier: tests.GetStringAddress("test")},
					Keys: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
										},
									},
								},
							},
						},
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{
												Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Accessor/empty/success",
			args: args{
				accessor:            &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewNumber(23),
			wantErr:        assert.NoError,
		},
		{
			name: "Accessor/empty/error",
			args: args{
				accessor: &parser.Accessor{
					Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotErr := translateAccessor(data.args.accessor, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateAtom(test *testing.T) {
	type args struct {
		atom                *parser.Atom
		declaredIdentifiers context.ValueNameGroup
	}

	for _, data := range []struct {
		name           string
		args           args
		wantExpression expressions.Expression
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "Atom/number",
			args: args{
				atom:                &parser.Atom{Number: tests.GetNumberAddress(23)},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewNumber(23),
			wantErr:        assert.NoError,
		},
		{
			name: "Atom/string",
			args: args{
				atom:                &parser.Atom{String: tests.GetStringAddress("test")},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewString("test"),
			wantErr:        assert.NoError,
		},
		{
			name: "Atom/list definition/success",
			args: args{
				atom: &parser.Atom{
					ListDefinition: &parser.ListDefinition{
						Items: []*parser.Expression{
							{
								ListConstruction: &parser.ListConstruction{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
											},
										},
									},
								},
							},
							{
								ListConstruction: &parser.ListConstruction{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
											},
										},
									},
								},
							},
							{
								ListConstruction: &parser.ListConstruction{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(42)}},
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
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
			wantErr: assert.NoError,
		},
		{
			name: "Atom/list definition/error",
			args: args{
				atom: &parser.Atom{
					ListDefinition: &parser.ListDefinition{
						Items: []*parser.Expression{
							{
								ListConstruction: &parser.ListConstruction{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
											},
										},
									},
								},
							},
							{
								ListConstruction: &parser.ListConstruction{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
											},
										},
									},
								},
							},
							{
								ListConstruction: &parser.ListConstruction{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{
													Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Atom/function call/success",
			args: args{
				atom: &parser.Atom{
					FunctionCall: &parser.FunctionCall{
						Name: "test",
						Arguments: []*parser.Expression{
							{
								ListConstruction: &parser.ListConstruction{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
											},
										},
									},
								},
							},
							{
								ListConstruction: &parser.ListConstruction{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
											},
										},
									},
								},
							},
							{
								ListConstruction: &parser.ListConstruction{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(42)}},
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewFunctionCall("test", []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewNumber(23),
				expressions.NewNumber(42),
			}),
			wantErr: assert.NoError,
		},
		{
			name: "Atom/function call/error",
			args: args{
				atom: &parser.Atom{
					FunctionCall: &parser.FunctionCall{
						Name: "test",
						Arguments: []*parser.Expression{
							{
								ListConstruction: &parser.ListConstruction{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
											},
										},
									},
								},
							},
							{
								ListConstruction: &parser.ListConstruction{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
											},
										},
									},
								},
							},
							{
								ListConstruction: &parser.ListConstruction{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{
													Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Atom/identifier/success",
			args: args{
				atom:                &parser.Atom{Identifier: tests.GetStringAddress("test")},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewIdentifier("test"),
			wantErr:        assert.NoError,
		},
		{
			name: "Atom/identifier/error",
			args: args{
				atom:                &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Atom/expression/success",
			args: args{
				atom: &parser.Atom{
					Expression: &parser.Expression{
						ListConstruction: &parser.ListConstruction{
							Addition: &parser.Addition{
								Multiplication: &parser.Multiplication{
									Unary: &parser.Unary{
										Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewNumber(23),
			wantErr:        assert.NoError,
		},
		{
			name: "Atom/expression/error",
			args: args{
				atom: &parser.Atom{
					Expression: &parser.Expression{
						ListConstruction: &parser.ListConstruction{
							Addition: &parser.Addition{
								Multiplication: &parser.Multiplication{
									Unary: &parser.Unary{
										Accessor: &parser.Accessor{
											Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotErr := translateAtom(data.args.atom, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateListDefinition(test *testing.T) {
	type args struct {
		listDefinition      *parser.ListDefinition
		declaredIdentifiers context.ValueNameGroup
	}

	for _, data := range []struct {
		name           string
		args           args
		wantExpression expressions.Expression
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "ListDefinition/success/few items",
			args: args{
				listDefinition: &parser.ListDefinition{
					Items: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
										},
									},
								},
							},
						},
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
										},
									},
								},
							},
						},
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(42)}},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
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
			wantErr: assert.NoError,
		},
		{
			name: "ListDefinition/success/no items",
			args: args{
				listDefinition: &parser.ListDefinition{
					Items: nil,
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewIdentifier(EmptyListConstantName),
			wantErr:        assert.NoError,
		},
		{
			name: "ListDefinition/error",
			args: args{
				listDefinition: &parser.ListDefinition{
					Items: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
										},
									},
								},
							},
						},
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
										},
									},
								},
							},
						},
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{
												Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotErr :=
				translateListDefinition(data.args.listDefinition, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateFunctionCall(test *testing.T) {
	type args struct {
		functionCall        *parser.FunctionCall
		declaredIdentifiers context.ValueNameGroup
	}

	for _, data := range []struct {
		name           string
		args           args
		wantExpression expressions.Expression
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "FunctionCall/success/few arguments",
			args: args{
				functionCall: &parser.FunctionCall{
					Name: "test",
					Arguments: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
										},
									},
								},
							},
						},
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
										},
									},
								},
							},
						},
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(42)}},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewFunctionCall("test", []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewNumber(23),
				expressions.NewNumber(42),
			}),
			wantErr: assert.NoError,
		},
		{
			name: "FunctionCall/success/no arguments",
			args: args{
				functionCall: &parser.FunctionCall{
					Name:      "test",
					Arguments: nil,
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: expressions.NewFunctionCall("test", nil),
			wantErr:        assert.NoError,
		},
		{
			name: "FunctionCall/error/unknown function",
			args: args{
				functionCall: &parser.FunctionCall{
					Name: "unknown",
					Arguments: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
										},
									},
								},
							},
						},
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
										},
									},
								},
							},
						},
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(42)}},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "FunctionCall/error/argument translating",
			args: args{
				functionCall: &parser.FunctionCall{
					Name: "test",
					Arguments: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(12)}},
										},
									},
								},
							},
						},
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
										},
									},
								},
							},
						},
						{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{
												Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: context.ValueNameGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotErr :=
				translateFunctionCall(data.args.functionCall, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			data.wantErr(test, gotErr)
		})
	}
}
