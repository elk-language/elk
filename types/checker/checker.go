// Package checker implements the Elk type checker
package checker

import (
	"fmt"

	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/types"
	typed "github.com/elk-language/elk/types/ast" // typed AST
	"github.com/elk-language/elk/value"
)

// Check the types of Elk source code.
func CheckSource(sourceName string, source string, globalEnv *types.GlobalEnvironment, headerMode bool) (typed.Node, errors.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	return CheckAST(sourceName, ast, globalEnv, headerMode)
}

// Check the types of an Elk AST.
func CheckAST(sourceName string, ast *ast.ProgramNode, globalEnv *types.GlobalEnvironment, headerMode bool) (typed.Node, errors.ErrorList) {
	checker := new(position.NewLocationWithSpan(sourceName, ast.Span()), globalEnv, headerMode)
	typedAst := checker.checkProgram(ast)
	return typedAst, checker.Errors
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

// Holds the state of the type checking process
type Checker struct {
	Location       *position.Location
	Errors         errors.ErrorList
	GlobalEnv      *types.GlobalEnvironment
	constantScopes []constantScope
	localEnvs      []*localEnvironment
	headerMode     bool
}

// Instantiate a new Checker instance.
func new(loc *position.Location, globalEnv *types.GlobalEnvironment, headerMode bool) *Checker {
	if globalEnv == nil {
		globalEnv = types.NewGlobalEnvironment()
	}
	return &Checker{
		Location:   loc,
		GlobalEnv:  globalEnv,
		headerMode: headerMode,
		constantScopes: []constantScope{
			makeConstantScope(globalEnv.Std()),
			makeLocalConstantScope(globalEnv.Root),
		},
		localEnvs: []*localEnvironment{
			newLocalEnvironment(nil),
		},
	}
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

// Create a new location struct with the given position.
func (c *Checker) newLocation(span *position.Span) *position.Location {
	return position.NewLocationWithSpan(c.Location.Filename, span)
}

func (c *Checker) checkProgram(node *ast.ProgramNode) *typed.ProgramNode {
	newStatements := c.checkStatements(node.Body)
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

func (c *Checker) checkExpression(node ast.ExpressionNode) typed.ExpressionNode {
	switch n := node.(type) {
	case *ast.FalseLiteralNode:
		return typed.NewFalseLiteralNode(n.Span())
	case *ast.TrueLiteralNode:
		return typed.NewTrueLiteralNode(n.Span())
	case *ast.NilLiteralNode:
		return typed.NewNilLiteralNode(n.Span())
	case *ast.IntLiteralNode:
		return typed.NewIntLiteralNode(
			n.Span(),
			n.Value,
		)
	case *ast.Int64LiteralNode:
		return typed.NewInt64LiteralNode(
			n.Span(),
			n.Value,
		)
	case *ast.Int32LiteralNode:
		return typed.NewInt32LiteralNode(
			n.Span(),
			n.Value,
		)
	case *ast.Int16LiteralNode:
		return typed.NewInt16LiteralNode(
			n.Span(),
			n.Value,
		)
	case *ast.Int8LiteralNode:
		return typed.NewInt8LiteralNode(
			n.Span(),
			n.Value,
		)
	case *ast.UInt64LiteralNode:
		return typed.NewUInt64LiteralNode(
			n.Span(),
			n.Value,
		)
	case *ast.UInt32LiteralNode:
		return typed.NewUInt32LiteralNode(
			n.Span(),
			n.Value,
		)
	case *ast.UInt16LiteralNode:
		return typed.NewUInt16LiteralNode(
			n.Span(),
			n.Value,
		)
	case *ast.UInt8LiteralNode:
		return typed.NewUInt8LiteralNode(
			n.Span(),
			n.Value,
		)
	case *ast.FloatLiteralNode:
		return typed.NewFloatLiteralNode(
			n.Span(),
			n.Value,
		)
	case *ast.Float64LiteralNode:
		return typed.NewFloat64LiteralNode(
			n.Span(),
			n.Value,
		)
	case *ast.Float32LiteralNode:
		return typed.NewFloat32LiteralNode(
			n.Span(),
			n.Value,
		)
	case *ast.BigFloatLiteralNode:
		return typed.NewBigFloatLiteralNode(
			n.Span(),
			n.Value,
		)
	case *ast.DoubleQuotedStringLiteralNode:
		return typed.NewDoubleQuotedStringLiteralNode(
			n.Span(),
			n.Value,
		)
	case *ast.VariableDeclarationNode:
		return c.variableDeclaration(n)
	case *ast.ValueDeclarationNode:
		return c.valueDeclaration(n)
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
	default:
		c.addError(
			fmt.Sprintf("invalid expression type %T", node),
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
	constant := constScope.container.Constant(name)
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
		constant := constScope.container.Constant(name)
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
		constant := constScope.container.Constant(name)
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
		constant := constScope.container.Subtype(name)
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

	constant := leftContainer.Subtype(rightName)
	if constant == nil {
		c.addError(
			fmt.Sprintf("undefined type `%s`", typeName),
			node.Right.Span(),
		)
		return nil, typeName
	}

	return constant, typeName
}

func (c *Checker) checkTypeNode(node ast.TypeNode) typed.TypeNode {
	switch n := node.(type) {
	case *ast.PublicConstantNode:
		typ, _ := c.resolveType(n.Value, n.Span())
		if typ == nil {
			typ = types.Void{}
		}
		return typed.NewPublicConstantNode(
			n.Span(),
			n.Value,
			typ,
		)
	case *ast.PrivateConstantNode:
		typ, _ := c.resolveType(n.Value, n.Span())
		if typ == nil {
			typ = types.Void{}
		}
		return typed.NewPrivateConstantNode(
			n.Span(),
			n.Value,
			typ,
		)
	case *ast.ConstantLookupNode:
		return c.constantLookupType(n)
	default:
		c.addError(
			fmt.Sprintf("invalid type node %T", node),
			node.Span(),
		)
		return typed.NewInvalidNode(node.Span(), nil)
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

	constant := leftContainer.Constant(rightName)
	if constant == nil {
		c.addError(
			fmt.Sprintf("undefined constant `%s`", constantName),
			node.Right.Span(),
		)
		return nil, constantName
	}

	return constant, constantName
}

func (c *Checker) _resolveConstantLookupForSetter(node *ast.ConstantLookupNode, firstCall bool) (types.Type, string) {
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
		leftContainerType, leftContainerName = c._resolveConstantLookupForSetter(l, false)
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

	constant := leftContainer.Constant(rightName)
	if constant == nil {
		if !firstCall {
			c.addError(
				fmt.Sprintf("undefined constant `%s`", constantName),
				node.Right.Span(),
			)
		}
		return nil, constantName
	}

	return constant, constantName
}

func (c *Checker) resolveConstantLookupForSetter(node *ast.ConstantLookupNode) (types.Type, string) {
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
		declaredType := typed.TypeOf(declaredTypeNode, c.GlobalEnv)
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
		actualType := typed.TypeOf(init, c.GlobalEnv)
		c.addLocal(node.Name.Value, local{typ: actualType, initialised: true})
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
	declaredType := typed.TypeOf(declaredTypeNode, c.GlobalEnv)
	init := c.checkExpression(node.Initialiser)
	actualType := typed.TypeOf(init, c.GlobalEnv)
	c.addLocal(node.Name.Value, local{typ: declaredType, initialised: true})
	if !declaredType.IsSupertypeOf(actualType) {
		c.addError(
			fmt.Sprintf("type `%s` cannot be assigned to type `%s`", actualType.Inspect(), declaredType.Inspect()),
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
		declaredType := typed.TypeOf(declaredTypeNode, c.GlobalEnv)
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
		actualType := typed.TypeOf(init, c.GlobalEnv)
		c.addLocal(node.Name.Value, local{typ: actualType, initialised: true, singleAssignment: true})
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
	declaredType := typed.TypeOf(declaredTypeNode, c.GlobalEnv)
	init := c.checkExpression(node.Initialiser)
	actualType := typed.TypeOf(init, c.GlobalEnv)
	c.addLocal(node.Name.Value, local{typ: declaredType, initialised: true, singleAssignment: true})
	if !declaredType.IsSupertypeOf(actualType) {
		c.addError(
			fmt.Sprintf("type `%s` cannot be assigned to type `%s`", actualType.Inspect(), declaredType.Inspect()),
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

func (c *Checker) module(node *ast.ModuleDeclarationNode) *typed.ModuleDeclarationNode {
	constScope := c.currentConstScope()
	var typedConstantNode typed.ExpressionNode
	var module *types.Module

	switch constant := node.Constant.(type) {
	case *ast.PublicConstantNode:
		constantType, constantName := c.resolveConstantForSetter(constant.Value)
		constantModule, constantIsModule := constantType.(*types.Module)
		if constantType != nil {
			if !constantIsModule {
				c.addError(
					fmt.Sprintf("cannot redeclare constant `%s`", constantName),
					node.Constant.Span(),
				)
			}
			module = constantModule
		} else {
			module = constScope.container.DefineModule(constant.Value, nil, nil)
		}
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			constantName,
			module,
		)
	case *ast.PrivateConstantNode:
		constantType, constantName := c.resolveConstantForSetter(constant.Value)
		constantModule, constantIsModule := constantType.(*types.Module)
		if constantType != nil {
			if !constantIsModule {
				c.addError(
					fmt.Sprintf("cannot redeclare constant `%s`", constantName),
					node.Constant.Span(),
				)
			}
			module = constantModule
		} else {
			module = constScope.container.DefineModule(constant.Value, nil, nil)
		}
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			constantName,
			module,
		)
	case *ast.ConstantLookupNode:
		constantType, constantName := c.resolveConstantLookupForSetter(constant)
		constantModule, constantIsModule := constantType.(*types.Module)
		if constantType != nil {
			if !constantIsModule {
				c.addError(
					fmt.Sprintf("cannot redeclare constant `%s`", constantName),
					node.Constant.Span(),
				)
			}
			module = constantModule
		} else {
			module = constScope.container.DefineModule(constantName, nil, nil)
		}
		typedConstantNode = typed.NewPublicConstantNode(
			constant.Span(),
			constantName,
			module,
		)
	case nil:
		module = types.NewModule("", nil, nil)
	default:
		c.addError(
			fmt.Sprintf("invalid module name node %T", node.Constant),
			node.Constant.Span(),
		)
	}

	c.pushConstScope(makeLocalConstantScope(module))
	newBody := c.checkStatements(node.Body)
	c.popConstScope()

	return typed.NewModuleDeclarationNode(
		node.Span(),
		typedConstantNode,
		newBody,
		module,
	)
}
