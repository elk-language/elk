// Package comparer implements
// comparer functions for Elk values.
package comparer

import (
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/value"
	"github.com/elk-language/elk/vm"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Options() cmp.Options {
	return *Comparer
}

var Comparer *cmp.Options

func init() {
	bigFloatComparer := cmp.Comparer(func(x, y *value.BigFloat) bool {
		if x.IsNaN() || y.IsNaN() {
			return x.IsNaN() && y.IsNaN()
		}
		return x.AsGoBigFloat().Cmp(y.AsGoBigFloat()) == 0 &&
			(x.IsInf(0) || y.IsInf(0) || x.Precision() == y.Precision())
	})
	floatComparer := cmp.Comparer(func(x, y value.Float) bool {
		if x.IsNaN() || y.IsNaN() {
			return x.IsNaN() && y.IsNaN()
		}
		return x == y
	})
	float64Comparer := cmp.Comparer(func(x, y value.Float64) bool {
		if x.IsNaN() || y.IsNaN() {
			return x.IsNaN() && y.IsNaN()
		}
		return x == y
	})
	float32Comparer := cmp.Comparer(func(x, y value.Float32) bool {
		if x.IsNaN() || y.IsNaN() {
			return x.IsNaN() && y.IsNaN()
		}
		return x == y
	})

	opts := make(cmp.Options, 0, 30)
	Comparer = &opts
	*Comparer = append(
		*Comparer,
		cmp.AllowUnexported(value.SymbolTableStruct{}),
		cmpopts.IgnoreFields(value.SymbolTableStruct{}, "mutex"),
		cmp.AllowUnexported(
			value.Object{},
			value.BigInt{},
			value.Class{},
			bitfield.BitField8{},
		),
		cmp.AllowUnexported(vm.BytecodeFunction{}, vm.GetterMethod{}, vm.SetterMethod{}),
		floatComparer,
		bigFloatComparer,
		float32Comparer,
		float64Comparer,
		value.NewSymbolTableComparer(),
		vm.NewNativeMethodComparer(),
		// value.NewArrayListComparer(Comparer),
		value.NewObjectComparer(Comparer),
		value.NewClassComparer(Comparer),
		value.NewModuleComparer(Comparer),
		value.NewRegexComparer(Comparer),
		value.NewReferenceComparer(),
		value.NewInlineValueComparer(Comparer),
		vm.NewHashSetComparer(Comparer),
		vm.NewHashMapComparer(Comparer),
		vm.NewHashRecordComparer(Comparer),
		cmp.AllowUnexported(
			value.Result{},
			ast.NodeBase{},
			ast.TypedNodeBase{},
			token.Token{},
			ast.MacroDefinitionNode{},
			ast.DocCommentableNodeBase{},
			ast.BinaryExpressionNode{},
			ast.LogicalExpressionNode{},
			ast.KeyValueExpressionNode{},
			ast.ArrayListLiteralNode{},
			ast.ArrayTupleLiteralNode{},
			ast.HashSetLiteralNode{},
			ast.HashMapLiteralNode{},
			ast.HashRecordLiteralNode{},
			ast.RangeLiteralNode{},
			ast.SubscriptExpressionNode{},
			ast.NilSafeSubscriptExpressionNode{},
			ast.WordArrayListLiteralNode{},
			ast.WordHashSetLiteralNode{},
			ast.SymbolArrayListLiteralNode{},
			ast.SymbolHashSetLiteralNode{},
			ast.BinArrayListLiteralNode{},
			ast.BinHashSetLiteralNode{},
			ast.HexArrayListLiteralNode{},
			ast.HexHashSetLiteralNode{},
			ast.UninterpolatedRegexLiteralNode{},
			bitfield.BitField8{},
		),
		cmpopts.IgnoreFields(
			ast.TypedNodeBase{}, "typ",
		),
	)
}
