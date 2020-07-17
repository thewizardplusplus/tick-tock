package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPairFromSlice(test *testing.T) {
	type args struct {
		items []interface{}
	}

	for _, data := range []struct {
		name string
		args args
		want *Pair
	}{
		{
			name: "nonempty slice",
			args: args{[]interface{}{"one", "two"}},
			want: &Pair{"one", &Pair{"two", nil}},
		},
		{
			name: "empty slice",
			args: args{[]interface{}{}},
			want: nil,
		},
		{
			name: "nil slice",
			args: args{nil},
			want: nil,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewPairFromSlice(data.args.items)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestNewPairFromText(test *testing.T) {
	type args struct {
		text string
	}

	for _, data := range []struct {
		name string
		args args
		want *Pair
	}{
		{
			name: "nonempty text/latin1",
			args: args{"test"},
			want: &Pair{float64('t'), &Pair{float64('e'), &Pair{float64('s'), &Pair{float64('t'), nil}}}},
		},
		{
			name: "nonempty text/not latin1",
			args: args{"тест"},
			want: &Pair{float64('т'), &Pair{float64('е'), &Pair{float64('с'), &Pair{float64('т'), nil}}}},
		},
		{
			name: "empty text",
			args: args{""},
			want: nil,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := NewPairFromText(data.args.text)

			assert.Equal(test, data.want, got)
		})
	}
}

func TestPair_Size(test *testing.T) {
	for _, data := range []struct {
		name string
		pair *Pair
		want int
	}{
		{
			name: "nonempty pair",
			pair: &Pair{"one", &Pair{"two", nil}},
			want: 2,
		},
		{
			name: "empty pair",
			pair: nil,
			want: 0,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.pair.Size()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestPair_Equals(test *testing.T) {
	type args struct {
		sample *Pair
	}

	for _, data := range []struct {
		name       string
		pair       *Pair
		args       args
		wantResult assert.BoolAssertionFunc
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success/equal",
			pair: &Pair{12.0, &Pair{23.0, nil}},
			args: args{
				sample: &Pair{12.0, &Pair{23.0, nil}},
			},
			wantResult: assert.True,
			wantErr:    assert.NoError,
		},
		{
			name: "success/not equal",
			pair: &Pair{12.0, &Pair{23.0, nil}},
			args: args{
				sample: &Pair{12.0, &Pair{42.0, nil}},
			},
			wantResult: assert.False,
			wantErr:    assert.NoError,
		},
		{
			name: "error",
			pair: &Pair{12.0, &Pair{23.0, nil}},
			args: args{
				sample: &Pair{12.0, &Pair{Nil{}, nil}},
			},
			wantResult: assert.False,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotResult, gotErr := data.pair.Equals(data.args.sample)

			data.wantResult(test, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}

func TestPair_Compare(test *testing.T) {
	type args struct {
		sample *Pair
	}

	for _, data := range []struct {
		name       string
		pair       *Pair
		args       args
		wantResult ComparisonResult
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "success/less",
			pair: &Pair{12.0, &Pair{23.0, nil}},
			args: args{
				sample: &Pair{12.0, &Pair{42.0, nil}},
			},
			wantResult: Less,
			wantErr:    assert.NoError,
		},
		{
			name: "success/shorter",
			pair: &Pair{12.0, &Pair{23.0, nil}},
			args: args{
				sample: &Pair{12.0, &Pair{23.0, &Pair{42.0, nil}}},
			},
			wantResult: Less,
			wantErr:    assert.NoError,
		},
		{
			name: "success/equal",
			pair: &Pair{12.0, &Pair{23.0, nil}},
			args: args{
				sample: &Pair{12.0, &Pair{23.0, nil}},
			},
			wantResult: Equal,
			wantErr:    assert.NoError,
		},
		{
			name: "success/greater",
			pair: &Pair{12.0, &Pair{42.0, nil}},
			args: args{
				sample: &Pair{12.0, &Pair{23.0, nil}},
			},
			wantResult: Greater,
			wantErr:    assert.NoError,
		},
		{
			name: "success/longer",
			pair: &Pair{12.0, &Pair{23.0, &Pair{42.0, nil}}},
			args: args{
				sample: &Pair{12.0, &Pair{23.0, nil}},
			},
			wantResult: Greater,
			wantErr:    assert.NoError,
		},
		{
			name: "error",
			pair: &Pair{12.0, &Pair{23.0, nil}},
			args: args{
				sample: &Pair{12.0, &Pair{Nil{}, nil}},
			},
			wantResult: 0,
			wantErr:    assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotResult, gotErr := data.pair.Compare(data.args.sample)

			assert.Equal(test, data.wantResult, gotResult)
			data.wantErr(test, gotErr)
		})
	}
}

func TestPair_Item(test *testing.T) {
	type args struct {
		index float64
	}

	for _, data := range []struct {
		name     string
		pair     *Pair
		args     args
		wantItem interface{}
		wantOk   assert.BoolAssertionFunc
	}{
		{
			name:     "success/nonempty pair/first item",
			pair:     &Pair{"one", &Pair{"two", nil}},
			args:     args{0},
			wantItem: "one",
			wantOk:   assert.True,
		},
		{
			name:     "success/nonempty pair/last item",
			pair:     &Pair{"one", &Pair{"two", nil}},
			args:     args{1},
			wantItem: "two",
			wantOk:   assert.True,
		},
		{
			name:     "error/nonempty pair/too large index",
			pair:     &Pair{"one", &Pair{"two", nil}},
			args:     args{5},
			wantItem: nil,
			wantOk:   assert.False,
		},
		{
			name:     "error/nonempty pair/negative index",
			pair:     &Pair{"one", &Pair{"two", nil}},
			args:     args{-5},
			wantItem: nil,
			wantOk:   assert.False,
		},
		{
			name:     "error/empty pair",
			pair:     nil,
			args:     args{5},
			wantItem: nil,
			wantOk:   assert.False,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotItem, gotOk := data.pair.Item(data.args.index)

			assert.Equal(test, data.wantItem, gotItem)
			data.wantOk(test, gotOk)
		})
	}
}

func TestPair_Append(test *testing.T) {
	type args struct {
		anotherPair *Pair
	}

	for _, data := range []struct {
		name         string
		pair         *Pair
		args         args
		action       func(pair *Pair)
		wantOriginal *Pair
		wantArgument *Pair
		wantResult   *Pair
	}{
		{
			name: "nonempty pair/nonempty another pair",
			pair: &Pair{1, &Pair{2, nil}},
			args: args{
				anotherPair: &Pair{3, &Pair{4, nil}},
			},
			action: func(pair *Pair) {
				for ; pair != nil; pair = pair.Tail {
					pair.Head = pair.Head.(int) * 2
				}
			},
			wantOriginal: &Pair{2, &Pair{4, nil}},
			wantArgument: &Pair{6, &Pair{8, nil}},
			wantResult:   &Pair{1, &Pair{2, &Pair{6, &Pair{8, nil}}}},
		},
		{
			name: "empty pair/nonempty another pair",
			pair: nil,
			args: args{
				anotherPair: &Pair{3, &Pair{4, nil}},
			},
			action: func(pair *Pair) {
				for ; pair != nil; pair = pair.Tail {
					pair.Head = pair.Head.(int) * 2
				}
			},
			wantOriginal: nil,
			wantArgument: &Pair{6, &Pair{8, nil}},
			wantResult:   &Pair{6, &Pair{8, nil}},
		},
		{
			name: "nonempty pair/empty another pair",
			pair: &Pair{1, &Pair{2, nil}},
			args: args{
				anotherPair: nil,
			},
			action: func(pair *Pair) {
				for ; pair != nil; pair = pair.Tail {
					pair.Head = pair.Head.(int) * 2
				}
			},
			wantOriginal: &Pair{2, &Pair{4, nil}},
			wantArgument: nil,
			wantResult:   &Pair{1, &Pair{2, nil}},
		},
		{
			name: "empty pair/empty another pair",
			pair: nil,
			args: args{
				anotherPair: nil,
			},
			action: func(pair *Pair) {
				for ; pair != nil; pair = pair.Tail {
					pair.Head = pair.Head.(int) * 2
				}
			},
			wantOriginal: nil,
			wantArgument: nil,
			wantResult:   nil,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			result := data.pair.Append(data.args.anotherPair)
			data.action(data.pair)
			data.action(data.args.anotherPair)

			assert.Equal(test, data.wantOriginal, data.pair)
			assert.Equal(test, data.wantArgument, data.args.anotherPair)
			assert.Equal(test, data.wantResult, result)
		})
	}
}

func TestPair_Slice(test *testing.T) {
	for _, data := range []struct {
		name string
		pair *Pair
		want []interface{}
	}{
		{
			name: "nonempty pair",
			pair: &Pair{"one", &Pair{"two", nil}},
			want: []interface{}{"one", "two"},
		},
		{
			name: "empty pair",
			pair: nil,
			want: nil,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.pair.Slice()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestPair_DeepSlice(test *testing.T) {
	for _, data := range []struct {
		name string
		pair *Pair
		want []interface{}
	}{
		{
			name: "nonempty pair/tree in the head",
			pair: &Pair{&Pair{"one", &Pair{"two", nil}}, &Pair{"three", nil}},
			want: []interface{}{[]interface{}{"one", "two"}, "three"},
		},
		{
			name: "nonempty pair/tree in the tail",
			pair: &Pair{"one", &Pair{&Pair{"two", &Pair{"three", nil}}, nil}},
			want: []interface{}{"one", []interface{}{"two", "three"}},
		},
		{
			name: "empty pair",
			pair: nil,
			want: nil,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			got := data.pair.DeepSlice()

			assert.Equal(test, data.want, got)
		})
	}
}

func TestPair_Text(test *testing.T) {
	for _, data := range []struct {
		name     string
		pair     *Pair
		wantText string
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success/nonempty pair/latin1",
			pair: &Pair{
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
			wantText: "test",
			wantErr:  assert.NoError,
		},
		{
			name: "success/nonempty pair/not latin1",
			pair: &Pair{
				Head: float64('т'),
				Tail: &Pair{
					Head: float64('е'),
					Tail: &Pair{
						Head: float64('с'),
						Tail: &Pair{
							Head: float64('т'),
							Tail: nil,
						},
					},
				},
			},
			wantText: "тест",
			wantErr:  assert.NoError,
		},
		{
			name:     "success/empty pair",
			pair:     nil,
			wantText: "",
			wantErr:  assert.NoError,
		},
		{
			name: "error/incorrect type",
			pair: &Pair{
				Head: float64('t'),
				Tail: &Pair{
					Head: &Pair{float64('h'), &Pair{float64('i'), nil}},
					Tail: &Pair{
						Head: float64('s'),
						Tail: &Pair{
							Head: float64('t'),
							Tail: nil,
						},
					},
				},
			},
			wantText: "",
			wantErr:  assert.Error,
		},
		{
			name: "error/incorrect rune",
			pair: &Pair{
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
			wantText: "",
			wantErr:  assert.Error,
		},
	} {
		test.Run(data.name, func(test *testing.T) {
			gotText, gotErr := data.pair.Text()

			assert.Equal(test, data.wantText, gotText)
			data.wantErr(test, gotErr)
		})
	}
}
