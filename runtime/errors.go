package runtime

import (
	"fmt"

	"github.com/pkg/errors"
)

type unknownStateError string

func newUnknownStateError(state string) error {
	return errors.WithStack(unknownStateError(state))
}

// TODO: remove special processing of an empty state name.
// TODO: move formatting logic to the newUnknownStateError() function.
func (err unknownStateError) Error() string {
	state := string(err)
	if len(state) == 0 {
		state = "<empty>"
	}
	return fmt.Sprintf("unknown state %s", state)
}
