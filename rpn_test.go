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
				t.Fatalf("diff: input %q, want %.8f, got %.8f.", tt.input, tt.want, got)
			}
			continue
		}

		// Here, we want to see an error.
		if err == nil {
			t.Errorf("Got no error, want error")
		}
	}
}
