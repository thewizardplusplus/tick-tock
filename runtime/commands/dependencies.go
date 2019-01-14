package commands

import "io"

// Dependencies ...
type Dependencies struct {
	OutWriter io.Writer
	Sleep     SleepDependencies
}
