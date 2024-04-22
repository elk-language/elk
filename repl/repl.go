package repl

import (
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/elk-language/elk/compiler"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/vm"
	"github.com/elk-language/go-prompt"
	pstrings "github.com/elk-language/go-prompt/strings"
	"github.com/k0kubun/pp"
)

// Adapter for `lexer.Lexer` that
// provides an interface compatible with go-prompt.
type Lexer struct {
	lexer.Lexer
}

func (l *Lexer) Init(input string) {
	l.Lexer = *lexer.New(input)
}

func (l *Lexer) Next() (prompt.Token, bool) {
	t := l.Lexer.Next()
	if t.Type == token.END_OF_FILE {
		return nil, false
	}

	return t, true
}

type evaluator struct {
	compiler     *compiler.Compiler
	vm           *vm.VM
	inspectStack bool
}

func (e *evaluator) evaluate(input string) {
	var currentCompiler *compiler.Compiler
	var compileErr errors.ErrorList

	if e.compiler == nil {
		currentCompiler, compileErr = compiler.CompileREPL(sourceName, input)
		e.vm = vm.New()
	} else {
		currentCompiler, compileErr = e.compiler.CompileREPL(input)
	}

	if compileErr != nil {
		fmt.Println()
		fmt.Println(compileErr.HumanStringWithSource(input, true))
		return
	}

	e.compiler = currentCompiler
	value, runtimeErr := e.vm.InterpretREPL(e.compiler.Bytecode)
	if runtimeErr != nil {
		e.vm.PrintError()
		e.vm.ResetError()
		return
	}
	fmt.Printf("=> %s\n\n", lexer.Colorize(value.Inspect()))

	if e.inspectStack {
		e.vm.InspectStack()
	}
}

// compiles the input to bytecode and dumps it to the output
func (e *evaluator) disassemble(input string) {
	var currentCompiler *compiler.Compiler
	var compileErr errors.ErrorList
	if e.compiler == nil {
		currentCompiler, compileErr = compiler.CompileREPL(sourceName, input)
	} else {
		currentCompiler, compileErr = e.compiler.CompileREPL(input)
	}

	// ast, _ := parser.Parse("<repl>", input)
	// pp.Println(ast)
	if compileErr != nil {
		fmt.Println()
		fmt.Println(compileErr.HumanStringWithSource(input, true))
		return
	}

	e.compiler = currentCompiler

	currentCompiler.Bytecode.Disassemble(os.Stdout)
}

// parses the input and prints it to the output
func (e *evaluator) parse(input string) {
	ast, err := parser.Parse("<repl>", input)

	if err != nil {
		fmt.Println()
		fmt.Println(err.HumanStringWithSource(input, true))
		return
	}
	pp.Println(ast)
}

// lexes the input and prints it to the output
func (e *evaluator) lex(input string) {
	tokens := lexer.Lex(input)
	pp.Println(tokens)
}

// Start the REPL.
func Run(disassemble, inspectStack, parse, lex bool) {
	p := prompt.New(
		executor(disassemble, inspectStack, parse, lex),
		prompt.WithLexer(&Lexer{}),
		prompt.WithExecuteOnEnterCallback(executeOnEnter),
		prompt.WithPrefix(">> "),
	)
	p.Run()
}

// A Set of keywords that end a block of code
var blockEndKeywords = map[string]bool{
	"end": true,
}

// A Set of keywords that separate multiple blocks of code
var blockSeparatorKeywords = map[string]bool{
	"else":    true,
	"elsif":   true,
	"case":    true,
	"catch":   true,
	"finally": true,
}

// Callback triggered when the Enter key is pressed.
// Decides whether the input is complete and should be executed
// or whether a newline with indentation should be added to the buffer.
func executeOnEnter(pr *prompt.Prompt, indentSize int) (indent int, execute bool) {
	doc := pr.Buffer().Document()
	var input string
	if doc.OnLastLine() {
		input = doc.Text
	} else {
		input = doc.TextBeforeCursor()
	}

	p := parser.New(sourceName, input)
	p.Parse()

	baseIndent := doc.CurrentLineIndentSpaces()
	currentLine := doc.CurrentLine()
	lex := lexer.New(currentLine)
	firstToken := lex.Next()
	firstWord := firstToken.StringValue()
	blockEnd := blockEndKeywords[firstWord]
	blockSeparator := blockSeparatorKeywords[firstWord]

	var movedBack bool
	if blockEnd || blockSeparator {
		var indentDiff int
		var nextIndentDiff int

		indentDiff = indentSize - (baseIndent % indentSize)
		if indentDiff > baseIndent {
			indentDiff = baseIndent
		}
		if blockEnd {
			nextIndentDiff = indentDiff
		}

		if indentDiff != 0 {
			movedBack = true
		}
		toLeft := pstrings.RuneNumber(utf8.RuneCountInString(currentLine) - baseIndent + indentDiff)
		pr.CursorLeftRunes(toLeft)
		pr.InsertTextMoveCursor(currentLine[baseIndent:], false)
		pr.DeleteRunes(toLeft)
		baseIndent -= nextIndentDiff
	}

	indent = baseIndent / indentSize
	if doc.OnLastLine() {
		if p.ShouldIndent() && !movedBack {
			return indent + 1, false
		}
		if p.IsIncomplete() {
			return indent, false
		}

		return 0, true
	}

	if p.ShouldIndent() && !movedBack {
		return indent + 1, false
	}

	return indent, false
}

const (
	sourceName = "REPL"
)

func executor(disassemble, inspectStack, parse, lex bool) prompt.Executor {
	eval := &evaluator{
		inspectStack: inspectStack,
	}
	if lex {
		if inspectStack {
			fmt.Println("Lex and inspect stack modes cannot be enabled at the same time")
			os.Exit(64)
		}
		if parse {
			fmt.Println("Lex and parse modes cannot be enabled at the same time")
			os.Exit(64)
		}
		if disassemble {
			fmt.Println("Lex and disassemble modes cannot be enabled at the same time")
			os.Exit(64)
		}
		return eval.lex
	}
	if disassemble {
		if inspectStack {
			fmt.Println("Disassemble and inspect stack modes cannot be enabled at the same time")
			os.Exit(64)
		}
		if parse {
			fmt.Println("Disassemble and parse modes cannot be enabled at the same time")
			os.Exit(64)
		}
		return eval.disassemble
	}
	if parse {
		if inspectStack {
			fmt.Println("Parse and inspect stack modes cannot be enabled at the same time")
			os.Exit(64)
		}

		return eval.parse
	}

	return eval.evaluate
}
