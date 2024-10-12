// This file is part of rpn, a simple and useful CLI RPN calculator.
// For further information, check https://github.com/marcopaganini/rpn
//
// (C) Sep/2024 by Marco Paganini <paganini AT paganini DOT net>

package main

import (
	"errors"
	"io"
	"os"
	"os/exec"
)

// struct pager contains information about a pager object.
type pager struct {
	w            io.WriteCloser
	cmd          *exec.Cmd
	colorSupport bool
}

// newPager creates a new pager object and executes the pager.  If no suitable
// pager binary is found, os.Writer will point to the standard output.
func newPager() (pager, error) {
	// Look for a pager and set output to stdout if none found.
	prog, colorSupport, err := findPager()
	if err != nil {
		return pager{
			w:            os.Stdout,
			colorSupport: colorSupport}, nil
	}

	cmd := exec.Command(prog[0], prog[1:]...)
	w, err := cmd.StdinPipe()
	if err != nil {
		return pager{}, err
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Start()

	return pager{
		w:            w,
		cmd:          cmd,
		colorSupport: colorSupport}, nil
}

// findPager returns a suitable pager program in the PATH whether it supports
// color input or not.
func findPager() ([]string, bool, error) {
	if p, err := exec.LookPath("less"); err == nil {
		return []string{p, "-R"}, true, nil
	}
	if p, err := exec.LookPath("more"); err == nil {
		return []string{p}, false, nil
	}
	return nil, false, errors.New("unable to find pager program (less, more, etc)")
}

// wait closes the input and waits for the command to finish.
func (x pager) wait() error {
	// Do nothing if we're outputting to stdout.
	if x.w == os.Stdout {
		return nil
	}
	x.w.Close()
	return x.cmd.Wait()
}
