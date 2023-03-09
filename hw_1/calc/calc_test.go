package calc

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Q string
	A string
}

func TestBasic(t *testing.T) {
	assert := assert.New(t)
	testCases := []TestCase{
		{
			Q: "(1+2)-3",
			A: "0",
		},
		{
			Q: "(1+2)*3",
			A: "9",
		},
	}
	for _, tc := range testCases {
		result, err := Calc(tc.Q)
		assert.Nil(err)
		assert.Equal(result, tc.A)
	}
}

func TestOperatorsAndBraces(t *testing.T) {
	assert := assert.New(t)
	testCases := []TestCase{
		{
			Q: "1+2",
			A: "3",
		},
		{
			Q: "1-2",
			A: "-1",
		},
		{
			Q: "1*2",
			A: "2",
		},
		{
			Q: "1/2",
			A: "0",
		},
		{
			Q: "((((3*(1+2)))))",
			A: "9",
		},
		{
			Q: "-0+9999*9999/9999+1",
			A: "10000",
		},
		{
			Q: "-1+-2----1",
			A: "-2",
		},
		{
			Q: "-1+(-10*7-1+--2*0+((6/5)))",
			A: "-71",
		},
	}
	for _, tc := range testCases {
		result, err := Calc(tc.Q)
		assert.Nil(err)
		assert.Equal(result, tc.A)
	}
}

func TestZeroDivision(t *testing.T) {
	assert := assert.New(t)
	for _, tc := range []string{"22341 / 0", "22341 / (1 - 1)"} {
		_, err := Calc(tc)
		assert.ErrorIs(err, ErrZeroDivision)
	}
}

func TestInvalid(t *testing.T) {
	assert := assert.New(t)
	for _, tc := range []struct {
		Q string
		E error
	}{
		{
			Q: "",
			E: strconv.ErrSyntax,
		},
		{
			Q: "trololo",
			E: strconv.ErrSyntax,
		},
		{
			Q: "cos(0)",
			E: strconv.ErrSyntax,
		},
		{
			Q: "1+2)",
			E: ErrUnexpectedClosingBrace,
		},
		{
			Q: "(1 + 2",
			E: ErrUnclosedBrace,
		},
		{
			Q: "1  1",
			E: ErrMissingOperator,
		},
	} {
		_, err := Calc(tc.Q)
		assert.ErrorIs(err, tc.E)
	}
}
