// Package checker implements the Elk type checker
package checker

import (
	"testing"

	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/ast"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/k0kubun/pp"
)

// Represents a single checker test case.
type testCase struct {
	before string
	input  string
	want   *ast.ProgramNode
	err    errors.ErrorList
}

// Type of the checker test table.
type testTable map[string]testCase

// Represents a single checker test case.
type simplifiedTestCase struct {
	input string
	err   errors.ErrorList
}

// Type of the checker test table.
type simplifiedTestTable map[string]simplifiedTestCase

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
	),
}

var ignoreConstantTypesOpts = []cmp.Option{
	cmp.AllowUnexported(
		ast.NodeBase{},
		token.Token{},
		bitfield.BitField8{},
		types.ConstantMap{},
		ast.VariableDeclarationNode{},
		ast.ValueDeclarationNode{},
		ast.ConstantDeclarationNode{},
		ast.PublicIdentifierNode{},
		ast.PrivateIdentifierNode{},
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
	),
	cmpopts.IgnoreUnexported(
		ast.ModuleDeclarationNode{},
		ast.MixinDeclarationNode{},
		ast.ClassDeclarationNode{},
		ast.PublicConstantNode{},
		ast.PrivateConstantNode{},
		ast.ConstructorCallNode{},
	),
}

// Function which powers all type checker tests.
// Inspects if the produced typed AST matches the expected one.
func checkerTest(tc testCase, t *testing.T, ignoreConstantTypes bool) {
	t.Helper()
	checker := New()
	if tc.before != "" {
		_, err := checker.CheckSource("<before>", tc.before)
		if err != nil {
			t.Fatal(err)
		}
	}
	got, err := checker.CheckSource("<main>", tc.input)
	var opts []cmp.Option
	if ignoreConstantTypes {
		opts = ignoreConstantTypesOpts
	} else {
		opts = cmpOpts
	}

	if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
		t.Log(pp.Sprint(got))
		t.Fatal(diff)
	}

	if diff := cmp.Diff(tc.err, err); diff != "" {
		t.Log(pp.Sprint(err))
		t.Fatal(diff)
	}
}

func simplifiedCheckerTest(tc simplifiedTestCase, t *testing.T) {
	t.Helper()
	_, err := CheckSource("<main>", tc.input, nil, false)

	if diff := cmp.Diff(tc.err, err); diff != "" {
		t.Log(pp.Sprint(err))
		t.Fatal(diff)
	}
}
