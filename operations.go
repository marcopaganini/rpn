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
	// indicates how many items at the top of the stack should be popped and
	// sent to the function. ignoreResult will cause the function return to be
	// ignored (i.e, not be repushed into the stack and the top of the stack
	// won't be printed afterwards. This is used for commands like "show stack"
	// that don't change the stack at all.
	ophandler struct {
		op           string // operator or command
		desc         string // operation description (used by help)
		numArgs      int    // Number of arguments to function
		ignoreResult bool   // Ignore results from function
		fn           func([]float64) (float64, error)
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
		ophandler{"+", "Add x to y", 2, false, func(a []float64) (float64, error) {
			return a[0] + a[1], nil
		}},
		ophandler{"-", "Subtract x from y", 2, false, func(a []float64) (float64, error) {
			return a[0] - a[1], nil
		}},
		ophandler{"*", "Multiply x and y", 2, false, func(a []float64) (float64, error) {
			return a[0] * a[1], nil
		}},
		ophandler{"/", "Divide y by x", 2, false, func(a []float64) (float64, error) {
			if a[1] == 0 {
				return 0, errors.New("can't divide by zero")
			}
			return a[0] / a[1], nil
		}},
		ophandler{"chs", "Change signal of x", 1, false, func(a []float64) (float64, error) {
			return a[0] * -1, nil
		}},
		ophandler{"inv", "Invert x (1/x)", 1, false, func(a []float64) (float64, error) {
			return 1 / a[0], nil
		}},
		ophandler{"^", "Raise y to the power of x", 2, false, func(a []float64) (float64, error) {
			return math.Pow(a[0], a[1]), nil
		}},
		ophandler{"mod", "Calculates y modulo x", 2, false, func(a []float64) (float64, error) {
			return math.Mod(a[0], a[1]), nil
		}},
		ophandler{"%", "Calculate x% of y", 2, false, func(a []float64) (float64, error) {
			return a[0] * a[1] / 100, nil
		}},

		ophandler{"fac", "Calculate factorial of x", 1, false, func(a []float64) (float64, error) {
			x := uint64(a[0])
			if x <= 0 {
				return 0, errors.New("factorial requires a positive number")
			}
			fact := uint64(1)
			for ix := uint64(1); ix <= x; ix++ {
				fact *= ix
			}
			return float64(fact), nil
		}},
		"",
		bold("Bitwise Operations"),
		ophandler{"and", "Logical AND between x and y", 2, false, func(a []float64) (float64, error) {
			return float64(uint64(a[0]) & uint64(a[1])), nil
		}},
		ophandler{"or", "Logical OR between x and y", 2, false, func(a []float64) (float64, error) {
			return float64(uint64(a[0]) | uint64(a[1])), nil
		}},
		ophandler{"xor", "Logical XOR between x and y", 2, false, func(a []float64) (float64, error) {
			return float64(uint64(a[0]) ^ uint64(a[1])), nil
		}},
		ophandler{"lshift", "Shift y left x times", 2, false, func(a []float64) (float64, error) {
			return float64(uint64(a[0]) << uint64(a[1])), nil
		}},
		ophandler{"rshift", "Shift y right x times", 2, false, func(a []float64) (float64, error) {
			return float64(uint64(a[0]) >> uint64(a[1])), nil
		}},

		"",
		bold("Stack Operations"),
		ophandler{"p", "Display stack", 0, true, func(_ []float64) (float64, error) {
			stack.print(ret.base)
			return 0, nil
		}},
		ophandler{"c", "Clear stack", 0, true, func(_ []float64) (float64, error) {
			stack.clear()
			return 0, nil
		}},
		ophandler{"=", "Print top of stack (x)", 0, true, func(_ []float64) (float64, error) {
			fmt.Println(stack.top())
			return 0, nil
		}},
		ophandler{"d", "Drop top of stack (x)", 1, true, func(_ []float64) (float64, error) {
			return 0, nil
		}},

		"",
		bold("Math and Physical constants"),
		ophandler{"PI", "The famous transcedental number", 0, false, func(_ []float64) (float64, error) {
			return math.Pi, nil
		}},
		ophandler{"E", "Another famous transcedental number", 0, false, func(_ []float64) (float64, error) {
			return math.E, nil
		}},
		ophandler{"C", "Speed of light in vacuum, in m/s", 0, false, func(_ []float64) (float64, error) {
			return 299792458, nil
		}},
		ophandler{"MOL", "Avogadro's number", 1, false, func(_ []float64) (float64, error) {
			return 6.02214154e23, nil
		}},

		"",
		bold("Computer constants"),
		ophandler{"KB", "Kilobyte", 0, false, func(_ []float64) (float64, error) {
			return math.Pow(10, 3), nil
		}},
		ophandler{"MB", "Megabyte", 0, false, func(_ []float64) (float64, error) {
			return math.Pow(10, 6), nil
		}},
		ophandler{"GB", "Gigabyte", 0, false, func(_ []float64) (float64, error) {
			return math.Pow(10, 9), nil
		}},
		ophandler{"MB", "Terabyte", 0, false, func(_ []float64) (float64, error) {
			return math.Pow(10, 12), nil
		}},
		ophandler{"KIB", "Kibibyte", 0, false, func(_ []float64) (float64, error) {
			return math.Pow(2, 10), nil
		}},
		ophandler{"MIB", "Mebibyte", 0, false, func(_ []float64) (float64, error) {
			return math.Pow(2, 20), nil
		}},
		ophandler{"GIB", "Gibibyte", 0, false, func(_ []float64) (float64, error) {
			return math.Pow(2, 30), nil
		}},
		ophandler{"TIB", "Tebibyte", 0, false, func(_ []float64) (float64, error) {
			return math.Pow(2, 40), nil
		}},

		"",
		bold("Program Control"),
		ophandler{"dec", "Output in decimal", 0, true, func(_ []float64) (float64, error) {
			ret.base = 10
			return 0, nil
		}},
		ophandler{"hex", "Output in hexadecimal", 0, true, func(_ []float64) (float64, error) {
			ret.base = 16
			return 0, nil
		}},
		ophandler{"oct", "Output in octal", 0, true, func(_ []float64) (float64, error) {
			ret.base = 8
			return 0, nil
		}},
		ophandler{"debug", "Toggle debugging", 0, true, func(_ []float64) (float64, error) {
			ret.debug = !ret.debug
			fmt.Printf("Debugging state: %v\n", ret.debug)
			return 0, nil
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
