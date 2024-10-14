# This file is part of rpn, a simple and useful CLI RPN calculator.
# For further information, check https://github.com/marcopaganini/rpn
#
# You need vhs (https://github.com/charmbracelet/vhs) to regenerate the
# animated gif output. Once vhs is installed, run:
#
# export PS1='$ '
# unset PROMPT_COMMAND
# vhs rpn.tape
#
# (C) Oct/2024 by Marco Paganini <paganini AT paganini DOT net>

Output "rpn.gif"
Require echo

Set Shell "bash"
Set FontSize 16
Set Width  960
Set Height 600
Set Padding 5

Sleep 1s
Type "# Welcome to rpn!" Enter
Type "# A simple and useful CLI RPN calculator." Enter
Type "# You can use rpn interactively (the default) or use rpn 'expr'" Enter
Type "# To evaluate an expression directly." Enter
Sleep 2s
Type "rpn" Sleep 500ms Enter
Sleep 2s
Type "# Basic operations." Sleep 2s Enter
Type "1 2 +" Sleep 500ms Enter
Sleep 2s
Type "# Note how the stack is automatically printed (unlike other CLI calculators)." Sleep 500ms Enter
Type "# You can also use readline() style arrows to edit and recall previous lines." Sleep 500ms Enter
Sleep 2s
Type "3 4 +" Sleep 500ms Enter
Sleep 2s
Type "# Most common operations are directly supported and more..." Sleep 500ms Enter
Sleep 2s
Type "1 2 3 + +" Sleep 500ms Enter
Type "2 /" Sleep 500ms Enter
Sleep 2s
Type "10 *" Sleep 500ms Enter
Type "2 ^" Sleep 500ms Enter
Type "sqr" Sleep 500ms Enter
Type "7 mod" Sleep 500ms Enter
Sleep 2s
Type "# RPN also supports bitwise operations." Sleep 500ms Enter
Sleep 2s
Type "1 2" Sleep 500ms Enter
Type "or" Sleep 500ms Enter
Type "3 1 and" Sleep 500ms Enter
Type "3 4 xor" Sleep 500ms Enter
Type "7 4 xor" Sleep 500ms Enter
Sleep 2s
Type "# You can use common trigonometric functions in radians or degrees." Sleep 500ms Enter
Sleep 2s
Type "PI 2 / sin" Sleep 500ms Enter
Type "deg" Sleep 500ms Enter
Type "30 sin" Sleep 500ms Enter
Type "60 cos" Sleep 500ms Enter
Type "rad" Sleep 500ms Enter
Sleep 2s
Type "# You can input and output in different (and independent) bases " Sleep 500ms Enter
Type "# Use the 0x prefix for hex, 0 for octal, 0b for binary." Sleep 500ms Enter
Sleep 2s
Type "0xffff" Sleep 500ms Enter
Type "=" Sleep 500ms Enter
Sleep 2s
Type "# Note how = prints the top of the stack (in the current default base)." Sleep 500ms Enter
Sleep 2s
Type "0b1111" Sleep 500ms Enter
Type "=" Sleep 500ms Enter
Type "0100" Sleep 500ms Enter
Type "=" Sleep 500ms Enter
Sleep 2s
Type "# By default, output is in decimal, but you can change that with" Sleep 500ms Enter
Type "# the dec, hex, oct, and bin commands." Sleep 500ms Enter
Sleep 2s
Type "hex" Sleep 500ms Enter
Type "1000" Sleep 500ms Enter
Type "=" Sleep 500ms Enter
Sleep 1s
Sleep 1s Type "oct" Sleep 500ms Enter
Type "=" Sleep 500ms Enter
Sleep 1s
Type "bin" Sleep 500ms Enter
Type "=" Sleep 500ms Enter
Sleep 1s
Type "dec" Sleep 500ms Enter
Sleep 2s
Type "# Math, physical, and computer constants are built-in." Sleep 500ms Enter
Sleep 2s
Type "PI" Sleep 500ms Enter
Type "2 32 ^" Sleep 500ms Enter
Type "MIB /" Sleep 500ms Enter
Sleep 2s
Type "# Stack and program specific operations are also available." Sleep 500ms Enter
Sleep 2s
Type "# Clear stack." Sleep 500ms Enter
Sleep 2s
Type "c" Sleep 500ms Enter
Sleep 2s
Type "# Display stack" Sleep 500ms Enter
Sleep 2s
Type "1 2 3 4" Sleep 500ms Enter
Type "p" Sleep 500ms Enter
Sleep 2s
Type "# Drop top of stack" Sleep 500ms Enter
Sleep 2s
Type "d" Sleep 500ms Enter
Type "p" Sleep 500ms Enter
Sleep 2s
Type "# Exchange x (top of the stack) with y (second from the top)." Sleep 500ms Enter
Sleep 2s
Type "x" Sleep 500ms Enter
Type "p" Sleep 500ms Enter
Sleep 2s
Type "# RPN also understands some simple and useful conversion functions (more to come)..."
Sleep 2s
Sleep 500ms Enter
Sleep 2s
Type "# Example: Fahrenheit to Celsius and back." Sleep 500ms Enter
Sleep 2s
Type "104 f2c" Sleep 500ms Enter
Type "c2f" Sleep 500ms Enter
Sleep 2s
Type "# No danger of getting lost. Just use the online help (help)." Sleep 500ms Enter
Sleep 2s
Type "help" Sleep 500ms Enter
Sleep 3s
Type "G q"
Sleep 3s
Type "q" Sleep 500ms Enter
Type "# Thanks!" Enter
Sleep 4s
