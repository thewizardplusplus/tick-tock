package translator

import (
	"io"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/internal/tests"
	testsmocks "github.com/thewizardplusplus/tick-tock/internal/tests/mocks"
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
		declaredIdentifiers declaredIdentifierGroup
	}

	for _, testData := range []struct {
		name           string
		args           args
		makeWantActors func(options Options, dependencies Dependencies) runtime.ConcurrentActorGroup
		wantErr        assert.ErrorAssertionFunc
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
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantActors: func(options Options, dependencies Dependencies) runtime.ConcurrentActorGroup {
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
					runtime.NewConcurrentActor(actorOne, options.InboxSize, dependencies.Runtime),
					runtime.NewConcurrentActor(actorTwo, options.InboxSize, dependencies.Runtime),
				}
			},
			wantErr: assert.NoError,
		},
		{
			name: "success without actors",
			args: args{
				makeActors:          func(options Options) []*parser.Actor { return nil },
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantActors: func(options Options, dependencies Dependencies) runtime.ConcurrentActorGroup {
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
															Addition: &parser.Addition{
																Multiplication: &parser.Multiplication{
																	Unary: &parser.Unary{
																		Accessor: &parser.Accessor{
																			Atom: &parser.Atom{Number: tests.GetNumberAddress(23)},
																		},
																	},
																},
															},
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
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantActors: func(options Options, dependencies Dependencies) runtime.ConcurrentActorGroup {
				actorOne, _ := runtime.NewActor(
					runtime.StateGroup{
						options.InitialState: runtime.MessageGroup{
							"message_0": runtime.CommandGroup{commands.NewExpressionCommand(expressions.NewNumber(23))},
						},
					},
					options.InitialState,
				)
				return runtime.ConcurrentActorGroup{
					runtime.NewConcurrentActor(actorOne, options.InboxSize, dependencies.Runtime),
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
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantActors: func(options Options, dependencies Dependencies) runtime.ConcurrentActorGroup {
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
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantActors: func(options Options, dependencies Dependencies) runtime.ConcurrentActorGroup {
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
									},
								},
							},
						},
					}
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantActors: func(options Options, dependencies Dependencies) runtime.ConcurrentActorGroup {
				return nil
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := make(declaredIdentifierGroup)
			for identifier := range testData.args.declaredIdentifiers {
				originDeclaredIdentifiers[identifier] = struct{}{}
			}

			options := Options{tests.BufferedInbox, "__initialization__"}
			outWriter := new(testsmocks.Writer)
			randomizer := new(testsmocks.Randomizer)
			sleeper := new(testsmocks.Sleeper)
			waiter := new(waitermocks.Waiter)
			errorHandler := new(runtimemocks.ErrorHandler)
			dependencies := Dependencies{
				Commands: commands.Dependencies{
					OutWriter: outWriter,
					Sleep: commands.SleepDependencies{
						Randomizer: randomizer.Randomize,
						Sleeper:    sleeper.Sleep,
					},
				},
				Runtime: runtime.Dependencies{Waiter: waiter, ErrorHandler: errorHandler},
			}
			gotActors, err := Translate(
				testData.args.makeActors(options),
				testData.args.declaredIdentifiers,
				options,
				dependencies,
			)

			mock.AssertExpectationsForObjects(test, outWriter, randomizer, sleeper, waiter, errorHandler)
			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(
				test,
				cleanInboxes(testData.makeWantActors(options, dependencies)),
				cleanInboxes(gotActors),
			)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateStates(test *testing.T) {
	type args struct {
		states              []*parser.State
		declaredIdentifiers declaredIdentifierGroup
	}

	for _, testData := range []struct {
		name           string
		args           args
		makeWantStates func(outWriter io.Writer) runtime.StateGroup
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "success with nonempty states",
			args: args{
				states: []*parser.State{
					{Name: "state_0", Messages: []*parser.Message{{Name: "message_0"}, {Name: "message_1"}}},
					{Name: "state_1", Messages: []*parser.Message{{Name: "message_2"}, {Name: "message_3"}}},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup {
				return runtime.StateGroup{
					"state_0": runtime.MessageGroup{"message_0": nil, "message_1": nil},
					"state_1": runtime.MessageGroup{"message_2": nil, "message_3": nil},
				}
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with empty states",
			args: args{
				states:              []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup {
				return runtime.StateGroup{"state_0": runtime.MessageGroup{}, "state_1": runtime.MessageGroup{}}
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
							},
						},
					},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup {
				return runtime.StateGroup{
					"state_0": runtime.MessageGroup{
						"message_0": runtime.CommandGroup{commands.NewExpressionCommand(expressions.NewNumber(23))},
					},
				}
			},
			wantErr: assert.NoError,
		},
		{
			name: "error without states",
			args: args{
				states:              nil,
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup { return nil },
			wantErr:        assert.Error,
		},
		{
			name: "error with duplicate states",
			args: args{
				states:              []*parser.State{{Name: "test"}, {Name: "test"}},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup { return nil },
			wantErr:        assert.Error,
		},
		{
			name: "error with messages translation",
			args: args{
				states: []*parser.State{
					{Name: "state_0", Messages: []*parser.Message{{Name: "message_0"}, {Name: "message_1"}}},
					{Name: "state_1", Messages: []*parser.Message{{Name: "test"}, {Name: "test"}}},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup { return nil },
			wantErr:        assert.Error,
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
									{Send: tests.GetStringAddress("command_0")},
									{Set: tests.GetStringAddress("state_unknown")},
								},
							},
							{
								Name: "message_1",
								Commands: []*parser.Command{
									{Send: tests.GetStringAddress("command_2")},
									{Set: tests.GetStringAddress("state_unknown")},
								},
							},
						},
					},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup { return nil },
			wantErr:        assert.Error,
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
						},
					},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup { return nil },
			wantErr:        assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := make(declaredIdentifierGroup)
			for identifier := range testData.args.declaredIdentifiers {
				originDeclaredIdentifiers[identifier] = struct{}{}
			}

			outWriter := new(testsmocks.Writer)
			randomizer := new(testsmocks.Randomizer)
			sleeper := new(testsmocks.Sleeper)
			dependencies := commands.Dependencies{
				OutWriter: outWriter,
				Sleep: commands.SleepDependencies{
					Randomizer: randomizer.Randomize,
					Sleeper:    sleeper.Sleep,
				},
			}
			gotStates, err :=
				translateStates(testData.args.states, testData.args.declaredIdentifiers, dependencies)

			mock.AssertExpectationsForObjects(test, outWriter, randomizer, sleeper)
			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.makeWantStates(outWriter), gotStates)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateMessages(test *testing.T) {
	type args struct {
		messages            []*parser.Message
		declaredIdentifiers declaredIdentifierGroup
	}

	for _, testData := range []struct {
		name             string
		args             args
		makeWantMessages func(outWriter io.Writer) runtime.MessageGroup
		wantStates       settedStateGroup
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "success with nonempty messages (without set commands)",
			args: args{
				messages: []*parser.Message{
					{
						Name: "message_0",
						Commands: []*parser.Command{
							{Send: tests.GetStringAddress("command_0")},
							{Send: tests.GetStringAddress("command_1")},
						},
					},
					{
						Name: "message_1",
						Commands: []*parser.Command{
							{Send: tests.GetStringAddress("command_2")},
							{Send: tests.GetStringAddress("command_3")},
						},
					},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantMessages: func(outWriter io.Writer) runtime.MessageGroup {
				return runtime.MessageGroup{
					"message_0": runtime.CommandGroup{
						commands.NewSendCommand("command_0"),
						commands.NewSendCommand("command_1"),
					},
					"message_1": runtime.CommandGroup{
						commands.NewSendCommand("command_2"),
						commands.NewSendCommand("command_3"),
					},
				}
			},
			wantStates: make(settedStateGroup),
			wantErr:    assert.NoError,
		},
		{
			name: "success with nonempty messages (with set commands)",
			args: args{
				messages: []*parser.Message{
					{
						Name: "message_0",
						Commands: []*parser.Command{
							{Send: tests.GetStringAddress("command_0")},
							{Set: tests.GetStringAddress("command_1")},
						},
					},
					{
						Name: "message_1",
						Commands: []*parser.Command{
							{Send: tests.GetStringAddress("command_2")},
							{Set: tests.GetStringAddress("command_3")},
						},
					},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantMessages: func(outWriter io.Writer) runtime.MessageGroup {
				return runtime.MessageGroup{
					"message_0": runtime.CommandGroup{
						commands.NewSendCommand("command_0"),
						commands.NewSetCommand("command_1"),
					},
					"message_1": runtime.CommandGroup{
						commands.NewSendCommand("command_2"),
						commands.NewSetCommand("command_3"),
					},
				}
			},
			wantStates: settedStateGroup{"message_0": "command_1", "message_1": "command_3"},
			wantErr:    assert.NoError,
		},
		{
			name: "success with empty messages",
			args: args{
				messages:            []*parser.Message{{Name: "message_0"}, {Name: "message_1"}},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantMessages: func(outWriter io.Writer) runtime.MessageGroup {
				return runtime.MessageGroup{"message_0": nil, "message_1": nil}
			},
			wantStates: make(settedStateGroup),
			wantErr:    assert.NoError,
		},
		{
			name: "success without messages",
			args: args{
				messages:            nil,
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantMessages: func(outWriter io.Writer) runtime.MessageGroup {
				return runtime.MessageGroup{}
			},
			wantStates: make(settedStateGroup),
			wantErr:    assert.NoError,
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
					},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantMessages: func(outWriter io.Writer) runtime.MessageGroup {
				return runtime.MessageGroup{
					"message_0": runtime.CommandGroup{commands.NewExpressionCommand(expressions.NewNumber(23))},
				}
			},
			wantStates: make(settedStateGroup),
			wantErr:    assert.NoError,
		},
		{
			name: "error with duplicate messages",
			args: args{
				messages:            []*parser.Message{{Name: "test"}, {Name: "test"}},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantMessages: func(outWriter io.Writer) runtime.MessageGroup { return nil },
			wantStates:       nil,
			wantErr:          assert.Error,
		},
		{
			name: "error with commands translation",
			args: args{
				messages: []*parser.Message{
					{
						Name: "message_0",
						Commands: []*parser.Command{
							{Send: tests.GetStringAddress("command_0")},
							{Send: tests.GetStringAddress("command_1")},
						},
					},
					{
						Name: "message_1",
						Commands: []*parser.Command{
							{Send: tests.GetStringAddress("command_2")},
							{Set: tests.GetStringAddress("command_3")},
							{Send: tests.GetStringAddress("command_4")},
							{Set: tests.GetStringAddress("command_5")},
						},
					},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantMessages: func(outWriter io.Writer) runtime.MessageGroup { return nil },
			wantStates:       nil,
			wantErr:          assert.Error,
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
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantMessages: func(outWriter io.Writer) runtime.MessageGroup { return nil },
			wantStates:       nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := make(declaredIdentifierGroup)
			for identifier := range testData.args.declaredIdentifiers {
				originDeclaredIdentifiers[identifier] = struct{}{}
			}

			outWriter := new(testsmocks.Writer)
			randomizer := new(testsmocks.Randomizer)
			sleeper := new(testsmocks.Sleeper)
			dependencies := commands.Dependencies{
				OutWriter: outWriter,
				Sleep: commands.SleepDependencies{
					Randomizer: randomizer.Randomize,
					Sleeper:    sleeper.Sleep,
				},
			}
			gotMessages, gotStates, err :=
				translateMessages(testData.args.messages, testData.args.declaredIdentifiers, dependencies)

			mock.AssertExpectationsForObjects(test, outWriter, randomizer, sleeper)
			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.makeWantMessages(outWriter), gotMessages)
			assert.Equal(test, testData.wantStates, gotStates)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateCommands(test *testing.T) {
	type args struct {
		commands            []*parser.Command
		declaredIdentifiers declaredIdentifierGroup
	}

	for _, testData := range []struct {
		name             string
		args             args
		makeWantCommands func(outWriter io.Writer) runtime.CommandGroup
		wantState        string
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "success with commands (without a set command)",
			args: args{
				commands: []*parser.Command{
					{Send: tests.GetStringAddress("one")},
					{Send: tests.GetStringAddress("two")},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantCommands: func(outWriter io.Writer) runtime.CommandGroup {
				return runtime.CommandGroup{commands.NewSendCommand("one"), commands.NewSendCommand("two")}
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with commands (with a set command)",
			args: args{
				commands: []*parser.Command{
					{Send: tests.GetStringAddress("one")},
					{Set: tests.GetStringAddress("two")},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantCommands: func(outWriter io.Writer) runtime.CommandGroup {
				return runtime.CommandGroup{commands.NewSendCommand("one"), commands.NewSetCommand("two")}
			},
			wantState: "two",
			wantErr:   assert.NoError,
		},
		{
			name: "success with commands (with using an existing identifier)",
			args: args{
				commands: []*parser.Command{
					{
						Expression: &parser.Expression{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{
												Atom: &parser.Atom{Identifier: tests.GetStringAddress("test")},
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantCommands: func(outWriter io.Writer) runtime.CommandGroup {
				return runtime.CommandGroup{commands.NewExpressionCommand(expressions.NewIdentifier("test"))}
			},
			wantState: "",
			wantErr:   assert.NoError,
		},
		{
			name: "success with commands (with using an nonexistent identifier)",
			args: args{
				commands: []*parser.Command{
					{
						Let: &parser.LetCommand{
							Identifier: "test2",
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
					},
					{
						Expression: &parser.Expression{
							ListConstruction: &parser.ListConstruction{
								Addition: &parser.Addition{
									Multiplication: &parser.Multiplication{
										Unary: &parser.Unary{
											Accessor: &parser.Accessor{
												Atom: &parser.Atom{Identifier: tests.GetStringAddress("test2")},
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantCommands: func(outWriter io.Writer) runtime.CommandGroup {
				return runtime.CommandGroup{
					commands.NewLetCommand("test2", expressions.NewNumber(23)),
					commands.NewExpressionCommand(expressions.NewIdentifier("test2")),
				}
			},
			wantState: "",
			wantErr:   assert.NoError,
		},
		{
			name: "success without commands",
			args: args{
				commands:            nil,
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantCommands: func(outWriter io.Writer) runtime.CommandGroup { return nil },
			wantErr:          assert.NoError,
		},
		{
			name: "error with command translation",
			args: args{
				commands: []*parser.Command{
					{
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
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantCommands: func(outWriter io.Writer) runtime.CommandGroup { return nil },
			wantErr:          assert.Error,
		},
		{
			name: "error with a second set command",
			args: args{
				commands: []*parser.Command{
					{Send: tests.GetStringAddress("one")},
					{Set: tests.GetStringAddress("two")},
					{Send: tests.GetStringAddress("three")},
					{Set: tests.GetStringAddress("four")},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			makeWantCommands: func(outWriter io.Writer) runtime.CommandGroup { return nil },
			wantErr:          assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := make(declaredIdentifierGroup)
			for identifier := range testData.args.declaredIdentifiers {
				originDeclaredIdentifiers[identifier] = struct{}{}
			}

			outWriter := new(testsmocks.Writer)
			randomizer := new(testsmocks.Randomizer)
			sleeper := new(testsmocks.Sleeper)
			dependencies := commands.Dependencies{
				OutWriter: outWriter,
				Sleep: commands.SleepDependencies{
					Randomizer: randomizer.Randomize,
					Sleeper:    sleeper.Sleep,
				},
			}
			gotCommands, gotState, err :=
				translateCommands(testData.args.commands, testData.args.declaredIdentifiers, dependencies)

			mock.AssertExpectationsForObjects(test, outWriter, randomizer, sleeper)
			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.makeWantCommands(outWriter), gotCommands)
			assert.Equal(test, testData.wantState, gotState)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateCommand(test *testing.T) {
	type args struct {
		command             *parser.Command
		declaredIdentifiers declaredIdentifierGroup
	}

	for _, testData := range []struct {
		name                    string
		args                    args
		wantDeclaredIdentifiers declaredIdentifierGroup
		makeWantCommand         func(outWriter io.Writer) runtime.Command
		wantState               string
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
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantDeclaredIdentifiers: declaredIdentifierGroup{"test": {}, "test2": {}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewLetCommand("test2", expressions.NewNumber(23))
			},
			wantState: "",
			wantErr:   assert.NoError,
		},
		{
			name: "Command/let/success/existing identifier",
			args: args{
				command: &parser.Command{
					Let: &parser.LetCommand{
						Identifier: "test",
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
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantDeclaredIdentifiers: declaredIdentifierGroup{"test": {}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewLetCommand("test", expressions.NewNumber(23))
			},
			wantState: "",
			wantErr:   assert.NoError,
		},
		{
			name: "Command/let/error",
			args: args{
				command: &parser.Command{
					Let: &parser.LetCommand{
						Identifier: "test2",
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
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantDeclaredIdentifiers: declaredIdentifierGroup{"test": {}},
			makeWantCommand:         func(outWriter io.Writer) runtime.Command { return nil },
			wantState:               "",
			wantErr:                 assert.Error,
		},
		{
			name: "Command/send",
			args: args{
				command:             &parser.Command{Send: tests.GetStringAddress("test")},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantDeclaredIdentifiers: declaredIdentifierGroup{"test": {}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewSendCommand("test")
			},
			wantState: "",
			wantErr:   assert.NoError,
		},
		{
			name: "Command/set",
			args: args{
				command:             &parser.Command{Set: tests.GetStringAddress("test")},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantDeclaredIdentifiers: declaredIdentifierGroup{"test": {}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewSetCommand("test")
			},
			wantState: "test",
			wantErr:   assert.NoError,
		},
		{
			name: "Command/out/nonempty",
			args: args{
				command:             &parser.Command{Out: tests.GetStringAddress("test")},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantDeclaredIdentifiers: declaredIdentifierGroup{"test": {}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewOutCommand("test", outWriter)
			},
			wantState: "",
			wantErr:   assert.NoError,
		},
		{
			name: "Command/out/empty",
			args: args{
				command:             &parser.Command{Out: tests.GetStringAddress("")},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantDeclaredIdentifiers: declaredIdentifierGroup{"test": {}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewOutCommand("", outWriter)
			},
			wantState: "",
			wantErr:   assert.NoError,
		},
		{
			name: "Command/sleep/success",
			args: args{
				command: &parser.Command{
					Sleep: &parser.SleepCommand{
						Minimum: tests.GetNumberAddress(1.2),
						Maximum: tests.GetNumberAddress(3.4),
					},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantDeclaredIdentifiers: declaredIdentifierGroup{"test": {}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				command, _ := commands.NewSleepCommand(1.2, 3.4, commands.SleepDependencies{})
				return command
			},
			wantState: "",
			wantErr:   assert.NoError,
		},
		{
			name: "Command/sleep/error",
			args: args{
				command: &parser.Command{
					Sleep: &parser.SleepCommand{
						Minimum: tests.GetNumberAddress(3.4),
						Maximum: tests.GetNumberAddress(1.2),
					},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantDeclaredIdentifiers: declaredIdentifierGroup{"test": {}},
			makeWantCommand:         func(outWriter io.Writer) runtime.Command { return nil },
			wantState:               "",
			wantErr:                 assert.Error,
		},
		{
			name: "Command/exit",
			args: args{
				command:             &parser.Command{Exit: true},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantDeclaredIdentifiers: declaredIdentifierGroup{"test": {}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.ExitCommand{}
			},
			wantState: "",
			wantErr:   assert.NoError,
		},
		{
			name: "Command/expression/success",
			args: args{
				command: &parser.Command{
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
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantDeclaredIdentifiers: declaredIdentifierGroup{"test": {}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewExpressionCommand(expressions.NewNumber(23))
			},
			wantState: "",
			wantErr:   assert.NoError,
		},
		{
			name: "Command/expression/error",
			args: args{
				command: &parser.Command{
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
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantDeclaredIdentifiers: declaredIdentifierGroup{"test": {}},
			makeWantCommand:         func(outWriter io.Writer) runtime.Command { return nil },
			wantState:               "",
			wantErr:                 assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			outWriter := new(testsmocks.Writer)
			randomizer := new(testsmocks.Randomizer)
			sleeper := new(testsmocks.Sleeper)
			dependencies := commands.Dependencies{
				OutWriter: outWriter,
				Sleep: commands.SleepDependencies{
					Randomizer: randomizer.Randomize,
					Sleeper:    sleeper.Sleep,
				},
			}
			gotCommand, gotState, err :=
				translateCommand(testData.args.command, testData.args.declaredIdentifiers, dependencies)
			if sleepCommand, ok := gotCommand.(commands.SleepCommand); ok {
				cleanSleepDependencies(&sleepCommand)
				gotCommand = sleepCommand
			}

			mock.AssertExpectationsForObjects(test, outWriter, randomizer, sleeper)
			assert.Equal(test, testData.wantDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.makeWantCommand(outWriter), gotCommand)
			assert.Equal(test, testData.wantState, gotState)
			testData.wantErr(test, err)
		})
	}
}

func cleanInboxes(actors runtime.ConcurrentActorGroup) runtime.ConcurrentActorGroup {
	actorsReflection := reflect.ValueOf(actors)
	for index := range actors {
		fieldAddress := getFieldAddress(actorsReflection.Index(index), "inbox")
		*(*chan string)(fieldAddress) = nil
	}

	return actors
}

func cleanSleepDependencies(command *commands.SleepCommand) {
	fieldAddress := getFieldAddress(reflect.ValueOf(command).Elem(), "dependencies")
	*(*commands.SleepDependencies)(fieldAddress) = commands.SleepDependencies{}
}

func getFieldAddress(value reflect.Value, name string) unsafe.Pointer {
	return unsafe.Pointer(value.FieldByName(name).UnsafeAddr())
}
