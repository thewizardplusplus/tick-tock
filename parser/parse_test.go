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

func TestParseToAST(test *testing.T) {
	type (
		args struct {
			code string
			ast  interface{}
		}
		testAST struct {
			Number int `parser:"@Int"`
		}
	)

	for _, testData := range []struct {
		name    string
		args    args
		wantAST interface{}
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "testAST/success",
			args:    args{"23", new(testAST)},
			wantAST: &testAST{23},
			wantErr: assert.NoError,
		},
		{
			name:    "testAST/error/building",
			args:    args{"23", "incorrect"},
			wantAST: "incorrect",
			wantErr: assert.Error,
		},
		{
			name:    "testAST/error/parsing",
			args:    args{"incorrect", new(testAST)},
			wantAST: new(testAST),
			wantErr: assert.Error,
		},
		{
			name: "comment/line",
			args: args{
				code: "actor One();\n// actor Two();\nactor Three();",
				ast:  new(Program),
			},
			wantAST: &Program{
				Definitions: []*Definition{{Actor: &Actor{Name: "One"}}, {Actor: &Actor{Name: "Three"}}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "comment/block",
			args: args{
				code: "actor One(); /* actor Two(); */ actor Three();",
				ast:  new(Program),
			},
			wantAST: &Program{
				Definitions: []*Definition{{Actor: &Actor{Name: "One"}}, {Actor: &Actor{Name: "Three"}}},
			},
			wantErr: assert.NoError,
		},
	} {
		test.Run(testData.name, func(test *testing.T) {
			err := parseToAST(testData.args.code, testData.args.ast)
			assert.Equal(test, testData.wantAST, testData.args.ast)
			testData.wantErr(test, err)
		})
	}
}
