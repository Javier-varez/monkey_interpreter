package repl

import (
	"fmt"
	"log"

	"github.com/javier-varez/monkey_interpreter/lexer"
	"github.com/javier-varez/monkey_interpreter/parser"
	"github.com/peterh/liner"
)

const PROMPT = ">> "

func Start() {
	linerState := liner.NewLiner()

	for {
		txt, err := linerState.Prompt(PROMPT)
		if err != nil {
			log.Fatalf("Error from liner: %v", err)
			return
		}

		linerState.AppendHistory(txt)

		lex := lexer.New(txt)
		p := parser.New(lex)

		program := p.ParseProgram()

		fmt.Println(program)
	}
}
