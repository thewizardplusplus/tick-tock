package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
