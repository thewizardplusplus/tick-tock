package translator

import (
	"testing"

	"github.com/AlekSi/pointer"
	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/commands"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

func TestTranslateExpression(test *testing.T) {
	type args struct {
		expression          *parser.Expression
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
				expression: &parser.Expression{
					ListConstruction: &parser.ListConstruction{
						Disjunction: &parser.Disjunction{
							Conjunction: &parser.Conjunction{
								Equality: &parser.Equality{
									Comparison: &parser.Comparison{
										BitwiseDisjunction: &parser.BitwiseDisjunction{
											BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
												BitwiseConjunction: &parser.BitwiseConjunction{
													Shift: &parser.Shift{
														Addition: &parser.Addition{
															Multiplication: &parser.Multiplication{
																Unary: &parser.Unary{
																	Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Expression/success/with setted states",
			args: args{
				expression: &parser.Expression{
					ListConstruction: &parser.ListConstruction{
						Disjunction: &parser.Disjunction{
							Conjunction: &parser.Conjunction{
								Equality: &parser.Equality{
									Comparison: &parser.Comparison{
										BitwiseDisjunction: &parser.BitwiseDisjunction{
											BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
												BitwiseConjunction: &parser.BitwiseConjunction{
													Shift: &parser.Shift{
														Addition: &parser.Addition{
															Multiplication: &parser.Multiplication{
																Unary: &parser.Unary{
																	Accessor: &parser.Accessor{
																		Atom: &parser.Atom{
																			ConditionalExpression: &parser.ConditionalExpression{
																				ConditionalCases: []*parser.ConditionalCase{
																					{
																						Condition: &parser.Expression{
																							ListConstruction: &parser.ListConstruction{
																								Disjunction: &parser.Disjunction{
																									Conjunction: &parser.Conjunction{
																										Equality: &parser.Equality{
																											Comparison: &parser.Comparison{
																												BitwiseDisjunction: &parser.BitwiseDisjunction{
																													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																														BitwiseConjunction: &parser.BitwiseConjunction{
																															Shift: &parser.Shift{
																																Addition: &parser.Addition{
																																	Multiplication: &parser.Multiplication{
																																		Unary: &parser.Unary{
																																			Accessor: &parser.Accessor{
																																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																						Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																					},
																					{
																						Condition: &parser.Expression{
																							ListConstruction: &parser.ListConstruction{
																								Disjunction: &parser.Disjunction{
																									Conjunction: &parser.Conjunction{
																										Equality: &parser.Equality{
																											Comparison: &parser.Comparison{
																												BitwiseDisjunction: &parser.BitwiseDisjunction{
																													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																														BitwiseConjunction: &parser.BitwiseConjunction{
																															Shift: &parser.Shift{
																																Addition: &parser.Addition{
																																	Multiplication: &parser.Multiplication{
																																		Unary: &parser.Unary{
																																			Accessor: &parser.Accessor{
																																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																						Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
				expression: &parser.Expression{
					ListConstruction: &parser.ListConstruction{
						Disjunction: &parser.Disjunction{
							Conjunction: &parser.Conjunction{
								Equality: &parser.Equality{
									Comparison: &parser.Comparison{
										BitwiseDisjunction: &parser.BitwiseDisjunction{
											BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
												BitwiseConjunction: &parser.BitwiseConjunction{
													Shift: &parser.Shift{
														Addition: &parser.Addition{
															Multiplication: &parser.Multiplication{
																Unary: &parser.Unary{
																	Accessor: &parser.Accessor{
																		Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateExpression(data.args.expression, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateListConstruction(test *testing.T) {
	type args struct {
		listConstruction    *parser.ListConstruction
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
				listConstruction: &parser.ListConstruction{
					Disjunction: &parser.Disjunction{
						Conjunction: &parser.Conjunction{
							Equality: &parser.Equality{
								Comparison: &parser.Comparison{
									BitwiseDisjunction: &parser.BitwiseDisjunction{
										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
											BitwiseConjunction: &parser.BitwiseConjunction{
												Shift: &parser.Shift{
													Addition: &parser.Addition{
														Multiplication: &parser.Multiplication{
															Unary: &parser.Unary{
																Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
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
					ListConstruction: &parser.ListConstruction{
						Disjunction: &parser.Disjunction{
							Conjunction: &parser.Conjunction{
								Equality: &parser.Equality{
									Comparison: &parser.Comparison{
										BitwiseDisjunction: &parser.BitwiseDisjunction{
											BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
												BitwiseConjunction: &parser.BitwiseConjunction{
													Shift: &parser.Shift{
														Addition: &parser.Addition{
															Multiplication: &parser.Multiplication{
																Unary: &parser.Unary{
																	Accessor: &parser.Accessor{
																		Atom: &parser.Atom{Identifier: pointer.ToString("test")},
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
				listConstruction: &parser.ListConstruction{
					Disjunction: &parser.Disjunction{
						Conjunction: &parser.Conjunction{
							Equality: &parser.Equality{
								Comparison: &parser.Comparison{
									BitwiseDisjunction: &parser.BitwiseDisjunction{
										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
											BitwiseConjunction: &parser.BitwiseConjunction{
												Shift: &parser.Shift{
													Addition: &parser.Addition{
														Multiplication: &parser.Multiplication{
															Unary: &parser.Unary{
																Accessor: &parser.Accessor{
																	Atom: &parser.Atom{
																		ConditionalExpression: &parser.ConditionalExpression{
																			ConditionalCases: []*parser.ConditionalCase{
																				{
																					Condition: &parser.Expression{
																						ListConstruction: &parser.ListConstruction{
																							Disjunction: &parser.Disjunction{
																								Conjunction: &parser.Conjunction{
																									Equality: &parser.Equality{
																										Comparison: &parser.Comparison{
																											BitwiseDisjunction: &parser.BitwiseDisjunction{
																												BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																													BitwiseConjunction: &parser.BitwiseConjunction{
																														Shift: &parser.Shift{
																															Addition: &parser.Addition{
																																Multiplication: &parser.Multiplication{
																																	Unary: &parser.Unary{
																																		Accessor: &parser.Accessor{
																																			Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																					Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																				},
																				{
																					Condition: &parser.Expression{
																						ListConstruction: &parser.ListConstruction{
																							Disjunction: &parser.Disjunction{
																								Conjunction: &parser.Conjunction{
																									Equality: &parser.Equality{
																										Comparison: &parser.Comparison{
																											BitwiseDisjunction: &parser.BitwiseDisjunction{
																												BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																													BitwiseConjunction: &parser.BitwiseConjunction{
																														Shift: &parser.Shift{
																															Addition: &parser.Addition{
																																Multiplication: &parser.Multiplication{
																																	Unary: &parser.Unary{
																																		Accessor: &parser.Accessor{
																																			Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																					Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
					ListConstruction: &parser.ListConstruction{
						Disjunction: &parser.Disjunction{
							Conjunction: &parser.Conjunction{
								Equality: &parser.Equality{
									Comparison: &parser.Comparison{
										BitwiseDisjunction: &parser.BitwiseDisjunction{
											BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
												BitwiseConjunction: &parser.BitwiseConjunction{
													Shift: &parser.Shift{
														Addition: &parser.Addition{
															Multiplication: &parser.Multiplication{
																Unary: &parser.Unary{
																	Accessor: &parser.Accessor{
																		Atom: &parser.Atom{
																			ConditionalExpression: &parser.ConditionalExpression{
																				ConditionalCases: []*parser.ConditionalCase{
																					{
																						Condition: &parser.Expression{
																							ListConstruction: &parser.ListConstruction{
																								Disjunction: &parser.Disjunction{
																									Conjunction: &parser.Conjunction{
																										Equality: &parser.Equality{
																											Comparison: &parser.Comparison{
																												BitwiseDisjunction: &parser.BitwiseDisjunction{
																													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																														BitwiseConjunction: &parser.BitwiseConjunction{
																															Shift: &parser.Shift{
																																Addition: &parser.Addition{
																																	Multiplication: &parser.Multiplication{
																																		Unary: &parser.Unary{
																																			Accessor: &parser.Accessor{
																																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
																						Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																					},
																					{
																						Condition: &parser.Expression{
																							ListConstruction: &parser.ListConstruction{
																								Disjunction: &parser.Disjunction{
																									Conjunction: &parser.Conjunction{
																										Equality: &parser.Equality{
																											Comparison: &parser.Comparison{
																												BitwiseDisjunction: &parser.BitwiseDisjunction{
																													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																														BitwiseConjunction: &parser.BitwiseConjunction{
																															Shift: &parser.Shift{
																																Addition: &parser.Addition{
																																	Multiplication: &parser.Multiplication{
																																		Unary: &parser.Unary{
																																			Accessor: &parser.Accessor{
																																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
																						Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				listConstruction: &parser.ListConstruction{
					Disjunction: &parser.Disjunction{
						Conjunction: &parser.Conjunction{
							Equality: &parser.Equality{
								Comparison: &parser.Comparison{
									BitwiseDisjunction: &parser.BitwiseDisjunction{
										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
											BitwiseConjunction: &parser.BitwiseConjunction{
												Shift: &parser.Shift{
													Addition: &parser.Addition{
														Multiplication: &parser.Multiplication{
															Unary: &parser.Unary{
																Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
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
					ListConstruction: &parser.ListConstruction{
						Disjunction: &parser.Disjunction{
							Conjunction: &parser.Conjunction{
								Equality: &parser.Equality{
									Comparison: &parser.Comparison{
										BitwiseDisjunction: &parser.BitwiseDisjunction{
											BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
												BitwiseConjunction: &parser.BitwiseConjunction{
													Shift: &parser.Shift{
														Addition: &parser.Addition{
															Multiplication: &parser.Multiplication{
																Unary: &parser.Unary{
																	Accessor: &parser.Accessor{
																		Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "ListConstruction/empty/success",
			args: args{
				listConstruction: &parser.ListConstruction{
					Disjunction: &parser.Disjunction{
						Conjunction: &parser.Conjunction{
							Equality: &parser.Equality{
								Comparison: &parser.Comparison{
									BitwiseDisjunction: &parser.BitwiseDisjunction{
										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
											BitwiseConjunction: &parser.BitwiseConjunction{
												Shift: &parser.Shift{
													Addition: &parser.Addition{
														Multiplication: &parser.Multiplication{
															Unary: &parser.Unary{
																Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "ListConstruction/empty/success/with setted states",
			args: args{
				listConstruction: &parser.ListConstruction{
					Disjunction: &parser.Disjunction{
						Conjunction: &parser.Conjunction{
							Equality: &parser.Equality{
								Comparison: &parser.Comparison{
									BitwiseDisjunction: &parser.BitwiseDisjunction{
										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
											BitwiseConjunction: &parser.BitwiseConjunction{
												Shift: &parser.Shift{
													Addition: &parser.Addition{
														Multiplication: &parser.Multiplication{
															Unary: &parser.Unary{
																Accessor: &parser.Accessor{
																	Atom: &parser.Atom{
																		ConditionalExpression: &parser.ConditionalExpression{
																			ConditionalCases: []*parser.ConditionalCase{
																				{
																					Condition: &parser.Expression{
																						ListConstruction: &parser.ListConstruction{
																							Disjunction: &parser.Disjunction{
																								Conjunction: &parser.Conjunction{
																									Equality: &parser.Equality{
																										Comparison: &parser.Comparison{
																											BitwiseDisjunction: &parser.BitwiseDisjunction{
																												BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																													BitwiseConjunction: &parser.BitwiseConjunction{
																														Shift: &parser.Shift{
																															Addition: &parser.Addition{
																																Multiplication: &parser.Multiplication{
																																	Unary: &parser.Unary{
																																		Accessor: &parser.Accessor{
																																			Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																					Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																				},
																				{
																					Condition: &parser.Expression{
																						ListConstruction: &parser.ListConstruction{
																							Disjunction: &parser.Disjunction{
																								Conjunction: &parser.Conjunction{
																									Equality: &parser.Equality{
																										Comparison: &parser.Comparison{
																											BitwiseDisjunction: &parser.BitwiseDisjunction{
																												BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																													BitwiseConjunction: &parser.BitwiseConjunction{
																														Shift: &parser.Shift{
																															Addition: &parser.Addition{
																																Multiplication: &parser.Multiplication{
																																	Unary: &parser.Unary{
																																		Accessor: &parser.Accessor{
																																			Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																					Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
				listConstruction: &parser.ListConstruction{
					Disjunction: &parser.Disjunction{
						Conjunction: &parser.Conjunction{
							Equality: &parser.Equality{
								Comparison: &parser.Comparison{
									BitwiseDisjunction: &parser.BitwiseDisjunction{
										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
											BitwiseConjunction: &parser.BitwiseConjunction{
												Shift: &parser.Shift{
													Addition: &parser.Addition{
														Multiplication: &parser.Multiplication{
															Unary: &parser.Unary{
																Accessor: &parser.Accessor{
																	Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateListConstruction(data.args.listConstruction, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateDisjunction(test *testing.T) {
	type args struct {
		disjunction         *parser.Disjunction
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
				disjunction: &parser.Disjunction{
					Conjunction: &parser.Conjunction{
						Equality: &parser.Equality{
							Comparison: &parser.Comparison{
								BitwiseDisjunction: &parser.BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
										BitwiseConjunction: &parser.BitwiseConjunction{
											Shift: &parser.Shift{
												Addition: &parser.Addition{
													Multiplication: &parser.Multiplication{
														Unary: &parser.Unary{
															Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
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
					Disjunction: &parser.Disjunction{
						Conjunction: &parser.Conjunction{
							Equality: &parser.Equality{
								Comparison: &parser.Comparison{
									BitwiseDisjunction: &parser.BitwiseDisjunction{
										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
											BitwiseConjunction: &parser.BitwiseConjunction{
												Shift: &parser.Shift{
													Addition: &parser.Addition{
														Multiplication: &parser.Multiplication{
															Unary: &parser.Unary{
																Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
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
						Disjunction: &parser.Disjunction{
							Conjunction: &parser.Conjunction{
								Equality: &parser.Equality{
									Comparison: &parser.Comparison{
										BitwiseDisjunction: &parser.BitwiseDisjunction{
											BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
												BitwiseConjunction: &parser.BitwiseConjunction{
													Shift: &parser.Shift{
														Addition: &parser.Addition{
															Multiplication: &parser.Multiplication{
																Unary: &parser.Unary{
																	Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
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
				disjunction: &parser.Disjunction{
					Conjunction: &parser.Conjunction{
						Equality: &parser.Equality{
							Comparison: &parser.Comparison{
								BitwiseDisjunction: &parser.BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
										BitwiseConjunction: &parser.BitwiseConjunction{
											Shift: &parser.Shift{
												Addition: &parser.Addition{
													Multiplication: &parser.Multiplication{
														Unary: &parser.Unary{
															Accessor: &parser.Accessor{
																Atom: &parser.Atom{
																	ConditionalExpression: &parser.ConditionalExpression{
																		ConditionalCases: []*parser.ConditionalCase{
																			{
																				Condition: &parser.Expression{
																					ListConstruction: &parser.ListConstruction{
																						Disjunction: &parser.Disjunction{
																							Conjunction: &parser.Conjunction{
																								Equality: &parser.Equality{
																									Comparison: &parser.Comparison{
																										BitwiseDisjunction: &parser.BitwiseDisjunction{
																											BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																												BitwiseConjunction: &parser.BitwiseConjunction{
																													Shift: &parser.Shift{
																														Addition: &parser.Addition{
																															Multiplication: &parser.Multiplication{
																																Unary: &parser.Unary{
																																	Accessor: &parser.Accessor{
																																		Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																				Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																			},
																			{
																				Condition: &parser.Expression{
																					ListConstruction: &parser.ListConstruction{
																						Disjunction: &parser.Disjunction{
																							Conjunction: &parser.Conjunction{
																								Equality: &parser.Equality{
																									Comparison: &parser.Comparison{
																										BitwiseDisjunction: &parser.BitwiseDisjunction{
																											BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																												BitwiseConjunction: &parser.BitwiseConjunction{
																													Shift: &parser.Shift{
																														Addition: &parser.Addition{
																															Multiplication: &parser.Multiplication{
																																Unary: &parser.Unary{
																																	Accessor: &parser.Accessor{
																																		Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																				Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
					Disjunction: &parser.Disjunction{
						Conjunction: &parser.Conjunction{
							Equality: &parser.Equality{
								Comparison: &parser.Comparison{
									BitwiseDisjunction: &parser.BitwiseDisjunction{
										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
											BitwiseConjunction: &parser.BitwiseConjunction{
												Shift: &parser.Shift{
													Addition: &parser.Addition{
														Multiplication: &parser.Multiplication{
															Unary: &parser.Unary{
																Accessor: &parser.Accessor{
																	Atom: &parser.Atom{
																		ConditionalExpression: &parser.ConditionalExpression{
																			ConditionalCases: []*parser.ConditionalCase{
																				{
																					Condition: &parser.Expression{
																						ListConstruction: &parser.ListConstruction{
																							Disjunction: &parser.Disjunction{
																								Conjunction: &parser.Conjunction{
																									Equality: &parser.Equality{
																										Comparison: &parser.Comparison{
																											BitwiseDisjunction: &parser.BitwiseDisjunction{
																												BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																													BitwiseConjunction: &parser.BitwiseConjunction{
																														Shift: &parser.Shift{
																															Addition: &parser.Addition{
																																Multiplication: &parser.Multiplication{
																																	Unary: &parser.Unary{
																																		Accessor: &parser.Accessor{
																																			Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
																					Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																				},
																				{
																					Condition: &parser.Expression{
																						ListConstruction: &parser.ListConstruction{
																							Disjunction: &parser.Disjunction{
																								Conjunction: &parser.Conjunction{
																									Equality: &parser.Equality{
																										Comparison: &parser.Comparison{
																											BitwiseDisjunction: &parser.BitwiseDisjunction{
																												BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																													BitwiseConjunction: &parser.BitwiseConjunction{
																														Shift: &parser.Shift{
																															Addition: &parser.Addition{
																																Multiplication: &parser.Multiplication{
																																	Unary: &parser.Unary{
																																		Accessor: &parser.Accessor{
																																			Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
																					Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				disjunction: &parser.Disjunction{
					Conjunction: &parser.Conjunction{
						Equality: &parser.Equality{
							Comparison: &parser.Comparison{
								BitwiseDisjunction: &parser.BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
										BitwiseConjunction: &parser.BitwiseConjunction{
											Shift: &parser.Shift{
												Addition: &parser.Addition{
													Multiplication: &parser.Multiplication{
														Unary: &parser.Unary{
															Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
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
					Disjunction: &parser.Disjunction{
						Conjunction: &parser.Conjunction{
							Equality: &parser.Equality{
								Comparison: &parser.Comparison{
									BitwiseDisjunction: &parser.BitwiseDisjunction{
										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
											BitwiseConjunction: &parser.BitwiseConjunction{
												Shift: &parser.Shift{
													Addition: &parser.Addition{
														Multiplication: &parser.Multiplication{
															Unary: &parser.Unary{
																Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
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
						Disjunction: &parser.Disjunction{
							Conjunction: &parser.Conjunction{
								Equality: &parser.Equality{
									Comparison: &parser.Comparison{
										BitwiseDisjunction: &parser.BitwiseDisjunction{
											BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
												BitwiseConjunction: &parser.BitwiseConjunction{
													Shift: &parser.Shift{
														Addition: &parser.Addition{
															Multiplication: &parser.Multiplication{
																Unary: &parser.Unary{
																	Accessor: &parser.Accessor{
																		Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Disjunction/empty/success",
			args: args{
				disjunction: &parser.Disjunction{
					Conjunction: &parser.Conjunction{
						Equality: &parser.Equality{
							Comparison: &parser.Comparison{
								BitwiseDisjunction: &parser.BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
										BitwiseConjunction: &parser.BitwiseConjunction{
											Shift: &parser.Shift{
												Addition: &parser.Addition{
													Multiplication: &parser.Multiplication{
														Unary: &parser.Unary{
															Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Disjunction/empty/success/with setted states",
			args: args{
				disjunction: &parser.Disjunction{
					Conjunction: &parser.Conjunction{
						Equality: &parser.Equality{
							Comparison: &parser.Comparison{
								BitwiseDisjunction: &parser.BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
										BitwiseConjunction: &parser.BitwiseConjunction{
											Shift: &parser.Shift{
												Addition: &parser.Addition{
													Multiplication: &parser.Multiplication{
														Unary: &parser.Unary{
															Accessor: &parser.Accessor{
																Atom: &parser.Atom{
																	ConditionalExpression: &parser.ConditionalExpression{
																		ConditionalCases: []*parser.ConditionalCase{
																			{
																				Condition: &parser.Expression{
																					ListConstruction: &parser.ListConstruction{
																						Disjunction: &parser.Disjunction{
																							Conjunction: &parser.Conjunction{
																								Equality: &parser.Equality{
																									Comparison: &parser.Comparison{
																										BitwiseDisjunction: &parser.BitwiseDisjunction{
																											BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																												BitwiseConjunction: &parser.BitwiseConjunction{
																													Shift: &parser.Shift{
																														Addition: &parser.Addition{
																															Multiplication: &parser.Multiplication{
																																Unary: &parser.Unary{
																																	Accessor: &parser.Accessor{
																																		Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																				Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																			},
																			{
																				Condition: &parser.Expression{
																					ListConstruction: &parser.ListConstruction{
																						Disjunction: &parser.Disjunction{
																							Conjunction: &parser.Conjunction{
																								Equality: &parser.Equality{
																									Comparison: &parser.Comparison{
																										BitwiseDisjunction: &parser.BitwiseDisjunction{
																											BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																												BitwiseConjunction: &parser.BitwiseConjunction{
																													Shift: &parser.Shift{
																														Addition: &parser.Addition{
																															Multiplication: &parser.Multiplication{
																																Unary: &parser.Unary{
																																	Accessor: &parser.Accessor{
																																		Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																				Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
				disjunction: &parser.Disjunction{
					Conjunction: &parser.Conjunction{
						Equality: &parser.Equality{
							Comparison: &parser.Comparison{
								BitwiseDisjunction: &parser.BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
										BitwiseConjunction: &parser.BitwiseConjunction{
											Shift: &parser.Shift{
												Addition: &parser.Addition{
													Multiplication: &parser.Multiplication{
														Unary: &parser.Unary{
															Accessor: &parser.Accessor{
																Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateDisjunction(data.args.disjunction, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateConjunction(test *testing.T) {
	type args struct {
		conjunction         *parser.Conjunction
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
				conjunction: &parser.Conjunction{
					Equality: &parser.Equality{
						Comparison: &parser.Comparison{
							BitwiseDisjunction: &parser.BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
									BitwiseConjunction: &parser.BitwiseConjunction{
										Shift: &parser.Shift{
											Addition: &parser.Addition{
												Multiplication: &parser.Multiplication{
													Unary: &parser.Unary{
														Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					Conjunction: &parser.Conjunction{
						Equality: &parser.Equality{
							Comparison: &parser.Comparison{
								BitwiseDisjunction: &parser.BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
										BitwiseConjunction: &parser.BitwiseConjunction{
											Shift: &parser.Shift{
												Addition: &parser.Addition{
													Multiplication: &parser.Multiplication{
														Unary: &parser.Unary{
															Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						Conjunction: &parser.Conjunction{
							Equality: &parser.Equality{
								Comparison: &parser.Comparison{
									BitwiseDisjunction: &parser.BitwiseDisjunction{
										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
											BitwiseConjunction: &parser.BitwiseConjunction{
												Shift: &parser.Shift{
													Addition: &parser.Addition{
														Multiplication: &parser.Multiplication{
															Unary: &parser.Unary{
																Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
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
				conjunction: &parser.Conjunction{
					Equality: &parser.Equality{
						Comparison: &parser.Comparison{
							BitwiseDisjunction: &parser.BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
									BitwiseConjunction: &parser.BitwiseConjunction{
										Shift: &parser.Shift{
											Addition: &parser.Addition{
												Multiplication: &parser.Multiplication{
													Unary: &parser.Unary{
														Accessor: &parser.Accessor{
															Atom: &parser.Atom{
																ConditionalExpression: &parser.ConditionalExpression{
																	ConditionalCases: []*parser.ConditionalCase{
																		{
																			Condition: &parser.Expression{
																				ListConstruction: &parser.ListConstruction{
																					Disjunction: &parser.Disjunction{
																						Conjunction: &parser.Conjunction{
																							Equality: &parser.Equality{
																								Comparison: &parser.Comparison{
																									BitwiseDisjunction: &parser.BitwiseDisjunction{
																										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																											BitwiseConjunction: &parser.BitwiseConjunction{
																												Shift: &parser.Shift{
																													Addition: &parser.Addition{
																														Multiplication: &parser.Multiplication{
																															Unary: &parser.Unary{
																																Accessor: &parser.Accessor{
																																	Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																			Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																		},
																		{
																			Condition: &parser.Expression{
																				ListConstruction: &parser.ListConstruction{
																					Disjunction: &parser.Disjunction{
																						Conjunction: &parser.Conjunction{
																							Equality: &parser.Equality{
																								Comparison: &parser.Comparison{
																									BitwiseDisjunction: &parser.BitwiseDisjunction{
																										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																											BitwiseConjunction: &parser.BitwiseConjunction{
																												Shift: &parser.Shift{
																													Addition: &parser.Addition{
																														Multiplication: &parser.Multiplication{
																															Unary: &parser.Unary{
																																Accessor: &parser.Accessor{
																																	Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																			Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
					Conjunction: &parser.Conjunction{
						Equality: &parser.Equality{
							Comparison: &parser.Comparison{
								BitwiseDisjunction: &parser.BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
										BitwiseConjunction: &parser.BitwiseConjunction{
											Shift: &parser.Shift{
												Addition: &parser.Addition{
													Multiplication: &parser.Multiplication{
														Unary: &parser.Unary{
															Accessor: &parser.Accessor{
																Atom: &parser.Atom{
																	ConditionalExpression: &parser.ConditionalExpression{
																		ConditionalCases: []*parser.ConditionalCase{
																			{
																				Condition: &parser.Expression{
																					ListConstruction: &parser.ListConstruction{
																						Disjunction: &parser.Disjunction{
																							Conjunction: &parser.Conjunction{
																								Equality: &parser.Equality{
																									Comparison: &parser.Comparison{
																										BitwiseDisjunction: &parser.BitwiseDisjunction{
																											BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																												BitwiseConjunction: &parser.BitwiseConjunction{
																													Shift: &parser.Shift{
																														Addition: &parser.Addition{
																															Multiplication: &parser.Multiplication{
																																Unary: &parser.Unary{
																																	Accessor: &parser.Accessor{
																																		Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
																				Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																			},
																			{
																				Condition: &parser.Expression{
																					ListConstruction: &parser.ListConstruction{
																						Disjunction: &parser.Disjunction{
																							Conjunction: &parser.Conjunction{
																								Equality: &parser.Equality{
																									Comparison: &parser.Comparison{
																										BitwiseDisjunction: &parser.BitwiseDisjunction{
																											BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																												BitwiseConjunction: &parser.BitwiseConjunction{
																													Shift: &parser.Shift{
																														Addition: &parser.Addition{
																															Multiplication: &parser.Multiplication{
																																Unary: &parser.Unary{
																																	Accessor: &parser.Accessor{
																																		Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
																				Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				conjunction: &parser.Conjunction{
					Equality: &parser.Equality{
						Comparison: &parser.Comparison{
							BitwiseDisjunction: &parser.BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
									BitwiseConjunction: &parser.BitwiseConjunction{
										Shift: &parser.Shift{
											Addition: &parser.Addition{
												Multiplication: &parser.Multiplication{
													Unary: &parser.Unary{
														Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					Conjunction: &parser.Conjunction{
						Equality: &parser.Equality{
							Comparison: &parser.Comparison{
								BitwiseDisjunction: &parser.BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
										BitwiseConjunction: &parser.BitwiseConjunction{
											Shift: &parser.Shift{
												Addition: &parser.Addition{
													Multiplication: &parser.Multiplication{
														Unary: &parser.Unary{
															Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
														},
													},
												},
											},
										},
									},
								},
							},
						},
						Conjunction: &parser.Conjunction{
							Equality: &parser.Equality{
								Comparison: &parser.Comparison{
									BitwiseDisjunction: &parser.BitwiseDisjunction{
										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
											BitwiseConjunction: &parser.BitwiseConjunction{
												Shift: &parser.Shift{
													Addition: &parser.Addition{
														Multiplication: &parser.Multiplication{
															Unary: &parser.Unary{
																Accessor: &parser.Accessor{
																	Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Conjunction/empty/success",
			args: args{
				conjunction: &parser.Conjunction{
					Equality: &parser.Equality{
						Comparison: &parser.Comparison{
							BitwiseDisjunction: &parser.BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
									BitwiseConjunction: &parser.BitwiseConjunction{
										Shift: &parser.Shift{
											Addition: &parser.Addition{
												Multiplication: &parser.Multiplication{
													Unary: &parser.Unary{
														Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Conjunction/empty/success/with setted states",
			args: args{
				conjunction: &parser.Conjunction{
					Equality: &parser.Equality{
						Comparison: &parser.Comparison{
							BitwiseDisjunction: &parser.BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
									BitwiseConjunction: &parser.BitwiseConjunction{
										Shift: &parser.Shift{
											Addition: &parser.Addition{
												Multiplication: &parser.Multiplication{
													Unary: &parser.Unary{
														Accessor: &parser.Accessor{
															Atom: &parser.Atom{
																ConditionalExpression: &parser.ConditionalExpression{
																	ConditionalCases: []*parser.ConditionalCase{
																		{
																			Condition: &parser.Expression{
																				ListConstruction: &parser.ListConstruction{
																					Disjunction: &parser.Disjunction{
																						Conjunction: &parser.Conjunction{
																							Equality: &parser.Equality{
																								Comparison: &parser.Comparison{
																									BitwiseDisjunction: &parser.BitwiseDisjunction{
																										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																											BitwiseConjunction: &parser.BitwiseConjunction{
																												Shift: &parser.Shift{
																													Addition: &parser.Addition{
																														Multiplication: &parser.Multiplication{
																															Unary: &parser.Unary{
																																Accessor: &parser.Accessor{
																																	Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																			Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																		},
																		{
																			Condition: &parser.Expression{
																				ListConstruction: &parser.ListConstruction{
																					Disjunction: &parser.Disjunction{
																						Conjunction: &parser.Conjunction{
																							Equality: &parser.Equality{
																								Comparison: &parser.Comparison{
																									BitwiseDisjunction: &parser.BitwiseDisjunction{
																										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																											BitwiseConjunction: &parser.BitwiseConjunction{
																												Shift: &parser.Shift{
																													Addition: &parser.Addition{
																														Multiplication: &parser.Multiplication{
																															Unary: &parser.Unary{
																																Accessor: &parser.Accessor{
																																	Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																			Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
				conjunction: &parser.Conjunction{
					Equality: &parser.Equality{
						Comparison: &parser.Comparison{
							BitwiseDisjunction: &parser.BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
									BitwiseConjunction: &parser.BitwiseConjunction{
										Shift: &parser.Shift{
											Addition: &parser.Addition{
												Multiplication: &parser.Multiplication{
													Unary: &parser.Unary{
														Accessor: &parser.Accessor{
															Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateConjunction(data.args.conjunction, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateEquality(test *testing.T) {
	type args struct {
		equality            *parser.Equality
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
				equality: &parser.Equality{
					Comparison: &parser.Comparison{
						BitwiseDisjunction: &parser.BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
								BitwiseConjunction: &parser.BitwiseConjunction{
									Shift: &parser.Shift{
										Addition: &parser.Addition{
											Multiplication: &parser.Multiplication{
												Unary: &parser.Unary{
													Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
												},
											},
										},
									},
								},
							},
						},
					},
					Operation: "==",
					Equality: &parser.Equality{
						Comparison: &parser.Comparison{
							BitwiseDisjunction: &parser.BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
									BitwiseConjunction: &parser.BitwiseConjunction{
										Shift: &parser.Shift{
											Addition: &parser.Addition{
												Multiplication: &parser.Multiplication{
													Unary: &parser.Unary{
														Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
													},
												},
											},
										},
									},
								},
							},
						},
						Operation: "!=",
						Equality: &parser.Equality{
							Comparison: &parser.Comparison{
								BitwiseDisjunction: &parser.BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
										BitwiseConjunction: &parser.BitwiseConjunction{
											Shift: &parser.Shift{
												Addition: &parser.Addition{
													Multiplication: &parser.Multiplication{
														Unary: &parser.Unary{
															Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
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
				equality: &parser.Equality{
					Comparison: &parser.Comparison{
						BitwiseDisjunction: &parser.BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
								BitwiseConjunction: &parser.BitwiseConjunction{
									Shift: &parser.Shift{
										Addition: &parser.Addition{
											Multiplication: &parser.Multiplication{
												Unary: &parser.Unary{
													Accessor: &parser.Accessor{
														Atom: &parser.Atom{
															ConditionalExpression: &parser.ConditionalExpression{
																ConditionalCases: []*parser.ConditionalCase{
																	{
																		Condition: &parser.Expression{
																			ListConstruction: &parser.ListConstruction{
																				Disjunction: &parser.Disjunction{
																					Conjunction: &parser.Conjunction{
																						Equality: &parser.Equality{
																							Comparison: &parser.Comparison{
																								BitwiseDisjunction: &parser.BitwiseDisjunction{
																									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																										BitwiseConjunction: &parser.BitwiseConjunction{
																											Shift: &parser.Shift{
																												Addition: &parser.Addition{
																													Multiplication: &parser.Multiplication{
																														Unary: &parser.Unary{
																															Accessor: &parser.Accessor{
																																Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																		Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																	},
																	{
																		Condition: &parser.Expression{
																			ListConstruction: &parser.ListConstruction{
																				Disjunction: &parser.Disjunction{
																					Conjunction: &parser.Conjunction{
																						Equality: &parser.Equality{
																							Comparison: &parser.Comparison{
																								BitwiseDisjunction: &parser.BitwiseDisjunction{
																									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																										BitwiseConjunction: &parser.BitwiseConjunction{
																											Shift: &parser.Shift{
																												Addition: &parser.Addition{
																													Multiplication: &parser.Multiplication{
																														Unary: &parser.Unary{
																															Accessor: &parser.Accessor{
																																Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																		Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
					Operation: "==",
					Equality: &parser.Equality{
						Comparison: &parser.Comparison{
							BitwiseDisjunction: &parser.BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
									BitwiseConjunction: &parser.BitwiseConjunction{
										Shift: &parser.Shift{
											Addition: &parser.Addition{
												Multiplication: &parser.Multiplication{
													Unary: &parser.Unary{
														Accessor: &parser.Accessor{
															Atom: &parser.Atom{
																ConditionalExpression: &parser.ConditionalExpression{
																	ConditionalCases: []*parser.ConditionalCase{
																		{
																			Condition: &parser.Expression{
																				ListConstruction: &parser.ListConstruction{
																					Disjunction: &parser.Disjunction{
																						Conjunction: &parser.Conjunction{
																							Equality: &parser.Equality{
																								Comparison: &parser.Comparison{
																									BitwiseDisjunction: &parser.BitwiseDisjunction{
																										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																											BitwiseConjunction: &parser.BitwiseConjunction{
																												Shift: &parser.Shift{
																													Addition: &parser.Addition{
																														Multiplication: &parser.Multiplication{
																															Unary: &parser.Unary{
																																Accessor: &parser.Accessor{
																																	Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
																			Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																		},
																		{
																			Condition: &parser.Expression{
																				ListConstruction: &parser.ListConstruction{
																					Disjunction: &parser.Disjunction{
																						Conjunction: &parser.Conjunction{
																							Equality: &parser.Equality{
																								Comparison: &parser.Comparison{
																									BitwiseDisjunction: &parser.BitwiseDisjunction{
																										BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																											BitwiseConjunction: &parser.BitwiseConjunction{
																												Shift: &parser.Shift{
																													Addition: &parser.Addition{
																														Multiplication: &parser.Multiplication{
																															Unary: &parser.Unary{
																																Accessor: &parser.Accessor{
																																	Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
																			Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				equality: &parser.Equality{
					Comparison: &parser.Comparison{
						BitwiseDisjunction: &parser.BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
								BitwiseConjunction: &parser.BitwiseConjunction{
									Shift: &parser.Shift{
										Addition: &parser.Addition{
											Multiplication: &parser.Multiplication{
												Unary: &parser.Unary{
													Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
												},
											},
										},
									},
								},
							},
						},
					},
					Operation: "==",
					Equality: &parser.Equality{
						Comparison: &parser.Comparison{
							BitwiseDisjunction: &parser.BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
									BitwiseConjunction: &parser.BitwiseConjunction{
										Shift: &parser.Shift{
											Addition: &parser.Addition{
												Multiplication: &parser.Multiplication{
													Unary: &parser.Unary{
														Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
													},
												},
											},
										},
									},
								},
							},
						},
						Operation: "!=",
						Equality: &parser.Equality{
							Comparison: &parser.Comparison{
								BitwiseDisjunction: &parser.BitwiseDisjunction{
									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
										BitwiseConjunction: &parser.BitwiseConjunction{
											Shift: &parser.Shift{
												Addition: &parser.Addition{
													Multiplication: &parser.Multiplication{
														Unary: &parser.Unary{
															Accessor: &parser.Accessor{
																Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Equality/empty/success",
			args: args{
				equality: &parser.Equality{
					Comparison: &parser.Comparison{
						BitwiseDisjunction: &parser.BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
								BitwiseConjunction: &parser.BitwiseConjunction{
									Shift: &parser.Shift{
										Addition: &parser.Addition{
											Multiplication: &parser.Multiplication{
												Unary: &parser.Unary{
													Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Equality/empty/success/with setted states",
			args: args{
				equality: &parser.Equality{
					Comparison: &parser.Comparison{
						BitwiseDisjunction: &parser.BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
								BitwiseConjunction: &parser.BitwiseConjunction{
									Shift: &parser.Shift{
										Addition: &parser.Addition{
											Multiplication: &parser.Multiplication{
												Unary: &parser.Unary{
													Accessor: &parser.Accessor{
														Atom: &parser.Atom{
															ConditionalExpression: &parser.ConditionalExpression{
																ConditionalCases: []*parser.ConditionalCase{
																	{
																		Condition: &parser.Expression{
																			ListConstruction: &parser.ListConstruction{
																				Disjunction: &parser.Disjunction{
																					Conjunction: &parser.Conjunction{
																						Equality: &parser.Equality{
																							Comparison: &parser.Comparison{
																								BitwiseDisjunction: &parser.BitwiseDisjunction{
																									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																										BitwiseConjunction: &parser.BitwiseConjunction{
																											Shift: &parser.Shift{
																												Addition: &parser.Addition{
																													Multiplication: &parser.Multiplication{
																														Unary: &parser.Unary{
																															Accessor: &parser.Accessor{
																																Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																		Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																	},
																	{
																		Condition: &parser.Expression{
																			ListConstruction: &parser.ListConstruction{
																				Disjunction: &parser.Disjunction{
																					Conjunction: &parser.Conjunction{
																						Equality: &parser.Equality{
																							Comparison: &parser.Comparison{
																								BitwiseDisjunction: &parser.BitwiseDisjunction{
																									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																										BitwiseConjunction: &parser.BitwiseConjunction{
																											Shift: &parser.Shift{
																												Addition: &parser.Addition{
																													Multiplication: &parser.Multiplication{
																														Unary: &parser.Unary{
																															Accessor: &parser.Accessor{
																																Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																		Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
				equality: &parser.Equality{
					Comparison: &parser.Comparison{
						BitwiseDisjunction: &parser.BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
								BitwiseConjunction: &parser.BitwiseConjunction{
									Shift: &parser.Shift{
										Addition: &parser.Addition{
											Multiplication: &parser.Multiplication{
												Unary: &parser.Unary{
													Accessor: &parser.Accessor{
														Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateEquality(data.args.equality, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateComparison(test *testing.T) {
	type args struct {
		comparison          *parser.Comparison
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
				comparison: &parser.Comparison{
					BitwiseDisjunction: &parser.BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
							BitwiseConjunction: &parser.BitwiseConjunction{
								Shift: &parser.Shift{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
											},
										},
									},
								},
							},
						},
					},
					Operation: "<",
					Comparison: &parser.Comparison{
						BitwiseDisjunction: &parser.BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
								BitwiseConjunction: &parser.BitwiseConjunction{
									Shift: &parser.Shift{
										Addition: &parser.Addition{
											Multiplication: &parser.Multiplication{
												Unary: &parser.Unary{
													Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
												},
											},
										},
									},
								},
							},
						},
						Operation: "<",
						Comparison: &parser.Comparison{
							BitwiseDisjunction: &parser.BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
									BitwiseConjunction: &parser.BitwiseConjunction{
										Shift: &parser.Shift{
											Addition: &parser.Addition{
												Multiplication: &parser.Multiplication{
													Unary: &parser.Unary{
														Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
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
				comparison: &parser.Comparison{
					BitwiseDisjunction: &parser.BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
							BitwiseConjunction: &parser.BitwiseConjunction{
								Shift: &parser.Shift{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
											},
										},
									},
								},
							},
						},
					},
					Operation: "<=",
					Comparison: &parser.Comparison{
						BitwiseDisjunction: &parser.BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
								BitwiseConjunction: &parser.BitwiseConjunction{
									Shift: &parser.Shift{
										Addition: &parser.Addition{
											Multiplication: &parser.Multiplication{
												Unary: &parser.Unary{
													Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
												},
											},
										},
									},
								},
							},
						},
						Operation: "<=",
						Comparison: &parser.Comparison{
							BitwiseDisjunction: &parser.BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
									BitwiseConjunction: &parser.BitwiseConjunction{
										Shift: &parser.Shift{
											Addition: &parser.Addition{
												Multiplication: &parser.Multiplication{
													Unary: &parser.Unary{
														Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
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
				comparison: &parser.Comparison{
					BitwiseDisjunction: &parser.BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
							BitwiseConjunction: &parser.BitwiseConjunction{
								Shift: &parser.Shift{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
											},
										},
									},
								},
							},
						},
					},
					Operation: ">",
					Comparison: &parser.Comparison{
						BitwiseDisjunction: &parser.BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
								BitwiseConjunction: &parser.BitwiseConjunction{
									Shift: &parser.Shift{
										Addition: &parser.Addition{
											Multiplication: &parser.Multiplication{
												Unary: &parser.Unary{
													Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
												},
											},
										},
									},
								},
							},
						},
						Operation: ">",
						Comparison: &parser.Comparison{
							BitwiseDisjunction: &parser.BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
									BitwiseConjunction: &parser.BitwiseConjunction{
										Shift: &parser.Shift{
											Addition: &parser.Addition{
												Multiplication: &parser.Multiplication{
													Unary: &parser.Unary{
														Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
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
				comparison: &parser.Comparison{
					BitwiseDisjunction: &parser.BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
							BitwiseConjunction: &parser.BitwiseConjunction{
								Shift: &parser.Shift{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
											},
										},
									},
								},
							},
						},
					},
					Operation: ">=",
					Comparison: &parser.Comparison{
						BitwiseDisjunction: &parser.BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
								BitwiseConjunction: &parser.BitwiseConjunction{
									Shift: &parser.Shift{
										Addition: &parser.Addition{
											Multiplication: &parser.Multiplication{
												Unary: &parser.Unary{
													Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
												},
											},
										},
									},
								},
							},
						},
						Operation: ">=",
						Comparison: &parser.Comparison{
							BitwiseDisjunction: &parser.BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
									BitwiseConjunction: &parser.BitwiseConjunction{
										Shift: &parser.Shift{
											Addition: &parser.Addition{
												Multiplication: &parser.Multiplication{
													Unary: &parser.Unary{
														Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
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
				comparison: &parser.Comparison{
					BitwiseDisjunction: &parser.BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
							BitwiseConjunction: &parser.BitwiseConjunction{
								Shift: &parser.Shift{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{
													Atom: &parser.Atom{
														ConditionalExpression: &parser.ConditionalExpression{
															ConditionalCases: []*parser.ConditionalCase{
																{
																	Condition: &parser.Expression{
																		ListConstruction: &parser.ListConstruction{
																			Disjunction: &parser.Disjunction{
																				Conjunction: &parser.Conjunction{
																					Equality: &parser.Equality{
																						Comparison: &parser.Comparison{
																							BitwiseDisjunction: &parser.BitwiseDisjunction{
																								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																									BitwiseConjunction: &parser.BitwiseConjunction{
																										Shift: &parser.Shift{
																											Addition: &parser.Addition{
																												Multiplication: &parser.Multiplication{
																													Unary: &parser.Unary{
																														Accessor: &parser.Accessor{
																															Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																	Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																},
																{
																	Condition: &parser.Expression{
																		ListConstruction: &parser.ListConstruction{
																			Disjunction: &parser.Disjunction{
																				Conjunction: &parser.Conjunction{
																					Equality: &parser.Equality{
																						Comparison: &parser.Comparison{
																							BitwiseDisjunction: &parser.BitwiseDisjunction{
																								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																									BitwiseConjunction: &parser.BitwiseConjunction{
																										Shift: &parser.Shift{
																											Addition: &parser.Addition{
																												Multiplication: &parser.Multiplication{
																													Unary: &parser.Unary{
																														Accessor: &parser.Accessor{
																															Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																	Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
					Operation: "<",
					Comparison: &parser.Comparison{
						BitwiseDisjunction: &parser.BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
								BitwiseConjunction: &parser.BitwiseConjunction{
									Shift: &parser.Shift{
										Addition: &parser.Addition{
											Multiplication: &parser.Multiplication{
												Unary: &parser.Unary{
													Accessor: &parser.Accessor{
														Atom: &parser.Atom{
															ConditionalExpression: &parser.ConditionalExpression{
																ConditionalCases: []*parser.ConditionalCase{
																	{
																		Condition: &parser.Expression{
																			ListConstruction: &parser.ListConstruction{
																				Disjunction: &parser.Disjunction{
																					Conjunction: &parser.Conjunction{
																						Equality: &parser.Equality{
																							Comparison: &parser.Comparison{
																								BitwiseDisjunction: &parser.BitwiseDisjunction{
																									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																										BitwiseConjunction: &parser.BitwiseConjunction{
																											Shift: &parser.Shift{
																												Addition: &parser.Addition{
																													Multiplication: &parser.Multiplication{
																														Unary: &parser.Unary{
																															Accessor: &parser.Accessor{
																																Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
																		Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																	},
																	{
																		Condition: &parser.Expression{
																			ListConstruction: &parser.ListConstruction{
																				Disjunction: &parser.Disjunction{
																					Conjunction: &parser.Conjunction{
																						Equality: &parser.Equality{
																							Comparison: &parser.Comparison{
																								BitwiseDisjunction: &parser.BitwiseDisjunction{
																									BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																										BitwiseConjunction: &parser.BitwiseConjunction{
																											Shift: &parser.Shift{
																												Addition: &parser.Addition{
																													Multiplication: &parser.Multiplication{
																														Unary: &parser.Unary{
																															Accessor: &parser.Accessor{
																																Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
																		Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				comparison: &parser.Comparison{
					BitwiseDisjunction: &parser.BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
							BitwiseConjunction: &parser.BitwiseConjunction{
								Shift: &parser.Shift{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
											},
										},
									},
								},
							},
						},
					},
					Operation: "<",
					Comparison: &parser.Comparison{
						BitwiseDisjunction: &parser.BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
								BitwiseConjunction: &parser.BitwiseConjunction{
									Shift: &parser.Shift{
										Addition: &parser.Addition{
											Multiplication: &parser.Multiplication{
												Unary: &parser.Unary{
													Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
												},
											},
										},
									},
								},
							},
						},
						Operation: "<",
						Comparison: &parser.Comparison{
							BitwiseDisjunction: &parser.BitwiseDisjunction{
								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
									BitwiseConjunction: &parser.BitwiseConjunction{
										Shift: &parser.Shift{
											Addition: &parser.Addition{
												Multiplication: &parser.Multiplication{
													Unary: &parser.Unary{
														Accessor: &parser.Accessor{
															Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Comparison/empty/success",
			args: args{
				comparison: &parser.Comparison{
					BitwiseDisjunction: &parser.BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
							BitwiseConjunction: &parser.BitwiseConjunction{
								Shift: &parser.Shift{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Comparison/empty/success/with setted states",
			args: args{
				comparison: &parser.Comparison{
					BitwiseDisjunction: &parser.BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
							BitwiseConjunction: &parser.BitwiseConjunction{
								Shift: &parser.Shift{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{
													Atom: &parser.Atom{
														ConditionalExpression: &parser.ConditionalExpression{
															ConditionalCases: []*parser.ConditionalCase{
																{
																	Condition: &parser.Expression{
																		ListConstruction: &parser.ListConstruction{
																			Disjunction: &parser.Disjunction{
																				Conjunction: &parser.Conjunction{
																					Equality: &parser.Equality{
																						Comparison: &parser.Comparison{
																							BitwiseDisjunction: &parser.BitwiseDisjunction{
																								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																									BitwiseConjunction: &parser.BitwiseConjunction{
																										Shift: &parser.Shift{
																											Addition: &parser.Addition{
																												Multiplication: &parser.Multiplication{
																													Unary: &parser.Unary{
																														Accessor: &parser.Accessor{
																															Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																	Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																},
																{
																	Condition: &parser.Expression{
																		ListConstruction: &parser.ListConstruction{
																			Disjunction: &parser.Disjunction{
																				Conjunction: &parser.Conjunction{
																					Equality: &parser.Equality{
																						Comparison: &parser.Comparison{
																							BitwiseDisjunction: &parser.BitwiseDisjunction{
																								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																									BitwiseConjunction: &parser.BitwiseConjunction{
																										Shift: &parser.Shift{
																											Addition: &parser.Addition{
																												Multiplication: &parser.Multiplication{
																													Unary: &parser.Unary{
																														Accessor: &parser.Accessor{
																															Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																	Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
				comparison: &parser.Comparison{
					BitwiseDisjunction: &parser.BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
							BitwiseConjunction: &parser.BitwiseConjunction{
								Shift: &parser.Shift{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: pointer.ToString("unknown")}},
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateComparison(data.args.comparison, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateBitwiseDisjunction(test *testing.T) {
	type args struct {
		bitwiseDisjunction  *parser.BitwiseDisjunction
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
				bitwiseDisjunction: &parser.BitwiseDisjunction{
					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
						BitwiseConjunction: &parser.BitwiseConjunction{
							Shift: &parser.Shift{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
										},
									},
								},
							},
						},
					},
					BitwiseDisjunction: &parser.BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
							BitwiseConjunction: &parser.BitwiseConjunction{
								Shift: &parser.Shift{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
											},
										},
									},
								},
							},
						},
						BitwiseDisjunction: &parser.BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
								BitwiseConjunction: &parser.BitwiseConjunction{
									Shift: &parser.Shift{
										Addition: &parser.Addition{
											Multiplication: &parser.Multiplication{
												Unary: &parser.Unary{
													Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
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
				bitwiseDisjunction: &parser.BitwiseDisjunction{
					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
						BitwiseConjunction: &parser.BitwiseConjunction{
							Shift: &parser.Shift{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{
												Atom: &parser.Atom{
													ConditionalExpression: &parser.ConditionalExpression{
														ConditionalCases: []*parser.ConditionalCase{
															{
																Condition: &parser.Expression{
																	ListConstruction: &parser.ListConstruction{
																		Disjunction: &parser.Disjunction{
																			Conjunction: &parser.Conjunction{
																				Equality: &parser.Equality{
																					Comparison: &parser.Comparison{
																						BitwiseDisjunction: &parser.BitwiseDisjunction{
																							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																								BitwiseConjunction: &parser.BitwiseConjunction{
																									Shift: &parser.Shift{
																										Addition: &parser.Addition{
																											Multiplication: &parser.Multiplication{
																												Unary: &parser.Unary{
																													Accessor: &parser.Accessor{
																														Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
															},
															{
																Condition: &parser.Expression{
																	ListConstruction: &parser.ListConstruction{
																		Disjunction: &parser.Disjunction{
																			Conjunction: &parser.Conjunction{
																				Equality: &parser.Equality{
																					Comparison: &parser.Comparison{
																						BitwiseDisjunction: &parser.BitwiseDisjunction{
																							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																								BitwiseConjunction: &parser.BitwiseConjunction{
																									Shift: &parser.Shift{
																										Addition: &parser.Addition{
																											Multiplication: &parser.Multiplication{
																												Unary: &parser.Unary{
																													Accessor: &parser.Accessor{
																														Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
					BitwiseDisjunction: &parser.BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
							BitwiseConjunction: &parser.BitwiseConjunction{
								Shift: &parser.Shift{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{
													Atom: &parser.Atom{
														ConditionalExpression: &parser.ConditionalExpression{
															ConditionalCases: []*parser.ConditionalCase{
																{
																	Condition: &parser.Expression{
																		ListConstruction: &parser.ListConstruction{
																			Disjunction: &parser.Disjunction{
																				Conjunction: &parser.Conjunction{
																					Equality: &parser.Equality{
																						Comparison: &parser.Comparison{
																							BitwiseDisjunction: &parser.BitwiseDisjunction{
																								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																									BitwiseConjunction: &parser.BitwiseConjunction{
																										Shift: &parser.Shift{
																											Addition: &parser.Addition{
																												Multiplication: &parser.Multiplication{
																													Unary: &parser.Unary{
																														Accessor: &parser.Accessor{
																															Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
																	Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																},
																{
																	Condition: &parser.Expression{
																		ListConstruction: &parser.ListConstruction{
																			Disjunction: &parser.Disjunction{
																				Conjunction: &parser.Conjunction{
																					Equality: &parser.Equality{
																						Comparison: &parser.Comparison{
																							BitwiseDisjunction: &parser.BitwiseDisjunction{
																								BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																									BitwiseConjunction: &parser.BitwiseConjunction{
																										Shift: &parser.Shift{
																											Addition: &parser.Addition{
																												Multiplication: &parser.Multiplication{
																													Unary: &parser.Unary{
																														Accessor: &parser.Accessor{
																															Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
																	Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				bitwiseDisjunction: &parser.BitwiseDisjunction{
					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
						BitwiseConjunction: &parser.BitwiseConjunction{
							Shift: &parser.Shift{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
										},
									},
								},
							},
						},
					},
					BitwiseDisjunction: &parser.BitwiseDisjunction{
						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
							BitwiseConjunction: &parser.BitwiseConjunction{
								Shift: &parser.Shift{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
											},
										},
									},
								},
							},
						},
						BitwiseDisjunction: &parser.BitwiseDisjunction{
							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
								BitwiseConjunction: &parser.BitwiseConjunction{
									Shift: &parser.Shift{
										Addition: &parser.Addition{
											Multiplication: &parser.Multiplication{
												Unary: &parser.Unary{
													Accessor: &parser.Accessor{
														Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "BitwiseDisjunction/empty/success",
			args: args{
				bitwiseDisjunction: &parser.BitwiseDisjunction{
					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
						BitwiseConjunction: &parser.BitwiseConjunction{
							Shift: &parser.Shift{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseDisjunction/empty/success/with setted states",
			args: args{
				bitwiseDisjunction: &parser.BitwiseDisjunction{
					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
						BitwiseConjunction: &parser.BitwiseConjunction{
							Shift: &parser.Shift{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{
												Atom: &parser.Atom{
													ConditionalExpression: &parser.ConditionalExpression{
														ConditionalCases: []*parser.ConditionalCase{
															{
																Condition: &parser.Expression{
																	ListConstruction: &parser.ListConstruction{
																		Disjunction: &parser.Disjunction{
																			Conjunction: &parser.Conjunction{
																				Equality: &parser.Equality{
																					Comparison: &parser.Comparison{
																						BitwiseDisjunction: &parser.BitwiseDisjunction{
																							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																								BitwiseConjunction: &parser.BitwiseConjunction{
																									Shift: &parser.Shift{
																										Addition: &parser.Addition{
																											Multiplication: &parser.Multiplication{
																												Unary: &parser.Unary{
																													Accessor: &parser.Accessor{
																														Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
															},
															{
																Condition: &parser.Expression{
																	ListConstruction: &parser.ListConstruction{
																		Disjunction: &parser.Disjunction{
																			Conjunction: &parser.Conjunction{
																				Equality: &parser.Equality{
																					Comparison: &parser.Comparison{
																						BitwiseDisjunction: &parser.BitwiseDisjunction{
																							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																								BitwiseConjunction: &parser.BitwiseConjunction{
																									Shift: &parser.Shift{
																										Addition: &parser.Addition{
																											Multiplication: &parser.Multiplication{
																												Unary: &parser.Unary{
																													Accessor: &parser.Accessor{
																														Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
				bitwiseDisjunction: &parser.BitwiseDisjunction{
					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
						BitwiseConjunction: &parser.BitwiseConjunction{
							Shift: &parser.Shift{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: pointer.ToString("unknown")}},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr :=
				translateBitwiseDisjunction(data.args.bitwiseDisjunction, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			assert.Equal(test, data.wantSettedStates, gotSettedStates)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateBitwiseExclusiveDisjunction(test *testing.T) {
	type args struct {
		bitwiseExclusiveDisjunction *parser.BitwiseExclusiveDisjunction
		declaredIdentifiers         mapset.Set
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
				bitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
					BitwiseConjunction: &parser.BitwiseConjunction{
						Shift: &parser.Shift{
							Addition: &parser.Addition{
								Multiplication: &parser.Multiplication{
									Unary: &parser.Unary{
										Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
									},
								},
							},
						},
					},
					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
						BitwiseConjunction: &parser.BitwiseConjunction{
							Shift: &parser.Shift{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
										},
									},
								},
							},
						},
						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
							BitwiseConjunction: &parser.BitwiseConjunction{
								Shift: &parser.Shift{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
											},
										},
									},
								},
							},
						},
					},
				},
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
				bitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
					BitwiseConjunction: &parser.BitwiseConjunction{
						Shift: &parser.Shift{
							Addition: &parser.Addition{
								Multiplication: &parser.Multiplication{
									Unary: &parser.Unary{
										Accessor: &parser.Accessor{
											Atom: &parser.Atom{
												ConditionalExpression: &parser.ConditionalExpression{
													ConditionalCases: []*parser.ConditionalCase{
														{
															Condition: &parser.Expression{
																ListConstruction: &parser.ListConstruction{
																	Disjunction: &parser.Disjunction{
																		Conjunction: &parser.Conjunction{
																			Equality: &parser.Equality{
																				Comparison: &parser.Comparison{
																					BitwiseDisjunction: &parser.BitwiseDisjunction{
																						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																							BitwiseConjunction: &parser.BitwiseConjunction{
																								Shift: &parser.Shift{
																									Addition: &parser.Addition{
																										Multiplication: &parser.Multiplication{
																											Unary: &parser.Unary{
																												Accessor: &parser.Accessor{
																													Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
															Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
														},
														{
															Condition: &parser.Expression{
																ListConstruction: &parser.ListConstruction{
																	Disjunction: &parser.Disjunction{
																		Conjunction: &parser.Conjunction{
																			Equality: &parser.Equality{
																				Comparison: &parser.Comparison{
																					BitwiseDisjunction: &parser.BitwiseDisjunction{
																						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																							BitwiseConjunction: &parser.BitwiseConjunction{
																								Shift: &parser.Shift{
																									Addition: &parser.Addition{
																										Multiplication: &parser.Multiplication{
																											Unary: &parser.Unary{
																												Accessor: &parser.Accessor{
																													Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
															Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
						BitwiseConjunction: &parser.BitwiseConjunction{
							Shift: &parser.Shift{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{
												Atom: &parser.Atom{
													ConditionalExpression: &parser.ConditionalExpression{
														ConditionalCases: []*parser.ConditionalCase{
															{
																Condition: &parser.Expression{
																	ListConstruction: &parser.ListConstruction{
																		Disjunction: &parser.Disjunction{
																			Conjunction: &parser.Conjunction{
																				Equality: &parser.Equality{
																					Comparison: &parser.Comparison{
																						BitwiseDisjunction: &parser.BitwiseDisjunction{
																							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																								BitwiseConjunction: &parser.BitwiseConjunction{
																									Shift: &parser.Shift{
																										Addition: &parser.Addition{
																											Multiplication: &parser.Multiplication{
																												Unary: &parser.Unary{
																													Accessor: &parser.Accessor{
																														Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
																Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
															},
															{
																Condition: &parser.Expression{
																	ListConstruction: &parser.ListConstruction{
																		Disjunction: &parser.Disjunction{
																			Conjunction: &parser.Conjunction{
																				Equality: &parser.Equality{
																					Comparison: &parser.Comparison{
																						BitwiseDisjunction: &parser.BitwiseDisjunction{
																							BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																								BitwiseConjunction: &parser.BitwiseConjunction{
																									Shift: &parser.Shift{
																										Addition: &parser.Addition{
																											Multiplication: &parser.Multiplication{
																												Unary: &parser.Unary{
																													Accessor: &parser.Accessor{
																														Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
																Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				bitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
					BitwiseConjunction: &parser.BitwiseConjunction{
						Shift: &parser.Shift{
							Addition: &parser.Addition{
								Multiplication: &parser.Multiplication{
									Unary: &parser.Unary{
										Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
									},
								},
							},
						},
					},
					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
						BitwiseConjunction: &parser.BitwiseConjunction{
							Shift: &parser.Shift{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
										},
									},
								},
							},
						},
						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
							BitwiseConjunction: &parser.BitwiseConjunction{
								Shift: &parser.Shift{
									Addition: &parser.Addition{
										Multiplication: &parser.Multiplication{
											Unary: &parser.Unary{
												Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: pointer.ToString("unknown")}},
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "BitwiseExclusiveDisjunction/empty/success",
			args: args{
				bitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
					BitwiseConjunction: &parser.BitwiseConjunction{
						Shift: &parser.Shift{
							Addition: &parser.Addition{
								Multiplication: &parser.Multiplication{
									Unary: &parser.Unary{
										Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseExclusiveDisjunction/empty/success/with setted states",
			args: args{
				bitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
					BitwiseConjunction: &parser.BitwiseConjunction{
						Shift: &parser.Shift{
							Addition: &parser.Addition{
								Multiplication: &parser.Multiplication{
									Unary: &parser.Unary{
										Accessor: &parser.Accessor{
											Atom: &parser.Atom{
												ConditionalExpression: &parser.ConditionalExpression{
													ConditionalCases: []*parser.ConditionalCase{
														{
															Condition: &parser.Expression{
																ListConstruction: &parser.ListConstruction{
																	Disjunction: &parser.Disjunction{
																		Conjunction: &parser.Conjunction{
																			Equality: &parser.Equality{
																				Comparison: &parser.Comparison{
																					BitwiseDisjunction: &parser.BitwiseDisjunction{
																						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																							BitwiseConjunction: &parser.BitwiseConjunction{
																								Shift: &parser.Shift{
																									Addition: &parser.Addition{
																										Multiplication: &parser.Multiplication{
																											Unary: &parser.Unary{
																												Accessor: &parser.Accessor{
																													Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
															Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
														},
														{
															Condition: &parser.Expression{
																ListConstruction: &parser.ListConstruction{
																	Disjunction: &parser.Disjunction{
																		Conjunction: &parser.Conjunction{
																			Equality: &parser.Equality{
																				Comparison: &parser.Comparison{
																					BitwiseDisjunction: &parser.BitwiseDisjunction{
																						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																							BitwiseConjunction: &parser.BitwiseConjunction{
																								Shift: &parser.Shift{
																									Addition: &parser.Addition{
																										Multiplication: &parser.Multiplication{
																											Unary: &parser.Unary{
																												Accessor: &parser.Accessor{
																													Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
															Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
				bitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
					BitwiseConjunction: &parser.BitwiseConjunction{
						Shift: &parser.Shift{
							Addition: &parser.Addition{
								Multiplication: &parser.Multiplication{
									Unary: &parser.Unary{
										Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: pointer.ToString("unknown")}},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotSettedStates, gotErr := translateBitwiseExclusiveDisjunction(
				data.args.bitwiseExclusiveDisjunction,
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
				bitwiseConjunction: &parser.BitwiseConjunction{
					Shift: &parser.Shift{
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
								},
							},
						},
					},
					BitwiseConjunction: &parser.BitwiseConjunction{
						Shift: &parser.Shift{
							Addition: &parser.Addition{
								Multiplication: &parser.Multiplication{
									Unary: &parser.Unary{
										Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
									},
								},
							},
						},
						BitwiseConjunction: &parser.BitwiseConjunction{
							Shift: &parser.Shift{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
										},
									},
								},
							},
						},
					},
				},
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
				bitwiseConjunction: &parser.BitwiseConjunction{
					Shift: &parser.Shift{
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{
										Atom: &parser.Atom{
											ConditionalExpression: &parser.ConditionalExpression{
												ConditionalCases: []*parser.ConditionalCase{
													{
														Condition: &parser.Expression{
															ListConstruction: &parser.ListConstruction{
																Disjunction: &parser.Disjunction{
																	Conjunction: &parser.Conjunction{
																		Equality: &parser.Equality{
																			Comparison: &parser.Comparison{
																				BitwiseDisjunction: &parser.BitwiseDisjunction{
																					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																						BitwiseConjunction: &parser.BitwiseConjunction{
																							Shift: &parser.Shift{
																								Addition: &parser.Addition{
																									Multiplication: &parser.Multiplication{
																										Unary: &parser.Unary{
																											Accessor: &parser.Accessor{
																												Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
														Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
													},
													{
														Condition: &parser.Expression{
															ListConstruction: &parser.ListConstruction{
																Disjunction: &parser.Disjunction{
																	Conjunction: &parser.Conjunction{
																		Equality: &parser.Equality{
																			Comparison: &parser.Comparison{
																				BitwiseDisjunction: &parser.BitwiseDisjunction{
																					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																						BitwiseConjunction: &parser.BitwiseConjunction{
																							Shift: &parser.Shift{
																								Addition: &parser.Addition{
																									Multiplication: &parser.Multiplication{
																										Unary: &parser.Unary{
																											Accessor: &parser.Accessor{
																												Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
														Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					BitwiseConjunction: &parser.BitwiseConjunction{
						Shift: &parser.Shift{
							Addition: &parser.Addition{
								Multiplication: &parser.Multiplication{
									Unary: &parser.Unary{
										Accessor: &parser.Accessor{
											Atom: &parser.Atom{
												ConditionalExpression: &parser.ConditionalExpression{
													ConditionalCases: []*parser.ConditionalCase{
														{
															Condition: &parser.Expression{
																ListConstruction: &parser.ListConstruction{
																	Disjunction: &parser.Disjunction{
																		Conjunction: &parser.Conjunction{
																			Equality: &parser.Equality{
																				Comparison: &parser.Comparison{
																					BitwiseDisjunction: &parser.BitwiseDisjunction{
																						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																							BitwiseConjunction: &parser.BitwiseConjunction{
																								Shift: &parser.Shift{
																									Addition: &parser.Addition{
																										Multiplication: &parser.Multiplication{
																											Unary: &parser.Unary{
																												Accessor: &parser.Accessor{
																													Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
															Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
														},
														{
															Condition: &parser.Expression{
																ListConstruction: &parser.ListConstruction{
																	Disjunction: &parser.Disjunction{
																		Conjunction: &parser.Conjunction{
																			Equality: &parser.Equality{
																				Comparison: &parser.Comparison{
																					BitwiseDisjunction: &parser.BitwiseDisjunction{
																						BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																							BitwiseConjunction: &parser.BitwiseConjunction{
																								Shift: &parser.Shift{
																									Addition: &parser.Addition{
																										Multiplication: &parser.Multiplication{
																											Unary: &parser.Unary{
																												Accessor: &parser.Accessor{
																													Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
															Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				bitwiseConjunction: &parser.BitwiseConjunction{
					Shift: &parser.Shift{
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
								},
							},
						},
					},
					BitwiseConjunction: &parser.BitwiseConjunction{
						Shift: &parser.Shift{
							Addition: &parser.Addition{
								Multiplication: &parser.Multiplication{
									Unary: &parser.Unary{
										Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
									},
								},
							},
						},
						BitwiseConjunction: &parser.BitwiseConjunction{
							Shift: &parser.Shift{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: pointer.ToString("unknown")}},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "BitwiseConjunction/empty/success",
			args: args{
				bitwiseConjunction: &parser.BitwiseConjunction{
					Shift: &parser.Shift{
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "BitwiseConjunction/empty/success/with setted states",
			args: args{
				bitwiseConjunction: &parser.BitwiseConjunction{
					Shift: &parser.Shift{
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{
										Atom: &parser.Atom{
											ConditionalExpression: &parser.ConditionalExpression{
												ConditionalCases: []*parser.ConditionalCase{
													{
														Condition: &parser.Expression{
															ListConstruction: &parser.ListConstruction{
																Disjunction: &parser.Disjunction{
																	Conjunction: &parser.Conjunction{
																		Equality: &parser.Equality{
																			Comparison: &parser.Comparison{
																				BitwiseDisjunction: &parser.BitwiseDisjunction{
																					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																						BitwiseConjunction: &parser.BitwiseConjunction{
																							Shift: &parser.Shift{
																								Addition: &parser.Addition{
																									Multiplication: &parser.Multiplication{
																										Unary: &parser.Unary{
																											Accessor: &parser.Accessor{
																												Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
														Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
													},
													{
														Condition: &parser.Expression{
															ListConstruction: &parser.ListConstruction{
																Disjunction: &parser.Disjunction{
																	Conjunction: &parser.Conjunction{
																		Equality: &parser.Equality{
																			Comparison: &parser.Comparison{
																				BitwiseDisjunction: &parser.BitwiseDisjunction{
																					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																						BitwiseConjunction: &parser.BitwiseConjunction{
																							Shift: &parser.Shift{
																								Addition: &parser.Addition{
																									Multiplication: &parser.Multiplication{
																										Unary: &parser.Unary{
																											Accessor: &parser.Accessor{
																												Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
														Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
				bitwiseConjunction: &parser.BitwiseConjunction{
					Shift: &parser.Shift{
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: pointer.ToString("unknown")}},
								},
							},
						},
					},
				},
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
				shift: &parser.Shift{
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(5)}},
							},
						},
					},
					Operation: "<<",
					Shift: &parser.Shift{
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
								},
							},
						},
						Operation: ">>",
						Shift: &parser.Shift{
							Addition: &parser.Addition{
								Multiplication: &parser.Multiplication{
									Unary: &parser.Unary{
										Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
									},
								},
							},
							Operation: ">>>",
							Shift: &parser.Shift{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
										},
									},
								},
							},
						},
					},
				},
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
				shift: &parser.Shift{
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{
									Atom: &parser.Atom{
										ConditionalExpression: &parser.ConditionalExpression{
											ConditionalCases: []*parser.ConditionalCase{
												{
													Condition: &parser.Expression{
														ListConstruction: &parser.ListConstruction{
															Disjunction: &parser.Disjunction{
																Conjunction: &parser.Conjunction{
																	Equality: &parser.Equality{
																		Comparison: &parser.Comparison{
																			BitwiseDisjunction: &parser.BitwiseDisjunction{
																				BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																					BitwiseConjunction: &parser.BitwiseConjunction{
																						Shift: &parser.Shift{
																							Addition: &parser.Addition{
																								Multiplication: &parser.Multiplication{
																									Unary: &parser.Unary{
																										Accessor: &parser.Accessor{
																											Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
													Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
												},
												{
													Condition: &parser.Expression{
														ListConstruction: &parser.ListConstruction{
															Disjunction: &parser.Disjunction{
																Conjunction: &parser.Conjunction{
																	Equality: &parser.Equality{
																		Comparison: &parser.Comparison{
																			BitwiseDisjunction: &parser.BitwiseDisjunction{
																				BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																					BitwiseConjunction: &parser.BitwiseConjunction{
																						Shift: &parser.Shift{
																							Addition: &parser.Addition{
																								Multiplication: &parser.Multiplication{
																									Unary: &parser.Unary{
																										Accessor: &parser.Accessor{
																											Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
													Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
												},
											},
										},
									},
								},
							},
						},
					},
					Operation: "<<",
					Shift: &parser.Shift{
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{
										Atom: &parser.Atom{
											ConditionalExpression: &parser.ConditionalExpression{
												ConditionalCases: []*parser.ConditionalCase{
													{
														Condition: &parser.Expression{
															ListConstruction: &parser.ListConstruction{
																Disjunction: &parser.Disjunction{
																	Conjunction: &parser.Conjunction{
																		Equality: &parser.Equality{
																			Comparison: &parser.Comparison{
																				BitwiseDisjunction: &parser.BitwiseDisjunction{
																					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																						BitwiseConjunction: &parser.BitwiseConjunction{
																							Shift: &parser.Shift{
																								Addition: &parser.Addition{
																									Multiplication: &parser.Multiplication{
																										Unary: &parser.Unary{
																											Accessor: &parser.Accessor{
																												Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
														Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
													},
													{
														Condition: &parser.Expression{
															ListConstruction: &parser.ListConstruction{
																Disjunction: &parser.Disjunction{
																	Conjunction: &parser.Conjunction{
																		Equality: &parser.Equality{
																			Comparison: &parser.Comparison{
																				BitwiseDisjunction: &parser.BitwiseDisjunction{
																					BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																						BitwiseConjunction: &parser.BitwiseConjunction{
																							Shift: &parser.Shift{
																								Addition: &parser.Addition{
																									Multiplication: &parser.Multiplication{
																										Unary: &parser.Unary{
																											Accessor: &parser.Accessor{
																												Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
														Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				shift: &parser.Shift{
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(5)}},
							},
						},
					},
					Operation: "<<",
					Shift: &parser.Shift{
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
								},
							},
						},
						Operation: ">>",
						Shift: &parser.Shift{
							Addition: &parser.Addition{
								Multiplication: &parser.Multiplication{
									Unary: &parser.Unary{
										Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
									},
								},
							},
							Operation: ">>>",
							Shift: &parser.Shift{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: pointer.ToString("unknown")}},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "Shift/empty/success",
			args: args{
				shift: &parser.Shift{
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Shift/empty/success/with setted states",
			args: args{
				shift: &parser.Shift{
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{
									Atom: &parser.Atom{
										ConditionalExpression: &parser.ConditionalExpression{
											ConditionalCases: []*parser.ConditionalCase{
												{
													Condition: &parser.Expression{
														ListConstruction: &parser.ListConstruction{
															Disjunction: &parser.Disjunction{
																Conjunction: &parser.Conjunction{
																	Equality: &parser.Equality{
																		Comparison: &parser.Comparison{
																			BitwiseDisjunction: &parser.BitwiseDisjunction{
																				BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																					BitwiseConjunction: &parser.BitwiseConjunction{
																						Shift: &parser.Shift{
																							Addition: &parser.Addition{
																								Multiplication: &parser.Multiplication{
																									Unary: &parser.Unary{
																										Accessor: &parser.Accessor{
																											Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
													Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
												},
												{
													Condition: &parser.Expression{
														ListConstruction: &parser.ListConstruction{
															Disjunction: &parser.Disjunction{
																Conjunction: &parser.Conjunction{
																	Equality: &parser.Equality{
																		Comparison: &parser.Comparison{
																			BitwiseDisjunction: &parser.BitwiseDisjunction{
																				BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																					BitwiseConjunction: &parser.BitwiseConjunction{
																						Shift: &parser.Shift{
																							Addition: &parser.Addition{
																								Multiplication: &parser.Multiplication{
																									Unary: &parser.Unary{
																										Accessor: &parser.Accessor{
																											Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
													Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
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
				shift: &parser.Shift{
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: pointer.ToString("unknown")}},
							},
						},
					},
				},
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
				addition: &parser.Addition{
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
						},
					},
					Operation: "+",
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
							},
						},
						Operation: "+",
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
								},
							},
						},
					},
				},
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
				addition: &parser.Addition{
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
						},
					},
					Operation: "-",
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
							},
						},
						Operation: "-",
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
								},
							},
						},
					},
				},
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
				addition: &parser.Addition{
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{
								Atom: &parser.Atom{
									ConditionalExpression: &parser.ConditionalExpression{
										ConditionalCases: []*parser.ConditionalCase{
											{
												Condition: &parser.Expression{
													ListConstruction: &parser.ListConstruction{
														Disjunction: &parser.Disjunction{
															Conjunction: &parser.Conjunction{
																Equality: &parser.Equality{
																	Comparison: &parser.Comparison{
																		BitwiseDisjunction: &parser.BitwiseDisjunction{
																			BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																				BitwiseConjunction: &parser.BitwiseConjunction{
																					Shift: &parser.Shift{
																						Addition: &parser.Addition{
																							Multiplication: &parser.Multiplication{
																								Unary: &parser.Unary{
																									Accessor: &parser.Accessor{
																										Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
												Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
											},
											{
												Condition: &parser.Expression{
													ListConstruction: &parser.ListConstruction{
														Disjunction: &parser.Disjunction{
															Conjunction: &parser.Conjunction{
																Equality: &parser.Equality{
																	Comparison: &parser.Comparison{
																		BitwiseDisjunction: &parser.BitwiseDisjunction{
																			BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																				BitwiseConjunction: &parser.BitwiseConjunction{
																					Shift: &parser.Shift{
																						Addition: &parser.Addition{
																							Multiplication: &parser.Multiplication{
																								Unary: &parser.Unary{
																									Accessor: &parser.Accessor{
																										Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
												Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
											},
										},
									},
								},
							},
						},
					},
					Operation: "+",
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{
									Atom: &parser.Atom{
										ConditionalExpression: &parser.ConditionalExpression{
											ConditionalCases: []*parser.ConditionalCase{
												{
													Condition: &parser.Expression{
														ListConstruction: &parser.ListConstruction{
															Disjunction: &parser.Disjunction{
																Conjunction: &parser.Conjunction{
																	Equality: &parser.Equality{
																		Comparison: &parser.Comparison{
																			BitwiseDisjunction: &parser.BitwiseDisjunction{
																				BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																					BitwiseConjunction: &parser.BitwiseConjunction{
																						Shift: &parser.Shift{
																							Addition: &parser.Addition{
																								Multiplication: &parser.Multiplication{
																									Unary: &parser.Unary{
																										Accessor: &parser.Accessor{
																											Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
													Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
												},
												{
													Condition: &parser.Expression{
														ListConstruction: &parser.ListConstruction{
															Disjunction: &parser.Disjunction{
																Conjunction: &parser.Conjunction{
																	Equality: &parser.Equality{
																		Comparison: &parser.Comparison{
																			BitwiseDisjunction: &parser.BitwiseDisjunction{
																				BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																					BitwiseConjunction: &parser.BitwiseConjunction{
																						Shift: &parser.Shift{
																							Addition: &parser.Addition{
																								Multiplication: &parser.Multiplication{
																									Unary: &parser.Unary{
																										Accessor: &parser.Accessor{
																											Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
													Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
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
				addition: &parser.Addition{
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
						},
					},
					Operation: "+",
					Addition: &parser.Addition{
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
							},
						},
						Operation: "+",
						Addition: &parser.Addition{
							Multiplication: &parser.Multiplication{
								Unary: &parser.Unary{
									Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: pointer.ToString("unknown")}},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
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
							Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Addition/empty/success/with setted states",
			args: args{
				addition: &parser.Addition{
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{
								Atom: &parser.Atom{
									ConditionalExpression: &parser.ConditionalExpression{
										ConditionalCases: []*parser.ConditionalCase{
											{
												Condition: &parser.Expression{
													ListConstruction: &parser.ListConstruction{
														Disjunction: &parser.Disjunction{
															Conjunction: &parser.Conjunction{
																Equality: &parser.Equality{
																	Comparison: &parser.Comparison{
																		BitwiseDisjunction: &parser.BitwiseDisjunction{
																			BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																				BitwiseConjunction: &parser.BitwiseConjunction{
																					Shift: &parser.Shift{
																						Addition: &parser.Addition{
																							Multiplication: &parser.Multiplication{
																								Unary: &parser.Unary{
																									Accessor: &parser.Accessor{
																										Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
												Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
											},
											{
												Condition: &parser.Expression{
													ListConstruction: &parser.ListConstruction{
														Disjunction: &parser.Disjunction{
															Conjunction: &parser.Conjunction{
																Equality: &parser.Equality{
																	Comparison: &parser.Comparison{
																		BitwiseDisjunction: &parser.BitwiseDisjunction{
																			BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																				BitwiseConjunction: &parser.BitwiseConjunction{
																					Shift: &parser.Shift{
																						Addition: &parser.Addition{
																							Multiplication: &parser.Multiplication{
																								Unary: &parser.Unary{
																									Accessor: &parser.Accessor{
																										Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
												Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
											},
										},
									},
								},
							},
						},
					},
				},
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
				addition: &parser.Addition{
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: pointer.ToString("unknown")}},
						},
					},
				},
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
				multiplication: &parser.Multiplication{
					Unary: &parser.Unary{
						Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
					},
					Operation: "*",
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
						},
						Operation: "*",
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
							},
						},
					},
				},
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
				multiplication: &parser.Multiplication{
					Unary: &parser.Unary{
						Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
					},
					Operation: "/",
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
						},
						Operation: "/",
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
							},
						},
					},
				},
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
				multiplication: &parser.Multiplication{
					Unary: &parser.Unary{
						Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
					},
					Operation: "%",
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
						},
						Operation: "%",
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)}},
							},
						},
					},
				},
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
				multiplication: &parser.Multiplication{
					Unary: &parser.Unary{
						Accessor: &parser.Accessor{
							Atom: &parser.Atom{
								ConditionalExpression: &parser.ConditionalExpression{
									ConditionalCases: []*parser.ConditionalCase{
										{
											Condition: &parser.Expression{
												ListConstruction: &parser.ListConstruction{
													Disjunction: &parser.Disjunction{
														Conjunction: &parser.Conjunction{
															Equality: &parser.Equality{
																Comparison: &parser.Comparison{
																	BitwiseDisjunction: &parser.BitwiseDisjunction{
																		BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																			BitwiseConjunction: &parser.BitwiseConjunction{
																				Shift: &parser.Shift{
																					Addition: &parser.Addition{
																						Multiplication: &parser.Multiplication{
																							Unary: &parser.Unary{
																								Accessor: &parser.Accessor{
																									Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
											Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
										},
										{
											Condition: &parser.Expression{
												ListConstruction: &parser.ListConstruction{
													Disjunction: &parser.Disjunction{
														Conjunction: &parser.Conjunction{
															Equality: &parser.Equality{
																Comparison: &parser.Comparison{
																	BitwiseDisjunction: &parser.BitwiseDisjunction{
																		BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																			BitwiseConjunction: &parser.BitwiseConjunction{
																				Shift: &parser.Shift{
																					Addition: &parser.Addition{
																						Multiplication: &parser.Multiplication{
																							Unary: &parser.Unary{
																								Accessor: &parser.Accessor{
																									Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
											Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
										},
									},
								},
							},
						},
					},
					Operation: "*",
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{
								Atom: &parser.Atom{
									ConditionalExpression: &parser.ConditionalExpression{
										ConditionalCases: []*parser.ConditionalCase{
											{
												Condition: &parser.Expression{
													ListConstruction: &parser.ListConstruction{
														Disjunction: &parser.Disjunction{
															Conjunction: &parser.Conjunction{
																Equality: &parser.Equality{
																	Comparison: &parser.Comparison{
																		BitwiseDisjunction: &parser.BitwiseDisjunction{
																			BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																				BitwiseConjunction: &parser.BitwiseConjunction{
																					Shift: &parser.Shift{
																						Addition: &parser.Addition{
																							Multiplication: &parser.Multiplication{
																								Unary: &parser.Unary{
																									Accessor: &parser.Accessor{
																										Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
												Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
											},
											{
												Condition: &parser.Expression{
													ListConstruction: &parser.ListConstruction{
														Disjunction: &parser.Disjunction{
															Conjunction: &parser.Conjunction{
																Equality: &parser.Equality{
																	Comparison: &parser.Comparison{
																		BitwiseDisjunction: &parser.BitwiseDisjunction{
																			BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																				BitwiseConjunction: &parser.BitwiseConjunction{
																					Shift: &parser.Shift{
																						Addition: &parser.Addition{
																							Multiplication: &parser.Multiplication{
																								Unary: &parser.Unary{
																									Accessor: &parser.Accessor{
																										Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
												Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
											},
										},
									},
								},
							},
						},
					},
				},
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
				multiplication: &parser.Multiplication{
					Unary: &parser.Unary{
						Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)}},
					},
					Operation: "*",
					Multiplication: &parser.Multiplication{
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
						},
						Operation: "*",
						Multiplication: &parser.Multiplication{
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: pointer.ToString("unknown")}},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Multiplication/empty/success",
			args: args{
				multiplication: &parser.Multiplication{
					Unary: &parser.Unary{
						Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Multiplication/empty/success/with setted states",
			args: args{
				multiplication: &parser.Multiplication{
					Unary: &parser.Unary{
						Accessor: &parser.Accessor{
							Atom: &parser.Atom{
								ConditionalExpression: &parser.ConditionalExpression{
									ConditionalCases: []*parser.ConditionalCase{
										{
											Condition: &parser.Expression{
												ListConstruction: &parser.ListConstruction{
													Disjunction: &parser.Disjunction{
														Conjunction: &parser.Conjunction{
															Equality: &parser.Equality{
																Comparison: &parser.Comparison{
																	BitwiseDisjunction: &parser.BitwiseDisjunction{
																		BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																			BitwiseConjunction: &parser.BitwiseConjunction{
																				Shift: &parser.Shift{
																					Addition: &parser.Addition{
																						Multiplication: &parser.Multiplication{
																							Unary: &parser.Unary{
																								Accessor: &parser.Accessor{
																									Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
											Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
										},
										{
											Condition: &parser.Expression{
												ListConstruction: &parser.ListConstruction{
													Disjunction: &parser.Disjunction{
														Conjunction: &parser.Conjunction{
															Equality: &parser.Equality{
																Comparison: &parser.Comparison{
																	BitwiseDisjunction: &parser.BitwiseDisjunction{
																		BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																			BitwiseConjunction: &parser.BitwiseConjunction{
																				Shift: &parser.Shift{
																					Addition: &parser.Addition{
																						Multiplication: &parser.Multiplication{
																							Unary: &parser.Unary{
																								Accessor: &parser.Accessor{
																									Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
											Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
										},
									},
								},
							},
						},
					},
				},
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
				multiplication: &parser.Multiplication{
					Unary: &parser.Unary{
						Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: pointer.ToString("unknown")}},
					},
				},
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
				unary: &parser.Unary{
					Operation: "-",
					Unary: &parser.Unary{
						Operation: "~",
						Unary: &parser.Unary{
							Operation: "!",
							Unary: &parser.Unary{
								Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
							},
						},
					},
				},
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
				unary: &parser.Unary{
					Operation: "-",
					Unary: &parser.Unary{
						Operation: "!",
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{
								Atom: &parser.Atom{
									ConditionalExpression: &parser.ConditionalExpression{
										ConditionalCases: []*parser.ConditionalCase{
											{
												Condition: &parser.Expression{
													ListConstruction: &parser.ListConstruction{
														Disjunction: &parser.Disjunction{
															Conjunction: &parser.Conjunction{
																Equality: &parser.Equality{
																	Comparison: &parser.Comparison{
																		BitwiseDisjunction: &parser.BitwiseDisjunction{
																			BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																				BitwiseConjunction: &parser.BitwiseConjunction{
																					Shift: &parser.Shift{
																						Addition: &parser.Addition{
																							Multiplication: &parser.Multiplication{
																								Unary: &parser.Unary{
																									Accessor: &parser.Accessor{
																										Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
												Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
											},
											{
												Condition: &parser.Expression{
													ListConstruction: &parser.ListConstruction{
														Disjunction: &parser.Disjunction{
															Conjunction: &parser.Conjunction{
																Equality: &parser.Equality{
																	Comparison: &parser.Comparison{
																		BitwiseDisjunction: &parser.BitwiseDisjunction{
																			BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																				BitwiseConjunction: &parser.BitwiseConjunction{
																					Shift: &parser.Shift{
																						Addition: &parser.Addition{
																							Multiplication: &parser.Multiplication{
																								Unary: &parser.Unary{
																									Accessor: &parser.Accessor{
																										Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
												Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: expressions.NewFunctionCall(
				ArithmeticNegationFunctionName,
				[]expressions.Expression{
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
				},
			),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Unary/nonempty/error",
			args: args{
				unary: &parser.Unary{
					Operation: "-",
					Unary: &parser.Unary{
						Operation: "!",
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: pointer.ToString("unknown")}},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Unary/empty/success",
			args: args{
				unary: &parser.Unary{
					Accessor: &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Unary/empty/success/with setted states",
			args: args{
				unary: &parser.Unary{
					Accessor: &parser.Accessor{
						Atom: &parser.Atom{
							ConditionalExpression: &parser.ConditionalExpression{
								ConditionalCases: []*parser.ConditionalCase{
									{
										Condition: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
																BitwiseDisjunction: &parser.BitwiseDisjunction{
																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																		BitwiseConjunction: &parser.BitwiseConjunction{
																			Shift: &parser.Shift{
																				Addition: &parser.Addition{
																					Multiplication: &parser.Multiplication{
																						Unary: &parser.Unary{
																							Accessor: &parser.Accessor{
																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
										Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
									},
									{
										Condition: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
																BitwiseDisjunction: &parser.BitwiseDisjunction{
																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																		BitwiseConjunction: &parser.BitwiseConjunction{
																			Shift: &parser.Shift{
																				Addition: &parser.Addition{
																					Multiplication: &parser.Multiplication{
																						Unary: &parser.Unary{
																							Accessor: &parser.Accessor{
																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
										Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
									},
								},
							},
						},
					},
				},
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
				unary: &parser.Unary{
					Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: pointer.ToString("unknown")}},
				},
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
			name: "Accessor/nonempty/success",
			args: args{
				accessor: &parser.Accessor{
					Atom: &parser.Atom{Identifier: pointer.ToString("test")},
					Keys: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
			name: "Accessor/nonempty/success/with setted states",
			args: args{
				accessor: &parser.Accessor{
					Atom: &parser.Atom{
						ConditionalExpression: &parser.ConditionalExpression{
							ConditionalCases: []*parser.ConditionalCase{
								{
									Condition: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
								},
								{
									Condition: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
								},
							},
						},
					},
					Keys: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{
																					ConditionalExpression: &parser.ConditionalExpression{
																						ConditionalCases: []*parser.ConditionalCase{
																							{
																								Condition: &parser.Expression{
																									ListConstruction: &parser.ListConstruction{
																										Disjunction: &parser.Disjunction{
																											Conjunction: &parser.Conjunction{
																												Equality: &parser.Equality{
																													Comparison: &parser.Comparison{
																														BitwiseDisjunction: &parser.BitwiseDisjunction{
																															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																BitwiseConjunction: &parser.BitwiseConjunction{
																																	Shift: &parser.Shift{
																																		Addition: &parser.Addition{
																																			Multiplication: &parser.Multiplication{
																																				Unary: &parser.Unary{
																																					Accessor: &parser.Accessor{
																																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
																								Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																							},
																							{
																								Condition: &parser.Expression{
																									ListConstruction: &parser.ListConstruction{
																										Disjunction: &parser.Disjunction{
																											Conjunction: &parser.Conjunction{
																												Equality: &parser.Equality{
																													Comparison: &parser.Comparison{
																														BitwiseDisjunction: &parser.BitwiseDisjunction{
																															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																BitwiseConjunction: &parser.BitwiseConjunction{
																																	Shift: &parser.Shift{
																																		Addition: &parser.Addition{
																																			Multiplication: &parser.Multiplication{
																																				Unary: &parser.Unary{
																																					Accessor: &parser.Accessor{
																																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
																								Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				accessor: &parser.Accessor{
					Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
					Keys: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Accessor/nonempty/error/key translating",
			args: args{
				accessor: &parser.Accessor{
					Atom: &parser.Atom{Identifier: pointer.ToString("test")},
					Keys: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Accessor/empty/success",
			args: args{
				accessor:            &parser.Accessor{Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)}},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Accessor/empty/success/with setted states",
			args: args{
				accessor: &parser.Accessor{
					Atom: &parser.Atom{
						ConditionalExpression: &parser.ConditionalExpression{
							ConditionalCases: []*parser.ConditionalCase{
								{
									Condition: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
								},
								{
									Condition: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
								},
							},
						},
					},
				},
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
				accessor: &parser.Accessor{
					Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
				},
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
				atom:                &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/number/floating-point",
			args: args{
				atom:                &parser.Atom{FloatingPointNumber: pointer.ToFloat64(2.3)},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(2.3),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/symbol/latin1",
			args: args{
				atom:                &parser.Atom{Symbol: pointer.ToString("t")},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(116),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/symbol/not latin1",
			args: args{
				atom:                &parser.Atom{Symbol: pointer.ToString("т")},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(1090),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/string",
			args: args{
				atom:                &parser.Atom{String: pointer.ToString("test")},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewString("test"),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/list definition/success",
			args: args{
				atom: &parser.Atom{
					ListDefinition: &parser.ListDefinition{
						Items: []*parser.Expression{
							{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
				atom: &parser.Atom{
					ListDefinition: &parser.ListDefinition{
						Items: []*parser.Expression{
							{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{
																						ConditionalExpression: &parser.ConditionalExpression{
																							ConditionalCases: []*parser.ConditionalCase{
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											Disjunction: &parser.Disjunction{
																												Conjunction: &parser.Conjunction{
																													Equality: &parser.Equality{
																														Comparison: &parser.Comparison{
																															BitwiseDisjunction: &parser.BitwiseDisjunction{
																																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																	BitwiseConjunction: &parser.BitwiseConjunction{
																																		Shift: &parser.Shift{
																																			Addition: &parser.Addition{
																																				Multiplication: &parser.Multiplication{
																																					Unary: &parser.Unary{
																																						Accessor: &parser.Accessor{
																																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											Disjunction: &parser.Disjunction{
																												Conjunction: &parser.Conjunction{
																													Equality: &parser.Equality{
																														Comparison: &parser.Comparison{
																															BitwiseDisjunction: &parser.BitwiseDisjunction{
																																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																	BitwiseConjunction: &parser.BitwiseConjunction{
																																		Shift: &parser.Shift{
																																			Addition: &parser.Addition{
																																				Multiplication: &parser.Multiplication{
																																					Unary: &parser.Unary{
																																						Accessor: &parser.Accessor{
																																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
					expressions.NewIdentifier(EmptyListConstantName),
				},
			),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/list definition/error",
			args: args{
				atom: &parser.Atom{
					ListDefinition: &parser.ListDefinition{
						Items: []*parser.Expression{
							{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Atom/hash table definition/success",
			args: args{
				atom: &parser.Atom{
					HashTableDefinition: &parser.HashTableDefinition{
						Entries: []*parser.HashTableEntry{
							{
								Name: pointer.ToString("x"),
								Value: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
								Name: pointer.ToString("y"),
								Value: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
								Name: pointer.ToString("z"),
								Value: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
				atom: &parser.Atom{
					HashTableDefinition: &parser.HashTableDefinition{
						Entries: []*parser.HashTableEntry{
							{
								Name: pointer.ToString("x"),
								Value: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{
																							ConditionalExpression: &parser.ConditionalExpression{
																								ConditionalCases: []*parser.ConditionalCase{
																									{
																										Condition: &parser.Expression{
																											ListConstruction: &parser.ListConstruction{
																												Disjunction: &parser.Disjunction{
																													Conjunction: &parser.Conjunction{
																														Equality: &parser.Equality{
																															Comparison: &parser.Comparison{
																																BitwiseDisjunction: &parser.BitwiseDisjunction{
																																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																		BitwiseConjunction: &parser.BitwiseConjunction{
																																			Shift: &parser.Shift{
																																				Addition: &parser.Addition{
																																					Multiplication: &parser.Multiplication{
																																						Unary: &parser.Unary{
																																							Accessor: &parser.Accessor{
																																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																										Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																									},
																									{
																										Condition: &parser.Expression{
																											ListConstruction: &parser.ListConstruction{
																												Disjunction: &parser.Disjunction{
																													Conjunction: &parser.Conjunction{
																														Equality: &parser.Equality{
																															Comparison: &parser.Comparison{
																																BitwiseDisjunction: &parser.BitwiseDisjunction{
																																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																		BitwiseConjunction: &parser.BitwiseConjunction{
																																			Shift: &parser.Shift{
																																				Addition: &parser.Addition{
																																					Multiplication: &parser.Multiplication{
																																						Unary: &parser.Unary{
																																							Accessor: &parser.Accessor{
																																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																										Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
							{
								Name: pointer.ToString("y"),
								Value: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{
																							ConditionalExpression: &parser.ConditionalExpression{
																								ConditionalCases: []*parser.ConditionalCase{
																									{
																										Condition: &parser.Expression{
																											ListConstruction: &parser.ListConstruction{
																												Disjunction: &parser.Disjunction{
																													Conjunction: &parser.Conjunction{
																														Equality: &parser.Equality{
																															Comparison: &parser.Comparison{
																																BitwiseDisjunction: &parser.BitwiseDisjunction{
																																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																		BitwiseConjunction: &parser.BitwiseConjunction{
																																			Shift: &parser.Shift{
																																				Addition: &parser.Addition{
																																					Multiplication: &parser.Multiplication{
																																						Unary: &parser.Unary{
																																							Accessor: &parser.Accessor{
																																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
																										Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																									},
																									{
																										Condition: &parser.Expression{
																											ListConstruction: &parser.ListConstruction{
																												Disjunction: &parser.Disjunction{
																													Conjunction: &parser.Conjunction{
																														Equality: &parser.Equality{
																															Comparison: &parser.Comparison{
																																BitwiseDisjunction: &parser.BitwiseDisjunction{
																																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																		BitwiseConjunction: &parser.BitwiseConjunction{
																																			Shift: &parser.Shift{
																																				Addition: &parser.Addition{
																																					Multiplication: &parser.Multiplication{
																																						Unary: &parser.Unary{
																																							Accessor: &parser.Accessor{
																																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
																										Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				atom: &parser.Atom{
					HashTableDefinition: &parser.HashTableDefinition{
						Entries: []*parser.HashTableEntry{
							{
								Name: pointer.ToString("x"),
								Value: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
								Name: pointer.ToString("y"),
								Value: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
								Name: pointer.ToString("z"),
								Value: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
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
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
				atom: &parser.Atom{
					FunctionCall: &parser.FunctionCall{
						Name: "test",
						Arguments: []*parser.Expression{
							{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{
																						ConditionalExpression: &parser.ConditionalExpression{
																							ConditionalCases: []*parser.ConditionalCase{
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											Disjunction: &parser.Disjunction{
																												Conjunction: &parser.Conjunction{
																													Equality: &parser.Equality{
																														Comparison: &parser.Comparison{
																															BitwiseDisjunction: &parser.BitwiseDisjunction{
																																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																	BitwiseConjunction: &parser.BitwiseConjunction{
																																		Shift: &parser.Shift{
																																			Addition: &parser.Addition{
																																				Multiplication: &parser.Multiplication{
																																					Unary: &parser.Unary{
																																						Accessor: &parser.Accessor{
																																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											Disjunction: &parser.Disjunction{
																												Conjunction: &parser.Conjunction{
																													Equality: &parser.Equality{
																														Comparison: &parser.Comparison{
																															BitwiseDisjunction: &parser.BitwiseDisjunction{
																																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																	BitwiseConjunction: &parser.BitwiseConjunction{
																																		Shift: &parser.Shift{
																																			Addition: &parser.Addition{
																																				Multiplication: &parser.Multiplication{
																																					Unary: &parser.Unary{
																																						Accessor: &parser.Accessor{
																																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
			}),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
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
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Atom/conditional expression/success",
			args: args{
				atom: &parser.Atom{
					ConditionalExpression: &parser.ConditionalExpression{
						ConditionalCases: []*parser.ConditionalCase{
							{
								Condition: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
								Commands: []*parser.Command{
									{
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
																BitwiseDisjunction: &parser.BitwiseDisjunction{
																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																		BitwiseConjunction: &parser.BitwiseConjunction{
																			Shift: &parser.Shift{
																				Addition: &parser.Addition{
																					Multiplication: &parser.Multiplication{
																						Unary: &parser.Unary{
																							Accessor: &parser.Accessor{
																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
																BitwiseDisjunction: &parser.BitwiseDisjunction{
																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																		BitwiseConjunction: &parser.BitwiseConjunction{
																			Shift: &parser.Shift{
																				Addition: &parser.Addition{
																					Multiplication: &parser.Multiplication{
																						Unary: &parser.Unary{
																							Accessor: &parser.Accessor{
																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
							{
								Condition: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(13)},
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
								Commands: []*parser.Command{
									{
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
																BitwiseDisjunction: &parser.BitwiseDisjunction{
																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																		BitwiseConjunction: &parser.BitwiseConjunction{
																			Shift: &parser.Shift{
																				Addition: &parser.Addition{
																					Multiplication: &parser.Multiplication{
																						Unary: &parser.Unary{
																							Accessor: &parser.Accessor{
																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
																BitwiseDisjunction: &parser.BitwiseDisjunction{
																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																		BitwiseConjunction: &parser.BitwiseConjunction{
																			Shift: &parser.Shift{
																				Addition: &parser.Addition{
																					Multiplication: &parser.Multiplication{
																						Unary: &parser.Unary{
																							Accessor: &parser.Accessor{
																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
							{
								Condition: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(14)},
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
								Commands: []*parser.Command{
									{
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
																BitwiseDisjunction: &parser.BitwiseDisjunction{
																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																		BitwiseConjunction: &parser.BitwiseConjunction{
																			Shift: &parser.Shift{
																				Addition: &parser.Addition{
																					Multiplication: &parser.Multiplication{
																						Unary: &parser.Unary{
																							Accessor: &parser.Accessor{
																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(25)},
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
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
																BitwiseDisjunction: &parser.BitwiseDisjunction{
																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																		BitwiseConjunction: &parser.BitwiseConjunction{
																			Shift: &parser.Shift{
																				Addition: &parser.Addition{
																					Multiplication: &parser.Multiplication{
																						Unary: &parser.Unary{
																							Accessor: &parser.Accessor{
																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(44)},
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
				atom: &parser.Atom{
					ConditionalExpression: &parser.ConditionalExpression{
						ConditionalCases: []*parser.ConditionalCase{
							{
								Condition: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
								Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
							},
							{
								Condition: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
								Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
							},
						},
					},
				},
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
				atom: &parser.Atom{
					ConditionalExpression: &parser.ConditionalExpression{
						ConditionalCases: []*parser.ConditionalCase{
							{
								Condition: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
								Commands: []*parser.Command{
									{
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
																BitwiseDisjunction: &parser.BitwiseDisjunction{
																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																		BitwiseConjunction: &parser.BitwiseConjunction{
																			Shift: &parser.Shift{
																				Addition: &parser.Addition{
																					Multiplication: &parser.Multiplication{
																						Unary: &parser.Unary{
																							Accessor: &parser.Accessor{
																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
																BitwiseDisjunction: &parser.BitwiseDisjunction{
																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																		BitwiseConjunction: &parser.BitwiseConjunction{
																			Shift: &parser.Shift{
																				Addition: &parser.Addition{
																					Multiplication: &parser.Multiplication{
																						Unary: &parser.Unary{
																							Accessor: &parser.Accessor{
																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
							{
								Condition: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(13)},
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
								Commands: []*parser.Command{
									{
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
																BitwiseDisjunction: &parser.BitwiseDisjunction{
																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																		BitwiseConjunction: &parser.BitwiseConjunction{
																			Shift: &parser.Shift{
																				Addition: &parser.Addition{
																					Multiplication: &parser.Multiplication{
																						Unary: &parser.Unary{
																							Accessor: &parser.Accessor{
																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
																BitwiseDisjunction: &parser.BitwiseDisjunction{
																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																		BitwiseConjunction: &parser.BitwiseConjunction{
																			Shift: &parser.Shift{
																				Addition: &parser.Addition{
																					Multiplication: &parser.Multiplication{
																						Unary: &parser.Unary{
																							Accessor: &parser.Accessor{
																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
							{
								Condition: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
														BitwiseDisjunction: &parser.BitwiseDisjunction{
															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																BitwiseConjunction: &parser.BitwiseConjunction{
																	Shift: &parser.Shift{
																		Addition: &parser.Addition{
																			Multiplication: &parser.Multiplication{
																				Unary: &parser.Unary{
																					Accessor: &parser.Accessor{
																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(14)},
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
								Commands: []*parser.Command{
									{
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
																BitwiseDisjunction: &parser.BitwiseDisjunction{
																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																		BitwiseConjunction: &parser.BitwiseConjunction{
																			Shift: &parser.Shift{
																				Addition: &parser.Addition{
																					Multiplication: &parser.Multiplication{
																						Unary: &parser.Unary{
																							Accessor: &parser.Accessor{
																								Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(25)},
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
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
																BitwiseDisjunction: &parser.BitwiseDisjunction{
																	BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																		BitwiseConjunction: &parser.BitwiseConjunction{
																			Shift: &parser.Shift{
																				Addition: &parser.Addition{
																					Multiplication: &parser.Multiplication{
																						Unary: &parser.Unary{
																							Accessor: &parser.Accessor{
																								Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Atom/identifier/success",
			args: args{
				atom:                &parser.Atom{Identifier: pointer.ToString("test")},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewIdentifier("test"),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/identifier/error",
			args: args{
				atom:                &parser.Atom{Identifier: pointer.ToString("unknown")},
				declaredIdentifiers: mapset.NewSet("test"),
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
							Disjunction: &parser.Disjunction{
								Conjunction: &parser.Conjunction{
									Equality: &parser.Equality{
										Comparison: &parser.Comparison{
											BitwiseDisjunction: &parser.BitwiseDisjunction{
												BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
													BitwiseConjunction: &parser.BitwiseConjunction{
														Shift: &parser.Shift{
															Addition: &parser.Addition{
																Multiplication: &parser.Multiplication{
																	Unary: &parser.Unary{
																		Accessor: &parser.Accessor{
																			Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewNumber(23),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "Atom/expression/success/with setted states",
			args: args{
				atom: &parser.Atom{
					Expression: &parser.Expression{
						ListConstruction: &parser.ListConstruction{
							Disjunction: &parser.Disjunction{
								Conjunction: &parser.Conjunction{
									Equality: &parser.Equality{
										Comparison: &parser.Comparison{
											BitwiseDisjunction: &parser.BitwiseDisjunction{
												BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
													BitwiseConjunction: &parser.BitwiseConjunction{
														Shift: &parser.Shift{
															Addition: &parser.Addition{
																Multiplication: &parser.Multiplication{
																	Unary: &parser.Unary{
																		Accessor: &parser.Accessor{
																			Atom: &parser.Atom{
																				ConditionalExpression: &parser.ConditionalExpression{
																					ConditionalCases: []*parser.ConditionalCase{
																						{
																							Condition: &parser.Expression{
																								ListConstruction: &parser.ListConstruction{
																									Disjunction: &parser.Disjunction{
																										Conjunction: &parser.Conjunction{
																											Equality: &parser.Equality{
																												Comparison: &parser.Comparison{
																													BitwiseDisjunction: &parser.BitwiseDisjunction{
																														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																															BitwiseConjunction: &parser.BitwiseConjunction{
																																Shift: &parser.Shift{
																																	Addition: &parser.Addition{
																																		Multiplication: &parser.Multiplication{
																																			Unary: &parser.Unary{
																																				Accessor: &parser.Accessor{
																																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																							Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																						},
																						{
																							Condition: &parser.Expression{
																								ListConstruction: &parser.ListConstruction{
																									Disjunction: &parser.Disjunction{
																										Conjunction: &parser.Conjunction{
																											Equality: &parser.Equality{
																												Comparison: &parser.Comparison{
																													BitwiseDisjunction: &parser.BitwiseDisjunction{
																														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																															BitwiseConjunction: &parser.BitwiseConjunction{
																																Shift: &parser.Shift{
																																	Addition: &parser.Addition{
																																		Multiplication: &parser.Multiplication{
																																			Unary: &parser.Unary{
																																				Accessor: &parser.Accessor{
																																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																							Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
				atom: &parser.Atom{
					Expression: &parser.Expression{
						ListConstruction: &parser.ListConstruction{
							Disjunction: &parser.Disjunction{
								Conjunction: &parser.Conjunction{
									Equality: &parser.Equality{
										Comparison: &parser.Comparison{
											BitwiseDisjunction: &parser.BitwiseDisjunction{
												BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
													BitwiseConjunction: &parser.BitwiseConjunction{
														Shift: &parser.Shift{
															Addition: &parser.Addition{
																Multiplication: &parser.Multiplication{
																	Unary: &parser.Unary{
																		Accessor: &parser.Accessor{
																			Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				listDefinition: &parser.ListDefinition{
					Items: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
				listDefinition: &parser.ListDefinition{
					Items: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{
																					ConditionalExpression: &parser.ConditionalExpression{
																						ConditionalCases: []*parser.ConditionalCase{
																							{
																								Condition: &parser.Expression{
																									ListConstruction: &parser.ListConstruction{
																										Disjunction: &parser.Disjunction{
																											Conjunction: &parser.Conjunction{
																												Equality: &parser.Equality{
																													Comparison: &parser.Comparison{
																														BitwiseDisjunction: &parser.BitwiseDisjunction{
																															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																BitwiseConjunction: &parser.BitwiseConjunction{
																																	Shift: &parser.Shift{
																																		Addition: &parser.Addition{
																																			Multiplication: &parser.Multiplication{
																																				Unary: &parser.Unary{
																																					Accessor: &parser.Accessor{
																																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																								Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																							},
																							{
																								Condition: &parser.Expression{
																									ListConstruction: &parser.ListConstruction{
																										Disjunction: &parser.Disjunction{
																											Conjunction: &parser.Conjunction{
																												Equality: &parser.Equality{
																													Comparison: &parser.Comparison{
																														BitwiseDisjunction: &parser.BitwiseDisjunction{
																															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																BitwiseConjunction: &parser.BitwiseConjunction{
																																	Shift: &parser.Shift{
																																		Addition: &parser.Addition{
																																			Multiplication: &parser.Multiplication{
																																				Unary: &parser.Unary{
																																					Accessor: &parser.Accessor{
																																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																								Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
						{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{
																					ConditionalExpression: &parser.ConditionalExpression{
																						ConditionalCases: []*parser.ConditionalCase{
																							{
																								Condition: &parser.Expression{
																									ListConstruction: &parser.ListConstruction{
																										Disjunction: &parser.Disjunction{
																											Conjunction: &parser.Conjunction{
																												Equality: &parser.Equality{
																													Comparison: &parser.Comparison{
																														BitwiseDisjunction: &parser.BitwiseDisjunction{
																															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																BitwiseConjunction: &parser.BitwiseConjunction{
																																	Shift: &parser.Shift{
																																		Addition: &parser.Addition{
																																			Multiplication: &parser.Multiplication{
																																				Unary: &parser.Unary{
																																					Accessor: &parser.Accessor{
																																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
																								Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																							},
																							{
																								Condition: &parser.Expression{
																									ListConstruction: &parser.ListConstruction{
																										Disjunction: &parser.Disjunction{
																											Conjunction: &parser.Conjunction{
																												Equality: &parser.Equality{
																													Comparison: &parser.Comparison{
																														BitwiseDisjunction: &parser.BitwiseDisjunction{
																															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																BitwiseConjunction: &parser.BitwiseConjunction{
																																	Shift: &parser.Shift{
																																		Addition: &parser.Addition{
																																			Multiplication: &parser.Multiplication{
																																				Unary: &parser.Unary{
																																					Accessor: &parser.Accessor{
																																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
																								Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				listDefinition: &parser.ListDefinition{
					Items: nil,
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewIdentifier(EmptyListConstantName),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "ListDefinition/error",
			args: args{
				listDefinition: &parser.ListDefinition{
					Items: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				hashTableDefinition: &parser.HashTableDefinition{
					Entries: []*parser.HashTableEntry{
						{
							Name: pointer.ToString("x"),
							Value: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							Name: pointer.ToString("y"),
							Value: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
							Name: pointer.ToString("z"),
							Value: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
				hashTableDefinition: &parser.HashTableDefinition{
					Entries: []*parser.HashTableEntry{
						{
							Name: pointer.ToString("x"),
							Value: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							Name: pointer.ToString("y"),
							Value: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
							Name: pointer.ToString("z"),
							Value: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
				hashTableDefinition: &parser.HashTableDefinition{
					Entries: []*parser.HashTableEntry{
						{
							Name: pointer.ToString("x"),
							Value: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{
																						ConditionalExpression: &parser.ConditionalExpression{
																							ConditionalCases: []*parser.ConditionalCase{
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											Disjunction: &parser.Disjunction{
																												Conjunction: &parser.Conjunction{
																													Equality: &parser.Equality{
																														Comparison: &parser.Comparison{
																															BitwiseDisjunction: &parser.BitwiseDisjunction{
																																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																	BitwiseConjunction: &parser.BitwiseConjunction{
																																		Shift: &parser.Shift{
																																			Addition: &parser.Addition{
																																				Multiplication: &parser.Multiplication{
																																					Unary: &parser.Unary{
																																						Accessor: &parser.Accessor{
																																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											Disjunction: &parser.Disjunction{
																												Conjunction: &parser.Conjunction{
																													Equality: &parser.Equality{
																														Comparison: &parser.Comparison{
																															BitwiseDisjunction: &parser.BitwiseDisjunction{
																																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																	BitwiseConjunction: &parser.BitwiseConjunction{
																																		Shift: &parser.Shift{
																																			Addition: &parser.Addition{
																																				Multiplication: &parser.Multiplication{
																																					Unary: &parser.Unary{
																																						Accessor: &parser.Accessor{
																																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
						{
							Name: pointer.ToString("y"),
							Value: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{
																						ConditionalExpression: &parser.ConditionalExpression{
																							ConditionalCases: []*parser.ConditionalCase{
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											Disjunction: &parser.Disjunction{
																												Conjunction: &parser.Conjunction{
																													Equality: &parser.Equality{
																														Comparison: &parser.Comparison{
																															BitwiseDisjunction: &parser.BitwiseDisjunction{
																																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																	BitwiseConjunction: &parser.BitwiseConjunction{
																																		Shift: &parser.Shift{
																																			Addition: &parser.Addition{
																																				Multiplication: &parser.Multiplication{
																																					Unary: &parser.Unary{
																																						Accessor: &parser.Accessor{
																																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											Disjunction: &parser.Disjunction{
																												Conjunction: &parser.Conjunction{
																													Equality: &parser.Equality{
																														Comparison: &parser.Comparison{
																															BitwiseDisjunction: &parser.BitwiseDisjunction{
																																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																	BitwiseConjunction: &parser.BitwiseConjunction{
																																		Shift: &parser.Shift{
																																			Addition: &parser.Addition{
																																				Multiplication: &parser.Multiplication{
																																					Unary: &parser.Unary{
																																						Accessor: &parser.Accessor{
																																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				hashTableDefinition: &parser.HashTableDefinition{
					Entries: []*parser.HashTableEntry{
						{
							Expression: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{Identifier: pointer.ToString("test")},
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
							Value: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
				hashTableDefinition: &parser.HashTableDefinition{
					Entries: []*parser.HashTableEntry{
						{
							Expression: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{
																						ConditionalExpression: &parser.ConditionalExpression{
																							ConditionalCases: []*parser.ConditionalCase{
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											Disjunction: &parser.Disjunction{
																												Conjunction: &parser.Conjunction{
																													Equality: &parser.Equality{
																														Comparison: &parser.Comparison{
																															BitwiseDisjunction: &parser.BitwiseDisjunction{
																																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																	BitwiseConjunction: &parser.BitwiseConjunction{
																																		Shift: &parser.Shift{
																																			Addition: &parser.Addition{
																																				Multiplication: &parser.Multiplication{
																																					Unary: &parser.Unary{
																																						Accessor: &parser.Accessor{
																																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											Disjunction: &parser.Disjunction{
																												Conjunction: &parser.Conjunction{
																													Equality: &parser.Equality{
																														Comparison: &parser.Comparison{
																															BitwiseDisjunction: &parser.BitwiseDisjunction{
																																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																	BitwiseConjunction: &parser.BitwiseConjunction{
																																		Shift: &parser.Shift{
																																			Addition: &parser.Addition{
																																				Multiplication: &parser.Multiplication{
																																					Unary: &parser.Unary{
																																						Accessor: &parser.Accessor{
																																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
							Value: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
				hashTableDefinition: &parser.HashTableDefinition{
					Entries: nil,
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewIdentifier(EmptyHashTableConstantName),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "HashTableDefinition/error/unknown identifier in the expression",
			args: args{
				hashTableDefinition: &parser.HashTableDefinition{
					Entries: []*parser.HashTableEntry{
						{
							Expression: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
							Value: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "HashTableDefinition/error/unknown identifier in the value",
			args: args{
				hashTableDefinition: &parser.HashTableDefinition{
					Entries: []*parser.HashTableEntry{
						{
							Name: pointer.ToString("x"),
							Value: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							Name: pointer.ToString("y"),
							Value: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
							Name: pointer.ToString("z"),
							Value: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				functionCall: &parser.FunctionCall{
					Name: "test",
					Arguments: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
				functionCall: &parser.FunctionCall{
					Name: "test",
					Arguments: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{
																					ConditionalExpression: &parser.ConditionalExpression{
																						ConditionalCases: []*parser.ConditionalCase{
																							{
																								Condition: &parser.Expression{
																									ListConstruction: &parser.ListConstruction{
																										Disjunction: &parser.Disjunction{
																											Conjunction: &parser.Conjunction{
																												Equality: &parser.Equality{
																													Comparison: &parser.Comparison{
																														BitwiseDisjunction: &parser.BitwiseDisjunction{
																															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																BitwiseConjunction: &parser.BitwiseConjunction{
																																	Shift: &parser.Shift{
																																		Addition: &parser.Addition{
																																			Multiplication: &parser.Multiplication{
																																				Unary: &parser.Unary{
																																					Accessor: &parser.Accessor{
																																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																								Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																							},
																							{
																								Condition: &parser.Expression{
																									ListConstruction: &parser.ListConstruction{
																										Disjunction: &parser.Disjunction{
																											Conjunction: &parser.Conjunction{
																												Equality: &parser.Equality{
																													Comparison: &parser.Comparison{
																														BitwiseDisjunction: &parser.BitwiseDisjunction{
																															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																BitwiseConjunction: &parser.BitwiseConjunction{
																																	Shift: &parser.Shift{
																																		Addition: &parser.Addition{
																																			Multiplication: &parser.Multiplication{
																																				Unary: &parser.Unary{
																																					Accessor: &parser.Accessor{
																																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																								Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
						{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{
																					ConditionalExpression: &parser.ConditionalExpression{
																						ConditionalCases: []*parser.ConditionalCase{
																							{
																								Condition: &parser.Expression{
																									ListConstruction: &parser.ListConstruction{
																										Disjunction: &parser.Disjunction{
																											Conjunction: &parser.Conjunction{
																												Equality: &parser.Equality{
																													Comparison: &parser.Comparison{
																														BitwiseDisjunction: &parser.BitwiseDisjunction{
																															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																BitwiseConjunction: &parser.BitwiseConjunction{
																																	Shift: &parser.Shift{
																																		Addition: &parser.Addition{
																																			Multiplication: &parser.Multiplication{
																																				Unary: &parser.Unary{
																																					Accessor: &parser.Accessor{
																																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
																								Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																							},
																							{
																								Condition: &parser.Expression{
																									ListConstruction: &parser.ListConstruction{
																										Disjunction: &parser.Disjunction{
																											Conjunction: &parser.Conjunction{
																												Equality: &parser.Equality{
																													Comparison: &parser.Comparison{
																														BitwiseDisjunction: &parser.BitwiseDisjunction{
																															BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																BitwiseConjunction: &parser.BitwiseConjunction{
																																	Shift: &parser.Shift{
																																		Addition: &parser.Addition{
																																			Multiplication: &parser.Multiplication{
																																				Unary: &parser.Unary{
																																					Accessor: &parser.Accessor{
																																						Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
																								Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				functionCall: &parser.FunctionCall{
					Name:      "test",
					Arguments: nil,
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewFunctionCall("test", nil),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "FunctionCall/error/unknown function",
			args: args{
				functionCall: &parser.FunctionCall{
					Name: "unknown",
					Arguments: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
				declaredIdentifiers: mapset.NewSet("test"),
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
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												BitwiseDisjunction: &parser.BitwiseDisjunction{
													BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
														BitwiseConjunction: &parser.BitwiseConjunction{
															Shift: &parser.Shift{
																Addition: &parser.Addition{
																	Multiplication: &parser.Multiplication{
																		Unary: &parser.Unary{
																			Accessor: &parser.Accessor{
																				Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
				conditionalExpression: &parser.ConditionalExpression{
					ConditionalCases: []*parser.ConditionalCase{
						{
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							Commands: []*parser.Command{
								{
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
				conditionalExpression: &parser.ConditionalExpression{
					ConditionalCases: []*parser.ConditionalCase{
						{
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
				conditionalExpression: &parser.ConditionalExpression{
					ConditionalCases: []*parser.ConditionalCase{
						{
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							Commands: []*parser.Command{
								{
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
						{
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(13)},
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
							Commands: []*parser.Command{
								{
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
						{
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(14)},
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
							Commands: []*parser.Command{
								{
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(25)},
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
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(44)},
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
				conditionalExpression: &parser.ConditionalExpression{
					ConditionalCases: []*parser.ConditionalCase{
						{
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(13)},
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
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(14)},
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
				conditionalExpression: &parser.ConditionalExpression{},
				declaredIdentifiers:   mapset.NewSet("test"),
			},
			wantExpression:   expressions.NewConditionalExpression(nil),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "ConditionalExpression/success/nonempty/with setted states",
			args: args{
				conditionalExpression: &parser.ConditionalExpression{
					ConditionalCases: []*parser.ConditionalCase{
						{
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{
																						ConditionalExpression: &parser.ConditionalExpression{
																							ConditionalCases: []*parser.ConditionalCase{
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											Disjunction: &parser.Disjunction{
																												Conjunction: &parser.Conjunction{
																													Equality: &parser.Equality{
																														Comparison: &parser.Comparison{
																															BitwiseDisjunction: &parser.BitwiseDisjunction{
																																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																	BitwiseConjunction: &parser.BitwiseConjunction{
																																		Shift: &parser.Shift{
																																			Addition: &parser.Addition{
																																				Multiplication: &parser.Multiplication{
																																					Unary: &parser.Unary{
																																						Accessor: &parser.Accessor{
																																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											Disjunction: &parser.Disjunction{
																												Conjunction: &parser.Conjunction{
																													Equality: &parser.Equality{
																														Comparison: &parser.Comparison{
																															BitwiseDisjunction: &parser.BitwiseDisjunction{
																																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																	BitwiseConjunction: &parser.BitwiseConjunction{
																																		Shift: &parser.Shift{
																																			Addition: &parser.Addition{
																																				Multiplication: &parser.Multiplication{
																																					Unary: &parser.Unary{
																																						Accessor: &parser.Accessor{
																																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
							Commands: []*parser.Command{
								{
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{
																								ConditionalExpression: &parser.ConditionalExpression{
																									ConditionalCases: []*parser.ConditionalCase{
																										{
																											Condition: &parser.Expression{
																												ListConstruction: &parser.ListConstruction{
																													Disjunction: &parser.Disjunction{
																														Conjunction: &parser.Conjunction{
																															Equality: &parser.Equality{
																																Comparison: &parser.Comparison{
																																	BitwiseDisjunction: &parser.BitwiseDisjunction{
																																		BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																			BitwiseConjunction: &parser.BitwiseConjunction{
																																				Shift: &parser.Shift{
																																					Addition: &parser.Addition{
																																						Multiplication: &parser.Multiplication{
																																							Unary: &parser.Unary{
																																								Accessor: &parser.Accessor{
																																									Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
																											Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																										},
																										{
																											Condition: &parser.Expression{
																												ListConstruction: &parser.ListConstruction{
																													Disjunction: &parser.Disjunction{
																														Conjunction: &parser.Conjunction{
																															Equality: &parser.Equality{
																																Comparison: &parser.Comparison{
																																	BitwiseDisjunction: &parser.BitwiseDisjunction{
																																		BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																			BitwiseConjunction: &parser.BitwiseConjunction{
																																				Shift: &parser.Shift{
																																					Addition: &parser.Addition{
																																						Multiplication: &parser.Multiplication{
																																							Unary: &parser.Unary{
																																								Accessor: &parser.Accessor{
																																									Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
																											Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "three"}}},
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
				conditionalExpression: &parser.ConditionalExpression{
					ConditionalCases: []*parser.ConditionalCase{
						{
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{
																						ConditionalExpression: &parser.ConditionalExpression{
																							ConditionalCases: []*parser.ConditionalCase{
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											Disjunction: &parser.Disjunction{
																												Conjunction: &parser.Conjunction{
																													Equality: &parser.Equality{
																														Comparison: &parser.Comparison{
																															BitwiseDisjunction: &parser.BitwiseDisjunction{
																																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																	BitwiseConjunction: &parser.BitwiseConjunction{
																																		Shift: &parser.Shift{
																																			Addition: &parser.Addition{
																																				Multiplication: &parser.Multiplication{
																																					Unary: &parser.Unary{
																																						Accessor: &parser.Accessor{
																																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											Disjunction: &parser.Disjunction{
																												Conjunction: &parser.Conjunction{
																													Equality: &parser.Equality{
																														Comparison: &parser.Comparison{
																															BitwiseDisjunction: &parser.BitwiseDisjunction{
																																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																																	BitwiseConjunction: &parser.BitwiseConjunction{
																																		Shift: &parser.Shift{
																																			Addition: &parser.Addition{
																																				Multiplication: &parser.Multiplication{
																																					Unary: &parser.Unary{
																																						Accessor: &parser.Accessor{
																																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
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
				conditionalExpression: &parser.ConditionalExpression{
					ConditionalCases: []*parser.ConditionalCase{
						{
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							Commands: []*parser.Command{
								{
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
						{
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(13)},
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
							Commands: []*parser.Command{
								{
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
						{
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
							Commands: []*parser.Command{
								{
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(25)},
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
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(44)},
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "ConditionalExpression/error/command translating",
			args: args{
				conditionalExpression: &parser.ConditionalExpression{
					ConditionalCases: []*parser.ConditionalCase{
						{
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(12)},
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
							Commands: []*parser.Command{
								{
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(23)},
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
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(42)},
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
						{
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(13)},
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
							Commands: []*parser.Command{
								{
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(24)},
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
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(43)},
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
						{
							Condition: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													BitwiseDisjunction: &parser.BitwiseDisjunction{
														BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
															BitwiseConjunction: &parser.BitwiseConjunction{
																Shift: &parser.Shift{
																	Addition: &parser.Addition{
																		Multiplication: &parser.Multiplication{
																			Unary: &parser.Unary{
																				Accessor: &parser.Accessor{
																					Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(14)},
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
							Commands: []*parser.Command{
								{
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{IntegerNumber: pointer.ToInt64(25)},
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
									Expression: &parser.Expression{
										ListConstruction: &parser.ListConstruction{
											Disjunction: &parser.Disjunction{
												Conjunction: &parser.Conjunction{
													Equality: &parser.Equality{
														Comparison: &parser.Comparison{
															BitwiseDisjunction: &parser.BitwiseDisjunction{
																BitwiseExclusiveDisjunction: &parser.BitwiseExclusiveDisjunction{
																	BitwiseConjunction: &parser.BitwiseConjunction{
																		Shift: &parser.Shift{
																			Addition: &parser.Addition{
																				Multiplication: &parser.Multiplication{
																					Unary: &parser.Unary{
																						Accessor: &parser.Accessor{
																							Atom: &parser.Atom{Identifier: pointer.ToString("unknown")},
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
