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

func ArgsString() string {
	return strings.Join(os.Args[1:], " ")
}

func ReadAllStdin() string {
	result, err := io.ReadAll(os.Stdin)
	common.Exit1OnError(err)
	return string(result)
}

func ReadLineStdin() string {
	reader := bufio.NewReader(os.Stdin)
	result, err := reader.ReadString('\n')
	common.Exit1OnError(err)
	return result
}

func Eval(s string) string {
	result, err := types.Eval(token.NewFileSet(), nil, token.NoPos, s)
	common.Exit1OnError(err)
	return result.Value.String()
}

func ChooseInput() (result func() string) {
	fi, err := os.Stdin.Stat()
	common.Exit1OnError(err)
	stdinUsed := fi.Mode()&os.ModeNamedPipe != 0
	argsUsed := len(os.Args) > 1

	switch {
	case stdinUsed && argsUsed:
		common.Exit1OnError(ErrMultipleInputSources)
	case stdinUsed:
		result = ReadAllStdin
	case argsUsed:
		result = ArgsString
	default:
		result = ReadLineStdin
	}
	return result
}

func main() {
	inputFunc := ChooseInput()
	fmt.Print(Eval(inputFunc()))
}
