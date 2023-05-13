package parser

import (
	"testing"

	"github.com/elk-language/elk/lexer"
	"github.com/elk-language/elk/parser/ast"
	"github.com/google/go-cmp/cmp"
)

// Represents a single parser test case.
type testCase struct {
	input string
	want  *ast.ProgramNode
	err   ErrorList
}

// Type of the parser test table.
type testTable map[string]testCase

// Create a new token in tests.
var T = lexer.NewToken

// Create a new token with value in tests.
var V = lexer.NewTokenWithValue

// Create a new source position in tests.
var P = lexer.NewPosition

// Function which powers all parser tests.
// Inspects if the produced AST matches the expected one.
func parserTest(tc testCase, t *testing.T) {
	ast, err := Parse([]byte(tc.input))

	if diff := cmp.Diff(tc.want, ast); diff != "" {
		t.Fatal(diff)
	}

	if diff := cmp.Diff(tc.err, err); diff != "" {
		t.Fatal(diff)
	}
}

func TestAddition(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 + 2 + 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), lexer.PlusToken),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), lexer.PlusToken),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), lexer.DecIntToken, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), lexer.DecIntToken, "3")),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 +\n2 +\n3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 2, 3), lexer.PlusToken),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), lexer.PlusToken),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
								ast.NewIntLiteralNode(P(4, 1, 2, 1), V(P(4, 1, 2, 1), lexer.DecIntToken, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 3, 1), V(P(8, 1, 3, 1), lexer.DecIntToken, "3")),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n+ 2\n+ 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 2, 1, 1),
						ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
					),
					ast.NewExpressionStatementNode(
						P(2, 4, 2, 1),
						ast.NewUnaryExpressionNode(
							P(2, 3, 2, 1),
							T(P(2, 1, 2, 1), lexer.PlusToken),
							ast.NewIntLiteralNode(P(4, 1, 2, 3), V(P(4, 1, 2, 3), lexer.DecIntToken, "2")),
						),
					),
					ast.NewExpressionStatementNode(
						P(6, 3, 3, 1),
						ast.NewUnaryExpressionNode(
							P(6, 3, 3, 1),
							T(P(6, 1, 3, 1), lexer.PlusToken),
							ast.NewIntLiteralNode(P(8, 1, 3, 3), V(P(8, 1, 3, 3), lexer.DecIntToken, "3")),
						),
					),
				},
			),
		},
		"has higher precedence than comparison operators": {
			input: "foo >= bar + baz",
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 16, 1, 1),
							T(P(4, 2, 1, 5), lexer.GreaterEqualToken),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewBinaryExpressionNode(
								P(7, 9, 1, 8),
								T(P(11, 1, 1, 12), lexer.PlusToken),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								ast.NewPublicIdentifierNode(P(13, 3, 1, 14), "baz"),
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

func TestSubtraction(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 - 2 - 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), lexer.MinusToken),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), lexer.MinusToken),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), lexer.DecIntToken, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), lexer.DecIntToken, "3")),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 -\n2 -\n3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 2, 3), lexer.MinusToken),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), lexer.MinusToken),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
								ast.NewIntLiteralNode(P(4, 1, 2, 1), V(P(4, 1, 2, 1), lexer.DecIntToken, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 3, 1), V(P(8, 1, 3, 1), lexer.DecIntToken, "3")),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n- 2\n- 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 2, 1, 1),
						ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
					),
					ast.NewExpressionStatementNode(
						P(2, 4, 2, 1),
						ast.NewUnaryExpressionNode(
							P(2, 3, 2, 1),
							T(P(2, 1, 2, 1), lexer.MinusToken),
							ast.NewIntLiteralNode(P(4, 1, 2, 3), V(P(4, 1, 2, 3), lexer.DecIntToken, "2")),
						),
					),
					ast.NewExpressionStatementNode(
						P(6, 3, 3, 1),
						ast.NewUnaryExpressionNode(
							P(6, 3, 3, 1),
							T(P(6, 1, 3, 1), lexer.MinusToken),
							ast.NewIntLiteralNode(P(8, 1, 3, 3), V(P(8, 1, 3, 3), lexer.DecIntToken, "3")),
						),
					),
				},
			),
		},
		"has the same precedence as addition": {
			input: "1 + 2 - 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), lexer.MinusToken),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), lexer.PlusToken),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), lexer.DecIntToken, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), lexer.DecIntToken, "3")),
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

func TestMultiplication(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 * 2 * 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), lexer.StarToken),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), lexer.StarToken),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), lexer.DecIntToken, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), lexer.DecIntToken, "3")),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 *\n2 *\n3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 2, 3), lexer.StarToken),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), lexer.StarToken),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
								ast.NewIntLiteralNode(P(4, 1, 2, 1), V(P(4, 1, 2, 1), lexer.DecIntToken, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 3, 1), V(P(8, 1, 3, 1), lexer.DecIntToken, "3")),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n* 2\n* 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 2, 1, 1),
						ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
					),
					ast.NewExpressionStatementNode(
						P(2, 4, 2, 1),
						ast.NewInvalidNode(P(2, 1, 2, 1), T(P(2, 1, 2, 1), lexer.StarToken)),
					),
					ast.NewExpressionStatementNode(
						P(6, 3, 3, 1),
						ast.NewInvalidNode(P(6, 1, 3, 1), T(P(6, 1, 3, 1), lexer.StarToken)),
					),
				},
			),
			err: ErrorList{
				NewError(P(2, 1, 2, 1), "unexpected *, expected an expression"),
				NewError(P(6, 1, 3, 1), "unexpected *, expected an expression"),
			},
		},
		"has higher precedence than addition": {
			input: "1 + 2 * 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(2, 1, 1, 3), lexer.PlusToken),
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
							ast.NewBinaryExpressionNode(
								P(4, 5, 1, 5),
								T(P(6, 1, 1, 7), lexer.StarToken),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), lexer.DecIntToken, "2")),
								ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), lexer.DecIntToken, "3")),
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

func TestDivision(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 / 2 / 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), lexer.SlashToken),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), lexer.SlashToken),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), lexer.DecIntToken, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), lexer.DecIntToken, "3")),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 /\n2 /\n3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 2, 3), lexer.SlashToken),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), lexer.SlashToken),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
								ast.NewIntLiteralNode(P(4, 1, 2, 1), V(P(4, 1, 2, 1), lexer.DecIntToken, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 3, 1), V(P(8, 1, 3, 1), lexer.DecIntToken, "3")),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n/ 2\n/ 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 2, 1, 1),
						ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
					),
					ast.NewExpressionStatementNode(
						P(2, 4, 2, 1),
						ast.NewInvalidNode(P(2, 1, 2, 1), T(P(2, 1, 2, 1), lexer.SlashToken)),
					),
					ast.NewExpressionStatementNode(
						P(6, 3, 3, 1),
						ast.NewInvalidNode(P(6, 1, 3, 1), T(P(6, 1, 3, 1), lexer.SlashToken)),
					),
				},
			),
			err: ErrorList{
				NewError(P(2, 1, 2, 1), "unexpected /, expected an expression"),
				NewError(P(6, 1, 3, 1), "unexpected /, expected an expression"),
			},
		},
		"has the same precedence as multiplication": {
			input: "1 * 2 / 3",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 9, 1, 1),
							T(P(6, 1, 1, 7), lexer.SlashToken),
							ast.NewBinaryExpressionNode(
								P(0, 5, 1, 1),
								T(P(2, 1, 1, 3), lexer.StarToken),
								ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
								ast.NewIntLiteralNode(P(4, 1, 1, 5), V(P(4, 1, 1, 5), lexer.DecIntToken, "2")),
							),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), lexer.DecIntToken, "3")),
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

func TestUnaryExpressions(t *testing.T) {
	tests := testTable{
		"plus can be nested": {
			input: "+++1.5",
			want: ast.NewProgramNode(
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewUnaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(0, 1, 1, 1), lexer.PlusToken),
							ast.NewUnaryExpressionNode(
								P(1, 5, 1, 2),
								T(P(1, 1, 1, 2), lexer.PlusToken),
								ast.NewUnaryExpressionNode(
									P(2, 4, 1, 3),
									T(P(2, 1, 1, 3), lexer.PlusToken),
									ast.NewFloatLiteralNode(P(3, 3, 1, 4), "1.5"),
								),
							),
						),
					),
				},
			),
		},
		"minus can be nested": {
			input: "---1.5",
			want: ast.NewProgramNode(
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewUnaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(0, 1, 1, 1), lexer.MinusToken),
							ast.NewUnaryExpressionNode(
								P(1, 5, 1, 2),
								T(P(1, 1, 1, 2), lexer.MinusToken),
								ast.NewUnaryExpressionNode(
									P(2, 4, 1, 3),
									T(P(2, 1, 1, 3), lexer.MinusToken),
									ast.NewFloatLiteralNode(P(3, 3, 1, 4), "1.5"),
								),
							),
						),
					),
				},
			),
		},
		"logical not can be nested": {
			input: "!!!1.5",
			want: ast.NewProgramNode(
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewUnaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(0, 1, 1, 1), lexer.BangToken),
							ast.NewUnaryExpressionNode(
								P(1, 5, 1, 2),
								T(P(1, 1, 1, 2), lexer.BangToken),
								ast.NewUnaryExpressionNode(
									P(2, 4, 1, 3),
									T(P(2, 1, 1, 3), lexer.BangToken),
									ast.NewFloatLiteralNode(P(3, 3, 1, 4), "1.5"),
								),
							),
						),
					),
				},
			),
		},
		"bitwise not can be nested": {
			input: "~~~1.5",
			want: ast.NewProgramNode(
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewUnaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(0, 1, 1, 1), lexer.TildeToken),
							ast.NewUnaryExpressionNode(
								P(1, 5, 1, 2),
								T(P(1, 1, 1, 2), lexer.TildeToken),
								ast.NewUnaryExpressionNode(
									P(2, 4, 1, 3),
									T(P(2, 1, 1, 3), lexer.TildeToken),
									ast.NewFloatLiteralNode(P(3, 3, 1, 4), "1.5"),
								),
							),
						),
					),
				},
			),
		},
		"all have the same precedence": {
			input: "!+~1.5",
			want: ast.NewProgramNode(
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewUnaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(0, 1, 1, 1), lexer.BangToken),
							ast.NewUnaryExpressionNode(
								P(1, 5, 1, 2),
								T(P(1, 1, 1, 2), lexer.PlusToken),
								ast.NewUnaryExpressionNode(
									P(2, 4, 1, 3),
									T(P(2, 1, 1, 3), lexer.TildeToken),
									ast.NewFloatLiteralNode(P(3, 3, 1, 4), "1.5"),
								),
							),
						),
					),
				},
			),
		},
		"have higher precedence than multiplicative expressions": {
			input: "!!1.5 * 2 + ~.5",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 15, 1, 1),
							T(P(10, 1, 1, 11), lexer.PlusToken),
							ast.NewBinaryExpressionNode(
								P(0, 9, 1, 1),
								T(P(6, 1, 1, 7), lexer.StarToken),
								ast.NewUnaryExpressionNode(
									P(0, 5, 1, 1),
									T(P(0, 1, 1, 1), lexer.BangToken),
									ast.NewUnaryExpressionNode(
										P(1, 4, 1, 2),
										T(P(1, 1, 1, 2), lexer.BangToken),
										ast.NewFloatLiteralNode(P(2, 3, 1, 3), "1.5"),
									),
								),
								ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), lexer.DecIntToken, "2")),
							),
							ast.NewUnaryExpressionNode(
								P(12, 3, 1, 13),
								T(P(12, 1, 1, 13), lexer.TildeToken),
								ast.NewFloatLiteralNode(P(13, 2, 1, 14), "0.5"),
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

func TestExponentiation(t *testing.T) {
	tests := testTable{
		"is evaluated from right to left": {
			input: "1 ** 2 ** 3",
			want: ast.NewProgramNode(
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 11, 1, 1),
							T(P(2, 2, 1, 3), lexer.StarStarToken),
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
							ast.NewBinaryExpressionNode(
								P(5, 6, 1, 6),
								T(P(7, 2, 1, 8), lexer.StarStarToken),
								ast.NewIntLiteralNode(P(5, 1, 1, 6), V(P(5, 1, 1, 6), lexer.DecIntToken, "2")),
								ast.NewIntLiteralNode(P(10, 1, 1, 11), V(P(10, 1, 1, 11), lexer.DecIntToken, "3")),
							),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 **\n2 **\n3",
			want: ast.NewProgramNode(
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 11, 1, 1),
							T(P(2, 2, 1, 3), lexer.StarStarToken),
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
							ast.NewBinaryExpressionNode(
								P(5, 6, 2, 1),
								T(P(7, 2, 2, 3), lexer.StarStarToken),
								ast.NewIntLiteralNode(P(5, 1, 2, 1), V(P(5, 1, 2, 1), lexer.DecIntToken, "2")),
								ast.NewIntLiteralNode(P(10, 1, 3, 1), V(P(10, 1, 3, 1), lexer.DecIntToken, "3")),
							),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n** 2\n** 3",
			want: ast.NewProgramNode(
				P(0, 11, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 2, 1, 1),
						ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
					),
					ast.NewExpressionStatementNode(
						P(2, 5, 2, 1),
						ast.NewInvalidNode(P(2, 2, 2, 1), T(P(2, 2, 2, 1), lexer.StarStarToken)),
					),
					ast.NewExpressionStatementNode(
						P(7, 4, 3, 1),
						ast.NewInvalidNode(P(7, 2, 3, 1), T(P(7, 2, 3, 1), lexer.StarStarToken)),
					),
				},
			),
			err: ErrorList{
				NewError(P(2, 2, 2, 1), "unexpected **, expected an expression"),
				NewError(P(7, 2, 3, 1), "unexpected **, expected an expression"),
			},
		},
		"has higher precedence than unary expressions": {
			input: "-2 ** 3",
			want: ast.NewProgramNode(
				P(0, 7, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 7, 1, 1),
						ast.NewUnaryExpressionNode(
							P(0, 7, 1, 1),
							T(P(0, 1, 1, 1), lexer.MinusToken),
							ast.NewBinaryExpressionNode(
								P(1, 6, 1, 2),
								T(P(3, 2, 1, 4), lexer.StarStarToken),
								ast.NewIntLiteralNode(P(1, 1, 1, 2), V(P(1, 1, 1, 2), lexer.DecIntToken, "2")),
								ast.NewIntLiteralNode(P(6, 1, 1, 7), V(P(6, 1, 1, 7), lexer.DecIntToken, "3")),
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

func TestStatement(t *testing.T) {
	tests := testTable{
		"semicolons can separate statements": {
			input: "1 ** 2; 5 * 8",
			want: ast.NewProgramNode(
				P(0, 13, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 7, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(2, 2, 1, 3), lexer.StarStarToken),
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
							ast.NewIntLiteralNode(P(5, 1, 1, 6), V(P(5, 1, 1, 6), lexer.DecIntToken, "2")),
						),
					),
					ast.NewExpressionStatementNode(
						P(8, 5, 1, 9),
						ast.NewBinaryExpressionNode(
							P(8, 5, 1, 9),
							T(P(10, 1, 1, 11), lexer.StarToken),
							ast.NewIntLiteralNode(P(8, 1, 1, 9), V(P(8, 1, 1, 9), lexer.DecIntToken, "5")),
							ast.NewIntLiteralNode(P(12, 1, 1, 13), V(P(12, 1, 1, 13), lexer.DecIntToken, "8")),
						),
					),
				},
			),
		},
		"endlines can separate statements": {
			input: "1 ** 2\n5 * 8",
			want: ast.NewProgramNode(
				P(0, 12, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 7, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(2, 2, 1, 3), lexer.StarStarToken),
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
							ast.NewIntLiteralNode(P(5, 1, 1, 6), V(P(5, 1, 1, 6), lexer.DecIntToken, "2")),
						),
					),
					ast.NewExpressionStatementNode(
						P(7, 5, 2, 1),
						ast.NewBinaryExpressionNode(
							P(7, 5, 2, 1),
							T(P(9, 1, 2, 3), lexer.StarToken),
							ast.NewIntLiteralNode(P(7, 1, 2, 1), V(P(7, 1, 2, 1), lexer.DecIntToken, "5")),
							ast.NewIntLiteralNode(P(11, 1, 2, 5), V(P(11, 1, 2, 5), lexer.DecIntToken, "8")),
						),
					),
				},
			),
		},
		"spaces can't separate statements": {
			input: "1 ** 2 \t 5 * 8",
			want: ast.NewProgramNode(
				P(0, 14, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 6, 1, 1),
							T(P(2, 2, 1, 3), lexer.StarStarToken),
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
							ast.NewIntLiteralNode(P(5, 1, 1, 6), V(P(5, 1, 1, 6), lexer.DecIntToken, "2")),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(9, 1, 1, 10), "unexpected DecInt, expected a statement separator `\\n`, `;` or end of file"),
			},
		},
		"can be empty with newlines": {
			input: "\n\n\n",
			want: ast.NewProgramNode(
				P(0, 3, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 3, 1, 1)),
				},
			),
		},
		"can be empty with semicolons": {
			input: ";;;",
			want: ast.NewProgramNode(
				P(0, 3, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewEmptyStatementNode(P(1, 1, 1, 2)),
					ast.NewEmptyStatementNode(P(2, 1, 1, 3)),
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

func TestAssignment(t *testing.T) {
	tests := testTable{
		"ints are not valid assignment targets": {
			input: "1 -= 2",
			want: ast.NewProgramNode(
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 6, 1, 1),
							T(P(2, 2, 1, 3), lexer.MinusEqualToken),
							ast.NewIntLiteralNode(P(0, 1, 1, 1), V(P(0, 1, 1, 1), lexer.DecIntToken, "1")),
							ast.NewIntLiteralNode(P(5, 1, 1, 6), V(P(5, 1, 1, 6), lexer.DecIntToken, "2")),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(0, 1, 1, 1), "invalid `-=` assignment target"),
			},
		},
		"strings are not valid assignment targets": {
			input: "'foo' -= 2",
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 10, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 10, 1, 1),
							T(P(6, 2, 1, 7), lexer.MinusEqualToken),
							ast.NewRawStringLiteralNode(P(0, 5, 1, 1), "foo"),
							ast.NewIntLiteralNode(P(9, 1, 1, 10), V(P(9, 1, 1, 10), lexer.DecIntToken, "2")),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(0, 5, 1, 1), "invalid `-=` assignment target"),
			},
		},
		"constants are not valid assignment targets": {
			input: "FooBa -= 2",
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 10, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 10, 1, 1),
							T(P(6, 2, 1, 7), lexer.MinusEqualToken),
							ast.NewPublicConstantNode(P(0, 5, 1, 1), "FooBa"),
							ast.NewIntLiteralNode(P(9, 1, 1, 10), V(P(9, 1, 1, 10), lexer.DecIntToken, "2")),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(0, 5, 1, 1), "constants can't be assigned, maybe you meant to declare it with `:=`"),
			},
		},
		"private constants are not valid assignment targets": {
			input: "_FooB -= 2",
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 10, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 10, 1, 1),
							T(P(6, 2, 1, 7), lexer.MinusEqualToken),
							ast.NewPrivateConstantNode(P(0, 5, 1, 1), "_FooB"),
							ast.NewIntLiteralNode(P(9, 1, 1, 10), V(P(9, 1, 1, 10), lexer.DecIntToken, "2")),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(0, 5, 1, 1), "constants can't be assigned, maybe you meant to declare it with `:=`"),
			},
		},
		"identifiers can be assigned": {
			input: "foo -= 2",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 8, 1, 1),
							T(P(4, 2, 1, 5), lexer.MinusEqualToken),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewIntLiteralNode(P(7, 1, 1, 8), V(P(7, 1, 1, 8), lexer.DecIntToken, "2")),
						),
					),
				},
			),
		},
		"private identifiers can be assigned": {
			input: "_fo -= 2",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 8, 1, 1),
							T(P(4, 2, 1, 5), lexer.MinusEqualToken),
							ast.NewPrivateIdentifierNode(P(0, 3, 1, 1), "_fo"),
							ast.NewIntLiteralNode(P(7, 1, 1, 8), V(P(7, 1, 1, 8), lexer.DecIntToken, "2")),
						),
					),
				},
			),
		},
		"can be nested": {
			input: "foo = bar = baz = 3",
			want: ast.NewProgramNode(
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 19, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 19, 1, 1),
							T(P(4, 1, 1, 5), lexer.EqualToken),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewAssignmentExpressionNode(
								P(6, 13, 1, 7),
								T(P(10, 1, 1, 11), lexer.EqualToken),
								ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
								ast.NewAssignmentExpressionNode(
									P(12, 7, 1, 13),
									T(P(16, 1, 1, 17), lexer.EqualToken),
									ast.NewPublicIdentifierNode(P(12, 3, 1, 13), "baz"),
									ast.NewIntLiteralNode(P(18, 1, 1, 19), V(P(18, 1, 1, 19), lexer.DecIntToken, "3")),
								),
							),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "foo =\nbar =\nbaz =\n3",
			want: ast.NewProgramNode(
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 19, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 19, 1, 1),
							T(P(4, 1, 1, 5), lexer.EqualToken),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewAssignmentExpressionNode(
								P(6, 13, 2, 1),
								T(P(10, 1, 2, 5), lexer.EqualToken),
								ast.NewPublicIdentifierNode(P(6, 3, 2, 1), "bar"),
								ast.NewAssignmentExpressionNode(
									P(12, 7, 3, 1),
									T(P(16, 1, 3, 5), lexer.EqualToken),
									ast.NewPublicIdentifierNode(P(12, 3, 3, 1), "baz"),
									ast.NewIntLiteralNode(P(18, 1, 4, 1), V(P(18, 1, 4, 1), lexer.DecIntToken, "3")),
								),
							),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "foo\n= bar\n= baz\n= 3",
			want: ast.NewProgramNode(
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
					),
					ast.NewExpressionStatementNode(
						P(4, 6, 2, 1),
						ast.NewInvalidNode(P(4, 1, 2, 1), T(P(4, 1, 2, 1), lexer.EqualToken)),
					),
					ast.NewExpressionStatementNode(
						P(10, 6, 3, 1),
						ast.NewInvalidNode(P(10, 1, 3, 1), T(P(10, 1, 3, 1), lexer.EqualToken)),
					),
					ast.NewExpressionStatementNode(
						P(16, 3, 4, 1),
						ast.NewInvalidNode(P(16, 1, 4, 1), T(P(16, 1, 4, 1), lexer.EqualToken)),
					),
				},
			),
			err: ErrorList{
				NewError(P(4, 1, 2, 1), "unexpected =, expected an expression"),
				NewError(P(10, 1, 3, 1), "unexpected =, expected an expression"),
				NewError(P(16, 1, 4, 1), "unexpected =, expected an expression"),
			},
		},
		"has lower precedence than other expressions": {
			input: "f = some && awesome || thing + 2 * 8 > 5 == false",
			want: ast.NewProgramNode(
				P(0, 49, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 49, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 49, 1, 1),
							T(P(2, 1, 1, 3), lexer.EqualToken),
							ast.NewPublicIdentifierNode(P(0, 1, 1, 1), "f"),
							ast.NewLogicalExpressionNode(
								P(4, 45, 1, 5),
								T(P(20, 2, 1, 21), lexer.OrOrToken),
								ast.NewLogicalExpressionNode(
									P(4, 15, 1, 5),
									T(P(9, 2, 1, 10), lexer.AndAndToken),
									ast.NewPublicIdentifierNode(P(4, 4, 1, 5), "some"),
									ast.NewPublicIdentifierNode(P(12, 7, 1, 13), "awesome"),
								),
								ast.NewBinaryExpressionNode(
									P(23, 26, 1, 24),
									T(P(41, 2, 1, 42), lexer.EqualEqualToken),
									ast.NewBinaryExpressionNode(
										P(23, 17, 1, 24),
										T(P(37, 1, 1, 38), lexer.GreaterToken),
										ast.NewBinaryExpressionNode(
											P(23, 13, 1, 24),
											T(P(29, 1, 1, 30), lexer.PlusToken),
											ast.NewPublicIdentifierNode(P(23, 5, 1, 24), "thing"),
											ast.NewBinaryExpressionNode(
												P(31, 5, 1, 32),
												T(P(33, 1, 1, 34), lexer.StarToken),
												ast.NewIntLiteralNode(P(31, 1, 1, 32), V(P(31, 1, 1, 32), lexer.DecIntToken, "2")),
												ast.NewIntLiteralNode(P(35, 1, 1, 36), V(P(35, 1, 1, 36), lexer.DecIntToken, "8")),
											),
										),
										ast.NewIntLiteralNode(P(39, 1, 1, 40), V(P(39, 1, 1, 40), lexer.DecIntToken, "5")),
									),
									ast.NewFalseLiteralNode(P(44, 5, 1, 45)),
								),
							),
						),
					),
				},
			),
		},
		"has many versions": {
			input: "a = b -= c += d *= e /= f **= g ~= h &&= i &= j ||= k |= l ^= m ??= n <<= o >>= p %= q",
			want: ast.NewProgramNode(
				P(0, 86, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 86, 1, 1),
						ast.NewAssignmentExpressionNode(
							P(0, 86, 1, 1),
							T(P(2, 1, 1, 3), lexer.EqualToken),
							ast.NewPublicIdentifierNode(P(0, 1, 1, 1), "a"),
							ast.NewAssignmentExpressionNode(
								P(4, 82, 1, 5),
								T(P(6, 2, 1, 7), lexer.MinusEqualToken),
								ast.NewPublicIdentifierNode(P(4, 1, 1, 5), "b"),
								ast.NewAssignmentExpressionNode(
									P(9, 77, 1, 10),
									T(P(11, 2, 1, 12), lexer.PlusEqualToken),
									ast.NewPublicIdentifierNode(P(9, 1, 1, 10), "c"),
									ast.NewAssignmentExpressionNode(
										P(14, 72, 1, 15),
										T(P(16, 2, 1, 17), lexer.StarEqualToken),
										ast.NewPublicIdentifierNode(P(14, 1, 1, 15), "d"),
										ast.NewAssignmentExpressionNode(
											P(19, 67, 1, 20),
											T(P(21, 2, 1, 22), lexer.SlashEqualToken),
											ast.NewPublicIdentifierNode(P(19, 1, 1, 20), "e"),
											ast.NewAssignmentExpressionNode(
												P(24, 62, 1, 25),
												T(P(26, 3, 1, 27), lexer.StarStarEqualToken),
												ast.NewPublicIdentifierNode(P(24, 1, 1, 25), "f"),
												ast.NewAssignmentExpressionNode(
													P(30, 56, 1, 31),
													T(P(32, 2, 1, 33), lexer.TildeEqualToken),
													ast.NewPublicIdentifierNode(P(30, 1, 1, 31), "g"),
													ast.NewAssignmentExpressionNode(
														P(35, 51, 1, 36),
														T(P(37, 3, 1, 38), lexer.AndAndEqualToken),
														ast.NewPublicIdentifierNode(P(35, 1, 1, 36), "h"),
														ast.NewAssignmentExpressionNode(
															P(41, 45, 1, 42),
															T(P(43, 2, 1, 44), lexer.AndEqualToken),
															ast.NewPublicIdentifierNode(P(41, 1, 1, 42), "i"),
															ast.NewAssignmentExpressionNode(
																P(46, 40, 1, 47),
																T(P(48, 3, 1, 49), lexer.OrOrEqualToken),
																ast.NewPublicIdentifierNode(P(46, 1, 1, 47), "j"),
																ast.NewAssignmentExpressionNode(
																	P(52, 34, 1, 53),
																	T(P(54, 2, 1, 55), lexer.OrEqualToken),
																	ast.NewPublicIdentifierNode(P(52, 1, 1, 53), "k"),
																	ast.NewAssignmentExpressionNode(
																		P(57, 29, 1, 58),
																		T(P(59, 2, 1, 60), lexer.XorEqualToken),
																		ast.NewPublicIdentifierNode(P(57, 1, 1, 58), "l"),
																		ast.NewAssignmentExpressionNode(
																			P(62, 24, 1, 63),
																			T(P(64, 3, 1, 65), lexer.QuestionQuestionEqualToken),
																			ast.NewPublicIdentifierNode(P(62, 1, 1, 63), "m"),
																			ast.NewAssignmentExpressionNode(
																				P(68, 18, 1, 69),
																				T(P(70, 3, 1, 71), lexer.LBitShiftEqualToken),
																				ast.NewPublicIdentifierNode(P(68, 1, 1, 69), "n"),
																				ast.NewAssignmentExpressionNode(
																					P(74, 12, 1, 75),
																					T(P(76, 3, 1, 77), lexer.RBitShiftEqualToken),
																					ast.NewPublicIdentifierNode(P(74, 1, 1, 75), "o"),
																					ast.NewAssignmentExpressionNode(
																						P(80, 6, 1, 81),
																						T(P(82, 2, 1, 83), lexer.PercentEqualToken),
																						ast.NewPublicIdentifierNode(P(80, 1, 1, 81), "p"),
																						ast.NewPublicIdentifierNode(P(85, 1, 1, 86), "q"),
																					),
																				),
																			),
																		),
																	),
																),
															),
														),
													),
												),
											),
										),
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

func TestBooleanLogic(t *testing.T) {
	tests := testTable{
		"has lower precedence than equality": {
			input: "foo && bar == baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(4, 2, 1, 5), lexer.AndAndToken),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewBinaryExpressionNode(
								P(7, 10, 1, 8),
								T(P(11, 2, 1, 12), lexer.EqualEqualToken),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
							),
						),
					),
				},
			),
		},
		"or has lower precedence than and": {
			input: "foo || bar && baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(4, 2, 1, 5), lexer.OrOrToken),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewLogicalExpressionNode(
								P(7, 10, 1, 8),
								T(P(11, 2, 1, 12), lexer.AndAndToken),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
							),
						),
					),
				},
			),
		},
		"nil coalescing operator has lower precedence than and": {
			input: "foo ?? bar && baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(4, 2, 1, 5), lexer.QuestionQuestionToken),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewLogicalExpressionNode(
								P(7, 10, 1, 8),
								T(P(11, 2, 1, 12), lexer.AndAndToken),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
							),
						),
					),
				},
			),
		},
		"nil coalescing operator has the same precedence as or": {
			input: "foo ?? bar || baz ?? boo",
			want: ast.NewProgramNode(
				P(0, 24, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 24, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 24, 1, 1),
							T(P(18, 2, 1, 19), lexer.QuestionQuestionToken),
							ast.NewLogicalExpressionNode(
								P(0, 17, 1, 1),
								T(P(11, 2, 1, 12), lexer.OrOrToken),
								ast.NewLogicalExpressionNode(
									P(0, 10, 1, 1),
									T(P(4, 2, 1, 5), lexer.QuestionQuestionToken),
									ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
									ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								),
								ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
							),
							ast.NewPublicIdentifierNode(P(21, 3, 1, 22), "boo"),
						),
					),
				},
			),
		},
		"or is evaluated from left to right": {
			input: "foo || bar || baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(11, 2, 1, 12), lexer.OrOrToken),
							ast.NewLogicalExpressionNode(
								P(0, 10, 1, 1),
								T(P(4, 2, 1, 5), lexer.OrOrToken),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
							),
							ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
						),
					),
				},
			),
		},
		"and is evaluated from left to right": {
			input: "foo && bar && baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(11, 2, 1, 12), lexer.AndAndToken),
							ast.NewLogicalExpressionNode(
								P(0, 10, 1, 1),
								T(P(4, 2, 1, 5), lexer.AndAndToken),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
							),
							ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
						),
					),
				},
			),
		},
		"nil coalescing operator is evaluated from left to right": {
			input: "foo ?? bar ?? baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(11, 2, 1, 12), lexer.QuestionQuestionToken),
							ast.NewLogicalExpressionNode(
								P(0, 10, 1, 1),
								T(P(4, 2, 1, 5), lexer.QuestionQuestionToken),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
							),
							ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
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

func TestStringLiteral(t *testing.T) {
	tests := testTable{
		"processes escape sequences": {
			input: `"foo\nbar\rbaz\\car\t\b\"\v\f\x12\a"`,
			want: ast.NewProgramNode(
				P(0, 36, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 36, 1, 1),
						ast.NewStringLiteralNode(
							P(0, 36, 1, 1),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(P(1, 34, 1, 2), "foo\nbar\rbaz\\car\t\b\"\v\f\x12\a"),
							},
						),
					),
				},
			),
		},
		"reports errors for invalid hex escapes": {
			input: `"foo \xgh bar"`,
			want: ast.NewProgramNode(
				P(0, 14, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 14, 1, 1),
						ast.NewStringLiteralNode(
							P(0, 14, 1, 1),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(P(1, 4, 1, 2), "foo "),
								ast.NewInvalidNode(P(5, 4, 1, 6), V(P(5, 4, 1, 6), lexer.ErrorToken, "invalid hex escape in string literal")),
								ast.NewStringLiteralContentSectionNode(P(9, 4, 1, 10), " bar"),
							},
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(5, 4, 1, 6), "invalid hex escape in string literal"),
			},
		},
		"reports errors for nonexistent escape sequences": {
			input: `"foo \q bar"`,
			want: ast.NewProgramNode(
				P(0, 12, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 12, 1, 1),
						ast.NewStringLiteralNode(
							P(0, 12, 1, 1),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(P(1, 4, 1, 2), "foo "),
								ast.NewInvalidNode(P(5, 2, 1, 6), V(P(5, 2, 1, 6), lexer.ErrorToken, "invalid escape sequence `\\q` in string literal")),
								ast.NewStringLiteralContentSectionNode(P(7, 4, 1, 8), " bar"),
							},
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(5, 2, 1, 6), "invalid escape sequence `\\q` in string literal"),
			},
		},
		"can contain interpolated expressions": {
			input: `"foo ${bar + 2} baz ${fudge}"`,
			want: ast.NewProgramNode(
				P(0, 29, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 29, 1, 1),
						ast.NewStringLiteralNode(
							P(0, 29, 1, 1),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(P(1, 4, 1, 2), "foo "),
								ast.NewStringInterpolationNode(
									P(5, 10, 1, 6),
									ast.NewBinaryExpressionNode(
										P(7, 7, 1, 8),
										T(P(11, 1, 1, 12), lexer.PlusToken),
										ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
										ast.NewIntLiteralNode(P(13, 1, 1, 14), V(P(13, 1, 1, 14), lexer.DecIntToken, "2")),
									),
								),
								ast.NewStringLiteralContentSectionNode(P(15, 5, 1, 16), " baz "),
								ast.NewStringInterpolationNode(
									P(20, 8, 1, 21),
									ast.NewPublicIdentifierNode(P(22, 5, 1, 23), "fudge"),
								),
							},
						),
					),
				},
			),
		},
		"can't contain string literals inside interpolation": {
			input: `"foo ${"bar" + 2} baza"`,
			want: ast.NewProgramNode(
				P(0, 23, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 23, 1, 1),
						ast.NewStringLiteralNode(
							P(0, 23, 1, 1),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(P(1, 4, 1, 2), "foo "),
								ast.NewStringInterpolationNode(
									P(5, 12, 1, 6),
									ast.NewBinaryExpressionNode(
										P(7, 9, 1, 8),
										T(P(13, 1, 1, 14), lexer.PlusToken),
										ast.NewInvalidNode(P(7, 5, 1, 8), V(P(7, 5, 1, 8), lexer.ErrorToken, "unexpected string literal in string interpolation, only raw strings delimited with `'` can be used in string interpolation")),
										ast.NewIntLiteralNode(P(15, 1, 1, 16), V(P(15, 1, 1, 16), lexer.DecIntToken, "2")),
									),
								),
								ast.NewStringLiteralContentSectionNode(P(17, 5, 1, 18), " baza"),
							},
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(7, 5, 1, 8), "unexpected string literal in string interpolation, only raw strings delimited with `'` can be used in string interpolation"),
			},
		},
		"can contain raw string literals inside interpolation": {
			input: `"foo ${'bar' + 2} baza"`,
			want: ast.NewProgramNode(
				P(0, 23, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 23, 1, 1),
						ast.NewStringLiteralNode(
							P(0, 23, 1, 1),
							[]ast.StringLiteralContentNode{
								ast.NewStringLiteralContentSectionNode(P(1, 4, 1, 2), "foo "),
								ast.NewStringInterpolationNode(
									P(5, 12, 1, 6),
									ast.NewBinaryExpressionNode(
										P(7, 9, 1, 8),
										T(P(13, 1, 1, 14), lexer.PlusToken),
										ast.NewRawStringLiteralNode(P(7, 5, 1, 8), "bar"),
										ast.NewIntLiteralNode(P(15, 1, 1, 16), V(P(15, 1, 1, 16), lexer.DecIntToken, "2")),
									),
								),
								ast.NewStringLiteralContentSectionNode(P(17, 5, 1, 18), " baza"),
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

func TestRawStringLiteral(t *testing.T) {
	tests := testTable{
		"doesn't process escape sequences": {
			input: `'foo\nbar\rbaz\\car\t\b\"\v\f\x12\a'`,
			want: ast.NewProgramNode(
				P(0, 36, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 36, 1, 1),
						ast.NewRawStringLiteralNode(P(0, 36, 1, 1), `foo\nbar\rbaz\\car\t\b\"\v\f\x12\a`),
					),
				},
			),
		},
		"can't contain interpolated expressions": {
			input: `'foo ${bar + 2} baz ${fudge}'`,
			want: ast.NewProgramNode(
				P(0, 29, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 29, 1, 1),
						ast.NewRawStringLiteralNode(P(0, 29, 1, 1), `foo ${bar + 2} baz ${fudge}`),
					),
				},
			),
		},
		"can contain double quotes": {
			input: `'foo ${"bar" + 2} baza'`,
			want: ast.NewProgramNode(
				P(0, 23, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 23, 1, 1),
						ast.NewRawStringLiteralNode(P(0, 23, 1, 1), `foo ${"bar" + 2} baza`),
					),
				},
			),
		},
		"doesn't allow escaping single quotes": {
			input: `'foo\'s house'`,
			want: ast.NewProgramNode(
				P(0, 14, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewRawStringLiteralNode(P(0, 6, 1, 1), "foo\\"),
					),
				},
			),
			err: ErrorList{
				NewError(P(6, 1, 1, 7), "unexpected PublicIdentifier, expected a statement separator `\\n`, `;` or end of file"),
				NewError(P(13, 1, 1, 14), "unterminated raw string literal, missing `'`"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestEquality(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "bar == baz == 1",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 15, 1, 1),
							T(P(11, 2, 1, 12), lexer.EqualEqualToken),
							ast.NewBinaryExpressionNode(
								P(0, 10, 1, 1),
								T(P(4, 2, 1, 5), lexer.EqualEqualToken),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "bar"),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "baz"),
							),
							ast.NewIntLiteralNode(P(14, 1, 1, 15), V(P(14, 1, 1, 15), lexer.DecIntToken, "1")),
						),
					),
				},
			),
		},
		"can have endlines after the operator": {
			input: "bar ==\nbaz ==\n1",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 15, 1, 1),
							T(P(11, 2, 2, 5), lexer.EqualEqualToken),
							ast.NewBinaryExpressionNode(
								P(0, 10, 1, 1),
								T(P(4, 2, 1, 5), lexer.EqualEqualToken),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "bar"),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 1), "baz"),
							),
							ast.NewIntLiteralNode(P(14, 1, 3, 1), V(P(14, 1, 3, 1), lexer.DecIntToken, "1")),
						),
					),
				},
			),
		},
		"can't have endlines before the operator": {
			input: "bar\n== baz\n== 1",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "bar"),
					),
					ast.NewExpressionStatementNode(
						P(4, 7, 2, 1),
						ast.NewInvalidNode(P(4, 2, 2, 1), T(P(4, 2, 2, 1), lexer.EqualEqualToken)),
					),
					ast.NewExpressionStatementNode(
						P(11, 4, 3, 1),
						ast.NewInvalidNode(P(11, 2, 3, 1), T(P(11, 2, 3, 1), lexer.EqualEqualToken)),
					),
				},
			),
			err: ErrorList{
				NewError(P(4, 2, 2, 1), "unexpected ==, expected an expression"),
				NewError(P(11, 2, 3, 1), "unexpected ==, expected an expression"),
			},
		},
		"has many versions": {
			input: "a == b != c === d !== e =:= f =!= g",
			want: ast.NewProgramNode(
				P(0, 35, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 35, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 35, 1, 1),
							T(P(30, 3, 1, 31), lexer.RefNotEqualToken),
							ast.NewBinaryExpressionNode(
								P(0, 29, 1, 1),
								T(P(24, 3, 1, 25), lexer.RefEqualToken),
								ast.NewBinaryExpressionNode(
									P(0, 23, 1, 1),
									T(P(18, 3, 1, 19), lexer.StrictNotEqualToken),
									ast.NewBinaryExpressionNode(
										P(0, 17, 1, 1),
										T(P(12, 3, 1, 13), lexer.StrictEqualToken),
										ast.NewBinaryExpressionNode(
											P(0, 11, 1, 1),
											T(P(7, 2, 1, 8), lexer.NotEqualToken),
											ast.NewBinaryExpressionNode(
												P(0, 6, 1, 1),
												T(P(2, 2, 1, 3), lexer.EqualEqualToken),
												ast.NewPublicIdentifierNode(P(0, 1, 1, 1), "a"),
												ast.NewPublicIdentifierNode(P(5, 1, 1, 6), "b"),
											),
											ast.NewPublicIdentifierNode(P(10, 1, 1, 11), "c"),
										),
										ast.NewPublicIdentifierNode(P(16, 1, 1, 17), "d"),
									),
									ast.NewPublicIdentifierNode(P(22, 1, 1, 23), "e"),
								),
								ast.NewPublicIdentifierNode(P(28, 1, 1, 29), "f"),
							),
							ast.NewPublicIdentifierNode(P(34, 1, 1, 35), "g"),
						),
					),
				},
			),
		},
		"has higher precedence than logical operators": {
			input: "foo && bar == baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 17, 1, 1),
							T(P(4, 2, 1, 5), lexer.AndAndToken),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewBinaryExpressionNode(
								P(7, 10, 1, 8),
								T(P(11, 2, 1, 12), lexer.EqualEqualToken),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
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

func TestComparison(t *testing.T) {
	tests := testTable{
		"is processed from left to right": {
			input: "foo > bar > baz",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 15, 1, 1),
							T(P(10, 1, 1, 11), lexer.GreaterToken),
							ast.NewBinaryExpressionNode(
								P(0, 9, 1, 1),
								T(P(4, 1, 1, 5), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
							),
							ast.NewPublicIdentifierNode(P(12, 3, 1, 13), "baz"),
						),
					),
				},
			),
		},
		"can have endlines after the operator": {
			input: "foo >\nbar >\nbaz",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 15, 1, 1),
							T(P(10, 1, 2, 5), lexer.GreaterToken),
							ast.NewBinaryExpressionNode(
								P(0, 9, 1, 1),
								T(P(4, 1, 1, 5), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(6, 3, 2, 1), "bar"),
							),
							ast.NewPublicIdentifierNode(P(12, 3, 3, 1), "baz"),
						),
					),
				},
			),
		},
		"can't have endlines before the operator": {
			input: "bar\n> baz\n> baz",
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "bar"),
					),
					ast.NewExpressionStatementNode(
						P(4, 6, 2, 1),
						ast.NewInvalidNode(P(4, 1, 2, 1), T(P(4, 1, 2, 1), lexer.GreaterToken)),
					),
					ast.NewExpressionStatementNode(
						P(10, 5, 3, 1),
						ast.NewInvalidNode(P(10, 1, 3, 1), T(P(10, 1, 3, 1), lexer.GreaterToken)),
					),
				},
			),
			err: ErrorList{
				NewError(P(4, 1, 2, 1), "unexpected >, expected an expression"),
				NewError(P(10, 1, 3, 1), "unexpected >, expected an expression"),
			},
		},
		"has many versions": {
			input: "a < b <= c > d >= e <: f :> g <<: h :>> i <=> j",
			want: ast.NewProgramNode(
				P(0, 47, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 47, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 47, 1, 1),
							T(P(42, 3, 1, 43), lexer.SpaceshipOpToken),
							ast.NewBinaryExpressionNode(
								P(0, 41, 1, 1),
								T(P(36, 3, 1, 37), lexer.ReverseInstanceOfToken),
								ast.NewBinaryExpressionNode(
									P(0, 35, 1, 1),
									T(P(30, 3, 1, 31), lexer.InstanceOfToken),
									ast.NewBinaryExpressionNode(
										P(0, 29, 1, 1),
										T(P(25, 2, 1, 26), lexer.ReverseSubtypeToken),
										ast.NewBinaryExpressionNode(
											P(0, 24, 1, 1),
											T(P(20, 2, 1, 21), lexer.SubtypeToken),
											ast.NewBinaryExpressionNode(
												P(0, 19, 1, 1),
												T(P(15, 2, 1, 16), lexer.GreaterEqualToken),
												ast.NewBinaryExpressionNode(
													P(0, 14, 1, 1),
													T(P(11, 1, 1, 12), lexer.GreaterToken),
													ast.NewBinaryExpressionNode(
														P(0, 10, 1, 1),
														T(P(6, 2, 1, 7), lexer.LessEqualToken),
														ast.NewBinaryExpressionNode(
															P(0, 5, 1, 1),
															T(P(2, 1, 1, 3), lexer.LessToken),
															ast.NewPublicIdentifierNode(P(0, 1, 1, 1), "a"),
															ast.NewPublicIdentifierNode(P(4, 1, 1, 5), "b"),
														),
														ast.NewPublicIdentifierNode(P(9, 1, 1, 10), "c"),
													),
													ast.NewPublicIdentifierNode(P(13, 1, 1, 14), "d"),
												),
												ast.NewPublicIdentifierNode(P(18, 1, 1, 19), "e"),
											),
											ast.NewPublicIdentifierNode(P(23, 1, 1, 24), "f"),
										),
										ast.NewPublicIdentifierNode(P(28, 1, 1, 29), "g"),
									),
									ast.NewPublicIdentifierNode(P(34, 1, 1, 35), "h"),
								),
								ast.NewPublicIdentifierNode(P(40, 1, 1, 41), "i"),
							),
							ast.NewPublicIdentifierNode(P(46, 1, 1, 47), "j"),
						),
					),
				},
			),
		},
		"has higher precedence than equality operators": {
			input: "foo == bar >= baz",
			want: ast.NewProgramNode(
				P(0, 17, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 17, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 17, 1, 1),
							T(P(4, 2, 1, 5), lexer.EqualEqualToken),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewBinaryExpressionNode(
								P(7, 10, 1, 8),
								T(P(11, 2, 1, 12), lexer.GreaterEqualToken),
								ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
								ast.NewPublicIdentifierNode(P(14, 3, 1, 15), "baz"),
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

func TestModifierExpression(t *testing.T) {
	tests := testTable{
		"has lower precedence than assignment": {
			input: "foo = bar if baz",
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewModifierNode(
							P(0, 16, 1, 1),
							T(P(10, 2, 1, 11), lexer.IfToken),
							ast.NewAssignmentExpressionNode(
								P(0, 9, 1, 1),
								T(P(4, 1, 1, 5), lexer.EqualToken),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
							),
							ast.NewPublicIdentifierNode(P(13, 3, 1, 14), "baz"),
						),
					),
				},
			),
		},
		"if can contain else": {
			input: "foo = bar if baz else car = red",
			want: ast.NewProgramNode(
				P(0, 31, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 31, 1, 1),
						ast.NewModifierIfElseNode(
							P(0, 31, 1, 1),
							ast.NewAssignmentExpressionNode(
								P(0, 9, 1, 1),
								T(P(4, 1, 1, 5), lexer.EqualToken),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
							),
							ast.NewPublicIdentifierNode(P(13, 3, 1, 14), "baz"),
							ast.NewAssignmentExpressionNode(
								P(22, 9, 1, 23),
								T(P(26, 1, 1, 27), lexer.EqualToken),
								ast.NewPublicIdentifierNode(P(22, 3, 1, 23), "car"),
								ast.NewPublicIdentifierNode(P(28, 3, 1, 29), "red"),
							),
						),
					),
				},
			),
		},
		"has many versions": {
			input: "foo if bar\nfoo unless bar\nfoo while bar\nfoo until bar",
			want: ast.NewProgramNode(
				P(0, 53, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 11, 1, 1),
						ast.NewModifierNode(
							P(0, 10, 1, 1),
							T(P(4, 2, 1, 5), lexer.IfToken),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewPublicIdentifierNode(P(7, 3, 1, 8), "bar"),
						),
					),
					ast.NewExpressionStatementNode(
						P(11, 15, 2, 1),
						ast.NewModifierNode(
							P(11, 14, 2, 1),
							T(P(15, 6, 2, 5), lexer.UnlessToken),
							ast.NewPublicIdentifierNode(P(11, 3, 2, 1), "foo"),
							ast.NewPublicIdentifierNode(P(22, 3, 2, 12), "bar"),
						),
					),
					ast.NewExpressionStatementNode(
						P(26, 14, 3, 1),
						ast.NewModifierNode(
							P(26, 13, 3, 1),
							T(P(30, 5, 3, 5), lexer.WhileToken),
							ast.NewPublicIdentifierNode(P(26, 3, 3, 1), "foo"),
							ast.NewPublicIdentifierNode(P(36, 3, 3, 11), "bar"),
						),
					),
					ast.NewExpressionStatementNode(
						P(40, 13, 4, 1),
						ast.NewModifierNode(
							P(40, 13, 4, 1),
							T(P(44, 5, 4, 5), lexer.UntilToken),
							ast.NewPublicIdentifierNode(P(40, 3, 4, 1), "foo"),
							ast.NewPublicIdentifierNode(P(50, 3, 4, 11), "bar"),
						),
					),
				},
			),
		},
		"can't be nested": {
			input: "foo = bar if baz if false\n3",
			want: ast.NewProgramNode(
				P(0, 27, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 16, 1, 1),
						ast.NewModifierNode(
							P(0, 16, 1, 1),
							T(P(10, 2, 1, 11), lexer.IfToken),
							ast.NewAssignmentExpressionNode(
								P(0, 9, 1, 1),
								T(P(4, 1, 1, 5), lexer.EqualToken),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicIdentifierNode(P(6, 3, 1, 7), "bar"),
							),
							ast.NewPublicIdentifierNode(P(13, 3, 1, 14), "baz"),
						),
					),
					ast.NewExpressionStatementNode(
						P(26, 1, 2, 1),
						ast.NewIntLiteralNode(P(26, 1, 2, 1), V(P(26, 1, 2, 1), lexer.DecIntToken, "3")),
					),
				},
			),
			err: ErrorList{
				NewError(P(17, 2, 1, 18), "unexpected if, expected a statement separator `\\n`, `;` or end of file"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestIf(t *testing.T) {
	tests := testTable{
		"can have one branch": {
			input: `
if foo > 0
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				P(0, 31, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 30, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 29, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(13, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(13, 8, 3, 2),
										T(P(17, 2, 3, 6), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(13, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(20, 1, 3, 9), V(P(20, 1, 3, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(23, 4, 4, 2),
									ast.NewNilLiteralNode(P(23, 3, 4, 2)),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
if foo > 0
end
`,
			want: ast.NewProgramNode(
				P(0, 16, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 15, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 14, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), lexer.DecIntToken, "0")),
							),
							nil,
							nil,
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
bar =
	if foo > 0
		foo += 2
	end
nil
`,
			want: ast.NewProgramNode(
				P(0, 39, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 34, 2, 1),
						ast.NewAssignmentExpressionNode(
							P(1, 33, 2, 1),
							T(P(5, 1, 2, 5), lexer.EqualToken),
							ast.NewPublicIdentifierNode(P(1, 3, 2, 1), "bar"),
							ast.NewIfExpressionNode(
								P(8, 26, 3, 2),
								ast.NewBinaryExpressionNode(
									P(11, 7, 3, 5),
									T(P(15, 1, 3, 9), lexer.GreaterToken),
									ast.NewPublicIdentifierNode(P(11, 3, 3, 5), "foo"),
									ast.NewIntLiteralNode(P(17, 1, 3, 11), V(P(17, 1, 3, 11), lexer.DecIntToken, "0")),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										P(21, 9, 4, 3),
										ast.NewAssignmentExpressionNode(
											P(21, 8, 4, 3),
											T(P(25, 2, 4, 7), lexer.PlusEqualToken),
											ast.NewPublicIdentifierNode(P(21, 3, 4, 3), "foo"),
											ast.NewIntLiteralNode(P(28, 1, 4, 10), V(P(28, 1, 4, 10), lexer.DecIntToken, "2")),
										),
									),
								},
								nil,
							),
						),
					),
					ast.NewExpressionStatementNode(
						P(35, 4, 6, 1),
						ast.NewNilLiteralNode(P(35, 3, 6, 1)),
					),
				},
			),
		},
		"can be single line with then and without end": {
			input: `
if foo > 0 then foo += 2
nil
`,
			want: ast.NewProgramNode(
				P(0, 30, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 25, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 24, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(17, 8, 2, 17),
									ast.NewAssignmentExpressionNode(
										P(17, 8, 2, 17),
										T(P(21, 2, 2, 21), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(17, 3, 2, 17), "foo"),
										ast.NewIntLiteralNode(P(24, 1, 2, 24), V(P(24, 1, 2, 24), lexer.DecIntToken, "2")),
									),
								),
							},
							nil,
						),
					),
					ast.NewExpressionStatementNode(
						P(26, 4, 3, 1),
						ast.NewNilLiteralNode(P(26, 3, 3, 1)),
					),
				},
			),
		},
		"can have else": {
			input: `
if foo > 0
	foo += 2
	nil
else
  foo -= 2
	nil
end
nil
`,
			want: ast.NewProgramNode(
				P(0, 56, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 51, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 50, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(13, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(13, 8, 3, 2),
										T(P(17, 2, 3, 6), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(13, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(20, 1, 3, 9), V(P(20, 1, 3, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(23, 4, 4, 2),
									ast.NewNilLiteralNode(P(23, 3, 4, 2)),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(34, 9, 6, 3),
									ast.NewAssignmentExpressionNode(
										P(34, 8, 6, 3),
										T(P(38, 2, 6, 7), lexer.MinusEqualToken),
										ast.NewPublicIdentifierNode(P(34, 3, 6, 3), "foo"),
										ast.NewIntLiteralNode(P(41, 1, 6, 10), V(P(41, 1, 6, 10), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(44, 4, 7, 2),
									ast.NewNilLiteralNode(P(44, 3, 7, 2)),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(52, 4, 9, 1),
						ast.NewNilLiteralNode(P(52, 3, 9, 1)),
					),
				},
			),
		},
		"can have else in short form": {
			input: `
if foo > 0 then foo += 2
else foo -= 2
nil
`,
			want: ast.NewProgramNode(
				P(0, 44, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 39, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 38, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(17, 8, 2, 17),
									ast.NewAssignmentExpressionNode(
										P(17, 8, 2, 17),
										T(P(21, 2, 2, 21), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(17, 3, 2, 17), "foo"),
										ast.NewIntLiteralNode(P(24, 1, 2, 24), V(P(24, 1, 2, 24), lexer.DecIntToken, "2")),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(31, 8, 3, 6),
									ast.NewAssignmentExpressionNode(
										P(31, 8, 3, 6),
										T(P(35, 2, 3, 10), lexer.MinusEqualToken),
										ast.NewPublicIdentifierNode(P(31, 3, 3, 6), "foo"),
										ast.NewIntLiteralNode(P(38, 1, 3, 13), V(P(38, 1, 3, 13), lexer.DecIntToken, "2")),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(40, 4, 4, 1),
						ast.NewNilLiteralNode(P(40, 3, 4, 1)),
					),
				},
			),
		},
		"can't have two elses": {
			input: `
if foo > 0 then foo += 2
else foo -= 2
else bar
nil
`,
			want: ast.NewProgramNode(
				P(0, 53, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 39, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 38, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(17, 8, 2, 17),
									ast.NewAssignmentExpressionNode(
										P(17, 8, 2, 17),
										T(P(21, 2, 2, 21), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(17, 3, 2, 17), "foo"),
										ast.NewIntLiteralNode(P(24, 1, 2, 24), V(P(24, 1, 2, 24), lexer.DecIntToken, "2")),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(31, 8, 3, 6),
									ast.NewAssignmentExpressionNode(
										P(31, 8, 3, 6),
										T(P(35, 2, 3, 10), lexer.MinusEqualToken),
										ast.NewPublicIdentifierNode(P(31, 3, 3, 6), "foo"),
										ast.NewIntLiteralNode(P(38, 1, 3, 13), V(P(38, 1, 3, 13), lexer.DecIntToken, "2")),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(40, 9, 4, 1),
						ast.NewInvalidNode(P(40, 4, 4, 1), T(P(40, 4, 4, 1), lexer.ElseToken)),
					),
					ast.NewExpressionStatementNode(
						P(49, 4, 5, 1),
						ast.NewNilLiteralNode(P(49, 3, 5, 1)),
					),
				},
			),
			err: ErrorList{
				NewError(P(40, 4, 4, 1), "unexpected else, expected an expression"),
			},
		},
		"can have many elsif blocks": {
			input: `
if foo > 0
	foo += 2
	nil
elsif foo < 5
	foo *= 10
elsif foo < 0
	foo %= 3
else
	foo -= 2
	nil
end
nil
`,
			want: ast.NewProgramNode(
				P(0, 104, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 99, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 98, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(13, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(13, 8, 3, 2),
										T(P(17, 2, 3, 6), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(13, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(20, 1, 3, 9), V(P(20, 1, 3, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(23, 4, 4, 2),
									ast.NewNilLiteralNode(P(23, 3, 4, 2)),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(27, 25, 5, 1),
									ast.NewIfExpressionNode(
										P(27, 25, 5, 1),
										ast.NewBinaryExpressionNode(
											P(33, 7, 5, 7),
											T(P(37, 1, 5, 11), lexer.LessToken),
											ast.NewPublicIdentifierNode(P(33, 3, 5, 7), "foo"),
											ast.NewIntLiteralNode(P(39, 1, 5, 13), V(P(39, 1, 5, 13), lexer.DecIntToken, "5")),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												P(42, 10, 6, 2),
												ast.NewAssignmentExpressionNode(
													P(42, 9, 6, 2),
													T(P(46, 2, 6, 6), lexer.StarEqualToken),
													ast.NewPublicIdentifierNode(P(42, 3, 6, 2), "foo"),
													ast.NewIntLiteralNode(P(49, 2, 6, 9), V(P(49, 2, 6, 9), lexer.DecIntToken, "10")),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												P(52, 47, 7, 1),
												ast.NewIfExpressionNode(
													P(52, 47, 7, 1),
													ast.NewBinaryExpressionNode(
														P(58, 7, 7, 7),
														T(P(62, 1, 7, 11), lexer.LessToken),
														ast.NewPublicIdentifierNode(P(58, 3, 7, 7), "foo"),
														ast.NewIntLiteralNode(P(64, 1, 7, 13), V(P(64, 1, 7, 13), lexer.DecIntToken, "0")),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															P(67, 9, 8, 2),
															ast.NewAssignmentExpressionNode(
																P(67, 8, 8, 2),
																T(P(71, 2, 8, 6), lexer.PercentEqualToken),
																ast.NewPublicIdentifierNode(P(67, 3, 8, 2), "foo"),
																ast.NewIntLiteralNode(P(74, 1, 8, 9), V(P(74, 1, 8, 9), lexer.DecIntToken, "3")),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															P(82, 9, 10, 2),
															ast.NewAssignmentExpressionNode(
																P(82, 8, 10, 2),
																T(P(86, 2, 10, 6), lexer.MinusEqualToken),
																ast.NewPublicIdentifierNode(P(82, 3, 10, 2), "foo"),
																ast.NewIntLiteralNode(P(89, 1, 10, 9), V(P(89, 1, 10, 9), lexer.DecIntToken, "2")),
															),
														),
														ast.NewExpressionStatementNode(
															P(92, 4, 11, 2),
															ast.NewNilLiteralNode(P(92, 3, 11, 2)),
														),
													},
												),
											),
										},
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(100, 4, 13, 1),
						ast.NewNilLiteralNode(P(100, 3, 13, 1)),
					),
				},
			),
		},
		"can have elsifs in short form": {
			input: `
if foo > 0 then foo += 2
elsif foo < 5 then foo *= 10
elsif foo < 0 then foo %= 3
else foo -= 2
nil
`,
			want: ast.NewProgramNode(
				P(0, 101, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 96, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 95, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(17, 8, 2, 17),
									ast.NewAssignmentExpressionNode(
										P(17, 8, 2, 17),
										T(P(21, 2, 2, 21), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(17, 3, 2, 17), "foo"),
										ast.NewIntLiteralNode(P(24, 1, 2, 24), V(P(24, 1, 2, 24), lexer.DecIntToken, "2")),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(26, 28, 3, 1),
									ast.NewIfExpressionNode(
										P(26, 28, 3, 1),
										ast.NewBinaryExpressionNode(
											P(32, 7, 3, 7),
											T(P(36, 1, 3, 11), lexer.LessToken),
											ast.NewPublicIdentifierNode(P(32, 3, 3, 7), "foo"),
											ast.NewIntLiteralNode(P(38, 1, 3, 13), V(P(38, 1, 3, 13), lexer.DecIntToken, "5")),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												P(45, 9, 3, 20),
												ast.NewAssignmentExpressionNode(
													P(45, 9, 3, 20),
													T(P(49, 2, 3, 24), lexer.StarEqualToken),
													ast.NewPublicIdentifierNode(P(45, 3, 3, 20), "foo"),
													ast.NewIntLiteralNode(P(52, 2, 3, 27), V(P(52, 2, 3, 27), lexer.DecIntToken, "10")),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												P(55, 41, 4, 1),
												ast.NewIfExpressionNode(
													P(55, 41, 4, 1),
													ast.NewBinaryExpressionNode(
														P(61, 7, 4, 7),
														T(P(65, 1, 4, 11), lexer.LessToken),
														ast.NewPublicIdentifierNode(P(61, 3, 4, 7), "foo"),
														ast.NewIntLiteralNode(P(67, 1, 4, 13), V(P(67, 1, 4, 13), lexer.DecIntToken, "0")),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															P(74, 8, 4, 20),
															ast.NewAssignmentExpressionNode(
																P(74, 8, 4, 20),
																T(P(78, 2, 4, 24), lexer.PercentEqualToken),
																ast.NewPublicIdentifierNode(P(74, 3, 4, 20), "foo"),
																ast.NewIntLiteralNode(P(81, 1, 4, 27), V(P(81, 1, 4, 27), lexer.DecIntToken, "3")),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															P(88, 8, 5, 6),
															ast.NewAssignmentExpressionNode(
																P(88, 8, 5, 6),
																T(P(92, 2, 5, 10), lexer.MinusEqualToken),
																ast.NewPublicIdentifierNode(P(88, 3, 5, 6), "foo"),
																ast.NewIntLiteralNode(P(95, 1, 5, 13), V(P(95, 1, 5, 13), lexer.DecIntToken, "2")),
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
					),
					ast.NewExpressionStatementNode(
						P(97, 4, 6, 1),
						ast.NewNilLiteralNode(P(97, 3, 6, 1)),
					),
				},
			),
		},
		"else if is also possible": {
			input: `
if foo > 0
	foo += 2
	nil
else if foo < 5
	foo *= 10
else if foo < 0
	foo %= 3
else
	foo -= 2
	nil
end
nil
`,
			want: ast.NewProgramNode(
				P(0, 108, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 103, 2, 1),
						ast.NewIfExpressionNode(
							P(1, 102, 2, 1),
							ast.NewBinaryExpressionNode(
								P(4, 7, 2, 4),
								T(P(8, 1, 2, 8), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(4, 3, 2, 4), "foo"),
								ast.NewIntLiteralNode(P(10, 1, 2, 10), V(P(10, 1, 2, 10), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(13, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(13, 8, 3, 2),
										T(P(17, 2, 3, 6), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(13, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(20, 1, 3, 9), V(P(20, 1, 3, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(23, 4, 4, 2),
									ast.NewNilLiteralNode(P(23, 3, 4, 2)),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(32, 71, 5, 6),
									ast.NewIfExpressionNode(
										P(32, 71, 5, 6),
										ast.NewBinaryExpressionNode(
											P(35, 7, 5, 9),
											T(P(39, 1, 5, 13), lexer.LessToken),
											ast.NewPublicIdentifierNode(P(35, 3, 5, 9), "foo"),
											ast.NewIntLiteralNode(P(41, 1, 5, 15), V(P(41, 1, 5, 15), lexer.DecIntToken, "5")),
										),
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												P(44, 10, 6, 2),
												ast.NewAssignmentExpressionNode(
													P(44, 9, 6, 2),
													T(P(48, 2, 6, 6), lexer.StarEqualToken),
													ast.NewPublicIdentifierNode(P(44, 3, 6, 2), "foo"),
													ast.NewIntLiteralNode(P(51, 2, 6, 9), V(P(51, 2, 6, 9), lexer.DecIntToken, "10")),
												),
											),
										},
										[]ast.StatementNode{
											ast.NewExpressionStatementNode(
												P(59, 44, 7, 6),
												ast.NewIfExpressionNode(
													P(59, 44, 7, 6),
													ast.NewBinaryExpressionNode(
														P(62, 7, 7, 9),
														T(P(66, 1, 7, 13), lexer.LessToken),
														ast.NewPublicIdentifierNode(P(62, 3, 7, 9), "foo"),
														ast.NewIntLiteralNode(P(68, 1, 7, 15), V(P(68, 1, 7, 15), lexer.DecIntToken, "0")),
													),
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															P(71, 9, 8, 2),
															ast.NewAssignmentExpressionNode(
																P(71, 8, 8, 2),
																T(P(75, 2, 8, 6), lexer.PercentEqualToken),
																ast.NewPublicIdentifierNode(P(71, 3, 8, 2), "foo"),
																ast.NewIntLiteralNode(P(78, 1, 8, 9), V(P(78, 1, 8, 9), lexer.DecIntToken, "3")),
															),
														),
													},
													[]ast.StatementNode{
														ast.NewExpressionStatementNode(
															P(86, 9, 10, 2),
															ast.NewAssignmentExpressionNode(
																P(86, 8, 10, 2),
																T(P(90, 2, 10, 6), lexer.MinusEqualToken),
																ast.NewPublicIdentifierNode(P(86, 3, 10, 2), "foo"),
																ast.NewIntLiteralNode(P(93, 1, 10, 9), V(P(93, 1, 10, 9), lexer.DecIntToken, "2")),
															),
														),
														ast.NewExpressionStatementNode(
															P(96, 4, 11, 2),
															ast.NewNilLiteralNode(P(96, 3, 11, 2)),
														),
													},
												),
											),
										},
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(104, 4, 13, 1),
						ast.NewNilLiteralNode(P(104, 3, 13, 1)),
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

func TestUnless(t *testing.T) {
	tests := testTable{
		"can have one branch": {
			input: `
unless foo > 0
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				P(0, 35, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 34, 2, 1),
						ast.NewUnlessExpressionNode(
							P(1, 33, 2, 1),
							ast.NewBinaryExpressionNode(
								P(8, 7, 2, 8),
								T(P(12, 1, 2, 12), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(8, 3, 2, 8), "foo"),
								ast.NewIntLiteralNode(P(14, 1, 2, 14), V(P(14, 1, 2, 14), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(17, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(17, 8, 3, 2),
										T(P(21, 2, 3, 6), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(17, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(24, 1, 3, 9), V(P(24, 1, 3, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(27, 4, 4, 2),
									ast.NewNilLiteralNode(P(27, 3, 4, 2)),
								),
							},
							nil,
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
unless foo > 0
end
`,
			want: ast.NewProgramNode(
				P(0, 20, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 19, 2, 1),
						ast.NewUnlessExpressionNode(
							P(1, 18, 2, 1),
							ast.NewBinaryExpressionNode(
								P(8, 7, 2, 8),
								T(P(12, 1, 2, 12), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(8, 3, 2, 8), "foo"),
								ast.NewIntLiteralNode(P(14, 1, 2, 14), V(P(14, 1, 2, 14), lexer.DecIntToken, "0")),
							),
							nil,
							nil,
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
bar =
	unless foo > 0
		foo += 2
	end
nil
`,
			want: ast.NewProgramNode(
				P(0, 43, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 38, 2, 1),
						ast.NewAssignmentExpressionNode(
							P(1, 37, 2, 1),
							T(P(5, 1, 2, 5), lexer.EqualToken),
							ast.NewPublicIdentifierNode(P(1, 3, 2, 1), "bar"),
							ast.NewUnlessExpressionNode(
								P(8, 30, 3, 2),
								ast.NewBinaryExpressionNode(
									P(15, 7, 3, 9),
									T(P(19, 1, 3, 13), lexer.GreaterToken),
									ast.NewPublicIdentifierNode(P(15, 3, 3, 9), "foo"),
									ast.NewIntLiteralNode(P(21, 1, 3, 15), V(P(21, 1, 3, 15), lexer.DecIntToken, "0")),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										P(25, 9, 4, 3),
										ast.NewAssignmentExpressionNode(
											P(25, 8, 4, 3),
											T(P(29, 2, 4, 7), lexer.PlusEqualToken),
											ast.NewPublicIdentifierNode(P(25, 3, 4, 3), "foo"),
											ast.NewIntLiteralNode(P(32, 1, 4, 10), V(P(32, 1, 4, 10), lexer.DecIntToken, "2")),
										),
									),
								},
								nil,
							),
						),
					),
					ast.NewExpressionStatementNode(
						P(39, 4, 6, 1),
						ast.NewNilLiteralNode(P(39, 3, 6, 1)),
					),
				},
			),
		},
		"can be single line with then and without end": {
			input: `
unless foo > 0 then foo += 2
nil
`,
			want: ast.NewProgramNode(
				P(0, 34, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 29, 2, 1),
						ast.NewUnlessExpressionNode(
							P(1, 28, 2, 1),
							ast.NewBinaryExpressionNode(
								P(8, 7, 2, 8),
								T(P(12, 1, 2, 12), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(8, 3, 2, 8), "foo"),
								ast.NewIntLiteralNode(P(14, 1, 2, 14), V(P(14, 1, 2, 14), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(21, 8, 2, 21),
									ast.NewAssignmentExpressionNode(
										P(21, 8, 2, 21),
										T(P(25, 2, 2, 25), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(21, 3, 2, 21), "foo"),
										ast.NewIntLiteralNode(P(28, 1, 2, 28), V(P(28, 1, 2, 28), lexer.DecIntToken, "2")),
									),
								),
							},
							nil,
						),
					),
					ast.NewExpressionStatementNode(
						P(30, 4, 3, 1),
						ast.NewNilLiteralNode(P(30, 3, 3, 1)),
					),
				},
			),
		},
		"can have else": {
			input: `
unless foo > 0
	foo += 2
	nil
else
	foo -= 2
	nil
end
nil
`,
			want: ast.NewProgramNode(
				P(0, 59, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 54, 2, 1),
						ast.NewUnlessExpressionNode(
							P(1, 53, 2, 1),
							ast.NewBinaryExpressionNode(
								P(8, 7, 2, 8),
								T(P(12, 1, 2, 12), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(8, 3, 2, 8), "foo"),
								ast.NewIntLiteralNode(P(14, 1, 2, 14), V(P(14, 1, 2, 14), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(17, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(17, 8, 3, 2),
										T(P(21, 2, 3, 6), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(17, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(24, 1, 3, 9), V(P(24, 1, 3, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(27, 4, 4, 2),
									ast.NewNilLiteralNode(P(27, 3, 4, 2)),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(37, 9, 6, 2),
									ast.NewAssignmentExpressionNode(
										P(37, 8, 6, 2),
										T(P(41, 2, 6, 6), lexer.MinusEqualToken),
										ast.NewPublicIdentifierNode(P(37, 3, 6, 2), "foo"),
										ast.NewIntLiteralNode(P(44, 1, 6, 9), V(P(44, 1, 6, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(47, 4, 7, 2),
									ast.NewNilLiteralNode(P(47, 3, 7, 2)),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(55, 4, 9, 1),
						ast.NewNilLiteralNode(P(55, 3, 9, 1)),
					),
				},
			),
		},
		"can have else in short form": {
			input: `
unless foo > 0 then foo += 2
else foo -= 2
nil
`,
			want: ast.NewProgramNode(
				P(0, 48, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 43, 2, 1),
						ast.NewUnlessExpressionNode(
							P(1, 42, 2, 1),
							ast.NewBinaryExpressionNode(
								P(8, 7, 2, 8),
								T(P(12, 1, 2, 12), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(8, 3, 2, 8), "foo"),
								ast.NewIntLiteralNode(P(14, 1, 2, 14), V(P(14, 1, 2, 14), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(21, 8, 2, 21),
									ast.NewAssignmentExpressionNode(
										P(21, 8, 2, 21),
										T(P(25, 2, 2, 25), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(21, 3, 2, 21), "foo"),
										ast.NewIntLiteralNode(P(28, 1, 2, 28), V(P(28, 1, 2, 28), lexer.DecIntToken, "2")),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(35, 8, 3, 6),
									ast.NewAssignmentExpressionNode(
										P(35, 8, 3, 6),
										T(P(39, 2, 3, 10), lexer.MinusEqualToken),
										ast.NewPublicIdentifierNode(P(35, 3, 3, 6), "foo"),
										ast.NewIntLiteralNode(P(42, 1, 3, 13), V(P(42, 1, 3, 13), lexer.DecIntToken, "2")),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(44, 4, 4, 1),
						ast.NewNilLiteralNode(P(44, 3, 4, 1)),
					),
				},
			),
		},
		"can't have two elses": {
			input: `
unless foo > 0 then foo += 2
else foo -= 2
else bar
nil
`,
			want: ast.NewProgramNode(
				P(0, 57, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 43, 2, 1),
						ast.NewUnlessExpressionNode(
							P(1, 42, 2, 1),
							ast.NewBinaryExpressionNode(
								P(8, 7, 2, 8),
								T(P(12, 1, 2, 12), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(8, 3, 2, 8), "foo"),
								ast.NewIntLiteralNode(P(14, 1, 2, 14), V(P(14, 1, 2, 14), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(21, 8, 2, 21),
									ast.NewAssignmentExpressionNode(
										P(21, 8, 2, 21),
										T(P(25, 2, 2, 25), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(21, 3, 2, 21), "foo"),
										ast.NewIntLiteralNode(P(28, 1, 2, 28), V(P(28, 1, 2, 28), lexer.DecIntToken, "2")),
									),
								),
							},
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(35, 8, 3, 6),
									ast.NewAssignmentExpressionNode(
										P(35, 8, 3, 6),
										T(P(39, 2, 3, 10), lexer.MinusEqualToken),
										ast.NewPublicIdentifierNode(P(35, 3, 3, 6), "foo"),
										ast.NewIntLiteralNode(P(42, 1, 3, 13), V(P(42, 1, 3, 13), lexer.DecIntToken, "2")),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(44, 9, 4, 1),
						ast.NewInvalidNode(P(44, 4, 4, 1), T(P(44, 4, 4, 1), lexer.ElseToken)),
					),
					ast.NewExpressionStatementNode(
						P(53, 4, 5, 1),
						ast.NewNilLiteralNode(P(53, 3, 5, 1)),
					),
				},
			),
			err: ErrorList{
				NewError(P(44, 4, 4, 1), "unexpected else, expected an expression"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestWhile(t *testing.T) {
	tests := testTable{
		"can have a multiline body": {
			input: `
while foo > 0
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				P(0, 34, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 33, 2, 1),
						ast.NewWhileExpressionNode(
							P(1, 32, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(16, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(16, 8, 3, 2),
										T(P(20, 2, 3, 6), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(16, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(23, 1, 3, 9), V(P(23, 1, 3, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(26, 4, 4, 2),
									ast.NewNilLiteralNode(P(26, 3, 4, 2)),
								),
							},
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
while foo > 0
end
`,
			want: ast.NewProgramNode(
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 18, 2, 1),
						ast.NewWhileExpressionNode(
							P(1, 17, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), lexer.DecIntToken, "0")),
							),
							nil,
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
bar =
	while foo > 0
		foo += 2
	end
nil
`,
			want: ast.NewProgramNode(
				P(0, 42, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 37, 2, 1),
						ast.NewAssignmentExpressionNode(
							P(1, 36, 2, 1),
							T(P(5, 1, 2, 5), lexer.EqualToken),
							ast.NewPublicIdentifierNode(P(1, 3, 2, 1), "bar"),
							ast.NewWhileExpressionNode(
								P(8, 29, 3, 2),
								ast.NewBinaryExpressionNode(
									P(14, 7, 3, 8),
									T(P(18, 1, 3, 12), lexer.GreaterToken),
									ast.NewPublicIdentifierNode(P(14, 3, 3, 8), "foo"),
									ast.NewIntLiteralNode(P(20, 1, 3, 14), V(P(20, 1, 3, 14), lexer.DecIntToken, "0")),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										P(24, 9, 4, 3),
										ast.NewAssignmentExpressionNode(
											P(24, 8, 4, 3),
											T(P(28, 2, 4, 7), lexer.PlusEqualToken),
											ast.NewPublicIdentifierNode(P(24, 3, 4, 3), "foo"),
											ast.NewIntLiteralNode(P(31, 1, 4, 10), V(P(31, 1, 4, 10), lexer.DecIntToken, "2")),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						P(38, 4, 6, 1),
						ast.NewNilLiteralNode(P(38, 3, 6, 1)),
					),
				},
			),
		},
		"can be single line with then and without end": {
			input: `
while foo > 0 then foo += 2
nil
`,
			want: ast.NewProgramNode(
				P(0, 33, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 28, 2, 1),
						ast.NewWhileExpressionNode(
							P(1, 27, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(20, 8, 2, 20),
									ast.NewAssignmentExpressionNode(
										P(20, 8, 2, 20),
										T(P(24, 2, 2, 24), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(20, 3, 2, 20), "foo"),
										ast.NewIntLiteralNode(P(27, 1, 2, 27), V(P(27, 1, 2, 27), lexer.DecIntToken, "2")),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(29, 4, 3, 1),
						ast.NewNilLiteralNode(P(29, 3, 3, 1)),
					),
				},
			),
		},
		"can't have else": {
			input: `
while foo > 0
	foo += 2
	nil
else
	foo -= 2
	nil
end
nil
`,
			want: ast.NewProgramNode(
				P(0, 58, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 53, 2, 1),
						ast.NewWhileExpressionNode(
							P(1, 52, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(16, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(16, 8, 3, 2),
										T(P(20, 2, 3, 6), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(16, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(23, 1, 3, 9), V(P(23, 1, 3, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(26, 4, 4, 2),
									ast.NewNilLiteralNode(P(26, 3, 4, 2)),
								),
								ast.NewExpressionStatementNode(
									P(30, 5, 5, 1),
									ast.NewInvalidNode(P(30, 4, 5, 1), T(P(30, 4, 5, 1), lexer.ElseToken)),
								),
								ast.NewExpressionStatementNode(
									P(36, 9, 6, 2),
									ast.NewAssignmentExpressionNode(
										P(36, 8, 6, 2),
										T(P(40, 2, 6, 6), lexer.MinusEqualToken),
										ast.NewPublicIdentifierNode(P(36, 3, 6, 2), "foo"),
										ast.NewIntLiteralNode(P(43, 1, 6, 9), V(P(43, 1, 6, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(46, 4, 7, 2),
									ast.NewNilLiteralNode(P(46, 3, 7, 2)),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(54, 4, 9, 1),
						ast.NewNilLiteralNode(P(54, 3, 9, 1)),
					),
				},
			),
			err: ErrorList{
				NewError(P(30, 4, 5, 1), "unexpected else, expected an expression"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestUntil(t *testing.T) {
	tests := testTable{
		"can have a multiline body": {
			input: `
until foo > 0
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				P(0, 34, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 33, 2, 1),
						ast.NewUntilExpressionNode(
							P(1, 32, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(16, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(16, 8, 3, 2),
										T(P(20, 2, 3, 6), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(16, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(23, 1, 3, 9), V(P(23, 1, 3, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(26, 4, 4, 2),
									ast.NewNilLiteralNode(P(26, 3, 4, 2)),
								),
							},
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
until foo > 0
end
`,
			want: ast.NewProgramNode(
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 18, 2, 1),
						ast.NewUntilExpressionNode(
							P(1, 17, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), lexer.DecIntToken, "0")),
							),
							nil,
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
bar =
	until foo > 0
		foo += 2
	end
nil
`,
			want: ast.NewProgramNode(
				P(0, 42, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 37, 2, 1),
						ast.NewAssignmentExpressionNode(
							P(1, 36, 2, 1),
							T(P(5, 1, 2, 5), lexer.EqualToken),
							ast.NewPublicIdentifierNode(P(1, 3, 2, 1), "bar"),
							ast.NewUntilExpressionNode(
								P(8, 29, 3, 2),
								ast.NewBinaryExpressionNode(
									P(14, 7, 3, 8),
									T(P(18, 1, 3, 12), lexer.GreaterToken),
									ast.NewPublicIdentifierNode(P(14, 3, 3, 8), "foo"),
									ast.NewIntLiteralNode(P(20, 1, 3, 14), V(P(20, 1, 3, 14), lexer.DecIntToken, "0")),
								),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										P(24, 9, 4, 3),
										ast.NewAssignmentExpressionNode(
											P(24, 8, 4, 3),
											T(P(28, 2, 4, 7), lexer.PlusEqualToken),
											ast.NewPublicIdentifierNode(P(24, 3, 4, 3), "foo"),
											ast.NewIntLiteralNode(P(31, 1, 4, 10), V(P(31, 1, 4, 10), lexer.DecIntToken, "2")),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						P(38, 4, 6, 1),
						ast.NewNilLiteralNode(P(38, 3, 6, 1)),
					),
				},
			),
		},
		"can be single line with then and without end": {
			input: `
until foo > 0 then foo += 2
nil
`,
			want: ast.NewProgramNode(
				P(0, 33, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 28, 2, 1),
						ast.NewUntilExpressionNode(
							P(1, 27, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(20, 8, 2, 20),
									ast.NewAssignmentExpressionNode(
										P(20, 8, 2, 20),
										T(P(24, 2, 2, 24), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(20, 3, 2, 20), "foo"),
										ast.NewIntLiteralNode(P(27, 1, 2, 27), V(P(27, 1, 2, 27), lexer.DecIntToken, "2")),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(29, 4, 3, 1),
						ast.NewNilLiteralNode(P(29, 3, 3, 1)),
					),
				},
			),
		},
		"can't have else": {
			input: `
until foo > 0
	foo += 2
	nil
else
	foo -= 2
	nil
end
nil
`,
			want: ast.NewProgramNode(
				P(0, 58, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 53, 2, 1),
						ast.NewUntilExpressionNode(
							P(1, 52, 2, 1),
							ast.NewBinaryExpressionNode(
								P(7, 7, 2, 7),
								T(P(11, 1, 2, 11), lexer.GreaterToken),
								ast.NewPublicIdentifierNode(P(7, 3, 2, 7), "foo"),
								ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), lexer.DecIntToken, "0")),
							),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(16, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(16, 8, 3, 2),
										T(P(20, 2, 3, 6), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(16, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(23, 1, 3, 9), V(P(23, 1, 3, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(26, 4, 4, 2),
									ast.NewNilLiteralNode(P(26, 3, 4, 2)),
								),
								ast.NewExpressionStatementNode(
									P(30, 5, 5, 1),
									ast.NewInvalidNode(P(30, 4, 5, 1), T(P(30, 4, 5, 1), lexer.ElseToken)),
								),
								ast.NewExpressionStatementNode(
									P(36, 9, 6, 2),
									ast.NewAssignmentExpressionNode(
										P(36, 8, 6, 2),
										T(P(40, 2, 6, 6), lexer.MinusEqualToken),
										ast.NewPublicIdentifierNode(P(36, 3, 6, 2), "foo"),
										ast.NewIntLiteralNode(P(43, 1, 6, 9), V(P(43, 1, 6, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(46, 4, 7, 2),
									ast.NewNilLiteralNode(P(46, 3, 7, 2)),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(54, 4, 9, 1),
						ast.NewNilLiteralNode(P(54, 3, 9, 1)),
					),
				},
			),
			err: ErrorList{
				NewError(P(30, 4, 5, 1), "unexpected else, expected an expression"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestLoop(t *testing.T) {
	tests := testTable{
		"can have a multiline body": {
			input: `
loop
	foo += 2
	nil
end
`,
			want: ast.NewProgramNode(
				P(0, 25, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 24, 2, 1),
						ast.NewLoopExpressionNode(
							P(1, 23, 2, 1),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(7, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(7, 8, 3, 2),
										T(P(11, 2, 3, 6), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(7, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(14, 1, 3, 9), V(P(14, 1, 3, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(17, 4, 4, 2),
									ast.NewNilLiteralNode(P(17, 3, 4, 2)),
								),
							},
						),
					),
				},
			),
		},
		"can have an empty body": {
			input: `
loop
end
`,
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 9, 2, 1),
						ast.NewLoopExpressionNode(
							P(1, 8, 2, 1),
							nil,
						),
					),
				},
			),
		},
		"is an expression": {
			input: `
bar =
	loop
		foo += 2
	end
nil
`,
			want: ast.NewProgramNode(
				P(0, 33, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 28, 2, 1),
						ast.NewAssignmentExpressionNode(
							P(1, 27, 2, 1),
							T(P(5, 1, 2, 5), lexer.EqualToken),
							ast.NewPublicIdentifierNode(P(1, 3, 2, 1), "bar"),
							ast.NewLoopExpressionNode(
								P(8, 20, 3, 2),
								[]ast.StatementNode{
									ast.NewExpressionStatementNode(
										P(15, 9, 4, 3),
										ast.NewAssignmentExpressionNode(
											P(15, 8, 4, 3),
											T(P(19, 2, 4, 7), lexer.PlusEqualToken),
											ast.NewPublicIdentifierNode(P(15, 3, 4, 3), "foo"),
											ast.NewIntLiteralNode(P(22, 1, 4, 10), V(P(22, 1, 4, 10), lexer.DecIntToken, "2")),
										),
									),
								},
							),
						),
					),
					ast.NewExpressionStatementNode(
						P(29, 4, 6, 1),
						ast.NewNilLiteralNode(P(29, 3, 6, 1)),
					),
				},
			),
		},
		"can be single line without end": {
			input: `
loop foo += 2
nil
`,
			want: ast.NewProgramNode(
				P(0, 19, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 14, 2, 1),
						ast.NewLoopExpressionNode(
							P(1, 13, 2, 1),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(6, 8, 2, 6),
									ast.NewAssignmentExpressionNode(
										P(6, 8, 2, 6),
										T(P(10, 2, 2, 10), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(6, 3, 2, 6), "foo"),
										ast.NewIntLiteralNode(P(13, 1, 2, 13), V(P(13, 1, 2, 13), lexer.DecIntToken, "2")),
									),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(15, 4, 3, 1),
						ast.NewNilLiteralNode(P(15, 3, 3, 1)),
					),
				},
			),
		},
		"can't have else": {
			input: `
loop
	foo += 2
	nil
else
	foo -= 2
	nil
end
nil
`,
			want: ast.NewProgramNode(
				P(0, 49, 1, 1),
				[]ast.StatementNode{
					ast.NewEmptyStatementNode(P(0, 1, 1, 1)),
					ast.NewExpressionStatementNode(
						P(1, 44, 2, 1),
						ast.NewLoopExpressionNode(
							P(1, 43, 2, 1),
							[]ast.StatementNode{
								ast.NewExpressionStatementNode(
									P(7, 9, 3, 2),
									ast.NewAssignmentExpressionNode(
										P(7, 8, 3, 2),
										T(P(11, 2, 3, 6), lexer.PlusEqualToken),
										ast.NewPublicIdentifierNode(P(7, 3, 3, 2), "foo"),
										ast.NewIntLiteralNode(P(14, 1, 3, 9), V(P(14, 1, 3, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(17, 4, 4, 2),
									ast.NewNilLiteralNode(P(17, 3, 4, 2)),
								),
								ast.NewExpressionStatementNode(
									P(21, 5, 5, 1),
									ast.NewInvalidNode(P(21, 4, 5, 1), T(P(21, 4, 5, 1), lexer.ElseToken)),
								),
								ast.NewExpressionStatementNode(
									P(27, 9, 6, 2),
									ast.NewAssignmentExpressionNode(
										P(27, 8, 6, 2),
										T(P(31, 2, 6, 6), lexer.MinusEqualToken),
										ast.NewPublicIdentifierNode(P(27, 3, 6, 2), "foo"),
										ast.NewIntLiteralNode(P(34, 1, 6, 9), V(P(34, 1, 6, 9), lexer.DecIntToken, "2")),
									),
								),
								ast.NewExpressionStatementNode(
									P(37, 4, 7, 2),
									ast.NewNilLiteralNode(P(37, 3, 7, 2)),
								),
							},
						),
					),
					ast.NewExpressionStatementNode(
						P(45, 4, 9, 1),
						ast.NewNilLiteralNode(P(45, 3, 9, 1)),
					),
				},
			),
			err: ErrorList{
				NewError(P(21, 4, 5, 1), "unexpected else, expected an expression"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestBreak(t *testing.T) {
	tests := testTable{
		"can stand alone": {
			input: `break`,
			want: ast.NewProgramNode(
				P(0, 5, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewBreakExpressionNode(P(0, 5, 1, 1)),
					),
				},
			),
		},
		"can't have an argument": {
			input: `break 2`,
			want: ast.NewProgramNode(
				P(0, 7, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewBreakExpressionNode(P(0, 5, 1, 1)),
					),
				},
			),
			err: ErrorList{
				NewError(
					P(6, 1, 1, 7),
					"unexpected DecInt, expected a statement separator `\\n`, `;` or end of file",
				),
			},
		},
		"is an expression": {
			input: `foo && break`,
			want: ast.NewProgramNode(
				P(0, 12, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 12, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 12, 1, 1),
							T(P(4, 2, 1, 5), lexer.AndAndToken),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewBreakExpressionNode(P(7, 5, 1, 8)),
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

func TestReturn(t *testing.T) {
	tests := testTable{
		"can stand alone at the end": {
			input: `return`,
			want: ast.NewProgramNode(
				P(0, 6, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewReturnExpressionNode(P(0, 6, 1, 1), nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "return\n1",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 7, 1, 1),
						ast.NewReturnExpressionNode(P(0, 6, 1, 1), nil),
					),
					ast.NewExpressionStatementNode(
						P(7, 1, 2, 1),
						ast.NewIntLiteralNode(P(7, 1, 2, 1), V(P(7, 1, 2, 1), lexer.DecIntToken, "1")),
					),
				},
			),
		},
		"can have an argument": {
			input: `return 2`,
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewReturnExpressionNode(
							P(0, 8, 1, 1),
							ast.NewIntLiteralNode(P(7, 1, 1, 8), V(P(7, 1, 1, 8), lexer.DecIntToken, "2")),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && return`,
			want: ast.NewProgramNode(
				P(0, 13, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 13, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 13, 1, 1),
							T(P(4, 2, 1, 5), lexer.AndAndToken),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewReturnExpressionNode(P(7, 6, 1, 8), nil),
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

func TestContinue(t *testing.T) {
	tests := testTable{
		"can stand alone at the end": {
			input: `continue`,
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewContinueExpressionNode(P(0, 8, 1, 1), nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "continue\n1",
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewContinueExpressionNode(P(0, 8, 1, 1), nil),
					),
					ast.NewExpressionStatementNode(
						P(9, 1, 2, 1),
						ast.NewIntLiteralNode(P(9, 1, 2, 1), V(P(9, 1, 2, 1), lexer.DecIntToken, "1")),
					),
				},
			),
		},
		"can have an argument": {
			input: `continue 2`,
			want: ast.NewProgramNode(
				P(0, 10, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 10, 1, 1),
						ast.NewContinueExpressionNode(
							P(0, 10, 1, 1),
							ast.NewIntLiteralNode(P(9, 1, 1, 10), V(P(9, 1, 1, 10), lexer.DecIntToken, "2")),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && continue`,
			want: ast.NewProgramNode(
				P(0, 15, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 15, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 15, 1, 1),
							T(P(4, 2, 1, 5), lexer.AndAndToken),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewContinueExpressionNode(P(7, 8, 1, 8), nil),
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

func TestThrow(t *testing.T) {
	tests := testTable{
		"can stand alone at the end": {
			input: `throw`,
			want: ast.NewProgramNode(
				P(0, 5, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 5, 1, 1),
						ast.NewThrowExpressionNode(P(0, 5, 1, 1), nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "throw\n1",
			want: ast.NewProgramNode(
				P(0, 7, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 6, 1, 1),
						ast.NewThrowExpressionNode(P(0, 5, 1, 1), nil),
					),
					ast.NewExpressionStatementNode(
						P(6, 1, 2, 1),
						ast.NewIntLiteralNode(P(6, 1, 2, 1), V(P(6, 1, 2, 1), lexer.DecIntToken, "1")),
					),
				},
			),
		},
		"can have an argument": {
			input: `throw 2`,
			want: ast.NewProgramNode(
				P(0, 7, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 7, 1, 1),
						ast.NewThrowExpressionNode(
							P(0, 7, 1, 1),
							ast.NewIntLiteralNode(P(6, 1, 1, 7), V(P(6, 1, 1, 7), lexer.DecIntToken, "2")),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && throw`,
			want: ast.NewProgramNode(
				P(0, 12, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 12, 1, 1),
						ast.NewLogicalExpressionNode(
							P(0, 12, 1, 1),
							T(P(4, 2, 1, 5), lexer.AndAndToken),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewThrowExpressionNode(P(7, 5, 1, 8), nil),
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

func TestVariableDeclaration(t *testing.T) {
	tests := testTable{
		"is valid without type or initialiser": {
			input: "var foo",
			want: ast.NewProgramNode(
				P(0, 7, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 7, 1, 1),
						ast.NewVariableDeclarationNode(
							P(0, 7, 1, 1),
							V(P(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							nil,
							nil,
						),
					),
				},
			),
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
							V(P(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							nil,
							ast.NewIntLiteralNode(P(10, 1, 1, 11), V(P(10, 1, 1, 11), lexer.DecIntToken, "5")),
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
							V(P(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							ast.NewPublicConstantNode(P(9, 3, 1, 10), "Int"),
							ast.NewIntLiteralNode(P(15, 1, 1, 16), V(P(15, 1, 1, 16), lexer.DecIntToken, "5")),
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
							V(P(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
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
							V(P(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
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
							V(P(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							ast.NewBinaryTypeExpressionNode(
								P(9, 12, 1, 10),
								T(P(13, 1, 1, 14), lexer.OrToken),
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
							V(P(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							ast.NewBinaryTypeExpressionNode(
								P(9, 21, 1, 10),
								T(P(22, 1, 1, 23), lexer.OrToken),
								ast.NewBinaryTypeExpressionNode(
									P(9, 12, 1, 10),
									T(P(13, 1, 1, 14), lexer.OrToken),
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
							V(P(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							ast.NewNilableTypeNode(
								P(10, 14, 1, 11),
								ast.NewBinaryTypeExpressionNode(
									P(10, 12, 1, 11),
									T(P(14, 1, 1, 15), lexer.OrToken),
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
							V(P(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							ast.NewBinaryTypeExpressionNode(
								P(9, 12, 1, 10),
								T(P(13, 1, 1, 14), lexer.AndToken),
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
							V(P(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							ast.NewBinaryTypeExpressionNode(
								P(9, 21, 1, 10),
								T(P(22, 1, 1, 23), lexer.AndToken),
								ast.NewBinaryTypeExpressionNode(
									P(9, 12, 1, 10),
									T(P(13, 1, 1, 14), lexer.AndToken),
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
							V(P(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							ast.NewNilableTypeNode(
								P(10, 14, 1, 11),
								ast.NewBinaryTypeExpressionNode(
									P(10, 12, 1, 11),
									T(P(14, 1, 1, 15), lexer.AndToken),
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			parserTest(tc, t)
		})
	}
}

func TestConstantLookup(t *testing.T) {
	tests := testTable{
		"is executed from left to right": {
			input: "Foo::Bar::Baz",
			want: ast.NewProgramNode(
				P(0, 13, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 13, 1, 1),
						ast.NewConstantLookupNode(
							P(0, 13, 1, 1),
							ast.NewConstantLookupNode(
								P(0, 8, 1, 1),
								ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
								ast.NewPublicConstantNode(P(5, 3, 1, 6), "Bar"),
							),
							ast.NewPublicConstantNode(P(10, 3, 1, 11), "Baz"),
						),
					),
				},
			),
		},
		"can't access private constants from outside": {
			input: "Foo::_Bar",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewConstantLookupNode(
							P(0, 9, 1, 1),
							ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
							ast.NewPrivateConstantNode(P(5, 4, 1, 6), "_Bar"),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(5, 4, 1, 6), "unexpected PrivateConstant, can't access a private constant from the outside"),
			},
		},
		"can have newlines after the operator": {
			input: "Foo::\nBar",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 9, 1, 1),
						ast.NewConstantLookupNode(
							P(0, 9, 1, 1),
							ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
							ast.NewPublicConstantNode(P(6, 3, 2, 1), "Bar"),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "Foo\n::Bar",
			want: ast.NewProgramNode(
				P(0, 9, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 4, 1, 1),
						ast.NewPublicConstantNode(P(0, 3, 1, 1), "Foo"),
					),
					ast.NewExpressionStatementNode(
						P(4, 5, 2, 1),
						ast.NewInvalidNode(P(4, 2, 2, 1), T(P(4, 2, 2, 1), lexer.ScopeResOpToken)),
					),
				},
			),
			err: ErrorList{
				NewError(P(4, 2, 2, 1), "unexpected ::, expected an expression"),
			},
		},
		"can have other primary expressions as the left side": {
			input: "foo::Bar",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewConstantLookupNode(
							P(0, 8, 1, 1),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewPublicConstantNode(P(5, 3, 1, 6), "Bar"),
						),
					),
				},
			),
		},
		"must have a constant as the right side": {
			input: "foo::123",
			want: ast.NewProgramNode(
				P(0, 8, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 8, 1, 1),
						ast.NewConstantLookupNode(
							P(0, 8, 1, 1),
							ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
							ast.NewInvalidNode(P(5, 3, 1, 6), V(P(5, 3, 1, 6), lexer.DecIntToken, "123")),
						),
					),
				},
			),
			err: ErrorList{
				NewError(P(5, 3, 1, 6), "unexpected DecInt, expected a constant"),
			},
		},
		"can be a part of an expression": {
			input: "foo::Bar + .3",
			want: ast.NewProgramNode(
				P(0, 13, 1, 1),
				[]ast.StatementNode{
					ast.NewExpressionStatementNode(
						P(0, 13, 1, 1),
						ast.NewBinaryExpressionNode(
							P(0, 13, 1, 1),
							T(P(9, 1, 1, 10), lexer.PlusToken),
							ast.NewConstantLookupNode(
								P(0, 8, 1, 1),
								ast.NewPublicIdentifierNode(P(0, 3, 1, 1), "foo"),
								ast.NewPublicConstantNode(P(5, 3, 1, 6), "Bar"),
							),
							ast.NewFloatLiteralNode(P(11, 2, 1, 12), "0.3"),
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

// func TestX(t *testing.T) {
// 	tests := testTable{}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			parserTest(tc, t)
// 		})
// 	}
// }
