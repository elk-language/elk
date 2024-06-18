// Package checker implements the Elk type checker
package checker

import (
	"testing"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/k0kubun/pp"
)

// Represents a single checker test case.
type testCase struct {
	input string
	err   error.ErrorList
}

// Type of the checker test table.
type testTable map[string]testCase

// Create a new token in tests.
var T = token.New

// Create a new token with value in tests.
var V = token.NewWithValue

// Create a new source position in tests.
var P = position.New

// Create a new span in tests.
var S = position.NewSpan

// Create a new source location in tests.
var L = position.NewLocation

var cmpOpts = []cmp.Option{
	cmp.AllowUnexported(
		ast.NodeBase{},
		token.Token{},
		bitfield.BitField8{},
		types.ConstantMap{},
		ast.VariableDeclarationNode{},
		ast.ValueDeclarationNode{},
		ast.ConstantDeclarationNode{},
		ast.PublicConstantNode{},
		ast.PrivateConstantNode{},
		ast.PublicIdentifierNode{},
		ast.PrivateIdentifierNode{},
		ast.ModuleDeclarationNode{},
		ast.MixinDeclarationNode{},
		ast.ClassDeclarationNode{},
		ast.RawStringLiteralNode{},
		ast.DoubleQuotedStringLiteralNode{},
		ast.RawCharLiteralNode{},
		ast.CharLiteralNode{},
		ast.SimpleSymbolLiteralNode{},
		ast.IntLiteralNode{},
		ast.Int64LiteralNode{},
		ast.Int32LiteralNode{},
		ast.Int16LiteralNode{},
		ast.Int8LiteralNode{},
		ast.UInt64LiteralNode{},
		ast.UInt32LiteralNode{},
		ast.UInt16LiteralNode{},
		ast.UInt8LiteralNode{},
		ast.FloatLiteralNode{},
		ast.Float64LiteralNode{},
		ast.Float32LiteralNode{},
		ast.BigFloatLiteralNode{},
		ast.UnionTypeNode{},
		ast.IntersectionTypeNode{},
		ast.MethodCallNode{},
		ast.ArrayTupleLiteralNode{},
		ast.ArrayListLiteralNode{},
		ast.HashRecordLiteralNode{},
		ast.ConstructorCallNode{},
		ast.ReceiverlessMethodCallNode{},
		types.Class{},
		types.Mixin{},
		types.MixinProxy{},
	),
	cmpopts.IgnoreUnexported(
		types.MethodMap{},
		types.TypeMap{},
	),
}

func checkerTest(tc testCase, t *testing.T) {
	t.Helper()
	_, err := CheckSource("<main>", tc.input, nil, false)

	if diff := cmp.Diff(tc.err, err, cmpOpts...); diff != "" {
		t.Log(pp.Sprint(err))
		t.Fatal(diff)
	}
}
