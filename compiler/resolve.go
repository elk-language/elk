package compiler

import (
	"slices"
	"strconv"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
)

// Create Elk runtime values from static AST nodes.
// Returns nil when no value could be created.
func resolve(node ast.ExpressionNode) value.Value {
	if !node.IsStatic() {
		return nil
	}

	switch n := node.(type) {
	case *ast.LabeledExpressionNode:
		return resolve(n.Expression)
	case *ast.TupleLiteralNode:
		return resolveTupleLiteral(n)
	case *ast.LogicalExpressionNode:
		return resolveLogicalExpression(n)
	case *ast.BinaryExpressionNode:
		return resolveBinaryExpression(n)
	case *ast.UnaryExpressionNode:
		return resolveUnaryExpression(n)
	case *ast.SimpleSymbolLiteralNode:
		return value.ToSymbol(n.Content)
	case *ast.RawStringLiteralNode:
		return value.String(n.Value)
	case *ast.DoubleQuotedStringLiteralNode:
		return value.String(n.Value)
	case *ast.RawCharLiteralNode:
		return value.Char(n.Value)
	case *ast.CharLiteralNode:
		return value.Char(n.Value)
	case *ast.NilLiteralNode:
		return value.Nil
	case *ast.TrueLiteralNode:
		return value.True
	case *ast.FalseLiteralNode:
		return value.False
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

	return nil
}

func resolveTupleLiteral(node *ast.TupleLiteralNode) value.Value {
	if !node.IsStatic() {
		return nil
	}

	newTuple := make(value.Tuple, 0, len(node.Elements))
	for _, elementNode := range node.Elements {
		switch e := elementNode.(type) {
		case *ast.KeyValueExpressionNode:
			key := resolve(e.Key)
			if key == nil {
				return nil
			}

			index, ok := value.ToGoInt(key)
			if !ok {
				return nil
			}

			val := resolve(e.Value)
			if val == nil {
				return nil
			}

			if index >= len(newTuple) {
				newElementsCount := (index + 1) - len(newTuple)
				newTuple = slices.Grow(newTuple, newElementsCount)
				for i := 0; i < newElementsCount; i++ {
					newTuple = append(newTuple, value.Nil)
				}
			}
			newTuple[index] = val
		default:
			element := resolve(elementNode)
			if element == nil {
				return nil
			}

			newTuple = append(newTuple, element)
		}
	}

	return &newTuple
}

func resolveLogicalExpression(node *ast.LogicalExpressionNode) value.Value {
	left := resolve(node.Left)
	if left == nil {
		return nil
	}
	right := resolve(node.Right)
	if right == nil {
		return nil
	}

	switch node.Op.Type {
	case token.AND_AND:
		if value.Truthy(left) {
			return right
		}
		return left
	case token.OR_OR:
		if value.Falsy(left) {
			return right
		}
		return left
	case token.QUESTION_QUESTION:
		if left == value.Nil {
			return right
		}
		return left
	}

	return nil
}

func resolveUnaryExpression(node *ast.UnaryExpressionNode) value.Value {
	right := resolve(node.Right)
	if right == nil {
		return nil
	}

	switch node.Op.Type {
	case token.PLUS:
		return right
	case token.MINUS:
		result := value.Negate(right)
		if result == nil {
			return nil
		}
		return result
	case token.BANG:
		return value.ToNotBool(right)
	case token.AND:
		singleton := right.SingletonClass()
		if singleton == nil {
			return nil
		}
		return singleton
	default:
		return nil
	}
}

func resolveBinaryExpression(node *ast.BinaryExpressionNode) value.Value {
	left := resolve(node.Left)
	if left == nil {
		return nil
	}
	right := resolve(node.Right)
	if right == nil {
		return nil
	}

	var result value.Value
	var err *value.Error

	switch node.Op.Type {
	case token.PLUS:
		result, err = value.Add(left, right)
	case token.MINUS:
		result, err = value.Subtract(left, right)
	case token.STAR:
		result, err = value.Multiply(left, right)
	case token.SLASH:
		result, err = value.Divide(left, right)
	case token.STAR_STAR:
		result, err = value.Exponentiate(left, right)
	case token.PERCENT:
		result, err = value.Modulo(left, right)
	case token.RBITSHIFT:
		result, err = value.RightBitshift(left, right)
	case token.RTRIPLE_BITSHIFT:
		result, err = value.LogicalRightBitshift(left, right)
	case token.LBITSHIFT:
		result, err = value.LeftBitshift(left, right)
	case token.LTRIPLE_BITSHIFT:
		result, err = value.LogicalLeftBitshift(left, right)
	case token.AND:
		result, err = value.BitwiseAnd(left, right)
	case token.OR:
		result, err = value.BitwiseOr(left, right)
	case token.XOR:
		result, err = value.BitwiseXor(left, right)
	case token.EQUAL_EQUAL:
		result = value.Equal(left, right)
	case token.NOT_EQUAL:
		result = value.NotEqual(left, right)
	case token.STRICT_EQUAL:
		result = value.StrictEqual(left, right)
	case token.STRICT_NOT_EQUAL:
		result = value.StrictNotEqual(left, right)
	case token.GREATER:
		result, err = value.GreaterThan(left, right)
	case token.GREATER_EQUAL:
		result, err = value.GreaterThanEqual(left, right)
	case token.LESS:
		result, err = value.LessThan(left, right)
	case token.LESS_EQUAL:
		result, err = value.LessThanEqual(left, right)
	default:
		return nil
	}

	if err != nil {
		return nil
	}
	return result
}

func resolveInt(node *ast.IntLiteralNode) value.Value {
	i, err := value.ParseBigInt(node.Value, 0)
	if err != nil {
		return nil
	}
	if i.IsSmallInt() {
		return i.ToSmallInt()
	}

	return i
}

func resolveInt64(node *ast.Int64LiteralNode) value.Value {
	i, err := value.StrictParseInt(node.Value, 0, 64)
	if err != nil {
		return nil
	}

	return value.Int64(i)
}

func resolveInt32(node *ast.Int32LiteralNode) value.Value {
	i, err := value.StrictParseInt(node.Value, 0, 32)
	if err != nil {
		return nil
	}

	return value.Int32(i)
}

func resolveInt16(node *ast.Int16LiteralNode) value.Value {
	i, err := value.StrictParseInt(node.Value, 0, 16)
	if err != nil {
		return nil
	}

	return value.Int16(i)
}

func resolveInt8(node *ast.Int8LiteralNode) value.Value {
	i, err := value.StrictParseInt(node.Value, 0, 8)
	if err != nil {
		return nil
	}

	return value.Int8(i)
}

func resolveUInt64(node *ast.UInt64LiteralNode) value.Value {
	i, err := value.StrictParseUint(node.Value, 0, 64)
	if err != nil {
		return nil
	}

	return value.UInt64(i)
}

func resolveUInt32(node *ast.UInt32LiteralNode) value.Value {
	i, err := value.StrictParseUint(node.Value, 0, 32)
	if err != nil {
		return nil
	}

	return value.UInt32(i)
}

func resolveUInt16(node *ast.UInt16LiteralNode) value.Value {
	i, err := value.StrictParseUint(node.Value, 0, 16)
	if err != nil {
		return nil
	}

	return value.UInt16(i)
}

func resolveUInt8(node *ast.UInt8LiteralNode) value.Value {
	i, err := value.StrictParseUint(node.Value, 0, 8)
	if err != nil {
		return nil
	}

	return value.UInt8(i)
}

func resolveBigFloat(node *ast.BigFloatLiteralNode) value.Value {
	f, err := value.ParseBigFloat(node.Value)
	if err != nil {
		return nil
	}

	return f
}

func resolveFloat64(node *ast.Float64LiteralNode) value.Value {
	f, err := strconv.ParseFloat(node.Value, 64)
	if err != nil {
		return nil
	}

	return value.Float64(f)
}

func resolveFloat32(node *ast.Float32LiteralNode) value.Value {
	f, err := strconv.ParseFloat(node.Value, 32)
	if err != nil {
		return nil
	}

	return value.Float32(f)
}

func resolveFloat(node *ast.FloatLiteralNode) value.Value {
	f, err := strconv.ParseFloat(node.Value, 64)
	if err != nil {
		return nil
	}

	return value.Float(f)
}
