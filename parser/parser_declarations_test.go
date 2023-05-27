package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/token"
)

func TestVariableDeclaration(t *testing.T) {
	tests := testTable{
		"is valid without a type or initialiser": {
			input: "var foo",
			want: ast.NewProgramNode(
				P(0, 7, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 7, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 7, 1, 1),
							V(P(4, 3, 1, 5), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a private identifier as the variable name": {
			input: "var _foo",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 8, 1, 1),
							V(P(4, 4, 1, 5), token.PRIVATE_IDENTIFIER, "_foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have an instance variable as the variable name": {
			input: "var @foo",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 8, 1, 1),
							V(P(4, 4, 1, 5), token.INSTANCE_VARIABLE, "foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can't have a constant as the variable name": {
			input: "var Foo",
			want: ast.NewProgramNode(
				P(0, 7, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(4, 3, 1, 5),
						ast.NewInvalidNode(
							P(4, 3, 1, 5),
							V(P(4, 3, 1, 5), token.PUBLIC_CONSTANT, "Foo"),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(4, 3, 1, 5), "unexpected PUBLIC_CONSTANT, expected an identifier as the name of the declared variable"),
			},
		},
		"can have an initialiser without a type": {
			input: "var foo = 5",
			want: ast.NewProgramNode(
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 11, 1, 1),
							V(P(4, 3, 1, 5), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewIntLiteralNode(P(10, 1, 1, 11), V(P(10, 1, 1, 11), token.DEC_INT, "5")),
						),
					),
				},
			),
		},
		"can have an initialiser with a type": {
			input: "var foo: Int = 5",
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 16, 1, 1),
							V(P(4, 3, 1, 5), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(P(9, 3, 1, 10), "Int"),
							ast.NewIntLiteralNode(P(15, 1, 1, 16), V(P(15, 1, 1, 16), token.DEC_INT, "5")),
						),
					),
				},
			),
		},
		"can have a type": {
			input: "var foo: Int",
			want: ast.NewProgramNode(
				P(0, 12, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 12, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 12, 1, 1),
							V(P(4, 3, 1, 5), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(P(9, 3, 1, 10), "Int"),
							nil,
						),
					),
				},
			),
		},
		"can have a nilable type": {
			input: "var foo: Int?",
			want: ast.NewProgramNode(
				P(0, 13, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 13, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 13, 1, 1),
							V(P(4, 3, 1, 5), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewNilableTypeNode(
								P(9, 4, 1, 10),
								ast.NewPublicConstantNode(P(9, 3, 1, 10), "Int"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a union type": {
			input: "var foo: Int | String",
			want: ast.NewProgramNode(
				P(0, 21, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 21, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 21, 1, 1),
							V(P(4, 3, 1, 5), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewBinaryTypeExpressionNode(
								P(9, 12, 1, 10),
								T(P(13, 1, 1, 14), token.OR),
								ast.NewPublicConstantNode(P(9, 3, 1, 10), "Int"),
								ast.NewPublicConstantNode(P(15, 6, 1, 16), "String"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nested union type": {
			input: "var foo: Int | String | Symbol",
			want: ast.NewProgramNode(
				P(0, 30, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 30, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 30, 1, 1),
							V(P(4, 3, 1, 5), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewBinaryTypeExpressionNode(
								P(9, 21, 1, 10),
								T(P(22, 1, 1, 23), token.OR),
								ast.NewBinaryTypeExpressionNode(
									P(9, 12, 1, 10),
									T(P(13, 1, 1, 14), token.OR),
									ast.NewPublicConstantNode(P(9, 3, 1, 10), "Int"),
									ast.NewPublicConstantNode(P(15, 6, 1, 16), "String"),
								),
								ast.NewPublicConstantNode(P(24, 6, 1, 25), "Symbol"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nilable union type": {
			input: "var foo: (Int | String)?",
			want: ast.NewProgramNode(
				P(0, 24, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 24, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 24, 1, 1),
							V(P(4, 3, 1, 5), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewNilableTypeNode(
								P(10, 14, 1, 11),
								ast.NewBinaryTypeExpressionNode(
									P(10, 12, 1, 11),
									T(P(14, 1, 1, 15), token.OR),
									ast.NewPublicConstantNode(P(10, 3, 1, 11), "Int"),
									ast.NewPublicConstantNode(P(16, 6, 1, 17), "String"),
								),
							),
							nil,
						),
					),
				},
			),
		},
		"can have an intersection type": {
			input: "var foo: Int & String",
			want: ast.NewProgramNode(
				P(0, 21, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 21, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 21, 1, 1),
							V(P(4, 3, 1, 5), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewBinaryTypeExpressionNode(
								P(9, 12, 1, 10),
								T(P(13, 1, 1, 14), token.AND),
								ast.NewPublicConstantNode(P(9, 3, 1, 10), "Int"),
								ast.NewPublicConstantNode(P(15, 6, 1, 16), "String"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nested intersection type": {
			input: "var foo: Int & String & Symbol",
			want: ast.NewProgramNode(
				P(0, 30, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 30, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 30, 1, 1),
							V(P(4, 3, 1, 5), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewBinaryTypeExpressionNode(
								P(9, 21, 1, 10),
								T(P(22, 1, 1, 23), token.AND),
								ast.NewBinaryTypeExpressionNode(
									P(9, 12, 1, 10),
									T(P(13, 1, 1, 14), token.AND),
									ast.NewPublicConstantNode(P(9, 3, 1, 10), "Int"),
									ast.NewPublicConstantNode(P(15, 6, 1, 16), "String"),
								),
								ast.NewPublicConstantNode(P(24, 6, 1, 25), "Symbol"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nilable intersection type": {
			input: "var foo: (Int & String)?",
			want: ast.NewProgramNode(
				P(0, 24, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 24, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 24, 1, 1),
							V(P(4, 3, 1, 5), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewNilableTypeNode(
								P(10, 14, 1, 11),
								ast.NewBinaryTypeExpressionNode(
									P(10, 12, 1, 11),
									T(P(14, 1, 1, 15), token.AND),
									ast.NewPublicConstantNode(P(10, 3, 1, 11), "Int"),
									ast.NewPublicConstantNode(P(16, 6, 1, 17), "String"),
								),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a generic type": {
			input: "var foo: Std::Map[Std::Symbol, List[String]]",
			want: ast.NewProgramNode(
				P(0, 44, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 44, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 44, 1, 1),
							V(P(4, 3, 1, 5), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewGenericConstantNode(
								P(9, 35, 1, 10),
								ast.NewConstantLookupNode(
									P(9, 8, 1, 10),
									ast.NewPublicConstantNode(P(9, 3, 1, 10), "Std"),
									ast.NewPublicConstantNode(P(14, 3, 1, 15), "Map"),
								),
								[]ast.ComplexConstantNode{
									ast.NewConstantLookupNode(
										P(18, 11, 1, 19),
										ast.NewPublicConstantNode(P(18, 3, 1, 19), "Std"),
										ast.NewPublicConstantNode(P(23, 6, 1, 24), "Symbol"),
									),
									ast.NewGenericConstantNode(
										P(31, 12, 1, 32),
										ast.NewPublicConstantNode(P(31, 4, 1, 32), "List"),
										[]ast.ComplexConstantNode{
											ast.NewPublicConstantNode(P(36, 6, 1, 37), "String"),
										},
									),
								},
							),
							nil,
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

func TestClassDeclaration(t *testing.T) {
	tests := testTable{
		"can be anonymous": {
			input: `class; end`,
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 10, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 10, 1, 1),
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can be a part of an expression": {
			input: `foo = class; end`,
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 16, 1, 1),
							T(P(4, 1, 1, 5), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewClassDeclarationNode(
								P(6, 10, 1, 7),
								nil,
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can be anonymous with a superclass": {
			input: `class < Foo; end`,
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 16, 1, 1),
							nil,
							nil,
							ast.NewPublicConstantNode(P(8, 3, 1, 9), "Foo"),
							nil,
						),
					),
				},
			),
		},
		"can have type variables": {
			input: `class Foo[V, +T, -Z]; end`,
			want: ast.NewProgramNode(
				P(0, 25, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 25, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 25, 1, 1),
							ast.NewPublicConstantNode(P(6, 3, 1, 7), "Foo"),
							[]ast.TypeVariableNode{
								ast.NewVariantTypeVariableNode(
									P(10, 1, 1, 11),
									ast.INVARIANT,
									"V",
									nil,
								),
								ast.NewVariantTypeVariableNode(
									P(13, 2, 1, 14),
									ast.COVARIANT,
									"T",
									nil,
								),
								ast.NewVariantTypeVariableNode(
									P(17, 2, 1, 18),
									ast.CONTRAVARIANT,
									"Z",
									nil,
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have type variables with upper bounds": {
			input: `class Foo[V < Std::String, +T < Foo, -Z < _Bar]; end`,
			want: ast.NewProgramNode(
				P(0, 52, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 52, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 52, 1, 1),
							ast.NewPublicConstantNode(P(6, 3, 1, 7), "Foo"),
							[]ast.TypeVariableNode{
								ast.NewVariantTypeVariableNode(
									P(10, 15, 1, 11),
									ast.INVARIANT,
									"V",
									ast.NewConstantLookupNode(
										P(14, 11, 1, 15),
										ast.NewPublicConstantNode(P(14, 3, 1, 15), "Std"),
										ast.NewPublicConstantNode(P(19, 6, 1, 20), "String"),
									),
								),
								ast.NewVariantTypeVariableNode(
									P(27, 8, 1, 28),
									ast.COVARIANT,
									"T",
									ast.NewPublicConstantNode(P(32, 3, 1, 33), "Foo"),
								),
								ast.NewVariantTypeVariableNode(
									P(37, 9, 1, 38),
									ast.CONTRAVARIANT,
									"Z",
									ast.NewPrivateConstantNode(P(42, 4, 1, 43), "_Bar"),
								),
							},
							nil,
							nil,
						),
					),
				},
			),
		},
		"can't have an empty type variable list": {
			input: `class Foo[]; end`,
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 16, 1, 1),
							ast.NewPublicConstantNode(P(6, 3, 1, 7), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(10, 1, 1, 11), "unexpected ], expected a list of type variables"),
			},
		},
		"can have a public constant as a name": {
			input: `class Foo; end`,
			want: ast.NewProgramNode(
				P(0, 14, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 14, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 14, 1, 1),
							ast.NewPublicConstantNode(P(6, 3, 1, 7), "Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a private constant as a name": {
			input: `class _Foo; end`,
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 15, 1, 1),
							ast.NewPrivateConstantNode(P(6, 4, 1, 7), "_Foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a constant lookup as a name": {
			input: `class Foo::Bar; end`,
			want: ast.NewProgramNode(
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 19, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 19, 1, 1),
							ast.NewConstantLookupNode(
								P(6, 8, 1, 7),
								ast.NewPublicConstantNode(P(6, 3, 1, 7), "Foo"),
								ast.NewPublicConstantNode(P(11, 3, 1, 12), "Bar"),
							),
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can't have an identifier as a name": {
			input: `class foo; end`,
			want: ast.NewProgramNode(
				P(0, 14, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 14, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 14, 1, 1),
							ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "foo"),
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(6, 3, 1, 7), "invalid class name, expected a constant"),
			},
		},
		"can have a public constant as a superclass": {
			input: `class Foo < Bar; end`,
			want: ast.NewProgramNode(
				P(0, 20, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 20, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 20, 1, 1),
							ast.NewPublicConstantNode(P(6, 3, 1, 7), "Foo"),
							nil,
							ast.NewPublicConstantNode(P(12, 3, 1, 13), "Bar"),
							nil,
						),
					),
				},
			),
		},
		"can have a private constant as a superclass": {
			input: `class Foo < _Bar; end`,
			want: ast.NewProgramNode(
				P(0, 21, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 21, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 21, 1, 1),
							ast.NewPublicConstantNode(P(6, 3, 1, 7), "Foo"),
							nil,
							ast.NewPrivateConstantNode(P(12, 4, 1, 13), "_Bar"),
							nil,
						),
					),
				},
			),
		},
		"can have a constant lookup as a superclass": {
			input: `class Foo < Bar::Baz; end`,
			want: ast.NewProgramNode(
				P(0, 25, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 25, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 25, 1, 1),
							ast.NewPublicConstantNode(P(6, 3, 1, 7), "Foo"),
							nil,
							ast.NewConstantLookupNode(
								P(12, 8, 1, 13),
								ast.NewPublicConstantNode(P(12, 3, 1, 13), "Bar"),
								ast.NewPublicConstantNode(P(17, 3, 1, 18), "Baz"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a generic constant as a superclass": {
			input: `class Foo < Std::Map[Symbol, String]; end`,
			want: ast.NewProgramNode(
				P(0, 41, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 41, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 41, 1, 1),
							ast.NewPublicConstantNode(P(6, 3, 1, 7), "Foo"),
							nil,
							ast.NewGenericConstantNode(
								P(12, 24, 1, 13),
								ast.NewConstantLookupNode(
									P(12, 8, 1, 13),
									ast.NewPublicConstantNode(P(12, 3, 1, 13), "Std"),
									ast.NewPublicConstantNode(P(17, 3, 1, 18), "Map"),
								),
								[]ast.ComplexConstantNode{
									ast.NewPublicConstantNode(P(21, 6, 1, 22), "Symbol"),
									ast.NewPublicConstantNode(P(29, 6, 1, 30), "String"),
								},
							),
							nil,
						),
					),
				},
			),
		},
		"can't have an identifier as a superclass": {
			input: `class Foo < bar; end`,
			want: ast.NewProgramNode(
				P(0, 20, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 20, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 20, 1, 1),
							ast.NewPublicConstantNode(P(6, 3, 1, 7), "Foo"),
							nil,
							ast.NewInvalidNode(P(12, 3, 1, 13), V(P(12, 3, 1, 13), token.PUBLIC_IDENTIFIER, "bar")),
							nil,
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(12, 3, 1, 13), "unexpected PUBLIC_IDENTIFIER, expected a constant"),
			},
		},
		"can have a multiline body": {
			input: `class Foo
	foo = 2
	nil
end`,
			want: ast.NewProgramNode(
				P(0, 27, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 27, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 27, 1, 1),
							ast.NewPublicConstantNode(P(6, 3, 1, 7), "Foo"),
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(11, 8, 2, 2),
									ast.NewAssignmentExpressionNode(
										P(11, 7, 2, 2),
										T(P(15, 1, 2, 6), token.EQUAL_OP),
										ast.NewPublicIdentifierNode(P(11, 3, 2, 2), "foo"),
										ast.NewIntLiteralNode(P(17, 1, 2, 8), V(P(17, 1, 2, 8), token.DEC_INT, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(20, 4, 3, 2),
									ast.NewNilLiteralNode(P(20, 3, 3, 2)),
								),
							},
						),
					),
				},
			),
		},
		"can have a single line body with then": {
			input: `class Foo then .1 * .2`,
			want: ast.NewProgramNode(
				P(0, 22, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 22, 1, 1),
						ast.NewClassDeclarationNode(
							P(0, 22, 1, 1),
							ast.NewPublicConstantNode(P(6, 3, 1, 7), "Foo"),
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(15, 7, 1, 16),
									ast.NewBinaryExpressionNode(
										P(15, 7, 1, 16),
										T(P(18, 1, 1, 19), token.STAR),
										ast.NewFloatLiteralNode(P(15, 2, 1, 16), "0.1"),
										ast.NewFloatLiteralNode(P(20, 2, 1, 21), "0.2"),
									),
								),
							},
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

func TestMethodDeclaration(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "bar = def foo; end",
			want: ast.NewProgramNode(
				P(0, 18, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 18, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 18, 1, 1),
							T(P(4, 1, 1, 5), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "bar"),
							ast.NewMethodDeclarationNode(
								P(6, 12, 1, 7),
								"foo",
								nil,
								nil,
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can have a public identifier as a name": {
			input: "def foo; end",
			want: ast.NewProgramNode(
				P(0, 12, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 12, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 12, 1, 1),
							"foo",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a private identifier as a name": {
			input: "def _foo; end",
			want: ast.NewProgramNode(
				P(0, 13, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 13, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 13, 1, 1),
							"_foo",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a keyword as a name": {
			input: "def class; end",
			want: ast.NewProgramNode(
				P(0, 14, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 14, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 14, 1, 1),
							"class",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have an overridable operator as a name": {
			input: "def +; end",
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 10, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 10, 1, 1),
							"+",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can't have a public constant as a name": {
			input: "def Foo; end",
			want: ast.NewProgramNode(
				P(0, 12, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 12, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 12, 1, 1),
							"Foo",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(4, 3, 1, 5), "unexpected PUBLIC_CONSTANT, expected a method name (identifier, overridable operator)"),
			},
		},
		"can't have a non overridable operator as a name": {
			input: "def &&; end",
			want: ast.NewProgramNode(
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 11, 1, 1),
							"&&",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(4, 2, 1, 5), "unexpected &&, expected a method name (identifier, overridable operator)"),
			},
		},
		"can't have a private constant as a name": {
			input: "def _Foo; end",
			want: ast.NewProgramNode(
				P(0, 13, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 13, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 13, 1, 1),
							"_Foo",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(4, 4, 1, 5), "unexpected PRIVATE_CONSTANT, expected a method name (identifier, overridable operator)"),
			},
		},
		"can have an empty argument list": {
			input: "def foo(); end",
			want: ast.NewProgramNode(
				P(0, 14, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 14, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 14, 1, 1),
							"foo",
							nil,
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a return type and omit arguments": {
			input: "def foo: String?; end",
			want: ast.NewProgramNode(
				P(0, 21, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 21, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 21, 1, 1),
							"foo",
							nil,
							ast.NewNilableTypeNode(
								P(9, 7, 1, 10),
								ast.NewPublicConstantNode(P(9, 6, 1, 10), "String"),
							),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a throw type and omit arguments": {
			input: "def foo! NoMethodError | TypeError; end",
			want: ast.NewProgramNode(
				P(0, 39, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 39, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 39, 1, 1),
							"foo",
							nil,
							nil,
							ast.NewBinaryTypeExpressionNode(
								P(9, 25, 1, 10),
								T(P(23, 1, 1, 24), token.OR),
								ast.NewPublicConstantNode(P(9, 13, 1, 10), "NoMethodError"),
								ast.NewPublicConstantNode(P(25, 9, 1, 26), "TypeError"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a return and throw type and omit arguments": {
			input: "def foo : String? ! NoMethodError | TypeError; end",
			want: ast.NewProgramNode(
				P(0, 50, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 50, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 50, 1, 1),
							"foo",
							nil,
							ast.NewNilableTypeNode(
								P(10, 7, 1, 11),
								ast.NewPublicConstantNode(P(10, 6, 1, 11), "String"),
							),
							ast.NewBinaryTypeExpressionNode(
								P(20, 25, 1, 21),
								T(P(34, 1, 1, 35), token.OR),
								ast.NewPublicConstantNode(P(20, 13, 1, 21), "NoMethodError"),
								ast.NewPublicConstantNode(P(36, 9, 1, 37), "TypeError"),
							),
							nil,
						),
					),
				},
			),
		},
		"can have arguments": {
			input: "def foo(a, b); end",
			want: ast.NewProgramNode(
				P(0, 18, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 18, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 18, 1, 1),
							"foo",
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									P(8, 1, 1, 9),
									"a",
									nil,
									nil,
								),
								ast.NewFormalParameterNode(
									P(11, 1, 1, 12),
									"b",
									nil,
									nil,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have arguments with types": {
			input: "def foo(a: Int, b: String?); end",
			want: ast.NewProgramNode(
				P(0, 32, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 32, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 32, 1, 1),
							"foo",
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									P(8, 6, 1, 9),
									"a",
									ast.NewPublicConstantNode(P(11, 3, 1, 12), "Int"),
									nil,
								),
								ast.NewFormalParameterNode(
									P(16, 10, 1, 17),
									"b",
									ast.NewNilableTypeNode(
										P(19, 7, 1, 20),
										ast.NewPublicConstantNode(P(19, 6, 1, 20), "String"),
									),
									nil,
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have arguments with initialisers": {
			input: "def foo(a = 32, b: String = 'foo'); end",
			want: ast.NewProgramNode(
				P(0, 39, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 39, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 39, 1, 1),
							"foo",
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									P(8, 6, 1, 9),
									"a",
									nil,
									ast.NewIntLiteralNode(P(12, 2, 1, 13), V(P(12, 2, 1, 13), token.DEC_INT, "32")),
								),
								ast.NewFormalParameterNode(
									P(16, 17, 1, 17),
									"b",
									ast.NewPublicConstantNode(P(19, 6, 1, 20), "String"),
									ast.NewRawStringLiteralNode(P(28, 5, 1, 29), "foo"),
								),
							},
							nil,
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have a multiline body": {
			input: `def foo
  a := .5
  a += .7
end`,
			want: ast.NewProgramNode(
				P(0, 31, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 31, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 31, 1, 1),
							"foo",
							nil,
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(10, 8, 2, 3),
									ast.NewAssignmentExpressionNode(
										P(10, 7, 2, 3),
										T(P(12, 2, 2, 5), token.COLON_EQUAL),
										ast.NewPublicIdentifierNode(P(10, 1, 2, 3), "a"),
										ast.NewFloatLiteralNode(P(15, 2, 2, 8), "0.5"),
									),
								),
								ast.NewExpressionStatementNode(
									P(20, 8, 3, 3),
									ast.NewAssignmentExpressionNode(
										P(20, 7, 3, 3),
										T(P(22, 2, 3, 5), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(P(20, 1, 3, 3), "a"),
										ast.NewFloatLiteralNode(P(25, 2, 3, 8), "0.7"),
									),
								),
							},
						),
					),
				},
			),
		},
		"can be single line with then": {
			input: `def foo then .3 + .4`,
			want: ast.NewProgramNode(
				P(0, 20, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 20, 1, 1),
						ast.NewMethodDeclarationNode(
							P(0, 20, 1, 1),
							"foo",
							nil,
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(13, 7, 1, 14),
									ast.NewBinaryExpressionNode(
										P(13, 7, 1, 14),
										T(P(16, 1, 1, 17), token.PLUS),
										ast.NewFloatLiteralNode(P(13, 2, 1, 14), "0.3"),
										ast.NewFloatLiteralNode(P(18, 2, 1, 19), "0.4"),
									),
								),
							},
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
