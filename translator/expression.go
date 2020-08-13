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
	argumentOne, settedStates, err :=
		translateBitwiseDisjunction(comparison.BitwiseDisjunction, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the bitwise disjunction")
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

func translateBitwiseDisjunction(
	bitwiseDisjunction *parser.BitwiseDisjunction,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	argumentOne, settedStates, err := translateBitwiseExclusiveDisjunction(
		bitwiseDisjunction.BitwiseExclusiveDisjunction,
		declaredIdentifiers,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the bitwise exclusive disjunction")
	}
	if bitwiseDisjunction.BitwiseDisjunction == nil {
		return argumentOne, settedStates, nil
	}

	argumentTwo, settedStates2, err :=
		translateBitwiseDisjunction(bitwiseDisjunction.BitwiseDisjunction, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the bitwise disjunction")
	}

	expression = expressions.NewFunctionCall(
		BitwiseDisjunctionFunctionName,
		[]expressions.Expression{argumentOne, argumentTwo},
	)
	settedStates = settedStates.Union(settedStates2)

	return expression, settedStates, nil
}

func translateBitwiseExclusiveDisjunction(
	bitwiseExclusiveDisjunction *parser.BitwiseExclusiveDisjunction,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	argumentOne, settedStates, err :=
		translateBitwiseConjunction(bitwiseExclusiveDisjunction.BitwiseConjunction, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the bitwise conjunction")
	}
	if bitwiseExclusiveDisjunction.BitwiseExclusiveDisjunction == nil {
		return argumentOne, settedStates, nil
	}

	argumentTwo, settedStates2, err := translateBitwiseExclusiveDisjunction(
		bitwiseExclusiveDisjunction.BitwiseExclusiveDisjunction,
		declaredIdentifiers,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the bitwise exclusive disjunction")
	}

	expression = expressions.NewFunctionCall(
		BitwiseExclusiveDisjunctionFunctionName,
		[]expressions.Expression{argumentOne, argumentTwo},
	)
	settedStates = settedStates.Union(settedStates2)

	return expression, settedStates, nil
}

func translateBitwiseConjunction(
	bitwiseConjunction *parser.BitwiseConjunction,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	argumentOne, settedStates, err := translateShift(bitwiseConjunction.Shift, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the shift")
	}
	if bitwiseConjunction.BitwiseConjunction == nil {
		return argumentOne, settedStates, nil
	}

	argumentTwo, settedStates2, err :=
		translateBitwiseConjunction(bitwiseConjunction.BitwiseConjunction, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the bitwise conjunction")
	}

	expression = expressions.NewFunctionCall(
		BitwiseConjunctionFunctionName,
		[]expressions.Expression{argumentOne, argumentTwo},
	)
	settedStates = settedStates.Union(settedStates2)

	return expression, settedStates, nil
}

func translateShift(
	shift *parser.Shift,
	declaredIdentifiers mapset.Set,
) (
	expression expressions.Expression,
	settedStates mapset.Set,
	err error,
) {
	argumentOne, settedStates, err := translateAddition(shift.Addition, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the addition")
	}
	if shift.Shift == nil {
		return argumentOne, settedStates, nil
	}

	argumentTwo, settedStates2, err := translateShift(shift.Shift, declaredIdentifiers)
	if err != nil {
		return nil, nil, errors.Wrap(err, "unable to translate the shift")
	}

	var functionName string
	switch shift.Operation {
	case "<<":
		functionName = BitwiseLeftShiftFunctionName
	case ">>":
		functionName = BitwiseRightShiftFunctionName
	case ">>>":
		functionName = BitwiseUnsignedRightShiftFunctionName
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
	settedStates = settedStates.Union(settedStates2)

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
		argumentTwo, settedStates2, err := translateExpression(key, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(err, "unable to translate the key #%d of the accessor", index)
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
		expression, settedStates, err = translateExpression(atom.Expression, declaredIdentifiers)
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
	argumentTwo := expressions.Expression(expressions.NewIdentifier(EmptyListConstantName))
	settedStates = mapset.NewSet()
	for index := len(listDefinition.Items) - 1; index >= 0; index-- {
		argumentOne, settedStates2, err :=
			translateExpression(listDefinition.Items[index], declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(
				err,
				"unable to translate the item #%d of the list definition",
				index,
			)
		}

		argumentTwo = expressions.NewFunctionCall(
			ListConstructionFunctionName,
			[]expressions.Expression{argumentOne, argumentTwo},
		)
		settedStates = settedStates.Union(settedStates2)
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
			identifier := *entry.Name
			if !declaredIdentifiers.Contains(identifier) {
				return nil, nil, errors.Errorf("unknown identifier %s", identifier)
			}

			argumentTwo = expressions.NewIdentifier(identifier)
			settedStates2 = mapset.NewSet()
		case entry.Expression != nil:
			argumentTwo, settedStates2, err = translateExpression(entry.Expression, declaredIdentifiers)
			if err != nil {
				return nil, nil, errors.Wrapf(
					err,
					"unable to translate the key #%d for the hash table definition",
					index,
				)
			}
		}

		argumentThree, settedStates3, err := translateExpression(entry.Value, declaredIdentifiers)
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

	var arguments []expressions.Expression
	settedStates = mapset.NewSet()
	for index, argument := range functionCall.Arguments {
		result, settedStates2, err := translateExpression(argument, declaredIdentifiers)
		if err != nil {
			return nil, nil, errors.Wrapf(
				err,
				"unable to translate the argument #%d for the function %s",
				index,
				functionCall.Name,
			)
		}

		arguments = append(arguments, result)
		settedStates = settedStates.Union(settedStates2)
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
			translateExpression(conditionalCase.Condition, declaredIdentifiers)
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
