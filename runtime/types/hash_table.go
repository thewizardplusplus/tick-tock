package types

import (
	"github.com/pkg/errors"
)

// HashTable ...
type HashTable map[interface{}]interface{}

// Size ...
func (table HashTable) Size() int {
	return len(table)
}

// Get ...
func (table HashTable) Get(key interface{}) (interface{}, error) {
	preparedKey, err := prepareKey(key)
	if err != nil {
		return nil, errors.Wrap(err, "unable to prepare the key for the hash table")
	}

	value, ok := table[preparedKey]
	if !ok {
		value = Nil{}
	}
	return value, nil
}

// Set ...
func (table HashTable) Set(key interface{}, value interface{}) error {
	preparedKey, err := prepareKey(key)
	if err != nil {
		return errors.Wrap(err, "unable to prepare the key for the hash table")
	}

	if value != (Nil{}) {
		table[preparedKey] = value
	} else {
		delete(table, preparedKey)
	}

	return nil
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
		return nil, errors.Errorf("unsupported type %T of the key for the hash table", key)
	}
}
