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

func TestArithmetic(t *testing.T) {
	tests := testTable{
		"addition is evaluated from left to right": {
			input: "1 + 2 + 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: &ast.BinaryExpressionNode{
							Position: P(0, 9, 1, 1),
							Left: &ast.BinaryExpressionNode{
								Position: P(0, 5, 1, 1),
								Left: &ast.IntLiteralNode{
									Position: P(0, 1, 1, 1),
									Token:    V(lexer.DecIntToken, "1", 0, 1, 1, 1),
								},
								Op: T(lexer.PlusToken, 2, 1, 1, 3),
								Right: &ast.IntLiteralNode{
									Position: P(4, 1, 1, 5),
									Token:    V(lexer.DecIntToken, "2", 4, 1, 1, 5),
								},
							},
							Op: T(lexer.PlusToken, 6, 1, 1, 7),
							Right: &ast.IntLiteralNode{
								Position: P(8, 1, 1, 9),
								Token:    V(lexer.DecIntToken, "3", 8, 1, 1, 9),
							},
						},
					},
				},
			},
		},
		"addition can have newlines after the operator": {
			input: "1 +\n2 +\n3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: &ast.BinaryExpressionNode{
							Position: P(0, 9, 1, 1),
							Left: &ast.BinaryExpressionNode{
								Position: P(0, 5, 1, 1),
								Left: &ast.IntLiteralNode{
									Position: P(0, 1, 1, 1),
									Token:    V(lexer.DecIntToken, "1", 0, 1, 1, 1),
								},
								Op: T(lexer.PlusToken, 2, 1, 1, 3),
								Right: &ast.IntLiteralNode{
									Position: P(4, 1, 2, 1),
									Token:    V(lexer.DecIntToken, "2", 4, 1, 2, 1),
								},
							},
							Op: T(lexer.PlusToken, 6, 1, 2, 3),
							Right: &ast.IntLiteralNode{
								Position: P(8, 1, 3, 1),
								Token:    V(lexer.DecIntToken, "3", 8, 1, 3, 1),
							},
						},
					},
				},
			},
		},
		"subtraction is evaluated from left to right": {
			input: "1 - 2 - 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: &ast.BinaryExpressionNode{
							Position: P(0, 9, 1, 1),
							Left: &ast.BinaryExpressionNode{
								Position: P(0, 5, 1, 1),
								Left: &ast.IntLiteralNode{
									Position: P(0, 1, 1, 1),
									Token:    V(lexer.DecIntToken, "1", 0, 1, 1, 1),
								},
								Op: T(lexer.MinusToken, 2, 1, 1, 3),
								Right: &ast.IntLiteralNode{
									Position: P(4, 1, 1, 5),
									Token:    V(lexer.DecIntToken, "2", 4, 1, 1, 5),
								},
							},
							Op: T(lexer.MinusToken, 6, 1, 1, 7),
							Right: &ast.IntLiteralNode{
								Position: P(8, 1, 1, 9),
								Token:    V(lexer.DecIntToken, "3", 8, 1, 1, 9),
							},
						},
					},
				},
			},
		},
		"subtraction can have newlines after the operator": {
			input: "1 -\n2 -\n3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: &ast.BinaryExpressionNode{
							Position: P(0, 9, 1, 1),
							Left: &ast.BinaryExpressionNode{
								Position: P(0, 5, 1, 1),
								Left: &ast.IntLiteralNode{
									Position: P(0, 1, 1, 1),
									Token:    V(lexer.DecIntToken, "1", 0, 1, 1, 1),
								},
								Op: T(lexer.MinusToken, 2, 1, 1, 3),
								Right: &ast.IntLiteralNode{
									Position: P(4, 1, 2, 1),
									Token:    V(lexer.DecIntToken, "2", 4, 1, 2, 1),
								},
							},
							Op: T(lexer.MinusToken, 6, 1, 2, 3),
							Right: &ast.IntLiteralNode{
								Position: P(8, 1, 3, 1),
								Token:    V(lexer.DecIntToken, "3", 8, 1, 3, 1),
							},
						},
					},
				},
			},
		},
		"subtraction and addition have the same precedence": {
			input: "1 + 2 - 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: &ast.BinaryExpressionNode{
							Position: P(0, 9, 1, 1),
							Left: &ast.BinaryExpressionNode{
								Position: P(0, 5, 1, 1),
								Left: &ast.IntLiteralNode{
									Position: P(0, 1, 1, 1),
									Token:    V(lexer.DecIntToken, "1", 0, 1, 1, 1),
								},
								Op: T(lexer.PlusToken, 2, 1, 1, 3),
								Right: &ast.IntLiteralNode{
									Position: P(4, 1, 1, 5),
									Token:    V(lexer.DecIntToken, "2", 4, 1, 1, 5),
								},
							},
							Op: T(lexer.MinusToken, 6, 1, 1, 7),
							Right: &ast.IntLiteralNode{
								Position: P(8, 1, 1, 9),
								Token:    V(lexer.DecIntToken, "3", 8, 1, 1, 9),
							},
						},
					},
				},
			},
		},
		"multiplication is evaluated from left to right": {
			input: "1 * 2 * 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: &ast.BinaryExpressionNode{
							Position: P(0, 9, 1, 1),
							Left: &ast.BinaryExpressionNode{
								Position: P(0, 5, 1, 1),
								Left: &ast.IntLiteralNode{
									Position: P(0, 1, 1, 1),
									Token:    V(lexer.DecIntToken, "1", 0, 1, 1, 1),
								},
								Op: T(lexer.StarToken, 2, 1, 1, 3),
								Right: &ast.IntLiteralNode{
									Position: P(4, 1, 1, 5),
									Token:    V(lexer.DecIntToken, "2", 4, 1, 1, 5),
								},
							},
							Op: T(lexer.StarToken, 6, 1, 1, 7),
							Right: &ast.IntLiteralNode{
								Position: P(8, 1, 1, 9),
								Token:    V(lexer.DecIntToken, "3", 8, 1, 1, 9),
							},
						},
					},
				},
			},
		},
		"multiplication can have newlines after the operator": {
			input: "1 *\n2 *\n3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: &ast.BinaryExpressionNode{
							Position: P(0, 9, 1, 1),
							Left: &ast.BinaryExpressionNode{
								Position: P(0, 5, 1, 1),
								Left: &ast.IntLiteralNode{
									Position: P(0, 1, 1, 1),
									Token:    V(lexer.DecIntToken, "1", 0, 1, 1, 1),
								},
								Op: T(lexer.StarToken, 2, 1, 1, 3),
								Right: &ast.IntLiteralNode{
									Position: P(4, 1, 2, 1),
									Token:    V(lexer.DecIntToken, "2", 4, 1, 2, 1),
								},
							},
							Op: T(lexer.StarToken, 6, 1, 2, 3),
							Right: &ast.IntLiteralNode{
								Position: P(8, 1, 3, 1),
								Token:    V(lexer.DecIntToken, "3", 8, 1, 3, 1),
							},
						},
					},
				},
			},
		},
		"multiplication has higher precedence than addition": {
			input: "1 + 2 * 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: &ast.BinaryExpressionNode{
							Position: P(0, 9, 1, 1),
							Left: &ast.IntLiteralNode{
								Position: P(0, 1, 1, 1),
								Token:    V(lexer.DecIntToken, "1", 0, 1, 1, 1),
							},
							Op: T(lexer.PlusToken, 2, 1, 1, 3),
							Right: &ast.BinaryExpressionNode{
								Position: P(4, 5, 1, 5),
								Left: &ast.IntLiteralNode{
									Position: P(4, 1, 1, 5),
									Token:    V(lexer.DecIntToken, "2", 4, 1, 1, 5),
								},
								Op: T(lexer.StarToken, 6, 1, 1, 7),
								Right: &ast.IntLiteralNode{
									Position: P(8, 1, 1, 9),
									Token:    V(lexer.DecIntToken, "3", 8, 1, 1, 9),
								},
							},
						},
					},
				},
			},
		},
		"division is evaluated from left to right": {
			input: "1 / 2 / 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: &ast.BinaryExpressionNode{
							Position: P(0, 9, 1, 1),
							Left: &ast.BinaryExpressionNode{
								Position: P(0, 5, 1, 1),
								Left: &ast.IntLiteralNode{
									Position: P(0, 1, 1, 1),
									Token:    V(lexer.DecIntToken, "1", 0, 1, 1, 1),
								},
								Op: T(lexer.SlashToken, 2, 1, 1, 3),
								Right: &ast.IntLiteralNode{
									Position: P(4, 1, 1, 5),
									Token:    V(lexer.DecIntToken, "2", 4, 1, 1, 5),
								},
							},
							Op: T(lexer.SlashToken, 6, 1, 1, 7),
							Right: &ast.IntLiteralNode{
								Position: P(8, 1, 1, 9),
								Token:    V(lexer.DecIntToken, "3", 8, 1, 1, 9),
							},
						},
					},
				},
			},
		},
		"division can have newlines after the operator": {
			input: "1 /\n2 /\n3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: &ast.BinaryExpressionNode{
							Position: P(0, 9, 1, 1),
							Left: &ast.BinaryExpressionNode{
								Position: P(0, 5, 1, 1),
								Left: &ast.IntLiteralNode{
									Position: P(0, 1, 1, 1),
									Token:    V(lexer.DecIntToken, "1", 0, 1, 1, 1),
								},
								Op: T(lexer.SlashToken, 2, 1, 1, 3),
								Right: &ast.IntLiteralNode{
									Position: P(4, 1, 2, 1),
									Token:    V(lexer.DecIntToken, "2", 4, 1, 2, 1),
								},
							},
							Op: T(lexer.SlashToken, 6, 1, 2, 3),
							Right: &ast.IntLiteralNode{
								Position: P(8, 1, 3, 1),
								Token:    V(lexer.DecIntToken, "3", 8, 1, 3, 1),
							},
						},
					},
				},
			},
		},
		"division and multiplication have the same precedence": {
			input: "1 * 2 / 3",
			want: &ast.ProgramNode{
				Position: P(0, 9, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 9, 1, 1),
						Expression: &ast.BinaryExpressionNode{
							Position: P(0, 9, 1, 1),
							Left: &ast.BinaryExpressionNode{
								Position: P(0, 5, 1, 1),
								Left: &ast.IntLiteralNode{
									Position: P(0, 1, 1, 1),
									Token:    V(lexer.DecIntToken, "1", 0, 1, 1, 1),
								},
								Op: T(lexer.StarToken, 2, 1, 1, 3),
								Right: &ast.IntLiteralNode{
									Position: P(4, 1, 1, 5),
									Token:    V(lexer.DecIntToken, "2", 4, 1, 1, 5),
								},
							},
							Op: T(lexer.SlashToken, 6, 1, 1, 7),
							Right: &ast.IntLiteralNode{
								Position: P(8, 1, 1, 9),
								Token:    V(lexer.DecIntToken, "3", 8, 1, 1, 9),
							},
						},
					},
				},
			},
		},
		"exponentiation is evaluated from right to left": {
			input: "1 ** 2 ** 3",
			want: &ast.ProgramNode{
				Position: P(0, 11, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 11, 1, 1),
						Expression: &ast.BinaryExpressionNode{
							Position: P(0, 11, 1, 1),
							Left: &ast.IntLiteralNode{
								Position: P(0, 1, 1, 1),
								Token:    V(lexer.DecIntToken, "1", 0, 1, 1, 1),
							},
							Op: T(lexer.PowerToken, 2, 2, 1, 3),
							Right: &ast.BinaryExpressionNode{
								Position: P(5, 6, 1, 6),
								Left: &ast.IntLiteralNode{
									Position: P(5, 1, 1, 6),
									Token:    V(lexer.DecIntToken, "2", 5, 1, 1, 6),
								},
								Op: T(lexer.PowerToken, 7, 2, 1, 8),
								Right: &ast.IntLiteralNode{
									Position: P(10, 1, 1, 11),
									Token:    V(lexer.DecIntToken, "3", 10, 1, 1, 11),
								},
							},
						},
					},
				},
			},
		},
		"exponentiation can have newlines after the operator": {
			input: "1 **\n2 **\n3",
			want: &ast.ProgramNode{
				Position: P(0, 11, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 11, 1, 1),
						Expression: &ast.BinaryExpressionNode{
							Position: P(0, 11, 1, 1),
							Left: &ast.IntLiteralNode{
								Position: P(0, 1, 1, 1),
								Token:    V(lexer.DecIntToken, "1", 0, 1, 1, 1),
							},
							Op: T(lexer.PowerToken, 2, 2, 1, 3),
							Right: &ast.BinaryExpressionNode{
								Position: P(5, 6, 2, 1),
								Left: &ast.IntLiteralNode{
									Position: P(5, 1, 2, 1),
									Token:    V(lexer.DecIntToken, "2", 5, 1, 2, 1),
								},
								Op: T(lexer.PowerToken, 7, 2, 2, 3),
								Right: &ast.IntLiteralNode{
									Position: P(10, 1, 3, 1),
									Token:    V(lexer.DecIntToken, "3", 10, 1, 3, 1),
								},
							},
						},
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
