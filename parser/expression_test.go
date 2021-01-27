package parser

import (
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
)

func TestParseToAST_withExpression(test *testing.T) {
	type args struct {
		code string
		ast  interface{}
	}

	for _, testData := range []struct {
		name    string
		args    args
		wantAST interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Atom/number/integer",
			args:    args{"23", new(Atom)},
			wantAST: &Atom{IntegerNumber: pointer.ToInt64(23)},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/number/integer/hexadecimal",
			args:    args{"0x23", new(Atom)},
			wantAST: &Atom{IntegerNumber: pointer.ToInt64(0x23)},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/number/floating-point",
			args:    args{"2.3", new(Atom)},
			wantAST: &Atom{FloatingPointNumber: pointer.ToFloat64(2.3)},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/symbol/latin1",
			args:    args{"'t'", new(Atom)},
			wantAST: &Atom{Symbol: pointer.ToString("t")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/symbol/not latin1",
			args:    args{"'т'", new(Atom)},
			wantAST: &Atom{Symbol: pointer.ToString("т")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/symbol/escape sequence",
			args:    args{`'\n'`, new(Atom)},
			wantAST: &Atom{Symbol: pointer.ToString("\n")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/string/interpreted/single-quoted",
			args:    args{`'line #1\nline #2'`, new(Atom)},
			wantAST: &Atom{String: pointer.ToString("line #1\nline #2")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/string/interpreted/double-quoted",
			args:    args{`"line #1\nline #2"`, new(Atom)},
			wantAST: &Atom{String: pointer.ToString("line #1\nline #2")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/string/raw/single line",
			args:    args{"`line #1\\nline #2`", new(Atom)},
			wantAST: &Atom{String: pointer.ToString("line #1\\nline #2")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/string/raw/few lines",
			args:    args{"`line #1\nline #2`", new(Atom)},
			wantAST: &Atom{String: pointer.ToString("line #1\nline #2")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/identifier",
			args:    args{"test", new(Atom)},
			wantAST: &Atom{Identifier: pointer.ToString("test")},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/list definition/no items",
			args:    args{"[]", new(Atom)},
			wantAST: &Atom{ListDefinition: &ListDefinition{Items: &ExpressionGroup{}}},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/list definition/few items",
			args: args{"[12, 23, 42]", new(Atom)},
			wantAST: &Atom{
				ListDefinition: &ListDefinition{
					Items: &ExpressionGroup{[]*Expression{
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression),
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(23)).(*Expression),
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(42)).(*Expression),
					}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/hash table definition/no items",
			args:    args{"{}", new(Atom)},
			wantAST: &Atom{HashTableDefinition: &HashTableDefinition{Entries: nil}},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/hash table definition/identifier/single entry",
			args: args{"{x: 12}", new(Atom)},
			wantAST: &Atom{
				HashTableDefinition: &HashTableDefinition{
					Entries: []*HashTableEntry{
						{
							Name:  pointer.ToString("x"),
							Value: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression),
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/hash table definition/identifier/single entry/trailing comma",
			args: args{"{x: 12,}", new(Atom)},
			wantAST: &Atom{
				HashTableDefinition: &HashTableDefinition{
					Entries: []*HashTableEntry{
						{
							Name:  pointer.ToString("x"),
							Value: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression),
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/hash table definition/identifier/few entries",
			args: args{"{x: 12, y: 23, z: 42}", new(Atom)},
			wantAST: &Atom{
				HashTableDefinition: &HashTableDefinition{
					Entries: []*HashTableEntry{
						{
							Name:  pointer.ToString("x"),
							Value: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression),
						},
						{
							Name:  pointer.ToString("y"),
							Value: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(23)).(*Expression),
						},
						{
							Name:  pointer.ToString("z"),
							Value: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(42)).(*Expression),
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/hash table definition/identifier/few entries/trailing comma",
			args: args{"{x: 12, y: 23, z: 42,}", new(Atom)},
			wantAST: &Atom{
				HashTableDefinition: &HashTableDefinition{
					Entries: []*HashTableEntry{
						{
							Name:  pointer.ToString("x"),
							Value: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression),
						},
						{
							Name:  pointer.ToString("y"),
							Value: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(23)).(*Expression),
						},
						{
							Name:  pointer.ToString("z"),
							Value: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(42)).(*Expression),
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/hash table definition/expression",
			args: args{"{[test()]: 12}", new(Atom)},
			wantAST: &Atom{
				HashTableDefinition: &HashTableDefinition{
					Entries: []*HashTableEntry{
						{
							Expression: SetInnerField(&Expression{}, "FunctionCall", &FunctionCall{
								Name:      "test",
								Arguments: &ExpressionGroup{},
							}).(*Expression),
							Value: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression),
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/function call/no arguments",
			args:    args{"test()", new(Atom)},
			wantAST: &Atom{FunctionCall: &FunctionCall{Name: "test", Arguments: &ExpressionGroup{}}},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/function call/few arguments",
			args: args{"test(12, 23, 42)", new(Atom)},
			wantAST: &Atom{
				FunctionCall: &FunctionCall{
					Name: "test",
					Arguments: &ExpressionGroup{[]*Expression{
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression),
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(23)).(*Expression),
						SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(42)).(*Expression),
					}},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/conditional expression/single conditional case/nonempty",
			args: args{"when => 12 23 42;", new(Atom)},
			wantAST: &Atom{
				ConditionalExpression: &ConditionalExpression{
					ConditionalCases: []*ConditionalCase{
						{
							Condition: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression),
							Commands: []*Command{
								{
									Expression: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(
										23,
									)).(*Expression),
								},
								{
									Expression: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(
										42,
									)).(*Expression),
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/conditional expression/single conditional case/empty",
			args: args{"when => 12;", new(Atom)},
			wantAST: &Atom{
				ConditionalExpression: &ConditionalExpression{
					ConditionalCases: []*ConditionalCase{
						{Condition: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression)},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/conditional expression/few conditional cases/nonempty",
			args: args{"when => 12 23 42 => 13 24 43 => 14 25 44;", new(Atom)},
			wantAST: &Atom{
				ConditionalExpression: &ConditionalExpression{
					ConditionalCases: []*ConditionalCase{
						{
							Condition: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression),
							Commands: []*Command{
								{
									Expression: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(
										23,
									)).(*Expression),
								},
								{
									Expression: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(
										42,
									)).(*Expression),
								},
							},
						},
						{
							Condition: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(13)).(*Expression),
							Commands: []*Command{
								{
									Expression: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(
										24,
									)).(*Expression),
								},
								{
									Expression: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(
										43,
									)).(*Expression),
								},
							},
						},
						{
							Condition: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(14)).(*Expression),
							Commands: []*Command{
								{
									Expression: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(
										25,
									)).(*Expression),
								},
								{
									Expression: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(
										44,
									)).(*Expression),
								},
							},
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/conditional expression/few conditional cases/empty",
			args: args{"when => 12 => 23 => 42;", new(Atom)},
			wantAST: &Atom{
				ConditionalExpression: &ConditionalExpression{
					ConditionalCases: []*ConditionalCase{
						{Condition: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression)},
						{Condition: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(23)).(*Expression)},
						{Condition: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(42)).(*Expression)},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Atom/conditional expression/without conditional cases",
			args:    args{"when;", new(Atom)},
			wantAST: &Atom{ConditionalExpression: &ConditionalExpression{}},
			wantErr: assert.NoError,
		},
		{
			name: "Atom/expression",
			args: args{"(23)", new(Atom)},
			wantAST: &Atom{
				Expression: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(23)).(*Expression),
			},
			wantErr: assert.NoError,
		},
		{
			name: "Accessor/nonempty/brackets",
			args: args{"test[12][23]", new(Accessor)},
			wantAST: &Accessor{
				Atom: &Atom{Identifier: pointer.ToString("test")},
				Keys: []*AccessorKey{
					{Expression: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression)},
					{Expression: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(23)).(*Expression)},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Accessor/nonempty/identifiers",
			args: args{"test.one.two", new(Accessor)},
			wantAST: &Accessor{
				Atom: &Atom{Identifier: pointer.ToString("test")},
				Keys: []*AccessorKey{{Name: pointer.ToString("one")}, {Name: pointer.ToString("two")}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Accessor/nonempty/identifiers and brackets",
			args: args{"test.one[12].two[23]", new(Accessor)},
			wantAST: &Accessor{
				Atom: &Atom{Identifier: pointer.ToString("test")},
				Keys: []*AccessorKey{
					{Name: pointer.ToString("one")},
					{Expression: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(12)).(*Expression)},
					{Name: pointer.ToString("two")},
					{Expression: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(23)).(*Expression)},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Accessor/empty",
			args:    args{"23", new(Accessor)},
			wantAST: &Accessor{Atom: &Atom{IntegerNumber: pointer.ToInt64(23)}},
			wantErr: assert.NoError,
		},
		{
			name: "Unary/nonempty",
			args: args{"-~!23", new(Unary)},
			wantAST: &Unary{
				Operation: "-",
				Unary: &Unary{
					Operation: "~",
					Unary: &Unary{
						Operation: "!",
						Unary:     SetInnerField(&Unary{}, "IntegerNumber", pointer.ToInt64(23)).(*Unary),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Unary/empty",
			args:    args{"23", new(Unary)},
			wantAST: SetInnerField(&Unary{}, "IntegerNumber", pointer.ToInt64(23)).(*Unary),
			wantErr: assert.NoError,
		},
		{
			name: "Multiplication/nonempty",
			args: args{"5 * 12 / 23 % 42", new(Multiplication)},
			wantAST: &Multiplication{
				Unary:     SetInnerField(&Unary{}, "IntegerNumber", pointer.ToInt64(5)).(*Unary),
				Operation: "*",
				Multiplication: &Multiplication{
					Unary:     SetInnerField(&Unary{}, "IntegerNumber", pointer.ToInt64(12)).(*Unary),
					Operation: "/",
					Multiplication: &Multiplication{
						Unary:     SetInnerField(&Unary{}, "IntegerNumber", pointer.ToInt64(23)).(*Unary),
						Operation: "%",
						Multiplication: &Multiplication{
							Unary: SetInnerField(&Unary{}, "IntegerNumber", pointer.ToInt64(42)).(*Unary),
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Multiplication/empty",
			args: args{"23", new(Multiplication)},
			wantAST: SetInnerField(&Multiplication{}, "IntegerNumber", pointer.ToInt64(
				23,
			)).(*Multiplication),
			wantErr: assert.NoError,
		},
		{
			name: "Addition/nonempty",
			args: args{"12 + 23 - 42", new(Addition)},
			wantAST: &Addition{
				Multiplication: SetInnerField(&Multiplication{}, "IntegerNumber", pointer.ToInt64(
					12,
				)).(*Multiplication),
				Operation: "+",
				Addition: &Addition{
					Multiplication: SetInnerField(&Multiplication{}, "IntegerNumber", pointer.ToInt64(
						23,
					)).(*Multiplication),
					Operation: "-",
					Addition: &Addition{
						Multiplication: SetInnerField(&Multiplication{}, "IntegerNumber", pointer.ToInt64(
							42,
						)).(*Multiplication),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Addition/empty",
			args:    args{"23", new(Addition)},
			wantAST: SetInnerField(&Addition{}, "IntegerNumber", pointer.ToInt64(23)).(*Addition),
			wantErr: assert.NoError,
		},
		{
			name: "Shift/nonempty",
			args: args{"5 << 12 >> 23 >>> 42", new(Shift)},
			wantAST: &Shift{
				Addition:  SetInnerField(&Addition{}, "IntegerNumber", pointer.ToInt64(5)).(*Addition),
				Operation: "<<",
				Shift: &Shift{
					Addition:  SetInnerField(&Addition{}, "IntegerNumber", pointer.ToInt64(12)).(*Addition),
					Operation: ">>",
					Shift: &Shift{
						Addition:  SetInnerField(&Addition{}, "IntegerNumber", pointer.ToInt64(23)).(*Addition),
						Operation: ">>>",
						Shift: &Shift{
							Addition: SetInnerField(&Addition{}, "IntegerNumber", pointer.ToInt64(42)).(*Addition),
						},
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Shift/empty",
			args:    args{"23", new(Shift)},
			wantAST: SetInnerField(&Shift{}, "IntegerNumber", pointer.ToInt64(23)).(*Shift),
			wantErr: assert.NoError,
		},
		{
			name: "BitwiseConjunction/nonempty",
			args: args{"12 & 23 & 42", new(BitwiseConjunction)},
			wantAST: &BitwiseConjunction{
				Shift:     SetInnerField(&Shift{}, "IntegerNumber", pointer.ToInt64(12)).(*Shift),
				Operation: "&",
				BitwiseConjunction: &BitwiseConjunction{
					Shift:     SetInnerField(&Shift{}, "IntegerNumber", pointer.ToInt64(23)).(*Shift),
					Operation: "&",
					BitwiseConjunction: &BitwiseConjunction{
						Shift: SetInnerField(&Shift{}, "IntegerNumber", pointer.ToInt64(42)).(*Shift),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "BitwiseConjunction/empty",
			args: args{"23", new(BitwiseConjunction)},
			wantAST: SetInnerField(&BitwiseConjunction{}, "IntegerNumber", pointer.ToInt64(
				23,
			)).(*BitwiseConjunction),
			wantErr: assert.NoError,
		},
		{
			name: "BitwiseExclusiveDisjunction/nonempty",
			args: args{"12 ^ 23 ^ 42", new(BitwiseExclusiveDisjunction)},
			wantAST: &BitwiseExclusiveDisjunction{
				BitwiseConjunction: SetInnerField(&BitwiseConjunction{}, "IntegerNumber", pointer.ToInt64(
					12,
				)).(*BitwiseConjunction),
				Operation: "^",
				BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
					BitwiseConjunction: SetInnerField(&BitwiseConjunction{}, "IntegerNumber", pointer.ToInt64(
						23,
					)).(*BitwiseConjunction),
					Operation: "^",
					BitwiseExclusiveDisjunction: &BitwiseExclusiveDisjunction{
						BitwiseConjunction: SetInnerField(&BitwiseConjunction{}, "IntegerNumber", pointer.ToInt64(
							42,
						)).(*BitwiseConjunction),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "BitwiseExclusiveDisjunction/empty",
			args: args{"23", new(BitwiseExclusiveDisjunction)},
			wantAST: SetInnerField(&BitwiseExclusiveDisjunction{}, "IntegerNumber", pointer.ToInt64(
				23,
			)).(*BitwiseExclusiveDisjunction),
			wantErr: assert.NoError,
		},
		{
			name: "BitwiseDisjunction/nonempty",
			args: args{"12 | 23 | 42", new(BitwiseDisjunction)},
			wantAST: &BitwiseDisjunction{
				BitwiseExclusiveDisjunction: SetInnerField(&BitwiseExclusiveDisjunction{}, "IntegerNumber", pointer.ToInt64(
					12,
				)).(*BitwiseExclusiveDisjunction),
				Operation: "|",
				BitwiseDisjunction: &BitwiseDisjunction{
					BitwiseExclusiveDisjunction: SetInnerField(&BitwiseExclusiveDisjunction{}, "IntegerNumber", pointer.ToInt64(
						23,
					)).(*BitwiseExclusiveDisjunction),
					Operation: "|",
					BitwiseDisjunction: &BitwiseDisjunction{
						BitwiseExclusiveDisjunction: SetInnerField(&BitwiseExclusiveDisjunction{}, "IntegerNumber", pointer.ToInt64(
							42,
						)).(*BitwiseExclusiveDisjunction),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "BitwiseDisjunction/empty",
			args: args{"23", new(BitwiseDisjunction)},
			wantAST: SetInnerField(&BitwiseDisjunction{}, "IntegerNumber", pointer.ToInt64(
				23,
			)).(*BitwiseDisjunction),
			wantErr: assert.NoError,
		},
		{
			name: "Comparison/nonempty/less",
			args: args{"12 < 23 <= 42", new(Comparison)},
			wantAST: &Comparison{
				BitwiseDisjunction: SetInnerField(&BitwiseDisjunction{}, "IntegerNumber", pointer.ToInt64(
					12,
				)).(*BitwiseDisjunction),
				Operation: "<",
				Comparison: &Comparison{
					BitwiseDisjunction: SetInnerField(&BitwiseDisjunction{}, "IntegerNumber", pointer.ToInt64(
						23,
					)).(*BitwiseDisjunction),
					Operation: "<=",
					Comparison: &Comparison{
						BitwiseDisjunction: SetInnerField(&BitwiseDisjunction{}, "IntegerNumber", pointer.ToInt64(
							42,
						)).(*BitwiseDisjunction),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "Comparison/nonempty/great",
			args: args{"12 > 23 >= 42", new(Comparison)},
			wantAST: &Comparison{
				BitwiseDisjunction: SetInnerField(&BitwiseDisjunction{}, "IntegerNumber", pointer.ToInt64(
					12,
				)).(*BitwiseDisjunction),
				Operation: ">",
				Comparison: &Comparison{
					BitwiseDisjunction: SetInnerField(&BitwiseDisjunction{}, "IntegerNumber", pointer.ToInt64(
						23,
					)).(*BitwiseDisjunction),
					Operation: ">=",
					Comparison: &Comparison{
						BitwiseDisjunction: SetInnerField(&BitwiseDisjunction{}, "IntegerNumber", pointer.ToInt64(
							42,
						)).(*BitwiseDisjunction),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Comparison/empty",
			args:    args{"23", new(Comparison)},
			wantAST: SetInnerField(&Comparison{}, "IntegerNumber", pointer.ToInt64(23)).(*Comparison),
			wantErr: assert.NoError,
		},
		{
			name: "Equality/nonempty",
			args: args{"12 == 23 != 42", new(Equality)},
			wantAST: &Equality{
				Comparison: SetInnerField(&Comparison{}, "IntegerNumber", pointer.ToInt64(12)).(*Comparison),
				Operation:  "==",
				Equality: &Equality{
					Comparison: SetInnerField(&Comparison{}, "IntegerNumber", pointer.ToInt64(23)).(*Comparison),
					Operation:  "!=",
					Equality: &Equality{
						Comparison: SetInnerField(&Comparison{}, "IntegerNumber", pointer.ToInt64(42)).(*Comparison),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Equality/empty",
			args:    args{"23", new(Equality)},
			wantAST: SetInnerField(&Equality{}, "IntegerNumber", pointer.ToInt64(23)).(*Equality),
			wantErr: assert.NoError,
		},
		{
			name: "Conjunction/nonempty",
			args: args{"12 && 23 && 42", new(Conjunction)},
			wantAST: &Conjunction{
				Equality:  SetInnerField(&Equality{}, "IntegerNumber", pointer.ToInt64(12)).(*Equality),
				Operation: "&&",
				Conjunction: &Conjunction{
					Equality:  SetInnerField(&Equality{}, "IntegerNumber", pointer.ToInt64(23)).(*Equality),
					Operation: "&&",
					Conjunction: &Conjunction{
						Equality: SetInnerField(&Equality{}, "IntegerNumber", pointer.ToInt64(42)).(*Equality),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Conjunction/empty",
			args:    args{"23", new(Conjunction)},
			wantAST: SetInnerField(&Conjunction{}, "IntegerNumber", pointer.ToInt64(23)).(*Conjunction),
			wantErr: assert.NoError,
		},
		{
			name: "Disjunction/nonempty",
			args: args{"12 || 23 || 42", new(Disjunction)},
			wantAST: &Disjunction{
				Conjunction: SetInnerField(&Conjunction{}, "IntegerNumber", pointer.ToInt64(
					12,
				)).(*Conjunction),
				Operation: "||",
				Disjunction: &Disjunction{
					Conjunction: SetInnerField(&Conjunction{}, "IntegerNumber", pointer.ToInt64(
						23,
					)).(*Conjunction),
					Operation: "||",
					Disjunction: &Disjunction{
						Conjunction: SetInnerField(&Conjunction{}, "IntegerNumber", pointer.ToInt64(
							42,
						)).(*Conjunction),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "Disjunction/empty",
			args:    args{"23", new(Disjunction)},
			wantAST: SetInnerField(&Disjunction{}, "IntegerNumber", pointer.ToInt64(23)).(*Disjunction),
			wantErr: assert.NoError,
		},
		{
			name: "NilCoalescing/nonempty",
			args: args{"12 ?? 23 ?? 42", new(NilCoalescing)},
			wantAST: &NilCoalescing{
				Disjunction: SetInnerField(&Disjunction{}, "IntegerNumber", pointer.ToInt64(
					12,
				)).(*Disjunction),
				Operation: "??",
				NilCoalescing: &NilCoalescing{
					Disjunction: SetInnerField(&Disjunction{}, "IntegerNumber", pointer.ToInt64(
						23,
					)).(*Disjunction),
					Operation: "??",
					NilCoalescing: &NilCoalescing{
						Disjunction: SetInnerField(&Disjunction{}, "IntegerNumber", pointer.ToInt64(
							42,
						)).(*Disjunction),
					},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "NilCoalescing/empty",
			args:    args{"23", new(NilCoalescing)},
			wantAST: SetInnerField(&NilCoalescing{}, "IntegerNumber", pointer.ToInt64(23)).(*NilCoalescing),
			wantErr: assert.NoError,
		},
		{
			name: "ListConstruction/nonempty",
			args: args{"5 : 12 : [23, 42]", new(ListConstruction)},
			wantAST: &ListConstruction{
				NilCoalescing: SetInnerField(&NilCoalescing{}, "IntegerNumber", pointer.ToInt64(
					5,
				)).(*NilCoalescing),
				Operation: ":",
				ListConstruction: &ListConstruction{
					NilCoalescing: SetInnerField(&NilCoalescing{}, "IntegerNumber", pointer.ToInt64(
						12,
					)).(*NilCoalescing),
					Operation: ":",
					ListConstruction: SetInnerField(&ListConstruction{}, "ListDefinition", &ListDefinition{
						Items: &ExpressionGroup{[]*Expression{
							SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(23)).(*Expression),
							SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(42)).(*Expression),
						}},
					}).(*ListConstruction),
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "ListConstruction/empty",
			args: args{"23", new(ListConstruction)},
			wantAST: SetInnerField(&ListConstruction{}, "IntegerNumber", pointer.ToInt64(
				23,
			)).(*ListConstruction),
			wantErr: assert.NoError,
		},
		{
			name:    "Expression",
			args:    args{"23", new(Expression)},
			wantAST: SetInnerField(&Expression{}, "IntegerNumber", pointer.ToInt64(23)).(*Expression),
			wantErr: assert.NoError,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			err := ParseToAST(testData.args.code, testData.args.ast)

			assert.Equal(test, testData.wantAST, testData.args.ast)
			testData.wantErr(test, err)
		})
	}
}
