package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"go_tp/hw_1/common"
	"log"
	"os"
	"strings"
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
	log.Println(c.CountEntries, c.OnlyRepeated, c.OnlyUnique, c.IgnoreFields, c.IgnoreChars, c.IgnoreCase, c.InputPath, c.OutputPath)

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
	InputFile *os.File
	AreEqual  func(string, string) bool
	CountIsOk func(int) bool
	Println   func(int, string)
	Teardown  func()
}

func (us *UniqStrategy) Execute() {
	prev, latest := "", ""
	count := 0
	scanner := bufio.NewScanner(us.InputFile)
	defer us.InputFile.Close()

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
			us.Println(count, prev)
		}
		prev = latest
		count = 1
	}
	if us.CountIsOk(count) {
		us.Println(count, prev)
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

	joinDropAndCast := func(words []string) string {
		return string(Drop([]rune(strings.Join(words, " ")), ignoreChars))
	}
	return strings.Compare(joinDropAndCast(first), joinDropAndCast(second)) == 0
}

func NewUniqStrategy(c Config) *UniqStrategy {
	us := &UniqStrategy{}

	if c.InputPath != "" {
		us.InputFile, _ = os.Open(c.InputPath)
	} else {
		us.InputFile = os.Stdin
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

	var outputFile *os.File
	if c.OutputPath != "" {
		outputFile, _ = os.Create(c.OutputPath)
	} else {
		outputFile = os.Stdout
	}
	if c.CountEntries {
		us.Println = func(count int, line string) {
			fmt.Fprintf(outputFile, "%d %s\n", count, line)
		}
	} else {
		us.Println = func(count int, line string) {
			fmt.Fprintln(outputFile, line)
		}
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
