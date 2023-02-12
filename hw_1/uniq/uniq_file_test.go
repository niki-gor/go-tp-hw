package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestFiles(t *testing.T) {
	for _, tc := range allTestCases {
		input, _ := os.CreateTemp("", "*")
		output, _ := os.CreateTemp("", "*")
		defer os.Remove(input.Name())
		defer os.Remove(output.Name())

		input.WriteString(tc.Q)

		tc.C.InputPath = input.Name()
		tc.C.OutputPath = output.Name()
		uniq := exec.Command("/usr/bin/go", append([]string{"run", "uniq.go"}, toCli(tc.C)...)...)
		uniq.Run()

		result, err := os.ReadFile(output.Name())
		if err != nil {
			t.Fatal()
		}

		if string(result) != tc.A {
			t.Errorf("result: %s\nexpected: %s\n", string(result), tc.A)
		}
	}
}
