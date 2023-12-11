package compiler

import (
	"strconv"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Create Elk runtime values from static AST nodes.
func resolve(node ast.ExpressionNode) (value.Value, bool) {
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
		return value.ToSymbol(n.Content), true
	case *ast.RawStringLiteralNode:
		return value.String(n.Value), true
	case *ast.DoubleQuotedStringLiteralNode:
		return value.String(n.Value), true
	case *ast.RawCharLiteralNode:
		return value.Char(n.Value), true
	case *ast.CharLiteralNode:
		return value.Char(n.Value), true
	case *ast.NilLiteralNode:
		return value.Nil, true
	case *ast.TrueLiteralNode:
		return value.True, true
	case *ast.FalseLiteralNode:
		return value.False, true
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

func resolveLogicalExpression(node *ast.LogicalExpressionNode) (value.Value, bool) {
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
		if value.Truthy(left) {
			return right, true
		}
		return left, true
	case token.OR_OR:
		if value.Falsy(left) {
			return right, true
		}
		return left, true
	case token.QUESTION_QUESTION:
		if left == value.Nil {
			return right, true
		}
		return left, true
	}

	return nil, false
}

func resolveUnaryExpression(node *ast.UnaryExpressionNode) (value.Value, bool) {
	right, ok := resolve(node.Right)
	if !ok {
		return nil, false
	}

	switch node.Op.Type {
	case token.PLUS:
		return right, true
	case token.MINUS:
		return value.Negate(right)
	case token.BANG:
		return value.ToNotBool(right), true
	default:
		return nil, false
	}
}

func resolveBinaryExpression(node *ast.BinaryExpressionNode) (value.Value, bool) {
	left, ok := resolve(node.Left)
	if !ok {
		return nil, false
	}
	right, ok := resolve(node.Right)
	if !ok {
		return nil, false
	}

	var result value.Value
	var err *value.Error

	switch node.Op.Type {
	case token.PLUS:
		result, err, ok = value.Add(left, right)
	case token.MINUS:
		result, err, ok = value.Subtract(left, right)
	case token.STAR:
		result, err, ok = value.Multiply(left, right)
	case token.SLASH:
		result, err, ok = value.Divide(left, right)
	case token.STAR_STAR:
		result, err, ok = value.Exponentiate(left, right)
	case token.PERCENT:
		result, err, ok = value.Modulo(left, right)
	case token.RBITSHIFT:
		result, err, ok = value.RightBitshift(left, right)
	case token.RTRIPLE_BITSHIFT:
		result, err, ok = value.LogicalRightBitshift(left, right)
	case token.LBITSHIFT:
		result, err, ok = value.LeftBitshift(left, right)
	case token.LTRIPLE_BITSHIFT:
		result, err, ok = value.LogicalLeftBitshift(left, right)
	case token.AND:
		result, err, ok = value.BitwiseAnd(left, right)
	case token.OR:
		result, err, ok = value.BitwiseOr(left, right)
	case token.XOR:
		result, err, ok = value.BitwiseXor(left, right)
	case token.EQUAL_EQUAL:
		result, ok = value.Equal(left, right)
	case token.NOT_EQUAL:
		result, ok = value.NotEqual(left, right)
	case token.STRICT_EQUAL:
		result, ok = value.StrictEqual(left, right)
	case token.STRICT_NOT_EQUAL:
		result, ok = value.StrictNotEqual(left, right)
	case token.GREATER:
		result, err, ok = value.GreaterThan(left, right)
	case token.GREATER_EQUAL:
		result, err, ok = value.GreaterThanEqual(left, right)
	case token.LESS:
		result, err, ok = value.LessThan(left, right)
	case token.LESS_EQUAL:
		result, err, ok = value.LessThanEqual(left, right)
	default:
		return nil, false
	}

	if err != nil || !ok {
		return nil, false
	}
	return result, true
}

func resolveInt(node *ast.IntLiteralNode) (value.Value, bool) {
	i, err := value.ParseBigInt(node.Value, 0)
	if err != nil {
		return nil, false
	}
	if i.IsSmallInt() {
		return i.ToSmallInt(), true
	}

	return i, true
}

func resolveInt64(node *ast.Int64LiteralNode) (value.Value, bool) {
	i, err := value.StrictParseInt(node.Value, 0, 64)
	if err != nil {
		return nil, false
	}

	return value.Int64(i), true
}

func resolveInt32(node *ast.Int32LiteralNode) (value.Value, bool) {
	i, err := value.StrictParseInt(node.Value, 0, 32)
	if err != nil {
		return nil, false
	}

	return value.Int32(i), true
}

func resolveInt16(node *ast.Int16LiteralNode) (value.Value, bool) {
	i, err := value.StrictParseInt(node.Value, 0, 16)
	if err != nil {
		return nil, false
	}

	return value.Int16(i), true
}

func resolveInt8(node *ast.Int8LiteralNode) (value.Value, bool) {
	i, err := value.StrictParseInt(node.Value, 0, 8)
	if err != nil {
		return nil, false
	}

	return value.Int8(i), true
}

func resolveUInt64(node *ast.UInt64LiteralNode) (value.Value, bool) {
	i, err := value.StrictParseUint(node.Value, 0, 64)
	if err != nil {
		return nil, false
	}

	return value.UInt64(i), true
}

func resolveUInt32(node *ast.UInt32LiteralNode) (value.Value, bool) {
	i, err := value.StrictParseUint(node.Value, 0, 32)
	if err != nil {
		return nil, false
	}

	return value.UInt32(i), true
}

func resolveUInt16(node *ast.UInt16LiteralNode) (value.Value, bool) {
	i, err := value.StrictParseUint(node.Value, 0, 16)
	if err != nil {
		return nil, false
	}

	return value.UInt16(i), true
}

func resolveUInt8(node *ast.UInt8LiteralNode) (value.Value, bool) {
	i, err := value.StrictParseUint(node.Value, 0, 8)
	if err != nil {
		return nil, false
	}

	return value.UInt8(i), true
}

func resolveBigFloat(node *ast.BigFloatLiteralNode) (value.Value, bool) {
	f, err := value.ParseBigFloat(node.Value)
	if err != nil {
		return nil, false
	}

	return f, true
}

func resolveFloat64(node *ast.Float64LiteralNode) (value.Value, bool) {
	f, err := strconv.ParseFloat(node.Value, 64)
	if err != nil {
		return nil, false
	}

	return value.Float64(f), true
}

func resolveFloat32(node *ast.Float32LiteralNode) (value.Value, bool) {
	f, err := strconv.ParseFloat(node.Value, 32)
	if err != nil {
		return nil, false
	}

	return value.Float32(f), true
}

func resolveFloat(node *ast.FloatLiteralNode) (value.Value, bool) {
	f, err := strconv.ParseFloat(node.Value, 64)
	if err != nil {
		return nil, false
	}

	return value.Float(f), true
}
