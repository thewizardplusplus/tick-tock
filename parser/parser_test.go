package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
			want:    &Program{[]*Actor{{[]*State{{false, "one", nil}}}, {[]*State{{false, "two", nil}}}}},
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
			wantAST: &Command{Send: getAddress("test")},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/set",
			args:    args{"set test", new(Command)},
			wantAST: &Command{Set: getAddress("test")},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/out/nonempty",
			args:    args{`out "test"`, new(Command)},
			wantAST: &Command{Out: getAddress("test")},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/out/empty",
			args:    args{`out ""`, new(Command)},
			wantAST: &Command{Out: getAddress("")},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/exit",
			args:    args{"exit", new(Command)},
			wantAST: &Command{Exit: true},
			wantErr: assert.NoError,
		},
		{
			name:    "Message/nonempty",
			args:    args{"message test send one send two;", new(Message)},
			wantAST: &Message{"test", []*Command{{Send: getAddress("one")}, {Send: getAddress("two")}}},
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
			wantAST: &State{false, "test", []*Message{{"one", nil}, {"two", nil}}},
			wantErr: assert.NoError,
		},
		{
			name:    "State/empty",
			args:    args{"state test;", new(State)},
			wantAST: &State{false, "test", nil},
			wantErr: assert.NoError,
		},
		{
			name:    "State/initial",
			args:    args{"initial state test;", new(State)},
			wantAST: &State{true, "test", nil},
			wantErr: assert.NoError,
		},
		{
			name:    "Actor/nonempty",
			args:    args{"actor state one; state two;;", new(Actor)},
			wantAST: &Actor{[]*State{{false, "one", nil}, {false, "two", nil}}},
			wantErr: assert.NoError,
		},
		{
			name:    "Actor/with an initial state",
			args:    args{"actor state one; initial state two;;", new(Actor)},
			wantAST: &Actor{[]*State{{false, "one", nil}, {true, "two", nil}}},
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
			wantAST: &Program{[]*Actor{{[]*State{{false, "one", nil}}}, {[]*State{{false, "two", nil}}}}},
			wantErr: assert.NoError,
		},
		{
			name:    "Program/empty",
			args:    args{"", new(Program)},
			wantAST: new(Program),
			wantErr: assert.NoError,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			err := ParseToAST(testData.args.code, testData.args.ast)
			assert.Equal(test, testData.wantAST, testData.args.ast)
			testData.wantErr(test, err)
		})
	}
}

func getAddress(s string) *string {
	return &s
}
