package checker

import (
	"testing"

	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/ast"
)

func TestUnionType(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"construct a simple union type": {
			input: "type Int | String",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewTypeExpressionNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							ast.NewUnionTypeNode(
								S(P(5, 1, 6), P(16, 1, 17)),
								[]ast.TypeNode{
									ast.NewPublicConstantNode(
										S(P(5, 1, 6), P(7, 1, 8)),
										"Int",
										globalEnv.StdSubtypeString("Int"),
									),
									ast.NewPublicConstantNode(
										S(P(11, 1, 12), P(16, 1, 17)),
										"String",
										globalEnv.StdSubtypeString("String"),
									),
								},
								types.NewUnion(
									globalEnv.StdSubtypeString("Int"),
									globalEnv.StdSubtypeString("String"),
								),
							),
						),
					),
				},
			),
		},
		"flatten nested union types": {
			input: "type Int | String | Float | nil",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(30, 1, 31)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(30, 1, 31)),
						ast.NewTypeExpressionNode(
							S(P(0, 1, 1), P(30, 1, 31)),
							ast.NewUnionTypeNode(
								S(P(5, 1, 6), P(30, 1, 31)),
								[]ast.TypeNode{
									ast.NewPublicConstantNode(
										S(P(5, 1, 6), P(7, 1, 8)),
										"Int",
										globalEnv.StdSubtypeString("Int"),
									),
									ast.NewPublicConstantNode(
										S(P(11, 1, 12), P(16, 1, 17)),
										"String",
										globalEnv.StdSubtypeString("String"),
									),
									ast.NewPublicConstantNode(
										S(P(20, 1, 21), P(24, 1, 25)),
										"Float",
										globalEnv.StdSubtypeString("Float"),
									),
									ast.NewNilLiteralNode(S(P(28, 1, 29), P(30, 1, 31))),
								},
								types.NewUnion(
									globalEnv.StdSubtypeString("Int"),
									globalEnv.StdSubtypeString("String"),
									globalEnv.StdSubtypeString("Float"),
									globalEnv.StdSubtypeString("Nil"),
								),
							),
						),
					),
				},
			),
		},
		"assign Int to union type with Int": {
			input: "var a = 3; var b: String | Int = a",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(33, 1, 34)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							V(S(P(4, 1, 5), P(4, 1, 5)), token.PUBLIC_IDENTIFIER, "a"),
							nil,
							ast.NewIntLiteralNode(
								S(P(8, 1, 9), P(8, 1, 9)),
								"3",
								types.NewIntLiteral("3"),
							),
							globalEnv.StdSubtypeString("Int"),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(11, 1, 12), P(33, 1, 34)),
						ast.NewVariableDeclarationNode(
							S(P(11, 1, 12), P(33, 1, 34)),
							V(S(P(15, 1, 16), P(15, 1, 16)), token.PUBLIC_IDENTIFIER, "b"),
							ast.NewUnionTypeNode(
								S(P(18, 1, 19), P(29, 1, 30)),
								[]ast.TypeNode{
									ast.NewPublicConstantNode(
										S(P(18, 1, 19), P(23, 1, 24)),
										"String",
										globalEnv.StdSubtypeString("String"),
									),
									ast.NewPublicConstantNode(
										S(P(27, 1, 28), P(29, 1, 30)),
										"Int",
										globalEnv.StdSubtypeString("Int"),
									),
								},
								types.NewUnion(
									globalEnv.StdSubtypeString("String"),
									globalEnv.StdSubtypeString("Int"),
								),
							),
							ast.NewPublicIdentifierNode(
								S(P(33, 1, 34), P(33, 1, 34)),
								"a",
								globalEnv.StdSubtypeString("Int"),
							),
							types.NewUnion(
								globalEnv.StdSubtypeString("String"),
								globalEnv.StdSubtypeString("Int"),
							),
						),
					),
				},
			),
		},
		"assign Int to union type without Int": {
			input: "var a = 3; var b: String | Float = a",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(35, 1, 36)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(9, 1, 10)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(8, 1, 9)),
							V(S(P(4, 1, 5), P(4, 1, 5)), token.PUBLIC_IDENTIFIER, "a"),
							nil,
							ast.NewIntLiteralNode(
								S(P(8, 1, 9), P(8, 1, 9)),
								"3",
								types.NewIntLiteral("3"),
							),
							globalEnv.StdSubtypeString("Int"),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(11, 1, 12), P(35, 1, 36)),
						ast.NewVariableDeclarationNode(
							S(P(11, 1, 12), P(35, 1, 36)),
							V(S(P(15, 1, 16), P(15, 1, 16)), token.PUBLIC_IDENTIFIER, "b"),
							ast.NewUnionTypeNode(
								S(P(18, 1, 19), P(31, 1, 32)),
								[]ast.TypeNode{
									ast.NewPublicConstantNode(
										S(P(18, 1, 19), P(23, 1, 24)),
										"String",
										globalEnv.StdSubtypeString("String"),
									),
									ast.NewPublicConstantNode(
										S(P(27, 1, 28), P(31, 1, 32)),
										"Float",
										globalEnv.StdSubtypeString("Float"),
									),
								},
								types.NewUnion(
									globalEnv.StdSubtypeString("String"),
									globalEnv.StdSubtypeString("Float"),
								),
							),
							ast.NewPublicIdentifierNode(
								S(P(35, 1, 36), P(35, 1, 36)),
								"a",
								globalEnv.StdSubtypeString("Int"),
							),
							types.NewUnion(
								globalEnv.StdSubtypeString("String"),
								globalEnv.StdSubtypeString("Float"),
							),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(35, 1, 36), P(35, 1, 36)), "type `Std::Int` cannot be assigned to type `Std::String | Std::Float`"),
			},
		},
		"assign union type to non union type": {
			input: "var a: String | Int = 3; var b: Int = a",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(38, 1, 39)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(23, 1, 24)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							V(S(P(4, 1, 5), P(4, 1, 5)), token.PUBLIC_IDENTIFIER, "a"),
							ast.NewUnionTypeNode(
								S(P(7, 1, 8), P(18, 1, 19)),
								[]ast.TypeNode{
									ast.NewPublicConstantNode(
										S(P(7, 1, 8), P(12, 1, 13)),
										"String",
										globalEnv.StdSubtypeString("String"),
									),
									ast.NewPublicConstantNode(
										S(P(16, 1, 17), P(18, 1, 19)),
										"Int",
										globalEnv.StdSubtypeString("Int"),
									),
								},
								types.NewUnion(
									globalEnv.StdSubtypeString("String"),
									globalEnv.StdSubtypeString("Int"),
								),
							),
							ast.NewIntLiteralNode(
								S(P(22, 1, 23), P(22, 1, 23)),
								"3",
								types.NewIntLiteral("3"),
							),
							types.NewUnion(
								globalEnv.StdSubtypeString("String"),
								globalEnv.StdSubtypeString("Int"),
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(25, 1, 26), P(38, 1, 39)),
						ast.NewVariableDeclarationNode(
							S(P(25, 1, 26), P(38, 1, 39)),
							V(S(P(29, 1, 30), P(29, 1, 30)), token.PUBLIC_IDENTIFIER, "b"),
							ast.NewPublicConstantNode(
								S(P(32, 1, 33), P(34, 1, 35)),
								"Int",
								globalEnv.StdSubtypeString("Int"),
							),
							ast.NewPublicIdentifierNode(
								S(P(38, 1, 39), P(38, 1, 39)),
								"a",
								types.NewUnion(
									globalEnv.StdSubtypeString("String"),
									globalEnv.StdSubtypeString("Int"),
								),
							),
							globalEnv.StdSubtypeString("Int"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(38, 1, 39), P(38, 1, 39)), "type `Std::String | Std::Int` cannot be assigned to type `Std::Int`"),
			},
		},
		"assign union type to more general union type": {
			input: "var a: String | Int = 3; var b: Int | Float | String = a",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(55, 1, 56)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(23, 1, 24)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							V(S(P(4, 1, 5), P(4, 1, 5)), token.PUBLIC_IDENTIFIER, "a"),
							ast.NewUnionTypeNode(
								S(P(7, 1, 8), P(18, 1, 19)),
								[]ast.TypeNode{
									ast.NewPublicConstantNode(
										S(P(7, 1, 8), P(12, 1, 13)),
										"String",
										globalEnv.StdSubtypeString("String"),
									),
									ast.NewPublicConstantNode(
										S(P(16, 1, 17), P(18, 1, 19)),
										"Int",
										globalEnv.StdSubtypeString("Int"),
									),
								},
								types.NewUnion(
									globalEnv.StdSubtypeString("String"),
									globalEnv.StdSubtypeString("Int"),
								),
							),
							ast.NewIntLiteralNode(
								S(P(22, 1, 23), P(22, 1, 23)),
								"3",
								types.NewIntLiteral("3"),
							),
							types.NewUnion(
								globalEnv.StdSubtypeString("String"),
								globalEnv.StdSubtypeString("Int"),
							),
						),
					),
					ast.NewExpressionStatementNode(
						S(P(25, 1, 26), P(55, 1, 56)),
						ast.NewVariableDeclarationNode(
							S(P(25, 1, 26), P(55, 1, 56)),
							V(S(P(29, 1, 30), P(29, 1, 30)), token.PUBLIC_IDENTIFIER, "b"),
							ast.NewUnionTypeNode(
								S(P(32, 1, 33), P(51, 1, 52)),
								[]ast.TypeNode{
									ast.NewPublicConstantNode(
										S(P(32, 1, 33), P(34, 1, 35)),
										"Int",
										globalEnv.StdSubtypeString("Int"),
									),
									ast.NewPublicConstantNode(
										S(P(38, 1, 39), P(42, 1, 43)),
										"Float",
										globalEnv.StdSubtypeString("Float"),
									),
									ast.NewPublicConstantNode(
										S(P(46, 1, 47), P(51, 1, 52)),
										"String",
										globalEnv.StdSubtypeString("String"),
									),
								},
								types.NewUnion(
									globalEnv.StdSubtypeString("Int"),
									globalEnv.StdSubtypeString("Float"),
									globalEnv.StdSubtypeString("String"),
								),
							),
							ast.NewPublicIdentifierNode(
								S(P(55, 1, 56), P(55, 1, 56)),
								"a",
								types.NewUnion(
									globalEnv.StdSubtypeString("String"),
									globalEnv.StdSubtypeString("Int"),
								),
							),
							types.NewUnion(
								globalEnv.StdSubtypeString("Int"),
								globalEnv.StdSubtypeString("Float"),
								globalEnv.StdSubtypeString("String"),
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

func TestIntersectionType(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"construct a simple intersection type": {
			input: "type Int & String",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(16, 1, 17)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(16, 1, 17)),
						ast.NewTypeExpressionNode(
							S(P(0, 1, 1), P(16, 1, 17)),
							ast.NewIntersectionTypeNode(
								S(P(5, 1, 6), P(16, 1, 17)),
								[]ast.TypeNode{
									ast.NewPublicConstantNode(
										S(P(5, 1, 6), P(7, 1, 8)),
										"Int",
										globalEnv.StdSubtypeString("Int"),
									),
									ast.NewPublicConstantNode(
										S(P(11, 1, 12), P(16, 1, 17)),
										"String",
										globalEnv.StdSubtypeString("String"),
									),
								},
								types.NewIntersection(
									globalEnv.StdSubtypeString("Int"),
									globalEnv.StdSubtypeString("String"),
								),
							),
						),
					),
				},
			),
		},
		"flatten nested intersection types": {
			input: "type Int & String & Float & nil",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(30, 1, 31)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(30, 1, 31)),
						ast.NewTypeExpressionNode(
							S(P(0, 1, 1), P(30, 1, 31)),
							ast.NewIntersectionTypeNode(
								S(P(5, 1, 6), P(30, 1, 31)),
								[]ast.TypeNode{
									ast.NewPublicConstantNode(
										S(P(5, 1, 6), P(7, 1, 8)),
										"Int",
										globalEnv.StdSubtypeString("Int"),
									),
									ast.NewPublicConstantNode(
										S(P(11, 1, 12), P(16, 1, 17)),
										"String",
										globalEnv.StdSubtypeString("String"),
									),
									ast.NewPublicConstantNode(
										S(P(20, 1, 21), P(24, 1, 25)),
										"Float",
										globalEnv.StdSubtypeString("Float"),
									),
									ast.NewNilLiteralNode(S(P(28, 1, 29), P(30, 1, 31))),
								},
								types.NewIntersection(
									globalEnv.StdSubtypeString("Int"),
									globalEnv.StdSubtypeString("String"),
									globalEnv.StdSubtypeString("Float"),
									globalEnv.StdSubtypeString("Nil"),
								),
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
