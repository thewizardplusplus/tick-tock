package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNil_MarshalJSON(test *testing.T) {
	// it's an example of an implicit call of the types.Nil.MarshalJSON() method
	gotBytes, gotErr := json.Marshal(Nil{})

	assert.Equal(test, []byte("null"), gotBytes)
	assert.NoError(test, gotErr)
}
