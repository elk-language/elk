package main

import (
	"fmt"
	"os"
	"path"

	"path/filepath"

	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/ext"
	"github.com/elk-language/elk/ext/std/test"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/repl"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/vm"
	"github.com/spf13/pflag"
)

// Main entry point to the interpreter.
func main() {
	command := os.Args[1]
	switch command {
	case "repl":
		fs := pflag.NewFlagSet("repl", pflag.ExitOnError)
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
		fs := pflag.NewFlagSet("test", pflag.ExitOnError)
		main := fs.String("main", "main.elk.test", "specify the main test file that loads tests")
		grep := fs.String("grep", "", "test name filter regex pattern")
		path := fs.StringSliceP("path", "p", []string{}, "test file name glob with an optional line number")
		fs.Parse(os.Args[2:])

		runTest(*main, *grep, *path)
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

	bytecode, diagnostics := checker.CheckFile(fileName, nil, bitfield.BitField16{}, nil)
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

func runTest(main string, grep string, paths []string) {
	if grep != "" {
		regexFilter, err := test.NewRegexFilter(grep)
		if err != nil {
			fmt.Printf("invalid grep: %s\n", err)
			os.Exit(1)
		}
		test.RegisterFilter(regexFilter)
	}
	for _, path := range paths {
		pathFilter, err := test.NewPathFilter(path)
		if err != nil {
			fmt.Printf("invalid path: %s\n", err)
			os.Exit(1)
		}
		test.RegisterFilter(pathFilter)
	}
	runTestFile(main)
}

func runTestFile(fileName string) {
	runFile(fileName)
	testExt := ext.Map["std/test"]
	if !testExt.Initialised {
		testExt.RuntimeInit()
	}

	report := test.Run()
	if report == nil || report.Status() != test.TEST_SUCCESS {
		os.Exit(1)
	}
}
