package translator

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
)

func translateExpressionGroup(expressions *parser.ExpressionGroup, declaredIdentifiers mapset.Set) (
	translatedExpressions []expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	settedStates = mapset.NewSet()
	for index, expression := range expressions.Expressions {
		translatedExpression, settedStatesByExpression, err :=
			TranslateExpression(expression, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to translate the expression #%d", index)
		}

		translatedExpressions = append(translatedExpressions, translatedExpression)
		settedStates = settedStates.Union(settedStatesByExpression)
	}

	return translatedExpressions, settedStates, nil
}
