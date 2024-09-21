// This file is part of rpn, a simple and useful CLI RPN calculator.
// For further information, check https://github.com/marcopaganini/rpn
//
// (C) Sep/2024 by Marco Paganini <paganini AT paganini DOT net>
package main

import (
	"fmt"
	"math"

	"github.com/fatih/color"
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
func (x *stackType) push(n ...float64) {
	x.list = append(x.list, n...)
}

// clear clears the stack.
func (x *stackType) clear() {
	x.list = []float64{}
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
	color.Cyan("= %s", formatNumber(x.top(), base))
}

// print displays the contents of the stack using the base indicated.
func (x *stackType) print(base int) {
	bold := color.New(color.Bold).SprintFunc()
	last := len(x.list) - 1

	fmt.Println(bold("===== Stack ====="))
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
	// default format when decimals are present.
	decfmt := "%.8f"
	if base == 10 && math.Floor(n) == n {
		decfmt = "%.f"
	}

	// Indicate possible truncation
	suffix := ""
	if base != 10 && math.Floor(n) != n {
		suffix = fmt.Sprintf("(truncated from %v)", n)
	}

	switch {
	case base == 2:
		return fmt.Sprintf("0b%b %s", uint64(n), suffix)
	case base == 8:
		return fmt.Sprintf("0%o %s", uint64(n), suffix)
	case base == 16:
		return fmt.Sprintf("0x%x %s", uint64(n), suffix)
	default:
		return fmt.Sprintf(decfmt+" %s", n, suffix)
	}
}
