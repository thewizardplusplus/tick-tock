package commands

import (
	"reflect"
	"testing"
	"testing/iotest"
	"unsafe"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	contextmocks "github.com/thewizardplusplus/tick-tock/runtime/context/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	expressionsmocks "github.com/thewizardplusplus/tick-tock/runtime/expressions/mocks"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
)

func TestStartCommand(test *testing.T) {
	type fields struct {
		actorFactory expressions.Expression
	}
	type args struct {
		context context.Context
	}

	for _, testData := range []struct {
		name       string
		fields     fields
		args       args
		wantResult interface{}
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			fields: fields{
				actorFactory: func() expressions.Expression {
					actorFactory, _ := runtime.NewActorFactory(
						"Test",
						runtime.ParameterizedStateGroup{StateGroup: runtime.StateGroup{"state_0": {}, "state_1": {}}},
						context.State{Name: "state_0"},
					)
					concurrentActorFactory :=
						runtime.NewConcurrentActorFactory(actorFactory, 0, runtime.Dependencies{})

					expression := new(expressionsmocks.Expression)
					expression.
						On("Evaluate", mock.AnythingOfType("*mocks.Context")).
						Return(concurrentActorFactory, nil)

					return expression
				}(),
			},
			args: args{
				context: func() context.Context {
					actorFactory, _ := runtime.NewActorFactory(
						"Test",
						runtime.ParameterizedStateGroup{StateGroup: runtime.StateGroup{"state_0": {}, "state_1": {}}},
						context.State{Name: "state_0"},
					)
					concurrentActorFactory :=
						runtime.NewConcurrentActorFactory(actorFactory, 0, runtime.Dependencies{})
					wantActor := concurrentActorFactory.CreateActor()
					cleanInbox(&wantActor)

					context := new(contextmocks.Context)
					context.
						On("RegisterActor", mock.MatchedBy(func(gotActor runtime.ConcurrentActor) bool {
							cleanInbox(&gotActor)
							return reflect.DeepEqual(wantActor, gotActor)
						})).
						Return()

					return context
				}(),
			},
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name: "error with actor class evaluation",
			fields: fields{
				actorFactory: func() expressions.Expression {
					expression := new(expressionsmocks.Expression)
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(nil, iotest.ErrTimeout)

					return expression
				}(),
			},
			args: args{
				context: new(contextmocks.Context),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "error with an incorrect actor class type",
			fields: fields{
				actorFactory: func() expressions.Expression {
					expression := new(expressionsmocks.Expression)
					expression.On("Evaluate", mock.AnythingOfType("*mocks.Context")).Return(2.3, nil)

					return expression
				}(),
			},
			args: args{
				context: new(contextmocks.Context),
			},
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			gotResult, gotErr := NewStartCommand(testData.fields.actorFactory).Run(testData.args.context)

			mock.AssertExpectationsForObjects(test, testData.args.context)
			assert.Equal(test, testData.wantResult, gotResult)
			testData.wantErr(test, gotErr)
		})
	}
}

func cleanInbox(actor *runtime.ConcurrentActor) {
	inboxField := reflect.ValueOf(actor).Elem().FieldByName("inbox")
	*(*chan string)(unsafe.Pointer(inboxField.UnsafeAddr())) = nil
}
