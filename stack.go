// This file is part of rpn, a simple and useful CLI RPN calculator.
// For further information, check https://github.com/marcopaganini/rpn
//
// (C) Sep/2024 by Marco Paganini <paganini AT paganini DOT net>
package main

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
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

// commafWithDigits comes straight from humanize, but modified to call
// strconv.Formatfloat with 0 as the precision. This will print the entire
// number, and not truncate to the exact precision.
func commafWithDigits(v float64, decimals int) string {
	buf := &bytes.Buffer{}
	if v < 0 {
		buf.Write([]byte{'-'})
		v = 0 - v
	}

	comma := []byte{','}

	// HACK: If this number has decimals, use -1 for precision so we get the
	// maximum exact representation possible. If it can be represented as an
	// integer, just use 0 to format the entire number.
	prec := -1
	if v == math.Trunc(v) {
		prec = 0
	}
	parts := strings.Split(strconv.FormatFloat(v, 'f', prec, 64), ".")

	pos := 0
	if len(parts[0])%3 != 0 {
		pos += len(parts[0]) % 3
		buf.WriteString(parts[0][:pos])
		buf.Write(comma)
	}
	for ; pos < len(parts[0]); pos += 3 {
		buf.WriteString(parts[0][pos : pos+3])
		buf.Write(comma)
	}
	buf.Truncate(buf.Len() - 1)

	if len(parts) > 1 {
		buf.Write([]byte{'.'})
		buf.WriteString(parts[1])
	}
	return stripTrailingDigits(buf.String(), decimals)
}

func stripTrailingDigits(s string, digits int) string {
	if i := strings.Index(s, "."); i >= 0 {
		if digits <= 0 {
			return s[:i]
		}
		i++
		if i+digits >= len(s) {
			return s
		}
		return s[:i+digits]
	}
	return s
}

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
	// Indicate possible truncation
	suffix := ""
	// clean = double as ascii, without non-significant decimal zeroes.
	clean := humanize.FtoaWithDigits(n, 6)

	if base != 10 && math.Floor(n) != n {
		suffix = fmt.Sprintf(" (truncated from %s)", clean)
	}

	switch {
	case base == 2:
		return fmt.Sprintf("0b%b%s", uint64(n), suffix)
	case base == 8:
		return fmt.Sprintf("0%o%s", uint64(n), suffix)
	case base == 16:
		return fmt.Sprintf("0x%x%s", uint64(n), suffix)
	default:
		h := commafWithDigits(n, 6)
		// Only print humanized format when it differs from original value.
		if h != clean {
			suffix = " (" + h + ")"
		}
		return clean + suffix
	}
}
