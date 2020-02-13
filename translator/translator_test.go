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
	runtimemocks "github.com/thewizardplusplus/tick-tock/runtime/mocks"
	waitermocks "github.com/thewizardplusplus/tick-tock/runtime/waiter/mocks"
)

func TestTranslate(test *testing.T) {
	for _, testData := range []struct {
		name       string
		makeActors func(options Options) []*parser.Actor
		makeWant   func(options Options, dependencies Dependencies) runtime.ConcurrentActorGroup
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success with actors",
			makeActors: func(options Options) []*parser.Actor {
				return []*parser.Actor{
					{States: []*parser.State{{Name: options.InitialState}, {Name: "one"}}},
					{States: []*parser.State{{Name: options.InitialState}, {Name: "two"}}},
				}
			},
			makeWant: func(options Options, dependencies Dependencies) runtime.ConcurrentActorGroup {
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
			name:       "success without actors",
			makeActors: func(options Options) []*parser.Actor { return nil },
			makeWant: func(options Options, dependencies Dependencies) runtime.ConcurrentActorGroup {
				return nil
			},
			wantErr: assert.NoError,
		},
		{
			name: "error with states translation",
			makeActors: func(options Options) []*parser.Actor {
				return []*parser.Actor{{States: []*parser.State{{Name: options.InitialState}}}, {}}
			},
			makeWant: func(options Options, dependencies Dependencies) runtime.ConcurrentActorGroup {
				return nil
			},
			wantErr: assert.Error,
		},
		{
			name: "error with actor construction",
			makeActors: func(options Options) []*parser.Actor {
				return []*parser.Actor{
					{States: []*parser.State{{Name: "one"}}},
					{States: []*parser.State{{Name: "two"}}},
				}
			},
			makeWant: func(options Options, dependencies Dependencies) runtime.ConcurrentActorGroup {
				return nil
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
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
			want := testData.makeWant(options, dependencies)
			got, err := Translate(testData.makeActors(options), options, dependencies)

			mock.AssertExpectationsForObjects(test, outWriter, randomizer, sleeper, waiter, errorHandler)
			assert.Equal(test, cleanInboxes(want), cleanInboxes(got))
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateStates(test *testing.T) {
	type args struct {
		states []*parser.State
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
			args: args{[]*parser.State{{Name: "state_0"}, {Name: "state_1"}}},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup {
				return runtime.StateGroup{"state_0": runtime.MessageGroup{}, "state_1": runtime.MessageGroup{}}
			},
			wantErr: assert.NoError,
		},
		{
			name:           "error without states",
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup { return nil },
			wantErr:        assert.Error,
		},
		{
			name:           "error with duplicate states",
			args:           args{[]*parser.State{{Name: "test"}, {Name: "test"}}},
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
			},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup { return nil },
			wantErr:        assert.Error,
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
			gotStates, err := translateStates(testData.args.states, dependencies)

			mock.AssertExpectationsForObjects(test, outWriter, randomizer, sleeper)
			assert.Equal(test, testData.makeWantStates(outWriter), gotStates)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateMessages(test *testing.T) {
	type args struct {
		messages []*parser.Message
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
			args: args{[]*parser.Message{{Name: "message_0"}, {Name: "message_1"}}},
			makeWantMessages: func(outWriter io.Writer) runtime.MessageGroup {
				return runtime.MessageGroup{"message_0": nil, "message_1": nil}
			},
			wantStates: make(settedStateGroup),
			wantErr:    assert.NoError,
		},
		{
			name: "success without messages",
			makeWantMessages: func(outWriter io.Writer) runtime.MessageGroup {
				return runtime.MessageGroup{}
			},
			wantStates: make(settedStateGroup),
			wantErr:    assert.NoError,
		},
		{
			name:             "error with duplicate messages",
			args:             args{[]*parser.Message{{Name: "test"}, {Name: "test"}}},
			makeWantMessages: func(outWriter io.Writer) runtime.MessageGroup { return nil },
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
			},
			makeWantMessages: func(outWriter io.Writer) runtime.MessageGroup { return nil },
			wantErr:          assert.Error,
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
			gotMessages, gotStates, err := translateMessages(testData.args.messages, dependencies)

			mock.AssertExpectationsForObjects(test, outWriter, randomizer, sleeper)
			assert.Equal(test, testData.makeWantMessages(outWriter), gotMessages)
			assert.Equal(test, testData.wantStates, gotStates)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateCommands(test *testing.T) {
	type args struct {
		commands []*parser.Command
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
			},
			makeWantCommands: func(outWriter io.Writer) runtime.CommandGroup {
				return runtime.CommandGroup{commands.NewSendCommand("one"), commands.NewSetCommand("two")}
			},
			wantState: "two",
			wantErr:   assert.NoError,
		},
		{
			name:             "success without commands",
			makeWantCommands: func(outWriter io.Writer) runtime.CommandGroup { return nil },
			wantErr:          assert.NoError,
		},
		{
			name: "error with command translation",
			args: args{
				commands: []*parser.Command{
					{
						Sleep: &parser.SleepCommand{
							Minimum: tests.GetNumberAddress(3.4),
							Maximum: tests.GetNumberAddress(1.2),
						},
					},
				},
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
			},
			makeWantCommands: func(outWriter io.Writer) runtime.CommandGroup { return nil },
			wantErr:          assert.Error,
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
			gotCommands, gotState, err := translateCommands(testData.args.commands, dependencies)

			mock.AssertExpectationsForObjects(test, outWriter, randomizer, sleeper)
			assert.Equal(test, testData.makeWantCommands(outWriter), gotCommands)
			assert.Equal(test, testData.wantState, gotState)
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateCommand(test *testing.T) {
	type args struct {
		command *parser.Command
	}

	for _, testData := range []struct {
		name            string
		args            args
		makeWantCommand func(outWriter io.Writer) runtime.Command
		wantState       string
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "Command/send",
			args: args{&parser.Command{Send: tests.GetStringAddress("test")}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewSendCommand("test")
			},
			wantErr: assert.NoError,
		},
		{
			name: "Command/set",
			args: args{&parser.Command{Set: tests.GetStringAddress("test")}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewSetCommand("test")
			},
			wantState: "test",
			wantErr:   assert.NoError,
		},
		{
			name: "Command/out/nonempty",
			args: args{&parser.Command{Out: tests.GetStringAddress("test")}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewOutCommand("test", outWriter)
			},
			wantErr: assert.NoError,
		},
		{
			name: "Command/out/empty",
			args: args{&parser.Command{Out: tests.GetStringAddress("")}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewOutCommand("", outWriter)
			},
			wantErr: assert.NoError,
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
			},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				command, _ := commands.NewSleepCommand(1.2, 3.4, commands.SleepDependencies{})
				return command
			},
			wantErr: assert.NoError,
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
			},
			makeWantCommand: func(outWriter io.Writer) runtime.Command { return nil },
			wantErr:         assert.Error,
		},
		{
			name:            "Command/exit",
			args:            args{&parser.Command{Exit: true}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command { return commands.ExitCommand{} },
			wantErr:         assert.NoError,
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
			gotCommand, gotState, err := translateCommand(testData.args.command, dependencies)
			if sleepCommand, ok := gotCommand.(commands.SleepCommand); ok {
				cleanSleepDependencies(&sleepCommand)
				gotCommand = sleepCommand
			}

			mock.AssertExpectationsForObjects(test, outWriter, randomizer, sleeper)
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
