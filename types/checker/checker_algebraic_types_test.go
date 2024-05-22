package checker

import (
	"testing"

	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/ast"
)

func TestNilableSubtype(t *testing.T) {
	tests := simplifiedTestTable{
		"assign String to nilable String": {
			input: `
				var a = "foo"
				var b: String? = a
			`,
		},
		"assign nil to nilable String": {
			input: `
				var a = nil
				var b: String? = a
			`,
		},
		"assign Int to nilable String": {
			input: `
				var a = 3
				var b: String? = a
			`,
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(36, 3, 22), P(36, 3, 22)), "type `Std::Int` cannot be assigned to type `Std::String?`"),
			},
		},
		"assign nilable String to union type with String and nil": {
			input: `
				var a: String? = "foo"
				var b: String | Float | nil = a
			`,
		},
		"assign nilable String to union type without nil": {
			input: `
				var a: String? = "foo"
				var b: String | Float = a
			`,
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(56, 3, 29), P(56, 3, 29)), "type `Std::String?` cannot be assigned to type `Std::String | Std::Float`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			simplifiedCheckerTest(tc, t)
		})
	}
}

func TestUnionTypeSubtype(t *testing.T) {
	tests := simplifiedTestTable{
		"assign Int to union type with Int": {
			input: `
				var a = 3
				var b: String | Int = a
			`,
		},
		"assign Int to union type without Int": {
			input: `
				var a = 3
				var b: String | Float = a
			`,
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(43, 3, 29), P(43, 3, 29)), "type `Std::Int` cannot be assigned to type `Std::String | Std::Float`"),
			},
		},
		"assign union type to non union type": {
			input: `
				var a = 3
				var b: String | Float = a
			`,
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(43, 3, 29), P(43, 3, 29)), "type `Std::Int` cannot be assigned to type `Std::String | Std::Float`"),
			},
		},
		"assign union type to more general union type": {
			input: `
				var a: String | Int = 3
				var b: Int | Float | String = a
			`,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			simplifiedCheckerTest(tc, t)
		})
	}
}

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
			checkerTest(tc, t, false)
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
			checkerTest(tc, t, false)
		})
	}
}
