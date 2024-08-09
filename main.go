package main

import (
	"fmt"
	"log"
	"os"

	"github.com/javier-varez/monkey_interpreter/compiler"
	"github.com/javier-varez/monkey_interpreter/evaluator"
	"github.com/javier-varez/monkey_interpreter/lexer"
	"github.com/javier-varez/monkey_interpreter/object"
	"github.com/javier-varez/monkey_interpreter/parser"
	"github.com/javier-varez/monkey_interpreter/repl"
	"github.com/javier-varez/monkey_interpreter/transpiler"
	"github.com/javier-varez/monkey_interpreter/vm"
	"github.com/spf13/cobra"
)

var rootCmd cobra.Command = cobra.Command{}

var replCmd cobra.Command = cobra.Command{
	Use: "repl",
	Run: runRepl,
}

var runCmd cobra.Command = cobra.Command{
	Use:  "run filename",
	Args: cobra.ExactArgs(1),
	Run:  runFile,
}

var compileCmd cobra.Command = cobra.Command{
	Use:  "compile filename",
	Args: cobra.ExactArgs(1),
	Run:  compileFile,
}

var useVm bool

func init() {
	replCmd.Flags().BoolVar(&useVm, "vm", false, "Instructs to use the VM instead of the interpreter")
	runCmd.Flags().BoolVar(&useVm, "vm", false, "Instructs to use the VM instead of the interpreter")
	rootCmd.AddCommand(&replCmd)
	rootCmd.AddCommand(&runCmd)
	rootCmd.AddCommand(&compileCmd)
}

func runRepl(c *cobra.Command, args []string) {
	repl.Start(useVm)
}

func runFile(c *cobra.Command, args []string) {
	fmt.Println("running file", args[0])

	txt, err := os.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
	}

	lex := lexer.New(string(txt))
	p := parser.New(lex)

	program := p.ParseProgram()

	if len(program.Diagnostics) != 0 {
		fmt.Print("Diagnostics:\n\n")
		for _, diag := range program.Diagnostics {
			fmt.Println(diag.ContextualError())
		}
	} else {
		if useVm {
			fmt.Println("Using VM")
			c := compiler.New()
			if err := c.Compile(program); err != nil {
				fmt.Println("Compilation error: ", err)
				return
			}

			bytecode := c.Bytecode()
			vm := vm.New(bytecode)
			vm.Run()
		} else {
			fmt.Println("Using interpreter")
			env := object.NewEnvironment()
			result := evaluator.Eval(program, env)
			if result != nil {
				if result.Type() == object.ERROR_VALUE_OBJ {
					err := result.(*object.Error)
					fmt.Printf("%s\n", err.ContextualError())
				}
			}
		}
	}
}

func compileFile(c *cobra.Command, args []string) {
	fmt.Println("compiling file", args[0])

	txt, err := os.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
	}

	lex := lexer.New(string(txt))
	p := parser.New(lex)

	program := p.ParseProgram()

	if len(program.Diagnostics) != 0 {
		fmt.Print("Diagnostics:\n\n")
		for _, diag := range program.Diagnostics {
			fmt.Println(diag.ContextualError())
		}
	} else {
		transpiledCode := transpiler.Transpile(program)
		runOut := transpiler.Compile(transpiledCode)
		fmt.Println(runOut)
	}
}

func main() {
	rootCmd.Execute()
}
