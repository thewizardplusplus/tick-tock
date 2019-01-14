package interpreter

import (
	"io"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/tests/mocks"
)

func TestReadCode(test *testing.T) {
	type args struct {
		filename string
	}

	for _, testData := range []struct {
		name                   string
		args                   args
		initializeDependencies func(
			defaultReader *mocks.Reader,
			fileSystem *mocks.FileSystem,
			file *mocks.File,
		)
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success with a default source",
			initializeDependencies: func(defaultReader *mocks.Reader, _ *mocks.FileSystem, _ *mocks.File) {
				defaultReader.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(func(buffer []byte) int { return copy(buffer, "test") }, io.EOF)
			},
			want:    "test",
			wantErr: assert.NoError,
		},
		{
			name: "success with a file source",
			args: args{"file"},
			initializeDependencies: func(_ *mocks.Reader, fileSystem *mocks.FileSystem, file *mocks.File) {
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
			args: args{"file"},
			initializeDependencies: func(_ *mocks.Reader, fileSystem *mocks.FileSystem, _ *mocks.File) {
				fileSystem.On("Open", "file").Return(nil, iotest.ErrTimeout)
			},
			wantErr: assert.Error,
		},
		{
			name: "error on a default source reading",
			initializeDependencies: func(defaultReader *mocks.Reader, _ *mocks.FileSystem, _ *mocks.File) {
				defaultReader.On("Read", mock.AnythingOfType("[]uint8")).Return(0, iotest.ErrTimeout)
			},
			wantErr: assert.Error,
		},
		{
			name: "error on a file reading",
			args: args{"file"},
			initializeDependencies: func(_ *mocks.Reader, fileSystem *mocks.FileSystem, file *mocks.File) {
				fileSystem.On("Open", "file").Return(file, nil)

				file.On("Read", mock.AnythingOfType("[]uint8")).Return(0, iotest.ErrTimeout)
				file.On("Close").Return(nil)
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			defaultReader, fileSystem, file := new(mocks.Reader), new(mocks.FileSystem), new(mocks.File)
			testData.initializeDependencies(defaultReader, fileSystem, file)

			got, err := readCode(testData.args.filename, ReaderDependencies{defaultReader, fileSystem})

			mock.AssertExpectationsForObjects(test, defaultReader, fileSystem, file)
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
