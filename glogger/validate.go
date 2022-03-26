package glogger

import (
	"reflect"
)

func fieldMatchesSchema(name string, typeof reflect.Type, schema interface{}) bool {
	t := reflect.TypeOf(schema)
	field, _ := t.FieldByName(name)
	return typeof == field.Type
}
