package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashTableGet(test *testing.T) {
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
			name:  "success/existing key",
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
			name:  "success/nonexistent key",
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
			wantValue: Nil{},
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
			},
			wantValue: nil,
			wantErr:   assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotValue, gotErr := data.table.Get(data.args.key)

			assert.Equal(test, data.wantValue, gotValue)
			data.wantErr(test, gotErr)
		})
	}
}

func TestHashTableSet(test *testing.T) {
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
			name:  "success/existing key",
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
			name:  "success/nonexistent key",
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
			name:  "success/Nil key",
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
			name:  "error/incorrect rune",
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
			wantTable: HashTable{"one": "two", "three": "four"},
			wantErr:   assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotErr := data.table.Set(data.args.key, data.args.value)

			assert.Equal(test, data.wantTable, data.table)
			data.wantErr(test, gotErr)
		})
	}
}

func TestPrepareKey(test *testing.T) {
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
