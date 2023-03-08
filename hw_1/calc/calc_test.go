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
	_, err := Calc("22341 / 0")
	assert.ErrorIs(err, ErrZeroDivision)
	_, err = Calc("22341 / (1 - 1)")
	assert.ErrorIs(err, ErrZeroDivision)
}

func TestInvalid(t *testing.T) {
	assert := assert.New(t)
	for _, tc := range []string{"", "trololo", "cos(0)"} {
		_, err := Calc(tc)
		assert.ErrorIs(err, strconv.ErrSyntax)
	}
	_, err := Calc("1+2)")
	assert.ErrorIs(err, ErrUnexpectedClosingBrace)
	_, err = Calc("(1 + 2")
	assert.ErrorIs(err, ErrUnclosedBrace)
	_, err = Calc("1  1")
	assert.ErrorIs(err, ErrMissingOperator)
}
