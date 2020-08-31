package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNil_MarshalText(test *testing.T) {
	// it's an example of an implicit call of the types.Nil.MarshalText() method;
	// you also can use json.Encoder with its method SetEscapeHTML() to avoid HTML escaping
	gotBytes, gotErr := json.Marshal(Nil{})

	assert.Equal(test, []byte(`"\u003cnil\u003e"`), gotBytes)
	assert.NoError(test, gotErr)
}
