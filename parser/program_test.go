package parser

import (
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestParseToAST_withProgram(test *testing.T) {
	type args struct {
		code string
		ast  interface{}
	}

	for _, testData := range []struct {
		name    string
		args    args
		wantAST interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "Command/let",
			args: args{"let number = 23", new(Command)},
			wantAST: &Command{
				Let: &LetCommand{
					Identifier: "number",
					Expression: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(23)).(*Expression),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Command/start/identifier/no arguments",
			args: args{"start Test()", new(Command)},
			wantAST: &Command{
				Start: &StartCommand{Name: pointer.ToString("Test"), Arguments: &ExpressionGroup{}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Command/start/identifier/few arguments",
			args: args{"start Test(12, 23, 42)", new(Command)},
			wantAST: &Command{
				Start: &StartCommand{
					Name: pointer.ToString("Test"),
					Arguments: &ExpressionGroup{[]*Expression{
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression),
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(23)).(*Expression),
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(42)).(*Expression),
					}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Command/start/expression",
			args: args{"start [test()]()", new(Command)},
			wantAST: &Command{
				Start: &StartCommand{
					Expression: SetInnerField(&Expression{}, "FunctionCall", &FunctionCall{
						Name:      "test",
						Arguments: &ExpressionGroup{},
					}).(*Expression),
					Arguments: &ExpressionGroup{},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/send/no arguments",
			args:    args{"send test()", new(Command)},
			wantAST: &Command{Send: &SendCommand{Name: "test", Arguments: &ExpressionGroup{}}},
			wantErr: assert.NoError,
		},
		{
			name: "Command/send/few arguments",
			args: args{"send test(12, 23, 42)", new(Command)},
			wantAST: &Command{
				Send: &SendCommand{
					Name: "test",
					Arguments: &ExpressionGroup{[]*Expression{
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression),
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(23)).(*Expression),
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(42)).(*Expression),
					}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/set/no arguments",
			args:    args{"set test()", new(Command)},
			wantAST: &Command{Set: &SetCommand{Name: "test", Arguments: &ExpressionGroup{}}},
			wantErr: assert.NoError,
		},
		{
			name: "Command/set/few arguments",
			args: args{"set test(12, 23, 42)", new(Command)},
			wantAST: &Command{
				Set: &SetCommand{
					Name: "test",
					Arguments: &ExpressionGroup{[]*Expression{
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression),
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(23)).(*Expression),
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(42)).(*Expression),
					}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Command/return",
			args:    args{"return", new(Command)},
			wantAST: &Command{Return: true},
			wantErr: assert.NoError,
		},
		{
			name: "Command/expression",
			args: args{"test()", new(Command)},
			wantAST: &Command{
				Expression: SetInnerField(&Expression{}, "FunctionCall", &FunctionCall{
					Name:      "test",
					Arguments: &ExpressionGroup{},
				}).(*Expression),
			},
			wantErr: assert.NoError,
		},
		{
			name: "Message/nonempty/no parameters",
			args: args{"message test() send one() send two();", new(Message)},
			wantAST: &Message{
				Name:       "test",
				Parameters: &IdentifierGroup{},
				Commands: []*Command{
					{Send: &SendCommand{Name: "one", Arguments: &ExpressionGroup{}}},
					{Send: &SendCommand{Name: "two", Arguments: &ExpressionGroup{}}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Message/nonempty/few parameters",
			args: args{"message test(x, y, z) send one() send two();", new(Message)},
			wantAST: &Message{
				Name:       "test",
				Parameters: &IdentifierGroup{Identifiers: []string{"x", "y", "z"}},
				Commands: []*Command{
					{Send: &SendCommand{Name: "one", Arguments: &ExpressionGroup{}}},
					{Send: &SendCommand{Name: "two", Arguments: &ExpressionGroup{}}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Message/empty",
			args:    args{"message test();", new(Message)},
			wantAST: &Message{"test", &IdentifierGroup{}, nil},
			wantErr: assert.NoError,
		},
		{
			name: "State/nonempty/no parameters",
			args: args{"state test() message one(); message two();;", new(State)},
			wantAST: &State{
				Name:       "test",
				Parameters: &IdentifierGroup{},
				Messages:   []*Message{{"one", &IdentifierGroup{}, nil}, {"two", &IdentifierGroup{}, nil}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "State/nonempty/few parameters",
			args: args{"state test(x, y, z) message one(); message two();;", new(State)},
			wantAST: &State{
				Name:       "test",
				Parameters: &IdentifierGroup{Identifiers: []string{"x", "y", "z"}},
				Messages:   []*Message{{"one", &IdentifierGroup{}, nil}, {"two", &IdentifierGroup{}, nil}},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "State/empty",
			args:    args{"state test();", new(State)},
			wantAST: &State{"test", &IdentifierGroup{}, nil},
			wantErr: assert.NoError,
		},
		{
			name: "Actor/nonempty/no parameters",
			args: args{"actor Main() state one(); state two();;", new(Actor)},
			wantAST: &Actor{
				Name:       "Main",
				Parameters: &IdentifierGroup{},
				States:     []*State{{"one", &IdentifierGroup{}, nil}, {"two", &IdentifierGroup{}, nil}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Actor/nonempty/few parameters",
			args: args{"actor Main(x, y, z) state one(); state two();;", new(Actor)},
			wantAST: &Actor{
				Name:       "Main",
				Parameters: &IdentifierGroup{Identifiers: []string{"x", "y", "z"}},
				States:     []*State{{"one", &IdentifierGroup{}, nil}, {"two", &IdentifierGroup{}, nil}},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Actor/empty",
			args:    args{"actor Main();", new(Actor)},
			wantAST: &Actor{"Main", &IdentifierGroup{}, nil},
			wantErr: assert.NoError,
		},
		{
			name: "ActorClass/nonempty/no parameters",
			args: args{"class Main() state one(); state two();;", new(ActorClass)},
			wantAST: &ActorClass{
				Name:       "Main",
				Parameters: &IdentifierGroup{},
				States:     []*State{{"one", &IdentifierGroup{}, nil}, {"two", &IdentifierGroup{}, nil}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "ActorClass/nonempty/few parameters",
			args: args{"class Main(x, y, z) state one(); state two();;", new(ActorClass)},
			wantAST: &ActorClass{
				Name:       "Main",
				Parameters: &IdentifierGroup{Identifiers: []string{"x", "y", "z"}},
				States:     []*State{{"one", &IdentifierGroup{}, nil}, {"two", &IdentifierGroup{}, nil}},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "ActorClass/empty",
			args:    args{"class Main();", new(ActorClass)},
			wantAST: &ActorClass{"Main", &IdentifierGroup{}, nil},
			wantErr: assert.NoError,
		},
		{
			name: "Definition/actor",
			args: args{"actor Main() state one(); state two();;", new(Definition)},
			wantAST: &Definition{
				Actor: &Actor{
					Name:       "Main",
					Parameters: &IdentifierGroup{},
					States:     []*State{{"one", &IdentifierGroup{}, nil}, {"two", &IdentifierGroup{}, nil}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Definition/actor class",
			args: args{"class Main() state one(); state two();;", new(Definition)},
			wantAST: &Definition{
				ActorClass: &ActorClass{
					Name:       "Main",
					Parameters: &IdentifierGroup{},
					States:     []*State{{"one", &IdentifierGroup{}, nil}, {"two", &IdentifierGroup{}, nil}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Program/nonempty",
			args: args{"actor One(); actor Two();", new(Program)},
			wantAST: &Program{
				Definitions: []*Definition{
					{Actor: &Actor{"One", &IdentifierGroup{}, nil}},
					{Actor: &Actor{"Two", &IdentifierGroup{}, nil}},
				},
			},
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
