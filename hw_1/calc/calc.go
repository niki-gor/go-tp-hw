package calc

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrUnexpectedClosingBrace = errors.New("unexpected closing brace")
	ErrUnclosedBrace          = errors.New("unclosed brace")
	ErrZeroDivision           = errors.New("division by zero")
	ErrMissingOperator        = errors.New("missing operator")

	spaceBetweenNumber = regexp.MustCompile(`\d\s+\d`)
	spaces             = regexp.MustCompile(`\s+`)
	plusesMinuses      = regexp.MustCompile(`[\+\-]+`)
	binaryPlusMinus    = regexp.MustCompile(`(\d)[+-]`) // после числа идет +-, следовательно бинарный +-
	multDiv            = regexp.MustCompile(`[\*\/]`)
)

func reducePlusesMinuses(s string) string {
	if strings.Count(s, "-")%2 == 0 {
		return "+"
	}
	return "-"
}

func calcMonomial(s string) (result int, err error) {
	literals := multDiv.Split(s, -1)
	operators := multDiv.FindAllString(s, -1)
	result, err = strconv.Atoi(literals[0])
	if err != nil {
		return
	}
	for i := range operators {
		var literal int
		literal, err = strconv.Atoi(literals[i+1])
		if err != nil {
			return
		}
		if operators[i] == "*" {
			result *= literal
		} else {
			if literal == 0 {
				err = ErrZeroDivision
				return
			}
			result /= literal
		}
	}
	return
}

// вычисляет значение выражения без скобок
func PlainCalc(s string) (string, error) {
	s = plusesMinuses.ReplaceAllStringFunc(s, reducePlusesMinuses) // сокращение записей +---+--х до -х и т.п.
	operators := binaryPlusMinus.FindAllString(s, -1)
	s = binaryPlusMinus.ReplaceAllString(s, "$1 ")
	monomials := strings.Split(s, " ")
	result, err := calcMonomial(monomials[0])
	if err != nil {
		return "", err
	}
	for i := range operators {
		monomial, err := calcMonomial(monomials[i+1])
		if err != nil {
			return "", err
		}
		if operators[i][1] == '+' {
			result += monomial
		} else {
			result -= monomial
		}
	}
	return strconv.Itoa(result), nil
}

func calc(s string) (result string, length int, err error) {
	plain := ""
	i := 0
	for i < len(s) {
		c := s[i]
		i += 1
		if c == '(' {
			var inner string
			var addLength int
			inner, addLength, err = calc(s[i:])
			if err != nil {
				return
			}
			i += addLength
			plain += inner
		} else if c == ')' {
			result, err = PlainCalc(plain)
			length = i
			return
		} else {
			plain += string(c)
		}
	}
	err = ErrUnclosedBrace
	return
}

func Calc(s string) (string, error) {
	if spaceBetweenNumber.MatchString(s) { // литералы вида "100 500", "19 84" не допускаются
		return "", ErrMissingOperator // считается, что между ними должен стоть оператор
	}
	s = spaces.ReplaceAllString(s, "")
	s += ")" // требуется для удобного парсинга
	result, length, err := calc(s)
	if err != nil { // вычисляет значение без скобок
		return "", err
	}
	if length != len(s) { // проверка, что последняя сбалансированная закрывающая скобка находится в конце строки
		return "", ErrUnexpectedClosingBrace
	}
	return result, nil
}
