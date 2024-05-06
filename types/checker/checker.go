// Package checker implements the Elk type checker
package checker

import (
	"fmt"

	"github.com/elk-language/elk/parser"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	typed "github.com/elk-language/elk/types/ast" // typed AST
)

// Check the types of Elk source code.
func CheckSource(sourceName string, source string) (typed.Node, errors.ErrorList) {
	ast, err := parser.Parse(sourceName, source)
	if err != nil {
		return nil, err
	}

	return CheckAST(sourceName, ast)
}

// Check the types of an Elk AST.
func CheckAST(sourceName string, ast *ast.ProgramNode) (typed.Node, errors.ErrorList) {
	checker := new(position.NewLocationWithSpan(sourceName, ast.Span()))
	typedAst := checker.checkProgram(ast)
	return typedAst, checker.Errors
}

// Holds the state of the type checking process
type Checker struct {
	Location *position.Location
	Errors   errors.ErrorList
}

// Instantiate a new Checker instance.
func new(loc *position.Location) *Checker {
	return &Checker{
		Location: loc,
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

func (c *Checker) checkExpression(node ast.Node) typed.ExpressionNode {
	switch node := node.(type) {
	case *ast.ProgramNode:
		var newStatements []typed.StatementNode
		for _, statement := range node.Body {
			newStatements = append(newStatements, c.checkStatement(statement))
		}

	}
}
