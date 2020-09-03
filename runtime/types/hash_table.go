package types

import (
	stderrors "errors"

	"github.com/pkg/errors"
)

// ...
var (
	ErrNotFound = stderrors.New("not found")
)

// HashTable ...
type HashTable map[interface{}]interface{}

// Size ...
func (table HashTable) Size() int {
	return len(table)
}

// Keys ...
func (table HashTable) Keys() []interface{} {
	var keys []interface{}
	for key := range table {
		if keyAsString, ok := key.(string); ok {
			key = NewPairFromText(keyAsString)
		}

		keys = append(keys, key)
	}

	return keys
}

// Equals ...
func (table HashTable) Equals(sample HashTable) (bool, error) {
	if len(table) != len(sample) {
		return false, nil
	}

	for key, tableValue := range table {
		sampleValue, ok := sample[key]
		if !ok {
			return false, nil
		}

		equals, err := Equals(tableValue, sampleValue)
		if err != nil {
			return false, errors.Wrap(err, "unable to compare some values for equality")
		}
		if !equals {
			return false, nil
		}
	}

	return true, nil
}

// Item ...
func (table HashTable) Item(key interface{}) (interface{}, error) {
	preparedKey, err := prepareKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "unable to prepare the key")
	}

	value, ok := table[preparedKey]
	if !ok {
		return nil, ErrNotFound
	}

	return value, nil
}

// Copy ...
func (table HashTable) Copy() HashTable {
	copiedTable := make(HashTable)
	for key, value := range table {
		copiedTable[key] = value
	}

	return copiedTable
}

// With ...
func (table HashTable) With(key interface{}, value interface{}) (HashTable, error) {
	preparedKey, err := prepareKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "unable to prepare the key")
	}

	updatedTable := table.Copy()
	if value != (Nil{}) {
		updatedTable[preparedKey] = value
	} else {
		delete(updatedTable, preparedKey)
	}

	return updatedTable, nil
}

// Merge ...
func (table HashTable) Merge(anotherTable HashTable) HashTable {
	unionTable := table.Copy()
	for key, value := range anotherTable {
		unionTable[key] = value
	}

	return unionTable
}

// DeepMap ...
func (table HashTable) DeepMap() (HashTable, error) {
	result := make(HashTable)
	for key, value := range table {
		var err error
		switch typedValue := value.(type) {
		case *Pair:
			value, err = typedValue.DeepSlice()
			if err != nil {
				return nil, errors.Wrap(err, "unable to get the deep list")
			}
		case HashTable:
			value, err = typedValue.DeepMap()
			if err != nil {
				return nil, errors.Wrap(err, "unable to get the deep hash table")
			}
		}

		result[key] = value
	}

	return result, nil
}

func prepareKey(key interface{}) (interface{}, error) {
	switch typedKey := key.(type) {
	case Nil, float64:
		return typedKey, nil
	case *Pair:
		keyAsString, err := typedKey.Text()
		if err != nil {
			return nil, errors.Wrap(err, "unable to convert the key to a string")
		}

		return keyAsString, nil
	default:
		return nil, errors.Errorf("unsupported type %T of the key", key)
	}
}
