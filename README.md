# Monkey interpreter

Implementation of a [Monkey](https://monkeylang.org) interpreter. This project is derived work from
the original language developed by [Thorsten Ball](http://interpreterbook.com).

## Features

Pretty much the regular monkey language, with a few customizations:
 - Supports range expressions like `0..123`.
 - Supports variable-length arguments to functions with `fn(a, ...) { }` syntax.
 - Supports builtin functions to turn a vararg object into an array, like `fn(a, ...) { a + len(toArray(...)) }`.
 - Support a `contains` builtin that returns a boolean indicating if a `Hash` object contains a key.
 - Closures capture the environment by value, not by reference, making it truly functional.
 - Implements nicer error reporting, giving contextual information of where the error happened.
 - Apart from the interpreter, it implements the bytecode VM and a transpiler to C++, which turns out to be the fastest.
 - Saves the repl history using `liner`.
