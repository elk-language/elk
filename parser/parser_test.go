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
func Pos(startByte, byteLength, line, column int) lexer.Position {
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
		Position: Pos(startByte, byteLength, line, column),
		Token:    VTok(tokenType, value, startByte, byteLength, line, column),
	}
}

// Create a new float node.
func Float(value string, startByte, byteLength, line, column int) *ast.FloatLiteralNode {
	return &ast.FloatLiteralNode{
		Position: Pos(startByte, byteLength, line, column),
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

// Create a program node.
func Prog(pos lexer.Position, body []ast.StatementNode) *ast.ProgramNode {
	return &ast.ProgramNode{
		Position: pos,
		Body:     body,
	}
}

// Create an expression statement node.
func ExprStmt(pos lexer.Position, expr ast.ExpressionNode) *ast.ExpressionStatementNode {
	return &ast.ExpressionStatementNode{
		Position:   pos,
		Expression: expr,
	}
}

// Create an assignment expression node.
func Asgmt(pos lexer.Position, op *lexer.Token, left ast.ExpressionNode, right ast.ExpressionNode) *ast.AssignmentExpressionNode {
	return &ast.AssignmentExpressionNode{
		Position: pos,
		Left:     left,
		Op:       op,
		Right:    right,
	}
}

// Create a raw string literal node.
func RawStr(value string, pos lexer.Position) *ast.RawStringLiteralNode {
	return &ast.RawStringLiteralNode{
		Position: pos,
		Value:    value,
	}
}

// Create a raw string literal node.
func Ident(value string, pos lexer.Position) *ast.IdentifierNode {
	return &ast.IdentifierNode{
		Position: pos,
		Value:    value,
	}
}

// Create a raw string literal node.
func PrivIdent(value string, pos lexer.Position) *ast.PrivateIdentifierNode {
	return &ast.PrivateIdentifierNode{
		Position: pos,
		Value:    value,
	}
}

// Create a raw string literal node.
func Const(value string, pos lexer.Position) *ast.ConstantNode {
	return &ast.ConstantNode{
		Position: pos,
		Value:    value,
	}
}

// Create a raw string literal node.
func PrivConst(value string, pos lexer.Position) *ast.PrivateConstantNode {
	return &ast.PrivateConstantNode{
		Position: pos,
		Value:    value,
	}
}

// Create a new token in tests.
var Tok = lexer.NewToken

// Create a new token with value in tests.
var VTok = lexer.NewTokenWithValue

// Slice of statement nodes.
type Stmts = []ast.StatementNode

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
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 9, 1, 1),
						Bin(
							Pos(0, 9, 1, 1),
							Tok(lexer.PlusToken, 6, 1, 1, 7),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(lexer.PlusToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 1, 5),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 1, 9),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 +\n2 +\n3",
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 9, 1, 1),
						Bin(
							Pos(0, 9, 1, 1),
							Tok(lexer.PlusToken, 6, 1, 2, 3),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(lexer.PlusToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 2, 1),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 3, 1),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n+ 2\n+ 3",
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 2, 1, 1),
						Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
					),
					ExprStmt(
						Pos(2, 4, 2, 1),
						Unary(
							Pos(2, 3, 2, 1),
							Tok(lexer.PlusToken, 2, 1, 2, 1),
							Int(lexer.DecIntToken, "2", 4, 1, 2, 3),
						),
					),
					ExprStmt(
						Pos(6, 3, 3, 1),
						Unary(
							Pos(6, 3, 3, 1),
							Tok(lexer.PlusToken, 6, 1, 3, 1),
							Int(lexer.DecIntToken, "3", 8, 1, 3, 3),
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
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 9, 1, 1),
						Bin(
							Pos(0, 9, 1, 1),
							Tok(lexer.MinusToken, 6, 1, 1, 7),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(lexer.MinusToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 1, 5),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 1, 9),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 -\n2 -\n3",
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 9, 1, 1),
						Bin(
							Pos(0, 9, 1, 1),
							Tok(lexer.MinusToken, 6, 1, 2, 3),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(lexer.MinusToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 2, 1),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 3, 1),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n- 2\n- 3",
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 2, 1, 1),
						Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
					),
					ExprStmt(
						Pos(2, 4, 2, 1),
						Unary(
							Pos(2, 3, 2, 1),
							Tok(lexer.MinusToken, 2, 1, 2, 1),
							Int(lexer.DecIntToken, "2", 4, 1, 2, 3),
						),
					),
					ExprStmt(
						Pos(6, 3, 3, 1),
						Unary(
							Pos(6, 3, 3, 1),
							Tok(lexer.MinusToken, 6, 1, 3, 1),
							Int(lexer.DecIntToken, "3", 8, 1, 3, 3),
						),
					),
				},
			),
		},
		"has the same precedence as addition": {
			input: "1 + 2 - 3",
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 9, 1, 1),
						Bin(
							Pos(0, 9, 1, 1),
							Tok(lexer.MinusToken, 6, 1, 1, 7),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(lexer.PlusToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 1, 5),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 1, 9),
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
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 9, 1, 1),
						Bin(
							Pos(0, 9, 1, 1),
							Tok(lexer.StarToken, 6, 1, 1, 7),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(lexer.StarToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 1, 5),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 1, 9),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 *\n2 *\n3",
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 9, 1, 1),
						Bin(
							Pos(0, 9, 1, 1),
							Tok(lexer.StarToken, 6, 1, 2, 3),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(lexer.StarToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 2, 1),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 3, 1),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n* 2\n* 3",
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 2, 1, 1),
						Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
					),
					ExprStmt(
						Pos(2, 1, 2, 1),
						Invalid(Tok(lexer.StarToken, 2, 1, 2, 1)),
					),
					ExprStmt(
						Pos(6, 1, 3, 1),
						Invalid(Tok(lexer.StarToken, 6, 1, 3, 1)),
					),
				},
			),
			err: ErrorList{
				&Error{Pos(2, 1, 2, 1), "Unexpected *, expected an expression"},
				&Error{Pos(6, 1, 3, 1), "Unexpected *, expected an expression"},
			},
		},
		"has higher precedence than addition": {
			input: "1 + 2 * 3",
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 9, 1, 1),
						Bin(
							Pos(0, 9, 1, 1),
							Tok(lexer.PlusToken, 2, 1, 1, 3),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Bin(
								Pos(4, 5, 1, 5),
								Tok(lexer.StarToken, 6, 1, 1, 7),
								Int(lexer.DecIntToken, "2", 4, 1, 1, 5),
								Int(lexer.DecIntToken, "3", 8, 1, 1, 9),
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
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 9, 1, 1),
						Bin(
							Pos(0, 9, 1, 1),
							Tok(lexer.SlashToken, 6, 1, 1, 7),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(lexer.SlashToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 1, 5),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 1, 9),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 /\n2 /\n3",
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 9, 1, 1),
						Bin(
							Pos(0, 9, 1, 1),
							Tok(lexer.SlashToken, 6, 1, 2, 3),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(lexer.SlashToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 2, 1),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 3, 1),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n/ 2\n/ 3",
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 2, 1, 1),
						Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
					),
					ExprStmt(
						Pos(2, 1, 2, 1),
						Invalid(Tok(lexer.SlashToken, 2, 1, 2, 1)),
					),
					ExprStmt(
						Pos(6, 1, 3, 1),
						Invalid(Tok(lexer.SlashToken, 6, 1, 3, 1)),
					),
				},
			),
			err: ErrorList{
				&Error{Pos(2, 1, 2, 1), "Unexpected /, expected an expression"},
				&Error{Pos(6, 1, 3, 1), "Unexpected /, expected an expression"},
			},
		},
		"has the same precedence as multiplication": {
			input: "1 * 2 / 3",
			want: Prog(
				Pos(0, 9, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 9, 1, 1),
						Bin(
							Pos(0, 9, 1, 1),
							Tok(lexer.SlashToken, 6, 1, 1, 7),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(lexer.StarToken, 2, 1, 1, 3),
								Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
								Int(lexer.DecIntToken, "2", 4, 1, 1, 5),
							),
							Int(lexer.DecIntToken, "3", 8, 1, 1, 9),
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
			want: Prog(
				Pos(0, 11, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 11, 1, 1),
						Bin(
							Pos(0, 11, 1, 1),
							Tok(lexer.StarStarToken, 2, 2, 1, 3),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Bin(
								Pos(5, 6, 1, 6),
								Tok(lexer.StarStarToken, 7, 2, 1, 8),
								Int(lexer.DecIntToken, "2", 5, 1, 1, 6),
								Int(lexer.DecIntToken, "3", 10, 1, 1, 11),
							),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "1 **\n2 **\n3",
			want: Prog(
				Pos(0, 11, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 11, 1, 1),
						Bin(
							Pos(0, 11, 1, 1),
							Tok(lexer.StarStarToken, 2, 2, 1, 3),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Bin(
								Pos(5, 6, 2, 1),
								Tok(lexer.StarStarToken, 7, 2, 2, 3),
								Int(lexer.DecIntToken, "2", 5, 1, 2, 1),
								Int(lexer.DecIntToken, "3", 10, 1, 3, 1),
							),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "1\n** 2\n** 3",
			want: Prog(
				Pos(0, 11, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 2, 1, 1),
						Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
					),
					ExprStmt(
						Pos(2, 2, 2, 1),
						Invalid(Tok(lexer.StarStarToken, 2, 2, 2, 1)),
					),
					ExprStmt(
						Pos(7, 2, 3, 1),
						Invalid(Tok(lexer.StarStarToken, 7, 2, 3, 1)),
					),
				},
			),
			err: ErrorList{
				&Error{Pos(2, 2, 2, 1), "Unexpected **, expected an expression"},
				&Error{Pos(7, 2, 3, 1), "Unexpected **, expected an expression"},
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
			want: Prog(
				Pos(0, 6, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 6, 1, 1),
						Unary(
							Pos(0, 6, 1, 1),
							Tok(lexer.PlusToken, 0, 1, 1, 1),
							Unary(
								Pos(1, 5, 1, 2),
								Tok(lexer.PlusToken, 1, 1, 1, 2),
								Unary(
									Pos(2, 4, 1, 3),
									Tok(lexer.PlusToken, 2, 1, 1, 3),
									Float("1.5", 3, 3, 1, 4),
								),
							),
						),
					),
				},
			),
		},
		"minus can be nested": {
			input: "---1.5",
			want: Prog(
				Pos(0, 6, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 6, 1, 1),
						Unary(
							Pos(0, 6, 1, 1),
							Tok(lexer.MinusToken, 0, 1, 1, 1),
							Unary(
								Pos(1, 5, 1, 2),
								Tok(lexer.MinusToken, 1, 1, 1, 2),
								Unary(
									Pos(2, 4, 1, 3),
									Tok(lexer.MinusToken, 2, 1, 1, 3),
									Float("1.5", 3, 3, 1, 4),
								),
							),
						),
					),
				},
			),
		},
		"logical not can be nested": {
			input: "!!!1.5",
			want: Prog(
				Pos(0, 6, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 6, 1, 1),
						Unary(
							Pos(0, 6, 1, 1),
							Tok(lexer.BangToken, 0, 1, 1, 1),
							Unary(
								Pos(1, 5, 1, 2),
								Tok(lexer.BangToken, 1, 1, 1, 2),
								Unary(
									Pos(2, 4, 1, 3),
									Tok(lexer.BangToken, 2, 1, 1, 3),
									Float("1.5", 3, 3, 1, 4),
								),
							),
						),
					),
				},
			),
		},
		"bitwise not can be nested": {
			input: "~~~1.5",
			want: Prog(
				Pos(0, 6, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 6, 1, 1),
						Unary(
							Pos(0, 6, 1, 1),
							Tok(lexer.TildeToken, 0, 1, 1, 1),
							Unary(
								Pos(1, 5, 1, 2),
								Tok(lexer.TildeToken, 1, 1, 1, 2),
								Unary(
									Pos(2, 4, 1, 3),
									Tok(lexer.TildeToken, 2, 1, 1, 3),
									Float("1.5", 3, 3, 1, 4),
								),
							),
						),
					),
				},
			),
		},
		"all have the same precedence": {
			input: "!+~1.5",
			want: Prog(
				Pos(0, 6, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 6, 1, 1),
						Unary(
							Pos(0, 6, 1, 1),
							Tok(lexer.BangToken, 0, 1, 1, 1),
							Unary(
								Pos(1, 5, 1, 2),
								Tok(lexer.PlusToken, 1, 1, 1, 2),
								Unary(
									Pos(2, 4, 1, 3),
									Tok(lexer.TildeToken, 2, 1, 1, 3),
									Float("1.5", 3, 3, 1, 4),
								),
							),
						),
					),
				},
			),
		},
		"have higher precedence than additive an multiplicative expression": {
			input: "!!1.5 * 2 + ~.5",
			want: Prog(
				Pos(0, 15, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 15, 1, 1),
						Bin(
							Pos(0, 15, 1, 1),
							Tok(lexer.PlusToken, 10, 1, 1, 11),
							Bin(
								Pos(0, 9, 1, 1),
								Tok(lexer.StarToken, 6, 1, 1, 7),
								Unary(
									Pos(0, 5, 1, 1),
									Tok(lexer.BangToken, 0, 1, 1, 1),
									Unary(
										Pos(1, 4, 1, 2),
										Tok(lexer.BangToken, 1, 1, 1, 2),
										Float("1.5", 2, 3, 1, 3),
									),
								),
								Int(lexer.DecIntToken, "2", 8, 1, 1, 9),
							),
							Unary(
								Pos(12, 3, 1, 13),
								Tok(lexer.TildeToken, 12, 1, 1, 13),
								Float("0.5", 13, 2, 1, 14),
							),
						),
					),
				},
			),
		},
		"have lower precedence than exponentiation": {
			input: "-2 ** 3",
			want: Prog(
				Pos(0, 7, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 7, 1, 1),
						Unary(
							Pos(0, 7, 1, 1),
							Tok(lexer.MinusToken, 0, 1, 1, 1),
							Bin(
								Pos(1, 6, 1, 2),
								Tok(lexer.StarStarToken, 3, 2, 1, 4),
								Int(lexer.DecIntToken, "2", 1, 1, 1, 2),
								Int(lexer.DecIntToken, "3", 6, 1, 1, 7),
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
			want: Prog(
				Pos(0, 13, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 7, 1, 1),
						Bin(
							Pos(0, 6, 1, 1),
							Tok(lexer.StarStarToken, 2, 2, 1, 3),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Int(lexer.DecIntToken, "2", 5, 1, 1, 6),
						),
					),
					ExprStmt(
						Pos(8, 5, 1, 9),
						Bin(
							Pos(8, 5, 1, 9),
							Tok(lexer.StarToken, 10, 1, 1, 11),
							Int(lexer.DecIntToken, "5", 8, 1, 1, 9),
							Int(lexer.DecIntToken, "8", 12, 1, 1, 13),
						),
					),
				},
			),
		},
		"endlines can separate statements": {
			input: "1 ** 2\n5 * 8",
			want: Prog(
				Pos(0, 12, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 7, 1, 1),
						Bin(
							Pos(0, 6, 1, 1),
							Tok(lexer.StarStarToken, 2, 2, 1, 3),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Int(lexer.DecIntToken, "2", 5, 1, 1, 6),
						),
					),
					ExprStmt(
						Pos(7, 5, 2, 1),
						Bin(
							Pos(7, 5, 2, 1),
							Tok(lexer.StarToken, 9, 1, 2, 3),
							Int(lexer.DecIntToken, "5", 7, 1, 2, 1),
							Int(lexer.DecIntToken, "8", 11, 1, 2, 5),
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

func TestAssignment(t *testing.T) {
	tests := testTable{
		"ints are not valid assignment targets": {
			input: "1 -= 2",
			want: Prog(
				Pos(0, 6, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 6, 1, 1),
						Asgmt(
							Pos(0, 6, 1, 1),
							Tok(lexer.MinusEqualToken, 2, 2, 1, 3),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Int(lexer.DecIntToken, "2", 5, 1, 1, 6),
						),
					),
				},
			),
			err: ErrorList{
				&Error{Position: Pos(0, 1, 1, 1), Message: "invalid `-=` assignment target"},
			},
		},
		"strings are not valid assignment targets": {
			input: "'foo' -= 2",
			want: Prog(
				Pos(0, 10, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 10, 1, 1),
						Asgmt(
							Pos(0, 10, 1, 1),
							Tok(lexer.MinusEqualToken, 6, 2, 1, 7),
							RawStr("foo", Pos(0, 5, 1, 1)),
							Int(lexer.DecIntToken, "2", 9, 1, 1, 10),
						),
					),
				},
			),
			err: ErrorList{
				&Error{Position: Pos(0, 5, 1, 1), Message: "invalid `-=` assignment target"},
			},
		},
		"constants are not valid assignment targets": {
			input: "FooBa -= 2",
			want: Prog(
				Pos(0, 10, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 10, 1, 1),
						Asgmt(
							Pos(0, 10, 1, 1),
							Tok(lexer.MinusEqualToken, 6, 2, 1, 7),
							Const("FooBa", Pos(0, 5, 1, 1)),
							Int(lexer.DecIntToken, "2", 9, 1, 1, 10),
						),
					),
				},
			),
			err: ErrorList{
				&Error{Position: Pos(0, 5, 1, 1), Message: "constants can't be assigned, maybe you meant to declare it with `:=`"},
			},
		},
		"private constants are not valid assignment targets": {
			input: "_FooB -= 2",
			want: Prog(
				Pos(0, 10, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 10, 1, 1),
						Asgmt(
							Pos(0, 10, 1, 1),
							Tok(lexer.MinusEqualToken, 6, 2, 1, 7),
							PrivConst("_FooB", Pos(0, 5, 1, 1)),
							Int(lexer.DecIntToken, "2", 9, 1, 1, 10),
						),
					),
				},
			),
			err: ErrorList{
				&Error{Position: Pos(0, 5, 1, 1), Message: "constants can't be assigned, maybe you meant to declare it with `:=`"},
			},
		},
		"identifiers can be assigned": {
			input: "foo -= 2",
			want: Prog(
				Pos(0, 8, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 8, 1, 1),
						Asgmt(
							Pos(0, 8, 1, 1),
							Tok(lexer.MinusEqualToken, 4, 2, 1, 5),
							Ident("foo", Pos(0, 3, 1, 1)),
							Int(lexer.DecIntToken, "2", 7, 1, 1, 8),
						),
					),
				},
			),
		},
		"private identifiers can be assigned": {
			input: "_fo -= 2",
			want: Prog(
				Pos(0, 8, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 8, 1, 1),
						Asgmt(
							Pos(0, 8, 1, 1),
							Tok(lexer.MinusEqualToken, 4, 2, 1, 5),
							PrivIdent("_fo", Pos(0, 3, 1, 1)),
							Int(lexer.DecIntToken, "2", 7, 1, 1, 8),
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