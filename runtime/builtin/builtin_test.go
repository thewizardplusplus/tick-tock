package builtin

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/thewizardplusplus/tick-tock/parser"
	"github.com/thewizardplusplus/tick-tock/runtime"
	"github.com/thewizardplusplus/tick-tock/runtime/context"
	"github.com/thewizardplusplus/tick-tock/runtime/expressions"
	"github.com/thewizardplusplus/tick-tock/runtime/types"
	"github.com/thewizardplusplus/tick-tock/translator"
)

func TestValues(test *testing.T) {
	for _, data := range []struct {
		name                  string
		additionalDefinitions context.ValueGroup
		code                  string
		wantResult            interface{}
		wantErr               assert.ErrorAssertionFunc
	}{
		{
			name:       "empty list",
			code:       "[]",
			wantResult: (*types.Pair)(nil),
			wantErr:    assert.NoError,
		},
		{
			name:       "empty hash table",
			code:       "{}",
			wantResult: (types.HashTable)(nil),
			wantErr:    assert.NoError,
		},
		{
			name:       "nil",
			code:       "nil",
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name:       "false",
			code:       "false",
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "true",
			code:       "true",
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "inf",
			code:       "inf",
			wantResult: math.Inf(+1),
			wantErr:    assert.NoError,
		},
		{
			name:       "pi",
			code:       "pi",
			wantResult: math.Pi,
			wantErr:    assert.NoError,
		},
		{
			name:       "e",
			code:       "e",
			wantResult: math.E,
			wantErr:    assert.NoError,
		},
		{
			name:       "list construction",
			code:       "[12, 23, 42]",
			wantResult: types.NewPairFromSlice([]interface{}{12.0, 23.0, 42.0}),
			wantErr:    assert.NoError,
		},
		{
			name:       "hash table construction/success",
			code:       "{x: 12, y: 23, z: 42}",
			wantResult: types.HashTable{"x": 12.0, "y": 23.0, "z": 42.0},
			wantErr:    assert.NoError,
		},
		{
			name:       "hash table construction/error",
			code:       "{x: 12, [[-23]]: 23, z: 42}",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "equal/success/false/same types",
			code:       "2 == 3",
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "equal/success/false/different types",
			code:       "2 == nil",
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "equal/success/true",
			code:       "2 == 2",
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "equal/error",
			code:       "__eq__ == nil",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "not equal/success/false",
			code:       "2 != 2",
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "not equal/success/true/same types",
			code:       "2 != 3",
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "not equal/success/true/different types",
			code:       "2 != nil",
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "not equal/error",
			code:       "__eq__ != nil",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "less/success/false",
			code:       "4 < 2",
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "less/success/true",
			code:       "2 < 3",
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "less/error",
			code:       "2 < nil",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "less or equal/success/false",
			code:       "4 <= 2",
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "less or equal/success/true/less",
			code:       "2 <= 3",
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "less or equal/success/true/equal",
			code:       "2 <= 2",
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "less or equal/error",
			code:       "2 <= nil",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "greater/success/false",
			code:       "2 > 3",
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "greater/success/true",
			code:       "4 > 2",
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "greater/error",
			code:       "2 > nil",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "greater or equal/success/false",
			code:       "2 >= 3",
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "greater or equal/success/true/greater",
			code:       "4 >= 2",
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "greater or equal/success/true/equal",
			code:       "2 >= 2",
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "greater or equal/error",
			code:       "2 >= nil",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "bitwise disjunction/positive",
			code:       "23 | 42",
			wantResult: 63.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "bitwise disjunction/negative",
			code:       "-23 | -42",
			wantResult: -1.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "bitwise exclusive disjunction/positive",
			code:       "23 ^ 42",
			wantResult: 61.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "bitwise exclusive disjunction/negative",
			code:       "-23 ^ -42",
			wantResult: 63.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "bitwise conjunction/positive",
			code:       "23 & 42",
			wantResult: 2.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "bitwise conjunction/negative",
			code:       "-23 & -42",
			wantResult: -64.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "left shift/positive",
			code:       "2 << 3",
			wantResult: 16.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "left shift/negative",
			code:       "-2 << 3",
			wantResult: -16.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "right shift/positive",
			code:       "16 >> 3",
			wantResult: 2.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "right shift/negative",
			code:       "-16 >> 3",
			wantResult: -2.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "unsigned right shift/positive",
			code:       "16 >>> 3",
			wantResult: 2.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "unsigned right shift/negative",
			code:       "-16 >>> 3",
			wantResult: 536870910.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "addition/success/float64",
			code:       "2 + 3",
			wantResult: 5.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "addition/success/*types.Pair",
			code:       `"te" + "st"`,
			wantResult: types.NewPairFromText("test"),
			wantErr:    assert.NoError,
		},
		{
			name: "addition/success/types.HashTable",
			code: `{[12]: "one", [23]: "two"} + {[23]: "three", [42]: "four"}`,
			wantResult: types.HashTable{
				12.0: types.NewPairFromText("one"),
				23.0: types.NewPairFromText("three"),
				42.0: types.NewPairFromText("four"),
			},
			wantErr: assert.NoError,
		},
		{
			name:       "addition/error/argument #0",
			code:       "__add__ + []",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "addition/error/argument #1",
			code:       "23 + []",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "subtraction",
			code:       "2 - 3",
			wantResult: -1.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "multiplication",
			code:       "2 * 3",
			wantResult: 6.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "division",
			code:       "10 / 2",
			wantResult: 5.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "modulo",
			code:       "10 % 3",
			wantResult: 1.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "arithmetic negation",
			code:       "-23",
			wantResult: -23.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "bitwise negation/positive",
			code:       "~23",
			wantResult: -24.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "bitwise negation/negative",
			code:       "~-23",
			wantResult: 22.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "logical negation/success/false",
			code:       "!false",
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "logical negation/success/true",
			code:       "!true",
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "logical negation/error",
			code:       "!__logical_not__",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "key accessor/success/*types.Pair/index in range",
			code:       "[12, 23, 42][1]",
			wantResult: 23.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "key accessor/success/*types.Pair/index out of range",
			code:       "[12, 23, 42][23]",
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name:       "key accessor/success/types.HashTable/existing key",
			code:       `{x: 12, y: 23, z: 42}["y"]`,
			wantResult: 23.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "key accessor/success/types.HashTable/nonexistent key",
			code:       `{x: 12, y: 23, z: 42}["nonexistent"]`,
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name:       "key accessor/error/incorrect index for *types.Pair",
			code:       `[12, 23, 42]["incorrect"]`,
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "key accessor/error/incorrect key for types.HashTable",
			code:       `{x: 12, y: 23, z: 42}[[-23]]`,
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "key accessor/error/unsupported type",
			code:       "__item__[1]",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "type/success/nil",
			code:       "type(nil)",
			wantResult: types.NewPairFromText("nil"),
			wantErr:    assert.NoError,
		},
		{
			name:       "type/success/float64",
			code:       "type(23)",
			wantResult: types.NewPairFromText("num"),
			wantErr:    assert.NoError,
		},
		{
			name:       "type/success/*types.Pair",
			code:       "type([12, 23, 42])",
			wantResult: types.NewPairFromText("list"),
			wantErr:    assert.NoError,
		},
		{
			name:       "type/success/types.HashTable",
			code:       "type({x: 12, y: 23, z: 42})",
			wantResult: types.NewPairFromText("hash"),
			wantErr:    assert.NoError,
		},
		{
			name: "type/success/actor class",
			code: "type(Test)",
			additionalDefinitions: context.ValueGroup{
				"Test": func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test",
						runtime.ParameterizedStateGroup{StateGroup: runtime.StateGroup{"state_0": {}}},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 0, runtime.Dependencies{})
				}(),
			},
			wantResult: types.NewPairFromText("class"),
			wantErr:    assert.NoError,
		},
		{
			name:       "type/error",
			code:       "type(type)",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name: "name",
			code: "name(Test)",
			additionalDefinitions: context.ValueGroup{
				"Test": func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test",
						runtime.ParameterizedStateGroup{StateGroup: runtime.StateGroup{"state_0": {}}},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 0, runtime.Dependencies{})
				}(),
			},
			wantResult: types.NewPairFromText("Test"),
			wantErr:    assert.NoError,
		},
		{
			name:       "size/success/*types.Pair",
			code:       "size([12, 23, 42])",
			wantResult: 3.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "size/success/types.HashTable",
			code:       "size({x: 12, y: 23, z: 42})",
			wantResult: 3.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "size/error",
			code:       "size(size)",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "bool/success/false",
			code:       `bool("")`,
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "bool/success/true",
			code:       `bool("test")`,
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "bool/error",
			code:       "bool(bool)",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "floor",
			code:       "floor(2.5)",
			wantResult: 2.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "ceil",
			code:       "ceil(2.5)",
			wantResult: 3.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "trunc",
			code:       "trunc(2.5)",
			wantResult: 2.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "round",
			code:       "round(2.5)",
			wantResult: 3.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "pow",
			code:       "pow(2, 3)",
			wantResult: 8.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "sqrt",
			code:       "sqrt(4)",
			wantResult: 2.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "abs",
			code:       "abs(-23)",
			wantResult: 23.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "is_nan/false",
			code:       "is_nan(23)",
			wantResult: types.False,
			wantErr:    assert.NoError,
		},
		{
			name:       "is_nan/true",
			code:       "is_nan(nan)",
			wantResult: types.True,
			wantErr:    assert.NoError,
		},
		{
			name:       "head/success",
			code:       "head([12, 23, 42])",
			wantResult: 12.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "head/error",
			code:       "head([])",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "tail/success/nonempty tail",
			code:       "tail([12, 23, 42])",
			wantResult: &types.Pair{Head: 23.0, Tail: &types.Pair{Head: 42.0, Tail: nil}},
			wantErr:    assert.NoError,
		},
		{
			name:       "tail/success/empty tail",
			code:       "tail([23])",
			wantResult: (*types.Pair)(nil),
			wantErr:    assert.NoError,
		},
		{
			name:       "tail/error",
			code:       "tail([])",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "num/success/correct number",
			code:       `num("23")`,
			wantResult: 23.0,
			wantErr:    assert.NoError,
		},
		{
			name:       "num/success/incorrect number",
			code:       `num("test")`,
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name:       "num/error",
			code:       `num(['t', "hi", 's', 't'])`,
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "str/success/nil",
			code:       "str(nil)",
			wantResult: types.NewPairFromText("null"),
			wantErr:    assert.NoError,
		},
		{
			name:       "str/success/float64",
			code:       "str(23)",
			wantResult: types.NewPairFromText("23"),
			wantErr:    assert.NoError,
		},
		{
			name:       "str/success/*types.Pair/tree in the head",
			code:       `str(["hi", 23, 42])`,
			wantResult: types.NewPairFromText("[[104,105],23,42]"),
			wantErr:    assert.NoError,
		},
		{
			name:       "str/success/*types.Pair/tree in the tail",
			code:       `str([12, "hi", 42])`,
			wantResult: types.NewPairFromText("[12,[104,105],42]"),
			wantErr:    assert.NoError,
		},
		{
			name:       "str/success/*types.Pair/with the nil type",
			code:       "str([12, nil, 42])",
			wantResult: types.NewPairFromText("[12,null,42]"),
			wantErr:    assert.NoError,
		},
		{
			name:       "str/success/*types.Pair/with the hash table",
			code:       "str([12, {x: 12, y: 23, z: 42}, 42])",
			wantResult: types.NewPairFromText(`[12,{"x":12,"y":23,"z":42},42]`),
			wantErr:    assert.NoError,
		},
		{
			name: "str/success/*types.Pair/with the actor class",
			code: "str([12, Test, 42])",
			additionalDefinitions: context.ValueGroup{
				"Test": func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test",
						runtime.ParameterizedStateGroup{StateGroup: runtime.StateGroup{"state_0": {}}},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 0, runtime.Dependencies{})
				}(),
			},
			wantResult: types.NewPairFromText(`[12,"<class Test>",42]`),
			wantErr:    assert.NoError,
		},
		{
			name:       "str/success/types.HashTable/tree",
			code:       "str({x: 12, y: {x: 12, y: 23, z: 42}, z: 42})",
			wantResult: types.NewPairFromText(`{"x":12,"y":{"x":12,"y":23,"z":42},"z":42}`),
			wantErr:    assert.NoError,
		},
		{
			name:       "str/success/types.HashTable/with the nil type",
			code:       "str({x: 12, y: nil, z: 42})",
			wantResult: types.NewPairFromText(`{"x":12,"z":42}`),
			wantErr:    assert.NoError,
		},
		{
			name:       "str/success/types.HashTable/with the list",
			code:       "str({x: 12, y: [12, 23, 42], z: 42})",
			wantResult: types.NewPairFromText(`{"x":12,"y":[12,23,42],"z":42}`),
			wantErr:    assert.NoError,
		},
		{
			name: "str/success/types.HashTable/with the actor class",
			code: "str({x: 12, y: Test, z: 42})",
			additionalDefinitions: context.ValueGroup{
				"Test": func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test",
						runtime.ParameterizedStateGroup{StateGroup: runtime.StateGroup{"state_0": {}}},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 0, runtime.Dependencies{})
				}(),
			},
			wantResult: types.NewPairFromText(`{"x":12,"y":"<class Test>","z":42}`),
			wantErr:    assert.NoError,
		},
		{
			name: "str/success/actor class",
			code: "str(Test)",
			additionalDefinitions: context.ValueGroup{
				"Test": func() runtime.ConcurrentActorFactory {
					actorFactory, _ := runtime.NewActorFactory(
						"Test",
						runtime.ParameterizedStateGroup{StateGroup: runtime.StateGroup{"state_0": {}}},
						context.State{Name: "state_0"},
					)
					return runtime.NewConcurrentActorFactory(actorFactory, 0, runtime.Dependencies{})
				}(),
			},
			wantResult: types.NewPairFromText("<class Test>"),
			wantErr:    assert.NoError,
		},
		{
			name:       "str/error/deep value/*types.Pair",
			code:       `str([12, {[12]: "one", [23]: "two", [42]: "three"}, 42])`,
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "str/error/deep value/types.HashTable",
			code:       `str({[12]: "one", [23]: "two", [42]: "three"})`,
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "str/error/JSON marshalling/*types.Pair",
			code:       "str([12, str, 42])",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "str/error/JSON marshalling/types.HashTable",
			code:       "str({x: 12, y: str, z: 42})",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "str/error/unsupported type",
			code:       "str(str)",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "strb/success/false",
			code:       `strb("")`,
			wantResult: types.NewPairFromText("false"),
			wantErr:    assert.NoError,
		},
		{
			name:       "strb/success/true",
			code:       `strb("test")`,
			wantResult: types.NewPairFromText("true"),
			wantErr:    assert.NoError,
		},
		{
			name:       "strb/error",
			code:       "strb(strb)",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "strs/success",
			code:       `strs("\"test\"")`,
			wantResult: types.NewPairFromText(`"\"test\""`),
			wantErr:    assert.NoError,
		},
		{
			name:       "strs/error",
			code:       `strs(['t', "hi", 's', 't'])`,
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "strl/success",
			code:       `strl(["\"one\"", "\"two\""])`,
			wantResult: types.NewPairFromText(`["\"one\"","\"two\""]`),
			wantErr:    assert.NoError,
		},
		{
			name:       "strl/error/incorrect type",
			code:       `strl(["\"one\"", 23])`,
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "strl/error/string conversion",
			code:       `strl(["\"one\"", ['\x22', "hi", 'w', 'o', '\x22']])`,
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "strh/success",
			code:       "strh({one: 12, two: 23, three: 42})",
			wantResult: types.NewPairFromText(`{"one":12,"three":42,"two":23}`),
			wantErr:    assert.NoError,
		},
		{
			name:       "strh/error/incorrect key",
			code:       "strh({one: 5, [12]: 23, three: 42})",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "strh/error/incorrect value",
			code:       "strh({one: 12, two: strh, three: 42})",
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "strhh/success",
			code:       `strhh({one: "two", three: "four", five: "six"})`,
			wantResult: types.NewPairFromText(`{"five":"six","one":"two","three":"four"}`),
			wantErr:    assert.NoError,
		},
		{
			name:       "strhh/error/incorrect key type",
			code:       `strhh({one: "two", [23]: "four", five: "six"})`,
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "strhh/error/incorrect value type",
			code:       `strhh({one: "two", three: 23, five: "six"})`,
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "strhh/error/string conversion",
			code:       `strhh({one: "two", three: ['t', "hi", 's', 't'], five: "six"})`,
			wantResult: nil,
			wantErr:    assert.Error,
		},
		{
			name:       "with/success",
			code:       `with({x: 12, y: 23}, "z", 42)`,
			wantResult: types.HashTable{"x": 12.0, "y": 23.0, "z": 42.0},
			wantErr:    assert.NoError,
		},
		{
			name:       "with/error",
			code:       `with({x: 12, [[-23]]: 23}, "z", 42)`,
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			ctx := context.NewDefaultContext()
			context.SetValues(ctx, Values)
			context.SetValues(ctx, data.additionalDefinitions)

			expressionAST := new(parser.Expression)
			err := parser.ParseToAST(data.code, expressionAST)
			require.NoError(test, err)

			expression, _, err := translator.TranslateExpression(expressionAST, ctx.ValuesNames())
			require.NoError(test, err)

			gotResult, gotErr := expression.Evaluate(ctx)

			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}

func TestValues_nan(test *testing.T) {
	ctx := context.NewDefaultContext()
	context.SetValues(ctx, Values)

	expression := expressions.NewIdentifier("nan")
	result, err := expression.Evaluate(ctx)

	if assert.NoError(test, err) {
		require.IsType(test, float64(0), result)
		assert.True(test, math.IsNaN(result.(float64)))
	}
}

func TestValues_inDelta(test *testing.T) {
	for _, data := range []struct {
		name string
		code string
		want float64
	}{
		{
			name: "sin",
			code: "sin(23)",
			want: -0.846220,
		},
		{
			name: "cos",
			code: "cos(23)",
			want: -0.532833,
		},
		{
			name: "tn",
			code: "tn(23)",
			want: 1.588153,
		},
		{
			name: "arcsin",
			code: "arcsin(0.5)",
			want: 0.523598,
		},
		{
			name: "arccos",
			code: "arccos(0.5)",
			want: 1.047197,
		},
		{
			name: "arctn",
			code: "arctn(0.5)",
			want: 0.463647,
		},
		{
			name: "angle",
			code: "angle(2, 3)",
			want: 0.982793,
		},
		{
			name: "exp",
			code: "exp(2.3)",
			want: 9.974182,
		},
		{
			name: "ln",
			code: "ln(23)",
			want: 3.135494,
		},
		{
			name: "lg",
			code: "lg(23)",
			want: 1.361727,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			ctx := context.NewDefaultContext()
			context.SetValues(ctx, Values)

			expressionAST := new(parser.Expression)
			err := parser.ParseToAST(data.code, expressionAST)
			require.NoError(test, err)

			expression, _, err := translator.TranslateExpression(expressionAST, ctx.ValuesNames())
			require.NoError(test, err)

			got, err := expression.Evaluate(ctx)

			if assert.NoError(test, err) {
				require.IsType(test, float64(0), got)
				assert.InDelta(test, data.want, got.(float64), 1e-6)
			}
		})
	}
}

func TestValues_random(test *testing.T) {
	const numberCount = 10

	ctx := context.NewDefaultContext()
	context.SetValues(ctx, Values)

	expression := expressions.NewFunctionCall("seed", []expressions.Expression{
		expressions.NewNumber(23),
	})
	got, err := expression.Evaluate(ctx)

	assert.Equal(test, types.Nil{}, got)
	assert.NoError(test, err)

	var numbers []float64
	for i := 0; i < numberCount; i++ {
		expression := expressions.NewFunctionCall("random", nil)
		result, err := expression.Evaluate(ctx)

		assert.IsType(test, float64(0), result)
		assert.NoError(test, err)

		if number, ok := result.(float64); ok {
			numbers = append(numbers, number)
		}
	}

	rand.Seed(23)

	var wantNumbers []float64
	for i := 0; i < numberCount; i++ {
		wantNumber := rand.Float64()
		wantNumbers = append(wantNumbers, wantNumber)
	}

	assert.InDeltaSlice(test, wantNumbers, numbers, 1e-6)
}

func TestValues_keys(test *testing.T) {
	for _, data := range []struct {
		name string
		code string
		want []interface{}
	}{
		{
			name: "empty",
			code: "keys({})",
			want: nil,
		},
		{
			name: "float64",
			code: `keys({[12]: "one", [23]: "two", [42]: "three"})`,
			want: []interface{}{12.0, 23.0, 42.0},
		},
		{
			name: "*types.Pair",
			code: `keys({one: 12, two: 23, three: 42})`,
			want: []interface{}{
				types.NewPairFromText("one"),
				types.NewPairFromText("two"),
				types.NewPairFromText("three"),
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			ctx := context.NewDefaultContext()
			context.SetValues(ctx, Values)

			expressionAST := new(parser.Expression)
			err := parser.ParseToAST(data.code, expressionAST)
			require.NoError(test, err)

			expression, _, err := translator.TranslateExpression(expressionAST, ctx.ValuesNames())
			require.NoError(test, err)

			got, err := expression.Evaluate(ctx)

			if assert.IsType(test, (*types.Pair)(nil), got) {
				assert.ElementsMatch(test, data.want, got.(*types.Pair).Slice())
			}
			assert.NoError(test, err)
		})
	}
}

func TestValues_env(test *testing.T) {
	const envName = "TEST"

	for _, data := range []struct {
		name       string
		prepare    func(test *testing.T)
		code       string
		wantResult interface{}
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success/existing variable/nonempty value",
			prepare: func(test *testing.T) {
				err := os.Setenv(envName, "test")
				require.NoError(test, err)
			},
			code:       fmt.Sprintf("env(%q)", envName),
			wantResult: types.NewPairFromText("test"),
			wantErr:    assert.NoError,
		},
		{
			name: "success/existing variable/empty value",
			prepare: func(test *testing.T) {
				err := os.Setenv(envName, "")
				require.NoError(test, err)
			},
			code:       fmt.Sprintf("env(%q)", envName),
			wantResult: (*types.Pair)(nil),
			wantErr:    assert.NoError,
		},
		{
			name: "success/nonexistent variable",
			prepare: func(test *testing.T) {
				err := os.Unsetenv(envName)
				require.NoError(test, err)
			},
			code:       fmt.Sprintf("env(%q)", envName),
			wantResult: types.Nil{},
			wantErr:    assert.NoError,
		},
		{
			name:       "error",
			prepare:    func(test *testing.T) {},
			code:       `env(['t', "hi", 's', 't'])`,
			wantResult: nil,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			previousValue, wasSet := os.LookupEnv(envName)
			defer func() {
				if wasSet {
					err := os.Setenv(envName, previousValue)
					require.NoError(test, err)
				}
			}()
			data.prepare(test)

			ctx := context.NewDefaultContext()
			context.SetValues(ctx, Values)

			expressionAST := new(parser.Expression)
			err := parser.ParseToAST(data.code, expressionAST)
			require.NoError(test, err)

			expression, _, err := translator.TranslateExpression(expressionAST, ctx.ValuesNames())
			require.NoError(test, err)

			gotResult, gotErr := expression.Evaluate(ctx)

			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}

func TestValues_time(test *testing.T) {
	ctx := context.NewDefaultContext()
	context.SetValues(ctx, Values)

	expression := expressions.NewFunctionCall("time", nil)
	result, err := expression.Evaluate(ctx)

	if assert.NoError(test, err) {
		require.IsType(test, float64(0), result)

		resultTime := time.Unix(0, int64(result.(float64)*1e9))
		assert.WithinDuration(test, time.Now(), resultTime, time.Minute)
	}
}

func TestValues_sleep(test *testing.T) {
	startTime := time.Now()

	ctx := context.NewDefaultContext()
	context.SetValues(ctx, Values)

	expression := expressions.NewFunctionCall("sleep", []expressions.Expression{
		expressions.NewNumber(2.3),
	})
	result, err := expression.Evaluate(ctx)

	elapsedTime := int64(time.Since(startTime))
	assert.GreaterOrEqual(test, elapsedTime, int64(2300*time.Millisecond))
	assert.Less(test, elapsedTime, int64(time.Minute))
	assert.Equal(test, types.Nil{}, result)
	assert.NoError(test, err)
}

// based on https://talks.golang.org/2014/testing.slide#23 by Andrew Gerrand
func TestValues_exit(test *testing.T) {
	if os.Getenv("EXIT_TEST") == "TRUE" {
		ctx := context.NewDefaultContext()
		context.SetValues(ctx, Values)

		expression := expressions.NewFunctionCall("exit", []expressions.Expression{
			expressions.NewNumber(23),
		})
		expression.Evaluate(ctx) // nolint: errcheck

		return
	}

	command := exec.Command(os.Args[0], "-test.run=TestValues_exit")
	command.Env = append(os.Environ(), "EXIT_TEST=TRUE")

	err := command.Run()

	assert.IsType(test, (*exec.ExitError)(nil), err)
	assert.EqualError(test, err, "exit status 23")
}

func TestValues_input(test *testing.T) {
	for _, data := range []struct {
		name       string
		prepare    func(test *testing.T, tempFile *os.File)
		code       string
		wantResult interface{}
	}{
		{
			name: "in/part of symbols/success/part of symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			code:       "in(2)",
			wantResult: types.NewPairFromText("te"),
		},
		{
			name: "in/part of symbols/success/all symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			code:       "in(4)",
			wantResult: types.NewPairFromText("test"),
		},
		{
			name: "in/part of symbols/success/sequential calls",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			// simulate sequential calls via wrapping into a list
			code: "[in(2), in(2)]",
			wantResult: types.NewPairFromSlice([]interface{}{
				types.NewPairFromText("te"),
				types.NewPairFromText("st"),
			}),
		},
		{
			name:       "in/part of symbols/error/without symbols",
			prepare:    func(test *testing.T, tempFile *os.File) {},
			code:       "in(2)",
			wantResult: types.Nil{},
		},
		{
			name: "in/part of symbols/error/with lack of symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			code:       "in(5)",
			wantResult: types.Nil{},
		},
		{
			name: "in/all symbols/with symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			code:       "in(-1)",
			wantResult: types.NewPairFromText("test"),
		},
		{
			name:       "in/all symbols/without symbols",
			prepare:    func(test *testing.T, tempFile *os.File) {},
			code:       "in(-1)",
			wantResult: (*types.Pair)(nil),
		},
		{
			name: "in/all symbols/sequential calls",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			// simulate sequential calls via wrapping into a list
			code: "[in(-1), in(-1)]",
			wantResult: types.NewPairFromSlice([]interface{}{
				types.NewPairFromText("test"),
				(*types.Pair)(nil),
			}),
		},
		{
			name: "inln/part of symbols/success/part of symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			code:       "inln(2)",
			wantResult: types.NewPairFromText("te"),
		},
		{
			name: "inln/part of symbols/success/all symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			code:       "inln(4)",
			wantResult: types.NewPairFromText("test"),
		},
		{
			name: "inln/part of symbols/success/sequential calls",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			// simulate sequential calls via wrapping into a list
			code: "[inln(2), inln(2)]",
			wantResult: types.NewPairFromSlice([]interface{}{
				types.NewPairFromText("te"),
				types.NewPairFromText("st"),
			}),
		},
		{
			name:       "inln/part of symbols/error/without symbols",
			prepare:    func(test *testing.T, tempFile *os.File) {},
			code:       "inln(2)",
			wantResult: types.Nil{},
		},
		{
			name: "inln/part of symbols/error/with lack of symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			code:       "inln(5)",
			wantResult: types.Nil{},
		},
		{
			name: "inln/all symbols/success/with symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test #1\ntest #2\n")
				require.NoError(test, err)
			},
			code:       "inln(-1)",
			wantResult: types.NewPairFromText("test #1"),
		},
		{
			name: "inln/all symbols/success/without symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("\ntest #2\n")
				require.NoError(test, err)
			},
			code:       "inln(-1)",
			wantResult: (*types.Pair)(nil),
		},
		{
			name: "inln/all symbols/success/sequential calls",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test #1\ntest #2\n")
				require.NoError(test, err)
			},
			// simulate sequential calls via wrapping into a list
			code: "[inln(-1), inln(-1)]",
			wantResult: types.NewPairFromSlice([]interface{}{
				types.NewPairFromText("test #1"),
				types.NewPairFromText("test #2"),
			}),
		},
		{
			name: "inln/all symbols/error/with symbols",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			code:       "inln(-1)",
			wantResult: types.Nil{},
		},
		{
			name:       "inln/all symbols/error/without symbols",
			prepare:    func(test *testing.T, tempFile *os.File) {},
			code:       "inln(-1)",
			wantResult: types.Nil{},
		},
		{
			name: "in & inln/part of symbols/sequential calls",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			// simulate sequential calls via wrapping into a list
			code: "[in(2), inln(2)]",
			wantResult: types.NewPairFromSlice([]interface{}{
				types.NewPairFromText("te"),
				types.NewPairFromText("st"),
			}),
		},
		{
			name: "in & inln/all symbols/sequential calls",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test #1\ntest #2\n")
				require.NoError(test, err)
			},
			// simulate sequential calls via wrapping into a list
			code: "[in(-1), inln(-1)]",
			wantResult: types.NewPairFromSlice([]interface{}{
				types.NewPairFromText("test #1\ntest #2\n"),
				types.Nil{},
			}),
		},
		{
			name: "inln & in/part of symbols/sequential calls",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test")
				require.NoError(test, err)
			},
			// simulate sequential calls via wrapping into a list
			code: "[inln(2), in(2)]",
			wantResult: types.NewPairFromSlice([]interface{}{
				types.NewPairFromText("te"),
				types.NewPairFromText("st"),
			}),
		},
		{
			name: "inln & in/all symbols/sequential calls",
			prepare: func(test *testing.T, tempFile *os.File) {
				_, err := tempFile.WriteString("test #1\ntest #2\n")
				require.NoError(test, err)
			},
			// simulate sequential calls via wrapping into a list
			code: "[inln(-1), in(-1)]",
			wantResult: types.NewPairFromSlice([]interface{}{
				types.NewPairFromText("test #1"),
				types.NewPairFromText("test #2\n"),
			}),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			previousStdin := bufferedStdin
			defer func() { bufferedStdin = previousStdin }()

			tempFile, err := ioutil.TempFile("", "test.*")
			require.NoError(test, err)
			defer os.Remove(tempFile.Name()) // nolint: errcheck
			defer tempFile.Close()           // nolint: errcheck

			data.prepare(test, tempFile)
			err = tempFile.Close()
			require.NoError(test, err)

			tempFile, err = os.Open(tempFile.Name())
			require.NoError(test, err)
			bufferedStdin = bufio.NewReader(tempFile)

			ctx := context.NewDefaultContext()
			context.SetValues(ctx, Values)

			expressionAST := new(parser.Expression)
			err = parser.ParseToAST(data.code, expressionAST)
			require.NoError(test, err)

			expression, _, err := translator.TranslateExpression(expressionAST, ctx.ValuesNames())
			require.NoError(test, err)

			gotResult, gotErr := expression.Evaluate(ctx)

			assert.Equal(test, data.wantResult, gotResult)
			assert.NoError(test, gotErr)
		})
	}
}

func TestValues_output(test *testing.T) {
	for _, data := range []struct {
		name       string
		prepare    func(test *testing.T, tempFile *os.File)
		code       string
		wantResult interface{}
		wantOutput string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:       "out/success",
			prepare:    func(test *testing.T, tempFile *os.File) { os.Stdout = tempFile },
			code:       `out("test")`,
			wantResult: types.Nil{},
			wantOutput: "test",
			wantErr:    assert.NoError,
		},
		{
			name:       "out/error",
			prepare:    func(test *testing.T, tempFile *os.File) { os.Stdout = tempFile },
			code:       `out(['t', "hi", 's', 't'])`,
			wantResult: nil,
			wantOutput: "",
			wantErr:    assert.Error,
		},
		{
			name:       "outln/success",
			prepare:    func(test *testing.T, tempFile *os.File) { os.Stdout = tempFile },
			code:       `outln("test")`,
			wantResult: types.Nil{},
			wantOutput: "test\n",
			wantErr:    assert.NoError,
		},
		{
			name:       "outln/error",
			prepare:    func(test *testing.T, tempFile *os.File) { os.Stdout = tempFile },
			code:       `outln(['t', "hi", 's', 't'])`,
			wantResult: nil,
			wantOutput: "",
			wantErr:    assert.Error,
		},
		{
			name:       "err/success",
			prepare:    func(test *testing.T, tempFile *os.File) { os.Stderr = tempFile },
			code:       `err("test")`,
			wantResult: types.Nil{},
			wantOutput: "test",
			wantErr:    assert.NoError,
		},
		{
			name:       "err/error",
			prepare:    func(test *testing.T, tempFile *os.File) { os.Stderr = tempFile },
			code:       `err(['t', "hi", 's', 't'])`,
			wantResult: nil,
			wantOutput: "",
			wantErr:    assert.Error,
		},
		{
			name:       "errln/success",
			prepare:    func(test *testing.T, tempFile *os.File) { os.Stderr = tempFile },
			code:       `errln("test")`,
			wantResult: types.Nil{},
			wantOutput: "test\n",
			wantErr:    assert.NoError,
		},
		{
			name:       "errln/error",
			prepare:    func(test *testing.T, tempFile *os.File) { os.Stderr = tempFile },
			code:       `errln(['t', "hi", 's', 't'])`,
			wantResult: nil,
			wantOutput: "",
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			previousStdout, previousStderr := os.Stdout, os.Stderr
			defer func() { os.Stdout, os.Stderr = previousStdout, previousStderr }()

			tempFile, err := ioutil.TempFile("", "test.*")
			require.NoError(test, err)
			defer os.Remove(tempFile.Name()) // nolint: errcheck
			defer tempFile.Close()           // nolint: errcheck
			data.prepare(test, tempFile)

			ctx := context.NewDefaultContext()
			context.SetValues(ctx, Values)

			expressionAST := new(parser.Expression)
			err = parser.ParseToAST(data.code, expressionAST)
			require.NoError(test, err)

			expression, _, err := translator.TranslateExpression(expressionAST, ctx.ValuesNames())
			require.NoError(test, err)

			gotResult, gotErr := expression.Evaluate(ctx)

			err = tempFile.Close()
			require.NoError(test, err)

			gotOutputBytes, err := ioutil.ReadFile(tempFile.Name())
			require.NoError(test, err)

			assert.Equal(test, data.wantResult, gotResult)
			assert.Equal(test, data.wantOutput, string(gotOutputBytes))
			data.wantErr(test, gotErr)
		})
	}
}
