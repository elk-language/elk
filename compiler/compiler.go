// Package Compiler implements
// compilers that turn Elk source code
// into bytecode or Go source code
package compiler

import (
	"fmt"
	"io"

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
	CompileMethodBody(node *ast.MethodDefinitionNode, name value.Symbol) Compiler
	Flush() // Outputs the compiled code to an output file
}

func CreateCompiler(parent Compiler, checker types.Checker, loc *position.Location, errors *diagnostic.SyncDiagnosticList) Compiler {
	switch parent := parent.(type) {
	case *BytecodeCompiler, nil:
		compiler := NewBytecodeCompiler(loc.FilePath, topLevelBytecodeCompilerMode, loc, checker)
		compiler.Errors = errors
		compiler.SetParent(parent)
		return compiler
	case *GoCompiler:
		compiler := NewGoCompiler(loc.FilePath, topLevelGoCompilerMode, loc, checker, parent.bigIntCache, parent.symbolCache, parent.output)
		compiler.Errors = errors
		compiler.SetParent(parent)
		return compiler
	default:
		panic(fmt.Sprintf("invalid parent compiler: %T", parent))
	}
}
