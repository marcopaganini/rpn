![GolangCI](https://github.com/marcopaganini/rpn/actions/workflows/golangci-lint.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/marcopaganini/rpn)](https://goreportcard.com/report/github.com/marcopaganini/rpn)

# RPN - A simple and useful CLI RPN calculator.

## Description

`rpn` is a simple but useful, CLI
[RPN](https://en.wikipedia.org/wiki/Reverse_Polish_notation) calculator for
Linux. `rpn` is written in Go and is statically compiled, requiring no external
libraries. Installation is as simple as running one command and should work on
any Linux distribution.

This is a work in progress, but `rpn` already supports a rich set of operations
and should be useful for daily work.

![Demo Video](assets/rpn.gif)

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

To download and install automatically (under `/usr/local/bin`), just run:

```bash
curl -s \
  'https://raw.githubusercontent.com/marcopaganini/rpn/master/install.sh' |
  sudo sh -s -- marcopaganini/rpn
```

This assumes you have root equivalence using `sudo` and will possibly require you
to enter your password.

To download and install under another directory (for example, `$HOME/.local/bin`), run:

```bash
curl -s \
  'https://raw.githubusercontent.com/marcopaganini/rpn/master/install.sh' |
  sh -s -- marcopaganini/rpn "$HOME/.local/bin"
```

Note that `sudo` is not required on the second command as the installation directory
is under your home. Whatever location you choose, make sure your PATH environment
variable contains that location.

### Homebrew

RPN is available on homebrew. To install, use:

```
brew tap marcopaganini/homebrew-tap
brew install rpn
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

## Limitations and Caveats

This projects uses the excellent
[decimal](https://github.com/EricLagergren/decimal) math library, by [Eric
Lagergren](https://github.com/ericlagergren), so it inherits its strengths and
limitations:

* RPN uses [General Decimal Arithmetic](https://speleotrove.com/decimal/).
* Internally, we use the IEEE 754R Decimal128 format, with a precision a maximum
  scale of `10^128`.
* Maximum precision is 34 digits.
* When operating on non-decimal numbers, input is truncated to a `uint64`
  (maximum = `2^64`).
* We currently trim trailing fractional zeroes. This means that, for example,
  `1.23 + 1.27` will give `2.5` as a result, and not `2.50`. There are valid
  applications that require the full precision and we may add an option for that
  if enough people need it.
* `n / 0 == Infinity`
* `0 / 0 == Nan`
* Anything above `10^128 == +Infinity`

## Similar projects

There are *many* other calculators for Linux, and some of them supporting RPN.
This is a list of a few that support RPN:

* [dc](https://www.wikiwand.com/en/articles/Dc_%28computer_program%29): The
  venerable `dc` calculator. One of the oldest Unix programs. Readily available
  in most distributions. dc is somewhat quirky and does not support many
  operations, but uses arbitrary precision math and supports truly gigantic
  numbers.
* [dcim](https://github.com/43615/dcim): An improved version of `dc`.
* [orpie](https://github.com/pelzlpj/orpie): A Curses-based RPN calculator
  (TUI).
* [qalculate](https://qalculate.github.io/): A very complete calculator with
  CLI, GTK, and QT versions. Focus is mostly on the UI part. The author doesn't
  seem to use RPN so the RPN mode is a bit rough around the edges.

## Contributions

Feel free to open issues, send ideas and PRs.
