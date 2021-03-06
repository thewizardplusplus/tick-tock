package translator

import (
	"reflect"
	"unicode/utf8"

	mapset "github.com/deckarep/golang-set"
	"github.com/pkg/errors"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

// ...
const (
	EmptyListConstantName      = "__empty_list__"
	EmptyHashTableConstantName = "__empty_hash__"

	ListConstructionFunctionName            = "__cons__"
	HashTableConstructionFunctionName       = "__with__"
	EqualFunctionName                       = "__eq__"
	NotEqualFunctionName                    = "__ne__"
	LessFunctionName                        = "__lt__"
	LessOrEqualFunctionName                 = "__le__"
	GreaterFunctionName                     = "__gt__"
	GreaterOrEqualFunctionName              = "__ge__"
	BitwiseDisjunctionFunctionName          = "__or__"
	BitwiseExclusiveDisjunctionFunctionName = "__xor__"
	BitwiseConjunctionFunctionName          = "__and__"
	BitwiseLeftShiftFunctionName            = "__lshift__"
	BitwiseRightShiftFunctionName           = "__rshift__"
	BitwiseUnsignedRightShiftFunctionName   = "__urshift__"
	AdditionFunctionName                    = "__add__"
	SubtractionFunctionName                 = "__sub__"
	MultiplicationFunctionName              = "__mul__"
	DivisionFunctionName                    = "__div__"
	ModuloFunctionName                      = "__mod__"
	ArithmeticNegationFunctionName          = "__neg__"
	BitwiseNegationFunctionName             = "__bitwise_not__"
	LogicalNegationFunctionName             = "__logical_not__"
	KeyAccessorFunctionName                 = "__item__"
)

var (
	binaryOperations = map[string]string{
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

// TranslateExpression ...
func TranslateExpression(
	expression *parser.Expression,
	declaredIdentifiers mapset.Set,
) (
	result expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	result, settedStates, err =
		translateBinaryOperation(expression.ListConstruction, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the list construction")
	}

	return result, settedStates, nil
}

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
		functionName := binaryOperations[operationName]
		translatedBinaryOperation = expressions.NewFunctionCall(
			functionName,
			[]expressions.Expression{translatedArgumentOne, translatedArgumentTwo},
		)
	}

	settedStates = settedStates.Union(settedStates2)
	return translatedBinaryOperation, settedStates, nil
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
	case "~":
		functionName = BitwiseNegationFunctionName
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
		var argumentTwo expressions.Expression
		var settedStates2 mapset.Set
		switch {
		case key.Name != nil:
			argumentTwo = expressions.NewString(*key.Name)
			settedStates2 = mapset.NewSet()
		case key.Expression != nil:
			argumentTwo, settedStates2, err = TranslateExpression(key.Expression, declaredIdentifiers)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "unable to translate the key #%d of the accessor", index)
			}
		}

		argumentOne = expressions.NewFunctionCall(
			KeyAccessorFunctionName,
			[]expressions.Expression{argumentOne, argumentTwo},
		)
		settedStates = settedStates.Union(settedStates2)
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
	settedStates = mapset.NewSet()
	switch {
	case atom.IntegerNumber != nil:
		expression = expressions.NewNumber(float64(*atom.IntegerNumber))
	case atom.FloatingPointNumber != nil:
		expression = expressions.NewNumber(*atom.FloatingPointNumber)
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
		expression, settedStates, err = translateListDefinition(atom.ListDefinition, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to translate the list definition")
		}
	case atom.HashTableDefinition != nil:
		expression, settedStates, err =
			translateHashTableDefinition(atom.HashTableDefinition, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to translate the hash table definition")
		}
	case atom.FunctionCall != nil:
		expression, settedStates, err = translateFunctionCall(atom.FunctionCall, declaredIdentifiers)
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
		expression, settedStates, err = TranslateExpression(atom.Expression, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrap(err, "unable to translate the expression")
		}
	}

	return expression, settedStates, nil
}

func translateListDefinition(
	listDefinition *parser.ListDefinition,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	items, settedStates, err := translateExpressionGroup(listDefinition.Items, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "unable to translate items for the list definition")
	}

	argumentTwo := expressions.Expression(expressions.NewIdentifier(EmptyListConstantName))
	for index := len(items) - 1; index >= 0; index-- {
		argumentOne := items[index]
		argumentTwo = expressions.NewFunctionCall(
			ListConstructionFunctionName,
			[]expressions.Expression{argumentOne, argumentTwo},
		)
	}

	return argumentTwo, settedStates, nil
}

func translateHashTableDefinition(
	hashTableDefinition *parser.HashTableDefinition,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	argumentOne := expressions.Expression(expressions.NewIdentifier(EmptyHashTableConstantName))
	settedStates = mapset.NewSet()
	for index, entry := range hashTableDefinition.Entries {
		var argumentTwo expressions.Expression
		var settedStates2 mapset.Set
		switch {
		case entry.Name != nil:
			argumentTwo = expressions.NewString(*entry.Name)
			settedStates2 = mapset.NewSet()
		case entry.Expression != nil:
			argumentTwo, settedStates2, err = TranslateExpression(entry.Expression, declaredIdentifiers)
			if err != nil {
				return nil, nil, errors.Wrapf(
					err,
					"unable to translate the key #%d for the hash table definition",
					index,
				)
			}
		}

		argumentThree, settedStates3, err := TranslateExpression(entry.Value, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(
				err,
				"unable to translate the value #%d for the hash table definition",
				index,
			)
		}

		argumentOne = expressions.NewFunctionCall(
			HashTableConstructionFunctionName,
			[]expressions.Expression{argumentOne, argumentTwo, argumentThree},
		)
		settedStates = settedStates.Union(settedStates2)
		settedStates = settedStates.Union(settedStates3)
	}

	return argumentOne, settedStates, nil
}

func translateFunctionCall(
	functionCall *parser.FunctionCall,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	if !declaredIdentifiers.Contains(functionCall.Name) {
		return nil, nil, errors.Errorf("unknown function %s", functionCall.Name)
	}

	arguments, settedStates, err :=
		translateExpressionGroup(functionCall.Arguments, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrapf(
			err,
			"unable to translate arguments for the function %s",
			functionCall.Name,
		)
	}

	expression = expressions.NewFunctionCall(functionCall.Name, arguments)
	return expression, settedStates, nil
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
			TranslateExpression(conditionalCase.Condition, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to translate the condition #%d", index)
		}

		commands, settedStates3, err := translateCommands(conditionalCase.Commands, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to translate commands of the condition #%d", index)
		}

		conditionalCases = append(conditionalCases, expressions.ConditionalCase{
			Condition: condition,
			Command:   commands,
		})
		settedStates = settedStates.Union(settedStates2)
		settedStates = settedStates.Union(settedStates3)
	}

	expression = expressions.NewConditionalExpression(conditionalCases)
	return expression, settedStates, nil
}
