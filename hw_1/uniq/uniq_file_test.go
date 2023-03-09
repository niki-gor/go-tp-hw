package uniq

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFiles(t *testing.T) {
	assert := assert.New(t)
	for _, tc := range allTestCases {
		input, _ := os.CreateTemp("", "*")
		output, _ := os.CreateTemp("", "*")
		defer os.Remove(input.Name())
		defer os.Remove(output.Name())

		input.WriteString(tc.Q) // nolint

		tc.C.InputPath = input.Name()
		tc.C.OutputPath = output.Name()

		us := NewUniqStrategy(tc.C)
		us.Execute()

		result, err := io.ReadAll(output)
		if err != nil {
			t.Fatal()
		}

		assert.Equal(string(result), tc.A)
	}
}
