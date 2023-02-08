package common

import (
	"fmt"
	"os"
)

func Exit1OnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
