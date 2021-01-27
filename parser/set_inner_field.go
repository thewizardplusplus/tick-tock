package parser

import (
	"reflect"
)

// SetInnerField ...
//
// It creates a hierarchy of structures and sets the specified field.
//
// First, it searches the specified field in the passed root value. If the field was found, then
// the specified value is set to it.
//
// Otherwise, the first field is selected and filled with a zero value of the corresponding type.
// After that, the search recursively continues in it with the same logic.
//
// The passed root value and the first field on each step of the search should be a pointer
// to a structure.
//
func SetInnerField(rootValue interface{}, fieldName string, fieldValue interface{}) interface{} {
	value := reflect.ValueOf(rootValue).Elem()
	for {
		field := value.FieldByName(fieldName)
		if field.IsValid() {
			field.Set(reflect.ValueOf(fieldValue))
			return rootValue
		}

		fieldIndex := 0
		if valueType := value.Type(); valueType == reflect.TypeOf(Unary{}) {
			fieldIndex = valueType.NumField() - 1
		}

		field = value.Field(fieldIndex)
		field.Set(reflect.New(field.Type().Elem()))

		value = field.Elem()
	}
}
