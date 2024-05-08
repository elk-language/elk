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

// Holds the state of the type checking process
type Checker struct {
	Location           *position.Location
	Errors             errors.ErrorList
	GlobalEnv          *types.GlobalEnvironment
	ConstantContainers []types.ConstantContainer
}

// Instantiate a new Checker instance.
func new(loc *position.Location, globalEnv *types.GlobalEnvironment) *Checker {
	if globalEnv == nil {
		globalEnv = types.NewGlobalEnvironment()
	}
	return &Checker{
		Location:  loc,
		GlobalEnv: globalEnv,
		ConstantContainers: []types.ConstantContainer{
			globalEnv.Root,
			globalEnv.Std(),
		},
	}
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

func (c *Checker) constructType(node ast.TypeNode) types.Type {
	switch n := node.(type) {
	case *ast.PublicConstantNode:
		for i := range len(c.ConstantContainers) {
			constContainer := c.ConstantContainers[len(c.ConstantContainers)-i-1]
			constant := constContainer.Constant(n.Value)
			if constant != nil {
				return constant
			}
		}
		c.Errors.Add(
			fmt.Sprintf("undefined constant %s", n.Value),
			c.newLocation(n.Span()),
		)
	case *ast.PrivateConstantNode:
		for i := range len(c.ConstantContainers) {
			constContainer := c.ConstantContainers[len(c.ConstantContainers)-i]
			constant := constContainer.Constant(n.Value)
			if constant != nil {
				return constant
			}
		}
		c.Errors.Add(
			fmt.Sprintf("undefined constant %s", n.Value),
			c.newLocation(n.Span()),
		)
	default:
		c.Errors.Add(
			fmt.Sprintf("invalid type node %T", node),
			c.newLocation(node.Span()),
		)
	}

	return types.Void{}
}

func (c *Checker) variableDeclaration(node *ast.VariableDeclarationNode) *typed.VariableDeclarationNode {
	if node.Type != nil {
		expectedType := c.constructType(node.Type)
		if node.Initialiser != nil {
			init := c.checkExpression(node.Initialiser)
			actualType := typed.TypeOf(init, c.GlobalEnv)
			if !expectedType.IsSupertypeOf(actualType) {
				c.Errors.Add(
					fmt.Sprintf("type %s cannot be assigned to type %s", actualType.Inspect(), expectedType.Inspect()),
					c.newLocation(node.Span()),
				)
			}
		}
	}

	return nil
}
