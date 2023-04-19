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

// Create a new position in tests.
func P(startByte, byteLength, line, column int) lexer.Position {
	return lexer.Position{
		StartByte:  startByte,
		ByteLength: byteLength,
		Line:       line,
		Column:     column,
	}
}

// Create a new integer node.
func Int(tokenType lexer.TokenType, value string, startByte, byteLength, line, column int) *ast.IntLiteralNode {
	return &ast.IntLiteralNode{
		Position: P(startByte, byteLength, line, column),
		Token:    V(tokenType, value, startByte, byteLength, line, column),
	}
}

// Create a new float node.
func Float(value string, startByte, byteLength, line, column int) *ast.FloatLiteralNode {
	return &ast.FloatLiteralNode{
		Position: P(startByte, byteLength, line, column),
		Value:    value,
	}
}

// Create an invalid node.
func Invalid(token *lexer.Token) *ast.InvalidNode {
	return &ast.InvalidNode{
		Position: token.Position,
		Token:    token,
	}
}

// Create a unary expression node.
func Unary(pos lexer.Position, op *lexer.Token, right ast.ExpressionNode) *ast.UnaryExpressionNode {
	return &ast.UnaryExpressionNode{
		Position: pos,
		Op:       op,
		Right:    right,
	}
}

// Create a binary expression node.
func Bin(pos lexer.Position, op *lexer.Token, left ast.ExpressionNode, right ast.ExpressionNode) *ast.BinaryExpressionNode {
	return &ast.BinaryExpressionNode{
		Position: pos,
		Left:     left,
		Op:       op,
		Right:    right,
	}
}

// Create a new token in tests.
var T = lexer.NewToken

// Create a new token with value in tests.
var V = lexer.NewTokenWithValue

// Function which powers all parser tests.
// Inspects if the produced AST matches the expected one.
func parserTest(tc testCase, t *testing.T) {
	ast, err := Parse([]byte(tc.input))

	if diff := cmp.Diff(tc.want, ast); diff != "" {
		t.Fatalf(diff)
	}

	if diff := cmp.Diff(tc.err, err); diff != "" {
		t.Fatalf(diff)
	}
}

func TestAddition(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "1 + 2 + 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: Bin(
							P(0, 9, 1, 1),
							T(lexer.PlusToken, 6, 1, 1, 7),
							Bin(
								P(0, 5, 1, 1),
								T(lexer.PlusToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 1, 5),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 1, 9),
						),
					},
				},
			},
		},
		"can have newlines after the operator": {
			input: "1 +\n2 +\n3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: Bin(
							P(0, 9, 1, 1),
							T(lexer.PlusToken, 6, 1, 2, 3),
							Bin(
								P(0, 5, 1, 1),
								T(lexer.PlusToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 2, 1),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 3, 1),
						),
					},
				},
			},
		},
		"can't have newlines before the operator": {
			input: "1\n+ 2\n+ 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position:   P(0, 2, 1, 1),
						Expression: Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
					},
					&ast.ExpressionStatementNode{
						Position: P(2, 4, 2, 1),
						Expression: Unary(
							P(2, 3, 2, 1),
							T(lexer.PlusToken, 2, 1, 2, 1),
							Int(lexer.DecIntToken, "2", 4, 1, 2, 3),
						),
					},
					&ast.ExpressionStatementNode{
						Position: P(6, 3, 3, 1),
						Expression: Unary(
							P(6, 3, 3, 1),
							T(lexer.PlusToken, 6, 1, 3, 1),
							Int(lexer.DecIntToken, "3", 8, 1, 3, 3),
						),
					},
				},
			},
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
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: Bin(
							P(0, 9, 1, 1),
							T(lexer.MinusToken, 6, 1, 1, 7),
							Bin(
								P(0, 5, 1, 1),
								T(lexer.MinusToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 1, 5),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 1, 9),
						),
					},
				},
			},
		},
		"can have newlines after the operator": {
			input: "1 -\n2 -\n3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: Bin(
							P(0, 9, 1, 1),
							T(lexer.MinusToken, 6, 1, 2, 3),
							Bin(
								P(0, 5, 1, 1),
								T(lexer.MinusToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 2, 1),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 3, 1),
						),
					},
				},
			},
		},
		"can't have newlines before the operator": {
			input: "1\n- 2\n- 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position:   P(0, 2, 1, 1),
						Expression: Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
					},
					&ast.ExpressionStatementNode{
						Position: P(2, 4, 2, 1),
						Expression: Unary(
							P(2, 3, 2, 1),
							T(lexer.MinusToken, 2, 1, 2, 1),
							Int(lexer.DecIntToken, "2", 4, 1, 2, 3),
						),
					},
					&ast.ExpressionStatementNode{
						Position: P(6, 3, 3, 1),
						Expression: Unary(
							P(6, 3, 3, 1),
							T(lexer.MinusToken, 6, 1, 3, 1),
							Int(lexer.DecIntToken, "3", 8, 1, 3, 3),
						),
					},
				},
			},
		},
		"has the same precedence as addition": {
			input: "1 + 2 - 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: Bin(
							P(0, 9, 1, 1),
							T(lexer.MinusToken, 6, 1, 1, 7),
							Bin(
								P(0, 5, 1, 1),
								T(lexer.PlusToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 1, 5),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 1, 9),
						),
					},
				},
			},
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
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: Bin(
							P(0, 9, 1, 1),
							T(lexer.StarToken, 6, 1, 1, 7),
							Bin(
								P(0, 5, 1, 1),
								T(lexer.StarToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 1, 5),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 1, 9),
						),
					},
				},
			},
		},
		"can have newlines after the operator": {
			input: "1 *\n2 *\n3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: Bin(
							P(0, 9, 1, 1),
							T(lexer.StarToken, 6, 1, 2, 3),
							Bin(
								P(0, 5, 1, 1),
								T(lexer.StarToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 2, 1),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 3, 1),
						),
					},
				},
			},
		},
		"can't have newlines before the operator": {
			input: "1\n* 2\n* 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position:   P(0, 2, 1, 1),
						Expression: Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
					},
					&ast.ExpressionStatementNode{
						Position:   P(2, 1, 2, 1),
						Expression: Invalid(T(lexer.StarToken, 2, 1, 2, 1)),
					},
					&ast.ExpressionStatementNode{
						Position:   P(6, 1, 3, 1),
						Expression: Invalid(T(lexer.StarToken, 6, 1, 3, 1)),
					},
				},
			},
			err: ErrorList{
				&Error{P(2, 1, 2, 1), "Unexpected *, expected an expression"},
				&Error{P(6, 1, 3, 1), "Unexpected *, expected an expression"},
			},
		},
		"has higher precedence than addition": {
			input: "1 + 2 * 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: Bin(
							P(0, 9, 1, 1),
							T(lexer.PlusToken, 2, 1, 1, 3),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Bin(
								P(4, 5, 1, 5),
								T(lexer.StarToken, 6, 1, 1, 7),
								Int(lexer.DecIntToken, "2", 4, 1, 1, 5),
								Int(lexer.DecIntToken, "3", 8, 1, 1, 9),
							),
						),
					},
				},
			},
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
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: Bin(
							P(0, 9, 1, 1),
							T(lexer.SlashToken, 6, 1, 1, 7),
							Bin(
								P(0, 5, 1, 1),
								T(lexer.SlashToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 1, 5),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 1, 9),
						),
					},
				},
			},
		},
		"can have newlines after the operator": {
			input: "1 /\n2 /\n3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: Bin(
							P(0, 9, 1, 1),
							T(lexer.SlashToken, 6, 1, 2, 3),
							Bin(
								P(0, 5, 1, 1),
								T(lexer.SlashToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 2, 1),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 3, 1),
						),
					},
				},
			},
		},
		"can't have newlines before the operator": {
			input: "1\n/ 2\n/ 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position:   P(0, 2, 1, 1),
						Expression: Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
					},
					&ast.ExpressionStatementNode{
						Position:   P(2, 1, 2, 1),
						Expression: Invalid(T(lexer.SlashToken, 2, 1, 2, 1)),
					},
					&ast.ExpressionStatementNode{
						Position:   P(6, 1, 3, 1),
						Expression: Invalid(T(lexer.SlashToken, 6, 1, 3, 1)),
					},
				},
			},
			err: ErrorList{
				&Error{P(2, 1, 2, 1), "Unexpected /, expected an expression"},
				&Error{P(6, 1, 3, 1), "Unexpected /, expected an expression"},
			},
		},
		"has the same precedence as multiplication": {
			input: "1 * 2 / 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: Bin(
							P(0, 9, 1, 1),
							T(lexer.SlashToken, 6, 1, 1, 7),
							Bin(
								P(0, 5, 1, 1),
								T(lexer.StarToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 1, 5),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 1, 9),
						),
					},
				},
			},
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
			want: &ast.ProgramNode{
				Position: P(0, 11, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 11, 1, 1),
						Expression: Bin(
							P(0, 11, 1, 1),
							T(lexer.StarStarToken, 2, 2, 1, 3),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Bin(
								P(5, 6, 1, 6),
								T(lexer.StarStarToken, 7, 2, 1, 8),
								Int(lexer.DecIntToken, "2", 5, 1, 1, 6),
								Int(lexer.DecIntToken, "3", 10, 1, 1, 11),
							),
						),
					},
				},
			},
		},
		"can have newlines after the operator": {
			input: "1 **\n2 **\n3",
			want: &ast.ProgramNode{
				Position: P(0, 11, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 11, 1, 1),
						Expression: Bin(
							P(0, 11, 1, 1),
							T(lexer.StarStarToken, 2, 2, 1, 3),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Bin(
								P(5, 6, 2, 1),
								T(lexer.StarStarToken, 7, 2, 2, 3),
								Int(lexer.DecIntToken, "2", 5, 1, 2, 1),
								Int(lexer.DecIntToken, "3", 10, 1, 3, 1),
							),
						),
					},
				},
			},
		},
		"can't have newlines before the operator": {
			input: "1\n** 2\n** 3",
			want: &ast.ProgramNode{
				Position: P(0, 11, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position:   P(0, 2, 1, 1),
						Expression: Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
					},
					&ast.ExpressionStatementNode{
						Position:   P(2, 2, 2, 1),
						Expression: Invalid(T(lexer.StarStarToken, 2, 2, 2, 1)),
					},
					&ast.ExpressionStatementNode{
						Position:   P(7, 2, 3, 1),
						Expression: Invalid(T(lexer.StarStarToken, 7, 2, 3, 1)),
					},
				},
			},
			err: ErrorList{
				&Error{P(2, 2, 2, 1), "Unexpected **, expected an expression"},
				&Error{P(7, 2, 3, 1), "Unexpected **, expected an expression"},
			},
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
			want: &ast.ProgramNode{
				Position: P(0, 6, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 6, 1, 1),
						Expression: Unary(
							P(0, 6, 1, 1),
							T(lexer.PlusToken, 0, 1, 1, 1),
							Unary(
								P(1, 5, 1, 2),
								T(lexer.PlusToken, 1, 1, 1, 2),
								Unary(
									P(2, 4, 1, 3),
									T(lexer.PlusToken, 2, 1, 1, 3),
									Float("1.5", 3, 3, 1, 4),
								),
							),
						),
					},
				},
			},
		},
		"minus can be nested": {
			input: "---1.5",
			want: &ast.ProgramNode{
				Position: P(0, 6, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 6, 1, 1),
						Expression: Unary(
							P(0, 6, 1, 1),
							T(lexer.MinusToken, 0, 1, 1, 1),
							Unary(
								P(1, 5, 1, 2),
								T(lexer.MinusToken, 1, 1, 1, 2),
								Unary(
									P(2, 4, 1, 3),
									T(lexer.MinusToken, 2, 1, 1, 3),
									Float("1.5", 3, 3, 1, 4),
								),
							),
						),
					},
				},
			},
		},
		"logical not can be nested": {
			input: "!!!1.5",
			want: &ast.ProgramNode{
				Position: P(0, 6, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 6, 1, 1),
						Expression: Unary(
							P(0, 6, 1, 1),
							T(lexer.BangToken, 0, 1, 1, 1),
							Unary(
								P(1, 5, 1, 2),
								T(lexer.BangToken, 1, 1, 1, 2),
								Unary(
									P(2, 4, 1, 3),
									T(lexer.BangToken, 2, 1, 1, 3),
									Float("1.5", 3, 3, 1, 4),
								),
							),
						),
					},
				},
			},
		},
		"bitwise not can be nested": {
			input: "~~~1.5",
			want: &ast.ProgramNode{
				Position: P(0, 6, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 6, 1, 1),
						Expression: Unary(
							P(0, 6, 1, 1),
							T(lexer.TildeToken, 0, 1, 1, 1),
							Unary(
								P(1, 5, 1, 2),
								T(lexer.TildeToken, 1, 1, 1, 2),
								Unary(
									P(2, 4, 1, 3),
									T(lexer.TildeToken, 2, 1, 1, 3),
									Float("1.5", 3, 3, 1, 4),
								),
							),
						),
					},
				},
			},
		},
		"all have the same precedence": {
			input: "!+~1.5",
			want: &ast.ProgramNode{
				Position: P(0, 6, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 6, 1, 1),
						Expression: Unary(
							P(0, 6, 1, 1),
							T(lexer.BangToken, 0, 1, 1, 1),
							Unary(
								P(1, 5, 1, 2),
								T(lexer.PlusToken, 1, 1, 1, 2),
								Unary(
									P(2, 4, 1, 3),
									T(lexer.TildeToken, 2, 1, 1, 3),
									Float("1.5", 3, 3, 1, 4),
								),
							),
						),
					},
				},
			},
		},
		"have higher precedence than additive an multiplicative expression": {
			input: "!!1.5 * 2 + ~.5",
			want: &ast.ProgramNode{
				Position: P(0, 15, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 15, 1, 1),
						Expression: Bin(
							P(0, 15, 1, 1),
							T(lexer.PlusToken, 10, 1, 1, 11),
							Bin(
								P(0, 9, 1, 1),
								T(lexer.StarToken, 6, 1, 1, 7),
								Unary(
									P(0, 5, 1, 1),
									T(lexer.BangToken, 0, 1, 1, 1),
									Unary(
										P(1, 4, 1, 2),
										T(lexer.BangToken, 1, 1, 1, 2),
										Float("1.5", 2, 3, 1, 3),
									),
								),
								Int(lexer.DecIntToken, "2", 8, 1, 1, 9),
							),
							Unary(
								P(12, 3, 1, 13),
								T(lexer.TildeToken, 12, 1, 1, 13),
								Float("0.5", 13, 2, 1, 14),
							),
						),
					},
				},
			},
		},
		"have lower precedence than exponentiation": {
			input: "-2 ** 3",
			want: &ast.ProgramNode{
				Position: P(0, 7, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 7, 1, 1),
						Expression: Unary(
							P(0, 7, 1, 1),
							T(lexer.MinusToken, 0, 1, 1, 1),
							Bin(
								P(1, 6, 1, 2),
								T(lexer.StarStarToken, 3, 2, 1, 4),
								Int(lexer.DecIntToken, "2", 1, 1, 1, 2),
								Int(lexer.DecIntToken, "3", 6, 1, 1, 7),
							),
						),
					},
				},
			},
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
