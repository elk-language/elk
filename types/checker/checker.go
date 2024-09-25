// Package checker implements the Elk type checker
package checker

import (
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/elk-language/elk/concurrent"
	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
	"github.com/rivo/uniseg"
)

// Check the types of Elk source code.
func CheckSource(sourceName string, source string, globalEnv *types.GlobalEnvironment, headerMode bool) (*vm.BytecodeFunction, error.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	return CheckAST(sourceName, ast, globalEnv, headerMode)
}

// Check the types of an Elk AST.
func CheckAST(sourceName string, ast *ast.ProgramNode, globalEnv *types.GlobalEnvironment, headerMode bool) (*vm.BytecodeFunction, error.ErrorList) {
	checker := newChecker(sourceName, globalEnv, headerMode)
	bytecode := checker.checkProgram(ast)
	return bytecode, checker.Errors.ErrorList
}

// Check the types of an Elk file.
func CheckFile(fileName string, globalEnv *types.GlobalEnvironment, headerMode bool) (*vm.BytecodeFunction, error.ErrorList) {
	checker := newChecker(fileName, globalEnv, headerMode)
	bytecode := checker.checkFile(fileName)
	return bytecode, checker.Errors.ErrorList
}

type mode uint8

const (
	topLevelMode mode = iota
	moduleMode
	classMode
	mixinMode
	interfaceMode
	methodMode
	implicitInterfaceSubtypeMode
	closureInferReturnTypeMode
	singletonMode
	namedGenericTypeDefinitionMode
	outputPositionTypeMode
	inputPositionTypeMode
	instanceVariableMode
	inheritanceMode // active when typechecking an included mixin, implemented interface, or superclass
	inferTypeArgumentMode
	methodCompatibilityInAlgebraicTypeMode
	variablePatternMode
	valuePatternMode
)

type phase uint8

const (
	initPhase phase = iota
	constantCheckPhase
	expressionPhase
)

// Holds the state of the type checking process
type Checker struct {
	Filename                string
	Errors                  *error.SyncErrorList
	GlobalEnv               *types.GlobalEnvironment
	IsHeader                bool
	phase                   phase
	mode                    mode
	constantScopes          []constantScope
	constantScopesCopyCache []constantScope
	methodScopes            []methodScope
	methodScopesCopyCache   []methodScope
	localEnvs               []*localEnvironment
	loops                   []*loop
	returnType              types.Type
	throwType               types.Type
	selfType                types.Type
	namespacePlaceholders   *concurrent.Slice[*types.NamespacePlaceholder]
	placeholders            *concurrent.Slice[*types.Placeholder]
	methodChecks            *concurrent.Slice[methodCheckEntry]
	constantChecks          *constantDefinitionChecks
	typeDefinitionChecks    *typeDefinitionChecks
	astCache                *concurrent.Map[string, *ast.ProgramNode]
}

// Instantiate a new Checker instance.
func newChecker(filename string, globalEnv *types.GlobalEnvironment, headerMode bool) *Checker {
	if globalEnv == nil {
		globalEnv = types.NewGlobalEnvironment()
	}
	return &Checker{
		Filename:   filename,
		GlobalEnv:  globalEnv,
		IsHeader:   headerMode,
		selfType:   globalEnv.StdSubtype(symbol.Object),
		returnType: types.Void{},
		mode:       topLevelMode,
		Errors:     new(error.SyncErrorList),
		constantScopes: []constantScope{
			makeConstantScope(globalEnv.Std()),
			makeLocalConstantScope(globalEnv.Root),
		},
		methodScopes: []methodScope{
			makeLocalMethodScope(globalEnv.StdSubtypeClass(symbol.Object)),
			makeUsingMethodScope(globalEnv.StdSubtypeModule(symbol.Kernel)),
		},
		localEnvs: []*localEnvironment{
			newLocalEnvironment(nil),
		},
		namespacePlaceholders: concurrent.NewSlice[*types.NamespacePlaceholder](),
		placeholders:          concurrent.NewSlice[*types.Placeholder](),
		methodChecks:          concurrent.NewSlice[methodCheckEntry](),
		typeDefinitionChecks:  newTypeDefinitionChecks(),
		constantChecks:        newConstantDefinitionChecks(),
		astCache:              concurrent.NewMap[string, *ast.ProgramNode](),
	}
}

// Instantiate a new Checker instance.
func New() *Checker {
	return newChecker("", nil, false)
}

func (c *Checker) newMethodChecker(filename string, constScopes []constantScope, methodScopes []methodScope, selfType, returnType, throwType types.Type) *Checker {
	return &Checker{
		GlobalEnv:      c.GlobalEnv,
		Filename:       filename,
		mode:           methodMode,
		phase:          c.phase,
		selfType:       selfType,
		returnType:     returnType,
		throwType:      throwType,
		constantScopes: constScopes,
		methodScopes:   methodScopes,
		Errors:         c.Errors,
		IsHeader:       c.IsHeader,
		localEnvs: []*localEnvironment{
			newLocalEnvironment(nil),
		},
		typeDefinitionChecks: newTypeDefinitionChecks(),
	}
}

func (c *Checker) CheckSource(sourceName string, source string) (*vm.BytecodeFunction, error.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	c.Filename = sourceName
	c.methodChecks = concurrent.NewSlice[methodCheckEntry]()
	bytecodeFunc := c.checkProgram(ast)
	return bytecodeFunc, c.Errors.ErrorList
}

func (c *Checker) checkNamespacePlaceholders() {
	for _, placeholder := range c.namespacePlaceholders.Slice {
		replacement := placeholder.Namespace
		if _, ok := replacement.(*types.ModulePlaceholder); !ok {
			continue
		}

		for _, location := range placeholder.Locations.Slice {
			c.addFailureWithLocation(
				fmt.Sprintf("undefined namespace `%s`", placeholder.Name()),
				location,
			)
		}
		placeholder.Locations.Slice = nil
	}
}

func (c *Checker) checkConstantAndTypePlaceholders() {
	for _, placeholder := range c.placeholders.Slice {
		if placeholder.Checked || placeholder.Sibling != nil && placeholder.Sibling.Checked {
			continue
		}
		placeholder.Checked = true
		if placeholder.Replaced || placeholder.Sibling != nil && placeholder.Sibling.Replaced {
			continue
		}

		c.addFailureWithLocation(
			fmt.Sprintf("undefined type or constant `%s`", lexer.Colorize(placeholder.FullName)),
			placeholder.Location,
		)
	}
	c.placeholders.Slice = nil
}

func (c *Checker) checkProgram(node *ast.ProgramNode) *vm.BytecodeFunction {
	statements := node.Body

	c.hoistNamespaceDefinitions(statements)
	c.checkNamespacePlaceholders()
	c.checkTypeDefinitions()
	c.hoistMethodDefinitions(statements)
	c.checkConstantAndTypePlaceholders()
	c.phase = expressionPhase
	c.checkMethods()
	c.checkConstants()
	c.checkStatements(statements)

	return nil
}

func (c *Checker) checkFile(filename string) *vm.BytecodeFunction {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		c.addFailure(
			fmt.Sprintf(
				"cannot read file: %s",
				filename,
			),
			position.NewSpan(position.New(0, 1, 1), position.New(0, 1, 1)),
		)
	}

	source := string(bytes)
	ast, errList := parser.Parse(filename, source)
	if errList != nil {
		c.Errors.ErrorList.Join(errList)
		return nil
	}

	var wg sync.WaitGroup
	// parse all files of the project concurrently
	c.checkImportsForFile(filename, ast, &wg)
	wg.Wait()

	c.hoistNamespaceDefinitionsInFile(filename, ast)
	c.checkNamespacePlaceholders()
	c.checkTypeDefinitions()
	c.hoistMethodDefinitionsInFile(filename, ast)
	c.checkConstantAndTypePlaceholders()
	c.phase = expressionPhase
	c.checkMethods()
	c.checkConstants()
	c.checkExpressionsInFile(filename, ast)

	return nil
}

func (c *Checker) hoistNamespaceDefinitionsInFile(filename string, node *ast.ProgramNode) {
	node.State = ast.CHECKING_NAMESPACES
	for _, importPath := range node.ImportPaths {
		importedAst, ok := c.astCache.GetUnsafe(importPath)
		if !ok {
			continue
		}
		switch importedAst.State {
		case ast.CHECKING_NAMESPACES, ast.CHECKED_NAMESPACES:
			continue
		}
		c.hoistNamespaceDefinitionsInFile(importPath, importedAst)
	}

	prevFilename := c.Filename
	c.Filename = filename
	c.hoistNamespaceDefinitions(node.Body)
	c.Filename = prevFilename
	node.State = ast.CHECKED_NAMESPACES
}

func (c *Checker) hoistMethodDefinitionsInFile(filename string, node *ast.ProgramNode) {
	node.State = ast.CHECKING_METHODS
	for _, importPath := range node.ImportPaths {
		importedAst, ok := c.astCache.GetUnsafe(importPath)
		if !ok {
			continue
		}
		switch importedAst.State {
		case ast.CHECKING_METHODS, ast.CHECKED_METHODS:
			continue
		}
		c.hoistMethodDefinitionsInFile(importPath, importedAst)
	}

	prevFilename := c.Filename
	c.Filename = filename
	c.hoistMethodDefinitions(node.Body)
	c.Filename = prevFilename
	node.State = ast.CHECKED_METHODS
}

func (c *Checker) checkExpressionsInFile(filename string, node *ast.ProgramNode) {
	node.State = ast.CHECKING_EXPRESSIONS
	for _, importPath := range node.ImportPaths {
		importedAst, ok := c.astCache.GetUnsafe(importPath)
		if !ok {
			continue
		}
		switch importedAst.State {
		case ast.CHECKING_EXPRESSIONS, ast.CHECKED_EXPRESSIONS:
			continue
		}
		c.checkExpressionsInFile(importPath, importedAst)
	}

	prevFilename := c.Filename
	c.Filename = filename
	c.checkStatements(node.Body)
	c.Filename = prevFilename
	node.State = ast.CHECKED_EXPRESSIONS
}

func (c *Checker) checkImportsForFile(fileName string, ast *ast.ProgramNode, wg *sync.WaitGroup) {
	c.astCache.Set(fileName, ast)

	imports := c.hoistImports(ast.Body)
	for _, importStmt := range imports {
		ast.ImportPaths = append(ast.ImportPaths, importStmt.FsPaths...)
		for _, importPath := range importStmt.FsPaths {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, ok := c.astCache.Get(importPath)
				if ok {
					return
				}
				bytes, err := os.ReadFile(importPath)
				if err != nil {
					c.addFailure(
						fmt.Sprintf(
							"cannot read file: %s (%s)",
							importPath,
							err,
						),
						importStmt.Span(),
					)
					return
				}

				ast, errList := parser.Parse(importPath, string(bytes))
				if errList != nil {
					c.Errors.JoinErrList(errList)
					return
				}

				c.checkImportsForFile(importPath, ast, wg)
			}()
		}
	}

}

func (c *Checker) setMode(mode mode) {
	c.mode = mode
}

func (c *Checker) ClearErrors() {
	c.Errors = new(error.SyncErrorList)
}

// Create a new location struct with the given position.
func (c *Checker) newLocation(span *position.Span) *position.Location {
	return position.NewLocationWithSpan(c.Filename, span)
}

func (c *Checker) checkStatements(stmts []ast.StatementNode) (types.Type, *position.Span) {
	var lastType types.Type
	var lastTypeSpan *position.Span
	var seenNever bool
	var unreachableCodeErrorReported bool
	for _, statement := range stmts {
		var t types.Type
		t, span := c.checkStatement(statement)
		if t == nil {
			continue
		}

		if seenNever {
			if !unreachableCodeErrorReported {
				unreachableCodeErrorReported = true
				c.addUnreachableCodeError(span)
			}
			continue
		}
		lastTypeSpan = span
		lastType = t

		if types.IsNever(t) {
			seenNever = true
		}
	}

	if lastType == nil {
		return types.Nil{}, nil
	} else {
		return lastType, lastTypeSpan
	}
}

func (c *Checker) checkStatement(node ast.Node) (types.Type, *position.Span) {
	switch node := node.(type) {
	case *ast.EmptyStatementNode:
		return nil, nil
	case *ast.ExpressionStatementNode:
		node.Expression = c.checkExpression(node.Expression)
		return c.typeOf(node.Expression), node.Expression.Span()
	case *ast.ImportStatementNode:
		return nil, nil
	default:
		c.addFailure(
			fmt.Sprintf("incorrect statement type %#v", node),
			node.Span(),
		)
		return nil, nil
	}
}

func (c *Checker) checkExpressions(exprs []ast.ExpressionNode) {
	for i, expr := range exprs {
		exprs[i] = c.checkExpression(expr)
	}
}

func (c *Checker) checkExpressionsWithinModule(node *ast.ModuleDeclarationNode) {
	module, ok := c.typeOf(node).(*types.Module)
	previousSelf := c.selfType
	if ok {
		c.pushConstScope(makeLocalConstantScope(module))
		c.pushMethodScope(makeLocalMethodScope(module))
		c.pushIsolatedLocalEnv()
		c.selfType = module
	} else {
		c.selfType = types.Nothing{}
		c.addFailure(
			"module definitions cannot appear in this context",
			node.Span(),
		)
	}

	c.checkStatements(node.Body)
	c.selfType = previousSelf

	if ok {
		c.popLocalConstScope()
		c.popMethodScope()
		c.popLocalEnv()
	}
}

func (c *Checker) checkExpressionsWithinClass(node *ast.ClassDeclarationNode) {
	class, ok := c.typeOf(node).(*types.Class)
	previousSelf := c.selfType
	if ok {
		c.checkAbstractMethods(class, node.Constant.Span())
		c.pushConstScope(makeLocalConstantScope(class))
		c.pushMethodScope(makeLocalMethodScope(class))
		c.pushIsolatedLocalEnv()
		c.selfType = class.Singleton()
	} else {
		c.selfType = types.Nothing{}
		c.addFailure(
			"class definitions cannot appear in this context",
			node.Span(),
		)
	}

	c.checkStatements(node.Body)
	c.selfType = previousSelf

	if ok {
		c.popLocalConstScope()
		c.popMethodScope()
		c.popLocalEnv()
	}
}

func (c *Checker) checkExpressionsWithinMixin(node *ast.MixinDeclarationNode) {
	mixin, ok := c.typeOf(node).(*types.Mixin)
	previousSelf := c.selfType
	if ok {
		c.checkAbstractMethods(mixin, node.Constant.Span())
		c.pushConstScope(makeLocalConstantScope(mixin))
		c.pushMethodScope(makeLocalMethodScope(mixin))
		c.pushIsolatedLocalEnv()
		c.selfType = mixin.Singleton()
	} else {
		c.selfType = types.Nothing{}
		c.addFailure(
			"mixin definitions cannot appear in this context",
			node.Span(),
		)
	}

	c.checkStatements(node.Body)
	c.selfType = previousSelf

	if ok {
		c.popLocalConstScope()
		c.popMethodScope()
		c.popLocalEnv()
	}
}

func (c *Checker) checkExpressionsWithinInterface(node *ast.InterfaceDeclarationNode) {
	iface, ok := c.typeOf(node).(*types.Interface)
	previousSelf := c.selfType
	if ok {
		c.pushConstScope(makeLocalConstantScope(iface))
		c.pushMethodScope(makeLocalMethodScope(iface))
		c.pushIsolatedLocalEnv()
		c.selfType = iface.Singleton()
	} else {
		c.selfType = types.Nothing{}
		c.addFailure(
			"interface definitions cannot appear in this context",
			node.Span(),
		)
	}

	c.checkStatements(node.Body)
	c.selfType = previousSelf

	if ok {
		c.popLocalConstScope()
		c.popMethodScope()
		c.popLocalEnv()
	}
}

func (c *Checker) checkExpressionsWithinSingleton(node *ast.SingletonBlockExpressionNode) {
	class, ok := c.typeOf(node).(*types.SingletonClass)
	previousSelf := c.selfType
	if ok {
		c.pushConstScope(makeLocalConstantScope(class))
		c.pushMethodScope(makeLocalMethodScope(class))
		c.pushIsolatedLocalEnv()
		c.selfType = c.GlobalEnv.StdSubtype(symbol.Class)
	} else {
		c.selfType = types.Nothing{}
		c.addFailure(
			"singleton definitions cannot appear in this context",
			node.Span(),
		)
	}

	c.checkStatements(node.Body)
	c.selfType = previousSelf

	if ok {
		c.popLocalConstScope()
		c.popMethodScope()
		c.popLocalEnv()
	}
}

func (c *Checker) checkExpressionWithType(node ast.ExpressionNode, typ types.Type) ast.ExpressionNode {
	switch n := node.(type) {
	case *ast.ClosureLiteralNode:
		switch t := typ.(type) {
		case *types.NamedType:
			return c.checkExpressionWithType(node, t.Type)
		case *types.Closure:
			return c.checkClosureLiteralNodeWithType(n, t)
		}
	case *ast.ArrayListLiteralNode:
		generic, ok := typ.(*types.Generic)
		if !ok || generic.TypeArguments.Len() != 1 || !c.isSubtype(c.StdArrayList(), generic.Namespace, nil) {
			break
		}
		return c.checkArrayListLiteralNodeWithType(n, generic)
	case *ast.ArrayTupleLiteralNode:
		generic, ok := typ.(*types.Generic)
		if !ok || generic.TypeArguments.Len() != 1 || !c.isSubtype(c.StdArrayTuple(), generic.Namespace, nil) {
			break
		}
		return c.checkArrayTupleLiteralNodeWithType(n, generic)
	case *ast.HashSetLiteralNode:
		generic, ok := typ.(*types.Generic)
		if !ok || generic.TypeArguments.Len() != 1 || !c.isSubtype(c.StdHashSet(), generic.Namespace, nil) {
			break
		}
		return c.checkHashSetLiteralNodeWithType(n, generic)
	case *ast.HashMapLiteralNode:
		generic, ok := typ.(*types.Generic)
		if !ok || generic.TypeArguments.Len() != 2 || !c.isSubtype(c.StdHashMap(), generic.Namespace, nil) {
			break
		}
		return c.checkHashMapLiteralNodeWithType(n, generic)
	case *ast.HashRecordLiteralNode:
		generic, ok := typ.(*types.Generic)
		if !ok || generic.TypeArguments.Len() != 2 || !c.isSubtype(c.StdHashRecord(), generic.Namespace, nil) {
			break
		}
		return c.checkHashRecordLiteralNodeWithType(n, generic)
	}

	return c.checkExpression(node)
}

func (c *Checker) checkExpression(node ast.ExpressionNode) ast.ExpressionNode {
	if node == nil {
		return nil
	}
	if node.SkipTypechecking() {
		return node
	}

	switch n := node.(type) {
	case *ast.FalseLiteralNode, *ast.TrueLiteralNode, *ast.NilLiteralNode,
		*ast.InterpolatedSymbolLiteralNode, *ast.ConstantDeclarationNode, *ast.UninterpolatedRegexLiteralNode:
		return n
	case *ast.ImplementExpressionNode:
		if c.typeOf(node) == nil {
			c.addFailure(
				"cannot implement interfaces in this context",
				node.Span(),
			)
			n.SetType(types.Nothing{})
		}
		return n
	case *ast.InstanceVariableDeclarationNode:
		if c.typeOf(node) == nil {
			c.addFailure(
				"instance variable definitions cannot appear in this context",
				n.Span(),
			)
			n.SetType(types.Nothing{})
		}
		return n
	case *ast.UsingExpressionNode:
		if c.typeOf(node) == nil {
			c.addFailure(
				"using declarations cannot appear in this context",
				n.Span(),
			)
			n.SetType(types.Nothing{})
			return n
		}
		c.resolveUsingExpression(n)
		return n
	case *ast.TypeDefinitionNode, *ast.GenericTypeDefinitionNode:
		if c.typeOf(node) == nil {
			c.addFailure(
				"type definitions cannot appear in this context",
				n.Span(),
			)
			n.SetType(types.Nothing{})
		}
		return n
	case *ast.StructDeclarationNode:
		c.addFailure(
			"struct definitions cannot appear in this context",
			n.Constant.Span(),
		)
		n.SetType(types.Nothing{})
		return n
	case *ast.MethodDefinitionNode, *ast.InitDefinitionNode,
		*ast.MethodSignatureDefinitionNode, *ast.SetterDeclarationNode,
		*ast.GetterDeclarationNode, *ast.AttrDeclarationNode, *ast.AliasDeclarationNode:
		if c.typeOf(node) == nil {
			c.addFailure(
				"method definitions cannot appear in this context",
				node.Span(),
			)
			node.SetType(types.Nothing{})
		}
		return n
	case *ast.SelfLiteralNode:
		switch c.mode {
		case methodMode, closureInferReturnTypeMode:
			n.SetType(types.Self{})
		default:
			n.SetType(c.selfType)
		}
		return n
	case *ast.IncludeExpressionNode:
		c.checkIncludeExpressionNode(n)
		return n
	case *ast.TypeExpressionNode:
		n.TypeNode = c.checkTypeNode(n.TypeNode)
		return n
	case *ast.IntLiteralNode:
		n.SetType(types.NewIntLiteral(n.Value))
		return n
	case *ast.Int64LiteralNode:
		n.SetType(types.NewInt64Literal(n.Value))
		return n
	case *ast.Int32LiteralNode:
		n.SetType(types.NewInt32Literal(n.Value))
		return n
	case *ast.Int16LiteralNode:
		n.SetType(types.NewInt16Literal(n.Value))
		return n
	case *ast.Int8LiteralNode:
		n.SetType(types.NewInt8Literal(n.Value))
		return n
	case *ast.UInt64LiteralNode:
		n.SetType(types.NewUInt64Literal(n.Value))
		return n
	case *ast.UInt32LiteralNode:
		n.SetType(types.NewUInt32Literal(n.Value))
		return n
	case *ast.UInt16LiteralNode:
		n.SetType(types.NewUInt16Literal(n.Value))
		return n
	case *ast.UInt8LiteralNode:
		n.SetType(types.NewUInt8Literal(n.Value))
		return n
	case *ast.FloatLiteralNode:
		n.SetType(types.NewFloatLiteral(n.Value))
		return n
	case *ast.Float64LiteralNode:
		n.SetType(types.NewFloat64Literal(n.Value))
		return n
	case *ast.Float32LiteralNode:
		n.SetType(types.NewFloat32Literal(n.Value))
		return n
	case *ast.BigFloatLiteralNode:
		n.SetType(types.NewBigFloatLiteral(n.Value))
		return n
	case *ast.DoubleQuotedStringLiteralNode:
		n.SetType(types.NewStringLiteral(n.Value))
		return n
	case *ast.RawStringLiteralNode:
		n.SetType(types.NewStringLiteral(n.Value))
		return n
	case *ast.RawCharLiteralNode:
		n.SetType(types.NewCharLiteral(n.Value))
		return n
	case *ast.CharLiteralNode:
		n.SetType(types.NewCharLiteral(n.Value))
		return n
	case *ast.InterpolatedStringLiteralNode:
		c.checkInterpolatedStringLiteralNode(n)
		return n
	case *ast.InterpolatedRegexLiteralNode:
		c.checkInterpolatedRegexLiteralNode(n)
		return n
	case *ast.SimpleSymbolLiteralNode:
		n.SetType(types.NewSymbolLiteral(n.Content))
		return n
	case *ast.VariableDeclarationNode:
		c.checkVariableDeclarationNode(n)
		return n
	case *ast.VariablePatternDeclarationNode:
		return c.checkVariablePatternDeclarationNode(n)
	case *ast.ValueDeclarationNode:
		c.checkValueDeclarationNode(n)
		return n
	case *ast.PublicIdentifierNode:
		c.checkPublicIdentifierNode(n)
		return n
	case *ast.PrivateIdentifierNode:
		c.checkPrivateIdentifierNode(n)
		return n
	case *ast.InstanceVariableNode:
		c.checkInstanceVariableNode(n)
		return n
	case *ast.PublicConstantNode:
		c.checkPublicConstantNode(n)
		return n
	case *ast.PrivateConstantNode:
		c.checkPrivateConstantNode(n)
		return n
	case *ast.ConstantLookupNode:
		return c.checkConstantLookupNode(n)
	case *ast.ModuleDeclarationNode:
		c.checkExpressionsWithinModule(n)
		return n
	case *ast.ClassDeclarationNode:
		c.checkExpressionsWithinClass(n)
		return n
	case *ast.MixinDeclarationNode:
		c.checkExpressionsWithinMixin(n)
		return n
	case *ast.InterfaceDeclarationNode:
		c.checkExpressionsWithinInterface(n)
		return n
	case *ast.SingletonBlockExpressionNode:
		c.checkExpressionsWithinSingleton(n)
		return n
	case *ast.AssignmentExpressionNode:
		return c.checkAssignmentExpressionNode(n)
	case *ast.ReceiverlessMethodCallNode:
		return c.checkReceiverlessMethodCallNode(n)
	case *ast.GenericReceiverlessMethodCallNode:
		return c.checkGenericReceiverlessMethodCallNode(n)
	case *ast.MethodCallNode:
		return c.checkMethodCallNode(n)
	case *ast.GenericMethodCallNode:
		return c.checkGenericMethodCallNode(n)
	case *ast.CallNode:
		return c.checkCallNode(n)
	case *ast.ClosureLiteralNode:
		return c.checkClosureLiteralNode(n)
	case *ast.NewExpressionNode:
		return c.checkNewExpressionNode(n)
	case *ast.ConstructorCallNode:
		return c.checkConstructorCallNode(n)
	case *ast.GenericConstructorCallNode:
		return c.checkGenericConstructorCallNode(n)
	case *ast.AttributeAccessNode:
		return c.checkAttributeAccessNode(n)
	case *ast.NilSafeSubscriptExpressionNode:
		return c.checkNilSafeSubscriptExpressionNode(n)
	case *ast.SubscriptExpressionNode:
		return c.checkSubscriptExpressionNode(n)
	case *ast.LogicalExpressionNode:
		return c.checkLogicalExpression(n)
	case *ast.BinaryExpressionNode:
		return c.checkBinaryExpression(n)
	case *ast.UnaryExpressionNode:
		return c.checkUnaryExpression(n)
	case *ast.PostfixExpressionNode:
		return c.checkPostfixExpressionNode(n)
	case *ast.DoExpressionNode:
		return c.checkDoExpressionNode(n)
	case *ast.IfExpressionNode:
		return c.checkIfExpressionNode(n)
	case *ast.ModifierIfElseNode:
		return c.checkModifierIfElseNode(n)
	case *ast.ModifierNode:
		return c.checkModifierNode(n)
	case *ast.UnlessExpressionNode:
		return c.checkUnlessExpressionNode(n)
	case *ast.WhileExpressionNode:
		return c.checkWhileExpressionNode("", n)
	case *ast.UntilExpressionNode:
		return c.checkUntilExpressionNode("", n)
	case *ast.LoopExpressionNode:
		return c.checkLoopExpressionNode("", n)
	case *ast.NumericForExpressionNode:
		return c.checkNumericForExpressionNode("", n)
	case *ast.LabeledExpressionNode:
		return c.checkLabeledExpressionNode(n)
	case *ast.ReturnExpressionNode:
		return c.checkReturnExpressionNode(n)
	case *ast.BreakExpressionNode:
		return c.checkBreakExpressionNode(n)
	case *ast.ContinueExpressionNode:
		return c.checkContinueExpressionNode(n)
	case *ast.ArrayListLiteralNode:
		return c.checkArrayListLiteralNode(n)
	case *ast.WordArrayListLiteralNode:
		return c.checkWordArrayListLiteralNode(n)
	case *ast.SymbolArrayListLiteralNode:
		return c.checkSymbolArrayListLiteralNode(n)
	case *ast.HexArrayListLiteralNode:
		return c.checkHexArrayListLiteralNode(n)
	case *ast.BinArrayListLiteralNode:
		return c.checkBinArrayListLiteralNode(n)
	case *ast.ArrayTupleLiteralNode:
		return c.checkArrayTupleLiteralNode(n)
	case *ast.WordArrayTupleLiteralNode:
		return c.checkWordArrayTupleLiteralNode(n)
	case *ast.SymbolArrayTupleLiteralNode:
		return c.checkSymbolArrayTupleLiteralNode(n)
	case *ast.HexArrayTupleLiteralNode:
		return c.checkHexArrayTupleLiteralNode(n)
	case *ast.BinArrayTupleLiteralNode:
		return c.checkBinArrayTupleLiteralNode(n)
	case *ast.HashSetLiteralNode:
		return c.checkHashSetLiteralNode(n)
	case *ast.WordHashSetLiteralNode:
		return c.checkWordHashSetLiteralNode(n)
	case *ast.SymbolHashSetLiteralNode:
		return c.checkSymbolHashSetLiteralNode(n)
	case *ast.HexHashSetLiteralNode:
		return c.checkHexHashSetLiteralNode(n)
	case *ast.BinHashSetLiteralNode:
		return c.checkBinHashSetLiteralNode(n)
	case *ast.HashMapLiteralNode:
		return c.checkHashMapLiteralNode(n)
	case *ast.HashRecordLiteralNode:
		return c.checkHashRecordLiteralNode(n)
	default:
		c.addFailure(
			fmt.Sprintf("invalid expression type %T", node),
			node.Span(),
		)
		return n
	}
}

func (c *Checker) StdInt() *types.Class {
	return c.GlobalEnv.StdSubtypeClass(symbol.Int)
}

func (c *Checker) StdFloat() *types.Class {
	return c.GlobalEnv.StdSubtypeClass(symbol.Float)
}

func (c *Checker) StdBigFloat() *types.Class {
	return c.GlobalEnv.StdSubtypeClass(symbol.BigFloat)
}

func (c *Checker) StdClass() *types.Class {
	return c.GlobalEnv.StdSubtypeClass(symbol.Class)
}

func (c *Checker) Std(name value.Symbol) types.Type {
	return c.GlobalEnv.StdSubtype(name)
}

func (c *Checker) StdString() types.Type {
	return c.GlobalEnv.StdSubtype(symbol.String)
}

func (c *Checker) StdStringConvertible() types.Type {
	return c.GlobalEnv.StdSubtype(symbol.StringConvertible)
}

func (c *Checker) StdInspectable() types.Type {
	return c.GlobalEnv.StdSubtype(symbol.Inspectable)
}

func (c *Checker) StdAnyInt() types.Type {
	return c.GlobalEnv.StdSubtype(symbol.AnyInt)
}

func (c *Checker) StdBool() *types.Class {
	return c.GlobalEnv.StdSubtypeClass(symbol.Bool)
}

func (c *Checker) StdArrayList() *types.Class {
	return c.GlobalEnv.StdSubtypeClass(symbol.ArrayList)
}

func (c *Checker) StdArrayTuple() *types.Class {
	return c.GlobalEnv.StdSubtypeClass(symbol.ArrayTuple)
}

func (c *Checker) StdHashSet() *types.Class {
	return c.GlobalEnv.StdSubtypeClass(symbol.HashSet)
}

func (c *Checker) StdHashMap() *types.Class {
	return c.GlobalEnv.StdSubtypeClass(symbol.HashMap)
}

func (c *Checker) StdHashRecord() *types.Class {
	return c.GlobalEnv.StdSubtypeClass(symbol.HashRecord)
}

func (c *Checker) StdNil() *types.Class {
	return c.GlobalEnv.StdSubtypeClass(symbol.Nil)
}

func (c *Checker) StdTrue() *types.Class {
	return c.GlobalEnv.StdSubtypeClass(symbol.True)
}

func (c *Checker) StdFalse() *types.Class {
	return c.GlobalEnv.StdSubtypeClass(symbol.False)
}

func (c *Checker) checkArithmeticBinaryOperator(
	left,
	right ast.ExpressionNode,
	methodName value.Symbol,
	span *position.Span,
) types.Type {
	leftType := c.toNonLiteral(c.typeOf(left), true)
	leftClassType, leftIsClass := leftType.(*types.Class)

	rightType := c.toNonLiteral(c.typeOf(right), true)
	rightClassType, rightIsClass := rightType.(*types.Class)
	if !leftIsClass || !rightIsClass {
		return c.checkBinaryOpMethodCall(
			left,
			right,
			methodName,
			span,
		)
	}

	switch leftClassType.Name() {
	case "Std::Int":
		switch rightClassType.Name() {
		case "Std::Int":
			return c.StdInt()
		case "Std::Float":
			return c.StdFloat()
		case "Std::BigFloat":
			return c.StdBigFloat()
		}
	case "Std::Float":
		switch rightClassType.Name() {
		case "Std::Int":
			return c.StdFloat()
		case "Std::Float":
			return c.StdFloat()
		case "Std::BigFloat":
			return c.StdBigFloat()
		}
	}

	return c.checkBinaryOpMethodCall(
		left,
		right,
		methodName,
		span,
	)
}

func (c *Checker) checkUnaryExpression(node *ast.UnaryExpressionNode) ast.ExpressionNode {
	switch node.Op.Type {
	case token.PLUS:
		receiver, _, typ := c.checkSimpleMethodCall(
			node.Right,
			token.DOT,
			symbol.OpUnaryPlus,
			nil,
			nil,
			nil,
			node.Span(),
		)
		node.Right = receiver
		node.SetType(typ)
		return node
	case token.MINUS:
		receiver, _, typ := c.checkSimpleMethodCall(
			node.Right,
			token.DOT,
			symbol.OpNegate,
			nil,
			nil,
			nil,
			node.Span(),
		)
		node.Right = receiver
		node.SetType(typ)
		return node
	case token.TILDE:
		receiver, _, typ := c.checkSimpleMethodCall(
			node.Right,
			token.DOT,
			value.ToSymbol(node.Op.StringValue()),
			nil,
			nil,
			nil,
			node.Span(),
		)
		node.Right = receiver
		node.SetType(typ)
		return node
	case token.BANG:
		return c.checkNotOperator(node)
	default:
		c.addFailure(
			fmt.Sprintf("invalid unary operator %s", node.Op.String()),
			node.Span(),
		)
		return node
	}
}

func (c *Checker) checkPostfixExpression(node *ast.PostfixExpressionNode, methodName string) ast.ExpressionNode {
	return c.checkExpression(
		ast.NewAssignmentExpressionNode(
			node.Span(),
			token.New(node.Op.Span(), token.EQUAL_OP),
			node.Expression,
			ast.NewMethodCallNode(
				node.Span(),
				node.Expression,
				token.New(node.Op.Span(), token.DOT),
				methodName,
				nil,
				nil,
			),
		),
	)
}

func (c *Checker) checkModifierInRecord(node *ast.ModifierNode) (keyType, valueType types.Type) {
	switch node.Modifier.Type {
	case token.IF:
		return c.checkRecordIfModifier(node)
	case token.UNLESS:
		return c.checkRecordUnlessModifier(node)
	default:
		panic(fmt.Sprintf("invalid collection modifier: %#v", node.Modifier))
	}
}

func (c *Checker) checkRecordIfElseModifier(node *ast.ModifierIfElseNode) (keyType, valueType types.Type) {
	c.pushNestedLocalEnv()
	node.Condition = c.checkExpression(node.Condition)
	conditionType := c.typeOfGuardVoid(node.Condition)

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Condition, assumptionTruthy)
	var thenKeyType, thenValueType types.Type
	switch l := node.ThenExpression.(type) {
	case *ast.KeyValueExpressionNode:
		l.Key = c.checkExpression(l.Key)
		thenKeyType = c.typeOf(l.Key)

		l.Value = c.checkExpression(l.Value)
		thenValueType = c.typeOf(l.Value)
	case *ast.SymbolKeyValueExpressionNode:
		thenKeyType = c.Std(symbol.Symbol)

		l.Value = c.checkExpression(l.Value)
		thenValueType = c.typeOf(l.Value)
	default:
		panic(fmt.Sprintf("invalid record element node: %#v", node.ThenExpression))
	}
	c.popLocalEnv()

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Condition, assumptionFalsy)
	var elseKeyType, elseValueType types.Type
	switch l := node.ElseExpression.(type) {
	case *ast.KeyValueExpressionNode:
		l.Key = c.checkExpression(l.Key)
		elseKeyType = c.typeOf(l.Key)

		l.Value = c.checkExpression(l.Value)
		elseValueType = c.typeOf(l.Value)
	case *ast.SymbolKeyValueExpressionNode:
		elseKeyType = c.Std(symbol.Symbol)

		l.Value = c.checkExpression(l.Value)
		elseValueType = c.typeOf(l.Value)
	default:
		panic(fmt.Sprintf("invalid record element node: %#v", node.ThenExpression))
	}
	c.popLocalEnv()

	c.popLocalEnv()

	if c.isTruthy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is truthy",
				types.InspectWithColor(conditionType),
			),
			node.Condition.Span(),
		)
		c.addUnreachableCodeError(node.ElseExpression.Span())
		return thenKeyType, thenValueType
	}
	if c.isFalsy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is falsy",
				types.InspectWithColor(conditionType),
			),
			node.Condition.Span(),
		)
		c.addUnreachableCodeError(node.ThenExpression.Span())
		return elseKeyType, elseValueType
	}

	return c.newNormalisedUnion(thenKeyType, elseKeyType), c.newNormalisedUnion(thenValueType, elseValueType)
}

func (c *Checker) checkMapPairs(pairs []ast.ExpressionNode) (keyTypes []types.Type, valueTypes []types.Type) {
	for _, pairNode := range pairs {
		switch p := pairNode.(type) {
		case *ast.KeyValueExpressionNode:
			p.Key = c.checkExpression(p.Key)
			keyTypes = append(keyTypes, c.toNonLiteral(c.typeOfGuardVoid(p.Key), false))

			p.Value = c.checkExpression(p.Value)
			valueTypes = append(valueTypes, c.toNonLiteral(c.typeOfGuardVoid(p.Value), false))
		case *ast.SymbolKeyValueExpressionNode:
			keyTypes = append(keyTypes, c.Std(symbol.Symbol))

			p.Value = c.checkExpression(p.Value)
			valueTypes = append(valueTypes, c.toNonLiteral(c.typeOfGuardVoid(p.Value), false))
		case *ast.ModifierNode:
			keyType, valueType := c.checkModifierInRecord(p)
			keyTypes = append(keyTypes, c.toNonLiteral(keyType, false))

			valueTypes = append(valueTypes, c.toNonLiteral(valueType, false))
		case *ast.ModifierIfElseNode:
			keyType, valueType := c.checkRecordIfElseModifier(p)

			keyTypes = append(keyTypes, c.toNonLiteral(keyType, false))
			valueTypes = append(valueTypes, c.toNonLiteral(valueType, false))
		default:
			panic(fmt.Sprintf("invalid map element node: %#v", pairNode))
		}
	}

	return keyTypes, valueTypes
}

func (c *Checker) checkRecordPairs(pairs []ast.ExpressionNode) (keyTypes []types.Type, valueTypes []types.Type) {
	for _, pairNode := range pairs {
		switch p := pairNode.(type) {
		case *ast.KeyValueExpressionNode:
			p.Key = c.checkExpression(p.Key)
			keyTypes = append(keyTypes, c.typeOfGuardVoid(p.Key))

			p.Value = c.checkExpression(p.Value)
			valueTypes = append(valueTypes, c.typeOfGuardVoid(p.Value))
		case *ast.SymbolKeyValueExpressionNode:
			keyTypes = append(keyTypes, types.NewSymbolLiteral(p.Key))

			p.Value = c.checkExpression(p.Value)
			valueTypes = append(valueTypes, c.typeOfGuardVoid(p.Value))
		case *ast.ModifierNode:
			keyType, valueType := c.checkModifierInRecord(p)

			keyTypes = append(keyTypes, keyType)
			valueTypes = append(valueTypes, valueType)
		case *ast.ModifierIfElseNode:
			keyType, valueType := c.checkRecordIfElseModifier(p)

			keyTypes = append(keyTypes, keyType)
			valueTypes = append(valueTypes, valueType)
		default:
			panic(fmt.Sprintf("invalid map element node: %#v", pairNode))
		}
	}

	return keyTypes, valueTypes
}

func (c *Checker) checkHashMapLiteralNode(node *ast.HashMapLiteralNode) ast.ExpressionNode {
	return c.checkHashMapLiteralNodeWithType(node, nil)
}

func (c *Checker) checkHashMapLiteralNodeWithType(node *ast.HashMapLiteralNode, typ *types.Generic) ast.ExpressionNode {
	keyTypes, valueTypes := c.checkMapPairs(node.Elements)

	keyType := c.newNormalisedUnion(keyTypes...)
	valueType := c.newNormalisedUnion(valueTypes...)
	if typ != nil {
		c.checkCanAssign(keyType, typ.TypeArguments.Get(0).Type, node.Span())
		c.checkCanAssign(valueType, typ.TypeArguments.Get(1).Type, node.Span())
		node.SetType(typ)
	} else if len(keyTypes) == 0 {
		generic := types.NewGenericWithTypeArgs(c.StdHashMap(), types.Any{}, types.Any{})
		node.SetType(generic)
	} else {
		generic := types.NewGenericWithTypeArgs(c.StdHashMap(), keyType, valueType)
		node.SetType(generic)
	}

	if node.Capacity != nil {
		node.Capacity = c.checkExpression(node.Capacity)
		capacityType := c.typeOf(node.Capacity)
		if !c.isSubtype(capacityType, c.StdAnyInt(), nil) {
			c.addFailure(
				fmt.Sprintf(
					"capacity must be an integer, got `%s`",
					types.InspectWithColor(capacityType),
				),
				node.Span(),
			)
		}
	}

	return node
}

func (c *Checker) checkHashRecordLiteralNode(node *ast.HashRecordLiteralNode) ast.ExpressionNode {
	return c.checkHashRecordLiteralNodeWithType(node, nil)
}

func (c *Checker) checkHashRecordLiteralNodeWithType(node *ast.HashRecordLiteralNode, typ *types.Generic) ast.ExpressionNode {
	keyTypes, valueTypes := c.checkRecordPairs(node.Elements)

	keyType := c.newNormalisedUnion(keyTypes...)
	valueType := c.newNormalisedUnion(valueTypes...)
	if typ != nil {
		c.checkCanAssign(keyType, typ.TypeArguments.Get(0).Type, node.Span())
		c.checkCanAssign(valueType, typ.TypeArguments.Get(1).Type, node.Span())
		node.SetType(typ)
	} else if len(keyTypes) == 0 {
		generic := types.NewGenericWithTypeArgs(c.StdHashRecord(), types.Any{}, types.Any{})
		node.SetType(generic)
	} else {
		generic := types.NewGenericWithTypeArgs(c.StdHashRecord(), keyType, valueType)
		node.SetType(generic)
	}

	return node
}

func (c *Checker) checkModifierInCollection(node *ast.ModifierNode) ast.ExpressionNode {
	switch node.Modifier.Type {
	case token.IF:
		return c.checkCollectionIfModifier(node)
	case token.UNLESS:
		return c.checkCollectionUnlessModifier(node)
	default:
		panic(fmt.Sprintf("invalid collection modifier: %#v", node.Modifier))
	}
}

func (c *Checker) checkMutableCollectionElements(elements []ast.ExpressionNode) []types.Type {
	var elementTypes []types.Type
	for i, elementNode := range elements {
		switch e := elementNode.(type) {
		case *ast.ModifierNode:
			elementNode := c.checkModifierInCollection(e)
			elements[i] = elementNode
			elementTypes = append(elementTypes, c.toNonLiteral(c.typeOfGuardVoid(elementNode), false))
		case *ast.ModifierIfElseNode:
			elementNode := c.checkCollectionIfElseModifier(e)
			elements[i] = elementNode
			elementTypes = append(elementTypes, c.toNonLiteral(c.typeOfGuardVoid(elementNode), false))
		default:
			elementNode := c.checkExpression(elementNode)
			elements[i] = elementNode
			elementTypes = append(elementTypes, c.toNonLiteral(c.typeOfGuardVoid(elementNode), false))
		}
	}

	return elementTypes
}

func (c *Checker) checkImmutableCollectionElements(elements []ast.ExpressionNode) []types.Type {
	var elementTypes []types.Type
	for i, elementNode := range elements {
		switch e := elementNode.(type) {
		case *ast.ModifierNode:
			elementNode := c.checkModifierInCollection(e)
			elements[i] = elementNode
			elementTypes = append(elementTypes, c.typeOfGuardVoid(elementNode))
		case *ast.ModifierIfElseNode:
			elementNode := c.checkCollectionIfElseModifier(e)
			elements[i] = elementNode
			elementTypes = append(elementTypes, c.typeOfGuardVoid(elementNode))
		default:
			elementNode := c.checkExpression(elementNode)
			elements[i] = elementNode
			elementTypes = append(elementTypes, c.typeOfGuardVoid(elementNode))
		}
	}

	return elementTypes
}

func (c *Checker) checkArrayListLiteralNode(node *ast.ArrayListLiteralNode) ast.ExpressionNode {
	return c.checkArrayListLiteralNodeWithType(node, nil)
}

func (c *Checker) checkArrayListLiteralNodeWithType(node *ast.ArrayListLiteralNode, typ *types.Generic) ast.ExpressionNode {
	elementTypes := c.checkMutableCollectionElements(node.Elements)

	elementType := c.newNormalisedUnion(elementTypes...)
	if typ != nil {
		c.checkCanAssign(elementType, typ.TypeArguments.Get(0).Type, node.Span())
		node.SetType(typ)
	} else if len(elementTypes) == 0 {
		node.SetType(types.NewGenericWithTypeArgs(c.StdArrayList(), types.Any{}))
	} else {
		node.SetType(types.NewGenericWithTypeArgs(c.StdArrayList(), elementType))
	}

	if node.Capacity != nil {
		node.Capacity = c.checkExpression(node.Capacity)
		capacityType := c.typeOf(node.Capacity)
		if !c.isSubtype(capacityType, c.StdAnyInt(), nil) {
			c.addFailure(
				fmt.Sprintf(
					"capacity must be an integer, got `%s`",
					types.InspectWithColor(capacityType),
				),
				node.Span(),
			)
		}
	}

	return node
}

func checkSpecialCollectionLiteralNode[E ast.ExpressionNode](c *Checker, collectionType types.Namespace, elementType types.Type, elements []E, capacity ast.ExpressionNode) types.Type {
	for _, elementNode := range elements {
		c.checkExpression(elementNode)
	}

	generic := types.NewGenericWithTypeArgs(collectionType, elementType)

	if capacity != nil {
		capacity = c.checkExpression(capacity)
		capacityType := c.typeOf(capacity)
		if !c.isSubtype(capacityType, c.StdAnyInt(), nil) {
			c.addFailure(
				fmt.Sprintf(
					"capacity must be an integer, got `%s`",
					types.InspectWithColor(capacityType),
				),
				capacity.Span(),
			)
		}
	}

	return generic
}

func (c *Checker) checkBinArrayListLiteralNode(node *ast.BinArrayListLiteralNode) ast.ExpressionNode {
	typ := checkSpecialCollectionLiteralNode(
		c,
		c.StdArrayList(),
		c.Std(symbol.Int),
		node.Elements,
		node.Capacity,
	)
	node.SetType(typ)

	return node
}

func (c *Checker) checkHexArrayListLiteralNode(node *ast.HexArrayListLiteralNode) ast.ExpressionNode {
	typ := checkSpecialCollectionLiteralNode(
		c,
		c.StdArrayList(),
		c.Std(symbol.Int),
		node.Elements,
		node.Capacity,
	)
	node.SetType(typ)

	return node
}

func (c *Checker) checkSymbolArrayListLiteralNode(node *ast.SymbolArrayListLiteralNode) ast.ExpressionNode {
	typ := checkSpecialCollectionLiteralNode(
		c,
		c.StdArrayList(),
		c.Std(symbol.Symbol),
		node.Elements,
		node.Capacity,
	)
	node.SetType(typ)

	return node
}

func (c *Checker) checkWordArrayListLiteralNode(node *ast.WordArrayListLiteralNode) ast.ExpressionNode {
	typ := checkSpecialCollectionLiteralNode(
		c,
		c.StdArrayList(),
		c.Std(symbol.String),
		node.Elements,
		node.Capacity,
	)
	node.SetType(typ)

	return node
}

func (c *Checker) checkBinArrayTupleLiteralNode(node *ast.BinArrayTupleLiteralNode) ast.ExpressionNode {
	typ := checkSpecialCollectionLiteralNode(
		c,
		c.StdArrayTuple(),
		c.Std(symbol.Int),
		node.Elements,
		nil,
	)
	node.SetType(typ)

	return node
}

func (c *Checker) checkHexArrayTupleLiteralNode(node *ast.HexArrayTupleLiteralNode) ast.ExpressionNode {
	typ := checkSpecialCollectionLiteralNode(
		c,
		c.StdArrayTuple(),
		c.Std(symbol.Int),
		node.Elements,
		nil,
	)
	node.SetType(typ)

	return node
}

func (c *Checker) checkSymbolArrayTupleLiteralNode(node *ast.SymbolArrayTupleLiteralNode) ast.ExpressionNode {
	typ := checkSpecialCollectionLiteralNode(
		c,
		c.StdArrayTuple(),
		c.Std(symbol.Symbol),
		node.Elements,
		nil,
	)
	node.SetType(typ)

	return node
}

func (c *Checker) checkWordArrayTupleLiteralNode(node *ast.WordArrayTupleLiteralNode) ast.ExpressionNode {
	typ := checkSpecialCollectionLiteralNode(
		c,
		c.StdArrayTuple(),
		c.Std(symbol.String),
		node.Elements,
		nil,
	)
	node.SetType(typ)

	return node
}

func (c *Checker) checkBinHashSetLiteralNode(node *ast.BinHashSetLiteralNode) ast.ExpressionNode {
	typ := checkSpecialCollectionLiteralNode(
		c,
		c.StdHashSet(),
		c.Std(symbol.Int),
		node.Elements,
		node.Capacity,
	)
	node.SetType(typ)

	return node
}

func (c *Checker) checkHexHashSetLiteralNode(node *ast.HexHashSetLiteralNode) ast.ExpressionNode {
	typ := checkSpecialCollectionLiteralNode(
		c,
		c.StdHashSet(),
		c.Std(symbol.Int),
		node.Elements,
		node.Capacity,
	)
	node.SetType(typ)

	return node
}

func (c *Checker) checkSymbolHashSetLiteralNode(node *ast.SymbolHashSetLiteralNode) ast.ExpressionNode {
	typ := checkSpecialCollectionLiteralNode(
		c,
		c.StdHashSet(),
		c.Std(symbol.Symbol),
		node.Elements,
		node.Capacity,
	)
	node.SetType(typ)

	return node
}

func (c *Checker) checkWordHashSetLiteralNode(node *ast.WordHashSetLiteralNode) ast.ExpressionNode {
	typ := checkSpecialCollectionLiteralNode(
		c,
		c.StdHashSet(),
		c.Std(symbol.String),
		node.Elements,
		node.Capacity,
	)
	node.SetType(typ)

	return node
}

func (c *Checker) checkHashSetLiteralNode(node *ast.HashSetLiteralNode) ast.ExpressionNode {
	return c.checkHashSetLiteralNodeWithType(node, nil)
}

func (c *Checker) checkHashSetLiteralNodeWithType(node *ast.HashSetLiteralNode, typ *types.Generic) ast.ExpressionNode {
	elementTypes := c.checkMutableCollectionElements(node.Elements)

	elementType := c.newNormalisedUnion(elementTypes...)
	if typ != nil {
		c.checkCanAssign(elementType, typ.TypeArguments.Get(0).Type, node.Span())
		node.SetType(typ)
	} else if len(elementTypes) == 0 {
		node.SetType(types.NewGenericWithTypeArgs(c.StdHashSet(), types.Any{}))
	} else {
		node.SetType(types.NewGenericWithTypeArgs(c.StdHashSet(), elementType))
	}

	if node.Capacity != nil {
		node.Capacity = c.checkExpression(node.Capacity)
		capacityType := c.typeOf(node.Capacity)
		if !c.isSubtype(capacityType, c.StdAnyInt(), nil) {
			c.addFailure(
				fmt.Sprintf(
					"capacity must be an integer, got `%s`",
					types.InspectWithColor(capacityType),
				),
				node.Span(),
			)
		}
	}

	return node
}

func (c *Checker) checkArrayTupleLiteralNode(node *ast.ArrayTupleLiteralNode) ast.ExpressionNode {
	return c.checkArrayTupleLiteralNodeWithType(node, nil)
}

func (c *Checker) checkArrayTupleLiteralNodeWithType(node *ast.ArrayTupleLiteralNode, typ *types.Generic) ast.ExpressionNode {
	elementTypes := c.checkImmutableCollectionElements(node.Elements)

	elementType := c.newNormalisedUnion(elementTypes...)
	if typ != nil {
		c.checkCanAssign(elementType, typ.TypeArguments.Get(0).Type, node.Span())
		node.SetType(typ)
	} else if len(elementTypes) == 0 {
		node.SetType(types.NewGenericWithTypeArgs(c.StdArrayTuple(), types.Any{}))
	} else {
		node.SetType(types.NewGenericWithTypeArgs(c.StdArrayTuple(), elementType))
	}

	return node
}

func (c *Checker) checkContinueExpressionNode(node *ast.ContinueExpressionNode) ast.ExpressionNode {
	var typ types.Type
	if node.Value == nil {
		typ = types.Nil{}
	} else {
		node.Value = c.checkExpression(node.Value)
		typ = c.typeOfGuardVoid(node.Value)
	}

	loop := c.findLoop(node.Label, node.Span())
	if loop != nil && !loop.endless {
		if loop.returnType == nil {
			loop.returnType = typ
		} else {
			loop.returnType = c.newNormalisedUnion(loop.returnType, typ)
		}
	}

	return node
}

func (c *Checker) checkBreakExpressionNode(node *ast.BreakExpressionNode) ast.ExpressionNode {
	var typ types.Type
	if node.Value == nil {
		typ = types.Nil{}
	} else {
		node.Value = c.checkExpression(node.Value)
		typ = c.typeOfGuardVoid(node.Value)
	}

	loop := c.findLoop(node.Label, node.Span())
	if loop != nil {
		if loop.returnType == nil {
			loop.returnType = typ
		} else {
			loop.returnType = c.newNormalisedUnion(loop.returnType, typ)
		}
	}

	return node
}

func (c *Checker) checkReturnExpressionNode(node *ast.ReturnExpressionNode) ast.ExpressionNode {
	if c.mode == closureInferReturnTypeMode {
		var typ types.Type
		if node.Value == nil {
			typ = types.Nil{}
		} else {
			node.Value = c.checkExpression(node.Value)
			typ = c.typeOfGuardVoid(node.Value)
		}

		c.addToReturnType(typ)
		return node
	}

	var typ types.Type
	if node.Value == nil {
		typ = types.Nil{}
	} else {
		if types.IsVoid(c.returnType) {
			c.addWarning(
				"values returned in void context will be ignored",
				node.Value.Span(),
			)
		}
		node.Value = c.checkExpression(node.Value)
		typ = c.typeOfGuardVoid(node.Value)
	}
	c.checkCanAssign(typ, c.returnType, node.Span())

	return node
}

func (c *Checker) checkLabeledExpressionNode(node *ast.LabeledExpressionNode) ast.ExpressionNode {
	switch expr := node.Expression.(type) {
	case *ast.LoopExpressionNode:
		node.Expression = c.checkLoopExpressionNode(node.Label, expr)
	case *ast.WhileExpressionNode:
		node.Expression = c.checkWhileExpressionNode(node.Label, expr)
	case *ast.UntilExpressionNode:
		node.Expression = c.checkUntilExpressionNode(node.Label, expr)
	case *ast.NumericForExpressionNode:
		node.Expression = c.checkNumericForExpressionNode(node.Label, expr)
	case *ast.ModifierNode:
		switch expr.Modifier.Type {
		case token.WHILE:
			node.Expression = c.checkWhileModifierNode(node.Label, expr)
		case token.UNTIL:
			node.Expression = c.checkUntilModifierNode(node.Label, expr)
		default:
			node.Expression = c.checkExpression(node.Expression)
		}
	// case *ast.ForInExpressionNode:
	// 	c.forInExpression(node.Label, expr)
	// case *ast.ModifierForInNode:
	// 	c.modifierForIn(node.Label, expr)
	default:
		node.Expression = c.checkExpression(node.Expression)
	}
	return node
}

func (c *Checker) checkNumericForExpressionNode(label string, node *ast.NumericForExpressionNode) ast.ExpressionNode {
	c.pushNestedLocalEnv()
	node.Initialiser = c.checkExpression(node.Initialiser)

	node.Condition = c.checkExpression(node.Condition)
	conditionType := c.typeOfGuardVoid(node.Condition)

	var typ types.Type
	var endless bool
	var noop bool
	if node.Condition == nil {
		endless = true
		typ = types.Never{}
	}
	if c.isTruthy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is truthy",
				types.InspectWithColor(conditionType),
			),
			node.Condition.Span(),
		)
		endless = true
		typ = types.Never{}
	}
	if c.isFalsy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this loop will never execute since type `%s` is falsy",
				types.InspectWithColor(conditionType),
			),
			node.Condition.Span(),
		)
		if len(node.ThenBody) > 0 {
			c.addUnreachableCodeError(node.ThenBody[0].Span())
		}
		noop = true
		typ = types.Nil{}
	}
	loop := c.registerLoop(label, endless)

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Condition, assumptionTruthy)

	node.Increment = c.checkExpression(node.Increment)

	thenType, _ := c.checkStatements(node.ThenBody)
	c.popLocalEnv()

	c.popLocalEnv()

	if noop {
		node.SetType(typ)
		return node
	}
	if typ == nil {
		typ = c.toNilable(thenType)
	}
	if loop.returnType != nil {
		typ = c.newNormalisedUnion(typ, loop.returnType)
	}
	node.SetType(typ)
	return node
}

func (c *Checker) checkLoopExpressionNode(label string, node *ast.LoopExpressionNode) ast.ExpressionNode {
	loop := c.registerLoop(label, true)
	c.pushNestedLocalEnv()
	c.checkStatements(node.ThenBody)
	c.popLocalEnv()
	c.popLoop()

	if loop.returnType == nil {
		node.SetType(types.Never{})
	} else {
		node.SetType(loop.returnType)
	}
	return node
}

func (c *Checker) checkUntilExpressionNode(label string, node *ast.UntilExpressionNode) ast.ExpressionNode {
	c.pushNestedLocalEnv()
	node.Condition = c.checkExpression(node.Condition)
	conditionType := c.typeOfGuardVoid(node.Condition)

	var typ types.Type
	var endless bool
	var noop bool
	if c.isTruthy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this loop will never execute since type `%s` is truthy",
				types.InspectWithColor(conditionType),
			),
			node.Condition.Span(),
		)
		if len(node.ThenBody) > 0 {
			c.addUnreachableCodeError(node.ThenBody[0].Span())
		}
		noop = true
		typ = types.Nil{}
	}
	if c.isFalsy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is falsy",
				types.InspectWithColor(conditionType),
			),
			node.Condition.Span(),
		)
		endless = true
		typ = types.Never{}
	}
	loop := c.registerLoop(label, endless)

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Condition, assumptionFalsy)
	thenType, _ := c.checkStatements(node.ThenBody)
	c.popLocalEnv()

	c.popLocalEnv()

	if noop {
		node.SetType(typ)
		return node
	}
	if typ == nil {
		typ = c.toNilable(thenType)
	}
	if loop.returnType != nil {
		typ = c.newNormalisedUnion(typ, loop.returnType)
	}
	node.SetType(typ)
	return node
}

func (c *Checker) checkUntilModifierNode(label string, node *ast.ModifierNode) ast.ExpressionNode {
	c.pushNestedLocalEnv()
	node.Right = c.checkExpression(node.Right)
	conditionType := c.typeOfGuardVoid(node.Right)

	var typ types.Type
	var endless bool
	if c.isTruthy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is truthy",
				types.InspectWithColor(conditionType),
			),
			node.Right.Span(),
		)
	}
	if c.isFalsy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is falsy",
				types.InspectWithColor(conditionType),
			),
			node.Right.Span(),
		)
		endless = true
		typ = types.Never{}
	}
	loop := c.registerLoop(label, endless)

	node.Left = c.checkExpression(node.Left)
	thenType := c.typeOf(node.Left)

	c.popLocalEnv()

	if typ == nil {
		typ = c.toNilable(thenType)
	}
	if loop.returnType != nil {
		typ = c.newNormalisedUnion(typ, loop.returnType)
	}
	node.SetType(typ)
	return node
}

func (c *Checker) checkWhileExpressionNode(label string, node *ast.WhileExpressionNode) ast.ExpressionNode {
	c.pushNestedLocalEnv()
	node.Condition = c.checkExpression(node.Condition)
	conditionType := c.typeOfGuardVoid(node.Condition)

	var typ types.Type
	var endless bool
	var noop bool
	if c.isTruthy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is truthy",
				types.InspectWithColor(conditionType),
			),
			node.Condition.Span(),
		)
		endless = true
		typ = types.Never{}
	}
	if c.isFalsy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this loop will never execute since type `%s` is falsy",
				types.InspectWithColor(conditionType),
			),
			node.Condition.Span(),
		)
		if len(node.ThenBody) > 0 {
			c.addUnreachableCodeError(node.ThenBody[0].Span())
		}
		noop = true
		typ = types.Nil{}
	}
	loop := c.registerLoop(label, endless)

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Condition, assumptionTruthy)
	thenType, _ := c.checkStatements(node.ThenBody)
	c.popLocalEnv()

	c.popLocalEnv()

	if noop {
		node.SetType(typ)
		return node
	}
	if typ == nil {
		typ = c.toNilable(thenType)
	}
	if loop.returnType != nil {
		typ = c.newNormalisedUnion(typ, loop.returnType)
	}
	node.SetType(typ)
	return node
}

func (c *Checker) checkWhileModifierNode(label string, node *ast.ModifierNode) ast.ExpressionNode {
	c.pushNestedLocalEnv()
	node.Right = c.checkExpression(node.Right)
	conditionType := c.typeOfGuardVoid(node.Right)

	var typ types.Type
	var endless bool
	if c.isTruthy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is truthy",
				types.InspectWithColor(conditionType),
			),
			node.Right.Span(),
		)
		endless = true
		typ = types.Never{}
	}
	if c.isFalsy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is falsy",
				types.InspectWithColor(conditionType),
			),
			node.Right.Span(),
		)
	}
	loop := c.registerLoop(label, endless)

	node.Left = c.checkExpression(node.Left)
	thenType := c.typeOf(node.Left)

	c.popLocalEnv()

	if typ == nil {
		typ = c.toNilable(thenType)
	}
	if loop.returnType != nil {
		typ = c.newNormalisedUnion(typ, loop.returnType)
	}
	node.SetType(typ)
	return node
}

func (c *Checker) checkUnlessExpressionNode(node *ast.UnlessExpressionNode) ast.ExpressionNode {
	c.pushNestedLocalEnv()
	node.Condition = c.checkExpression(node.Condition)
	conditionType := c.typeOfGuardVoid(node.Condition)

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Condition, assumptionFalsy)
	thenType, _ := c.checkStatements(node.ThenBody)
	c.popLocalEnv()

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Condition, assumptionTruthy)
	elseType, _ := c.checkStatements(node.ElseBody)
	c.popLocalEnv()

	c.popLocalEnv()

	if c.isTruthy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is truthy",
				types.InspectWithColor(conditionType),
			),
			node.Condition.Span(),
		)
		if len(node.ThenBody) > 0 {
			c.addUnreachableCodeError(node.ThenBody[0].Span())
		}
		node.SetType(elseType)
		return node
	}
	if c.isFalsy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is falsy",
				types.InspectWithColor(conditionType),
			),
			node.Condition.Span(),
		)
		if len(node.ElseBody) > 0 {
			c.addUnreachableCodeError(node.ElseBody[0].Span())
		}
		node.SetType(thenType)
		return node
	}

	node.SetType(c.newNormalisedUnion(thenType, elseType))
	return node
}

func (c *Checker) checkModifierNode(node *ast.ModifierNode) ast.ExpressionNode {
	switch node.Modifier.Type {
	case token.IF:
		return c.checkIfExpressionNode(
			ast.NewIfExpressionNode(
				node.Span(),
				node.Right,
				ast.ExpressionToStatements(node.Left),
				nil,
			),
		)
	case token.UNLESS:
		return c.checkUnlessExpressionNode(
			ast.NewUnlessExpressionNode(
				node.Span(),
				node.Right,
				ast.ExpressionToStatements(node.Left),
				nil,
			),
		)
	case token.WHILE:
		return c.checkWhileModifierNode("", node)
	case token.UNTIL:
		return c.checkUntilModifierNode("", node)
	default:
		c.addFailure(
			fmt.Sprintf("illegal modifier: %s", node.Modifier.StringValue()),
			node.Span(),
		)
		return node
	}
}

func (c *Checker) checkModifierIfElseNode(node *ast.ModifierIfElseNode) ast.ExpressionNode {
	return c.checkIfExpressionNode(
		ast.NewIfExpressionNode(
			node.Span(),
			node.Condition,
			ast.ExpressionToStatements(node.ThenExpression),
			ast.ExpressionToStatements(node.ElseExpression),
		),
	)
}

func (c *Checker) checkIfExpressionNode(node *ast.IfExpressionNode) ast.ExpressionNode {
	c.pushNestedLocalEnv()
	node.Condition = c.checkExpression(node.Condition)
	conditionType := c.typeOfGuardVoid(node.Condition)

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Condition, assumptionTruthy)
	thenType, _ := c.checkStatements(node.ThenBody)
	c.popLocalEnv()

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Condition, assumptionFalsy)
	elseType, _ := c.checkStatements(node.ElseBody)
	c.popLocalEnv()

	c.popLocalEnv()

	if c.isTruthy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is truthy",
				types.InspectWithColor(conditionType),
			),
			node.Condition.Span(),
		)
		if len(node.ElseBody) > 0 {
			c.addUnreachableCodeError(node.ElseBody[0].Span())
		}
		node.SetType(thenType)
		return node
	}
	if c.isFalsy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is falsy",
				types.InspectWithColor(conditionType),
			),
			node.Condition.Span(),
		)
		if len(node.ThenBody) > 0 {
			c.addUnreachableCodeError(node.ThenBody[0].Span())
		}
		node.SetType(elseType)
		return node
	}

	node.SetType(c.newNormalisedUnion(thenType, elseType))
	return node
}

func (c *Checker) checkCollectionIfElseModifier(node *ast.ModifierIfElseNode) ast.ExpressionNode {
	c.pushNestedLocalEnv()
	node.Condition = c.checkExpression(node.Condition)
	conditionType := c.typeOfGuardVoid(node.Condition)

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Condition, assumptionTruthy)
	node.ThenExpression = c.checkExpression(node.ThenExpression)
	thenType := c.typeOf(node.ThenExpression)
	c.popLocalEnv()

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Condition, assumptionFalsy)
	node.ElseExpression = c.checkExpression(node.ElseExpression)
	elseType := c.typeOf(node.ElseExpression)
	c.popLocalEnv()

	c.popLocalEnv()

	if c.isTruthy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is truthy",
				types.InspectWithColor(conditionType),
			),
			node.Condition.Span(),
		)
		c.addUnreachableCodeError(node.ElseExpression.Span())
		node.SetType(thenType)
		return node
	}
	if c.isFalsy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is falsy",
				types.InspectWithColor(conditionType),
			),
			node.Condition.Span(),
		)
		c.addUnreachableCodeError(node.ThenExpression.Span())
		node.SetType(elseType)
		return node
	}

	node.SetType(c.newNormalisedUnion(thenType, elseType))
	return node
}

func (c *Checker) checkCollectionIfModifier(node *ast.ModifierNode) ast.ExpressionNode {
	c.pushNestedLocalEnv()
	node.Right = c.checkExpression(node.Right)
	conditionType := c.typeOfGuardVoid(node.Right)

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Right, assumptionTruthy)
	node.Left = c.checkExpression(node.Left)
	thenType := c.typeOfGuardVoid(node.Left)
	c.popLocalEnv()

	c.popLocalEnv()

	if c.isTruthy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is truthy",
				types.InspectWithColor(conditionType),
			),
			node.Right.Span(),
		)
		node.SetType(thenType)
		return node
	}
	if c.isFalsy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is falsy",
				types.InspectWithColor(conditionType),
			),
			node.Right.Span(),
		)
		c.addUnreachableCodeError(node.Left.Span())
		node.SetType(types.Never{})
		return node
	}

	node.SetType(thenType)
	return node
}

func (c *Checker) checkCollectionUnlessModifier(node *ast.ModifierNode) ast.ExpressionNode {
	c.pushNestedLocalEnv()
	node.Right = c.checkExpression(node.Right)
	conditionType := c.typeOfGuardVoid(node.Right)

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Right, assumptionFalsy)
	node.Left = c.checkExpression(node.Left)
	thenType := c.typeOfGuardVoid(node.Left)
	c.popLocalEnv()

	c.popLocalEnv()

	if c.isTruthy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is truthy",
				types.InspectWithColor(conditionType),
			),
			node.Right.Span(),
		)
		c.addUnreachableCodeError(node.Left.Span())
		node.SetType(types.Never{})
		return node
	}
	if c.isFalsy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is falsy",
				types.InspectWithColor(conditionType),
			),
			node.Right.Span(),
		)
		node.SetType(thenType)
		return node
	}

	node.SetType(thenType)
	return node
}

func (c *Checker) checkRecordIfModifier(node *ast.ModifierNode) (keyType, valueType types.Type) {
	c.pushNestedLocalEnv()
	node.Right = c.checkExpression(node.Right)
	conditionType := c.typeOfGuardVoid(node.Right)

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Right, assumptionTruthy)
	switch l := node.Left.(type) {
	case *ast.KeyValueExpressionNode:
		l.Key = c.checkExpression(l.Key)
		keyType = c.typeOf(l.Key)

		l.Value = c.checkExpression(l.Value)
		valueType = c.typeOf(l.Value)
	case *ast.SymbolKeyValueExpressionNode:
		keyType = c.Std(symbol.Symbol)

		l.Value = c.checkExpression(l.Value)
		valueType = c.typeOf(l.Value)
	default:
		panic(fmt.Sprintf("invalid record element node: %#v", node.Left))
	}
	c.popLocalEnv()

	c.popLocalEnv()

	if c.isTruthy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is truthy",
				types.InspectWithColor(conditionType),
			),
			node.Right.Span(),
		)
		return keyType, valueType
	}
	if c.isFalsy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is falsy",
				types.InspectWithColor(conditionType),
			),
			node.Right.Span(),
		)
		c.addUnreachableCodeError(node.Left.Span())
		return types.Never{}, types.Never{}
	}

	return keyType, valueType
}

func (c *Checker) checkRecordUnlessModifier(node *ast.ModifierNode) (keyType, valueType types.Type) {
	c.pushNestedLocalEnv()
	node.Right = c.checkExpression(node.Right)
	conditionType := c.typeOfGuardVoid(node.Right)

	c.pushNestedLocalEnv()
	c.narrowCondition(node.Right, assumptionFalsy)
	switch l := node.Left.(type) {
	case *ast.KeyValueExpressionNode:
		l.Key = c.checkExpression(l.Key)
		keyType = c.typeOf(l.Key)

		l.Value = c.checkExpression(l.Value)
		valueType = c.typeOf(l.Value)
	case *ast.SymbolKeyValueExpressionNode:
		keyType = c.Std(symbol.Symbol)

		l.Value = c.checkExpression(l.Value)
		valueType = c.typeOf(l.Value)
	default:
		panic(fmt.Sprintf("invalid record element node: %#v", node.Left))
	}
	c.popLocalEnv()

	c.popLocalEnv()

	if c.isTruthy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is truthy",
				types.InspectWithColor(conditionType),
			),
			node.Right.Span(),
		)
		c.addUnreachableCodeError(node.Left.Span())
		return types.Never{}, types.Never{}
	}
	if c.isFalsy(conditionType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is falsy",
				types.InspectWithColor(conditionType),
			),
			node.Right.Span(),
		)
		return keyType, valueType
	}

	return keyType, valueType
}

func (c *Checker) checkDoExpressionNode(node *ast.DoExpressionNode) ast.ExpressionNode {
	c.pushNestedLocalEnv()

	typ, _ := c.checkStatements(node.Body)
	node.SetType(typ)

	c.popLocalEnv()
	return node
}

func (c *Checker) checkPostfixExpressionNode(node *ast.PostfixExpressionNode) ast.ExpressionNode {
	switch node.Op.Type {
	case token.PLUS_PLUS:
		return c.checkPostfixExpression(node, "++")
	case token.MINUS_MINUS:
		return c.checkPostfixExpression(node, "--")
	default:
		c.addFailure(
			fmt.Sprintf("invalid unary operator %s", node.Op.String()),
			node.Span(),
		)
		return node
	}
}

func (c *Checker) checkNotOperator(node *ast.UnaryExpressionNode) ast.ExpressionNode {
	node.Right = c.checkExpression(node.Right)
	node.SetType(c.StdBool())
	return node
}

func (c *Checker) checkNilSafeSubscriptExpressionNode(node *ast.NilSafeSubscriptExpressionNode) ast.ExpressionNode {
	receiver, args, typ := c.checkSimpleMethodCall(
		node.Receiver,
		token.QUESTION_DOT,
		symbol.OpSubscript,
		nil,
		[]ast.ExpressionNode{node.Key},
		nil,
		node.Span(),
	)
	node.Receiver = receiver
	node.Key = args[0]

	node.SetType(typ)
	return node
}

func (c *Checker) checkSubscriptExpressionNode(node *ast.SubscriptExpressionNode) ast.ExpressionNode {
	receiver, args, typ := c.checkSimpleMethodCall(
		node.Receiver,
		token.DOT,
		symbol.OpSubscript,
		nil,
		[]ast.ExpressionNode{node.Key},
		nil,
		node.Span(),
	)
	node.Receiver = receiver
	node.Key = args[0]

	node.SetType(typ)
	return node
}

func (c *Checker) checkLogicalExpression(node *ast.LogicalExpressionNode) ast.ExpressionNode {
	switch node.Op.Type {
	case token.AND_AND:
		return c.checkLogicalAnd(node)
	case token.OR_OR:
		return c.checkLogicalOr(node)
	case token.QUESTION_QUESTION:
		return c.checkNilCoalescingOperator(node)
	default:
		node.SetType(types.Nothing{})
		c.addFailure(
			fmt.Sprintf(
				"invalid logical operator: `%s`",
				node.Op.String(),
			),
			node.Op.Span(),
		)
		return node
	}
}

func (c *Checker) checkNilCoalescingOperator(node *ast.LogicalExpressionNode) ast.ExpressionNode {
	node.Left = c.checkExpression(node.Left)
	c.pushNestedLocalEnv()
	c.narrowCondition(node.Left, assumptionNil)

	node.Right = c.checkExpression(node.Right)
	c.popLocalEnv()

	leftType := c.typeOfGuardVoid(node.Left)
	rightType := c.typeOf(node.Right)

	if c.isNil(leftType) {
		c.addWarning(
			"this condition will always have the same result",
			node.Left.Span(),
		)
		node.SetType(rightType)
		return node
	}
	if c.isNotNilable(leftType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` can never be nil",
				types.InspectWithColor(leftType),
			),
			node.Left.Span(),
		)
		c.addUnreachableCodeError(node.Right.Span())
		node.SetType(leftType)
		return node
	}
	node.SetType(c.newNormalisedUnion(c.toNonNilable(leftType), rightType))

	return node
}

func (c *Checker) checkLogicalOr(node *ast.LogicalExpressionNode) ast.ExpressionNode {
	node.Left = c.checkExpression(node.Left)
	c.pushNestedLocalEnv()
	c.narrowCondition(node.Left, assumptionFalsy)

	node.Right = c.checkExpression(node.Right)
	c.popLocalEnv()

	leftType := c.typeOfGuardVoid(node.Left)
	rightType := c.typeOf(node.Right)

	if c.isTruthy(leftType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is truthy",
				types.InspectWithColor(leftType),
			),
			node.Left.Span(),
		)
		c.addUnreachableCodeError(node.Right.Span())
		node.SetType(leftType)
		return node
	}
	if c.isFalsy(leftType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is falsy",
				types.InspectWithColor(leftType),
			),
			node.Left.Span(),
		)
		node.SetType(rightType)
		return node
	}
	union := c.newNormalisedUnion(c.toNonFalsy(leftType), rightType)
	node.SetType(union)

	return node
}

func (c *Checker) checkLogicalAnd(node *ast.LogicalExpressionNode) ast.ExpressionNode {
	node.Left = c.checkExpression(node.Left)
	c.pushNestedLocalEnv()
	c.narrowCondition(node.Left, assumptionTruthy)

	node.Right = c.checkExpression(node.Right)
	c.popLocalEnv()

	leftType := c.typeOfGuardVoid(node.Left)
	rightType := c.typeOf(node.Right)

	if c.isTruthy(leftType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is truthy",
				types.InspectWithColor(leftType),
			),
			node.Left.Span(),
		)
		node.SetType(rightType)
		return node
	}
	if c.isFalsy(leftType) {
		c.addWarning(
			fmt.Sprintf(
				"this condition will always have the same result since type `%s` is falsy",
				types.InspectWithColor(leftType),
			),
			node.Left.Span(),
		)
		c.addUnreachableCodeError(node.Right.Span())
		node.SetType(leftType)
		return node
	}

	node.SetType(c.newNormalisedUnion(c.toNonTruthy(leftType), rightType))

	return node
}

func (c *Checker) checkBinaryExpression(node *ast.BinaryExpressionNode) ast.ExpressionNode {
	switch node.Op.Type {
	case token.PLUS:
		node.Left = c.checkExpression(node.Left)
		node.Right = c.checkExpression(node.Right)
		node.SetType(c.checkArithmeticBinaryOperator(node.Left, node.Right, symbol.OpAdd, node.Span()))
	case token.MINUS:
		node.Left = c.checkExpression(node.Left)
		node.Right = c.checkExpression(node.Right)
		node.SetType(c.checkArithmeticBinaryOperator(node.Left, node.Right, symbol.OpSubtract, node.Span()))
	case token.STAR:
		node.Left = c.checkExpression(node.Left)
		node.Right = c.checkExpression(node.Right)
		node.SetType(c.checkArithmeticBinaryOperator(node.Left, node.Right, symbol.OpMultiply, node.Span()))
	case token.SLASH:
		node.Left = c.checkExpression(node.Left)
		node.Right = c.checkExpression(node.Right)
		node.SetType(c.checkArithmeticBinaryOperator(node.Left, node.Right, symbol.OpDivide, node.Span()))
	case token.STAR_STAR:
		node.Left = c.checkExpression(node.Left)
		node.Right = c.checkExpression(node.Right)
		node.SetType(c.checkArithmeticBinaryOperator(node.Left, node.Right, symbol.OpExponentiate, node.Span()))
	case token.PIPE_OP:
		return c.checkPipeExpression(node)
	case token.INSTANCE_OF_OP:
		c.checkInstanceOf(node, false)
	case token.REVERSE_INSTANCE_OF_OP:
		c.checkInstanceOf(node, true)
	case token.ISA_OP:
		c.checkIsA(node, false)
	case token.REVERSE_ISA_OP:
		c.checkIsA(node, true)
	case token.STRICT_EQUAL, token.STRICT_NOT_EQUAL:
		c.checkStrictEqual(node)
	case token.LAX_NOT_EQUAL:
		_, _, typ := c.checkSimpleMethodCall(
			node.Left,
			token.DOT,
			symbol.OpLaxEqual,
			nil,
			[]ast.ExpressionNode{node.Right},
			nil,
			node.Span(),
		)
		node.SetType(typ)
	case token.NOT_EQUAL:
		_, _, typ := c.checkSimpleMethodCall(
			node.Left,
			token.DOT,
			symbol.OpEqual,
			nil,
			[]ast.ExpressionNode{node.Right},
			nil,
			node.Span(),
		)
		node.SetType(typ)
	case token.LBITSHIFT, token.LTRIPLE_BITSHIFT,
		token.RBITSHIFT, token.RTRIPLE_BITSHIFT, token.AND,
		token.AND_TILDE, token.OR, token.XOR, token.PERCENT,
		token.LAX_EQUAL, token.EQUAL_EQUAL,
		token.GREATER, token.GREATER_EQUAL,
		token.LESS, token.LESS_EQUAL, token.SPACESHIP_OP:
		_, _, typ := c.checkSimpleMethodCall(
			node.Left,
			token.DOT,
			value.ToSymbol(node.Op.StringValue()),
			nil,
			[]ast.ExpressionNode{node.Right},
			nil,
			node.Span(),
		)
		node.SetType(typ)
	default:
		node.Left = c.checkExpression(node.Left)
		node.Right = c.checkExpression(node.Right)
		node.SetType(types.Nothing{})
		c.addFailure(
			fmt.Sprintf(
				"invalid binary operator: `%s`",
				node.Op.String(),
			),
			node.Op.Span(),
		)
	}

	return node
}

func (c *Checker) checkPipeExpression(node *ast.BinaryExpressionNode) ast.ExpressionNode {
	switch r := node.Right.(type) {
	case *ast.MethodCallNode:
		r.PositionalArguments = slices.Insert(r.PositionalArguments, 0, node.Left)
		return c.checkExpression(r)
	case *ast.CallNode:
		r.PositionalArguments = slices.Insert(r.PositionalArguments, 0, node.Left)
		return c.checkExpression(r)
	case *ast.ConstructorCallNode:
		r.PositionalArguments = slices.Insert(r.PositionalArguments, 0, node.Left)
		return c.checkExpression(r)
	case *ast.ReceiverlessMethodCallNode:
		r.PositionalArguments = slices.Insert(r.PositionalArguments, 0, node.Left)
		return c.checkExpression(r)
	case *ast.AttributeAccessNode:
		return c.checkExpression(
			ast.NewMethodCallNode(
				node.Span(),
				r.Receiver,
				token.New(r.Span(), token.DOT),
				r.AttributeName,
				[]ast.ExpressionNode{node.Left},
				nil,
			),
		)
	default:
		c.addFailure(
			fmt.Sprintf(
				"invalid right hand side of a pipe expression: `%T`",
				node.Right,
			),
			node.Right.Span(),
		)
		node.SetType(types.Nothing{})
		return node
	}
}

func (c *Checker) checkStrictEqual(node *ast.BinaryExpressionNode) {
	node.SetType(c.StdBool())
	node.Left = c.checkExpression(node.Left)
	node.Right = c.checkExpression(node.Right)
	leftType := c.typeOfGuardVoid(node.Left)
	rightType := c.typeOfGuardVoid(node.Right)

	if !c.typesIntersect(leftType, rightType) {
		c.addWarning(
			fmt.Sprintf(
				"this strict equality check is impossible, `%s` cannot ever be equal to `%s`",
				types.InspectWithColor(leftType),
				types.InspectWithColor(rightType),
			),
			node.Left.Span(),
		)
		return
	}
}

func (c *Checker) checkInstanceOf(node *ast.BinaryExpressionNode, reverse bool) {
	node.SetType(c.StdBool())
	node.Left = c.checkExpression(node.Left)
	node.Right = c.checkExpression(node.Right)
	var left ast.ExpressionNode
	var right ast.ExpressionNode
	if reverse {
		left = node.Right
		right = node.Left
	} else {
		left = node.Left
		right = node.Right
	}
	leftType := c.typeOfGuardVoid(left)
	rightType := c.typeOf(right)

	rightSingleton, ok := rightType.(*types.SingletonClass)
	if !ok {
		c.addFailure(
			"only classes are allowed as the right operand of the instance of operator `<<:`",
			right.Span(),
		)
		return
	}

	class, ok := rightSingleton.AttachedObject.(*types.Class)
	if !ok {
		c.addFailure(
			"only classes are allowed as the right operand of the instance of operator `<<:`",
			right.Span(),
		)
		return
	}

	if c.isSubtype(leftType, class, nil) {
		c.addWarning(
			fmt.Sprintf(
				"this \"instance of\" check is always true, `%s` will always be an instance of `%s`",
				types.InspectWithColor(leftType),
				types.InspectWithColor(class),
			),
			left.Span(),
		)
		return
	}

	if !c.isSubtype(class, leftType, nil) {
		c.addWarning(
			fmt.Sprintf(
				"impossible \"instance of\" check, `%s` cannot ever be an instance of `%s`",
				types.InspectWithColor(leftType),
				types.InspectWithColor(class),
			),
			left.Span(),
		)
	}
}

func (c *Checker) checkIsA(node *ast.BinaryExpressionNode, reverse bool) {
	node.SetType(c.StdBool())
	node.Left = c.checkExpression(node.Left)
	node.Right = c.checkExpression(node.Right)
	var left ast.ExpressionNode
	var right ast.ExpressionNode
	if reverse {
		left = node.Right
		right = node.Left
	} else {
		left = node.Left
		right = node.Right
	}
	leftType := c.typeOfGuardVoid(left)
	rightType := c.typeOf(right)

	rightSingleton, ok := rightType.(*types.SingletonClass)
	if !ok {
		c.addFailure(
			"only classes and mixins are allowed as the right operand of the is a operator `<:`",
			right.Span(),
		)
		return
	}

	switch rightSingleton.AttachedObject.(type) {
	case *types.Class, *types.Mixin:
	default:
		c.addFailure(
			"only classes and mixins are allowed as the right operand of the is a operator `<:`",
			right.Span(),
		)
		return
	}

	if c.isSubtype(leftType, rightSingleton.AttachedObject, nil) {
		c.addWarning(
			fmt.Sprintf(
				"this \"is a\" check is always true, `%s` will always be an instance of `%s`",
				types.InspectWithColor(leftType),
				types.InspectWithColor(rightSingleton.AttachedObject),
			),
			left.Span(),
		)
		return
	}

	if !c.canBeIsA(leftType, rightSingleton.AttachedObject) {
		c.addWarning(
			fmt.Sprintf(
				"impossible \"is a\" check, `%s` cannot ever be an instance of a descendant of `%s`",
				types.InspectWithColor(leftType),
				types.InspectWithColor(rightSingleton.AttachedObject),
			),
			left.Span(),
		)
	}
}

func (c *Checker) checkAbstractMethods(namespace types.Namespace, span *position.Span) {
	if namespace.IsAbstract() {
		return
	}

	errDetailsBuff := new(strings.Builder)
	for _, parentMethod := range c.abstractMethodsInNamespace(namespace.Parent()) {
		method := c.resolveNonAbstractMethodInNamespace(namespace, parentMethod.Name)
		if method == nil || method.IsAbstract() {
			fmt.Fprintf(
				errDetailsBuff,
				"\n  - method `%s`:\n      `%s`\n",
				types.InspectWithColor(parentMethod),
				parentMethod.InspectSignatureWithColor(false),
			)
		}
	}

	if errDetailsBuff.Len() > 0 {
		c.addFailure(
			fmt.Sprintf(
				"missing abstract method implementations in `%s`:\n%s",
				types.InspectWithColor(namespace),
				errDetailsBuff.String(),
			),
			span,
		)
	}
}

// Search through the ancestor chain of the current namespace
// looking for the direct parent of the proxy representing the given mixin.
func (c *Checker) findParentOfMixinProxy(mixin types.Namespace) types.Namespace {
	currentNamespace := c.currentConstScope().container
	currentParent := currentNamespace.Parent()

	for ; currentParent != nil; currentParent = currentParent.Parent() {
		if types.NamespacesAreEqual(currentParent, mixin) {
			return currentParent.Parent()
		}
	}

	return nil
}

type instanceVariableOverride struct {
	name value.Symbol

	super          types.Type
	superNamespace types.Namespace

	override          types.Type
	overrideNamespace types.Namespace
}

func (c *Checker) checkIncludeExpressionNode(node *ast.IncludeExpressionNode) {
	if c.typeOf(node) == nil {
		c.addFailure(
			"cannot include mixins in this context",
			node.Span(),
		)
		return
	}
	targetNamespace := c.currentMethodScope().container

	for _, constantNode := range node.Constants {
		constantType := c.typeOf(constantNode)
		var includedMixin types.Namespace
		switch t := constantType.(type) {
		case *types.Mixin:
			includedMixin = t
		case *types.Generic:
			includedMixin = t
		default:
			continue
		}

		parentOfMixin := c.findParentOfMixinProxy(includedMixin)
		if parentOfMixin == nil {
			continue
		}

		var incompatibleMethods []methodOverride
		for name, includedMethod := range c.methodsInNamespace(includedMixin) {
			superMethod := c.resolveMethodInNamespace(parentOfMixin, name)
			if !c.checkMethodCompatibility(superMethod, includedMethod, nil) {
				incompatibleMethods = append(incompatibleMethods, methodOverride{
					superMethod: superMethod,
					override:    includedMethod,
				})
			}
		}

		var incompatibleIvars []instanceVariableOverride
		for name, ivar := range c.instanceVariablesInNamespace(includedMixin) {
			includedIvar := ivar.Type
			includedNamespace := ivar.Namespace
			superIvar, superNamespace := types.GetInstanceVariableInNamespace(parentOfMixin, name)
			if superIvar == nil {
				continue
			}
			if !c.isTheSameType(superIvar, includedIvar, nil) {
				incompatibleIvars = append(incompatibleIvars, instanceVariableOverride{
					name:              name,
					super:             superIvar,
					superNamespace:    superNamespace,
					override:          includedIvar,
					overrideNamespace: includedNamespace,
				})
			}
		}

		if len(incompatibleMethods) == 0 && len(incompatibleIvars) == 0 {
			continue
		}

		detailsBuff := new(strings.Builder)
		for _, incompatibleMethod := range incompatibleMethods {
			override := incompatibleMethod.override
			superMethod := incompatibleMethod.superMethod

			overrideNamespaceName := types.InspectWithColor(override.DefinedUnder)
			superNamespaceName := types.InspectWithColor(superMethod.DefinedUnder)
			overrideNamespaceWidth := uniseg.StringWidth(types.Inspect(override.DefinedUnder))
			superNamespaceWidth := uniseg.StringWidth(types.Inspect(superMethod.DefinedUnder))
			var overrideWidthDiff int
			var superWidthDiff int
			if overrideNamespaceWidth < superNamespaceWidth {
				overrideWidthDiff = overrideWidthDiff - superNamespaceWidth
			} else {
				superWidthDiff = superNamespaceWidth - overrideNamespaceWidth
			}

			fmt.Fprintf(
				detailsBuff,
				"\n  - incompatible definitions of method `%s`\n      `%s`% *s has: `%s`\n      `%s`% *s has: `%s`\n",
				override.Name,
				overrideNamespaceName,
				overrideWidthDiff,
				"",
				override.InspectSignatureWithColor(false),
				superNamespaceName,
				superWidthDiff,
				"",
				superMethod.InspectSignatureWithColor(false),
			)
		}
		for _, incompatibleIvar := range incompatibleIvars {
			override := incompatibleIvar.override
			overrideNamespace := incompatibleIvar.overrideNamespace
			super := incompatibleIvar.super
			superNamespace := incompatibleIvar.superNamespace
			name := incompatibleIvar.name

			overrideNamespaceName := types.InspectWithColor(overrideNamespace)
			superNamespaceName := types.InspectWithColor(superNamespace)
			overrideNamespaceWidth := uniseg.StringWidth(types.Inspect(overrideNamespace))
			superNamespaceWidth := uniseg.StringWidth(types.Inspect(superNamespace))
			var overrideWidthDiff int
			var superWidthDiff int
			if overrideNamespaceWidth < superNamespaceWidth {
				overrideWidthDiff = overrideWidthDiff - superNamespaceWidth
			} else {
				superWidthDiff = superNamespaceWidth - overrideNamespaceWidth
			}

			fmt.Fprintf(
				detailsBuff,
				"\n  - incompatible definitions of instance variable `%s`\n      `%s`% *s has: `%s`\n      `%s`% *s has: `%s`\n",
				types.InspectInstanceVariableWithColor(name.String()),
				overrideNamespaceName,
				overrideWidthDiff,
				"",
				types.InspectInstanceVariableDeclarationWithColor(name.String(), override),
				superNamespaceName,
				superWidthDiff,
				"",
				types.InspectInstanceVariableDeclarationWithColor(name.String(), super),
			)
		}

		c.addFailure(
			fmt.Sprintf(
				"cannot include `%s` in `%s`:\n%s",
				types.InspectWithColor(includedMixin),
				types.InspectWithColor(targetNamespace),
				detailsBuff.String(),
			),
			constantNode.Span(),
		)

	}
}

func (c *Checker) typeOf(node ast.Node) types.Type {
	if node == nil {
		return types.Void{}
	}
	return node.Type(c.GlobalEnv)
}

func (c *Checker) typeOfGuardVoid(node ast.Node) types.Type {
	typ := c.typeOf(node)
	if node != nil && types.IsVoid(typ) {
		c.addFailure(
			fmt.Sprintf(
				"cannot use type `%s` as a value in this context",
				types.InspectWithColor(types.Void{}),
			),
			node.Span(),
		)
		return types.Nothing{}
	}
	return typ
}

func (c *Checker) checkGenericReceiverlessMethodCallNode(node *ast.GenericReceiverlessMethodCallNode) ast.ExpressionNode {
	method := c.getMethod(c.selfType, value.ToSymbol(node.MethodName), node.Span())
	if method == nil {
		c.checkExpressions(node.PositionalArguments)
		c.checkNamedArguments(node.NamedArguments)
		node.SetType(types.Nothing{})
		return node
	}

	if len(node.TypeArguments) > 0 {
		typeArgs, ok := c.checkTypeArguments(
			method,
			node.TypeArguments,
			method.TypeParameters,
			node.Span(),
		)
		if !ok {
			c.checkExpressions(node.PositionalArguments)
			c.checkNamedArguments(node.NamedArguments)
			node.SetType(types.Nothing{})
			return node
		}

		method = c.replaceTypeParametersInMethod(c.deepCopyMethod(method), typeArgs.ArgumentMap)
	} else if len(method.TypeParameters) > 0 {
		c.addTypeArgumentCountError(types.InspectWithColor(method), len(method.TypeParameters), len(node.TypeArguments), node.Span())
		node.SetType(types.Nothing{})
		return node
	}

	typedPositionalArguments := c.checkMethodArguments(method, node.PositionalArguments, node.NamedArguments, node.Span())
	node.PositionalArguments = typedPositionalArguments
	node.NamedArguments = nil
	node.SetType(method.ReturnType)
	return node
}

func (c *Checker) checkReceiverlessMethodCallNode(node *ast.ReceiverlessMethodCallNode) ast.ExpressionNode {
	method := c.getReceiverlessMethod(value.ToSymbol(node.MethodName), node.Span())
	if method == nil {
		c.checkExpressions(node.PositionalArguments)
		c.checkNamedArguments(node.NamedArguments)
		node.SetType(types.Nothing{})
		return node
	}

	var typedPositionalArguments []ast.ExpressionNode
	if len(method.TypeParameters) > 0 {
		var typeArgMap map[value.Symbol]*types.TypeArgument
		method = c.deepCopyMethod(method)
		typedPositionalArguments, typeArgMap = c.checkMethodArgumentsAndInferTypeArguments(
			method,
			node.PositionalArguments,
			node.NamedArguments,
			method.TypeParameters,
			node.Span(),
		)
		if typedPositionalArguments == nil {
			node.SetType(types.Nothing{})
			return node
		}
		if len(typeArgMap) != len(method.TypeParameters) {
			node.SetType(types.Nothing{})
			return node
		}
		method.ReturnType = c.replaceTypeParameters(method.ReturnType, typeArgMap)
		method.ThrowType = c.replaceTypeParameters(method.ThrowType, typeArgMap)
	} else {
		typedPositionalArguments = c.checkMethodArguments(method, node.PositionalArguments, node.NamedArguments, node.Span())
	}

	var receiver ast.ExpressionNode
	switch method.DefinedUnder.(type) {
	case *types.Module, *types.SingletonClass:
		// from using
		receiver = ast.NewPublicConstantNode(node.Span(), method.DefinedUnder.Name())
	default:
		// from self
		receiver = ast.NewSelfLiteralNode(node.Span())
	}

	newNode := ast.NewMethodCallNode(
		node.Span(),
		receiver,
		token.New(node.Span(), token.DOT),
		method.Name.String(),
		typedPositionalArguments,
		nil,
	)
	newNode.SetType(method.ReturnType)
	return newNode
}

func (c *Checker) checkNamedArguments(args []ast.NamedArgumentNode) {
	for _, arg := range args {
		arg, ok := arg.(*ast.NamedCallArgumentNode)
		if !ok {
			continue
		}

		c.checkExpression(arg.Value)
	}
}

func (c *Checker) addTypeParamBoundError(arg, boundType types.Type, boundName string, span *position.Span) {
	c.addFailure(
		fmt.Sprintf(
			"type `%s` does not satisfy the %s bound `%s`",
			types.InspectWithColor(arg),
			boundName,
			types.InspectWithColor(boundType),
		),
		span,
	)
}

func (c *Checker) addUpperBoundError(arg, bound types.Type, span *position.Span) {
	c.addTypeParamBoundError(arg, bound, "upper", span)
}

func (c *Checker) addLowerBoundError(arg, bound types.Type, span *position.Span) {
	c.addTypeParamBoundError(arg, bound, "lower", span)
}

func (c *Checker) checkTypeArguments(typ types.Type, typeArgs []ast.TypeNode, typeParams []*types.TypeParameter, span *position.Span) (*types.TypeArguments, bool) {
	if len(typeArgs) != len(typeParams) {
		c.addTypeArgumentCountError(types.InspectWithColor(typ), len(typeParams), len(typeArgs), span)
		return nil, false
	}

	typeArgumentMap := make(map[value.Symbol]*types.TypeArgument, len(typeParams))
	typeArgumentOrder := make([]value.Symbol, 0, len(typeParams))
	var fail bool
	for i := range len(typeParams) {
		typeParameter := typeParams[i]

		typeArgs[i] = c.checkTypeNode(typeArgs[i])
		typeArgumentNode := typeArgs[i]
		typeArgument := c.typeOf(typeArgumentNode)
		typeArgumentMap[typeParameter.Name] = types.NewTypeArgument(
			typeArgument,
			typeParameter.Variance,
		)
		typeArgumentOrder = append(typeArgumentOrder, typeParameter.Name)

		if !c.isSubtype(typeArgument, typeParameter.UpperBound, typeArgumentNode.Span()) {
			c.addUpperBoundError(typeArgument, typeParameter.UpperBound, typeArgumentNode.Span())
			fail = true
		}
		if !c.isSubtype(typeParameter.LowerBound, typeArgument, typeArgumentNode.Span()) {
			c.addLowerBoundError(typeArgument, typeParameter.LowerBound, typeArgumentNode.Span())
			fail = true
		}
		switch t := typeArgument.(type) {
		case *types.TypeParameter:
			if t.Variance != types.INVARIANT && t.Variance != typeParameter.Variance {
				c.addFailure(
					fmt.Sprintf(
						"%s type `%s` cannot appear in %s position",
						t.Variance.Name(),
						types.InspectWithColor(typeArgument),
						typeParameter.Variance.Name(),
					),
					typeArgumentNode.Span(),
				)
			}
		}
	}
	if fail {
		return nil, false
	}

	return types.NewTypeArguments(typeArgumentMap, typeArgumentOrder), true
}

func (c *Checker) checkNewExpressionNode(node *ast.NewExpressionNode) ast.ExpressionNode {
	var class *types.Class
	var isSingleton bool
	switch t := c.selfType.(type) {
	case *types.Class:
		class = t
	case *types.SingletonClass:
		if attached, ok := t.AttachedObject.(*types.Class); ok {
			class = attached
		}
		isSingleton = true
	}

	if class == nil {
		c.addFailure(
			fmt.Sprintf("`%s` cannot be instantiated", types.InspectWithColor(c.selfType)),
			node.Span(),
		)
		c.checkExpressions(node.PositionalArguments)
		c.checkNamedArguments(node.NamedArguments)
		node.SetType(types.Nothing{})
		return node
	}

	var typeArgs *types.TypeArguments
	var method *types.Method
	if len(class.TypeParameters()) > 0 {
		typeArgumentMap := make(map[value.Symbol]*types.TypeArgument, len(class.TypeParameters()))
		typeArgumentOrder := make([]value.Symbol, len(class.TypeParameters()))
		for i, param := range class.TypeParameters() {
			typeArgumentMap[param.Name] = types.NewTypeArgument(
				param,
				param.Variance,
			)
			typeArgumentOrder[i] = param.Name
		}
		typeArgs = types.NewTypeArguments(
			typeArgumentMap,
			typeArgumentOrder,
		)
		generic := types.NewGeneric(
			class,
			typeArgs,
		)
		method = c.getMethod(generic, symbol.M_init, nil)
	} else {
		method = c.getMethod(class, symbol.M_init, nil)
	}

	if method == nil {
		method = types.NewMethod(
			"",
			false,
			false,
			true,
			symbol.M_init,
			nil,
			nil,
			nil,
			nil,
			class,
		)
	}

	typedPositionalArguments := c.checkMethodArguments(method, node.PositionalArguments, node.NamedArguments, node.Span())

	node.PositionalArguments = typedPositionalArguments
	node.NamedArguments = nil
	if isSingleton {
		node.SetType(types.NewInstanceOf(types.Self{}))
	} else {
		node.SetType(types.Self{})
	}
	return node
}

func (c *Checker) checkGenericConstructorCallNode(node *ast.GenericConstructorCallNode) ast.ExpressionNode {
	classType, _ := c.resolveConstantType(node.Class)
	if classType == nil {
		classType = types.Nothing{}
	}

	if types.IsNothing(classType) {
		c.checkExpressions(node.PositionalArguments)
		c.checkNamedArguments(node.NamedArguments)
		node.SetType(types.Nothing{})
		return node
	}
	class, isClass := classType.(*types.Class)
	if !isClass {
		c.addFailure(
			fmt.Sprintf("`%s` cannot be instantiated", types.InspectWithColor(classType)),
			node.Span(),
		)
		c.checkExpressions(node.PositionalArguments)
		c.checkNamedArguments(node.NamedArguments)
		node.SetType(types.Nothing{})
		return node
	}

	if class.IsAbstract() {
		c.addFailure(
			fmt.Sprintf("cannot instantiate abstract class `%s`", types.InspectWithColor(class)),
			node.Span(),
		)
	}

	typeArgs, ok := c.checkTypeArguments(
		class,
		node.TypeArguments,
		class.TypeParameters(),
		node.Class.Span(),
	)
	if !ok {
		c.checkExpressions(node.PositionalArguments)
		c.checkNamedArguments(node.NamedArguments)
		node.SetType(types.Nothing{})
		return node
	}

	generic := types.NewGeneric(class, typeArgs)
	method := c.getMethod(generic, symbol.M_init, nil)
	if method == nil {
		method = types.NewMethod(
			"",
			false,
			false,
			true,
			symbol.M_init,
			nil,
			nil,
			nil,
			nil,
			class,
		)
	}

	typedPositionalArguments := c.checkMethodArguments(method, node.PositionalArguments, node.NamedArguments, node.Span())

	node.PositionalArguments = typedPositionalArguments
	node.NamedArguments = nil
	node.SetType(generic)
	return node
}

func (c *Checker) checkConstructorCallNode(node *ast.ConstructorCallNode) ast.ExpressionNode {
	classType, _ := c.resolveConstantType(node.Class)
	if classType == nil {
		classType = types.Nothing{}
	}

	if types.IsNothing(classType) {
		c.checkExpressions(node.PositionalArguments)
		c.checkNamedArguments(node.NamedArguments)
		node.SetType(types.Nothing{})
		return node
	}
	class, isClass := classType.(*types.Class)
	if !isClass {
		c.addFailure(
			fmt.Sprintf("`%s` cannot be instantiated", types.InspectWithColor(classType)),
			node.Span(),
		)
		c.checkExpressions(node.PositionalArguments)
		c.checkNamedArguments(node.NamedArguments)
		node.SetType(types.Nothing{})
		return node
	}

	if class.IsAbstract() {
		c.addFailure(
			fmt.Sprintf("cannot instantiate abstract class `%s`", types.InspectWithColor(class)),
			node.Span(),
		)
	}

	if !class.IsGeneric() {
		method := c.getMethod(class, symbol.M_init, nil)
		if method == nil {
			method = types.NewMethod(
				"",
				false,
				false,
				true,
				symbol.M_init,
				nil,
				nil,
				nil,
				nil,
				class,
			)
		}

		typedPositionalArguments := c.checkMethodArguments(method, node.PositionalArguments, node.NamedArguments, node.Span())

		node.PositionalArguments = typedPositionalArguments
		node.NamedArguments = nil
		node.SetType(class)
		return node
	}

	method := c._getMethodInNamespace(class, class, symbol.M_init, nil, false)
	if method == nil {
		method = types.NewMethod(
			"",
			false,
			false,
			true,
			symbol.M_init,
			nil,
			nil,
			nil,
			nil,
			class,
		)
	} else {
		method = c.deepCopyMethod(method)
	}

	typedPositionalArguments, typeArgMap := c.checkMethodArgumentsAndInferTypeArguments(
		method,
		node.PositionalArguments,
		node.NamedArguments,
		class.TypeParameters(),
		node.Span(),
	)
	if len(typeArgMap) != len(class.TypeParameters()) {
		node.SetType(types.Nothing{})
		return node
	}
	method.ReturnType = c.replaceTypeParameters(method.ReturnType, typeArgMap)
	method.ThrowType = c.replaceTypeParameters(method.ThrowType, typeArgMap)
	typeArgOrder := make([]value.Symbol, len(class.TypeParameters()))
	for i, param := range class.TypeParameters() {
		typeArgOrder[i] = param.Name
	}
	generic := types.NewGeneric(
		class,
		types.NewTypeArguments(
			typeArgMap,
			typeArgOrder,
		),
	)
	node.PositionalArguments = typedPositionalArguments
	node.NamedArguments = nil
	node.SetType(generic)
	return node
}

func (c *Checker) checkCallNode(node *ast.CallNode) ast.ExpressionNode {
	var typ types.Type
	var op token.Type
	if node.NilSafe {
		op = token.QUESTION_DOT
	} else {
		op = token.DOT
	}
	node.Receiver, node.PositionalArguments, typ = c.checkSimpleMethodCall(
		node.Receiver,
		op,
		value.ToSymbol("call"),
		nil,
		node.PositionalArguments,
		node.NamedArguments,
		node.Span(),
	)
	node.SetType(typ)
	return node
}

func (c *Checker) checkMethodCallNode(node *ast.MethodCallNode) ast.ExpressionNode {
	var typ types.Type
	node.Receiver, node.PositionalArguments, typ = c.checkSimpleMethodCall(
		node.Receiver,
		node.Op.Type,
		value.ToSymbol(node.MethodName),
		nil,
		node.PositionalArguments,
		node.NamedArguments,
		node.Span(),
	)
	node.SetType(typ)
	return node
}

func (c *Checker) checkGenericMethodCallNode(node *ast.GenericMethodCallNode) ast.ExpressionNode {
	var typ types.Type
	node.Receiver, node.PositionalArguments, typ = c.checkSimpleMethodCall(
		node.Receiver,
		node.Op.Type,
		value.ToSymbol(node.MethodName),
		node.TypeArguments,
		node.PositionalArguments,
		node.NamedArguments,
		node.Span(),
	)
	node.SetType(typ)
	return node
}

func (c *Checker) checkClosureLiteralNodeWithType(node *ast.ClosureLiteralNode, closureType *types.Closure) ast.ExpressionNode {
	baseMethod := closureType.Method(symbol.M_call)
	closure := types.NewClosure(nil)
	method, mod := c.declareMethod(
		baseMethod,
		closure,
		"",
		false,
		false,
		true,
		symbol.M_call,
		nil,
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Span(),
	)
	returnTypeNode, throwTypeNode := c.checkMethod(
		closure,
		method,
		nil,
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Body,
		node.Span(),
	)
	closure.Body = method
	node.ReturnType = returnTypeNode
	node.ThrowType = throwTypeNode
	node.SetType(closure)
	if mod != nil {
		c.popConstScope()
	}
	return node
}

func (c *Checker) checkClosureLiteralNode(node *ast.ClosureLiteralNode) ast.ExpressionNode {
	closure := types.NewClosure(nil)
	method, mod := c.declareMethod(
		nil,
		closure,
		"",
		false,
		false,
		true,
		symbol.M_call,
		nil,
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Span(),
	)
	returnTypeNode, throwTypeNode := c.checkMethod(
		closure,
		method,
		nil,
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Body,
		node.Span(),
	)
	node.ReturnType = returnTypeNode
	node.ThrowType = throwTypeNode
	closure.Body = method
	node.SetType(closure)
	if mod != nil {
		c.popConstScope()
	}
	return node
}

func (c *Checker) checkAttributeAccessNode(node *ast.AttributeAccessNode) ast.ExpressionNode {
	var newNode ast.ExpressionNode = ast.NewMethodCallNode(
		node.Span(),
		node.Receiver,
		token.New(node.Span(), token.DOT),
		node.AttributeName,
		nil,
		nil,
	)
	return c.checkExpression(newNode)
}

func (c *Checker) checkLogicalOperatorAssignmentExpression(node *ast.AssignmentExpressionNode, operator token.Type) ast.ExpressionNode {
	return c.checkExpression(
		ast.NewAssignmentExpressionNode(
			node.Span(),
			token.New(node.Op.Span(), token.EQUAL_OP),
			node.Left,
			ast.NewLogicalExpressionNode(
				node.Right.Span(),
				token.New(node.Op.Span(), operator),
				node.Left,
				node.Right,
			),
		),
	)
}

func (c *Checker) checkBinaryOperatorAssignmentExpression(node *ast.AssignmentExpressionNode, operator token.Type) ast.ExpressionNode {
	return c.checkExpression(
		ast.NewAssignmentExpressionNode(
			node.Span(),
			token.New(node.Op.Span(), token.EQUAL_OP),
			node.Left,
			ast.NewBinaryExpressionNode(
				node.Right.Span(),
				token.New(node.Op.Span(), operator),
				node.Left,
				node.Right,
			),
		),
	)
}

func (c *Checker) checkAssignmentExpressionNode(node *ast.AssignmentExpressionNode) ast.ExpressionNode {
	span := node.Span()
	switch node.Op.Type {
	case token.EQUAL_OP:
		return c.checkAssignment(node)
	case token.QUESTION_QUESTION_EQUAL:
		return c.checkLogicalOperatorAssignmentExpression(node, token.QUESTION_QUESTION)
	case token.OR_OR_EQUAL:
		return c.checkLogicalOperatorAssignmentExpression(node, token.OR_OR)
	case token.AND_AND_EQUAL:
		return c.checkLogicalOperatorAssignmentExpression(node, token.AND_AND)
	case token.PLUS_EQUAL:
		return c.checkBinaryOperatorAssignmentExpression(node, token.PLUS)
	case token.MINUS_EQUAL:
		return c.checkBinaryOperatorAssignmentExpression(node, token.MINUS)
	case token.STAR_EQUAL:
		return c.checkBinaryOperatorAssignmentExpression(node, token.STAR)
	case token.SLASH_EQUAL:
		return c.checkBinaryOperatorAssignmentExpression(node, token.SLASH)
	case token.STAR_STAR_EQUAL:
		return c.checkBinaryOperatorAssignmentExpression(node, token.STAR_STAR)
	case token.PERCENT_EQUAL:
		return c.checkBinaryOperatorAssignmentExpression(node, token.PERCENT)
	case token.AND_EQUAL:
		return c.checkBinaryOperatorAssignmentExpression(node, token.AND)
	case token.OR_EQUAL:
		return c.checkBinaryOperatorAssignmentExpression(node, token.OR)
	case token.XOR_EQUAL:
		return c.checkBinaryOperatorAssignmentExpression(node, token.XOR)
	case token.LBITSHIFT_EQUAL:
		return c.checkBinaryOperatorAssignmentExpression(node, token.LBITSHIFT)
	case token.LTRIPLE_BITSHIFT_EQUAL:
		return c.checkBinaryOperatorAssignmentExpression(node, token.LTRIPLE_BITSHIFT)
	case token.RBITSHIFT_EQUAL:
		return c.checkBinaryOperatorAssignmentExpression(node, token.RBITSHIFT)
	case token.RTRIPLE_BITSHIFT_EQUAL:
		return c.checkBinaryOperatorAssignmentExpression(node, token.RTRIPLE_BITSHIFT)
	case token.COLON_EQUAL:
		return c.checkShortVariableDeclaration(node)
	default:
		c.addFailure(
			fmt.Sprintf("assignment using this operator has not been implemented: %s", node.Op.Type.String()),
			span,
		)
		return node
	}
}

func (c *Checker) checkShortVariableDeclaration(node *ast.AssignmentExpressionNode) ast.ExpressionNode {
	var name string
	switch left := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		name = left.Value
	case *ast.PrivateIdentifierNode:
		name = left.Value
	}
	init, _, typ := c.checkVariableDeclaration(name, node.Right, nil, node.Span())
	node.Right = init
	node.SetType(typ)
	return node
}

func (c *Checker) checkAssignment(node *ast.AssignmentExpressionNode) ast.ExpressionNode {
	switch n := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		return c.checkLocalVariableAssignment(n.Value, node)
	case *ast.PrivateIdentifierNode:
		return c.checkLocalVariableAssignment(n.Value, node)
	case *ast.SubscriptExpressionNode:
		return c.checkSubscriptAssignment(n, node)
	case *ast.InstanceVariableNode:
		return c.checkInstanceVariableAssignment(n.Value, node)
	case *ast.AttributeAccessNode:
		return c.checkAttributeAssignment(n, node)
	default:
		c.addFailure(
			fmt.Sprintf("cannot assign to: %T", node.Left),
			node.Span(),
		)
		return node
	}
}

func (c *Checker) checkSubscriptAssignment(subscriptNode *ast.SubscriptExpressionNode, assignmentNode *ast.AssignmentExpressionNode) ast.ExpressionNode {
	receiver, args, _ := c.checkSimpleMethodCall(
		subscriptNode.Receiver,
		token.DOT,
		symbol.OpSubscriptSet,
		nil,
		[]ast.ExpressionNode{subscriptNode.Key, assignmentNode.Right},
		nil,
		assignmentNode.Span(),
	)
	subscriptNode.Receiver = receiver
	subscriptNode.Key = args[0]
	assignmentNode.Right = args[1]

	assignmentNode.SetType(c.typeOf(assignmentNode.Right))
	return assignmentNode
}

func (c *Checker) checkAttributeAssignment(attributeNode *ast.AttributeAccessNode, assignmentNode *ast.AssignmentExpressionNode) ast.ExpressionNode {
	receiver, args, _ := c.checkSimpleMethodCall(
		attributeNode.Receiver,
		token.DOT,
		value.ToSymbol(attributeNode.AttributeName+"="),
		nil,
		[]ast.ExpressionNode{assignmentNode.Right},
		nil,
		assignmentNode.Span(),
	)
	attributeNode.Receiver = receiver
	assignmentNode.Right = args[0]

	assignmentNode.SetType(c.typeOf(assignmentNode.Right))
	return assignmentNode
}

func (c *Checker) checkInstanceVariableAssignment(name string, node *ast.AssignmentExpressionNode) ast.ExpressionNode {
	ivarType := c.checkInstanceVariable(name, node.Left.Span())

	node.Right = c.checkExpression(node.Right)
	assignedType := c.typeOfGuardVoid(node.Right)
	c.checkCanAssign(assignedType, ivarType, node.Right.Span())
	return node
}

func (c *Checker) checkLocalVariableAssignment(name string, node *ast.AssignmentExpressionNode) ast.ExpressionNode {
	var variableType types.Type
	variable, _ := c.resolveLocal(name, node.Left.Span())
	if variable == nil {
		node.Left.SetType(types.Nothing{})
		variableType = types.Nothing{}
	} else if variable.singleAssignment && variable.initialised {
		c.addFailure(
			fmt.Sprintf("local value `%s` cannot be reassigned", name),
			node.Left.Span(),
		)
		variableType = variable.typ
	} else {
		variableType = variable.typ
	}

	node.Right = c.checkExpression(node.Right)
	assignedType := c.typeOfGuardVoid(node.Right)
	c.checkCanAssign(assignedType, variableType, node.Right.Span())
	node.SetType(assignedType)
	return node
}

func (c *Checker) checkInterpolatedRegexLiteralNode(node *ast.InterpolatedRegexLiteralNode) {
	for _, contentSection := range node.Content {
		c.checkRegexContent(contentSection)
	}
}

func (c *Checker) checkRegexContent(node ast.RegexLiteralContentNode) {
	switch n := node.(type) {
	case *ast.RegexInterpolationNode:
		expr := c.checkExpression(n.Expression)
		n.Expression = expr
		c.isSubtype(c.typeOf(n.Expression), c.StdStringConvertible(), expr.Span())
	case *ast.RegexLiteralContentSectionNode:
	default:
		c.addFailure(
			fmt.Sprintf("invalid regex content %T", node),
			node.Span(),
		)
	}
}

func (c *Checker) checkInterpolatedStringLiteralNode(node *ast.InterpolatedStringLiteralNode) {
	for _, contentSection := range node.Content {
		c.checkStringContent(contentSection)
	}
}

func (c *Checker) checkStringContent(node ast.StringLiteralContentNode) {
	switch n := node.(type) {
	case *ast.StringInspectInterpolationNode:
		expr := c.checkExpression(n.Expression)
		n.Expression = expr
		c.isSubtype(c.typeOf(n.Expression), c.StdInspectable(), expr.Span())
	case *ast.StringInterpolationNode:
		expr := c.checkExpression(n.Expression)
		n.Expression = expr
		c.isSubtype(c.typeOf(n.Expression), c.StdStringConvertible(), expr.Span())
	case *ast.StringLiteralContentSectionNode:
	default:
		c.addFailure(
			fmt.Sprintf("invalid string content %T", node),
			node.Span(),
		)
	}
}

func (c *Checker) addFailureWithLocation(message string, loc *position.Location) {
	c.Errors.AddFailure(
		message,
		loc,
	)
}

func (c *Checker) addWarning(message string, span *position.Span) {
	if span == nil {
		return
	}
	c.Errors.AddWarning(
		message,
		c.newLocation(span),
	)
}

func (c *Checker) addFailure(message string, span *position.Span) {
	if span == nil {
		return
	}
	c.Errors.AddFailure(
		message,
		c.newLocation(span),
	)
}

func (c *Checker) addUnreachableCodeError(span *position.Span) {
	c.addWarning(
		"unreachable code",
		span,
	)
}

// Get the type of the subtype with the given name
func (c *Checker) resolveTypeInRoot(constantExpression ast.ExpressionNode) (types.Namespace, types.Type, string, string) {
	switch constant := constantExpression.(type) {
	case *ast.PublicConstantNode:
		return c.GlobalEnv.Root, c.resolveSimpleTypeInRoot(constant.Value), constant.Value, constant.Value
	case *ast.PrivateConstantNode:
		return c.GlobalEnv.Root, c.resolveSimpleTypeInRoot(constant.Value), constant.Value, constant.Value
	case *ast.ConstantLookupNode:
		return c.resolveTypeLookupInRoot(constant)
	default:
		panic(fmt.Sprintf("invalid constant node: %T", constantExpression))
	}
}

// Get the type of the constant with the given name
func (c *Checker) resolveSimpleTypeInRoot(name string) types.Type {
	root := c.GlobalEnv.Root
	constant, ok := root.SubtypeString(name)
	if ok {
		return constant.Type
	}
	return nil
}

// Get the type of the constant with the given name
func (c *Checker) resolveSimpleConstantInRoot(name string) types.Type {
	root := c.GlobalEnv.Root
	constant, ok := root.ConstantString(name)
	if ok {
		return constant.Type
	}
	return nil
}

func (c *Checker) resolveTypeLookupInRoot(node *ast.ConstantLookupNode) (_parentNamespace types.Namespace, _typ types.Type, _fullName, _constName string) {
	return c._resolveConstantLookupTypeInRoot(node, true)
}

func (c *Checker) _resolveConstantLookupTypeInRoot(node *ast.ConstantLookupNode, firstCall bool) (_parentNamespace types.Namespace, _typ types.Type, _fullName, _constName string) {
	var leftContainerType types.Type
	var fullName, leftContainerName string

	switch l := node.Left.(type) {
	case *ast.PublicConstantNode:
		namespace := c.GlobalEnv.Root
		leftConstant, ok := namespace.ConstantString(l.Value)
		leftContainerType = leftConstant.Type
		leftContainerName = types.MakeFullConstantName(namespace.Name(), l.Value)
		if !ok {
			placeholder := types.NewNamespacePlaceholder(leftContainerName)
			placeholder.Locations.Append(c.newLocation(l.Span()))
			leftContainerType = placeholder
			c.registerPlaceholderNamespace(placeholder)
			namespace.DefineConstant(value.ToSymbol(l.Value), types.NewSingletonClass(placeholder, nil))
		} else if placeholder, ok := leftContainerType.(*types.NamespacePlaceholder); ok {
			placeholder.Locations.Append(c.newLocation(l.Span()))
		}
	case *ast.PrivateConstantNode:
		namespace := c.GlobalEnv.Root
		leftConstant, ok := namespace.ConstantString(l.Value)
		leftContainerType = leftConstant.Type
		leftContainerName = types.MakeFullConstantName(namespace.Name(), l.Value)
		if !ok {
			placeholder := types.NewNamespacePlaceholder(leftContainerName)
			placeholder.Locations.Append(c.newLocation(l.Span()))
			leftContainerType = placeholder
			c.registerPlaceholderNamespace(placeholder)
			namespace.DefineConstant(value.ToSymbol(l.Value), types.NewSingletonClass(placeholder, nil))
		} else if placeholder, ok := leftContainerType.(*types.NamespacePlaceholder); ok {
			placeholder.Locations.Append(c.newLocation(l.Span()))
		}
	case nil:
		leftContainerType = c.GlobalEnv.Root
	case *ast.ConstantLookupNode:
		_, leftContainerType, leftContainerName, _ = c._resolveConstantLookupTypeInRoot(l, false)
	default:
		c.addFailure(
			fmt.Sprintf("invalid constant node %T", node),
			node.Span(),
		)
		return nil, nil, "", ""
	}

	var rightName string
	switch r := node.Right.(type) {
	case *ast.PublicConstantNode:
		rightName = r.Value
	case *ast.PrivateConstantNode:
		rightName = r.Value
		c.addFailure(
			fmt.Sprintf("cannot read private constant `%s`", rightName),
			node.Span(),
		)
	default:
		c.addFailure(
			fmt.Sprintf("invalid constant node %T", node),
			node.Span(),
		)
		return nil, nil, "", ""
	}

	fullName = types.MakeFullConstantName(leftContainerName, rightName)
	if leftContainerType == nil {
		return nil, nil, fullName, rightName
	}
	var leftContainer types.Namespace
	switch l := leftContainerType.(type) {
	case *types.Module:
		leftContainer = l
	case *types.Class:
		leftContainer = l
	case *types.Mixin:
		leftContainer = l
	case *types.NamespacePlaceholder:
		leftContainer = l
	case *types.SingletonClass:
		leftContainer = l.AttachedObject
	default:
		c.addFailure(
			fmt.Sprintf("cannot read constants from `%s`, it is not a constant container", leftContainerName),
			node.Span(),
		)
		return nil, nil, fullName, rightName
	}

	rightSymbol := value.ToSymbol(rightName)
	constant, ok := leftContainer.Constant(rightSymbol)
	if len(constant.FullName) > 0 {
		fullName = constant.FullName
	}
	constantType := constant.Type
	if !ok && !firstCall {
		placeholder := types.NewNamespacePlaceholder(fullName)
		placeholder.Locations.Append(c.newLocation(node.Right.Span()))
		constantType = placeholder
		c.registerPlaceholderNamespace(placeholder)
		leftContainer.DefineConstant(rightSymbol, types.NewSingletonClass(placeholder, nil))
	} else if placeholder, ok := constantType.(*types.NamespacePlaceholder); ok {
		placeholder.Locations.Append(c.newLocation(node.Right.Span()))
	}

	return leftContainer, constantType, fullName, rightName
}

// Get the type of the constant with the given name
func (c *Checker) resolveConstantForDeclaration(constantExpression ast.ExpressionNode) (types.Namespace, types.Type, string) {
	switch constant := constantExpression.(type) {
	case *ast.PublicConstantNode:
		return c.resolveSimpleConstantForSetter(constant.Value)
	case *ast.PrivateConstantNode:
		return c.resolveSimpleConstantForSetter(constant.Value)
	case *ast.ConstantLookupNode:
		return c.resolveConstantLookupForDeclaration(constant)
	default:
		panic(fmt.Sprintf("invalid constant node: %T", constantExpression))
	}
}

func (c *Checker) registerPlaceholderNamespace(placeholder *types.NamespacePlaceholder) {
	c.namespacePlaceholders.Append(placeholder)
}

func (c *Checker) registerPlaceholder(placeholder *types.Placeholder) {
	c.placeholders.Append(placeholder)
}

func (c *Checker) resolveConstantLookupForDeclaration(node *ast.ConstantLookupNode) (types.Namespace, types.Type, string) {
	return c._resolveConstantLookupForDeclaration(node, true)
}

func (c *Checker) _resolveConstantLookupForDeclaration(node *ast.ConstantLookupNode, firstCall bool) (types.Namespace, types.Type, string) {
	var leftContainerType types.Type
	var leftContainerName string

	switch l := node.Left.(type) {
	case *ast.PublicConstantNode:
		namespace := c.currentConstScope().container
		leftConstant, ok := namespace.ConstantString(l.Value)
		leftContainerType = leftConstant.Type
		leftContainerName = types.MakeFullConstantName(namespace.Name(), l.Value)
		if !ok {
			placeholder := types.NewNamespacePlaceholder(leftContainerName)
			placeholder.Locations.Append(c.newLocation(l.Span()))
			leftContainerType = placeholder
			c.registerPlaceholderNamespace(placeholder)
			namespace.DefineConstant(value.ToSymbol(l.Value), types.NewSingletonClass(placeholder, nil))
		} else if placeholder, ok := leftContainerType.(*types.NamespacePlaceholder); ok {
			placeholder.Locations.Append(c.newLocation(l.Span()))
		}
	case *ast.PrivateConstantNode:
		namespace := c.currentConstScope().container
		leftConstant, ok := namespace.ConstantString(l.Value)
		leftContainerType = leftConstant.Type
		leftContainerName = types.MakeFullConstantName(namespace.Name(), l.Value)
		if !ok {
			placeholder := types.NewNamespacePlaceholder(leftContainerName)
			placeholder.Locations.Append(c.newLocation(l.Span()))
			leftContainerType = placeholder
			c.registerPlaceholderNamespace(placeholder)
			namespace.DefineConstant(value.ToSymbol(l.Value), types.NewSingletonClass(placeholder, nil))
		} else if placeholder, ok := leftContainerType.(*types.NamespacePlaceholder); ok {
			placeholder.Locations.Append(c.newLocation(l.Span()))
		}
	case nil:
		leftContainerType = c.GlobalEnv.Root
	case *ast.ConstantLookupNode:
		_, leftContainerType, leftContainerName = c._resolveConstantLookupForDeclaration(l, false)
	default:
		c.addFailure(
			fmt.Sprintf("invalid constant node %T", node),
			node.Span(),
		)
		return nil, nil, ""
	}

	var rightName string
	switch r := node.Right.(type) {
	case *ast.PublicConstantNode:
		rightName = r.Value
	case *ast.PrivateConstantNode:
		rightName = r.Value
		c.addFailure(
			fmt.Sprintf("cannot read private constant `%s`", rightName),
			node.Span(),
		)
	default:
		c.addFailure(
			fmt.Sprintf("invalid constant node %T", node),
			node.Span(),
		)
		return nil, nil, ""
	}

	constantName := types.MakeFullConstantName(leftContainerName, rightName)
	if leftContainerType == nil {
		return nil, nil, constantName
	}
	var leftContainer types.Namespace
	switch l := leftContainerType.(type) {
	case *types.Module:
		leftContainer = l
	case *types.Class:
		leftContainer = l
	case *types.Mixin:
		leftContainer = l
	case *types.NamespacePlaceholder:
		leftContainer = l
	case *types.SingletonClass:
		leftContainer = l.AttachedObject
	default:
		c.addFailure(
			fmt.Sprintf("cannot read constants from `%s`, it is not a constant container", leftContainerName),
			node.Span(),
		)
		return nil, nil, constantName
	}

	rightSymbol := value.ToSymbol(rightName)
	constant, ok := leftContainer.Constant(rightSymbol)
	constantType := constant.Type
	if !ok && !firstCall {
		placeholder := types.NewNamespacePlaceholder(constantName)
		placeholder.Locations.Append(c.newLocation(node.Right.Span()))
		constantType = placeholder
		c.registerPlaceholderNamespace(placeholder)
		leftContainer.DefineConstant(rightSymbol, types.NewSingletonClass(placeholder, nil))
	} else if placeholder, ok := constantType.(*types.NamespacePlaceholder); ok {
		placeholder.Locations.Append(c.newLocation(node.Right.Span()))
	}

	return leftContainer, constantType, constantName
}

// Get the type of the constant with the given name
func (c *Checker) resolveSimpleConstantForSetter(name string) (types.Namespace, types.Type, string) {
	namespace := c.currentConstScope().container
	constant, ok := namespace.ConstantString(name)
	var fullName string
	if len(constant.FullName) > 0 {
		fullName = constant.FullName
	} else {
		fullName = types.MakeFullConstantName(namespace.Name(), name)
	}

	if ok {
		return namespace, constant.Type, constant.FullName
	}
	return namespace, nil, fullName
}

func (c *Checker) checkIfTypeParameterIsAllowed(typ types.Type, span *position.Span) bool {
	t, ok := typ.(*types.TypeParameter)
	if !ok {
		return true
	}

	switch c.mode {
	case inputPositionTypeMode:
		if t.Variance == types.COVARIANT {
			c.addFailure(
				fmt.Sprintf("covariant type parameter `%s` cannot appear in input positions", types.InspectWithColor(t)),
				span,
			)
			return false
		}
		enclosingScope := c.enclosingConstScope().container
		if e, ok := enclosingScope.(*types.TypeParamNamespace); ok {
			if _, ok := e.Subtype(t.Name); ok {
				break
			}
		}
		currentScope := c.currentConstScope().container
		if _, ok := currentScope.Subtype(t.Name); ok {
			break
		}

		c.addFailure(
			fmt.Sprintf("undefined type `%s`", types.InspectWithColor(t)),
			span,
		)
		return false
	case outputPositionTypeMode:
		if t.Variance == types.CONTRAVARIANT {
			c.addFailure(
				fmt.Sprintf("contravariant type parameter `%s` cannot appear in output positions", types.InspectWithColor(t)),
				span,
			)
			return false
		}
		enclosingScope := c.enclosingConstScope().container
		if e, ok := enclosingScope.(*types.TypeParamNamespace); ok {
			if _, ok := e.Subtype(t.Name); ok {
				break
			}
		}
		currentScope := c.currentConstScope().container
		if _, ok := currentScope.Subtype(t.Name); ok {
			break
		}

		c.addFailure(
			fmt.Sprintf("undefined type `%s`", types.InspectWithColor(t)),
			span,
		)
		return false
	case namedGenericTypeDefinitionMode, inheritanceMode, instanceVariableMode:
		enclosingScope := c.enclosingConstScope().container
		if _, ok := enclosingScope.Subtype(t.Name); !ok {
			c.addFailure(
				fmt.Sprintf("undefined type `%s`", types.InspectWithColor(t)),
				span,
			)
			return false
		}
	case methodMode, closureInferReturnTypeMode:
		enclosingScope := c.enclosingConstScope().container
		if e, ok := enclosingScope.(*types.TypeParamNamespace); ok {
			if _, ok := e.Subtype(t.Name); ok {
				break
			}
		}
		currentScope := c.currentConstScope().container
		if _, ok := currentScope.Subtype(t.Name); ok {
			break
		}

		c.addFailure(
			fmt.Sprintf("undefined type `%s`", types.InspectWithColor(t)),
			span,
		)
		return false
	default:
		c.addFailure(
			fmt.Sprintf("type parameter `%s` cannot be used in this context", types.InspectWithColor(t)),
			span,
		)
		return false
	}
	return true
}

// Get the type with the given name
func (c *Checker) resolveType(name string, span *position.Span) (types.Type, string) {
	for i := range len(c.constantScopes) {
		constScope := c.constantScopes[len(c.constantScopes)-i-1]
		constant, ok := constScope.container.SubtypeString(name)
		if !ok {
			continue
		}

		var fullName string
		if len(constant.FullName) > 0 {
			fullName = constant.FullName
		} else {
			fullName = types.MakeFullConstantName(constScope.container.Name(), name)
		}
		if !c.checkTypeIfNecessary(fullName, span) {
			return nil, fullName
		}

		if types.IsNoValue(constant.Type) || types.IsPlaceholder(constant.Type) {
			c.addFailure(
				fmt.Sprintf("undefined type `%s`", lexer.Colorize(fullName)),
				span,
			)
			return nil, fullName
		}

		if c.checkIfTypeParameterIsAllowed(constant.Type, span) {
			return constant.Type, fullName
		}
		return nil, fullName
	}

	c.addFailure(
		fmt.Sprintf("undefined type `%s`", name),
		span,
	)
	return nil, name
}

// Add the local with the given name to the current local environment
func (c *Checker) addLocal(name string, l *local) {
	env := c.currentLocalEnv()
	env.locals[value.ToSymbol(name)] = l
}

// Get the local with the specified name from the current local environment
func (c *Checker) getLocal(name string) *local {
	env := c.currentLocalEnv()
	local := env.getLocal(name)
	if local == nil || local.shadow {
		return nil
	}

	return local
}

// Get the instance variable with the specified name
func (c *Checker) getInstanceVariableIn(name value.Symbol, typ types.Namespace) (types.Type, types.Namespace) {
	if typ == nil {
		return nil, typ
	}

	var generics []*types.Generic

	for parent := range types.Parents(typ) {
		if generic, ok := parent.(*types.Generic); ok {
			generics = append(generics, generic)
		}
		ivar := parent.InstanceVariable(name)
		if ivar == nil {
			continue
		}

		if len(generics) < 1 {
			return ivar, parent
		}

		for i := len(generics) - 1; i >= 0; i-- {
			generic := generics[i]
			ivar = c.replaceTypeParameters(ivar, generic.ArgumentMap)
		}
		return ivar, parent
	}

	return nil, typ
}

func (c *Checker) instanceVariablesInNamespace(namespace types.Namespace) iter.Seq2[value.Symbol, types.InstanceVariable] {
	return func(yield func(name value.Symbol, ivar types.InstanceVariable) bool) {
		var generics []*types.Generic
		seenIvars := make(map[value.Symbol]bool)

		for parent := range types.Parents(namespace) {
			if generic, ok := parent.(*types.Generic); ok {
				generics = append(generics, generic)
			}
			for name, ivar := range parent.InstanceVariables() {
				if seenIvars[name] {
					continue
				}
				if len(generics) < 1 {
					ivarStruct := types.InstanceVariable{Type: ivar, Namespace: parent}
					if !yield(name, ivarStruct) {
						return
					}
					seenIvars[name] = true
					continue
				}

				for i := len(generics) - 1; i >= 0; i-- {
					generic := generics[i]
					ivar = c.replaceTypeParameters(ivar, generic.ArgumentMap)
				}
				ivarStruct := types.InstanceVariable{Type: ivar, Namespace: parent}
				if !yield(name, ivarStruct) {
					return
				}
				seenIvars[name] = true
			}
		}
	}
}

// Get the instance variable with the specified name
func (c *Checker) getInstanceVariable(name value.Symbol) (types.Type, types.Namespace) {
	container, ok := c.selfType.(types.Namespace)
	if !ok {
		return nil, nil
	}

	typ, _ := c.getInstanceVariableIn(name, container)
	return typ, container
}

// Resolve the local with the given name from the current local environment or any parent environment
func (c *Checker) resolveLocal(name string, span *position.Span) (*local, bool) {
	env := c.currentLocalEnv()
	local, inCurrentEnv := env.resolveLocal(name)
	if local == nil {
		c.addFailure(
			fmt.Sprintf("undefined local `%s`", name),
			span,
		)
	}
	return local, inCurrentEnv
}

func (c *Checker) resolveConstantLookupType(node *ast.ConstantLookupNode) (types.Type, string) {
	var leftContainerType types.Type
	var leftContainerName string

	switch l := node.Left.(type) {
	case *ast.PublicConstantNode:
		leftContainerType, leftContainerName = c.resolveType(l.Value, l.Span())
	case *ast.PrivateConstantNode:
		leftContainerType, leftContainerName = c.resolveType(l.Value, l.Span())
	case nil:
		leftContainerType = c.GlobalEnv.Root
	case *ast.ConstantLookupNode:
		leftContainerType, leftContainerName = c.resolveConstantLookupType(l)
	default:
		c.addFailure(
			fmt.Sprintf("invalid type node %T", node),
			node.Span(),
		)
		return nil, ""
	}

	var rightName string
	switch r := node.Right.(type) {
	case *ast.PublicConstantNode:
		rightName = r.Value
	case *ast.PrivateConstantNode:
		rightName = r.Value
		c.addFailure(
			fmt.Sprintf("cannot use private type `%s`", rightName),
			node.Span(),
		)
	default:
		c.addFailure(
			fmt.Sprintf("invalid type node %T", node),
			node.Span(),
		)
		return nil, ""
	}

	typeName := types.MakeFullConstantName(leftContainerName, rightName)
	if leftContainerType == nil {
		return nil, typeName
	}
	leftContainer, ok := leftContainerType.(types.Namespace)
	if !ok {
		c.addFailure(
			fmt.Sprintf("cannot read subtypes from `%s`, it is not a type container", leftContainerName),
			node.Span(),
		)
		return nil, typeName
	}

	subtype, ok := leftContainer.SubtypeString(rightName)
	if !ok {
		c.addFailure(
			fmt.Sprintf("undefined type `%s`", typeName),
			node.Right.Span(),
		)
		return nil, typeName
	}
	if len(subtype.FullName) > 0 {
		typeName = subtype.FullName
	}

	if types.IsNoValue(subtype.Type) || types.IsPlaceholder(subtype.Type) {
		c.addFailure(
			fmt.Sprintf("undefined type `%s`", lexer.Colorize(typeName)),
			node.Right.Span(),
		)
		return nil, typeName
	}

	if !c.checkIfTypeParameterIsAllowed(subtype.Type, node.Right.Span()) {
		return nil, typeName
	}

	if !c.checkTypeIfNecessary(typeName, node.Right.Span()) {
		return types.Nothing{}, typeName
	}
	return subtype.Type, typeName
}

func (c *Checker) checkComplexConstantType(node ast.ExpressionNode) ast.ExpressionNode {
	switch n := node.(type) {
	case *ast.PublicConstantNode:
		return c.checkPublicConstantType(n)
	case *ast.PrivateConstantNode:
		return c.checkPrivateConstantType(n)
	case *ast.ConstantLookupNode:
		return c.constantLookupType(n)
	case *ast.GenericConstantNode:
		c.checkGenericConstantType(n)
		return n
	default:
		c.addFailure(
			fmt.Sprintf("invalid constant type node %T", node),
			node.Span(),
		)
		return n
	}
}

func (c *Checker) addTypeArgumentCountError(name string, paramCount, argCount int, span *position.Span) {
	c.addFailure(
		fmt.Sprintf("`%s` requires %d type argument(s), got: %d", name, paramCount, argCount),
		span,
	)
}

func (c *Checker) checkGenericConstantType(node *ast.GenericConstantNode) (ast.TypeNode, string) {
	constantType, fullName := c.resolveConstantType(node.Constant)
	if constantType == nil {
		node.SetType(types.Nothing{})
		return node, fullName
	}

	switch t := constantType.(type) {
	case *types.GenericNamedType:
		typeArgumentMap, ok := c.checkTypeArguments(
			constantType,
			node.TypeArguments,
			t.TypeParameters,
			node.Constant.Span(),
		)
		if !ok {
			node.SetType(types.Nothing{})
			return node, fullName
		}

		node.SetType(c.replaceTypeParameters(t.Type, typeArgumentMap.ArgumentMap))
		return node, fullName
	case *types.Class:
		typeArgumentMap, ok := c.checkTypeArguments(
			constantType,
			node.TypeArguments,
			t.TypeParameters(),
			node.Constant.Span(),
		)
		if !ok {
			node.SetType(types.Nothing{})
			return node, fullName
		}

		generic := types.NewGeneric(t, typeArgumentMap)
		node.SetType(generic)
		return node, fullName
	case *types.Mixin:
		typeArgumentMap, ok := c.checkTypeArguments(
			constantType,
			node.TypeArguments,
			t.TypeParameters(),
			node.Constant.Span(),
		)
		if !ok {
			node.SetType(types.Nothing{})
			return node, fullName
		}

		generic := types.NewGeneric(t, typeArgumentMap)
		node.SetType(generic)
		return node, fullName
	case *types.Interface:
		typeArgumentMap, ok := c.checkTypeArguments(
			constantType,
			node.TypeArguments,
			t.TypeParameters(),
			node.Constant.Span(),
		)
		if !ok {
			node.SetType(types.Nothing{})
			return node, fullName
		}

		generic := types.NewGeneric(t, typeArgumentMap)
		node.SetType(generic)
		return node, fullName
	case types.Nothing:
		node.SetType(types.Nothing{})
		return node, fullName
	default:
		c.addFailure(
			fmt.Sprintf("type `%s` is not generic", types.InspectWithColor(constantType)),
			node.Constant.Span(),
		)
		node.SetType(types.Nothing{})
		return node, fullName
	}
}

func (c *Checker) checkSimpleConstantType(name string, span *position.Span) types.Type {
	typ, _ := c.resolveType(name, span)
	switch t := typ.(type) {
	case *types.GenericNamedType:
		c.addTypeArgumentCountError(types.InspectWithColor(typ), len(t.TypeParameters), 0, span)
		typ = types.Nothing{}
	case *types.Class:
		if t.IsGeneric() {
			c.addTypeArgumentCountError(types.InspectWithColor(typ), len(t.TypeParameters()), 0, span)
			typ = types.Nothing{}
		}
	case *types.Mixin:
		if t.IsGeneric() {
			c.addTypeArgumentCountError(types.InspectWithColor(typ), len(t.TypeParameters()), 0, span)
			typ = types.Nothing{}
		}
	case *types.Interface:
		if t.IsGeneric() {
			c.addTypeArgumentCountError(types.InspectWithColor(typ), len(t.TypeParameters()), 0, span)
			typ = types.Nothing{}
		}
	case nil:
		typ = types.Nothing{}
	}
	return typ
}

func (c *Checker) checkPublicConstantType(node *ast.PublicConstantNode) *ast.PublicConstantNode {
	typ := c.checkSimpleConstantType(node.Value, node.Span())
	node.SetType(typ)
	return node
}

func (c *Checker) checkPrivateConstantType(node *ast.PrivateConstantNode) *ast.PrivateConstantNode {
	typ := c.checkSimpleConstantType(node.Value, node.Span())
	node.SetType(typ)
	return node
}

func (c *Checker) checkTypeNode(node ast.TypeNode) ast.TypeNode {
	if node.SkipTypechecking() {
		return node
	}
	switch n := node.(type) {
	case *ast.PublicConstantNode:
		c.checkPublicConstantType(n)
		return n
	case *ast.PrivateConstantNode:
		c.checkPrivateConstantType(n)
		return n
	case *ast.GenericConstantNode:
		typeNode, _ := c.checkGenericConstantType(n)
		return typeNode
	case *ast.ConstantLookupNode:
		return c.constantLookupType(n)
	case *ast.ClosureTypeNode:
		return c.checkClosureTypeNode(n)
	case *ast.RawStringLiteralNode:
		n.SetType(types.NewStringLiteral(n.Value))
		return n
	case *ast.DoubleQuotedStringLiteralNode:
		n.SetType(types.NewStringLiteral(n.Value))
		return n
	case *ast.RawCharLiteralNode:
		n.SetType(types.NewCharLiteral(n.Value))
		return n
	case *ast.CharLiteralNode:
		n.SetType(types.NewCharLiteral(n.Value))
		return n
	case *ast.SimpleSymbolLiteralNode:
		n.SetType(types.NewSymbolLiteral(n.Content))
		return n
	case *ast.BinaryTypeExpressionNode:
		return c.checkBinaryTypeExpressionNode(n)
	case *ast.IntLiteralNode:
		n.SetType(types.NewIntLiteral(n.Value))
		return n
	case *ast.Int64LiteralNode:
		n.SetType(types.NewInt64Literal(n.Value))
		return n
	case *ast.Int32LiteralNode:
		n.SetType(types.NewInt32Literal(n.Value))
		return n
	case *ast.Int16LiteralNode:
		n.SetType(types.NewInt16Literal(n.Value))
		return n
	case *ast.Int8LiteralNode:
		n.SetType(types.NewInt8Literal(n.Value))
		return n
	case *ast.UInt64LiteralNode:
		n.SetType(types.NewUInt64Literal(n.Value))
		return n
	case *ast.UInt32LiteralNode:
		n.SetType(types.NewUInt32Literal(n.Value))
		return n
	case *ast.UInt16LiteralNode:
		n.SetType(types.NewUInt16Literal(n.Value))
		return n
	case *ast.UInt8LiteralNode:
		n.SetType(types.NewUInt8Literal(n.Value))
		return n
	case *ast.FloatLiteralNode:
		n.SetType(types.NewFloatLiteral(n.Value))
		return n
	case *ast.Float64LiteralNode:
		n.SetType(types.NewFloat64Literal(n.Value))
		return n
	case *ast.Float32LiteralNode:
		n.SetType(types.NewFloat32Literal(n.Value))
		return n
	case *ast.SelfLiteralNode:
		switch c.mode {
		case closureInferReturnTypeMode, methodMode, outputPositionTypeMode:
			n.SetType(types.Self{})
		default:
			c.addFailure(
				fmt.Sprintf(
					"type `%s` can appear only in method throw types, method return types and method bodies",
					types.InspectWithColor(types.Self{}),
				),
				n.Span(),
			)
			n.SetType(types.Nothing{})
		}
		return n
	case *ast.TrueLiteralNode, *ast.FalseLiteralNode, *ast.VoidTypeNode,
		*ast.NeverTypeNode, *ast.AnyTypeNode, *ast.NilLiteralNode, *ast.BoolLiteralNode:
		return n
	case *ast.NilableTypeNode:
		n.TypeNode = c.checkTypeNode(n.TypeNode)
		typ := c.toNilable(c.typeOf(n.TypeNode))
		n.SetType(typ)
		return n
	case *ast.NotTypeNode:
		n.TypeNode = c.checkTypeNode(n.TypeNode)
		typ := c.normaliseType(types.NewNot(c.typeOf(n.TypeNode)))
		n.SetType(typ)
		return n
	case *ast.SingletonTypeNode:
		return c.checkSingletonTypeNode(n)
	case *ast.InstanceOfTypeNode:
		return c.checkInstanceOfTypeNode(n)
	default:
		c.addFailure(
			fmt.Sprintf("invalid type node %T", node),
			node.Span(),
		)
		return n
	}
}

func (c *Checker) checkSingletonTypeNode(node *ast.SingletonTypeNode) ast.TypeNode {
	node.TypeNode = c.checkTypeNode(node.TypeNode)
	typ := c.typeOf(node.TypeNode)
	var singleton *types.SingletonClass
	switch t := typ.(type) {
	case *types.Class:
		singleton = t.Singleton()
	case *types.Mixin:
		singleton = t.Singleton()
	case *types.Interface:
		singleton = t.Singleton()
	case *types.TypeParameter:
		switch t.UpperBound.(type) {
		case *types.Class, *types.Interface, *types.Mixin:
		default:
			c.addFailure(
				fmt.Sprintf("type parameter `%s` must have an upper bound that is a class, mixin or interface to be used with the singleton type", types.InspectWithColor(typ)),
				node.Span(),
			)
			node.SetType(types.Nothing{})
			return node
		}
		singletonOf := types.NewSingletonOf(t)
		node.SetType(singletonOf)
		return node
	case types.Self:
		switch c.selfType.(type) {
		case *types.Class, *types.Mixin:
		default:
			c.addFailure(
				fmt.Sprintf("type `%s` must be a class or mixin to be used with the singleton type", types.InspectWithColor(c.selfType)),
				node.Span(),
			)
			node.SetType(types.Nothing{})
			return node
		}
		singletonOf := types.NewSingletonOf(t)
		node.SetType(singletonOf)
		return node
	case types.Nothing:
		node.SetType(types.Nothing{})
		return node
	}

	if singleton == nil {
		c.addFailure(
			fmt.Sprintf("cannot get singleton class of `%s`", types.InspectWithColor(typ)),
			node.Span(),
		)
		node.SetType(types.Nothing{})
		return node
	}

	node.SetType(singleton)
	return node
}

func (c *Checker) checkInstanceOfTypeNode(node *ast.InstanceOfTypeNode) ast.TypeNode {
	node.TypeNode = c.checkTypeNode(node.TypeNode)
	typ := c.typeOf(node.TypeNode)
	var namespace types.Namespace

	switch t := typ.(type) {
	case *types.SingletonClass:
		namespace = t.AttachedObject
	case *types.TypeParameter:
		switch upper := t.UpperBound.(type) {
		case *types.SingletonClass:
		case *types.Class:
			switch upper.Name() {
			case "Std::Class", "Std::Mixin", "Std::Interface":
			default:
				c.addFailure(
					fmt.Sprintf("type parameter `%s` must have an upper bound that is a singleton class to be used with the instance of type", types.InspectWithColor(typ)),
					node.Span(),
				)
				node.SetType(types.Nothing{})
				return node
			}
		default:
			c.addFailure(
				fmt.Sprintf("type parameter `%s` must have an upper bound that is a singleton class to be used with the instance of type", types.InspectWithColor(typ)),
				node.Span(),
			)
			node.SetType(types.Nothing{})
			return node
		}
		instanceOf := types.NewInstanceOf(t)
		node.SetType(instanceOf)
		return node
	case types.Self:
		switch c.selfType.(type) {
		case *types.SingletonClass:
		default:
			c.addFailure(
				fmt.Sprintf("type `%s` must be a singleton class to be used with the instance of type", types.InspectWithColor(c.selfType)),
				node.Span(),
			)
			node.SetType(types.Nothing{})
			return node
		}
		instanceOf := types.NewInstanceOf(t)
		node.SetType(instanceOf)
		return node
	case types.Nothing:
		node.SetType(types.Nothing{})
		return node
	}

	if namespace == nil {
		c.addFailure(
			fmt.Sprintf("cannot get instance of `%s`", types.InspectWithColor(typ)),
			node.Span(),
		)
		node.SetType(types.Nothing{})
		return node
	}

	node.SetType(namespace)
	return node
}

func (c *Checker) checkBinaryTypeExpressionNode(node *ast.BinaryTypeExpressionNode) ast.TypeNode {
	switch node.Op.Type {
	case token.OR:
		return c.constructUnionType(node)
	case token.AND:
		return c.constructIntersectionType(node)
	case token.SLASH:
		node.Left = c.checkTypeNode(node.Left)
		node.Right = c.checkTypeNode(node.Right)
		typ := c.differenceType(c.typeOf(node.Left), c.typeOf(node.Right))
		node.SetType(typ)
		return node
	default:
		panic("invalid binary type expression operator")
	}
}

func (c *Checker) checkClosureTypeNode(node *ast.ClosureTypeNode) ast.TypeNode {
	closure := types.NewClosure(nil)
	method, mod := c.declareMethod(
		nil,
		closure,
		"",
		false,
		false,
		false,
		symbol.M_call,
		nil,
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Span(),
	)
	if mod != nil {
		c.popConstScope()
	}
	closure.Body = method
	node.SetType(closure)
	return node
}

func (c *Checker) checkPublicIdentifierNode(node *ast.PublicIdentifierNode) *ast.PublicIdentifierNode {
	local, _ := c.resolveLocal(node.Value, node.Span())
	if local == nil {
		node.SetType(types.Nothing{})
		return node
	}
	if !local.initialised {
		c.addFailure(
			fmt.Sprintf("cannot access uninitialised local `%s`", node.Value),
			node.Span(),
		)
	}
	node.SetType(local.typ)
	return node
}

func (c *Checker) checkPrivateIdentifierNode(node *ast.PrivateIdentifierNode) *ast.PrivateIdentifierNode {
	local, _ := c.resolveLocal(node.Value, node.Span())
	if local == nil {
		node.SetType(types.Nothing{})
		return node
	}
	if !local.initialised {
		c.addFailure(
			fmt.Sprintf("cannot access uninitialised local `%s`", node.Value),
			node.Span(),
		)
	}
	node.SetType(local.typ)
	return node
}

func (c *Checker) checkInstanceVariable(name string, span *position.Span) types.Type {
	typ, container := c.getInstanceVariable(value.ToSymbol(name))
	self, ok := c.selfType.(types.Namespace)
	if !ok || self.IsPrimitive() {
		c.addFailure(
			"cannot use instance variables in this context",
			span,
		)
	}

	if typ == nil {
		c.addFailure(
			fmt.Sprintf(
				"undefined instance variable `%s` in type `%s`",
				types.InspectInstanceVariableWithColor(name),
				types.InspectWithColor(container),
			),
			span,
		)
		return types.Nothing{}
	}

	return typ
}

func (c *Checker) checkInstanceVariableNode(node *ast.InstanceVariableNode) {
	typ := c.checkInstanceVariable(node.Value, node.Span())
	node.SetType(typ)
}

func (c *Checker) declareInstanceVariableForAttribute(name value.Symbol, typ types.Type, span *position.Span) {
	methodNamespace := c.currentMethodScope().container
	currentIvar, ivarNamespace := c.getInstanceVariableIn(name, methodNamespace)

	if currentIvar != nil {
		if !c.isTheSameType(typ, currentIvar, span) {
			c.addFailure(
				fmt.Sprintf(
					"cannot redeclare instance variable `%s` with a different type, is `%s`, should be `%s`, previous definition found in `%s`",
					types.InspectInstanceVariableWithColor(name.String()),
					types.InspectWithColor(typ),
					types.InspectWithColor(currentIvar),
					types.InspectWithColor(ivarNamespace),
				),
				span,
			)
		}
	} else {
		c.declareInstanceVariable(name, typ, span)
	}
}

func (c *Checker) hoistGetterDeclaration(node *ast.GetterDeclarationNode) {
	node.SetType(types.Nothing{})
	for _, entry := range node.Entries {
		attribute, ok := entry.(*ast.AttributeParameterNode)
		if !ok {
			continue
		}

		c.declareMethodForGetter(attribute, node.DocComment())
		c.declareInstanceVariableForAttribute(value.ToSymbol(attribute.Name), c.typeOf(attribute.TypeNode), attribute.Span())
	}
}

func (c *Checker) hoistSetterDeclaration(node *ast.SetterDeclarationNode) {
	node.SetType(types.Nothing{})
	for _, entry := range node.Entries {
		attribute, ok := entry.(*ast.AttributeParameterNode)
		if !ok {
			continue
		}

		c.declareMethodForSetter(attribute, node.DocComment())
		c.declareInstanceVariableForAttribute(value.ToSymbol(attribute.Name), c.typeOf(attribute.TypeNode), attribute.Span())
	}
}

func (c *Checker) hoistAttrDeclaration(node *ast.AttrDeclarationNode) {
	node.SetType(types.Nothing{})
	for _, entry := range node.Entries {
		attribute, ok := entry.(*ast.AttributeParameterNode)
		if !ok {
			continue
		}

		c.declareMethodForSetter(attribute, node.DocComment())
		c.declareMethodForGetter(attribute, node.DocComment())
		c.declareInstanceVariableForAttribute(value.ToSymbol(attribute.Name), c.typeOf(attribute.TypeNode), attribute.Span())
	}
}

func (c *Checker) hoistInstanceVariableDeclaration(node *ast.InstanceVariableDeclarationNode) {
	methodNamespace := c.currentMethodScope().container
	ivar, ivarNamespace := c.getInstanceVariableIn(value.ToSymbol(node.Name), methodNamespace)
	var declaredType types.Type

	if node.TypeNode == nil {
		c.addFailure(
			fmt.Sprintf(
				"cannot declare instance variable `%s` without a type",
				types.InspectInstanceVariableWithColor(node.Name),
			),
			node.Span(),
		)

		declaredType = types.Nothing{}
	} else {
		prevMode := c.mode
		c.mode = instanceVariableMode

		declaredTypeNode := c.checkTypeNode(node.TypeNode)

		c.mode = prevMode

		declaredType = c.typeOf(declaredTypeNode)
		node.TypeNode = declaredTypeNode
		if ivar != nil && !c.isTheSameType(ivar, declaredType, nil) {
			c.addFailure(
				fmt.Sprintf(
					"cannot redeclare instance variable `%s` with a different type, is `%s`, should be `%s`, previous definition found in `%s`",
					types.InspectInstanceVariableWithColor(node.Name),
					types.InspectWithColor(declaredType),
					types.InspectWithColor(ivar),
					types.InspectWithColor(ivarNamespace),
				),
				node.Span(),
			)
			node.SetType(types.Nothing{})
			return
		}
	}

	switch c.mode {
	case mixinMode, classMode, moduleMode, singletonMode:
	default:
		c.addFailure(
			fmt.Sprintf(
				"cannot declare instance variable `%s` in this context",
				types.InspectInstanceVariableWithColor(node.Name),
			),
			node.Span(),
		)
		node.SetType(types.Nothing{})
		return
	}

	node.SetType(declaredType)
	c.declareInstanceVariable(value.ToSymbol(node.Name), declaredType, node.Span())
}

func (c *Checker) checkVariableDeclaration(
	name string,
	initialiser ast.ExpressionNode,
	typeNode ast.TypeNode,
	span *position.Span,
) (
	_init ast.ExpressionNode,
	_typeNode ast.TypeNode,
	typ types.Type,
) {
	if variable := c.getLocal(name); variable != nil {
		c.addFailure(
			fmt.Sprintf("cannot redeclare local `%s`", name),
			span,
		)
	}
	if initialiser == nil {
		if typeNode == nil {
			c.addFailure(
				fmt.Sprintf("cannot declare a variable without a type `%s`", name),
				span,
			)
			c.addLocal(name, newLocal(types.Nothing{}, false, false))
			return initialiser, typeNode, types.Nothing{}
		}

		// without an initialiser but with a type
		declaredTypeNode := c.checkTypeNode(typeNode)
		declaredType := c.typeOf(declaredTypeNode)
		c.addLocal(name, newLocal(declaredType, false, false))
		return initialiser, declaredTypeNode, types.Nothing{}
	}

	// with an initialiser
	if typeNode == nil {
		// without a type, inference
		init := c.checkExpression(initialiser)
		actualType := c.toNonLiteral(c.typeOfGuardVoid(init), false)
		c.addLocal(name, newLocal(actualType, true, false))
		return init, nil, actualType
	}

	// with a type and an initializer

	declaredTypeNode := c.checkTypeNode(typeNode)
	declaredType := c.typeOf(declaredTypeNode)
	init := c.checkExpressionWithType(initialiser, declaredType)
	actualType := c.typeOfGuardVoid(init)
	c.addLocal(name, newLocal(declaredType, true, false))
	c.checkCanAssign(actualType, declaredType, init.Span())

	return init, declaredTypeNode, declaredType
}

func (c *Checker) checkVariablePatternDeclarationNode(node *ast.VariablePatternDeclarationNode) *ast.VariablePatternDeclarationNode {
	node.Initialiser = c.checkExpression(node.Initialiser)
	initType := c.typeOf(node.Initialiser)

	prevMode := c.mode
	c.mode = variablePatternMode
	c.checkPattern(node.Pattern, initType)
	c.mode = prevMode
	return node
}

func (c *Checker) checkPattern(node ast.PatternNode, typ types.Type) {
	switch n := node.(type) {
	case *ast.PublicIdentifierNode:
		c.checkIdentifierPattern(n.Value, typ, n.Span())
	case *ast.PrivateIdentifierNode:
		c.checkIdentifierPattern(n.Value, typ, n.Span())
	case *ast.RestPatternNode:
		c.checkPattern(n.Identifier, typ)
	case *ast.ListPatternNode:
		typeGeneric, ok := typ.(*types.Generic)
		if !ok || !c.isExplicitSubtypeOfInterface(typeGeneric, c.Std(symbol.List).(*types.Interface)) || typeGeneric.TypeArguments.Len() != 1 {
			c.addFailure(
				fmt.Sprintf("type `%s` cannot be matched against a list pattern", types.InspectWithColor(typ)),
				node.Span(),
			)
			return
		}
		typeArg := typeGeneric.TypeArguments.Get(0)
		for _, element := range n.Elements {
			c.checkPattern(element, typeArg.Type)
		}
	}
}

func (c *Checker) checkIdentifierPattern(name string, typ types.Type, span *position.Span) {
	variable := c.getLocal(name)
	if variable == nil {
		c.addLocal(name, newLocal(typ, true, c.mode == valuePatternMode))
		return
	}

	variable.initialised = true
	c.checkCanAssign(typ, variable.typ, span)
}

func (c *Checker) checkVariableDeclarationNode(node *ast.VariableDeclarationNode) {
	init, typeNode, typ := c.checkVariableDeclaration(node.Name, node.Initialiser, node.TypeNode, node.Span())
	node.Initialiser = init
	node.TypeNode = typeNode
	node.SetType(typ)
}

func (c *Checker) checkValueDeclarationNode(node *ast.ValueDeclarationNode) {
	if variable := c.getLocal(node.Name); variable != nil {
		c.addFailure(
			fmt.Sprintf("cannot redeclare local `%s`", node.Name),
			node.Span(),
		)
	}
	if node.Initialiser == nil {
		if node.TypeNode == nil {
			c.addFailure(
				fmt.Sprintf("cannot declare a value without a type `%s`", node.Name),
				node.Span(),
			)
			c.addLocal(node.Name, newLocal(types.Nothing{}, false, true))
			node.SetType(types.Nothing{})
			return
		}

		// without an initialiser but with a type
		declaredTypeNode := c.checkTypeNode(node.TypeNode)
		declaredType := c.typeOf(declaredTypeNode)
		c.addLocal(node.Name, newLocal(declaredType, false, true))
		node.TypeNode = declaredTypeNode
		node.SetType(types.Nothing{})
		return
	}

	// with an initialiser
	if node.TypeNode == nil {
		// without a type, inference
		init := c.checkExpression(node.Initialiser)
		actualType := c.typeOfGuardVoid(init)
		c.addLocal(node.Name, newLocal(actualType, true, true))
		node.Initialiser = init
		node.SetType(actualType)
		return
	}

	// with a type and an initializer

	declaredTypeNode := c.checkTypeNode(node.TypeNode)
	declaredType := c.typeOf(declaredTypeNode)
	init := c.checkExpressionWithType(node.Initialiser, declaredType)
	actualType := c.typeOfGuardVoid(init)
	c.addLocal(node.Name, newLocal(declaredType, true, true))
	c.checkCanAssign(actualType, declaredType, init.Span())

	node.TypeNode = declaredTypeNode
	node.Initialiser = init
	node.SetType(declaredType)
}

func extractConstantName(node ast.ExpressionNode) string {
	switch n := node.(type) {
	case *ast.PublicConstantNode:
		return n.Value
	case *ast.PrivateConstantNode:
		return n.Value
	case *ast.ConstantLookupNode:
		return extractConstantNameFromLookup(n)
	default:
		panic(fmt.Sprintf("invalid constant node: %T", node))
	}
}

func extractConstantNameFromLookup(lookup *ast.ConstantLookupNode) string {
	switch r := lookup.Right.(type) {
	case *ast.PublicConstantNode:
		return r.Value
	case *ast.PrivateConstantNode:
		return r.Value
	default:
		panic(fmt.Sprintf("invalid right side of constant lookup node: %T", lookup.Right))
	}
}

func (c *Checker) declareModule(docComment string, namespace types.Namespace, constantType types.Type, fullConstantName string, constantName value.Symbol, span *position.Span) *types.Module {
	if constantType != nil {
		ct, ok := constantType.(*types.SingletonClass)
		if ok {
			constantType = ct.AttachedObject
		}

		switch t := constantType.(type) {
		case *types.Module:
			t.AppendDocComment(docComment)
			return t
		case *types.NamespacePlaceholder:
			module := types.NewModuleWithDetails(
				docComment,
				t.Name(),
				t.Constants(),
				t.Subtypes(),
				t.Methods(),
				c.GlobalEnv,
			)
			t.Namespace = module
			namespace.DefineConstant(constantName, module)
			namespace.DefineSubtype(constantName, module)
			return module
		case *types.Placeholder:
			module := types.NewModule(
				docComment,
				fullConstantName,
				c.GlobalEnv,
			)
			c.replaceSimpleNamespacePlaceholder(t, module, module, constantName)
			namespace.DefineConstant(constantName, module)
			namespace.DefineSubtype(constantName, module)
			return module
		default:
			c.addRedeclaredConstantError(fullConstantName, span)
			return types.NewModule(docComment, fullConstantName, c.GlobalEnv)
		}
	}

	if namespace == nil {
		return types.NewModule(docComment, fullConstantName, c.GlobalEnv)
	}

	return namespace.DefineModule(docComment, constantName, c.GlobalEnv)
}

func (c *Checker) declareInstanceVariable(name value.Symbol, typ types.Type, errSpan *position.Span) {
	container := c.currentConstScope().container
	if container.IsPrimitive() {
		c.addFailure(
			fmt.Sprintf("cannot declare instance variable `%s` in a primitive `%s`", name, types.InspectWithColor(container)),
			errSpan,
		)
	}
	container.DefineInstanceVariable(name, typ)
}

func (c *Checker) declareClass(docComment string, abstract, sealed, primitive bool, namespace types.Namespace, constantType types.Type, fullConstantName string, constantName value.Symbol, span *position.Span) *types.Class {
	if constantType != nil {
		switch ct := constantType.(type) {
		case *types.SingletonClass:
			constantType = ct.AttachedObject
		case *types.Placeholder:
			class := types.NewClass(
				docComment,
				abstract,
				sealed,
				primitive,
				fullConstantName,
				nil,
				c.GlobalEnv,
			)
			classSingleton := class.Singleton()
			c.replaceSimpleNamespacePlaceholder(ct, class, classSingleton, constantName)
			namespace.DefineConstant(constantName, class.Singleton())
			namespace.DefineSubtype(constantName, class)
			return class
		default:
			c.addRedeclaredConstantError(fullConstantName, span)
			return types.NewClass(
				docComment,
				abstract,
				sealed,
				primitive,
				fullConstantName,
				nil,
				c.GlobalEnv,
			)
		}

		switch t := constantType.(type) {
		case *types.Class:
			if abstract != t.IsAbstract() || sealed != t.IsSealed() || primitive != t.IsPrimitive() {
				c.addFailure(
					fmt.Sprintf(
						"cannot redeclare class `%s` with a different modifier, is `%s`, should be `%s`",
						fullConstantName,
						types.InspectModifier(abstract, sealed, primitive),
						types.InspectModifier(t.IsAbstract(), t.IsSealed(), t.IsPrimitive()),
					),
					span,
				)
			}
			t.AppendDocComment(docComment)
			return t
		case *types.NamespacePlaceholder:
			class := types.NewClassWithDetails(
				docComment,
				abstract,
				sealed,
				primitive,
				t.Name(),
				nil,
				t.Constants(),
				t.Subtypes(),
				t.Methods(),
				c.GlobalEnv,
			)
			t.Namespace = class
			namespace.DefineConstant(constantName, class.Singleton())
			namespace.DefineSubtype(constantName, class)
			return class
		default:
			c.addRedeclaredConstantError(fullConstantName, span)
			return types.NewClass(
				docComment,
				abstract,
				sealed,
				primitive,
				fullConstantName,
				nil,
				c.GlobalEnv,
			)
		}
	}

	if namespace == nil {
		return types.NewClass(
			docComment,
			abstract,
			sealed,
			primitive,
			fullConstantName,
			nil,
			c.GlobalEnv,
		)
	}

	return namespace.DefineClass(
		docComment,
		abstract,
		sealed,
		primitive,
		constantName,
		nil,
		c.GlobalEnv,
	)
}

func (c *Checker) checkCanAssign(assignedType types.Type, targetType types.Type, span *position.Span) {
	if !c.isSubtype(assignedType, targetType, span) {
		c.addFailure(
			fmt.Sprintf(
				"type `%s` cannot be assigned to type `%s`",
				types.InspectWithColor(assignedType),
				types.InspectWithColor(targetType),
			),
			span,
		)
	}
}

func (c *Checker) checkCanAssignInstanceVariable(name string, assignedType types.Type, targetType types.Type, span *position.Span) {
	if !c.isSubtype(assignedType, targetType, span) {
		c.addFailure(
			fmt.Sprintf(
				"type `%s` cannot be assigned to instance variable `%s` of type `%s`",
				types.InspectWithColor(assignedType),
				types.InspectInstanceVariableWithColor(name),
				types.InspectWithColor(targetType),
			),
			span,
		)
	}
}

func (c *Checker) hoistStructDeclaration(structNode *ast.StructDeclarationNode) ast.ExpressionNode {
	switch c.mode {
	case topLevelMode, classMode, interfaceMode,
		moduleMode, mixinMode:
	default:
		return structNode
	}

	container, constant, fullConstantName := c.resolveConstantForDeclaration(structNode.Constant)
	constantName := value.ToSymbol(extractConstantName(structNode.Constant))
	class := c.declareClass(
		structNode.DocComment(),
		false,
		false,
		false,
		container,
		constant,
		fullConstantName,
		constantName,
		structNode.Span(),
	)
	structNode.SetType(class)
	structNode.Constant = ast.NewPublicConstantNode(structNode.Constant.Span(), fullConstantName)

	init := ast.NewInitDefinitionNode(
		structNode.Span(),
		nil,
		nil,
		nil,
	)
	attrDeclaration := ast.NewAttrDeclarationNode(
		structNode.Span(),
		"",
		nil,
	)
	newStatements := []ast.StatementNode{
		ast.ExpressionToStatement(init),
		ast.ExpressionToStatement(attrDeclaration),
	}

	var optionalParamSeen bool

	for _, stmt := range structNode.Body {
		switch s := stmt.(type) {
		case *ast.ParameterStatementNode:
			param := s.Parameter.(*ast.AttributeParameterNode)
			if optionalParamSeen && param.Initialiser == nil {
				c.addFailure(
					fmt.Sprintf(
						"required struct attribute `%s` cannot appear after optional attributes",
						param.Name,
					),
					param.Span(),
				)
			}
			if param.Initialiser != nil {
				optionalParamSeen = true
			}
			init.Parameters = append(
				init.Parameters,
				ast.NewMethodParameterNode(
					param.Span(),
					param.Name,
					true,
					param.TypeNode,
					param.Initialiser,
					ast.NormalParameterKind,
				),
			)
			attrDeclaration.Entries = append(
				attrDeclaration.Entries,
				ast.NewAttributeParameterNode(
					param.Span(),
					param.Name,
					param.TypeNode,
					nil,
				),
			)
		}
	}

	classNode := ast.NewClassDeclarationNode(
		structNode.Span(),
		structNode.DocComment(),
		false,
		false,
		false,
		structNode.Constant,
		nil,
		nil,
		newStatements,
	)
	classNode.SetType(class)
	c.registerNamespaceDeclarationCheck(fullConstantName, classNode, class)
	return classNode
}

func (c *Checker) hoistModuleDeclaration(node *ast.ModuleDeclarationNode) {
	switch c.mode {
	case topLevelMode, classMode, interfaceMode,
		moduleMode, mixinMode:
	default:
		return
	}

	container, constant, fullConstantName := c.resolveConstantForDeclaration(node.Constant)
	constantName := value.ToSymbol(extractConstantName(node.Constant))
	module := c.declareModule(
		node.DocComment(),
		container,
		constant,
		fullConstantName,
		constantName,
		node.Span(),
	)
	node.SetType(module)
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)

	prevMode := c.mode
	c.mode = moduleMode
	c.pushConstScope(makeLocalConstantScope(module))
	c.pushMethodScope(makeLocalMethodScope(module))

	c.hoistNamespaceDefinitions(node.Body)

	c.popLocalConstScope()
	c.popMethodScope()
	c.mode = prevMode
}

func (c *Checker) hoistClassDeclaration(node *ast.ClassDeclarationNode) {
	switch c.mode {
	case topLevelMode, classMode, interfaceMode,
		moduleMode, mixinMode:
	default:
		return
	}

	container, constant, fullConstantName := c.resolveConstantForDeclaration(node.Constant)
	constantName := value.ToSymbol(extractConstantName(node.Constant))
	class := c.declareClass(
		node.DocComment(),
		node.Abstract,
		node.Sealed,
		node.Primitive,
		container,
		constant,
		fullConstantName,
		constantName,
		node.Span(),
	)
	node.SetType(class)
	c.registerNamespaceDeclarationCheck(fullConstantName, node, class)
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)

	prevMode := c.mode
	c.mode = classMode
	c.pushConstScope(makeLocalConstantScope(class))
	c.pushMethodScope(makeLocalMethodScope(class))

	c.hoistNamespaceDefinitions(node.Body)

	c.popLocalConstScope()
	c.popMethodScope()
	c.mode = prevMode
}

func (c *Checker) hoistMixinDeclaration(node *ast.MixinDeclarationNode) {
	switch c.mode {
	case topLevelMode, classMode, interfaceMode,
		moduleMode, mixinMode:
	default:
		return
	}

	container, constant, fullConstantName := c.resolveConstantForDeclaration(node.Constant)
	constantName := value.ToSymbol(extractConstantName(node.Constant))
	mixin := c.declareMixin(
		node.DocComment(),
		node.Abstract,
		container,
		constant,
		fullConstantName,
		constantName,
		node.Span(),
	)
	node.SetType(mixin)
	c.registerNamespaceDeclarationCheck(fullConstantName, node, mixin)
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)

	prevMode := c.mode
	c.mode = mixinMode
	c.pushConstScope(makeLocalConstantScope(mixin))
	c.pushMethodScope(makeLocalMethodScope(mixin))

	c.hoistNamespaceDefinitions(node.Body)

	c.popLocalConstScope()
	c.popMethodScope()
	c.mode = prevMode
}

func (c *Checker) hoistInterfaceDeclaration(node *ast.InterfaceDeclarationNode) {
	switch c.mode {
	case topLevelMode, classMode, interfaceMode,
		moduleMode, mixinMode:
	default:
		return
	}

	container, constant, fullConstantName := c.resolveConstantForDeclaration(node.Constant)
	constantName := value.ToSymbol(extractConstantName(node.Constant))
	iface := c.declareInterface(
		node.DocComment(),
		container,
		constant,
		fullConstantName,
		constantName,
		node.Span(),
	)
	node.SetType(iface)
	c.registerNamespaceDeclarationCheck(fullConstantName, node, iface)
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)

	prevMode := c.mode
	c.mode = interfaceMode
	c.pushConstScope(makeLocalConstantScope(iface))
	c.pushMethodScope(makeLocalMethodScope(iface))

	c.hoistNamespaceDefinitions(node.Body)

	c.popLocalConstScope()
	c.popMethodScope()
	c.mode = prevMode
}

func (c *Checker) hoistSingletonDeclaration(node *ast.SingletonBlockExpressionNode) {
	switch c.mode {
	case classMode, mixinMode, interfaceMode:
	default:
		return
	}

	currentNamespace := c.currentConstScope().container
	singleton := currentNamespace.Singleton()
	node.SetType(singleton)

	prevMode := c.mode
	c.mode = singletonMode
	c.pushConstScope(makeLocalConstantScope(singleton))
	c.pushMethodScope(makeLocalMethodScope(singleton))

	c.hoistNamespaceDefinitions(node.Body)

	c.popLocalConstScope()
	c.popMethodScope()
	c.mode = prevMode
}

func (c *Checker) hoistImports(statements []ast.StatementNode) []*ast.ImportStatementNode {
	var imports []*ast.ImportStatementNode

	for _, statement := range statements {
		switch stmt := statement.(type) {
		case *ast.ImportStatementNode:
			imports = append(imports, stmt)
			c.checkImport(stmt)
		}
	}

	return imports
}

func (c *Checker) hoistNamespaceDefinitions(statements []ast.StatementNode) {
	for _, statement := range statements {
		switch stmt := statement.(type) {
		case *ast.ExpressionStatementNode:
			expression := stmt.Expression

			switch expr := expression.(type) {
			case *ast.StructDeclarationNode:
				stmt.Expression = c.hoistStructDeclaration(expr)
			case *ast.ModuleDeclarationNode:
				c.hoistModuleDeclaration(expr)
			case *ast.ClassDeclarationNode:
				c.hoistClassDeclaration(expr)
			case *ast.MixinDeclarationNode:
				c.hoistMixinDeclaration(expr)
			case *ast.InterfaceDeclarationNode:
				c.hoistInterfaceDeclaration(expr)
			case *ast.SingletonBlockExpressionNode:
				c.hoistSingletonDeclaration(expr)
			case *ast.TypeDefinitionNode:
				c.registerNamedTypeCheck(expr)
			case *ast.GenericTypeDefinitionNode:
				c.registerGenericNamedTypeCheck(expr)
			case *ast.UsingExpressionNode:
				c.checkUsingExpressionForNamespaces(expr)
			case *ast.ImplementExpressionNode:
				switch c.mode {
				case classMode, mixinMode, interfaceMode:
				default:
					continue
				}
				namespace := c.currentMethodScope().container
				expr.SetType(types.Nothing{})
				c.registerNamespaceDeclarationCheck(namespace.Name(), expr, namespace)
			case *ast.IncludeExpressionNode:
				switch c.mode {
				case classMode, mixinMode, singletonMode:
				default:
					continue
				}
				namespace := c.currentMethodScope().container
				expr.SetType(types.Nothing{})
				c.registerNamespaceDeclarationCheck(namespace.Name(), expr, namespace)
			case *ast.InstanceVariableDeclarationNode, *ast.GetterDeclarationNode,
				*ast.SetterDeclarationNode, *ast.AttrDeclarationNode:
				namespace := c.currentMethodScope().container
				namespace.DefineInstanceVariable(symbol.Empty, nil) // placeholder
			}
		}
	}
}

func (c *Checker) checkUsingExpressionForNamespaces(node *ast.UsingExpressionNode) {
	node.SetType(types.Nothing{})
	for i, entry := range node.Entries {
		node.Entries[i] = c.checkUsingEntryNodeForNamespaces(entry)
	}
}

func (c *Checker) checkUsingEntryNodeForNamespaces(node ast.UsingEntryNode) ast.UsingEntryNode {
	switch n := node.(type) {
	case *ast.UsingAllEntryNode:
		return c.checkUsingAllEntryNode(n)
	case *ast.PublicConstantNode, *ast.PrivateConstantNode:
		c.addFailure(
			"this using statement will have no effect",
			node.Span(),
		)
		return node
	case *ast.ConstantLookupNode:
		return c.checkUsingConstantLookupEntryNodeForNamespace(n)
	default:
		panic(fmt.Sprintf("invalid using entry node: %T", node))
	}
}

func (c *Checker) checkUsingConstantLookupEntryNodeForNamespace(node *ast.ConstantLookupNode) ast.UsingEntryNode {
	container, constant, fullConstantName, constantName := c.resolveTypeLookupInRoot(node)
	constantSymbol := value.ToSymbol(constantName)
	node.Left.SetType(container)
	usingNamespace := c.getUsingBufferNamespace()
	switch n := constant.(type) {
	case *types.Class:
		usingNamespace.DefineSubtypeWithFullName(constantSymbol, fullConstantName, constant)
		usingNamespace.DefineConstantWithFullName(constantSymbol, fullConstantName, n.Singleton())
		return node
	case *types.Mixin:
		usingNamespace.DefineSubtypeWithFullName(constantSymbol, fullConstantName, constant)
		usingNamespace.DefineConstantWithFullName(constantSymbol, fullConstantName, n.Singleton())
		return node
	case *types.Interface:
		usingNamespace.DefineSubtypeWithFullName(constantSymbol, fullConstantName, constant)
		usingNamespace.DefineConstantWithFullName(constantSymbol, fullConstantName, n.Singleton())
		return node
	case *types.Module:
		usingNamespace.DefineSubtypeWithFullName(constantSymbol, fullConstantName, constant)
		usingNamespace.DefineConstantWithFullName(constantSymbol, fullConstantName, constant)
		return node
	default:
		usingNamespace.DefineSubtypeWithFullName(constantSymbol, fullConstantName, constant)
		return node
	case nil: // continue
	}

	placeholderType := types.NewPlaceholder(fullConstantName, usingNamespace.Subtypes(), c.newLocation(node.Span()))
	c.registerPlaceholder(placeholderType)

	placeholderConstant := types.NewPlaceholder(fullConstantName, usingNamespace.Constants(), c.newLocation(node.Span()))
	c.registerPlaceholder(placeholderConstant)

	placeholderConstant.Sibling = placeholderType
	placeholderType.Sibling = placeholderConstant

	container.DefineSubtype(constantSymbol, placeholderType)
	usingNamespace.DefineSubtypeWithFullName(constantSymbol, fullConstantName, placeholderType)

	container.DefineConstant(constantSymbol, placeholderConstant)
	usingNamespace.DefineConstantWithFullName(constantSymbol, fullConstantName, placeholderConstant)
	node.SetType(usingNamespace)

	return node
}

func (c *Checker) checkSimpleUsingEntry(typ types.Type, constName, fullName string, parentNamespace types.Namespace, span *position.Span) types.Namespace {
	if typ != nil {
		var namespace types.Namespace
		switch t := typ.(type) {
		case *types.Module:
			namespace = t
		case *types.Mixin:
			namespace = t
		case *types.Class:
			namespace = t
		case *types.Interface:
			namespace = t
		case *types.NamespacePlaceholder:
			namespace = t
			t.Locations.Append(c.newLocation(span))
		default:
			c.addFailure(
				fmt.Sprintf("type `%s` is not a namespace", types.InspectWithColor(typ)),
				span,
			)
		}

		return namespace
	}

	placeholder := types.NewNamespacePlaceholder(fullName)
	placeholder.Locations.Append(c.newLocation(span))
	c.registerPlaceholderNamespace(placeholder)
	parentNamespace.DefineSubtype(value.ToSymbol(constName), placeholder)
	parentNamespace.DefineConstant(value.ToSymbol(constName), types.NewSingletonClass(placeholder, nil))
	return placeholder
}

// func (c *Checker) checkUsingConstantLookupEntryNode(node *ast.ConstantLookupNode) ast.UsingEntryNode {
// 	parentNamespace, typ, fullName, constName := c.resolveTypeInRoot(node)
// 	newNode := ast.NewPublicConstantNode(node.Span(), fullName)
// 	namespace := c.checkSimpleUsingEntry(typ, constName, fullName, parentNamespace, node.Span())
// 	node.SetType(namespace)
// 	c.pushConstScope(makeUsingConstantScope(namespace))
// 	singleton := namespace.Singleton()
// 	if singleton == nil {
// 		c.pushMethodScope(makeUsingMethodScope(namespace))
// 	} else {
// 		c.pushMethodScope(makeUsingMethodScope(singleton))
// 	}

// 	return newNode
// }

func (c *Checker) checkUsingAllEntryNode(node *ast.UsingAllEntryNode) *ast.UsingAllEntryNode {
	parentNamespace, typ, fullName, constName := c.resolveTypeInRoot(node.Namespace)
	node.Namespace = ast.NewPublicConstantNode(node.Namespace.Span(), fullName)
	namespace := c.checkSimpleUsingEntry(typ, constName, fullName, parentNamespace, node.Span())
	if namespace == nil {
		return node
	}
	node.SetType(namespace)
	c.pushConstScope(makeUsingConstantScope(namespace))
	singleton := namespace.Singleton()
	if singleton == nil {
		c.pushMethodScope(makeUsingMethodScope(namespace))
	} else {
		c.pushMethodScope(makeUsingMethodScope(singleton))
	}
	return node
}

func (c *Checker) checkImport(node *ast.ImportStatementNode) {
	var path string

	switch pathNode := node.Path.(type) {
	case *ast.DoubleQuotedStringLiteralNode:
		path = pathNode.Value
	case *ast.RawStringLiteralNode:
		path = pathNode.Value
	}

	if !filepath.IsAbs(path) {
		dirPath := filepath.Dir(c.Filename)
		path = filepath.Join(dirPath, path)
	}

	filePaths, err := filepath.Glob(path)
	if err != nil {
		c.addFailure(
			fmt.Sprintf(
				"invalid glob pattern: %s (%s)",
				path,
				err,
			),
			node.Span(),
		)
		return
	}

	for _, filePath := range filePaths {
		if !fileExists(filePath) {
			c.addFailure(
				fmt.Sprintf(
					"cannot find file: %s",
					path,
				),
				node.Span(),
			)
			continue
		}
		node.FsPaths = append(node.FsPaths, filePath)
	}
}

func fileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (c *Checker) hoistInitDefinition(initNode *ast.InitDefinitionNode) *ast.MethodDefinitionNode {
	switch c.mode {
	case classMode:
	default:
		c.addFailure(
			"init definitions cannot appear outside of classes",
			initNode.Span(),
		)
	}
	method, mod := c.declareMethod(
		nil,
		c.currentMethodScope().container,
		initNode.DocComment(),
		false,
		false,
		false,
		symbol.M_init,
		nil,
		initNode.Parameters,
		nil,
		initNode.ThrowType,
		initNode.Span(),
	)
	initNode.SetType(method)
	newNode := ast.NewMethodDefinitionNode(
		initNode.Span(),
		initNode.DocComment(),
		false,
		false,
		"#init",
		nil,
		initNode.Parameters,
		nil,
		initNode.ThrowType,
		initNode.Body,
	)
	newNode.SetType(method)
	c.registerMethodCheck(method, newNode)
	if mod != nil {
		c.popConstScope()
	}
	return newNode
}

func (c *Checker) hoistAliasDeclaration(node *ast.AliasDeclarationNode) {
	node.SetType(types.Nothing{})
	namespace := c.currentMethodScope().container
	for _, entry := range node.Entries {
		method := c.resolveMethodInNamespace(namespace, value.ToSymbol(entry.OldName))
		if method == nil {
			c.addMissingMethodError(namespace, entry.OldName, entry.Span())
			continue
		}
		namespace.SetMethod(value.ToSymbol(entry.NewName), method)
	}
}

func (c *Checker) hoistMethodDefinition(node *ast.MethodDefinitionNode) {
	method, mod := c.declareMethod(
		nil,
		c.currentMethodScope().container,
		node.DocComment(),
		node.Abstract,
		node.Sealed,
		false,
		value.ToSymbol(node.Name),
		node.TypeParameters,
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Span(),
	)
	node.SetType(method)
	c.registerMethodCheck(method, node)
	if mod != nil {
		c.popConstScope()
	}
}

func (c *Checker) hoistMethodSignatureDefinition(node *ast.MethodSignatureDefinitionNode) {
	method, mod := c.declareMethod(
		nil,
		c.currentMethodScope().container,
		node.DocComment(),
		true,
		false,
		false,
		value.ToSymbol(node.Name),
		node.TypeParameters,
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Span(),
	)
	if mod != nil {
		c.popConstScope()
	}
	node.SetType(method)
}

func (c *Checker) hoistMethodDefinitionsWithinClass(node *ast.ClassDeclarationNode) {
	class, ok := c.typeOf(node).(*types.Class)
	if ok {
		c.pushConstScope(makeLocalConstantScope(class))
		c.pushMethodScope(makeLocalMethodScope(class))
	}

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = classMode
	c.selfType = class
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)
	c.selfType = previousSelf
	if ok {
		c.popLocalConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) hoistMethodDefinitionsWithinModule(node *ast.ModuleDeclarationNode) {
	module, ok := c.typeOf(node).(*types.Module)
	if ok {
		c.pushConstScope(makeLocalConstantScope(module))
		c.pushMethodScope(makeLocalMethodScope(module))
	}

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = moduleMode
	c.selfType = module
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)
	c.selfType = previousSelf

	if ok {
		c.popLocalConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) hoistMethodDefinitionsWithinMixin(node *ast.MixinDeclarationNode) {
	mixin, ok := c.typeOf(node).(*types.Mixin)
	if ok {
		c.pushConstScope(makeLocalConstantScope(mixin))
		c.pushMethodScope(makeLocalMethodScope(mixin))
	}

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = mixinMode
	c.selfType = mixin
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)
	c.selfType = previousSelf

	if ok {
		c.popLocalConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) hoistMethodDefinitionsWithinInterface(node *ast.InterfaceDeclarationNode) {
	iface, ok := c.typeOf(node).(*types.Interface)
	if ok {
		c.pushConstScope(makeLocalConstantScope(iface))
		c.pushMethodScope(makeLocalMethodScope(iface))
	}

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = interfaceMode
	c.selfType = iface
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)
	c.selfType = previousSelf

	if ok {
		c.popLocalConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) hoistMethodDefinitionsWithinSingleton(expr *ast.SingletonBlockExpressionNode) {
	namespace := c.currentConstScope().container
	singleton := namespace.Singleton()

	c.pushConstScope(makeLocalConstantScope(singleton))
	c.pushMethodScope(makeLocalMethodScope(singleton))

	previousMode := c.mode
	previousSelf := c.selfType
	c.mode = singletonMode
	c.selfType = singleton
	c.hoistMethodDefinitions(expr.Body)
	c.setMode(previousMode)
	c.selfType = previousSelf

	c.popLocalConstScope()
	c.popMethodScope()
}

func (c *Checker) hoistMethodDefinitions(statements []ast.StatementNode) {
	for _, statement := range statements {
		stmt, ok := statement.(*ast.ExpressionStatementNode)
		if !ok {
			continue
		}

		expression := stmt.Expression

		switch expr := expression.(type) {
		case *ast.ConstantDeclarationNode:
			c.hoistConstantDeclaration(expr)
		case *ast.AliasDeclarationNode:
			c.hoistAliasDeclaration(expr)
		case *ast.MethodDefinitionNode:
			c.hoistMethodDefinition(expr)
		case *ast.MethodSignatureDefinitionNode:
			c.hoistMethodSignatureDefinition(expr)
		case *ast.InitDefinitionNode:
			stmt.Expression = c.hoistInitDefinition(expr)
		case *ast.InstanceVariableDeclarationNode:
			c.hoistInstanceVariableDeclaration(expr)
		case *ast.GetterDeclarationNode:
			c.hoistGetterDeclaration(expr)
		case *ast.SetterDeclarationNode:
			c.hoistSetterDeclaration(expr)
		case *ast.AttrDeclarationNode:
			c.hoistAttrDeclaration(expr)
		case *ast.UsingExpressionNode:
			c.resolveUsingExpression(expr)
		case *ast.ModuleDeclarationNode:
			c.hoistMethodDefinitionsWithinModule(expr)
		case *ast.ClassDeclarationNode:
			c.hoistMethodDefinitionsWithinClass(expr)
		case *ast.MixinDeclarationNode:
			c.hoistMethodDefinitionsWithinMixin(expr)
		case *ast.InterfaceDeclarationNode:
			c.hoistMethodDefinitionsWithinInterface(expr)
		case *ast.SingletonBlockExpressionNode:
			c.hoistMethodDefinitionsWithinSingleton(expr)
		}
	}
}

func (c *Checker) resolveUsingExpression(node *ast.UsingExpressionNode) {
	for _, constNode := range node.Entries {
		typ := c.typeOf(constNode)
		switch t := typ.(type) {
		case *types.Module:
			c.pushConstScope(makeUsingConstantScope(t))
			c.pushMethodScope(makeUsingMethodScope(t))
		case *types.Mixin:
			c.pushConstScope(makeUsingConstantScope(t))
			c.pushMethodScope(makeUsingMethodScope(t))
		case *types.Class:
			c.pushConstScope(makeUsingConstantScope(t))
			c.pushMethodScope(makeUsingMethodScope(t))
		case *types.Interface:
			c.pushConstScope(makeUsingConstantScope(t))
			c.pushMethodScope(makeUsingMethodScope(t))
		case *types.NamespacePlaceholder:
			c.pushConstScope(makeUsingConstantScope(t))
			c.pushMethodScope(makeUsingMethodScope(t))
			t.Locations.Append(c.newLocation(constNode.Span()))
		case *types.UsingBufferNamespace:
			if c.enclosingScopeIsAUsingBuffer() {
				continue
			}
			c.pushConstScope(makeUsingBufferConstantScope(t))
			c.pushMethodScope(makeUsingBufferMethodScope(t))
		}
	}
}

func (c *Checker) declareMixin(docComment string, abstract bool, namespace types.Namespace, constantType types.Type, fullConstantName string, constantName value.Symbol, span *position.Span) *types.Mixin {
	if constantType != nil {
		switch ct := constantType.(type) {
		case *types.SingletonClass:
			constantType = ct.AttachedObject
		case *types.Placeholder:
			mixin := types.NewMixin(
				docComment,
				abstract,
				fullConstantName,
				c.GlobalEnv,
			)
			mixinSingleton := mixin.Singleton()
			c.replaceSimpleNamespacePlaceholder(ct, mixin, mixinSingleton, constantName)
			namespace.DefineConstant(constantName, mixinSingleton)
			namespace.DefineSubtype(constantName, mixin)
			return mixin
		default:
			c.addRedeclaredConstantError(fullConstantName, span)
			return types.NewMixin(docComment, abstract, fullConstantName, c.GlobalEnv)
		}

		switch t := constantType.(type) {
		case *types.Mixin:
			if abstract != t.IsAbstract() {
				c.addFailure(
					fmt.Sprintf(
						"cannot redeclare mixin `%s` with a different modifier, is `%s`, should be `%s`",
						fullConstantName,
						types.InspectModifier(abstract, false, false),
						types.InspectModifier(t.IsAbstract(), false, false),
					),
					span,
				)
			}
			t.AppendDocComment(docComment)
			return t
		case *types.NamespacePlaceholder:
			mixin := types.NewMixinWithDetails(
				docComment,
				abstract,
				t.Name(),
				nil,
				t.Constants(),
				t.Subtypes(),
				t.Methods(),
				c.GlobalEnv,
			)
			t.Namespace = mixin
			namespace.DefineConstant(constantName, mixin.Singleton())
			namespace.DefineSubtype(constantName, mixin)
			return mixin
		default:
			c.addRedeclaredConstantError(fullConstantName, span)
			return types.NewMixin(docComment, abstract, fullConstantName, c.GlobalEnv)
		}
	}

	if namespace == nil {
		return types.NewMixin(docComment, abstract, fullConstantName, c.GlobalEnv)
	}

	return namespace.DefineMixin(docComment, abstract, constantName, c.GlobalEnv)
}

func (c *Checker) declareInterface(docComment string, namespace types.Namespace, constantType types.Type, fullConstantName string, constantName value.Symbol, span *position.Span) *types.Interface {
	if constantType != nil {
		switch ct := constantType.(type) {
		case *types.SingletonClass:
			constantType = ct.AttachedObject
		case *types.Placeholder:
			iface := types.NewInterface(
				docComment,
				fullConstantName,
				c.GlobalEnv,
			)
			ifaceSingleton := iface.Singleton()
			c.replaceSimpleNamespacePlaceholder(ct, iface, ifaceSingleton, constantName)
			namespace.DefineConstant(constantName, ifaceSingleton)
			namespace.DefineSubtype(constantName, iface)
			return iface
		default:
			c.addRedeclaredConstantError(fullConstantName, span)
			return types.NewInterface(docComment, fullConstantName, c.GlobalEnv)
		}

		switch t := constantType.(type) {
		case *types.Interface:
			t.AppendDocComment(docComment)
			return t
		case *types.NamespacePlaceholder:
			iface := types.NewInterfaceWithDetails(
				t.Name(),
				nil,
				t.Constants(),
				t.Subtypes(),
				t.Methods(),
			)
			t.Namespace = iface
			namespace.DefineConstant(constantName, iface.Singleton())
			namespace.DefineSubtype(constantName, iface)
			return iface
		default:
			c.addRedeclaredConstantError(fullConstantName, span)
			return types.NewInterface(docComment, fullConstantName, c.GlobalEnv)
		}
	} else if namespace == nil {
		return types.NewInterface(docComment, fullConstantName, c.GlobalEnv)
	} else {
		return namespace.DefineInterface(docComment, constantName, c.GlobalEnv)
	}
}
