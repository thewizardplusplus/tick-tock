package interpreter

import (
	"fmt"
	"io"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	runtimemocks "github.com/thewizardplusplus/tick-tock/runtime/mocks"
	"github.com/thewizardplusplus/tick-tock/internal/tests"
	testsmocks "github.com/thewizardplusplus/tick-tock/internal/tests/mocks"
	"github.com/thewizardplusplus/tick-tock/translator"
)

func TestInterpret(test *testing.T) {
	for _, testData := range []struct {
		name                   string
		initializeDependencies func(
			context *contextmocks.Context,
			waiter tests.SynchronousWaiter,
			outWriter *testsmocks.Writer,
			defaultReader *testsmocks.Reader,
		)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			initializeDependencies: func(
				context *contextmocks.Context,
				waiter tests.SynchronousWaiter,
				outWriter *testsmocks.Writer,
				defaultReader *testsmocks.Reader,
			) {
				context.On("SetMessageSender", mock.AnythingOfType("runtime.ConcurrentActorGroup")).Return()
				context.On("SetStateHolder", mock.AnythingOfType("*runtime.Actor")).Return()

				waiter.On("Add", 1).Return()
				waiter.On("Done").Return()

				const message = "Hello, world!"
				outWriter.On("Write", []byte(message)).Return(len(message), nil)

				defaultReader.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(func(buffer []byte) int {
						code := fmt.Sprintf(`actor state test message __initialization__ out "%s";;;`, message)
						return copy(buffer, code)
					}, io.EOF)
			},
			wantErr: assert.NoError,
		},
		{
			name: "error on a code reading",
			initializeDependencies: func(
				context *contextmocks.Context,
				waiter tests.SynchronousWaiter,
				outWriter *testsmocks.Writer,
				defaultReader *testsmocks.Reader,
			) {
				defaultReader.On("Read", mock.AnythingOfType("[]uint8")).Return(0, iotest.ErrTimeout)
			},
			wantErr: assert.Error,
		},
		{
			name: "error on a code parsing",
			initializeDependencies: func(
				context *contextmocks.Context,
				waiter tests.SynchronousWaiter,
				outWriter *testsmocks.Writer,
				defaultReader *testsmocks.Reader,
			) {
				defaultReader.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(func(buffer []byte) int { return copy(buffer, "incorrect") }, io.EOF)
			},
			wantErr: assert.Error,
		},
		{
			name: "error on a code translating",
			initializeDependencies: func(
				context *contextmocks.Context,
				waiter tests.SynchronousWaiter,
				outWriter *testsmocks.Writer,
				defaultReader *testsmocks.Reader,
			) {
				defaultReader.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(func(buffer []byte) int { return copy(buffer, "actor state test;; actor;") }, io.EOF)
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			context := new(contextmocks.Context)
			waiter := tests.NewSynchronousWaiter()
			errorHandler := new(runtimemocks.ErrorHandler)
			outWriter := new(testsmocks.Writer)
			defaultReader := new(testsmocks.Reader)
			fileSystem := new(testsmocks.FileSystem)
			testData.initializeDependencies(context, waiter, outWriter, defaultReader)

			inboxSize := tests.BufferedInbox
			dependencies := Dependencies{
				Dependencies: translator.Dependencies{
					Dependencies: runtime.Dependencies{Waiter: waiter, ErrorHandler: errorHandler},
					OutWriter:    outWriter,
				},
				ReaderDependencies: ReaderDependencies{defaultReader, fileSystem},
			}
			err := Interpret(context, "", inboxSize, dependencies)
			waiter.Wait()

			mock.AssertExpectationsForObjects(
				test,
				context,
				waiter,
				errorHandler,
				outWriter,
				defaultReader,
				fileSystem,
			)
			testData.wantErr(test, err)
		})
	}
}
