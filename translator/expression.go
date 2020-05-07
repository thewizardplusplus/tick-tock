package translator

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// ...
const (
	EmptyListConstantName = "__empty_list__"

	ListConstructionFunctionName = "__cons__"
	EqualFunctionName            = "__eq__"
	NotEqualFunctionName         = "__ne__"
	LessFunctionName             = "__lt__"
	LessOrEqualFunctionName      = "__le__"
	GreatFunctionName            = "__gt__"
	GreatOrEqualFunctionName     = "__ge__"
	AdditionFunctionName         = "__add__"
	SubtractionFunctionName      = "__sub__"
	MultiplicationFunctionName   = "__mul__"
	DivisionFunctionName         = "__div__"
	ModuloFunctionName           = "__mod__"
	NegationFunctionName         = "__neg__"
	KeyAccessorFunctionName      = "__item__"
)

func translateExpression(
	expression *parser.Expression,
	declaredIdentifiers mapset.Set,
) (expressions.Expression, error) {
	result, err := translateListConstruction(expression.ListConstruction, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the list construction")
	}

	return result, nil
}

func translateListConstruction(
	listConstruction *parser.ListConstruction,
	declaredIdentifiers mapset.Set,
) (expressions.Expression, error) {
	conjunction := listConstruction.Disjunction.Conjunction
	argumentOne, err := translateConjunction(conjunction, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the conjunction")
	}
	if listConstruction.ListConstruction == nil {
		return argumentOne, nil
	}

	argumentTwo, err :=
		translateListConstruction(listConstruction.ListConstruction, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the list construction")
	}

	expression := expressions.NewFunctionCall(
		ListConstructionFunctionName,
		[]expressions.Expression{argumentOne, argumentTwo},
	)
	return expression, nil
}

func translateConjunction(
	conjunction *parser.Conjunction,
	declaredIdentifiers mapset.Set,
) (expressions.Expression, error) {
	argumentOne, err := translateEquality(conjunction.Equality, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the equality")
	}
	if conjunction.Conjunction == nil {
		return argumentOne, nil
	}

	argumentTwo, err := translateConjunction(conjunction.Conjunction, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the conjunction")
	}

	expression := expressions.NewBooleanOperator(argumentOne, argumentTwo, types.False)
	return expression, nil
}

func translateEquality(
	equality *parser.Equality,
	declaredIdentifiers mapset.Set,
) (expressions.Expression, error) {
	argumentOne, err := translateComparison(equality.Comparison, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the comparison")
	}
	if equality.Equality == nil {
		return argumentOne, nil
	}

	argumentTwo, err := translateEquality(equality.Equality, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the equality")
	}

	var functionName string
	switch equality.Operation {
	case "==":
		functionName = EqualFunctionName
	case "!=":
		functionName = NotEqualFunctionName
	}

	expression :=
		expressions.NewFunctionCall(functionName, []expressions.Expression{argumentOne, argumentTwo})
	return expression, nil
}

func translateComparison(
	comparison *parser.Comparison,
	declaredIdentifiers mapset.Set,
) (expressions.Expression, error) {
	argumentOne, err := translateAddition(comparison.Addition, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the addition")
	}
	if comparison.Comparison == nil {
		return argumentOne, nil
	}

	argumentTwo, err := translateComparison(comparison.Comparison, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the comparison")
	}

	var functionName string
	switch comparison.Operation {
	case "<":
		functionName = LessFunctionName
	case "<=":
		functionName = LessOrEqualFunctionName
	case ">":
		functionName = GreatFunctionName
	case ">=":
		functionName = GreatOrEqualFunctionName
	}

	expression :=
		expressions.NewFunctionCall(functionName, []expressions.Expression{argumentOne, argumentTwo})
	return expression, nil
}

func translateAddition(
	addition *parser.Addition,
	declaredIdentifiers mapset.Set,
) (expressions.Expression, error) {
	argumentOne, err := translateMultiplication(addition.Multiplication, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the multiplication")
	}
	if addition.Addition == nil {
		return argumentOne, nil
	}

	argumentTwo, err := translateAddition(addition.Addition, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the addition")
	}

	var functionName string
	switch addition.Operation {
	case "+":
		functionName = AdditionFunctionName
	case "-":
		functionName = SubtractionFunctionName
	}

	expression :=
		expressions.NewFunctionCall(functionName, []expressions.Expression{argumentOne, argumentTwo})
	return expression, nil
}

func translateMultiplication(
	multiplication *parser.Multiplication,
	declaredIdentifiers mapset.Set,
) (expressions.Expression, error) {
	argumentOne, err := translateUnary(multiplication.Unary, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the unary")
	}
	if multiplication.Multiplication == nil {
		return argumentOne, nil
	}

	argumentTwo, err := translateMultiplication(multiplication.Multiplication, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the multiplication")
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

	expression :=
		expressions.NewFunctionCall(functionName, []expressions.Expression{argumentOne, argumentTwo})
	return expression, nil
}

func translateUnary(
	unary *parser.Unary,
	declaredIdentifiers mapset.Set,
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
	declaredIdentifiers mapset.Set,
) (expressions.Expression, error) {
	argumentOne, err := translateAtom(accessor.Atom, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the atom")
	}

	for index, key := range accessor.Keys {
		argumentTwo, err := translateExpression(key, declaredIdentifiers)
		if err != nil {
			return nil, errors.Wrapf(err, "unable to translate the key #%d of the accessor", index)
		}

		argumentOne = expressions.NewFunctionCall(
			KeyAccessorFunctionName,
			[]expressions.Expression{argumentOne, argumentTwo},
		)
	}

	return argumentOne, nil
}

func translateAtom(
	atom *parser.Atom,
	declaredIdentifiers mapset.Set,
) (expression expressions.Expression, err error) {
	switch {
	case atom.Number != nil:
		expression = expressions.NewNumber(*atom.Number)
	case atom.String != nil:
		expression = expressions.NewString(*atom.String)
	case atom.Identifier != nil:
		identifier := *atom.Identifier
		if !declaredIdentifiers.Contains(identifier) {
			return nil, errors.Errorf("unknown identifier %s", identifier)
		}

		expression = expressions.NewIdentifier(identifier)
	case atom.ListDefinition != nil:
		expression, err = translateListDefinition(atom.ListDefinition, declaredIdentifiers)
		if err != nil {
			return nil, errors.Wrap(err, "unable to translate the list definition")
		}
	case atom.FunctionCall != nil:
		expression, err = translateFunctionCall(atom.FunctionCall, declaredIdentifiers)
		if err != nil {
			return nil, errors.Wrap(err, "unable to translate the function call")
		}
	case atom.Expression != nil:
		expression, err = translateExpression(atom.Expression, declaredIdentifiers)
		if err != nil {
			return nil, errors.Wrap(err, "unable to translate the expression")
		}
	}

	return expression, nil
}

func translateListDefinition(
	listDefinition *parser.ListDefinition,
	declaredIdentifiers mapset.Set,
) (expressions.Expression, error) {
	argumentTwo := expressions.Expression(expressions.NewIdentifier(EmptyListConstantName))
	for index := len(listDefinition.Items) - 1; index >= 0; index-- {
		argumentOne, err := translateExpression(listDefinition.Items[index], declaredIdentifiers)
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
		result, err := translateExpression(argument, declaredIdentifiers)
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
