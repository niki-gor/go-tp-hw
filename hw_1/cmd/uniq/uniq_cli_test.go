package main

import (
	"bytes"
	"go_tp/hw_1/uniq"
	"os/exec"
	"strconv"
	"testing"
)

type TestCase struct {
	Q string
	A string
	C uniq.Options
}

func basicTestCases() []TestCase {
	defaultQ := `I love music.
I love music.
I love music.

I love music of Kartik.
I love music of Kartik.
Thanks.
I love music of Kartik.
I love music of Kartik.`
	countEntriesC := uniq.Options{}
	countEntriesC.CountEntries = true
	onlyRepeatedC := uniq.Options{}
	onlyRepeatedC.OnlyRepeated = true
	onlyUniqueC := uniq.Options{}
	onlyUniqueC.OnlyUnique = true
	ignoreCaseC := uniq.Options{}
	ignoreCaseC.IgnoreCase = true
	ignoreOneField := uniq.Options{}
	ignoreOneField.IgnoreFields = 1
	ignoreOneChar := uniq.Options{}
	ignoreOneChar.IgnoreChars = 1
	testCases := []TestCase{
		{
			Q: defaultQ,
			A: `I love music.

I love music of Kartik.
Thanks.
I love music of Kartik.
`,
			C: uniq.Options{},
		},
		{
			Q: defaultQ,
			A: `3 I love music.
1 
2 I love music of Kartik.
1 Thanks.
2 I love music of Kartik.
`,
			C: countEntriesC,
		},
		{
			Q: defaultQ,
			A: `I love music.
I love music of Kartik.
I love music of Kartik.
`,
			C: onlyRepeatedC,
		},
		{
			Q: defaultQ,
			A: `
Thanks.
`,
			C: onlyUniqueC,
		},
		{
			Q: `I LOVE MUSIC.
I love music.
I LoVe MuSiC.

I love MuSIC of Kartik.
I love music of kartik.
Thanks.
I love music of kartik.
I love MuSIC of Kartik.`,
			A: `I LOVE MUSIC.

I love MuSIC of Kartik.
Thanks.
I love music of kartik.
`,
			C: ignoreCaseC,
		},
		{
			Q: `We love music.
I love music.
They love music.

I love music of Kartik.
We love music of Kartik.
Thanks.`,
			A: `We love music.

I love music of Kartik.
Thanks.
`,
			C: ignoreOneField,
		},
		{
			Q: `I love music.
A love music.
C love music.

I love music of Kartik.
We love music of Kartik.
Thanks.`,
			A: `I love music.

I love music of Kartik.
We love music of Kartik.
Thanks.
`,
			C: ignoreOneChar,
		},
	}
	return testCases
}

func combinedTestCase() TestCase {
	c := uniq.Options{}
	c.IgnoreFields = 1 // не учитывается первое поле
	c.IgnoreChars = 1  // и первый символ 2-го поля
	tc := TestCase{
		Q: `We love music.
I move music.
They dove music.

I love music of Kartik.
We move music of Kartik.
Thanks.`,
		A: `We love music.

I love music of Kartik.
Thanks.
`,
		C: c,
	}
	return tc
}

var allTestCases = append(basicTestCases(), combinedTestCase())

func toCli(c uniq.Options) []string {
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
		// тест *работает* на *nix с установленным компилятором go
		echo := exec.Command("echo", tc.Q)
		uniq := exec.Command("go", append([]string{"run", "main.go"}, toCli(tc.C)...)...)

		uniq.Stdin, _ = echo.StdoutPipe()
		result := bytes.Buffer{}
		uniq.Stdout = &result

		uniq.Start() // nolint
		echo.Run()   // nolint
		uniq.Wait()  // nolint

		if result.String() != tc.A {
			t.Errorf("result: %s\nexpected: %s\n", result.String(), tc.A)
		}
	}
}
