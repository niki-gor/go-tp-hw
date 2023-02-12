package main

import (
	"bytes"
	"os/exec"
	"strconv"
	"testing"
)

func toCli(c Config) []string {
	result := []string{}
	if c.CountEntries {
		result = append(result, "-c")
	}
	if c.OnlyRepeated {
		result = append(result, "-d")
	}
	if c.OnlyUnique {
		result = append(result, "-u")
	}
	result = append(result, "-f", strconv.Itoa(c.IgnoreFields))
	result = append(result, "-s", strconv.Itoa(c.IgnoreChars))
	if c.IgnoreCase {
		result = append(result, "-i")
	}
	if c.InputPath != "" {
		result = append(result, c.InputPath)
		if c.OutputPath != "" {
			result = append(result, c.OutputPath)
		}
	}
	return result
}

func TestCli(t *testing.T) {
	for _, tc := range allTestCases {
		echo := exec.Command("echo", tc.Q)
		uniq := exec.Command("/usr/bin/go", append([]string{"run", "uniq.go"}, toCli(tc.C)...)...)

		uniq.Stdin, _ = echo.StdoutPipe()
		result := bytes.Buffer{}
		uniq.Stdout = &result

		uniq.Start()
		echo.Run()
		uniq.Wait()

		if result.String() != tc.A {
			t.Errorf("result: %s\nexpected: %s\n", result.String(), tc.A)
		}
	}
}
