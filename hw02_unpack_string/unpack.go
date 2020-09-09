// A function that performs primitive unpacking of repeated characters/runes,
// for example:
//		"a4bc2d5e" => "aaaabccddddde"

package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")
var ErrWrite = errors.New("write error")

type state uint8

const (
	start state = iota
	print
	backSlash
	digit
	exit
)

func Unpack(str string) (string, error) {
	var resultStr strings.Builder
	state := start
	var prevCh rune
	for _, ch := range str {
		switch state {
		case start:
			startState(&state, ch)
		case print:
			if err := printState(&state, ch, prevCh, &resultStr); err != nil {
				return "", err
			}
		case backSlash:
			backSlashState(&state, ch)
		case digit:
			digitState(&state, ch)
		case exit:
			return "", ErrInvalidString
		}
		prevCh = ch
	}

	if state != digit && str != "" {
		if _, err := resultStr.WriteRune(prevCh); err != nil {
			return "", ErrWrite
		}
	}

	return resultStr.String(), nil
}

func startState(state *state, ch rune) {
	switch {
	case unicode.IsDigit(ch):
		*state = exit
	case ch == '\\':
		*state = backSlash
	default:
		*state = print
	}
}

func printState(state *state, ch rune, prevCh rune, resultStr *strings.Builder) error {
	switch {
	case unicode.IsDigit(ch):
		num := int(ch - '0')
		if _, err := resultStr.WriteString(strings.Repeat(string(prevCh), num)); err != nil {
			return ErrWrite
		}
		*state = digit
	case ch == '\\':
		if _, err := resultStr.WriteRune(prevCh); err != nil {
			return ErrWrite
		}
		*state = backSlash
	default:
		if _, err := resultStr.WriteRune(prevCh); err != nil {
			return ErrWrite
		}
		*state = print
	}
	return nil
}

func backSlashState(state *state, ch rune) {
	if unicode.IsDigit(ch) || ch == '\\' {
		*state = print
	} else {
		*state = exit
	}
}

func digitState(state *state, ch rune) {
	switch {
	case unicode.IsDigit(ch):
		*state = exit
	case ch == '\\':
		*state = backSlash
	default:
		*state = print
	}
}
