package runtime

import "github.com/pkg/errors"

func newUnknownStateError(state string) error {
	return errors.Errorf("unknown state %s", state)
}
