package translator

import (
	"reflect"

	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

var (
	operations = map[string]string{
		"*":   MultiplicationFunctionName,
		"/":   DivisionFunctionName,
		"%":   ModuloFunctionName,
		"+":   AdditionFunctionName,
		"-":   SubtractionFunctionName,
		"<<":  BitwiseLeftShiftFunctionName,
		">>":  BitwiseRightShiftFunctionName,
		">>>": BitwiseUnsignedRightShiftFunctionName,
		"&":   BitwiseConjunctionFunctionName,
		"^":   BitwiseExclusiveDisjunctionFunctionName,
		"|":   BitwiseDisjunctionFunctionName,
		"<":   LessFunctionName,
		"<=":  LessOrEqualFunctionName,
		">":   GreaterFunctionName,
		">=":  GreaterOrEqualFunctionName,
		"==":  EqualFunctionName,
		"!=":  NotEqualFunctionName,
		":":   ListConstructionFunctionName,
	}
)

func translateBinaryOperation(
	binaryOperation interface{},
	declaredIdentifiers mapset.Set,
) (
	translatedBinaryOperation expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	binaryOperationReflection := reflect.ValueOf(binaryOperation).Elem()
	argumentOneReflection := binaryOperationReflection.Field(0)
	argumentTwoReflection := binaryOperationReflection.Field(2)
	operationNameReflection := binaryOperationReflection.FieldByName("Operation")

	var translatedArgumentOne expressions.Expression
	if argumentOneReflection.Type() == reflect.TypeOf(&parser.Unary{}) {
		translatedArgumentOne, settedStates, err =
			translateUnary(argumentOneReflection.Interface().(*parser.Unary), declaredIdentifiers)
	} else {
		translatedArgumentOne, settedStates, err =
			translateBinaryOperation(argumentOneReflection.Interface(), declaredIdentifiers)
	}
	if err != nil {
		return nil, nil, errors.Wrapf(
			err,
			"unable to translate the %T value",
			argumentOneReflection.Interface(),
		)
	}
	if argumentTwoReflection.IsNil() {
		return translatedArgumentOne, settedStates, nil
	}

	translatedArgumentTwo, settedStates2, err :=
		translateBinaryOperation(argumentTwoReflection.Interface(), declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrapf(
			err,
			"unable to translate the %T value",
			argumentOneReflection.Interface(),
		)
	}

	switch operationName := operationNameReflection.Interface().(string); operationName {
	case "&&": // parser.Conjunction
		translatedBinaryOperation =
			expressions.NewBooleanOperator(translatedArgumentOne, translatedArgumentTwo, types.False)
	case "||": // parser.Disjunction
		translatedBinaryOperation =
			expressions.NewBooleanOperator(translatedArgumentOne, translatedArgumentTwo, types.True)
	case "??": // parser.NilCoalescing
		translatedBinaryOperation =
			expressions.NewNilCoalescingOperator(translatedArgumentOne, translatedArgumentTwo)
	default:
		functionName := operations[operationName]
		translatedBinaryOperation = expressions.NewFunctionCall(
			functionName,
			[]expressions.Expression{translatedArgumentOne, translatedArgumentTwo},
		)
	}

	settedStates = settedStates.Union(settedStates2)
	return translatedBinaryOperation, settedStates, nil
}

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
