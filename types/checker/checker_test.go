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
	"github.com/k0kubun/pp"
)

// Represents a single parser test case.
type testCase struct {
	input string
	want  *ast.ProgramNode
	err   errors.ErrorList
}

// Type of the parser test table.
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

// Function which powers all type checker tests.
// Inspects if the produced typed AST matches the expected one.
func checkerTest(tc testCase, t *testing.T) {
	t.Helper()
	got, err := CheckSource("<main>", tc.input, nil, false)

	opts := []cmp.Option{
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
		),
	}
	if diff := cmp.Diff(tc.want, got, opts...); diff != "" {
		pp.Println(got)
		t.Fatal(diff)
	}

	if diff := cmp.Diff(tc.err, err); diff != "" {
		t.Fatal(diff)
	}
}
