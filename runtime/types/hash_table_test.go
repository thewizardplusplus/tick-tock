package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashTable_Size(test *testing.T) {
	for _, data := range []struct {
		name  string
		table HashTable
		want  int
	}{
		{
			name:  "empty",
			table: nil,
			want:  0,
		},
		{
			name:  "nonempty",
			table: HashTable{"one": "two", "three": "four"},
			want:  2,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.table.Size()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHashTable_Keys(test *testing.T) {
	for _, data := range []struct {
		name  string
		table HashTable
		want  []interface{}
	}{
		{
			name:  "Nil",
			table: HashTable{Nil{}: "test"},
			want:  []interface{}{Nil{}},
		},
		{
			name:  "float64",
			table: HashTable{23.0: "one", 42.0: "two"},
			want:  []interface{}{23.0, 42.0},
		},
		{
			name:  "string",
			table: HashTable{"one": "two", "three": "four"},
			want: []interface{}{
				&Pair{
					Head: float64('o'),
					Tail: &Pair{
						Head: float64('n'),
						Tail: &Pair{
							Head: float64('e'),
							Tail: nil,
						},
					},
				},
				&Pair{
					Head: float64('t'),
					Tail: &Pair{
						Head: float64('h'),
						Tail: &Pair{
							Head: float64('r'),
							Tail: &Pair{
								Head: float64('e'),
								Tail: &Pair{
									Head: float64('e'),
									Tail: nil,
								},
							},
						},
					},
				},
			},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.table.Keys()

			assert.ElementsMatch(test, data.want, got)
		})
	}
}

func TestHashTable_Equals(test *testing.T) {
	type args struct {
		sample HashTable
	}

	for _, data := range []struct {
		name       string
		table      HashTable
		args       args
		wantResult assert.BoolAssertionFunc
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name:  "success/equal",
			table: HashTable{"one": 12.0, "two": 23.0},
			args: args{
				sample: HashTable{"one": 12.0, "two": 23.0},
			},
			wantResult: assert.True,
			wantErr:    assert.NoError,
		},
		{
			name:  "success/not equal/by keys",
			table: HashTable{"one": 12.0, "two": 23.0},
			args: args{
				sample: HashTable{"one": 12.0, "three": 23.0},
			},
			wantResult: assert.False,
			wantErr:    assert.NoError,
		},
		{
			name:  "success/not equal/by values",
			table: HashTable{"one": 12.0, "two": 23.0},
			args: args{
				sample: HashTable{"one": 12.0, "two": 42.0},
			},
			wantResult: assert.False,
			wantErr:    assert.NoError,
		},
		{
			name:  "success/not equal/shorter",
			table: HashTable{"one": 12.0, "two": 23.0},
			args: args{
				sample: HashTable{"one": 12.0, "two": 23.0, "three": 42.0},
			},
			wantResult: assert.False,
			wantErr:    assert.NoError,
		},
		{
			name:  "success/not equal/longer",
			table: HashTable{"one": 12.0, "two": 23.0, "three": 42.0},
			args: args{
				sample: HashTable{"one": 12.0, "two": 23.0},
			},
			wantResult: assert.False,
			wantErr:    assert.NoError,
		},
		{
			name:  "error",
			table: HashTable{"one": 12.0, "two": func() {}},
			args: args{
				sample: HashTable{"one": 12.0, "two": 23.0},
			},
			wantResult: assert.False,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotResult, gotErr := data.table.Equals(data.args.sample)

			data.wantResult(test, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}

func TestHashTable_Item(test *testing.T) {
	type args struct {
		key interface{}
	}

	for _, data := range []struct {
		name      string
		table     HashTable
		args      args
		wantValue interface{}
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:  "empty",
			table: nil,
			args: args{
				key: &Pair{
					Head: float64('o'),
					Tail: &Pair{
						Head: float64('n'),
						Tail: &Pair{
							Head: float64('e'),
							Tail: nil,
						},
					},
				},
			},
			wantValue: nil,
			wantErr: func(test assert.TestingT, err error, args ...interface{}) bool {
				return assert.Equal(test, ErrNotFound, err)
			},
		},
		{
			name:  "nonempty/existing key",
			table: HashTable{"one": "two", "three": "four"},
			args: args{
				key: &Pair{
					Head: float64('o'),
					Tail: &Pair{
						Head: float64('n'),
						Tail: &Pair{
							Head: float64('e'),
							Tail: nil,
						},
					},
				},
			},
			wantValue: "two",
			wantErr:   assert.NoError,
		},
		{
			name:  "nonempty/nonexistent key",
			table: HashTable{"one": "two", "three": "four"},
			args: args{
				key: &Pair{
					Head: float64('f'),
					Tail: &Pair{
						Head: float64('i'),
						Tail: &Pair{
							Head: float64('v'),
							Tail: &Pair{
								Head: float64('e'),
								Tail: nil,
							},
						},
					},
				},
			},
			wantValue: nil,
			wantErr: func(test assert.TestingT, err error, args ...interface{}) bool {
				return assert.Equal(test, ErrNotFound, err)
			},
		},
		{
			name:  "incorrect key",
			table: HashTable{"one": "two", "three": "four"},
			args: args{
				key: &Pair{
					Head: -23.0,
					Tail: &Pair{
						Head: float64('n'),
						Tail: &Pair{
							Head: float64('e'),
							Tail: nil,
						},
					},
				},
			},
			wantValue: nil,
			wantErr:   assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotValue, gotErr := data.table.Item(data.args.key)

			assert.Equal(test, data.wantValue, gotValue)
			data.wantErr(test, gotErr)
		})
	}
}

func TestHashTable_With(test *testing.T) {
	type args struct {
		key   interface{}
		value interface{}
	}

	for _, data := range []struct {
		name      string
		table     HashTable
		args      args
		wantTable HashTable
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:  "success/empty",
			table: nil,
			args: args{
				key: &Pair{
					Head: float64('o'),
					Tail: &Pair{
						Head: float64('n'),
						Tail: &Pair{
							Head: float64('e'),
							Tail: nil,
						},
					},
				},
				value: "two",
			},
			wantTable: HashTable{"one": "two"},
			wantErr:   assert.NoError,
		},
		{
			name:  "success/nonempty/existing key",
			table: HashTable{"one": "two", "three": "four"},
			args: args{
				key: &Pair{
					Head: float64('o'),
					Tail: &Pair{
						Head: float64('n'),
						Tail: &Pair{
							Head: float64('e'),
							Tail: nil,
						},
					},
				},
				value: "five",
			},
			wantTable: HashTable{"one": "five", "three": "four"},
			wantErr:   assert.NoError,
		},
		{
			name:  "success/nonempty/nonexistent key",
			table: HashTable{"one": "two", "three": "four"},
			args: args{
				key: &Pair{
					Head: float64('f'),
					Tail: &Pair{
						Head: float64('i'),
						Tail: &Pair{
							Head: float64('v'),
							Tail: &Pair{
								Head: float64('e'),
								Tail: nil,
							},
						},
					},
				},
				value: "six",
			},
			wantTable: HashTable{"one": "two", "three": "four", "five": "six"},
			wantErr:   assert.NoError,
		},
		{
			name:  "success/nonempty/Nil value",
			table: HashTable{"one": "two", "three": "four"},
			args: args{
				key: &Pair{
					Head: float64('o'),
					Tail: &Pair{
						Head: float64('n'),
						Tail: &Pair{
							Head: float64('e'),
							Tail: nil,
						},
					},
				},
				value: Nil{},
			},
			wantTable: HashTable{"three": "four"},
			wantErr:   assert.NoError,
		},
		{
			name:  "error",
			table: HashTable{"one": "two", "three": "four"},
			args: args{
				key: &Pair{
					Head: -23.0,
					Tail: &Pair{
						Head: float64('n'),
						Tail: &Pair{
							Head: float64('e'),
							Tail: nil,
						},
					},
				},
				value: "five",
			},
			wantTable: nil,
			wantErr:   assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotTable, gotErr := data.table.With(data.args.key, data.args.value)

			assert.Equal(test, data.wantTable, gotTable)
			data.wantErr(test, gotErr)
		})
	}
}

func TestHashTable_Merge(test *testing.T) {
	type args struct {
		anotherTable HashTable
	}

	for _, data := range []struct {
		name  string
		table HashTable
		args  args
		want  HashTable
	}{
		{
			name:  "both are empty",
			table: nil,
			args: args{
				anotherTable: nil,
			},
			want: HashTable{},
		},
		{
			name:  "first is nonempty",
			table: HashTable{"one": "two", "three": "four"},
			args: args{
				anotherTable: nil,
			},
			want: HashTable{"one": "two", "three": "four"},
		},
		{
			name:  "second is nonempty",
			table: nil,
			args: args{
				anotherTable: HashTable{"five": "six", "seven": "eight"},
			},
			want: HashTable{"five": "six", "seven": "eight"},
		},
		{
			name:  "both are nonempty",
			table: HashTable{"one": "two", "three": "four"},
			args: args{
				anotherTable: HashTable{"five": "six", "seven": "eight"},
			},
			want: HashTable{"one": "two", "three": "four", "five": "six", "seven": "eight"},
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.table.Merge(data.args.anotherTable)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestHashTable_DeepMap(test *testing.T) {
	for _, data := range []struct {
		name      string
		table     HashTable
		wantTable map[string]interface{}
		wantErr   assert.ErrorAssertionFunc
	}{
		{
			name:      "success",
			table:     HashTable{"one": "two", "three": "four"},
			wantTable: map[string]interface{}{"one": "two", "three": "four"},
			wantErr:   assert.NoError,
		},
		{
			name: "success with a pair",
			table: HashTable{
				"test": &Pair{
					Head: float64('t'),
					Tail: &Pair{
						Head: float64('e'),
						Tail: &Pair{
							Head: float64('s'),
							Tail: &Pair{
								Head: float64('t'),
								Tail: nil,
							},
						},
					},
				},
			},
			wantTable: map[string]interface{}{
				"test": []interface{}{float64('t'), float64('e'), float64('s'), float64('t')},
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with a hash table",
			table: HashTable{
				"test": HashTable{"one": "two", "three": "four"},
			},
			wantTable: map[string]interface{}{
				"test": map[string]interface{}{"one": "two", "three": "four"},
			},
			wantErr: assert.NoError,
		},
		{
			name: "success with a hash table that contains a pair",
			table: HashTable{
				"one": HashTable{
					"two": &Pair{
						Head: float64('t'),
						Tail: &Pair{
							Head: float64('h'),
							Tail: &Pair{
								Head: float64('r'),
								Tail: &Pair{
									Head: float64('e'),
									Tail: &Pair{
										Head: float64('e'),
										Tail: nil,
									},
								},
							},
						},
					},
				},
			},
			wantTable: map[string]interface{}{
				"one": map[string]interface{}{
					"two": []interface{}{float64('t'), float64('h'), float64('r'), float64('e'), float64('e')},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name:      "error with a key",
			table:     HashTable{23.0: "one", 42.0: "two"},
			wantTable: nil,
			wantErr:   assert.Error,
		},
		{
			name: "error with a value (list)",
			table: HashTable{
				"one":   "two",
				"three": &Pair{"four", &Pair{HashTable{23.0: "five", 42.0: "six"}, nil}},
			},
			wantTable: nil,
			wantErr:   assert.Error,
		},
		{
			name:      "error with a value (hash table)",
			table:     HashTable{"one": "two", "three": HashTable{23.0: "four", 42.0: "five"}},
			wantTable: nil,
			wantErr:   assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotTable, gotErr := data.table.DeepMap()

			assert.Equal(test, data.wantTable, gotTable)
			data.wantErr(test, gotErr)
		})
	}
}

func Test_prepareKey(test *testing.T) {
	type args struct {
		key interface{}
	}

	for _, data := range []struct {
		name           string
		args           args
		wantPrepareKey interface{}
		wantErr        assert.ErrorAssertionFunc
	}{
		{
			name: "success/Nil",
			args: args{
				key: Nil{},
			},
			wantPrepareKey: Nil{},
			wantErr:        assert.NoError,
		},
		{
			name: "success/float64",
			args: args{
				key: 23.0,
			},
			wantPrepareKey: 23.0,
			wantErr:        assert.NoError,
		},
		{
			name: "success/*Pair",
			args: args{
				key: &Pair{
					Head: float64('t'),
					Tail: &Pair{
						Head: float64('e'),
						Tail: &Pair{
							Head: float64('s'),
							Tail: &Pair{
								Head: float64('t'),
								Tail: nil,
							},
						},
					},
				},
			},
			wantPrepareKey: "test",
			wantErr:        assert.NoError,
		},
		{
			name: "error/incorrect type",
			args: args{
				key: func() {},
			},
			wantPrepareKey: nil,
			wantErr:        assert.Error,
		},
		{
			name: "error/incorrect rune",
			args: args{
				key: &Pair{
					Head: float64('t'),
					Tail: &Pair{
						Head: -23.0,
						Tail: &Pair{
							Head: float64('s'),
							Tail: &Pair{
								Head: float64('t'),
								Tail: nil,
							},
						},
					},
				},
			},
			wantPrepareKey: nil,
			wantErr:        assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotPrepareKey, gotErr := prepareKey(data.args.key)

			assert.Equal(test, data.wantPrepareKey, gotPrepareKey)
			data.wantErr(test, gotErr)
		})
	}
}
