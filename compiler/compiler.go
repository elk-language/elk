// Package Compiler implements
// compilers that turn Elk source code
// into bytecode or Go source code
package compiler

import (
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

type Compiler interface {
	Bytecode() *vm.BytecodeFunction
	Method() value.Method
	Parent() Compiler
	SetParent(Compiler)
	InitMainCompiler()
	InitGlobalEnv() Compiler
	FinishGlobalEnvCompiler()
	CreateMainCompiler(checker types.Checker, loc *position.Location, errors *diagnostic.SyncDiagnosticList, output io.Writer) Compiler
	CompileClassInheritance(*types.Class, *position.Location)
	CompileIvarIndices(target types.NamespaceWithIvarIndices, location *position.Location)
	CompileInclude(target types.Namespace, mixin *types.Mixin, location *position.Location)
	InitExpressionCompiler(location *position.Location) Compiler
	CompileExpressionsInFile(node *ast.ProgramNode)
	InitMethodCompiler(location *position.Location) (Compiler, int)
	CompileMethods(location *position.Location, execOffset int)
	InitIvarIndicesCompiler(location *position.Location) (Compiler, int)
	FinishIvarIndicesCompiler(location *position.Location, execOffset int) Compiler
	CompileConstantDeclaration(node *ast.ConstantDeclarationNode, namespace types.Namespace, constName value.Symbol)
	RegisterMethod(node *ast.MethodDefinitionNode)
	CompileMethodBody(node *ast.MethodDefinitionNode, name value.Symbol) Compiler
	Flush() // Outputs the compiled code to an output file
}

func CreateCompiler(funcName string, parent Compiler, checker types.Checker, loc *position.Location, errors *diagnostic.SyncDiagnosticList) Compiler {
	switch parent := parent.(type) {
	case *BytecodeCompiler, nil:
		cmp := NewBytecodeCompiler(funcName, topLevelBytecodeCompilerMode, loc, checker)
		cmp.Errors = errors
		cmp.SetParent(parent)
		return cmp
	case *GoCompiler:
		goName := MangleGoIdentifier(funcName)
		cmp := NewGoCompiler(funcName, goName, topLevelGoCompilerMode, loc, checker, parent.globalData, parent.output)
		cmp.Errors = errors
		cmp.SetParent(parent)
		return cmp
	default:
		panic(fmt.Sprintf("invalid parent compiler: %T", parent))
	}
}

// sort elements by their `inspect` string
func inspectSort[V value.Inspectable](elements []V) []V {
	slices.SortStableFunc(elements, func(a, b V) int {
		return strings.Compare(a.Inspect(), b.Inspect())
	})
	return elements
}

func mergeRegexDiagnostics(target *diagnostic.SyncDiagnosticList, src error, loc *position.Location) bool {
	if src == nil {
		return false
	}

	errList, ok := src.(diagnostic.DiagnosticList)
	if !ok {
		target.AddFailure(src.Error(), loc)
		return true
	}

	regexStartPos := loc.StartPos
	for _, err := range errList {
		errStartPos := err.Span.StartPos
		errEndPos := err.Span.EndPos

		columnDifference := regexStartPos.Column - 1 + 2 // add 2 to account for `%/`
		byteDifference := regexStartPos.ByteOffset + 2   // add 2 to account for `%/`
		lineDifference := regexStartPos.Line - 1

		if errStartPos.Line == 1 {
			errStartPos.Column += columnDifference
		}
		errStartPos.Line += lineDifference
		errStartPos.ByteOffset += byteDifference

		if errEndPos != errStartPos {
			if errEndPos.Line == 1 {
				errEndPos.Column += columnDifference
			}
			errEndPos.Line += lineDifference
			errEndPos.ByteOffset += byteDifference
		}
		err.Location.FilePath = loc.FilePath

		target.Append(err)
	}

	return true
}
