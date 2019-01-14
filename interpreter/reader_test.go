package interpreter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsEmptyFilename(test *testing.T) {
	type args struct {
		filename string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "usual filename",
			args: args{"test"},
		},
		{
			name: "special filename",
			args: args{"-"},
			want: true,
		},
		{
			name: "empty filename",
			want: true,
		},
	}
	for _, testData := range tests {
		test.Run(testData.name, func(test *testing.T) {
			got := isEmptyFilename(testData.args.filename)
			assert.Equal(test, testData.want, got)
		})
	}
}
