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

func TestTranslateCommands(test *testing.T) {
	type args struct {
		commands []*parser.Command
	}

	for _, testData := range []struct {
		name     string
		args     args
		makeWant func(writer io.Writer) runtime.CommandGroup
	}{
		{
			name: "success with commands",
			args: args{[]*parser.Command{{Send: tests.GetAddress("one")}, {Send: tests.GetAddress("two")}}},
			makeWant: func(writer io.Writer) runtime.CommandGroup {
				return runtime.CommandGroup{commands.NewSendCommand("one"), commands.NewSendCommand("two")}
			},
		},
		{
			name:     "success without commands",
			makeWant: func(writer io.Writer) runtime.CommandGroup { return nil },
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			writer := new(mocks.Writer)
			want := testData.makeWant(writer)
			got := TranslateCommands(writer, testData.args.commands)

			writer.AssertExpectations(test)
			assert.Equal(test, want, got)
		})
	}
}

func TestTranslateCommand(test *testing.T) {
	type args struct {
		command *parser.Command
	}

	for _, testData := range []struct {
		name     string
		args     args
		makeWant func(writer io.Writer) runtime.Command
	}{
		{
			name:     "Command/send",
			args:     args{&parser.Command{Send: tests.GetAddress("test")}},
			makeWant: func(writer io.Writer) runtime.Command { return commands.NewSendCommand("test") },
		},
		{
			name:     "Command/set",
			args:     args{&parser.Command{Set: tests.GetAddress("test")}},
			makeWant: func(writer io.Writer) runtime.Command { return commands.NewSetCommand("test") },
		},
		{
			name: "Command/out/nonempty",
			args: args{&parser.Command{Out: tests.GetAddress("test")}},
			makeWant: func(writer io.Writer) runtime.Command {
				return commands.NewOutCommand(writer, "test")
			},
		},
		{
			name:     "Command/out/empty",
			args:     args{&parser.Command{Out: tests.GetAddress("")}},
			makeWant: func(writer io.Writer) runtime.Command { return commands.NewOutCommand(writer, "") },
		},
		{
			name:     "Command/exit",
			args:     args{&parser.Command{Exit: true}},
			makeWant: func(writer io.Writer) runtime.Command { return commands.ExitCommand{} },
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			writer := new(mocks.Writer)
			want := testData.makeWant(writer)
			got := TranslateCommand(writer, testData.args.command)

			writer.AssertExpectations(test)
			assert.Equal(test, want, got)
		})
	}
}
