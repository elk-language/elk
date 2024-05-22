package checker

import (
	"testing"

	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/ast"
	"github.com/elk-language/elk/value/symbol"
)

func TestConstantAccess(t *testing.T) {
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
							globalEnv.StdConst(symbol.Int),
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
							globalEnv.StdConst(symbol.Int),
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
							globalEnv.StdConst(symbol.Int),
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

func TestConstantDeclarations(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"infer type from initialiser": {
			input: "const Foo = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							V(S(P(6, 1, 7), P(8, 1, 9)), token.PUBLIC_CONSTANT, "Foo"),
							nil,
							ast.NewIntLiteralNode(S(P(12, 1, 13), P(12, 1, 13)), "5", types.NewIntLiteral("5")),
							types.NewIntLiteral("5"),
						),
					),
				},
			),
		},
		"declare with explicit type": {
			input: "const Foo: Int = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							V(S(P(6, 1, 7), P(8, 1, 9)), token.PUBLIC_CONSTANT, "Foo"),
							ast.NewPublicConstantNode(
								S(P(11, 1, 12), P(13, 1, 14)),
								"Int",
								globalEnv.StdSubtype(symbol.Int),
							),
							ast.NewIntLiteralNode(S(P(17, 1, 18), P(17, 1, 18)), "5", types.NewIntLiteral("5")),
							globalEnv.StdSubtype(symbol.Int),
						),
					),
				},
			),
		},
		"declare with incorrect explicit type": {
			input: "const Foo: String = 5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewConstantDeclarationNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							V(S(P(6, 1, 7), P(8, 1, 9)), token.PUBLIC_CONSTANT, "Foo"),
							ast.NewPublicConstantNode(
								S(P(11, 1, 12), P(16, 1, 17)),
								"String",
								globalEnv.StdSubtype(symbol.String),
							),
							ast.NewIntLiteralNode(S(P(20, 1, 21), P(20, 1, 21)), "5", types.NewIntLiteral("5")),
							globalEnv.StdSubtype(symbol.String),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(20, 1, 21), P(20, 1, 21)), "type `Std::Int(5)` cannot be assigned to type `Std::String`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, false)
		})
	}
}
