package checker

import (
	"testing"

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
