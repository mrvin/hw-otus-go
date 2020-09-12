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

type state uint8

const (
	start state = iota
	print
	backSlash
	exit
)

func Unpack(str string) (string, error) {
	var resultStr strings.Builder

	state := start
	var prevCh rune
	for _, ch := range str {
	splitStr := SplitAfterDigit(str)
			switch state { //nolint:exhaustive // Warns about missing cases in switch of type state: exit
			case start:
				var err error
				if state, err = startState(ch); err != nil {
					return "", err
				}
			case print:
				state = printState(ch, prevCh, &resultStr)
			case backSlash:
				var err error
				if state, err = backSlashState(ch); err != nil {
					return "", err
				}
			}
		}
		prevCh = ch
	}

		}
	}

	return resultStr.String(), nil
}

func startState(ch rune) (state, error) {
	switch {
	case unicode.IsDigit(ch):
		return exit, ErrInvalidString
	case ch == '\\':
		return backSlash, nil
	default:
		return print, nil
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
		return print
	}
}

func backSlashState(ch rune) (state, error) {
	if unicode.IsDigit(ch) || ch == '\\' {
		return print, nil
	}
	return exit, ErrInvalidString
}

// SplitAfterDigit splits the string into digit-terminated substrings, not counting the escaped digits.
func SplitAfterDigit(str string) []string {
	var result []string
	var substr strings.Builder
	var prevCh rune
	var countBS int

	for _, ch := range str {
		switch {
		case !unicode.IsDigit(ch) || (unicode.IsDigit(ch) && prevCh == '\\' && countBS == 1):
			substr.WriteRune(ch)
			if ch == '\\' {
				countBS++
			} else {
				countBS = 0
			}
		default:
			substr.WriteRune(ch)
			result = append(result, substr.String())
			substr.Reset()
			countBS = 0
		}
		prevCh = ch
	}

	if substr.Len() != 0 {
		result = append(result, substr.String())
	}

	return result
}
