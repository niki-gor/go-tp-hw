package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"go_tp/hw_1/common"
	"io"
	"os"
	"strings"

	"golang.org/x/exp/slices"
)

var (
	ErrIncompatibleFlags = errors.New("флаги -c, -d и -u несовместимы - вы можете выбрать только один из них")
)

type Config struct {
	CountEntries bool
	OnlyRepeated bool
	OnlyUnique   bool

	IgnoreFields int
	IgnoreChars  int

	IgnoreCase bool

	InputPath  string
	OutputPath string
}

func (c *Config) ParseFromFlags() {
	flag.BoolVar(&c.CountEntries, "c", false, "подсчитать количество встречаний строки во входных данных. Вывести это число перед строкой отделив пробелом.")
	flag.BoolVar(&c.OnlyRepeated, "d", false, "вывести только те строки, которые повторились во входных данных.")
	flag.BoolVar(&c.OnlyUnique, "u", false, "вывести только те строки, которые не повторились во входных данных.")

	flag.IntVar(&c.IgnoreFields, "f", 0, "не учитывать первые num_fields полей в строке. Полем в строке является непустой набор символов отделённый пробелом.")
	flag.IntVar(&c.IgnoreChars, "s", 0, "не учитывать первые num_chars символов в строке. При использовании вместе с параметром -f учитываются первые символы после num_fields полей (не учитывая пробел-разделитель после последнего поля).")

	flag.BoolVar(&c.IgnoreCase, "i", false, "не учитывать регистр букв.")

	flag.Parse()

	c.InputPath = flag.Arg(0)
	c.OutputPath = flag.Arg(1)
}

func (c *Config) Validate() error {
	b2i := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}
	if b2i(c.CountEntries)+b2i(c.OnlyRepeated)+b2i(c.OnlyUnique) > 1 {
		return ErrIncompatibleFlags
	}

	if c.InputPath != "" {
		f, err := os.Open(c.InputPath)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	if c.OutputPath != "" {
		f, err := os.Create(c.OutputPath)
		if err != nil {
			return err
		}
		defer f.Close()
	}

	return nil
}

type UniqStrategy struct {
	Reader    io.Reader
	AreEqual  func(string, string) bool
	CountIsOk func(int) bool
	Format    func(int, string) string
	Writer    io.Writer
}

func (us *UniqStrategy) Execute() {
	prev, latest := "", ""
	count := 0
	scanner := bufio.NewScanner(us.Reader)

	if scanner.Scan() {
		prev = scanner.Text()
		count = 1
	}
	for scanner.Scan() {
		latest = scanner.Text()
		if us.AreEqual(prev, latest) {
			count++
			continue
		}
		if us.CountIsOk(count) {
			fmt.Fprintln(us.Writer, us.Format(count, prev))
		}
		prev = latest
		count = 1
	}
	if us.CountIsOk(count) {
		fmt.Fprintln(us.Writer, us.Format(count, prev))
	}
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func Drop[T any](arr []T, n int) []T {
	return arr[Min(len(arr), n):]
}

func AreEqual(s1, s2 string, ignoreFields, ignoreChars int) bool {
	first := Drop(strings.Fields(s1), ignoreFields)
	second := Drop(strings.Fields(s2), ignoreFields)

	dropRunes := func(words []string) []rune {
		return Drop([]rune(strings.Join(words, " ")), ignoreChars)
	}
	return slices.Equal(dropRunes(first), dropRunes(second))
}

func NewUniqStrategy(c Config) *UniqStrategy {
	us := &UniqStrategy{}

	if c.InputPath != "" {
		us.Reader, _ = os.Open(c.InputPath)
	} else {
		us.Reader = os.Stdin
	}

	if c.IgnoreCase {
		us.AreEqual = func(s1, s2 string) bool {
			return AreEqual(strings.ToLower(s1), strings.ToLower(s2), c.IgnoreFields, c.IgnoreChars)
		}
	} else {
		us.AreEqual = func(s1, s2 string) bool {
			return AreEqual(s1, s2, c.IgnoreFields, c.IgnoreChars)
		}
	}

	switch {
	case c.OnlyUnique:
		us.CountIsOk = func(count int) bool {
			return count == 1
		}
	case c.OnlyRepeated:
		us.CountIsOk = func(count int) bool {
			return count > 1
		}
	default:
		us.CountIsOk = func(count int) bool {
			return count > 0
		}
	}

	if c.CountEntries {
		us.Format = func(count int, line string) string {
			return fmt.Sprintf("%d %s", count, line)
		}
	} else {
		us.Format = func(count int, line string) string {
			return fmt.Sprint(line)
		}
	}

	if c.OutputPath != "" {
		us.Writer, _ = os.Create(c.OutputPath)
	} else {
		us.Writer = os.Stdout
	}

	return us
}

func main() {
	config := Config{}
	config.ParseFromFlags()
	if err := config.Validate(); err != nil {
		common.Exit1OnError(err)
	}
	us := NewUniqStrategy(config)
	us.Execute()
}
