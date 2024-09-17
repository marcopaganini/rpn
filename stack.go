// This file is part of rpn, a silly RPN calculator for the CLI.
// For further information, check https://github.com/marcopaganini/rcalc
//
// (C) 2024 by Marco Paganini <paganini AT paganini DOT net>
package main

import (
	"fmt"
	"math"
)

type (
	// stackType holds the representation of the RPN stack. It contains
	// two stacks, "list" (the main stack), and "savedList", which is
	// used to save the stack and later restore it in case of error.
	stackType struct {
		list      []float64
		savedList []float64
	}
)

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

// clear clears the stack.
func (x *stackType) clear() {
	x.list = []float64{}
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
	// Remove the number of arguments this operation consumes if needed.
	if len(x.list) > 0 {
		x.list = x.list[0 : length-handler.numArgs]
	}
	// Add the return from the function to the stack if ignoreResult = false.
	// Re-check stack length first, as the function may have changed it.
	if !handler.ignoreResult {
		x.push(ret)
	}
	return nil
}

// top returns the topmost element on the stack (without popping it).
func (x *stackType) top() float64 {
	if len(x.list) == 0 {
		return 0
	}
	return x.list[len(x.list)-1]
}

// printTop displays the top of the stack using the base indicated.
func (x *stackType) printTop(base int) {
	last := len(x.list) - 1
	fmt.Println("=", formatNumber(x.list[last], base))
}

// print displays the contents of the stack using the base indicated.
func (x *stackType) print(base int) {
	last := len(x.list) - 1
	fmt.Println("===== Stack =====")
	for ix := last; ix >= 0; ix-- {
		tag := fmt.Sprintf("%2d", ix)
		switch ix {
		case last:
			tag = " x"
		case last - 1:
			tag = " y"
		}
		fmt.Printf("%s: %s\n", tag, formatNumber(x.list[ix], base))
	}
}

// formatNumber formats the number using base. For bases different than 10,
// non-integer floating numbers are truncated.
func formatNumber(n float64, base int) string {
	// Indicate possible truncation
	suffix := ""
	if base != 10 && math.Floor(n) != n {
		suffix = "(truncated)"
	}

	switch {
	case base == 8:
		return fmt.Sprintf("0%o %s", uint64(n), suffix)
	case base == 16:
		return fmt.Sprintf("0x%x %s", uint64(n), suffix)
	default:
		return fmt.Sprintf("%v %s", n, suffix)
	}
}
