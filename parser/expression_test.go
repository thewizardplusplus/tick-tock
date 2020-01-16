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
	} {
		test.Run(testData.name, func(test *testing.T) {
			err := parseToAST(testData.args.code, testData.args.ast)
			assert.Equal(test, testData.wantAST, testData.args.ast)
			testData.wantErr(test, err)
		})
	}
}
