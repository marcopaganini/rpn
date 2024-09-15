![GolangCI](https://github.com/marcopaganini/rpn/actions/workflows/golangci-lint.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/marcopaganini/rpn)](https://goreportcard.com/report/github.com/marcopaganini/rpn)

# RPN - A simple and useful CLI RPN calculator.

## Description

`rpn` is a simple, CLI
[RPN](https://en.wikipedia.org/wiki/Reverse_Polish_notation) calculator for
Linux. `rpn` was written in Go and is statically compiled, requiring no external
libraries. Installation is as simple as running one command and should work on
any Linux distribution.

This is a work in progress, but `rpn` already supports a number of operations and
should be useful for daily work.

## Motivation

I've been using RPN calculators for a long time and the lack of a simple and
convenient CLI version (think `bc`, but using RPN) frustrated me. I've used
`dc` many times, but it has some quirks that quickly stand in the way of
productivity, like having to print the results of the stack every single time,
no editing in the command-line, etc.

Looking around the Debian repos and Github revealed some possible alternatives,
but they also come with their own set of inconveniences, like:

* Requiring a GUI.
* Requiring the installation of many dependencies (npm, Java runtime, etc).
* Not allowing editing of the command-line, or recall of a previous line.
* Using a TUI (and making editing more complicated).
* Buggy or incomplete implementations.
* Too complex.

The idea of `rpn` is to be simple and functional enough to be your daily driver
RPN calculator. :)

## Installation

### Automatic process

To download and install automatically, just run:

```bash
wget -q -O/tmp/install \
  'https://raw.githubusercontent.com/marcopaganini/installer/master/install.sh' && \
  sudo sh /tmp/install marcopaganini/rpn
```

### Manual process

Just navigate to the [releases page](https://github.com/marcopaganini/rpn/releases) and download the desired
version. Unpack the tar file into `/usr/local/bin` and run a `chmod 755
/usr/local/bin/rpn`.  If typing `rpn` doesn't work after that, make sure
`/usr/local/bin` is in your PATH. In some distros you may need to create
`/usr/local/bin` first.

### Using go

If you have go installed, just run:

```
go install github.com/marcopaganini/rpn@latest
```

## Similar projects

There are *many* other calculators for Linux, and some of them supporting RPN.
This is a list of a few that support RPN:

* [dc](https://www.wikiwand.com/en/articles/Dc_%28computer_program%29): The
  venerable `dc` calculator. One of the oldest Unix programs. Readily available
  in most distributions.
* [dcim](https://github.com/43615/dcim): An improved version of `dc`.
* [orpie](https://github.com/pelzlpj/orpie): A Curses-based RPN calculator
  (TUI).
* [qalculate](https://qalculate.github.io/): A very complete calculator with
  CLI, GTK, and QT versions. Focus is mostly on the UI part. The author doesn't
  seem to use RPN so the RPN mode is a bit rough around the edges.

## Contributions

Feel free to open issues, send ideas and PRs.
