package translator

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
)

type declaredIdentifierGroup map[string]struct{}

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
