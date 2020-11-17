package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(test *testing.T) {
	type args struct {
		code string
	}

	for _, testData := range []struct {
		name    string
		args    args
		want    *Program
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{"actor One(); actor Two();"},
			want: &Program{
				Definitions: []*Definition{{Actor: &Actor{"One", nil, nil}}, {Actor: &Actor{"Two", nil, nil}}},
			},
			wantErr: assert.NoError,
		},
		{
			name:    "error",
			args:    args{"incorrect"},
			wantErr: assert.Error,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			got, err := Parse(testData.args.code)

			assert.Equal(test, testData.want, got)
			testData.wantErr(test, err)
		})
	}
}
