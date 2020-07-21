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
							Disjunction: &Disjunction{
								Conjunction: &Conjunction{
									Equality: &Equality{
										Comparison: &Comparison{
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
			name:    "Command/start/identifier/no arguments",
			args:    args{"start Test()", new(Command)},
			wantAST: &Command{Start: &StartCommand{Name: pointer.ToString("Test")}},
			wantErr: assert.NoError,
		},
		{
			name: "Command/start/identifier/single argument",
			args: args{"start Test(12)", new(Command)},
			wantAST: &Command{
				Start: &StartCommand{
					Name: pointer.ToString("Test"),
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
			wantErr: assert.NoError,
		},
		{
			name: "Command/start/identifier/single argument/trailing comma",
			args: args{"start Test(12,)", new(Command)},
			wantAST: &Command{
				Start: &StartCommand{
					Name: pointer.ToString("Test"),
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
			wantErr: assert.NoError,
		},
		{
			name: "Command/start/identifier/few arguments",
			args: args{"start Test(12, 23, 42)", new(Command)},
			wantAST: &Command{
				Start: &StartCommand{
					Name: pointer.ToString("Test"),
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
			name: "Command/start/identifier/few arguments/trailing comma",
			args: args{"start Test(12, 23, 42,)", new(Command)},
			wantAST: &Command{
				Start: &StartCommand{
					Name: pointer.ToString("Test"),
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
			name: "Command/start/expression",
			args: args{"start [test()]()", new(Command)},
			wantAST: &Command{
				Start: &StartCommand{
					Expression: &Expression{
						ListConstruction: &ListConstruction{
							Disjunction: &Disjunction{
								Conjunction: &Conjunction{
									Equality: &Equality{
										Comparison: &Comparison{
											Addition: &Addition{
												Multiplication: &Multiplication{
													Unary: &Unary{
														Accessor: &Accessor{Atom: &Atom{FunctionCall: &FunctionCall{Name: "test"}}},
													},
												},
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
			name:    "Command/send/no arguments",
			args:    args{"send test()", new(Command)},
			wantAST: &Command{Send: &SendCommand{Name: "test"}},
			wantErr: assert.NoError,
		},
		{
			name: "Command/send/single argument",
			args: args{"send test(12)", new(Command)},
			wantAST: &Command{
				Send: &SendCommand{
					Name: "test",
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
			wantErr: assert.NoError,
		},
		{
			name: "Command/send/single argument/trailing comma",
			args: args{"send test(12,)", new(Command)},
			wantAST: &Command{
				Send: &SendCommand{
					Name: "test",
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
			wantErr: assert.NoError,
		},
		{
			name: "Command/send/few arguments",
			args: args{"send test(12, 23, 42)", new(Command)},
			wantAST: &Command{
				Send: &SendCommand{
					Name: "test",
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
			name: "Command/send/few arguments/trailing comma",
			args: args{"send test(12, 23, 42,)", new(Command)},
			wantAST: &Command{
				Send: &SendCommand{
					Name: "test",
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
			name:    "Command/set/no arguments",
			args:    args{"set test()", new(Command)},
			wantAST: &Command{Set: &SetCommand{Name: "test"}},
			wantErr: assert.NoError,
		},
		{
			name: "Command/set/single argument",
			args: args{"set test(12)", new(Command)},
			wantAST: &Command{
				Set: &SetCommand{
					Name: "test",
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
			wantErr: assert.NoError,
		},
		{
			name: "Command/set/single argument/trailing comma",
			args: args{"set test(12,)", new(Command)},
			wantAST: &Command{
				Set: &SetCommand{
					Name: "test",
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
			wantErr: assert.NoError,
		},
		{
			name: "Command/set/few arguments",
			args: args{"set test(12, 23, 42)", new(Command)},
			wantAST: &Command{
				Set: &SetCommand{
					Name: "test",
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
			name: "Command/set/few arguments/trailing comma",
			args: args{"set test(12, 23, 42,)", new(Command)},
			wantAST: &Command{
				Set: &SetCommand{
					Name: "test",
					Arguments: []*Expression{
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
						{
							ListConstruction: &ListConstruction{
								Disjunction: &Disjunction{
									Conjunction: &Conjunction{
										Equality: &Equality{
											Comparison: &Comparison{
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
						Disjunction: &Disjunction{
							Conjunction: &Conjunction{
								Equality: &Equality{
									Comparison: &Comparison{
										Addition: &Addition{
											Multiplication: &Multiplication{
												Unary: &Unary{
													Accessor: &Accessor{Atom: &Atom{FunctionCall: &FunctionCall{Name: "test"}}},
												},
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
				Name:     "test",
				Commands: []*Command{{Send: &SendCommand{Name: "one"}}, {Send: &SendCommand{Name: "two"}}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Message/nonempty/single parameter",
			args: args{"message test(x) send one() send two();", new(Message)},
			wantAST: &Message{
				Name:       "test",
				Parameters: []string{"x"},
				Commands:   []*Command{{Send: &SendCommand{Name: "one"}}, {Send: &SendCommand{Name: "two"}}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Message/nonempty/single parameter/trailing comma",
			args: args{"message test(x,) send one() send two();", new(Message)},
			wantAST: &Message{
				Name:       "test",
				Parameters: []string{"x"},
				Commands:   []*Command{{Send: &SendCommand{Name: "one"}}, {Send: &SendCommand{Name: "two"}}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Message/nonempty/few parameters",
			args: args{"message test(x, y, z) send one() send two();", new(Message)},
			wantAST: &Message{
				Name:       "test",
				Parameters: []string{"x", "y", "z"},
				Commands:   []*Command{{Send: &SendCommand{Name: "one"}}, {Send: &SendCommand{Name: "two"}}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Message/nonempty/few parameters/trailing comma",
			args: args{"message test(x, y, z,) send one() send two();", new(Message)},
			wantAST: &Message{
				Name:       "test",
				Parameters: []string{"x", "y", "z"},
				Commands:   []*Command{{Send: &SendCommand{Name: "one"}}, {Send: &SendCommand{Name: "two"}}},
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
			name:    "Actor/nonempty",
			args:    args{"actor Main state one(); state two();;", new(Actor)},
			wantAST: &Actor{"Main", []*State{{"one", nil, nil}, {"two", nil, nil}}},
			wantErr: assert.NoError,
		},
		{
			name:    "Actor/empty",
			args:    args{"actor Main;", new(Actor)},
			wantAST: &Actor{"Main", nil},
			wantErr: assert.NoError,
		},
		{
			name:    "ActorClass/nonempty",
			args:    args{"class Main state one(); state two();;", new(ActorClass)},
			wantAST: &ActorClass{"Main", []*State{{"one", nil, nil}, {"two", nil, nil}}},
			wantErr: assert.NoError,
		},
		{
			name:    "ActorClass/empty",
			args:    args{"class Main;", new(ActorClass)},
			wantAST: &ActorClass{"Main", nil},
			wantErr: assert.NoError,
		},
		{
			name:    "Definition/actor",
			args:    args{"actor Main state one(); state two();;", new(Definition)},
			wantAST: &Definition{Actor: &Actor{"Main", []*State{{"one", nil, nil}, {"two", nil, nil}}}},
			wantErr: assert.NoError,
		},
		{
			name: "Definition/actor class",
			args: args{"class Main state one(); state two();;", new(Definition)},
			wantAST: &Definition{
				ActorClass: &ActorClass{"Main", []*State{{"one", nil, nil}, {"two", nil, nil}}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Program/nonempty",
			args: args{"actor One; actor Two;", new(Program)},
			wantAST: &Program{
				Definitions: []*Definition{{Actor: &Actor{"One", nil}}, {Actor: &Actor{"Two", nil}}},
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
			err := parseToAST(testData.args.code, testData.args.ast)
			assert.Equal(test, testData.wantAST, testData.args.ast)
			testData.wantErr(test, err)
		})
	}
}
