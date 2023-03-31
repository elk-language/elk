package main

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/elk-language/elk/lexer"
	"github.com/k0kubun/pp"
)

func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("Usage: elk [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		runPrompt()
	}
}

func runFile(fileName string) {
	absFileName, err := filepath.Abs(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not find file `%s`\n", fileName)
		os.Exit(1)
	}
	_, err = os.Stat(absFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not find file `%s`\n", absFileName)
		os.Exit(1)
	}
	source, err := os.ReadFile(absFileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file `%s`\n", absFileName)
		os.Exit(1)
	}

	runSource(source)
}

func runSource(source []byte) {
	lex := lexer.New(source)
	for {
		lexeme, err := lex.Next()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		pp.Println(lexeme)

		if lexeme.Type == lexer.LexEOF {
			break
		}
	}
}

func runPrompt() {
	var input []byte

	for {
		fmt.Print(">>> ")
		_, err := fmt.Scanln(&input)
		if err != nil {
			fmt.Println(err)
			os.Exit(65)
		}
		runSource(input)
	}
}
