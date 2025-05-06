package parser

import (
	"testing"

	"github.com/elk-language/elk/parser/ast"
	"github.com/elk-language/elk/position/diagnostic"
	"github.com/elk-language/elk/token"
)

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
								"bar",
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
							"foo",
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
				L(S(P(0, 1, 1), P(6, 1, 7))),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						L(S(P(4, 1, 5), P(6, 1, 7))),
						ast.NewInvalidNode(
							L(S(P(4, 1, 5), P(6, 1, 7))),
							T(L(S(P(4, 1, 5), P(6, 1, 7))), token.COLON_COLON_LBRACKET),
						),
					),
				},
			),
			err: diagnostic.DiagnosticList{
				diagnostic.NewFailure(L(S(P(4, 1, 5), P(6, 1, 7))), "unexpected ::[, expected macro arguments"),
				diagnostic.NewFailure(L(S(P(7, 1, 8), P(12, 1, 13))), "unexpected PUBLIC_CONSTANT, expected a statement separator `\\n`, `;`"),
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
							"foo",
							[]ast.ExpressionNode{
								ast.NewClosureLiteralNode(
									L(S(P(7, 1, 8), P(18, 1, 19))),
									[]ast.ParameterNode{
										ast.NewFormalParameterNode(L(S(P(8, 1, 9), P(8, 1, 9))), "i", nil, nil, ast.NormalParameterKind),
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
							"foo",
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
							"foo",
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
		// "can omit the receiver and have a trailing closure": {
		// 	input: "foo(1, 5) |i| -> i * 2",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(21, 1, 22))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(21, 1, 22))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(21, 1, 22))),
		// 					"foo",
		// 					[]ast.ExpressionNode{
		// 						ast.NewIntLiteralNode(L(S(P(4, 1, 5), P(4, 1, 5))), "1"),
		// 						ast.NewIntLiteralNode(L(S(P(7, 1, 8), P(7, 1, 8))), "5"),
		// 						ast.NewClosureLiteralNode(
		// 							L(S(P(10, 1, 11), P(21, 1, 22))),
		// 							[]ast.ParameterNode{
		// 								ast.NewFormalParameterNode(L(S(P(11, 1, 12), P(11, 1, 12))), "i", nil, nil, ast.NormalParameterKind),
		// 							},
		// 							nil,
		// 							nil,
		// 							[]ast.StatementNode{
		// 								ast.NewExpressionStatementNode(
		// 									L(S(P(17, 1, 18), P(21, 1, 22))),
		// 									ast.NewBinaryExpressionNode(
		// 										L(S(P(17, 1, 18), P(21, 1, 22))),
		// 										T(L(S(P(19, 1, 20), P(19, 1, 20))), token.STAR),
		// 										ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "i"),
		// 										ast.NewIntLiteralNode(L(S(P(21, 1, 22), P(21, 1, 22))), "2"),
		// 									),
		// 								),
		// 							},
		// 						),
		// 					},
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can omit the receiver have named arguments and a trailing closure": {
		// 	input: "foo(f: 5) |i| -> i * 2",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(21, 1, 22))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(21, 1, 22))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(21, 1, 22))),
		// 					"foo",
		// 					nil,
		// 					[]ast.NamedArgumentNode{
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(4, 1, 5), P(7, 1, 8))),
		// 							"f",
		// 							ast.NewIntLiteralNode(L(S(P(7, 1, 8), P(7, 1, 8))), "5"),
		// 						),
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(10, 1, 11), P(21, 1, 22))),
		// 							"func",
		// 							ast.NewClosureLiteralNode(
		// 								L(S(P(10, 1, 11), P(21, 1, 22))),
		// 								[]ast.ParameterNode{
		// 									ast.NewFormalParameterNode(L(S(P(11, 1, 12), P(11, 1, 12))), "i", nil, nil, ast.NormalParameterKind),
		// 								},
		// 								nil,
		// 								nil,
		// 								[]ast.StatementNode{
		// 									ast.NewExpressionStatementNode(
		// 										L(S(P(17, 1, 18), P(21, 1, 22))),
		// 										ast.NewBinaryExpressionNode(
		// 											L(S(P(17, 1, 18), P(21, 1, 22))),
		// 											T(L(S(P(19, 1, 20), P(19, 1, 20))), token.STAR),
		// 											ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "i"),
		// 											ast.NewIntLiteralNode(L(S(P(21, 1, 22), P(21, 1, 22))), "2"),
		// 										),
		// 									),
		// 								},
		// 							),
		// 						),
		// 					},
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have a private name with an implicit receiver": {
		// 	input: "_foo()",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(5, 1, 6))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(5, 1, 6))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(5, 1, 6))),
		// 					"_foo",
		// 					nil,
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have an explicit receiver": {
		// 	input: "foo.bar()",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(8, 1, 9))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(8, 1, 9))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(8, 1, 9))),
		// 					ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 					T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
		// 					"bar",
		// 					nil,
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have an explicit receiver with type arguments": {
		// 	input: "foo.bar::[String]()",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(18, 1, 19))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(18, 1, 19))),
		// 				ast.NewGenericMethodCallNode(
		// 					L(S(P(0, 1, 1), P(18, 1, 19))),
		// 					ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 					T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
		// 					"bar",
		// 					[]ast.TypeNode{
		// 						ast.NewPublicConstantNode(L(S(P(10, 1, 11), P(15, 1, 16))), "String"),
		// 					},
		// 					nil,
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can be multiline, have an explicit receiver with type arguments": {
		// 	input: `
		// 		foo.bar::[
		// 			String,
		// 		]()
		// 	`,
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(36, 4, 8))),
		// 		[]ast.StatementNode{
		// 			ast.NewEmptyStatementNode(
		// 				L(S(P(0, 1, 1), P(0, 1, 1))),
		// 			),
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(5, 2, 5), P(36, 4, 8))),
		// 				ast.NewGenericMethodCallNode(
		// 					L(S(P(5, 2, 5), P(35, 4, 7))),
		// 					ast.NewPublicIdentifierNode(L(S(P(5, 2, 5), P(7, 2, 7))), "foo"),
		// 					T(L(S(P(8, 2, 8), P(8, 2, 8))), token.DOT),
		// 					"bar",
		// 					[]ast.TypeNode{
		// 						ast.NewPublicConstantNode(L(S(P(21, 3, 6), P(26, 3, 11))), "String"),
		// 					},
		// 					nil,
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have an explicit receiver and a trailing closure without pipes": {
		// 	input: "foo.bar() -> i * 2",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(17, 1, 18))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(17, 1, 18))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(17, 1, 18))),
		// 					ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 					T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
		// 					"bar",
		// 					[]ast.ExpressionNode{
		// 						ast.NewClosureLiteralNode(
		// 							L(S(P(10, 1, 11), P(17, 1, 18))),
		// 							nil,
		// 							nil,
		// 							nil,
		// 							[]ast.StatementNode{
		// 								ast.NewExpressionStatementNode(
		// 									L(S(P(13, 1, 14), P(17, 1, 18))),
		// 									ast.NewBinaryExpressionNode(
		// 										L(S(P(13, 1, 14), P(17, 1, 18))),
		// 										T(L(S(P(15, 1, 16), P(15, 1, 16))), token.STAR),
		// 										ast.NewPublicIdentifierNode(L(S(P(13, 1, 14), P(13, 1, 14))), "i"),
		// 										ast.NewIntLiteralNode(L(S(P(17, 1, 18), P(17, 1, 18))), "2"),
		// 									),
		// 								),
		// 							},
		// 						),
		// 					},
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have an explicit receiver and a trailing closure without arguments": {
		// 	input: "foo.bar() || -> i * 2",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(20, 1, 21))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(20, 1, 21))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(20, 1, 21))),
		// 					ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 					T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
		// 					"bar",
		// 					[]ast.ExpressionNode{
		// 						ast.NewClosureLiteralNode(
		// 							L(S(P(10, 1, 11), P(20, 1, 21))),
		// 							nil,
		// 							nil,
		// 							nil,
		// 							[]ast.StatementNode{
		// 								ast.NewExpressionStatementNode(
		// 									L(S(P(16, 1, 17), P(20, 1, 21))),
		// 									ast.NewBinaryExpressionNode(
		// 										L(S(P(16, 1, 17), P(20, 1, 21))),
		// 										T(L(S(P(18, 1, 19), P(18, 1, 19))), token.STAR),
		// 										ast.NewPublicIdentifierNode(L(S(P(16, 1, 17), P(16, 1, 17))), "i"),
		// 										ast.NewIntLiteralNode(L(S(P(20, 1, 21), P(20, 1, 21))), "2"),
		// 									),
		// 								),
		// 							},
		// 						),
		// 					},
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have an explicit receiver and a trailing closure": {
		// 	input: "foo.bar() |i| -> i * 2",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(21, 1, 22))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(21, 1, 22))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(21, 1, 22))),
		// 					ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 					T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
		// 					"bar",
		// 					[]ast.ExpressionNode{
		// 						ast.NewClosureLiteralNode(
		// 							L(S(P(10, 1, 11), P(21, 1, 22))),
		// 							[]ast.ParameterNode{
		// 								ast.NewFormalParameterNode(L(S(P(11, 1, 12), P(11, 1, 12))), "i", nil, nil, ast.NormalParameterKind),
		// 							},
		// 							nil,
		// 							nil,
		// 							[]ast.StatementNode{
		// 								ast.NewExpressionStatementNode(
		// 									L(S(P(17, 1, 18), P(21, 1, 22))),
		// 									ast.NewBinaryExpressionNode(
		// 										L(S(P(17, 1, 18), P(21, 1, 22))),
		// 										T(L(S(P(19, 1, 20), P(19, 1, 20))), token.STAR),
		// 										ast.NewPublicIdentifierNode(L(S(P(17, 1, 18), P(17, 1, 18))), "i"),
		// 										ast.NewIntLiteralNode(L(S(P(21, 1, 22), P(21, 1, 22))), "2"),
		// 									),
		// 								),
		// 							},
		// 						),
		// 					},
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have an explicit receiver, arguments and a trailing closure": {
		// 	input: "foo.bar(1, 5) |i| -> i * 2",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(25, 1, 26))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(25, 1, 26))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(25, 1, 26))),
		// 					ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 					T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
		// 					"bar",
		// 					[]ast.ExpressionNode{
		// 						ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "1"),
		// 						ast.NewIntLiteralNode(L(S(P(11, 1, 12), P(11, 1, 12))), "5"),
		// 						ast.NewClosureLiteralNode(
		// 							L(S(P(14, 1, 15), P(25, 1, 26))),
		// 							[]ast.ParameterNode{
		// 								ast.NewFormalParameterNode(L(S(P(15, 1, 16), P(15, 1, 16))), "i", nil, nil, ast.NormalParameterKind),
		// 							},
		// 							nil,
		// 							nil,
		// 							[]ast.StatementNode{
		// 								ast.NewExpressionStatementNode(
		// 									L(S(P(21, 1, 22), P(25, 1, 26))),
		// 									ast.NewBinaryExpressionNode(
		// 										L(S(P(21, 1, 22), P(25, 1, 26))),
		// 										T(L(S(P(23, 1, 24), P(23, 1, 24))), token.STAR),
		// 										ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(21, 1, 22))), "i"),
		// 										ast.NewIntLiteralNode(L(S(P(25, 1, 26), P(25, 1, 26))), "2"),
		// 									),
		// 								),
		// 							},
		// 						),
		// 					},
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have an explicit receiver, named arguments and a trailing closure": {
		// 	input: "foo.bar(f: 5) |i| -> i * 2",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(25, 1, 26))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(25, 1, 26))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(25, 1, 26))),
		// 					ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 					T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
		// 					"bar",
		// 					nil,
		// 					[]ast.NamedArgumentNode{
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(8, 1, 9), P(11, 1, 12))),
		// 							"f",
		// 							ast.NewIntLiteralNode(L(S(P(11, 1, 12), P(11, 1, 12))), "5"),
		// 						),
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(14, 1, 15), P(25, 1, 26))),
		// 							"func",
		// 							ast.NewClosureLiteralNode(
		// 								L(S(P(14, 1, 15), P(25, 1, 26))),
		// 								[]ast.ParameterNode{
		// 									ast.NewFormalParameterNode(L(S(P(15, 1, 16), P(15, 1, 16))), "i", nil, nil, ast.NormalParameterKind),
		// 								},
		// 								nil,
		// 								nil,
		// 								[]ast.StatementNode{
		// 									ast.NewExpressionStatementNode(
		// 										L(S(P(21, 1, 22), P(25, 1, 26))),
		// 										ast.NewBinaryExpressionNode(
		// 											L(S(P(21, 1, 22), P(25, 1, 26))),
		// 											T(L(S(P(23, 1, 24), P(23, 1, 24))), token.STAR),
		// 											ast.NewPublicIdentifierNode(L(S(P(21, 1, 22), P(21, 1, 22))), "i"),
		// 											ast.NewIntLiteralNode(L(S(P(25, 1, 26), P(25, 1, 26))), "2"),
		// 										),
		// 									),
		// 								},
		// 							),
		// 						),
		// 					},
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can omit parentheses": {
		// 	input: "foo.bar 1",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(8, 1, 9))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(8, 1, 9))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(8, 1, 9))),
		// 					ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 					T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
		// 					"bar",
		// 					[]ast.ExpressionNode{
		// 						ast.NewIntLiteralNode(L(S(P(8, 1, 9), P(8, 1, 9))), "1"),
		// 					},
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can use the safe navigation operator": {
		// 	input: "foo?.bar",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(7, 1, 8))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(7, 1, 8))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(7, 1, 8))),
		// 					ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 					T(L(S(P(3, 1, 4), P(4, 1, 5))), token.QUESTION_DOT),
		// 					"bar",
		// 					nil,
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can use the cascade call operator": {
		// 	input: "foo..bar",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(7, 1, 8))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(7, 1, 8))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(7, 1, 8))),
		// 					ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 					T(L(S(P(3, 1, 4), P(4, 1, 5))), token.DOT_DOT),
		// 					"bar",
		// 					nil,
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can use the safe cascade call operator": {
		// 	input: "foo?..bar",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(8, 1, 9))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(8, 1, 9))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(8, 1, 9))),
		// 					ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 					T(L(S(P(3, 1, 4), P(5, 1, 6))), token.QUESTION_DOT_DOT),
		// 					"bar",
		// 					nil,
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can be nested with parentheses": {
		// 	input: "foo.bar().baz()",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(14, 1, 15))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(14, 1, 15))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(14, 1, 15))),
		// 					ast.NewMethodCallNode(
		// 						L(S(P(0, 1, 1), P(8, 1, 9))),
		// 						ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 						T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
		// 						"bar",
		// 						nil,
		// 						nil,
		// 					),
		// 					T(L(S(P(9, 1, 10), P(9, 1, 10))), token.DOT),
		// 					"baz",
		// 					nil,
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have any expression as the receiver": {
		// 	input: "(foo + 2).bar()",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(14, 1, 15))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(1, 1, 2), P(14, 1, 15))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(1, 1, 2), P(14, 1, 15))),
		// 					ast.NewBinaryExpressionNode(
		// 						L(S(P(1, 1, 2), P(7, 1, 8))),
		// 						T(L(S(P(5, 1, 6), P(5, 1, 6))), token.PLUS),
		// 						ast.NewPublicIdentifierNode(L(S(P(1, 1, 2), P(3, 1, 4))), "foo"),
		// 						ast.NewIntLiteralNode(L(S(P(7, 1, 8), P(7, 1, 8))), "2"),
		// 					),
		// 					T(L(S(P(9, 1, 10), P(9, 1, 10))), token.DOT),
		// 					"bar",
		// 					nil,
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "cannot call a private method on an explicit receiver": {
		// 	input: "foo._bar()",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(9, 1, 10))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(9, 1, 10))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(9, 1, 10))),
		// 					ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 					T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
		// 					"_bar",
		// 					nil,
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// 	err: diagnostic.DiagnosticList{
		// 		diagnostic.NewFailure(L(S(P(4, 1, 5), P(7, 1, 8))), "unexpected PRIVATE_IDENTIFIER, expected a public method name (public identifier, keyword or overridable operator)"),
		// 	},
		// },
		// "can have an overridable operator as the method name with an explicit receiver": {
		// 	input: "foo.+()",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(6, 1, 7))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(6, 1, 7))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(6, 1, 7))),
		// 					ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 					T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
		// 					"+",
		// 					nil,
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "cannot have a non overridable operator as the method name with an explicit receiver": {
		// 	input: "foo.&&()",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(7, 1, 8))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(7, 1, 8))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(7, 1, 8))),
		// 					ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 					T(L(S(P(3, 1, 4), P(3, 1, 4))), token.DOT),
		// 					"&&",
		// 					nil,
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// 	err: diagnostic.DiagnosticList{
		// 		diagnostic.NewFailure(L(S(P(4, 1, 5), P(5, 1, 6))), "unexpected &&, expected a public method name (public identifier, keyword or overridable operator)"),
		// 	},
		// },
		// "can call a private method on self": {
		// 	input: "self._foo()",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(10, 1, 11))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(10, 1, 11))),
		// 				ast.NewMethodCallNode(
		// 					L(S(P(0, 1, 1), P(10, 1, 11))),
		// 					ast.NewSelfLiteralNode(L(S(P(0, 1, 1), P(3, 1, 4)))),
		// 					T(L(S(P(4, 1, 5), P(4, 1, 5))), token.DOT),
		// 					"_foo",
		// 					nil,
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have positional arguments": {
		// 	input: "foo(.1, 'foo', :bar)",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(19, 1, 20))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(19, 1, 20))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(19, 1, 20))),
		// 					"foo",
		// 					[]ast.ExpressionNode{
		// 						ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
		// 						ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
		// 						ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
		// 					},
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have splat arguments": {
		// 	input: "foo(*baz, 'foo', *bar)",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(21, 1, 22))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(21, 1, 22))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(21, 1, 22))),
		// 					"foo",
		// 					[]ast.ExpressionNode{
		// 						ast.NewSplatExpressionNode(
		// 							L(S(P(4, 1, 5), P(7, 1, 8))),
		// 							ast.NewPublicIdentifierNode(
		// 								L(S(P(5, 1, 6), P(7, 1, 8))),
		// 								"baz",
		// 							),
		// 						),
		// 						ast.NewRawStringLiteralNode(L(S(P(10, 1, 11), P(14, 1, 15))), "foo"),
		// 						ast.NewSplatExpressionNode(
		// 							L(S(P(17, 1, 18), P(20, 1, 21))),
		// 							ast.NewPublicIdentifierNode(
		// 								L(S(P(18, 1, 19), P(20, 1, 21))),
		// 								"bar",
		// 							),
		// 						),
		// 					},
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have a trailing comma": {
		// 	input: "foo(.1, 'foo', :bar,)",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(20, 1, 21))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(20, 1, 21))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(20, 1, 21))),
		// 					"foo",
		// 					[]ast.ExpressionNode{
		// 						ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
		// 						ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
		// 						ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
		// 					},
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have named arguments": {
		// 	input: "foo(bar: :baz, elk: true)",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(24, 1, 25))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(24, 1, 25))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(24, 1, 25))),
		// 					"foo",
		// 					nil,
		// 					[]ast.NamedArgumentNode{
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(4, 1, 5), P(12, 1, 13))),
		// 							"bar",
		// 							ast.NewSimpleSymbolLiteralNode(L(S(P(9, 1, 10), P(12, 1, 13))), "baz"),
		// 						),
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(15, 1, 16), P(23, 1, 24))),
		// 							"elk",
		// 							ast.NewTrueLiteralNode(L(S(P(20, 1, 21), P(23, 1, 24)))),
		// 						),
		// 					},
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have double splat arguments": {
		// 	input: "foo(**f, bar: :baz, **dupa(), elk: true)",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(39, 1, 40))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(39, 1, 40))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(39, 1, 40))),
		// 					"foo",
		// 					nil,
		// 					[]ast.NamedArgumentNode{
		// 						ast.NewDoubleSplatExpressionNode(
		// 							L(S(P(4, 1, 5), P(6, 1, 7))),
		// 							ast.NewPublicIdentifierNode(L(S(P(6, 1, 7), P(6, 1, 7))), "f"),
		// 						),
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(9, 1, 10), P(17, 1, 18))),
		// 							"bar",
		// 							ast.NewSimpleSymbolLiteralNode(L(S(P(14, 1, 15), P(17, 1, 18))), "baz"),
		// 						),
		// 						ast.NewDoubleSplatExpressionNode(
		// 							L(S(P(20, 1, 21), P(27, 1, 28))),
		// 							ast.NewReceiverlessMethodCallNode(
		// 								L(S(P(22, 1, 23), P(27, 1, 28))),
		// 								"dupa",
		// 								nil,
		// 								nil,
		// 							),
		// 						),
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(30, 1, 31), P(38, 1, 39))),
		// 							"elk",
		// 							ast.NewTrueLiteralNode(L(S(P(35, 1, 36), P(38, 1, 39)))),
		// 						),
		// 					},
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have positional and named arguments": {
		// 	input: "foo(.1, 'foo', :bar, bar: :baz, elk: true)",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(41, 1, 42))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(41, 1, 42))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(41, 1, 42))),
		// 					"foo",
		// 					[]ast.ExpressionNode{
		// 						ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
		// 						ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
		// 						ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
		// 					},
		// 					[]ast.NamedArgumentNode{
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(21, 1, 22), P(29, 1, 30))),
		// 							"bar",
		// 							ast.NewSimpleSymbolLiteralNode(L(S(P(26, 1, 27), P(29, 1, 30))), "baz"),
		// 						),
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(32, 1, 33), P(40, 1, 41))),
		// 							"elk",
		// 							ast.NewTrueLiteralNode(L(S(P(37, 1, 38), P(40, 1, 41)))),
		// 						),
		// 					},
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have newlines after commas": {
		// 	input: "foo(.1,\n'foo',\n:bar, bar: :baz,\nelk: true)",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(41, 4, 10))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(41, 4, 10))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(41, 4, 10))),
		// 					"foo",
		// 					[]ast.ExpressionNode{
		// 						ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
		// 						ast.NewRawStringLiteralNode(L(S(P(8, 2, 1), P(12, 2, 5))), "foo"),
		// 						ast.NewSimpleSymbolLiteralNode(L(S(P(15, 3, 1), P(18, 3, 4))), "bar"),
		// 					},
		// 					[]ast.NamedArgumentNode{
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(21, 3, 7), P(29, 3, 15))),
		// 							"bar",
		// 							ast.NewSimpleSymbolLiteralNode(L(S(P(26, 3, 12), P(29, 3, 15))), "baz"),
		// 						),
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(32, 4, 1), P(40, 4, 9))),
		// 							"elk",
		// 							ast.NewTrueLiteralNode(L(S(P(37, 4, 6), P(40, 4, 9)))),
		// 						),
		// 					},
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have newlines around parentheses": {
		// 	input: "foo(\n.1, 'foo', :bar, bar: :baz, elk: true\n)",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(43, 3, 1))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(43, 3, 1))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(43, 3, 1))),
		// 					"foo",
		// 					[]ast.ExpressionNode{
		// 						ast.NewFloatLiteralNode(L(S(P(5, 2, 1), P(6, 2, 2))), "0.1"),
		// 						ast.NewRawStringLiteralNode(L(S(P(9, 2, 5), P(13, 2, 9))), "foo"),
		// 						ast.NewSimpleSymbolLiteralNode(L(S(P(16, 2, 12), P(19, 2, 15))), "bar"),
		// 					},
		// 					[]ast.NamedArgumentNode{
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(22, 2, 18), P(30, 2, 26))),
		// 							"bar",
		// 							ast.NewSimpleSymbolLiteralNode(L(S(P(27, 2, 23), P(30, 2, 26))), "baz"),
		// 						),
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(33, 2, 29), P(41, 2, 37))),
		// 							"elk",
		// 							ast.NewTrueLiteralNode(L(S(P(38, 2, 34), P(41, 2, 37)))),
		// 						),
		// 					},
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "cannot have newlines before the opening parenthesis": {
		// 	input: "foo\n(.1, 'foo', :bar, bar: :baz, elk: true)",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(42, 2, 39))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(3, 1, 4))),
		// 				ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 			),
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(5, 2, 2), P(42, 2, 39))),
		// 				ast.NewFloatLiteralNode(L(S(P(5, 2, 2), P(6, 2, 3))), "0.1"),
		// 			),
		// 		},
		// 	),
		// 	err: diagnostic.DiagnosticList{
		// 		diagnostic.NewFailure(L(S(P(7, 2, 4), P(7, 2, 4))), "unexpected ,, expected )"),
		// 	},
		// },
		// "can have positional arguments without parentheses": {
		// 	input: "foo .1, 'foo', :bar",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(18, 1, 19))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(18, 1, 19))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(18, 1, 19))),
		// 					"foo",
		// 					[]ast.ExpressionNode{
		// 						ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
		// 						ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
		// 						ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
		// 					},
		// 					nil,
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have named arguments without parentheses": {
		// 	input: "foo bar: :baz, elk: true",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(23, 1, 24))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(23, 1, 24))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(23, 1, 24))),
		// 					"foo",
		// 					nil,
		// 					[]ast.NamedArgumentNode{
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(4, 1, 5), P(12, 1, 13))),
		// 							"bar",
		// 							ast.NewSimpleSymbolLiteralNode(L(S(P(9, 1, 10), P(12, 1, 13))), "baz"),
		// 						),
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(15, 1, 16), P(23, 1, 24))),
		// 							"elk",
		// 							ast.NewTrueLiteralNode(L(S(P(20, 1, 21), P(23, 1, 24)))),
		// 						),
		// 					},
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have positional and named arguments without parentheses": {
		// 	input: "foo .1, 'foo', :bar, bar: :baz, elk: true",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(40, 1, 41))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(40, 1, 41))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(40, 1, 41))),
		// 					"foo",
		// 					[]ast.ExpressionNode{
		// 						ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
		// 						ast.NewRawStringLiteralNode(L(S(P(8, 1, 9), P(12, 1, 13))), "foo"),
		// 						ast.NewSimpleSymbolLiteralNode(L(S(P(15, 1, 16), P(18, 1, 19))), "bar"),
		// 					},
		// 					[]ast.NamedArgumentNode{
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(21, 1, 22), P(29, 1, 30))),
		// 							"bar",
		// 							ast.NewSimpleSymbolLiteralNode(L(S(P(26, 1, 27), P(29, 1, 30))), "baz"),
		// 						),
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(32, 1, 33), P(40, 1, 41))),
		// 							"elk",
		// 							ast.NewTrueLiteralNode(L(S(P(37, 1, 38), P(40, 1, 41)))),
		// 						),
		// 					},
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "can have newlines after commas without parentheses": {
		// 	input: "foo .1,\n'foo',\n:bar, bar: :baz,\nelk: true",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(40, 4, 9))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(40, 4, 9))),
		// 				ast.NewReceiverlessMethodCallNode(
		// 					L(S(P(0, 1, 1), P(40, 4, 9))),
		// 					"foo",
		// 					[]ast.ExpressionNode{
		// 						ast.NewFloatLiteralNode(L(S(P(4, 1, 5), P(5, 1, 6))), "0.1"),
		// 						ast.NewRawStringLiteralNode(L(S(P(8, 2, 1), P(12, 2, 5))), "foo"),
		// 						ast.NewSimpleSymbolLiteralNode(L(S(P(15, 3, 1), P(18, 3, 4))), "bar"),
		// 					},
		// 					[]ast.NamedArgumentNode{
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(21, 3, 7), P(29, 3, 15))),
		// 							"bar",
		// 							ast.NewSimpleSymbolLiteralNode(L(S(P(26, 3, 12), P(29, 3, 15))), "baz"),
		// 						),
		// 						ast.NewNamedCallArgumentNode(
		// 							L(S(P(32, 4, 1), P(40, 4, 9))),
		// 							"elk",
		// 							ast.NewTrueLiteralNode(L(S(P(37, 4, 6), P(40, 4, 9)))),
		// 						),
		// 					},
		// 				),
		// 			),
		// 		},
		// 	),
		// },
		// "cannot have newlines before the arguments without parentheses": {
		// 	input: "foo\n.1, 'foo', :bar, bar: :baz, elk: true",
		// 	want: ast.NewProgramNode(
		// 		L(S(P(0, 1, 1), P(5, 2, 2))),
		// 		[]ast.StatementNode{
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(0, 1, 1), P(3, 1, 4))),
		// 				ast.NewPublicIdentifierNode(L(S(P(0, 1, 1), P(2, 1, 3))), "foo"),
		// 			),
		// 			ast.NewExpressionStatementNode(
		// 				L(S(P(4, 2, 1), P(5, 2, 2))),
		// 				ast.NewFloatLiteralNode(L(S(P(4, 2, 1), P(5, 2, 2))), "0.1"),
		// 			),
		// 		},
		// 	),
		// 	err: diagnostic.DiagnosticList{
		// 		diagnostic.NewFailure(L(S(P(6, 2, 3), P(6, 2, 3))), "unexpected ,, expected a statement separator `\\n`, `;`"),
		// 	},
		// },
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
