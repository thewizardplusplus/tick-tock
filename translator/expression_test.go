package translator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/thewizardplusplus/tick-tock/internal/tests"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
)

func TestTranslateUnary(test *testing.T) {
	type args struct {
		unary               *parser.Unary
		declaredIdentifiers declaredIdentifierGroup
	}

	for _, data := range []struct {
		name           string
		args           args
		wantExpression expressions.Expression
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "Unary/nonempty/success",
			args: args{
				unary: &parser.Unary{
					Operation: "-",
					Unary: &parser.Unary{
						Operation: "-",
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
						},
					},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantExpression: func() expressions.Expression {
				var expression expressions.Expression
				expression = expressions.NewNumber(23)
				expression = expressions.NewFunctionCall(
					NegationFunctionName,
					[]expressions.Expression{expression},
				)
				expression = expressions.NewFunctionCall(
					NegationFunctionName,
					[]expressions.Expression{expression},
				)

				return expression
			}(),
			wantErr: assert.NoError,
		},
		{
			name: "Unary/nonempty/error",
			args: args{
				unary: &parser.Unary{
					Operation: "-",
					Unary: &parser.Unary{
						Operation: "-",
						Unary: &parser.Unary{
							Accessor: &parser.Accessor{
								Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
							},
						},
					},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
		{
			name: "Unary/empty/success",
			args: args{
				unary: &parser.Unary{
					Accessor: &parser.Accessor{Atom: &parser.Atom{Number: tests.GetNumberAddress(23)}},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantExpression: expressions.NewNumber(23),
			wantErr:        assert.NoError,
		},
		{
			name: "Unary/empty/error",
			args: args{
				unary: &parser.Unary{
					Accessor: &parser.Accessor{Atom: &parser.Atom{Identifier: tests.GetStringAddress("unknown")}},
				},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotErr := translateUnary(data.args.unary, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			data.wantErr(test, gotErr)
		})
	}
}

func TestTranslateAtom(test *testing.T) {
	type args struct {
		atom                *parser.Atom
		declaredIdentifiers declaredIdentifierGroup
	}

	for _, data := range []struct {
		name           string
		args           args
		wantExpression expressions.Expression
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "Atom/number",
			args: args{
				atom:                &parser.Atom{Number: tests.GetNumberAddress(23)},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantExpression: expressions.NewNumber(23),
			wantErr:        assert.NoError,
		},
		{
			name: "Atom/string",
			args: args{
				atom:                &parser.Atom{String: tests.GetStringAddress("test")},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantExpression: expressions.NewString("test"),
			wantErr:        assert.NoError,
		},
		{
			name: "Atom/identifier/success",
			args: args{
				atom:                &parser.Atom{Identifier: tests.GetStringAddress("test")},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantExpression: expressions.NewIdentifier("test"),
			wantErr:        assert.NoError,
		},
		{
			name: "Atom/identifier/error",
			args: args{
				atom:                &parser.Atom{Identifier: tests.GetStringAddress("unknown")},
				declaredIdentifiers: declaredIdentifierGroup{"test": {}},
			},
			wantExpression: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotExpression, gotErr := translateAtom(data.args.atom, data.args.declaredIdentifiers)

			assert.Equal(test, data.wantExpression, gotExpression)
			data.wantErr(test, gotErr)
		})
	}
}
