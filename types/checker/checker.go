// Package checker implements the Elk type checker
package checker

import (
	"fmt"

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

// Contains definitions of local variables and values
type localEnvironment struct {
	parent *localEnvironment
	locals map[value.Symbol]local
}

// Get the local with the specified name from this local environment
func (l *localEnvironment) getLocal(name string) (local, bool) {
	local, ok := l.locals[value.ToSymbol(name)]
	return local, ok
}

// Resolve the local with the given name from this local environment or any parent environment
func (l *localEnvironment) resolveLocal(name string) (local, bool) {
	nameSymbol := value.ToSymbol(name)
	currentEnv := l
	for {
		if currentEnv == nil {
			return local{}, false
		}
		loc, ok := currentEnv.locals[nameSymbol]
		if ok {
			return loc, true
		}
		currentEnv = currentEnv.parent
	}
}

func newLocalEnvironment(parent *localEnvironment) *localEnvironment {
	return &localEnvironment{
		parent: parent,
		locals: make(map[value.Symbol]local),
	}
}

type constantScope struct {
	container types.ConstantContainer
	local     bool
}

func makeLocalConstantScope(container types.ConstantContainer) constantScope {
	return constantScope{
		container: container,
		local:     true,
	}
}

func makeConstantScope(container types.ConstantContainer) constantScope {
	return constantScope{
		container: container,
		local:     false,
	}
}

type methodScope struct {
	container types.ConstantContainer
	local     bool
}

func makeLocalMethodScope(container types.ConstantContainer) methodScope {
	return methodScope{
		container: container,
		local:     true,
	}
}

func makeMethodScope(container types.ConstantContainer) methodScope {
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
	methodMode
)

type methodCheckEntry struct {
	filename       string
	method         *types.Method
	constantScopes []constantScope
	node           *ast.MethodDefinitionNode
}

type initCheckEntry struct {
	filename       string
	method         *types.Method
	constantScopes []constantScope
	node           *ast.InitDefinitionNode
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
	localEnvs               []*localEnvironment
	returnType              types.Type
	throwType               types.Type
	selfType                types.Type
	mode                    mode
	placeholderNamespaces   *concurrent.Slice[*types.PlaceholderNamespace]
	methodChecks            *concurrent.Slice[methodCheckEntry]
	initChecks              *concurrent.Slice[initCheckEntry]
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
		initChecks:            concurrent.NewSlice[initCheckEntry](),
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
		initChecks:            concurrent.NewSlice[initCheckEntry](),
	}
}

func (c *Checker) newChecker(filename string, constScopes []constantScope, selfType, returnType, throwType types.Type) *Checker {
	return &Checker{
		GlobalEnv:      c.GlobalEnv,
		Filename:       filename,
		mode:           methodMode,
		selfType:       selfType,
		returnType:     returnType,
		throwType:      throwType,
		constantScopes: constScopes,
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
	c.initChecks = concurrent.NewSlice[initCheckEntry]()
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
}

func (c *Checker) pushMethodScope(methodScope methodScope) {
	c.methodScopes = append(c.methodScopes, methodScope)
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

func (c *Checker) registerInitCheck(method *types.Method, node *ast.InitDefinitionNode) {
	c.initChecks.Append(initCheckEntry{
		method:         method,
		constantScopes: c.constantScopesCopy(),
		node:           node,
		filename:       c.Filename,
	})
}

func (c *Checker) registerMethodCheck(method *types.Method, node *ast.MethodDefinitionNode) {
	c.methodChecks.Append(methodCheckEntry{
		method:         method,
		constantScopes: c.constantScopesCopy(),
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
		c.addError(
			fmt.Sprintf("incorrect statement type %#v", node),
			node.Span(),
		)
	}
}

func (c *Checker) isSubtype(a, b types.Type) bool {
	if types.IsNever(a) {
		return true
	}
	if types.IsAny(b) || types.IsVoid(b) {
		return true
	}
	if types.IsAny(a) || types.IsVoid(a) {
		return false
	}
	aNonLiteral := a.ToNonLiteral(c.GlobalEnv)
	if a != aNonLiteral && c.isSubtype(aNonLiteral, b) {
		return true
	}

	if aNilable, aIsNilable := a.(*types.Nilable); aIsNilable {
		return c.isSubtype(aNilable.Type, b) && c.isSubtype(c.GlobalEnv.StdSubtype(symbol.Nil), b)
	}

	if bNilable, bIsNilable := b.(*types.Nilable); bIsNilable {
		return c.isSubtype(a, bNilable.Type) || c.isSubtype(a, c.GlobalEnv.StdSubtype(symbol.Nil))
	}

	if aUnion, aIsUnion := a.(*types.Union); aIsUnion {
		for _, aElement := range aUnion.Elements {
			if !c.isSubtype(aElement, b) {
				return false
			}
		}
		return true
	}

	if bUnion, bIsUnion := b.(*types.Union); bIsUnion {
		for _, bElement := range bUnion.Elements {
			if c.isSubtype(a, bElement) {
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
		return c.classIsSubtype(a, b)
	case *types.Mixin:
		return c.mixinIsSubtype(a, b)
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

func (c *Checker) classIsSubtype(a *types.Class, b types.Type) bool {
	switch b := b.(type) {
	case *types.Class:
		var currentClass types.ConstantContainer = a
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
		return c.isSubtypeOfMixin(a, b)
	default:
		return false
	}
}

func (c *Checker) isSubtypeOfMixin(a types.ConstantContainer, b *types.Mixin) bool {
	var currentContainer types.ConstantContainer = a
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

func (c *Checker) mixinIsSubtype(a *types.Mixin, b types.Type) bool {
	bMixin, ok := b.(*types.Mixin)
	if !ok {
		return false
	}

	return c.isSubtypeOfMixin(a, bMixin)
}

func (c *Checker) checkExpressions(exprs []ast.ExpressionNode) {
	for i, expr := range exprs {
		exprs[i] = c.checkExpression(expr)
	}
}

func (c *Checker) checkExpression(node ast.ExpressionNode) ast.ExpressionNode {
	switch n := node.(type) {
	case *ast.FalseLiteralNode, *ast.TrueLiteralNode, *ast.NilLiteralNode,
		*ast.InterpolatedSymbolLiteralNode, *ast.ConstantDeclarationNode,
		*ast.InitDefinitionNode:
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
		c.checkStatements(n.Body)
		return n
	case *ast.ClassDeclarationNode:
		c.checkStatements(n.Body)
		return n
	case *ast.MixinDeclarationNode:
		c.checkStatements(n.Body)
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
	default:
		c.addError(
			fmt.Sprintf("invalid expression type %T", node),
			node.Span(),
		)
		return n
	}
}

func (c *Checker) includeMixin(node ast.ComplexConstantNode) {
	constantType, _ := c.resolveConstantType(node)

	constantMixin, constantIsMixin := constantType.(*types.Mixin)
	if !constantIsMixin {
		c.addError(
			"only mixins can be included",
			node.Span(),
		)
		return
	}

	switch c.mode {
	case classMode, mixinMode:
	default:
		c.addError(
			"cannot include mixins in this context",
			node.Span(),
		)
		return
	}

	headProxy, tailProxy := constantMixin.CreateProxy()
	target := c.currentConstScope().container

	switch t := target.(type) {
	case *types.Class:
		tailProxy.SetParent(t.Parent())
		t.SetParent(headProxy)
	case *types.Mixin:
		tailProxy.SetParent(t.Parent())
		t.SetParent(headProxy)
	default:
		c.addError(
			fmt.Sprintf(
				"cannot include `%s` in `%s`",
				types.Inspect(constantType),
				types.Inspect(t),
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
		c.addError(
			fmt.Sprintf("invalid constant type %T", node),
			node.Span(),
		)
		return node
	}
}

// Checks whether two methods are compatible.
func (c *Checker) checkMethodCompatibility(baseMethod, overrideMethod *types.Method, span *position.Span) bool {
	areCompatible := true
	if baseMethod != nil {
		if baseMethod.ReturnType != nil && overrideMethod.ReturnType != baseMethod.ReturnType {
			c.addError(
				fmt.Sprintf(
					"method `%s` has a different return type than `%s`, has `%s`, should have `%s`",
					types.Inspect(overrideMethod),
					types.Inspect(baseMethod),
					types.Inspect(overrideMethod.ReturnType),
					types.Inspect(baseMethod.ReturnType),
				),
				span,
			)
			areCompatible = false
		}
		if overrideMethod.ThrowType != baseMethod.ThrowType {
			c.addError(
				fmt.Sprintf(
					"method `%s` has a different throw type than `%s`, has `%s`, should have `%s`",
					types.Inspect(overrideMethod),
					types.Inspect(baseMethod),
					types.Inspect(overrideMethod.ThrowType),
					types.Inspect(baseMethod.ThrowType),
				),
				span,
			)
			areCompatible = false
		}

		if len(baseMethod.Params) > len(overrideMethod.Params) {
			c.addError(
				fmt.Sprintf(
					"method `%s` has less parameters than `%s`, has `%d`, should have `%d`",
					types.Inspect(overrideMethod),
					types.Inspect(baseMethod),
					len(overrideMethod.Params),
					len(baseMethod.Params),
				),
				span,
			)
			areCompatible = false
		} else {
			for i := range len(baseMethod.Params) {
				oldParam := baseMethod.Params[i]
				newParam := overrideMethod.Params[i]
				if oldParam.Name != newParam.Name {
					c.addError(
						fmt.Sprintf(
							"method `%s` has a different parameter name than `%s`, has `%s`, should have `%s`",
							types.Inspect(overrideMethod),
							types.Inspect(baseMethod),
							newParam.Name,
							oldParam.Name,
						),
						span,
					)
					areCompatible = false
					continue
				}
				if oldParam.Kind != newParam.Kind {
					c.addError(
						fmt.Sprintf(
							"method `%s` has a different parameter kind than `%s`, has `%s`, should have `%s`",
							types.Inspect(overrideMethod),
							types.Inspect(baseMethod),
							newParam.NameWithKind(),
							oldParam.NameWithKind(),
						),
						span,
					)
					areCompatible = false
					continue
				}
				if oldParam.Type != newParam.Type {
					c.addError(
						fmt.Sprintf(
							"method `%s` has a different type for parameter `%s` than `%s`, has `%s`, should have `%s`",
							types.Inspect(overrideMethod),
							newParam.Name,
							types.Inspect(baseMethod),
							types.Inspect(newParam.Type),
							types.Inspect(oldParam.Type),
						),
						span,
					)
					areCompatible = false
					continue
				}
			}

			for i := len(baseMethod.Params); i < len(overrideMethod.Params); i++ {
				param := overrideMethod.Params[i]
				if !param.IsOptional() {
					c.addError(
						fmt.Sprintf(
							"method `%s` has a required parameter missing in `%s`, got `%s`",
							types.Inspect(overrideMethod),
							types.Inspect(baseMethod),
							param.Name,
						),
						span,
					)
					areCompatible = false
				}
			}
		}

	}

	return areCompatible
}

func (c *Checker) getMethod(typ types.Type, name string, span *position.Span, reportErrors bool) *types.Method {
	return c._getMethod(typ, name, span, false, reportErrors)
}

func (c *Checker) getMethodInContainer(container types.ConstantContainer, typ types.Type, name string, span *position.Span, inParent, reportErrors bool) *types.Method {
	method := container.MethodString(name)
	if method == nil {
		parent := container.Parent()
		if parent != nil {
			method := c._getMethod(parent, name, span, true, reportErrors)
			if reportErrors && method == nil && !inParent {
				c.addMissingMethodError(typ, name, span)
			}
			return method
		}
		if reportErrors && !inParent {
			c.addMissingMethodError(typ, name, span)
		}
	}
	return method
}

func (c *Checker) _getMethod(typ types.Type, name string, span *position.Span, inParent, reportErrors bool) *types.Method {
	typ = typ.ToNonLiteral(c.GlobalEnv)

	switch t := typ.(type) {
	case *types.Class:
		return c.getMethodInContainer(t, typ, name, span, inParent, reportErrors)
	case *types.Module:
		return c.getMethodInContainer(t, typ, name, span, inParent, reportErrors)
	case *types.Mixin:
		return c.getMethodInContainer(t, typ, name, span, inParent, reportErrors)
	case *types.MixinProxy:
		return c.getMethodInContainer(t, typ, name, span, inParent, reportErrors)
	case *types.Nilable:
		nilType := c.GlobalEnv.StdSubtype(symbol.Nil).(*types.Class)
		nilMethod := nilType.MethodString(name)
		if reportErrors && nilMethod == nil {
			c.addMissingMethodError(nilType, name, span)
		}
		nonNilMethod := c.getMethod(t.Type, name, span, reportErrors)
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

		if c.checkMethodCompatibility(baseMethod, overrideMethod, span) {
			return baseMethod
		}
		return nil
	case *types.Union:
		var methods []*types.Method
		var baseMethod *types.Method

		for _, element := range t.Elements {
			elementMethod := c.getMethod(element, name, span, reportErrors)
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

			if !c.checkMethodCompatibility(baseMethod, method, span) {
				isCompatible = false
			}
		}

		if isCompatible {
			return baseMethod
		}

		return nil
	default:
		if reportErrors {
			c.addMissingMethodError(typ, name, span)
		}
		return nil
	}
}

func (c *Checker) addMissingMethodError(typ types.Type, name string, span *position.Span) {
	c.addError(
		fmt.Sprintf("method `%s` is not defined on type `%s`", name, types.Inspect(typ)),
		span,
	)
}

func (c *Checker) typeOf(node ast.Node) types.Type {
	return node.Type(c.GlobalEnv)
}

func (c *Checker) isNilable(typ types.Type) bool {
	return types.IsNilable(typ, c.GlobalEnv)
}

func (c *Checker) toNonNilable(typ types.Type) types.Type {
	return types.ToNonNilable(typ, c.GlobalEnv)
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
		if !c.isSubtype(posArgType, param.Type) {
			c.addError(
				fmt.Sprintf(
					"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
					types.Inspect(param.Type),
					param.Name,
					method.Name,
					types.Inspect(posArgType),
				),
				posArg.Span(),
			)
		}
	}

	if method.HasPositionalRestParam() {
		if len(positionalArguments) < requiredPosParamCount {
			c.addError(
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
			if !c.isSubtype(posArgType, posRestParam.Type) {
				c.addError(
					fmt.Sprintf(
						"expected type `%s` for rest parameter `*%s` in call to `%s`, got type `%s`",
						types.Inspect(posRestParam.Type),
						posRestParam.Name,
						method.Name,
						types.Inspect(posArgType),
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
			if !c.isSubtype(posArgType, param.Type) {
				c.addError(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.Inspect(param.Type),
						param.Name,
						method.Name,
						types.Inspect(posArgType),
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
				c.addError(
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
			if !c.isSubtype(namedArgType, param.Type) {
				c.addError(
					fmt.Sprintf(
						"expected type `%s` for parameter `%s` in call to `%s`, got type `%s`",
						types.Inspect(param.Type),
						param.Name,
						method.Name,
						types.Inspect(namedArgType),
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
			c.addError(
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
			if !c.isSubtype(namedArgType, namedRestParam.Type) {
				c.addError(
					fmt.Sprintf(
						"expected type `%s` for named rest parameter `**%s` in call to `%s`, got type `%s`",
						types.Inspect(namedRestParam.Type),
						namedRestParam.Name,
						method.Name,
						types.Inspect(namedArgType),
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
			c.addError(
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
	method := c.getMethod(c.selfType, node.MethodName, node.Span(), true)
	if method == nil {
		c.checkExpressions(node.PositionalArguments)
		node.SetType(types.Void{})
	}

	typedPositionalArguments := c.checkMethodArguments(method, node.PositionalArguments, node.NamedArguments, node.Span())
	node.PositionalArguments = typedPositionalArguments
	node.NamedArguments = nil
	node.SetType(method)
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
		c.addError(
			fmt.Sprintf("`%s` cannot be instantiated", className),
			node.Span(),
		)
		c.checkExpressions(node.PositionalArguments)
		node.SetType(types.Void{})
		return
	}
	method := c.getMethod(classType, "#init", node.Span(), false)
	if method == nil {
		method = types.NewMethod(
			"#init",
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

func (c *Checker) methodCall(node *ast.MethodCallNode) {
	receiver := c.checkExpression(node.Receiver)
	node.Receiver = receiver
	receiverType := c.typeOf(receiver)
	var method *types.Method
	if node.NilSafe {
		nonNilableReceiverType := c.toNonNilable(receiverType)
		method = c.getMethod(nonNilableReceiverType, node.MethodName, node.Span(), true)
	} else {
		method = c.getMethod(receiverType, node.MethodName, node.Span(), true)
	}
	if method == nil {
		c.checkExpressions(node.PositionalArguments)
		node.SetType(types.Void{})
		return
	}

	typedPositionalArguments := c.checkMethodArguments(method, node.PositionalArguments, node.NamedArguments, node.Span())
	returnType := method.ReturnType
	if node.NilSafe {
		if !c.isNilable(receiverType) {
			c.addError(
				fmt.Sprintf("cannot make a nil-safe call on type `%s` which is not nilable", types.Inspect(receiverType)),
				node.Span(),
			)
		} else {
			returnType = c.toNilable(returnType)
		}
	}

	node.PositionalArguments = typedPositionalArguments
	node.NamedArguments = nil
	node.SetType(returnType)
}

func (c *Checker) attributeAccess(node *ast.AttributeAccessNode) ast.ExpressionNode {
	receiver := c.checkExpression(node.Receiver)
	receiverType := c.typeOf(receiver)
	method := c.getMethod(receiverType, node.AttributeName, node.Span(), true)
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
	c.addError(
		fmt.Sprintf("expected %s arguments in call to `%s`, got %d", method.ExpectedParamCountString(), method.Name, got),
		span,
	)
}

func (c *Checker) checkMethod(
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	body []ast.StatementNode,
) (ast.TypeNode, ast.TypeNode) {
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
		c.addLocal(p.Name, local{typ: declaredType, initialised: true})
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

	previousMode := c.mode
	c.mode = topLevelMode
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
		c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
	case *ast.PrivateIdentifierNode:
		c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
	// case *ast.SubscriptExpressionNode:
	// case *ast.ConstantLookupNode:
	// case *ast.PublicConstantNode:
	// case *ast.PrivateConstantNode:
	// case *ast.InstanceVariableNode:
	// case *ast.AttributeAccessNode:
	default:
		c.Errors.Add(
			fmt.Sprintf("cannot assign to: %T", node.Left),
			c.newLocation(node.Span()),
		)
	}
}

func (c *Checker) localVariableAssignment(name string, operator *token.Token, right ast.ExpressionNode, span *position.Span) {
	switch operator.Type {
	case token.EQUAL_OP:
	// case token.COLON_EQUAL:
	// case token.OR_OR_EQUAL:
	// case token.AND_AND_EQUAL:
	// case token.QUESTION_QUESTION_EQUAL:
	// case token.PLUS_EQUAL:
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
	default:
		c.Errors.Add(
			fmt.Sprintf("assignment using this operator has not been implemented: %s", operator.Type.String()),
			c.newLocation(span),
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
		c.addError(
			fmt.Sprintf("invalid string content %T", node),
			node.Span(),
		)
	}
}

func (c *Checker) addErrorWithLocation(message string, loc *position.Location) {
	c.Errors.Add(
		message,
		loc,
	)
}

func (c *Checker) addError(message string, span *position.Span) {
	c.Errors.Add(
		message,
		c.newLocation(span),
	)
}

func (c *Checker) resolveDeclaredNamespace(constantExpression ast.ExpressionNode) types.ConstantContainer {
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
		constContainer := c.currentConstScope().container
		return constContainer.ConstantString(constant.Value)
	case *ast.PrivateConstantNode:
		constContainer := c.currentConstScope().container
		return constContainer.ConstantString(constant.Value)
	case *ast.ConstantLookupNode:
		_, typ, _ := c.resolveConstantLookupForDeclaration(constant)
		return typ
	default:
		panic(fmt.Sprintf("invalid constant node: %T", constantExpression))
	}
}

// Get the type of the constant with the given name
func (c *Checker) resolveConstantForDeclaration(constantExpression ast.ExpressionNode) (types.ConstantContainer, types.Type, string) {
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

func (c *Checker) resolveConstantLookupForDeclaration(node *ast.ConstantLookupNode) (types.ConstantContainer, types.Type, string) {
	return c._resolveConstantLookupForDeclaration(node, true)
}

func (c *Checker) _resolveConstantLookupForDeclaration(node *ast.ConstantLookupNode, firstCall bool) (types.ConstantContainer, types.Type, string) {
	var leftContainerType types.Type
	var leftContainerName string

	switch l := node.Left.(type) {
	case *ast.PublicConstantNode:
		constContainer := c.currentConstScope().container
		leftContainerType = constContainer.ConstantString(l.Value)
		leftContainerName = types.MakeFullConstantName(constContainer.Name(), l.Value)
		if leftContainerType == nil {
			placeholder := types.NewPlaceholderNamespace(leftContainerName)
			placeholder.Locations.Append(c.newLocation(l.Span()))
			leftContainerType = placeholder
			c.registerPlaceholderNamespace(placeholder)
			constContainer.DefineConstant(l.Value, leftContainerType)
		} else if placeholder, ok := leftContainerType.(*types.PlaceholderNamespace); ok {
			placeholder.Locations.Append(c.newLocation(l.Span()))
		}
	case *ast.PrivateConstantNode:
		constContainer := c.currentConstScope().container
		leftContainerType = constContainer.ConstantString(l.Value)
		leftContainerName = types.MakeFullConstantName(constContainer.Name(), l.Value)
		if leftContainerType == nil {
			placeholder := types.NewPlaceholderNamespace(leftContainerName)
			placeholder.Locations.Append(c.newLocation(l.Span()))
			leftContainerType = placeholder
			c.registerPlaceholderNamespace(placeholder)
			constContainer.DefineConstant(l.Value, leftContainerType)
		} else if placeholder, ok := leftContainerType.(*types.PlaceholderNamespace); ok {
			placeholder.Locations.Append(c.newLocation(l.Span()))
		}
	case nil:
		leftContainerType = c.GlobalEnv.Root
	case *ast.ConstantLookupNode:
		_, leftContainerType, leftContainerName = c._resolveConstantLookupForDeclaration(l, false)
	default:
		c.addError(
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
		c.addError(
			fmt.Sprintf("cannot read private constant `%s`", rightName),
			node.Span(),
		)
	default:
		c.addError(
			fmt.Sprintf("invalid constant node %T", node),
			node.Span(),
		)
		return nil, nil, ""
	}

	constantName := types.MakeFullConstantName(leftContainerName, rightName)
	if leftContainerType == nil {
		return nil, nil, constantName
	}
	var leftContainer types.ConstantContainer
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
		c.addError(
			fmt.Sprintf("cannot read constants from `%s`, it is not a constant container", leftContainerName),
			node.Span(),
		)
		return nil, nil, constantName
	}

	constant := leftContainer.ConstantString(rightName)
	if constant == nil && !firstCall {
		placeholder := types.NewPlaceholderNamespace(constantName)
		placeholder.Locations.Append(c.newLocation(node.Right.Span()))
		constant = placeholder
		c.registerPlaceholderNamespace(placeholder)
		leftContainer.DefineConstant(rightName, constant)
	} else if placeholder, ok := constant.(*types.PlaceholderNamespace); ok {
		placeholder.Locations.Append(c.newLocation(node.Right.Span()))
	}

	return leftContainer, constant, constantName
}

// Get the type of the constant with the given name
func (c *Checker) resolveSimpleConstantForSetter(name string) (types.ConstantContainer, types.Type, string) {
	constContainer := c.currentConstScope().container
	constant := constContainer.ConstantString(name)
	fullName := types.MakeFullConstantName(constContainer.Name(), name)
	if constant != nil {
		return constContainer, constant, fullName
	}
	return constContainer, nil, fullName
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

	c.addError(
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

	c.addError(
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

	c.addError(
		fmt.Sprintf("undefined type `%s`", name),
		span,
	)
	return nil, name
}

// Add the local with the given name to the current local environment
func (c *Checker) addLocal(name string, l local) {
	env := c.currentLocalEnv()
	env.locals[value.ToSymbol(name)] = l
}

// Get the local with the specified name from the current local environment
func (c *Checker) getLocal(name string) (local, bool) {
	env := c.currentLocalEnv()
	return env.getLocal(name)
}

// Get the instance variable with the specified name
func (c *Checker) getInstanceVariable(name string) (types.Type, types.ConstantContainer) {
	container := c.currentConstScope().container

	for container != nil {
		ivar := container.InstanceVariableString(name)
		if ivar != nil {
			return ivar, container
		}

		container = container.Parent()
	}

	return nil, nil
}

// Resolve the local with the given name from the current local environment or any parent environment
func (c *Checker) resolveLocal(name string, span *position.Span) (local, bool) {
	env := c.currentLocalEnv()
	local, ok := env.resolveLocal(name)
	if !ok {
		c.addError(
			fmt.Sprintf("undefined local `%s`", name),
			span,
		)
	}
	return local, ok
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
		c.addError(
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
		c.addError(
			fmt.Sprintf("cannot use private type `%s`", rightName),
			node.Span(),
		)
	default:
		c.addError(
			fmt.Sprintf("invalid type node %T", node),
			node.Span(),
		)
		return nil, ""
	}

	typeName := types.MakeFullConstantName(leftContainerName, rightName)
	if leftContainerType == nil {
		return nil, typeName
	}
	leftContainer, ok := leftContainerType.(types.ConstantContainer)
	if !ok {
		c.addError(
			fmt.Sprintf("cannot read subtypes from `%s`, it is not a type container", leftContainerName),
			node.Span(),
		)
		return nil, typeName
	}

	constant := leftContainer.SubtypeString(rightName)
	if constant == nil {
		c.addError(
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
		c.addError(
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
	default:
		c.addError(
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
		c.addError(
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
		c.addError(
			fmt.Sprintf("cannot read private constant `%s`", rightName),
			node.Span(),
		)
	default:
		c.addError(
			fmt.Sprintf("invalid constant node %T", node),
			node.Span(),
		)
		return nil, ""
	}

	constantName := types.MakeFullConstantName(leftContainerName, rightName)
	if leftContainerType == nil {
		return nil, constantName
	}

	var leftContainer types.ConstantContainer
	switch l := leftContainerType.(type) {
	case *types.Module:
		leftContainer = l
	case *types.SingletonClass:
		leftContainer = l.AttachedObject
	default:
		c.addError(
			fmt.Sprintf("cannot read constants from `%s`, it is not a constant container", leftContainerName),
			node.Span(),
		)
		return nil, constantName
	}

	constant := leftContainer.ConstantString(rightName)
	if constant == nil {
		c.addError(
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
	local, ok := c.resolveLocal(node.Value, node.Span())
	if ok && !local.initialised {
		c.addError(
			fmt.Sprintf("cannot access uninitialised local `%s`", node.Value),
			node.Span(),
		)
	}
	node.SetType(local.typ)
	return node
}

func (c *Checker) privateIdentifier(node *ast.PrivateIdentifierNode) *ast.PrivateIdentifierNode {
	local, ok := c.resolveLocal(node.Value, node.Span())
	if ok && !local.initialised {
		c.addError(
			fmt.Sprintf("cannot access uninitialised local `%s`", node.Value),
			node.Span(),
		)
	}
	node.SetType(local.typ)
	return node
}

func (c *Checker) instanceVariable(node *ast.InstanceVariableNode) {
	typ, _ := c.getInstanceVariable(node.Value)
	if typ == nil {
		c.addError(
			fmt.Sprintf("undefined instance variable `@%s`", node.Value),
			node.Span(),
		)
		node.SetType(types.Void{})
		return
	}

	node.SetType(typ)
}

func (c *Checker) instanceVariableDeclaration(node *ast.InstanceVariableDeclarationNode) {
	if ivar, container := c.getInstanceVariable(node.Name); ivar != nil {
		c.addError(
			fmt.Sprintf("cannot redeclare instance variable `@%s`, previous definition found in `%s`", node.Name, container.Name()),
			node.Span(),
		)
	}
	var declaredType types.Type

	if node.TypeNode == nil {
		c.addError(
			fmt.Sprintf("cannot declare instance variable `@%s` without a type", node.Name),
			node.Span(),
		)

		declaredType = types.Void{}
	} else {
		// without an initialiser but with a type
		declaredTypeNode := c.checkTypeNode(node.TypeNode)
		declaredType = c.typeOf(declaredTypeNode)
		node.TypeNode = declaredTypeNode
	}

	switch c.mode {
	case mixinMode, classMode, moduleMode:
	default:
		c.addError(
			fmt.Sprintf("cannot declare instance variable `@%s` in this context", node.Name),
			node.Span(),
		)
		return
	}

	c.declareInstanceVariable(node.Name, declaredType)
}

func (c *Checker) variableDeclaration(node *ast.VariableDeclarationNode) {
	if _, ok := c.getLocal(node.Name); ok {
		c.addError(
			fmt.Sprintf("cannot redeclare local `%s`", node.Name),
			node.Span(),
		)
	}
	if node.Initialiser == nil {
		if node.TypeNode == nil {
			c.addError(
				fmt.Sprintf("cannot declare a variable without a type `%s`", node.Name),
				node.Span(),
			)
			c.addLocal(node.Name, local{typ: types.Void{}})
			node.SetType(types.Void{})
			return
		}

		// without an initialiser but with a type
		declaredTypeNode := c.checkTypeNode(node.TypeNode)
		declaredType := c.typeOf(declaredTypeNode)
		c.addLocal(node.Name, local{typ: declaredType})
		node.TypeNode = declaredTypeNode
		node.SetType(types.Void{})
		return
	}

	// with an initialiser
	if node.TypeNode == nil {
		// without a type, inference
		init := c.checkExpression(node.Initialiser)
		actualType := c.typeOf(init).ToNonLiteral(c.GlobalEnv)
		c.addLocal(node.Name, local{typ: actualType, initialised: true})
		if types.IsVoid(actualType) {
			c.addError(
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
	c.addLocal(node.Name, local{typ: declaredType, initialised: true})
	c.checkCanAssign(actualType, declaredType, init.Span())

	node.TypeNode = declaredTypeNode
	node.Initialiser = init
	node.SetType(declaredType)
}

func (c *Checker) valueDeclaration(node *ast.ValueDeclarationNode) {
	if _, ok := c.getLocal(node.Name); ok {
		c.addError(
			fmt.Sprintf("cannot redeclare local `%s`", node.Name),
			node.Span(),
		)
	}
	if node.Initialiser == nil {
		if node.TypeNode == nil {
			c.addError(
				fmt.Sprintf("cannot declare a value without a type `%s`", node.Name),
				node.Span(),
			)
			c.addLocal(node.Name, local{typ: types.Void{}})
			node.SetType(types.Void{})
			return
		}

		// without an initialiser but with a type
		declaredTypeNode := c.checkTypeNode(node.TypeNode)
		declaredType := c.typeOf(declaredTypeNode)
		c.addLocal(node.Name, local{typ: declaredType, singleAssignment: true})
		node.TypeNode = declaredTypeNode
		node.SetType(types.Void{})
	}

	// with an initialiser
	if node.TypeNode == nil {
		// without a type, inference
		init := c.checkExpression(node.Initialiser)
		actualType := c.typeOf(init)
		c.addLocal(node.Name, local{typ: actualType, initialised: true, singleAssignment: true})
		if types.IsVoid(actualType) {
			c.addError(
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
	c.addLocal(node.Name, local{typ: declaredType, initialised: true, singleAssignment: true})
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

func (c *Checker) declareModule(constantContainer types.ConstantContainer, constantType types.Type, fullConstantName, constantName string, span *position.Span) *types.Module {
	if constantType != nil {
		ct, ok := constantType.(*types.SingletonClass)
		if ok {
			constantType = ct.AttachedObject
		}

		switch t := constantType.(type) {
		case *types.Module:
			return t
		case *types.PlaceholderNamespace:
			module := types.NewModuleWithDetails(
				t.Name(),
				t.Constants(),
				t.Subtypes(),
				t.Methods(),
			)
			t.Replacement = module
			constantContainer.DefineConstant(constantName, module)
			return module
		default:
			c.addError(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				span,
			)
			return types.NewModule(fullConstantName)
		}
	} else if constantContainer == nil {
		return types.NewModule(fullConstantName)
	} else {
		return constantContainer.DefineModule(constantName)
	}
}

func (c *Checker) declareInstanceVariable(name string, typ types.Type) {
	container := c.currentConstScope().container
	container.DefineInstanceVariable(name, typ)
}

func (c *Checker) declareClass(constantContainer types.ConstantContainer, constantType types.Type, fullConstantName, constantName string, span *position.Span) *types.Class {
	if constantType != nil {
		ct, ok := constantType.(*types.SingletonClass)
		if !ok {
			c.addError(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				span,
			)
			return types.NewClass(fullConstantName, nil)
		}
		constantType = ct.AttachedObject

		switch t := constantType.(type) {
		case *types.Class:
			return t
		case *types.PlaceholderNamespace:
			class := types.NewClassWithDetails(
				t.Name(),
				nil,
				t.Constants(),
				t.Subtypes(),
				t.Methods(),
			)
			t.Replacement = class
			constantContainer.DefineConstant(constantName, class)
			return class
		default:
			c.addError(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				span,
			)
			return types.NewClass(fullConstantName, nil)
		}
	} else if constantContainer == nil {
		return types.NewClass(fullConstantName, nil)
	} else {
		return constantContainer.DefineClass(constantName, nil)
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
			c.addErrorWithLocation(
				fmt.Sprintf("undefined namespace `%s`", placeholder.Name()),
				location,
			)
		}
	}

	c.hoistMethodDefinitions(statements)
	c.checkMethods()

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
	constantName := extractConstantName(node.Constant)
	node.Constant = ast.NewPublicConstantNode(node.Constant.Span(), fullConstantName)
	if constant != nil {
		c.addError(
			fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
			node.Span(),
		)
	}
	init := c.checkExpression(node.Initialiser)

	if !init.IsStatic() {
		c.addError(
			"values assigned to constants must be static, known at compile time",
			init.Span(),
		)
	}

	// with an initialiser
	if node.TypeNode == nil {
		// without a type, inference
		actualType := c.typeOf(init)
		if types.IsVoid(actualType) {
			c.addError(
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
	if !c.isSubtype(assignedType, targetType) {
		c.addError(
			fmt.Sprintf(
				"type `%s` cannot be assigned to type `%s`",
				types.Inspect(assignedType),
				types.Inspect(targetType),
			),
			span,
		)
	}
}

func (c *Checker) hoistTypeDefinitions(statements []ast.StatementNode) {
	for _, statement := range statements {
		stmt, ok := statement.(*ast.ExpressionStatementNode)
		if !ok {
			continue
		}

		expression := stmt.Expression
		if docNode, ok := stmt.Expression.(*ast.DocCommentNode); ok {
			expression = docNode.Expression
		}

		switch expr := expression.(type) {
		case *ast.ModuleDeclarationNode:
			container, constant, fullConstantName := c.resolveConstantForDeclaration(expr.Constant)
			constantName := extractConstantName(expr.Constant)
			module := c.declareModule(container, constant, fullConstantName, constantName, expr.Span())

			c.pushConstScope(makeLocalConstantScope(module))
			c.pushMethodScope(makeLocalMethodScope(module))

			c.hoistTypeDefinitions(expr.Body)

			c.popMethodScope()
			c.popConstScope()
		case *ast.ClassDeclarationNode:
			container, constant, fullConstantName := c.resolveConstantForDeclaration(expr.Constant)
			constantName := extractConstantName(expr.Constant)
			class := c.declareClass(container, constant, fullConstantName, constantName, expr.Span())

			c.pushConstScope(makeLocalConstantScope(class))
			c.pushMethodScope(makeLocalMethodScope(class))

			c.hoistTypeDefinitions(expr.Body)

			c.popMethodScope()
			c.popConstScope()
		case *ast.MixinDeclarationNode:
			container, constant, fullConstantName := c.resolveConstantForDeclaration(expr.Constant)
			constantName := extractConstantName(expr.Constant)
			mixin := c.declareMixin(container, constant, fullConstantName, constantName, expr.Span())
			c.pushConstScope(makeLocalConstantScope(mixin))
			c.pushMethodScope(makeLocalMethodScope(mixin))

			c.hoistTypeDefinitions(expr.Body)

			c.popMethodScope()
			c.popConstScope()
		case *ast.ConstantDeclarationNode:
			c.checkConstantDeclaration(expr)
		}
	}
}

func (c *Checker) hoistMethodDefinitions(statements []ast.StatementNode) {
	for _, statement := range statements {
		stmt, ok := statement.(*ast.ExpressionStatementNode)
		if !ok {
			continue
		}

		expression := stmt.Expression
		if docNode, ok := stmt.Expression.(*ast.DocCommentNode); ok {
			expression = docNode.Expression
		}

		switch expr := expression.(type) {
		case *ast.MethodDefinitionNode:
			method := c.declareMethod(
				expr.Name,
				expr.Parameters,
				expr.ReturnType,
				expr.ThrowType,
				expr.Span(),
			)
			c.registerMethodCheck(method, expr)
		case *ast.InitDefinitionNode:
			method := c.declareMethod(
				"#init",
				expr.Parameters,
				nil,
				expr.ThrowType,
				expr.Span(),
			)
			c.registerInitCheck(method, expr)
		case *ast.InstanceVariableDeclarationNode:
			c.instanceVariableDeclaration(expr)
		case *ast.IncludeExpressionNode:
			for _, constant := range expr.Constants {
				c.includeMixin(constant)
			}
		case *ast.ModuleDeclarationNode:
			module := c.resolveDeclaredNamespace(expr.Constant)

			c.pushConstScope(makeLocalConstantScope(module))
			c.pushMethodScope(makeLocalMethodScope(module))

			previousMode := c.mode
			c.mode = moduleMode
			c.hoistMethodDefinitions(expr.Body)
			c.setMode(previousMode)

			c.popMethodScope()
			c.popConstScope()
		case *ast.ClassDeclarationNode:
			classNamespace := c.resolveDeclaredNamespace(expr.Constant)

			if class, ok := classNamespace.(*types.Class); ok {
				var superclass *types.Class

				if expr.Superclass == nil {
					superclass = c.GlobalEnv.StdSubtypeClass(symbol.Object)
				} else {
					superclassType, _ := c.resolveConstantType(expr.Superclass)
					var ok bool
					superclass, ok = superclassType.(*types.Class)
					if !ok {
						c.addError(
							fmt.Sprintf("`%s` is not a class", types.Inspect(superclassType)),
							expr.Superclass.Span(),
						)
					}
				}

				parent := class.Parent()
				if parent == nil && superclass != nil {
					class.SetParent(superclass)
				} else if parent != nil && parent != superclass {
					var span *position.Span
					if expr.Superclass == nil {
						span = expr.Span()
					} else {
						span = expr.Superclass.Span()
					}

					c.addError(
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

			c.pushConstScope(makeLocalConstantScope(classNamespace))
			c.pushMethodScope(makeLocalMethodScope(classNamespace))

			previousMode := c.mode
			c.mode = classMode
			c.hoistMethodDefinitions(expr.Body)
			c.setMode(previousMode)

			c.popMethodScope()
			c.popConstScope()
		case *ast.MixinDeclarationNode:
			mixin := c.resolveDeclaredNamespace(expr.Constant)

			c.pushConstScope(makeLocalConstantScope(mixin))
			c.pushMethodScope(makeLocalMethodScope(mixin))

			previousMode := c.mode
			c.mode = mixinMode
			c.hoistMethodDefinitions(expr.Body)
			c.setMode(previousMode)

			c.popMethodScope()
			c.popConstScope()
		}
	}
}

func (c *Checker) checkMethodDefinition(node *ast.MethodDefinitionNode) {
	returnType, throwType := c.checkMethod(
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Body,
	)
	node.ReturnType = returnType
	node.ThrowType = throwType
}

func (c *Checker) declareMethod(
	name string,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	span *position.Span,
) *types.Method {
	methodScope := c.currentMethodScope()
	oldMethod := c.getMethod(methodScope.container, name, nil, false)

	var params []*types.Parameter
	for _, param := range paramNodes {
		p, ok := param.(*ast.MethodParameterNode)
		if !ok {
			c.addError(
				fmt.Sprintf("invalid param type %T", param),
				param.Span(),
			)
		}
		var declaredType types.Type
		var declaredTypeNode ast.TypeNode
		if p.TypeNode == nil {
			c.addError(
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
		name,
		params,
		returnType,
		throwType,
		methodScope.container,
	)

	methodScope.container.SetMethod(name, newMethod)

	if oldMethod != nil {
		if typedReturnTypeNode != nil && oldMethod.ReturnType != nil && newMethod.ReturnType != oldMethod.ReturnType {
			c.addError(
				fmt.Sprintf(
					"cannot redeclare method `%s` with a different return type, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: %s",
					name,
					types.Inspect(newMethod.ReturnType),
					types.Inspect(oldMethod.ReturnType),
					oldMethod.DefinedUnder.Name(),
					oldMethod.InspectSignature(),
				),
				typedReturnTypeNode.Span(),
			)
		}
		if newMethod.ThrowType != oldMethod.ThrowType {
			var throwSpan *position.Span
			if typedThrowTypeNode != nil {
				throwSpan = typedThrowTypeNode.Span()
			} else {
				throwSpan = span
			}
			c.addError(
				fmt.Sprintf(
					"cannot redeclare method `%s` with a different throw type, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: %s",
					name,
					types.Inspect(newMethod.ThrowType),
					types.Inspect(oldMethod.ThrowType),
					oldMethod.DefinedUnder.Name(),
					oldMethod.InspectSignature(),
				),
				throwSpan,
			)
		}

		if len(oldMethod.Params) > len(newMethod.Params) {
			paramSpan := position.JoinSpanOfCollection(paramNodes)
			if paramSpan == nil {
				paramSpan = span
			}
			c.addError(
				fmt.Sprintf(
					"cannot redeclare method `%s` with less parameters\n  previous definition found in `%s`, with signature: %s",
					name,
					oldMethod.DefinedUnder.Name(),
					oldMethod.InspectSignature(),
				),
				paramSpan,
			)
		} else {
			for i := range len(oldMethod.Params) {
				oldParam := oldMethod.Params[i]
				newParam := newMethod.Params[i]
				if oldParam.Name != newParam.Name {
					c.addError(
						fmt.Sprintf(
							"cannot redeclare method `%s` with invalid parameter name, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: %s",
							name,
							newParam.Name,
							oldParam.Name,
							oldMethod.DefinedUnder.Name(),
							oldMethod.InspectSignature(),
						),
						paramNodes[i].Span(),
					)
					continue
				}
				if oldParam.Kind != newParam.Kind {
					c.addError(
						fmt.Sprintf(
							"cannot redeclare method `%s` with invalid parameter kind, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: %s",
							name,
							newParam.NameWithKind(),
							oldParam.NameWithKind(),
							oldMethod.DefinedUnder.Name(),
							oldMethod.InspectSignature(),
						),
						paramNodes[i].Span(),
					)
					continue
				}
				if oldParam.Type != newParam.Type {
					c.addError(
						fmt.Sprintf(
							"cannot redeclare method `%s` with invalid parameter type, is `%s`, should be `%s`\n  previous definition found in `%s`, with signature: %s",
							name,
							types.Inspect(newParam.Type),
							types.Inspect(oldParam.Type),
							oldMethod.DefinedUnder.Name(),
							oldMethod.InspectSignature(),
						),
						paramNodes[i].Span(),
					)
					continue
				}
			}

			for i := len(oldMethod.Params); i < len(newMethod.Params); i++ {
				param := newMethod.Params[i]
				if !param.IsOptional() {
					c.addError(
						fmt.Sprintf(
							"cannot redeclare method `%s` with additional parameter `%s`\n  previous definition found in `%s`, with signature: %s",
							name,
							param.Name,
							oldMethod.DefinedUnder.Name(),
							oldMethod.InspectSignature(),
						),
						paramNodes[i].Span(),
					)
				}
			}
		}

	}

	return newMethod
}

func (c *Checker) declareMixin(constantContainer types.ConstantContainer, constantType types.Type, fullConstantName, constantName string, span *position.Span) *types.Mixin {
	if constantType != nil {
		ct, ok := constantType.(*types.SingletonClass)
		if !ok {
			c.addError(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				span,
			)
			return types.NewMixin(fullConstantName)
		}
		constantType = ct.AttachedObject

		switch t := constantType.(type) {
		case *types.Mixin:
			return t
		case *types.PlaceholderNamespace:
			mixin := types.NewMixinWithDetails(
				t.Name(),
				nil,
				t.Constants(),
				t.Subtypes(),
				t.Methods(),
			)
			t.Replacement = mixin
			constantContainer.DefineConstant(constantName, mixin)
			return mixin
		default:
			c.addError(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				span,
			)
			return types.NewMixin(fullConstantName)
		}
	} else if constantContainer == nil {
		return types.NewMixin(fullConstantName)
	} else {
		return constantContainer.DefineMixin(constantName)
	}
}

// func (c *Checker) mixin(node *ast.MixinDeclarationNode) *typed.MixinDeclarationNode {
// 	var typedConstantNode typed.ExpressionNode
// 	var mixin *types.Mixin

// 	switch constant := node.Constant.(type) {
// 	case *ast.PublicConstantNode:
// 		constScope := c.currentConstScope()
// 		constantType, fullConstantName := c.resolveConstantForSetter(constant.Value)
// 		mixin = c.declareMixin(constScope.container, constantType, fullConstantName, constant.Value, node.Span())
// 		typedConstantNode = typed.NewPublicConstantNode(
// 			constant.Span(),
// 			fullConstantName,
// 			mixin,
// 		)
// 	case *ast.PrivateConstantNode:
// 		constScope := c.currentConstScope()
// 		constantType, fullConstantName := c.resolveConstantForSetter(constant.Value)
// 		mixin = c.declareMixin(constScope.container, constantType, fullConstantName, constant.Value, node.Span())
// 		typedConstantNode = typed.NewPublicConstantNode(
// 			constant.Span(),
// 			fullConstantName,
// 			mixin,
// 		)
// 	case *ast.ConstantLookupNode:
// 		constantContainer, constantType, fullConstantName := c.resolveConstantLookupForSetter(constant)
// 		constantName := extractConstantNameFromLookup(constant)
// 		mixin = c.declareMixin(constantContainer, constantType, fullConstantName, constantName, node.Span())
// 		typedConstantNode = typed.NewPublicConstantNode(
// 			constant.Span(),
// 			fullConstantName,
// 			mixin,
// 		)
// 	case nil:
// 		mixin = types.NewMixin("")
// 	default:
// 		c.addError(
// 			fmt.Sprintf("invalid mixin name node %T", node.Constant),
// 			node.Constant.Span(),
// 		)
// 	}

// 	c.pushConstScope(makeLocalConstantScope(mixin))
// 	c.pushMethodScope(makeLocalMethodScope(mixin))
// 	prevSelfType := c.selfType
// 	c.selfType = types.NewSingletonClass(mixin)

// 	previousMode := c.mode
// 	c.mode = mixinMode
// 	defer c.setMode(previousMode)

// 	newBody := c.hoistStatements(node.Body)
// 	c.selfType = prevSelfType
// 	c.popConstScope()
// 	c.popMethodScope()

// 	return typed.NewMixinDeclarationNode(
// 		node.Span(),
// 		typedConstantNode,
// 		nil,
// 		newBody,
// 		mixin,
// 	)
// }
