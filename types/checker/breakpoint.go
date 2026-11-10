package checker

import (
	"slices"

	"github.com/elk-language/elk/compiler"
	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

type BreakpointContext struct {
	selfType       types.Type
	runtimeEnv     *types.GlobalEnvironment
	macroEnv       *types.GlobalEnvironment
	localEnvs      []*localEnvironment
	constantScopes []constantScope
	methodScopes   []methodScope
	value.ValueBase
}

var _ value.Reference = &BreakpointContext{}

func (c *BreakpointContext) CopyReference() value.Reference {
	return &BreakpointContext{
		selfType:       c.selfType,
		runtimeEnv:     c.runtimeEnv,
		macroEnv:       c.macroEnv,
		localEnvs:      c.localEnvs,
		constantScopes: c.constantScopes,
		methodScopes:   c.methodScopes,
	}
}

func (c *BreakpointContext) ToValue() value.Value {
	return value.Ref(c)
}

func (c *BreakpointContext) Inspect() string {
	return "<typechecker breakpoint context>"
}

func (c *Checker) createBreakpointContext() *BreakpointContext {
	return &BreakpointContext{
		selfType:       c.selfType,
		runtimeEnv:     c.runtimeEnv,
		macroEnv:       c.macroEnv,
		localEnvs:      deepCloneLocalEnvsForBreakpoint(c.localEnvs),
		constantScopes: c.constantScopesCopy(),
		methodScopes:   c.methodScopesCopy(),
	}
}

func (c *Checker) DumpVariablesNode(sourceName string, loc *position.Location) *ast.ProgramNode {
	var variableEntries []ast.ExpressionNode

	for name := range c.allLocals() {
		nameStr := name.String()
		pair := ast.NewSymbolKeyValueExpressionNode(
			loc,
			ast.NewPublicIdentifierNode(loc, nameStr),
			ast.NewPublicIdentifierNode(loc, nameStr),
		)
		variableEntries = append(variableEntries, pair)
	}

	node := ast.NewProgramNode(
		loc,
		ast.ExpressionToStatements(
			ast.NewHashMapLiteralNode(loc, variableEntries, nil),
		),
	)

	return node
}

func (c *Checker) DumpVariablesForBreakpoint(sourceName string, loc *position.Location) (*vm.BytecodeFunction, diagnostic.DiagnosticList) {
	return c.CheckBreakpointNode(sourceName, c.DumpVariablesNode(sourceName, loc))
}

func (c *Checker) DumpVariablesForREPL(sourceName string, loc *position.Location) (*vm.BytecodeFunction, diagnostic.DiagnosticList) {
	compiler, diagnostics := c.CheckProgramREPL(sourceName, c.DumpVariablesNode(sourceName, loc))
	return compiler.Bytecode(), diagnostics
}

func NewBreakpointChecker(context *compiler.BytecodeBreakpointContext) *Checker {
	checkerContext := context.TypecheckerContext.(*BreakpointContext)
	c := &Checker{
		Filename:       context.Location.FilePath,
		selfType:       checkerContext.selfType,
		returnType:     types.Void{},
		throwType:      types.Any{},
		mode:           methodMode,
		Errors:         new(diagnostic.SyncDiagnosticList),
		localEnvs:      deepCloneLocalEnvsForBreakpoint(checkerContext.localEnvs),
		constantScopes: slices.Clone(checkerContext.constantScopes),
		methodScopes:   slices.Clone(checkerContext.methodScopes),
		ASTCache:       concurrent.NewMap[string, *ast.ProgramNode](),
		macroEnv:       checkerContext.macroEnv,
		runtimeEnv:     checkerContext.runtimeEnv,
		threadPool:     vm.DefaultThreadPool,
	}

	c.compiler = compiler.CreateBreakpointCompiler(c, context, c.Errors)

	return c
}
