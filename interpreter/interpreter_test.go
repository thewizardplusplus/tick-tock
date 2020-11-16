package interpreter

import (
	"fmt"
	"go/types"
	"io"
	"sync"
	"testing"
	"testing/iotest"

	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	syncutils "github.com/thewizardplusplus/go-sync-utils"
	testutilsmocks "github.com/thewizardplusplus/tick-tock/internal/test-utils/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime"
)

func TestInterpret(test *testing.T) {
	for _, testData := range []struct {
		name                   string
		initializeDependencies func(
			options Options,
			context *MockContext,
			waiter *MockWaitGroup,
			defaultReader *testutilsmocks.Reader,
		)
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			initializeDependencies: func(
				options Options,
				context *MockContext,
				waiter *MockWaitGroup,
				defaultReader *testutilsmocks.Reader,
			) {
				context.On("ValuesNames").Return(mapset.NewSet("test"))
				context.On("Value", "test").Return(types.Nil{}, true)
				context.On("SetValue", "Main", mock.AnythingOfType("runtime.ConcurrentActorFactory")).Return()
				context.On("SetMessageSender", mock.AnythingOfType("*runtime.ConcurrentActorGroup")).Return()
				context.On("SetActorRegister", mock.AnythingOfType("*runtime.ConcurrentActorGroup")).Return()
				context.On("SetStateHolder", mock.AnythingOfType("*runtime.Actor")).Return()
				context.On("Copy").Return(context)

				waiter.On("Add", 1).Return()
				waiter.On("Done").Return()

				defaultReader.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(func(buffer []byte) int {
						return copy(buffer, fmt.Sprintf(
							`actor Main() state %s() message %s() test;;;`,
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
				context *MockContext,
				waiter *MockWaitGroup,
				defaultReader *testutilsmocks.Reader,
			) {
				context.On("ValuesNames").Return(mapset.NewSet("test"))
				context.On("Value", "test").Return(float64(23), true)
				context.On("SetValue", "Main", mock.AnythingOfType("runtime.ConcurrentActorFactory")).Return()
				context.On("SetMessageSender", mock.AnythingOfType("*runtime.ConcurrentActorGroup")).Return()
				context.On("SetActorRegister", mock.AnythingOfType("*runtime.ConcurrentActorGroup")).Return()
				context.On("SetStateHolder", mock.AnythingOfType("*runtime.Actor")).Return()
				context.On("Copy").Return(context)

				waiter.On("Add", 1).Return()
				waiter.On("Done").Return()

				defaultReader.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(func(buffer []byte) int {
						return copy(buffer, fmt.Sprintf(
							`actor Main() state %s() message %s() test;;;`,
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
				context *MockContext,
				waiter *MockWaitGroup,
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
				context *MockContext,
				waiter *MockWaitGroup,
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
				context *MockContext,
				waiter *MockWaitGroup,
				defaultReader *testutilsmocks.Reader,
			) {
				context.On("ValuesNames").Return(mapset.NewSet("test"))

				defaultReader.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(func(buffer []byte) int {
						return copy(buffer, fmt.Sprintf(
							"actor Main() state %s();; actor Incorrect();",
							options.InitialState,
						))
					}, io.EOF)
			},
			wantErr: assert.Error,
		},
		{
			name: "error with the expression",
			initializeDependencies: func(
				options Options,
				context *MockContext,
				waiter *MockWaitGroup,
				defaultReader *testutilsmocks.Reader,
			) {
				context.On("ValuesNames").Return(mapset.NewSet("test"))

				defaultReader.
					On("Read", mock.AnythingOfType("[]uint8")).
					Return(func(buffer []byte) int {
						return copy(buffer, fmt.Sprintf(
							`actor Main() state %s() message %s() unknown;;;`,
							options.InitialState,
							options.InitialMessage,
						))
					}, io.EOF)
			},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			waiter := new(MockWaitGroup)
			waiter.On("Wait").Times(1)

			options := Options{
				InboxSize:      0,
				InitialState:   "__initialization__",
				InitialMessage: "__initialize__",
			}
			context := new(MockContext)
			defaultReader := new(testutilsmocks.Reader)
			fileSystem := new(testutilsmocks.FileSystem)
			errorHandler := new(MockErrorHandler)
			testData.initializeDependencies(options, context, waiter, defaultReader)

			synchronousWaiter := syncutils.MultiWaitGroup{waiter, new(sync.WaitGroup)}
			dependencies := Dependencies{
				Reader:  ReaderDependencies{defaultReader, fileSystem},
				Runtime: runtime.Dependencies{WaitGroup: synchronousWaiter, ErrorHandler: errorHandler},
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
