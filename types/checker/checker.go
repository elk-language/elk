// Package checker implements the Elk type checker
package checker

import (
	"fmt"
	"sync"

	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	typed "github.com/elk-language/elk/types/ast" // typed AST
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
)

// Check the types of Elk source code.
func CheckSource(sourceName string, source string, globalEnv *types.GlobalEnvironment, headerMode bool) (typed.Node, error.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	return CheckAST(sourceName, ast, globalEnv, headerMode)
}

// Check the types of an Elk AST.
func CheckAST(sourceName string, ast *ast.ProgramNode, globalEnv *types.GlobalEnvironment, headerMode bool) (typed.Node, error.ErrorList) {
	checker := newChecker(position.NewLocationWithSpan(sourceName, ast.Span()), globalEnv, headerMode)
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

// Holds the state of the type checking process
type Checker struct {
	Location       *position.Location
	Errors         *error.SyncErrorList
	GlobalEnv      *types.GlobalEnvironment
	HeaderMode     bool
	constantScopes []constantScope
	methodScopes   []methodScope
	localEnvs      []*localEnvironment
	returnType     types.Type
	throwType      types.Type
	selfType       types.Type
	mode           mode
}

// Instantiate a new Checker instance.
func newChecker(loc *position.Location, globalEnv *types.GlobalEnvironment, headerMode bool) *Checker {
	if globalEnv == nil {
		globalEnv = types.NewGlobalEnvironment()
	}
	return &Checker{
		Location:   loc,
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
	}
}

func (c *Checker) CheckSource(sourceName string, source string) (typed.Node, error.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	loc := position.NewLocationWithSpan(sourceName, ast.Span())
	c.Location = loc
	typedAst := c.checkProgram(ast)
	return typedAst, c.Errors.ErrorList
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
}

func (c *Checker) pushConstScope(constScope constantScope) {
	c.constantScopes = append(c.constantScopes, constScope)
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

// Create a new location struct with the given position.
func (c *Checker) newLocation(span *position.Span) *position.Location {
	return position.NewLocationWithSpan(c.Location.Filename, span)
}

func (c *Checker) checkProgram(node *ast.ProgramNode) *typed.ProgramNode {
	newStatements := c.hoistStatements(node.Body)
	return typed.NewProgramNode(
		node.Span(),
		newStatements,
	)
}

func (c *Checker) checkStatements(stmts []ast.StatementNode) []typed.StatementNode {
	var newStatements []typed.StatementNode
	for _, statement := range stmts {
		switch statement.(type) {
		case *ast.EmptyStatementNode:
			continue
		}
		newStatements = append(newStatements, c.checkStatement(statement))
	}

	return newStatements
}

func (c *Checker) checkStatement(node ast.Node) typed.StatementNode {
	switch node := node.(type) {
	case *ast.ExpressionStatementNode:
		expr := c.checkExpression(node.Expression)
		return typed.NewExpressionStatementNode(
			node.Span(),
			expr,
		)
	default:
		c.addError(
			fmt.Sprintf("incorrect statement type %#v", node),
			node.Span(),
		)
		return typed.NewInvalidNode(node.Span(), nil)
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

func (c *Checker) checkExpressions(exprs []ast.ExpressionNode) []typed.ExpressionNode {
	var newExpressions []typed.ExpressionNode
	for _, expr := range exprs {
		newExpressions = append(newExpressions, c.checkExpression(expr))
	}

	return newExpressions
}

func (c *Checker) checkExpression(node ast.ExpressionNode) typed.ExpressionNode {
	switch n := node.(type) {
	case *ast.FalseLiteralNode:
		return typed.NewFalseLiteralNode(n.Span())
	case *ast.TrueLiteralNode:
		return typed.NewTrueLiteralNode(n.Span())
	case *ast.NilLiteralNode:
		return typed.NewNilLiteralNode(n.Span())
	case *ast.TypeExpressionNode:
		typeNode := c.checkTypeNode(n.TypeNode)
		return typed.NewTypeExpressionNode(n.Span(), typeNode)
	case *ast.IntLiteralNode:
		return typed.NewIntLiteralNode(
			n.Span(),
			n.Value,
			types.NewIntLiteral(n.Value),
		)
	case *ast.Int64LiteralNode:
		return typed.NewInt64LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt64Literal(n.Value),
		)
	case *ast.Int32LiteralNode:
		return typed.NewInt32LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt32Literal(n.Value),
		)
	case *ast.Int16LiteralNode:
		return typed.NewInt16LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt16Literal(n.Value),
		)
	case *ast.Int8LiteralNode:
		return typed.NewInt8LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt8Literal(n.Value),
		)
	case *ast.UInt64LiteralNode:
		return typed.NewUInt64LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt64Literal(n.Value),
		)
	case *ast.UInt32LiteralNode:
		return typed.NewUInt32LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt32Literal(n.Value),
		)
	case *ast.UInt16LiteralNode:
		return typed.NewUInt16LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt16Literal(n.Value),
		)
	case *ast.UInt8LiteralNode:
		return typed.NewUInt8LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt8Literal(n.Value),
		)
	case *ast.FloatLiteralNode:
		return typed.NewFloatLiteralNode(
			n.Span(),
			n.Value,
			types.NewFloatLiteral(n.Value),
		)
	case *ast.Float64LiteralNode:
		return typed.NewFloat64LiteralNode(
			n.Span(),
			n.Value,
			types.NewFloat64Literal(n.Value),
		)
	case *ast.Float32LiteralNode:
		return typed.NewFloat32LiteralNode(
			n.Span(),
			n.Value,
			types.NewFloat32Literal(n.Value),
		)
	case *ast.BigFloatLiteralNode:
		return typed.NewBigFloatLiteralNode(
			n.Span(),
			n.Value,
			types.NewBigFloatLiteral(n.Value),
		)
	case *ast.DoubleQuotedStringLiteralNode:
		return typed.NewDoubleQuotedStringLiteralNode(
			n.Span(),
			n.Value,
			types.NewStringLiteral(n.Value),
		)
	case *ast.RawStringLiteralNode:
		return typed.NewRawStringLiteralNode(
			n.Span(),
			n.Value,
			types.NewStringLiteral(n.Value),
		)
	case *ast.RawCharLiteralNode:
		return typed.NewRawCharLiteralNode(
			n.Span(),
			n.Value,
			types.NewCharLiteral(n.Value),
		)
	case *ast.CharLiteralNode:
		return typed.NewCharLiteralNode(
			n.Span(),
			n.Value,
			types.NewCharLiteral(n.Value),
		)
	case *ast.InterpolatedStringLiteralNode:
		return c.interpolatedStringLiteral(n)
	case *ast.SimpleSymbolLiteralNode:
		return typed.NewSimpleSymbolLiteralNode(
			n.Span(),
			n.Content,
			types.NewSymbolLiteral(n.Content),
		)
	case *ast.InterpolatedSymbolLiteralNode:
		return typed.NewInterpolatedSymbolLiteralNode(
			n.Span(),
			c.interpolatedStringLiteral(n.Content),
		)
	case *ast.VariableDeclarationNode:
		return c.variableDeclaration(n)
	case *ast.ValueDeclarationNode:
		return c.valueDeclaration(n)
	case *ast.ConstantDeclarationNode:
		return c.constantDeclaration(n)
	case *ast.PublicIdentifierNode:
		return c.publicIdentifier(n)
	case *ast.PrivateIdentifierNode:
		return c.privateIdentifier(n)
	case *ast.PublicConstantNode:
		return c.publicConstant(n)
	case *ast.PrivateConstantNode:
		return c.privateConstant(n)
	case *ast.ConstantLookupNode:
		return c.constantLookup(n)
	case *ast.ModuleDeclarationNode:
		return c.module(n)
	case *ast.ClassDeclarationNode:
		return c.class(n)
	case *ast.MixinDeclarationNode:
		return c.mixin(n)
	case *ast.InitDefinitionNode:
		return c.initDefinition(n)
	case *ast.AssignmentExpressionNode:
		return c.assignmentExpression(n)
	case *ast.ReceiverlessMethodCallNode:
		return c.receiverlessMethodCall(n)
	case *ast.MethodCallNode:
		return c.methodCall(n)
	case *ast.ConstructorCallNode:
		return c.constructorCall(n)
	case *ast.AttributeAccessNode:
		return c.attributeAccess(n)
	case *ast.IncludeExpressionNode:
		return c.include(n)
	default:
		c.addError(
			fmt.Sprintf("invalid expression type %T", node),
			node.Span(),
		)
		return typed.NewInvalidNode(node.Span(), nil)
	}
}

func (c *Checker) include(node *ast.IncludeExpressionNode) *typed.IncludeExpressionNode {
	var constants []typed.ComplexConstantNode
	for _, constant := range node.Constants {
		constants = append(constants, c.includeMixin(constant))
	}

	return typed.NewIncludeExpressionNode(
		node.Span(),
		constants,
	)
}

func (c *Checker) includeMixin(node ast.ComplexConstantNode) typed.ComplexConstantNode {
	constantNode := c.checkComplexConstantType(node)
	constantType := c.typeOf(constantNode)

	constantMixin, constantIsMixin := constantType.(*types.Mixin)
	if !constantIsMixin {
		c.addError(
			"only mixins can be included",
			node.Span(),
		)

		return constantNode
	}

	switch c.mode {
	case classMode, mixinMode:
	default:
		c.addError(
			"cannot include mixins in this context",
			node.Span(),
		)

		return constantNode
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

	return constantNode
}

func (c *Checker) checkComplexConstant(node ast.ComplexConstantNode) typed.ComplexConstantNode {
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
		return typed.NewInvalidNode(node.Span(), nil)
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

func (c *Checker) typeOf(node typed.Node) types.Type {
	return typed.TypeOf(node, c.GlobalEnv)
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

func (c *Checker) checkMethodArguments(method *types.Method, positionalArguments []ast.ExpressionNode, namedArguments []ast.NamedArgumentNode, span *position.Span) []typed.ExpressionNode {
	reqParamCount := method.RequiredParamCount()
	requiredPosParamCount := len(method.Params) - method.OptionalParamCount
	if method.PostParamCount != -1 {
		requiredPosParamCount -= method.PostParamCount + 1
	}
	positionalRestParamIndex := method.PositionalRestParamIndex()
	var typedPositionalArguments []typed.ExpressionNode

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
		restPositionalArguments := typed.NewArrayTupleLiteralNode(
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
				typed.NewUndefinedLiteralNode(span),
			)
		}
	}

	if method.HasNamedRestParam {
		namedRestArgs := typed.NewHashRecordLiteralNode(
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
				typed.NewSymbolKeyValueExpressionNode(
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

func (c *Checker) receiverlessMethodCall(node *ast.ReceiverlessMethodCallNode) *typed.ReceiverlessMethodCallNode {
	method := c.getMethod(c.selfType, node.MethodName, node.Span(), true)
	if method == nil {
		return typed.NewReceiverlessMethodCallNode(
			node.Span(),
			node.MethodName,
			c.checkExpressions(node.PositionalArguments),
			types.Void{},
		)
	}

	typedPositionalArguments := c.checkMethodArguments(method, node.PositionalArguments, node.NamedArguments, node.Span())

	return typed.NewReceiverlessMethodCallNode(
		node.Span(),
		node.MethodName,
		typedPositionalArguments,
		method.ReturnType,
	)
}

func (c *Checker) constructorCall(node *ast.ConstructorCallNode) *typed.ConstructorCallNode {
	classNode := c.checkComplexConstantType(node.Class)
	classType := c.typeOf(classNode)
	var className string

	switch cn := classNode.(type) {
	case *typed.PublicConstantNode:
		className = cn.Value
	case *typed.PrivateConstantNode:
		className = cn.Value
	}

	class, isClass := classType.(*types.Class)
	if !isClass {
		c.addError(
			fmt.Sprintf("`%s` cannot be instantiated", className),
			node.Span(),
		)
		return typed.NewConstructorCallNode(
			node.Span(),
			classNode,
			c.checkExpressions(node.PositionalArguments),
			types.Void{},
		)
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

	return typed.NewConstructorCallNode(
		node.Span(),
		classNode,
		typedPositionalArguments,
		class,
	)
}

func (c *Checker) methodCall(node *ast.MethodCallNode) *typed.MethodCallNode {
	receiver := c.checkExpression(node.Receiver)
	receiverType := c.typeOf(receiver)
	var method *types.Method
	if node.NilSafe {
		nonNilableReceiverType := c.toNonNilable(receiverType)
		method = c.getMethod(nonNilableReceiverType, node.MethodName, node.Span(), true)
	} else {
		method = c.getMethod(receiverType, node.MethodName, node.Span(), true)
	}
	if method == nil {
		return typed.NewMethodCallNode(
			node.Span(),
			receiver,
			node.NilSafe,
			node.MethodName,
			c.checkExpressions(node.PositionalArguments),
			types.Void{},
		)
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

	return typed.NewMethodCallNode(
		node.Span(),
		receiver,
		node.NilSafe,
		node.MethodName,
		typedPositionalArguments,
		returnType,
	)
}

func (c *Checker) attributeAccess(node *ast.AttributeAccessNode) typed.ExpressionNode {
	receiver := c.checkExpression(node.Receiver)
	receiverType := c.typeOf(receiver)
	method := c.getMethod(receiverType, node.AttributeName, node.Span(), true)
	if method == nil {
		return typed.NewAttributeAccessNode(
			node.Span(),
			receiver,
			node.AttributeName,
			types.Void{},
		)
	}

	typedPositionalArguments := c.checkMethodArguments(method, nil, nil, node.Span())

	return typed.NewMethodCallNode(
		node.Span(),
		receiver,
		false,
		node.AttributeName,
		typedPositionalArguments,
		method.ReturnType,
	)
}

func (c *Checker) addWrongArgumentCountError(got int, method *types.Method, span *position.Span) {
	c.addError(
		fmt.Sprintf("expected %s arguments in call to `%s`, got %d", method.ExpectedParamCountString(), method.Name, got),
		span,
	)
}

func (c *Checker) checkMethod(
	oldMethod *types.Method,
	newMethod *types.Method,
	name string,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	body []ast.StatementNode,
	span *position.Span,
) (
	[]typed.ParameterNode,
	typed.TypeNode,
	typed.TypeNode,
	[]typed.StatementNode,
) {
	methodScope := c.currentMethodScope()
	env := newLocalEnvironment(nil)
	c.pushLocalEnv(env)
	defer c.popLocalEnv()

	var typedParamNodes []typed.ParameterNode
	for _, param := range paramNodes {
		p, _ := param.(*ast.MethodParameterNode)
		var declaredType types.Type
		var declaredTypeNode typed.TypeNode
		if p.Type != nil {
			declaredTypeNode = c.checkTypeNode(p.Type)
			declaredType = c.typeOf(declaredTypeNode)
		}
		var initNode typed.ExpressionNode
		if p.Initialiser != nil {
			initNode = c.checkExpression(p.Initialiser)
			initType := c.typeOf(initNode)
			if !c.isSubtype(initType, declaredType) {
				c.addError(
					fmt.Sprintf(
						"type `%s` cannot be assigned to type `%s`",
						types.Inspect(initType),
						types.Inspect(declaredType),
					),
					initNode.Span(),
				)
			}
		}
		c.addLocal(p.Name, local{typ: declaredType, initialised: true})
		typedParamNodes = append(typedParamNodes, typed.NewMethodParameterNode(
			p.Span(),
			p.Name,
			p.SetInstanceVariable,
			declaredTypeNode,
			initNode,
			typed.ParameterKind(p.Kind),
		))
	}

	var returnType types.Type
	var typedReturnTypeNode typed.TypeNode
	if returnTypeNode != nil {
		typedReturnTypeNode = c.checkTypeNode(returnTypeNode)
		returnType = c.typeOf(typedReturnTypeNode)
	} else {
		returnType = types.Void{}
	}

	var throwType types.Type
	var typedThrowTypeNode typed.TypeNode
	if throwTypeNode != nil {
		typedThrowTypeNode = c.checkTypeNode(throwTypeNode)
		throwType = c.typeOf(typedThrowTypeNode)
	}

	setMethod := true
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
			setMethod = false
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
			setMethod = false
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
			setMethod = false
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
						typedParamNodes[i].Span(),
					)
					setMethod = false
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
						typedParamNodes[i].Span(),
					)
					setMethod = false
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
						typedParamNodes[i].Span(),
					)
					setMethod = false
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
						typedParamNodes[i].Span(),
					)
					setMethod = false
				}
			}
		}

	}

	if setMethod {
		methodScope.container.SetMethod(name, newMethod)
	}

	previousMode := c.mode
	c.mode = topLevelMode
	defer c.setMode(previousMode)
	c.returnType = returnType
	c.throwType = throwType
	typedBody := c.checkStatements(body)
	c.returnType = nil
	c.throwType = nil

	return typedParamNodes,
		typedReturnTypeNode,
		typedThrowTypeNode,
		typedBody
}

func (c *Checker) defineMethod(
	name string,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
	body []ast.StatementNode,
	span *position.Span,
) (
	[]typed.ParameterNode,
	typed.TypeNode,
	typed.TypeNode,
	[]typed.StatementNode,
) {
	methodScope := c.currentMethodScope()
	env := newLocalEnvironment(nil)
	c.pushLocalEnv(env)
	defer c.popLocalEnv()

	oldMethod := c.getMethod(methodScope.container, name, span, false)

	var typedParamNodes []typed.ParameterNode
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
		var declaredTypeNode typed.TypeNode
		if p.Type == nil {
			c.addError(
				fmt.Sprintf("cannot declare parameter `%s` without a type", p.Name),
				param.Span(),
			)
		} else {
			declaredTypeNode = c.checkTypeNode(p.Type)
			declaredType = c.typeOf(declaredTypeNode)
		}
		var initNode typed.ExpressionNode
		if p.Initialiser != nil {
			initNode = c.checkExpression(p.Initialiser)
			initType := c.typeOf(initNode)
			if !c.isSubtype(initType, declaredType) {
				c.addError(
					fmt.Sprintf(
						"type `%s` cannot be assigned to type `%s`",
						types.Inspect(initType),
						types.Inspect(declaredType),
					),
					initNode.Span(),
				)
			}
		}
		c.addLocal(p.Name, local{typ: declaredType, initialised: true})
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
		typedParamNodes = append(typedParamNodes, typed.NewMethodParameterNode(
			p.Span(),
			p.Name,
			p.SetInstanceVariable,
			declaredTypeNode,
			initNode,
			typed.ParameterKind(p.Kind),
		))
	}

	var returnType types.Type
	var typedReturnTypeNode typed.TypeNode
	if returnTypeNode != nil {
		typedReturnTypeNode = c.checkTypeNode(returnTypeNode)
		returnType = c.typeOf(typedReturnTypeNode)
	} else {
		returnType = types.Void{}
	}

	var throwType types.Type
	var typedThrowTypeNode typed.TypeNode
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

	setMethod := true
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
			setMethod = false
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
			setMethod = false
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
			setMethod = false
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
						typedParamNodes[i].Span(),
					)
					setMethod = false
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
						typedParamNodes[i].Span(),
					)
					setMethod = false
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
						typedParamNodes[i].Span(),
					)
					setMethod = false
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
						typedParamNodes[i].Span(),
					)
					setMethod = false
				}
			}
		}

	}

	if setMethod {
		methodScope.container.SetMethod(name, newMethod)
	}

	previousMode := c.mode
	c.mode = topLevelMode
	defer c.setMode(previousMode)
	c.returnType = returnType
	c.throwType = throwType
	typedBody := c.checkStatements(body)
	c.returnType = nil
	c.throwType = nil

	return typedParamNodes,
		typedReturnTypeNode,
		typedThrowTypeNode,
		typedBody
}

func (c *Checker) initDefinition(node *ast.InitDefinitionNode) typed.ExpressionNode {
	constScope := c.currentConstScope()

	switch constScope.container.(type) {
	case *types.Class:
	default:
		c.addError(
			"init definitions cannot appear outside of classes",
			node.Span(),
		)
		return typed.NewInvalidNode(node.Span(), nil)
	}

	paramNodes, _, throwTypeNode, body := c.defineMethod(
		"#init",
		node.Parameters,
		nil,
		node.ThrowType,
		node.Body,
		node.Span(),
	)
	return typed.NewInitDefinitionNode(
		node.Span(),
		paramNodes,
		throwTypeNode,
		body,
	)
}

func (c *Checker) assignmentExpression(node *ast.AssignmentExpressionNode) *typed.AssignmentExpressionNode {
	switch n := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		return c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
	case *ast.PrivateIdentifierNode:
		return c.localVariableAssignment(n.Value, node.Op, node.Right, node.Span())
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

	return nil
}

func (c *Checker) localVariableAssignment(name string, operator *token.Token, right ast.ExpressionNode, span *position.Span) *typed.AssignmentExpressionNode {
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
	return nil
}

func (c *Checker) interpolatedStringLiteral(node *ast.InterpolatedStringLiteralNode) *typed.InterpolatedStringLiteralNode {
	var newContent []typed.StringLiteralContentNode
	for _, contentSection := range node.Content {
		newContent = append(newContent, c.checkStringContent(contentSection))
	}
	return typed.NewInterpolatedStringLiteralNode(
		node.Span(),
		newContent,
	)
}

func (c *Checker) checkStringContent(node ast.StringLiteralContentNode) typed.StringLiteralContentNode {
	switch n := node.(type) {
	case *ast.StringInspectInterpolationNode:
		return typed.NewStringInspectInterpolationNode(
			n.Span(),
			c.checkExpression(n.Expression),
		)
	case *ast.StringInterpolationNode:
		return typed.NewStringInterpolationNode(
			n.Span(),
			c.checkExpression(n.Expression),
		)
	case *ast.StringLiteralContentSectionNode:
		return typed.NewStringLiteralContentSectionNode(
			n.Span(),
			n.Value,
		)
	default:
		c.addError(
			fmt.Sprintf("invalid string content %T", node),
			node.Span(),
		)
		return typed.NewInvalidNode(node.Span(), nil)
	}
}

func (c *Checker) addError(message string, span *position.Span) {
	c.Errors.Add(
		message,
		c.newLocation(span),
	)
}

// Get the type of the constant with the given name
func (c *Checker) resolveConstantForSetter(name string) (types.Type, string) {
	constScope := c.currentConstScope()
	constant := constScope.container.ConstantString(name)
	fullName := types.MakeFullConstantName(constScope.container.Name(), name)
	if constant != nil {
		return constant, fullName
	}
	return nil, fullName
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

func (c *Checker) checkComplexConstantType(node ast.ComplexConstantNode) typed.ComplexConstantNode {
	switch n := node.(type) {
	case *ast.PublicConstantNode:
		return c.checkPublicConstantType(n)
	case *ast.PrivateConstantNode:
		return c.checkPrivateConstantType(n)
	case *ast.ConstantLookupNode:
		return c.constantLookupType(n)
	default:
		c.addError(
			fmt.Sprintf("invalid constant type node %T", node),
			node.Span(),
		)
		return typed.NewInvalidNode(node.Span(), nil)
	}
}

func (c *Checker) checkPublicConstantType(node *ast.PublicConstantNode) *typed.PublicConstantNode {
	typ, _ := c.resolveType(node.Value, node.Span())
	if typ == nil {
		typ = types.Void{}
	}
	return typed.NewPublicConstantNode(
		node.Span(),
		node.Value,
		typ,
	)
}

func (c *Checker) checkPrivateConstantType(node *ast.PrivateConstantNode) *typed.PrivateConstantNode {
	typ, _ := c.resolveType(node.Value, node.Span())
	if typ == nil {
		typ = types.Void{}
	}
	return typed.NewPrivateConstantNode(
		node.Span(),
		node.Value,
		typ,
	)
}

func (c *Checker) checkTypeNode(node ast.TypeNode) typed.TypeNode {
	switch n := node.(type) {
	case *ast.PublicConstantNode:
		return c.checkPublicConstantType(n)
	case *ast.PrivateConstantNode:
		return c.checkPrivateConstantType(n)
	case *ast.ConstantLookupNode:
		return c.constantLookupType(n)
	case *ast.RawStringLiteralNode:
		return typed.NewRawStringLiteralNode(
			n.Span(),
			n.Value,
			types.NewStringLiteral(n.Value),
		)
	case *ast.DoubleQuotedStringLiteralNode:
		return typed.NewDoubleQuotedStringLiteralNode(
			n.Span(),
			n.Value,
			types.NewStringLiteral(n.Value),
		)
	case *ast.RawCharLiteralNode:
		return typed.NewRawCharLiteralNode(
			n.Span(),
			n.Value,
			types.NewCharLiteral(n.Value),
		)
	case *ast.CharLiteralNode:
		return typed.NewCharLiteralNode(
			n.Span(),
			n.Value,
			types.NewCharLiteral(n.Value),
		)
	case *ast.NilLiteralNode:
		return typed.NewNilLiteralNode(n.Span())
	case *ast.SimpleSymbolLiteralNode:
		return typed.NewSimpleSymbolLiteralNode(
			n.Span(),
			n.Content,
			types.NewSymbolLiteral(n.Content),
		)
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
		return typed.NewIntLiteralNode(
			n.Span(),
			n.Value,
			types.NewIntLiteral(n.Value),
		)
	case *ast.Int64LiteralNode:
		return typed.NewInt64LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt64Literal(n.Value),
		)
	case *ast.Int32LiteralNode:
		return typed.NewInt32LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt32Literal(n.Value),
		)
	case *ast.Int16LiteralNode:
		return typed.NewInt16LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt16Literal(n.Value),
		)
	case *ast.Int8LiteralNode:
		return typed.NewInt8LiteralNode(
			n.Span(),
			n.Value,
			types.NewInt8Literal(n.Value),
		)
	case *ast.UInt64LiteralNode:
		return typed.NewUInt64LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt64Literal(n.Value),
		)
	case *ast.UInt32LiteralNode:
		return typed.NewUInt32LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt32Literal(n.Value),
		)
	case *ast.UInt16LiteralNode:
		return typed.NewUInt16LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt16Literal(n.Value),
		)
	case *ast.UInt8LiteralNode:
		return typed.NewUInt8LiteralNode(
			n.Span(),
			n.Value,
			types.NewUInt8Literal(n.Value),
		)
	case *ast.FloatLiteralNode:
		return typed.NewFloatLiteralNode(
			n.Span(),
			n.Value,
			types.NewFloatLiteral(n.Value),
		)
	case *ast.Float64LiteralNode:
		return typed.NewFloat64LiteralNode(
			n.Span(),
			n.Value,
			types.NewFloat64Literal(n.Value),
		)
	case *ast.Float32LiteralNode:
		return typed.NewFloat32LiteralNode(
			n.Span(),
			n.Value,
			types.NewFloat32Literal(n.Value),
		)
	case *ast.TrueLiteralNode:
		return typed.NewTrueLiteralNode(n.Span())
	case *ast.FalseLiteralNode:
		return typed.NewFalseLiteralNode(n.Span())
	case *ast.VoidTypeNode:
		return typed.NewVoidTypeNode(n.Span())
	case *ast.NeverTypeNode:
		return typed.NewNeverTypeNode(n.Span())
	case *ast.AnyTypeNode:
		return typed.NewAnyTypeNode(n.Span())
	case *ast.NilableTypeNode:
		typeNode := c.checkTypeNode(n.Type)
		typ := c.toNilable(c.typeOf(typeNode))
		return typed.NewNilableTypeNode(
			n.Span(),
			c.checkTypeNode(n.Type),
			typ,
		)
	default:
		c.addError(
			fmt.Sprintf("invalid type node %T", node),
			node.Span(),
		)
		return typed.NewInvalidNode(node.Span(), nil)
	}
}

func (c *Checker) constructUnionType(node *ast.BinaryTypeExpressionNode) *typed.UnionTypeNode {
	union := types.NewUnion()
	elements := new([]typed.TypeNode)
	c._constructUnionType(node, elements, union)

	return typed.NewUnionTypeNode(
		node.Span(),
		*elements,
		union,
	)
}

func (c *Checker) _constructUnionType(node *ast.BinaryTypeExpressionNode, elements *[]typed.TypeNode, union *types.Union) {
	leftBinaryType, leftIsBinaryType := node.Left.(*ast.BinaryTypeExpressionNode)
	if leftIsBinaryType && leftBinaryType.Op.Type == token.OR {
		c._constructUnionType(leftBinaryType, elements, union)
	} else {
		leftTypeNode := c.checkTypeNode(node.Left)
		*elements = append(*elements, leftTypeNode)

		leftType := c.typeOf(leftTypeNode)
		union.Elements = append(union.Elements, leftType)
	}

	rightBinaryType, rightIsBinaryType := node.Right.(*ast.BinaryTypeExpressionNode)
	if rightIsBinaryType && rightBinaryType.Op.Type == token.OR {
		c._constructUnionType(rightBinaryType, elements, union)
	} else {
		rightTypeNode := c.checkTypeNode(node.Right)
		*elements = append(*elements, rightTypeNode)

		rightType := c.typeOf(rightTypeNode)
		union.Elements = append(union.Elements, rightType)
	}
}

func (c *Checker) constructIntersectionType(node *ast.BinaryTypeExpressionNode) *typed.IntersectionTypeNode {
	intersection := types.NewIntersection()
	elements := new([]typed.TypeNode)
	c._constructIntersectionType(node, elements, intersection)

	return typed.NewIntersectionTypeNode(
		node.Span(),
		*elements,
		intersection,
	)
}

func (c *Checker) _constructIntersectionType(node *ast.BinaryTypeExpressionNode, elements *[]typed.TypeNode, intersection *types.Intersection) {
	leftBinaryType, leftIsBinaryType := node.Left.(*ast.BinaryTypeExpressionNode)
	if leftIsBinaryType && leftBinaryType.Op.Type == token.AND {
		c._constructIntersectionType(leftBinaryType, elements, intersection)
	} else {
		leftTypeNode := c.checkTypeNode(node.Left)
		*elements = append(*elements, leftTypeNode)

		leftType := c.typeOf(leftTypeNode)
		intersection.Elements = append(intersection.Elements, leftType)
	}

	rightBinaryType, rightIsBinaryType := node.Right.(*ast.BinaryTypeExpressionNode)
	if rightIsBinaryType && rightBinaryType.Op.Type == token.AND {
		c._constructIntersectionType(rightBinaryType, elements, intersection)
	} else {
		rightTypeNode := c.checkTypeNode(node.Right)
		*elements = append(*elements, rightTypeNode)

		rightType := c.typeOf(rightTypeNode)
		intersection.Elements = append(intersection.Elements, rightType)
	}
}

func (c *Checker) constantLookupType(node *ast.ConstantLookupNode) *typed.PublicConstantNode {
	typ, name := c.resolveConstantLookupType(node)
	if typ == nil {
		typ = types.Void{}
	}

	return typed.NewPublicConstantNode(
		node.Span(),
		name,
		typ,
	)
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

func (c *Checker) _resolveConstantLookupForSetter(node *ast.ConstantLookupNode, firstCall bool) (types.ConstantContainer, types.Type, string) {
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
		_, leftContainerType, leftContainerName = c._resolveConstantLookupForSetter(l, false)
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
	if constant == nil {
		if !firstCall {
			c.addError(
				fmt.Sprintf("undefined constant `%s`", constantName),
				node.Right.Span(),
			)
		}
		return leftContainer, nil, constantName
	}

	return leftContainer, constant, constantName
}

func (c *Checker) resolveConstantLookupForSetter(node *ast.ConstantLookupNode) (types.ConstantContainer, types.Type, string) {
	return c._resolveConstantLookupForSetter(node, true)
}

func (c *Checker) constantLookup(node *ast.ConstantLookupNode) *typed.PublicConstantNode {
	typ, name := c.resolveConstantLookup(node)
	if typ == nil {
		typ = types.Void{}
	}

	return typed.NewPublicConstantNode(
		node.Span(),
		name,
		typ,
	)
}

func (c *Checker) publicConstant(node *ast.PublicConstantNode) *typed.PublicConstantNode {
	typ, name := c.resolvePublicConstant(node.Value, node.Span())
	if typ == nil {
		typ = types.Void{}
	}

	return typed.NewPublicConstantNode(
		node.Span(),
		name,
		typ,
	)
}

func (c *Checker) privateConstant(node *ast.PrivateConstantNode) *typed.PrivateConstantNode {
	typ, name := c.resolvePrivateConstant(node.Value, node.Span())
	if typ == nil {
		typ = types.Void{}
	}

	return typed.NewPrivateConstantNode(
		node.Span(),
		name,
		typ,
	)
}

func (c *Checker) publicIdentifier(node *ast.PublicIdentifierNode) *typed.PublicIdentifierNode {
	local, ok := c.resolveLocal(node.Value, node.Span())
	if ok && !local.initialised {
		c.addError(
			fmt.Sprintf("cannot access uninitialised local `%s`", node.Value),
			node.Span(),
		)
	}
	return typed.NewPublicIdentifierNode(
		node.Span(),
		node.Value,
		local.typ,
	)
}

func (c *Checker) privateIdentifier(node *ast.PrivateIdentifierNode) *typed.PrivateIdentifierNode {
	local, ok := c.resolveLocal(node.Value, node.Span())
	if ok && !local.initialised {
		c.addError(
			fmt.Sprintf("cannot access uninitialised local `%s`", node.Value),
			node.Span(),
		)
	}
	return typed.NewPrivateIdentifierNode(
		node.Span(),
		node.Value,
		local.typ,
	)
}

func (c *Checker) variableDeclaration(node *ast.VariableDeclarationNode) *typed.VariableDeclarationNode {
	if _, ok := c.getLocal(node.Name.Value); ok {
		c.addError(
			fmt.Sprintf("cannot redeclare local `%s`", node.Name.Value),
			node.Span(),
		)
	}
	if node.Initialiser == nil {
		if node.Type == nil {
			c.addError(
				fmt.Sprintf("cannot declare a variable without a type `%s`", node.Name.Value),
				node.Span(),
			)
			c.addLocal(node.Name.Value, local{typ: types.Void{}})
			return typed.NewVariableDeclarationNode(
				node.Span(),
				node.Name,
				nil,
				nil,
				types.Void{},
			)
		}

		// without an initialiser but with a type
		declaredTypeNode := c.checkTypeNode(node.Type)
		declaredType := c.typeOf(declaredTypeNode)
		c.addLocal(node.Name.Value, local{typ: declaredType})
		return typed.NewVariableDeclarationNode(
			node.Span(),
			node.Name,
			declaredTypeNode,
			nil,
			types.Void{},
		)
	}

	// with an initialiser
	if node.Type == nil {
		// without a type, inference
		init := c.checkExpression(node.Initialiser)
		actualType := c.typeOf(init).ToNonLiteral(c.GlobalEnv)
		c.addLocal(node.Name.Value, local{typ: actualType, initialised: true})
		if types.IsVoid(actualType) {
			c.addError(
				fmt.Sprintf("cannot declare variable `%s` with type `void`", node.Name.Value),
				init.Span(),
			)
		}
		return typed.NewVariableDeclarationNode(
			node.Span(),
			node.Name,
			nil,
			init,
			actualType,
		)
	}

	// with a type and an initializer

	declaredTypeNode := c.checkTypeNode(node.Type)
	declaredType := c.typeOf(declaredTypeNode)
	init := c.checkExpression(node.Initialiser)
	actualType := c.typeOf(init)
	c.addLocal(node.Name.Value, local{typ: declaredType, initialised: true})
	if !c.isSubtype(actualType, declaredType) {
		c.addError(
			fmt.Sprintf("type `%s` cannot be assigned to type `%s`", types.Inspect(actualType), types.Inspect(declaredType)),
			init.Span(),
		)
	}

	return typed.NewVariableDeclarationNode(
		node.Span(),
		node.Name,
		declaredTypeNode,
		init,
		declaredType,
	)
}

func (c *Checker) valueDeclaration(node *ast.ValueDeclarationNode) *typed.ValueDeclarationNode {
	if _, ok := c.getLocal(node.Name.Value); ok {
		c.addError(
			fmt.Sprintf("cannot redeclare local `%s`", node.Name.Value),
			node.Span(),
		)
	}
	if node.Initialiser == nil {
		if node.Type == nil {
			c.addError(
				fmt.Sprintf("cannot declare a value without a type `%s`", node.Name.Value),
				node.Span(),
			)
			c.addLocal(node.Name.Value, local{typ: types.Void{}})
			return typed.NewValueDeclarationNode(
				node.Span(),
				node.Name,
				nil,
				nil,
				types.Void{},
			)
		}

		// without an initialiser but with a type
		declaredTypeNode := c.checkTypeNode(node.Type)
		declaredType := c.typeOf(declaredTypeNode)
		c.addLocal(node.Name.Value, local{typ: declaredType, singleAssignment: true})
		return typed.NewValueDeclarationNode(
			node.Span(),
			node.Name,
			declaredTypeNode,
			nil,
			types.Void{},
		)
	}

	// with an initialiser
	if node.Type == nil {
		// without a type, inference
		init := c.checkExpression(node.Initialiser)
		actualType := c.typeOf(init)
		c.addLocal(node.Name.Value, local{typ: actualType, initialised: true, singleAssignment: true})
		if types.IsVoid(actualType) {
			c.addError(
				fmt.Sprintf("cannot declare value `%s` with type `void`", node.Name.Value),
				init.Span(),
			)
		}
		return typed.NewValueDeclarationNode(
			node.Span(),
			node.Name,
			nil,
			init,
			actualType,
		)
	}

	// with a type and an initializer

	declaredTypeNode := c.checkTypeNode(node.Type)
	declaredType := c.typeOf(declaredTypeNode)
	init := c.checkExpression(node.Initialiser)
	actualType := c.typeOf(init)
	c.addLocal(node.Name.Value, local{typ: declaredType, initialised: true, singleAssignment: true})
	if !c.isSubtype(actualType, declaredType) {
		c.addError(
			fmt.Sprintf("type `%s` cannot be assigned to type `%s`", types.Inspect(actualType), types.Inspect(declaredType)),
			init.Span(),
		)
	}

	return typed.NewValueDeclarationNode(
		node.Span(),
		node.Name,
		declaredTypeNode,
		init,
		declaredType,
	)
}

func (c *Checker) constantDeclaration(node *ast.ConstantDeclarationNode) *typed.ConstantDeclarationNode {
	if constant, constantName := c.resolveConstantForSetter(node.Name.Value); constant != nil {
		c.addError(
			fmt.Sprintf("cannot redeclare constant `%s`", constantName),
			node.Span(),
		)
	}

	scope := c.currentConstScope()
	if node.Type == nil {
		// without a type, inference
		init := c.checkExpression(node.Initialiser)
		actualType := c.typeOf(init)
		scope.container.DefineConstant(node.Name.Value, actualType)
		return typed.NewConstantDeclarationNode(
			node.Span(),
			node.Name,
			nil,
			init,
			actualType,
		)
	}

	// with a type

	declaredTypeNode := c.checkTypeNode(node.Type)
	declaredType := c.typeOf(declaredTypeNode)
	init := c.checkExpression(node.Initialiser)
	actualType := c.typeOf(init)
	scope.container.DefineConstant(node.Name.Value, declaredType)
	if !c.isSubtype(actualType, declaredType) {
		c.addError(
			fmt.Sprintf("type `%s` cannot be assigned to type `%s`", types.Inspect(actualType), types.Inspect(declaredType)),
			init.Span(),
		)
	}

	return typed.NewConstantDeclarationNode(
		node.Span(),
		node.Name,
		declaredTypeNode,
		init,
		declaredType,
	)
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
	constantModule, constantIsModule := constantType.(*types.Module)
	if constantType != nil {
		if !constantIsModule {
			c.addError(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				span,
			)
			return types.NewModule(fullConstantName, nil, nil, nil)
		} else {
			return constantModule
		}
	} else if constantContainer == nil {
		return types.NewModule(fullConstantName, nil, nil, nil)
	} else {
		return constantContainer.DefineModule(constantName, nil, nil, nil)
	}
}

func (c *Checker) module(node *ast.ModuleDeclarationNode) *typed.ModuleDeclarationNode {
	var typedConstantNode typed.ExpressionNode
	var module *types.Module

	switch constant := node.Constant.(type) {
	case *ast.PublicConstantNode:
		constScope := c.currentConstScope()
		constantType, fullConstantName := c.resolveConstantForSetter(constant.Value)
		module = c.declareModule(constScope.container, constantType, fullConstantName, constant.Value, node.Span())
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			fullConstantName,
			module,
		)
	case *ast.PrivateConstantNode:
		constScope := c.currentConstScope()
		constantType, fullConstantName := c.resolveConstantForSetter(constant.Value)
		module = c.declareModule(constScope.container, constantType, fullConstantName, constant.Value, node.Span())
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			fullConstantName,
			module,
		)
	case *ast.ConstantLookupNode:
		constantContainer, constantType, fullConstantName := c.resolveConstantLookupForSetter(constant)
		constantName := extractConstantNameFromLookup(constant)
		module = c.declareModule(constantContainer, constantType, fullConstantName, constantName, node.Span())
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			fullConstantName,
			module,
		)
	case nil:
		module = types.NewModule("", nil, nil, nil)
	default:
		c.addError(
			fmt.Sprintf("invalid module name node %T", node.Constant),
			node.Constant.Span(),
		)
	}

	c.pushConstScope(makeLocalConstantScope(module))
	c.pushMethodScope(makeLocalMethodScope(module))
	prevSelfType := c.selfType
	c.selfType = module

	previousMode := c.mode
	c.mode = moduleMode
	defer c.setMode(previousMode)

	newBody := c.hoistStatements(node.Body)
	c.selfType = prevSelfType
	c.popConstScope()
	c.popMethodScope()

	return typed.NewModuleDeclarationNode(
		node.Span(),
		typedConstantNode,
		newBody,
		module,
	)
}

func (c *Checker) declareClass(superclass *types.Class, constantContainer types.ConstantContainer, constantType types.Type, fullConstantName, constantName string, span *position.Span) *types.Class {
	constantSingleton, constantIsSingleton := constantType.(*types.SingletonClass)
	var superclassConstantContainer types.ConstantContainer
	if superclass == nil {
		superclassConstantContainer = c.GlobalEnv.StdSubtypeClass(symbol.Object)
	} else {
		superclassConstantContainer = superclass
	}

	if constantType == nil {
		return constantContainer.DefineClass(constantName, superclassConstantContainer, nil, nil)
	}

	if !constantIsSingleton {
		c.addError(
			fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
			span,
		)
		return types.NewClass(fullConstantName, superclassConstantContainer, nil, nil)
	}

	constantClass, constantIsClass := constantSingleton.AttachedObject.(*types.Class)
	if !constantIsClass {
		c.addError(
			fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
			span,
		)
		return types.NewClass(constantName, superclassConstantContainer, nil, nil)
	}

	if superclass == nil {
		superclass = c.GlobalEnv.StdSubtypeClass(symbol.Object)
	}

	if constantClass.Parent() != superclass {
		c.addError(
			fmt.Sprintf("superclass mismatch in `%s`, got `%s`, expected `%s`", fullConstantName, superclass.Name(), constantClass.Parent().Name()),
			span,
		)
		return constantClass
	}

	return constantClass
}

func (c *Checker) class(node *ast.ClassDeclarationNode) *typed.ClassDeclarationNode {
	var typedConstantNode typed.ExpressionNode
	var class *types.Class

	var superclassNode typed.ComplexConstantNode
	var superclass *types.Class

	if node.Superclass != nil {
		superclassNode = c.checkComplexConstantType(node.Superclass.(ast.ComplexConstantNode))
		superclassType := c.typeOf(superclassNode)
		var superclassIsClass bool
		superclass, superclassIsClass = superclassType.(*types.Class)
		if !superclassIsClass {
			c.addError(
				fmt.Sprintf("`%s` is not a class", types.Inspect(superclassType)),
				node.Superclass.Span(),
			)
		}
	}

	switch constant := node.Constant.(type) {
	case *ast.PublicConstantNode:
		constScope := c.currentConstScope()
		constantType, fullConstantName := c.resolveConstantForSetter(constant.Value)
		class = c.declareClass(superclass, constScope.container, constantType, fullConstantName, constant.Value, node.Constant.Span())
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			fullConstantName,
			class,
		)
	case *ast.PrivateConstantNode:
		constScope := c.currentConstScope()
		constantType, fullConstantName := c.resolveConstantForSetter(constant.Value)
		class = c.declareClass(superclass, constScope.container, constantType, fullConstantName, constant.Value, node.Constant.Span())
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			fullConstantName,
			class,
		)
	case *ast.ConstantLookupNode:
		constantContainer, constantType, fullConstantName := c.resolveConstantLookupForSetter(constant)
		constantName := extractConstantNameFromLookup(constant)
		class = c.declareClass(superclass, constantContainer, constantType, fullConstantName, constantName, node.Constant.Span())
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			fullConstantName,
			class,
		)
	case nil:
		class = types.NewClass("", superclass, nil, nil)
	default:
		c.addError(
			fmt.Sprintf("invalid class name node %T", node.Constant),
			node.Constant.Span(),
		)
	}

	c.pushConstScope(makeLocalConstantScope(class))
	c.pushMethodScope(makeLocalMethodScope(class))
	prevSelfType := c.selfType
	c.selfType = types.NewSingletonClass(class)

	previousMode := c.mode
	c.mode = classMode
	defer c.setMode(previousMode)

	newBody := c.hoistStatements(node.Body)

	c.selfType = prevSelfType
	c.popConstScope()
	c.popMethodScope()

	return typed.NewClassDeclarationNode(
		node.Span(),
		node.Abstract,
		node.Sealed,
		typedConstantNode,
		nil,
		superclassNode,
		newBody,
		class,
	)
}

// Extract method, class, module, mixin definitions etc
// and hoist them to the top of the block.
func (c *Checker) hoistStatements(statements []ast.StatementNode) []typed.StatementNode {
	var methodStatements []*ast.ExpressionStatementNode
	var oldMethods []*types.Method
	var newMethods []*types.Method
	var otherStatements []ast.StatementNode

	for _, statement := range statements {
		s, ok := statement.(*ast.ExpressionStatementNode)
		if !ok {
			otherStatements = append(otherStatements, statement)
			continue
		}

		switch expr := s.Expression.(type) {
		case *ast.MethodDefinitionNode:
			oldMethod, newMethod := c.declareMethod(
				expr.Name,
				expr.Parameters,
				expr.ReturnType,
				expr.ThrowType,
			)
			newMethods = append(newMethods, newMethod)
			oldMethods = append(oldMethods, oldMethod)
			methodStatements = append(methodStatements, s)
		case *ast.DocCommentNode:
			switch e := expr.Expression.(type) {
			case *ast.MethodDefinitionNode:
				oldMethod, newMethod := c.declareMethod(
					e.Name,
					e.Parameters,
					e.ReturnType,
					e.ThrowType,
				)
				newMethods = append(newMethods, newMethod)
				oldMethods = append(oldMethods, oldMethod)
				methodStatements = append(methodStatements, s)
			default:
				otherStatements = append(otherStatements, s)
			}
		default:
			otherStatements = append(otherStatements, s)
		}
	}

	typedStatements := make([]typed.StatementNode, len(methodStatements))
	var wg sync.WaitGroup
	wg.Add(len(methodStatements))

	for i, methodStatement := range methodStatements {
		go func() {
			defer wg.Done()
			typedStatements[i] = c.checkMethodStatement(methodStatement, oldMethods[i], newMethods[i])
		}()
	}
	wg.Wait()

	typedOtherStatements := c.checkStatements(otherStatements)
	typedStatements = append(typedStatements, typedOtherStatements...)

	if len(typedStatements) == 0 {
		return nil
	}
	return typedStatements
}

func (c *Checker) checkMethodDefinition(node *ast.MethodDefinitionNode, oldMethod, newMethod *types.Method) *typed.MethodDefinitionNode {
	paramNodes, returnTypeNode, throwTypeNode, body := c.checkMethod(
		oldMethod,
		newMethod,
		node.Name,
		node.Parameters,
		node.ReturnType,
		node.ThrowType,
		node.Body,
		node.Span(),
	)

	return typed.NewMethodDefinitionNode(
		node.Span(),
		node.Name,
		paramNodes,
		returnTypeNode,
		throwTypeNode,
		body,
	)
}

func (c *Checker) checkMethodStatement(methodStatement *ast.ExpressionStatementNode, oldMethod, newMethod *types.Method) *typed.ExpressionStatementNode {
	var typedExpression typed.ExpressionNode

	switch expr := methodStatement.Expression.(type) {
	case *ast.MethodDefinitionNode:
		typedExpression = c.checkMethodDefinition(expr, oldMethod, newMethod)
	case *ast.DocCommentNode:
		switch e := expr.Expression.(type) {
		case *ast.MethodDefinitionNode:
			typedExpression = typed.NewDocCommentNode(
				expr.Span(),
				expr.Comment,
				c.checkMethodDefinition(e, oldMethod, newMethod),
			)
		default:
			panic(
				fmt.Sprintf("invalid doc comment method expression node: %T", methodStatement.Expression),
			)
		}
	default:
		panic(
			fmt.Sprintf("invalid method expression node: %T", methodStatement.Expression),
		)
	}

	return typed.NewExpressionStatementNode(
		methodStatement.Span(),
		typedExpression,
	)
}

func (c *Checker) declareMethod(
	name string,
	paramNodes []ast.ParameterNode,
	returnTypeNode,
	throwTypeNode ast.TypeNode,
) (oldMethod, newMethod *types.Method) {
	methodScope := c.currentMethodScope()
	oldMethod = c.getMethod(methodScope.container, name, nil, false)

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
		var declaredTypeNode typed.TypeNode
		if p.Type == nil {
			c.addError(
				fmt.Sprintf("cannot declare parameter `%s` without a type", p.Name),
				param.Span(),
			)
		} else {
			declaredTypeNode = c.checkTypeNode(p.Type)
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
	var typedReturnTypeNode typed.TypeNode
	if returnTypeNode != nil {
		typedReturnTypeNode = c.checkTypeNode(returnTypeNode)
		returnType = c.typeOf(typedReturnTypeNode)
	} else {
		returnType = types.Void{}
	}

	var throwType types.Type
	var typedThrowTypeNode typed.TypeNode
	if throwTypeNode != nil {
		typedThrowTypeNode = c.checkTypeNode(throwTypeNode)
		throwType = c.typeOf(typedThrowTypeNode)
	}
	newMethod = types.NewMethod(
		name,
		params,
		returnType,
		throwType,
		methodScope.container,
	)

	methodScope.container.SetMethod(name, newMethod)
	return oldMethod, newMethod
}

func (c *Checker) declareMixin(constantContainer types.ConstantContainer, constantType types.Type, fullConstantName, constantName string, span *position.Span) *types.Mixin {
	constantMixin, constantIsMixin := constantType.(*types.Mixin)
	if constantType != nil {
		if !constantIsMixin {
			c.addError(
				fmt.Sprintf("cannot redeclare constant `%s`", fullConstantName),
				span,
			)
			return types.NewMixin(fullConstantName, nil, nil, nil, nil)
		} else {
			return constantMixin
		}
	} else if constantContainer == nil {
		return types.NewMixin(fullConstantName, nil, nil, nil, nil)
	} else {
		return constantContainer.DefineMixin(constantName, nil, nil, nil, nil)
	}
}

func (c *Checker) mixin(node *ast.MixinDeclarationNode) *typed.MixinDeclarationNode {
	var typedConstantNode typed.ExpressionNode
	var mixin *types.Mixin

	switch constant := node.Constant.(type) {
	case *ast.PublicConstantNode:
		constScope := c.currentConstScope()
		constantType, fullConstantName := c.resolveConstantForSetter(constant.Value)
		mixin = c.declareMixin(constScope.container, constantType, fullConstantName, constant.Value, node.Span())
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			fullConstantName,
			mixin,
		)
	case *ast.PrivateConstantNode:
		constScope := c.currentConstScope()
		constantType, fullConstantName := c.resolveConstantForSetter(constant.Value)
		mixin = c.declareMixin(constScope.container, constantType, fullConstantName, constant.Value, node.Span())
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			fullConstantName,
			mixin,
		)
	case *ast.ConstantLookupNode:
		constantContainer, constantType, fullConstantName := c.resolveConstantLookupForSetter(constant)
		constantName := extractConstantNameFromLookup(constant)
		mixin = c.declareMixin(constantContainer, constantType, fullConstantName, constantName, node.Span())
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			fullConstantName,
			mixin,
		)
	case nil:
		mixin = types.NewMixin("", nil, nil, nil, nil)
	default:
		c.addError(
			fmt.Sprintf("invalid mixin name node %T", node.Constant),
			node.Constant.Span(),
		)
	}

	c.pushConstScope(makeLocalConstantScope(mixin))
	c.pushMethodScope(makeLocalMethodScope(mixin))
	prevSelfType := c.selfType
	c.selfType = types.NewSingletonClass(mixin)

	previousMode := c.mode
	c.mode = mixinMode
	defer c.setMode(previousMode)

	newBody := c.hoistStatements(node.Body)
	c.selfType = prevSelfType
	c.popConstScope()
	c.popMethodScope()

	return typed.NewMixinDeclarationNode(
		node.Span(),
		typedConstantNode,
		nil,
		newBody,
		mixin,
	)
}
