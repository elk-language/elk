package checker

import (
	"testing"

	"github.com/elk-language/elk/position/errors"
	"github.com/elk-language/elk/token"
	"github.com/elk-language/elk/types"
	"github.com/elk-language/elk/types/ast"
	"github.com/elk-language/elk/value/symbol"
)

func TestStringLiteral(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"infer raw string": {
			input: "var foo = 'str'",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewRawStringLiteralNode(
								S(P(10, 1, 11), P(14, 1, 15)),
								"str",
								types.NewStringLiteral("str"),
							),
							globalEnv.StdSubtype(symbol.String),
						),
					),
				},
			),
		},
		"assign string literal to String": {
			input: "var foo: String = 'str'",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(14, 1, 15)),
								"String",
								globalEnv.StdSubtype(symbol.String),
							),
							ast.NewRawStringLiteralNode(
								S(P(18, 1, 19), P(22, 1, 23)),
								"str",
								types.NewStringLiteral("str"),
							),
							globalEnv.StdSubtype(symbol.String),
						),
					),
				},
			),
		},
		"assign string literal to matching literal type": {
			input: "var foo: 'str' = 'str'",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewRawStringLiteralNode(
								S(P(9, 1, 10), P(13, 1, 14)),
								"str",
								types.NewStringLiteral("str"),
							),
							ast.NewRawStringLiteralNode(
								S(P(17, 1, 18), P(21, 1, 22)),
								"str",
								types.NewStringLiteral("str"),
							),
							types.NewStringLiteral("str"),
						),
					),
				},
			),
		},
		"assign string literal to non matching literal type": {
			input: "var foo: 'str' = 'foo'",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewRawStringLiteralNode(
								S(P(9, 1, 10), P(13, 1, 14)),
								"str",
								types.NewStringLiteral("str"),
							),
							ast.NewRawStringLiteralNode(
								S(P(17, 1, 18), P(21, 1, 22)),
								"foo",
								types.NewStringLiteral("foo"),
							),
							types.NewStringLiteral("str"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(17, 1, 18), P(21, 1, 22)), "type `Std::String(\"foo\")` cannot be assigned to type `Std::String(\"str\")`"),
			},
		},
		"infer double quoted string": {
			input: `var foo = "str"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewDoubleQuotedStringLiteralNode(
								S(P(10, 1, 11), P(14, 1, 15)),
								"str",
								types.NewStringLiteral("str"),
							),
							globalEnv.StdSubtype(symbol.String),
						),
					),
				},
			),
		},
		"infer interpolated string": {
			input: `var foo = "${1} str #{5.2}"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(26, 1, 27)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(26, 1, 27)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(26, 1, 27)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewInterpolatedStringLiteralNode(
								S(P(10, 1, 11), P(26, 1, 27)),
								[]ast.StringLiteralContentNode{
									ast.NewStringInterpolationNode(
										S(P(11, 1, 12), P(14, 1, 15)),
										ast.NewIntLiteralNode(
											S(P(13, 1, 14), P(13, 1, 14)),
											"1",
											types.NewIntLiteral("1"),
										),
									),
									ast.NewStringLiteralContentSectionNode(
										S(P(15, 1, 16), P(19, 1, 20)),
										" str ",
									),
									ast.NewStringInspectInterpolationNode(
										S(P(20, 1, 21), P(25, 1, 26)),
										ast.NewFloatLiteralNode(
											S(P(22, 1, 23), P(24, 1, 25)),
											"5.2",
											types.NewFloatLiteral("5.2"),
										),
									),
								},
							),
							globalEnv.StdSubtype(symbol.String),
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

func TestSymbolLiteral(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"infer simple symbol": {
			input: "var foo = :str",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewSimpleSymbolLiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								"str",
								types.NewSymbolLiteral("str"),
							),
							globalEnv.StdSubtype(symbol.Symbol),
						),
					),
				},
			),
		},
		"infer double quoted symbol": {
			input: `var foo = :"str"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(15, 1, 16)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(15, 1, 16)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(15, 1, 16)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewSimpleSymbolLiteralNode(
								S(P(10, 1, 11), P(15, 1, 16)),
								"str",
								types.NewSymbolLiteral("str"),
							),
							globalEnv.StdSubtype(symbol.Symbol),
						),
					),
				},
			),
		},
		"infer interpolated symbol": {
			input: `var foo = :"${1} str #{5.2}"`,
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(27, 1, 28)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(27, 1, 28)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(27, 1, 28)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewInterpolatedSymbolLiteralNode(
								S(P(10, 1, 11), P(27, 1, 28)),
								ast.NewInterpolatedStringLiteralNode(
									S(P(11, 1, 12), P(27, 1, 28)),
									[]ast.StringLiteralContentNode{
										ast.NewStringInterpolationNode(
											S(P(12, 1, 13), P(15, 1, 16)),
											ast.NewIntLiteralNode(
												S(P(14, 1, 15), P(14, 1, 15)),
												"1",
												types.NewIntLiteral("1"),
											),
										),
										ast.NewStringLiteralContentSectionNode(
											S(P(16, 1, 17), P(20, 1, 21)),
											" str ",
										),
										ast.NewStringInspectInterpolationNode(
											S(P(21, 1, 22), P(26, 1, 27)),
											ast.NewFloatLiteralNode(
												S(P(23, 1, 24), P(25, 1, 26)),
												"5.2",
												types.NewFloatLiteral("5.2"),
											),
										),
									},
								),
							),
							globalEnv.StdSubtype(symbol.Symbol),
						),
					),
				},
			),
		},
		"assign symbol literal to Symbol": {
			input: "var foo: Symbol = :symb",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(14, 1, 15)),
								"Symbol",
								globalEnv.StdSubtype(symbol.Symbol),
							),
							ast.NewSimpleSymbolLiteralNode(
								S(P(18, 1, 19), P(22, 1, 23)),
								"symb",
								types.NewSymbolLiteral("symb"),
							),
							globalEnv.StdSubtype(symbol.Symbol),
						),
					),
				},
			),
		},
		"assign symbol literal to matching literal type": {
			input: "var foo: :symb = :symb",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewSimpleSymbolLiteralNode(
								S(P(9, 1, 10), P(13, 1, 14)),
								"symb",
								types.NewSymbolLiteral("symb"),
							),
							ast.NewSimpleSymbolLiteralNode(
								S(P(17, 1, 18), P(21, 1, 22)),
								"symb",
								types.NewSymbolLiteral("symb"),
							),
							types.NewSymbolLiteral("symb"),
						),
					),
				},
			),
		},
		"assign symbol literal to non matching literal type": {
			input: "var foo: :symb = :foob",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewSimpleSymbolLiteralNode(
								S(P(9, 1, 10), P(13, 1, 14)),
								"symb",
								types.NewSymbolLiteral("symb"),
							),
							ast.NewSimpleSymbolLiteralNode(
								S(P(17, 1, 18), P(21, 1, 22)),
								"foob",
								types.NewSymbolLiteral("foob"),
							),
							types.NewSymbolLiteral("symb"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(17, 1, 18), P(21, 1, 22)), "type `Std::Symbol(:foob)` cannot be assigned to type `Std::Symbol(:symb)`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, false)
		})
	}
}

func TestCharLiteral(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"infer raw char": {
			input: "var foo = r`s`",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewRawCharLiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								's',
								types.NewCharLiteral('s'),
							),
							globalEnv.StdSubtype(symbol.Char),
						),
					),
				},
			),
		},
		"infer char": {
			input: "var foo = `\\n`",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewCharLiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								'\n',
								types.NewCharLiteral('\n'),
							),
							globalEnv.StdSubtype(symbol.Char),
						),
					),
				},
			),
		},
		"assign char literal to Char": {
			input: "var foo: Char = `f`",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(18, 1, 19)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(18, 1, 19)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(18, 1, 19)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(12, 1, 13)),
								"Char",
								globalEnv.StdSubtype(symbol.Char),
							),
							ast.NewCharLiteralNode(
								S(P(16, 1, 17), P(18, 1, 19)),
								'f',
								types.NewCharLiteral('f'),
							),
							globalEnv.StdSubtype(symbol.Char),
						),
					),
				},
			),
		},
		"assign char literal to matching literal type": {
			input: "var foo: `f` = `f`",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewCharLiteralNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								'f',
								types.NewCharLiteral('f'),
							),
							ast.NewCharLiteralNode(
								S(P(15, 1, 16), P(17, 1, 18)),
								'f',
								types.NewCharLiteral('f'),
							),
							types.NewCharLiteral('f'),
						),
					),
				},
			),
		},
		"assign char literal to non matching literal type": {
			input: "var foo: `b` = `f`",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewCharLiteralNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								'b',
								types.NewCharLiteral('b'),
							),
							ast.NewCharLiteralNode(
								S(P(15, 1, 16), P(17, 1, 18)),
								'f',
								types.NewCharLiteral('f'),
							),
							types.NewCharLiteral('b'),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(15, 1, 16), P(17, 1, 18)), "type `Std::Char(`f`)` cannot be assigned to type `Std::Char(`b`)`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, false)
		})
	}
}

func TestIntLiteral(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"infer int": {
			input: "var foo = 1234",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewIntLiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								"1234",
								types.NewIntLiteral("1234"),
							),
							globalEnv.StdSubtype(symbol.Int),
						),
					),
				},
			),
		},
		"assign int literal to Int": {
			input: "var foo: Int = 12345678",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(11, 1, 12)),
								"Int",
								globalEnv.StdSubtype(symbol.Int),
							),
							ast.NewIntLiteralNode(
								S(P(15, 1, 16), P(22, 1, 23)),
								"12345678",
								types.NewIntLiteral("12345678"),
							),
							globalEnv.StdSubtype(symbol.Int),
						),
					),
				},
			),
		},
		"assign int literal to matching literal type": {
			input: "var foo: 12345 = 12345",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewIntLiteralNode(
								S(P(9, 1, 10), P(13, 1, 14)),
								"12345",
								types.NewIntLiteral("12345"),
							),
							ast.NewIntLiteralNode(
								S(P(17, 1, 18), P(21, 1, 22)),
								"12345",
								types.NewIntLiteral("12345"),
							),
							types.NewIntLiteral("12345"),
						),
					),
				},
			),
		},
		"assign int literal to non matching literal type": {
			input: "var foo: 23456 = 12345",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewIntLiteralNode(
								S(P(9, 1, 10), P(13, 1, 14)),
								"23456",
								types.NewIntLiteral("23456"),
							),
							ast.NewIntLiteralNode(
								S(P(17, 1, 18), P(21, 1, 22)),
								"12345",
								types.NewIntLiteral("12345"),
							),
							types.NewIntLiteral("23456"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(17, 1, 18), P(21, 1, 22)), "type `Std::Int(12345)` cannot be assigned to type `Std::Int(23456)`"),
			},
		},
		"infer int64": {
			input: "var foo = 1i64",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewInt64LiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								"1",
								types.NewInt64Literal("1"),
							),
							globalEnv.StdSubtype(symbol.Int64),
						),
					),
				},
			),
		},
		"infer int32": {
			input: "var foo = 1i32",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewInt32LiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								"1",
								types.NewInt32Literal("1"),
							),
							globalEnv.StdSubtype(symbol.Int32),
						),
					),
				},
			),
		},
		"infer int16": {
			input: "var foo = 1i16",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewInt16LiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								"1",
								types.NewInt16Literal("1"),
							),
							globalEnv.StdSubtype(symbol.Int16),
						),
					),
				},
			),
		},
		"infer int8": {
			input: "var foo = 12i8",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewInt8LiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								"12",
								types.NewInt8Literal("12"),
							),
							globalEnv.StdSubtype(symbol.Int8),
						),
					),
				},
			),
		},
		"infer uint64": {
			input: "var foo = 1u64",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewUInt64LiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								"1",
								types.NewUInt64Literal("1"),
							),
							globalEnv.StdSubtype(symbol.UInt64),
						),
					),
				},
			),
		},
		"infer uint32": {
			input: "var foo = 1u32",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewUInt32LiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								"1",
								types.NewUInt32Literal("1"),
							),
							globalEnv.StdSubtype(symbol.UInt32),
						),
					),
				},
			),
		},
		"infer uint16": {
			input: "var foo = 1u16",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewUInt16LiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								"1",
								types.NewUInt16Literal("1"),
							),
							globalEnv.StdSubtype(symbol.UInt16),
						),
					),
				},
			),
		},
		"infer uint8": {
			input: "var foo = 12u8",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewUInt8LiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								"12",
								types.NewUInt8Literal("12"),
							),
							globalEnv.StdSubtype(symbol.UInt8),
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

func TestFloatLiteral(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"infer float": {
			input: "var foo = 12.5",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewFloatLiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								"12.5",
								types.NewFloatLiteral("12.5"),
							),
							globalEnv.StdSubtype(symbol.Float),
						),
					),
				},
			),
		},
		"assign float literal to Float": {
			input: "var foo: Float = 1234.6",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(22, 1, 23)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(22, 1, 23)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(22, 1, 23)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(13, 1, 14)),
								"Float",
								globalEnv.StdSubtype(symbol.Float),
							),
							ast.NewFloatLiteralNode(
								S(P(17, 1, 18), P(22, 1, 23)),
								"1234.6",
								types.NewFloatLiteral("1234.6"),
							),
							globalEnv.StdSubtype(symbol.Float),
						),
					),
				},
			),
		},
		"assign float literal to matching literal type": {
			input: "var foo: 12.45 = 12.45",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewFloatLiteralNode(
								S(P(9, 1, 10), P(13, 1, 14)),
								"12.45",
								types.NewFloatLiteral("12.45"),
							),
							ast.NewFloatLiteralNode(
								S(P(17, 1, 18), P(21, 1, 22)),
								"12.45",
								types.NewFloatLiteral("12.45"),
							),
							types.NewFloatLiteral("12.45"),
						),
					),
				},
			),
		},
		"assign Float literal to non matching literal type": {
			input: "var foo: 23.56 = 12.45",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewFloatLiteralNode(
								S(P(9, 1, 10), P(13, 1, 14)),
								"23.56",
								types.NewFloatLiteral("23.56"),
							),
							ast.NewFloatLiteralNode(
								S(P(17, 1, 18), P(21, 1, 22)),
								"12.45",
								types.NewFloatLiteral("12.45"),
							),
							types.NewFloatLiteral("23.56"),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(17, 1, 18), P(21, 1, 22)), "type `Std::Float(12.45)` cannot be assigned to type `Std::Float(23.56)`"),
			},
		},
		"infer float64": {
			input: "var foo = 1f64",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewFloat64LiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								"1",
								types.NewFloat64Literal("1"),
							),
							globalEnv.StdSubtype(symbol.Float64),
						),
					),
				},
			),
		},
		"infer float32": {
			input: "var foo = 1f32",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewFloat32LiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								"1",
								types.NewFloat32Literal("1"),
							),
							globalEnv.StdSubtype(symbol.Float32),
						),
					),
				},
			),
		},
		"infer big float": {
			input: "var foo = 12bf",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewBigFloatLiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
								"12",
								types.NewBigFloatLiteral("12"),
							),
							globalEnv.StdSubtype(symbol.BigFloat),
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

func TestBoolLiteral(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"infer true": {
			input: "var foo = true",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(13, 1, 14)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(13, 1, 14)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(13, 1, 14)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewTrueLiteralNode(
								S(P(10, 1, 11), P(13, 1, 14)),
							),
							globalEnv.StdSubtype(symbol.True),
						),
					),
				},
			),
		},
		"assign true literal to True": {
			input: "var foo: True = true",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(12, 1, 13)),
								"True",
								globalEnv.StdSubtype(symbol.True),
							),
							ast.NewTrueLiteralNode(
								S(P(16, 1, 17), P(19, 1, 20)),
							),
							globalEnv.StdSubtype(symbol.True),
						),
					),
				},
			),
		},
		"assign true literal to Bool": {
			input: "var foo: Bool = true",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(12, 1, 13)),
								"Bool",
								globalEnv.StdSubtype(symbol.Bool),
							),
							ast.NewTrueLiteralNode(
								S(P(16, 1, 17), P(19, 1, 20)),
							),
							globalEnv.StdSubtype(symbol.Bool),
						),
					),
				},
			),
		},
		"assign true literal to matching literal type": {
			input: "var foo: true = true",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(19, 1, 20)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(19, 1, 20)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(19, 1, 20)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewTrueLiteralNode(
								S(P(9, 1, 10), P(12, 1, 13)),
							),
							ast.NewTrueLiteralNode(
								S(P(16, 1, 17), P(19, 1, 20)),
							),
							globalEnv.StdSubtype(symbol.True),
						),
					),
				},
			),
		},
		"assign true literal to non matching literal type": {
			input: "var foo: false = true",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewFalseLiteralNode(
								S(P(9, 1, 10), P(13, 1, 14)),
							),
							ast.NewTrueLiteralNode(
								S(P(17, 1, 18), P(20, 1, 21)),
							),
							globalEnv.StdSubtype(symbol.False),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(17, 1, 18), P(20, 1, 21)), "type `Std::True` cannot be assigned to type `Std::False`"),
			},
		},
		"infer false": {
			input: "var foo = false",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(14, 1, 15)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(14, 1, 15)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(14, 1, 15)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewFalseLiteralNode(
								S(P(10, 1, 11), P(14, 1, 15)),
							),
							globalEnv.StdSubtype(symbol.False),
						),
					),
				},
			),
		},
		"assign false literal to False": {
			input: "var foo: False = false",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(13, 1, 14)),
								"False",
								globalEnv.StdSubtype(symbol.False),
							),
							ast.NewFalseLiteralNode(
								S(P(17, 1, 18), P(21, 1, 22)),
							),
							globalEnv.StdSubtype(symbol.False),
						),
					),
				},
			),
		},
		"assign false literal to Bool": {
			input: "var foo: Bool = false",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewPublicConstantNode(
								S(P(9, 1, 10), P(12, 1, 13)),
								"Bool",
								globalEnv.StdSubtype(symbol.Bool),
							),
							ast.NewFalseLiteralNode(
								S(P(16, 1, 17), P(20, 1, 21)),
							),
							globalEnv.StdSubtype(symbol.Bool),
						),
					),
				},
			),
		},
		"assign false literal to matching literal type": {
			input: "var foo: false = false",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(21, 1, 22)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(21, 1, 22)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(21, 1, 22)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewFalseLiteralNode(
								S(P(9, 1, 10), P(13, 1, 14)),
							),
							ast.NewFalseLiteralNode(
								S(P(17, 1, 18), P(21, 1, 22)),
							),
							globalEnv.StdSubtype(symbol.False),
						),
					),
				},
			),
		},
		"assign false literal to non matching literal type": {
			input: "var foo: true = false",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(20, 1, 21)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(20, 1, 21)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(20, 1, 21)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewTrueLiteralNode(
								S(P(9, 1, 10), P(12, 1, 13)),
							),
							ast.NewFalseLiteralNode(
								S(P(16, 1, 17), P(20, 1, 21)),
							),
							globalEnv.StdSubtype(symbol.True),
						),
					),
				},
			),
			err: errors.ErrorList{
				errors.NewError(L("<main>", P(16, 1, 17), P(20, 1, 21)), "type `Std::False` cannot be assigned to type `Std::True`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			checkerTest(tc, t, false)
		})
	}
}

func TestNilLiteral(t *testing.T) {
	globalEnv := types.NewGlobalEnvironment()

	tests := testTable{
		"infer nil": {
			input: "var foo = nil",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(12, 1, 13)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(12, 1, 13)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(12, 1, 13)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							nil,
							ast.NewNilLiteralNode(
								S(P(10, 1, 11), P(12, 1, 13)),
							),
							globalEnv.StdSubtype(symbol.Nil),
						),
					),
				},
			),
		},
		"assign nil literal to Nil": {
			input: "var foo: Nil = nil",
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
								"Nil",
								globalEnv.StdSubtype(symbol.Nil),
							),
							ast.NewNilLiteralNode(
								S(P(15, 1, 16), P(17, 1, 18)),
							),
							globalEnv.StdSubtype(symbol.Nil),
						),
					),
				},
			),
		},
		"assign nil literal to matching literal type": {
			input: "var foo: nil = nil",
			want: ast.NewProgramNode(
				S(P(0, 1, 1), P(17, 1, 18)),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						S(P(0, 1, 1), P(17, 1, 18)),
						ast.NewVariableDeclarationNode(
							S(P(0, 1, 1), P(17, 1, 18)),
							V(S(P(4, 1, 5), P(6, 1, 7)), token.PUBLIC_IDENTIFIER, "foo"),
							ast.NewNilLiteralNode(
								S(P(9, 1, 10), P(11, 1, 12)),
							),
							ast.NewNilLiteralNode(
								S(P(15, 1, 16), P(17, 1, 18)),
							),
							globalEnv.StdSubtype(symbol.Nil),
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
