package compiler

import (
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
	case *ast.ArrayTupleLiteralNode:
		return resolveArrayTupleLiteral(n)
	case *ast.WordArrayTupleLiteralNode:
		return resolveWordArrayTupleLiteral(n)
	case *ast.SymbolArrayTupleLiteralNode:
		return resolveSymbolArrayTupleLiteral(n)
	case *ast.BinArrayTupleLiteralNode:
		return resolveBinArrayTupleLiteral(n)
	case *ast.HexArrayTupleLiteralNode:
		return resolveHexArrayTupleLiteral(n)
	case *ast.LogicalExpressionNode:
		return resolveLogicalExpression(n)
	case *ast.BinaryExpressionNode:
		return resolveBinaryExpression(n)
	case *ast.UnaryExpressionNode:
		return resolveUnaryExpression(n)
	case *ast.SubscriptExpressionNode:
		return resolveSubscript(n)
	case *ast.NilSafeSubscriptExpressionNode:
		return resolveNilSafeSubscript(n)
	case *ast.ArrayListLiteralNode:
		return resolveArrayListLiteral(n)
	case *ast.WordArrayListLiteralNode:
		return resolveWordArrayListLiteral(n)
	case *ast.SymbolArrayListLiteralNode:
		return resolveSymbolArrayListLiteral(n)
	case *ast.BinArrayListLiteralNode:
		return resolveBinArrayListLiteral(n)
	case *ast.HexArrayListLiteralNode:
		return resolveHexArrayListLiteral(n)
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

func resolveArrayListLiteral(node *ast.ArrayListLiteralNode) value.Value {
	if !node.IsStatic() {
		return nil
	}

	newList := make(value.ArrayList, 0, len(node.Elements))
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

			if index >= len(newList) {
				newElementsCount := (index + 1) - len(newList)
				newList.Expand(newElementsCount)
			}
			newList[index] = val
		default:
			element := resolve(elementNode)
			if element == nil {
				return nil
			}

			newList = append(newList, element)
		}
	}

	return &newList
}

func resolveWordArrayListLiteral(node *ast.WordArrayListLiteralNode) value.Value {
	if !node.IsStatic() {
		return nil
	}

	newList := make(value.ArrayList, 0, len(node.Elements))
	for _, elementNode := range node.Elements {
		element := resolve(elementNode)
		if element == nil {
			return nil
		}
		newList = append(newList, element)
	}

	return &newList
}

func resolveHexArrayListLiteral(node *ast.HexArrayListLiteralNode) value.Value {
	if !node.IsStatic() {
		return nil
	}

	newList := make(value.ArrayList, 0, len(node.Elements))
	for _, elementNode := range node.Elements {
		element := resolve(elementNode)
		if element == nil {
			return nil
		}
		newList = append(newList, element)
	}

	return &newList
}

func resolveSymbolArrayListLiteral(node *ast.SymbolArrayListLiteralNode) value.Value {
	if !node.IsStatic() {
		return nil
	}

	newList := make(value.ArrayList, 0, len(node.Elements))
	for _, elementNode := range node.Elements {
		element := resolve(elementNode)
		if element == nil {
			return nil
		}
		newList = append(newList, element)
	}

	return &newList
}

func resolveBinArrayListLiteral(node *ast.BinArrayListLiteralNode) value.Value {
	if !node.IsStatic() {
		return nil
	}

	newList := make(value.ArrayList, 0, len(node.Elements))
	for _, elementNode := range node.Elements {
		element := resolve(elementNode)
		if element == nil {
			return nil
		}
		newList = append(newList, element)
	}

	return &newList
}

func resolveWordArrayTupleLiteral(node *ast.WordArrayTupleLiteralNode) value.Value {
	if !node.IsStatic() {
		return nil
	}

	newArrayTuple := make(value.ArrayTuple, 0, len(node.Elements))
	for _, elementNode := range node.Elements {
		element := resolve(elementNode)
		if element == nil {
			return nil
		}
		newArrayTuple = append(newArrayTuple, element)
	}

	return &newArrayTuple
}

func resolveHexArrayTupleLiteral(node *ast.HexArrayTupleLiteralNode) value.Value {
	if !node.IsStatic() {
		return nil
	}

	newArrayTuple := make(value.ArrayTuple, 0, len(node.Elements))
	for _, elementNode := range node.Elements {
		element := resolve(elementNode)
		if element == nil {
			return nil
		}
		newArrayTuple = append(newArrayTuple, element)
	}

	return &newArrayTuple
}

func resolveSymbolArrayTupleLiteral(node *ast.SymbolArrayTupleLiteralNode) value.Value {
	if !node.IsStatic() {
		return nil
	}

	newArrayTuple := make(value.ArrayTuple, 0, len(node.Elements))
	for _, elementNode := range node.Elements {
		element := resolve(elementNode)
		if element == nil {
			return nil
		}
		newArrayTuple = append(newArrayTuple, element)
	}

	return &newArrayTuple
}

func resolveBinArrayTupleLiteral(node *ast.BinArrayTupleLiteralNode) value.Value {
	if !node.IsStatic() {
		return nil
	}

	newArrayTuple := make(value.ArrayTuple, 0, len(node.Elements))
	for _, elementNode := range node.Elements {
		element := resolve(elementNode)
		if element == nil {
			return nil
		}
		newArrayTuple = append(newArrayTuple, element)
	}

	return &newArrayTuple
}

func resolveArrayTupleLiteral(node *ast.ArrayTupleLiteralNode) value.Value {
	if !node.IsStatic() {
		return nil
	}

	newArrayTuple := make(value.ArrayTuple, 0, len(node.Elements))
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

			if index >= len(newArrayTuple) {
				newElementsCount := (index + 1) - len(newArrayTuple)
				newArrayTuple.Expand(newElementsCount)
			}
			newArrayTuple[index] = val
		default:
			element := resolve(elementNode)
			if element == nil {
				return nil
			}

			newArrayTuple = append(newArrayTuple, element)
		}
	}

	return &newArrayTuple
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

func resolveNilSafeSubscript(node *ast.NilSafeSubscriptExpressionNode) value.Value {
	receiver := resolve(node.Receiver)
	key := resolve(node.Key)

	if receiver == value.Nil {
		return value.Nil
	}

	result, err := value.Subscript(receiver, key)
	if err != nil {
		return nil
	}

	return result
}

func resolveSubscript(node *ast.SubscriptExpressionNode) value.Value {
	receiver := resolve(node.Receiver)
	key := resolve(node.Key)

	result, err := value.Subscript(receiver, key)
	if err != nil {
		return nil
	}

	return result
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
