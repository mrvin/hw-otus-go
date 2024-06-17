// A function that performs primitive unpacking of repeated characters/runes,
// for example:
//		"a4bc2d5e" => "aaaabccddddde"

package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

type state uint8

const (
	start state = iota
	printSt
	backSlash
	exit
)

func Unpack(str string) (string, error) {
	var resultStr strings.Builder

	var prevCh rune

	state := start
	for _, ch := range str {
		switch state {
		case start:
			state = startState(ch)
		case printSt:
			state = printState(ch, prevCh, &resultStr)
		case backSlash:
			state = backSlashState(ch)
		case exit:
			return "", ErrInvalidString
		}
		prevCh = ch
	}

	if state == backSlash || state == exit {
		return "", ErrInvalidString
	}

	if state != start {
		resultStr.WriteRune(prevCh)
	}

	return resultStr.String(), nil
}

func startState(ch rune) state {
	switch {
	case unicode.IsDigit(ch):
		return exit
	case ch == '\\':
		return backSlash
	default:
		return printSt
	}
}

func printState(ch rune, prevCh rune, resultStr *strings.Builder) state {
	switch {
	case unicode.IsDigit(ch):
		num := int(ch - '0')
		resultStr.WriteString(strings.Repeat(string(prevCh), num))
		return start
	case ch == '\\':
		resultStr.WriteRune(prevCh)
		return backSlash
	default:
		resultStr.WriteRune(prevCh)
		return printSt
	}
}

func backSlashState(ch rune) state {
	if unicode.IsDigit(ch) || ch == '\\' {
		return printSt
	}
	return exit
}
