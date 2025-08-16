package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/token"
)

func TestTypeof(t *testing.T) {
	tests := testTable{
		"with an argument": {
			input: "typeof 1",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewTypeofExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewIntLiteralNode(
								L(S(P(7, 1, 8), P(7, 1, 8))),
								"1",
							),
						),
					),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestCallableType(t *testing.T) {
	tests := testTable{
		"void callable without arguments": {
			input: "type ||",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							ast.NewCallableTypeNode(
								L(S(P(5, 1, 6), P(6, 1, 7))),
								nil,
								nil,
								nil,
								false,
							),
						),
					),
				},
			),
		},
		"callable with arguments, return type and throw type": {
			input: "type |a: String, b?: Int|: Int ! :dupa",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(37, 1, 38))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(37, 1, 38))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(37, 1, 38))),
							ast.NewCallableTypeNode(
								L(S(P(5, 1, 6), P(37, 1, 38))),
								[]ast.ParameterNode{
									ast.NewSignatureParameterNode(
										L(S(P(6, 1, 7), P(14, 1, 15))),
										ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(6, 1, 7))), "a"),
										ast.NewPublicConstantNode(L(S(P(9, 1, 10), P(14, 1, 15))), "String"),
										false,
										ast.NormalParameterKind,
									),
									ast.NewSignatureParameterNode(
										L(S(P(17, 1, 18), P(23, 1, 24))),
										ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "b"),
										ast.NewPublicConstantNode(L(S(P(21, 1, 22), P(23, 1, 24))), "Int"),
										true,
										ast.NormalParameterKind,
									),
								},
								ast.NewPublicConstantNode(L(S(P(27, 1, 28), P(29, 1, 30))), "Int"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(33, 1, 34), P(37, 1, 38))), "dupa"),
								false,
							),
						),
					),
				},
			),
		},
		"closure with arguments, return type and throw type": {
			input: "type %|a: String, b?: Int|: Int ! :dupa",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(38, 1, 39))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(38, 1, 39))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(38, 1, 39))),
							ast.NewCallableTypeNode(
								L(S(P(5, 1, 6), P(38, 1, 39))),
								[]ast.ParameterNode{
									ast.NewSignatureParameterNode(
										L(S(P(7, 1, 8), P(15, 1, 16))),
										ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(7, 1, 8))), "a"),
										ast.NewPublicConstantNode(L(S(P(10, 1, 11), P(15, 1, 16))), "String"),
										false,
										ast.NormalParameterKind,
									),
									ast.NewSignatureParameterNode(
										L(S(P(18, 1, 19), P(24, 1, 25))),
										ast.NewPublicIdentifierNode(L(S(P(18, 1, 19), P(18, 1, 19))), "b"),
										ast.NewPublicConstantNode(L(S(P(22, 1, 23), P(24, 1, 25))), "Int"),
										true,
										ast.NormalParameterKind,
									),
								},
								ast.NewPublicConstantNode(L(S(P(28, 1, 29), P(30, 1, 31))), "Int"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(34, 1, 35), P(38, 1, 39))), "dupa"),
								true,
							),
						),
					),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestConstantType(t *testing.T) {
	tests := testTable{
		"type can be a public constant": {
			input: "type String",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 1, 11))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 1, 11))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(10, 1, 11))),
							ast.NewPublicConstantNode(
								L(S(P(5, 1, 6), P(10, 1, 11))),
								"String",
							),
						),
					),
				},
			),
		},
		"type can be a generic public constant": {
			input: "type List[Int]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(13, 1, 14))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(13, 1, 14))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(13, 1, 14))),
							ast.NewGenericConstantNode(
								L(S(P(5, 1, 6), P(13, 1, 14))),
								ast.NewPublicConstantNode(
									L(S(P(5, 1, 6), P(8, 1, 9))),
									"List",
								),
								[]ast.TypeNode{
									ast.NewPublicConstantNode(
										L(S(P(10, 1, 11), P(12, 1, 13))),
										"Int",
									),
								},
							),
						),
					),
				},
			),
		},
		"type can be a private": {
			input: "type _FooBa",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 1, 11))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 1, 11))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(10, 1, 11))),
							ast.NewPrivateConstantNode(
								L(S(P(5, 1, 6), P(10, 1, 11))),
								"_FooBa",
							),
						),
					),
				},
			),
		},
		"type can be a constant lookup": {
			input: "type ::Foo::Bar",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(14, 1, 15))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(14, 1, 15))),
							ast.NewConstantLookupNode(
								L(S(P(5, 1, 6), P(14, 1, 15))),
								ast.NewConstantLookupNode(
									L(S(P(5, 1, 6), P(9, 1, 10))),
									nil,
									ast.NewPublicConstantNode(
										L(S(P(7, 1, 8), P(9, 1, 10))),
										"Foo",
									),
								),
								ast.NewPublicConstantNode(
									L(S(P(12, 1, 13), P(14, 1, 15))),
									"Bar",
								),
							),
						),
					),
				},
			),
		},
		"type can be a generic constant lookup": {
			input: "type ::Foo::Bar[Int, String]",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(27, 1, 28))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(27, 1, 28))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(27, 1, 28))),
							ast.NewGenericConstantNode(
								L(S(P(5, 1, 6), P(27, 1, 28))),
								ast.NewConstantLookupNode(
									L(S(P(5, 1, 6), P(14, 1, 15))),
									ast.NewConstantLookupNode(
										L(S(P(5, 1, 6), P(9, 1, 10))),
										nil,
										ast.NewPublicConstantNode(
											L(S(P(7, 1, 8), P(9, 1, 10))),
											"Foo",
										),
									),
									ast.NewPublicConstantNode(
										L(S(P(12, 1, 13), P(14, 1, 15))),
										"Bar",
									),
								),
								[]ast.TypeNode{
									ast.NewPublicConstantNode(
										L(S(P(16, 1, 17), P(18, 1, 19))),
										"Int",
									),
									ast.NewPublicConstantNode(
										L(S(P(21, 1, 22), P(26, 1, 27))),
										"String",
									),
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
			parserTest(tc, t)
		})
	}
}

func TestInstanceOfType(t *testing.T) {
	tests := testTable{
		"with a constant": {
			input: "type %String",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							ast.NewInstanceOfTypeNode(
								L(S(P(5, 1, 6), P(11, 1, 12))),
								ast.NewPublicConstantNode(
									L(S(P(6, 1, 7), P(11, 1, 12))),
									"String",
								),
							),
						),
					),
				},
			),
		},
		"with a private constant": {
			input: "type %_FooBa",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							ast.NewInstanceOfTypeNode(
								L(S(P(5, 1, 6), P(11, 1, 12))),
								ast.NewPrivateConstantNode(
									L(S(P(6, 1, 7), P(11, 1, 12))),
									"_FooBa",
								),
							),
						),
					),
				},
			),
		},
		"with constant lookup": {
			input: "type %::Foo::Bar",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							ast.NewInstanceOfTypeNode(
								L(S(P(5, 1, 6), P(15, 1, 16))),
								ast.NewConstantLookupNode(
									L(S(P(6, 1, 7), P(15, 1, 16))),
									ast.NewConstantLookupNode(
										L(S(P(6, 1, 7), P(10, 1, 11))),
										nil,
										ast.NewPublicConstantNode(
											L(S(P(8, 1, 9), P(10, 1, 11))),
											"Foo",
										),
									),
									ast.NewPublicConstantNode(
										L(S(P(13, 1, 14), P(15, 1, 16))),
										"Bar",
									),
								),
							),
						),
					),
				},
			),
		},
		"with literal": {
			input: "type %1",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							ast.NewInstanceOfTypeNode(
								L(S(P(5, 1, 6), P(6, 1, 7))),
								ast.NewIntLiteralNode(
									L(S(P(6, 1, 7), P(6, 1, 7))),
									"1",
								),
							),
						),
					),
				},
			),
		},
		"with an expression expression": {
			input: "type %(1)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewInstanceOfTypeNode(
								L(S(P(5, 1, 6), P(7, 1, 8))),
								ast.NewIntLiteralNode(
									L(S(P(7, 1, 8), P(7, 1, 8))),
									"1",
								),
							),
						),
					),
				},
			),
		},
		"can be nested": {
			input: "type %%1",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewInstanceOfTypeNode(
								L(S(P(5, 1, 6), P(7, 1, 8))),
								ast.NewInstanceOfTypeNode(
									L(S(P(6, 1, 7), P(7, 1, 8))),
									ast.NewIntLiteralNode(
										L(S(P(7, 1, 8), P(7, 1, 8))),
										"1",
									),
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
			parserTest(tc, t)
		})
	}
}
func TestNotType(t *testing.T) {
	tests := testTable{
		"type can be a not type with a constant": {
			input: "type ~String",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							ast.NewNotTypeNode(
								L(S(P(5, 1, 6), P(11, 1, 12))),
								ast.NewPublicConstantNode(
									L(S(P(6, 1, 7), P(11, 1, 12))),
									"String",
								),
							),
						),
					),
				},
			),
		},
		"type can be a not type with a private constant": {
			input: "type ~_FooBa",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							ast.NewNotTypeNode(
								L(S(P(5, 1, 6), P(11, 1, 12))),
								ast.NewPrivateConstantNode(
									L(S(P(6, 1, 7), P(11, 1, 12))),
									"_FooBa",
								),
							),
						),
					),
				},
			),
		},
		"type can be a not constant lookup": {
			input: "type ~::Foo::Bar",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							ast.NewNotTypeNode(
								L(S(P(5, 1, 6), P(15, 1, 16))),
								ast.NewConstantLookupNode(
									L(S(P(6, 1, 7), P(15, 1, 16))),
									ast.NewConstantLookupNode(
										L(S(P(6, 1, 7), P(10, 1, 11))),
										nil,
										ast.NewPublicConstantNode(
											L(S(P(8, 1, 9), P(10, 1, 11))),
											"Foo",
										),
									),
									ast.NewPublicConstantNode(
										L(S(P(13, 1, 14), P(15, 1, 16))),
										"Bar",
									),
								),
							),
						),
					),
				},
			),
		},
		"type can be a not literal": {
			input: "type ~1",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							ast.NewNotTypeNode(
								L(S(P(5, 1, 6), P(6, 1, 7))),
								ast.NewIntLiteralNode(
									L(S(P(6, 1, 7), P(6, 1, 7))),
									"1",
								),
							),
						),
					),
				},
			),
		},
		"type can be a not literal with expression": {
			input: "type ~(1)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewNotTypeNode(
								L(S(P(5, 1, 6), P(7, 1, 8))),
								ast.NewIntLiteralNode(
									L(S(P(7, 1, 8), P(7, 1, 8))),
									"1",
								),
							),
						),
					),
				},
			),
		},
		"can be nested": {
			input: "type ~~1",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewNotTypeNode(
								L(S(P(5, 1, 6), P(7, 1, 8))),
								ast.NewNotTypeNode(
									L(S(P(6, 1, 7), P(7, 1, 8))),
									ast.NewIntLiteralNode(
										L(S(P(7, 1, 8), P(7, 1, 8))),
										"1",
									),
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
			parserTest(tc, t)
		})
	}
}

func TestNilableType(t *testing.T) {
	tests := testTable{
		"can have unquote": {
			input: "type unquote(foo + 2)?",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(21, 1, 22))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(21, 1, 22))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(21, 1, 22))),
							ast.NewNilableTypeNode(
								L(S(P(5, 1, 6), P(21, 1, 22))),
								ast.NewUnquoteNode(
									L(S(P(5, 1, 6), P(20, 1, 21))),
									ast.UNQUOTE_TYPE_KIND,
									ast.NewBinaryExpressionNode(
										L(S(P(13, 1, 14), P(19, 1, 20))),
										T(L(S(P(17, 1, 18), P(17, 1, 18))), token.PLUS),
										ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(15, 1, 16))), "foo"),
										ast.NewIntLiteralNode(L(S(P(19, 1, 20), P(19, 1, 20))), "2"),
									),
								),
							),
						),
					),
				},
			),
		},
		"can have unquote_type": {
			input: "type unquote_type(foo + 2)?",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(26, 1, 27))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(26, 1, 27))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(26, 1, 27))),
							ast.NewNilableTypeNode(
								L(S(P(5, 1, 6), P(26, 1, 27))),
								ast.NewUnquoteNode(
									L(S(P(5, 1, 6), P(25, 1, 26))),
									ast.UNQUOTE_TYPE_KIND,
									ast.NewBinaryExpressionNode(
										L(S(P(18, 1, 19), P(24, 1, 25))),
										T(L(S(P(22, 1, 23), P(22, 1, 23))), token.PLUS),
										ast.NewPublicIdentifierNode(L(S(P(18, 1, 19), P(20, 1, 21))), "foo"),
										ast.NewIntLiteralNode(L(S(P(24, 1, 25), P(24, 1, 25))), "2"),
									),
								),
							),
						),
					),
				},
			),
		},
		"can have short unquote": {
			input: "type !{foo / 2}?",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							ast.NewNilableTypeNode(
								L(S(P(5, 1, 6), P(15, 1, 16))),
								ast.NewUnquoteNode(
									L(S(P(5, 1, 6), P(14, 1, 15))),
									ast.UNQUOTE_TYPE_KIND,
									ast.NewBinaryExpressionNode(
										L(S(P(7, 1, 8), P(13, 1, 14))),
										T(L(S(P(11, 1, 12), P(11, 1, 12))), token.SLASH),
										ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(9, 1, 10))), "foo"),
										ast.NewIntLiteralNode(L(S(P(13, 1, 14), P(13, 1, 14))), "2"),
									),
								),
							),
						),
					),
				},
			),
		},
		"type can be a nilable type with a constant": {
			input: "type String?",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							ast.NewNilableTypeNode(
								L(S(P(5, 1, 6), P(11, 1, 12))),
								ast.NewPublicConstantNode(
									L(S(P(5, 1, 6), P(10, 1, 11))),
									"String",
								),
							),
						),
					),
				},
			),
		},
		"type can be a nilable type with a private constant": {
			input: "type _FooBa?",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							ast.NewNilableTypeNode(
								L(S(P(5, 1, 6), P(11, 1, 12))),
								ast.NewPrivateConstantNode(
									L(S(P(5, 1, 6), P(10, 1, 11))),
									"_FooBa",
								),
							),
						),
					),
				},
			),
		},
		"type can be a nilable constant lookup": {
			input: "type ::Foo::Bar?",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							ast.NewNilableTypeNode(
								L(S(P(5, 1, 6), P(15, 1, 16))),
								ast.NewConstantLookupNode(
									L(S(P(5, 1, 6), P(14, 1, 15))),
									ast.NewConstantLookupNode(
										L(S(P(5, 1, 6), P(9, 1, 10))),
										nil,
										ast.NewPublicConstantNode(
											L(S(P(7, 1, 8), P(9, 1, 10))),
											"Foo",
										),
									),
									ast.NewPublicConstantNode(
										L(S(P(12, 1, 13), P(14, 1, 15))),
										"Bar",
									),
								),
							),
						),
					),
				},
			),
		},
		"type can be a nilable literal": {
			input: "type 1?",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							ast.NewNilableTypeNode(
								L(S(P(5, 1, 6), P(6, 1, 7))),
								ast.NewIntLiteralNode(
									L(S(P(5, 1, 6), P(5, 1, 6))),
									"1",
								),
							),
						),
					),
				},
			),
		},
		"type can be a nilable literal with expression": {
			input: "type (1)?",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewNilableTypeNode(
								L(S(P(6, 1, 7), P(8, 1, 9))),
								ast.NewIntLiteralNode(
									L(S(P(6, 1, 7), P(6, 1, 7))),
									"1",
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
			parserTest(tc, t)
		})
	}
}

func TestSingletonType(t *testing.T) {
	tests := testTable{
		"public constant": {
			input: "type &String",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							ast.NewSingletonTypeNode(
								L(S(P(5, 1, 6), P(11, 1, 12))),
								ast.NewPublicConstantNode(
									L(S(P(6, 1, 7), P(11, 1, 12))),
									"String",
								),
							),
						),
					),
				},
			),
		},
		"private constant": {
			input: "type &_FooBa",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							ast.NewSingletonTypeNode(
								L(S(P(5, 1, 6), P(11, 1, 12))),
								ast.NewPrivateConstantNode(
									L(S(P(6, 1, 7), P(11, 1, 12))),
									"_FooBa",
								),
							),
						),
					),
				},
			),
		},
		"type can be a nilable constant lookup": {
			input: "type &::Foo::Bar",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							ast.NewSingletonTypeNode(
								L(S(P(5, 1, 6), P(15, 1, 16))),
								ast.NewConstantLookupNode(
									L(S(P(6, 1, 7), P(15, 1, 16))),
									ast.NewConstantLookupNode(
										L(S(P(6, 1, 7), P(10, 1, 11))),
										nil,
										ast.NewPublicConstantNode(
											L(S(P(8, 1, 9), P(10, 1, 11))),
											"Foo",
										),
									),
									ast.NewPublicConstantNode(
										L(S(P(13, 1, 14), P(15, 1, 16))),
										"Bar",
									),
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
			parserTest(tc, t)
		})
	}
}

func TestBinaryType(t *testing.T) {
	tests := testTable{
		"union": {
			input: "type String | 4",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(14, 1, 15))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(14, 1, 15))),
							ast.NewBinaryTypeNode(
								L(S(P(5, 1, 6), P(14, 1, 15))),
								T(L(S(P(12, 1, 13), P(12, 1, 13))), token.OR),
								ast.NewPublicConstantNode(
									L(S(P(5, 1, 6), P(10, 1, 11))),
									"String",
								),
								ast.NewIntLiteralNode(
									L(S(P(14, 1, 15), P(14, 1, 15))),
									"4",
								),
							),
						),
					),
				},
			),
		},
		"intersection": {
			input: "type String & 4",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(14, 1, 15))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(14, 1, 15))),
							ast.NewBinaryTypeNode(
								L(S(P(5, 1, 6), P(14, 1, 15))),
								T(L(S(P(12, 1, 13), P(12, 1, 13))), token.AND),
								ast.NewPublicConstantNode(
									L(S(P(5, 1, 6), P(10, 1, 11))),
									"String",
								),
								ast.NewIntLiteralNode(
									L(S(P(14, 1, 15), P(14, 1, 15))),
									"4",
								),
							),
						),
					),
				},
			),
		},
		"difference": {
			input: "type String / 4",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(14, 1, 15))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(14, 1, 15))),
							ast.NewBinaryTypeNode(
								L(S(P(5, 1, 6), P(14, 1, 15))),
								T(L(S(P(12, 1, 13), P(12, 1, 13))), token.SLASH),
								ast.NewPublicConstantNode(
									L(S(P(5, 1, 6), P(10, 1, 11))),
									"String",
								),
								ast.NewIntLiteralNode(
									L(S(P(14, 1, 15), P(14, 1, 15))),
									"4",
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
			parserTest(tc, t)
		})
	}
}

func TestLiteralTypes(t *testing.T) {
	tests := testTable{
		"bool": {
			input: "type bool",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewBoolLiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
							),
						),
					),
				},
			),
		},
		"true": {
			input: "type true",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewTrueLiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
							),
						),
					),
				},
			),
		},
		"false": {
			input: "type false",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							ast.NewFalseLiteralNode(
								L(S(P(5, 1, 6), P(9, 1, 10))),
							),
						),
					),
				},
			),
		},
		"nil": {
			input: "type nil",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewNilLiteralNode(
								L(S(P(5, 1, 6), P(7, 1, 8))),
							),
						),
					),
				},
			),
		},
		"void": {
			input: "type void",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewVoidTypeNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
							),
						),
					),
				},
			),
		},
		"never": {
			input: "type never",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							ast.NewNeverTypeNode(
								L(S(P(5, 1, 6), P(9, 1, 10))),
							),
						),
					),
				},
			),
		},
		"any": {
			input: "type any",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewAnyTypeNode(
								L(S(P(5, 1, 6), P(7, 1, 8))),
							),
						),
					),
				},
			),
		},
		"raw char": {
			input: "type r`i`",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewRawCharLiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
								'i',
							),
						),
					),
				},
			),
		},
		"char": {
			input: "type `i`",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewCharLiteralNode(
								L(S(P(5, 1, 6), P(7, 1, 8))),
								'i',
							),
						),
					),
				},
			),
		},
		"raw string": {
			input: "type 'foo'",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							ast.NewRawStringLiteralNode(
								L(S(P(5, 1, 6), P(9, 1, 10))),
								"foo",
							),
						),
					),
				},
			),
		},
		"double quoted string": {
			input: `type "foo"`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							ast.NewDoubleQuotedStringLiteralNode(
								L(S(P(5, 1, 6), P(9, 1, 10))),
								"foo",
							),
						),
					),
				},
			),
		},
		"interpolated string": {
			input: `type "foo ${1}"`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(14, 1, 15))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(14, 1, 15))),
							ast.NewInterpolatedStringLiteralNode(
								L(S(P(5, 1, 6), P(14, 1, 15))),
								[]ast.StringLiteralContentNode{
									ast.NewStringLiteralContentSectionNode(
										L(S(P(6, 1, 7), P(9, 1, 10))),
										"foo ",
									),
									ast.NewStringInterpolationNode(
										L(S(P(10, 1, 11), P(13, 1, 14))),
										ast.NewIntLiteralNode(L(S(P(12, 1, 13), P(12, 1, 13))), "1"),
									),
								},
							),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(5, 1, 6), P(14, 1, 15))), "cannot interpolate strings in this context"),
			},
		},
		"simple symbol": {
			input: `type :foo`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewSimpleSymbolLiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
								"foo",
							),
						),
					),
				},
			),
		},
		"simple symbol with double quoted string": {
			input: `type :"foo"`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 1, 11))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 1, 11))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(10, 1, 11))),
							ast.NewSimpleSymbolLiteralNode(
								L(S(P(5, 1, 6), P(10, 1, 11))),
								"foo",
							),
						),
					),
				},
			),
		},
		"simple symbol with raw string": {
			input: `type :'foo'`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 1, 11))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 1, 11))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(10, 1, 11))),
							ast.NewSimpleSymbolLiteralNode(
								L(S(P(5, 1, 6), P(10, 1, 11))),
								"foo",
							),
						),
					),
				},
			),
		},
		"simple symbol with interpolated string": {
			input: `type :"foo ${1}"`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							ast.NewInterpolatedSymbolLiteralNode(
								L(S(P(5, 1, 6), P(15, 1, 16))),
								ast.NewInterpolatedStringLiteralNode(
									L(S(P(6, 1, 7), P(15, 1, 16))),
									[]ast.StringLiteralContentNode{
										ast.NewStringLiteralContentSectionNode(
											L(S(P(7, 1, 8), P(10, 1, 11))),
											"foo ",
										),
										ast.NewStringInterpolationNode(
											L(S(P(11, 1, 12), P(14, 1, 15))),
											ast.NewIntLiteralNode(L(S(P(13, 1, 14), P(13, 1, 14))), "1"),
										),
									},
								),
							),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(6, 1, 7), P(15, 1, 16))), "cannot interpolate strings in this context"),
			},
		},
		"int": {
			input: "type 1234",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewIntLiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
								"1234",
							),
						),
					),
				},
			),
		},
		"negative int": {
			input: "type -1234",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							ast.NewUnaryTypeNode(
								L(S(P(5, 1, 6), P(9, 1, 10))),
								T(L(S(P(5, 1, 6), P(5, 1, 6))), token.MINUS),
								ast.NewIntLiteralNode(
									L(S(P(6, 1, 7), P(9, 1, 10))),
									"1234",
								),
							),
						),
					),
				},
			),
		},
		"positive int": {
			input: "type +1234",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							ast.NewUnaryTypeNode(
								L(S(P(5, 1, 6), P(9, 1, 10))),
								T(L(S(P(5, 1, 6), P(5, 1, 6))), token.PLUS),
								ast.NewIntLiteralNode(
									L(S(P(6, 1, 7), P(9, 1, 10))),
									"1234",
								),
							),
						),
					),
				},
			),
		},
		"int64": {
			input: "type 1i64",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewInt64LiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
								"1",
							),
						),
					),
				},
			),
		},
		"int32": {
			input: "type 1i32",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewInt32LiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
								"1",
							),
						),
					),
				},
			),
		},
		"int16": {
			input: "type 1i16",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewInt16LiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
								"1",
							),
						),
					),
				},
			),
		},
		"int8": {
			input: "type 1i8",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewInt8LiteralNode(
								L(S(P(5, 1, 6), P(7, 1, 8))),
								"1",
							),
						),
					),
				},
			),
		},
		"uint64": {
			input: "type 1u64",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewUInt64LiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
								"1",
							),
						),
					),
				},
			),
		},
		"uint32": {
			input: "type 1u32",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewUInt32LiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
								"1",
							),
						),
					),
				},
			),
		},
		"uint16": {
			input: "type 1u16",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewUInt16LiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
								"1",
							),
						),
					),
				},
			),
		},
		"uint8": {
			input: "type 1u8",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(7, 1, 8))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(7, 1, 8))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(7, 1, 8))),
							ast.NewUInt8LiteralNode(
								L(S(P(5, 1, 6), P(7, 1, 8))),
								"1",
							),
						),
					),
				},
			),
		},
		"float": {
			input: "type 1.56",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewFloatLiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
								"1.56",
							),
						),
					),
				},
			),
		},
		"float64": {
			input: "type 1f64",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewFloat64LiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
								"1",
							),
						),
					),
				},
			),
		},
		"float32": {
			input: "type 1f32",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewFloat32LiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
								"1",
							),
						),
					),
				},
			),
		},
		"big float": {
			input: "type 12bf",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewTypeExpressionNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewBigFloatLiteralNode(
								L(S(P(5, 1, 6), P(8, 1, 9))),
								"12",
							),
						),
					),
				},
			),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}
