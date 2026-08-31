package main

import (
	"context"
	"fmt"
	"os"
	"path"

	"path/filepath"

	"github.com/elk-language/elk"
	_ "github.com/elk-language/elk"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/ext"
	"github.com/elk-language/elk/ext/std/test"
	"github.com/elk-language/elk/info"
	"github.com/elk-language/elk/info/banner"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/repl"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/vm"
	"github.com/spf13/cobra"
)

// Main entry point to the interpreter.
func main() {
	// Every subcommand terminates the process itself when the Elk program
	// fails, so an error here can only be a malformed invocation.
	if err := rootCommand().Execute(); err != nil {
		os.Exit(64)
	}
}

// Build the `elk` command tree.
func rootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:     "elk",
		Short:   "The Elk programming language toolchain",
		Long:    banner.Render(),
		Version: info.Version,
	}

	root.AddCommand(
		replCommand(),
		runCommand(),
		compileCommand(),
		testCommand(),
	)
	return root
}

// Build the `elk repl` command.
func replCommand() *cobra.Command {
	opts := &repl.Options{}
	cmd := &cobra.Command{
		Use:   "repl",
		Short: "Start an interactive Elk session",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			repl.Run(context.Background(), opts)
		},
	}

	flags := cmd.Flags()
	flags.BoolVar(&opts.Disassemble, "disassemble", false, "run the REPL in disassembler mode; compiles Elk to bytecode and prints it in a human-readable format")
	flags.BoolVar(&opts.Transpile, "transpile", false, "run the REPL in Go transpiler mode; compiles Elk to Go source code and prints it")
	flags.BoolVar(&opts.Native, "native", false, "run the REPL in native Go mode; transpiles Elk to Go, compiles it to a binary and runs it")
	flags.BoolVar(&opts.InspectStack, "inspect-stack", false, "print the stack after each iteration of the REPL")
	flags.BoolVar(&opts.Parse, "parse", false, "run the REPL in parser mode; parses Elk to an AST and prints it in a human-readable format")
	flags.BoolVar(&opts.Lex, "lex", false, "run the REPL in lexer mode; tokenizes Elk and prints the tokens in a human-readable format")
	flags.BoolVar(&opts.Typecheck, "typecheck", false, "run the REPL in type checker mode; typechecks Elk and prints the typechecked AST in a human-readable format")
	flags.BoolVar(&opts.Expand, "expand", false, "run the REPL in macro expansion mode; typechecks Elk, expands macros and prints the resulting AST in a human-readable format")

	return cmd
}

// Build the `elk run` command.
func runCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "run [file]",
		Short: "Execute an Elk file",
		Long:  "Execute an Elk file using the bytecode VM.\nWhen no file is given, `main.elk` in the current working directory is executed.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				runMain()
				return
			}
			runFile(args[0])
		},
	}
}

// Build the `elk compile` command.
func compileCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "compile [file]",
		Short: "Compile an Elk file",
		Long:  "Compile an Elk file to a binary.\nWhen no file is given, `main.elk` in the current working directory is compiled.",
		Args:  cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				compileMain()
				return
			}
			compileFile(args[0])
		},
	}
}

// Build the `elk test` command.
func testCommand() *cobra.Command {
	var (
		mainFile string
		grep     string
		paths    []string
	)

	cmd := &cobra.Command{
		Use:   "test",
		Short: "Run Elk tests",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			runTest(mainFile, grep, paths)
		},
	}

	flags := cmd.Flags()
	flags.StringVar(&mainFile, "main", "main.elk.test", "specify the main test file that loads tests")
	flags.StringVar(&grep, "grep", "", "test name filter regex pattern")
	flags.StringSliceVarP(&paths, "path", "p", []string{}, "test file name glob with an optional line number")

	return cmd
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

// Attempt to compile the given file.
func compileFile(fileName string) {
	info.CurrentMode = info.NativeMode
	_, err := elk.CompileFile(fileName)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// Attempt to compile the main file in the current working directory
func compileMain() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	mainPath := path.Join(cwd, "main.elk")
	compileFile(mainPath)
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
