package interpreter

import (
	"bytes"
	"io"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/tests/mocks"
)

func TestReadCode(test *testing.T) {
	type args struct {
		filename      string
		defaultReader io.Reader
	}

	for _, testData := range []struct {
		name                   string
		args                   args
		initializeDependencies func(fileSystem *mocks.FileSystem, file *mocks.File)
		want                   string
		wantErr                assert.ErrorAssertionFunc
	}{
		{
			name: "success with a default source",
			args: args{"", bytes.NewReader([]byte("test"))},
			initializeDependencies: func(fileSystem *mocks.FileSystem, file *mocks.File) {},
			want:    "test",
			wantErr: assert.NoError,
		},
		{
			name: "success with a file source",
			args: args{"file", nil},
			initializeDependencies: func(fileSystem *mocks.FileSystem, file *mocks.File) {
				fileSystem.On("Open", "file").Return(file, nil)

				file.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(func(buffer []byte) int { return copy(buffer, "content") }, io.EOF)
				file.On("Close").Return(nil)
			},
			want:    "content",
			wantErr: assert.NoError,
		},
		{
			name: "error on a file opening",
			args: args{"file", nil},
			initializeDependencies: func(fileSystem *mocks.FileSystem, file *mocks.File) {
				fileSystem.On("Open", "file").Return(nil, iotest.ErrTimeout)
			},
			wantErr: assert.Error,
		},
		{
			name: "error on a default source reading",
			args: args{"", iotest.TimeoutReader(bytes.NewReader([]byte("test")))},
			initializeDependencies: func(fileSystem *mocks.FileSystem, file *mocks.File) {},
			wantErr:                assert.Error,
		},
		{
			name: "error on a file reading",
			args: args{"file", nil},
			initializeDependencies: func(fileSystem *mocks.FileSystem, file *mocks.File) {
				fileSystem.On("Open", "file").Return(file, nil)

				file.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(0, iotest.ErrTimeout)
				file.On("Close").Return(nil)
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			fileSystem := new(mocks.FileSystem)
			file := new(mocks.File)
			testData.initializeDependencies(fileSystem, file)

			dependencies := ReaderDependencies{testData.args.defaultReader, fileSystem}
			got, err := readCode(testData.args.filename, dependencies)

			mock.AssertExpectationsForObjects(test, fileSystem, file)
			assert.Equal(test, testData.want, got)
			testData.wantErr(test, err)
		})
	}
}

func TestIsEmptyFilename(test *testing.T) {
	type args struct {
		filename string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "usual filename",
			args: args{"test"},
		},
		{
			name: "special filename",
			args: args{"-"},
			want: true,
		},
		{
			name: "empty filename",
			want: true,
		},
	}
	for _, testData := range tests {
		test.Run(testData.name, func(test *testing.T) {
			got := isEmptyFilename(testData.args.filename)
			assert.Equal(test, testData.want, got)
		})
	}
}
