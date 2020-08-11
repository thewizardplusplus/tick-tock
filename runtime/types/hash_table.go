package types

import (
	"github.com/pkg/errors"
)

// HashTable ...
type HashTable map[interface{}]interface{}

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
