package repl

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"os"

	"github.com/elk-language/elk"
	"github.com/elk-language/elk/compiler/colorize"
	"github.com/elk-language/elk/compiler/types"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/repl/prompt"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/vm"
	goprompt "github.com/elk-language/go-prompt"
	"github.com/k0kubun/pp/v3"
)

type evaluator struct {
	vm             *vm.Thread
	inspectStack   bool
	sourceMap      map[string]string
	inputIndex     int
	elkTypechecker *checker.Checker
	goTypechecker  *types.GoTypechecker
}

func (e *evaluator) sourceName() string {
	return fmt.Sprintf("<repl:%d>", e.inputIndex)
}

func (e *evaluator) deleteSource(sourceName string) {
	if e.elkTypechecker != nil && !e.elkTypechecker.DefinedMacros() {
		delete(e.sourceMap, sourceName)
	}
}

func (e *evaluator) addSource(input string) string {
	sourceName := e.sourceName()
	e.inputIndex++
	e.sourceMap[sourceName] = input
	return sourceName
}

func (e *evaluator) evaluate(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	if e.elkTypechecker == nil {
		e.elkTypechecker = checker.New()
		e.elkTypechecker.SetIncremental(true)
		e.vm = vm.New()
	}

	fn, dl := e.elkTypechecker.CheckSourceBytecode(sourceName, input)

	if dl != nil {
		fmt.Println()

		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.sourceMap)
		if err != nil {
			panic(err)
		}
		fmt.Println(str)
		isFailure := dl.IsFailure()
		e.elkTypechecker.ClearErrors()
		if isFailure {
			return
		}
	}

	value, runtimeErr := e.vm.InterpretREPL(fn)
	if !runtimeErr.IsUndefined() {
		e.vm.PrintError()
		e.vm.ResetError()
		return
	}
	fmt.Printf("=> %s\n\n", lexer.Colorize(value.Inspect()))

	if e.inspectStack {
		e.vm.InspectValueStack()
	}
}

// parses the input and prints it to the output
func (e *evaluator) parse(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	ast, dl := parser.Parse(sourceName, input)

	pp.Println(ast)
	if dl != nil {
		fmt.Println()

		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.sourceMap)
		if err != nil {
			panic(err)
		}

		fmt.Println(str)
		if dl.IsFailure() {
			return
		}
	}

}

// compiles the input to bytecode and dumps it to the output
func (e *evaluator) disassemble(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	if e.elkTypechecker == nil {
		e.elkTypechecker = checker.New()
		e.elkTypechecker.SetIncremental(true)
	}
	fn, dl := e.elkTypechecker.CheckSourceBytecode(sourceName, input)

	if dl != nil {
		fmt.Println()

		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.sourceMap)
		if err != nil {
			panic(err)
		}

		fmt.Println(str)
		isFailure := dl.IsFailure()
		e.elkTypechecker.ClearErrors()
		if isFailure {
			return
		}
	}

	fn.Disassemble(os.Stdout)
}

// compiles the input to Go source code and dumps it to the output
func (e *evaluator) transpile(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	var buff bytes.Buffer
	cmp, dl := checker.CheckSourceNative(sourceName, input, nil, &buff, vm.DefaultThreadPool)

	if dl != nil {
		fmt.Println()

		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.sourceMap)
		if err != nil {
			panic(err)
		}

		fmt.Println(str)
		if dl.IsFailure() {
			return
		}
	}

	cmp.Flush()
	result, err := format.Source(buff.Bytes())
	if err != nil {
		fmt.Println(buff.String())
		fmt.Fprintf(os.Stderr, "cannot format target go file: %s\n", err)
		return
	}

	_, err = os.Stdout.Write(colorize.ColorizeGo(result))
	if err != nil {
		panic(err)
	}
	fmt.Println()

	err = e.goTypechecker.CheckBytes(result)
	if err != nil {
		panic(err)
	}
}

// compiles the input to Go source code and executes it
func (e *evaluator) native(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	err := elk.CompileRunSource(sourceName, input)
	if err != nil {
		var dl diagnostic.DiagnosticList
		if errors.As(err, &dl) {
			fmt.Println()

			str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.sourceMap)
			if err != nil {
				panic(err)
			}

			fmt.Println(str)
		} else {
			fmt.Println(err.Error())
		}
		return
	}
}

// parses, typechecks the input and prints it to the output
func (e *evaluator) typecheck(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	if e.elkTypechecker == nil {
		e.elkTypechecker = checker.New()
		e.elkTypechecker.SetIncremental(true)
	}
	_, dl := e.elkTypechecker.CheckSource(sourceName, input)

	for _, ast := range e.elkTypechecker.ASTCache.Map {
		fmt.Println(lexer.Colorize(ast.Inspect()))
	}
	if dl != nil {
		fmt.Println()

		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.sourceMap)
		if err != nil {
			panic(err)
		}

		fmt.Println(str)
		isFailure := dl.IsFailure()
		e.elkTypechecker.ClearErrors()
		if isFailure {
			return
		}
	}

	fmt.Println("OK")
}

// parses, typechecks, expands macros and prints the AST to the output
func (e *evaluator) expand(input string) {
	sourceName := e.addSource(input)
	defer e.deleteSource(sourceName)

	if e.elkTypechecker == nil {
		e.elkTypechecker = checker.New()
		e.elkTypechecker.SetIncremental(true)
	}
	_, dl := e.elkTypechecker.CheckSource(sourceName, input)

	if dl != nil {
		fmt.Println()

		str, err := dl.HumanStringWithSourceMap(true, lexer.Colorizer{}, e.sourceMap)
		if err != nil {
			panic(err)
		}

		fmt.Println(str)
		isFailure := dl.IsFailure()
		e.elkTypechecker.ClearErrors()
		if isFailure {
			return
		}
	}

	ast, ok := e.elkTypechecker.ASTCache.GetUnsafe(sourceName)
	if !ok {
		panic(fmt.Sprintf("cannot get AST of %s in REPL", sourceName))
	}

	fmt.Println(lexer.Colorize(ast.String()))
}

// lexes the input and prints it to the output
func (e *evaluator) lex(input string) {
	tokens := lexer.Lex(input)
	pp.Println(tokens)
}

// Start the REPL.
func Run(disassemble, transpile, native, inspectStack, parse, lex, typecheck, expand bool) {
	prompt.Run(executor(disassemble, transpile, native, inspectStack, parse, lex, typecheck, expand))
}

func executor(disassemble, transpile, native, inspectStack, parse, lex, typecheck, expand bool) goprompt.Executor {
	eval := &evaluator{
		inspectStack: inspectStack,
		sourceMap:    make(map[string]string),
	}
	if lex {
		return eval.lex
	}
	if disassemble {
		return eval.disassemble
	}
	if transpile {
		checker, err := types.NewGoTypechecker()
		if err != nil {
			panic(fmt.Sprintf("go typechecker error: %s\n", err))
		}

		eval.goTypechecker = checker
		return eval.transpile
	}
	if native {
		return eval.native
	}
	if parse {
		return eval.parse
	}
	if typecheck {
		return eval.typecheck
	}
	if expand {
		return eval.expand
	}

	return eval.evaluate
}
