package main

import (
	"os"

	"github.com/javier-varez/monkey_interpreter/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
