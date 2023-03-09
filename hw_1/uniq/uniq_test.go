package uniq

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestCase struct {
	Q string
	A string
	C Options
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
	countEntriesC := Options{}
	countEntriesC.CountEntries = true
	onlyRepeatedC := Options{}
	onlyRepeatedC.OnlyRepeated = true
	onlyUniqueC := Options{}
	onlyUniqueC.OnlyUnique = true
	ignoreCaseC := Options{}
	ignoreCaseC.IgnoreCase = true
	ignoreOneField := Options{}
	ignoreOneField.IgnoreFields = 1
	ignoreOneChar := Options{}
	ignoreOneChar.IgnoreChars = 1
	testCases := []TestCase{
		{
			Q: defaultQ,
			A: `I love music.

I love music of Kartik.
Thanks.
I love music of Kartik.
`,
			C: Options{},
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
	c := Options{}
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

func TestLogic(t *testing.T) {
	assert := assert.New(t)
	for _, tc := range allTestCases {
		us := NewUniqStrategy(tc.C)
		us.Input = strings.NewReader(tc.Q)
		result := bytes.Buffer{}
		us.Output = &result
		us.Execute()
		assert.Equal(result.String(), tc.A)
	}
}
