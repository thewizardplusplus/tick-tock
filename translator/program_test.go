package translator

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/AlekSi/pointer"
	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	testutils "github.com/thewizardplusplus/tick-tock/internal/test-utils"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/commands"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	runtimemocks "github.com/thewizardplusplus/tick-tock/runtime/mocks"
	waitermocks "github.com/thewizardplusplus/tick-tock/runtime/waiter/mocks"
)

func TestTranslate(test *testing.T) {
	type args struct {
		makeActors          func(options Options) []*parser.Actor
		declaredIdentifiers mapset.Set
	}

	for _, testData := range []struct {
		name           string
		args           args
		makeWantActors func(
			options Options,
			dependencies runtime.Dependencies,
		) runtime.ConcurrentActorGroup
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success with actors",
			args: args{
				makeActors: func(options Options) []*parser.Actor {
					return []*parser.Actor{
						{States: []*parser.State{{Name: options.InitialState}, {Name: "one"}}},
						{States: []*parser.State{{Name: options.InitialState}, {Name: "two"}}},
					}
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			makeWantActors: func(
				options Options,
				dependencies runtime.Dependencies,
			) runtime.ConcurrentActorGroup {
				actorOne, _ := runtime.NewActor(
					runtime.StateGroup{
						options.InitialState: runtime.MessageGroup{},
						"one":                runtime.MessageGroup{},
					},
					options.InitialState,
				)
				actorTwo, _ := runtime.NewActor(
					runtime.StateGroup{
						options.InitialState: runtime.MessageGroup{},
						"two":                runtime.MessageGroup{},
					},
					options.InitialState,
				)
				return runtime.ConcurrentActorGroup{
					runtime.NewConcurrentActor(actorOne, options.InboxSize, dependencies),
					runtime.NewConcurrentActor(actorTwo, options.InboxSize, dependencies),
				}
			},
			wantErr: assert.NoError,
		},
		{
			name: "success without actors",
			args: args{
				makeActors:          func(options Options) []*parser.Actor { return nil },
				declaredIdentifiers: mapset.NewSet("test"),
			},
			makeWantActors: func(
				options Options,
				dependencies runtime.Dependencies,
			) runtime.ConcurrentActorGroup {
				return nil
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with the expression",
			args: args{
				makeActors: func(options Options) []*parser.Actor {
					return []*parser.Actor{
						{
							States: []*parser.State{
								{
									Name: options.InitialState,
									Messages: []*parser.Message{
										{
											Name: "message_0",
											Commands: []*parser.Command{
												{
													Expression: &parser.Expression{
														ListConstruction: &parser.ListConstruction{
															Disjunction: &parser.Disjunction{
																Conjunction: &parser.Conjunction{
																	Equality: &parser.Equality{
																		Comparison: &parser.Comparison{
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
								},
							},
						},
					}
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			makeWantActors: func(
				options Options,
				dependencies runtime.Dependencies,
			) runtime.ConcurrentActorGroup {
				actorOne, _ := runtime.NewActor(
					runtime.StateGroup{
						options.InitialState: runtime.MessageGroup{
							"message_0": runtime.CommandGroup{
								commands.NewExpressionCommand(expressions.NewIdentifier("test")),
							},
						},
					},
					options.InitialState,
				)
				return runtime.ConcurrentActorGroup{
					runtime.NewConcurrentActor(actorOne, options.InboxSize, dependencies),
				}
			},
			wantErr: assert.NoError,
		},
		{
			name: "error with states translation",
			args: args{
				makeActors: func(options Options) []*parser.Actor {
					return []*parser.Actor{{States: []*parser.State{{Name: options.InitialState}}}, {}}
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			makeWantActors: func(
				options Options,
				dependencies runtime.Dependencies,
			) runtime.ConcurrentActorGroup {
				return nil
			},
			wantErr: assert.Error,
		},
		{
			name: "error with actor construction",
			args: args{
				makeActors: func(options Options) []*parser.Actor {
					return []*parser.Actor{
						{States: []*parser.State{{Name: "one"}}},
						{States: []*parser.State{{Name: "two"}}},
					}
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			makeWantActors: func(
				options Options,
				dependencies runtime.Dependencies,
			) runtime.ConcurrentActorGroup {
				return nil
			},
			wantErr: assert.Error,
		},
		{
			name: "error with the expression",
			args: args{
				makeActors: func(options Options) []*parser.Actor {
					return []*parser.Actor{
						{
							States: []*parser.State{
								{
									Name: options.InitialState,
									Messages: []*parser.Message{
										{
											Name: "message_0",
											Commands: []*parser.Command{
												{
													Expression: &parser.Expression{
														ListConstruction: &parser.ListConstruction{
															Disjunction: &parser.Disjunction{
																Conjunction: &parser.Conjunction{
																	Equality: &parser.Equality{
																		Comparison: &parser.Comparison{
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
					}
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			makeWantActors: func(
				options Options,
				dependencies runtime.Dependencies,
			) runtime.ConcurrentActorGroup {
				return nil
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := testData.args.declaredIdentifiers.Clone()

			options := Options{testutils.BufferedInbox, "__initialization__"}
			waiter := new(waitermocks.Waiter)
			errorHandler := new(runtimemocks.ErrorHandler)
			dependencies := runtime.Dependencies{Waiter: waiter, ErrorHandler: errorHandler}
			gotActors, err := Translate(
				testData.args.makeActors(options),
				testData.args.declaredIdentifiers,
				options,
				dependencies,
			)
			cleanInboxes(gotActors)

			wantActors := testData.makeWantActors(options, dependencies)
			cleanInboxes(wantActors)

			mock.AssertExpectationsForObjects(test, waiter, errorHandler)
			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, wantActors, gotActors)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateStates(test *testing.T) {
	type args struct {
		states              []*parser.State
		declaredIdentifiers mapset.Set
	}

	for _, testData := range []struct {
		name       string
		args       args
		wantStates runtime.StateGroup
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success with nonempty states",
			args: args{
				states: []*parser.State{
					{Name: "state_0", Messages: []*parser.Message{{Name: "message_0"}, {Name: "message_1"}}},
					{Name: "state_1", Messages: []*parser.Message{{Name: "message_2"}, {Name: "message_3"}}},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantStates: runtime.StateGroup{
				"state_0": runtime.MessageGroup{"message_0": nil, "message_1": nil},
				"state_1": runtime.MessageGroup{"message_2": nil, "message_3": nil},
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with empty states",
			args: args{
				states:              []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantStates: runtime.StateGroup{
				"state_0": runtime.MessageGroup{},
				"state_1": runtime.MessageGroup{},
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with the expression",
			args: args{
				states: []*parser.State{
					{
						Name: "state_0",
						Messages: []*parser.Message{
							{
								Name: "message_0",
								Commands: []*parser.Command{
									{
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
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
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantStates: runtime.StateGroup{
				"state_0": runtime.MessageGroup{
					"message_0": runtime.CommandGroup{
						commands.NewExpressionCommand(expressions.NewIdentifier("test")),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "error without states",
			args: args{
				states:              nil,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantStates: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with duplicate states",
			args: args{
				states:              []*parser.State{{Name: "test"}, {Name: "test"}},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantStates: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with messages translation",
			args: args{
				states: []*parser.State{
					{Name: "state_0", Messages: []*parser.Message{{Name: "message_0"}, {Name: "message_1"}}},
					{Name: "state_1", Messages: []*parser.Message{{Name: "test"}, {Name: "test"}}},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantStates: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with an unknown state",
			args: args{
				states: []*parser.State{
					{
						Name: "state_0",
						Messages: []*parser.Message{
							{
								Name: "message_0",
								Commands: []*parser.Command{
									{Send: pointer.ToString("command_0")},
									{Set: pointer.ToString("state_unknown")},
								},
							},
							{
								Name: "message_1",
								Commands: []*parser.Command{
									{Send: pointer.ToString("command_2")},
									{Set: pointer.ToString("state_unknown")},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantStates: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with the expression",
			args: args{
				states: []*parser.State{
					{
						Name: "state_0",
						Messages: []*parser.Message{
							{
								Name: "message_0",
								Commands: []*parser.Command{
									{
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												Disjunction: &parser.Disjunction{
													Conjunction: &parser.Conjunction{
														Equality: &parser.Equality{
															Comparison: &parser.Comparison{
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
			wantStates: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := testData.args.declaredIdentifiers.Clone()

			gotStates, err := translateStates(testData.args.states, testData.args.declaredIdentifiers)

			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.wantStates, gotStates)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateMessages(test *testing.T) {
	type args struct {
		messages            []*parser.Message
		declaredIdentifiers mapset.Set
	}

	for _, testData := range []struct {
		name                       string
		args                       args
		wantMessages               runtime.MessageGroup
		wantSettedStatesByMessages settedStateGroup
		wantErr                    assert.ErrorAssertionFunc
	}{
		{
			name: "success with nonempty messages (without set commands)",
			args: args{
				messages: []*parser.Message{
					{
						Name: "message_0",
						Commands: []*parser.Command{
							{Send: pointer.ToString("command_0")},
							{Send: pointer.ToString("command_1")},
						},
					},
					{
						Name: "message_1",
						Commands: []*parser.Command{
							{Send: pointer.ToString("command_2")},
							{Send: pointer.ToString("command_3")},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages: runtime.MessageGroup{
				"message_0": runtime.CommandGroup{
					commands.NewSendCommand("command_0"),
					commands.NewSendCommand("command_1"),
				},
				"message_1": runtime.CommandGroup{
					commands.NewSendCommand("command_2"),
					commands.NewSendCommand("command_3"),
				},
			},
			wantSettedStatesByMessages: settedStateGroup{
				"message_0": mapset.NewSet(),
				"message_1": mapset.NewSet(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with nonempty messages (with different set commands)",
			args: args{
				messages: []*parser.Message{
					{
						Name: "message_0",
						Commands: []*parser.Command{
							{Send: pointer.ToString("command_0")},
							{Set: pointer.ToString("command_1")},
						},
					},
					{
						Name: "message_1",
						Commands: []*parser.Command{
							{Send: pointer.ToString("command_2")},
							{Set: pointer.ToString("command_3")},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages: runtime.MessageGroup{
				"message_0": runtime.CommandGroup{
					commands.NewSendCommand("command_0"),
					commands.NewSetCommand("command_1"),
				},
				"message_1": runtime.CommandGroup{
					commands.NewSendCommand("command_2"),
					commands.NewSetCommand("command_3"),
				},
			},
			wantSettedStatesByMessages: settedStateGroup{
				"message_0": mapset.NewSet("command_1"),
				"message_1": mapset.NewSet("command_3"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with nonempty messages (with same set commands)",
			args: args{
				messages: []*parser.Message{
					{
						Name: "message_0",
						Commands: []*parser.Command{
							{Send: pointer.ToString("command_1")},
							{Set: pointer.ToString("command_0")},
						},
					},
					{
						Name: "message_1",
						Commands: []*parser.Command{
							{Send: pointer.ToString("command_2")},
							{Set: pointer.ToString("command_0")},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages: runtime.MessageGroup{
				"message_0": runtime.CommandGroup{
					commands.NewSendCommand("command_1"),
					commands.NewSetCommand("command_0"),
				},
				"message_1": runtime.CommandGroup{
					commands.NewSendCommand("command_2"),
					commands.NewSetCommand("command_0"),
				},
			},
			wantSettedStatesByMessages: settedStateGroup{
				"message_0": mapset.NewSet("command_0"),
				"message_1": mapset.NewSet("command_0"),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with empty messages",
			args: args{
				messages:            []*parser.Message{{Name: "message_0"}, {Name: "message_1"}},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages: runtime.MessageGroup{"message_0": nil, "message_1": nil},
			wantSettedStatesByMessages: settedStateGroup{
				"message_0": mapset.NewSet(),
				"message_1": mapset.NewSet(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success without messages",
			args: args{
				messages:            nil,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages:               runtime.MessageGroup{},
			wantSettedStatesByMessages: make(settedStateGroup),
			wantErr:                    assert.NoError,
		},
		{
			name: "success with the expression",
			args: args{
				messages: []*parser.Message{
					{
						Name: "message_0",
						Commands: []*parser.Command{
							{
								Expression: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
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
			wantMessages: runtime.MessageGroup{
				"message_0": runtime.CommandGroup{
					commands.NewExpressionCommand(expressions.NewIdentifier("test")),
				},
			},
			wantSettedStatesByMessages: settedStateGroup{"message_0": mapset.NewSet()},
			wantErr:                    assert.NoError,
		},
		{
			name: "error with duplicate messages",
			args: args{
				messages:            []*parser.Message{{Name: "test"}, {Name: "test"}},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages:               nil,
			wantSettedStatesByMessages: nil,
			wantErr:                    assert.Error,
		},
		{
			name: "error with commands translation",
			args: args{
				messages: []*parser.Message{
					{
						Name: "message_0",
						Commands: []*parser.Command{
							{Send: pointer.ToString("command_0")},
							{Send: pointer.ToString("command_1")},
						},
					},
					{
						Name: "message_1",
						Commands: []*parser.Command{
							{Send: pointer.ToString("command_2")},
							{Set: pointer.ToString("command_3")},
							{Send: pointer.ToString("command_4")},
							{Set: pointer.ToString("command_5")},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages:               nil,
			wantSettedStatesByMessages: nil,
			wantErr:                    assert.Error,
		},
		{
			name: "error with the expression",
			args: args{
				messages: []*parser.Message{
					{
						Name: "message_0",
						Commands: []*parser.Command{
							{
								Expression: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										Disjunction: &parser.Disjunction{
											Conjunction: &parser.Conjunction{
												Equality: &parser.Equality{
													Comparison: &parser.Comparison{
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
			wantMessages:               nil,
			wantSettedStatesByMessages: nil,
			wantErr:                    assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := testData.args.declaredIdentifiers.Clone()

			gotMessages, gotSettedStatesByMessages, err :=
				translateMessages(testData.args.messages, testData.args.declaredIdentifiers)

			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.wantMessages, gotMessages)
			assert.Equal(test, testData.wantSettedStatesByMessages, gotSettedStatesByMessages)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateCommands(test *testing.T) {
	type args struct {
		commands            []*parser.Command
		declaredIdentifiers mapset.Set
	}

	for _, testData := range []struct {
		name             string
		args             args
		wantCommands     runtime.CommandGroup
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "success with commands (without a set command)",
			args: args{
				commands: []*parser.Command{
					{Send: pointer.ToString("one")},
					{Send: pointer.ToString("two")},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: runtime.CommandGroup{
				commands.NewSendCommand("one"),
				commands.NewSendCommand("two"),
			},
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "success with commands (with a set command)",
			args: args{
				commands: []*parser.Command{
					{Send: pointer.ToString("one")},
					{Set: pointer.ToString("two")},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: runtime.CommandGroup{
				commands.NewSendCommand("one"),
				commands.NewSetCommand("two"),
			},
			wantSettedStates: mapset.NewSet("two"),
			wantErr:          assert.NoError,
		},
		{
			name: "success with the return command",
			args: args{
				commands: []*parser.Command{
					{Send: pointer.ToString("one")},
					{Send: pointer.ToString("two")},
					{Return: true},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: runtime.CommandGroup{
				commands.NewSendCommand("one"),
				commands.NewSendCommand("two"),
				commands.ReturnCommand{},
			},
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "success with commands (with an expression command and an existing identifier)",
			args: args{
				commands: []*parser.Command{
					{
						Expression: &parser.Expression{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: runtime.CommandGroup{
				commands.NewExpressionCommand(expressions.NewIdentifier("test")),
			},
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "success with commands (with an expression command and an nonexistent identifier)",
			args: args{
				commands: []*parser.Command{
					{
						Let: &parser.LetCommand{
							Identifier: "test2",
							Expression: &parser.Expression{
								ListConstruction: &parser.ListConstruction{
									Disjunction: &parser.Disjunction{
										Conjunction: &parser.Conjunction{
											Equality: &parser.Equality{
												Comparison: &parser.Comparison{
													Addition: &parser.Addition{
														Multiplication: &parser.Multiplication{
															Unary: &parser.Unary{
																Accessor: &parser.Accessor{Atom: &parser.Atom{Number: pointer.ToFloat64(23)}},
															},
														},
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
												Addition: &parser.Addition{
													Multiplication: &parser.Multiplication{
														Unary: &parser.Unary{
															Accessor: &parser.Accessor{
																Atom: &parser.Atom{Identifier: pointer.ToString("test2")},
															},
														},
													},
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
			wantCommands: runtime.CommandGroup{
				commands.NewLetCommand("test2", expressions.NewNumber(23)),
				commands.NewExpressionCommand(expressions.NewIdentifier("test2")),
			},
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "success without commands",
			args: args{
				commands:            nil,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands:     nil,
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "success with commands (with setted states)",
			args: args{
				commands: []*parser.Command{
					{
						Expression: &parser.Expression{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
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
																										Addition: &parser.Addition{
																											Multiplication: &parser.Multiplication{
																												Unary: &parser.Unary{
																													Accessor: &parser.Accessor{
																														Atom: &parser.Atom{Number: pointer.ToFloat64(23)},
																													},
																												},
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																				Commands: []*parser.Command{{Set: pointer.ToString("one")}},
																			},
																			{
																				Condition: &parser.Expression{
																					ListConstruction: &parser.ListConstruction{
																						Disjunction: &parser.Disjunction{
																							Conjunction: &parser.Conjunction{
																								Equality: &parser.Equality{
																									Comparison: &parser.Comparison{
																										Addition: &parser.Addition{
																											Multiplication: &parser.Multiplication{
																												Unary: &parser.Unary{
																													Accessor: &parser.Accessor{
																														Atom: &parser.Atom{Number: pointer.ToFloat64(42)},
																													},
																												},
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																				Commands: []*parser.Command{{Set: pointer.ToString("two")}},
																			},
																		},
																	},
																},
															},
														},
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
																										Addition: &parser.Addition{
																											Multiplication: &parser.Multiplication{
																												Unary: &parser.Unary{
																													Accessor: &parser.Accessor{
																														Atom: &parser.Atom{Number: pointer.ToFloat64(24)},
																													},
																												},
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																				Commands: []*parser.Command{{Set: pointer.ToString("two")}},
																			},
																			{
																				Condition: &parser.Expression{
																					ListConstruction: &parser.ListConstruction{
																						Disjunction: &parser.Disjunction{
																							Conjunction: &parser.Conjunction{
																								Equality: &parser.Equality{
																									Comparison: &parser.Comparison{
																										Addition: &parser.Addition{
																											Multiplication: &parser.Multiplication{
																												Unary: &parser.Unary{
																													Accessor: &parser.Accessor{
																														Atom: &parser.Atom{Number: pointer.ToFloat64(43)},
																													},
																												},
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																				Commands: []*parser.Command{{Set: pointer.ToString("three")}},
																			},
																		},
																	},
																},
															},
														},
													},
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
			wantCommands: runtime.CommandGroup{
				commands.NewExpressionCommand(
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(23),
							Command:   runtime.CommandGroup{commands.NewSetCommand("one")},
						},
						{
							Condition: expressions.NewNumber(42),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two")},
						},
					}),
				),
				commands.NewExpressionCommand(
					expressions.NewConditionalExpression([]expressions.ConditionalCase{
						{
							Condition: expressions.NewNumber(24),
							Command:   runtime.CommandGroup{commands.NewSetCommand("two")},
						},
						{
							Condition: expressions.NewNumber(43),
							Command:   runtime.CommandGroup{commands.NewSetCommand("three")},
						},
					}),
				),
			},
			wantSettedStates: mapset.NewSet("one", "two", "three"),
			wantErr:          assert.NoError,
		},
		{
			name: "error with expression command translation",
			args: args{
				commands: []*parser.Command{
					{
						Expression: &parser.Expression{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
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
			wantCommands: nil,
			wantErr:      assert.Error,
		},
		{
			name: "error with the return command",
			args: args{
				commands: []*parser.Command{
					{Send: pointer.ToString("one")},
					{Return: true},
					{Send: pointer.ToString("two")},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: nil,
			wantErr:      assert.Error,
		},
		{
			name: "error with a second set command",
			args: args{
				commands: []*parser.Command{
					{Send: pointer.ToString("one")},
					{Set: pointer.ToString("two")},
					{Send: pointer.ToString("three")},
					{Set: pointer.ToString("four")},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: nil,
			wantErr:      assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := testData.args.declaredIdentifiers.Clone()

			gotCommands, gotSettedStates, err :=
				translateCommands(testData.args.commands, testData.args.declaredIdentifiers)

			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.wantCommands, gotCommands)
			assert.Equal(test, testData.wantSettedStates, gotSettedStates)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateCommand(test *testing.T) {
	type args struct {
		command             *parser.Command
		declaredIdentifiers mapset.Set
	}

	for _, testData := range []struct {
		name                    string
		args                    args
		wantDeclaredIdentifiers mapset.Set
		wantCommand             runtime.Command
		wantTopLevelSettedState string
		wantSettedStates        mapset.Set
		wantReturn              assert.BoolAssertionFunc
		wantErr                 assert.ErrorAssertionFunc
	}{
		{
			name: "Command/let/success/nonexistent identifier",
			args: args{
				command: &parser.Command{
					Let: &parser.LetCommand{
						Identifier: "test2",
						Expression: &parser.Expression{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												Addition: &parser.Addition{
													Multiplication: &parser.Multiplication{
														Unary: &parser.Unary{
															Accessor: &parser.Accessor{Atom: &parser.Atom{Number: pointer.ToFloat64(23)}},
														},
													},
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
			wantDeclaredIdentifiers: mapset.NewSet("test", "test2"),
			wantCommand:             commands.NewLetCommand("test2", expressions.NewNumber(23)),
			wantTopLevelSettedState: "",
			wantSettedStates:        mapset.NewSet(),
			wantReturn:              assert.False,
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/let/success/existing identifier",
			args: args{
				command: &parser.Command{
					Let: &parser.LetCommand{
						Identifier: "test",
						Expression: &parser.Expression{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
												Addition: &parser.Addition{
													Multiplication: &parser.Multiplication{
														Unary: &parser.Unary{
															Accessor: &parser.Accessor{Atom: &parser.Atom{Number: pointer.ToFloat64(23)}},
														},
													},
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
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             commands.NewLetCommand("test", expressions.NewNumber(23)),
			wantTopLevelSettedState: "",
			wantSettedStates:        mapset.NewSet(),
			wantReturn:              assert.False,
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/let/success/with setted states",
			args: args{
				command: &parser.Command{
					Let: &parser.LetCommand{
						Identifier: "test2",
						Expression: &parser.Expression{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
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
																										Addition: &parser.Addition{
																											Multiplication: &parser.Multiplication{
																												Unary: &parser.Unary{
																													Accessor: &parser.Accessor{
																														Atom: &parser.Atom{Number: pointer.ToFloat64(23)},
																													},
																												},
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																				Commands: []*parser.Command{{Set: pointer.ToString("one")}},
																			},
																			{
																				Condition: &parser.Expression{
																					ListConstruction: &parser.ListConstruction{
																						Disjunction: &parser.Disjunction{
																							Conjunction: &parser.Conjunction{
																								Equality: &parser.Equality{
																									Comparison: &parser.Comparison{
																										Addition: &parser.Addition{
																											Multiplication: &parser.Multiplication{
																												Unary: &parser.Unary{
																													Accessor: &parser.Accessor{
																														Atom: &parser.Atom{Number: pointer.ToFloat64(42)},
																													},
																												},
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																				Commands: []*parser.Command{{Set: pointer.ToString("two")}},
																			},
																		},
																	},
																},
															},
														},
													},
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
			wantDeclaredIdentifiers: mapset.NewSet("test", "test2"),
			wantCommand: commands.NewLetCommand(
				"test2",
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(23),
						Command:   runtime.CommandGroup{commands.NewSetCommand("one")},
					},
					{
						Condition: expressions.NewNumber(42),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two")},
					},
				}),
			),
			wantTopLevelSettedState: "",
			wantSettedStates:        mapset.NewSet("one", "two"),
			wantReturn:              assert.False,
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/let/error",
			args: args{
				command: &parser.Command{
					Let: &parser.LetCommand{
						Identifier: "test2",
						Expression: &parser.Expression{
							ListConstruction: &parser.ListConstruction{
								Disjunction: &parser.Disjunction{
									Conjunction: &parser.Conjunction{
										Equality: &parser.Equality{
											Comparison: &parser.Comparison{
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
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             nil,
			wantTopLevelSettedState: "",
			wantReturn:              assert.False,
			wantErr:                 assert.Error,
		},
		{
			name: "Command/send",
			args: args{
				command:             &parser.Command{Send: pointer.ToString("test")},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             commands.NewSendCommand("test"),
			wantTopLevelSettedState: "",
			wantSettedStates:        mapset.NewSet(),
			wantReturn:              assert.False,
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/set",
			args: args{
				command:             &parser.Command{Set: pointer.ToString("test")},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             commands.NewSetCommand("test"),
			wantTopLevelSettedState: "test",
			wantSettedStates:        mapset.NewSet("test"),
			wantReturn:              assert.False,
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/return",
			args: args{
				command:             &parser.Command{Return: true},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             commands.ReturnCommand{},
			wantTopLevelSettedState: "",
			wantSettedStates:        mapset.NewSet(),
			wantReturn:              assert.True,
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/expression/success",
			args: args{
				command: &parser.Command{
					Expression: &parser.Expression{
						ListConstruction: &parser.ListConstruction{
							Disjunction: &parser.Disjunction{
								Conjunction: &parser.Conjunction{
									Equality: &parser.Equality{
										Comparison: &parser.Comparison{
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             commands.NewExpressionCommand(expressions.NewIdentifier("test")),
			wantTopLevelSettedState: "",
			wantSettedStates:        mapset.NewSet(),
			wantReturn:              assert.False,
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/expression/success/with setted states",
			args: args{
				command: &parser.Command{
					Expression: &parser.Expression{
						ListConstruction: &parser.ListConstruction{
							Disjunction: &parser.Disjunction{
								Conjunction: &parser.Conjunction{
									Equality: &parser.Equality{
										Comparison: &parser.Comparison{
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
																									Addition: &parser.Addition{
																										Multiplication: &parser.Multiplication{
																											Unary: &parser.Unary{
																												Accessor: &parser.Accessor{
																													Atom: &parser.Atom{Number: pointer.ToFloat64(23)},
																												},
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																			Commands: []*parser.Command{{Set: pointer.ToString("one")}},
																		},
																		{
																			Condition: &parser.Expression{
																				ListConstruction: &parser.ListConstruction{
																					Disjunction: &parser.Disjunction{
																						Conjunction: &parser.Conjunction{
																							Equality: &parser.Equality{
																								Comparison: &parser.Comparison{
																									Addition: &parser.Addition{
																										Multiplication: &parser.Multiplication{
																											Unary: &parser.Unary{
																												Accessor: &parser.Accessor{
																													Atom: &parser.Atom{Number: pointer.ToFloat64(42)},
																												},
																											},
																										},
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																			Commands: []*parser.Command{{Set: pointer.ToString("two")}},
																		},
																	},
																},
															},
														},
													},
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
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand: commands.NewExpressionCommand(
				expressions.NewConditionalExpression([]expressions.ConditionalCase{
					{
						Condition: expressions.NewNumber(23),
						Command:   runtime.CommandGroup{commands.NewSetCommand("one")},
					},
					{
						Condition: expressions.NewNumber(42),
						Command:   runtime.CommandGroup{commands.NewSetCommand("two")},
					},
				}),
			),
			wantTopLevelSettedState: "",
			wantSettedStates:        mapset.NewSet("one", "two"),
			wantReturn:              assert.False,
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/expression/error",
			args: args{
				command: &parser.Command{
					Expression: &parser.Expression{
						ListConstruction: &parser.ListConstruction{
							Disjunction: &parser.Disjunction{
								Conjunction: &parser.Conjunction{
									Equality: &parser.Equality{
										Comparison: &parser.Comparison{
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
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             nil,
			wantTopLevelSettedState: "",
			wantReturn:              assert.False,
			wantErr:                 assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			gotCommand, gotTopLevelSettedState, gotSettedStates, gotReturn, err :=
				translateCommand(testData.args.command, testData.args.declaredIdentifiers)

			assert.Equal(test, testData.wantDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.wantCommand, gotCommand)
			assert.Equal(test, testData.wantTopLevelSettedState, gotTopLevelSettedState)
			assert.Equal(test, testData.wantSettedStates, gotSettedStates)
			testData.wantReturn(test, gotReturn)
			testData.wantErr(test, err)
		})
	}
}

func cleanInboxes(actors runtime.ConcurrentActorGroup) {
	for index := range actors {
		actorPointer := &actors[index]
		inboxField := reflect.ValueOf(actorPointer).Elem().FieldByName("inbox")
		*(*chan string)(unsafe.Pointer(inboxField.UnsafeAddr())) = nil
	}
}
