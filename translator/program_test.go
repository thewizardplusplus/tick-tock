package translator

import (
	"testing"

	"github.com/AlekSi/pointer"
	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/commands"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
)

func TestTranslate(test *testing.T) {
	type args struct {
		program             *parser.Program
		declaredIdentifiers mapset.Set
		options             Options
		dependencies        runtime.Dependencies
	}

	for _, testData := range []struct {
		name                 string
		args                 args
		wantDefinitions      context.ValueGroup
		wantTranslatedActors []runtime.ConcurrentActorFactory
		wantErr              assert.ErrorAssertionFunc
	}{
		{
			name: "success without definitions",
			args: args{
				program: &parser.Program{
					Definitions: nil,
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantDefinitions:      context.ValueGroup{},
			wantTranslatedActors: nil,
			wantErr:              assert.NoError,
		},
		{
			name: "success with only few actors",
			args: args{
				program: &parser.Program{
					Definitions: []*parser.Definition{
						{
							Actor: &parser.Actor{
								Name:   "Test0",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
						{
							Actor: &parser.Actor{
								Name:   "Test1",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantDefinitions: context.ValueGroup{
				"Test0": func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test0",
						runtime.ParameterizedStateGroup{
							StateGroup: runtime.StateGroup{
								"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
								"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							},
						},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
						WaitGroup:    new(MockWaiter),
						ErrorHandler: new(MockErrorHandler),
					})
				}(),
				"Test1": func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test1",
						runtime.ParameterizedStateGroup{
							StateGroup: runtime.StateGroup{
								"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
								"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							},
						},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
						WaitGroup:    new(MockWaiter),
						ErrorHandler: new(MockErrorHandler),
					})
				}(),
			},
			wantTranslatedActors: []runtime.ConcurrentActorFactory{
				func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test0",
						runtime.ParameterizedStateGroup{
							StateGroup: runtime.StateGroup{
								"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
								"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							},
						},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
						WaitGroup:    new(MockWaiter),
						ErrorHandler: new(MockErrorHandler),
					})
				}(),
				func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test1",
						runtime.ParameterizedStateGroup{
							StateGroup: runtime.StateGroup{
								"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
								"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							},
						},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
						WaitGroup:    new(MockWaiter),
						ErrorHandler: new(MockErrorHandler),
					})
				}(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with only few actor classes",
			args: args{
				program: &parser.Program{
					Definitions: []*parser.Definition{
						{
							ActorClass: &parser.ActorClass{
								Name:   "Test0",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
						{
							ActorClass: &parser.ActorClass{
								Name:   "Test1",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantDefinitions: context.ValueGroup{
				"Test0": func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test0",
						runtime.ParameterizedStateGroup{
							StateGroup: runtime.StateGroup{
								"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
								"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							},
						},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
						WaitGroup:    new(MockWaiter),
						ErrorHandler: new(MockErrorHandler),
					})
				}(),
				"Test1": func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test1",
						runtime.ParameterizedStateGroup{
							StateGroup: runtime.StateGroup{
								"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
								"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							},
						},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
						WaitGroup:    new(MockWaiter),
						ErrorHandler: new(MockErrorHandler),
					})
				}(),
			},
			wantTranslatedActors: nil,
			wantErr:              assert.NoError,
		},
		{
			name: "success with few actors and actor classes",
			args: args{
				program: &parser.Program{
					Definitions: []*parser.Definition{
						{
							Actor: &parser.Actor{
								Name:   "Test0",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
						{
							ActorClass: &parser.ActorClass{
								Name:   "Test1",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
						{
							Actor: &parser.Actor{
								Name:   "Test2",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
						{
							ActorClass: &parser.ActorClass{
								Name:   "Test3",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantDefinitions: context.ValueGroup{
				"Test0": func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test0",
						runtime.ParameterizedStateGroup{
							StateGroup: runtime.StateGroup{
								"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
								"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							},
						},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
						WaitGroup:    new(MockWaiter),
						ErrorHandler: new(MockErrorHandler),
					})
				}(),
				"Test1": func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test1",
						runtime.ParameterizedStateGroup{
							StateGroup: runtime.StateGroup{
								"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
								"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							},
						},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
						WaitGroup:    new(MockWaiter),
						ErrorHandler: new(MockErrorHandler),
					})
				}(),
				"Test2": func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test2",
						runtime.ParameterizedStateGroup{
							StateGroup: runtime.StateGroup{
								"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
								"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							},
						},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
						WaitGroup:    new(MockWaiter),
						ErrorHandler: new(MockErrorHandler),
					})
				}(),
				"Test3": func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test3",
						runtime.ParameterizedStateGroup{
							StateGroup: runtime.StateGroup{
								"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
								"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							},
						},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
						WaitGroup:    new(MockWaiter),
						ErrorHandler: new(MockErrorHandler),
					})
				}(),
			},
			wantTranslatedActors: []runtime.ConcurrentActorFactory{
				func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test0",
						runtime.ParameterizedStateGroup{
							StateGroup: runtime.StateGroup{
								"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
								"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							},
						},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
						WaitGroup:    new(MockWaiter),
						ErrorHandler: new(MockErrorHandler),
					})
				}(),
				func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test2",
						runtime.ParameterizedStateGroup{
							StateGroup: runtime.StateGroup{
								"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
								"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							},
						},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
						WaitGroup:    new(MockWaiter),
						ErrorHandler: new(MockErrorHandler),
					})
				}(),
			},
			wantErr: assert.NoError,
		},
		{
			name: "error with definition translation",
			args: args{
				program: &parser.Program{
					Definitions: []*parser.Definition{
						{
							Actor: &parser.Actor{
								Name:   "Test0",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
						{
							Actor: &parser.Actor{
								Name:   "Test1",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_0"}},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantDefinitions:      nil,
			wantTranslatedActors: nil,
			wantErr:              assert.Error,
		},
		{
			name: "error with duplicate actors",
			args: args{
				program: &parser.Program{
					Definitions: []*parser.Definition{
						{
							Actor: &parser.Actor{
								Name:   "Test0",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
						{
							Actor: &parser.Actor{
								Name:   "Test0",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantDefinitions:      nil,
			wantTranslatedActors: nil,
			wantErr:              assert.Error,
		},
		{
			name: "error with duplicate actor classes",
			args: args{
				program: &parser.Program{
					Definitions: []*parser.Definition{
						{
							ActorClass: &parser.ActorClass{
								Name:   "Test0",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
						{
							ActorClass: &parser.ActorClass{
								Name:   "Test0",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantDefinitions:      nil,
			wantTranslatedActors: nil,
			wantErr:              assert.Error,
		},
		{
			name: "error with the actor and the actor class with the same name",
			args: args{
				program: &parser.Program{
					Definitions: []*parser.Definition{
						{
							Actor: &parser.Actor{
								Name:   "Test0",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
						{
							ActorClass: &parser.ActorClass{
								Name:   "Test0",
								States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantDefinitions:      nil,
			wantTranslatedActors: nil,
			wantErr:              assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := testData.args.declaredIdentifiers.Clone()

			gotDefinitions, gotTranslatedActors, err := Translate(
				testData.args.program,
				testData.args.declaredIdentifiers,
				testData.args.options,
				testData.args.dependencies,
			)

			mock.AssertExpectationsForObjects(
				test,
				testData.args.dependencies.WaitGroup,
				testData.args.dependencies.ErrorHandler,
			)
			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.wantDefinitions, gotDefinitions)
			assert.Equal(test, testData.wantTranslatedActors, gotTranslatedActors)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateDefinition(test *testing.T) {
	type args struct {
		definition          *parser.Definition
		declaredIdentifiers mapset.Set
		options             Options
		dependencies        runtime.Dependencies
	}

	for _, testData := range []struct {
		name                     string
		args                     args
		wantDeclaredIdentifiers  mapset.Set
		wantTranslatedActorClass runtime.ConcurrentActorFactory
		wantActor                assert.BoolAssertionFunc
		wantErr                  assert.ErrorAssertionFunc
	}{
		{
			name: "Definition/actor/success",
			args: args{
				definition: &parser.Definition{
					Actor: &parser.Actor{
						Name:   "Test",
						States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantDeclaredIdentifiers: mapset.NewSet("test", "Test"),
			wantTranslatedActorClass: func() runtime.ConcurrentActorFactory {
				actorFactory, _ := runtime.NewActorFactory(
					"Test",
					runtime.ParameterizedStateGroup{
						StateGroup: runtime.StateGroup{
							"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
						},
					},
					context.State{Name: "state_0"},
				)
				return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				})
			}(),
			wantActor: assert.True,
			wantErr:   assert.NoError,
		},
		{
			name: "Definition/actor/error",
			args: args{
				definition: &parser.Definition{
					Actor: &parser.Actor{
						Name:   "Test",
						States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "unknown"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantDeclaredIdentifiers:  mapset.NewSet("test"),
			wantTranslatedActorClass: runtime.ConcurrentActorFactory{},
			wantActor:                assert.False,
			wantErr:                  assert.Error,
		},
		{
			name: "Definition/actor class/success",
			args: args{
				definition: &parser.Definition{
					ActorClass: &parser.ActorClass{
						Name:   "Test",
						States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantDeclaredIdentifiers: mapset.NewSet("test", "Test"),
			wantTranslatedActorClass: func() runtime.ConcurrentActorFactory {
				actorFactory, _ := runtime.NewActorFactory(
					"Test",
					runtime.ParameterizedStateGroup{
						StateGroup: runtime.StateGroup{
							"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
						},
					},
					context.State{Name: "state_0"},
				)
				return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				})
			}(),
			wantActor: assert.False,
			wantErr:   assert.NoError,
		},
		{
			name: "Definition/actor class/error",
			args: args{
				definition: &parser.Definition{
					ActorClass: &parser.ActorClass{
						Name:   "Test",
						States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "unknown"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantDeclaredIdentifiers:  mapset.NewSet("test"),
			wantTranslatedActorClass: runtime.ConcurrentActorFactory{},
			wantActor:                assert.False,
			wantErr:                  assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			gotTranslatedActorClass, gotActor, err := translateDefinition(
				testData.args.definition,
				testData.args.declaredIdentifiers,
				testData.args.options,
				testData.args.dependencies,
			)

			mock.AssertExpectationsForObjects(
				test,
				testData.args.dependencies.WaitGroup,
				testData.args.dependencies.ErrorHandler,
			)
			assert.Equal(test, testData.wantDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.wantTranslatedActorClass, gotTranslatedActorClass)
			testData.wantActor(test, gotActor)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateActorClass(test *testing.T) {
	type args struct {
		actorClass          *parser.ActorClass
		declaredIdentifiers mapset.Set
		options             Options
		dependencies        runtime.Dependencies
	}

	for _, testData := range []struct {
		name                     string
		args                     args
		wantTranslatedActorClass runtime.ConcurrentActorFactory
		wantErr                  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				actorClass: &parser.ActorClass{
					Name:   "Test",
					States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantTranslatedActorClass: func() runtime.ConcurrentActorFactory {
				actorFactory, _ := runtime.NewActorFactory(
					"Test",
					runtime.ParameterizedStateGroup{
						StateGroup: runtime.StateGroup{
							"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
							"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
						},
					},
					context.State{Name: "state_0"},
				)
				return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				})
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "success with the expression",
			args: args{
				actorClass: &parser.ActorClass{
					Name: "Test",
					States: []*parser.State{
						{
							Name: "state_0",
							Messages: []*parser.Message{
								{
									Name: "message_0",
									Commands: []*parser.Command{
										{
											Expression: &parser.Expression{
												ListConstruction: &parser.ListConstruction{
													NilCoalescing: &parser.NilCoalescing{
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
											},
										},
									},
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantTranslatedActorClass: func() runtime.ConcurrentActorFactory {
				actorFactory, _ := runtime.NewActorFactory(
					"Test",
					runtime.ParameterizedStateGroup{
						StateGroup: runtime.StateGroup{
							"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{
								"message_0": runtime.NewParameterizedCommandGroup(nil, runtime.CommandGroup{
									commands.NewExpressionCommand(expressions.NewIdentifier("test")),
								}),
							}),
						},
					},
					context.State{Name: "state_0"},
				)
				return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				})
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "success with parameters",
			args: args{
				actorClass: &parser.ActorClass{
					Name:       "Test",
					Parameters: []string{"one", "two"},
					States: []*parser.State{
						{
							Name: "state_0",
							Messages: []*parser.Message{
								{
									Name: "message_0",
									Commands: []*parser.Command{
										{
											Expression: &parser.Expression{
												ListConstruction: &parser.ListConstruction{
													NilCoalescing: &parser.NilCoalescing{
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
																										Atom: &parser.Atom{Identifier: pointer.ToString("one")},
																									},
																								},
																							},
																						},
																					},
																				},
																			},
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
													NilCoalescing: &parser.NilCoalescing{
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
																										Atom: &parser.Atom{Identifier: pointer.ToString("two")},
																									},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
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
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantTranslatedActorClass: func() runtime.ConcurrentActorFactory {
				actorFactory, _ := runtime.NewActorFactory(
					"Test",
					runtime.NewParameterizedStateGroup([]string{"one", "two"}, runtime.StateGroup{
						"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{
							"message_0": runtime.NewParameterizedCommandGroup(nil, runtime.CommandGroup{
								commands.NewExpressionCommand(expressions.NewIdentifier("one")),
								commands.NewExpressionCommand(expressions.NewIdentifier("two")),
							}),
						}),
					}),
					context.State{Name: "state_0"},
				)
				return runtime.NewConcurrentActorFactory(actorFactory, 23, runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				})
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "error with states translation",
			args: args{
				actorClass: &parser.ActorClass{
					Name:   "Test",
					States: []*parser.State{{Name: "state_0"}, {Name: "state_0"}},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantTranslatedActorClass: runtime.ConcurrentActorFactory{},
			wantErr:                  assert.Error,
		},
		{
			name: "error with factory construction",
			args: args{
				actorClass: &parser.ActorClass{
					Name:   "Test",
					States: []*parser.State{{Name: "state_0"}, {Name: "state_1"}},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "unknown"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantTranslatedActorClass: runtime.ConcurrentActorFactory{},
			wantErr:                  assert.Error,
		},
		{
			name: "error with the expression",
			args: args{
				actorClass: &parser.ActorClass{
					Name: "Test",
					States: []*parser.State{
						{
							Name: "state_0",
							Messages: []*parser.Message{
								{
									Name: "message_0",
									Commands: []*parser.Command{
										{
											Expression: &parser.Expression{
												ListConstruction: &parser.ListConstruction{
													NilCoalescing: &parser.NilCoalescing{
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
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
				options:             Options{InboxSize: 23, InitialState: context.State{Name: "state_0"}},
				dependencies: runtime.Dependencies{
					WaitGroup:    new(MockWaiter),
					ErrorHandler: new(MockErrorHandler),
				},
			},
			wantTranslatedActorClass: runtime.ConcurrentActorFactory{},
			wantErr:                  assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := testData.args.declaredIdentifiers.Clone()

			gotTranslatedActorClass, err := translateActorClass(
				testData.args.actorClass,
				testData.args.declaredIdentifiers,
				testData.args.options,
				testData.args.dependencies,
			)

			mock.AssertExpectationsForObjects(
				test,
				testData.args.dependencies.WaitGroup,
				testData.args.dependencies.ErrorHandler,
			)
			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.wantTranslatedActorClass, gotTranslatedActorClass)
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
				"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{
					"message_0": {},
					"message_1": {},
				}),
				"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{
					"message_2": {},
					"message_3": {},
				}),
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
				"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
				"state_1": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{}),
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
												NilCoalescing: &parser.NilCoalescing{
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
				"state_0": runtime.NewParameterizedMessageGroup(nil, runtime.MessageGroup{
					"message_0": runtime.NewParameterizedCommandGroup(nil, runtime.CommandGroup{
						commands.NewExpressionCommand(expressions.NewIdentifier("test")),
					}),
				}),
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with parameters",
			args: args{
				states: []*parser.State{
					{
						Name:       "state_0",
						Parameters: []string{"one", "two"},
						Messages: []*parser.Message{
							{
								Name: "message_0",
								Commands: []*parser.Command{
									{
										Expression: &parser.Expression{
											ListConstruction: &parser.ListConstruction{
												NilCoalescing: &parser.NilCoalescing{
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
																									Atom: &parser.Atom{Identifier: pointer.ToString("one")},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
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
												NilCoalescing: &parser.NilCoalescing{
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
																									Atom: &parser.Atom{Identifier: pointer.ToString("two")},
																								},
																							},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
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
				"state_0": runtime.NewParameterizedMessageGroup([]string{"one", "two"}, runtime.MessageGroup{
					"message_0": runtime.NewParameterizedCommandGroup(nil, runtime.CommandGroup{
						commands.NewExpressionCommand(expressions.NewIdentifier("one")),
						commands.NewExpressionCommand(expressions.NewIdentifier("two")),
					}),
				}),
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
									{Send: &parser.SendCommand{Name: "command_0"}},
									{Set: &parser.SetCommand{Name: "state_unknown"}},
								},
							},
							{
								Name: "message_1",
								Commands: []*parser.Command{
									{Send: &parser.SendCommand{Name: "command_2"}},
									{Set: &parser.SetCommand{Name: "state_unknown"}},
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
												NilCoalescing: &parser.NilCoalescing{
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
							{Send: &parser.SendCommand{Name: "command_0"}},
							{Send: &parser.SendCommand{Name: "command_1"}},
						},
					},
					{
						Name: "message_1",
						Commands: []*parser.Command{
							{Send: &parser.SendCommand{Name: "command_2"}},
							{Send: &parser.SendCommand{Name: "command_3"}},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages: runtime.MessageGroup{
				"message_0": runtime.NewParameterizedCommandGroup(nil, runtime.CommandGroup{
					commands.NewSendCommand("command_0", nil),
					commands.NewSendCommand("command_1", nil),
				}),
				"message_1": runtime.NewParameterizedCommandGroup(nil, runtime.CommandGroup{
					commands.NewSendCommand("command_2", nil),
					commands.NewSendCommand("command_3", nil),
				}),
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
							{Send: &parser.SendCommand{Name: "command_0"}},
							{Set: &parser.SetCommand{Name: "command_1"}},
						},
					},
					{
						Name: "message_1",
						Commands: []*parser.Command{
							{Send: &parser.SendCommand{Name: "command_2"}},
							{Set: &parser.SetCommand{Name: "command_3"}},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages: runtime.MessageGroup{
				"message_0": runtime.NewParameterizedCommandGroup(nil, runtime.CommandGroup{
					commands.NewSendCommand("command_0", nil),
					commands.NewSetCommand("command_1", nil),
				}),
				"message_1": runtime.NewParameterizedCommandGroup(nil, runtime.CommandGroup{
					commands.NewSendCommand("command_2", nil),
					commands.NewSetCommand("command_3", nil),
				}),
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
							{Send: &parser.SendCommand{Name: "command_1"}},
							{Set: &parser.SetCommand{Name: "command_0"}},
						},
					},
					{
						Name: "message_1",
						Commands: []*parser.Command{
							{Send: &parser.SendCommand{Name: "command_2"}},
							{Set: &parser.SetCommand{Name: "command_0"}},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages: runtime.MessageGroup{
				"message_0": runtime.NewParameterizedCommandGroup(nil, runtime.CommandGroup{
					commands.NewSendCommand("command_1", nil),
					commands.NewSetCommand("command_0", nil),
				}),
				"message_1": runtime.NewParameterizedCommandGroup(nil, runtime.CommandGroup{
					commands.NewSendCommand("command_2", nil),
					commands.NewSetCommand("command_0", nil),
				}),
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
			wantMessages: runtime.MessageGroup{"message_0": {}, "message_1": {}},
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
										NilCoalescing: &parser.NilCoalescing{
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
								},
							},
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages: runtime.MessageGroup{
				"message_0": runtime.NewParameterizedCommandGroup(nil, runtime.CommandGroup{
					commands.NewExpressionCommand(expressions.NewIdentifier("test")),
				}),
			},
			wantSettedStatesByMessages: settedStateGroup{"message_0": mapset.NewSet()},
			wantErr:                    assert.NoError,
		},
		{
			name: "success with parameters",
			args: args{
				messages: []*parser.Message{
					{
						Name:       "message_0",
						Parameters: []string{"one", "two"},
						Commands: []*parser.Command{
							{
								Expression: &parser.Expression{
									ListConstruction: &parser.ListConstruction{
										NilCoalescing: &parser.NilCoalescing{
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
																							Atom: &parser.Atom{Identifier: pointer.ToString("one")},
																						},
																					},
																				},
																			},
																		},
																	},
																},
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
										NilCoalescing: &parser.NilCoalescing{
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
																							Atom: &parser.Atom{Identifier: pointer.ToString("two")},
																						},
																					},
																				},
																			},
																		},
																	},
																},
															},
														},
													},
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
				"message_0": runtime.NewParameterizedCommandGroup(
					[]string{"one", "two"},
					runtime.CommandGroup{
						commands.NewExpressionCommand(expressions.NewIdentifier("one")),
						commands.NewExpressionCommand(expressions.NewIdentifier("two")),
					},
				),
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
							{Send: &parser.SendCommand{Name: "command_0"}},
							{Send: &parser.SendCommand{Name: "command_1"}},
						},
					},
					{
						Name: "message_1",
						Commands: []*parser.Command{
							{Send: &parser.SendCommand{Name: "command_2"}},
							{Set: &parser.SetCommand{Name: "command_3"}},
							{Send: &parser.SendCommand{Name: "command_4"}},
							{Set: &parser.SetCommand{Name: "command_5"}},
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
										NilCoalescing: &parser.NilCoalescing{
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
					{Send: &parser.SendCommand{Name: "one"}},
					{Send: &parser.SendCommand{Name: "two"}},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: runtime.CommandGroup{
				commands.NewSendCommand("one", nil),
				commands.NewSendCommand("two", nil),
			},
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "success with commands (with a set command)",
			args: args{
				commands: []*parser.Command{
					{Send: &parser.SendCommand{Name: "one"}},
					{Set: &parser.SetCommand{Name: "two"}},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: runtime.CommandGroup{
				commands.NewSendCommand("one", nil),
				commands.NewSetCommand("two", nil),
			},
			wantSettedStates: mapset.NewSet("two"),
			wantErr:          assert.NoError,
		},
		{
			name: "success with the return command",
			args: args{
				commands: []*parser.Command{
					{Send: &parser.SendCommand{Name: "one"}},
					{Send: &parser.SendCommand{Name: "two"}},
					{Return: true},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: runtime.CommandGroup{
				commands.NewSendCommand("one", nil),
				commands.NewSendCommand("two", nil),
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
								NilCoalescing: &parser.NilCoalescing{
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
									NilCoalescing: &parser.NilCoalescing{
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
					{
						Expression: &parser.Expression{
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
								NilCoalescing: &parser.NilCoalescing{
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
																											NilCoalescing: &parser.NilCoalescing{
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											NilCoalescing: &parser.NilCoalescing{
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
					{
						Expression: &parser.Expression{
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
																											NilCoalescing: &parser.NilCoalescing{
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											NilCoalescing: &parser.NilCoalescing{
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
			wantCommands: runtime.CommandGroup{
				commands.NewExpressionCommand(
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
				),
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
								NilCoalescing: &parser.NilCoalescing{
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
			wantCommands: nil,
			wantErr:      assert.Error,
		},
		{
			name: "error with the return command",
			args: args{
				commands: []*parser.Command{
					{Send: &parser.SendCommand{Name: "one"}},
					{Return: true},
					{Send: &parser.SendCommand{Name: "two"}},
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
					{Send: &parser.SendCommand{Name: "one"}},
					{Set: &parser.SetCommand{Name: "two"}},
					{Send: &parser.SendCommand{Name: "three"}},
					{Set: &parser.SetCommand{Name: "four"}},
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
								NilCoalescing: &parser.NilCoalescing{
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
								NilCoalescing: &parser.NilCoalescing{
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
								NilCoalescing: &parser.NilCoalescing{
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
																											NilCoalescing: &parser.NilCoalescing{
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											NilCoalescing: &parser.NilCoalescing{
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
			wantDeclaredIdentifiers: mapset.NewSet("test", "test2"),
			wantCommand: commands.NewLetCommand(
				"test2",
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
								NilCoalescing: &parser.NilCoalescing{
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
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             nil,
			wantTopLevelSettedState: "",
			wantReturn:              assert.False,
			wantErr:                 assert.Error,
		},
		{
			name: "Command/start/success",
			args: args{
				command: &parser.Command{
					Start: &parser.StartCommand{
						Expression: &parser.Expression{
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
						},
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             commands.NewStartCommand(expressions.NewIdentifier("test"), nil),
			wantTopLevelSettedState: "",
			wantSettedStates:        mapset.NewSet(),
			wantReturn:              assert.False,
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/start/success/with setted states",
			args: args{
				command: &parser.Command{
					Start: &parser.StartCommand{
						Expression: &parser.Expression{
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
																											NilCoalescing: &parser.NilCoalescing{
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											NilCoalescing: &parser.NilCoalescing{
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
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand: commands.NewStartCommand(
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
				nil,
			),
			wantTopLevelSettedState: "",
			wantSettedStates:        mapset.NewSet("one", "two"),
			wantReturn:              assert.False,
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/start/error",
			args: args{
				command: &parser.Command{
					Start: &parser.StartCommand{
						Expression: &parser.Expression{
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             nil,
			wantTopLevelSettedState: "",
			wantSettedStates:        nil,
			wantReturn:              assert.False,
			wantErr:                 assert.Error,
		},
		{
			name: "Command/send/success",
			args: args{
				command: &parser.Command{
					Send: &parser.SendCommand{
						Name: "test",
						Arguments: []*parser.Expression{
							{
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand: commands.NewSendCommand("test", []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewNumber(23),
				expressions.NewNumber(42),
			}),
			wantTopLevelSettedState: "",
			wantSettedStates:        mapset.NewSet(),
			wantReturn:              assert.False,
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/send/success/with setted states",
			args: args{
				command: &parser.Command{
					Send: &parser.SendCommand{
						Name: "test",
						Arguments: []*parser.Expression{
							{
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
																												NilCoalescing: &parser.NilCoalescing{
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
																										Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																									},
																									{
																										Condition: &parser.Expression{
																											ListConstruction: &parser.ListConstruction{
																												NilCoalescing: &parser.NilCoalescing{
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
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
																												NilCoalescing: &parser.NilCoalescing{
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
																										Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																									},
																									{
																										Condition: &parser.Expression{
																											ListConstruction: &parser.ListConstruction{
																												NilCoalescing: &parser.NilCoalescing{
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
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand: commands.NewSendCommand("test", []expressions.Expression{
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
			wantTopLevelSettedState: "",
			wantSettedStates:        mapset.NewSet("one", "two", "three"),
			wantReturn:              assert.False,
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/send/error",
			args: args{
				command: &parser.Command{
					Send: &parser.SendCommand{
						Name: "test",
						Arguments: []*parser.Expression{
							{
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             nil,
			wantTopLevelSettedState: "",
			wantSettedStates:        nil,
			wantReturn:              assert.False,
			wantErr:                 assert.Error,
		},
		{
			name: "Command/set/success",
			args: args{
				command: &parser.Command{
					Set: &parser.SetCommand{
						Name: "test",
						Arguments: []*parser.Expression{
							{
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand: commands.NewSetCommand("test", []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewNumber(23),
				expressions.NewNumber(42),
			}),
			wantTopLevelSettedState: "test",
			wantSettedStates:        mapset.NewSet("test"),
			wantReturn:              assert.False,
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/set/success/with setted states",
			args: args{
				command: &parser.Command{
					Set: &parser.SetCommand{
						Name: "test",
						Arguments: []*parser.Expression{
							{
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
																												NilCoalescing: &parser.NilCoalescing{
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
																										Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																									},
																									{
																										Condition: &parser.Expression{
																											ListConstruction: &parser.ListConstruction{
																												NilCoalescing: &parser.NilCoalescing{
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
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
																												NilCoalescing: &parser.NilCoalescing{
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
																										Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																									},
																									{
																										Condition: &parser.Expression{
																											ListConstruction: &parser.ListConstruction{
																												NilCoalescing: &parser.NilCoalescing{
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
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand: commands.NewSetCommand("test", []expressions.Expression{
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
			wantTopLevelSettedState: "test",
			wantSettedStates:        mapset.NewSet("one", "two", "three", "test"),
			wantReturn:              assert.False,
			wantErr:                 assert.NoError,
		},
		{
			name: "Command/set/error",
			args: args{
				command: &parser.Command{
					Set: &parser.SetCommand{
						Name: "test",
						Arguments: []*parser.Expression{
							{
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
								ListConstruction: &parser.ListConstruction{
									NilCoalescing: &parser.NilCoalescing{
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
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand:             nil,
			wantTopLevelSettedState: "",
			wantSettedStates:        nil,
			wantReturn:              assert.False,
			wantErr:                 assert.Error,
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
							NilCoalescing: &parser.NilCoalescing{
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
							NilCoalescing: &parser.NilCoalescing{
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
																										NilCoalescing: &parser.NilCoalescing{
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
																								Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																							},
																							{
																								Condition: &parser.Expression{
																									ListConstruction: &parser.ListConstruction{
																										NilCoalescing: &parser.NilCoalescing{
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantDeclaredIdentifiers: mapset.NewSet("test"),
			wantCommand: commands.NewExpressionCommand(
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
							NilCoalescing: &parser.NilCoalescing{
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

func TestTranslateStartCommand(test *testing.T) {
	type args struct {
		startCommand        *parser.StartCommand
		declaredIdentifiers mapset.Set
	}

	for _, testData := range []struct {
		name             string
		args             args
		wantCommand      runtime.Command
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "StartCommand/success/name",
			args: args{
				startCommand:        &parser.StartCommand{Name: pointer.ToString("test")},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand:      commands.NewStartCommand(expressions.NewIdentifier("test"), nil),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "StartCommand/success/name/few arguments",
			args: args{
				startCommand: &parser.StartCommand{
					Name: pointer.ToString("test"),
					Arguments: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
			wantCommand: commands.NewStartCommand(
				expressions.NewIdentifier("test"),
				[]expressions.Expression{
					expressions.NewNumber(12),
					expressions.NewNumber(23),
					expressions.NewNumber(42),
				},
			),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "StartCommand/success/name/few arguments/with setted states",
			args: args{
				startCommand: &parser.StartCommand{
					Name: pointer.ToString("test"),
					Arguments: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
																											NilCoalescing: &parser.NilCoalescing{
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											NilCoalescing: &parser.NilCoalescing{
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
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
																											NilCoalescing: &parser.NilCoalescing{
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											NilCoalescing: &parser.NilCoalescing{
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
			wantCommand: commands.NewStartCommand(
				expressions.NewIdentifier("test"),
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
			name: "StartCommand/success/expression",
			args: args{
				startCommand: &parser.StartCommand{
					Expression: &parser.Expression{
						ListConstruction: &parser.ListConstruction{
							NilCoalescing: &parser.NilCoalescing{
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
					},
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand:      commands.NewStartCommand(expressions.NewIdentifier("test"), nil),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "StartCommand/success/expression/with setted states",
			args: args{
				startCommand: &parser.StartCommand{
					Expression: &parser.Expression{
						ListConstruction: &parser.ListConstruction{
							NilCoalescing: &parser.NilCoalescing{
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
																										NilCoalescing: &parser.NilCoalescing{
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
																								Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																							},
																							{
																								Condition: &parser.Expression{
																									ListConstruction: &parser.ListConstruction{
																										NilCoalescing: &parser.NilCoalescing{
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
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand: commands.NewStartCommand(
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
				nil,
			),
			wantSettedStates: mapset.NewSet("one", "two"),
			wantErr:          assert.NoError,
		},
		{
			name: "StartCommand/error/unknown identifier in the name",
			args: args{
				startCommand:        &parser.StartCommand{Name: pointer.ToString("unknown")},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand:      nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "StartCommand/error/unknown identifier in the expression",
			args: args{
				startCommand: &parser.StartCommand{
					Expression: &parser.Expression{
						ListConstruction: &parser.ListConstruction{
							NilCoalescing: &parser.NilCoalescing{
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
			wantCommand:      nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "StartCommand/error/unknown identifier in the argument",
			args: args{
				startCommand: &parser.StartCommand{
					Name: pointer.ToString("test"),
					Arguments: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
			wantCommand:      nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			gotCommand, gotSettedStates, err :=
				translateStartCommand(testData.args.startCommand, testData.args.declaredIdentifiers)

			assert.Equal(test, testData.wantCommand, gotCommand)
			assert.Equal(test, testData.wantSettedStates, gotSettedStates)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateSendCommand(test *testing.T) {
	type args struct {
		sendCommand         *parser.SendCommand
		declaredIdentifiers mapset.Set
	}

	for _, testData := range []struct {
		name             string
		args             args
		wantCommand      runtime.Command
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "SendCommand/success/few arguments",
			args: args{
				sendCommand: &parser.SendCommand{
					Name: "test",
					Arguments: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
			wantCommand: commands.NewSendCommand("test", []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewNumber(23),
				expressions.NewNumber(42),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "SendCommand/success/few arguments/with setted states",
			args: args{
				sendCommand: &parser.SendCommand{
					Name: "test",
					Arguments: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
																											NilCoalescing: &parser.NilCoalescing{
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "one"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											NilCoalescing: &parser.NilCoalescing{
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
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
																											NilCoalescing: &parser.NilCoalescing{
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
																									Commands: []*parser.Command{{Set: &parser.SetCommand{Name: "two"}}},
																								},
																								{
																									Condition: &parser.Expression{
																										ListConstruction: &parser.ListConstruction{
																											NilCoalescing: &parser.NilCoalescing{
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
			wantCommand: commands.NewSendCommand("test", []expressions.Expression{
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
			name: "SendCommand/success/no arguments",
			args: args{
				sendCommand: &parser.SendCommand{
					Name:      "test",
					Arguments: nil,
				},
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand:      commands.NewSendCommand("test", nil),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "SendCommand/error",
			args: args{
				sendCommand: &parser.SendCommand{
					Name: "test",
					Arguments: []*parser.Expression{
						{
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
							ListConstruction: &parser.ListConstruction{
								NilCoalescing: &parser.NilCoalescing{
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
			wantCommand: nil,
			wantErr:     assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			gotCommand, gotSettedStates, err :=
				translateSendCommand(testData.args.sendCommand, testData.args.declaredIdentifiers)

			assert.Equal(test, testData.wantCommand, gotCommand)
			assert.Equal(test, testData.wantSettedStates, gotSettedStates)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateSetCommand(test *testing.T) {
	type args struct {
		code                string
		declaredIdentifiers mapset.Set
	}

	for _, testData := range []struct {
		name             string
		args             args
		wantCommand      runtime.Command
		wantSettedStates mapset.Set
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "SetCommand/success/few arguments",
			args: args{
				code:                "set test(12, 23, 42)",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand: commands.NewSetCommand("test", []expressions.Expression{
				expressions.NewNumber(12),
				expressions.NewNumber(23),
				expressions.NewNumber(42),
			}),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "SetCommand/success/few arguments/with setted states",
			args: args{
				code: `set test(
					when
						=> 23
							set one()
						=> 42
							set two()
					;,
					when
						=> 24
							set two()
						=> 43
							set three()
					;,
				)`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand: commands.NewSetCommand("test", []expressions.Expression{
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
			name: "SetCommand/success/no arguments",
			args: args{
				code:                "set test()",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand:      commands.NewSetCommand("test", nil),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "SetCommand/error",
			args: args{
				code:                "set test(12, 23, unknown)",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand: nil,
			wantErr:     assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			setCommand := new(parser.SetCommand)
			err := parser.ParseToAST(testData.args.code, setCommand)
			require.NoError(test, err)

			gotCommand, gotSettedStates, err :=
				translateSetCommand(setCommand, testData.args.declaredIdentifiers)

			assert.Equal(test, testData.wantCommand, gotCommand)
			assert.Equal(test, testData.wantSettedStates, gotSettedStates)
			testData.wantErr(test, err)
		})
	}
}
