package main

import (
	"flag"
	"fmt"
	"os"

	"path/filepath"

	"github.com/elk-language/elk/comparer"
	"github.com/elk-language/elk/repl"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
)

// Main entry point to the interpreter.
func main() {
	testClass := value.NewClassWithOptions(value.ClassWithName("TestClass"))
	vm.Def(&testClass.MethodContainer, "hash", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
		return value.UInt64(10), nil
	})
	vm.Def(&testClass.MethodContainer, "===", func(vm *vm.VM, args []value.Value) (returnVal value.Value, err value.Value) {
		if _, ok := args[1].(*value.Object); ok {
			return value.True, nil
		}
		return value.False, nil
	}, vm.DefWithParameters("other"))

	v := vm.New()
	hmap := vm.MustNewHashMapWithCapacityAndElements(
		v,
		5,
		value.Pair{Key: value.SmallInt(5), Value: value.String("foo")},
		value.Pair{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.String("bar")},
	)
	expected := vm.MustNewHashMapWithCapacityAndElements(
		v,
		10,
		value.Pair{Key: value.SmallInt(5), Value: value.String("foo")},
		value.Pair{Key: value.NewObject(value.ObjectWithClass(testClass)), Value: value.String("bar")},
	)

	vm.HashMapSetCapacity(v, hmap, 10)
	fmt.Println(
		cmp.Equal(
			expected,
			hmap,
			comparer.Options()...,
		),
	)
	os.Exit(20)
	if len(os.Args) < 2 {
		fmt.Println("You must specify a command")
		os.Exit(64)
	}

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
