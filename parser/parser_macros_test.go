package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/token"
)

func TestMacroDefinition(t *testing.T) {
	tests := testTable{
		"cannot be a part of an expression": {
			input: "bar = macro foo; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(19, 1, 20))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(19, 1, 20))),
						ast.NewAssignmentExpressionNode(
							L(S(P(0, 1, 1), P(19, 1, 20))),
							T(L(S(P(4, 1, 5), P(4, 1, 5))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "bar"),
							ast.NewMacroDefinitionNode(
								L(S(P(6, 1, 7), P(19, 1, 20))),
								"",
								false,
								ast.NewPublicIdentifierNode(L(S(P(12, 1, 13), P(14, 1, 15))), "foo"),
								nil,
								nil,
							),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(6, 1, 7), P(19, 1, 20))), "macro definitions cannot appear in expressions"),
			},
		},
		"can have a public identifier as a name": {
			input: "macro foo; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(13, 1, 14))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(13, 1, 14))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(13, 1, 14))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot be a generator": {
			input: "macro* foo; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(5, 1, 6), P(10, 1, 11))),
						ast.NewInvalidNode(
							L(S(P(5, 1, 6), P(5, 1, 6))),
							T(L(S(P(5, 1, 6), P(5, 1, 6))), token.STAR),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(12, 1, 13), P(14, 1, 15))),
						ast.NewInvalidNode(
							L(S(P(12, 1, 13), P(14, 1, 15))),
							T(L(S(P(12, 1, 13), P(14, 1, 15))), token.END),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(5, 1, 6), P(5, 1, 6))), "unexpected *, expected a macro name (public identifier, keyword)"),
				diagnostic.NewFailure(L(S(P(12, 1, 13), P(14, 1, 15))), "unexpected end, expected an expression"),
			},
		},
		"cannot have type variables": {
			input: "macro foo[V]; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(9, 1, 10), P(12, 1, 13))),
									ast.NewArrayListLiteralNode(
										L(S(P(9, 1, 10), P(11, 1, 12))),
										[]ast.ExpressionNode{
											ast.NewPublicConstantNode(L(S(P(10, 1, 11), P(10, 1, 11))), "V"),
										},
										nil,
									),
								),
							},
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(9, 1, 10), P(9, 1, 10))), "unexpected [, expected a statement separator `\\n`, `;`"),
			},
		},
		"can be sealed": {
			input: "sealed macro foo; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(20, 1, 21))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(20, 1, 21))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(20, 1, 21))),
							"",
							true,
							ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(15, 1, 16))), "foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot repeat sealed": {
			input: "sealed sealed macro foo; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(27, 1, 28))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(27, 1, 28))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(27, 1, 28))),
							"",
							true,
							ast.NewPublicIdentifierNode(L(S(P(20, 1, 21), P(22, 1, 23))), "foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(0, 1, 1), P(5, 1, 6))), "the sealed modifier can only be attached once"),
			},
		},
		"cannot be async": {
			input: "async macro foo; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(19, 1, 20))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(6, 1, 7), P(19, 1, 20))),
						ast.NewMacroDefinitionNode(
							L(S(P(6, 1, 7), P(19, 1, 20))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(12, 1, 13), P(14, 1, 15))), "foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(6, 1, 7), P(19, 1, 20))), "the async modifier can only be attached to methods"),
			},
		},
		"cannot repeat async": {
			input: "async async macro foo; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(25, 1, 26))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(12, 1, 13), P(25, 1, 26))),
						ast.NewMacroDefinitionNode(
							L(S(P(12, 1, 13), P(25, 1, 26))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(18, 1, 19), P(20, 1, 21))), "foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(12, 1, 13), P(25, 1, 26))), "the async modifier can only be attached to methods"),
				diagnostic.NewFailure(L(S(P(12, 1, 13), P(25, 1, 26))), "the async modifier can only be attached to methods"),
			},
		},
		"cannot be abstract": {
			input: "abstract macro foo; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(22, 1, 23))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(9, 1, 10), P(22, 1, 23))),
						ast.NewMacroDefinitionNode(
							L(S(P(9, 1, 10), P(22, 1, 23))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(17, 1, 18))), "foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(9, 1, 10), P(22, 1, 23))), "the abstract modifier can only be attached to classes, mixins and methods"),
			},
		},
		"cannot have a setter as a name": {
			input: "macro foo=(v); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(17, 1, 18))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(17, 1, 18))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(17, 1, 18))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(9, 1, 10), P(13, 1, 14))),
									ast.NewInvalidNode(
										L(S(P(9, 1, 10), P(9, 1, 10))),
										T(L(S(P(9, 1, 10), P(9, 1, 10))), token.EQUAL_OP),
									),
								),
							},
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(9, 1, 10), P(9, 1, 10))), "unexpected =, expected a statement separator `\\n`, `;`"),
				diagnostic.NewFailure(L(S(P(9, 1, 10), P(9, 1, 10))), "unexpected =, expected an expression"),
			},
		},
		"cannot have a private identifier as a name": {
			input: "macro _foo; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(6, 1, 7), P(10, 1, 11))),
						ast.NewInvalidNode(
							L(S(P(6, 1, 7), P(9, 1, 10))),
							V(L(S(P(6, 1, 7), P(9, 1, 10))), token.PRIVATE_IDENTIFIER, "_foo"),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(12, 1, 13), P(14, 1, 15))),
						ast.NewInvalidNode(
							L(S(P(12, 1, 13), P(14, 1, 15))),
							T(L(S(P(12, 1, 13), P(14, 1, 15))), token.END),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(6, 1, 7), P(9, 1, 10))), "unexpected PRIVATE_IDENTIFIER, expected a macro name (public identifier, keyword)"),
				diagnostic.NewFailure(L(S(P(12, 1, 13), P(14, 1, 15))), "unexpected end, expected an expression"),
			},
		},
		"can have a keyword as a name": {
			input: "macro class; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(10, 1, 11))), "class"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have an overridable operator as a name": {
			input: "macro +; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(6, 1, 7), P(7, 1, 8))),
						ast.NewInvalidNode(
							L(S(P(6, 1, 7), P(6, 1, 7))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.PLUS),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(9, 1, 10), P(11, 1, 12))),
						ast.NewInvalidNode(
							L(S(P(9, 1, 10), P(11, 1, 12))),
							T(L(S(P(9, 1, 10), P(11, 1, 12))), token.END),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(6, 1, 7), P(6, 1, 7))), "unexpected +, expected a macro name (public identifier, keyword)"),
				diagnostic.NewFailure(L(S(P(9, 1, 10), P(11, 1, 12))), "unexpected end, expected an expression"),
			},
		},
		"cannot have brackets as a name": {
			input: "macro []; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(12, 1, 13))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(6, 1, 7), P(8, 1, 9))),
						ast.NewInvalidNode(
							L(S(P(6, 1, 7), P(6, 1, 7))),
							T(L(S(P(6, 1, 7), P(6, 1, 7))), token.LBRACKET),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(10, 1, 11), P(12, 1, 13))),
						ast.NewInvalidNode(
							L(S(P(10, 1, 11), P(12, 1, 13))),
							T(L(S(P(10, 1, 11), P(12, 1, 13))), token.END),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(6, 1, 7), P(6, 1, 7))), "unexpected [, expected a macro name (public identifier, keyword)"),
				diagnostic.NewFailure(L(S(P(10, 1, 11), P(12, 1, 13))), "unexpected end, expected an expression"),
			},
		},
		"cannot have a public constant as a name": {
			input: "macro Foo; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(13, 1, 14))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(6, 1, 7), P(9, 1, 10))),
						ast.NewInvalidNode(
							L(S(P(6, 1, 7), P(8, 1, 9))),
							V(L(S(P(6, 1, 7), P(8, 1, 9))), token.PUBLIC_CONSTANT, "Foo"),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(11, 1, 12), P(13, 1, 14))),
						ast.NewInvalidNode(
							L(S(P(11, 1, 12), P(13, 1, 14))),
							T(L(S(P(11, 1, 12), P(13, 1, 14))), token.END),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(6, 1, 7), P(8, 1, 9))), "unexpected PUBLIC_CONSTANT, expected a macro name (public identifier, keyword)"),
				diagnostic.NewFailure(L(S(P(11, 1, 12), P(13, 1, 14))), "unexpected end, expected an expression"),
			},
		},
		"cannot have a non overridable operator as a name": {
			input: "macro &&; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(12, 1, 13))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(6, 1, 7), P(8, 1, 9))),
						ast.NewInvalidNode(
							L(S(P(6, 1, 7), P(7, 1, 8))),
							T(L(S(P(6, 1, 7), P(7, 1, 8))), token.AND_AND),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(10, 1, 11), P(12, 1, 13))),
						ast.NewInvalidNode(
							L(S(P(10, 1, 11), P(12, 1, 13))),
							T(L(S(P(10, 1, 11), P(12, 1, 13))), token.END),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(6, 1, 7), P(7, 1, 8))), "unexpected &&, expected a macro name (public identifier, keyword)"),
				diagnostic.NewFailure(L(S(P(10, 1, 11), P(12, 1, 13))), "unexpected end, expected an expression"),
			},
		},
		"cannot have a private constant as a name": {
			input: "macro _Foo; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(6, 1, 7), P(10, 1, 11))),
						ast.NewInvalidNode(
							L(S(P(6, 1, 7), P(9, 1, 10))),
							V(L(S(P(6, 1, 7), P(9, 1, 10))), token.PRIVATE_CONSTANT, "_Foo"),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(12, 1, 13), P(14, 1, 15))),
						ast.NewInvalidNode(
							L(S(P(12, 1, 13), P(14, 1, 15))),
							T(L(S(P(12, 1, 13), P(14, 1, 15))), token.END),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(6, 1, 7), P(9, 1, 10))), "unexpected PRIVATE_CONSTANT, expected a macro name (public identifier, keyword)"),
				diagnostic.NewFailure(L(S(P(12, 1, 13), P(14, 1, 15))), "unexpected end, expected an expression"),
			},
		},
		"can have an empty argument list": {
			input: "macro foo(); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(15, 1, 16))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(15, 1, 16))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have a return type": {
			input: "macro foo: String?; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(22, 1, 23))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(22, 1, 23))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(22, 1, 23))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(9, 1, 10), P(16, 1, 17))),
									ast.NewSimpleSymbolLiteralNode(
										L(S(P(9, 1, 10), P(16, 1, 17))),
										"String",
									),
								),
							},
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(9, 1, 10), P(9, 1, 10))), "unexpected :, expected a statement separator `\\n`, `;`"),
				diagnostic.NewFailure(L(S(P(17, 1, 18), P(17, 1, 18))), "unexpected ?, expected a statement separator `\\n`, `;`"),
			},
		},
		"cannot have a throw type": {
			input: "macro foo! NoMethodError | TypeError; end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(40, 1, 41))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(40, 1, 41))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(40, 1, 41))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(9, 1, 10), P(36, 1, 37))),
									ast.NewBinaryExpressionNode(
										L(S(P(9, 1, 10), P(35, 1, 36))),
										T(L(S(P(25, 1, 26), P(25, 1, 26))), token.OR),
										ast.NewUnaryExpressionNode(
											L(S(P(9, 1, 10), P(23, 1, 24))),
											T(L(S(P(9, 1, 10), P(9, 1, 10))), token.BANG),
											ast.NewPublicConstantNode(L(S(P(11, 1, 12), P(23, 1, 24))), "NoMethodError"),
										),
										ast.NewPublicConstantNode(L(S(P(27, 1, 28), P(35, 1, 36))), "TypeError"),
									),
								),
							},
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(9, 1, 10), P(9, 1, 10))), "unexpected !, expected a statement separator `\\n`, `;`"),
			},
		},
		"can have arguments": {
			input: "macro foo(a, b); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(19, 1, 20))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(19, 1, 20))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(19, 1, 20))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(10, 1, 11))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(13, 1, 14), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have a trailing comma in parameters": {
			input: "macro foo(a, b,); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(20, 1, 21))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(20, 1, 21))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(20, 1, 21))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(10, 1, 11))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(13, 1, 14), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have multiline parameters": {
			input: "macro foo(\na,\nb\n); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(21, 4, 6))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(21, 4, 6))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(21, 4, 6))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(11, 2, 1), P(11, 2, 1))),
									ast.NewPublicIdentifierNode(L(S(P(11, 2, 1), P(11, 2, 1))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(14, 3, 1), P(14, 3, 1))),
									ast.NewPublicIdentifierNode(L(S(P(14, 3, 1), P(14, 3, 1))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have a trailing comma in multiline parameters": {
			input: "macro foo(\na,\nb,\n); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(22, 4, 6))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(22, 4, 6))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(22, 4, 6))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(11, 2, 1), P(11, 2, 1))),
									ast.NewPublicIdentifierNode(L(S(P(11, 2, 1), P(11, 2, 1))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(14, 3, 1), P(14, 3, 1))),
									ast.NewPublicIdentifierNode(L(S(P(14, 3, 1), P(14, 3, 1))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have a positional rest parameter": {
			input: "macro foo(a, b, *c); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(23, 1, 24))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(23, 1, 24))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(23, 1, 24))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(10, 1, 11))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(13, 1, 14), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(16, 1, 17), P(17, 1, 18))),
									ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "c"),
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"cannot have a positional rest parameter with a default value": {
			input: "macro foo(a, b, *c = 3); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(27, 1, 28))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(27, 1, 28))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(27, 1, 28))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(10, 1, 11))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(13, 1, 14), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(16, 1, 17), P(21, 1, 22))),
									ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "c"),
									nil,
									ast.NewIntLiteralNode(L(S(P(21, 1, 22), P(21, 1, 22))), "3"),
									ast.PositionalRestParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(21, 1, 22), P(21, 1, 22))), "rest parameters cannot have default values"),
			},
		},
		"can have a positional rest parameter in the middle": {
			input: "macro foo(a, b, *c, d); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(26, 1, 27))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(26, 1, 27))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(26, 1, 27))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(10, 1, 11))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(13, 1, 14), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(16, 1, 17), P(17, 1, 18))),
									ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "c"),
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(20, 1, 21), P(20, 1, 21))),
									ast.NewPublicIdentifierNode(L(S(P(20, 1, 21), P(20, 1, 21))), "d"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"cannot have an optional parameter after a positional rest parameter": {
			input: "macro foo(a, b, *c, d = 3); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(30, 1, 31))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(30, 1, 31))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(30, 1, 31))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(10, 1, 11))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(13, 1, 14), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(16, 1, 17), P(17, 1, 18))),
									ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "c"),
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(20, 1, 21), P(24, 1, 25))),
									ast.NewPublicIdentifierNode(L(S(P(20, 1, 21), P(20, 1, 21))), "d"),
									nil,
									ast.NewIntLiteralNode(L(S(P(24, 1, 25), P(24, 1, 25))), "3"),
									ast.NormalParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(20, 1, 21), P(24, 1, 25))), "optional parameters cannot appear after rest parameters"),
			},
		},
		"cannot have multiple positional rest parameters": {
			input: "macro foo(a, b, *c, *d); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(27, 1, 28))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(27, 1, 28))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(27, 1, 28))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(10, 1, 11))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(13, 1, 14), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(16, 1, 17), P(17, 1, 18))),
									ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "c"),
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(20, 1, 21), P(21, 1, 22))),
									ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(21, 1, 22))), "d"),
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(20, 1, 21), P(21, 1, 22))), "there should be only a single positional rest parameter"),
			},
		},
		"can have a positional rest parameter with a type": {
			input: "macro foo(a, b, *c: String); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(31, 1, 32))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(31, 1, 32))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(31, 1, 32))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(10, 1, 11))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(13, 1, 14), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(16, 1, 17), P(25, 1, 26))),
									ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "c"),
									ast.NewPublicConstantNode(L(S(P(20, 1, 21), P(25, 1, 26))), "String"),
									nil,
									ast.PositionalRestParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have a named rest parameter": {
			input: "macro foo(a, b, **c); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(24, 1, 25))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(24, 1, 25))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(24, 1, 25))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(10, 1, 11))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(13, 1, 14), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(16, 1, 17), P(18, 1, 19))),
									ast.NewPublicIdentifierNode(L(S(P(18, 1, 19), P(18, 1, 19))), "c"),
									nil,
									nil,
									ast.NamedRestParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"cannot have a named rest parameter with a default value": {
			input: "macro foo(a, b, **c = 3); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(28, 1, 29))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(28, 1, 29))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(28, 1, 29))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(10, 1, 11))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(13, 1, 14), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(16, 1, 17), P(22, 1, 23))),
									ast.NewPublicIdentifierNode(L(S(P(18, 1, 19), P(18, 1, 19))), "c"),
									nil,
									ast.NewIntLiteralNode(L(S(P(22, 1, 23), P(22, 1, 23))), "3"),
									ast.NamedRestParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(22, 1, 23), P(22, 1, 23))), "rest parameters cannot have default values"),
			},
		},
		"can have a named rest parameter with a type": {
			input: "macro foo(a, b, **c: String); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(32, 1, 33))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(32, 1, 33))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(32, 1, 33))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(10, 1, 11))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(13, 1, 14), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(16, 1, 17), P(26, 1, 27))),
									ast.NewPublicIdentifierNode(L(S(P(18, 1, 19), P(18, 1, 19))), "c"),
									ast.NewPublicConstantNode(L(S(P(21, 1, 22), P(26, 1, 27))), "String"),
									nil,
									ast.NamedRestParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"cannot have parameters after a named rest parameter": {
			input: "macro foo(a, b, **c, d); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(27, 1, 28))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(27, 1, 28))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(27, 1, 28))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(10, 1, 11))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(13, 1, 14), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(16, 1, 17), P(18, 1, 19))),
									ast.NewPublicIdentifierNode(L(S(P(18, 1, 19), P(18, 1, 19))), "c"),
									nil,
									nil,
									ast.NamedRestParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(21, 1, 22), P(21, 1, 22))),
									ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(21, 1, 22))), "d"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(21, 1, 22), P(21, 1, 22))), "named rest parameters should appear last"),
			},
		},
		"can have a positional and named rest parameter": {
			input: "macro foo(a, b, *c, **d); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(28, 1, 29))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(28, 1, 29))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(28, 1, 29))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(10, 1, 11))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(13, 1, 14), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(16, 1, 17), P(17, 1, 18))),
									ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "c"),
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(20, 1, 21), P(22, 1, 23))),
									ast.NewPublicIdentifierNode(L(S(P(22, 1, 23), P(22, 1, 23))), "d"),
									nil,
									nil,
									ast.NamedRestParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"cannot have a post parameter and a named rest parameter": {
			input: "macro foo(a, b, *c, d, **e); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(31, 1, 32))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(31, 1, 32))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(31, 1, 32))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(10, 1, 11))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(13, 1, 14), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "b"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(16, 1, 17), P(17, 1, 18))),
									ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "c"),
									nil,
									nil,
									ast.PositionalRestParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(20, 1, 21), P(20, 1, 21))),
									ast.NewPublicIdentifierNode(L(S(P(20, 1, 21), P(20, 1, 21))), "d"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(23, 1, 24), P(25, 1, 26))),
									ast.NewPublicIdentifierNode(L(S(P(25, 1, 26), P(25, 1, 26))), "e"),
									nil,
									nil,
									ast.NamedRestParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(23, 1, 24), P(25, 1, 26))), "named rest parameters cannot appear after a post parameter"),
			},
		},
		"can have arguments with types": {
			input: "macro foo(a: Int, b: String?); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(33, 1, 34))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(33, 1, 34))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(33, 1, 34))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(15, 1, 16))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									ast.NewPublicConstantNode(L(S(P(13, 1, 14), P(15, 1, 16))), "Int"),
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(18, 1, 19), P(27, 1, 28))),
									ast.NewPublicIdentifierNode(L(S(P(18, 1, 19), P(18, 1, 19))), "b"),
									ast.NewNilableTypeNode(
										L(S(P(21, 1, 22), P(27, 1, 28))),
										ast.NewPublicConstantNode(L(S(P(21, 1, 22), P(26, 1, 27))), "String"),
									),
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have arguments with initialisers": {
			input: "macro foo(a = 32, b: String = 'foo'); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(40, 1, 41))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(40, 1, 41))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(40, 1, 41))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(15, 1, 16))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									ast.NewIntLiteralNode(L(S(P(14, 1, 15), P(15, 1, 16))), "32"),
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(18, 1, 19), P(34, 1, 35))),
									ast.NewPublicIdentifierNode(L(S(P(18, 1, 19), P(18, 1, 19))), "b"),
									ast.NewPublicConstantNode(L(S(P(21, 1, 22), P(26, 1, 27))), "String"),
									ast.NewRawStringLiteralNode(L(S(P(30, 1, 31), P(34, 1, 35))), "foo"),
									ast.NormalParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
		},
		"cannot have arguments that set instance variables": {
			input: "macro foo(@a = 32, @b: String = 'foo'); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(42, 1, 43))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(13, 1, 14), P(38, 1, 39))),
						ast.NewInvalidNode(
							L(S(P(13, 1, 14), P(13, 1, 14))),
							T(L(S(P(13, 1, 14), P(13, 1, 14))), token.EQUAL_OP),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(40, 1, 41), P(42, 1, 43))),
						ast.NewInvalidNode(
							L(S(P(40, 1, 41), P(42, 1, 43))),
							T(L(S(P(40, 1, 41), P(42, 1, 43))), token.END),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(10, 1, 11), P(11, 1, 12))), "unexpected INSTANCE_VARIABLE, expected an identifier as the name of the declared signature parameter"),
				diagnostic.NewFailure(L(S(P(13, 1, 14), P(13, 1, 14))), "unexpected =, expected )"),
				diagnostic.NewFailure(L(S(P(40, 1, 41), P(42, 1, 43))), "unexpected end, expected an expression"),
			},
		},
		"cannot have required arguments after optional ones": {
			input: "macro foo(a = 32, b: String, c = true, d); end",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(45, 1, 46))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(45, 1, 46))),
						ast.NewMacroDefinitionNode(
							L(S(P(0, 1, 1), P(45, 1, 46))),
							"",
							false,
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "foo"),
							[]ast.ParameterNode{
								ast.NewFormalParameterNode(
									L(S(P(10, 1, 11), P(15, 1, 16))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "a"),
									nil,
									ast.NewIntLiteralNode(L(S(P(14, 1, 15), P(15, 1, 16))), "32"),
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(18, 1, 19), P(26, 1, 27))),
									ast.NewPublicIdentifierNode(L(S(P(18, 1, 19), P(18, 1, 19))), "b"),
									ast.NewPublicConstantNode(L(S(P(21, 1, 22), P(26, 1, 27))), "String"),
									nil,
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(29, 1, 30), P(36, 1, 37))),
									ast.NewPublicIdentifierNode(L(S(P(29, 1, 30), P(29, 1, 30))), "c"),
									nil,
									ast.NewTrueLiteralNode(L(S(P(33, 1, 34), P(36, 1, 37)))),
									ast.NormalParameterKind,
								),
								ast.NewFormalParameterNode(
									L(S(P(39, 1, 40), P(39, 1, 40))),
									ast.NewPublicIdentifierNode(L(S(P(39, 1, 40), P(39, 1, 40))), "d"),
									nil,
									nil,
									ast.NormalParameterKind,
								),
							},
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(18, 1, 19), P(26, 1, 27))), "required parameters cannot appear after optional parameters"),
				diagnostic.NewFailure(L(S(P(39, 1, 40), P(39, 1, 40))), "required parameters cannot appear after optional parameters"),
			},
		},
		"can have a multiline body": {
			input: `def foo
  a := .5
  a += .7
end`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(30, 4, 3))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(30, 4, 3))),
						ast.NewMethodDefinitionNode(
							L(S(P(0, 1, 1), P(30, 4, 3))),
							"",
							0,
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "foo"),
							nil,
							nil,
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(10, 2, 3), P(17, 2, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(10, 2, 3), P(16, 2, 9))),
										T(L(S(P(12, 2, 5), P(13, 2, 6))), token.COLON_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(10, 2, 3), P(10, 2, 3))), "a"),
										ast.NewFloatLiteralNode(L(S(P(15, 2, 8), P(16, 2, 9))), "0.5"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(20, 3, 3), P(27, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(20, 3, 3), P(26, 3, 9))),
										T(L(S(P(22, 3, 5), P(23, 3, 6))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(20, 3, 3), P(20, 3, 3))), "a"),
										ast.NewFloatLiteralNode(L(S(P(25, 3, 8), P(26, 3, 9))), "0.7"),
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
				L(S(P(0, 1, 1), P(19, 1, 20))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(19, 1, 20))),
						ast.NewMethodDefinitionNode(
							L(S(P(0, 1, 1), P(19, 1, 20))),
							"",
							0,
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "foo"),
							nil,
							nil,
							nil,
							nil,
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(13, 1, 14), P(19, 1, 20))),
									ast.NewBinaryExpressionNode(
										L(S(P(13, 1, 14), P(19, 1, 20))),
										T(L(S(P(16, 1, 17), P(16, 1, 17))), token.PLUS),
										ast.NewFloatLiteralNode(L(S(P(13, 1, 14), P(14, 1, 15))), "0.3"),
										ast.NewFloatLiteralNode(L(S(P(18, 1, 19), P(19, 1, 20))), "0.4"),
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

func TestMacroCall(t *testing.T) {
	tests := testTable{
		"can be a part of an expression": {
			input: "foo = bar!()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(11, 1, 12))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(11, 1, 12))),
						ast.NewAssignmentExpressionNode(
							L(S(P(0, 1, 1), P(11, 1, 12))),
							T(L(S(P(4, 1, 5), P(4, 1, 5))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewReceiverlessMacroCallNode(
								L(S(P(6, 1, 7), P(11, 1, 12))),
								ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "bar"),
								nil,
								nil,
							),
						),
					),
				},
			),
		},
		"can omit the receiver and have an empty argument list": {
			input: "foo!()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(5, 1, 6))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(5, 1, 6))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(5, 1, 6))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can omit the receiver, parens and arguments": {
			input: "foo!",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(3, 1, 4))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(3, 1, 4))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(2, 1, 3))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot omit the receiver and have type arguments": {
			input: "foo!::[String]()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(2, 1, 3))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(2, 1, 3))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(2, 1, 3))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							nil,
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(4, 1, 5), P(6, 1, 7))), "unexpected ::[, expected a statement separator `\\n`, `;`"),
			},
		},
		"can omit the receiver, arguments and have a trailing closure": {
			input: "foo!() |i| -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(18, 1, 19))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(18, 1, 19))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(18, 1, 19))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewClosureLiteralNode(
									L(S(P(7, 1, 8), P(18, 1, 19))),
									[]ast.ParameterNode{
										ast.NewFormalParameterNode(
											L(S(P(8, 1, 9), P(8, 1, 9))),
											ast.NewPublicIdentifierNode(L(S(P(8, 1, 9), P(8, 1, 9))), "i"),
											nil,
											nil,
											ast.NormalParameterKind,
										),
									},
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(14, 1, 15), P(18, 1, 19))),
											ast.NewBinaryExpressionNode(
												L(S(P(14, 1, 15), P(18, 1, 19))),
												T(L(S(P(16, 1, 17), P(16, 1, 17))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(14, 1, 15))), "i"),
												ast.NewIntLiteralNode(L(S(P(18, 1, 19), P(18, 1, 19))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can omit the receiver, arguments and have a trailing closure without pipes": {
			input: "foo!() -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(14, 1, 15))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(14, 1, 15))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(14, 1, 15))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewClosureLiteralNode(
									L(S(P(7, 1, 8), P(14, 1, 15))),
									nil,
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(10, 1, 11), P(14, 1, 15))),
											ast.NewBinaryExpressionNode(
												L(S(P(10, 1, 11), P(14, 1, 15))),
												T(L(S(P(12, 1, 13), P(12, 1, 13))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(10, 1, 11))), "i"),
												ast.NewIntLiteralNode(L(S(P(14, 1, 15), P(14, 1, 15))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can omit the receiver, arguments and have a trailing closure without arguments": {
			input: "foo!() || -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(17, 1, 18))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(17, 1, 18))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(17, 1, 18))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewClosureLiteralNode(
									L(S(P(7, 1, 8), P(17, 1, 18))),
									nil,
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(13, 1, 14), P(17, 1, 18))),
											ast.NewBinaryExpressionNode(
												L(S(P(13, 1, 14), P(17, 1, 18))),
												T(L(S(P(15, 1, 16), P(15, 1, 16))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "i"),
												ast.NewIntLiteralNode(L(S(P(17, 1, 18), P(17, 1, 18))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can omit the receiver and have a trailing closure": {
			input: "foo!(1, 5) |i| -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(22, 1, 23))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(22, 1, 23))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(22, 1, 23))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(L(S(P(5, 1, 6), P(5, 1, 6))), "1"),
								ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "5"),
								ast.NewClosureLiteralNode(
									L(S(P(11, 1, 12), P(22, 1, 23))),
									[]ast.ParameterNode{
										ast.NewFormalParameterNode(
											L(S(P(12, 1, 13), P(12, 1, 13))),
											ast.NewPublicIdentifierNode(L(S(P(12, 1, 13), P(12, 1, 13))), "i"),
											nil,
											nil,
											ast.NormalParameterKind,
										),
									},
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(18, 1, 19), P(22, 1, 23))),
											ast.NewBinaryExpressionNode(
												L(S(P(18, 1, 19), P(22, 1, 23))),
												T(L(S(P(20, 1, 21), P(20, 1, 21))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(18, 1, 19), P(18, 1, 19))), "i"),
												ast.NewIntLiteralNode(L(S(P(22, 1, 23), P(22, 1, 23))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can omit the receiver have named arguments and a trailing closure": {
			input: "foo!(f: 5) |i| -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(22, 1, 23))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(22, 1, 23))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(22, 1, 23))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(5, 1, 6), P(8, 1, 9))),
									ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(5, 1, 6))), "f"),
									ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "5"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(11, 1, 12), P(22, 1, 23))),
									ast.NewPublicIdentifierNode(L(S(P(11, 1, 12), P(22, 1, 23))), "func"),
									ast.NewClosureLiteralNode(
										L(S(P(11, 1, 12), P(22, 1, 23))),
										[]ast.ParameterNode{
											ast.NewFormalParameterNode(
												L(S(P(12, 1, 13), P(12, 1, 13))),
												ast.NewPublicIdentifierNode(L(S(P(12, 1, 13), P(12, 1, 13))), "i"),
												nil,
												nil,
												ast.NormalParameterKind,
											),
										},
										nil,
										nil,
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												L(S(P(18, 1, 19), P(22, 1, 23))),
												ast.NewBinaryExpressionNode(
													L(S(P(18, 1, 19), P(22, 1, 23))),
													T(L(S(P(20, 1, 21), P(20, 1, 21))), token.STAR),
													ast.NewPublicIdentifierNode(L(S(P(18, 1, 19), P(18, 1, 19))), "i"),
													ast.NewIntLiteralNode(L(S(P(22, 1, 23), P(22, 1, 23))), "2"),
												),
											),
										},
									),
								),
							},
						),
					),
				},
			),
		},
		"can have a private name with an implicit receiver": {
			input: "_foo!()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							ast.NewPrivateIdentifierNode(L(S(P(0, 1, 1), P(3, 1, 4))), "_foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have an explicit receiver": {
			input: "foo.bar!()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewMacroCallNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot have an explicit receiver with type arguments": {
			input: "foo.bar!::[String]()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(19, 1, 20))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(19, 1, 20))),
						ast.NewMacroCallNode(
							L(S(P(0, 1, 1), P(19, 1, 20))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							nil,
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				&diagnostic.Diagnostic{
					Location: L(S(P(8, 1, 9), P(17, 1, 18))),
					Message:  "macro calls cannot be generic",
				},
			},
		},
		"can have an explicit receiver and a trailing closure without pipes": {
			input: "foo.bar!() -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(18, 1, 19))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(18, 1, 19))),
						ast.NewMacroCallNode(
							L(S(P(0, 1, 1), P(18, 1, 19))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							[]ast.ExpressionNode{
								ast.NewClosureLiteralNode(
									L(S(P(11, 1, 12), P(18, 1, 19))),
									nil,
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(14, 1, 15), P(18, 1, 19))),
											ast.NewBinaryExpressionNode(
												L(S(P(14, 1, 15), P(18, 1, 19))),
												T(L(S(P(16, 1, 17), P(16, 1, 17))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(14, 1, 15), P(14, 1, 15))), "i"),
												ast.NewIntLiteralNode(L(S(P(18, 1, 19), P(18, 1, 19))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have an explicit receiver and a trailing closure without arguments": {
			input: "foo.bar!() || -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(21, 1, 22))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(21, 1, 22))),
						ast.NewMacroCallNode(
							L(S(P(0, 1, 1), P(21, 1, 22))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							[]ast.ExpressionNode{
								ast.NewClosureLiteralNode(
									L(S(P(11, 1, 12), P(21, 1, 22))),
									nil,
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(17, 1, 18), P(21, 1, 22))),
											ast.NewBinaryExpressionNode(
												L(S(P(17, 1, 18), P(21, 1, 22))),
												T(L(S(P(19, 1, 20), P(19, 1, 20))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "i"),
												ast.NewIntLiteralNode(L(S(P(21, 1, 22), P(21, 1, 22))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have an explicit receiver and a trailing closure": {
			input: "foo.bar!() |i| -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(22, 1, 23))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(22, 1, 23))),
						ast.NewMacroCallNode(
							L(S(P(0, 1, 1), P(22, 1, 23))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							[]ast.ExpressionNode{
								ast.NewClosureLiteralNode(
									L(S(P(11, 1, 12), P(22, 1, 23))),
									[]ast.ParameterNode{
										ast.NewFormalParameterNode(
											L(S(P(12, 1, 13), P(12, 1, 13))),
											ast.NewPublicIdentifierNode(L(S(P(12, 1, 13), P(12, 1, 13))), "i"),
											nil,
											nil,
											ast.NormalParameterKind,
										),
									},
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(18, 1, 19), P(22, 1, 23))),
											ast.NewBinaryExpressionNode(
												L(S(P(18, 1, 19), P(22, 1, 23))),
												T(L(S(P(20, 1, 21), P(20, 1, 21))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(18, 1, 19), P(18, 1, 19))), "i"),
												ast.NewIntLiteralNode(L(S(P(22, 1, 23), P(22, 1, 23))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have an explicit receiver, arguments and a trailing closure": {
			input: "foo.bar!(1, 5) |i| -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(26, 1, 27))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(26, 1, 27))),
						ast.NewMacroCallNode(
							L(S(P(0, 1, 1), P(26, 1, 27))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(L(S(P(9, 1, 10), P(9, 1, 10))), "1"),
								ast.NewIntLiteralNode(L(S(P(12, 1, 13), P(12, 1, 13))), "5"),
								ast.NewClosureLiteralNode(
									L(S(P(15, 1, 16), P(26, 1, 27))),
									[]ast.ParameterNode{
										ast.NewFormalParameterNode(
											L(S(P(16, 1, 17), P(16, 1, 17))),
											ast.NewPublicIdentifierNode(L(S(P(16, 1, 17), P(16, 1, 17))), "i"),
											nil,
											nil,
											ast.NormalParameterKind,
										),
									},
									nil,
									nil,
									[]ast.StatementNode{
										ast.NewExpressionStatementNode(
											L(S(P(22, 1, 23), P(26, 1, 27))),
											ast.NewBinaryExpressionNode(
												L(S(P(22, 1, 23), P(26, 1, 27))),
												T(L(S(P(24, 1, 25), P(24, 1, 25))), token.STAR),
												ast.NewPublicIdentifierNode(L(S(P(22, 1, 23), P(22, 1, 23))), "i"),
												ast.NewIntLiteralNode(L(S(P(26, 1, 27), P(26, 1, 27))), "2"),
											),
										),
									},
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have an explicit receiver, named arguments and a trailing closure": {
			input: "foo.bar!(f: 5) |i| -> i * 2",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(26, 1, 27))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(26, 1, 27))),
						ast.NewMacroCallNode(
							L(S(P(0, 1, 1), P(26, 1, 27))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(9, 1, 10), P(12, 1, 13))),
									ast.NewPublicIdentifierNode(L(S(P(9, 1, 10), P(9, 1, 10))), "f"),
									ast.NewIntLiteralNode(L(S(P(12, 1, 13), P(12, 1, 13))), "5"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(15, 1, 16), P(26, 1, 27))),
									ast.NewPublicIdentifierNode(L(S(P(15, 1, 16), P(26, 1, 27))), "func"),
									ast.NewClosureLiteralNode(
										L(S(P(15, 1, 16), P(26, 1, 27))),
										[]ast.ParameterNode{
											ast.NewFormalParameterNode(
												L(S(P(16, 1, 17), P(16, 1, 17))),
												ast.NewPublicIdentifierNode(L(S(P(16, 1, 17), P(16, 1, 17))), "i"),
												nil,
												nil,
												ast.NormalParameterKind,
											),
										},
										nil,
										nil,
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												L(S(P(22, 1, 23), P(26, 1, 27))),
												ast.NewBinaryExpressionNode(
													L(S(P(22, 1, 23), P(26, 1, 27))),
													T(L(S(P(24, 1, 25), P(24, 1, 25))), token.STAR),
													ast.NewPublicIdentifierNode(L(S(P(22, 1, 23), P(22, 1, 23))), "i"),
													ast.NewIntLiteralNode(L(S(P(26, 1, 27), P(26, 1, 27))), "2"),
												),
											),
										},
									),
								),
							},
						),
					),
				},
			),
		},
		"can omit parentheses": {
			input: "foo.bar! 1",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewMacroCallNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
							[]ast.ExpressionNode{
								ast.NewIntLiteralNode(L(S(P(9, 1, 10), P(9, 1, 10))), "1"),
							},
							nil,
						),
					),
				},
			),
		},
		"cannot use the safe navigation operator": {
			input: "foo?.bar!",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewMacroCallNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "bar"),
							nil,
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(S(P(3, 1, 4), P(4, 1, 5))),
					"invalid macro call operator",
				),
			},
		},
		"cannot use the cascade call operator": {
			input: "foo..bar!",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewMacroCallNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "bar"),
							nil,
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(S(P(3, 1, 4), P(4, 1, 5))),
					"invalid macro call operator",
				),
			},
		},
		"cannot use the safe cascade call operator": {
			input: "foo?..bar!",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewMacroCallNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(8, 1, 9))), "bar"),
							nil,
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(
					L(S(P(3, 1, 4), P(5, 1, 6))),
					"invalid macro call operator",
				),
			},
		},
		"can be nested with parentheses": {
			input: "foo.bar!().baz!()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(16, 1, 17))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(16, 1, 17))),
						ast.NewMacroCallNode(
							L(S(P(0, 1, 1), P(16, 1, 17))),
							ast.NewMacroCallNode(
								L(S(P(0, 1, 1), P(9, 1, 10))),
								ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
								ast.NewPublicIdentifierNode(L(S(P(4, 1, 5), P(6, 1, 7))), "bar"),
								nil,
								nil,
							),
							ast.NewPublicIdentifierNode(L(S(P(11, 1, 12), P(13, 1, 14))), "baz"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have any expression as the receiver": {
			input: "(foo + 2).bar!()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(15, 1, 16))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(1, 1, 2), P(15, 1, 16))),
						ast.NewMacroCallNode(
							L(S(P(1, 1, 2), P(15, 1, 16))),
							ast.NewBinaryExpressionNode(
								L(S(P(1, 1, 2), P(7, 1, 8))),
								T(L(S(P(5, 1, 6), P(5, 1, 6))), token.PLUS),
								ast.NewPublicIdentifierNode(L(S(P(1, 1, 2), P(3, 1, 4))), "foo"),
								ast.NewIntLiteralNode(L(S(P(7, 1, 8), P(7, 1, 8))), "2"),
							),
							ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(12, 1, 13))), "bar"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"cannot call a private method on an explicit receiver": {
			input: "foo._bar!()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(10, 1, 11))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(10, 1, 11))),
						ast.NewMacroCallNode(
							L(S(P(0, 1, 1), P(10, 1, 11))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewPrivateIdentifierNode(L(S(P(4, 1, 5), P(7, 1, 8))), "_bar"),
							nil,
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(4, 1, 5), P(7, 1, 8))), "unexpected PRIVATE_IDENTIFIER, expected a public method name (public identifier, keyword or overridable operator)"),
			},
		},
		"cannot have a non overridable operator as the macro name": {
			input: "foo.&&!()",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(8, 1, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(8, 1, 9))),
						ast.NewMacroCallNode(
							L(S(P(0, 1, 1), P(8, 1, 9))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							ast.NewInvalidNode(
								L(S(P(4, 1, 5), P(5, 1, 6))),
								T(L(S(P(4, 1, 5), P(5, 1, 6))), token.AND_AND),
							),
							nil,
							nil,
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(4, 1, 5), P(5, 1, 6))), "unexpected &&, expected a public method name (public identifier, keyword or overridable operator)"),
				diagnostic.NewFailure(L(S(P(4, 1, 5), P(5, 1, 6))), "invalid macro name"),
			},
		},
		"can have positional arguments": {
			input: "foo!(.1, 'foo', :bar)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(20, 1, 21))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(20, 1, 21))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(20, 1, 21))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(5, 1, 6), P(6, 1, 7))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(9, 1, 10), P(13, 1, 14))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have splat arguments": {
			input: "foo!(*baz, 'foo', *bar)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(22, 1, 23))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(22, 1, 23))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(22, 1, 23))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewSplatExpressionNode(
									L(S(P(5, 1, 6), P(8, 1, 9))),
									ast.NewPublicIdentifierNode(
										L(S(P(6, 1, 7), P(8, 1, 9))),
										"baz",
									),
								),
								ast.NewRawStringLiteralNode(L(S(P(11, 1, 12), P(15, 1, 16))), "foo"),
								ast.NewSplatExpressionNode(
									L(S(P(18, 1, 19), P(21, 1, 22))),
									ast.NewPublicIdentifierNode(
										L(S(P(19, 1, 20), P(21, 1, 22))),
										"bar",
									),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have a trailing comma": {
			input: "foo!(.1, 'foo', :bar,)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(21, 1, 22))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(21, 1, 22))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(21, 1, 22))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(5, 1, 6), P(6, 1, 7))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(9, 1, 10), P(13, 1, 14))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have named arguments": {
			input: "foo!(bar: :baz, elk: true)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(25, 1, 26))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(25, 1, 26))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(25, 1, 26))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(5, 1, 6), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(10, 1, 11), P(13, 1, 14))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(16, 1, 17), P(24, 1, 25))),
									ast.NewPublicIdentifierNode(L(S(P(16, 1, 17), P(18, 1, 19))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(21, 1, 22), P(24, 1, 25)))),
								),
							},
						),
					),
				},
			),
		},
		"can have double splat arguments": {
			input: "foo!(**f, bar: :baz, **dupa(), elk: true)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(40, 1, 41))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(40, 1, 41))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(40, 1, 41))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewDoubleSplatExpressionNode(
									L(S(P(5, 1, 6), P(7, 1, 8))),
									ast.NewPublicIdentifierNode(L(S(P(7, 1, 8), P(7, 1, 8))), "f"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(10, 1, 11), P(18, 1, 19))),
									ast.NewPublicIdentifierNode(L(S(P(10, 1, 11), P(12, 1, 13))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "baz"),
								),
								ast.NewDoubleSplatExpressionNode(
									L(S(P(21, 1, 22), P(28, 1, 29))),
									ast.NewReceiverlessMethodCallNode(
										L(S(P(23, 1, 24), P(28, 1, 29))),
										ast.NewPublicIdentifierNode(L(S(P(23, 1, 24), P(26, 1, 27))), "dupa"),
										nil,
										nil,
									),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(31, 1, 32), P(39, 1, 40))),
									ast.NewPublicIdentifierNode(L(S(P(31, 1, 32), P(33, 1, 34))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(36, 1, 37), P(39, 1, 40)))),
								),
							},
						),
					),
				},
			),
		},
		"can have positional and named arguments": {
			input: "foo!(.1, 'foo', :bar, bar: :baz, elk: true)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(42, 1, 43))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(42, 1, 43))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(42, 1, 43))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(5, 1, 6), P(6, 1, 7))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(9, 1, 10), P(13, 1, 14))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(22, 1, 23), P(30, 1, 31))),
									ast.NewPublicIdentifierNode(L(S(P(22, 1, 23), P(24, 1, 25))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(27, 1, 28), P(30, 1, 31))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(33, 1, 34), P(41, 1, 42))),
									ast.NewPublicIdentifierNode(L(S(P(33, 1, 34), P(35, 1, 36))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(38, 1, 39), P(41, 1, 42)))),
								),
							},
						),
					),
				},
			),
		},
		"can have newlines after commas": {
			input: "foo!(.1,\n'foo',\n:bar, bar: :baz,\nelk: true)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(42, 4, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(42, 4, 10))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(42, 4, 10))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(5, 1, 6), P(6, 1, 7))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(9, 2, 1), P(13, 2, 5))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(16, 3, 1), P(19, 3, 4))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(22, 3, 7), P(30, 3, 15))),
									ast.NewPublicIdentifierNode(L(S(P(22, 3, 7), P(24, 3, 9))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(27, 3, 12), P(30, 3, 15))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(33, 4, 1), P(41, 4, 9))),
									ast.NewPublicIdentifierNode(L(S(P(33, 4, 1), P(35, 4, 3))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(38, 4, 6), P(41, 4, 9)))),
								),
							},
						),
					),
				},
			),
		},
		"can have newlines around parentheses": {
			input: "foo!(\n.1, 'foo', :bar, bar: :baz, elk: true\n)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(44, 3, 1))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(44, 3, 1))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(44, 3, 1))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(6, 2, 1), P(7, 2, 2))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(10, 2, 5), P(14, 2, 9))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(17, 2, 12), P(20, 2, 15))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(23, 2, 18), P(31, 2, 26))),
									ast.NewPublicIdentifierNode(L(S(P(23, 2, 18), P(25, 2, 20))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(28, 2, 23), P(31, 2, 26))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(34, 2, 29), P(42, 2, 37))),
									ast.NewPublicIdentifierNode(L(S(P(34, 2, 29), P(36, 2, 31))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(39, 2, 34), P(42, 2, 37)))),
								),
							},
						),
					),
				},
			),
		},
		"cannot have newlines before the opening parenthesis": {
			input: "foo!\n(.1, 'foo', :bar, bar: :baz, elk: true)",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(43, 2, 39))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(4, 1, 5))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(2, 1, 3))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							nil,
							nil,
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(6, 2, 2), P(43, 2, 39))),
						ast.NewFloatLiteralNode(L(S(P(6, 2, 2), P(7, 2, 3))), "0.1"),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(8, 2, 4), P(8, 2, 4))), "unexpected ,, expected )"),
			},
		},
		"can have positional arguments without parentheses": {
			input: "foo! .1, 'foo', :bar",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(19, 1, 20))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(19, 1, 20))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(19, 1, 20))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(5, 1, 6), P(6, 1, 7))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(9, 1, 10), P(13, 1, 14))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "bar"),
							},
							nil,
						),
					),
				},
			),
		},
		"can have named arguments without parentheses": {
			input: "foo! bar: :baz, elk: true",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(24, 1, 25))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(24, 1, 25))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(24, 1, 25))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							nil,
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(5, 1, 6), P(13, 1, 14))),
									ast.NewPublicIdentifierNode(L(S(P(5, 1, 6), P(7, 1, 8))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(10, 1, 11), P(13, 1, 14))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(16, 1, 17), P(24, 1, 25))),
									ast.NewPublicIdentifierNode(L(S(P(16, 1, 17), P(18, 1, 19))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(21, 1, 22), P(24, 1, 25)))),
								),
							},
						),
					),
				},
			),
		},
		"can have positional and named arguments without parentheses": {
			input: "foo! .1, 'foo', :bar, bar: :baz, elk: true",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(41, 1, 42))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(41, 1, 42))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(41, 1, 42))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(5, 1, 6), P(6, 1, 7))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(9, 1, 10), P(13, 1, 14))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(16, 1, 17), P(19, 1, 20))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(22, 1, 23), P(30, 1, 31))),
									ast.NewPublicIdentifierNode(L(S(P(22, 1, 23), P(24, 1, 25))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(27, 1, 28), P(30, 1, 31))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(33, 1, 34), P(41, 1, 42))),
									ast.NewPublicIdentifierNode(L(S(P(33, 1, 34), P(35, 1, 36))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(38, 1, 39), P(41, 1, 42)))),
								),
							},
						),
					),
				},
			),
		},
		"can have newlines after commas without parentheses": {
			input: "foo! .1,\n'foo',\n:bar, bar: :baz,\nelk: true",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(41, 4, 9))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(41, 4, 9))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(41, 4, 9))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							[]ast.ExpressionNode{
								ast.NewFloatLiteralNode(L(S(P(5, 1, 6), P(6, 1, 7))), "0.1"),
								ast.NewRawStringLiteralNode(L(S(P(9, 2, 1), P(13, 2, 5))), "foo"),
								ast.NewSimpleSymbolLiteralNode(L(S(P(16, 3, 1), P(19, 3, 4))), "bar"),
							},
							[]ast.NamedArgumentNode{
								ast.NewNamedCallArgumentNode(
									L(S(P(22, 3, 7), P(30, 3, 15))),
									ast.NewPublicIdentifierNode(L(S(P(22, 3, 7), P(24, 3, 9))), "bar"),
									ast.NewSimpleSymbolLiteralNode(L(S(P(27, 3, 12), P(30, 3, 15))), "baz"),
								),
								ast.NewNamedCallArgumentNode(
									L(S(P(33, 4, 1), P(41, 4, 9))),
									ast.NewPublicIdentifierNode(L(S(P(33, 4, 1), P(35, 4, 3))), "elk"),
									ast.NewTrueLiteralNode(L(S(P(38, 4, 6), P(41, 4, 9)))),
								),
							},
						),
					),
				},
			),
		},
		"cannot have newlines before the arguments without parentheses": {
			input: "foo!\n.1, 'foo', :bar, bar: :baz, elk: true",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 2, 2))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(4, 1, 5))),
						ast.NewReceiverlessMacroCallNode(
							L(S(P(0, 1, 1), P(2, 1, 3))),
							ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
							nil,
							nil,
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(5, 2, 1), P(6, 2, 2))),
						ast.NewFloatLiteralNode(L(S(P(5, 2, 1), P(6, 2, 2))), "0.1"),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(7, 2, 3), P(7, 2, 3))), "unexpected ,, expected a statement separator `\\n`, `;`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestMacroBoundary(t *testing.T) {
	tests := testTable{
		"can have a multiline body": {
			input: `
do macro
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(28, 5, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(28, 5, 4))),
						ast.NewMacroBoundaryNode(
							L(S(P(1, 2, 1), P(27, 5, 3))),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(11, 3, 2), P(19, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(11, 3, 2), P(18, 3, 9))),
										T(L(S(P(15, 3, 6), P(16, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(11, 3, 2), P(13, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(18, 3, 9), P(18, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(21, 4, 2), P(24, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(21, 4, 2), P(23, 4, 4)))),
								),
							},
							"",
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
				do macro
				end
			`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(21, 3, 8))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(5, 2, 5), P(21, 3, 8))),
						ast.NewMacroBoundaryNode(
							L(S(P(5, 2, 5), P(20, 3, 7))),
							nil,
							"",
						),
					),
				},
			),
		},
		"can have a name": {
			input: `
				do macro 'foo'
				end
			`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(27, 3, 8))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(5, 2, 5), P(27, 3, 8))),
						ast.NewMacroBoundaryNode(
							L(S(P(5, 2, 5), P(26, 3, 7))),
							nil,
							"foo",
						),
					),
				},
			),
		},
		"can be a one-liner": {
			input: "do macro 5",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(9, 1, 10))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(9, 1, 10))),
						ast.NewMacroBoundaryNode(
							L(S(P(0, 1, 1), P(9, 1, 10))),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(9, 1, 10), P(9, 1, 10))),
									ast.NewIntLiteralNode(
										L(S(P(9, 1, 10), P(9, 1, 10))),
										"5",
									),
								),
							},
							"",
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
				bar =
					do macro
						foo += 2
					end
				nil
			`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(56, 6, 8))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(5, 2, 5), P(48, 5, 9))),
						ast.NewAssignmentExpressionNode(
							L(S(P(5, 2, 5), P(47, 5, 8))),
							T(L(S(P(9, 2, 9), P(9, 2, 9))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(5, 2, 5), P(7, 2, 7))), "bar"),
							ast.NewMacroBoundaryNode(
								L(S(P(16, 3, 6), P(47, 5, 8))),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										L(S(P(31, 4, 7), P(39, 4, 15))),
										ast.NewAssignmentExpressionNode(
											L(S(P(31, 4, 7), P(38, 4, 14))),
											T(L(S(P(35, 4, 11), P(36, 4, 12))), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(L(S(P(31, 4, 7), P(33, 4, 9))), "foo"),
											ast.NewIntLiteralNode(L(S(P(38, 4, 14), P(38, 4, 14))), "2"),
										),
									),
								},
								"",
							),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(53, 6, 5), P(56, 6, 8))),
						ast.NewNilLiteralNode(L(S(P(53, 6, 5), P(55, 6, 7)))),
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

func TestQuoteExpression(t *testing.T) {
	tests := testTable{
		"can have a multiline body": {
			input: `
quote
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(25, 5, 4))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(1, 2, 1), P(25, 5, 4))),
						ast.NewQuoteExpressionNode(
							L(S(P(1, 2, 1), P(24, 5, 3))),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(8, 3, 2), P(16, 3, 10))),
									ast.NewAssignmentExpressionNode(
										L(S(P(8, 3, 2), P(15, 3, 9))),
										T(L(S(P(12, 3, 6), P(13, 3, 7))), token.PLUS_EQUAL),
										ast.NewPublicIdentifierNode(L(S(P(8, 3, 2), P(10, 3, 4))), "foo"),
										ast.NewIntLiteralNode(L(S(P(15, 3, 9), P(15, 3, 9))), "2"),
									),
								),
								ast.NewExpressionStatementNode(
									L(S(P(18, 4, 2), P(21, 4, 5))),
									ast.NewNilLiteralNode(L(S(P(18, 4, 2), P(20, 4, 4)))),
								),
							},
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
				quote
				end
			`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(18, 3, 8))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(5, 2, 5), P(18, 3, 8))),
						ast.NewQuoteExpressionNode(
							L(S(P(5, 2, 5), P(17, 3, 7))),
							nil,
						),
					),
				},
			),
		},
		"can be a one-liner": {
			input: "quote 5",
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(0, 1, 1), P(6, 1, 7))),
						ast.NewQuoteExpressionNode(
							L(S(P(0, 1, 1), P(6, 1, 7))),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									L(S(P(6, 1, 7), P(6, 1, 7))),
									ast.NewIntLiteralNode(
										L(S(P(6, 1, 7), P(6, 1, 7))),
										"5",
									),
								),
							},
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
				bar =
					quote
						foo += 2
					end
				nil
			`,
			want: ast.NewProgramNode(
				L(S(P(0, 1, 1), P(53, 6, 8))),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(L(S(P(0, 1, 1), P(0, 1, 1)))),
					ast.NewExpressionStatementNode(
						L(S(P(5, 2, 5), P(45, 5, 9))),
						ast.NewAssignmentExpressionNode(
							L(S(P(5, 2, 5), P(44, 5, 8))),
							T(L(S(P(9, 2, 9), P(9, 2, 9))), token.EQUAL_OP),
							ast.NewPublicIdentifierNode(L(S(P(5, 2, 5), P(7, 2, 7))), "bar"),
							ast.NewQuoteExpressionNode(
								L(S(P(16, 3, 6), P(44, 5, 8))),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										L(S(P(28, 4, 7), P(36, 4, 15))),
										ast.NewAssignmentExpressionNode(
											L(S(P(28, 4, 7), P(35, 4, 14))),
											T(L(S(P(32, 4, 11), P(33, 4, 12))), token.PLUS_EQUAL),
											ast.NewPublicIdentifierNode(L(S(P(28, 4, 7), P(30, 4, 9))), "foo"),
											ast.NewIntLiteralNode(L(S(P(35, 4, 14), P(35, 4, 14))), "2"),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						L(S(P(50, 6, 5), P(53, 6, 8))),
						ast.NewNilLiteralNode(L(S(P(50, 6, 5), P(52, 6, 7)))),
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
