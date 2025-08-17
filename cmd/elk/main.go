package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"path/filepath"

	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/ext"
	"github.com/elk-language/elk/ext/std/test"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/repl"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/vm"
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
		typecheck := fs.Bool("typecheck", false, "run the REPL in type checker mode")
		expand := fs.Bool("expand", false, "run the REPL in macro expansion mode")
		fs.Parse(os.Args[2:])
		repl.Run(*disassemble, *inspectStack, *parse, *lex, *typecheck, *expand)
	case "run":
		if len(os.Args) < 3 {
			runMain()
		} else {
			runFile(os.Args[2])
		}
	case "test":
		if len(os.Args) < 3 {
			runMainTestFile()
		} else {
			runTestFile(os.Args[2])
		}
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

	bytecode, diagnostics := checker.CheckFile(fileName, nil, false, nil)
	if diagnostics != nil {
		fmt.Println()

		diagnosticString, err := diagnostics.HumanString(true, lexer.Colorizer{})
		if err != nil {
			panic(err)
		}
		fmt.Println(diagnosticString)
		if diagnostics.IsFailure() {
			os.Exit(1)
		}
	}

	v := vm.New()
	_, elkErr := v.InterpretTopLevel(bytecode)
	if !elkErr.IsUndefined() {
		vm.PrintError(os.Stderr, v.ErrStackTrace(), elkErr)
		os.Exit(1)
	}
}

// Attempt to execute the main file in the current working directory
func runMain() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	mainPath := path.Join(cwd, "main.elk")
	runFile(mainPath)
}

func runTestFile(fileName string) {
	runFile(fileName)
	testExt := ext.Map["std/test"]
	if !testExt.Initialised {
		testExt.RuntimeInit()
	}

	report := test.Run()
	if report.Status() != test.TEST_SUCCESS {
		os.Exit(1)
	}
}

func runMainTestFile() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	mainPath := path.Join(cwd, "main.elk.test")
	runTestFile(mainPath)
}
