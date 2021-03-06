// The MIT License (MIT)
//
// Copyright (C) 2017 Alec Thomas
// Copyright (C) 2020-2021 thewizardplusplus <thewizardplusplus@yandex.ru>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is furnished to do
// so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package parser

// Expression ...
type Expression struct {
	ListConstruction *ListConstruction `parser:"@@"`
}

// ListConstruction ...
type ListConstruction struct {
	NilCoalescing    *NilCoalescing    `parser:"@@"`
	Operation        string            `parser:"[ @\":\""`
	ListConstruction *ListConstruction `parser:"@@ ]"`
}

// NilCoalescing ...
type NilCoalescing struct {
	Disjunction   *Disjunction   `parser:"@@"`
	Operation     string         `parser:"[ @( \"?\" \"?\" )"`
	NilCoalescing *NilCoalescing `parser:"@@ ]"`
}

// Disjunction ...
type Disjunction struct {
	Conjunction *Conjunction `parser:"@@"`
	Operation   string       `parser:"[ @( \"|\" \"|\" )"`
	Disjunction *Disjunction `parser:"@@ ]"`
}

// Conjunction ...
type Conjunction struct {
	Equality    *Equality    `parser:"@@"`
	Operation   string       `parser:"[ @( \"&\" \"&\" )"`
	Conjunction *Conjunction `parser:"@@ ]"`
}

// Equality ...
type Equality struct {
	Comparison *Comparison `parser:"@@"`
	Operation  string      `parser:"[ @( \"=\" \"=\" | \"!\" \"=\" )"`
	Equality   *Equality   `parser:"@@ ]"`
}

// Comparison ...
type Comparison struct {
	BitwiseDisjunction *BitwiseDisjunction `parser:"@@"`
	Operation          string              `parser:"[ @( \"<\" \"=\" | \"<\" | \">\" \"=\" | \">\" )"`
	Comparison         *Comparison         `parser:"@@ ]"`
}

// BitwiseDisjunction ...
type BitwiseDisjunction struct {
	BitwiseExclusiveDisjunction *BitwiseExclusiveDisjunction `parser:"@@"`
	Operation                   string                       `parser:"[ @\"|\""`
	BitwiseDisjunction          *BitwiseDisjunction          `parser:"@@ ]"`
}

// BitwiseExclusiveDisjunction ...
type BitwiseExclusiveDisjunction struct {
	BitwiseConjunction          *BitwiseConjunction          `parser:"@@"`
	Operation                   string                       `parser:"[ @\"^\""`
	BitwiseExclusiveDisjunction *BitwiseExclusiveDisjunction `parser:"@@ ]"`
}

// BitwiseConjunction ...
type BitwiseConjunction struct {
	Shift              *Shift              `parser:"@@"`
	Operation          string              `parser:"[ @\"&\""`
	BitwiseConjunction *BitwiseConjunction `parser:"@@ ]"`
}

// Shift ...
type Shift struct {
	Addition  *Addition `parser:"@@"`
	Operation string    `parser:"[ @( \"<\" \"<\" | \">\" \">\" [ \">\" ] )"`
	Shift     *Shift    `parser:"@@ ]"`
}

// Addition ...
type Addition struct {
	Multiplication *Multiplication `parser:"@@"`
	Operation      string          `parser:"[ @( \"+\" | \"-\" )"`
	Addition       *Addition       `parser:"@@ ]"`
}

// Multiplication ...
type Multiplication struct {
	Unary          *Unary          `parser:"@@"`
	Operation      string          `parser:"[ @( \"*\" | \"/\" | \"%\" )"`
	Multiplication *Multiplication `parser:"@@ ]"`
}

// Unary ...
type Unary struct {
	Operation string    `parser:"( @( \"-\" | \"~\" | \"!\" )"`
	Unary     *Unary    `parser:"@@ )"`
	Accessor  *Accessor `parser:"| @@"`
}

// Accessor ...
type Accessor struct {
	Atom *Atom          `parser:"@@"`
	Keys []*AccessorKey `parser:"{ @@ }"`
}

// AccessorKey ...
type AccessorKey struct {
	Name       *string     `parser:"\".\" @Ident"`
	Expression *Expression `parser:"| \"[\" @@ \"]\""`
}

// Atom ...
type Atom struct {
	IntegerNumber         *int64                 `parser:"@Int"`
	FloatingPointNumber   *float64               `parser:"| @Float"`
	Symbol                *string                `parser:"| @Char"`
	String                *string                `parser:"| @String | @RawString"`
	ListDefinition        *ListDefinition        `parser:"| @@"`
	HashTableDefinition   *HashTableDefinition   `parser:"| @@"`
	FunctionCall          *FunctionCall          `parser:"| @@"`
	ConditionalExpression *ConditionalExpression `parser:"| @@"`
	Identifier            *string                `parser:"| @Ident"`
	Expression            *Expression            `parser:"| \"(\" @@ \")\""`
}

// ListDefinition ...
type ListDefinition struct {
	Items *ExpressionGroup `parser:"\"[\" @@ \"]\""`
}

// HashTableDefinition ...
type HashTableDefinition struct {
	Entries []*HashTableEntry `parser:"\"{\" [ @@ { \",\" @@ } [ \",\" ] ] \"}\""`
}

// HashTableEntry ...
type HashTableEntry struct {
	Name       *string     `parser:"( @Ident"`
	Expression *Expression `parser:"| \"[\" @@ \"]\" )"`
	Value      *Expression `parser:"\":\" @@"`
}

// FunctionCall ...
type FunctionCall struct {
	Name      string           `parser:"@Ident"`
	Arguments *ExpressionGroup `parser:"\"(\" @@ \")\""`
}

// ConditionalExpression ...
type ConditionalExpression struct {
	ConditionalCases []*ConditionalCase `parser:"\"when\" { @@ } \";\""`
}

// ConditionalCase ...
type ConditionalCase struct {
	Condition *Expression `parser:"\"=\" \">\" @@"`
	Commands  []*Command  `parser:"{ @@ }"`
}
