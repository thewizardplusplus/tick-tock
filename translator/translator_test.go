package translator

import (
	"io"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/commands"
	runtimemocks "github.com/thewizardplusplus/tick-tock/runtime/mocks"
	"github.com/thewizardplusplus/tick-tock/tests"
	testsmocks "github.com/thewizardplusplus/tick-tock/tests/mocks"
)

func TestTranslate(test *testing.T) {
	type args struct {
		actors []*parser.Actor
	}

	for _, testData := range []struct {
		name     string
		args     args
		makeWant func(dependencies Dependencies) runtime.ConcurrentActorGroup
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success with actors",
			args: args{
				actors: []*parser.Actor{
					{[]*parser.State{{false, "one", nil}}},
					{[]*parser.State{{false, "two", nil}}},
				},
			},
			makeWant: func(dependencies Dependencies) runtime.ConcurrentActorGroup {
				actorOne, _ := runtime.NewActor(runtime.StateGroup{"one": runtime.MessageGroup{}}, "one")
				actorTwo, _ := runtime.NewActor(runtime.StateGroup{"two": runtime.MessageGroup{}}, "two")
				return runtime.ConcurrentActorGroup{
					runtime.NewConcurrentActor(tests.BufferedInbox, actorOne, dependencies.Dependencies),
					runtime.NewConcurrentActor(tests.BufferedInbox, actorTwo, dependencies.Dependencies),
				}
			},
			wantErr: assert.NoError,
		},
		{
			name:     "success without actors",
			args:     args{},
			makeWant: func(dependencies Dependencies) runtime.ConcurrentActorGroup { return nil },
			wantErr:  assert.NoError,
		},
		{
			name:     "error",
			args:     args{[]*parser.Actor{{[]*parser.State{{false, "test", nil}}}, {}}},
			makeWant: func(dependencies Dependencies) runtime.ConcurrentActorGroup { return nil },
			wantErr:  assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			waiter := new(runtimemocks.Waiter)
			errorHandler := new(runtimemocks.ErrorHandler)
			outWriter := new(testsmocks.Writer)
			dependencies := Dependencies{
				Dependencies: runtime.Dependencies{Waiter: waiter, ErrorHandler: errorHandler},
				OutWriter:    outWriter,
			}
			got, err := Translate(tests.BufferedInbox, testData.args.actors, dependencies)

			mock.AssertExpectationsForObjects(test, waiter, errorHandler, outWriter)
			assert.Equal(test, cleanInboxes(testData.makeWant(dependencies)), cleanInboxes(got))
			testData.wantErr(test, err)
		})
	}
}

func TestTranslateStates(test *testing.T) {
	type args struct {
		states []*parser.State
	}

	for _, testData := range []struct {
		name             string
		args             args
		makeWantStates   func(outWriter io.Writer) runtime.StateGroup
		wantInitialState string
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "success with nonempty states",
			args: args{
				states: []*parser.State{
					{false, "state_0", []*parser.Message{{"message_0", nil}, {"message_1", nil}}},
					{false, "state_1", []*parser.Message{{"message_2", nil}, {"message_3", nil}}},
				},
			},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup {
				return runtime.StateGroup{
					"state_0": runtime.MessageGroup{"message_0": nil, "message_1": nil},
					"state_1": runtime.MessageGroup{"message_2": nil, "message_3": nil},
				}
			},
			wantInitialState: "state_0",
			wantErr:          assert.NoError,
		},
		{
			name: "success with empty states (with an implicit initial state)",
			args: args{[]*parser.State{{false, "state_0", nil}, {false, "state_1", nil}}},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup {
				return runtime.StateGroup{"state_0": runtime.MessageGroup{}, "state_1": runtime.MessageGroup{}}
			},
			wantInitialState: "state_0",
			wantErr:          assert.NoError,
		},
		{
			name: "success with empty states (with an explicit initial state)",
			args: args{[]*parser.State{{false, "state_0", nil}, {true, "state_1", nil}}},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup {
				return runtime.StateGroup{"state_0": runtime.MessageGroup{}, "state_1": runtime.MessageGroup{}}
			},
			wantInitialState: "state_1",
			wantErr:          assert.NoError,
		},
		{
			name:           "error without states",
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup { return nil },
			wantErr:        assert.Error,
		},
		{
			name:           "error with duplicate states",
			args:           args{[]*parser.State{{false, "test", nil}, {false, "test", nil}}},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup { return nil },
			wantErr:        assert.Error,
		},
		{
			name: "error with few initial states",
			args: args{
				states: []*parser.State{
					{false, "state_0", nil},
					{true, "state_1", nil},
					{false, "state_2", nil},
					{true, "state_3", nil},
				},
			},
			makeWantStates: func(outWriter io.Writer) runtime.StateGroup { return nil },
			wantErr:        assert.Error,
		},
		{
			name: "error with messages translation",
			args: args{
				states: []*parser.State{
					{false, "state_0", []*parser.Message{{"message_0", nil}, {"message_1", nil}}},
					{false, "state_1", []*parser.Message{{"test", nil}, {"test", nil}}},
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
						Initial: false,
						Name:    "state_0",
						Messages: []*parser.Message{
							{
								Name: "message_0",
								Commands: []*parser.Command{
									{Send: tests.GetAddress("command_0")},
									{Set: tests.GetAddress("state_unknown")},
								},
							},
							{
								Name: "message_1",
								Commands: []*parser.Command{
									{Send: tests.GetAddress("command_2")},
									{Set: tests.GetAddress("state_unknown")},
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
			gotStates, gotInitialState, err := translateStates(testData.args.states, outWriter)

			mock.AssertExpectationsForObjects(test, outWriter)
			assert.Equal(test, testData.makeWantStates(outWriter), gotStates)
			assert.Equal(test, testData.wantInitialState, gotInitialState)
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
							{Send: tests.GetAddress("command_0")},
							{Send: tests.GetAddress("command_1")},
						},
					},
					{
						Name: "message_1",
						Commands: []*parser.Command{
							{Send: tests.GetAddress("command_2")},
							{Send: tests.GetAddress("command_3")},
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
							{Send: tests.GetAddress("command_0")},
							{Set: tests.GetAddress("command_1")},
						},
					},
					{
						Name: "message_1",
						Commands: []*parser.Command{
							{Send: tests.GetAddress("command_2")},
							{Set: tests.GetAddress("command_3")},
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
			args: args{[]*parser.Message{{"message_0", nil}, {"message_1", nil}}},
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
			args:             args{[]*parser.Message{{"test", nil}, {"test", nil}}},
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
							{Send: tests.GetAddress("command_0")},
							{Send: tests.GetAddress("command_1")},
						},
					},
					{
						Name: "message_1",
						Commands: []*parser.Command{
							{Send: tests.GetAddress("command_2")},
							{Set: tests.GetAddress("command_3")},
							{Send: tests.GetAddress("command_4")},
							{Set: tests.GetAddress("command_5")},
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
			gotMessages, gotStates, err := translateMessages(testData.args.messages, outWriter)

			mock.AssertExpectationsForObjects(test, outWriter)
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
			args: args{[]*parser.Command{{Send: tests.GetAddress("one")}, {Send: tests.GetAddress("two")}}},
			makeWantCommands: func(outWriter io.Writer) runtime.CommandGroup {
				return runtime.CommandGroup{commands.NewSendCommand("one"), commands.NewSendCommand("two")}
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with commands (with a set command)",
			args: args{[]*parser.Command{{Send: tests.GetAddress("one")}, {Set: tests.GetAddress("two")}}},
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
			name: "error",
			args: args{
				commands: []*parser.Command{
					{Send: tests.GetAddress("one")},
					{Set: tests.GetAddress("two")},
					{Send: tests.GetAddress("three")},
					{Set: tests.GetAddress("four")},
				},
			},
			makeWantCommands: func(outWriter io.Writer) runtime.CommandGroup { return nil },
			wantErr:          assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			outWriter := new(testsmocks.Writer)
			gotCommands, gotState, err := translateCommands(testData.args.commands, outWriter)

			mock.AssertExpectationsForObjects(test, outWriter)
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
	}{
		{
			name: "Command/send",
			args: args{&parser.Command{Send: tests.GetAddress("test")}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewSendCommand("test")
			},
		},
		{
			name: "Command/set",
			args: args{&parser.Command{Set: tests.GetAddress("test")}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewSetCommand("test")
			},
			wantState: "test",
		},
		{
			name: "Command/out/nonempty",
			args: args{&parser.Command{Out: tests.GetAddress("test")}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewOutCommand("test", outWriter)
			},
		},
		{
			name: "Command/out/empty",
			args: args{&parser.Command{Out: tests.GetAddress("")}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command {
				return commands.NewOutCommand("", outWriter)
			},
		},
		{
			name:            "Command/exit",
			args:            args{&parser.Command{Exit: true}},
			makeWantCommand: func(outWriter io.Writer) runtime.Command { return commands.ExitCommand{} },
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			outWriter := new(testsmocks.Writer)
			gotCommand, gotState := translateCommand(testData.args.command, outWriter)

			mock.AssertExpectationsForObjects(test, outWriter)
			assert.Equal(test, testData.makeWantCommand(outWriter), gotCommand)
			assert.Equal(test, testData.wantState, gotState)
		})
	}
}

func cleanInboxes(actors runtime.ConcurrentActorGroup) runtime.ConcurrentActorGroup {
	actorsReflection := reflect.ValueOf(actors)
	for index := range actors {
		field := actorsReflection.Index(index).FieldByName("inbox")
		*(*chan string)(unsafe.Pointer(field.UnsafeAddr())) = nil
	}

	return actors
}
