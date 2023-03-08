package calc

import (
	"testing"
)

type TestCase struct {
	Q string
	A string
}

func TestBasic(t *testing.T) {
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
		if err != nil {
			t.Fatal("error where it should not be")
		}
		if result != tc.A {
			t.Error("not eq")
		}
	}
}

func TestOperatorsAndBraces(t *testing.T) {
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
		if err != nil {
			t.Fatal("error where it should not be")
		}
		if result != tc.A {
			t.Error("not eq")
		}
	}
}

func TestZeroDivision(t *testing.T) {
	_, err := Calc("22341 / 0")
	if err == nil {
		t.Fatal("no error when zero division occurred")
	}
	_, err = Calc("22341 / (1 - 1)")
	if err == nil {
		t.Fatal("no error when zero division occurred")
	}
}

func TestInvalid(t *testing.T) {
	for _, tc := range []string{"", "trololo", "1+2)", "cos(0)"} {
		_, err := Calc(tc)
		if err == nil {
			t.Fatal("no error when it should be")
		}
	}
}
