package repl

import (
	"fmt"
	"io"
	"log"

	"github.com/javier-varez/monkey_interpreter/lexer"
	"github.com/javier-varez/monkey_interpreter/token"
	"github.com/peterh/liner"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	linerState := liner.NewLiner()

	line := 0
	for {
		txt, err := linerState.Prompt(PROMPT)
		if err != nil {
			log.Fatalf("Error from liner: %v", err)
			return
		}

		linerState.AppendHistory(txt)

		lex := lexer.New(txt)
		for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}

		line += 1
	}
}
