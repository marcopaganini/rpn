// This file is part of rpn, a simple and useful CLI RPN calculator.
// For further information, check https://github.com/marcopaganini/rpn
//
// (C) Sep/2024 by Marco Paganini <paganini AT paganini DOT net>
package main

import (
	"errors"
	"fmt"
	"math"
	"strings"

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
		base    int           // Base for printing (default = 10)
		debug   bool          // Debug state
		degmode bool          // Degrees mode (default = Radians)
		stack   *stackType    // stack object to use
		ops     []interface{} // list of ophandlers & descriptions
	}

	// opmapType is a handler to operation map, used to find the right
	// operation function to call.
	opmapType map[string]ophandler
)

// radOrDeg converts the value passed to radians if degmode (degrees
// mode) is set. Otherwise, it just returns the same value (radians).
func radOrDeg(n float64, degmode bool) float64 {
	if degmode {
		return n * math.Pi / 180
	}
	return n
}

func newOpsType(stack *stackType) *opsType {
	ret := &opsType{
		base:  10,
		stack: stack,
	}
	var build string
	if Build == "" {
		build = "no version info"
	} else {
		build = "v" + Build
	}

	ret.ops = []interface{}{
		// Header
		"BOLD:Online help for " + programTitle + " (" + build + ").",
		"BOLD:See http://github.com/marcopaganini/rpn for full details.",
		"",
		"BOLD:Data entry:",
		"  number <ENTER> - push a number on top of the stack.",
		"  operation <ENTER> - perform an operation on the stack (see below).",
		"",
		"  It's also possible to separate multiple operations with space:",
		"    10 2 3 * - (result = 4)",
		"",
		"  Prefix numbers with 0x to indicate hexadecimal, 0 for octal.",
		"",
		"BOLD:Operations:",
		"",
		"BOLD:Basic Operations",
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
		ophandler{"sqr", "Calculate square root of x", 1, func(a []float64) ([]float64, int, error) {
			return []float64{math.Sqrt(a[0])}, 1, nil
		}},
		ophandler{"cbr", "Calculate cubic root of x", 1, func(a []float64) ([]float64, int, error) {
			return []float64{math.Cbrt(a[0])}, 1, nil
		}},
		ophandler{"%", "Calculate x% of y", 2, func(a []float64) ([]float64, int, error) {
			return []float64{a[1] * a[0] / 100}, 2, nil
		}},

		ophandler{"sum", "Sum all elements in stack", 1, func(a []float64) ([]float64, int, error) {
			var sum float64
			for _, v := range a {
				sum += v
			}
			return []float64{sum}, len(a), nil
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
		"BOLD:Bitwise Operations",
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
		"BOLD:Trigonometric Operations",
		ophandler{"sin", "Sine of x", 1, func(a []float64) ([]float64, int, error) {
			return []float64{math.Sin(radOrDeg(a[0], ret.degmode))}, 1, nil
		}},
		ophandler{"cos", "Cosine of x", 1, func(a []float64) ([]float64, int, error) {
			return []float64{math.Cos(radOrDeg(a[0], ret.degmode))}, 1, nil
		}},
		ophandler{"tan", "Tangent of x", 1, func(a []float64) ([]float64, int, error) {
			return []float64{math.Tan(radOrDeg(a[0], ret.degmode))}, 1, nil
		}},
		ophandler{"asin", "Arcsine of x", 1, func(a []float64) ([]float64, int, error) {
			return []float64{math.Asin(radOrDeg(a[0], ret.degmode))}, 1, nil
		}},
		ophandler{"acos", "Arccosine of x", 1, func(a []float64) ([]float64, int, error) {
			return []float64{math.Acos(radOrDeg(a[0], ret.degmode))}, 1, nil
		}},
		ophandler{"atan", "Arctangent of x", 1, func(a []float64) ([]float64, int, error) {
			return []float64{math.Atan(radOrDeg(a[0], ret.degmode))}, 1, nil
		}},

		"",
		"BOLD:Miscellaneous Operations",
		ophandler{"f2c", "Convert x in Fahrenheit to Celsius", 1, func(a []float64) ([]float64, int, error) {
			return []float64{(a[0] - 32) * 5 / 9}, 1, nil
		}},
		ophandler{"c2f", "Convert x in Celsius to Fahrenheit", 1, func(a []float64) ([]float64, int, error) {
			return []float64{a[0]*9/5 + 32}, 1, nil
		}},
		"",
		"BOLD:Stack Operations",
		ophandler{"p", "Display stack", 0, func(_ []float64) ([]float64, int, error) {
			stack.print(ret.base)
			return nil, 0, nil
		}},
		ophandler{"c", "Clear stack", 0, func(_ []float64) ([]float64, int, error) {
			stack.clear()
			return nil, 0, nil
		}},
		ophandler{"=", "Print top of stack (x)", 0, func(_ []float64) ([]float64, int, error) {
			stack.printTop(ret.base)
			return nil, 0, nil
		}},
		ophandler{"d", "Drop top of stack (x)", 1, func(_ []float64) ([]float64, int, error) {
			return nil, 1, nil
		}},
		ophandler{"x", "Exchange x and y", 2, func(a []float64) ([]float64, int, error) {
			return []float64{a[0], a[1]}, 2, nil
		}},

		"",
		"BOLD:Math and Physical constants",
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
		"BOLD:Computer constants",
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
		"BOLD:Program Control",
		ophandler{"dec", "Output in decimal", 0, func(_ []float64) ([]float64, int, error) {
			ret.base = 10
			ret.degmode = false
			return nil, 0, nil
		}},
		ophandler{"bin", "Output in binary", 0, func(_ []float64) ([]float64, int, error) {
			ret.base = 2
			ret.degmode = false
			return nil, 0, nil
		}},
		ophandler{"oct", "Output in octal", 0, func(_ []float64) ([]float64, int, error) {
			ret.base = 8
			ret.degmode = false
			return nil, 0, nil
		}},
		ophandler{"hex", "Output in hexadecimal", 0, func(_ []float64) ([]float64, int, error) {
			ret.base = 16
			ret.degmode = false
			return nil, 0, nil
		}},
		ophandler{"deg", "All angles in degrees", 0, func(_ []float64) ([]float64, int, error) {
			ret.base = 10
			ret.degmode = true
			return nil, 0, nil
		}},
		ophandler{"rad", "All angles in radians", 0, func(_ []float64) ([]float64, int, error) {
			ret.base = 10
			ret.degmode = false
			return nil, 0, nil
		}},
		ophandler{"debug", "Toggle debugging", 0, func(_ []float64) ([]float64, int, error) {
			ret.debug = !ret.debug
			fmt.Printf("Debugging state: %v\n", ret.debug)
			return nil, 0, nil
		}},
		"",
		"BOLD:Please Note:",
		"  - x means the number at the top of the stack",
		"  - y means the second number from the top of the stack",
	}
	return ret
}

// operation performs an operation on the stack and returns a slice of elements
// added to the stack and the number of elements removed from the stack.
func operation(handler ophandler, stack *stackType) ([]float64, int, error) {
	// Make sure we have enough arguments in the list.
	length := len(stack.list)
	if length < handler.numArgs {
		return nil, 0, fmt.Errorf("this operation requires at least %d items in the stack", handler.numArgs)
	}

	// args contains a copy of all elements in the stack reversed.  This makes
	// it easier for functions to use x as a[0], y as a[1], etc.
	args := []float64{}
	for ix := length - 1; ix >= 0; ix-- {
		args = append(args, stack.list[ix])
	}

	ret, remove, err := handler.fn(args)
	if err != nil {
		return nil, 0, err
	}
	// Remove the number of arguments this operation consumes if needed.
	if remove > 0 && len(stack.list) < remove {
		return nil, 0, fmt.Errorf("(internal) operation %q wants to pop %d items, but we only have %d", handler.op, remove, len(stack.list))
	}

	stack.list = stack.list[0 : len(stack.list)-remove]

	// Add the return values from the function to the stack if we have any.
	if len(ret) > 0 {
		stack.push(ret...)
	}
	return ret, remove, nil
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

// help displays the help message to the screen based on the contents of opmap.
func (x opsType) help() error {
	pager, err := newPager()
	if err != nil {
		return err
	}
	if !pager.colorSupport {
		color.NoColor = true
	}
	bold := color.New(color.Bold).SprintFunc()

	for _, v := range x.ops {
		// ophandler lines.
		if handler, ok := v.(ophandler); ok {
			fmt.Fprintf(pager.w, "  - %s: %s\n", bold(handler.op), handler.desc)
			continue
		}
		// Regular strings.
		// Anything starting with "BOLD:" is printed in bold.
		if s, ok := v.(string); ok {
			if strings.HasPrefix(s, "BOLD:") {
				s = bold(s[5:])
			}
			fmt.Fprintln(pager.w, s)
		}
	}
	// Turn color support back on.
	color.NoColor = false
	return pager.wait()
}
