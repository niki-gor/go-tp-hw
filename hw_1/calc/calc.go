package calc

import (
	"errors"
	"go/token"
	"go/types"
)

var ErrInvalidQuery = errors.New("invalid query")

func Eval(s string) (result string, err error) {
	defer func() {
		if r := recover(); r != nil {
			result, err = "NaN", ErrInvalidQuery
		}
	}()
	calculation, err := types.Eval(token.NewFileSet(), nil, token.NoPos, s)
	result = calculation.Value.String()
	return
}
