package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/internal/tests"
)

func TestParse(test *testing.T) {
	type args struct {
		code string
	}

	for _, testData := range []struct {
		name    string
		args    args
		want    *Program
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "success",
			args:    args{"actor state one;; actor state two;;"},
			want:    &Program{[]*Actor{{[]*State{{"one", nil}}}, {[]*State{{"two", nil}}}}},
			wantErr: assert.NoError,
		},
		{
			name:    "error",
			args:    args{"incorrect"},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			got, err := Parse(testData.args.code)
			assert.Equal(test, testData.want, got)
			testData.wantErr(test, err)
		})
	}
}

func TestParseToAST(test *testing.T) {
	type (
		args struct {
			code string
			ast  interface{}
		}
		testAST struct {
			Number int `parser:"@Int"`
		}
	)

	for _, testData := range []struct {
		name    string
		args    args
		wantAST interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "testAST/success",
			args:    args{"23", new(testAST)},
			wantAST: &testAST{23},
			wantErr: assert.NoError,
		},
		{
			name:    "testAST/error/building",
			args:    args{"23", "incorrect"},
			wantAST: "incorrect",
			wantErr: assert.Error,
		},
		{
			name:    "testAST/error/parsing",
			args:    args{"incorrect", new(testAST)},
			wantAST: new(testAST),
			wantErr: assert.Error,
		},
		{
			name:    "Command/send",
			args:    args{"send test", new(Command)},
			wantAST: &Command{Send: tests.GetStringAddress("test")},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/set",
			args:    args{"set test", new(Command)},
			wantAST: &Command{Set: tests.GetStringAddress("test")},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/out/nonempty/interpreted",
			args:    args{`out "test"`, new(Command)},
			wantAST: &Command{Out: tests.GetStringAddress("test")},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/out/nonempty/raw",
			args:    args{"out `test`", new(Command)},
			wantAST: &Command{Out: tests.GetStringAddress("test")},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/out/empty",
			args:    args{`out ""`, new(Command)},
			wantAST: &Command{Out: tests.GetStringAddress("")},
			wantErr: assert.NoError,
		},
		{
			name: "Command/sleep/integer",
			args: args{`sleep 1, 2`, new(Command)},
			wantAST: &Command{
				Sleep: &SleepCommand{tests.GetNumberAddress(1), tests.GetNumberAddress(2)},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Command/sleep/floating-point",
			args: args{`sleep 1.2, 3.4`, new(Command)},
			wantAST: &Command{
				Sleep: &SleepCommand{tests.GetNumberAddress(1.2), tests.GetNumberAddress(3.4)},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/exit",
			args:    args{"exit", new(Command)},
			wantAST: &Command{Exit: true},
			wantErr: assert.NoError,
		},
		{
			name: "Message/nonempty",
			args: args{"message test send one send two;", new(Message)},
			wantAST: &Message{
				Name: "test",
				Commands: []*Command{
					{Send: tests.GetStringAddress("one")},
					{Send: tests.GetStringAddress("two")},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Message/empty",
			args:    args{"message test;", new(Message)},
			wantAST: &Message{"test", nil},
			wantErr: assert.NoError,
		},
		{
			name:    "State/nonempty",
			args:    args{"state test message one; message two;;", new(State)},
			wantAST: &State{"test", []*Message{{"one", nil}, {"two", nil}}},
			wantErr: assert.NoError,
		},
		{
			name:    "State/empty",
			args:    args{"state test;", new(State)},
			wantAST: &State{"test", nil},
			wantErr: assert.NoError,
		},
		{
			name:    "Actor/nonempty",
			args:    args{"actor state one; state two;;", new(Actor)},
			wantAST: &Actor{[]*State{{"one", nil}, {"two", nil}}},
			wantErr: assert.NoError,
		},
		{
			name:    "Actor/empty",
			args:    args{"actor;", new(Actor)},
			wantAST: new(Actor),
			wantErr: assert.NoError,
		},
		{
			name:    "Program/nonempty",
			args:    args{"actor state one;; actor state two;;", new(Program)},
			wantAST: &Program{[]*Actor{{[]*State{{"one", nil}}}, {[]*State{{"two", nil}}}}},
			wantErr: assert.NoError,
		},
		{
			name:    "Program/empty",
			args:    args{"", new(Program)},
			wantAST: new(Program),
			wantErr: assert.NoError,
		},
		{
			name:    "comment/line",
			args:    args{"actor state one;;\n// actor state two;;\nactor state three;;", new(Program)},
			wantAST: &Program{[]*Actor{{[]*State{{"one", nil}}}, {[]*State{{"three", nil}}}}},
			wantErr: assert.NoError,
		},
		{
			name:    "comment/block",
			args:    args{"actor state one;; /* actor state two;; */ actor state three;;", new(Program)},
			wantAST: &Program{[]*Actor{{[]*State{{"one", nil}}}, {[]*State{{"three", nil}}}}},
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
