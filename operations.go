// This file is part of rpn, a simple and useful CLI RPN calculator.
// For further information, check https://github.com/marcopaganini/rpn
//
// (C) Sep/2024 by Marco Paganini <paganini AT paganini DOT net>
package main

import (
	"errors"
	"fmt"
	"math"

	"github.com/fatih/color"
)

type (
	// ophandler contains the handler for a single operation.  numArgs
	// indicates how many arguments the function needs in the stack.
	ophandler struct {
		op      string // operator or command
		desc    string // operation description (used by help)
		numArgs int    // Number of arguments to function

		// Function receives the entire inverted stack (x=0, y=1, etc) and
		// returns the number of elements to be popped from the stack, and a
		// list of elements to be added pushes to the stack.
		fn func([]float64) ([]float64, int, error)
	}

	// opsType contains the base information for a list of operations and
	// their descriptions. The operations go in a list of interfaces so
	// we can also use strings and print them in the help() function.
	opsType struct {
		base  int           // Base for printing (default = 10)
		debug bool          // Debug state
		stack *stackType    // stack object to use
		ops   []interface{} // list of ophandlers & descriptions
	}

	// opmapType is a handler to operation map, used to find the right
	// operation function to call.
	opmapType map[string]ophandler
)

func newOpsType(stack *stackType) *opsType {
	bold := color.New(color.Bold).SprintFunc()

	ret := &opsType{
		base:  10,
		stack: stack,
	}
	ret.ops = []interface{}{
		// Header
		bold("Online help for ", programTitle, "."),
		bold("See http://github.com/marcopaganini/rpn for full details."),
		"",
		bold("Data entry:"),
		"  number <ENTER> - push a number on top of the stack.",
		"  operation <ENTER> - perform an operation on the stack (see below).",
		"",
		"  It's also possible to separate multiple operations with space:",
		"    10 2 3 * - (result = 4)",
		"",
		"  Prefix numbers with 0x to indicate hexadecimal, 0 for octal.",
		"",
		bold("Operations:"),
		"",
		bold("Basic Operations"),
		ophandler{"+", "Add x to y", 2, func(a []float64) ([]float64, int, error) {
			return []float64{a[1] + a[0]}, 2, nil
		}},
		ophandler{"-", "Subtract x from y", 2, func(a []float64) ([]float64, int, error) {
			return []float64{a[1] - a[0]}, 2, nil
		}},
		ophandler{"*", "Multiply x and y", 2, func(a []float64) ([]float64, int, error) {
			return []float64{a[1] * a[0]}, 2, nil
		}},
		ophandler{"/", "Divide y by x", 2, func(a []float64) ([]float64, int, error) {
			if a[0] == 0 {
				return nil, 2, errors.New("can't divide by zero")
			}
			return []float64{a[1] / a[0]}, 2, nil
		}},
		ophandler{"chs", "Change signal of x", 1, func(a []float64) ([]float64, int, error) {
			return []float64{a[0] * -1}, 1, nil
		}},
		ophandler{"inv", "Invert x (1/x)", 1, func(a []float64) ([]float64, int, error) {
			return []float64{1 / a[0]}, 1, nil
		}},
		ophandler{"^", "Raise y to the power of x", 2, func(a []float64) ([]float64, int, error) {
			return []float64{math.Pow(a[1], a[0])}, 2, nil
		}},
		ophandler{"mod", "Calculates y modulo x", 2, func(a []float64) ([]float64, int, error) {
			return []float64{math.Mod(a[1], a[0])}, 2, nil
		}},
		ophandler{"%", "Calculate x% of y", 2, func(a []float64) ([]float64, int, error) {
			return []float64{a[1] * a[0] / 100}, 2, nil
		}},

		ophandler{"fac", "Calculate factorial of x", 1, func(a []float64) ([]float64, int, error) {
			x := uint64(a[0])
			if x <= 0 {
				return nil, 1, errors.New("factorial requires a positive number")
			}
			fact := uint64(1)
			for ix := uint64(1); ix <= x; ix++ {
				fact *= ix
			}
			return []float64{float64(fact)}, 1, nil
		}},
		"",
		bold("Bitwise Operations"),
		ophandler{"and", "Logical AND between x and y", 2, func(a []float64) ([]float64, int, error) {
			return []float64{float64(uint64(a[1]) & uint64(a[0]))}, 2, nil
		}},
		ophandler{"or", "Logical OR between x and y", 2, func(a []float64) ([]float64, int, error) {
			return []float64{float64(uint64(a[1]) | uint64(a[0]))}, 2, nil
		}},
		ophandler{"xor", "Logical XOR between x and y", 2, func(a []float64) ([]float64, int, error) {
			return []float64{float64(uint64(a[1]) ^ uint64(a[0]))}, 2, nil
		}},
		ophandler{"lshift", "Shift y left x times", 2, func(a []float64) ([]float64, int, error) {
			return []float64{float64(uint64(a[1]) << uint64(a[0]))}, 2, nil
		}},
		ophandler{"rshift", "Shift y right x times", 2, func(a []float64) ([]float64, int, error) {
			return []float64{float64(uint64(a[1]) >> uint64(a[0]))}, 2, nil
		}},

		"",
		bold("Stack Operations"),
		ophandler{"p", "Display stack", 0, func(_ []float64) ([]float64, int, error) {
			stack.print(ret.base)
			return nil, 0, nil
		}},
		ophandler{"c", "Clear stack", 0, func(_ []float64) ([]float64, int, error) {
			stack.clear()
			return nil, 0, nil
		}},
		ophandler{"=", "Print top of stack (x)", 0, func(_ []float64) ([]float64, int, error) {
			fmt.Println(stack.top())
			return nil, 0, nil
		}},
		ophandler{"d", "Drop top of stack (x)", 1, func(_ []float64) ([]float64, int, error) {
			return nil, 1, nil
		}},

		"",
		bold("Math and Physical constants"),
		ophandler{"PI", "The famous transcedental number", 0, func(_ []float64) ([]float64, int, error) {
			return []float64{math.Pi}, 0, nil
		}},
		ophandler{"E", "Another famous transcedental number", 0, func(_ []float64) ([]float64, int, error) {
			return []float64{math.E}, 0, nil
		}},
		ophandler{"C", "Speed of light in vacuum, in m/s", 0, func(_ []float64) ([]float64, int, error) {
			return []float64{299792458}, 0, nil
		}},
		ophandler{"MOL", "Avogadro's number", 1, func(_ []float64) ([]float64, int, error) {
			return []float64{6.02214154e23}, 0, nil
		}},

		"",
		bold("Computer constants"),
		ophandler{"KB", "Kilobyte", 0, func(_ []float64) ([]float64, int, error) {
			return []float64{math.Pow(10, 3)}, 0, nil
		}},
		ophandler{"MB", "Megabyte", 0, func(_ []float64) ([]float64, int, error) {
			return []float64{math.Pow(10, 6)}, 0, nil
		}},
		ophandler{"GB", "Gigabyte", 0, func(_ []float64) ([]float64, int, error) {
			return []float64{math.Pow(10, 9)}, 0, nil
		}},
		ophandler{"MB", "Terabyte", 0, func(_ []float64) ([]float64, int, error) {
			return []float64{math.Pow(10, 12)}, 0, nil
		}},
		ophandler{"KIB", "Kibibyte", 0, func(_ []float64) ([]float64, int, error) {
			return []float64{math.Pow(2, 10)}, 0, nil
		}},
		ophandler{"MIB", "Mebibyte", 0, func(_ []float64) ([]float64, int, error) {
			return []float64{math.Pow(2, 20)}, 0, nil
		}},
		ophandler{"GIB", "Gibibyte", 0, func(_ []float64) ([]float64, int, error) {
			return []float64{math.Pow(2, 30)}, 0, nil
		}},
		ophandler{"TIB", "Tebibyte", 0, func(_ []float64) ([]float64, int, error) {
			return []float64{math.Pow(2, 40)}, 0, nil
		}},

		"",
		bold("Program Control"),
		ophandler{"dec", "Output in decimal", 0, func(_ []float64) ([]float64, int, error) {
			ret.base = 10
			return nil, 0, nil
		}},
		ophandler{"bin", "Output in binary", 0, func(_ []float64) ([]float64, int, error) {
			ret.base = 2
			return nil, 0, nil
		}},
		ophandler{"oct", "Output in octal", 0, func(_ []float64) ([]float64, int, error) {
			ret.base = 8
			return nil, 0, nil
		}},
		ophandler{"hex", "Output in hexadecimal", 0, func(_ []float64) ([]float64, int, error) {
			ret.base = 16
			return nil, 0, nil
		}},
		ophandler{"debug", "Toggle debugging", 0, func(_ []float64) ([]float64, int, error) {
			ret.debug = !ret.debug
			fmt.Printf("Debugging state: %v\n", ret.debug)
			return nil, 0, nil
		}},
		"",
		bold("Please Note:"),
		"  - x means the number at the top of the stack",
		"  - y means the second number from the top of the stack",
	}
	return ret
}

// help displays the help message to the screen based on the contents of opmap.
func (x opsType) help() {
	bold := color.New(color.Bold).SprintFunc()
	for _, v := range x.ops {
		if handler, ok := v.(ophandler); ok {
			fmt.Printf("  - %s: %s\n", bold(handler.op), handler.desc)
			continue
		}
		fmt.Println(v)
	}
}

// opmap returns a map of op (command) -> ophandler that can be easily used
// later to find the function to be executed. It takes a slice of interfaces
// and returns a map[string][ophandler].
func (x opsType) opmap() opmapType {
	ret := map[string]ophandler{}

	for _, v := range x.ops {
		if h, ok := v.(ophandler); ok {
			ret[h.op] = h
		}
	}
	return ret
}
