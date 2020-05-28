package types

import (
	"unicode/utf8"

	"github.com/pkg/errors"
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

// NewPairFromText ...
func NewPairFromText(text string) *Pair {
	if len(text) == 0 {
		return nil
	}

	head, headSize := utf8.DecodeRuneInString(text)
	tail := NewPairFromText(text[headSize:])
	return &Pair{float64(head), tail}
}

// Size ...
func (pair *Pair) Size() int {
	if pair == nil {
		return 0
	}

	return 1 + pair.Tail.Size()
}

// Equals ...
func (pair *Pair) Equals(sample *Pair) (bool, error) {
	result, err := Equals(pair, sample)
	if err != nil {
		return false, errors.Wrap(err, "unable to compare pairs for equality")
	}

	return result, nil
}

// Compare ...
func (pair *Pair) Compare(sample *Pair) (ComparisonResult, error) {
	if pair == nil || sample == nil {
		if pair == nil && sample != nil {
			return Less, nil
		} else if pair == nil && sample == nil {
			return Equal, nil
		} else if pair != nil && sample == nil {
			return Greater, nil
		}
	}

	result, err := Compare(pair.Head, sample.Head)
	if err != nil {
		return 0, errors.Wrap(err, "unable to compare some items")
	}
	if result != Equal {
		return result, nil
	}

	// heads are equal, continue
	return pair.Tail.Compare(sample.Tail)
}

// Item ...
func (pair *Pair) Item(index float64) (item interface{}, ok bool) {
	if pair == nil {
		return nil, false
	}
	if index == 0 {
		return pair.Head, true
	}

	return pair.Tail.Item(index - 1)
}

// Append ...
func (pair *Pair) Append(anotherPair *Pair) *Pair {
	if pair == nil {
		return anotherPair
	}

	head, tail := pair.Head, pair.Tail.Append(anotherPair)
	return &Pair{head, tail}
}

// Slice ...
func (pair *Pair) Slice() []interface{} {
	if pair == nil {
		return nil
	}

	items := []interface{}{pair.Head}
	items = append(items, pair.Tail.Slice()...)

	return items
}

// DeepSlice ...
func (pair *Pair) DeepSlice() []interface{} {
	if pair == nil {
		return nil
	}

	var head interface{}
	switch typedHead := pair.Head.(type) {
	case Nil:
	case *Pair:
		head = typedHead.DeepSlice()
	default:
		head = typedHead
	}

	items := []interface{}{head}
	items = append(items, pair.Tail.DeepSlice()...)

	return items
}

// Text ...
func (pair *Pair) Text() (string, error) {
	if pair == nil {
		return "", nil
	}

	head, ok := pair.Head.(float64)
	if !ok {
		return "", errors.Errorf(
			"incorrect type of some item for conversion to a string (%T instead float64)",
			pair.Head,
		)
	}

	runeHead := rune(head)
	if !utf8.ValidRune(runeHead) {
		return "", errors.New("incorrect rune in some item")
	}

	tail, err := pair.Tail.Text()
	if err != nil {
		return "", err
	}

	return string(runeHead) + tail, nil
}
