package compiler

import (
	"strconv"

	"github.com/elk-language/elk/object"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
)

// Create Elk runtime objects from static AST nodes.
func resolve(node ast.ExpressionNode) (object.Value, bool) {
	if !node.IsStatic() {
		return nil, false
	}

	switch n := node.(type) {
	case *ast.LogicalExpressionNode:
		return resolveLogicalExpression(n)
	case *ast.BinaryExpressionNode:
		return resolveBinaryExpression(n)
	case *ast.UnaryExpressionNode:
		return resolveUnaryExpression(n)
	case *ast.SimpleSymbolLiteralNode:
		return object.SymbolTable.Add(n.Content), true
	case *ast.RawStringLiteralNode:
		return object.String(n.Value), true
	case *ast.DoubleQuotedStringLiteralNode:
		return object.String(n.Value), true
	case *ast.RawCharLiteralNode:
		return object.Char(n.Value), true
	case *ast.CharLiteralNode:
		return object.Char(n.Value), true
	case *ast.NilLiteralNode:
		return object.Nil, true
	case *ast.TrueLiteralNode:
		return object.True, true
	case *ast.FalseLiteralNode:
		return object.False, true
	case *ast.IntLiteralNode:
		return resolveInt(n)
	case *ast.Int64LiteralNode:
		return resolveInt64(n)
	case *ast.Int32LiteralNode:
		return resolveInt32(n)
	case *ast.Int16LiteralNode:
		return resolveInt16(n)
	case *ast.Int8LiteralNode:
		return resolveInt8(n)
	case *ast.UInt64LiteralNode:
		return resolveUInt64(n)
	case *ast.UInt32LiteralNode:
		return resolveUInt32(n)
	case *ast.UInt16LiteralNode:
		return resolveUInt16(n)
	case *ast.UInt8LiteralNode:
		return resolveUInt8(n)
	case *ast.BigFloatLiteralNode:
		return resolveBigFloat(n)
	case *ast.Float64LiteralNode:
		return resolveFloat64(n)
	case *ast.Float32LiteralNode:
		return resolveFloat32(n)
	case *ast.FloatLiteralNode:
		return resolveFloat(n)
	}

	return nil, false
}

func resolveLogicalExpression(node *ast.LogicalExpressionNode) (object.Value, bool) {
	left, ok := resolve(node.Left)
	if !ok {
		return nil, false
	}
	right, ok := resolve(node.Right)
	if !ok {
		return nil, false
	}

	switch node.Op.Type {
	case token.AND_AND:
		if object.Truthy(left) {
			return right, true
		}
		return left, true
	case token.OR_OR:
		if object.Falsy(left) {
			return right, true
		}
		return left, true
	case token.QUESTION_QUESTION:
		if left == object.Nil {
			return right, true
		}
		return left, true
	}

	return nil, false
}

func resolveUnaryExpression(node *ast.UnaryExpressionNode) (object.Value, bool) {
	right, ok := resolve(node.Right)
	if !ok {
		return nil, false
	}

	switch node.Op.Type {
	case token.PLUS:
		return right, true
	case token.MINUS:
		return object.Negate(right)
	case token.BANG:
		return object.ToNotBool(right), true
	default:
		return nil, false
	}
}

func resolveBinaryExpression(node *ast.BinaryExpressionNode) (object.Value, bool) {
	left, ok := resolve(node.Left)
	if !ok {
		return nil, false
	}
	right, ok := resolve(node.Right)
	if !ok {
		return nil, false
	}

	var result object.Value
	var err *object.Error

	switch node.Op.Type {
	case token.PLUS:
		result, err, ok = object.Add(left, right)
	case token.MINUS:
		result, err, ok = object.Subtract(left, right)
	case token.STAR:
		result, err, ok = object.Multiply(left, right)
	case token.SLASH:
		result, err, ok = object.Divide(left, right)
	default:
		return nil, false
	}

	if err != nil || !ok {
		return nil, false
	}
	return result, true
}

func resolveInt(node *ast.IntLiteralNode) (object.Value, bool) {
	i, err := object.ParseBigInt(node.Value, 0)
	if err != nil {
		return nil, false
	}
	if i.IsSmallInt() {
		return i.ToSmallInt(), true
	}

	return i, true
}

func resolveInt64(node *ast.Int64LiteralNode) (object.Value, bool) {
	i, err := object.StrictParseInt(node.Value, 0, 64)
	if err != nil {
		return nil, false
	}

	return object.Int64(i), true
}

func resolveInt32(node *ast.Int32LiteralNode) (object.Value, bool) {
	i, err := object.StrictParseInt(node.Value, 0, 32)
	if err != nil {
		return nil, false
	}

	return object.Int32(i), true
}

func resolveInt16(node *ast.Int16LiteralNode) (object.Value, bool) {
	i, err := object.StrictParseInt(node.Value, 0, 16)
	if err != nil {
		return nil, false
	}

	return object.Int16(i), true
}

func resolveInt8(node *ast.Int8LiteralNode) (object.Value, bool) {
	i, err := object.StrictParseInt(node.Value, 0, 8)
	if err != nil {
		return nil, false
	}

	return object.Int8(i), true
}

func resolveUInt64(node *ast.UInt64LiteralNode) (object.Value, bool) {
	i, err := object.StrictParseUint(node.Value, 0, 64)
	if err != nil {
		return nil, false
	}

	return object.UInt64(i), true
}

func resolveUInt32(node *ast.UInt32LiteralNode) (object.Value, bool) {
	i, err := object.StrictParseUint(node.Value, 0, 32)
	if err != nil {
		return nil, false
	}

	return object.UInt32(i), true
}

func resolveUInt16(node *ast.UInt16LiteralNode) (object.Value, bool) {
	i, err := object.StrictParseUint(node.Value, 0, 16)
	if err != nil {
		return nil, false
	}

	return object.UInt16(i), true
}

func resolveUInt8(node *ast.UInt8LiteralNode) (object.Value, bool) {
	i, err := object.StrictParseUint(node.Value, 0, 8)
	if err != nil {
		return nil, false
	}

	return object.UInt8(i), true
}

func resolveBigFloat(node *ast.BigFloatLiteralNode) (object.Value, bool) {
	f, err := object.ParseBigFloat(node.Value)
	if err != nil {
		return nil, false
	}

	return f, true
}

func resolveFloat64(node *ast.Float64LiteralNode) (object.Value, bool) {
	f, err := strconv.ParseFloat(node.Value, 64)
	if err != nil {
		return nil, false
	}

	return object.Float64(f), true
}

func resolveFloat32(node *ast.Float32LiteralNode) (object.Value, bool) {
	f, err := strconv.ParseFloat(node.Value, 32)
	if err != nil {
		return nil, false
	}

	return object.Float32(f), true
}

func resolveFloat(node *ast.FloatLiteralNode) (object.Value, bool) {
	f, err := strconv.ParseFloat(node.Value, 64)
	if err != nil {
		return nil, false
	}

	return object.Float(f), true
}
