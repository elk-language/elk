package elk

import (
	"fmt"
	"os"

	_ "github.com/elk-language/elk/ext/std"
	"github.com/elk-language/elk/lexer"
	lexerRuntime "github.com/elk-language/elk/lexer/runtime"
	astRuntime "github.com/elk-language/elk/parser/ast/runtime"
	parserRuntime "github.com/elk-language/elk/parser/runtime"
	diagnosticRuntime "github.com/elk-language/elk/position/diagnostic/runtime"
	"github.com/elk-language/elk/types/checker"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

func InitGlobalEnvironment() {
	value.InitGlobalEnvironment()
	vm.InitGlobalEnvironment()
	astRuntime.InitGlobalEnvironment()
	diagnosticRuntime.InitGlobalEnvironment()
	lexerRuntime.InitGlobalEnvironment()
	parserRuntime.InitGlobalEnvironment()
}

// Interpret the given file.
// Panics when the file does not compile.
func Interpret(fileName string) (result value.Value, stackTrace *value.StackTrace, err value.Value) {
	bytecode, diagnostics := checker.CheckFile(fileName, nil, false, nil)
	if diagnostics != nil {
		fmt.Println()

		diagnosticString, err := diagnostics.HumanString(true, lexer.Colorizer{})
		if err != nil {
			panic(err)
		}
		fmt.Println(diagnosticString)
		if diagnostics.IsFailure() {
			panic("failed compilation")
		}
	}

	v := vm.New()
	result, err = v.InterpretTopLevel(bytecode)
	stackTrace = v.ErrStackTrace()
	return result, stackTrace, err
}

// Interpret the given file.
// Panics when the file does not compile or encounters a runtime error
func MustInterpret(fileName string) (result value.Value) {
	result, stackTrace, elkErr := Interpret(fileName)
	if !elkErr.IsUndefined() {
		vm.PrintError(os.Stderr, stackTrace, elkErr)
		panic("failed execution")
	}
	return result
}
