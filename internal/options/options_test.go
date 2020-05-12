package options

import (
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/internal/test-utils/mocks"
	"github.com/thewizardplusplus/tick-tock/interpreter"
	"github.com/thewizardplusplus/tick-tock/translator"
)

func TestParse(test *testing.T) {
	type args struct {
		args []string
	}

	const executablePath = "path/to/an/executable/file"
	const versionUsage = Version + "\n"
	const helpUsage = `usage: file [<flags>] [<filename>]

Flags:
  -h, --help      Show context-sensitive help (also try --help-long and
                  --help-man).
  -v, --version   Show application version.
  -i, --inbox=10  Inbox buffer size.
  -s, --state="__initialization__"  ` + `
                  Initial state.
  -m, --message="__initialize__"  ` + `
                  Initial message.

Args:
  [<filename>]  Source file name. Empty or "-" means stdin.

`
	defaultOptions := interpreter.Options{
		InitialMessage: DefaultInitialMessage,
		Translator: translator.Options{
			InboxSize:    DefaultInboxSize,
			InitialState: DefaultInitialState,
		},
	}
	for _, testData := range []struct {
		name                   string
		args                   args
		initializeDependencies func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter)
		wantUsage              []byte
		want                   interpreter.Options
		wantErr                assert.ErrorAssertionFunc
	}{
		{
			name:                   "success without flags and arguments",
			args:                   args{[]string{executablePath}},
			initializeDependencies: func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {},
			want:                   defaultOptions,
			wantErr:                assert.NoError,
		},
		{
			name:                   "success with the -v flag",
			args:                   args{[]string{executablePath, "-v"}},
			initializeDependencies: initializeForUsage,
			wantUsage:              []byte(versionUsage),
			want:                   defaultOptions,
			wantErr:                assert.NoError,
		},
		{
			name:                   "success with the --version flag",
			args:                   args{[]string{executablePath, "--version"}},
			initializeDependencies: initializeForUsage,
			wantUsage:              []byte(versionUsage),
			want:                   defaultOptions,
			wantErr:                assert.NoError,
		},
		{
			name:                   "success with the -h flag",
			args:                   args{[]string{executablePath, "-h"}},
			initializeDependencies: initializeForUsage,
			wantUsage:              []byte(helpUsage),
			want:                   interpreter.Options{},
			wantErr:                assert.NoError,
		},
		{
			name:                   "success with the --help flag",
			args:                   args{[]string{executablePath, "--help"}},
			initializeDependencies: initializeForUsage,
			wantUsage:              []byte(helpUsage),
			want:                   interpreter.Options{},
			wantErr:                assert.NoError,
		},
		{
			name:                   "success with the -i flag",
			args:                   args{[]string{executablePath, "-i", "1000"}},
			initializeDependencies: func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {},
			want:                   setOption(defaultOptions, "Translator.InboxSize", 1000),
			wantErr:                assert.NoError,
		},
		{
			name:                   "success with the --inbox flag",
			args:                   args{[]string{executablePath, "--inbox", "1000"}},
			initializeDependencies: func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {},
			want:                   setOption(defaultOptions, "Translator.InboxSize", 1000),
			wantErr:                assert.NoError,
		},
		{
			name:                   "success with the -s flag",
			args:                   args{[]string{executablePath, "-s", "test"}},
			initializeDependencies: func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {},
			want:                   setOption(defaultOptions, "Translator.InitialState", "test"),
			wantErr:                assert.NoError,
		},
		{
			name:                   "success with the --state flag",
			args:                   args{[]string{executablePath, "--state", "test"}},
			initializeDependencies: func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {},
			want:                   setOption(defaultOptions, "Translator.InitialState", "test"),
			wantErr:                assert.NoError,
		},
		{
			name:                   "success with the -m flag",
			args:                   args{[]string{executablePath, "-m", "test"}},
			initializeDependencies: func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {},
			want:                   setOption(defaultOptions, "InitialMessage", "test"),
			wantErr:                assert.NoError,
		},
		{
			name:                   "success with the --message flag",
			args:                   args{[]string{executablePath, "--message", "test"}},
			initializeDependencies: func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {},
			want:                   setOption(defaultOptions, "InitialMessage", "test"),
			wantErr:                assert.NoError,
		},
		{
			name:                   "success with the filename argument",
			args:                   args{[]string{executablePath, "test"}},
			initializeDependencies: func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {},
			want:                   setOption(defaultOptions, "Filename", "test"),
			wantErr:                assert.NoError,
		},
		{
			name:                   "error with an unknown flag",
			args:                   args{[]string{executablePath, "--unknown"}},
			initializeDependencies: func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {},
			want:                   interpreter.Options{},
			wantErr:                assert.Error,
		},
		{
			name:                   "error with the --inbox flag (missed argument)",
			args:                   args{[]string{executablePath, "--inbox"}},
			initializeDependencies: func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {},
			want:                   interpreter.Options{},
			wantErr:                assert.Error,
		},
		{
			name:                   "error with the --inbox flag (incorrect type)",
			args:                   args{[]string{executablePath, "--inbox", "test"}},
			initializeDependencies: func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {},
			want:                   interpreter.Options{},
			wantErr:                assert.Error,
		},
		{
			name:                   "error with the --state flag (missed argument)",
			args:                   args{[]string{executablePath, "--state"}},
			initializeDependencies: func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {},
			want:                   interpreter.Options{},
			wantErr:                assert.Error,
		},
		{
			name:                   "error with the --message flag (missed argument)",
			args:                   args{[]string{executablePath, "--message"}},
			initializeDependencies: func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {},
			want:                   interpreter.Options{},
			wantErr:                assert.Error,
		},
		{
			name:                   "error with an extra argument",
			args:                   args{[]string{executablePath, "one", "two"}},
			initializeDependencies: func(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {},
			want:                   interpreter.Options{},
			wantErr:                assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			var usage []byte
			usageWriter := new(mocks.Writer)
			errorWriter := new(mocks.Writer)
			exiter := new(mocks.Exiter)
			testData.initializeDependencies(&usage, usageWriter, exiter)

			dependencies := Dependencies{usageWriter, errorWriter, exiter.Exit}
			got, err := Parse(testData.args.args, dependencies)

			mock.AssertExpectationsForObjects(test, usageWriter, errorWriter, exiter)
			assert.Equal(test, testData.wantUsage, usage)
			assert.Equal(test, testData.want, got)
			testData.wantErr(test, err)
		})
	}
}

func initializeForUsage(usage *[]byte, writer *mocks.Writer, exiter *mocks.Exiter) {
	writer.On("Write", mock.AnythingOfType("[]uint8")).Return(func(buffer []byte) int {
		*usage = append(*usage, buffer...)
		return len(buffer)
	}, nil)

	exiter.On("Exit", 0).Return()
}

func setOption(options interpreter.Options, path string, value interface{}) interpreter.Options {
	optionReflection := reflect.ValueOf(&options).Elem()
	for _, field := range strings.Split(path, ".") {
		optionReflection = optionReflection.FieldByName(field)
	}

	optionReflection.Set(reflect.ValueOf(value))

	return options
}
