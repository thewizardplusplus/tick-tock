package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/internal/tests"
)

func TestParseToAST_withExpression(test *testing.T) {
	type args struct {
		code string
		ast  interface{}
	}

	for _, testData := range []struct {
		name    string
		args    args
		wantAST interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Atom/number/integer",
			args:    args{"23", new(Atom)},
			wantAST: &Atom{Number: tests.GetNumberAddress(23)},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/number/floating-point",
			args:    args{"2.3", new(Atom)},
			wantAST: &Atom{Number: tests.GetNumberAddress(2.3)},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/string/interpreted",
			args:    args{`"test"`, new(Atom)},
			wantAST: &Atom{String: tests.GetStringAddress("test")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/string/raw",
			args:    args{"`test`", new(Atom)},
			wantAST: &Atom{String: tests.GetStringAddress("test")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/identifier",
			args:    args{"test", new(Atom)},
			wantAST: &Atom{Identifier: tests.GetStringAddress("test")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/list definition/no items",
			args:    args{"[]", new(Atom)},
			wantAST: &Atom{ListDefinition: &ListDefinition{Items: nil}},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/list definition/single item",
			args: args{"[12]", new(Atom)},
			wantAST: &Atom{
				ListDefinition: &ListDefinition{
					Items: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(12)}}},
									},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/list definition/single item/trailing comma",
			args: args{"[12,]", new(Atom)},
			wantAST: &Atom{
				ListDefinition: &ListDefinition{
					Items: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(12)}}},
									},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/list definition/few items",
			args: args{"[12, 23, 42]", new(Atom)},
			wantAST: &Atom{
				ListDefinition: &ListDefinition{
					Items: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(12)}}},
									},
								},
							},
						},
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
									},
								},
							},
						},
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(42)}}},
									},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/list definition/few items/trailing comma",
			args: args{"[12, 23, 42,]", new(Atom)},
			wantAST: &Atom{
				ListDefinition: &ListDefinition{
					Items: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(12)}}},
									},
								},
							},
						},
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
									},
								},
							},
						},
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(42)}}},
									},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/function call/no arguments",
			args:    args{"test()", new(Atom)},
			wantAST: &Atom{FunctionCall: &FunctionCall{Name: "test"}},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/function call/single argument",
			args: args{"test(12)", new(Atom)},
			wantAST: &Atom{
				FunctionCall: &FunctionCall{
					Name: "test",
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(12)}}},
									},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/function call/single argument/trailing comma",
			args: args{"test(12,)", new(Atom)},
			wantAST: &Atom{
				FunctionCall: &FunctionCall{
					Name: "test",
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(12)}}},
									},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/function call/few arguments",
			args: args{"test(12, 23, 42)", new(Atom)},
			wantAST: &Atom{
				FunctionCall: &FunctionCall{
					Name: "test",
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(12)}}},
									},
								},
							},
						},
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
									},
								},
							},
						},
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(42)}}},
									},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/function call/few arguments/trailing comma",
			args: args{"test(12, 23, 42,)", new(Atom)},
			wantAST: &Atom{
				FunctionCall: &FunctionCall{
					Name: "test",
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(12)}}},
									},
								},
							},
						},
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
									},
								},
							},
						},
						{
							ListConstruction: &ListConstruction{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(42)}}},
									},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/expression",
			args: args{"(23)", new(Atom)},
			wantAST: &Atom{
				Expression: &Expression{
					ListConstruction: &ListConstruction{
						Addition: &Addition{
							Multiplication: &Multiplication{
								Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Accessor/nonempty",
			args: args{"test[12][23]", new(Accessor)},
			wantAST: &Accessor{
				Atom: &Atom{Identifier: tests.GetStringAddress("test")},
				Key: []*Expression{
					{
						ListConstruction: &ListConstruction{
							Addition: &Addition{
								Multiplication: &Multiplication{
									Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(12)}}},
								},
							},
						},
					},
					{
						ListConstruction: &ListConstruction{
							Addition: &Addition{
								Multiplication: &Multiplication{
									Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Accessor/empty",
			args:    args{"23", new(Accessor)},
			wantAST: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}},
			wantErr: assert.NoError,
		},
		{
			name: "Unary/nonempty",
			args: args{"--23", new(Unary)},
			wantAST: &Unary{
				Operation: "-",
				Unary: &Unary{
					Operation: "-",
					Unary:     &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Unary/empty",
			args:    args{"23", new(Unary)},
			wantAST: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
			wantErr: assert.NoError,
		},
		{
			name: "Multiplication/nonempty",
			args: args{"12 * 23 / 42", new(Multiplication)},
			wantAST: &Multiplication{
				Unary:     &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(12)}}},
				Operation: "*",
				Multiplication: &Multiplication{
					Unary:     &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
					Operation: "/",
					Multiplication: &Multiplication{
						Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(42)}}},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Multiplication/empty",
			args: args{"23", new(Multiplication)},
			wantAST: &Multiplication{
				Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Addition/nonempty",
			args: args{"12 + 23 - 42", new(Addition)},
			wantAST: &Addition{
				Multiplication: &Multiplication{
					Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(12)}}},
				},
				Operation: "+",
				Addition: &Addition{
					Multiplication: &Multiplication{
						Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
					},
					Operation: "-",
					Addition: &Addition{
						Multiplication: &Multiplication{
							Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(42)}}},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Addition/empty",
			args: args{"23", new(Addition)},
			wantAST: &Addition{
				Multiplication: &Multiplication{
					Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "ListConstruction/nonempty",
			args: args{"12 : [23, 42]", new(ListConstruction)},
			wantAST: &ListConstruction{
				Addition: &Addition{
					Multiplication: &Multiplication{
						Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(12)}}},
					},
				},
				ListConstruction: &ListConstruction{
					Addition: &Addition{
						Multiplication: &Multiplication{
							Unary: &Unary{
								Accessor: &Accessor{
									Atom: &Atom{
										ListDefinition: &ListDefinition{
											Items: []*Expression{
												{
													ListConstruction: &ListConstruction{
														Addition: &Addition{
															Multiplication: &Multiplication{
																Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
															},
														},
													},
												},
												{
													ListConstruction: &ListConstruction{
														Addition: &Addition{
															Multiplication: &Multiplication{
																Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(42)}}},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "ListConstruction/empty",
			args: args{"23", new(ListConstruction)},
			wantAST: &ListConstruction{
				Addition: &Addition{
					Multiplication: &Multiplication{
						Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Expression",
			args: args{"23", new(Expression)},
			wantAST: &Expression{
				ListConstruction: &ListConstruction{
					Addition: &Addition{
						Multiplication: &Multiplication{
							Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: tests.GetNumberAddress(23)}}},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			err := parseToAST(testData.args.code, testData.args.ast)
			assert.Equal(test, testData.wantAST, testData.args.ast)
			testData.wantErr(test, err)
		})
	}
}
