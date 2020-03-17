package translator

import (
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
)

type declaredIdentifierGroup map[string]struct{}

// ...
const (
	ListConstructionFunctionName = "__cons__"
	AdditionFunctionName         = "__add__"
	SubtractionFunctionName      = "__sub__"
	MultiplicationFunctionName   = "__mul__"
	DivisionFunctionName         = "__div__"
	ModuloFunctionName           = "__mod__"
	NegationFunctionName         = "__neg__"
)

func translateExpression(
	expression *parser.Expression,
	declaredIdentifiers declaredIdentifierGroup,
) (expressions.Expression, error) {
	result, err := translateListConstruction(expression.ListConstruction, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the list construction")
	}

	return result, nil
}

func translateListConstruction(
	listConstruction *parser.ListConstruction,
	declaredIdentifiers declaredIdentifierGroup,
) (expressions.Expression, error) {
	argumentOne, err := translateAddition(listConstruction.Addition, declaredIdentifiers)
	if err != nil {
		return nil, errors.Wrap(err, "unable to translate the addition")
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

func translateAddition(
	addition *parser.Addition,
	declaredIdentifiers declaredIdentifierGroup,
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
	declaredIdentifiers declaredIdentifierGroup,
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
	case atom.FunctionCall != nil:
		result, err := translateFunctionCall(atom.FunctionCall, declaredIdentifiers)
		if err != nil {
			return nil, errors.Wrap(err, "unable to translate the function call")
		}

		expression = result
	case atom.Expression != nil:
		result, err := translateExpression(atom.Expression, declaredIdentifiers)
		if err != nil {
			return nil, errors.Wrap(err, "unable to translate the expression")
		}

		expression = result
	}

	return expression, nil
}

func translateFunctionCall(
	functionCall *parser.FunctionCall,
	declaredIdentifiers declaredIdentifierGroup,
) (expressions.Expression, error) {
	if _, ok := declaredIdentifiers[functionCall.Name]; !ok {
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
