package glogger

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var (
	Int    = reflect.TypeOf(0)
	String = reflect.TypeOf("")
)

type TestStruct struct {
	value string
}

func TestFieldMatchesSchema(t *testing.T) {
	propname := "value"
	teststruct := &TestStruct{}
	assert.Equal(
		t,
		true,
		fieldMatchesSchema(propname, String, *teststruct),
		"should return true when provided type matches expected type in schema",
	)
	assert.NotEqual(
		t,
		true,
		fieldMatchesSchema(propname, Int, *teststruct),
		"should return false when provided type does not match expected type in schema",
	)
}
