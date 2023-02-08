package main

import (
	"flag"
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

func (c *Config) Validate() {
	b2i := func(b bool) int {
		if b {
			return 1
		}
		return 0
	}
	if b2i(c.CountEntries)+b2i(c.OnlyRepeated)+b2i(c.OnlyUnique) > 1 {

	}
}

func main() {
	config := Config{}
	config.ParseFromFlags()
	config.Validate()

}
