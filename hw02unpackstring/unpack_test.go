package hw02unpackstring

import "testing"

type test struct {
	input    string
	expected string
	err      error
}

func TestUnpack(t *testing.T) {
	var tests = []test{
		{
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			input:    "abccd",
			expected: "abccd",
		},
		{
			input:    "3abc",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "45",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "aaa10b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "aaa0b",
			expected: "aab",
		},
		{
			input:    "d\n5abc",
			expected: "d\n\n\n\n\nabc",
		},
		{
			input:    "При3вет!",
			expected: "Прииивет!",
		},
	}

	for _, tst := range tests {
		result, err := Unpack(tst.input)
		if tst.err != err || tst.expected != result {
			t.Errorf("Unpack(%q) = %q, %v; want: %q, %v", tst.input, result, err, tst.expected, tst.err)
		}
	}
}

func TestUnpackWithEscape(t *testing.T) {
	var tests = []test{
		{
			input:    `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			input:    `qwe\45`,
			expected: `qwe44444`,
		},
		{
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			input:    `qwe\\\3`,
			expected: `qwe\3`,
		},
		{
			input:    `qwe\\5a`,
			expected: `qwe\\\\\a`,
		},
		{
			input:    `qw\ne`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{
			input:    `\\`,
			expected: `\`,
		},
		{
			input:    `\n`,
			expected: ``,
			err:      ErrInvalidString,
		},
		{
			input:    `\`,
			expected: ``,
			err:      ErrInvalidString,
		},
	}

	for _, tst := range tests {
		result, err := Unpack(tst.input)
		if tst.err != err || tst.expected != result {
			t.Errorf("Unpack(%q) = %q, %v; want: %q, %v", tst.input, result, err, tst.expected, tst.err)
		}
	}
}
