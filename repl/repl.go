package repl

import (
	"bytes"
	"fmt"
	"log"

	"github.com/javier-varez/monkey_interpreter/lexer"
	"github.com/javier-varez/monkey_interpreter/parser"
	"github.com/peterh/liner"
	"gopkg.in/yaml.v2"
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

func Start() {
	linerState := PromptReader{linerState: liner.NewLiner()}

	for {
		txt, err := linerState.Prompt()
		if err != nil {
			log.Fatalf("Error from liner: %v", err)
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
		} else {
			fmt.Println("Successful ast parsing")
			fmt.Println(program)

			marshalled, err := yaml.Marshal(program)
			if err != nil {
				fmt.Println("Error marshalling program")
				continue
			}

			fmt.Println(string(marshalled))
		}
	}
}
