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
	// stackType holds the representation of the RPN stack. It contains
	// two stacks, "list" (the main stack), and "savedList", which is
	// used to save the stack and later restore it in case of error.
	stackType struct {
		list      []float64
		savedList []float64
	}

	// ophandler contains a single operation item.
	ophandler struct {
		numArgs int
		fn      func([]float64) (float64, error)
	}
)

// Stack functions.

// save saves the current stack in a separate structure.
func (x *stackType) save() {
	x.savedList = append([]float64{}, x.list...)
}

// restore restores the saved stack back into the main one.
func (x *stackType) restore() {
	x.list = append([]float64{}, x.savedList...)
}

// push adds a new element to the stack.
func (x *stackType) push(n float64) {
	x.list = append(x.list, n)
}

// operation performs an operation on the stack.
func (x *stackType) operation(handler ophandler) error {
	// Make sure we have enough arguments in the list.
	length := len(x.list)
	if length < handler.numArgs {
		return fmt.Errorf("this operation requires at least %d items in the stack", handler.numArgs)
	}

	// args contains a copy of the last numArgs in the stack.
	args := append([]float64{}, x.list[length-handler.numArgs:]...)

	ret, err := handler.fn(args)
	if err != nil {
		return err
	}
	// Remove the number of arguments this operation consumes and adds the return
	// from the function to the stack.
	x.list = x.list[0 : length-handler.numArgs]
	x.push(ret)
	return nil
}

// top returns the topmost element on the stack (without popping it).
func (x *stackType) top() float64 {
	if len(x.list) == 0 {
		return 0
	}
	return x.list[len(x.list)-1]
}

// printStacks prints the stacks (primary and backup).
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

// isNumber returns true if the string appears to be a number.
func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// atof converts a string to a float.
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
	ops := map[string]ophandler{
		"+": {2, func(a []float64) (float64, error) { return a[0] + a[1], nil }},
		"-": {2, func(a []float64) (float64, error) { return a[0] - a[1], nil }},
		"*": {2, func(a []float64) (float64, error) { return a[0] * a[1], nil }},
		"/": {2, func(a []float64) (float64, error) {
			if a[1] == 0 {
				return 0, errors.New("can't divide by zero")
			}
			return a[0] / a[1], nil
		}},
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
			handler, ok := ops[token]
			if ok {
				err := stack.operation(handler)
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
