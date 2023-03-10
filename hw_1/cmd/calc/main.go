package main

import (
	"bufio"
	"errors"
	"fmt"
	"go_tp/hw_1/calc"
	"io"
	"log"
	"os"
	"strings"
)

var (
	ErrMultipleInputSources = errors.New("both args and stdin are used - you can choose only one option")
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

func SelectInput() (func() (string, error), error) {
	fi, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
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
	if err != nil {
		log.Fatal(err)
	}
	s, err := input()
	if err != nil {
		log.Fatal(err)
	}
	result, err := calc.Calc(s)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(result)
}
