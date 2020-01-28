// The MIT License (MIT)
//
// Copyright (C) 2017 Alec Thomas
// Copyright (C) 2020 thewizardplusplus <thewizardplusplus@yandex.ru>
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
	Addition *Addition `parser:"@@"`
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
	Operation string    `parser:"( @( \"-\" )"`
	Unary     *Unary    `parser:"@@ )"`
	Accessor  *Accessor `parser:"| @@"`
}

// Accessor ...
type Accessor struct {
	Atom *Atom         `parser:"@@"`
	Key  []*Expression `parser:"{ \"[\" @@ \"]\" }"`
}

// Atom ...
type Atom struct {
	Number         *float64        `parser:"@Int | @Float"`
	String         *string         `parser:"| @String | @RawString"`
	ListDefinition *ListDefinition `parser:"| @@"`
	FunctionCall   *FunctionCall   `parser:"| @@"`
	Identifier     *string         `parser:"| @Ident"`
	Expression     *Expression     `parser:"| \"(\" @@ \")\""`
}

// ListDefinition ...
type ListDefinition struct {
	Items []*Expression `parser:"\"[\" [ @@ { \",\" @@ } [ \",\" ] ] \"]\""`
}

// FunctionCall ...
type FunctionCall struct {
	Name      string        `parser:"@Ident"`
	Arguments []*Expression `parser:"\"(\" [ @@ { \",\" @@ } [ \",\" ] ] \")\""`
}
