package main

import (
	"bytes"
	"strings"
	"testing"
)

type TestCase struct {
	Q string
	A string
	C Config
}

func TestBasic(t *testing.T) {
	defaultQ := `I love music.
I love music.
I love music.

I love music of Kartik.
I love music of Kartik.
Thanks.
I love music of Kartik.
I love music of Kartik.`
	countEntriesC := Config{}
	countEntriesC.CountEntries = true
	onlyRepeatedC := Config{}
	onlyRepeatedC.OnlyRepeated = true
	onlyUniqueC := Config{}
	onlyUniqueC.OnlyUnique = true
	ignoreCaseC := Config{}
	ignoreCaseC.IgnoreCase = true
	ignoreOneField := Config{}
	ignoreOneField.IgnoreFields = 1
	ignoreOneChar := Config{}
	ignoreOneChar.IgnoreChars = 1
	testCases := []TestCase{
		{
			Q: defaultQ,
			A: `I love music.

I love music of Kartik.
Thanks.
I love music of Kartik.`,
			C: Config{},
		},
		{
			Q: defaultQ,
			A: `3 I love music.
1 
2 I love music of Kartik.
1 Thanks.
2 I love music of Kartik.`,
			C: countEntriesC,
		},
		{
			Q: defaultQ,
			A: `I love music.
I love music of Kartik.
I love music of Kartik.`,
			C: onlyRepeatedC,
		},
		{
			Q: defaultQ,
			A: `
Thanks.`,
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
I love music of kartik.`,
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
Thanks.`,
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
Thanks.`,
			C: ignoreOneChar,
		},
	}
	for i, tc := range testCases {
		us := NewUniqStrategy(tc.C)
		us.Reader = strings.NewReader(tc.Q)
		result := bytes.Buffer{}
		us.Writer = &result
		us.Execute()
		if result.String() != tc.A+"\n" {
			t.Errorf("testcase %d\nquery: %s\nresult:\n%s\nexpected:\n%s", i+1, tc.Q, result.String(), tc.A+"\n")
		}
	}
}
