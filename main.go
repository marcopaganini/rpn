// This file is part of rcalc, a silly RPN calculator for the CLI.
// For further information, check https://github.com/marcopaganini/rcalc
//
// (C) 2024 by Marco Paganini <paganini AT paganini DOT net>
package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/chzyer/readline"
)

var (
	// Build is filled by go build -ldflags during build.
	Build string
)

type (
	stackType struct {
		list      []float64
		savedList []float64
	}
	opfunc func(float64, float64) (float64, error)
)

// Stack functions.

func (x *stackType) save() {
	x.savedList = append([]float64{}, x.list...)
}

func (x *stackType) restore() {
	x.list = append([]float64{}, x.savedList...)
}

func (x *stackType) push(n float64) {
	x.list = append(x.list, n)
}

func (x *stackType) operation(fn func(float64, float64) (float64, error)) error {
	length := len(x.list)
	if length < 2 {
		return errors.New("this operation needs at least two items in the stack")
	}
	last := x.list[length-1]
	beforeLast := x.list[length-2]

	ret, err := fn(beforeLast, last)
	if err != nil {
		return err
	}
	x.list = x.list[0 : length-2]
	x.push(ret)
	return nil
}

func (x *stackType) top() float64 {
	if len(x.list) == 0 {
		return 0
	}
	return x.list[len(x.list)-1]
}

func (x *stackType) printStacks() {
	length := len(x.list)
	fmt.Println("===== Stack =====")
	for ix := length - 1; ix >= 0; ix-- {
		fmt.Printf("%d: %f\n", ix, x.list[ix])
	}
	fmt.Println("== Saved Stack ==")
	for ix := length - 1; ix >= 0; ix-- {
		fmt.Printf("%d: %f\n", ix, x.savedList[ix])
	}
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func atof(s string) (float64, error) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func main() {
	var (
		debug bool
	)

	// Operations
	ops := map[string]opfunc{
		"+": func(a, b float64) (float64, error) { return a + b, nil },
		"-": func(a, b float64) (float64, error) { return a - b, nil },
		"*": func(a, b float64) (float64, error) { return a * b, nil },
		"/": func(a, b float64) (float64, error) {
			if b == 0 {
				return 0, errors.New("can't divide by zero")
			}
			return a / b, nil
		},
	}

	stack := &stackType{}

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
			stack.printStacks()
		}

		line, err := rl.Readline()
		if err != nil { // io.EOF
			break
		}
		line = strings.TrimSpace(line)

		// Split into fields and process
		opdone := false
		for _, token := range strings.Fields(line) {
			if isNumber(token) {
				n, err := atof(token)
				if err != nil {
					fmt.Printf("Not a number or operator: %q", token)
					stack.restore()
					break
				}
				stack.push(n)
				continue
			}

			// Check operator map
			fn, ok := ops[token]
			if ok {
				err := stack.operation(fn)
				if err != nil {
					fmt.Println("ERROR:", err)
					stack.restore()
					break
				}
				opdone = true
			}

			// General commands.
			if token == "debug" {
				debug = !debug
				fmt.Printf("Debugging state: %v\n", debug)
			}
			if token == "quit" || token == "q" {
				fmt.Printf("Bye\n")
				break
			}

		}
		if opdone {
			fmt.Println("=", stack.top())
		}
	}
}
