package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
