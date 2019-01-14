package interpreter

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// ReaderDependencies ...
type ReaderDependencies struct {
	DefaultReader io.Reader
	FileSystem    afero.Fs
}

func readCode(filename string, dependencies ReaderDependencies) (string, error) {
	var reader io.Reader
	if isEmptyFilename(filename) {
		reader = dependencies.DefaultReader
	} else {
		file, err := dependencies.FileSystem.Open(filename)
		if err != nil {
			return "", errors.Wrapf(err, "unable to open the file %s", filename)
		}
		defer file.Close() // nolint: errcheck

		reader = file
	}

	code, err := ioutil.ReadAll(reader)
	if err != nil {
		var message string
		if isEmptyFilename(filename) {
			message = "unable to read the default source"
		} else {
			message = fmt.Sprintf("unable to read the file %s", filename)
		}

		return "", errors.Wrap(err, message)
	}

	return string(code), nil
}

func isEmptyFilename(filename string) bool {
	return len(filename) == 0 || filename == "-"
}
