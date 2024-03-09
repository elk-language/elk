package main

import (
	"flag"
	"fmt"
	"os"

	"path/filepath"

	"github.com/elk-language/elk/repl"
)

// Main entry point to the interpreter.
func main() {
	command := os.Args[1]
	switch command {
	case "repl":
		fs := flag.NewFlagSet("repl", flag.ContinueOnError)
		disassemble := fs.Bool("disassemble", false, "run the REPL in disassembler mode")
		inspectStack := fs.Bool("inspect-stack", false, "print the stack after each iteration of the REPL")
		parse := fs.Bool("parse", false, "run the REPL in parser mode")
		lex := fs.Bool("lex", false, "run the REPL in lexer mode")
		fs.Parse(os.Args[2:])
		repl.Run(*disassemble, *inspectStack, *parse, *lex)
	case "run":
		runFile(os.Args[2])
	default:
		os.Exit(64)
	}
}

// Attempt to execute the given file.
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

	runSourceWithName(absFileName, source)
}

// Run the given string of source code with
// the specified name.
func runSourceWithName(sourceName string, source []byte) {
	// ast, err := parser.Parse(source)
	// pp.Println(ast)
	// pp.Println(err)

	// lex := lexer.NewWithName(sourceName, source)
	// for {
	// 	token := lex.Next()
	// 	pp.Println(token)

	// 	if token.Type == token.END_OF_FILE {
	// 		break
	// 	}
	// }

}

// Run the given slice of bytes containing
// Elk source code.
func runSource(source []byte) {
	runSourceWithName("(eval)", source)
}
