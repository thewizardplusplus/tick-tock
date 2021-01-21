package parser

import (
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestParseToAST_withProgram(test *testing.T) {
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
			name: "Command/let",
			args: args{"let number = 23", new(Command)},
			wantAST: &Command{
				Let: &LetCommand{
					Identifier: "number",
					Expression: &Expression{
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
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Command/start/identifier/no arguments",
			args: args{"start Test()", new(Command)},
			wantAST: &Command{
				Start: &StartCommand{Name: pointer.ToString("Test"), Arguments: &ExpressionGroup{}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Command/start/identifier/few arguments",
			args: args{"start Test(12, 23, 42)", new(Command)},
			wantAST: &Command{
				Start: &StartCommand{
					Name: pointer.ToString("Test"),
					Arguments: &ExpressionGroup{
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
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Command/start/expression",
			args: args{"start [test()]()", new(Command)},
			wantAST: &Command{
				Start: &StartCommand{
					Expression: &Expression{
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
																			Accessor: &Accessor{
																				Atom: &Atom{
																					FunctionCall: &FunctionCall{Name: "test", Arguments: &ExpressionGroup{}},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					Arguments: &ExpressionGroup{},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/send/no arguments",
			args:    args{"send test()", new(Command)},
			wantAST: &Command{Send: &SendCommand{Name: "test", Arguments: &ExpressionGroup{}}},
			wantErr: assert.NoError,
		},
		{
			name: "Command/send/few arguments",
			args: args{"send test(12, 23, 42)", new(Command)},
			wantAST: &Command{
				Send: &SendCommand{
					Name: "test",
					Arguments: &ExpressionGroup{
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
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/set/no arguments",
			args:    args{"set test()", new(Command)},
			wantAST: &Command{Set: &SetCommand{Name: "test", Arguments: &ExpressionGroup{}}},
			wantErr: assert.NoError,
		},
		{
			name: "Command/set/few arguments",
			args: args{"set test(12, 23, 42)", new(Command)},
			wantAST: &Command{
				Set: &SetCommand{
					Name: "test",
					Arguments: &ExpressionGroup{
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
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/return",
			args:    args{"return", new(Command)},
			wantAST: &Command{Return: true},
			wantErr: assert.NoError,
		},
		{
			name: "Command/expression",
			args: args{"test()", new(Command)},
			wantAST: &Command{
				Expression: &Expression{
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
																		Accessor: &Accessor{
																			Atom: &Atom{
																				FunctionCall: &FunctionCall{Name: "test", Arguments: &ExpressionGroup{}},
																			},
																		},
																	},
																},
															},
														},
													},
												},
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
			name: "Message/nonempty/no parameters",
			args: args{"message test() send one() send two();", new(Message)},
			wantAST: &Message{
				Name: "test",
				Commands: []*Command{
					{Send: &SendCommand{Name: "one", Arguments: &ExpressionGroup{}}},
					{Send: &SendCommand{Name: "two", Arguments: &ExpressionGroup{}}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Message/nonempty/single parameter",
			args: args{"message test(x) send one() send two();", new(Message)},
			wantAST: &Message{
				Name:       "test",
				Parameters: []string{"x"},
				Commands: []*Command{
					{Send: &SendCommand{Name: "one", Arguments: &ExpressionGroup{}}},
					{Send: &SendCommand{Name: "two", Arguments: &ExpressionGroup{}}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Message/nonempty/single parameter/trailing comma",
			args: args{"message test(x,) send one() send two();", new(Message)},
			wantAST: &Message{
				Name:       "test",
				Parameters: []string{"x"},
				Commands: []*Command{
					{Send: &SendCommand{Name: "one", Arguments: &ExpressionGroup{}}},
					{Send: &SendCommand{Name: "two", Arguments: &ExpressionGroup{}}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Message/nonempty/few parameters",
			args: args{"message test(x, y, z) send one() send two();", new(Message)},
			wantAST: &Message{
				Name:       "test",
				Parameters: []string{"x", "y", "z"},
				Commands: []*Command{
					{Send: &SendCommand{Name: "one", Arguments: &ExpressionGroup{}}},
					{Send: &SendCommand{Name: "two", Arguments: &ExpressionGroup{}}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Message/nonempty/few parameters/trailing comma",
			args: args{"message test(x, y, z,) send one() send two();", new(Message)},
			wantAST: &Message{
				Name:       "test",
				Parameters: []string{"x", "y", "z"},
				Commands: []*Command{
					{Send: &SendCommand{Name: "one", Arguments: &ExpressionGroup{}}},
					{Send: &SendCommand{Name: "two", Arguments: &ExpressionGroup{}}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Message/empty",
			args:    args{"message test();", new(Message)},
			wantAST: &Message{"test", nil, nil},
			wantErr: assert.NoError,
		},
		{
			name:    "State/nonempty/no parameters",
			args:    args{"state test() message one(); message two();;", new(State)},
			wantAST: &State{"test", nil, []*Message{{"one", nil, nil}, {"two", nil, nil}}},
			wantErr: assert.NoError,
		},
		{
			name:    "State/nonempty/single parameter",
			args:    args{"state test(x) message one(); message two();;", new(State)},
			wantAST: &State{"test", []string{"x"}, []*Message{{"one", nil, nil}, {"two", nil, nil}}},
			wantErr: assert.NoError,
		},
		{
			name:    "State/nonempty/single parameter/trailing comma",
			args:    args{"state test(x,) message one(); message two();;", new(State)},
			wantAST: &State{"test", []string{"x"}, []*Message{{"one", nil, nil}, {"two", nil, nil}}},
			wantErr: assert.NoError,
		},
		{
			name: "State/nonempty/few parameters",
			args: args{"state test(x, y, z) message one(); message two();;", new(State)},
			wantAST: &State{
				Name:       "test",
				Parameters: []string{"x", "y", "z"},
				Messages:   []*Message{{"one", nil, nil}, {"two", nil, nil}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "State/nonempty/few parameters/trailing comma",
			args: args{"state test(x, y, z,) message one(); message two();;", new(State)},
			wantAST: &State{
				Name:       "test",
				Parameters: []string{"x", "y", "z"},
				Messages:   []*Message{{"one", nil, nil}, {"two", nil, nil}},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "State/empty",
			args:    args{"state test();", new(State)},
			wantAST: &State{"test", nil, nil},
			wantErr: assert.NoError,
		},
		{
			name:    "Actor/nonempty/no parameters",
			args:    args{"actor Main() state one(); state two();;", new(Actor)},
			wantAST: &Actor{"Main", &IdentifierGroup{}, []*State{{"one", nil, nil}, {"two", nil, nil}}},
			wantErr: assert.NoError,
		},
		{
			name: "Actor/nonempty/single parameter",
			args: args{"actor Main(x) state one(); state two();;", new(Actor)},
			wantAST: &Actor{
				Name:       "Main",
				Parameters: &IdentifierGroup{Identifiers: []string{"x"}},
				States:     []*State{{"one", nil, nil}, {"two", nil, nil}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Actor/nonempty/single parameter/trailing comma",
			args: args{"actor Main(x,) state one(); state two();;", new(Actor)},
			wantAST: &Actor{
				Name:       "Main",
				Parameters: &IdentifierGroup{Identifiers: []string{"x"}},
				States:     []*State{{"one", nil, nil}, {"two", nil, nil}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Actor/nonempty/few parameters",
			args: args{"actor Main(x, y, z) state one(); state two();;", new(Actor)},
			wantAST: &Actor{
				Name:       "Main",
				Parameters: &IdentifierGroup{Identifiers: []string{"x", "y", "z"}},
				States:     []*State{{"one", nil, nil}, {"two", nil, nil}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Actor/nonempty/few parameters/trailing comma",
			args: args{"actor Main(x, y, z,) state one(); state two();;", new(Actor)},
			wantAST: &Actor{
				Name:       "Main",
				Parameters: &IdentifierGroup{Identifiers: []string{"x", "y", "z"}},
				States:     []*State{{"one", nil, nil}, {"two", nil, nil}},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Actor/empty",
			args:    args{"actor Main();", new(Actor)},
			wantAST: &Actor{"Main", &IdentifierGroup{}, nil},
			wantErr: assert.NoError,
		},
		{
			name:    "ActorClass/nonempty/no parameters",
			args:    args{"class Main() state one(); state two();;", new(ActorClass)},
			wantAST: &ActorClass{"Main", nil, []*State{{"one", nil, nil}, {"two", nil, nil}}},
			wantErr: assert.NoError,
		},
		{
			name:    "ActorClass/nonempty/single parameter",
			args:    args{"class Main(x) state one(); state two();;", new(ActorClass)},
			wantAST: &ActorClass{"Main", []string{"x"}, []*State{{"one", nil, nil}, {"two", nil, nil}}},
			wantErr: assert.NoError,
		},
		{
			name:    "ActorClass/nonempty/single parameter/trailing comma",
			args:    args{"class Main(x,) state one(); state two();;", new(ActorClass)},
			wantAST: &ActorClass{"Main", []string{"x"}, []*State{{"one", nil, nil}, {"two", nil, nil}}},
			wantErr: assert.NoError,
		},
		{
			name: "ActorClass/nonempty/few parameters",
			args: args{"class Main(x, y, z) state one(); state two();;", new(ActorClass)},
			wantAST: &ActorClass{
				Name:       "Main",
				Parameters: []string{"x", "y", "z"},
				States:     []*State{{"one", nil, nil}, {"two", nil, nil}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "ActorClass/nonempty/few parameters/trailing comma",
			args: args{"class Main(x, y, z,) state one(); state two();;", new(ActorClass)},
			wantAST: &ActorClass{
				Name:       "Main",
				Parameters: []string{"x", "y", "z"},
				States:     []*State{{"one", nil, nil}, {"two", nil, nil}},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "ActorClass/empty",
			args:    args{"class Main();", new(ActorClass)},
			wantAST: &ActorClass{"Main", nil, nil},
			wantErr: assert.NoError,
		},
		{
			name: "Definition/actor",
			args: args{"actor Main() state one(); state two();;", new(Definition)},
			wantAST: &Definition{
				Actor: &Actor{"Main", &IdentifierGroup{}, []*State{{"one", nil, nil}, {"two", nil, nil}}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Definition/actor class",
			args: args{"class Main() state one(); state two();;", new(Definition)},
			wantAST: &Definition{
				ActorClass: &ActorClass{"Main", nil, []*State{{"one", nil, nil}, {"two", nil, nil}}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Program/nonempty",
			args: args{"actor One(); actor Two();", new(Program)},
			wantAST: &Program{
				Definitions: []*Definition{
					{Actor: &Actor{"One", &IdentifierGroup{}, nil}},
					{Actor: &Actor{"Two", &IdentifierGroup{}, nil}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Program/empty",
			args:    args{"", new(Program)},
			wantAST: new(Program),
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
