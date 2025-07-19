package compiler

import (
	"regexp"
	"strconv"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/regex"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
)

// Create Elk runtime values from static AST nodes.
// Returns undefined when no value could be created.
func resolve(node ast.Node) value.Value {
	if !node.IsStatic() {
		return value.Undefined
	}

	switch n := node.(type) {
	case *ast.LabeledExpressionNode:
		return resolve(n.Expression)
	case *ast.UninterpolatedRegexLiteralNode:
		return resolveUninterpolatedRegexLiteral(n)
	case *ast.RangeLiteralNode:
		return resolveRangeLiteral(n)
	case *ast.HashSetLiteralNode:
		return resolveHashSetLiteral(n)
	case *ast.WordHashSetLiteralNode:
		return resolveSpecialHashSetLiteral(n.Elements, n.IsStatic())
	case *ast.SymbolHashSetLiteralNode:
		return resolveSpecialHashSetLiteral(n.Elements, n.IsStatic())
	case *ast.BinHashSetLiteralNode:
		return resolveSpecialHashSetLiteral(n.Elements, n.IsStatic())
	case *ast.HexHashSetLiteralNode:
		return resolveSpecialHashSetLiteral(n.Elements, n.IsStatic())
	case *ast.HashMapLiteralNode:
		return resolveHashMapLiteral(n)
	case *ast.HashRecordLiteralNode:
		return resolveHashRecordLiteral(n)
	case *ast.ArrayListLiteralNode:
		return resolveArrayListLiteral(n)
	case *ast.WordArrayListLiteralNode:
		return resolveSpecialArrayListLiteral(n.Elements, n.IsStatic())
	case *ast.SymbolArrayListLiteralNode:
		return resolveSpecialArrayListLiteral(n.Elements, n.IsStatic())
	case *ast.BinArrayListLiteralNode:
		return resolveSpecialArrayListLiteral(n.Elements, n.IsStatic())
	case *ast.HexArrayListLiteralNode:
		return resolveSpecialArrayListLiteral(n.Elements, n.IsStatic())
	case *ast.ArrayTupleLiteralNode:
		return resolveArrayTupleLiteral(n)
	case *ast.WordArrayTupleLiteralNode:
		return resolveSpecialArrayTupleLiteral(n.Elements, n.IsStatic())
	case *ast.SymbolArrayTupleLiteralNode:
		return resolveSpecialArrayTupleLiteral(n.Elements, n.IsStatic())
	case *ast.BinArrayTupleLiteralNode:
		return resolveSpecialArrayTupleLiteral(n.Elements, n.IsStatic())
	case *ast.HexArrayTupleLiteralNode:
		return resolveSpecialArrayTupleLiteral(n.Elements, n.IsStatic())
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
	case *ast.SimpleSymbolLiteralNode:
		return value.ToSymbol(n.Content).ToValue()
	case *ast.RawStringLiteralNode:
		return value.Ref(value.String(n.Value))
	case *ast.DoubleQuotedStringLiteralNode:
		return value.Ref(value.String(n.Value))
	case *ast.RawCharLiteralNode:
		return value.Char(n.Value).ToValue()
	case *ast.CharLiteralNode:
		return value.Char(n.Value).ToValue()
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
	case *ast.MacroBoundaryNode:
		return resolveMacroBoundary(n)
	}

	return value.Undefined
}

func resolveMacroBoundary(n *ast.MacroBoundaryNode) value.Value {
	if len(n.Body) != 1 {
		return value.Undefined
	}

	stmt := n.Body[0]
	exprStmt, ok := stmt.(*ast.ExpressionStatementNode)
	if !ok {
		return value.Undefined
	}

	return resolve(exprStmt.Expression)
}

func resolveUninterpolatedRegexLiteral(node *ast.UninterpolatedRegexLiteralNode) value.Value {
	goRegexString, errList := regex.Transpile(node.Content, node.Flags)
	if errList != nil {
		return value.Undefined
	}

	re, err := regexp.Compile(goRegexString)
	if err != nil {
		return value.Undefined
	}

	return value.Ref(value.NewRegex(*re, node.Content, node.Flags))
}

func resolveRangeLiteral(node *ast.RangeLiteralNode) value.Value {
	if node.Start == nil {
		switch node.Op.Type {
		case token.CLOSED_RANGE_OP, token.LEFT_OPEN_RANGE_OP:
			to := resolve(node.End)
			if to.IsUndefined() {
				return value.Undefined
			}
			return value.Ref(value.NewBeginlessClosedRange(to))
		case token.RIGHT_OPEN_RANGE_OP, token.OPEN_RANGE_OP:
			to := resolve(node.End)
			if to.IsUndefined() {
				return value.Undefined
			}
			return value.Ref(value.NewBeginlessOpenRange(to))
		default:
			return value.Undefined
		}
	}

	if node.End == nil {
		switch node.Op.Type {
		case token.CLOSED_RANGE_OP, token.RIGHT_OPEN_RANGE_OP:
			from := resolve(node.Start)
			if from.IsUndefined() {
				return value.Undefined
			}
			return value.Ref(value.NewEndlessClosedRange(from))
		case token.LEFT_OPEN_RANGE_OP, token.OPEN_RANGE_OP:
			from := resolve(node.Start)
			if from.IsUndefined() {
				return value.Undefined
			}
			return value.Ref(value.NewEndlessOpenRange(from))
		default:
			return value.Undefined
		}
	}

	from := resolve(node.Start)
	if from.IsUndefined() {
		return value.Undefined
	}
	to := resolve(node.End)
	if to.IsUndefined() {
		return value.Undefined
	}

	switch node.Op.Type {
	case token.CLOSED_RANGE_OP:
		return value.Ref(value.NewClosedRange(from, to))
	case token.OPEN_RANGE_OP:
		return value.Ref(value.NewOpenRange(from, to))
	case token.LEFT_OPEN_RANGE_OP:
		return value.Ref(value.NewLeftOpenRange(from, to))
	case token.RIGHT_OPEN_RANGE_OP:
		return value.Ref(value.NewRightOpenRange(from, to))
	default:
		return value.Undefined
	}

}

func resolveHashSetLiteral(node *ast.HashSetLiteralNode) value.Value {
	if !node.IsStatic() || node.Capacity != nil {
		return value.Undefined
	}

	newTable := make([]value.Value, len(node.Elements))
	newSet := &value.HashSet{
		Table: newTable,
	}
	for _, elementNode := range node.Elements {
		val := resolve(elementNode)
		if val.IsUndefined() {
			return value.Undefined
		}
		err := vm.HashSetAppend(nil, newSet, val)
		if !err.IsUndefined() {
			return value.Undefined
		}
	}

	return value.Ref(newSet)
}

func resolveHashMapLiteral(node *ast.HashMapLiteralNode) value.Value {
	if !node.IsStatic() || node.Capacity != nil {
		return value.Undefined
	}

	newTable := make([]value.Pair, len(node.Elements))
	newMap := &value.HashMap{
		Table: newTable,
	}
	for _, elementNode := range node.Elements {
		switch element := elementNode.(type) {
		case *ast.SymbolKeyValueExpressionNode:
			key := value.ToSymbol(identifierToName(element.Key)).ToValue()
			val := resolve(element.Value)
			if val.IsUndefined() {
				return value.Undefined
			}

			err := vm.HashMapSet(nil, newMap, key, val)
			if !err.IsUndefined() {
				return value.Undefined
			}
		case *ast.KeyValueExpressionNode:
			key := resolve(element.Key)
			if key.IsUndefined() {
				return value.Undefined
			}
			val := resolve(element.Value)
			if val.IsUndefined() {
				return value.Undefined
			}

			err := vm.HashMapSet(nil, newMap, key, val)
			if !err.IsUndefined() {
				return value.Undefined
			}
		default:
			return value.Undefined
		}
	}

	return value.Ref(newMap)
}

func resolveHashRecordLiteral(node *ast.HashRecordLiteralNode) value.Value {
	if !node.IsStatic() {
		return value.Undefined
	}

	newTable := make([]value.Pair, len(node.Elements))
	newRecord := &value.HashRecord{
		Table: newTable,
	}
	for _, elementNode := range node.Elements {
		switch element := elementNode.(type) {
		case *ast.SymbolKeyValueExpressionNode:
			key := value.ToSymbol(identifierToName(element.Key)).ToValue()
			val := resolve(element.Value)
			if val.IsUndefined() {
				return value.Undefined
			}

			err := vm.HashRecordSet(nil, newRecord, key, val)
			if !err.IsUndefined() {
				return value.Undefined
			}
		case *ast.KeyValueExpressionNode:
			key := resolve(element.Key)
			if key.IsUndefined() {
				return value.Undefined
			}
			val := resolve(element.Value)
			if val.IsUndefined() {
				return value.Undefined
			}

			err := vm.HashRecordSet(nil, newRecord, key, val)
			if !err.IsUndefined() {
				return value.Undefined
			}
		default:
			return value.Undefined
		}
	}

	return value.Ref(newRecord)
}

func resolveSpecialHashSetLiteral[T ast.ExpressionNode](elements []T, static bool) value.Value {
	if !static {
		return value.Undefined
	}

	newSet := value.NewHashSet(len(elements))
	for _, elementNode := range elements {
		element := resolve(elementNode)
		if element.IsUndefined() {
			return value.Undefined
		}
		err := vm.HashSetAppend(nil, newSet, element)
		if !err.IsUndefined() {
			return value.Undefined
		}
	}

	return value.Ref(newSet)
}
func resolveArrayListLiteral(node *ast.ArrayListLiteralNode) value.Value {
	if !node.IsStatic() || node.Capacity != nil {
		return value.Undefined
	}

	newList := make(value.ArrayList, 0, len(node.Elements))
	for _, elementNode := range node.Elements {
		switch e := elementNode.(type) {
		case *ast.KeyValueExpressionNode:
			key := resolve(e.Key)
			if key.IsUndefined() {
				return value.Undefined
			}

			index, ok := value.ToGoInt(key)
			if !ok {
				return value.Undefined
			}

			val := resolve(e.Value)
			if val.IsUndefined() {
				return value.Undefined
			}

			if index >= len(newList) {
				newElementsCount := (index + 1) - len(newList)
				newList.Expand(newElementsCount)
			}
			newList[index] = val
		default:
			element := resolve(elementNode)
			if element.IsUndefined() {
				return value.Undefined
			}

			newList = append(newList, element)
		}
	}

	return value.Ref(&newList)
}

func resolveSpecialArrayListLiteral[T ast.ExpressionNode](elements []T, static bool) value.Value {
	if !static {
		return value.Undefined
	}

	newList := make(value.ArrayList, 0, len(elements))
	for _, elementNode := range elements {
		element := resolve(elementNode)
		if element.IsUndefined() {
			return value.Undefined
		}
		newList = append(newList, element)
	}

	return value.Ref(&newList)
}

func resolveSpecialArrayTupleLiteral[T ast.ExpressionNode](elements []T, static bool) value.Value {
	if !static {
		return value.Undefined
	}

	newList := make(value.ArrayTuple, 0, len(elements))
	for _, elementNode := range elements {
		element := resolve(elementNode)
		if element.IsUndefined() {
			return value.Undefined
		}
		newList = append(newList, element)
	}

	return value.Ref(&newList)
}

func resolveArrayTupleLiteral(node *ast.ArrayTupleLiteralNode) value.Value {
	if !node.IsStatic() {
		return value.Undefined
	}

	newArrayTuple := make(value.ArrayTuple, 0, len(node.Elements))
	for _, elementNode := range node.Elements {
		switch e := elementNode.(type) {
		case *ast.KeyValueExpressionNode:
			key := resolve(e.Key)
			if key.IsUndefined() {
				return value.Undefined
			}

			index, ok := value.ToGoInt(key)
			if !ok {
				return value.Undefined
			}

			val := resolve(e.Value)
			if val.IsUndefined() {
				return value.Undefined
			}

			if index >= len(newArrayTuple) {
				newElementsCount := (index + 1) - len(newArrayTuple)
				newArrayTuple.Expand(newElementsCount)
			}
			newArrayTuple[index] = val
		default:
			element := resolve(elementNode)
			if element.IsUndefined() {
				return value.Undefined
			}

			newArrayTuple = append(newArrayTuple, element)
		}
	}

	return value.Ref(&newArrayTuple)
}

func resolveLogicalExpression(node *ast.LogicalExpressionNode) value.Value {
	left := resolve(node.Left)
	if left.IsUndefined() {
		return value.Undefined
	}
	right := resolve(node.Right)
	if right.IsUndefined() {
		return value.Undefined
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

	return value.Undefined
}

func resolveNilSafeSubscript(node *ast.NilSafeSubscriptExpressionNode) value.Value {
	receiver := resolve(node.Receiver)
	key := resolve(node.Key)

	if receiver == value.Nil {
		return value.Nil
	}

	result, err := value.SubscriptVal(receiver, key)
	if !err.IsUndefined() {
		return value.Undefined
	}

	return result
}

func resolveSubscript(node *ast.SubscriptExpressionNode) value.Value {
	receiver := resolve(node.Receiver)
	key := resolve(node.Key)

	result, err := value.SubscriptVal(receiver, key)
	if !err.IsUndefined() {
		return value.Undefined
	}

	return result
}

func resolveUnaryExpression(node *ast.UnaryExpressionNode) value.Value {
	right := resolve(node.Right)
	if right.IsUndefined() {
		return value.Undefined
	}

	switch node.Op.Type {
	case token.TILDE:
		result := value.BitwiseNotVal(right)
		if result.IsUndefined() {
			return value.Undefined
		}
		return result
	case token.PLUS:
		result := value.UnaryPlusVal(right)
		if result.IsUndefined() {
			return value.Undefined
		}
		return result
	case token.MINUS:
		result := value.NegateVal(right)
		if result.IsUndefined() {
			return value.Undefined
		}
		return result
	case token.BANG:
		return value.ToNotBool(right)
	case token.AND:
		singleton := right.SingletonClass()
		if singleton == nil {
			return value.Undefined
		}
		return value.Ref(singleton)
	default:
		return value.Undefined
	}
}

func resolveBinaryExpression(node *ast.BinaryExpressionNode) value.Value {
	left := resolve(node.Left)
	if left.IsUndefined() {
		return value.Undefined
	}
	right := resolve(node.Right)
	if right.IsUndefined() {
		return value.Undefined
	}

	var result value.Value
	var err value.Value

	switch node.Op.Type {
	case token.PLUS:
		result, err = value.AddVal(left, right)
	case token.MINUS:
		result, err = value.SubtractVal(left, right)
	case token.STAR:
		result, err = value.MultiplyVal(left, right)
	case token.SLASH:
		result, err = value.DivideVal(left, right)
	case token.STAR_STAR:
		result, err = value.ExponentiateVal(left, right)
	case token.PERCENT:
		result, err = value.ModuloVal(left, right)
	case token.RBITSHIFT:
		result, err = value.RightBitshiftVal(left, right)
	case token.RTRIPLE_BITSHIFT:
		result, err = value.LogicalRightBitshiftVal(left, right)
	case token.LBITSHIFT:
		result, err = value.LeftBitshiftVal(left, right)
	case token.LTRIPLE_BITSHIFT:
		result, err = value.LogicalLeftBitshiftVal(left, right)
	case token.AND:
		result, err = value.BitwiseAndVal(left, right)
	case token.AND_TILDE:
		result, err = value.BitwiseAndNotVal(left, right)
	case token.OR:
		result, err = value.BitwiseOrVal(left, right)
	case token.XOR:
		result, err = value.BitwiseXorVal(left, right)
	case token.EQUAL_EQUAL:
		result = value.EqualVal(left, right)
	case token.LAX_EQUAL:
		result = value.LaxEqualVal(left, right)
	case token.LAX_NOT_EQUAL:
		result = value.LaxNotEqualVal(left, right)
	case token.NOT_EQUAL:
		result = value.NotEqualVal(left, right)
	case token.STRICT_EQUAL:
		result = value.StrictEqualVal(left, right)
	case token.STRICT_NOT_EQUAL:
		result = value.StrictNotEqualVal(left, right)
	case token.GREATER:
		result, err = value.GreaterThanVal(left, right)
	case token.GREATER_EQUAL:
		result, err = value.GreaterThanEqualVal(left, right)
	case token.LESS:
		result, err = value.LessThanVal(left, right)
	case token.LESS_EQUAL:
		result, err = value.LessThanEqualVal(left, right)
	default:
		return value.Undefined
	}

	if !err.IsUndefined() {
		return value.Undefined
	}
	return result
}

func resolveInt(node *ast.IntLiteralNode) value.Value {
	i, err := value.ParseBigInt(node.Value, 0)
	if !err.IsUndefined() {
		return value.Undefined
	}
	if i.IsSmallInt() {
		return i.ToSmallInt().ToValue()
	}

	return value.Ref(i)
}

func resolveInt64(node *ast.Int64LiteralNode) value.Value {
	i, err := value.StrictParseInt(node.Value, 0, 64)
	if !err.IsUndefined() {
		return value.Undefined
	}

	return value.Int64(i).ToValue()
}

func resolveInt32(node *ast.Int32LiteralNode) value.Value {
	i, err := value.StrictParseInt(node.Value, 0, 32)
	if !err.IsUndefined() {
		return value.Undefined
	}

	return value.Int32(i).ToValue()
}

func resolveInt16(node *ast.Int16LiteralNode) value.Value {
	i, err := value.StrictParseInt(node.Value, 0, 16)
	if !err.IsUndefined() {
		return value.Undefined
	}

	return value.Int16(i).ToValue()
}

func resolveInt8(node *ast.Int8LiteralNode) value.Value {
	i, err := value.StrictParseInt(node.Value, 0, 8)
	if !err.IsUndefined() {
		return value.Undefined
	}

	return value.Int8(i).ToValue()
}

func resolveUInt64(node *ast.UInt64LiteralNode) value.Value {
	i, err := value.StrictParseUint(node.Value, 0, 64)
	if !err.IsUndefined() {
		return value.Undefined
	}

	return value.UInt64(i).ToValue()
}

func resolveUInt32(node *ast.UInt32LiteralNode) value.Value {
	i, err := value.StrictParseUint(node.Value, 0, 32)
	if !err.IsUndefined() {
		return value.Undefined
	}

	return value.UInt32(i).ToValue()
}

func resolveUInt16(node *ast.UInt16LiteralNode) value.Value {
	i, err := value.StrictParseUint(node.Value, 0, 16)
	if !err.IsUndefined() {
		return value.Undefined
	}

	return value.UInt16(i).ToValue()
}

func resolveUInt8(node *ast.UInt8LiteralNode) value.Value {
	i, err := value.StrictParseUint(node.Value, 0, 8)
	if !err.IsUndefined() {
		return value.Undefined
	}

	return value.UInt8(i).ToValue()
}

func resolveBigFloat(node *ast.BigFloatLiteralNode) value.Value {
	f, err := value.ParseBigFloat(node.Value)
	if !err.IsUndefined() {
		return value.Undefined
	}

	return value.Ref(f)
}

func resolveFloat64(node *ast.Float64LiteralNode) value.Value {
	f, err := strconv.ParseFloat(node.Value, 64)
	if err != nil {
		return value.Undefined
	}

	return value.Float64(f).ToValue()
}

func resolveFloat32(node *ast.Float32LiteralNode) value.Value {
	f, err := strconv.ParseFloat(node.Value, 32)
	if err != nil {
		return value.Undefined
	}

	return value.Float32(f).ToValue()
}

func resolveFloat(node *ast.FloatLiteralNode) value.Value {
	f, err := strconv.ParseFloat(node.Value, 64)
	if err != nil {
		return value.Undefined
	}

	return value.Float(f).ToValue()
}
