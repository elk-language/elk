// Package Compiler implements
// compilers that turn Elk source code
// into bytecode or Go source code
package compiler

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/types"
)

type Compiler interface {
	InitGlobalEnv() Compiler
	EmitExecInParent()
	CompileClassInheritance(*types.Class, *position.Location)
	CompileIvarIndices(target types.NamespaceWithIvarIndices, location *position.Location)
	CompileInclude(target types.Namespace, mixin *types.Mixin, location *position.Location)
	InitExpressionCompiler(location *position.Location) Compiler
	CompileExpressionsInFile(node *ast.ProgramNode)
}
