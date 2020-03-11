package translator

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
)

type declaredIdentifierGroup map[string]struct{}

// ...
const (
	NegationFunctionName = "__neg__"
)

func translateUnary(
	unary *parser.Unary,
	declaredIdentifiers declaredIdentifierGroup,
) (expressions.Expression, error) {
	if unary.Accessor != nil {
		expression, err := translateAccessor(unary.Accessor, declaredIdentifiers)
		if err != nil {
			return nil, errors.Wrap(err, "unable to translate the accessor")
		}

		return expression, nil
	}

	argument, err := translateUnary(unary.Unary, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the unary")
	}

	var functionName string
	switch unary.Operation {
	case "-":
		functionName = NegationFunctionName
	}

	expression := expressions.NewFunctionCall(functionName, []expressions.Expression{argument})
	return expression, nil
}

func translateAccessor(
	accessor *parser.Accessor,
	declaredIdentifiers declaredIdentifierGroup,
) (expressions.Expression, error) {
	return translateAtom(accessor.Atom, declaredIdentifiers)
}

func translateAtom(
	atom *parser.Atom,
	declaredIdentifiers declaredIdentifierGroup,
) (expressions.Expression, error) {
	var expression expressions.Expression
	switch {
	case atom.Number != nil:
		expression = expressions.NewNumber(*atom.Number)
	case atom.String != nil:
		expression = expressions.NewString(*atom.String)
	case atom.Identifier != nil:
		identifier := *atom.Identifier
		if _, ok := declaredIdentifiers[identifier]; !ok {
			return nil, errors.Errorf("unknown identifier %s", identifier)
		}

		expression = expressions.NewIdentifier(identifier)
	}

	return expression, nil
}
