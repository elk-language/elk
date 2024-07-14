// Package checker implements the Elk type checker
package checker

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/concurrent"
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
	typedAst := checker.checkProgram(ast)
	return typedAst, checker.Errors.ErrorList
}

// Represents a single local variable or local value
type local struct {
	typ              types.Type
	initialised      bool
	singleAssignment bool
}

func newLocal(typ types.Type, initialised, singleAssignment bool) *local {
	return &local{
		typ:              typ,
		initialised:      initialised,
		singleAssignment: singleAssignment,
	}
}

// Contains definitions of local variables and values
type localEnvironment struct {
	parent *localEnvironment
	locals map[value.Symbol]*local
}

// Get the local with the specified name from this local environment
func (l *localEnvironment) getLocal(name string) *local {
	local := l.locals[value.ToSymbol(name)]
	return local
}

// Resolve the local with the given name from this local environment or any parent environment
func (l *localEnvironment) resolveLocal(name string) *local {
	nameSymbol := value.ToSymbol(name)
	currentEnv := l
	for {
		if currentEnv == nil {
			return nil
		}
		loc, ok := currentEnv.locals[nameSymbol]
		if ok {
			return loc
		}
		currentEnv = currentEnv.parent
	}
}

func newLocalEnvironment(parent *localEnvironment) *localEnvironment {
	return &localEnvironment{
		parent: parent,
		locals: make(map[value.Symbol]*local),
	}
}

type constantScope struct {
	container types.Namespace
	local     bool
}

func makeLocalConstantScope(container types.Namespace) constantScope {
	return constantScope{
		container: container,
		local:     true,
	}
}

func makeConstantScope(container types.Namespace) constantScope {
	return constantScope{
		container: container,
		local:     false,
	}
}

type methodScope struct {
	container types.Namespace
	local     bool
}

func makeLocalMethodScope(container types.Namespace) methodScope {
	return methodScope{
		container: container,
		local:     true,
	}
}

func makeMethodScope(container types.Namespace) methodScope {
	return methodScope{
		container: container,
		local:     false,
	}
}

type mode uint8

const (
	topLevelMode mode = iota
	moduleMode
	classMode
	mixinMode
	interfaceMode
	methodMode
	singletonMode
)

type methodCheckEntry struct {
	filename       string
	method         *types.Method
	constantScopes []constantScope
	methodScopes   []methodScope
	node           *ast.MethodDefinitionNode
}

// Holds the state of the type checking process
type Checker struct {
	Filename                string
	Errors                  *error.SyncErrorList
	GlobalEnv               *types.GlobalEnvironment
	HeaderMode              bool
	constantScopes          []constantScope
	constantScopesCopyCache []constantScope
	methodScopes            []methodScope
	methodScopesCopyCache   []methodScope
	localEnvs               []*localEnvironment
	returnType              types.Type
	throwType               types.Type
	selfType                types.Type
	mode                    mode
	placeholderNamespaces   *concurrent.Slice[*types.PlaceholderNamespace]
	methodChecks            *concurrent.Slice[methodCheckEntry]
}

// Instantiate a new Checker instance.
func newChecker(filename string, globalEnv *types.GlobalEnvironment, headerMode bool) *Checker {
	if globalEnv == nil {
		globalEnv = types.NewGlobalEnvironment()
	}
	return &Checker{
		Filename:   filename,
		GlobalEnv:  globalEnv,
		HeaderMode: headerMode,
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
		},
		localEnvs: []*localEnvironment{
			newLocalEnvironment(nil),
		},
		placeholderNamespaces: concurrent.NewSlice[*types.PlaceholderNamespace](),
		methodChecks:          concurrent.NewSlice[methodCheckEntry](),
	}
}

// Instantiate a new Checker instance.
func New() *Checker {
	globalEnv := types.NewGlobalEnvironment()
	return &Checker{
		GlobalEnv:  globalEnv,
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
		},
		localEnvs: []*localEnvironment{
			newLocalEnvironment(nil),
		},
		placeholderNamespaces: concurrent.NewSlice[*types.PlaceholderNamespace](),
		methodChecks:          concurrent.NewSlice[methodCheckEntry](),
	}
}

func (c *Checker) newChecker(filename string, constScopes []constantScope, methodScopes []methodScope, selfType, returnType, throwType types.Type) *Checker {
	return &Checker{
		GlobalEnv:      c.GlobalEnv,
		Filename:       filename,
		mode:           methodMode,
		selfType:       selfType,
		returnType:     returnType,
		throwType:      throwType,
		constantScopes: constScopes,
		methodScopes:   methodScopes,
		Errors:         c.Errors,
		localEnvs: []*localEnvironment{
			newLocalEnvironment(nil),
		},
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

func (c *Checker) setMode(mode mode) {
	c.mode = mode
}

func (c *Checker) ClearErrors() {
	c.Errors = new(error.SyncErrorList)
}

func (c *Checker) popLocalEnv() {
	c.localEnvs = c.localEnvs[:len(c.localEnvs)-1]
}

func (c *Checker) pushLocalEnv(env *localEnvironment) {
	c.localEnvs = append(c.localEnvs, env)
}

func (c *Checker) currentLocalEnv() *localEnvironment {
	return c.localEnvs[len(c.localEnvs)-1]
}

func (c *Checker) popConstScope() {
	c.constantScopes = c.constantScopes[:len(c.constantScopes)-1]
	c.clearConstScopeCopyCache()
}

func (c *Checker) pushConstScope(constScope constantScope) {
	c.constantScopes = append(c.constantScopes, constScope)
	c.clearConstScopeCopyCache()
}

func (c *Checker) clearConstScopeCopyCache() {
	c.constantScopesCopyCache = nil
}

func (c *Checker) constantScopesCopy() []constantScope {
	if c.constantScopesCopyCache != nil {
		return c.constantScopesCopyCache
	}

	scopesCopy := make([]constantScope, len(c.constantScopes))
	copy(scopesCopy, c.constantScopes)
	c.constantScopesCopyCache = scopesCopy
	return scopesCopy
}

func (c *Checker) methodScopesCopy() []methodScope {
	if c.methodScopesCopyCache != nil {
		return c.methodScopesCopyCache
	}

	scopesCopy := make([]methodScope, len(c.methodScopes))
	copy(scopesCopy, c.methodScopes)
	c.methodScopesCopyCache = scopesCopy
	return scopesCopy
}

func (c *Checker) clearMethodScopeCopyCache() {
	c.methodScopesCopyCache = nil
}

func (c *Checker) currentConstScope() constantScope {
	for i := range len(c.constantScopes) {
		constScope := c.constantScopes[len(c.constantScopes)-i-1]
		if constScope.local {
			return constScope
		}
	}

	panic("no local constant scopes!")
}

func (c *Checker) popMethodScope() {
	c.methodScopes = c.methodScopes[:len(c.methodScopes)-1]
	c.clearMethodScopeCopyCache()
}

func (c *Checker) pushMethodScope(methodScope methodScope) {
	c.methodScopes = append(c.methodScopes, methodScope)
	c.clearMethodScopeCopyCache()
}

func (c *Checker) currentMethodScope() methodScope {
	for i := range len(c.methodScopes) {
		methodScope := c.methodScopes[len(c.methodScopes)-i-1]
		if methodScope.local {
			return methodScope
		}
	}

	panic("no local method scopes!")
}

func (c *Checker) registerMethodCheck(method *types.Method, node *ast.MethodDefinitionNode) {
	if c.HeaderMode {
		return
	}

	c.methodChecks.Append(methodCheckEntry{
		method:         method,
		constantScopes: c.constantScopesCopy(),
		methodScopes:   c.methodScopesCopy(),
		node:           node,
		filename:       c.Filename,
	})
}

// Create a new location struct with the given position.
func (c *Checker) newLocation(span *position.Span) *position.Location {
	return position.NewLocationWithSpan(c.Filename, span)
}

func (c *Checker) checkStatements(stmts []ast.StatementNode) {
	for _, statement := range stmts {
		c.checkStatement(statement)
	}
}

func (c *Checker) checkStatement(node ast.Node) {
	switch node := node.(type) {
	case *ast.EmptyStatementNode:
	case *ast.ExpressionStatementNode:
		node.Expression = c.checkExpression(node.Expression)
	default:
		c.addFailure(
			fmt.Sprintf("incorrect statement type %#v", node),
			node.Span(),
		)
	}
}

func (c *Checker) isTheSameType(a, b types.Type, errSpan *position.Span) bool {
	return c.isSubtype(a, b, errSpan) && c.isSubtype(b, a, errSpan)
}

func (c *Checker) isSubtype(a, b types.Type, errSpan *position.Span) bool {
	if a == nil && b != nil || a != nil && b == nil {
		return false
	}
	if a == nil && b == nil {
		return true
	}

	if aNamedType, ok := a.(*types.NamedType); ok {
		a = aNamedType.Type
	}
	if bNamedType, ok := b.(*types.NamedType); ok {
		b = bNamedType.Type
	}

	if types.IsNever(a) {
		return true
	}
	if types.IsAny(b) || types.IsVoid(b) {
		return true
	}
	if types.IsAny(a) || types.IsVoid(a) {
		return false
	}
	aNonLiteral := c.toNonLiteral(a)
	if a != aNonLiteral && c.isSubtype(aNonLiteral, b, errSpan) {
		return true
	}

	if aNilable, aIsNilable := a.(*types.Nilable); aIsNilable {
		return c.isSubtype(aNilable.Type, b, errSpan) && c.isSubtype(c.GlobalEnv.StdSubtype(symbol.Nil), b, errSpan)
	}

	if bNilable, bIsNilable := b.(*types.Nilable); bIsNilable {
		return c.isSubtype(a, bNilable.Type, errSpan) || c.isSubtype(a, c.GlobalEnv.StdSubtype(symbol.Nil), errSpan)
	}

	if aUnion, aIsUnion := a.(*types.Union); aIsUnion {
		for _, aElement := range aUnion.Elements {
			if !c.isSubtype(aElement, b, errSpan) {
				return false
			}
		}
		return true
	}

	if bUnion, bIsUnion := b.(*types.Union); bIsUnion {
		for _, bElement := range bUnion.Elements {
			if c.isSubtype(a, bElement, errSpan) {
				return true
			}
		}
		return false
	}

	originalA := a
	switch a := a.(type) {
	case types.Any:
		return types.IsAny(b)
	case *types.SingletonClass:
		b, ok := b.(*types.SingletonClass)
		if !ok {
			return false
		}
		return a.AttachedObject == b.AttachedObject
	case *types.Class:
		return c.classIsSubtype(a, b, errSpan)
	case *types.Mixin:
		return c.mixinIsSubtype(a, b, errSpan)
	case *types.Interface:
		return c.interfaceIsSubtype(a, b, errSpan)
	case *types.Module:
		b, ok := b.(*types.Module)
		if !ok {
			return false
		}
		return a == b
	case *types.Method:
		b, ok := b.(*types.Method)
		if !ok {
			return false
		}
		return a == b
	case *types.CharLiteral:
		b, ok := b.(*types.CharLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.StringLiteral:
		b, ok := b.(*types.StringLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.SymbolLiteral:
		b, ok := b.(*types.SymbolLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.FloatLiteral:
		b, ok := b.(*types.FloatLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Float64Literal:
		b, ok := b.(*types.Float64Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Float32Literal:
		b, ok := b.(*types.Float32Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.BigFloatLiteral:
		b, ok := b.(*types.BigFloatLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.IntLiteral:
		b, ok := b.(*types.IntLiteral)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Int64Literal:
		b, ok := b.(*types.Int64Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Int32Literal:
		b, ok := b.(*types.Int32Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Int16Literal:
		b, ok := b.(*types.Int16Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.Int8Literal:
		b, ok := b.(*types.Int8Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.UInt64Literal:
		b, ok := b.(*types.UInt64Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.UInt32Literal:
		b, ok := b.(*types.UInt32Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.UInt16Literal:
		b, ok := b.(*types.UInt16Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	case *types.UInt8Literal:
		b, ok := b.(*types.UInt8Literal)
		if !ok {
			return false
		}
		return a.Value == b.Value
	default:
		panic(fmt.Sprintf("invalid type: %T", originalA))
	}
}

func (c *Checker) classIsSubtype(a *types.Class, b types.Type, errSpan *position.Span) bool {
	switch b := b.(type) {
	case *types.Class:
		var currentClass types.Namespace = a
		for {
			if currentClass == nil {
				return false
			}
			if currentClass == b {
				return true
			}

			currentClass = currentClass.Parent()
		}
	case *types.Mixin:
		return c.isSubtypeOfMixin(a, b, errSpan)
	case *types.Interface:
		return c.isSubtypeOfInterface(a, b, errSpan)
	default:
		return false
	}
}

func (c *Checker) isSubtypeOfMixin(a types.Namespace, b *types.Mixin, errSpan *position.Span) bool {
	var currentContainer types.Namespace = a
	for {
		switch cont := currentContainer.(type) {
		case *types.Mixin:
			if cont == b {
				return true
			}
		case *types.MixinProxy:
			if cont.Mixin == b {
				return true
			}
		case nil:
			return false
		}

		currentContainer = currentContainer.Parent()
	}
}

func (c *Checker) mixinIsSubtype(a *types.Mixin, b types.Type, errSpan *position.Span) bool {
	bMixin, ok := b.(*types.Mixin)
	if !ok {
		return false
	}

	return c.isSubtypeOfMixin(a, bMixin, errSpan)
}

type methodOverride struct {
	superMethod *types.Method
	override    *types.Method
}

func (c *Checker) implementsInterface(a types.Namespace, b *types.Interface) bool {
	var currentContainer types.Namespace = a
loop:
	for {
		switch cont := currentContainer.(type) {
		case *types.Interface:
			if cont == b {
				return true
			}
		case *types.InterfaceProxy:
			if cont.Interface == b {
				return true
			}
		case nil:
			break loop
		}

		currentContainer = currentContainer.Parent()
	}

	return false
}

func (c *Checker) isSubtypeOfInterface(a types.Namespace, b *types.Interface, errSpan *position.Span) bool {
	if c.implementsInterface(a, b) {
		return true
	}

	var incorrectMethods []methodOverride
	var currentInterface types.Namespace = b
	for currentInterface != nil {
		for _, abstractMethod := range currentInterface.Methods().Map {
			method := types.GetMethodInNamespace(a, abstractMethod.Name)
			if method == nil || !c.checkMethodCompatibility(abstractMethod, method, nil) {
				incorrectMethods = append(incorrectMethods, methodOverride{
					superMethod: abstractMethod,
					override:    method,
				})
			}
		}

		currentInterface = currentInterface.Parent()
	}

	if len(incorrectMethods) > 0 {
		methodDetailsBuff := new(strings.Builder)
		for _, incorrectMethod := range incorrectMethods {
			implementation := incorrectMethod.override
			abstractMethod := incorrectMethod.superMethod
			if implementation == nil {
				fmt.Fprintf(
					methodDetailsBuff,
					"\n  - missing method `%s` with signature: `%s`\n",
					types.InspectWithColor(abstractMethod),
					abstractMethod.InspectSignatureWithColor(false),
				)
				continue
			}

			fmt.Fprintf(
				methodDetailsBuff,
				"\n  - incorrect implementation of `%s`\n      is:        `%s`\n      should be: `%s`\n",
				types.InspectWithColor(abstractMethod),
				implementation.InspectSignatureWithColor(false),
				abstractMethod.InspectSignatureWithColor(false),
			)
		}

		c.addFailure(
			fmt.Sprintf(
				"type `%s` does not implement interface `%s`:\n%s",
				types.InspectWithColor(a),
				types.InspectWithColor(b),
				methodDetailsBuff.String(),
			),
			errSpan,
		)

		return false
	}

	return true
}

func (c *Checker) interfaceIsSubtype(a *types.Interface, b types.Type, errSpan *position.Span) bool {
	bInterface, ok := b.(*types.Interface)
	if !ok {
		return false
	}

	return c.isSubtypeOfInterface(a, bInterface, errSpan)
}

func (c *Checker) checkExpressions(exprs []ast.ExpressionNode) {
	for i, expr := range exprs {
		exprs[i] = c.checkExpression(expr)
	}
}

func (c *Checker) checkExpressionsWithinModule(node *ast.ModuleDeclarationNode) {
	module, ok := c.typeOf(node).(*types.Module)
	if ok {
		c.pushConstScope(makeLocalConstantScope(module))
		c.pushMethodScope(makeLocalMethodScope(module))
	}

	previousSelf := c.selfType
	c.selfType = module
	c.checkStatements(node.Body)
	c.selfType = previousSelf

	if ok {
		c.popConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) checkExpressionsWithinClass(node *ast.ClassDeclarationNode) {
	class, ok := c.typeOf(node).(*types.Class)
	if ok {
		c.checkAbstractMethods(class, node.Constant.Span())
		c.pushConstScope(makeLocalConstantScope(class))
		c.pushMethodScope(makeLocalMethodScope(class))
	}

	previousSelf := c.selfType
	c.selfType = class.Singleton()
	c.checkStatements(node.Body)
	c.selfType = previousSelf

	if ok {
		c.popConstScope()
		c.popMethodScope()
	}
}
func (c *Checker) checkExpressionsWithinMixin(node *ast.MixinDeclarationNode) {
	mixin, ok := c.typeOf(node).(*types.Mixin)
	if ok {
		c.checkAbstractMethods(mixin, node.Constant.Span())
		c.pushConstScope(makeLocalConstantScope(mixin))
		c.pushMethodScope(makeLocalMethodScope(mixin))
	}

	previousSelf := c.selfType
	c.selfType = mixin.Singleton()
	c.checkStatements(node.Body)
	c.selfType = previousSelf

	if ok {
		c.popConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) checkExpressionsWithinInterface(node *ast.InterfaceDeclarationNode) {
	iface, ok := c.typeOf(node).(*types.Interface)
	if ok {
		c.pushConstScope(makeLocalConstantScope(iface))
		c.pushMethodScope(makeLocalMethodScope(iface))
	}

	previousSelf := c.selfType
	c.selfType = iface.Singleton()
	c.checkStatements(node.Body)
	c.selfType = previousSelf

	if ok {
		c.popConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) checkExpressionsWithinSingleton(node *ast.SingletonBlockExpressionNode) {
	class, ok := c.typeOf(node).(*types.SingletonClass)
	if ok {
		c.pushConstScope(makeLocalConstantScope(class))
		c.pushMethodScope(makeLocalMethodScope(class))
	}

	previousSelf := c.selfType
	c.selfType = c.GlobalEnv.StdSubtype(symbol.Class)
	c.checkStatements(node.Body)
	c.selfType = previousSelf

	if ok {
		c.popConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) checkExpression(node ast.ExpressionNode) ast.ExpressionNode {
	switch n := node.(type) {
	case *ast.FalseLiteralNode, *ast.TrueLiteralNode, *ast.NilLiteralNode,
		*ast.InterpolatedSymbolLiteralNode, *ast.ConstantDeclarationNode,
		*ast.InitDefinitionNode, *ast.MethodDefinitionNode, *ast.TypeDefinitionNode,
		*ast.ImplementExpressionNode, *ast.MethodSignatureDefinitionNode,
		*ast.InstanceVariableDeclarationNode, *ast.GetterDeclarationNode,
		*ast.SetterDeclarationNode, *ast.AttrDeclarationNode, *ast.AliasDeclarationNode:
		return n
	case *ast.IncludeExpressionNode:
		c.checkIncludeExpression(n)
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
		c.interpolatedStringLiteral(n)
		return n
	case *ast.SimpleSymbolLiteralNode:
		n.SetType(types.NewSymbolLiteral(n.Content))
		return n
	case *ast.VariableDeclarationNode:
		c.variableDeclaration(n)
		return n
	case *ast.ValueDeclarationNode:
		c.valueDeclaration(n)
		return n
	case *ast.PublicIdentifierNode:
		c.publicIdentifier(n)
		return n
	case *ast.PrivateIdentifierNode:
		c.privateIdentifier(n)
		return n
	case *ast.InstanceVariableNode:
		c.instanceVariable(n)
		return n
	case *ast.PublicConstantNode:
		c.publicConstant(n)
		return n
	case *ast.PrivateConstantNode:
		c.privateConstant(n)
		return n
	case *ast.ConstantLookupNode:
		return c.constantLookup(n)
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
		c.assignmentExpression(n)
		return n
	case *ast.ReceiverlessMethodCallNode:
		c.receiverlessMethodCall(n)
		return n
	case *ast.MethodCallNode:
		c.methodCall(n)
		return n
	case *ast.ConstructorCallNode:
		c.constructorCall(n)
		return n
	case *ast.AttributeAccessNode:
		return c.attributeAccess(n)
	case *ast.BinaryExpressionNode:
		return c.checkBinaryExpression(n)
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

func (c *Checker) checkArithmeticBinaryOperator(
	left,
	right ast.ExpressionNode,
	methodName value.Symbol,
	span *position.Span,
) types.Type {
	leftType := c.toNonLiteral(c.typeOf(left))
	leftClassType, leftIsClass := leftType.(*types.Class)

	rightType := c.toNonLiteral(c.typeOf(right))
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

func (c *Checker) checkBinaryExpression(node *ast.BinaryExpressionNode) ast.ExpressionNode {
	node.Left = c.checkExpression(node.Left)
	node.Right = c.checkExpression(node.Right)

	switch node.Op.Type {
	case token.PLUS:
		node.SetType(c.checkArithmeticBinaryOperator(node.Left, node.Right, symbol.OpAdd, node.Span()))
	case token.MINUS:
		node.SetType(c.checkArithmeticBinaryOperator(node.Left, node.Right, symbol.OpSubtract, node.Span()))
	case token.STAR:
		node.SetType(c.checkArithmeticBinaryOperator(node.Left, node.Right, symbol.OpMultiply, node.Span()))
	case token.SLASH:
		node.SetType(c.checkArithmeticBinaryOperator(node.Left, node.Right, symbol.OpDivide, node.Span()))
	case token.STAR_STAR:
		node.SetType(c.checkArithmeticBinaryOperator(node.Left, node.Right, symbol.OpExponentiate, node.Span()))
	// case token.LBITSHIFT:

	// case token.LTRIPLE_BITSHIFT:

	// case token.RBITSHIFT:

	// case token.RTRIPLE_BITSHIFT:

	// case token.AND:

	// case token.AND_TILDE:

	// case token.OR:

	// case token.XOR:

	// case token.PERCENT:

	// case token.LAX_EQUAL:

	// case token.LAX_NOT_EQUAL:

	// case token.EQUAL_EQUAL:

	// case token.NOT_EQUAL:

	// case token.STRICT_EQUAL:

	// case token.STRICT_NOT_EQUAL:

	// case token.GREATER:

	// case token.GREATER_EQUAL:

	// case token.LESS:

	// case token.LESS_EQUAL:

	// case token.SPACESHIP_OP:

	// case token.INSTANCE_OF_OP:

	// case token.ISA_OP:
	default:
		node.SetType(types.Void{})
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

func (c *Checker) checkAbstractMethods(namespace types.Namespace, span *position.Span) {
	if namespace.IsAbstract() {
		return
	}

	parent := namespace.Parent()

	for parent != nil {
		if !parent.IsAbstract() {
			parent = parent.Parent()
			continue
		}
		for _, parentMethod := range parent.Methods().Map {
			if !parentMethod.IsAbstract() {
				continue
			}

			method := types.GetMethodInNamespace(namespace, parentMethod.Name)
			if method == nil || method.IsAbstract() {
				c.addFailure(
					fmt.Sprintf(
						"missing abstract method implementation `%s` with signature: `%s`",
						types.InspectWithColor(parentMethod),
						parentMethod.InspectSignatureWithColor(false),
					),
					span,
				)
			}
		}
		parent = parent.Parent()
	}
}

// Search through the ancestor chain of the current namespace
// looking for the direct parent of the proxy representing the given mixin.
func (c *Checker) findParentOfMixinProxy(mixin *types.Mixin) types.Namespace {
	currentNamespace := c.currentConstScope().container
	currentParent := currentNamespace.Parent()
	mixinRootParent := types.FindRootParent(mixin)

	for ; currentParent != nil; currentParent = currentParent.Parent() {
		if types.NamespacesAreEqual(currentParent, mixinRootParent) {
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

func (c *Checker) checkIncludeExpression(node *ast.IncludeExpressionNode) {
	targetNamespace := c.currentMethodScope().container

	for _, constantNode := range node.Constants {
		constantType := c.typeOf(constantNode)
		includedMixin, ok := constantType.(*types.Mixin)
		if !ok || includedMixin == nil {
			continue
		}

		parentOfMixin := c.findParentOfMixinProxy(includedMixin)
		if parentOfMixin == nil {
			continue
		}

		var incompatibleMethods []methodOverride
		types.ForeachMethod(includedMixin, func(name value.Symbol, includedMethod *types.Method) {
			superMethod := types.GetMethodInNamespace(parentOfMixin, name)
			if !c.checkMethodCompatibility(superMethod, includedMethod, nil) {
				incompatibleMethods = append(incompatibleMethods, methodOverride{
					superMethod: superMethod,
					override:    includedMethod,
				})
			}
		})

		var incompatibleIvars []instanceVariableOverride
		types.ForeachInstanceVariable(includedMixin, func(name value.Symbol, includedIvar types.Type, includedNamespace types.Namespace) {
			superIvar, superNamespace := types.GetInstanceVariableInNamespace(parentOfMixin, name)
			if superIvar == nil {
				return
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
		})

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

func (c *Checker) includeMixin(node ast.ComplexConstantNode) {
	constantType, _ := c.resolveConstantType(node)

	constantMixin, constantIsMixin := constantType.(*types.Mixin)
	if !constantIsMixin {
		c.addFailure(
			"only mixins can be included",
			node.Span(),
		)
		return
	}
	node.SetType(constantMixin)

	switch c.mode {
	case classMode, mixinMode, singletonMode:
	default:
		c.addFailure(
			"cannot include mixins in this context",
			node.Span(),
		)
		return
	}

	target := c.currentConstScope().container
	if c.isSubtypeOfMixin(target, constantMixin, nil) {
		return
	}

	if target.IsPrimitive() && constantMixin.DeclaresInstanceVariables() {
		c.addFailure(
			fmt.Sprintf(
				"cannot include mixin with instance variables `%s` in primitive `%s`",
				types.InspectWithColor(constantType),
				types.InspectWithColor(target),
			),
			node.Span(),
		)
	}

	switch t := target.(type) {
	case *types.Class:
		t.IncludeMixin(constantMixin)
	case *types.SingletonClass:
		t.IncludeMixin(constantMixin)
	case *types.Mixin:
		t.IncludeMixin(constantMixin)
	default:
		c.addFailure(
			fmt.Sprintf(
				"cannot include `%s` in `%s`",
				types.InspectWithColor(constantType),
				types.InspectWithColor(t),
			),
			node.Span(),
		)
	}
}

func (c *Checker) implementInterface(node ast.ComplexConstantNode) {
	constantType, _ := c.resolveConstantType(node)

	constantInterface, constantIsInterface := constantType.(*types.Interface)
	if !constantIsInterface {
		c.addFailure(
			"only interfaces can be implemented",
			node.Span(),
		)
		return
	}

	switch c.mode {
	case classMode, mixinMode, interfaceMode:
	default:
		c.addFailure(
			"cannot implement interfaces in this context",
			node.Span(),
		)
		return
	}

	target := c.currentConstScope().container
	if c.implementsInterface(target, constantInterface) {
		return
	}

	switch t := target.(type) {
	case *types.Class:
		t.ImplementInterface(constantInterface)
	case *types.Mixin:
		t.ImplementInterface(constantInterface)
	case *types.Interface:
		t.ImplementInterface(constantInterface)
	default:
		c.addFailure(
			fmt.Sprintf(
				"cannot implement `%s` in `%s`",
				types.InspectWithColor(constantType),
				types.InspectWithColor(t),
			),
			node.Span(),
		)
	}
}

func (c *Checker) checkComplexConstant(node ast.ComplexConstantNode) ast.ComplexConstantNode {
	switch n := node.(type) {
	case *ast.PublicConstantNode:
		return c.publicConstant(n)
	case *ast.PrivateConstantNode:
		return c.privateConstant(n)
	case *ast.ConstantLookupNode:
		return c.constantLookup(n)
	default:
		c.addFailure(
			fmt.Sprintf("invalid constant type %T", node),
			node.Span(),
		)
		return node
	}
}

// Checks whether two methods are compatible.
func (c *Checker) checkMethodCompatibility(baseMethod, overrideMethod *types.Method, errSpan *position.Span) bool {
	areCompatible := true
	if baseMethod != nil {
		if !c.isSubtype(overrideMethod.ReturnType, baseMethod.ReturnType, errSpan) {
			c.addFailure(
				fmt.Sprintf(
					"method `%s` has a different return type than `%s`, has `%s`, should have `%s`",
					types.InspectWithColor(overrideMethod),
					types.InspectWithColor(baseMethod),
					types.InspectWithColor(overrideMethod.ReturnType),
					types.InspectWithColor(baseMethod.ReturnType),
				),
				errSpan,
			)
			areCompatible = false
		}
		if !c.isSubtype(overrideMethod.ThrowType, baseMethod.ThrowType, errSpan) {
			c.addFailure(
				fmt.Sprintf(
					"method `%s` has a different throw type than `%s`, has `%s`, should have `%s`",
					types.InspectWithColor(overrideMethod),
					types.InspectWithColor(baseMethod),
					types.InspectWithColor(overrideMethod.ThrowType),
					types.InspectWithColor(baseMethod.ThrowType),
				),
				errSpan,
			)
			areCompatible = false
		}

		if len(baseMethod.Params) > len(overrideMethod.Params) {
			c.addFailure(
				fmt.Sprintf(
					"method `%s` has less parameters than `%s`, has `%d`, should have `%d`",
					types.InspectWithColor(overrideMethod),
					types.InspectWithColor(baseMethod),
					len(overrideMethod.Params),
					len(baseMethod.Params),
				),
				errSpan,
			)
			areCompatible = false
		} else {
			for i := range len(baseMethod.Params) {
				oldParam := baseMethod.Params[i]
				newParam := overrideMethod.Params[i]
				if oldParam.Name != newParam.Name {
					c.addFailure(
						fmt.Sprintf(
							"method `%s` has a different parameter name than `%s`, has `%s`, should have `%s`",
							types.InspectWithColor(overrideMethod),
							types.InspectWithColor(baseMethod),
							newParam.Name,
							oldParam.Name,
						),
						errSpan,
					)
					areCompatible = false
					continue
				}
				if oldParam.Kind != newParam.Kind {
					c.addFailure(
						fmt.Sprintf(
							"method `%s` has a different parameter kind than `%s`, has `%s`, should have `%s`",
							types.InspectWithColor(overrideMethod),
							types.InspectWithColor(baseMethod),
							newParam.NameWithKind(),
							oldParam.NameWithKind(),
						),
						errSpan,
					)
					areCompatible = false
					continue
				}
				if !c.isSubtype(oldParam.Type, newParam.Type, errSpan) {
					c.addFailure(
						fmt.Sprintf(
							"method `%s` has a different type for parameter `%s` than `%s`, has `%s`, should have `%s`",
							types.InspectWithColor(overrideMethod),
							newParam.Name,
							types.InspectWithColor(baseMethod),
							types.InspectWithColor(newParam.Type),
							types.InspectWithColor(oldParam.Type),
						),
						errSpan,
					)
					areCompatible = false
					continue
				}
			}

			for i := len(baseMethod.Params); i < len(overrideMethod.Params); i++ {
				param := overrideMethod.Params[i]
				if !param.IsOptional() {
					c.addFailure(
						fmt.Sprintf(
							"method `%s` has a required parameter missing in `%s`, got `%s`",
							types.InspectWithColor(overrideMethod),
							types.InspectWithColor(baseMethod),
							param.Name,
						),
						errSpan,
					)
					areCompatible = false
				}
			}
		}

	}

	return areCompatible
}

func (c *Checker) getMethod(typ types.Type, name value.Symbol, errSpan *position.Span) *types.Method {
	return c._getMethod(typ, name, errSpan, false)
}

func (c *Checker) getMethodInContainer(container types.Namespace, typ types.Type, name value.Symbol, errSpan *position.Span, inParent bool) *types.Method {
	method := types.GetMethodInNamespace(container, name)
	if method != nil {
		return method
	}
	if !inParent {
		c.addMissingMethodError(typ, name.String(), errSpan)
	}
	return nil
}

func (c *Checker) _getMethod(typ types.Type, name value.Symbol, errSpan *position.Span, inParent bool) *types.Method {
	typ = c.toNonLiteral(typ)

	switch t := typ.(type) {
	case *types.NamedType:
		return c._getMethod(t.Type, name, errSpan, inParent)
	case *types.Class:
		return c.getMethodInContainer(t, typ, name, errSpan, inParent)
	case *types.SingletonClass:
		return c.getMethodInContainer(t, typ, name, errSpan, inParent)
	case *types.Interface:
		return c.getMethodInContainer(t, typ, name, errSpan, inParent)
	case *types.InterfaceProxy:
		return c.getMethodInContainer(t, typ, name, errSpan, inParent)
	case *types.Module:
		return c.getMethodInContainer(t, typ, name, errSpan, inParent)
	case *types.Mixin:
		return c.getMethodInContainer(t, typ, name, errSpan, inParent)
	case *types.MixinProxy:
		return c.getMethodInContainer(t, typ, name, errSpan, inParent)
	case *types.Nilable:
		nilType := c.GlobalEnv.StdSubtype(symbol.Nil).(*types.Class)
		nilMethod := nilType.Method(name)
		if nilMethod == nil {
			c.addMissingMethodError(nilType, name.String(), errSpan)
		}
		nonNilMethod := c.getMethod(t.Type, name, errSpan)
		if nilMethod == nil || nonNilMethod == nil {
			return nil
		}

		var baseMethod *types.Method
		var overrideMethod *types.Method
		if len(nilMethod.Params) < len(nonNilMethod.Params) {
			baseMethod = nilMethod
			overrideMethod = nonNilMethod
		} else {
			baseMethod = nonNilMethod
			overrideMethod = nilMethod
		}

		if c.checkMethodCompatibility(baseMethod, overrideMethod, errSpan) {
			return baseMethod
		}
		return nil
	case *types.Union:
		var methods []*types.Method
		var baseMethod *types.Method

		for _, element := range t.Elements {
			elementMethod := c.getMethod(element, name, errSpan)
			if elementMethod == nil {
				continue
			}
			methods = append(methods, elementMethod)
			if baseMethod == nil || len(baseMethod.Params) > len(elementMethod.Params) {
				baseMethod = elementMethod
			}
		}

		if len(methods) < len(t.Elements) {
			return nil
		}

		isCompatible := true
		for i := range len(methods) {
			method := methods[i]

			if !c.checkMethodCompatibility(baseMethod, method, errSpan) {
				isCompatible = false
			}
		}

		if isCompatible {
			return baseMethod
		}

		return nil
	default:
		c.addMissingMethodError(typ, name.String(), errSpan)
		return nil
	}
}

func (c *Checker) addMissingMethodError(typ types.Type, name string, span *position.Span) {
	c.addFailure(
		fmt.Sprintf("method `%s` is not defined on type `%s`", name, types.InspectWithColor(typ)),
		span,
	)
}

func (c *Checker) typeOf(node ast.Node) types.Type {
	return node.Type(c.GlobalEnv)
}

func (c *Checker) isNilable(typ types.Type) bool {
	return types.IsNilable(typ, c.GlobalEnv)
}

func (c *Checker) canBeFalsy(typ types.Type) bool {
	return types.CanBeFalsy(typ, c.GlobalEnv)
}

func (c *Checker) canBeTruthy(typ types.Type) bool {
	return types.CanBeTruthy(typ, c.GlobalEnv)
}

func (c *Checker) toNonNilable(typ types.Type) types.Type {
	return types.ToNonNilable(typ, c.GlobalEnv)
}
func (c *Checker) toNonLiteral(typ types.Type) types.Type {
	return typ.ToNonLiteral(c.GlobalEnv)
}

func (c *Checker) toNilable(typ types.Type) types.Type {
	return types.ToNilable(typ, c.GlobalEnv)
}

func (c *Checker) checkMethodArguments(method *types.Method, positionalArguments []ast.ExpressionNode, namedArguments []ast.NamedArgumentNode, span *position.Span) []ast.ExpressionNode {
	reqParamCount := method.RequiredParamCount()
	requiredPosParamCount := len(method.Params) - method.OptionalParamCount
	if method.PostParamCount != -1 {
		requiredPosParamCount -= method.PostParamCount + 1
	}
	positionalRestParamIndex := method.PositionalRestParamIndex()
	var typedPositionalArguments []ast.ExpressionNode

	var currentParamIndex int
	for ; currentParamIndex < len(positionalArguments); currentParamIndex++ {
		posArg := positionalArguments[currentParamIndex]
		if currentParamIndex == positionalRestParamIndex {
			break
		}
		if currentParamIndex >= len(method.Params) {
			c.addWrongArgumentCountError(
				len(positionalArguments)+len(namedArguments),
				method,
				span,
			)
			break
		}
		param := method.Params[currentParamIndex]

		typedPosArg := c.checkExpression(posArg)
		typedPositionalArguments = append(typedPositionalArguments, typedPosArg)
		posArgType := c.typeOf(typedPosArg)
		if !c.isSubtype(posArgType, param.Type, posArg.Span()) {
			c.addFailure(
				fmt.Sprintf(
					"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
					types.InspectWithColor(param.Type),
					param.Name,
					method.Name,
					types.InspectWithColor(posArgType),
				),
				posArg.Span(),
			)
		}
	}

	if method.HasPositionalRestParam() {
		if len(positionalArguments) < requiredPosParamCount {
			c.addFailure(
				fmt.Sprintf(
					"expected %d... positional arguments in call to `%s`, got %d",
					requiredPosParamCount,
					method.Name,
					len(positionalArguments),
				),
				span,
			)
			return nil
		}
		restPositionalArguments := ast.NewArrayTupleLiteralNode(
			span,
			nil,
		)
		posRestParam := method.Params[positionalRestParamIndex]

		currentArgIndex := currentParamIndex
		for ; currentArgIndex < len(positionalArguments)-method.PostParamCount; currentArgIndex++ {
			posArg := positionalArguments[currentArgIndex]
			typedPosArg := c.checkExpression(posArg)
			restPositionalArguments.Elements = append(restPositionalArguments.Elements, typedPosArg)
			posArgType := c.typeOf(typedPosArg)
			if !c.isSubtype(posArgType, posRestParam.Type, posArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for rest parameter `*%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(posRestParam.Type),
						posRestParam.Name,
						method.Name,
						types.InspectWithColor(posArgType),
					),
					posArg.Span(),
				)
			}
		}
		typedPositionalArguments = append(typedPositionalArguments, restPositionalArguments)

		currentParamIndex = positionalRestParamIndex
		for ; currentArgIndex < len(positionalArguments); currentArgIndex++ {
			posArg := positionalArguments[currentArgIndex]
			currentParamIndex++
			param := method.Params[currentParamIndex]

			typedPosArg := c.checkExpression(posArg)
			typedPositionalArguments = append(typedPositionalArguments, typedPosArg)
			posArgType := c.typeOf(typedPosArg)
			if !c.isSubtype(posArgType, param.Type, posArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(param.Type),
						param.Name,
						method.Name,
						types.InspectWithColor(posArgType),
					),
					posArg.Span(),
				)
			}
		}
		currentParamIndex++

		if method.PostParamCount > 0 {
			reqParamCount++
		}
	}

	firstNamedParamIndex := currentParamIndex
	definedNamedArgumentsSlice := make([]bool, len(namedArguments))

	for i := 0; i < len(method.Params); i++ {
		param := method.Params[i]
		switch param.Kind {
		case types.PositionalRestParameterKind, types.NamedRestParameterKind:
			continue
		}
		paramName := param.Name.String()
		var found bool

		for namedArgIndex, namedArgI := range namedArguments {
			namedArg := namedArgI.(*ast.NamedCallArgumentNode)
			if namedArg.Name != paramName {
				continue
			}
			if found || i < firstNamedParamIndex {
				c.addFailure(
					fmt.Sprintf(
						"duplicated argument `%s` in call to `%s`",
						paramName,
						method.Name,
					),
					namedArg.Span(),
				)
			}
			found = true
			definedNamedArgumentsSlice[namedArgIndex] = true
			typedNamedArgValue := c.checkExpression(namedArg.Value)
			namedArgType := c.typeOf(typedNamedArgValue)
			typedPositionalArguments = append(typedPositionalArguments, typedNamedArgValue)
			if !c.isSubtype(namedArgType, param.Type, namedArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(param.Type),
						param.Name,
						method.Name,
						types.InspectWithColor(namedArgType),
					),
					namedArg.Span(),
				)
			}
		}

		if i < firstNamedParamIndex {
			continue
		}
		if found {
			continue
		}

		if i < reqParamCount {
			// the parameter is required
			// but is not present in the call
			c.addFailure(
				fmt.Sprintf(
					"argument `%s` is missing in call to `%s`",
					paramName,
					method.Name,
				),
				span,
			)
		} else {
			// the parameter is missing and is optional
			// we push undefined as its value
			typedPositionalArguments = append(
				typedPositionalArguments,
				ast.NewUndefinedLiteralNode(span),
			)
		}
	}

	if method.HasNamedRestParam {
		namedRestArgs := ast.NewHashRecordLiteralNode(
			span,
			nil,
		)
		namedRestParam := method.Params[len(method.Params)-1]
		for i, defined := range definedNamedArgumentsSlice {
			if defined {
				continue
			}

			namedArgI := namedArguments[i]
			namedArg := namedArgI.(*ast.NamedCallArgumentNode)
			typedNamedArgValue := c.checkExpression(namedArg.Value)
			namedRestArgs.Elements = append(
				namedRestArgs.Elements,
				ast.NewSymbolKeyValueExpressionNode(
					namedArg.Span(),
					namedArg.Name,
					typedNamedArgValue,
				),
			)
			namedArgType := c.typeOf(typedNamedArgValue)
			if !c.isSubtype(namedArgType, namedRestParam.Type, namedArg.Span()) {
				c.addFailure(
					fmt.Sprintf(
						"expected type `%s` for named rest parameter `**%s` in call to `%s`, got type `%s`",
						types.InspectWithColor(namedRestParam.Type),
						namedRestParam.Name,
						method.Name,
						types.InspectWithColor(namedArgType),
					),
					namedArg.Span(),
				)
			}
		}

		typedPositionalArguments = append(typedPositionalArguments, namedRestArgs)
	} else {
		for i, defined := range definedNamedArgumentsSlice {
			if defined {
				continue
			}

			namedArgI := namedArguments[i]
			namedArg := namedArgI.(*ast.NamedCallArgumentNode)
			c.addFailure(
				fmt.Sprintf(
					"nonexistent parameter `%s` given in call to `%s`",
					namedArg.Name,
					method.Name,
				),
				namedArg.Span(),
			)
		}
	}

	return typedPositionalArguments
}

func (c *Checker) receiverlessMethodCall(node *ast.ReceiverlessMethodCallNode) {
	method := c.getMethod(c.selfType, value.ToSymbol(node.MethodName), node.Span())
	if method == nil {
		c.checkExpressions(node.PositionalArguments)
		node.SetType(types.Void{})
		return
	}

	typedPositionalArguments := c.checkMethodArguments(method, node.PositionalArguments, node.NamedArguments, node.Span())
	node.PositionalArguments = typedPositionalArguments
	node.NamedArguments = nil
	node.SetType(method.ReturnType)
}

func (c *Checker) constructorCall(node *ast.ConstructorCallNode) {
	classNode := c.checkComplexConstantType(node.Class)
	node.Class = classNode
	classType := c.typeOf(classNode)
	var className string

	switch cn := classNode.(type) {
	case *ast.PublicConstantNode:
		className = cn.Value
	case *ast.PrivateConstantNode:
		className = cn.Value
	}

	class, isClass := classType.(*types.Class)
	if !isClass {
		c.addFailure(
			fmt.Sprintf("`%s` cannot be instantiated", className),
			node.Span(),
		)
		c.checkExpressions(node.PositionalArguments)
		node.SetType(types.Void{})
		return
	}

	if class.IsAbstract() {
		c.addFailure(
			fmt.Sprintf("cannot instantiate abstract class `%s`", className),
			node.Span(),
		)
	}

	method := types.GetMethodInNamespace(class, symbol.M_init)
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
			class,
		)
	}

	typedPositionalArguments := c.checkMethodArguments(method, node.PositionalArguments, node.NamedArguments, node.Span())

	node.PositionalArguments = typedPositionalArguments
	node.NamedArguments = nil
	node.SetType(class)
}

func (c *Checker) checkSimpleMethodCall(
	receiver ast.ExpressionNode,
	nilSafe bool,
	methodName value.Symbol,
	positionalArguments []ast.ExpressionNode,
	namedArguments []ast.NamedArgumentNode,
	span *position.Span,
) (
	_receiver ast.ExpressionNode,
	_positionalArguments []ast.ExpressionNode,
	typ types.Type,
) {
	receiver = c.checkExpression(receiver)
	receiverType := c.typeOf(receiver)
	var method *types.Method
	if nilSafe {
		nonNilableReceiverType := c.toNonNilable(receiverType)
		method = c.getMethod(nonNilableReceiverType, methodName, span)
	} else {
		method = c.getMethod(receiverType, methodName, span)
	}
	if method == nil {
		c.checkExpressions(positionalArguments)
		return receiver, positionalArguments, types.Void{}
	}

	typedPositionalArguments := c.checkMethodArguments(method, positionalArguments, namedArguments, span)
	returnType := method.ReturnType
	if nilSafe {
		if !c.isNilable(receiverType) {
			c.addFailure(
				fmt.Sprintf("cannot make a nil-safe call on type `%s` which is not nilable", types.InspectWithColor(receiverType)),
				span,
			)
		} else {
			returnType = c.toNilable(returnType)
		}
	}

	return receiver, typedPositionalArguments, returnType
}

func (c *Checker) checkBinaryOpMethodCall(
	left ast.ExpressionNode,
	right ast.ExpressionNode,
	methodName value.Symbol,
	span *position.Span,
) types.Type {
	_, _, returnType := c.checkSimpleMethodCall(
		left,
		false,
		methodName,
		[]ast.ExpressionNode{right},
		nil,
		span,
	)

	return returnType
}

func (c *Checker) methodCall(node *ast.MethodCallNode) {
	var typ types.Type
	node.Receiver, node.PositionalArguments, typ = c.checkSimpleMethodCall(
		node.Receiver,
		node.NilSafe,
		value.ToSymbol(node.MethodName),
		node.PositionalArguments,
		node.NamedArguments,
		node.Span(),
	)
	node.SetType(typ)
}

func (c *Checker) attributeAccess(node *ast.AttributeAccessNode) ast.ExpressionNode {
	receiver := c.checkExpression(node.Receiver)
	receiverType := c.typeOf(receiver)
	method := c.getMethod(receiverType, value.ToSymbol(node.AttributeName), node.Span())
	if method == nil {
		node.Receiver = receiver
		node.SetType(types.Void{})
		return node
	}

	typedPositionalArguments := c.checkMethodArguments(method, nil, nil, node.Span())

	newNode := ast.NewMethodCallNode(
		node.Span(),
		receiver,
		false,
		node.AttributeName,
		typedPositionalArguments,
		nil,
	)
	newNode.SetType(method.ReturnType)
	return newNode
}

func (c *Checker) addWrongArgumentCountError(got int, method *types.Method, span *position.Span) {
	c.addFailure(
		fmt.Sprintf("expected %s arguments in call to `%s`, got %d", method.ExpectedParamCountString(), method.Name, got),
		span,
	)
}

func (c *Checker) addOverrideSealedMethodError(baseMethod *types.Method, span *position.Span) {
	c.addFailure(
		fmt.Sprintf(
			"cannot override sealed method `%s`\n  previous definition found in `%s`, with signature: `%s`",
			baseMethod.Name,
			types.InspectWithColor(baseMethod.DefinedUnder),
			baseMethod.InspectSignatureWithColor(true),
		),
		span,
	)
}

func (c *Checker) checkMethodOverride(
	overrideMethod,
	baseMethod *types.Method,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	span *position.Span,
) {
	name := baseMethod.Name
	if baseMethod.IsSealed() {
		c.addOverrideSealedMethodError(baseMethod, span)
	}
	if !baseMethod.IsAbstract() && overrideMethod.IsAbstract() {
		c.addFailure(
			fmt.Sprintf(
				"cannot override method `%s` with a different modifier, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: `%s`",
				name,
				types.InspectModifier(overrideMethod.IsAbstract(), overrideMethod.IsSealed(), false),
				types.InspectModifier(baseMethod.IsAbstract(), baseMethod.IsSealed(), false),
				types.InspectWithColor(baseMethod.DefinedUnder),
				baseMethod.InspectSignatureWithColor(true),
			),
			span,
		)
	}

	if !c.isSubtype(overrideMethod.ReturnType, baseMethod.ReturnType, nil) {
		var returnSpan *position.Span
		if returnTypeNode != nil {
			returnSpan = returnTypeNode.Span()
		} else {
			returnSpan = span
		}
		c.addFailure(
			fmt.Sprintf(
				"cannot override method `%s` with a different return type, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: `%s`",
				name,
				types.InspectWithColor(overrideMethod.ReturnType),
				types.InspectWithColor(baseMethod.ReturnType),
				types.InspectWithColor(baseMethod.DefinedUnder),
				baseMethod.InspectSignatureWithColor(true),
			),
			returnSpan,
		)
	}
	if !c.isSubtype(overrideMethod.ThrowType, baseMethod.ThrowType, nil) {
		var throwSpan *position.Span
		if throwTypeNode != nil {
			throwSpan = throwTypeNode.Span()
		} else {
			throwSpan = span
		}
		c.addFailure(
			fmt.Sprintf(
				"cannot override method `%s` with a different throw type, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: `%s`",
				name,
				types.InspectWithColor(overrideMethod.ThrowType),
				types.InspectWithColor(baseMethod.ThrowType),
				types.InspectWithColor(baseMethod.DefinedUnder),
				baseMethod.InspectSignatureWithColor(true),
			),
			throwSpan,
		)
	}

	if len(baseMethod.Params) > len(overrideMethod.Params) {
		paramSpan := position.JoinSpanOfCollection(paramNodes)
		if paramSpan == nil {
			paramSpan = span
		}
		c.addFailure(
			fmt.Sprintf(
				"cannot override method `%s` with less parameters\n  previous definition found in `%s`, with signature: `%s`",
				name,
				types.InspectWithColor(baseMethod.DefinedUnder),
				baseMethod.InspectSignatureWithColor(true),
			),
			paramSpan,
		)
	} else {
		for i := range len(baseMethod.Params) {
			oldParam := baseMethod.Params[i]
			newParam := overrideMethod.Params[i]
			if oldParam.Name != newParam.Name {
				c.addFailure(
					fmt.Sprintf(
						"cannot override method `%s` with invalid parameter name, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: `%s`",
						name,
						newParam.Name,
						oldParam.Name,
						types.InspectWithColor(baseMethod.DefinedUnder),
						baseMethod.InspectSignatureWithColor(true),
					),
					paramNodes[i].Span(),
				)
				continue
			}
			if oldParam.Kind != newParam.Kind {
				c.addFailure(
					fmt.Sprintf(
						"cannot override method `%s` with invalid parameter kind, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: `%s`",
						name,
						newParam.NameWithKind(),
						oldParam.NameWithKind(),
						types.InspectWithColor(baseMethod.DefinedUnder),
						baseMethod.InspectSignatureWithColor(true),
					),
					paramNodes[i].Span(),
				)
				continue
			}
			if !c.isSubtype(oldParam.Type, newParam.Type, paramNodes[i].Span()) {
				c.addFailure(
					fmt.Sprintf(
						"cannot override method `%s` with invalid parameter type, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: `%s`",
						name,
						types.InspectWithColor(newParam.Type),
						types.InspectWithColor(oldParam.Type),
						types.InspectWithColor(baseMethod.DefinedUnder),
						baseMethod.InspectSignatureWithColor(true),
					),
					paramNodes[i].Span(),
				)
				continue
			}
		}

		for i := len(baseMethod.Params); i < len(overrideMethod.Params); i++ {
			param := overrideMethod.Params[i]
			if !param.IsOptional() {
				c.addFailure(
					fmt.Sprintf(
						"cannot override method `%s` with additional parameter `%s`\n  previous definition found in `%s`, with signature: `%s`",
						name,
						param.Name,
						types.InspectWithColor(baseMethod.DefinedUnder),
						baseMethod.InspectSignatureWithColor(true),
					),
					paramNodes[i].Span(),
				)
			}
		}
	}

}

func (c *Checker) checkMethod(
	checkedMethod *types.Method,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	body []ast.StatementNode,
	span *position.Span,
) (ast.TypeNode, ast.TypeNode) {
	methodNamespace := c.currentMethodScope().container
	name := checkedMethod.Name

	currentMethod := types.GetMethodInNamespace(methodNamespace, name)
	if checkedMethod != currentMethod && checkedMethod.IsSealed() {
		c.addOverrideSealedMethodError(checkedMethod, currentMethod.Span())
	}

	parent := methodNamespace.Parent()

	if parent != nil {
		baseMethod := types.GetMethodInNamespace(parent, name)
		if baseMethod != nil {
			c.checkMethodOverride(
				checkedMethod,
				baseMethod,
				paramNodes,
				returnTypeNode,
				throwTypeNode,
				span,
			)
		}
	}

	env := newLocalEnvironment(nil)
	c.pushLocalEnv(env)
	defer c.popLocalEnv()

	for _, param := range paramNodes {
		p, _ := param.(*ast.MethodParameterNode)
		var declaredType types.Type
		var declaredTypeNode ast.TypeNode
		if p.TypeNode != nil {
			declaredTypeNode = p.TypeNode
			declaredType = c.typeOf(declaredTypeNode)
		}
		var initNode ast.ExpressionNode
		if p.Initialiser != nil {
			initNode = c.checkExpression(p.Initialiser)
			initType := c.typeOf(initNode)
			c.checkCanAssign(initType, declaredType, initNode.Span())
		}
		c.addLocal(p.Name, newLocal(declaredType, true, false))
		p.Initialiser = initNode
		p.TypeNode = declaredTypeNode
	}

	var returnType types.Type
	var typedReturnTypeNode ast.TypeNode
	if returnTypeNode != nil {
		typedReturnTypeNode = c.checkTypeNode(returnTypeNode)
		returnType = c.typeOf(typedReturnTypeNode)
	} else {
		returnType = types.Void{}
	}

	var throwType types.Type
	var typedThrowTypeNode ast.TypeNode
	if throwTypeNode != nil {
		typedThrowTypeNode = c.checkTypeNode(throwTypeNode)
		throwType = c.typeOf(typedThrowTypeNode)
	}

	if len(body) > 0 && checkedMethod.IsAbstract() {
		c.addFailure(
			fmt.Sprintf(
				"method `%s` cannot have a body because it is abstract",
				name,
			),
			span,
		)
	}

	previousMode := c.mode
	c.mode = methodMode
	defer c.setMode(previousMode)
	c.returnType = returnType
	c.throwType = throwType
	c.checkStatements(body)
	c.returnType = nil
	c.throwType = nil
	return typedReturnTypeNode, typedThrowTypeNode
}

func (c *Checker) assignmentExpression(node *ast.AssignmentExpressionNode) {
	switch n := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		c.localVariableAssignment(n.Value, node)
	case *ast.PrivateIdentifierNode:
		c.localVariableAssignment(n.Value, node)
	// case *ast.SubscriptExpressionNode:
	// case *ast.ConstantLookupNode:
	// case *ast.PublicConstantNode:
	// case *ast.PrivateConstantNode:
	// case *ast.InstanceVariableNode:
	// case *ast.AttributeAccessNode:
	default:
		c.addFailure(
			fmt.Sprintf("cannot assign to: %T", node.Left),
			node.Span(),
		)
	}
}

func (c *Checker) checkAddMethod(receiver types.Type, argument types.Type) (returnType types.Type) {
	// c.toNonLiteral(receiver)
	// receiverClass, ok := receiver.(*types.Class)
	return nil
}

func (c *Checker) localVariableAssignment(name string, node *ast.AssignmentExpressionNode) {
	span := node.Span()
	switch node.Op.Type {
	case token.EQUAL_OP:
		variable := c.getLocal(name)
		if variable == nil {
			c.addFailure(
				fmt.Sprintf("undefined local `%s`", name),
				span,
			)
			break
		}
		if variable.singleAssignment && variable.initialised {
			c.addFailure(
				fmt.Sprintf("local value `%s` cannot be reassigned", name),
				span,
			)
		}

		node.Right = c.checkExpression(node.Right)
		assignedType := c.typeOf(node.Right)
		c.checkCanAssign(assignedType, variable.typ, node.Right.Span())
	case token.QUESTION_QUESTION_EQUAL:
		variable := c.getLocal(name)
		if variable == nil {
			c.addFailure(
				fmt.Sprintf("undefined local `%s`", name),
				span,
			)
			break
		}
		if variable.singleAssignment && variable.initialised {
			c.addFailure(
				fmt.Sprintf("local value `%s` cannot be reassigned", name),
				span,
			)
		}

		if !variable.initialised {
			c.addFailure(
				fmt.Sprintf("cannot access uninitialised local `%s`", name),
				node.Left.Span(),
			)
		}

		if !c.isNilable(variable.typ) {
			c.addFailure(
				fmt.Sprintf("local `%s` is not nilable", name),
				node.Left.Span(),
			)
		}

		node.Right = c.checkExpression(node.Right)
		assignedType := c.typeOf(node.Right)
		c.checkCanAssign(assignedType, variable.typ, node.Right.Span())
	case token.OR_OR_EQUAL:
		variable := c.getLocal(name)
		if variable == nil {
			c.addFailure(
				fmt.Sprintf("undefined local `%s`", name),
				span,
			)
			break
		}
		if variable.singleAssignment && variable.initialised {
			c.addFailure(
				fmt.Sprintf("local value `%s` cannot be reassigned", name),
				span,
			)
		}

		if !variable.initialised {
			c.addFailure(
				fmt.Sprintf("cannot access uninitialised local `%s`", name),
				node.Left.Span(),
			)
		}

		if !c.canBeFalsy(variable.typ) {
			c.addFailure(
				fmt.Sprintf("local `%s` cannot be falsy", name),
				node.Left.Span(),
			)
		}

		node.Right = c.checkExpression(node.Right)
		assignedType := c.typeOf(node.Right)
		c.checkCanAssign(assignedType, variable.typ, node.Right.Span())
	case token.AND_AND_EQUAL:
		variable := c.getLocal(name)
		if variable == nil {
			c.addFailure(
				fmt.Sprintf("undefined local `%s`", name),
				span,
			)
			break
		}
		if variable.singleAssignment && variable.initialised {
			c.addFailure(
				fmt.Sprintf("local value `%s` cannot be reassigned", name),
				span,
			)
		}

		if !variable.initialised {
			c.addFailure(
				fmt.Sprintf("cannot access uninitialised local `%s`", name),
				node.Left.Span(),
			)
		}

		if !c.canBeTruthy(variable.typ) {
			c.addFailure(
				fmt.Sprintf("local `%s` cannot be truthy", name),
				node.Left.Span(),
			)
		}

		node.Right = c.checkExpression(node.Right)
		assignedType := c.typeOf(node.Right)
		c.checkCanAssign(assignedType, variable.typ, node.Right.Span())
	case token.PLUS_EQUAL:
		variable := c.getLocal(name)
		if variable == nil {
			c.addFailure(
				fmt.Sprintf("undefined local `%s`", name),
				span,
			)
			break
		}
		if variable.singleAssignment && variable.initialised {
			c.addFailure(
				fmt.Sprintf("local value `%s` cannot be reassigned", name),
				span,
			)
		}

		if !variable.initialised {
			c.addFailure(
				fmt.Sprintf("cannot access uninitialised local `%s`", name),
				node.Left.Span(),
			)
		}

		c.getMethod(variable.typ, symbol.OpAdd, node.Op.Span())

		// if !c.canBeTruthy(variable.typ) {
		// 	c.addFailure(
		// 		fmt.Sprintf("local `%s` cannot be truthy", name),
		// 		node.Left.Span(),
		// 	)
		// }

		// node.Right = c.checkExpression(node.Right)
		// assignedType := c.typeOf(node.Right)
		// c.checkCanAssign(assignedType, variable.typ, node.Right.Span())
	// case token.MINUS_EQUAL:
	// case token.STAR_EQUAL:
	// case token.SLASH_EQUAL:
	// case token.STAR_STAR_EQUAL:
	// case token.PERCENT_EQUAL:
	// case token.AND_EQUAL:
	// case token.OR_EQUAL:
	// case token.XOR_EQUAL:
	// case token.LBITSHIFT_EQUAL:
	// case token.LTRIPLE_BITSHIFT_EQUAL:
	// case token.RBITSHIFT_EQUAL:
	// case token.RTRIPLE_BITSHIFT_EQUAL:
	// case token.COLON_EQUAL:
	default:
		c.addFailure(
			fmt.Sprintf("assignment using this operator has not been implemented: %s", node.Op.Type.String()),
			span,
		)
	}
}

func (c *Checker) interpolatedStringLiteral(node *ast.InterpolatedStringLiteralNode) {
	for _, contentSection := range node.Content {
		c.checkStringContent(contentSection)
	}
}

func (c *Checker) checkStringContent(node ast.StringLiteralContentNode) {
	switch n := node.(type) {
	case *ast.StringInspectInterpolationNode:
		c.checkExpression(n.Expression)
	case *ast.StringInterpolationNode:
		c.checkExpression(n.Expression)
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

func (c *Checker) addFailure(message string, span *position.Span) {
	if span == nil {
		return
	}
	c.Errors.AddFailure(
		message,
		c.newLocation(span),
	)
}

func (c *Checker) resolveDeclaredNamespace(constantExpression ast.ExpressionNode) types.Namespace {
	typ := c.resolveDeclaredConstant(constantExpression)
	switch t := typ.(type) {
	case *types.SingletonClass:
		return t.AttachedObject
	case *types.Module:
		return t
	}

	return nil
}

func (c *Checker) resolveDeclaredConstant(constantExpression ast.ExpressionNode) types.Type {
	switch constant := constantExpression.(type) {
	case *ast.PublicConstantNode:
		namespace := c.currentConstScope().container
		return namespace.ConstantString(constant.Value)
	case *ast.PrivateConstantNode:
		namespace := c.currentConstScope().container
		return namespace.ConstantString(constant.Value)
	case *ast.ConstantLookupNode:
		_, typ, _ := c.resolveConstantLookupForDeclaration(constant)
		return typ
	default:
		panic(fmt.Sprintf("invalid constant node: %T", constantExpression))
	}
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

func (c *Checker) registerPlaceholderNamespace(placeholder *types.PlaceholderNamespace) {
	c.placeholderNamespaces.Append(placeholder)
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
		leftContainerType = namespace.ConstantString(l.Value)
		leftContainerName = types.MakeFullConstantName(namespace.Name(), l.Value)
		if leftContainerType == nil {
			placeholder := types.NewPlaceholderNamespace(leftContainerName)
			placeholder.Locations.Append(c.newLocation(l.Span()))
			leftContainerType = placeholder
			c.registerPlaceholderNamespace(placeholder)
			namespace.DefineConstant(value.ToSymbol(l.Value), leftContainerType)
		} else if placeholder, ok := leftContainerType.(*types.PlaceholderNamespace); ok {
			placeholder.Locations.Append(c.newLocation(l.Span()))
		}
	case *ast.PrivateConstantNode:
		namespace := c.currentConstScope().container
		leftContainerType = namespace.ConstantString(l.Value)
		leftContainerName = types.MakeFullConstantName(namespace.Name(), l.Value)
		if leftContainerType == nil {
			placeholder := types.NewPlaceholderNamespace(leftContainerName)
			placeholder.Locations.Append(c.newLocation(l.Span()))
			leftContainerType = placeholder
			c.registerPlaceholderNamespace(placeholder)
			namespace.DefineConstant(value.ToSymbol(l.Value), leftContainerType)
		} else if placeholder, ok := leftContainerType.(*types.PlaceholderNamespace); ok {
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
	case *types.PlaceholderNamespace:
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
	constant := leftContainer.Constant(rightSymbol)
	if constant == nil && !firstCall {
		placeholder := types.NewPlaceholderNamespace(constantName)
		placeholder.Locations.Append(c.newLocation(node.Right.Span()))
		constant = placeholder
		c.registerPlaceholderNamespace(placeholder)
		leftContainer.DefineConstant(rightSymbol, constant)
	} else if placeholder, ok := constant.(*types.PlaceholderNamespace); ok {
		placeholder.Locations.Append(c.newLocation(node.Right.Span()))
	}

	return leftContainer, constant, constantName
}

// Get the type of the constant with the given name
func (c *Checker) resolveSimpleConstantForSetter(name string) (types.Namespace, types.Type, string) {
	namespace := c.currentConstScope().container
	constant := namespace.ConstantString(name)
	fullName := types.MakeFullConstantName(namespace.Name(), name)
	if constant != nil {
		return namespace, constant, fullName
	}
	return namespace, nil, fullName
}

// Get the type of the public constant with the given name
func (c *Checker) resolvePublicConstant(name string, span *position.Span) (types.Type, string) {
	for i := range len(c.constantScopes) {
		constScope := c.constantScopes[len(c.constantScopes)-i-1]
		constant := constScope.container.ConstantString(name)
		if constant != nil {
			return constant, types.MakeFullConstantName(constScope.container.Name(), name)
		}
	}

	c.addFailure(
		fmt.Sprintf("undefined constant `%s`", name),
		span,
	)
	return nil, name
}

// Get the type of the private constant with the given name
func (c *Checker) resolvePrivateConstant(name string, span *position.Span) (types.Type, string) {
	for i := range len(c.constantScopes) {
		constScope := c.constantScopes[len(c.constantScopes)-i-1]
		if !constScope.local {
			continue
		}
		constant := constScope.container.ConstantString(name)
		if constant != nil {
			return constant, types.MakeFullConstantName(constScope.container.Name(), name)
		}
	}

	c.addFailure(
		fmt.Sprintf("undefined constant `%s`", name),
		span,
	)
	return nil, name
}

// Get the type with the given name
func (c *Checker) resolveType(name string, span *position.Span) (types.Type, string) {
	for i := range len(c.constantScopes) {
		constScope := c.constantScopes[len(c.constantScopes)-i-1]
		constant := constScope.container.SubtypeString(name)
		if constant != nil {
			return constant, types.MakeFullConstantName(constScope.container.Name(), name)
		}
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
	return env.getLocal(name)
}

// Get the instance variable with the specified name
func (c *Checker) getInstanceVariableIn(name value.Symbol, typ types.Namespace) (types.Type, types.Namespace) {
	currentContainer := typ
	for currentContainer != nil {
		ivar := currentContainer.InstanceVariable(name)
		if ivar != nil {
			return ivar, currentContainer
		}

		currentContainer = currentContainer.Parent()
	}

	return nil, typ
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
func (c *Checker) resolveLocal(name string, span *position.Span) *local {
	env := c.currentLocalEnv()
	local := env.resolveLocal(name)
	if local == nil {
		c.addFailure(
			fmt.Sprintf("undefined local `%s`", name),
			span,
		)
	}
	return local
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

	constant := leftContainer.SubtypeString(rightName)
	if constant == nil {
		c.addFailure(
			fmt.Sprintf("undefined type `%s`", typeName),
			node.Right.Span(),
		)
		return nil, typeName
	}

	return constant, typeName
}

func (c *Checker) checkComplexConstantType(node ast.ComplexConstantNode) ast.ComplexConstantNode {
	switch n := node.(type) {
	case *ast.PublicConstantNode:
		c.checkPublicConstantType(n)
		return n
	case *ast.PrivateConstantNode:
		c.checkPrivateConstantType(n)
		return n
	case *ast.ConstantLookupNode:
		return c.constantLookupType(n)
	default:
		c.addFailure(
			fmt.Sprintf("invalid constant type node %T", node),
			node.Span(),
		)
		return n
	}
}

func (c *Checker) checkPublicConstantType(node *ast.PublicConstantNode) {
	typ, _ := c.resolveType(node.Value, node.Span())
	if typ == nil {
		typ = types.Void{}
	}
	node.SetType(typ)
}

func (c *Checker) checkPrivateConstantType(node *ast.PrivateConstantNode) {
	typ, _ := c.resolveType(node.Value, node.Span())
	if typ == nil {
		typ = types.Void{}
	}
	node.SetType(typ)
}

func (c *Checker) checkTypeNode(node ast.TypeNode) ast.TypeNode {
	switch n := node.(type) {
	case *ast.PublicConstantNode:
		c.checkPublicConstantType(n)
		return n
	case *ast.PrivateConstantNode:
		c.checkPrivateConstantType(n)
		return n
	case *ast.ConstantLookupNode:
		c.constantLookupType(n)
		return n
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
		switch n.Op.Type {
		case token.OR:
			return c.constructUnionType(n)
		case token.AND:
			return c.constructIntersectionType(n)
		default:
			panic("invalid binary type expression operator")
		}
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
	case *ast.TrueLiteralNode, *ast.FalseLiteralNode, *ast.VoidTypeNode,
		*ast.NeverTypeNode, *ast.AnyTypeNode, *ast.NilLiteralNode:
		return n
	case *ast.NilableTypeNode:
		n.TypeNode = c.checkTypeNode(n.TypeNode)
		typ := c.toNilable(c.typeOf(n.TypeNode))
		n.SetType(typ)
		return n
	case *ast.SingletonTypeNode:
		n.TypeNode = c.checkTypeNode(n.TypeNode)
		typ := c.typeOf(n.TypeNode)
		t, ok := typ.(types.Namespace)
		if !ok {
			c.addFailure(
				fmt.Sprintf("cannot get singleton class of `%s`", types.InspectWithColor(typ)),
				n.Span(),
			)
			n.SetType(types.Void{})
			return n
		}

		singleton := t.Singleton()
		if singleton == nil {
			c.addFailure(
				fmt.Sprintf("cannot get singleton class of `%s`", types.InspectWithColor(typ)),
				n.Span(),
			)
			n.SetType(types.Void{})
			return n
		}

		n.SetType(singleton)
		return n
	default:
		c.addFailure(
			fmt.Sprintf("invalid type node %T", node),
			node.Span(),
		)
		return n
	}
}

func (c *Checker) constructUnionType(node *ast.BinaryTypeExpressionNode) *ast.UnionTypeNode {
	union := types.NewUnion()
	elements := new([]ast.TypeNode)
	c._constructUnionType(node, elements, union)

	newNode := ast.NewUnionTypeNode(
		node.Span(),
		*elements,
	)
	newNode.SetType(union)
	return newNode
}

func (c *Checker) _constructUnionType(node *ast.BinaryTypeExpressionNode, elements *[]ast.TypeNode, union *types.Union) {
	leftBinaryType, leftIsBinaryType := node.Left.(*ast.BinaryTypeExpressionNode)
	if leftIsBinaryType && leftBinaryType.Op.Type == token.OR {
		c._constructUnionType(leftBinaryType, elements, union)
	} else {
		leftTypeNode := node.Left
		leftTypeNode = c.checkTypeNode(leftTypeNode)
		*elements = append(*elements, leftTypeNode)

		leftType := c.typeOf(leftTypeNode)
		union.Elements = append(union.Elements, leftType)
	}

	rightBinaryType, rightIsBinaryType := node.Right.(*ast.BinaryTypeExpressionNode)
	if rightIsBinaryType && rightBinaryType.Op.Type == token.OR {
		c._constructUnionType(rightBinaryType, elements, union)
	} else {
		rightTypeNode := node.Right
		rightTypeNode = c.checkTypeNode(rightTypeNode)
		*elements = append(*elements, rightTypeNode)

		rightType := c.typeOf(rightTypeNode)
		union.Elements = append(union.Elements, rightType)
	}
}

func (c *Checker) constructIntersectionType(node *ast.BinaryTypeExpressionNode) *ast.IntersectionTypeNode {
	intersection := types.NewIntersection()
	elements := new([]ast.TypeNode)
	c._constructIntersectionType(node, elements, intersection)
	newNode := ast.NewIntersectionTypeNode(
		node.Span(),
		*elements,
	)
	newNode.SetType(intersection)
	return newNode
}

func (c *Checker) _constructIntersectionType(node *ast.BinaryTypeExpressionNode, elements *[]ast.TypeNode, intersection *types.Intersection) {
	leftBinaryType, leftIsBinaryType := node.Left.(*ast.BinaryTypeExpressionNode)
	if leftIsBinaryType && leftBinaryType.Op.Type == token.AND {
		c._constructIntersectionType(leftBinaryType, elements, intersection)
	} else {
		leftTypeNode := node.Left
		leftTypeNode = c.checkTypeNode(leftTypeNode)
		*elements = append(*elements, leftTypeNode)

		leftType := c.typeOf(leftTypeNode)
		intersection.Elements = append(intersection.Elements, leftType)
	}

	rightBinaryType, rightIsBinaryType := node.Right.(*ast.BinaryTypeExpressionNode)
	if rightIsBinaryType && rightBinaryType.Op.Type == token.AND {
		c._constructIntersectionType(rightBinaryType, elements, intersection)
	} else {
		rightTypeNode := node.Right
		rightTypeNode = c.checkTypeNode(rightTypeNode)
		*elements = append(*elements, rightTypeNode)

		rightType := c.typeOf(rightTypeNode)
		intersection.Elements = append(intersection.Elements, rightType)
	}
}

func (c *Checker) constantLookupType(node *ast.ConstantLookupNode) *ast.PublicConstantNode {
	typ, name := c.resolveConstantLookupType(node)
	if typ == nil {
		typ = types.Void{}
	}

	newNode := ast.NewPublicConstantNode(
		node.Span(),
		name,
	)
	newNode.SetType(typ)
	return newNode
}

func (c *Checker) resolveConstantType(constantExpression ast.ExpressionNode) (types.Type, string) {
	switch constant := constantExpression.(type) {
	case *ast.PublicConstantNode:
		return c.resolveType(constant.Value, constant.Span())
	case *ast.PrivateConstantNode:
		return c.resolveType(constant.Value, constant.Span())
	case *ast.ConstantLookupNode:
		return c.resolveConstantLookupType(constant)
	default:
		panic(fmt.Sprintf("invalid constant node: %T", constantExpression))
	}
}

func (c *Checker) resolveConstant(constantExpression ast.ExpressionNode) (types.Type, string) {
	switch constant := constantExpression.(type) {
	case *ast.PublicConstantNode:
		return c.resolvePublicConstant(constant.Value, constant.Span())
	case *ast.PrivateConstantNode:
		return c.resolvePrivateConstant(constant.Value, constant.Span())
	case *ast.ConstantLookupNode:
		return c.resolveConstantLookup(constant)
	default:
		panic(fmt.Sprintf("invalid constant node: %T", constantExpression))
	}
}

func (c *Checker) resolveConstantLookup(node *ast.ConstantLookupNode) (types.Type, string) {
	var leftContainerType types.Type
	var leftContainerName string

	switch l := node.Left.(type) {
	case *ast.PublicConstantNode:
		leftContainerType, leftContainerName = c.resolvePublicConstant(l.Value, l.Span())
	case *ast.PrivateConstantNode:
		leftContainerType, leftContainerName = c.resolvePrivateConstant(l.Value, l.Span())
	case nil:
		leftContainerType = c.GlobalEnv.Root
	case *ast.ConstantLookupNode:
		leftContainerType, leftContainerName = c.resolveConstantLookup(l)
	default:
		c.addFailure(
			fmt.Sprintf("invalid constant node %T", node),
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
			fmt.Sprintf("cannot read private constant `%s`", rightName),
			node.Span(),
		)
	default:
		c.addFailure(
			fmt.Sprintf("invalid constant node %T", node),
			node.Span(),
		)
		return nil, ""
	}

	constantName := types.MakeFullConstantName(leftContainerName, rightName)
	if leftContainerType == nil {
		return nil, constantName
	}

	var leftContainer types.Namespace
	switch l := leftContainerType.(type) {
	case *types.Module:
		leftContainer = l
	case *types.SingletonClass:
		leftContainer = l.AttachedObject
	default:
		c.addFailure(
			fmt.Sprintf("cannot read constants from `%s`, it is not a constant container", leftContainerName),
			node.Span(),
		)
		return nil, constantName
	}

	constant := leftContainer.ConstantString(rightName)
	if constant == nil {
		c.addFailure(
			fmt.Sprintf("undefined constant `%s`", constantName),
			node.Right.Span(),
		)
		return nil, constantName
	}

	return constant, constantName
}

func (c *Checker) constantLookup(node *ast.ConstantLookupNode) *ast.PublicConstantNode {
	typ, name := c.resolveConstantLookup(node)
	if typ == nil {
		typ = types.Void{}
	}

	newNode := ast.NewPublicConstantNode(
		node.Span(),
		name,
	)
	newNode.SetType(typ)
	return newNode
}

func (c *Checker) publicConstant(node *ast.PublicConstantNode) *ast.PublicConstantNode {
	typ, name := c.resolvePublicConstant(node.Value, node.Span())
	if typ == nil {
		typ = types.Void{}
	}

	node.Value = name
	node.SetType(typ)
	return node
}

func (c *Checker) privateConstant(node *ast.PrivateConstantNode) *ast.PrivateConstantNode {
	typ, name := c.resolvePrivateConstant(node.Value, node.Span())
	if typ == nil {
		typ = types.Void{}
	}

	node.Value = name
	node.SetType(typ)
	return node
}

func (c *Checker) publicIdentifier(node *ast.PublicIdentifierNode) *ast.PublicIdentifierNode {
	local := c.resolveLocal(node.Value, node.Span())
	if local == nil {
		node.SetType(types.Void{})
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

func (c *Checker) privateIdentifier(node *ast.PrivateIdentifierNode) *ast.PrivateIdentifierNode {
	local := c.resolveLocal(node.Value, node.Span())
	if local == nil {
		node.SetType(types.Void{})
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

func (c *Checker) instanceVariable(node *ast.InstanceVariableNode) {
	typ, container := c.getInstanceVariable(value.ToSymbol(node.Value))
	self, ok := c.selfType.(types.Namespace)
	if !ok || self.IsPrimitive() {
		c.addFailure(
			"cannot use instance variables in this context",
			node.Span(),
		)
	}

	if typ == nil {
		c.addFailure(
			fmt.Sprintf(
				"undefined instance variable `%s` in type `%s`",
				types.InspectInstanceVariableWithColor(node.Value),
				types.InspectWithColor(container),
			),
			node.Span(),
		)
		node.SetType(types.Void{})
		return
	}

	node.SetType(typ)
}

func (c *Checker) declareMethodForGetter(node *ast.AttributeParameterNode, docComment string) {
	method := c.declareMethod(
		docComment,
		false,
		false,
		value.ToSymbol(node.Name),
		nil,
		node.TypeNode,
		nil,
		node.Span(),
	)

	init := node.Initialiser
	var body []ast.StatementNode

	if init == nil {
		body = ast.ExpressionToStatements(
			ast.NewInstanceVariableNode(node.Span(), node.Name),
		)
	} else {
		body = ast.ExpressionToStatements(
			ast.NewAssignmentExpressionNode(
				node.Span(),
				token.New(init.Span(), token.QUESTION_QUESTION_EQUAL),
				ast.NewInstanceVariableNode(node.Span(), node.Name),
				init,
			),
		)
	}

	methodNode := ast.NewMethodDefinitionNode(
		node.Span(),
		"",
		false,
		false,
		node.Name,
		nil,
		node.TypeNode,
		nil,
		body,
	)
	methodNode.SetType(method)
	c.registerMethodCheck(
		method,
		methodNode,
	)
}

func (c *Checker) declareMethodForSetter(node *ast.AttributeParameterNode, docComment string) {
	setterName := node.Name + "="

	params := []ast.ParameterNode{
		ast.NewMethodParameterNode(
			node.TypeNode.Span(),
			node.Name,
			true,
			node.TypeNode,
			nil,
			ast.NormalParameterKind,
		),
	}
	method := c.declareMethod(
		docComment,
		false,
		false,
		value.ToSymbol(setterName),
		params,
		nil,
		nil,
		node.Span(),
	)

	methodNode := ast.NewMethodDefinitionNode(
		node.Span(),
		docComment,
		false,
		false,
		setterName,
		params,
		nil,
		nil,
		nil,
	)
	methodNode.SetType(method)
	c.registerMethodCheck(
		method,
		methodNode,
	)
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

		declaredType = types.Void{}
	} else {
		declaredTypeNode := c.checkTypeNode(node.TypeNode)
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
		return
	}

	c.declareInstanceVariable(value.ToSymbol(node.Name), declaredType, node.Span())
}

func (c *Checker) variableDeclaration(node *ast.VariableDeclarationNode) {
	if variable := c.getLocal(node.Name); variable != nil {
		c.addFailure(
			fmt.Sprintf("cannot redeclare local `%s`", node.Name),
			node.Span(),
		)
	}
	if node.Initialiser == nil {
		if node.TypeNode == nil {
			c.addFailure(
				fmt.Sprintf("cannot declare a variable without a type `%s`", node.Name),
				node.Span(),
			)
			c.addLocal(node.Name, newLocal(types.Void{}, false, false))
			node.SetType(types.Void{})
			return
		}

		// without an initialiser but with a type
		declaredTypeNode := c.checkTypeNode(node.TypeNode)
		declaredType := c.typeOf(declaredTypeNode)
		c.addLocal(node.Name, newLocal(declaredType, false, false))
		node.TypeNode = declaredTypeNode
		node.SetType(types.Void{})
		return
	}

	// with an initialiser
	if node.TypeNode == nil {
		// without a type, inference
		init := c.checkExpression(node.Initialiser)
		actualType := c.toNonLiteral(c.typeOf(init))
		c.addLocal(node.Name, newLocal(actualType, true, false))
		if types.IsVoid(actualType) {
			c.addFailure(
				fmt.Sprintf("cannot declare variable `%s` with type `void`", node.Name),
				init.Span(),
			)
		}
		node.Initialiser = init
		node.SetType(actualType)
		return
	}

	// with a type and an initializer

	declaredTypeNode := c.checkTypeNode(node.TypeNode)
	declaredType := c.typeOf(declaredTypeNode)
	init := c.checkExpression(node.Initialiser)
	actualType := c.typeOf(init)
	c.addLocal(node.Name, newLocal(declaredType, true, false))
	c.checkCanAssign(actualType, declaredType, init.Span())

	node.TypeNode = declaredTypeNode
	node.Initialiser = init
	node.SetType(declaredType)
}

func (c *Checker) valueDeclaration(node *ast.ValueDeclarationNode) {
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
			c.addLocal(node.Name, newLocal(types.Void{}, false, true))
			node.SetType(types.Void{})
			return
		}

		// without an initialiser but with a type
		declaredTypeNode := c.checkTypeNode(node.TypeNode)
		declaredType := c.typeOf(declaredTypeNode)
		c.addLocal(node.Name, newLocal(declaredType, false, true))
		node.TypeNode = declaredTypeNode
		node.SetType(types.Void{})
		return
	}

	// with an initialiser
	if node.TypeNode == nil {
		// without a type, inference
		init := c.checkExpression(node.Initialiser)
		actualType := c.typeOf(init)
		c.addLocal(node.Name, newLocal(actualType, true, true))
		if types.IsVoid(actualType) {
			c.addFailure(
				fmt.Sprintf("cannot declare value `%s` with type `void`", node.Name),
				init.Span(),
			)
		}
		node.Initialiser = init
		node.SetType(actualType)
		return
	}

	// with a type and an initializer

	declaredTypeNode := c.checkTypeNode(node.TypeNode)
	declaredType := c.typeOf(declaredTypeNode)
	init := c.checkExpression(node.Initialiser)
	actualType := c.typeOf(init)
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
		case *types.PlaceholderNamespace:
			module := types.NewModuleWithDetails(
				docComment,
				t.Name(),
				t.Constants(),
				t.Subtypes(),
				t.Methods(),
			)
			t.Replacement = module
			namespace.DefineConstant(constantName, module)
			return module
		default:
			c.addFailure(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				span,
			)
			return types.NewModule(docComment, fullConstantName)
		}
	} else if namespace == nil {
		return types.NewModule(docComment, fullConstantName)
	} else {
		return namespace.DefineModule(docComment, constantName)
	}
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
		ct, ok := constantType.(*types.SingletonClass)
		if !ok {
			c.addFailure(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				span,
			)
			return types.NewClass(docComment, abstract, sealed, primitive, fullConstantName, nil, c.GlobalEnv)
		}
		constantType = ct.AttachedObject

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
		case *types.PlaceholderNamespace:
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
			t.Replacement = class
			namespace.DefineConstant(constantName, class)
			return class
		default:
			c.addFailure(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				span,
			)
			return types.NewClass(docComment, abstract, sealed, primitive, fullConstantName, nil, c.GlobalEnv)
		}
	} else if namespace == nil {
		return types.NewClass(docComment, abstract, sealed, primitive, fullConstantName, nil, c.GlobalEnv)
	} else {
		return namespace.DefineClass(docComment, abstract, sealed, primitive, constantName, nil, c.GlobalEnv)
	}
}

func (c *Checker) checkProgram(node *ast.ProgramNode) *vm.BytecodeFunction {
	statements := node.Body

	c.hoistTypeDefinitions(statements)

	for _, placeholder := range c.placeholderNamespaces.Slice {
		replacement := placeholder.Replacement
		if replacement != nil {
			continue
		}

		for _, location := range placeholder.Locations.Slice {
			c.addFailureWithLocation(
				fmt.Sprintf("undefined namespace `%s`", placeholder.Name()),
				location,
			)
		}
	}

	c.hoistMethodDefinitions(statements)
	c.checkMethods()
	c.checkStatements(statements)

	return nil
}

const concurrencyLimit = 10_000

func (c *Checker) checkMethods() {
	concurrent.Foreach(
		concurrencyLimit,
		c.methodChecks.Slice,
		func(methodCheck methodCheckEntry) {
			methodChecker := c.newChecker(
				methodCheck.filename,
				methodCheck.constantScopes,
				methodCheck.methodScopes,
				methodCheck.method.DefinedUnder,
				methodCheck.method.ReturnType,
				methodCheck.method.ThrowType,
			)
			node := methodCheck.node
			methodChecker.checkMethodDefinition(node)
		},
	)
}

func (c *Checker) checkConstantDeclaration(node *ast.ConstantDeclarationNode) {
	container, constant, fullConstantName := c.resolveConstantForDeclaration(node.Constant)
	constantName := value.ToSymbol(extractConstantName(node.Constant))
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)

	if constant != nil {
		c.addFailure(
			fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
			node.Span(),
		)
	}

	if node.Initialiser == nil {
		if !c.HeaderMode {
			c.addFailure(
				fmt.Sprintf("constant `%s` has to be initialised", constantName),
				node.Span(),
			)
		}
		if node.TypeNode == nil {
			c.addFailure(
				fmt.Sprintf("cannot declare a constant without a type `%s`", constantName),
				node.Span(),
			)
			container.DefineConstant(constantName, types.Void{})
			node.SetType(types.Void{})
			return
		}

		// without an initialiser but with a type
		declaredTypeNode := c.checkTypeNode(node.TypeNode)
		declaredType := c.typeOf(declaredTypeNode)
		container.DefineConstant(constantName, declaredType)
		node.TypeNode = declaredTypeNode
		node.SetType(types.Void{})
		return
	}
	init := c.checkExpression(node.Initialiser)

	if !init.IsStatic() {
		c.addFailure(
			"values assigned to constants must be static, known at compile time",
			init.Span(),
		)
	}

	// with an initialiser
	if node.TypeNode == nil {
		// without a type, inference
		actualType := c.typeOf(init)
		if types.IsVoid(actualType) {
			c.addFailure(
				fmt.Sprintf("cannot declare constant `%s` with type `void`", fullConstantName),
				init.Span(),
			)
		}
		node.Initialiser = init
		node.SetType(actualType)
		container.DefineConstant(constantName, actualType)
		return
	}

	// with a type and an initializer

	declaredTypeNode := c.checkTypeNode(node.TypeNode)
	declaredType := c.typeOf(declaredTypeNode)
	actualType := c.typeOf(init)
	c.checkCanAssign(actualType, declaredType, init.Span())

	container.DefineConstant(constantName, declaredType)
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

func (c *Checker) hoistStructDeclaration(structNode *ast.StructDeclarationNode) *ast.ClassDeclarationNode {
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
	return classNode
}

func (c *Checker) hoistModuleDeclaration(node *ast.ModuleDeclarationNode) {
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

	c.pushConstScope(makeLocalConstantScope(module))
	c.pushMethodScope(makeLocalMethodScope(module))

	c.hoistTypeDefinitions(node.Body)

	c.popConstScope()
	c.popMethodScope()
}

func (c *Checker) hoistClassDeclaration(node *ast.ClassDeclarationNode) {
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
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)

	c.pushConstScope(makeLocalConstantScope(class))
	c.pushMethodScope(makeLocalMethodScope(class))

	c.hoistTypeDefinitions(node.Body)

	c.popConstScope()
	c.popMethodScope()
}

func (c *Checker) hoistMixinDeclaration(node *ast.MixinDeclarationNode) {
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
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)

	c.pushConstScope(makeLocalConstantScope(mixin))
	c.pushMethodScope(makeLocalMethodScope(mixin))

	c.hoistTypeDefinitions(node.Body)

	c.popConstScope()
	c.popMethodScope()
}

func (c *Checker) hoistInterfaceDeclaration(node *ast.InterfaceDeclarationNode) {
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
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)

	c.pushConstScope(makeLocalConstantScope(iface))
	c.pushMethodScope(makeLocalMethodScope(iface))

	c.hoistTypeDefinitions(node.Body)

	c.popConstScope()
	c.popMethodScope()
}

func (c *Checker) hoistTypeDefinition(node *ast.TypeDefinitionNode) {
	container, constant, fullConstantName := c.resolveConstantForDeclaration(node.Constant)
	constantName := value.ToSymbol(extractConstantName(node.Constant))
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)
	if constant != nil {
		c.addFailure(
			fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
			node.Constant.Span(),
		)
	}

	node.TypeNode = c.checkTypeNode(node.TypeNode)

	typ := c.typeOf(node.TypeNode)
	container.DefineConstant(constantName, types.Void{})
	namedType := types.NewNamedType(fullConstantName, typ)
	container.DefineSubtype(constantName, namedType)
}

func (c *Checker) hoistTypeDefinitions(statements []ast.StatementNode) {
	for _, statement := range statements {
		stmt, ok := statement.(*ast.ExpressionStatementNode)
		if !ok {
			continue
		}
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
		case *ast.ConstantDeclarationNode:
			c.checkConstantDeclaration(expr)
		case *ast.TypeDefinitionNode:
			c.hoistTypeDefinition(expr)
		}
	}
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
	method := c.declareMethod(
		initNode.DocComment(),
		false,
		false,
		symbol.M_init,
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
		initNode.Parameters,
		nil,
		initNode.ThrowType,
		initNode.Body,
	)
	newNode.SetType(method)
	c.registerMethodCheck(method, newNode)
	return newNode
}

func (c *Checker) hoistAliasDeclaration(node *ast.AliasDeclarationNode) {
	namespace := c.currentMethodScope().container
	for _, entry := range node.Entries {
		method := types.GetMethodInNamespace(namespace, value.ToSymbol(entry.OldName))
		if method == nil {
			c.addMissingMethodError(namespace, entry.OldName, entry.Span())
			continue
		}
		namespace.SetMethod(value.ToSymbol(entry.NewName), method)
	}
}

func (c *Checker) hoistMethodDefinition(node *ast.MethodDefinitionNode) {
	method := c.declareMethod(
		node.DocComment(),
		node.Abstract,
		node.Sealed,
		value.ToSymbol(node.Name),
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Span(),
	)
	node.SetType(method)
	c.registerMethodCheck(method, node)
}

func (c *Checker) hoistMethodSignatureDefinition(node *ast.MethodSignatureDefinitionNode) {
	method := c.declareMethod(
		node.DocComment(),
		true,
		false,
		value.ToSymbol(node.Name),
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Span(),
	)
	node.SetType(method)
}

func (c *Checker) hoistMethodDefinitionsWithinClass(node *ast.ClassDeclarationNode) {
	class, ok := c.typeOf(node).(*types.Class)
	if ok {
		c.pushConstScope(makeLocalConstantScope(class))
		c.pushMethodScope(makeLocalMethodScope(class))

		var superclass *types.Class

		if node.Superclass == nil {
			superclass = c.GlobalEnv.StdSubtypeClass(symbol.Object)
		} else {
			superclassType, _ := c.resolveConstantType(node.Superclass)
			var ok bool
			superclass, ok = superclassType.(*types.Class)
			if !ok {
				c.addFailure(
					fmt.Sprintf("`%s` is not a class", types.InspectWithColor(superclassType)),
					node.Superclass.Span(),
				)
			} else {
				if superclass.IsSealed() {
					c.addFailure(
						fmt.Sprintf("cannot inherit from sealed class `%s`", types.InspectWithColor(superclassType)),
						node.Superclass.Span(),
					)
				}
				if superclass.IsPrimitive() && !class.IsPrimitive() {
					c.addFailure(
						fmt.Sprintf("class `%s` must be primitive to inherit from primitive class `%s`", types.InspectWithColor(class), types.InspectWithColor(superclassType)),
						node.Superclass.Span(),
					)
				}
			}
		}

		parent := class.Superclass()
		if parent == nil && superclass != nil {
			class.SetParent(superclass)
		} else if parent != nil && parent != superclass {
			var span *position.Span
			if node.Superclass == nil {
				span = node.Span()
			} else {
				span = node.Superclass.Span()
			}

			c.addFailure(
				fmt.Sprintf(
					"superclass mismatch in `%s`, got `%s`, expected `%s`",
					class.Name(),
					superclass.Name(),
					parent.Name(),
				),
				span,
			)
		}

	}

	previousMode := c.mode
	c.mode = classMode
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)
	if ok {
		c.popConstScope()
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
	c.mode = moduleMode
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)

	if ok {
		c.popConstScope()
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
	c.mode = mixinMode
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)

	if ok {
		c.popConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) hoistMethodDefinitionsWithinInterface(node *ast.InterfaceDeclarationNode) {
	mixin, ok := c.typeOf(node).(*types.Interface)
	if ok {
		c.pushConstScope(makeLocalConstantScope(mixin))
		c.pushMethodScope(makeLocalMethodScope(mixin))
	}

	previousMode := c.mode
	c.mode = interfaceMode
	c.hoistMethodDefinitions(node.Body)
	c.setMode(previousMode)

	if ok {
		c.popConstScope()
		c.popMethodScope()
	}
}

func (c *Checker) hoistMethodDefinitionsWithinSingleton(expr *ast.SingletonBlockExpressionNode) {
	namespace := c.currentConstScope().container
	singleton := namespace.Singleton()
	if singleton == nil {
		c.addFailure(
			"cannot declare a singleton class in this context",
			expr.Span(),
		)
		return
	}
	expr.SetType(singleton)

	c.pushConstScope(makeLocalConstantScope(singleton))
	c.pushMethodScope(makeLocalMethodScope(singleton))

	previousMode := c.mode
	c.mode = singletonMode
	c.hoistMethodDefinitions(expr.Body)
	c.setMode(previousMode)

	c.popConstScope()
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
		case *ast.IncludeExpressionNode:
			for _, constant := range expr.Constants {
				c.includeMixin(constant)
			}
		case *ast.ImplementExpressionNode:
			for _, constant := range expr.Constants {
				c.implementInterface(constant)
			}
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

func (c *Checker) checkMethodDefinition(node *ast.MethodDefinitionNode) {
	returnType, throwType := c.checkMethod(
		c.typeOf(node).(*types.Method),
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Body,
		node.Span(),
	)
	node.ReturnType = returnType
	node.ThrowType = throwType
}

func (c *Checker) declareMethod(
	docComment string,
	abstract bool,
	sealed bool,
	name value.Symbol,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	span *position.Span,
) *types.Method {
	if c.mode == interfaceMode {
		abstract = true
	}
	methodScope := c.currentMethodScope()
	methodNamespace := methodScope.container
	oldMethod := methodNamespace.Method(name)
	if oldMethod != nil {
		if oldMethod.IsNative() && oldMethod.IsSealed() {
			c.addOverrideSealedMethodError(oldMethod, span)
		} else if sealed && !oldMethod.IsSealed() {
			c.addFailure(
				fmt.Sprintf(
					"cannot redeclare method `%s` with a different modifier, is `%s`, should be `%s`",
					name,
					types.InspectModifier(abstract, sealed, false),
					types.InspectModifier(oldMethod.IsAbstract(), oldMethod.IsSealed(), false),
				),
				span,
			)
		}
	}

	switch namespace := methodNamespace.(type) {
	case *types.Interface:
	case *types.Class:
		if abstract && !namespace.IsAbstract() {
			c.addFailure(
				fmt.Sprintf(
					"cannot declare abstract method `%s` in non-abstract class `%s`",
					name,
					types.InspectWithColor(methodNamespace),
				),
				span,
			)
		}
	case *types.Mixin:
		if abstract && !namespace.IsAbstract() {
			c.addFailure(
				fmt.Sprintf(
					"cannot declare abstract method `%s` in non-abstract mixin `%s`",
					name,
					types.InspectWithColor(methodNamespace),
				),
				span,
			)
		}
	default:
		if abstract {
			c.addFailure(
				fmt.Sprintf(
					"cannot declare abstract method `%s` in this context",
					name,
				),
				span,
			)
		}
	}

	var params []*types.Parameter
	for _, param := range paramNodes {
		switch p := param.(type) {
		case *ast.MethodParameterNode:
			var declaredType types.Type
			var declaredTypeNode ast.TypeNode
			if p.SetInstanceVariable {
				currentIvar, _ := c.getInstanceVariableIn(value.ToSymbol(p.Name), methodNamespace)
				if p.TypeNode == nil {
					if currentIvar == nil {
						c.addFailure(
							fmt.Sprintf(
								"cannot infer the type of instance variable `%s`",
								p.Name,
							),
							p.Span(),
						)
					}

					declaredType = currentIvar
				} else {
					declaredTypeNode = c.checkTypeNode(p.TypeNode)
					declaredType = c.typeOf(declaredTypeNode)
					if currentIvar != nil {
						c.checkCanAssignInstanceVariable(p.Name, declaredType, currentIvar, declaredTypeNode.Span())
					} else {
						c.declareInstanceVariable(value.ToSymbol(p.Name), declaredType, p.Span())
					}
				}
			} else if p.TypeNode == nil {
				c.addFailure(
					fmt.Sprintf("cannot declare parameter `%s` without a type", p.Name),
					param.Span(),
				)
			} else {
				declaredTypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.typeOf(declaredTypeNode)
			}

			var kind types.ParameterKind
			switch p.Kind {
			case ast.NormalParameterKind:
				kind = types.NormalParameterKind
			case ast.PositionalRestParameterKind:
				kind = types.PositionalRestParameterKind
			case ast.NamedRestParameterKind:
				kind = types.NamedRestParameterKind
			}
			if p.Initialiser != nil {
				kind = types.DefaultValueParameterKind
			}
			name := value.ToSymbol(p.Name)
			params = append(params, types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			))
		case *ast.SignatureParameterNode:
			var declaredType types.Type
			var declaredTypeNode ast.TypeNode
			if p.TypeNode == nil {
				c.addFailure(
					fmt.Sprintf("cannot declare parameter `%s` without a type", p.Name),
					param.Span(),
				)
			} else {
				declaredTypeNode = c.checkTypeNode(p.TypeNode)
				declaredType = c.typeOf(declaredTypeNode)
			}

			var kind types.ParameterKind
			switch p.Kind {
			case ast.NormalParameterKind:
				kind = types.NormalParameterKind
			case ast.PositionalRestParameterKind:
				kind = types.PositionalRestParameterKind
			case ast.NamedRestParameterKind:
				kind = types.NamedRestParameterKind
			}
			if p.Optional {
				kind = types.DefaultValueParameterKind
			}
			name := value.ToSymbol(p.Name)
			params = append(params, types.NewParameter(
				name,
				declaredType,
				kind,
				false,
			))
		default:
			c.addFailure(
				fmt.Sprintf("invalid param type %T", param),
				param.Span(),
			)
		}
	}

	var returnType types.Type
	var typedReturnTypeNode ast.TypeNode
	if returnTypeNode != nil {
		typedReturnTypeNode = c.checkTypeNode(returnTypeNode)
		returnType = c.typeOf(typedReturnTypeNode)
	} else {
		returnType = types.Void{}
	}

	var throwType types.Type
	var typedThrowTypeNode ast.TypeNode
	if throwTypeNode != nil {
		typedThrowTypeNode = c.checkTypeNode(throwTypeNode)
		throwType = c.typeOf(typedThrowTypeNode)
	}
	newMethod := types.NewMethod(
		docComment,
		abstract,
		sealed,
		c.HeaderMode,
		name,
		params,
		returnType,
		throwType,
		methodScope.container,
	)
	newMethod.SetSpan(span)

	methodScope.container.SetMethod(name, newMethod)

	return newMethod
}

func (c *Checker) declareMixin(docComment string, abstract bool, namespace types.Namespace, constantType types.Type, fullConstantName string, constantName value.Symbol, span *position.Span) *types.Mixin {
	if constantType != nil {
		ct, ok := constantType.(*types.SingletonClass)
		if !ok {
			c.addFailure(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				span,
			)
			return types.NewMixin(docComment, abstract, fullConstantName, c.GlobalEnv)
		}
		constantType = ct.AttachedObject

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
		case *types.PlaceholderNamespace:
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
			t.Replacement = mixin
			namespace.DefineConstant(constantName, mixin)
			return mixin
		default:
			c.addFailure(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				span,
			)
			return types.NewMixin(docComment, abstract, fullConstantName, c.GlobalEnv)
		}
	} else if namespace == nil {
		return types.NewMixin(docComment, abstract, fullConstantName, c.GlobalEnv)
	} else {
		return namespace.DefineMixin(docComment, abstract, constantName, c.GlobalEnv)
	}
}

func (c *Checker) declareInterface(docComment string, namespace types.Namespace, constantType types.Type, fullConstantName string, constantName value.Symbol, span *position.Span) *types.Interface {
	if constantType != nil {
		ct, ok := constantType.(*types.SingletonClass)
		if !ok {
			c.addFailure(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				span,
			)
			return types.NewInterface(docComment, fullConstantName, c.GlobalEnv)
		}
		constantType = ct.AttachedObject

		switch t := constantType.(type) {
		case *types.Interface:
			t.AppendDocComment(docComment)
			return t
		case *types.PlaceholderNamespace:
			iface := types.NewInterfaceWithDetails(
				t.Name(),
				nil,
				t.Constants(),
				t.Subtypes(),
				t.Methods(),
			)
			t.Replacement = iface
			namespace.DefineConstant(constantName, iface)
			return iface
		default:
			c.addFailure(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				span,
			)
			return types.NewInterface(docComment, fullConstantName, c.GlobalEnv)
		}
	} else if namespace == nil {
		return types.NewInterface(docComment, fullConstantName, c.GlobalEnv)
	} else {
		return namespace.DefineInterface(docComment, constantName, c.GlobalEnv)
	}
}
