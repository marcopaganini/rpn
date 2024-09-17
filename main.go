// This file is part of rpn, a silly RPN calculator for the CLI.
// For further information, check https://github.com/marcopaganini/rcalc
//
// (C) 2024 by Marco Paganini <paganini AT paganini DOT net>
package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
)

var (
	// Build is filled by go build -ldflags during build.
	Build        string
	programTitle = "rpn - a simple CLI RPN calculator"
)

type (
	// ophandler contains the handler for a single operation.  numArgs
	// indicates how many items at the top of the stack should be popped and
	// sent to the function. ignoreResult will cause the function return to be
	// ignored (i.e, not be repushed into the stack and the top of the stack
	// won't be printed afterwards. This is used for commands like "show stack"
	// that don't change the stack at all.
	ophandler struct {
		desc         string // operation description (used by help)
		numArgs      int    // Number of arguments to function
		ignoreResult bool   // Ignore results from function
		fn           func([]float64) (float64, error)
	}
)

// atof takes a string as an argument and return a float64 representing that
// string. Strings starting in 0x or 0X are treated as hex strings.  Strings
// starting in o or 0 are treated as octal strings.
func atof(s string) (float64, error) {
	base := 10
	switch {
	case (strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X")) && len(s) > 2:
		s = s[2:]
		base = 16
	case (strings.HasPrefix(s, "0") || strings.HasPrefix(s, "o")) && len(s) > 1:
		s = s[1:]
		base = 8
	}

	if base == 10 {
		return strconv.ParseFloat(s, 64)
	}
	ret, err := strconv.ParseUint(s, base, 64)
	return float64(ret), err
}

func main() {
	var debug bool

	// Default == decimal for printouts
	base := 10

	stack := &stackType{}

	// Operations
	ops := map[string]ophandler{
		// Basic operations
		"+": {"Add x to y", 2, false, func(a []float64) (float64, error) { return a[0] + a[1], nil }},
		"-": {"Subtract x from y", 2, false, func(a []float64) (float64, error) { return a[0] - a[1], nil }},
		"*": {"Multiply x and y", 2, false, func(a []float64) (float64, error) { return a[0] * a[1], nil }},
		"/": {"Divide y by x", 2, false, func(a []float64) (float64, error) {
			if a[1] == 0 {
				return 0, errors.New("can't divide by zero")
			}
			return a[0] / a[1], nil
		}},
		"chs": {"Change signal of x", 1, false, func(a []float64) (float64, error) { return a[0] * -1, nil }},
		"inv": {"Invert x (1/x)", 1, false, func(a []float64) (float64, error) { return 1 / a[0], nil }},
		"^":   {"Raise y to the power of x", 2, false, func(a []float64) (float64, error) { return math.Pow(a[0], a[1]), nil }},
		"mod": {"Calculates y modulo x", 2, false, func(a []float64) (float64, error) { return math.Mod(a[0], a[1]), nil }},
		"%":   {"Calculate x% of y", 2, false, func(a []float64) (float64, error) { return a[0] * a[1] / 100, nil }},

		"fac": {"Calculate factorial of x", 1, false, func(a []float64) (float64, error) {
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

		// stack operations
		"p": {"Display stack", 0, true, func(_ []float64) (float64, error) { stack.print(base); return 0, nil }},
		"c": {"Clear stack", 0, true, func(_ []float64) (float64, error) { stack.clear(); return 0, nil }},
		"=": {"Print top of stack (x)", 0, true, func(_ []float64) (float64, error) { stack.printTop(base); return 0, nil }},
		"d": {"Drop top of stack (x)", 1, true, func(_ []float64) (float64, error) { return 0, nil }},

		// math & physical constants
		"PI":  {"The famous transcedental number", 0, false, func(_ []float64) (float64, error) { return math.Pi, nil }},
		"E":   {"Another famous transcedental number", 0, false, func(_ []float64) (float64, error) { return math.E, nil }},
		"C":   {"Speed of light in vacuum, in m/s", 0, false, func(_ []float64) (float64, error) { return 299792458, nil }},
		"MOL": {"Avogadro's number", 1, false, func(_ []float64) (float64, error) { return 6.02214154e23, nil }},

		// program control
		"dec": {"Output in decimal", 0, true, func(_ []float64) (float64, error) { base = 10; return 0, nil }},
		"hex": {"Output in hexadecimal", 0, true, func(_ []float64) (float64, error) { base = 16; return 0, nil }},
		"oct": {"Output in octal", 0, true, func(_ []float64) (float64, error) { base = 8; return 0, nil }},

		"debug": {"Toggle debugging", 0, true, func(_ []float64) (float64, error) {
			debug = !debug
			fmt.Printf("Debugging state: %v\n", debug)
			return 0, nil
		}},
	}

	rl, err := readline.New("> ")
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	// Wait for entry until Ctrl-D or q is issued
	for {
		// Save a copy of the stack so we can restore it to the previous state
		// before this line was processed (in case of errors.)
		stack.save()

		if debug {
			stack.print(base)
		}

		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		line = strings.TrimSpace(line)

		// Split into fields and process
		autoprint := false
		for _, token := range strings.Fields(line) {
			// Check operator map
			handler, ok := ops[token]
			if ok {
				err := stack.operation(handler)
				if err != nil {
					color.Red("ERROR:", err)
					stack.restore()
					break
				}
				// If the particular handler does not ignore results from the
				// function, set autoprint to true. This will cause the top of
				// the stack results to be printed.
				autoprint = !handler.ignoreResult
				continue
			}

			// Help
			if token == "help" || token == "h" || token == "?" {
				bold := color.New(color.Bold).SprintFunc()

				fmt.Println(bold("Online help for ", programTitle, "."))
				fmt.Println(bold("See http://github.com/marcopaganini/rpn for full details."))
				fmt.Println()
				fmt.Println(bold("Data entry:"))
				fmt.Println("  number <ENTER> - push a number on top of the stack.")
				fmt.Println("  operation <ENTER> - perform an operation on the stack (see below).")
				fmt.Println()
				fmt.Println("  It's also possible to separate multiple operations with space:")
				fmt.Println("    10 2 3 * - (result = 4)")
				fmt.Println()
				fmt.Println("  Prefix numbers with 0x to indicate hexadecimal, 0 for octal.")
				fmt.Println()
				fmt.Println(bold("Operations:"))

				var keys []string
				for k := range ops {
					keys = append(keys, k)
				}
				sort.Strings(keys)
				for _, k := range keys {
					fmt.Printf("  - %s: %s\n", bold(k), ops[k].desc)
				}
				fmt.Println()
				fmt.Println(bold("Please Note:"))
				fmt.Println("  - x means the number at the top of the stack")
				fmt.Println("  - y means the second number from the top of the stack")
				continue
			}

			if token == "quit" || token == "exit" || token == "q" {
				fmt.Printf("Bye.\n")
				os.Exit(0)
			}

			// At this point, it's either a number or not recognized.
			n, err := atof(token)
			if err != nil {
				fmt.Printf("Not a number or operator: %q. Use \"help\" for online help.\n", token)
				stack.restore()
				continue
			}
			// Valid number
			stack.push(n)
			continue
		}
		if autoprint {
			stack.printTop(base)
		}
	}
}
