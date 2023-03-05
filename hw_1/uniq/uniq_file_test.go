package uniq

import (
	"io"
	"os"
	"testing"
)

func TestFiles(t *testing.T) {
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

		if string(result) != tc.A {
			t.Errorf("result: %s\nexpected: %s\n", string(result), tc.A)
		}
	}
}
