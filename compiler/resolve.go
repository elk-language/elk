package compiler

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/regex"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/value/symbol"
	"github.com/elk-language/elk/vm"
)

// Create Elk runtime values from static AST nodes.
// Returns undefined when no value could be created.
func resolve(node ast.Node, checker types.Checker) value.Value {
	if !node.IsStatic() {
		return value.Undefined
	}

	switch n := node.(type) {
	case *ast.LabeledExpressionNode:
		return resolve(n.Expression, checker)
	case *ast.UninterpolatedRegexLiteralNode:
		return resolveUninterpolatedRegexLiteral(n)
	case *ast.RangeLiteralNode:
		return resolveRangeLiteral(n, checker)
	case *ast.HashSetLiteralNode:
		return resolveHashSetLiteral(n, checker)
	case *ast.WordHashSetLiteralNode:
		return resolveSpecialNativeHashSetLiteral[ast.WordCollectionContentNode, value.String](n.Elements, checker, n.IsStatic())
	case *ast.SymbolHashSetLiteralNode:
		return resolveSpecialNativeHashSetLiteral[ast.SymbolCollectionContentNode, value.Symbol](n.Elements, checker, n.IsStatic())
	case *ast.BinHashSetLiteralNode:
		return resolveSpecialHashSetLiteral(n.Elements, checker, n.IsStatic())
	case *ast.HexHashSetLiteralNode:
		return resolveSpecialHashSetLiteral(n.Elements, checker, n.IsStatic())
	case *ast.HashMapLiteralNode:
		return resolveHashMapLiteral(n, checker)
	case *ast.HashRecordLiteralNode:
		return resolveHashRecordLiteral(n, checker)
	case *ast.ArrayListLiteralNode:
		return resolveArrayListLiteral(n, checker)
	case *ast.WordArrayListLiteralNode:
		return resolveSpecialNativeArrayListLiteral[ast.WordCollectionContentNode, value.String](n.Elements, checker, n.IsStatic())
	case *ast.SymbolArrayListLiteralNode:
		return resolveSpecialNativeArrayListLiteral[ast.SymbolCollectionContentNode, value.Symbol](n.Elements, checker, n.IsStatic())
	case *ast.BinArrayListLiteralNode:
		return resolveIntArrayListLiteral(n.Elements, n.Type(checker.Env()), checker, n.IsStatic())
	case *ast.HexArrayListLiteralNode:
		return resolveIntArrayListLiteral(n.Elements, n.Type(checker.Env()), checker, n.IsStatic())
	case *ast.ArrayTupleLiteralNode:
		return resolveArrayTupleLiteral(n, checker)
	case *ast.WordArrayTupleLiteralNode:
		return resolveSpecialNativeArrayTupleLiteral[ast.WordCollectionContentNode, value.String](n.Elements, checker, n.IsStatic())
	case *ast.SymbolArrayTupleLiteralNode:
		return resolveSpecialNativeArrayTupleLiteral[ast.SymbolCollectionContentNode, value.Symbol](n.Elements, checker, n.IsStatic())
	case *ast.BinArrayTupleLiteralNode:
		return resolveIntArrayTupleLiteral(n.Elements, n.Type(checker.Env()), checker, n.IsStatic())
	case *ast.HexArrayTupleLiteralNode:
		return resolveIntArrayTupleLiteral(n.Elements, n.Type(checker.Env()), checker, n.IsStatic())
	case *ast.LogicalExpressionNode:
		return resolveLogicalExpression(n, checker)
	case *ast.BinaryExpressionNode:
		return resolveBinaryExpression(n, checker)
	case *ast.UnaryExpressionNode:
		return resolveUnaryExpression(n, checker)
	case *ast.SubscriptExpressionNode:
		return resolveSubscript(n, checker)
	case *ast.NilSafeSubscriptExpressionNode:
		return resolveNilSafeSubscript(n, checker)
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
		return value.True.ToValue()
	case *ast.FalseLiteralNode:
		return value.False.ToValue()
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
		return resolveMacroBoundary(n, checker)
	}

	return value.Undefined
}

func resolveMacroBoundary(n *ast.MacroBoundaryNode, checker types.Checker) value.Value {
	if len(n.Body) != 1 {
		return value.Undefined
	}

	stmt := n.Body[0]
	exprStmt, ok := stmt.(*ast.ExpressionStatementNode)
	if !ok {
		return value.Undefined
	}

	return resolve(exprStmt.Expression, checker)
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

func resolveRangeLiteral(node *ast.RangeLiteralNode, checker types.Checker) value.Value {
	if node.Start == nil {
		switch node.Op.Type {
		case token.CLOSED_RANGE_OP, token.LEFT_OPEN_RANGE_OP:
			to := resolve(node.End, checker)
			if to.IsUndefined() {
				return value.Undefined
			}
			return value.Ref(value.NewBeginlessClosedRange(to))
		case token.RIGHT_OPEN_RANGE_OP, token.OPEN_RANGE_OP:
			to := resolve(node.End, checker)
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
			from := resolve(node.Start, checker)
			if from.IsUndefined() {
				return value.Undefined
			}
			return value.Ref(value.NewEndlessClosedRange(from))
		case token.LEFT_OPEN_RANGE_OP, token.OPEN_RANGE_OP:
			from := resolve(node.Start, checker)
			if from.IsUndefined() {
				return value.Undefined
			}
			return value.Ref(value.NewEndlessOpenRange(from))
		default:
			return value.Undefined
		}
	}

	from := resolve(node.Start, checker)
	if from.IsUndefined() {
		return value.Undefined
	}
	to := resolve(node.End, checker)
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

func resolveHashSetLiteral(node *ast.HashSetLiteralNode, checker types.Checker) value.Value {
	if !node.IsStatic() || node.Capacity != nil {
		return value.Undefined
	}

	typ := node.Type(checker.Env())
	elementType, _ := checker.GetIteratorElementType(typ)
	if types.IsUntyped(elementType) {
		return resolveHashSetOfValue(node, checker)
	}

	if checker.IsSubtype(elementType, checker.Std(symbol.String)) {
		return resolveNativeHashSet[value.String](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Symbol)) {
		return resolveNativeHashSet[value.Symbol](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt)) {
		return resolveNativeHashSet[value.UInt](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt64)) {
		return resolveNativeHashSet[value.UInt64](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Int64)) {
		return resolveNativeHashSet[value.Int64](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt32)) {
		return resolveNativeHashSet[value.UInt32](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Int32)) {
		return resolveNativeHashSet[value.Int32](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt16)) {
		return resolveNativeHashSet[value.UInt16](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Int16)) {
		return resolveNativeHashSet[value.Int16](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt8)) {
		return resolveNativeHashSet[value.UInt8](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Int8)) {
		return resolveNativeHashSet[value.Int8](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Float)) {
		return resolveNativeHashSet[value.Float](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Float64)) {
		return resolveNativeHashSet[value.Float64](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Float32)) {
		return resolveNativeHashSet[value.Float32](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Char)) {
		return resolveNativeHashSet[value.Char](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Date)) {
		return resolveNativeHashSet[value.Date](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Time)) {
		return resolveNativeHashSet[value.Time](node, checker)
	}

	return resolveHashSetOfValue(node, checker)
}

func resolveNativeHashSet[V value.ComparableValueInterface](node *ast.HashSetLiteralNode, checker types.Checker) value.Value {
	newSet := vm.NewNativeHashSet[V](len(node.Elements))
	for _, elementNode := range node.Elements {
		val := resolve(elementNode, checker)
		if val.IsUndefined() {
			return value.Undefined
		}
		v, ok := value.Downcast[V](val)
		if !ok {
			return value.Undefined
		}

		newSet.Append(v)
	}

	return value.Ref(newSet)
}

func resolveHashSetOfValue(node *ast.HashSetLiteralNode, checker types.Checker) value.Value {
	newSet := vm.NewHashSetOfValue(len(node.Elements))
	for _, elementNode := range node.Elements {
		val := resolve(elementNode, checker)
		if val.IsUndefined() {
			return value.Undefined
		}
		_, err := vm.HashSetOfValueAppend(nil, newSet, val)
		if !err.IsUndefined() {
			return value.Undefined
		}
	}

	return value.Ref(newSet)
}

func resolveHashMapLiteral(node *ast.HashMapLiteralNode, checker types.Checker) value.Value {
	if !node.IsStatic() || node.Capacity != nil {
		return value.Undefined
	}

	typ := node.Type(checker.Env())
	elementType, _ := checker.GetIteratorElementType(typ)
	g, ok := elementType.(*types.Generic)
	if !ok {
		return resolveHashMapOfValue(node, checker)
	}
	if !checker.IsTheSameNamespace(g.Namespace, checker.Std(symbol.Pair).(*types.Class)) {
		return value.Undefined
	}

	keyType := g.Get(0).Type
	valType := g.Get(1).Type

	if checker.IsSubtype(keyType, checker.Std(symbol.String)) {
		return resolveNativeHashMapOfString(node, valType, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Symbol)) {
		return resolveNativeHashMapOfSymbol(node, valType, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Char)) {
		return resolveNativeHashMapOfChar(node, valType, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Float)) {
		return resolveNativeHashMapOfFloat(node, valType, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Float64)) {
		return resolveNativeKeyHashMap[value.Float64](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Float32)) {
		return resolveNativeKeyHashMap[value.Float32](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.UInt)) {
		return resolveNativeKeyHashMap[value.UInt](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.UInt64)) {
		return resolveNativeKeyHashMap[value.UInt64](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Int64)) {
		return resolveNativeKeyHashMap[value.Int64](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.UInt32)) {
		return resolveNativeKeyHashMap[value.UInt32](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Int32)) {
		return resolveNativeKeyHashMap[value.Int32](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.UInt16)) {
		return resolveNativeKeyHashMap[value.UInt16](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Int16)) {
		return resolveNativeKeyHashMap[value.Int16](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.UInt8)) {
		return resolveNativeKeyHashMap[value.UInt8](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Int8)) {
		return resolveNativeKeyHashMap[value.Int8](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Date)) {
		return resolveNativeKeyHashMap[value.Date](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Time)) {
		return resolveNativeKeyHashMap[value.Time](node, checker)
	}

	return resolveHashMapOfValue(node, checker)
}

func resolveNativeHashMapOfString(node *ast.HashMapLiteralNode, valType types.Type, checker types.Checker) value.Value {
	if checker.IsSubtype(valType, checker.Std(symbol.String)) {
		return resolveNativeHashMap[value.String, value.String](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Symbol)) {
		return resolveNativeHashMap[value.String, value.Symbol](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt)) {
		return resolveNativeHashMap[value.String, value.UInt](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt64)) {
		return resolveNativeHashMap[value.String, value.UInt64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int64)) {
		return resolveNativeHashMap[value.String, value.Int64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt32)) {
		return resolveNativeHashMap[value.String, value.UInt32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int32)) {
		return resolveNativeHashMap[value.String, value.Int32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt16)) {
		return resolveNativeHashMap[value.String, value.UInt16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int16)) {
		return resolveNativeHashMap[value.String, value.Int16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt8)) {
		return resolveNativeHashMap[value.String, value.UInt8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int8)) {
		return resolveNativeHashMap[value.String, value.Int8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float)) {
		return resolveNativeHashMap[value.String, value.Float](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float64)) {
		return resolveNativeHashMap[value.String, value.Float64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float32)) {
		return resolveNativeHashMap[value.String, value.Float32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Char)) {
		return resolveNativeHashMap[value.String, value.Char](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Bool)) {
		return resolveNativeHashMap[value.String, value.Bool](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Date)) {
		return resolveNativeHashMap[value.String, value.Date](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Time)) {
		return resolveNativeHashMap[value.String, value.Time](node, checker)
	}

	return resolveNativeKeyHashMap[value.String](node, checker)
}

func resolveNativeHashMapOfFloat(node *ast.HashMapLiteralNode, valType types.Type, checker types.Checker) value.Value {
	if checker.IsSubtype(valType, checker.Std(symbol.String)) {
		return resolveNativeHashMap[value.Float, value.String](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Symbol)) {
		return resolveNativeHashMap[value.Float, value.Symbol](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt)) {
		return resolveNativeHashMap[value.Float, value.UInt](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt64)) {
		return resolveNativeHashMap[value.Float, value.UInt64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int64)) {
		return resolveNativeHashMap[value.Float, value.Int64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt32)) {
		return resolveNativeHashMap[value.Float, value.UInt32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int32)) {
		return resolveNativeHashMap[value.Float, value.Int32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt16)) {
		return resolveNativeHashMap[value.Float, value.UInt16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int16)) {
		return resolveNativeHashMap[value.Float, value.Int16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt8)) {
		return resolveNativeHashMap[value.Float, value.UInt8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int8)) {
		return resolveNativeHashMap[value.Float, value.Int8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float)) {
		return resolveNativeHashMap[value.Float, value.Float](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float64)) {
		return resolveNativeHashMap[value.Float, value.Float64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float32)) {
		return resolveNativeHashMap[value.Float, value.Float32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Char)) {
		return resolveNativeHashMap[value.Float, value.Char](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Bool)) {
		return resolveNativeHashMap[value.Float, value.Bool](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Date)) {
		return resolveNativeHashMap[value.Float, value.Date](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Time)) {
		return resolveNativeHashMap[value.Float, value.Time](node, checker)
	}

	return resolveNativeKeyHashMap[value.Float](node, checker)
}

func resolveNativeHashMapOfChar(node *ast.HashMapLiteralNode, valType types.Type, checker types.Checker) value.Value {
	if checker.IsSubtype(valType, checker.Std(symbol.String)) {
		return resolveNativeHashMap[value.Char, value.String](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Symbol)) {
		return resolveNativeHashMap[value.Char, value.Symbol](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt)) {
		return resolveNativeHashMap[value.Char, value.UInt](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt64)) {
		return resolveNativeHashMap[value.Char, value.UInt64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int64)) {
		return resolveNativeHashMap[value.Char, value.Int64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt32)) {
		return resolveNativeHashMap[value.Char, value.UInt32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int32)) {
		return resolveNativeHashMap[value.Char, value.Int32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt16)) {
		return resolveNativeHashMap[value.Char, value.UInt16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int16)) {
		return resolveNativeHashMap[value.Char, value.Int16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt8)) {
		return resolveNativeHashMap[value.Char, value.UInt8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int8)) {
		return resolveNativeHashMap[value.Char, value.Int8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float)) {
		return resolveNativeHashMap[value.Char, value.Float](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float64)) {
		return resolveNativeHashMap[value.Char, value.Float64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float32)) {
		return resolveNativeHashMap[value.Char, value.Float32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Char)) {
		return resolveNativeHashMap[value.Char, value.Char](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Bool)) {
		return resolveNativeHashMap[value.Char, value.Bool](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Date)) {
		return resolveNativeHashMap[value.Char, value.Date](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Time)) {
		return resolveNativeHashMap[value.Char, value.Time](node, checker)
	}

	return resolveNativeKeyHashMap[value.Char](node, checker)
}

func resolveNativeHashMapOfSymbol(node *ast.HashMapLiteralNode, valType types.Type, checker types.Checker) value.Value {
	if checker.IsSubtype(valType, checker.Std(symbol.String)) {
		return resolveNativeHashMap[value.Symbol, value.String](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Symbol)) {
		return resolveNativeHashMap[value.Symbol, value.Symbol](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt)) {
		return resolveNativeHashMap[value.Symbol, value.UInt](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt64)) {
		return resolveNativeHashMap[value.Symbol, value.UInt64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int64)) {
		return resolveNativeHashMap[value.Symbol, value.Int64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt32)) {
		return resolveNativeHashMap[value.Symbol, value.UInt32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int32)) {
		return resolveNativeHashMap[value.Symbol, value.Int32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt16)) {
		return resolveNativeHashMap[value.Symbol, value.UInt16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int16)) {
		return resolveNativeHashMap[value.Symbol, value.Int16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt8)) {
		return resolveNativeHashMap[value.Symbol, value.UInt8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int8)) {
		return resolveNativeHashMap[value.Symbol, value.Int8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float)) {
		return resolveNativeHashMap[value.Symbol, value.Float](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float64)) {
		return resolveNativeHashMap[value.Symbol, value.Float64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float32)) {
		return resolveNativeHashMap[value.Symbol, value.Float32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Char)) {
		return resolveNativeHashMap[value.Symbol, value.Char](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Bool)) {
		return resolveNativeHashMap[value.Symbol, value.Bool](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Date)) {
		return resolveNativeHashMap[value.Symbol, value.Date](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Time)) {
		return resolveNativeHashMap[value.Symbol, value.Time](node, checker)
	}

	return resolveNativeKeyHashMap[value.Symbol](node, checker)
}

func resolveNativeHashMap[K value.ComparableValueInterface, V value.ValueInterface](node *ast.HashMapLiteralNode, checker types.Checker) value.Value {
	newMap := vm.NewNativeHashMap[K, V](len(node.Elements))
	for _, elementNode := range node.Elements {
		switch element := elementNode.(type) {
		case *ast.SymbolKeyValueExpressionNode:
			key := value.ToSymbol(identifierToName(element.Key)).ToValue()
			k, ok := value.Downcast[K](key)
			if !ok {
				return value.Undefined
			}
			val := resolve(element.Value, checker)
			v, ok := value.Downcast[V](val)
			if !ok {
				return value.Undefined
			}

			newMap.Set(k, v)
		case *ast.KeyValueExpressionNode:
			key := resolve(element.Key, checker)
			k, ok := value.Downcast[K](key)
			if !ok {
				return value.Undefined
			}
			val := resolve(element.Value, checker)
			v, ok := value.Downcast[V](val)
			if !ok {
				return value.Undefined
			}

			newMap.Set(k, v)
		default:
			return value.Undefined
		}
	}

	return newMap.ToValue()
}

func resolveNativeKeyHashMap[K value.ComparableValueInterface](node *ast.HashMapLiteralNode, checker types.Checker) value.Value {
	newMap := vm.NewNativeKeyHashMap[K](len(node.Elements))
	for _, elementNode := range node.Elements {
		switch element := elementNode.(type) {
		case *ast.SymbolKeyValueExpressionNode:
			key := value.ToSymbol(identifierToName(element.Key)).ToValue()
			k, ok := value.Downcast[K](key)
			if !ok {
				return value.Undefined
			}
			val := resolve(element.Value, checker)
			if val.IsUndefined() {
				return value.Undefined
			}

			newMap.Set(k, val)
		case *ast.KeyValueExpressionNode:
			key := resolve(element.Key, checker)
			k, ok := value.Downcast[K](key)
			if !ok {
				return value.Undefined
			}
			val := resolve(element.Value, checker)
			if val.IsUndefined() {
				return value.Undefined
			}

			newMap.Set(k, val)
		default:
			return value.Undefined
		}
	}

	return newMap.ToValue()
}

func resolveHashMapOfValue(node *ast.HashMapLiteralNode, checker types.Checker) value.Value {
	newMap := vm.NewHashMapOfValue(len(node.Elements))
	for _, elementNode := range node.Elements {
		switch element := elementNode.(type) {
		case *ast.SymbolKeyValueExpressionNode:
			key := value.ToSymbol(identifierToName(element.Key)).ToValue()
			val := resolve(element.Value, checker)
			if val.IsUndefined() {
				return value.Undefined
			}

			err := vm.HashMapOfValueSet(nil, newMap, key, val)
			if !err.IsUndefined() {
				return value.Undefined
			}
		case *ast.KeyValueExpressionNode:
			key := resolve(element.Key, checker)
			if key.IsUndefined() {
				return value.Undefined
			}
			val := resolve(element.Value, checker)
			if val.IsUndefined() {
				return value.Undefined
			}

			err := vm.HashMapOfValueSet(nil, newMap, key, val)
			if !err.IsUndefined() {
				return value.Undefined
			}
		default:
			return value.Undefined
		}
	}

	return newMap.ToValue()
}

func resolveHashRecordLiteral(node *ast.HashRecordLiteralNode, checker types.Checker) value.Value {
	if !node.IsStatic() {
		return value.Undefined
	}

	typ := node.Type(checker.Env())
	elementType, _ := checker.GetIteratorElementType(typ)
	g, ok := elementType.(*types.Generic)
	if !ok {
		return resolveHashRecordOfValue(node, checker)
	}
	if !checker.IsTheSameNamespace(g.Namespace, checker.Std(symbol.Pair).(*types.Class)) {
		return value.Undefined
	}

	keyType := g.Get(0).Type
	valType := g.Get(1).Type

	if checker.IsSubtype(keyType, checker.Std(symbol.String)) {
		return resolveNativeHashRecordOfString(node, valType, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Symbol)) {
		return resolveNativeHashRecordOfSymbol(node, valType, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Char)) {
		return resolveNativeHashRecordOfChar(node, valType, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Float)) {
		return resolveNativeHashRecordOfFloat(node, valType, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Float64)) {
		return resolveNativeKeyHashRecord[value.Float64](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Float32)) {
		return resolveNativeKeyHashRecord[value.Float32](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.UInt)) {
		return resolveNativeKeyHashRecord[value.UInt](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.UInt64)) {
		return resolveNativeKeyHashRecord[value.UInt64](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Int64)) {
		return resolveNativeKeyHashRecord[value.Int64](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.UInt32)) {
		return resolveNativeKeyHashRecord[value.UInt32](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Int32)) {
		return resolveNativeKeyHashRecord[value.Int32](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.UInt16)) {
		return resolveNativeKeyHashRecord[value.UInt16](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Int16)) {
		return resolveNativeKeyHashRecord[value.Int16](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.UInt8)) {
		return resolveNativeKeyHashRecord[value.UInt8](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Int8)) {
		return resolveNativeKeyHashRecord[value.Int8](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Date)) {
		return resolveNativeKeyHashRecord[value.Date](node, checker)
	}
	if checker.IsSubtype(keyType, checker.Std(symbol.Time)) {
		return resolveNativeKeyHashRecord[value.Time](node, checker)
	}

	return resolveHashRecordOfValue(node, checker)
}

func resolveNativeHashRecordOfString(node *ast.HashRecordLiteralNode, valType types.Type, checker types.Checker) value.Value {
	if checker.IsSubtype(valType, checker.Std(symbol.String)) {
		return resolveNativeHashRecord[value.String, value.String](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Symbol)) {
		return resolveNativeHashRecord[value.String, value.Symbol](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt)) {
		return resolveNativeHashRecord[value.String, value.UInt](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt64)) {
		return resolveNativeHashRecord[value.String, value.UInt64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int64)) {
		return resolveNativeHashRecord[value.String, value.Int64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt32)) {
		return resolveNativeHashRecord[value.String, value.UInt32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int32)) {
		return resolveNativeHashRecord[value.String, value.Int32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt16)) {
		return resolveNativeHashRecord[value.String, value.UInt16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int16)) {
		return resolveNativeHashRecord[value.String, value.Int16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt8)) {
		return resolveNativeHashRecord[value.String, value.UInt8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int8)) {
		return resolveNativeHashRecord[value.String, value.Int8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float)) {
		return resolveNativeHashRecord[value.String, value.Float](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float64)) {
		return resolveNativeHashRecord[value.String, value.Float64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float32)) {
		return resolveNativeHashRecord[value.String, value.Float32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Char)) {
		return resolveNativeHashRecord[value.String, value.Char](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Bool)) {
		return resolveNativeHashRecord[value.String, value.Bool](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Date)) {
		return resolveNativeHashRecord[value.String, value.Date](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Time)) {
		return resolveNativeHashRecord[value.String, value.Time](node, checker)
	}

	return resolveNativeKeyHashRecord[value.String](node, checker)
}

func resolveNativeHashRecordOfFloat(node *ast.HashRecordLiteralNode, valType types.Type, checker types.Checker) value.Value {
	if checker.IsSubtype(valType, checker.Std(symbol.String)) {
		return resolveNativeHashRecord[value.Float, value.String](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Symbol)) {
		return resolveNativeHashRecord[value.Float, value.Symbol](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt)) {
		return resolveNativeHashRecord[value.Float, value.UInt](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt64)) {
		return resolveNativeHashRecord[value.Float, value.UInt64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int64)) {
		return resolveNativeHashRecord[value.Float, value.Int64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt32)) {
		return resolveNativeHashRecord[value.Float, value.UInt32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int32)) {
		return resolveNativeHashRecord[value.Float, value.Int32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt16)) {
		return resolveNativeHashRecord[value.Float, value.UInt16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int16)) {
		return resolveNativeHashRecord[value.Float, value.Int16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt8)) {
		return resolveNativeHashRecord[value.Float, value.UInt8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int8)) {
		return resolveNativeHashRecord[value.Float, value.Int8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float)) {
		return resolveNativeHashRecord[value.Float, value.Float](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float64)) {
		return resolveNativeHashRecord[value.Float, value.Float64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float32)) {
		return resolveNativeHashRecord[value.Float, value.Float32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Char)) {
		return resolveNativeHashRecord[value.Float, value.Char](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Bool)) {
		return resolveNativeHashRecord[value.Float, value.Bool](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Date)) {
		return resolveNativeHashRecord[value.Float, value.Date](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Time)) {
		return resolveNativeHashRecord[value.Float, value.Time](node, checker)
	}

	return resolveNativeKeyHashRecord[value.Float](node, checker)
}

func resolveNativeHashRecordOfChar(node *ast.HashRecordLiteralNode, valType types.Type, checker types.Checker) value.Value {
	if checker.IsSubtype(valType, checker.Std(symbol.String)) {
		return resolveNativeHashRecord[value.Char, value.String](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Symbol)) {
		return resolveNativeHashRecord[value.Char, value.Symbol](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt)) {
		return resolveNativeHashRecord[value.Char, value.UInt](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt64)) {
		return resolveNativeHashRecord[value.Char, value.UInt64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int64)) {
		return resolveNativeHashRecord[value.Char, value.Int64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt32)) {
		return resolveNativeHashRecord[value.Char, value.UInt32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int32)) {
		return resolveNativeHashRecord[value.Char, value.Int32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt16)) {
		return resolveNativeHashRecord[value.Char, value.UInt16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int16)) {
		return resolveNativeHashRecord[value.Char, value.Int16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt8)) {
		return resolveNativeHashRecord[value.Char, value.UInt8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int8)) {
		return resolveNativeHashRecord[value.Char, value.Int8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float)) {
		return resolveNativeHashRecord[value.Char, value.Float](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float64)) {
		return resolveNativeHashRecord[value.Char, value.Float64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float32)) {
		return resolveNativeHashRecord[value.Char, value.Float32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Char)) {
		return resolveNativeHashRecord[value.Char, value.Char](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Bool)) {
		return resolveNativeHashRecord[value.Char, value.Bool](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Date)) {
		return resolveNativeHashRecord[value.Char, value.Date](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Time)) {
		return resolveNativeHashRecord[value.Char, value.Time](node, checker)
	}

	return resolveNativeKeyHashRecord[value.Char](node, checker)
}

func resolveNativeHashRecordOfSymbol(node *ast.HashRecordLiteralNode, valType types.Type, checker types.Checker) value.Value {
	if checker.IsSubtype(valType, checker.Std(symbol.String)) {
		return resolveNativeHashRecord[value.Symbol, value.String](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Symbol)) {
		return resolveNativeHashRecord[value.Symbol, value.Symbol](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt)) {
		return resolveNativeHashRecord[value.Symbol, value.UInt](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt64)) {
		return resolveNativeHashRecord[value.Symbol, value.UInt64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int64)) {
		return resolveNativeHashRecord[value.Symbol, value.Int64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt32)) {
		return resolveNativeHashRecord[value.Symbol, value.UInt32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int32)) {
		return resolveNativeHashRecord[value.Symbol, value.Int32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt16)) {
		return resolveNativeHashRecord[value.Symbol, value.UInt16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int16)) {
		return resolveNativeHashRecord[value.Symbol, value.Int16](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.UInt8)) {
		return resolveNativeHashRecord[value.Symbol, value.UInt8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Int8)) {
		return resolveNativeHashRecord[value.Symbol, value.Int8](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float)) {
		return resolveNativeHashRecord[value.Symbol, value.Float](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float64)) {
		return resolveNativeHashRecord[value.Symbol, value.Float64](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Float32)) {
		return resolveNativeHashRecord[value.Symbol, value.Float32](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Char)) {
		return resolveNativeHashRecord[value.Symbol, value.Char](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Bool)) {
		return resolveNativeHashRecord[value.Symbol, value.Bool](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Date)) {
		return resolveNativeHashRecord[value.Symbol, value.Date](node, checker)
	}
	if checker.IsSubtype(valType, checker.Std(symbol.Time)) {
		return resolveNativeHashRecord[value.Symbol, value.Time](node, checker)
	}

	return resolveNativeKeyHashRecord[value.Symbol](node, checker)
}

func resolveNativeHashRecord[K value.ComparableValueInterface, V value.ValueInterface](node *ast.HashRecordLiteralNode, checker types.Checker) value.Value {
	newMap := vm.MakeNativeHashRecord[K, V](len(node.Elements))
	for _, elementNode := range node.Elements {
		switch element := elementNode.(type) {
		case *ast.SymbolKeyValueExpressionNode:
			key := value.ToSymbol(identifierToName(element.Key)).ToValue()
			k, ok := value.Downcast[K](key)
			if !ok {
				return value.Undefined
			}
			val := resolve(element.Value, checker)
			v, ok := value.Downcast[V](val)
			if !ok {
				return value.Undefined
			}

			newMap[k] = v
		case *ast.KeyValueExpressionNode:
			key := resolve(element.Key, checker)
			k, ok := value.Downcast[K](key)
			if !ok {
				return value.Undefined
			}
			val := resolve(element.Value, checker)
			v, ok := value.Downcast[V](val)
			if !ok {
				return value.Undefined
			}

			newMap[k] = v
		default:
			return value.Undefined
		}
	}

	return newMap.ToValue()
}

func resolveNativeKeyHashRecord[K value.ComparableValueInterface](node *ast.HashRecordLiteralNode, checker types.Checker) value.Value {
	newRecord := vm.MakeNativeKeyHashRecord[K](len(node.Elements))
	for _, elementNode := range node.Elements {
		switch element := elementNode.(type) {
		case *ast.SymbolKeyValueExpressionNode:
			key := value.ToSymbol(identifierToName(element.Key)).ToValue()
			k, ok := value.Downcast[K](key)
			if !ok {
				return value.Undefined
			}
			val := resolve(element.Value, checker)
			if val.IsUndefined() {
				return value.Undefined
			}

			newRecord[k] = val
		case *ast.KeyValueExpressionNode:
			key := resolve(element.Key, checker)
			k, ok := value.Downcast[K](key)
			if !ok {
				return value.Undefined
			}
			val := resolve(element.Value, checker)
			if val.IsUndefined() {
				return value.Undefined
			}

			newRecord[k] = val
		default:
			return value.Undefined
		}
	}

	return newRecord.ToValue()
}

func resolveHashRecordOfValue(node *ast.HashRecordLiteralNode, checker types.Checker) value.Value {
	newRecord := vm.NewHashRecordOfValue(len(node.Elements))
	for _, elementNode := range node.Elements {
		switch element := elementNode.(type) {
		case *ast.SymbolKeyValueExpressionNode:
			key := value.ToSymbol(identifierToName(element.Key)).ToValue()
			val := resolve(element.Value, checker)
			if val.IsUndefined() {
				return value.Undefined
			}

			err := vm.HashRecordOfValueSet(nil, newRecord, key, val)
			if !err.IsUndefined() {
				return value.Undefined
			}
		case *ast.KeyValueExpressionNode:
			key := resolve(element.Key, checker)
			if key.IsUndefined() {
				return value.Undefined
			}
			val := resolve(element.Value, checker)
			if val.IsUndefined() {
				return value.Undefined
			}

			err := vm.HashRecordOfValueSet(nil, newRecord, key, val)
			if !err.IsUndefined() {
				return value.Undefined
			}
		default:
			return value.Undefined
		}
	}

	return value.Ref(newRecord)
}

func resolveSpecialHashSetLiteral[T ast.ExpressionNode](elements []T, checker types.Checker, static bool) value.Value {
	if !static {
		return value.Undefined
	}

	newSet := vm.NewHashSetOfValue(len(elements))
	for _, elementNode := range elements {
		element := resolve(elementNode, checker)
		if element.IsUndefined() {
			return value.Undefined
		}
		_, err := vm.HashSetOfValueAppend(nil, newSet, element)
		if !err.IsUndefined() {
			return value.Undefined
		}
	}

	return value.Ref(newSet)
}

func resolveSpecialNativeHashSetLiteral[N ast.ExpressionNode, T value.ComparableValueInterface](elements []N, checker types.Checker, static bool) value.Value {
	if !static {
		return value.Undefined
	}
	var t T

	newSet := vm.NewNativeHashSet[T](len(elements))
	for _, elementNode := range elements {
		element := resolve(elementNode, checker)
		if element.IsUndefined() {
			return value.Undefined
		}
		e, ok := value.Downcast[T](element)
		if !ok {
			panic(fmt.Sprintf("cannot cast %s to %T while resolving a hash set", element.Inspect(), t))
		}
		newSet.Append(e)
	}

	return value.Ref(newSet)
}

func resolveArrayListLiteral(node *ast.ArrayListLiteralNode, checker types.Checker) value.Value {
	if !node.IsStatic() || node.Capacity != nil {
		return value.Undefined
	}

	typ := node.Type(checker.Env())
	elementType, _ := checker.GetIteratorElementType(typ)
	if types.IsUntyped(elementType) {
		return resolveArrayListOfValue(node, checker)
	}

	if checker.IsSubtype(elementType, checker.Std(symbol.String)) {
		return resolveNativeArrayList[value.String](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Symbol)) {
		return resolveNativeArrayList[value.Symbol](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt)) {
		return resolveNativeArrayList[value.UInt](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt64)) {
		return resolveNativeArrayList[value.UInt64](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Int64)) {
		return resolveNativeArrayList[value.Int64](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt32)) {
		return resolveNativeArrayList[value.UInt32](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Int32)) {
		return resolveNativeArrayList[value.Int32](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt16)) {
		return resolveNativeArrayList[value.UInt16](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Int16)) {
		return resolveNativeArrayList[value.Int16](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt8)) {
		return resolveNativeArrayList[value.UInt8](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Int8)) {
		return resolveNativeArrayList[value.Int8](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Float)) {
		return resolveNativeArrayList[value.Float](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Float64)) {
		return resolveNativeArrayList[value.Float64](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Float32)) {
		return resolveNativeArrayList[value.Float32](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Char)) {
		return resolveNativeArrayList[value.Char](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Bool)) {
		return resolveNativeArrayList[value.Bool](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Date)) {
		return resolveNativeArrayList[value.Date](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Time)) {
		return resolveNativeArrayList[value.Time](node, checker)
	}

	return resolveArrayListOfValue(node, checker)
}

func resolveArrayListOfValue(node *ast.ArrayListLiteralNode, checker types.Checker) value.Value {
	newList := make(value.ArrayListOfValue, 0, len(node.Elements))
	for _, elementNode := range node.Elements {
		switch e := elementNode.(type) {
		case *ast.KeyValueExpressionNode:
			key := resolve(e.Key, checker)
			if key.IsUndefined() {
				return value.Undefined
			}

			index, ok := value.ToGoInt(key)
			if !ok {
				return value.Undefined
			}

			val := resolve(e.Value, checker)
			if val.IsUndefined() {
				return value.Undefined
			}

			if index >= len(newList) {
				newElementsCount := (index + 1) - len(newList)
				newList.Expand(newElementsCount)
			}
			newList[index] = val
		default:
			element := resolve(elementNode, checker)
			if element.IsUndefined() {
				return value.Undefined
			}

			newList = append(newList, element)
		}
	}

	return value.Ref(&newList)
}

func resolveNativeArrayList[T value.ValueInterface](node *ast.ArrayListLiteralNode, checker types.Checker) value.Value {
	newList := value.NewNativeArrayList[T](len(node.Elements))
	for _, elementNode := range node.Elements {
		element := resolve(elementNode, checker)
		if element.IsUndefined() {
			return value.Undefined
		}

		e, ok := value.Downcast[T](element)
		if !ok {
			return value.Undefined
		}

		newList.Append(e)
	}

	return newList.ToValue()
}

func resolveNativeArrayTuple[T value.ValueInterface](node *ast.ArrayTupleLiteralNode, checker types.Checker) value.Value {
	newTuple := value.NewNativeArrayTuple[T](len(node.Elements))
	for _, elementNode := range node.Elements {
		element := resolve(elementNode, checker)
		if element.IsUndefined() {
			return value.Undefined
		}

		e, ok := value.Downcast[T](element)
		if !ok {
			return value.Undefined
		}

		newTuple.Append(e)
	}

	return newTuple.ToValue()
}

func resolveSpecialNativeArrayListLiteral[N ast.ExpressionNode, T value.ValueInterface](elements []N, checker types.Checker, static bool) value.Value {
	if !static {
		return value.Undefined
	}
	var t T

	newList := value.NewNativeArrayList[T](len(elements))
	for _, elementNode := range elements {
		element := resolve(elementNode, checker)
		if element.IsUndefined() {
			return value.Undefined
		}
		e, ok := value.Downcast[T](element)
		if !ok {
			panic(fmt.Sprintf("cannot cast %s to %T while resolving an array list", element.Inspect(), t))
		}
		newList.Append(e)
	}

	return newList.ToValue()
}

func resolveSpecialNativeArrayTupleLiteral[N ast.ExpressionNode, T value.ValueInterface](elements []N, checker types.Checker, static bool) value.Value {
	if !static {
		return value.Undefined
	}
	var t T

	newTuple := value.NewNativeArrayTuple[T](len(elements))
	for _, elementNode := range elements {
		element := resolve(elementNode, checker)
		if element.IsUndefined() {
			return value.Undefined
		}
		e, ok := value.Downcast[T](element)
		if !ok {
			panic(fmt.Sprintf("cannot cast %s to %T while resolving an array list", element.Inspect(), t))
		}
		newTuple.Append(e)
	}

	return newTuple.ToValue()
}

func resolveIntArrayListLiteral(elements []ast.IntCollectionContentNode, typ types.Type, checker types.Checker, static bool) value.Value {
	if !static {
		return value.Undefined
	}

	tmpList := make([]*value.BigInt, 0, len(elements))
	for _, elementNode := range elements {
		n, ok := elementNode.(*ast.IntLiteralNode)
		if !ok {
			continue
		}

		val := value.ParseBigIntPanic(n.Value, 0)
		tmpList = append(tmpList, val)
	}

	g, ok := typ.(*types.Generic)
	if !ok {
		return value.Undefined
	}

	elementType := g.Get(0).Type
	if checker.IsSubtype(elementType, checker.Std(symbol.Int)) {
		return resolveBigIntSliceToArrayListOfValue(tmpList).ToValue()
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt64)) {
		return resolveBigIntSliceToNativeArrayList[value.UInt64](tmpList).ToValue()
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt32)) {
		return resolveBigIntSliceToNativeArrayList[value.UInt32](tmpList).ToValue()
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt16)) {
		return resolveBigIntSliceToNativeArrayList[value.UInt16](tmpList).ToValue()
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt8)) {
		return resolveBigIntSliceToNativeArrayList[value.UInt8](tmpList).ToValue()
	}

	return value.Undefined
}

func resolveIntArrayTupleLiteral(elements []ast.IntCollectionContentNode, typ types.Type, checker types.Checker, static bool) value.Value {
	if !static {
		return value.Undefined
	}

	tmpTuple := make([]*value.BigInt, 0, len(elements))
	for _, elementNode := range elements {
		n, ok := elementNode.(*ast.IntLiteralNode)
		if !ok {
			continue
		}

		val := value.ParseBigIntPanic(n.Value, 0)
		tmpTuple = append(tmpTuple, val)
	}

	g, ok := typ.(*types.Generic)
	if !ok {
		return value.Undefined
	}

	elementType := g.Get(0).Type
	if checker.IsSubtype(elementType, checker.Std(symbol.Int)) {
		return resolveBigIntSliceToArrayTupleOfValue(tmpTuple).ToValue()
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt64)) {
		return resolveBigIntSliceToNativeArrayTuple[value.UInt64](tmpTuple).ToValue()
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt32)) {
		return resolveBigIntSliceToNativeArrayTuple[value.UInt32](tmpTuple).ToValue()
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt16)) {
		return resolveBigIntSliceToNativeArrayTuple[value.UInt16](tmpTuple).ToValue()
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt8)) {
		return resolveBigIntSliceToNativeArrayTuple[value.UInt8](tmpTuple).ToValue()
	}

	return value.Undefined
}

func resolveBigIntSliceToNativeArrayList[I value.StrictNumeric](elements []*value.BigInt) *value.NativeArrayList[I] {
	newList := value.NewNativeArrayList[I](len(elements))
	for _, element := range elements {
		nativeVal := I(element.ToUInt64())
		newList.Append(nativeVal)
	}

	return newList
}

func resolveBigIntSliceToNativeArrayTuple[I value.StrictNumeric](elements []*value.BigInt) *value.NativeArrayTuple[I] {
	newTuple := value.NewNativeArrayTuple[I](len(elements))
	for _, element := range elements {
		nativeVal := I(element.ToUInt64())
		newTuple.Append(nativeVal)
	}

	return newTuple
}

func resolveBigIntSliceToArrayListOfValue(elements []*value.BigInt) *value.ArrayListOfValue {
	newList := value.NewArrayListOfValue(len(elements))
	for _, element := range elements {
		newList.Append(element.Normalize())
	}

	return newList
}

func resolveBigIntSliceToArrayTupleOfValue(elements []*value.BigInt) *value.ArrayTupleOfValue {
	newTuple := value.NewArrayTupleOfValue(len(elements))
	for _, element := range elements {
		newTuple.Append(element.Normalize())
	}

	return newTuple
}

func resolveArrayTupleLiteral(node *ast.ArrayTupleLiteralNode, checker types.Checker) value.Value {
	if !node.IsStatic() {
		return value.Undefined
	}

	typ := node.Type(checker.Env())
	elementType, _ := checker.GetIteratorElementType(typ)
	if types.IsUntyped(elementType) {
		return resolveArrayTupleOfValue(node, checker)
	}

	if checker.IsSubtype(elementType, checker.Std(symbol.String)) {
		return resolveNativeArrayTuple[value.String](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Symbol)) {
		return resolveNativeArrayTuple[value.Symbol](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt)) {
		return resolveNativeArrayTuple[value.UInt](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt64)) {
		return resolveNativeArrayTuple[value.UInt64](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Int64)) {
		return resolveNativeArrayTuple[value.Int64](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt32)) {
		return resolveNativeArrayTuple[value.UInt32](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Int32)) {
		return resolveNativeArrayTuple[value.Int32](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt16)) {
		return resolveNativeArrayTuple[value.UInt16](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Int16)) {
		return resolveNativeArrayTuple[value.Int16](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.UInt8)) {
		return resolveNativeArrayTuple[value.UInt8](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Int8)) {
		return resolveNativeArrayTuple[value.Int8](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Float)) {
		return resolveNativeArrayTuple[value.Float](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Float64)) {
		return resolveNativeArrayTuple[value.Float64](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Float32)) {
		return resolveNativeArrayTuple[value.Float32](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Char)) {
		return resolveNativeArrayTuple[value.Char](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Bool)) {
		return resolveNativeArrayTuple[value.Bool](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Time)) {
		return resolveNativeArrayTuple[value.Time](node, checker)
	}
	if checker.IsSubtype(elementType, checker.Std(symbol.Date)) {
		return resolveNativeArrayTuple[value.Date](node, checker)
	}

	return resolveArrayTupleOfValue(node, checker)
}

func resolveArrayTupleOfValue(node *ast.ArrayTupleLiteralNode, checker types.Checker) value.Value {
	newArrayTuple := make(value.ArrayTupleOfValue, 0, len(node.Elements))
	for _, elementNode := range node.Elements {
		switch e := elementNode.(type) {
		case *ast.KeyValueExpressionNode:
			key := resolve(e.Key, checker)
			if key.IsUndefined() {
				return value.Undefined
			}

			index, ok := value.ToGoInt(key)
			if !ok {
				return value.Undefined
			}

			val := resolve(e.Value, checker)
			if val.IsUndefined() {
				return value.Undefined
			}

			if index >= len(newArrayTuple) {
				newElementsCount := (index + 1) - len(newArrayTuple)
				newArrayTuple.Expand(newElementsCount)
			}
			newArrayTuple[index] = val
		default:
			element := resolve(elementNode, checker)
			if element.IsUndefined() {
				return value.Undefined
			}

			newArrayTuple = append(newArrayTuple, element)
		}
	}

	return value.Ref(&newArrayTuple)
}

func resolveLogicalExpression(node *ast.LogicalExpressionNode, checker types.Checker) value.Value {
	left := resolve(node.Left, checker)
	if left.IsUndefined() {
		return value.Undefined
	}
	right := resolve(node.Right, checker)
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

func resolveNilSafeSubscript(node *ast.NilSafeSubscriptExpressionNode, checker types.Checker) value.Value {
	receiver := resolve(node.Receiver, checker)
	key := resolve(node.Key, checker)

	if receiver == value.Nil {
		return value.Nil
	}

	result, err := value.SubscriptVal(receiver, key)
	if !err.IsUndefined() {
		return value.Undefined
	}

	return result
}

func resolveSubscript(node *ast.SubscriptExpressionNode, checker types.Checker) value.Value {
	receiver := resolve(node.Receiver, checker)
	key := resolve(node.Key, checker)

	result, err := value.SubscriptVal(receiver, key)
	if !err.IsUndefined() {
		return value.Undefined
	}

	return result
}

func resolveUnaryExpression(node *ast.UnaryExpressionNode, checker types.Checker) value.Value {
	right := resolve(node.Right, checker)
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
		return value.ToNotBool(right).ToValue()
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

func resolveBinaryExpression(node *ast.BinaryExpressionNode, checker types.Checker) value.Value {
	left := resolve(node.Left, checker)
	if left.IsUndefined() {
		return value.Undefined
	}
	right := resolve(node.Right, checker)
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
