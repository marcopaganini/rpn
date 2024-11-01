// This file is part of rpn, a simple and useful CLI RPN calculator.
// For further information, check https://github.com/marcopaganini/rpn
//
// (C) Sep/2024 by Marco Paganini <paganini AT paganini DOT net>

package main

import (
	"os"
	"testing"

	"github.com/ericlagergren/decimal"
)

func TestRPN(t *testing.T) {
	ctx := decimal.Context128

	casetests := []struct {
		input     string
		want      *decimal.Big
		wantError bool
	}{
		// Note: We use a "continuous" stack across operations.
		// Basic operations.
		{input: "1 2 +", want: bigUint(3)},
		{input: "4 5 8 + + +", want: bigUint(20)},
		{input: "5 -", want: bigUint(15)},
		{input: "3 /", want: bigUint(5)},
		{input: "3.5 *", want: bigFloat("17.5")},
		{input: "2 *", want: bigFloat("35.0")},
		{input: "chs", want: bigFloat("-35.0")},
		{input: "inv", want: bigFloat("-0.02857142857142857142857142857142857")},
		{input: "inv", want: bigFloat("-35.00000000000000000000000000000000")},
		{input: "chs", want: bigUint(35)},
		{input: "3 ^", want: bigUint(42875)},
		{input: "10 mod", want: bigUint(5)},
		{input: "fac", want: bigUint(120)},
		{input: "10 %", want: bigUint(12)},
		{input: "2 ^", want: bigUint(144)},
		{input: "sqr", want: bigUint(12)},
		{input: "3 ^", want: bigUint(1728)},
		{input: "cbr", want: bigUint(12)},
		{input: "c 1 2 3 4 sum", want: bigUint(10)},
		{input: "c 1 2 x", want: bigUint(1)},
		{input: "x", want: bigUint(2)},
		{input: "c", want: bigUint(0)},
		{input: "0.25 4 *", want: bigUint(1)},
		{input: ".2 5 *", want: bigUint(1)},
		{input: "$2,500.00 €3,500.00 +", want: bigUint(6000)},
		{input: "c", want: bigUint(0)},

		// Invalid operator should not cause changes to stack.
		{input: "foobar", want: bigUint(0)},

		// Trigonometric functions.
		{input: "deg 90 sin", want: bigFloat("0.9999999999999999999999999999999989")},
		{input: "rad 90 PI * 180 / sin", want: bigUint(1)},
		{input: "deg 0 cos", want: bigUint(1)},
		{input: "rad 60 PI * 180 / cos", want: bigFloat("0.5")},
		{input: "deg 30 tan", want: func() *decimal.Big { // Sqrt(3) / 3)
			square := big().Quo(bigUint(1), bigUint(2))
			z := ctx.Pow(big(), bigUint(3), square)
			return z.Quo(z, bigUint(3))
		}()},
		{input: "rad 60 PI * 180 / tan", want: ctx.Pow(big(), bigUint(3), big().Quo(bigUint(1), bigUint(2)))}, // Sqrt(3)
		{input: "deg -1 180 * PI / asin", want: big().Quo(ctx.Pi(big()), bigUint(2)).SetSignbit(true)},        // -Pi / 2
		{input: "rad -1 asin", want: big().Quo(ctx.Pi(big()), bigUint(2)).SetSignbit(true)},                   // -Pi / 2
		{input: "rad 0 acos", want: big().Quo(ctx.Pi(big()), bigUint(2))},                                     // Pi / 2
		{input: "rad 1 acos", want: bigUint(0)},
		{input: "deg 0.5 180 * PI / acos", want: big().Quo(ctx.Pi(big()), bigUint(3))}, // Pi / 3
		{input: "rad 3 sqr atan", want: big().Quo(ctx.Pi(big()), bigUint(3))},          // Pi / 3
		{input: "deg 1 180 * PI / atan", want: big().Quo(ctx.Pi(big()), bigUint(4))},   // Pi / 4

		// Bitwise operations and base input modes.
		{input: "0x00ff 0xff00 or", want: bigUint(0xffff)},
		{input: "0x0ff0 and", want: bigUint(0x0ff0)},
		{input: "0x1ee1 xor", want: bigUint(0x1111)},
		{input: "1 lshift", want: bigUint(0x2222)},
		{input: "1 lshift", want: bigUint(0x4444)},
		{input: "2 rshift", want: bigUint(0x1111)},
		{input: "0b00100010 0B01000100 015 o20 0x1000 0x2000 + + + + +", want: bigUint(12419)},
		{input: "c", want: bigUint(0)},

		// Bitwise operations and base input modes.
		{input: "0x00ff 0xff00 or", want: bigUint(0xffff)},
		{input: "0x0ff0 and", want: bigUint(0x0ff0)},
		{input: "0x1ee1 xor", want: bigUint(0x1111)},
		{input: "1 lshift", want: bigUint(0x2222)},
		{input: "1 lshift", want: bigUint(0x4444)},
		{input: "2 rshift", want: bigUint(0x1111)},
		{input: "0b00100010 0B01000100 015 o20 0x1000 0x2000 + + + + +", want: bigUint(12419)},
		{input: "d d", want: bigUint(0)},

		// Miscellaneous operations
		{input: "212 f2c", want: bigUint(100)},
		{input: "c2f", want: bigUint(212)},
		{input: "-40 f2c", want: bigFloat("-40")},
		{input: "-10 c2f", want: bigUint(14)},
		{input: "0 c2f", want: bigUint(32)},
		{input: "c", want: bigUint(0)},
	}

	stack := &stackType{}

	for _, tt := range casetests {
		err := calc(stack, tt.input)
		if !tt.wantError {
			if err != nil {
				t.Fatalf("Got error %q, want no error", err)
			}
			// Round to 6 decimals to make comparing floating point a bit easier.
			got := decimal.WithPrecision(16).Set(stack.top())
			want := decimal.WithPrecision(16).Set(tt.want)

			if got.Cmp(want) != 0 {
				t.Fatalf("diff: input: %s, want: %s, got: %s", tt.input, want, got)
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
	ctx := decimal.Context128

	casetests := []struct {
		base  int
		input *decimal.Big
		want  string
	}{
		// Decimal
		{10, bigUint(0), "0"},
		{10, bigUint(1), "1"},
		{10, bigUint(999), "999"},
		{10, bigUint(1000), "1000 (1,000)"},
		{10, bigUint(1000000), "1000000 (1,000,000)"},
		{10, bigUint(1000000000000000), "1000000000000000 (1,000,000,000,000,000)"},
		{10, bigFloat("10000.3333333333"), "10000.3333333333 (10,000.3333333333)"},
		{10, bigFloat("-10000.3333333333"), "-10000.3333333333 (-10,000.3333333333)"},
		{10, ctx.Pow(big(), bigUint(2), bigUint(64)), "18446744073709551616 (18,446,744,073,709,551,616)"},
		{10, ctx.Pow(big(), bigUint(2), bigUint(1234567890)), "Infinity"},
		{10, ctx.Quo(big(), bigUint(0), bigUint(0)), "NaN"},
		{10, ctx.Quo(big(), bigUint(1), bigUint(0)), "Infinity"},
		{10, ctx.Quo(big(), bigFloat("-1"), bigUint(0)), "-Infinity"},

		// Binary
		{2, bigUint(0b11111111), "0b11111111"},
		{2, big().Add(bigUint(0b11111111), bigFloat("0.5")), "0b11111111 (truncated from 255.5)"},
		{2, big().Add(bigUint(0b11111111), bigFloat("0.5")).SetSignbit(true), "-0b11111111 (truncated from -255.5)"},

		// Octal
		{8, bigUint(0377), "0377"},
		{8, big().Add(bigUint(0377), bigFloat("0.5")), "0377 (truncated from 255.5)"},
		{8, big().Add(bigUint(0377), bigFloat("0.5")).SetSignbit(true), "-0377 (truncated from -255.5)"},

		// Hex
		{16, bigUint(0xff), "0xff"},
		{16, big().Add(bigUint(0xff), bigFloat("0.5")).SetSignbit(true), "-0xff (truncated from -255.5)"},
	}
	for _, tt := range casetests {
		got := formatNumber(ctx, tt.input, tt.base)
		if got != tt.want {
			t.Fatalf("diff: base: %d, input: %v, want: %q, got: %q", tt.base, tt.input, tt.want, got)
		}
		continue
	}
}

func Example_main() {
	os.Args = []string{"rpn", "1", "2", "3", "+", "+", "6", "-"}
	main()
	// Output: 0
}
