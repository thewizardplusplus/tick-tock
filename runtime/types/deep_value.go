package types

import (
	"github.com/pkg/errors"
)

// GetDeepValue ...
func GetDeepValue(value interface{}) (interface{}, error) {
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

	return value, nil
}
