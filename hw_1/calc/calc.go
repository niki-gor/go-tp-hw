package main

import (
	"bufio"
	"errors"
	"fmt"
	"go/token"
	"go/types"
	"go_tp/hw_1/common"
	"io"
	"os"
	"strings"
)

var (
	ErrMultipleInputSources = errors.New("Both args and stdin are used - you can choose only one option")
)

func ArgsString() (string, error) {
	return strings.Join(os.Args[1:], " "), nil
}

func ReadAllStdin() (string, error) {
	result, err := io.ReadAll(os.Stdin)
	return string(result), err
}

func ReadLineStdin() (string, error) {
	return bufio.NewReader(os.Stdin).ReadString('\n')
}

func Eval(s string) (string, error) {
	result, err := types.Eval(token.NewFileSet(), nil, token.NoPos, s)
	return result.Value.String(), err
}

func SelectInput() (func() (string, error), error) {
	fi, err := os.Stdin.Stat()
	common.Exit1OnError(err)
	stdinUsed := fi.Mode()&os.ModeNamedPipe != 0
	argsUsed := len(os.Args) > 1

	switch {
	case stdinUsed && argsUsed:
		return nil, ErrMultipleInputSources
	case stdinUsed:
		return ReadAllStdin, nil
	case argsUsed:
		return ArgsString, nil
	}
	return ReadAllStdin, nil
}

func main() {
	input, err := SelectInput()
	common.Exit1OnError(err)
	s, err := input()
	common.Exit1OnError(err)
	result, err := Eval(s)
	common.Exit1OnError(err)
	fmt.Print(result)
}
