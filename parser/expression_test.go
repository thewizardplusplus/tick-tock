package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/internal/tests"
)

func TestParseToAST_withExpression(test *testing.T) {
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
			name:    "Atom/number/integer",
			args:    args{"23", new(Atom)},
			wantAST: &Atom{Number: tests.GetNumberAddress(23)},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/number/floating-point",
			args:    args{"2.3", new(Atom)},
			wantAST: &Atom{Number: tests.GetNumberAddress(2.3)},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/identifier",
			args:    args{"test", new(Atom)},
			wantAST: &Atom{Identifier: tests.GetStringAddress("test")},
			wantErr: assert.NoError,
		},
		{
			name: "Unary/nonempty",
			args: args{"--23", new(Unary)},
			wantAST: &Unary{
				Operation: "-",
				Unary: &Unary{
					Operation: "-",
					Unary:     &Unary{Atom: &Atom{Number: tests.GetNumberAddress(23)}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Unary/empty",
			args:    args{"23", new(Unary)},
			wantAST: &Unary{Atom: &Atom{Number: tests.GetNumberAddress(23)}},
			wantErr: assert.NoError,
		},
		{
			name: "Multiplication/nonempty",
			args: args{"12 * 23 / 42", new(Multiplication)},
			wantAST: &Multiplication{
				Unary:     &Unary{Atom: &Atom{Number: tests.GetNumberAddress(12)}},
				Operation: "*",
				Multiplication: &Multiplication{
					Unary:     &Unary{Atom: &Atom{Number: tests.GetNumberAddress(23)}},
					Operation: "/",
					Multiplication: &Multiplication{
						Unary: &Unary{Atom: &Atom{Number: tests.GetNumberAddress(42)}},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Multiplication/empty",
			args: args{"23", new(Multiplication)},
			wantAST: &Multiplication{
				Unary: &Unary{Atom: &Atom{Number: tests.GetNumberAddress(23)}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Addition/nonempty",
			args: args{"12 + 23 - 42", new(Addition)},
			wantAST: &Addition{
				Multiplication: &Multiplication{
					Unary: &Unary{Atom: &Atom{Number: tests.GetNumberAddress(12)}},
				},
				Operation: "+",
				Addition: &Addition{
					Multiplication: &Multiplication{
						Unary: &Unary{Atom: &Atom{Number: tests.GetNumberAddress(23)}},
					},
					Operation: "-",
					Addition: &Addition{
						Multiplication: &Multiplication{
							Unary: &Unary{Atom: &Atom{Number: tests.GetNumberAddress(42)}},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Addition/empty",
			args: args{"23", new(Addition)},
			wantAST: &Addition{
				Multiplication: &Multiplication{
					Unary: &Unary{Atom: &Atom{Number: tests.GetNumberAddress(23)}},
				},
			},
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
