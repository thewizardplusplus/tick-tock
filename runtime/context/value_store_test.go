package context

import (
	"testing"

	mapset "github.com/deckarep/golang-set"
	"github.com/stretchr/testify/assert"
)

func TestDefaultValueStore_ValuesNames(test *testing.T) {
	for _, data := range []struct {
		name  string
		store DefaultValueStore
		want  mapset.Set
	}{
		{
			name:  "with values",
			store: DefaultValueStore{"one": 1, "two": 2},
			want:  mapset.NewSet("one", "two"),
		},
		{
			name:  "without values",
			store: DefaultValueStore{},
			want:  mapset.NewSet(),
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.store.ValuesNames()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestDefaultValueStore_Value(test *testing.T) {
	type args struct {
		name string
	}

	for _, data := range []struct {
		name      string
		store     DefaultValueStore
		args      args
		wantValue interface{}
		wantOk    assert.BoolAssertionFunc
	}{
		{
			name:      "existent value",
			store:     DefaultValueStore{"one": 1, "two": 2},
			args:      args{"two"},
			wantValue: 2,
			wantOk:    assert.True,
		},
		{
			name:      "not existent value",
			store:     DefaultValueStore{"one": 1, "two": 2},
			args:      args{"three"},
			wantValue: nil,
			wantOk:    assert.False,
		},
		{
			name:      "without values",
			store:     DefaultValueStore{},
			args:      args{"one"},
			wantValue: nil,
			wantOk:    assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotValue, gotOk := data.store.Value(data.args.name)

			assert.Equal(test, data.wantValue, gotValue)
			data.wantOk(test, gotOk)
		})
	}
}

func TestDefaultValueStore_SetValue(test *testing.T) {
	type args struct {
		name  string
		value interface{}
	}

	for _, data := range []struct {
		name      string
		store     DefaultValueStore
		args      args
		wantStore DefaultValueStore
	}{
		{
			name:      "add a new value with the same type",
			store:     DefaultValueStore{"one": 1, "two": 2},
			args:      args{"three", 3},
			wantStore: DefaultValueStore{"one": 1, "two": 2, "three": 3},
		},
		{
			name:      "add a new value with a different type",
			store:     DefaultValueStore{"one": 1, "two": 2},
			args:      args{"three", "3"},
			wantStore: DefaultValueStore{"one": 1, "two": 2, "three": "3"},
		},
		{
			name:      "update an old value with the same type",
			store:     DefaultValueStore{"one": 1, "two": 2},
			args:      args{"two", -2},
			wantStore: DefaultValueStore{"one": 1, "two": -2},
		},
		{
			name:      "update an old value with a different type",
			store:     DefaultValueStore{"one": 1, "two": 2},
			args:      args{"two", "2"},
			wantStore: DefaultValueStore{"one": 1, "two": "2"},
		},
		{
			name:      "without values",
			store:     DefaultValueStore{},
			args:      args{"one", 1},
			wantStore: DefaultValueStore{"one": 1},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			data.store.SetValue(data.args.name, data.args.value)

			assert.Equal(test, data.wantStore, data.store)
		})
	}
}

func TestDefaultValueStore_Copy(test *testing.T) {
	for _, data := range []struct {
		name         string
		store        DefaultValueStore
		postprocess  func(store DefaultValueStore)
		wantOriginal DefaultValueStore
		wantCopy     DefaultValueStore
	}{
		{
			name:         "with values",
			store:        DefaultValueStore{"one": 1, "two": 2},
			postprocess:  func(store DefaultValueStore) { store.SetValue("three", 3) },
			wantOriginal: DefaultValueStore{"one": 1, "two": 2, "three": 3},
			wantCopy:     DefaultValueStore{"one": 1, "two": 2},
		},
		{
			name:         "without values",
			store:        DefaultValueStore{},
			postprocess:  func(store DefaultValueStore) { store.SetValue("one", 1) },
			wantOriginal: DefaultValueStore{"one": 1},
			wantCopy:     DefaultValueStore{},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.store.Copy()
			data.postprocess(data.store)

			assert.Equal(test, data.wantOriginal, data.store)
			assert.Equal(test, data.wantCopy, got)
		})
	}
}
