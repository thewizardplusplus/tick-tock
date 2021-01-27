package parser

import (
	"reflect"
)

// SetInnerField ...
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
