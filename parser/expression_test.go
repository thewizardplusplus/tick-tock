package parser

import (
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
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
			wantAST: &Atom{Number: pointer.ToFloat64(23)},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/number/floating-point",
			args:    args{"2.3", new(Atom)},
			wantAST: &Atom{Number: pointer.ToFloat64(2.3)},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/symbol/latin1",
			args:    args{"'t'", new(Atom)},
			wantAST: &Atom{Symbol: pointer.ToString("t")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/symbol/not latin1",
			args:    args{"'т'", new(Atom)},
			wantAST: &Atom{Symbol: pointer.ToString("т")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/symbol/escape sequence",
			args:    args{`'\n'`, new(Atom)},
			wantAST: &Atom{Symbol: pointer.ToString("\n")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/string/interpreted/single-quoted",
			args:    args{`'line #1\nline #2'`, new(Atom)},
			wantAST: &Atom{String: pointer.ToString("line #1\nline #2")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/string/interpreted/double-quoted",
			args:    args{`"line #1\nline #2"`, new(Atom)},
			wantAST: &Atom{String: pointer.ToString("line #1\nline #2")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/string/raw/single line",
			args:    args{"`line #1\\nline #2`", new(Atom)},
			wantAST: &Atom{String: pointer.ToString("line #1\\nline #2")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/string/raw/few lines",
			args:    args{"`line #1\nline #2`", new(Atom)},
			wantAST: &Atom{String: pointer.ToString("line #1\nline #2")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/identifier",
			args:    args{"test", new(Atom)},
			wantAST: &Atom{Identifier: pointer.ToString("test")},
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
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
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
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
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
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
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
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
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
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
												BitwiseDisjunction: &BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
														BitwiseConjunction: &BitwiseConjunction{
															Shift: &Shift{
																Addition: &Addition{
																	Multiplication: &Multiplication{
																		Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/conditional expression/single conditional case/nonempty",
			args: args{"when => 12 23 42;", new(Atom)},
			wantAST: &Atom{
				ConditionalExpression: &ConditionalExpression{
					ConditionalCases: []*ConditionalCase{
						{
							Condition: &Expression{
								ListConstruction: &ListConstruction{
									Disjunction: &Disjunction{
										Conjunction: &Conjunction{
											Equality: &Equality{
												Comparison: &Comparison{
													BitwiseDisjunction: &BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
															BitwiseConjunction: &BitwiseConjunction{
																Shift: &Shift{
																	Addition: &Addition{
																		Multiplication: &Multiplication{
																			Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
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
							Commands: []*Command{
								{
									Expression: &Expression{
										ListConstruction: &ListConstruction{
											Disjunction: &Disjunction{
												Conjunction: &Conjunction{
													Equality: &Equality{
														Comparison: &Comparison{
															BitwiseDisjunction: &BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &BitwiseConjunction{
																		Shift: &Shift{
																			Addition: &Addition{
																				Multiplication: &Multiplication{
																					Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
								{
									Expression: &Expression{
										ListConstruction: &ListConstruction{
											Disjunction: &Disjunction{
												Conjunction: &Conjunction{
													Equality: &Equality{
														Comparison: &Comparison{
															BitwiseDisjunction: &BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &BitwiseConjunction{
																		Shift: &Shift{
																			Addition: &Addition{
																				Multiplication: &Multiplication{
																					Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/conditional expression/single conditional case/empty",
			args: args{"when => 12;", new(Atom)},
			wantAST: &Atom{
				ConditionalExpression: &ConditionalExpression{
					ConditionalCases: []*ConditionalCase{
						{
							Condition: &Expression{
								ListConstruction: &ListConstruction{
									Disjunction: &Disjunction{
										Conjunction: &Conjunction{
											Equality: &Equality{
												Comparison: &Comparison{
													BitwiseDisjunction: &BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
															BitwiseConjunction: &BitwiseConjunction{
																Shift: &Shift{
																	Addition: &Addition{
																		Multiplication: &Multiplication{
																			Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
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
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/conditional expression/few conditional cases/nonempty",
			args: args{"when => 12 23 42 => 13 24 43 => 14 25 44;", new(Atom)},
			wantAST: &Atom{
				ConditionalExpression: &ConditionalExpression{
					ConditionalCases: []*ConditionalCase{
						{
							Condition: &Expression{
								ListConstruction: &ListConstruction{
									Disjunction: &Disjunction{
										Conjunction: &Conjunction{
											Equality: &Equality{
												Comparison: &Comparison{
													BitwiseDisjunction: &BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
															BitwiseConjunction: &BitwiseConjunction{
																Shift: &Shift{
																	Addition: &Addition{
																		Multiplication: &Multiplication{
																			Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
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
							Commands: []*Command{
								{
									Expression: &Expression{
										ListConstruction: &ListConstruction{
											Disjunction: &Disjunction{
												Conjunction: &Conjunction{
													Equality: &Equality{
														Comparison: &Comparison{
															BitwiseDisjunction: &BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &BitwiseConjunction{
																		Shift: &Shift{
																			Addition: &Addition{
																				Multiplication: &Multiplication{
																					Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
								{
									Expression: &Expression{
										ListConstruction: &ListConstruction{
											Disjunction: &Disjunction{
												Conjunction: &Conjunction{
													Equality: &Equality{
														Comparison: &Comparison{
															BitwiseDisjunction: &BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &BitwiseConjunction{
																		Shift: &Shift{
																			Addition: &Addition{
																				Multiplication: &Multiplication{
																					Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
							},
						},
						{
							Condition: &Expression{
								ListConstruction: &ListConstruction{
									Disjunction: &Disjunction{
										Conjunction: &Conjunction{
											Equality: &Equality{
												Comparison: &Comparison{
													BitwiseDisjunction: &BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
															BitwiseConjunction: &BitwiseConjunction{
																Shift: &Shift{
																	Addition: &Addition{
																		Multiplication: &Multiplication{
																			Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(13)}}},
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
							Commands: []*Command{
								{
									Expression: &Expression{
										ListConstruction: &ListConstruction{
											Disjunction: &Disjunction{
												Conjunction: &Conjunction{
													Equality: &Equality{
														Comparison: &Comparison{
															BitwiseDisjunction: &BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &BitwiseConjunction{
																		Shift: &Shift{
																			Addition: &Addition{
																				Multiplication: &Multiplication{
																					Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(24)}}},
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
								{
									Expression: &Expression{
										ListConstruction: &ListConstruction{
											Disjunction: &Disjunction{
												Conjunction: &Conjunction{
													Equality: &Equality{
														Comparison: &Comparison{
															BitwiseDisjunction: &BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &BitwiseConjunction{
																		Shift: &Shift{
																			Addition: &Addition{
																				Multiplication: &Multiplication{
																					Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(43)}}},
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
							},
						},
						{
							Condition: &Expression{
								ListConstruction: &ListConstruction{
									Disjunction: &Disjunction{
										Conjunction: &Conjunction{
											Equality: &Equality{
												Comparison: &Comparison{
													BitwiseDisjunction: &BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
															BitwiseConjunction: &BitwiseConjunction{
																Shift: &Shift{
																	Addition: &Addition{
																		Multiplication: &Multiplication{
																			Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(14)}}},
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
							Commands: []*Command{
								{
									Expression: &Expression{
										ListConstruction: &ListConstruction{
											Disjunction: &Disjunction{
												Conjunction: &Conjunction{
													Equality: &Equality{
														Comparison: &Comparison{
															BitwiseDisjunction: &BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &BitwiseConjunction{
																		Shift: &Shift{
																			Addition: &Addition{
																				Multiplication: &Multiplication{
																					Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(25)}}},
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
								{
									Expression: &Expression{
										ListConstruction: &ListConstruction{
											Disjunction: &Disjunction{
												Conjunction: &Conjunction{
													Equality: &Equality{
														Comparison: &Comparison{
															BitwiseDisjunction: &BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &BitwiseConjunction{
																		Shift: &Shift{
																			Addition: &Addition{
																				Multiplication: &Multiplication{
																					Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(44)}}},
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
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/conditional expression/few conditional cases/empty",
			args: args{"when => 12 => 23 => 42;", new(Atom)},
			wantAST: &Atom{
				ConditionalExpression: &ConditionalExpression{
					ConditionalCases: []*ConditionalCase{
						{
							Condition: &Expression{
								ListConstruction: &ListConstruction{
									Disjunction: &Disjunction{
										Conjunction: &Conjunction{
											Equality: &Equality{
												Comparison: &Comparison{
													BitwiseDisjunction: &BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
															BitwiseConjunction: &BitwiseConjunction{
																Shift: &Shift{
																	Addition: &Addition{
																		Multiplication: &Multiplication{
																			Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
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
						{
							Condition: &Expression{
								ListConstruction: &ListConstruction{
									Disjunction: &Disjunction{
										Conjunction: &Conjunction{
											Equality: &Equality{
												Comparison: &Comparison{
													BitwiseDisjunction: &BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
															BitwiseConjunction: &BitwiseConjunction{
																Shift: &Shift{
																	Addition: &Addition{
																		Multiplication: &Multiplication{
																			Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
						{
							Condition: &Expression{
								ListConstruction: &ListConstruction{
									Disjunction: &Disjunction{
										Conjunction: &Conjunction{
											Equality: &Equality{
												Comparison: &Comparison{
													BitwiseDisjunction: &BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
															BitwiseConjunction: &BitwiseConjunction{
																Shift: &Shift{
																	Addition: &Addition{
																		Multiplication: &Multiplication{
																			Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/conditional expression/without conditional cases",
			args:    args{"when;", new(Atom)},
			wantAST: &Atom{ConditionalExpression: &ConditionalExpression{}},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/expression",
			args: args{"(23)", new(Atom)},
			wantAST: &Atom{
				Expression: &Expression{
					ListConstruction: &ListConstruction{
						Disjunction: &Disjunction{
							Conjunction: &Conjunction{
								Equality: &Equality{
									Comparison: &Comparison{
										BitwiseDisjunction: &BitwiseDisjunction{
											BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
												BitwiseConjunction: &BitwiseConjunction{
													Shift: &Shift{
														Addition: &Addition{
															Multiplication: &Multiplication{
																Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
			name: "Accessor/nonempty",
			args: args{"test[12][23]", new(Accessor)},
			wantAST: &Accessor{
				Atom: &Atom{Identifier: pointer.ToString("test")},
				Keys: []*Expression{
					{
						ListConstruction: &ListConstruction{
							Disjunction: &Disjunction{
								Conjunction: &Conjunction{
									Equality: &Equality{
										Comparison: &Comparison{
											BitwiseDisjunction: &BitwiseDisjunction{
												BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
													BitwiseConjunction: &BitwiseConjunction{
														Shift: &Shift{
															Addition: &Addition{
																Multiplication: &Multiplication{
																	Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
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
					{
						ListConstruction: &ListConstruction{
							Disjunction: &Disjunction{
								Conjunction: &Conjunction{
									Equality: &Equality{
										Comparison: &Comparison{
											BitwiseDisjunction: &BitwiseDisjunction{
												BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
													BitwiseConjunction: &BitwiseConjunction{
														Shift: &Shift{
															Addition: &Addition{
																Multiplication: &Multiplication{
																	Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Accessor/empty",
			args:    args{"23", new(Accessor)},
			wantAST: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}},
			wantErr: assert.NoError,
		},
		{
			name: "Unary/nonempty",
			args: args{"-~!23", new(Unary)},
			wantAST: &Unary{
				Operation: "-",
				Unary: &Unary{
					Operation: "~",
					Unary: &Unary{
						Operation: "!",
						Unary:     &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Unary/empty",
			args:    args{"23", new(Unary)},
			wantAST: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
			wantErr: assert.NoError,
		},
		{
			name: "Multiplication/nonempty",
			args: args{"5 * 12 / 23 % 42", new(Multiplication)},
			wantAST: &Multiplication{
				Unary:     &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(5)}}},
				Operation: "*",
				Multiplication: &Multiplication{
					Unary:     &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
					Operation: "/",
					Multiplication: &Multiplication{
						Unary:     &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
						Operation: "%",
						Multiplication: &Multiplication{
							Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Multiplication/empty",
			args: args{"23", new(Multiplication)},
			wantAST: &Multiplication{
				Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Addition/nonempty",
			args: args{"12 + 23 - 42", new(Addition)},
			wantAST: &Addition{
				Multiplication: &Multiplication{
					Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
				},
				Operation: "+",
				Addition: &Addition{
					Multiplication: &Multiplication{
						Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
					},
					Operation: "-",
					Addition: &Addition{
						Multiplication: &Multiplication{
							Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
					Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Shift/nonempty",
			args: args{"5 << 12 >> 23 >>> 42", new(Shift)},
			wantAST: &Shift{
				Addition: &Addition{
					Multiplication: &Multiplication{
						Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(5)}}},
					},
				},
				Operation: "<<",
				Shift: &Shift{
					Addition: &Addition{
						Multiplication: &Multiplication{
							Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
						},
					},
					Operation: ">>",
					Shift: &Shift{
						Addition: &Addition{
							Multiplication: &Multiplication{
								Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
							},
						},
						Operation: ">>>",
						Shift: &Shift{
							Addition: &Addition{
								Multiplication: &Multiplication{
									Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Shift/empty",
			args: args{"23", new(Shift)},
			wantAST: &Shift{
				Addition: &Addition{
					Multiplication: &Multiplication{
						Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "BitwiseConjunction/nonempty",
			args: args{"12 & 23 & 42", new(BitwiseConjunction)},
			wantAST: &BitwiseConjunction{
				Shift: &Shift{
					Addition: &Addition{
						Multiplication: &Multiplication{
							Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
						},
					},
				},
				BitwiseConjunction: &BitwiseConjunction{
					Shift: &Shift{
						Addition: &Addition{
							Multiplication: &Multiplication{
								Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
							},
						},
					},
					BitwiseConjunction: &BitwiseConjunction{
						Shift: &Shift{
							Addition: &Addition{
								Multiplication: &Multiplication{
									Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "BitwiseConjunction/empty",
			args: args{"23", new(BitwiseConjunction)},
			wantAST: &BitwiseConjunction{
				Shift: &Shift{
					Addition: &Addition{
						Multiplication: &Multiplication{
							Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "BitwiseExclusiveDisjunction/nonempty",
			args: args{"12 ^ 23 ^ 42", new(BitwiseExclusiveDisjunction)},
			wantAST: &BitwiseExclusiveDisjunction{
				BitwiseConjunction: &BitwiseConjunction{
					Shift: &Shift{
						Addition: &Addition{
							Multiplication: &Multiplication{
								Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
							},
						},
					},
				},
				BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
					BitwiseConjunction: &BitwiseConjunction{
						Shift: &Shift{
							Addition: &Addition{
								Multiplication: &Multiplication{
									Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
								},
							},
						},
					},
					BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
						BitwiseConjunction: &BitwiseConjunction{
							Shift: &Shift{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
			name: "BitwiseExclusiveDisjunction/empty",
			args: args{"23", new(BitwiseExclusiveDisjunction)},
			wantAST: &BitwiseExclusiveDisjunction{
				BitwiseConjunction: &BitwiseConjunction{
					Shift: &Shift{
						Addition: &Addition{
							Multiplication: &Multiplication{
								Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "BitwiseDisjunction/nonempty",
			args: args{"12 | 23 | 42", new(BitwiseDisjunction)},
			wantAST: &BitwiseDisjunction{
				BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
					BitwiseConjunction: &BitwiseConjunction{
						Shift: &Shift{
							Addition: &Addition{
								Multiplication: &Multiplication{
									Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
								},
							},
						},
					},
				},
				BitwiseDisjunction: &BitwiseDisjunction{
					BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
						BitwiseConjunction: &BitwiseConjunction{
							Shift: &Shift{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
									},
								},
							},
						},
					},
					BitwiseDisjunction: &BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
							BitwiseConjunction: &BitwiseConjunction{
								Shift: &Shift{
									Addition: &Addition{
										Multiplication: &Multiplication{
											Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
			name: "BitwiseDisjunction/empty",
			args: args{"23", new(BitwiseDisjunction)},
			wantAST: &BitwiseDisjunction{
				BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
					BitwiseConjunction: &BitwiseConjunction{
						Shift: &Shift{
							Addition: &Addition{
								Multiplication: &Multiplication{
									Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Comparison/nonempty/less",
			args: args{"12 < 23 <= 42", new(Comparison)},
			wantAST: &Comparison{
				BitwiseDisjunction: &BitwiseDisjunction{
					BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
						BitwiseConjunction: &BitwiseConjunction{
							Shift: &Shift{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
									},
								},
							},
						},
					},
				},
				Operation: "<",
				Comparison: &Comparison{
					BitwiseDisjunction: &BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
							BitwiseConjunction: &BitwiseConjunction{
								Shift: &Shift{
									Addition: &Addition{
										Multiplication: &Multiplication{
											Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
										},
									},
								},
							},
						},
					},
					Operation: "<=",
					Comparison: &Comparison{
						BitwiseDisjunction: &BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
								BitwiseConjunction: &BitwiseConjunction{
									Shift: &Shift{
										Addition: &Addition{
											Multiplication: &Multiplication{
												Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
			name: "Comparison/nonempty/great",
			args: args{"12 > 23 >= 42", new(Comparison)},
			wantAST: &Comparison{
				BitwiseDisjunction: &BitwiseDisjunction{
					BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
						BitwiseConjunction: &BitwiseConjunction{
							Shift: &Shift{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
									},
								},
							},
						},
					},
				},
				Operation: ">",
				Comparison: &Comparison{
					BitwiseDisjunction: &BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
							BitwiseConjunction: &BitwiseConjunction{
								Shift: &Shift{
									Addition: &Addition{
										Multiplication: &Multiplication{
											Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
										},
									},
								},
							},
						},
					},
					Operation: ">=",
					Comparison: &Comparison{
						BitwiseDisjunction: &BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
								BitwiseConjunction: &BitwiseConjunction{
									Shift: &Shift{
										Addition: &Addition{
											Multiplication: &Multiplication{
												Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
			name: "Comparison/empty",
			args: args{"23", new(Comparison)},
			wantAST: &Comparison{
				BitwiseDisjunction: &BitwiseDisjunction{
					BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
						BitwiseConjunction: &BitwiseConjunction{
							Shift: &Shift{
								Addition: &Addition{
									Multiplication: &Multiplication{
										Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
			name: "Equality/nonempty",
			args: args{"12 == 23 != 42", new(Equality)},
			wantAST: &Equality{
				Comparison: &Comparison{
					BitwiseDisjunction: &BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
							BitwiseConjunction: &BitwiseConjunction{
								Shift: &Shift{
									Addition: &Addition{
										Multiplication: &Multiplication{
											Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
										},
									},
								},
							},
						},
					},
				},
				Operation: "==",
				Equality: &Equality{
					Comparison: &Comparison{
						BitwiseDisjunction: &BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
								BitwiseConjunction: &BitwiseConjunction{
									Shift: &Shift{
										Addition: &Addition{
											Multiplication: &Multiplication{
												Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
											},
										},
									},
								},
							},
						},
					},
					Operation: "!=",
					Equality: &Equality{
						Comparison: &Comparison{
							BitwiseDisjunction: &BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
									BitwiseConjunction: &BitwiseConjunction{
										Shift: &Shift{
											Addition: &Addition{
												Multiplication: &Multiplication{
													Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
			name: "Equality/empty",
			args: args{"23", new(Equality)},
			wantAST: &Equality{
				Comparison: &Comparison{
					BitwiseDisjunction: &BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
							BitwiseConjunction: &BitwiseConjunction{
								Shift: &Shift{
									Addition: &Addition{
										Multiplication: &Multiplication{
											Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
			name: "Conjunction/nonempty",
			args: args{"12 && 23 && 42", new(Conjunction)},
			wantAST: &Conjunction{
				Equality: &Equality{
					Comparison: &Comparison{
						BitwiseDisjunction: &BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
								BitwiseConjunction: &BitwiseConjunction{
									Shift: &Shift{
										Addition: &Addition{
											Multiplication: &Multiplication{
												Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
											},
										},
									},
								},
							},
						},
					},
				},
				Conjunction: &Conjunction{
					Equality: &Equality{
						Comparison: &Comparison{
							BitwiseDisjunction: &BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
									BitwiseConjunction: &BitwiseConjunction{
										Shift: &Shift{
											Addition: &Addition{
												Multiplication: &Multiplication{
													Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
												},
											},
										},
									},
								},
							},
						},
					},
					Conjunction: &Conjunction{
						Equality: &Equality{
							Comparison: &Comparison{
								BitwiseDisjunction: &BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
										BitwiseConjunction: &BitwiseConjunction{
											Shift: &Shift{
												Addition: &Addition{
													Multiplication: &Multiplication{
														Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
			name: "Conjunction/empty",
			args: args{"23", new(Conjunction)},
			wantAST: &Conjunction{
				Equality: &Equality{
					Comparison: &Comparison{
						BitwiseDisjunction: &BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
								BitwiseConjunction: &BitwiseConjunction{
									Shift: &Shift{
										Addition: &Addition{
											Multiplication: &Multiplication{
												Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
			name: "Disjunction/nonempty",
			args: args{"12 || 23 || 42", new(Disjunction)},
			wantAST: &Disjunction{
				Conjunction: &Conjunction{
					Equality: &Equality{
						Comparison: &Comparison{
							BitwiseDisjunction: &BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
									BitwiseConjunction: &BitwiseConjunction{
										Shift: &Shift{
											Addition: &Addition{
												Multiplication: &Multiplication{
													Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				Disjunction: &Disjunction{
					Conjunction: &Conjunction{
						Equality: &Equality{
							Comparison: &Comparison{
								BitwiseDisjunction: &BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
										BitwiseConjunction: &BitwiseConjunction{
											Shift: &Shift{
												Addition: &Addition{
													Multiplication: &Multiplication{
														Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					Disjunction: &Disjunction{
						Conjunction: &Conjunction{
							Equality: &Equality{
								Comparison: &Comparison{
									BitwiseDisjunction: &BitwiseDisjunction{
										BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
											BitwiseConjunction: &BitwiseConjunction{
												Shift: &Shift{
													Addition: &Addition{
														Multiplication: &Multiplication{
															Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}}},
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
			name: "Disjunction/empty",
			args: args{"23", new(Disjunction)},
			wantAST: &Disjunction{
				Conjunction: &Conjunction{
					Equality: &Equality{
						Comparison: &Comparison{
							BitwiseDisjunction: &BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
									BitwiseConjunction: &BitwiseConjunction{
										Shift: &Shift{
											Addition: &Addition{
												Multiplication: &Multiplication{
													Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
			name: "ListConstruction/nonempty",
			args: args{"12 : [23, 42]", new(ListConstruction)},
			wantAST: &ListConstruction{
				Disjunction: &Disjunction{
					Conjunction: &Conjunction{
						Equality: &Equality{
							Comparison: &Comparison{
								BitwiseDisjunction: &BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
										BitwiseConjunction: &BitwiseConjunction{
											Shift: &Shift{
												Addition: &Addition{
													Multiplication: &Multiplication{
														Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(12)}}},
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
				ListConstruction: &ListConstruction{
					Disjunction: &Disjunction{
						Conjunction: &Conjunction{
							Equality: &Equality{
								Comparison: &Comparison{
									BitwiseDisjunction: &BitwiseDisjunction{
										BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
											BitwiseConjunction: &BitwiseConjunction{
												Shift: &Shift{
													Addition: &Addition{
														Multiplication: &Multiplication{
															Unary: &Unary{
																Accessor: &Accessor{
																	Atom: &Atom{
																		ListDefinition: &ListDefinition{
																			Items: []*Expression{
																				{
																					ListConstruction: &ListConstruction{
																						Disjunction: &Disjunction{
																							Conjunction: &Conjunction{
																								Equality: &Equality{
																									Comparison: &Comparison{
																										BitwiseDisjunction: &BitwiseDisjunction{
																											BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
																												BitwiseConjunction: &BitwiseConjunction{
																													Shift: &Shift{
																														Addition: &Addition{
																															Multiplication: &Multiplication{
																																Unary: &Unary{
																																	Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}},
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
																				{
																					ListConstruction: &ListConstruction{
																						Disjunction: &Disjunction{
																							Conjunction: &Conjunction{
																								Equality: &Equality{
																									Comparison: &Comparison{
																										BitwiseDisjunction: &BitwiseDisjunction{
																											BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
																												BitwiseConjunction: &BitwiseConjunction{
																													Shift: &Shift{
																														Addition: &Addition{
																															Multiplication: &Multiplication{
																																Unary: &Unary{
																																	Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(42)}},
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
				Disjunction: &Disjunction{
					Conjunction: &Conjunction{
						Equality: &Equality{
							Comparison: &Comparison{
								BitwiseDisjunction: &BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
										BitwiseConjunction: &BitwiseConjunction{
											Shift: &Shift{
												Addition: &Addition{
													Multiplication: &Multiplication{
														Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
			name: "Expression",
			args: args{"23", new(Expression)},
			wantAST: &Expression{
				ListConstruction: &ListConstruction{
					Disjunction: &Disjunction{
						Conjunction: &Conjunction{
							Equality: &Equality{
								Comparison: &Comparison{
									BitwiseDisjunction: &BitwiseDisjunction{
										BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
											BitwiseConjunction: &BitwiseConjunction{
												Shift: &Shift{
													Addition: &Addition{
														Multiplication: &Multiplication{
															Unary: &Unary{Accessor: &Accessor{Atom: &Atom{Number: pointer.ToFloat64(23)}}},
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
	} {
		test.Run(testData.name, func(test *testing.T) {
			err := parseToAST(testData.args.code, testData.args.ast)
			assert.Equal(test, testData.wantAST, testData.args.ast)
			testData.wantErr(test, err)
		})
	}
}
