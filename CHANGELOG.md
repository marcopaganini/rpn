# Changelog

## v1.0.1 (Dec/2024)

- NEW: Limit precision back to 128 digits. The previous attempt to use 6144
  bits caused all sorts of problems with the decimal library, including some
  serious rounding issues.
- NEW: Add exp/log/ln operations (thanks https://github.com/kpbuk)
- BUGFIX: Fix precision bug on formatted numbers.
- BUGFIX: The MOL constant required a bogus argument to work.
- BUGFIX: Removed upx (binary compression) from goreleaser. This breaks Mac
  binaries.
- BUGFIX: Make sure the original value is not changed when formatting numbers.
  This caused a latent bug under some specific conditions.

## v1.0.0 (Nov/2024)

- NEW: Big change! RPN now uses decimal math. Previously, RPN used plain
  float64s for calculations. With this version, we move to decimal math, which
  is more precise and supports larger numbers. Maximum precision is 34 digits,
  and maximum representable number is 10^6144. Please note that non-decimal
  numbers are limited to the maximum a uint64 can represent (2^64).
- NEW: Added the "fmt" operator to set the output to a number of decimals
  (E.g.: 3 fmt will set the output formatting to 3 decimals.)
- NEW: Another new operator, "dup", duplicates the top of the stack.
- NEW: Now available on homebrew for installation! Look at the [README.md] file
  for instructions.

## v0.3.1 (Oct/2024)

- BUGFIX: Introduced a small hack to properly pretty print large numbers.
  Please note that rpn uses 64-bit floats internally, so numbers above 2^53
  cannot be reliably represented with full precision.
- BUGFIX: Small fix when printing the version number in the help page.

## v0.3.0 (Oct/2024)

- NEW: Added comment support. Anything starting with '#' is a comment.
- NEW: Added a demo GIF to the main github page.
- NEW: The `help` command now prints the binary version as well.

## v0.2.0 (Oct/2024)

- NEW: Added pagination to the help screen.
- NEW: Added trigonometric functions.
- NEW: Remove all irrelevant characters from input. It's now possible to type numbers with
  commas, dollar signs, etc. Please note that this assumes the decimal point to be "."
- NEW: Added f2c and c2f (Fahrenheit/Celsius conversion).
- BUGFIX: Fix parsing of decimal fractions (0.xx).

## v0.1.0

- Initial release
