package translator

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/commands"
	"github.com/thewizardplusplus/tick-tock/tests"
	"github.com/thewizardplusplus/tick-tock/tests/mocks"
)

func TestTranslateStates(test *testing.T) {
	type args struct {
		states []*parser.State
	}

	for _, testData := range []struct {
		name     string
		args     args
		makeWant func(writer io.Writer) runtime.StateGroup
	}{
		{
			name: "success with nonempty states",
			args: args{
				states: []*parser.State{
					{false, "state_0", []*parser.Message{{"message_0", nil}, {"message_1", nil}}},
					{false, "state_1", []*parser.Message{{"message_2", nil}, {"message_3", nil}}},
				},
			},
			makeWant: func(writer io.Writer) runtime.StateGroup {
				return runtime.StateGroup{
					"state_0": runtime.MessageGroup{"message_0": nil, "message_1": nil},
					"state_1": runtime.MessageGroup{"message_2": nil, "message_3": nil},
				}
			},
		},
		{
			name: "success with empty states",
			args: args{[]*parser.State{{false, "state_0", nil}, {false, "state_1", nil}}},
			makeWant: func(writer io.Writer) runtime.StateGroup {
				return runtime.StateGroup{"state_0": runtime.MessageGroup{}, "state_1": runtime.MessageGroup{}}
			},
		},
		{
			name:     "success without states",
			makeWant: func(writer io.Writer) runtime.StateGroup { return runtime.StateGroup{} },
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			writer := new(mocks.Writer)
			want := testData.makeWant(writer)
			got := TranslateStates(writer, testData.args.states)

			writer.AssertExpectations(test)
			assert.Equal(test, want, got)
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
		makeWantMessages func(writer io.Writer) runtime.MessageGroup
		wantStates       SettedStateGroup
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
			makeWantMessages: func(writer io.Writer) runtime.MessageGroup {
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
			wantStates: make(SettedStateGroup),
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
			makeWantMessages: func(writer io.Writer) runtime.MessageGroup {
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
			wantStates: SettedStateGroup{"message_0": "command_1", "message_1": "command_3"},
			wantErr:    assert.NoError,
		},
		{
			name: "success with empty messages",
			args: args{[]*parser.Message{{"message_0", nil}, {"message_1", nil}}},
			makeWantMessages: func(writer io.Writer) runtime.MessageGroup {
				return runtime.MessageGroup{"message_0": nil, "message_1": nil}
			},
			wantStates: make(SettedStateGroup),
			wantErr:    assert.NoError,
		},
		{
			name:             "success without messages",
			makeWantMessages: func(writer io.Writer) runtime.MessageGroup { return runtime.MessageGroup{} },
			wantStates:       make(SettedStateGroup),
			wantErr:          assert.NoError,
		},
		{
			name:             "error with duplicate messages",
			args:             args{[]*parser.Message{{"test", nil}, {"test", nil}}},
			makeWantMessages: func(writer io.Writer) runtime.MessageGroup { return nil },
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
			makeWantMessages: func(writer io.Writer) runtime.MessageGroup { return nil },
			wantErr:          assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			writer := new(mocks.Writer)
			wantMessages := testData.makeWantMessages(writer)
			gotMessages, gotStates, err := TranslateMessages(writer, testData.args.messages)

			writer.AssertExpectations(test)
			assert.Equal(test, wantMessages, gotMessages)
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
		makeWantCommands func(writer io.Writer) runtime.CommandGroup
		wantState        string
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "success with commands (without a set command)",
			args: args{[]*parser.Command{{Send: tests.GetAddress("one")}, {Send: tests.GetAddress("two")}}},
			makeWantCommands: func(writer io.Writer) runtime.CommandGroup {
				return runtime.CommandGroup{commands.NewSendCommand("one"), commands.NewSendCommand("two")}
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with commands (with a set command)",
			args: args{[]*parser.Command{{Send: tests.GetAddress("one")}, {Set: tests.GetAddress("two")}}},
			makeWantCommands: func(writer io.Writer) runtime.CommandGroup {
				return runtime.CommandGroup{commands.NewSendCommand("one"), commands.NewSetCommand("two")}
			},
			wantState: "two",
			wantErr:   assert.NoError,
		},
		{
			name:             "success without commands",
			makeWantCommands: func(writer io.Writer) runtime.CommandGroup { return nil },
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
			makeWantCommands: func(writer io.Writer) runtime.CommandGroup { return nil },
			wantErr:          assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			writer := new(mocks.Writer)
			wantCommands := testData.makeWantCommands(writer)
			gotCommands, gotState, err := TranslateCommands(writer, testData.args.commands)

			writer.AssertExpectations(test)
			assert.Equal(test, wantCommands, gotCommands)
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
		makeWantCommand func(writer io.Writer) runtime.Command
		wantState       string
	}{
		{
			name: "Command/send",
			args: args{&parser.Command{Send: tests.GetAddress("test")}},
			makeWantCommand: func(writer io.Writer) runtime.Command {
				return commands.NewSendCommand("test")
			},
		},
		{
			name: "Command/set",
			args: args{&parser.Command{Set: tests.GetAddress("test")}},
			makeWantCommand: func(writer io.Writer) runtime.Command {
				return commands.NewSetCommand("test")
			},
			wantState: "test",
		},
		{
			name: "Command/out/nonempty",
			args: args{&parser.Command{Out: tests.GetAddress("test")}},
			makeWantCommand: func(writer io.Writer) runtime.Command {
				return commands.NewOutCommand(writer, "test")
			},
		},
		{
			name: "Command/out/empty",
			args: args{&parser.Command{Out: tests.GetAddress("")}},
			makeWantCommand: func(writer io.Writer) runtime.Command {
				return commands.NewOutCommand(writer, "")
			},
		},
		{
			name:            "Command/exit",
			args:            args{&parser.Command{Exit: true}},
			makeWantCommand: func(writer io.Writer) runtime.Command { return commands.ExitCommand{} },
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			writer := new(mocks.Writer)
			wantCommand := testData.makeWantCommand(writer)
			gotCommand, gotState := TranslateCommand(writer, testData.args.command)

			writer.AssertExpectations(test)
			assert.Equal(test, wantCommand, gotCommand)
			assert.Equal(test, testData.wantState, gotState)
		})
	}
}
