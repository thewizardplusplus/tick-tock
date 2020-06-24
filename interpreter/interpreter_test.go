package interpreter

import (
	"fmt"
	"go/types"
	"io"
	"testing"
	"testing/iotest"

	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	testutils "github.com/thewizardplusplus/tick-tock/internal/test-utils"
	testutilsmocks "github.com/thewizardplusplus/tick-tock/internal/test-utils/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	runtimemocks "github.com/thewizardplusplus/tick-tock/runtime/mocks"
	waitermocks "github.com/thewizardplusplus/tick-tock/runtime/waiter/mocks"
)

func TestInterpret(test *testing.T) {
	for _, testData := range []struct {
		name                   string
		initializeDependencies func(
			options Options,
			context *contextmocks.Context,
			waiter *waitermocks.Waiter,
			defaultReader *testutilsmocks.Reader,
		)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			initializeDependencies: func(
				options Options,
				context *contextmocks.Context,
				waiter *waitermocks.Waiter,
				defaultReader *testutilsmocks.Reader,
			) {
				context.On("ValuesNames").Return(mapset.NewSet("test"))
				context.On("Value", "test").Return(types.Nil{}, true)
				context.On("SetMessageSender", mock.AnythingOfType("runtime.ConcurrentActorGroup")).Return()
				context.On("SetStateHolder", mock.AnythingOfType("*runtime.Actor")).Return()
				context.On("Copy").Return(context)

				waiter.On("Add", 1).Return()
				waiter.On("Done").Return()

				defaultReader.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(func(buffer []byte) int {
						return copy(buffer, fmt.Sprintf(
							`actor state %s() message %s() test;;;`,
							options.InitialState,
							options.InitialMessage,
						))
					}, io.EOF)
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with the expression",
			initializeDependencies: func(
				options Options,
				context *contextmocks.Context,
				waiter *waitermocks.Waiter,
				defaultReader *testutilsmocks.Reader,
			) {
				context.On("ValuesNames").Return(mapset.NewSet("test"))
				context.On("Value", "test").Return(float64(23), true)
				context.On("SetMessageSender", mock.AnythingOfType("runtime.ConcurrentActorGroup")).Return()
				context.On("SetStateHolder", mock.AnythingOfType("*runtime.Actor")).Return()
				context.On("Copy").Return(context)

				waiter.On("Add", 1).Return()
				waiter.On("Done").Return()

				defaultReader.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(func(buffer []byte) int {
						return copy(buffer, fmt.Sprintf(
							`actor state %s() message %s() test;;;`,
							options.InitialState,
							options.InitialMessage,
						))
					}, io.EOF)
			},
			wantErr: assert.NoError,
		},
		{
			name: "error on code reading",
			initializeDependencies: func(
				options Options,
				context *contextmocks.Context,
				waiter *waitermocks.Waiter,
				defaultReader *testutilsmocks.Reader,
			) {
				defaultReader.On("Read", mock.AnythingOfType("[]uint8")).Return(0, iotest.ErrTimeout)
			},
			wantErr: assert.Error,
		},
		{
			name: "error on code parsing",
			initializeDependencies: func(
				options Options,
				context *contextmocks.Context,
				waiter *waitermocks.Waiter,
				defaultReader *testutilsmocks.Reader,
			) {
				defaultReader.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(func(buffer []byte) int { return copy(buffer, "incorrect") }, io.EOF)
			},
			wantErr: assert.Error,
		},
		{
			name: "error on code translation",
			initializeDependencies: func(
				options Options,
				context *contextmocks.Context,
				waiter *waitermocks.Waiter,
				defaultReader *testutilsmocks.Reader,
			) {
				context.On("ValuesNames").Return(mapset.NewSet("test"))

				defaultReader.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(func(buffer []byte) int {
						return copy(buffer, fmt.Sprintf("actor state %s();; actor;", options.InitialState))
					}, io.EOF)
			},
			wantErr: assert.Error,
		},
		{
			name: "error with the expression",
			initializeDependencies: func(
				options Options,
				context *contextmocks.Context,
				waiter *waitermocks.Waiter,
				defaultReader *testutilsmocks.Reader,
			) {
				context.On("ValuesNames").Return(mapset.NewSet("test"))

				defaultReader.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(func(buffer []byte) int {
						return copy(buffer, fmt.Sprintf(
							`actor state %s() message %s() unknown;;;`,
							options.InitialState,
							options.InitialMessage,
						))
					}, io.EOF)
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			options := Options{
				InboxSize:      testutils.BufferedInbox,
				InitialState:   "__initialization__",
				InitialMessage: "__initialize__",
			}
			context := new(contextmocks.Context)
			waiter := new(waitermocks.Waiter)
			defaultReader := new(testutilsmocks.Reader)
			fileSystem := new(testutilsmocks.FileSystem)
			errorHandler := new(runtimemocks.ErrorHandler)
			testData.initializeDependencies(options, context, waiter, defaultReader)

			synchronousWaiter := testutils.NewSynchronousWaiter(waiter)
			dependencies := Dependencies{
				Reader:  ReaderDependencies{defaultReader, fileSystem},
				Runtime: runtime.Dependencies{Waiter: synchronousWaiter, ErrorHandler: errorHandler},
			}
			err := Interpret(context, options, dependencies)
			synchronousWaiter.Wait()

			mock.AssertExpectationsForObjects(
				test,
				context,
				waiter,
				defaultReader,
				fileSystem,
				errorHandler,
			)
			testData.wantErr(test, err)
		})
	}
}
