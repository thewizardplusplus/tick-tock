package translator

import (
	"testing"

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

func TestTranslateProgram(test *testing.T) {
	type args struct {
		code                string
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
				code:                "",
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
				code: `
					actor Test0()
						state state_0();
						state state_1();
					;
					actor Test1()
						state state_0();
						state state_1();
					;
				`,
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
				code: `
					class Test0()
						state state_0();
						state state_1();
					;
					class Test1()
						state state_0();
						state state_1();
					;
				`,
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
				code: `
					actor Test0()
						state state_0();
						state state_1();
					;
					class Test1()
						state state_0();
						state state_1();
					;
					actor Test2()
						state state_0();
						state state_1();
					;
					class Test3()
						state state_0();
						state state_1();
					;
				`,
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
				code: `
					actor Test0()
						state state_0();
						state state_1();
					;
					actor Test1()
						state state_0();
						state state_0();
					;
				`,
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
				code: `
					actor Test0()
						state state_0();
						state state_1();
					;
					actor Test0()
						state state_0();
						state state_1();
					;
				`,
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
				code: `
					class Test0()
						state state_0();
						state state_1();
					;
					class Test0()
						state state_0();
						state state_1();
					;
				`,
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
				code: `
					actor Test0()
						state state_0();
						state state_1();
					;
					class Test0()
						state state_0();
						state state_1();
					;
				`,
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

			program := new(parser.Program)
			err := parser.ParseToAST(testData.args.code, program)
			require.NoError(test, err)

			gotDefinitions, gotTranslatedActors, err := TranslateProgram(
				program,
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
		code                string
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
				code: `actor Test()
					state state_0();
					state state_1();
				;`,
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
				code: `actor Test()
					state state_0();
					state state_1();
				;`,
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
				code: `class Test()
					state state_0();
					state state_1();
				;`,
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
				code: `class Test()
					state state_0();
					state state_1();
				;`,
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
			definition := new(parser.Definition)
			err := parser.ParseToAST(testData.args.code, definition)
			require.NoError(test, err)

			gotTranslatedActorClass, gotActor, err := translateDefinition(
				definition,
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
		code                string
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
				code: `class Test()
					state state_0();
					state state_1();
				;`,
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
				code: `class Test()
					state state_0()
						message message_0()
							test
						;
					;
				;`,
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
				code: `class Test(one, two)
					state state_0()
						message message_0()
							one
							two
						;
					;
				;`,
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
				code: `class Test()
					state state_0();
					state state_0();
				;`,
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
				code: `class Test()
					state state_0();
					state state_1();
				;`,
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
				code: `class Test()
					state state_0()
						message message_0()
							unknown
						;
					;
				;`,
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

			actorClass := new(parser.ActorClass)
			err := parser.ParseToAST(testData.args.code, actorClass)
			require.NoError(test, err)

			gotTranslatedActorClass, err := translateActorClass(
				actorClass,
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
	type statesWrapper struct {
		States []*parser.State `parser:"{ @@ }"`
	}
	type args struct {
		code                string
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
				code: `
					state state_0()
						message message_0();
						message message_1();
					;
					state state_1()
						message message_2();
						message message_3();
					;
				`,
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
				code:                "state state_0(); state state_1();",
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
				code: `state state_0()
					message message_0()
						test
					;
				;`,
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
				code: `state state_0(one, two)
					message message_0()
						one
						two
					;
				;`,
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
				code:                "",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantStates: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with duplicate states",
			args: args{
				code:                "state test(); state test();",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantStates: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with messages translation",
			args: args{
				code: `
					state state_0()
						message message_0();
						message message_1();
					;
					state state_1()
						message test();
						message test();
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantStates: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with an unknown state",
			args: args{
				code: `state state_0()
					message message_0()
						send command_0()
						set state_unknown()
					;
					message message_1()
						send command_2()
						set state_unknown()
					;
				;`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantStates: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with the expression",
			args: args{
				code: `state state_0()
					message message_0()
						unknown
					;
				;`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantStates: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := testData.args.declaredIdentifiers.Clone()

			statesWrapper := new(statesWrapper)
			err := parser.ParseToAST(testData.args.code, statesWrapper)
			require.NoError(test, err)

			gotStates, err := translateStates(statesWrapper.States, testData.args.declaredIdentifiers)

			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.wantStates, gotStates)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateMessages(test *testing.T) {
	type messagesWrapper struct {
		Messages []*parser.Message `parser:"{ @@ }"`
	}
	type args struct {
		code                string
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
				code: `
					message message_0()
						send command_0()
						send command_1()
					;
					message message_1()
						send command_2()
						send command_3()
					;
				`,
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
				code: `
					message message_0()
						send command_0()
						set command_1()
					;
					message message_1()
						send command_2()
						set command_3()
					;
				`,
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
				code: `
					message message_0()
						send command_1()
						set command_0()
					;
					message message_1()
						send command_2()
						set command_0()
					;
				`,
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
				code:                "message message_0(); message message_1();",
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
				code:                "",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages:               runtime.MessageGroup{},
			wantSettedStatesByMessages: make(settedStateGroup),
			wantErr:                    assert.NoError,
		},
		{
			name: "success with the expression",
			args: args{
				code:                "message message_0() test;",
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
				code:                "message message_0(one, two) one two;",
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
				code:                "message test(); message test();",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages:               nil,
			wantSettedStatesByMessages: nil,
			wantErr:                    assert.Error,
		},
		{
			name: "error with commands translation",
			args: args{
				code: `
					message message_0()
						send command_0()
						send command_1()
					;
					message message_1()
						send command_2()
						set command_3()
						send command_4()
						set command_5()
					;
				`,
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages:               nil,
			wantSettedStatesByMessages: nil,
			wantErr:                    assert.Error,
		},
		{
			name: "error with the expression",
			args: args{
				code:                "message message_0() unknown;",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantMessages:               nil,
			wantSettedStatesByMessages: nil,
			wantErr:                    assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := testData.args.declaredIdentifiers.Clone()

			messagesWrapper := new(messagesWrapper)
			err := parser.ParseToAST(testData.args.code, messagesWrapper)
			require.NoError(test, err)

			gotMessages, gotSettedStatesByMessages, err :=
				translateMessages(messagesWrapper.Messages, testData.args.declaredIdentifiers)

			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.wantMessages, gotMessages)
			assert.Equal(test, testData.wantSettedStatesByMessages, gotSettedStatesByMessages)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateCommands(test *testing.T) {
	type commandsWrapper struct {
		Commands []*parser.Command `parser:"{ @@ }"`
	}
	type args struct {
		code                string
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
				code:                "send one() send two()",
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
				code:                "send one() set two()",
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
				code:                "send one() send two() return",
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
				code:                "test",
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
				code:                "let test2 = 23 test2",
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
				code:                "",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands:     nil,
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "success with commands (with setted states)",
			args: args{
				code: `when
					=> 23
						set one()
					=> 42
						set two()
				;
				when
					=> 24
						set two()
					=> 43
						set three()
				;`,
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
				code:                "unknown",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: nil,
			wantErr:      assert.Error,
		},
		{
			name: "error with the return command",
			args: args{
				code:                "send one() return send two()",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: nil,
			wantErr:      assert.Error,
		},
		{
			name: "error with a second set command",
			args: args{
				code:                "send one() set two() send three() set four()",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommands: nil,
			wantErr:      assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			originDeclaredIdentifiers := testData.args.declaredIdentifiers.Clone()

			commandsWrapper := new(commandsWrapper)
			err := parser.ParseToAST(testData.args.code, commandsWrapper)
			require.NoError(test, err)

			gotCommands, gotSettedStates, err :=
				translateCommands(commandsWrapper.Commands, testData.args.declaredIdentifiers)

			assert.Equal(test, originDeclaredIdentifiers, testData.args.declaredIdentifiers)
			assert.Equal(test, testData.wantCommands, gotCommands)
			assert.Equal(test, testData.wantSettedStates, gotSettedStates)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateCommand(test *testing.T) {
	type args struct {
		code                string
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
				code:                "let test2 = 23",
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
				code:                "let test = 23",
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
				code: `let test2 = when
					=> 23
						set one()
					=> 42
						set two()
				;`,
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
				code:                "let test2 = unknown",
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
				code:                "start test()",
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
				code: `start [
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				]()`,
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
				code:                "start unknown()",
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
				code:                "send test(12, 23, 42)",
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
				code: `send test(
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
				code:                "send test(12, 23, unknown)",
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
				code:                "set test(12, 23, 42)",
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
				code:                "set test(12, 23, unknown)",
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
				code:                "return",
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
				code:                "test",
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
				code: `when
					=> 23
						set one()
					=> 42
						set two()
				;`,
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
				code:                "unknown",
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
			command := new(parser.Command)
			err := parser.ParseToAST(testData.args.code, command)
			require.NoError(test, err)

			gotCommand, gotTopLevelSettedState, gotSettedStates, gotReturn, err :=
				translateCommand(command, testData.args.declaredIdentifiers)

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
			name: "StartCommand/success/name",
			args: args{
				code:                "start test()",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand:      commands.NewStartCommand(expressions.NewIdentifier("test"), nil),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "StartCommand/success/name/few arguments",
			args: args{
				code:                "start test(12, 23, 42)",
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
				code: `start test(
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
				code:                "start [test()]()",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand:      commands.NewStartCommand(expressions.NewFunctionCall("test", nil), nil),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "StartCommand/success/expression/with setted states",
			args: args{
				code: `start [
					when
						=> 23
							set one()
						=> 42
							set two()
					;
				]()`,
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
				code:                "start unknown()",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand:      nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "StartCommand/error/unknown identifier in the expression",
			args: args{
				code:                "start [unknown()]()",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand:      nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
		{
			name: "StartCommand/error/unknown identifier in the argument",
			args: args{
				code:                "start test(12, 23, unknown)",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand:      nil,
			wantSettedStates: nil,
			wantErr:          assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			startCommand := new(parser.StartCommand)
			err := parser.ParseToAST(testData.args.code, startCommand)
			require.NoError(test, err)

			gotCommand, gotSettedStates, err :=
				translateStartCommand(startCommand, testData.args.declaredIdentifiers)

			assert.Equal(test, testData.wantCommand, gotCommand)
			assert.Equal(test, testData.wantSettedStates, gotSettedStates)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateSendCommand(test *testing.T) {
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
			name: "SendCommand/success/few arguments",
			args: args{
				code:                "send test(12, 23, 42)",
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
				code: `send test(
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
				code:                "send test()",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand:      commands.NewSendCommand("test", nil),
			wantSettedStates: mapset.NewSet(),
			wantErr:          assert.NoError,
		},
		{
			name: "SendCommand/error",
			args: args{
				code:                "send test(12, 23, unknown)",
				declaredIdentifiers: mapset.NewSet("test"),
			},
			wantCommand: nil,
			wantErr:     assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			sendCommand := new(parser.SendCommand)
			err := parser.ParseToAST(testData.args.code, sendCommand)
			require.NoError(test, err)

			gotCommand, gotSettedStates, err :=
				translateSendCommand(sendCommand, testData.args.declaredIdentifiers)

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
