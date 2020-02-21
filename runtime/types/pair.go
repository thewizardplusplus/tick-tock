package types

import (
	"unicode/utf8"
)

// Pair ...
type Pair struct {
	Head interface{}
	Tail *Pair
}

// NewPairFromSlice ...
func NewPairFromSlice(items []interface{}) *Pair {
	if len(items) == 0 {
		return nil
	}

	head, tail := items[0], NewPairFromSlice(items[1:])
	return &Pair{head, tail}
}

// NewPairFromString ...
func NewPairFromString(text string) *Pair {
	if len(text) == 0 {
		return nil
	}

	head, headSize := utf8.DecodeRuneInString(text)
	tail := NewPairFromString(text[headSize:])
	return &Pair{head, tail}
}

// Size ...
func (pair *Pair) Size() int {
	if pair == nil {
		return 0
	}

	return 1 + pair.Tail.Size()
}

// Item ...
func (pair *Pair) Item(index int) (item interface{}, ok bool) {
	if pair == nil {
		return nil, false
	}
	if index == 0 {
		return pair.Head, true
	}

	return pair.Tail.Item(index - 1)
}

// Copy ...
func (pair *Pair) Copy() *Pair {
	if pair == nil {
		return nil
	}

	head, tail := pair.Head, pair.Tail.Copy()
	return &Pair{head, tail}
}

// Append ...
func (pair *Pair) Append(anotherPair *Pair) *Pair {
	if pair == nil {
		return anotherPair
	}

	head, tail := pair.Head, pair.Tail.Append(anotherPair)
	return &Pair{head, tail}
}
