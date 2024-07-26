// Package checker implements the Elk type checker
package checker

import (
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
)

func (c *Checker) narrowCondition(node ast.ExpressionNode, assumeTruthy bool) {
	switch n := node.(type) {
	case *ast.UnaryExpressionNode:
		c.narrowUnary(n, assumeTruthy)
	case *ast.BinaryExpressionNode:
		c.narrowBinary(n, assumeTruthy)
	case *ast.LogicalExpressionNode:
		c.narrowLogical(n, assumeTruthy)
	case *ast.PublicIdentifierNode:
		c.narrowLocal(n.Value, assumeTruthy)
	case *ast.PrivateIdentifierNode:
		c.narrowLocal(n.Value, assumeTruthy)
	}
}

func (c *Checker) narrowLogical(node *ast.LogicalExpressionNode, assumeTruthy bool) {
	switch node.Op.Type {
	case token.AND_AND:
		c.narrowLogicalAnd(node, assumeTruthy)
	case token.OR_OR:
		c.narrowLogicalOr(node, assumeTruthy)
		// case token.QUESTION_QUESTION:
		// 	c.narrowNilCoalescing(node, assumeTruthy)
	}
}

func (c *Checker) narrowLogicalAnd(node *ast.LogicalExpressionNode, assumeTruthy bool) {
	if assumeTruthy {
		// the whole condition is truthy, so the entire expression must be truthy
		c.narrowCondition(node.Left, true)
		c.narrowCondition(node.Right, true)
		return
	}

	// the whole condition is falsy
	leftType := c.typeOf(node.Left)
	rightType := c.typeOf(node.Right)
	if c.isTruthy(leftType) {
		// left is truthy, so right must be falsy
		c.narrowCondition(node.Right, false)
		return
	}
	if c.isTruthy(rightType) {
		// left is falsy, right is truthy
		c.narrowCondition(node.Left, false)
	}
}

func (c *Checker) narrowLogicalOr(node *ast.LogicalExpressionNode, assumeTruthy bool) {
	if !assumeTruthy {
		// the whole condition is falsy, so the entire expression must be falsy
		c.narrowCondition(node.Left, false)
		c.narrowCondition(node.Right, false)
		return
	}

	// the whole condition is truthy
	leftType := c.typeOf(node.Left)
	rightType := c.typeOf(node.Right)
	if c.isFalsy(leftType) {
		// left is falsy, so right must be truthy
		c.narrowCondition(node.Right, true)
		return
	}
	if c.isTruthy(rightType) {
		// left is truthy, right is falsy
		c.narrowCondition(node.Left, true)
	}
}

func (c *Checker) narrowBinary(node *ast.BinaryExpressionNode, assumeTruthy bool) {
	switch node.Op.Type {
	case token.INSTANCE_OF_OP:
		c.narrowInstanceOf(node, assumeTruthy)
	case token.ISA_OP:
		c.narrowIsA(node, assumeTruthy)
	}
}

func (c *Checker) narrowIsA(node *ast.BinaryExpressionNode, assumeTruthy bool) {
	var localName string
	switch l := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		localName = l.Value
	case *ast.PrivateIdentifierNode:
		localName = l.Value
	default:
		return
	}

	rightSingleton, ok := c.typeOf(node.Right).(*types.SingletonClass)
	if !ok {
		return
	}
	namespace := rightSingleton.AttachedObject

	local := c.resolveLocal(localName, nil)
	if local == nil {
		return
	}
	newLocal := local.copy()
	newLocal.shadow = true
	if assumeTruthy {
		newLocal.typ = namespace
	} else {
		newLocal.typ = c.newNormalisedIntersection(local.typ, types.NewNot(namespace))
	}
	c.addLocal(localName, newLocal)
}

func (c *Checker) narrowInstanceOf(node *ast.BinaryExpressionNode, assumeTruthy bool) {
	var localName string
	switch l := node.Left.(type) {
	case *ast.PublicIdentifierNode:
		localName = l.Value
	case *ast.PrivateIdentifierNode:
		localName = l.Value
	default:
		return
	}

	rightSingleton, ok := c.typeOf(node.Right).(*types.SingletonClass)
	if !ok {
		return
	}
	class, ok := rightSingleton.AttachedObject.(*types.Class)
	if !ok {
		return
	}

	local := c.resolveLocal(localName, nil)
	if local == nil {
		return
	}
	newLocal := local.copy()
	newLocal.shadow = true
	if assumeTruthy {
		newLocal.typ = class
	} else {
		newLocal.typ = c.newNormalisedIntersection(local.typ, types.NewNot(class))
	}
	c.addLocal(localName, newLocal)
}

func (c *Checker) narrowUnary(node *ast.UnaryExpressionNode, assumeTruthy bool) {
	switch node.Op.Type {
	case token.BANG:
		c.narrowCondition(node.Right, !assumeTruthy)
	}
}

func (c *Checker) narrowLocal(name string, assumeTruthy bool) {
	local := c.resolveLocal(name, nil)
	if local == nil {
		return
	}

	newLocal := local.copy()
	newLocal.shadow = true
	if assumeTruthy {
		newLocal.typ = c.toNonFalsy(local.typ)
	} else {
		newLocal.typ = c.toNonTruthy(local.typ)
	}
	c.addLocal(name, newLocal)
}

func (c *Checker) toNonNilable(typ types.Type) types.Type {
	return c.newNormalisedIntersection(typ, types.NewNot(types.Nil{}))
}

func (c *Checker) toNonFalsy(typ types.Type) types.Type {
	return c.newNormalisedIntersection(typ, types.NewNot(types.Nil{}), types.NewNot(types.False{}))
}

func (c *Checker) toNonTruthy(typ types.Type) types.Type {
	return c.newNormalisedIntersection(typ, types.NewUnion(types.Nil{}, types.False{}))
}

func (c *Checker) toNonLiteral(typ types.Type, widenSingletonTypes bool) types.Type {
	if !widenSingletonTypes {
		switch typ.(type) {
		case types.Nil, types.Bool:
			return typ
		case types.False, types.True:
			return types.Bool{}
		}
	}

	return typ.ToNonLiteral(c.GlobalEnv)
}

func (c *Checker) toNilable(typ types.Type) types.Type {
	return c.normaliseType(types.NewNilable(typ))
}
