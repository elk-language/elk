package main

import (
	"fmt"
	"os"

	"path/filepath"

	"github.com/elk-language/elk/bytecode"
	"github.com/elk-language/elk/object"
	"github.com/elk-language/elk/vm"
)

// Main entry point to the interpreter.
func main() {
	args := os.Args[1:]
	if len(args) > 1 {
		fmt.Println("Usage: elk [script]")
		os.Exit(64)
	} else if len(args) == 1 {
		runFile(args[0])
	} else {
		// runRepl()
		chunk := &bytecode.Chunk{
			Instructions: []byte{
				// byte(bytecode.CONSTANT8),
				// 0,
				// byte(bytecode.CONSTANT8),
				// 1,
				// byte(bytecode.ADD),
				// byte(bytecode.RETURN),
				byte(bytecode.CONSTANT8),
				0x0,
				byte(bytecode.NEGATE),
				byte(bytecode.RETURN),
			},
			Constants: []object.Value{
				object.Int8(5),
				object.String("foo"),
				object.String("bar"),
			},
		}

		// chunk.Disassemble(os.Stdout)
		v := vm.New()
		v.InterpretBytecode(chunk)
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

// Start the Elk Read Evaluate Print Loop.
func runRepl() {
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
