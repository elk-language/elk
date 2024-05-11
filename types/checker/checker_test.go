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
			ast.PublicConstantNode{},
			ast.PrivateConstantNode{},
			ast.ValueDeclarationNode{},
			ast.PublicIdentifierNode{},
			ast.PrivateIdentifierNode{},
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

func TestVariableDeclaration(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"accept variable declaration with matching initializer and type": {
			input: "var foo: Int = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Int",
								globalEnv.StdSubtype("Int"),
							),
							ast.NewIntLiteralNode(
								S(P(15, 1, 16), P(15, 1, 16)),
								"5",
							),
							globalEnv.StdSubtype("Int"),
						),
					),
				},
			),
		},
		"reject variable declaration without matching initializer and type": {
			input: "var foo: Int = 5.2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Int",
								globalEnv.StdSubtype("Int"),
							),
							ast.NewFloatLiteralNode(
								S(P(15, 1, 16), P(17, 1, 18)),
								"5.2",
							),
							globalEnv.StdSubtype("Int"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(15, 1, 16), P(17, 1, 18)), "type `Std::Float` cannot be assigned to type `Std::Int`"),
			},
		},
		"accept variable declaration without initializer": {
			input: "var foo: Int",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Int",
								globalEnv.StdSubtype("Int"),
							),
							nil,
							types.Void{},
						),
					),
				},
			),
		},
		"reject variable declaration with invalid type": {
			input: "var foo: Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Foo",
								types.Void{},
							),
							nil,
							types.Void{},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(9, 1, 10), P(11, 1, 12)), "undefined type `Foo`"),
			},
		},
		"reject variable declaration without initializer and type": {
			input: "var foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							nil,
							types.Void{},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(0, 1, 1), P(6, 1, 7)), "cannot declare a variable without a type `foo`"),
			},
		},
		"reject redeclared variable": {
			input: "var foo: Int; var foo: String",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 1, 29)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Int",
								globalEnv.StdSubtype("Int"),
							),
							nil,
							types.Void{},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(14, 1, 15), P(28, 1, 29)),
						ast.NewVariableDeclarationNode(
							S(P(14, 1, 15), P(28, 1, 29)),
							V(S(P(18, 1, 19), P(20, 1, 21)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(23, 1, 24), P(28, 1, 29)),
								"String",
								globalEnv.StdSubtype("String"),
							),
							nil,
							types.Void{},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(14, 1, 15), P(28, 1, 29)), "cannot redeclare local `foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestValueDeclaration(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"accept value declaration with matching initializer and type": {
			input: "val foo: Int = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Int",
								globalEnv.StdSubtype("Int"),
							),
							ast.NewIntLiteralNode(
								S(P(15, 1, 16), P(15, 1, 16)),
								"5",
							),
							globalEnv.StdSubtype("Int"),
						),
					),
				},
			),
		},
		"reject value declaration without matching initializer and type": {
			input: "val foo: Int = 5.2",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Int",
								globalEnv.StdSubtype("Int"),
							),
							ast.NewFloatLiteralNode(
								S(P(15, 1, 16), P(17, 1, 18)),
								"5.2",
							),
							globalEnv.StdSubtype("Int"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(15, 1, 16), P(17, 1, 18)), "type `Std::Float` cannot be assigned to type `Std::Int`"),
			},
		},
		"accept value declaration without initializer": {
			input: "val foo: Int",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Int",
								globalEnv.StdSubtype("Int"),
							),
							nil,
							types.Void{},
						),
					),
				},
			),
		},
		"reject value declaration with invalid type": {
			input: "val foo: Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(11, 1, 12)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(11, 1, 12)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Foo",
								types.Void{},
							),
							nil,
							types.Void{},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(9, 1, 10), P(11, 1, 12)), "undefined type `Foo`"),
			},
		},
		"reject value declaration without initializer and type": {
			input: "val foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(6, 1, 7)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(6, 1, 7)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(6, 1, 7)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							nil,
							types.Void{},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(0, 1, 1), P(6, 1, 7)), "cannot declare a value without a type `foo`"),
			},
		},
		"reject redeclared value": {
			input: "val foo: Int; val foo: String",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(28, 1, 29)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Int",
								globalEnv.StdSubtype("Int"),
							),
							nil,
							types.Void{},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(14, 1, 15), P(28, 1, 29)),
						ast.NewValueDeclarationNode(
							S(P(14, 1, 15), P(28, 1, 29)),
							V(S(P(18, 1, 19), P(20, 1, 21)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(23, 1, 24), P(28, 1, 29)),
								"String",
								globalEnv.StdSubtype("String"),
							),
							nil,
							types.Void{},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(14, 1, 15), P(28, 1, 29)), "cannot redeclare local `foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestLocalAccess(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"access initialised variable": {
			input: "var foo: Int = 5; foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Int",
								globalEnv.StdSubtype("Int"),
							),
							ast.NewIntLiteralNode(
								S(P(15, 1, 16), P(15, 1, 16)),
								"5",
							),
							globalEnv.StdSubtype("Int"),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(18, 1, 19), P(20, 1, 21)),
						ast.NewPublicIdentifierNode(
							S(P(18, 1, 19), P(20, 1, 21)),
							"foo",
							globalEnv.StdSubtype("Int"),
						),
					),
				},
			),
		},
		"access uninitialised variable": {
			input: "var foo: Int; foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Int",
								globalEnv.StdSubtype("Int"),
							),
							nil,
							types.Void{},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(14, 1, 15), P(16, 1, 17)),
						ast.NewPublicIdentifierNode(
							S(P(14, 1, 15), P(16, 1, 17)),
							"foo",
							globalEnv.StdSubtype("Int"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(14, 1, 15), P(16, 1, 17)), "cannot access uninitialised local `foo`"),
			},
		},
		"access initialised value": {
			input: "val foo: Int = 5; foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Int",
								globalEnv.StdSubtype("Int"),
							),
							ast.NewIntLiteralNode(
								S(P(15, 1, 16), P(15, 1, 16)),
								"5",
							),
							globalEnv.StdSubtype("Int"),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(18, 1, 19), P(20, 1, 21)),
						ast.NewPublicIdentifierNode(
							S(P(18, 1, 19), P(20, 1, 21)),
							"foo",
							globalEnv.StdSubtype("Int"),
						),
					),
				},
			),
		},
		"access uninitialised value": {
			input: "val foo: Int; foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(11, 1, 12)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Int",
								globalEnv.StdSubtype("Int"),
							),
							nil,
							types.Void{},
						),
					),
					ast.NewExpressionStatementNode(
						S(P(14, 1, 15), P(16, 1, 17)),
						ast.NewPublicIdentifierNode(
							S(P(14, 1, 15), P(16, 1, 17)),
							"foo",
							globalEnv.StdSubtype("Int"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(14, 1, 15), P(16, 1, 17)), "cannot access uninitialised local `foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestConstants(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"access class constant": {
			input: "Int",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(2, 1, 3)),
						ast.NewPublicConstantNode(
							S(P(0, 1, 1), P(2, 1, 3)),
							"Std::Int",
							globalEnv.StdConst("Int"),
						),
					),
				},
			),
		},
		"access module constant": {
			input: "Std",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(2, 1, 3)),
						ast.NewPublicConstantNode(
							S(P(0, 1, 1), P(2, 1, 3)),
							"Std",
							globalEnv.Std(),
						),
					),
				},
			),
		},
		"access undefined constant": {
			input: "Foo",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(2, 1, 3)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(2, 1, 3)),
						ast.NewPublicConstantNode(
							S(P(0, 1, 1), P(2, 1, 3)),
							"Foo",
							types.Void{},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(0, 1, 1), P(2, 1, 3)), "undefined constant `Foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
