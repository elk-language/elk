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
func CheckSource(sourceName string, source string, globalEnv *types.GlobalEnvironment) (typed.Node, errors.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	return CheckAST(sourceName, ast, globalEnv)
}

// Check the types of an Elk AST.
func CheckAST(sourceName string, ast *ast.ProgramNode, globalEnv *types.GlobalEnvironment) (typed.Node, errors.ErrorList) {
	checker := new(position.NewLocationWithSpan(sourceName, ast.Span()), globalEnv)
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

func newLocalEnvironment(parent *localEnvironment) *localEnvironment {
	return &localEnvironment{
		parent: parent,
		locals: make(map[value.Symbol]local),
	}
}

// Holds the state of the type checking process
type Checker struct {
	Location           *position.Location
	Errors             errors.ErrorList
	GlobalEnv          *types.GlobalEnvironment
	constantContainers []types.ConstantContainer
	localEnvs          []*localEnvironment
}

// Instantiate a new Checker instance.
func new(loc *position.Location, globalEnv *types.GlobalEnvironment) *Checker {
	if globalEnv == nil {
		globalEnv = types.NewGlobalEnvironment()
	}
	return &Checker{
		Location:  loc,
		GlobalEnv: globalEnv,
		constantContainers: []types.ConstantContainer{
			globalEnv.Root,
			globalEnv.Std(),
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

// Create a new location struct with the given position.
func (c *Checker) newLocation(span *position.Span) *position.Location {
	return position.NewLocationWithSpan(c.Location.Filename, span)
}

func (c *Checker) checkProgram(node *ast.ProgramNode) *typed.ProgramNode {
	var newStatements []typed.StatementNode
	for _, statement := range node.Body {
		newStatements = append(newStatements, c.checkStatement(statement))
	}
	return typed.NewProgramNode(
		node.Span(),
		newStatements,
	)
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
		c.Errors.Add(
			fmt.Sprintf("incorrect statement type %#v", node),
			c.newLocation(node.Span()),
		)
		return typed.NewInvalidNode(node.Span(), nil)
	}
}

func (c *Checker) checkExpression(node ast.ExpressionNode) typed.ExpressionNode {
	switch node := node.(type) {
	case *ast.FalseLiteralNode:
		return typed.NewFalseLiteralNode(node.Span())
	case *ast.TrueLiteralNode:
		return typed.NewTrueLiteralNode(node.Span())
	case *ast.NilLiteralNode:
		return typed.NewNilLiteralNode(node.Span())
	case *ast.IntLiteralNode:
		return typed.NewIntLiteralNode(
			node.Span(),
			node.Value,
		)
	case *ast.Int64LiteralNode:
		return typed.NewInt64LiteralNode(
			node.Span(),
			node.Value,
		)
	case *ast.Int32LiteralNode:
		return typed.NewInt32LiteralNode(
			node.Span(),
			node.Value,
		)
	case *ast.Int16LiteralNode:
		return typed.NewInt16LiteralNode(
			node.Span(),
			node.Value,
		)
	case *ast.Int8LiteralNode:
		return typed.NewInt8LiteralNode(
			node.Span(),
			node.Value,
		)
	case *ast.UInt64LiteralNode:
		return typed.NewUInt64LiteralNode(
			node.Span(),
			node.Value,
		)
	case *ast.UInt32LiteralNode:
		return typed.NewUInt32LiteralNode(
			node.Span(),
			node.Value,
		)
	case *ast.UInt16LiteralNode:
		return typed.NewUInt16LiteralNode(
			node.Span(),
			node.Value,
		)
	case *ast.UInt8LiteralNode:
		return typed.NewUInt8LiteralNode(
			node.Span(),
			node.Value,
		)
	case *ast.FloatLiteralNode:
		return typed.NewFloatLiteralNode(
			node.Span(),
			node.Value,
		)
	case *ast.Float64LiteralNode:
		return typed.NewFloat64LiteralNode(
			node.Span(),
			node.Value,
		)
	case *ast.Float32LiteralNode:
		return typed.NewFloat32LiteralNode(
			node.Span(),
			node.Value,
		)
	case *ast.BigFloatLiteralNode:
		return typed.NewBigFloatLiteralNode(
			node.Span(),
			node.Value,
		)
	case *ast.DoubleQuotedStringLiteralNode:
		return typed.NewDoubleQuotedStringLiteralNode(
			node.Span(),
			node.Value,
		)
	case *ast.VariableDeclarationNode:
		return c.variableDeclaration(node)
	default:
		c.Errors.Add(
			fmt.Sprintf("invalid expression type %T", node),
			c.newLocation(node.Span()),
		)
		return typed.NewInvalidNode(node.Span(), nil)
	}
}

func (c *Checker) resolveConstant(name string, span *position.Span) types.Type {
	for i := range len(c.constantContainers) {
		constContainer := c.constantContainers[len(c.constantContainers)-i-1]
		constant := constContainer.Constant(name)
		if constant != nil {
			return constant
		}
	}

	c.Errors.Add(
		fmt.Sprintf("undefined constant `%s`", name),
		c.newLocation(span),
	)
	return types.Void{}
}

func (c *Checker) addLocal(name string, l local) {
	env := c.currentLocalEnv()
	env.locals[value.ToSymbol(name)] = l
}

func (c *Checker) checkTypeNode(node ast.TypeNode) typed.TypeNode {
	switch n := node.(type) {
	case *ast.PublicConstantNode:
		typ := c.resolveConstant(n.Value, n.Span())
		return typed.NewPublicConstantNode(
			n.Span(),
			n.Value,
			typ,
		)
	case *ast.PrivateConstantNode:
		typ := c.resolveConstant(n.Value, n.Span())
		return typed.NewPrivateConstantNode(
			n.Span(),
			n.Value,
			typ,
		)
	default:
		c.Errors.Add(
			fmt.Sprintf("invalid type node %T", node),
			c.newLocation(node.Span()),
		)
		return typed.NewInvalidNode(node.Span(), nil)
	}
}

func (c *Checker) variableDeclaration(node *ast.VariableDeclarationNode) *typed.VariableDeclarationNode {
	if node.Initialiser == nil {
		if node.Type == nil {
			c.Errors.Add(
				fmt.Sprintf("cannot define a variable without a type `%s`", node.Name.Value),
				c.newLocation(node.Span()),
			)
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
		c.addLocal(node.Name.Value, local{typ: actualType})
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
	if !declaredType.IsSupertypeOf(actualType) {
		c.Errors.Add(
			fmt.Sprintf("type `%s` cannot be assigned to type `%s`", actualType.Inspect(), declaredType.Inspect()),
			c.newLocation(init.Span()),
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
