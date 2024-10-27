// Package checker implements the Elk type checker
package checker

import (
	"fmt"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
)

type assumption uint8

func (a assumption) negate() assumption {
	switch a {
	case assumptionTruthy:
		return assumptionFalsy
	case assumptionFalsy:
		return assumptionTruthy
	case assumptionNil:
		return assumptionNotNil
	case assumptionNotNil:
		return assumptionNil
	case assumptionNever:
		return a
	default:
		panic(fmt.Sprintf("invalid assumption: %#v", a))
	}
}

func (a assumption) toNilable() assumption {
	switch a {
	case assumptionTruthy:
		return assumptionNotNil
	case assumptionFalsy:
		return assumptionNil
	default:
		return a
	}
}

const (
	assumptionTruthy assumption = iota
	assumptionFalsy
	assumptionNil
	assumptionNotNil
	assumptionNever
)

func (c *Checker) narrowToType(node ast.ExpressionNode, typ types.Type) {
	switch n := node.(type) {
	case *ast.PublicIdentifierNode:
		c.narrowLocalToType(n.Value, c.typeOf(n), typ)
	case *ast.PrivateIdentifierNode:
		c.narrowLocalToType(n.Value, c.typeOf(n), typ)
	case *ast.VariableDeclarationNode:
		c.narrowLocalToType(n.Name, c.typeOf(n), typ)
	case *ast.ValueDeclarationNode:
		c.narrowLocalToType(n.Name, c.typeOf(n), typ)
	case *ast.AssignmentExpressionNode:
		c.narrowAssignmentToType(n, typ)
	}
}

func (c *Checker) narrowAssignmentToType(node *ast.AssignmentExpressionNode, typ types.Type) {
	switch node.Op.Type {
	case token.EQUAL_OP, token.COLON_EQUAL:
		nodeType := c.typeOf(node)
		switch l := node.Left.(type) {
		case *ast.PublicIdentifierNode:
			c.narrowLocalToType(l.Value, nodeType, typ)
		case *ast.PrivateIdentifierNode:
			c.narrowLocalToType(l.Value, nodeType, typ)
		}
		return
	}

	c.narrowToType(node.Left, typ)
}

func (c *Checker) narrowLocalToType(name string, localType, typ types.Type) types.Type {
	local, inCurrentEnv := c.resolveLocal(name, nil)
	if local == nil {
		return types.Untyped{}
	}

	if !inCurrentEnv {
		local = local.copy()
		c.addLocal(name, local)
		local.shadow = true
	}
	narrowedType := c.newNormalisedIntersection(localType, typ)
	local.typ = narrowedType
	return narrowedType
}

func (c *Checker) narrowCondition(node ast.ExpressionNode, assume assumption) {
	switch n := node.(type) {
	case *ast.UnaryExpressionNode:
		c.narrowUnary(n, assume)
	case *ast.BinaryExpressionNode:
		c.narrowBinary(n, assume)
	case *ast.LogicalExpressionNode:
		c.narrowLogical(n, assume)
	case *ast.PublicIdentifierNode:
		c.narrowLocal(n.Value, assume)
	case *ast.PrivateIdentifierNode:
		c.narrowLocal(n.Value, assume)
	}
}

func (c *Checker) narrowLogical(node *ast.LogicalExpressionNode, assume assumption) {
	switch node.Op.Type {
	case token.AND_AND:
		c.narrowLogicalAnd(node, assume)
	case token.OR_OR:
		c.narrowLogicalOr(node, assume)
	case token.QUESTION_QUESTION:
		c.narrowNilCoalescing(node, assume)
	}
}

func (c *Checker) narrowLogicalAnd(node *ast.LogicalExpressionNode, assume assumption) {
	leftType := c.typeOf(node.Left)
	rightType := c.typeOf(node.Right)

	switch assume {
	case assumptionTruthy:
		// the whole condition is truthy, so the entire expression must be truthy
		if c.isFalsy(leftType) || c.isFalsy(rightType) {
			// any side is falsy, so the condition is impossible
			c.narrowCondition(node.Left, assumptionNever)
			c.narrowCondition(node.Right, assumptionNever)
			return
		}
		c.narrowCondition(node.Left, assumptionTruthy)
		c.narrowCondition(node.Right, assumptionTruthy)
		return
	case assumptionNotNil:
		// the whole condition is not nil, so the entire expression must be not nil
		if c.isNil(leftType) || c.isNil(rightType) {
			// any side is nil, so the condition is impossible
			c.narrowCondition(node.Left, assumptionNever)
			c.narrowCondition(node.Right, assumptionNever)
			return
		}
		c.narrowCondition(node.Left, assumptionNotNil)
		c.narrowCondition(node.Right, assumptionNotNil)
		return
	case assumptionNever:
		c.narrowCondition(node.Left, assumptionNever)
		c.narrowCondition(node.Right, assumptionNever)
		return
	case assumptionFalsy:
		// the whole condition is falsy
		if c.isTruthy(leftType) {
			// left is truthy, so right must be falsy
			c.narrowCondition(node.Right, assumptionFalsy)
			return
		}
		if c.isTruthy(rightType) {
			// left is falsy, right is truthy
			c.narrowCondition(node.Left, assumptionFalsy)
		}
	}

}

func (c *Checker) narrowLogicalOr(node *ast.LogicalExpressionNode, assume assumption) {
	leftType := c.typeOf(node.Left)
	rightType := c.typeOf(node.Right)

	switch assume {
	case assumptionFalsy:
		// the whole condition is falsy, so the entire expression must be falsy
		if c.isTruthy(leftType) || c.isTruthy(rightType) {
			// any side is truthy, so the condition is impossible
			c.narrowCondition(node.Left, assumptionNever)
			c.narrowCondition(node.Right, assumptionNever)
			return
		}
		c.narrowCondition(node.Left, assumptionFalsy)
		c.narrowCondition(node.Right, assumptionFalsy)
		return
	case assumptionNil:
		// the whole condition is nil, so the entire expression must be nil
		if c.isNotNilable(leftType) || c.isNotNilable(rightType) {
			// any side is not nilable, so the condition is impossible
			c.narrowCondition(node.Left, assumptionNever)
			c.narrowCondition(node.Right, assumptionNever)
			return
		}
		c.narrowCondition(node.Left, assumptionNil)
		c.narrowCondition(node.Right, assumptionNil)
		return
	case assumptionNever:
		c.narrowCondition(node.Left, assumptionNever)
		c.narrowCondition(node.Right, assumptionNever)
		return
	case assumptionTruthy:
		// the whole condition is truthy
		if c.isFalsy(leftType) {
			// left is falsy, so right must be truthy
			c.narrowCondition(node.Right, assumptionTruthy)
			return
		}
		if c.isFalsy(rightType) {
			// right is falsy, left be truthy
			c.narrowCondition(node.Left, assumptionTruthy)
		}
	case assumptionNotNil:
		// the whole condition is not nil
		if c.isFalsy(leftType) {
			// left is falsy, so right must not be nil
			c.narrowCondition(node.Right, assumptionNotNil)
			return
		}
	}
}

func (c *Checker) narrowNilCoalescing(node *ast.LogicalExpressionNode, assume assumption) {
	leftType := c.typeOf(node.Left)
	rightType := c.typeOf(node.Right)

	switch assume {
	case assumptionNil:
		// the whole condition is nil, so the entire expression must be nil
		if c.isNotNilable(leftType) || c.isNotNilable(rightType) {
			// any side is not nilable, so the condition is impossible
			c.narrowCondition(node.Left, assumptionNever)
			c.narrowCondition(node.Right, assumptionNever)
			return
		}
		c.narrowCondition(node.Left, assumptionNil)
		c.narrowCondition(node.Right, assumptionNil)
		return
	case assumptionNever:
		c.narrowCondition(node.Left, assumptionNever)
		c.narrowCondition(node.Right, assumptionNever)
		return
	case assumptionNotNil:
		// the whole condition is not nil
		if c.isNil(leftType) {
			// left is nil, so right must not be nil
			c.narrowCondition(node.Right, assumptionNotNil)
			return
		}
		if c.isNil(rightType) {
			// right is nil, left must not be nil
			c.narrowCondition(node.Left, assumptionNotNil)
		}
	}
}

func (c *Checker) narrowBinary(node *ast.BinaryExpressionNode, assume assumption) {
	switch node.Op.Type {
	case token.INSTANCE_OF_OP:
		c.narrowInstanceOf(node.Left, node.Right, assume)
	case token.REVERSE_INSTANCE_OF_OP:
		c.narrowInstanceOf(node.Right, node.Left, assume)
	case token.ISA_OP:
		c.narrowIsA(node.Left, node.Right, assume)
	case token.REVERSE_ISA_OP:
		c.narrowIsA(node.Right, node.Left, assume)
	case token.STRICT_EQUAL, token.EQUAL_EQUAL:
		c.narrowEqual(node, assume)
	case token.STRICT_NOT_EQUAL, token.NOT_EQUAL:
		c.narrowEqual(node, assume.negate())
	}
}

func (c *Checker) narrowEqual(node *ast.BinaryExpressionNode, assume assumption) {
	switch assume {
	case assumptionTruthy:
		c.narrowToIntersectWith(node.Left, c.typeOf(node.Right))
		c.narrowToIntersectWith(node.Right, c.typeOf(node.Left))
	case assumptionNil:
		c.narrowCondition(node.Left, assumptionNever)
		c.narrowCondition(node.Right, assumptionNever)
	}
}

func (c *Checker) narrowToIntersectWith(node ast.ExpressionNode, typ types.Type) {
	var localName string
	switch n := node.(type) {
	case *ast.PublicIdentifierNode:
		localName = n.Value
	case *ast.PrivateIdentifierNode:
		localName = n.Value
	default:
		return
	}

	local, inCurrentEnv := c.resolveLocal(localName, nil)
	if local == nil {
		return
	}

	if !inCurrentEnv {
		local = local.copy()
		local.shadow = true
		c.addLocal(localName, local)
	}
	local.typ = c.newNormalisedIntersection(local.typ, typ)
}

func (c *Checker) narrowIsA(left, right ast.ExpressionNode, assume assumption) {
	var localName string
	switch l := left.(type) {
	case *ast.PublicIdentifierNode:
		localName = l.Value
	case *ast.PrivateIdentifierNode:
		localName = l.Value
	default:
		return
	}

	rightSingleton, ok := c.typeOf(right).(*types.SingletonClass)
	if !ok {
		return
	}
	namespace := rightSingleton.AttachedObject

	local, inCurrentEnv := c.resolveLocal(localName, nil)
	if local == nil {
		return
	}

	if !inCurrentEnv {
		local = local.copy()
		local.shadow = true
		c.addLocal(localName, local)
	}
	switch assume {
	case assumptionTruthy:
		local.typ = namespace
	case assumptionFalsy:
		local.typ = c.differenceType(local.typ, namespace)
	case assumptionNotNil:
	case assumptionNever, assumptionNil:
		local.typ = types.Never{}
	}
}

func (c *Checker) differenceType(a, b types.Type) types.Type {
	return c.newNormalisedIntersection(a, types.NewNot(b))
}

func (c *Checker) narrowInstanceOf(left, right ast.ExpressionNode, assume assumption) {
	var localName string
	switch l := left.(type) {
	case *ast.PublicIdentifierNode:
		localName = l.Value
	case *ast.PrivateIdentifierNode:
		localName = l.Value
	default:
		return
	}

	rightSingleton, ok := c.typeOf(right).(*types.SingletonClass)
	if !ok {
		return
	}
	class, ok := rightSingleton.AttachedObject.(*types.Class)
	if !ok {
		return
	}

	local, inCurrentEnv := c.resolveLocal(localName, nil)
	if local == nil {
		return
	}

	if !inCurrentEnv {
		local = local.copy()
		local.shadow = true
		c.addLocal(localName, local)
	}
	switch assume {
	case assumptionTruthy:
		local.typ = class
	case assumptionFalsy:
		local.typ = c.differenceType(local.typ, class)
	case assumptionNotNil:
	case assumptionNever, assumptionNil:
		local.typ = types.Never{}
	}
}

func (c *Checker) narrowUnary(node *ast.UnaryExpressionNode, assume assumption) {
	switch node.Op.Type {
	case token.BANG:
		c.narrowCondition(node.Right, assume.negate())
	}
}

func (c *Checker) narrowLocal(name string, assume assumption) {
	local, inCurrentEnv := c.resolveLocal(name, nil)
	if local == nil {
		return
	}

	if !inCurrentEnv {
		local = local.copy()
		c.addLocal(name, local)
		local.shadow = true
	}
	switch assume {
	case assumptionTruthy:
		local.typ = c.toNonFalsy(local.typ)
	case assumptionFalsy:
		local.typ = c.toNonTruthy(local.typ)
	case assumptionNever:
		local.typ = types.Never{}
	case assumptionNil:
		local.typ = types.Nil{}
	case assumptionNotNil:
		local.typ = c.toNonNilable(local.typ)
	}
}

func (c *Checker) toNonNilable(typ types.Type) types.Type {
	return c.differenceType(typ, types.Nil{})
}

func (c *Checker) toNonFalsy(typ types.Type) types.Type {
	return c.newNormalisedIntersection(typ, types.NewNot(types.Nil{}), types.NewNot(types.False{}))
}

func (c *Checker) toNonTruthy(typ types.Type) types.Type {
	return c.newNormalisedIntersection(typ, types.NewUnion(types.Nil{}, types.False{}))
}

func (c *Checker) toNonLiteral(typ types.Type, widenSingletonTypes bool) types.Type {
	if !widenSingletonTypes {
		switch t := typ.(type) {
		case types.Nil, types.Bool:
			return typ
		case types.False, types.True:
			return types.Bool{}
		case *types.Union:
			newElements := make([]types.Type, len(t.Elements))
			for i, element := range t.Elements {
				newElements[i] = c.toNonLiteral(element, widenSingletonTypes)
			}
			return c.newNormalisedUnion(newElements...)
		}
	}

	return typ.ToNonLiteral(c.GlobalEnv)
}

func (c *Checker) toNilable(typ types.Type) types.Type {
	return c.normaliseType(types.NewNilable(typ))
}
