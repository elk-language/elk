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
	"github.com/elk-language/elk/value"
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
		"declaration with type lookup": {
			input: "val foo: Std::Int",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(16, 1, 17)),
								"Std::Int",
								globalEnv.StdSubtype("Int"),
							),
							nil,
							types.Void{},
						),
					),
				},
			),
		},
		"declaration with type lookup and error in the middle": {
			input: "val foo: Std::Foo::Bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(21, 1, 22)),
								"Std::Foo::Bar",
								types.Void{},
							),
							nil,
							types.Void{},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(14, 1, 15), P(16, 1, 17)), "undefined type `Std::Foo`"),
			},
		},
		"declaration with type lookup and error at the start": {
			input: "val foo: Foo::Bar::Baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(21, 1, 22)),
								"Foo::Bar::Baz",
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
		"declaration with absolute type lookup": {
			input: "val foo: ::Std::Int",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(18, 1, 19)),
								"Std::Int",
								globalEnv.StdSubtype("Int"),
							),
							nil,
							types.Void{},
						),
					),
				},
			),
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
		"constant lookup": {
			input: "Std::Int",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(7, 1, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(7, 1, 8)),
						ast.NewPublicConstantNode(
							S(P(0, 1, 1), P(7, 1, 8)),
							"Std::Int",
							globalEnv.StdConst("Int"),
						),
					),
				},
			),
		},
		"constant lookup with error in the middle": {
			input: "Std::Foo::Bar",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewPublicConstantNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							"Std::Foo::Bar",
							types.Void{},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(5, 1, 6), P(7, 1, 8)), "undefined constant `Std::Foo`"),
			},
		},
		"constant lookup with error at the start": {
			input: "Foo::Bar::Baz",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewPublicConstantNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							"Foo::Bar::Baz",
							types.Void{},
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(0, 1, 1), P(2, 1, 3)), "undefined constant `Foo`"),
			},
		},
		"absolute constant lookup": {
			input: "::Std::Int",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(9, 1, 10)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewPublicConstantNode(
							S(P(0, 1, 1), P(9, 1, 10)),
							"Std::Int",
							globalEnv.StdConst("Int"),
						),
					),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}

func TestModule(t *testing.T) {
	// globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"module with public constant": {
			input: `module Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							ast.NewPublicConstantNode(
								S(P(7, 1, 8), P(9, 1, 10)),
								"Foo",
								types.NewModule("Foo", nil, nil),
							),
							nil,
							types.NewModule("Foo", nil, nil),
						),
					),
				},
			),
		},
		"module with conflicting constant with Std": {
			input: `module Int; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							ast.NewPublicConstantNode(
								S(P(7, 1, 8), P(9, 1, 10)),
								"Int",
								types.NewModule("Int", nil, nil),
							),
							nil,
							types.NewModule("Int", nil, nil),
						),
					),
				},
			),
		},
		"module with private constant": {
			input: `module _Fo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							ast.NewPublicConstantNode(
								S(P(7, 1, 8), P(9, 1, 10)),
								"_Fo",
								types.NewModule("_Fo", nil, nil),
							),
							nil,
							types.NewModule("_Fo", nil, nil),
						),
					),
				},
			),
		},
		"module with simple constant lookup": {
			input: `module Std::Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							ast.NewPublicConstantNode(
								S(P(7, 1, 8), P(14, 1, 15)),
								"Std::Foo",
								types.NewModule("Std::Foo", nil, nil),
							),
							nil,
							types.NewModule("Std::Foo", nil, nil),
						),
					),
				},
			),
		},
		"module with non obvious constant lookup": {
			input: `module Int::Foo; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							ast.NewPublicConstantNode(
								S(P(7, 1, 8), P(14, 1, 15)),
								"Std::Int::Foo",
								types.NewModule("Std::Int::Foo", nil, nil),
							),
							nil,
							types.NewModule("Std::Int::Foo", nil, nil),
						),
					),
				},
			),
		},
		"module with undefined root constant": {
			input: `module Foo::Bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							ast.NewPublicConstantNode(
								S(P(7, 1, 8), P(14, 1, 15)),
								"Foo::Bar",
								types.NewModule("Foo::Bar", nil, nil),
							),
							nil,
							types.NewModule("Foo::Bar", nil, nil),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(7, 1, 8), P(9, 1, 10)), "undefined constant `Foo`"),
			},
		},
		"module with undefined constant in the middle": {
			input: `module Std::Foo::Bar; end`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(24, 1, 25)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(24, 1, 25)),
						ast.NewModuleDeclarationNode(
							S(P(0, 1, 1), P(24, 1, 25)),
							ast.NewPublicConstantNode(
								S(P(7, 1, 8), P(19, 1, 20)),
								"Std::Foo::Bar",
								types.NewModule("Std::Foo::Bar", nil, nil),
							),
							nil,
							types.NewModule("Std::Foo::Bar", nil, nil),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(12, 1, 13), P(14, 1, 15)), "undefined constant `Std::Foo`"),
			},
		},
		"nested modules": {
			input: `
				module Foo
					module Bar; end
				end
			`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(44, 4, 8)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(5, 2, 5), P(44, 4, 8)),
						ast.NewModuleDeclarationNode(
							S(P(5, 2, 5), P(43, 4, 7)),
							ast.NewPublicConstantNode(
								S(P(12, 2, 12), P(14, 2, 14)),
								"Foo",
								types.NewModule(
									"Foo",
									map[value.Symbol]types.Type{
										value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil),
									},
									map[value.Symbol]types.Type{
										value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil),
									},
								),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									S(P(21, 3, 6), P(36, 3, 21)),
									ast.NewModuleDeclarationNode(
										S(P(21, 3, 6), P(35, 3, 20)),
										ast.NewPublicConstantNode(
											S(P(28, 3, 13), P(30, 3, 15)),
											"Foo::Bar",
											types.NewModule("Foo::Bar", nil, nil),
										),
										nil,
										types.NewModule("Foo::Bar", nil, nil),
									),
								),
							},
							types.NewModule(
								"Foo",
								map[value.Symbol]types.Type{
									value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil),
								},
								map[value.Symbol]types.Type{
									value.ToSymbol("Bar"): types.NewModule("Foo::Bar", nil, nil),
								},
							),
						),
					),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t)
		})
	}
}
