package uniq

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
	Input             io.Reader
	AreEqual          func(string, string) bool
	CountIsAcceptable func(int) bool
	Format            func(int, string) string
	Output            io.Writer
}

func (us *UniqStrategy) Execute() {
	prev, latest := "", ""
	count := 0
	scanner := bufio.NewScanner(us.Input)

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
		if us.CountIsAcceptable(count) {
			fmt.Fprintln(us.Output, us.Format(count, prev))
		}
		prev = latest
		count = 1
	}
	if us.CountIsAcceptable(count) {
		fmt.Fprintln(us.Output, us.Format(count, prev))
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func drop[T any](arr []T, n int) []T {
	return arr[min(len(arr), n):]
}

func areEqual(s1, s2 string, ignoreFields, ignoreChars int) bool {
	first := drop(strings.Fields(s1), ignoreFields)
	second := drop(strings.Fields(s2), ignoreFields)

	dropRunes := func(words []string) string {
		joined := strings.Join(words, " ")
		dropped := drop([]rune(joined), ignoreChars)
		return string(dropped)
	}
	return strings.Compare(dropRunes(first), dropRunes(second)) == 0
}

func NewUniqStrategy(c Config) *UniqStrategy {
	us := &UniqStrategy{}

	if c.InputPath != "" {
		us.Input, _ = os.Open(c.InputPath)
	} else {
		us.Input = os.Stdin
	}

	if c.IgnoreCase {
		us.AreEqual = func(s1, s2 string) bool {
			return areEqual(strings.ToLower(s1), strings.ToLower(s2), c.IgnoreFields, c.IgnoreChars)
		}
	} else {
		us.AreEqual = func(s1, s2 string) bool {
			return areEqual(s1, s2, c.IgnoreFields, c.IgnoreChars)
		}
	}

	switch {
	case c.OnlyUnique:
		us.CountIsAcceptable = func(count int) bool {
			return count == 1
		}
	case c.OnlyRepeated:
		us.CountIsAcceptable = func(count int) bool {
			return count > 1
		}
	default:
		us.CountIsAcceptable = func(count int) bool {
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
		us.Output, _ = os.Create(c.OutputPath)
	} else {
		us.Output = os.Stdout
	}

	return us
}
