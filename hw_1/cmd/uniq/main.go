package main

import (
	"flag"
	"go_tp/hw_1/uniq"
	"log"
)

func ParseFromFlags(c *uniq.Options) {
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

func main() {
	config := uniq.Options{}
	ParseFromFlags(&config)
	if err := config.Validate(); err != nil {
		log.Fatal(err)
	}
	us := uniq.NewUniqStrategy(config)
	us.Execute()
}
