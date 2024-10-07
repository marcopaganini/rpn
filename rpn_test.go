// This file is part of rpn, a simple and useful CLI RPN calculator.
// For further information, check https://github.com/marcopaganini/rpn
//
// (C) Sep/2024 by Marco Paganini <paganini AT paganini DOT net>

package main

import (
	"testing"
)

func TestRPN(t *testing.T) {
	casetests := []struct {
		input     string
		want      float64
		wantError bool
	}{
		// Note: We use a "continuous" stack across operations.
		// Basic operations.
		{input: "1 2 +", want: 3},
		{input: "4 5 8 + + +", want: 20},
		{input: "5 -", want: 15},
		{input: "3 /", want: 5},
		{input: "3.5 *", want: 17.5},
		{input: "2 *", want: 35},
		{input: "chs", want: -35},
		{input: "inv", want: -0.02857142857142857},
		{input: "inv", want: -35},
		{input: "chs", want: 35},
		{input: "3 ^", want: 42875},
		{input: "10 mod", want: 5},
		{input: "fac", want: 120},
		{input: "10 %", want: 12},
		{input: "2 ^", want: 144},
		{input: "sqr", want: 12},
		{input: "3 ^", want: 1728},
		{input: "cbr", want: 12},
		{input: "c 1 2 3 4 sum", want: 10},
		{input: "c 1 2 x", want: 1},
		{input: "x", want: 2},
		{input: "c", want: 0},
		{input: "0.25 4 *", want: 1},
		{input: ".2 5 *", want: 1},
		{input: "c", want: 0},

		// Bitwise operations and base input modes.
		{input: "0x00ff 0xff00 or", want: 0xffff},
		{input: "0x0ff0 and", want: 0x0ff0},
		{input: "0x1ee1 xor", want: 0x1111},
		{input: "1 lshift", want: 0x2222},
		{input: "1 lshift", want: 0x4444},
		{input: "2 rshift", want: 0x1111},
		{input: "0b00100010 0B01000100 015 o20 0x1000 0x2000 + + + + +", want: 12419},
		{input: "c", want: 0},

		// Bitwise operations and base input modes.
		{input: "0x00ff 0xff00 or", want: 0xffff},
		{input: "0x0ff0 and", want: 0x0ff0},
		{input: "0x1ee1 xor", want: 0x1111},
		{input: "1 lshift", want: 0x2222},
		{input: "1 lshift", want: 0x4444},
		{input: "2 rshift", want: 0x1111},
		{input: "0b00100010 0B01000100 015 o20 0x1000 0x2000 + + + + +", want: 12419},
		{input: "d d", want: 0},

		// Miscellaneous operations
		{input: "212 f2c", want: 100},
		{input: "c2f", want: 212},
		{input: "-40 f2c", want: -40},
		{input: "-10 c2f", want: 14},
		{input: "0 c2f", want: 32},
		{input: "c", want: 0},
	}

	stack := &stackType{}

	for _, tt := range casetests {
		err := calc(stack, tt.input)
		if !tt.wantError {
			if err != nil {
				t.Fatalf("Got error %q, want no error", err)
			}
			got := stack.top()
			if got != tt.want {
				t.Fatalf("diff: input: %q, want: %.8f, got: %.8f", tt.input, tt.want, got)
			}
			continue
		}

		// Here, we want to see an error.
		if err == nil {
			t.Errorf("Got no error, want error")
		}
	}
}

func TestFormatNumber(t *testing.T) {
	casetests := []struct {
		base  int
		input float64
		want  string
	}{
		// Decimal
		{10, 0, "0"},
		{10, 1, "1"},
		{10, 999, "999"},
		{10, 1000, "1000 (1,000)"},
		{10, 1000000, "1000000 (1,000,000)"},
		{10, 1000000000000000, "1000000000000000 (1,000,000,000,000,000)"},
		{10, 10000.3333333333, "10000.333333 (10,000.333333)"},
		{10, -10000.3333333333, "-10000.333333 (-10,000.333333)"},

		// Binary
		{2, 0b11111111, "0b11111111"},
		{2, 0b11111111 + 0.5, "0b11111111 (truncated from 255.5)"},
		{2, (0b11111111 + 0.5) * -1, "0b1111111111111111111111111111111111111111111111111111111100000001 (truncated from -255.5)"},

		// Octal
		{8, 0377, "0377"},
		{8, 0377 + 0.5, "0377 (truncated from 255.5)"},
		{8, (0377 + 0.5) * -1, "01777777777777777777401 (truncated from -255.5)"},

		// Hex
		{16, 0xff, "0xff"},
		{16, (0xff + 0.5) * -1, "0xffffffffffffff01 (truncated from -255.5)"},
	}
	for _, tt := range casetests {
		got := formatNumber(tt.input, tt.base)
		if got != tt.want {
			t.Fatalf("diff: base: %d, input: %v, want: %q, got: %q", tt.base, tt.input, tt.want, got)
		}
		continue
	}
}
