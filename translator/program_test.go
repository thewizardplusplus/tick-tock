package translator

import (
	"reflect"
	"testing"
	"unsafe"

	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/internal/tests"
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

			options := Options{tests.BufferedInbox, "__initialization__"}
			waiter := new(waitermocks.Waiter)
			errorHandler := new(runtimemocks.ErrorHandler)
			dependencies := runtime.Dependencies{Waiter: waiter, ErrorHandler: errorHandler}
			gotActors, err := Translate(
				testData.args.makeActors(options),
				testData.args.declaredIdentifiers,
				options,
				dependencies,
			)

			mock.AssertExpectationsForObjects(test, waiter, errorHandler)
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
		name         string
		args         args
		wantMessages runtime.MessageGroup
		wantStates   settedStateGroup
		wantErr      assert.ErrorAssertionFunc
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
			wantStates: settedStateGroup{"message_0": "command_1", "message_1": "command_3"},
			wantErr:    assert.NoError,
		},
		{
			name: "success with empty messages",
			args: args{
				messages:            []*parser.Message{{Name: "message_0"}, {Name: "message_1"}},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages: runtime.MessageGroup{"message_0": nil, "message_1": nil},
			wantStates:   make(settedStateGroup),
			wantErr:      assert.NoError,
		},
		{
			name: "success without messages",
			args: args{
				messages:            nil,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages: runtime.MessageGroup{},
			wantStates:   make(settedStateGroup),
			wantErr:      assert.NoError,
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
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages: runtime.MessageGroup{
				"message_0": runtime.CommandGroup{
					commands.NewExpressionCommand(expressions.NewIdentifier("test")),
				},
			},
			wantStates: make(settedStateGroup),
			wantErr:    assert.NoError,
		},
		{
			name: "error with duplicate messages",
			args: args{
				messages:            []*parser.Message{{Name: "test"}, {Name: "test"}},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages: nil,
			wantStates:   nil,
			wantErr:      assert.Error,
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages: nil,
			wantStates:   nil,
			wantErr:      assert.Error,
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages: nil,
			wantStates:   nil,
			wantErr:      assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := testData.args.declaredIdentifiers.Clone()

			gotMessages, gotStates, err :=
				translateMessages(testData.args.messages, testData.args.declaredIdentifiers)

			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.wantMessages, gotMessages)
			assert.Equal(test, testData.wantStates, gotStates)
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
		name         string
		args         args
		wantCommands runtime.CommandGroup
		wantState    string
		wantErr      assert.ErrorAssertionFunc
	}{
		{
			name: "success with commands (without a set command)",
			args: args{
				commands: []*parser.Command{
					{Send: tests.GetStringAddress("one")},
					{Send: tests.GetStringAddress("two")},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: runtime.CommandGroup{
				commands.NewSendCommand("one"),
				commands.NewSendCommand("two"),
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: runtime.CommandGroup{
				commands.NewSendCommand("one"),
				commands.NewSetCommand("two"),
			},
			wantState: "two",
			wantErr:   assert.NoError,
		},
		{
			name: "success with commands (with an expression command and an existing identifier)",
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: runtime.CommandGroup{
				commands.NewExpressionCommand(expressions.NewIdentifier("test")),
			},
			wantState: "",
			wantErr:   assert.NoError,
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: runtime.CommandGroup{
				commands.NewLetCommand("test2", expressions.NewNumber(23)),
				commands.NewExpressionCommand(expressions.NewIdentifier("test2")),
			},
			wantState: "",
			wantErr:   assert.NoError,
		},
		{
			name: "success without commands",
			args: args{
				commands:            nil,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: nil,
			wantErr:      assert.NoError,
		},
		{
			name: "error with expression command translation",
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: nil,
			wantErr:      assert.Error,
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: nil,
			wantErr:      assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := testData.args.declaredIdentifiers.Clone()

			gotCommands, gotState, err :=
				translateCommands(testData.args.commands, testData.args.declaredIdentifiers)

			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.wantCommands, gotCommands)
			assert.Equal(test, testData.wantState, gotState)
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantDeclaredIdentifiers: mapset.NewSet("test", "test2"),
			wantCommand:             commands.NewLetCommand("test2", expressions.NewNumber(23)),
			wantState:               "",
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             commands.NewLetCommand("test", expressions.NewNumber(23)),
			wantState:               "",
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             nil,
			wantState:               "",
			wantErr:                 assert.Error,
		},
		{
			name: "Command/send",
			args: args{
				command:             &parser.Command{Send: tests.GetStringAddress("test")},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             commands.NewSendCommand("test"),
			wantState:               "",
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/set",
			args: args{
				command:             &parser.Command{Set: tests.GetStringAddress("test")},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             commands.NewSetCommand("test"),
			wantState:               "test",
			wantErr:                 assert.NoError,
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
										Accessor: &parser.Accessor{
											Atom: &parser.Atom{Identifier: tests.GetStringAddress("test")},
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
			wantState:               "",
			wantErr:                 assert.NoError,
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             nil,
			wantState:               "",
			wantErr:                 assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			gotCommand, gotState, err :=
				translateCommand(testData.args.command, testData.args.declaredIdentifiers)

			assert.Equal(test, testData.wantDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.wantCommand, gotCommand)
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

func getFieldAddress(value reflect.Value, name string) unsafe.Pointer {
	return unsafe.Pointer(value.FieldByName(name).UnsafeAddr())
}
