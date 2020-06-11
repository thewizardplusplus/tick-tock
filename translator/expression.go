package translator

import (
	"unicode/utf8"

	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// ...
const (
	EmptyListConstantName = "__empty_list__"

	ListConstructionFunctionName   = "__cons__"
	EqualFunctionName              = "__eq__"
	NotEqualFunctionName           = "__ne__"
	LessFunctionName               = "__lt__"
	LessOrEqualFunctionName        = "__le__"
	GreaterFunctionName            = "__gt__"
	GreaterOrEqualFunctionName     = "__ge__"
	AdditionFunctionName           = "__add__"
	SubtractionFunctionName        = "__sub__"
	MultiplicationFunctionName     = "__mul__"
	DivisionFunctionName           = "__div__"
	ModuloFunctionName             = "__mod__"
	ArithmeticNegationFunctionName = "__neg__"
	LogicalNegationFunctionName    = "__not__"
	KeyAccessorFunctionName        = "__item__"
)

func translateExpression(
	expression *parser.Expression,
	declaredIdentifiers mapset.Set,
) (
	result expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	result, settedStates, err =
		translateListConstruction(expression.ListConstruction, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the list construction")
	}

	return result, settedStates, nil
}

func translateListConstruction(
	listConstruction *parser.ListConstruction,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	argumentOne, settedStates, err :=
		translateDisjunction(listConstruction.Disjunction, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the disjunction")
	}
	if listConstruction.ListConstruction == nil {
		return argumentOne, settedStates, nil
	}

	argumentTwo, settedStates2, err :=
		translateListConstruction(listConstruction.ListConstruction, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the list construction")
	}

	expression = expressions.NewFunctionCall(
		ListConstructionFunctionName,
		[]expressions.Expression{argumentOne, argumentTwo},
	)
	settedStates = settedStates.Union(settedStates2)

	return expression, settedStates, nil
}

func translateDisjunction(
	disjunction *parser.Disjunction,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	argumentOne, settedStates, err :=
		translateConjunction(disjunction.Conjunction, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the conjunction")
	}
	if disjunction.Disjunction == nil {
		return argumentOne, settedStates, nil
	}

	argumentTwo, settedStates2, err :=
		translateDisjunction(disjunction.Disjunction, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the disjunction")
	}

	expression = expressions.NewBooleanOperator(argumentOne, argumentTwo, types.True)
	settedStates = settedStates.Union(settedStates2)

	return expression, settedStates, nil
}

func translateConjunction(
	conjunction *parser.Conjunction,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	argumentOne, settedStates, err := translateEquality(conjunction.Equality, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the equality")
	}
	if conjunction.Conjunction == nil {
		return argumentOne, settedStates, nil
	}

	argumentTwo, settedStates2, err :=
		translateConjunction(conjunction.Conjunction, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the conjunction")
	}

	expression = expressions.NewBooleanOperator(argumentOne, argumentTwo, types.False)
	settedStates = settedStates.Union(settedStates2)

	return expression, settedStates, nil
}

func translateEquality(
	equality *parser.Equality,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	argumentOne, settedStates, err := translateComparison(equality.Comparison, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the comparison")
	}
	if equality.Equality == nil {
		return argumentOne, settedStates, nil
	}

	argumentTwo, settedStates2, err := translateEquality(equality.Equality, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the equality")
	}

	var functionName string
	switch equality.Operation {
	case "==":
		functionName = EqualFunctionName
	case "!=":
		functionName = NotEqualFunctionName
	}

	expression =
		expressions.NewFunctionCall(functionName, []expressions.Expression{argumentOne, argumentTwo})
	settedStates = settedStates.Union(settedStates2)

	return expression, settedStates, nil
}

func translateComparison(
	comparison *parser.Comparison,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	argumentOne, settedStates, err := translateAddition(comparison.Addition, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the addition")
	}
	if comparison.Comparison == nil {
		return argumentOne, settedStates, nil
	}

	argumentTwo, settedStates2, err := translateComparison(comparison.Comparison, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the comparison")
	}

	var functionName string
	switch comparison.Operation {
	case "<":
		functionName = LessFunctionName
	case "<=":
		functionName = LessOrEqualFunctionName
	case ">":
		functionName = GreaterFunctionName
	case ">=":
		functionName = GreaterOrEqualFunctionName
	}

	expression =
		expressions.NewFunctionCall(functionName, []expressions.Expression{argumentOne, argumentTwo})
	settedStates = settedStates.Union(settedStates2)

	return expression, settedStates, nil
}

func translateAddition(
	addition *parser.Addition,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	argumentOne, settedStates, err :=
		translateMultiplication(addition.Multiplication, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the multiplication")
	}
	if addition.Addition == nil {
		return argumentOne, settedStates, nil
	}

	argumentTwo, settedStates2, err := translateAddition(addition.Addition, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the addition")
	}

	var functionName string
	switch addition.Operation {
	case "+":
		functionName = AdditionFunctionName
	case "-":
		functionName = SubtractionFunctionName
	}

	expression =
		expressions.NewFunctionCall(functionName, []expressions.Expression{argumentOne, argumentTwo})
	settedStates = settedStates.Union(settedStates2)

	return expression, settedStates, nil
}

func translateMultiplication(
	multiplication *parser.Multiplication,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	argumentOne, settedStates, err := translateUnary(multiplication.Unary, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the unary")
	}
	if multiplication.Multiplication == nil {
		return argumentOne, settedStates, nil
	}

	argumentTwo, settedStates2, err :=
		translateMultiplication(multiplication.Multiplication, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the multiplication")
	}

	var functionName string
	switch multiplication.Operation {
	case "*":
		functionName = MultiplicationFunctionName
	case "/":
		functionName = DivisionFunctionName
	case "%":
		functionName = ModuloFunctionName
	}

	expression =
		expressions.NewFunctionCall(functionName, []expressions.Expression{argumentOne, argumentTwo})
	settedStates.Union(settedStates2)

	return expression, settedStates, nil
}

func translateUnary(
	unary *parser.Unary,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	if unary.Accessor != nil {
		expression, settedStates, err = translateAccessor(unary.Accessor, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to translate the accessor")
		}

		return expression, settedStates, nil
	}

	argument, settedStates, err := translateUnary(unary.Unary, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the unary")
	}

	var functionName string
	switch unary.Operation {
	case "-":
		functionName = ArithmeticNegationFunctionName
	case "!":
		functionName = LogicalNegationFunctionName
	}

	expression = expressions.NewFunctionCall(functionName, []expressions.Expression{argument})
	return expression, settedStates, nil
}

func translateAccessor(
	accessor *parser.Accessor,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	argumentOne, settedStates, err := translateAtom(accessor.Atom, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the atom")
	}

	for index, key := range accessor.Keys {
		argumentTwo, settedStates2, err := translateExpression(key, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to translate the key #%d of the accessor", index)
		}

		argumentOne = expressions.NewFunctionCall(
			KeyAccessorFunctionName,
			[]expressions.Expression{argumentOne, argumentTwo},
		)
		settedStates.Union(settedStates2)
	}

	return argumentOne, settedStates, nil
}

func translateAtom(
	atom *parser.Atom,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	switch {
	case atom.Number != nil:
		expression = expressions.NewNumber(*atom.Number)
	case atom.Symbol != nil:
		symbol, _ := utf8.DecodeRuneInString(*atom.Symbol)
		expression = expressions.NewNumber(float64(symbol))
	case atom.String != nil:
		expression = expressions.NewString(*atom.String)
	case atom.Identifier != nil:
		identifier := *atom.Identifier
		if !declaredIdentifiers.Contains(identifier) {
			return nil, nil, errors.Errorf("unknown identifier %s", identifier)
		}

		expression = expressions.NewIdentifier(identifier)
	case atom.ListDefinition != nil:
		expression, err = translateListDefinition(atom.ListDefinition, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to translate the list definition")
		}
	case atom.FunctionCall != nil:
		expression, err = translateFunctionCall(atom.FunctionCall, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to translate the function call")
		}
	case atom.ConditionalExpression != nil:
		expression, settedStates, err =
			translateConditionalExpression(atom.ConditionalExpression, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to translate the conditional expression")
		}
	case atom.Expression != nil:
		expression, settedStates, err = translateExpression(atom.Expression, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to translate the expression")
		}
	}
	if settedStates == nil {
		settedStates = mapset.NewSet()
	}

	return expression, settedStates, nil
}

func translateListDefinition(
	listDefinition *parser.ListDefinition,
	declaredIdentifiers mapset.Set,
) (expressions.Expression, error) {
	argumentTwo := expressions.Expression(expressions.NewIdentifier(EmptyListConstantName))
	for index := len(listDefinition.Items) - 1; index >= 0; index-- {
		argumentOne, _, err := translateExpression(listDefinition.Items[index], declaredIdentifiers)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to translate the item #%d of the list definition", index)
		}

		argumentTwo = expressions.NewFunctionCall(
			ListConstructionFunctionName,
			[]expressions.Expression{argumentOne, argumentTwo},
		)
	}

	return argumentTwo, nil
}

func translateFunctionCall(
	functionCall *parser.FunctionCall,
	declaredIdentifiers mapset.Set,
) (expressions.Expression, error) {
	if !declaredIdentifiers.Contains(functionCall.Name) {
		return nil, errors.Errorf("unknown function %s", functionCall.Name)
	}

	var arguments []expressions.Expression
	for index, argument := range functionCall.Arguments {
		result, _, err := translateExpression(argument, declaredIdentifiers)
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"unable to translate the argument #%d for the function %s",
				index,
				functionCall.Name,
			)
		}

		arguments = append(arguments, result)
	}

	expression := expressions.NewFunctionCall(functionCall.Name, arguments)
	return expression, nil
}

func translateConditionalExpression(
	conditionalExpression *parser.ConditionalExpression,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	var conditionalCases []expressions.ConditionalCase
	settedStates = mapset.NewSet()
	for index, conditionalCase := range conditionalExpression.ConditionalCases {
		condition, settedStates2, err :=
			translateExpression(conditionalCase.Condition, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to translate the condition #%d", index)
		}

		commands, _, err := translateCommands(conditionalCase.Commands, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to translate commands of the condition #%d", index)
		}

		conditionalCases = append(conditionalCases, expressions.ConditionalCase{
			Condition: condition,
			Command:   commands,
		})
		settedStates.Union(settedStates2)
	}

	expression = expressions.NewConditionalExpression(conditionalCases)
	return expression, settedStates, nil
}
