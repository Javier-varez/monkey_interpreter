package repl

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/javier-varez/monkey_interpreter/compiler"
	"github.com/javier-varez/monkey_interpreter/evaluator"
	"github.com/javier-varez/monkey_interpreter/lexer"
	"github.com/javier-varez/monkey_interpreter/object"
	"github.com/javier-varez/monkey_interpreter/parser"
	"github.com/javier-varez/monkey_interpreter/vm"
	"github.com/peterh/liner"
)

const PROMPT = ">> "

type PromptReader struct {
	linerState *liner.State
	buffer     bytes.Buffer
}

func (p *PromptReader) Prompt() (string, error) {
	p.buffer.Reset()

	for {
		txt, err := p.linerState.Prompt(PROMPT)
		if err != nil {
			return "", err
		}

		if len(txt) > 0 && txt[len(txt)-1] == '\\' {
			p.buffer.WriteString(txt[:len(txt)-1])
			p.buffer.WriteByte('\n')
		} else {
			p.buffer.WriteString(txt)
			break
		}
	}

	entry := p.buffer.String()
	p.linerState.AppendHistory(entry)
	return entry, nil
}

func (p *PromptReader) SaveHistoryFile() {
	pwd, err := os.Getwd()
	if err != nil {
		return
	}

	monkeyFilePath := filepath.Join(pwd, ".monkey")
	monkeyFile, err := os.Create(monkeyFilePath)
	if err != nil {
		return
	}

	p.linerState.WriteHistory(monkeyFile)
}

func newPromptReader() PromptReader {
	s := liner.NewLiner()

	pwd, err := os.Getwd()
	if err != nil {
		return PromptReader{linerState: s}
	}

	monkeyFilePath := filepath.Join(pwd, ".monkey")
	stat, err := os.Stat(monkeyFilePath)
	if err != nil || stat.IsDir() {
		return PromptReader{linerState: s}
	}

	monkeyFile, err := os.Open(monkeyFilePath)
	if err != nil {
		return PromptReader{linerState: s}
	}

	s.ReadHistory(monkeyFile)
	return PromptReader{linerState: s}
}

func Start(useVm bool) {
	linerState := newPromptReader()
	defer linerState.SaveHistoryFile()

	env := object.NewEnvironment()
	constants := []object.Object{}
	symbolTable := compiler.NewSymbolTable()
	for i, builtin := range object.Builtins {
		symbolTable.DefineBuiltin(i, builtin.Name)
	}

	globals := make([]object.Object, vm.GLOBALS_SIZE)

	if useVm {
		fmt.Println("Started REPL with VM")
	} else {
		fmt.Println("Started REPL with interpreter")
	}

	for {
		txt, err := linerState.Prompt()
		if err != nil {
			fmt.Printf("Error from liner: %v\n", err)
			return
		}

		lex := lexer.New(txt)
		p := parser.New(lex)

		program := p.ParseProgram()

		if len(program.Diagnostics) != 0 {
			fmt.Print("Diagnostics:\n\n")
			for _, diag := range program.Diagnostics {
				fmt.Println(diag.ContextualError())
			}
			continue
		}

		if useVm {
			c := compiler.NewWithState(constants, symbolTable)
			err := c.Compile(program)
			if err != nil {
				fmt.Printf("Error from compiler: %s\n", err)
				continue
			}

			bytecode := c.Bytecode()
			constants = bytecode.Constants
			vmInst := vm.NewWithGlobalKeyStore(bytecode, globals)
			err = vmInst.Run()
			if err != nil {
				fmt.Printf("Error from vm: %s\n", err)
				continue
			}

			lastPopped := vmInst.LastPoppedStackElem()
			fmt.Printf("%v\n", lastPopped.Inspect())
		} else {
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
