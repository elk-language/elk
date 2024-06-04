package checker

import (
	"testing"

	"github.com/elk-language/elk/position/error"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/ast"
	"github.com/elk-language/elk/value/symbol"
)

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
								globalEnv.StdSubtype(symbol.Int),
							),
							ast.NewIntLiteralNode(
								S(P(15, 1, 16), P(15, 1, 16)),
								"5",
								types.NewIntLiteral("5"),
							),
							globalEnv.StdSubtype(symbol.Int),
						),
					),
				},
			),
		},
		"accept variable declaration with inference": {
			input: "var foo = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewIntLiteralNode(
								S(P(10, 1, 11), P(10, 1, 11)),
								"5",
								types.NewIntLiteral("5"),
							),
							globalEnv.StdSubtype(symbol.Int),
						),
					),
				},
			),
		},
		"cannot declare variable with type void": {
			before: "def bar; end",
			input:  "var foo = bar()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewReceiverlessMethodCallNode(
								S(P(10, 1, 11), P(14, 1, 15)),
								"bar",
								nil,
								types.Void{},
							),
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(10, 1, 11), P(14, 1, 15)), "cannot declare variable `foo` with type `void`"),
			},
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
								globalEnv.StdSubtype(symbol.Int),
							),
							ast.NewFloatLiteralNode(
								S(P(15, 1, 16), P(17, 1, 18)),
								"5.2",
								types.NewFloatLiteral("5.2"),
							),
							globalEnv.StdSubtype(symbol.Int),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(15, 1, 16), P(17, 1, 18)), "type `Std::Float(5.2)` cannot be assigned to type `Std::Int`"),
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
								globalEnv.StdSubtype(symbol.Int),
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
			err: error.ErrorList{
				error.NewError(L("<main>", P(9, 1, 10), P(11, 1, 12)), "undefined type `Foo`"),
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
			err: error.ErrorList{
				error.NewError(L("<main>", P(0, 1, 1), P(6, 1, 7)), "cannot declare a variable without a type `foo`"),
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
								globalEnv.StdSubtype(symbol.Int),
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
								globalEnv.StdSubtype(symbol.String),
							),
							nil,
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(14, 1, 15), P(28, 1, 29)), "cannot redeclare local `foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, false)
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
								globalEnv.StdSubtype(symbol.Int),
							),
							ast.NewIntLiteralNode(
								S(P(15, 1, 16), P(15, 1, 16)),
								"5",
								types.NewIntLiteral("5"),
							),
							globalEnv.StdSubtype(symbol.Int),
						),
					),
				},
			),
		},
		"accept variable declaration with inference": {
			input: "val foo = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(10, 1, 11)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(10, 1, 11)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(10, 1, 11)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewIntLiteralNode(
								S(P(10, 1, 11), P(10, 1, 11)),
								"5",
								types.NewIntLiteral("5"),
							),
							types.NewIntLiteral("5"),
						),
					),
				},
			),
		},
		"cannot declare value with type void": {
			before: "def bar; end",
			input:  "val foo = bar()",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewValueDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewReceiverlessMethodCallNode(
								S(P(10, 1, 11), P(14, 1, 15)),
								"bar",
								nil,
								types.Void{},
							),
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(10, 1, 11), P(14, 1, 15)), "cannot declare value `foo` with type `void`"),
			},
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
								globalEnv.StdSubtype(symbol.Int),
							),
							ast.NewFloatLiteralNode(
								S(P(15, 1, 16), P(17, 1, 18)),
								"5.2",
								types.NewFloatLiteral("5.2"),
							),
							globalEnv.StdSubtype(symbol.Int),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(15, 1, 16), P(17, 1, 18)), "type `Std::Float(5.2)` cannot be assigned to type `Std::Int`"),
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
								globalEnv.StdSubtype(symbol.Int),
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
			err: error.ErrorList{
				error.NewError(L("<main>", P(9, 1, 10), P(11, 1, 12)), "undefined type `Foo`"),
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
			err: error.ErrorList{
				error.NewError(L("<main>", P(0, 1, 1), P(6, 1, 7)), "cannot declare a value without a type `foo`"),
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
								globalEnv.StdSubtype(symbol.Int),
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
								globalEnv.StdSubtype(symbol.String),
							),
							nil,
							types.Void{},
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(14, 1, 15), P(28, 1, 29)), "cannot redeclare local `foo`"),
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
								globalEnv.StdSubtype(symbol.Int),
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
			err: error.ErrorList{
				error.NewError(L("<main>", P(14, 1, 15), P(16, 1, 17)), "undefined type `Std::Foo`"),
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
			err: error.ErrorList{
				error.NewError(L("<main>", P(9, 1, 10), P(11, 1, 12)), "undefined type `Foo`"),
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
								globalEnv.StdSubtype(symbol.Int),
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
			checkerTest(tc, t, false)
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
								globalEnv.StdSubtype(symbol.Int),
							),
							ast.NewIntLiteralNode(
								S(P(15, 1, 16), P(15, 1, 16)),
								"5",
								types.NewIntLiteral("5"),
							),
							globalEnv.StdSubtype(symbol.Int),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(18, 1, 19), P(20, 1, 21)),
						ast.NewPublicIdentifierNode(
							S(P(18, 1, 19), P(20, 1, 21)),
							"foo",
							globalEnv.StdSubtype(symbol.Int),
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
								globalEnv.StdSubtype(symbol.Int),
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
							globalEnv.StdSubtype(symbol.Int),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(14, 1, 15), P(16, 1, 17)), "cannot access uninitialised local `foo`"),
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
								globalEnv.StdSubtype(symbol.Int),
							),
							ast.NewIntLiteralNode(
								S(P(15, 1, 16), P(15, 1, 16)),
								"5",
								types.NewIntLiteral("5"),
							),
							globalEnv.StdSubtype(symbol.Int),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(18, 1, 19), P(20, 1, 21)),
						ast.NewPublicIdentifierNode(
							S(P(18, 1, 19), P(20, 1, 21)),
							"foo",
							globalEnv.StdSubtype(symbol.Int),
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
								globalEnv.StdSubtype(symbol.Int),
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
							globalEnv.StdSubtype(symbol.Int),
						),
					),
				},
			),
			err: error.ErrorList{
				error.NewError(L("<main>", P(14, 1, 15), P(16, 1, 17)), "cannot access uninitialised local `foo`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, false)
		})
	}
}
