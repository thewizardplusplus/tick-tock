package parser

import (
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestParseToAST_withCommon(test *testing.T) {
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
			name:    "ExpressionGroup/no items",
			args:    args{"", new(ExpressionGroup)},
			wantAST: &ExpressionGroup{},
			wantErr: assert.NoError,
		},
		{
			name: "ExpressionGroup/single item",
			args: args{"12", new(ExpressionGroup)},
			wantAST: &ExpressionGroup{
				Expressions: []*Expression{
					{
						ListConstruction: &ListConstruction{
							NilCoalescing: &NilCoalescing{
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
																			Accessor: &Accessor{Atom: &Atom{IntegerNumber: pointer.ToInt64(12)}},
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
			name: "ExpressionGroup/single item/trailing comma",
			args: args{"12,", new(ExpressionGroup)},
			wantAST: &ExpressionGroup{
				Expressions: []*Expression{
					{
						ListConstruction: &ListConstruction{
							NilCoalescing: &NilCoalescing{
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
																			Accessor: &Accessor{Atom: &Atom{IntegerNumber: pointer.ToInt64(12)}},
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
			name: "ExpressionGroup/few items",
			args: args{"12, 23, 42", new(ExpressionGroup)},
			wantAST: &ExpressionGroup{
				Expressions: []*Expression{
					{
						ListConstruction: &ListConstruction{
							NilCoalescing: &NilCoalescing{
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
																			Accessor: &Accessor{Atom: &Atom{IntegerNumber: pointer.ToInt64(12)}},
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
						ListConstruction: &ListConstruction{
							NilCoalescing: &NilCoalescing{
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
																			Accessor: &Accessor{Atom: &Atom{IntegerNumber: pointer.ToInt64(23)}},
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
						ListConstruction: &ListConstruction{
							NilCoalescing: &NilCoalescing{
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
																			Accessor: &Accessor{Atom: &Atom{IntegerNumber: pointer.ToInt64(42)}},
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
			name: "ExpressionGroup/few items/trailing comma",
			args: args{"12, 23, 42,", new(ExpressionGroup)},
			wantAST: &ExpressionGroup{
				Expressions: []*Expression{
					{
						ListConstruction: &ListConstruction{
							NilCoalescing: &NilCoalescing{
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
																			Accessor: &Accessor{Atom: &Atom{IntegerNumber: pointer.ToInt64(12)}},
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
						ListConstruction: &ListConstruction{
							NilCoalescing: &NilCoalescing{
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
																			Accessor: &Accessor{Atom: &Atom{IntegerNumber: pointer.ToInt64(23)}},
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
						ListConstruction: &ListConstruction{
							NilCoalescing: &NilCoalescing{
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
																			Accessor: &Accessor{Atom: &Atom{IntegerNumber: pointer.ToInt64(42)}},
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
	} {
		test.Run(testData.name, func(test *testing.T) {
			err := ParseToAST(testData.args.code, testData.args.ast)

			assert.Equal(test, testData.wantAST, testData.args.ast)
			testData.wantErr(test, err)
		})
	}
}
