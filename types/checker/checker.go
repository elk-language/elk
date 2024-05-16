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
	HeaderMode     bool
	constantScopes []constantScope
	localEnvs      []*localEnvironment
	returnType     types.Type
	throwType      types.Type
}

// Instantiate a new Checker instance.
func new(loc *position.Location, globalEnv *types.GlobalEnvironment, headerMode bool) *Checker {
	if globalEnv == nil {
		globalEnv = types.NewGlobalEnvironment()
	}
	return &Checker{
		Location:   loc,
		GlobalEnv:  globalEnv,
		HeaderMode: headerMode,
		constantScopes: []constantScope{
			makeConstantScope(globalEnv.Std()),
			makeLocalConstantScope(globalEnv.Root),
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
		GlobalEnv: globalEnv,
		constantScopes: []constantScope{
			makeConstantScope(globalEnv.Std()),
			makeLocalConstantScope(globalEnv.Root),
		},
		localEnvs: []*localEnvironment{
			newLocalEnvironment(nil),
		},
	}
}

func (c *Checker) CheckSource(sourceName string, source string) (typed.Node, errors.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	loc := position.NewLocationWithSpan(sourceName, ast.Span())
	c.Location = loc
	typedAst := c.checkProgram(ast)
	return typedAst, c.Errors
}

func (c *Checker) ClearErrors() {
	c.Errors = nil
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
	case *ast.MethodDefinitionNode:
		return c.methodDefinition(n)
	default:
		c.addError(
			fmt.Sprintf("invalid expression type %T", node),
			node.Span(),
		)
		return typed.NewInvalidNode(node.Span(), nil)
	}
}

func (c *Checker) typeOf(node typed.Node) types.Type {
	return typed.TypeOf(node, c.GlobalEnv)
}

func (c *Checker) methodDefinition(node *ast.MethodDefinitionNode) *typed.MethodDefinitionNode {
	constScope := c.currentConstScope()
	env := newLocalEnvironment(nil)
	c.pushLocalEnv(env)
	defer c.popLocalEnv()

	oldMethod := constScope.container.MethodString(node.Name)

	var paramNodes []typed.ParameterNode
	var params []*types.Parameter
	for _, param := range node.Parameters {
		p, ok := param.(*ast.MethodParameterNode)
		if !ok {
			c.addError(
				fmt.Sprintf("invalid param type %T", node),
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
			if !initType.IsSubtypeOf(declaredType, c.GlobalEnv) {
				c.addError(
					fmt.Sprintf("type `%s` cannot be assigned to type `%s`", initType.Inspect(), declaredType.Inspect()),
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
		paramNodes = append(paramNodes, typed.NewMethodParameterNode(
			p.Span(),
			p.Name,
			p.SetInstanceVariable,
			declaredTypeNode,
			initNode,
			typed.ParameterKind(p.Kind),
		))
	}

	var returnType types.Type
	var returnTypeNode typed.TypeNode
	if node.ReturnType != nil {
		returnTypeNode = c.checkTypeNode(node.ReturnType)
		returnType = c.typeOf(returnTypeNode)
	} else {
		returnType = types.Void{}
	}

	var throwType types.Type
	var throwTypeNode typed.TypeNode
	if node.ThrowType != nil {
		throwTypeNode = c.checkTypeNode(node.ThrowType)
		throwType = c.typeOf(throwTypeNode)
	}
	newMethod := types.NewMethod(
		node.Name,
		params,
		returnType,
		throwType,
	)

	setMethod := true
	if oldMethod != nil {
		if returnTypeNode != nil && oldMethod.ReturnType != nil && newMethod.ReturnType != oldMethod.ReturnType {
			c.addError(
				fmt.Sprintf("cannot redeclare method `%s` with a different return type, is `%s`, should be `%s`", node.Name, types.Inspect(newMethod.ReturnType), types.Inspect(oldMethod.ReturnType)),
				returnTypeNode.Span(),
			)
			setMethod = false
		}
		if newMethod.ThrowType != oldMethod.ThrowType {
			var span *position.Span
			if throwTypeNode != nil {
				span = throwTypeNode.Span()
			} else {
				span = node.Span()
			}
			c.addError(
				fmt.Sprintf("cannot redeclare method `%s` with a different throw type, is `%s`, should be `%s`", node.Name, types.Inspect(newMethod.ThrowType), types.Inspect(oldMethod.ThrowType)),
				span,
			)
			setMethod = false
		}

		if len(oldMethod.Params) > len(newMethod.Params) {
			c.addError(
				fmt.Sprintf("cannot redeclare method `%s` with less parameters", node.Name),
				position.JoinSpanOfCollection(node.Parameters),
			)
			setMethod = false
		} else {
			for i := range len(oldMethod.Params) {
				oldParam := oldMethod.Params[i]
				newParam := newMethod.Params[i]
				if oldParam.Name != newParam.Name {
					c.addError(
						fmt.Sprintf("cannot redeclare method `%s` with invalid parameter name, is `%s`, should be `%s`", node.Name, newParam.Name, oldParam.Name),
						paramNodes[i].Span(),
					)
					setMethod = false
					continue
				}
				if oldParam.Kind != newParam.Kind {
					c.addError(
						fmt.Sprintf("cannot redeclare method `%s` with invalid parameter kind, is `%s`, should be `%s`", node.Name, newParam.NameWithKind(), oldParam.NameWithKind()),
						paramNodes[i].Span(),
					)
					setMethod = false
					continue
				}
				if oldParam.Type != newParam.Type {
					c.addError(
						fmt.Sprintf("cannot redeclare method `%s` with invalid parameter type, is `%s`, should be `%s`", node.Name, types.Inspect(newParam.Type), types.Inspect(oldParam.Type)),
						paramNodes[i].Span(),
					)
					setMethod = false
					continue
				}
			}

			for i := len(oldMethod.Params); i < len(newMethod.Params); i++ {
				param := newMethod.Params[i]
				if !param.IsOptional() {
					c.addError(
						fmt.Sprintf("cannot redeclare method `%s` with additional required parameter `%s`", node.Name, param.Name),
						paramNodes[i].Span(),
					)
					setMethod = false
				}
			}
		}

	}

	if setMethod {
		constScope.container.SetMethod(node.Name, newMethod)
	}

	c.returnType = returnType
	c.throwType = throwType
	body := c.checkStatements(node.Body)
	c.returnType = nil
	c.throwType = nil

	return typed.NewMethodDefinitionNode(
		node.Span(),
		node.Name,
		paramNodes,
		returnTypeNode,
		throwTypeNode,
		body,
	)
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

	constant := leftContainer.ConstantString(rightName)
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
	if !actualType.IsSubtypeOf(declaredType, c.GlobalEnv) {
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
	if !actualType.IsSubtypeOf(declaredType, c.GlobalEnv) {
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
	scope.container.DefineConstant(node.Name.Value, actualType)
	if !actualType.IsSubtypeOf(declaredType, c.GlobalEnv) {
		c.addError(
			fmt.Sprintf("type `%s` cannot be assigned to type `%s`", actualType.Inspect(), declaredType.Inspect()),
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
