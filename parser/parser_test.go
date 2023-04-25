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

// Create a string literal node.
func Str(pos lexer.Position, content ...ast.StringLiteralContentNode) *ast.StringLiteralNode {
	return &ast.StringLiteralNode{
		Position: pos,
		Content:  content,
	}
}

// Create a string literal interpolation node.
func StrInterp(pos lexer.Position, expr ast.ExpressionNode) *ast.StringInterpolationNode {
	return &ast.StringInterpolationNode{
		Position:   pos,
		Expression: expr,
	}
}

// Create a string literal content section node.
func StrCont(value string, pos lexer.Position) *ast.StringLiteralContentSectionNode {
	return &ast.StringLiteralContentSectionNode{
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

// Create a false literal node.
func False(pos lexer.Position) *ast.FalseLiteralNode {
	return &ast.FalseLiteralNode{
		Position: pos,
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
				&Error{Pos(2, 1, 2, 1), "unexpected *, expected an expression"},
				&Error{Pos(6, 1, 3, 1), "unexpected *, expected an expression"},
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
				&Error{Pos(2, 1, 2, 1), "unexpected /, expected an expression"},
				&Error{Pos(6, 1, 3, 1), "unexpected /, expected an expression"},
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
				&Error{Pos(2, 2, 2, 1), "unexpected **, expected an expression"},
				&Error{Pos(7, 2, 3, 1), "unexpected **, expected an expression"},
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
		"have higher precedence than additive and multiplicative expression": {
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
		"can be nested": {
			input: "foo = bar = baz = 3",
			want: Prog(
				Pos(0, 19, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 19, 1, 1),
						Asgmt(
							Pos(0, 19, 1, 1),
							Tok(lexer.EqualToken, 4, 1, 1, 5),
							Ident("foo", Pos(0, 3, 1, 1)),
							Asgmt(
								Pos(6, 13, 1, 7),
								Tok(lexer.EqualToken, 10, 1, 1, 11),
								Ident("bar", Pos(6, 3, 1, 7)),
								Asgmt(
									Pos(12, 7, 1, 13),
									Tok(lexer.EqualToken, 16, 1, 1, 17),
									Ident("baz", Pos(12, 3, 1, 13)),
									Int(lexer.DecIntToken, "3", 18, 1, 1, 19),
								),
							),
						),
					),
				},
			),
		},
		"can have newlines after the operator": {
			input: "foo =\nbar =\nbaz =\n3",
			want: Prog(
				Pos(0, 19, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 19, 1, 1),
						Asgmt(
							Pos(0, 19, 1, 1),
							Tok(lexer.EqualToken, 4, 1, 1, 5),
							Ident("foo", Pos(0, 3, 1, 1)),
							Asgmt(
								Pos(6, 13, 2, 1),
								Tok(lexer.EqualToken, 10, 1, 2, 5),
								Ident("bar", Pos(6, 3, 2, 1)),
								Asgmt(
									Pos(12, 7, 3, 1),
									Tok(lexer.EqualToken, 16, 1, 3, 5),
									Ident("baz", Pos(12, 3, 3, 1)),
									Int(lexer.DecIntToken, "3", 18, 1, 4, 1),
								),
							),
						),
					),
				},
			),
		},
		"can't have newlines before the operator": {
			input: "foo\n= bar\n= baz\n= 3",
			want: Prog(
				Pos(0, 19, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 4, 1, 1),
						Ident("foo", Pos(0, 3, 1, 1)),
					),
					ExprStmt(
						Pos(4, 1, 2, 1),
						Invalid(Tok(lexer.EqualToken, 4, 1, 2, 1)),
					),
					ExprStmt(
						Pos(10, 1, 3, 1),
						Invalid(Tok(lexer.EqualToken, 10, 1, 3, 1)),
					),
					ExprStmt(
						Pos(16, 1, 4, 1),
						Invalid(Tok(lexer.EqualToken, 16, 1, 4, 1)),
					),
				},
			),
			err: ErrorList{
				&Error{Pos(4, 1, 2, 1), "unexpected =, expected an expression"},
				&Error{Pos(10, 1, 3, 1), "unexpected =, expected an expression"},
				&Error{Pos(16, 1, 4, 1), "unexpected =, expected an expression"},
			},
		},
		"has lower precedence than other expressions": {
			input: "f = some && awesome || thing + 2 * 8 > 5 == false",
			want: Prog(
				Pos(0, 49, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 49, 1, 1),
						Asgmt(
							Pos(0, 49, 1, 1),
							Tok(lexer.EqualToken, 2, 1, 1, 3),
							Ident("f", Pos(0, 1, 1, 1)),
							Bin(
								Pos(4, 45, 1, 5),
								Tok(lexer.OrOrToken, 20, 2, 1, 21),
								Bin(
									Pos(4, 15, 1, 5),
									Tok(lexer.AndAndToken, 9, 2, 1, 10),
									Ident("some", Pos(4, 4, 1, 5)),
									Ident("awesome", Pos(12, 7, 1, 13)),
								),
								Bin(
									Pos(23, 26, 1, 24),
									Tok(lexer.EqualEqualToken, 41, 2, 1, 42),
									Bin(
										Pos(23, 17, 1, 24),
										Tok(lexer.GreaterToken, 37, 1, 1, 38),
										Bin(
											Pos(23, 13, 1, 24),
											Tok(lexer.PlusToken, 29, 1, 1, 30),
											Ident("thing", Pos(23, 5, 1, 24)),
											Bin(
												Pos(31, 5, 1, 32),
												Tok(lexer.StarToken, 33, 1, 1, 34),
												Int(lexer.DecIntToken, "2", 31, 1, 1, 32),
												Int(lexer.DecIntToken, "8", 35, 1, 1, 36),
											),
										),
										Int(lexer.DecIntToken, "5", 39, 1, 1, 40),
									),
									False(Pos(44, 5, 1, 45)),
								),
							),
						),
					),
				},
			),
		},
		"has many versions": {
			input: "a = b -= c += d *= e /= f **= g ~= h &&= i &= j ||= k |= l ^= m ??= n <<= o >>= p %= q",
			want: Prog(
				Pos(0, 86, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 86, 1, 1),
						Asgmt(
							Pos(0, 86, 1, 1),
							Tok(lexer.EqualToken, 2, 1, 1, 3),
							Ident("a", Pos(0, 1, 1, 1)),
							Asgmt(
								Pos(4, 82, 1, 5),
								Tok(lexer.MinusEqualToken, 6, 2, 1, 7),
								Ident("b", Pos(4, 1, 1, 5)),
								Asgmt(
									Pos(9, 77, 1, 10),
									Tok(lexer.PlusEqualToken, 11, 2, 1, 12),
									Ident("c", Pos(9, 1, 1, 10)),
									Asgmt(
										Pos(14, 72, 1, 15),
										Tok(lexer.StarEqualToken, 16, 2, 1, 17),
										Ident("d", Pos(14, 1, 1, 15)),
										Asgmt(
											Pos(19, 67, 1, 20),
											Tok(lexer.SlashEqualToken, 21, 2, 1, 22),
											Ident("e", Pos(19, 1, 1, 20)),
											Asgmt(
												Pos(24, 62, 1, 25),
												Tok(lexer.StarStarEqualToken, 26, 3, 1, 27),
												Ident("f", Pos(24, 1, 1, 25)),
												Asgmt(
													Pos(30, 56, 1, 31),
													Tok(lexer.TildeEqualToken, 32, 2, 1, 33),
													Ident("g", Pos(30, 1, 1, 31)),
													Asgmt(
														Pos(35, 51, 1, 36),
														Tok(lexer.AndAndEqualToken, 37, 3, 1, 38),
														Ident("h", Pos(35, 1, 1, 36)),
														Asgmt(
															Pos(41, 45, 1, 42),
															Tok(lexer.AndEqualToken, 43, 2, 1, 44),
															Ident("i", Pos(41, 1, 1, 42)),
															Asgmt(
																Pos(46, 40, 1, 47),
																Tok(lexer.OrOrEqualToken, 48, 3, 1, 49),
																Ident("j", Pos(46, 1, 1, 47)),
																Asgmt(
																	Pos(52, 34, 1, 53),
																	Tok(lexer.OrEqualToken, 54, 2, 1, 55),
																	Ident("k", Pos(52, 1, 1, 53)),
																	Asgmt(
																		Pos(57, 29, 1, 58),
																		Tok(lexer.XorEqualToken, 59, 2, 1, 60),
																		Ident("l", Pos(57, 1, 1, 58)),
																		Asgmt(
																			Pos(62, 24, 1, 63),
																			Tok(lexer.QuestionQuestionEqualToken, 64, 3, 1, 65),
																			Ident("m", Pos(62, 1, 1, 63)),
																			Asgmt(
																				Pos(68, 18, 1, 69),
																				Tok(lexer.LBitShiftEqualToken, 70, 3, 1, 71),
																				Ident("n", Pos(68, 1, 1, 69)),
																				Asgmt(
																					Pos(74, 12, 1, 75),
																					Tok(lexer.RBitShiftEqualToken, 76, 3, 1, 77),
																					Ident("o", Pos(74, 1, 1, 75)),
																					Asgmt(
																						Pos(80, 6, 1, 81),
																						Tok(lexer.PercentEqualToken, 82, 2, 1, 83),
																						Ident("p", Pos(80, 1, 1, 81)),
																						Ident("q", Pos(85, 1, 1, 86)),
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
		"or has lower precedence than and": {
			input: "foo || bar && baz",
			want: Prog(
				Pos(0, 17, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 17, 1, 1),
						Bin(
							Pos(0, 17, 1, 1),
							Tok(lexer.OrOrToken, 4, 2, 1, 5),
							Ident("foo", Pos(0, 3, 1, 1)),
							Bin(
								Pos(7, 10, 1, 8),
								Tok(lexer.AndAndToken, 11, 2, 1, 12),
								Ident("bar", Pos(7, 3, 1, 8)),
								Ident("baz", Pos(14, 3, 1, 15)),
							),
						),
					),
				},
			),
		},
		"nil coalescing operator has lower precedence than and": {
			input: "foo ?? bar && baz",
			want: Prog(
				Pos(0, 17, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 17, 1, 1),
						Bin(
							Pos(0, 17, 1, 1),
							Tok(lexer.QuestionQuestionToken, 4, 2, 1, 5),
							Ident("foo", Pos(0, 3, 1, 1)),
							Bin(
								Pos(7, 10, 1, 8),
								Tok(lexer.AndAndToken, 11, 2, 1, 12),
								Ident("bar", Pos(7, 3, 1, 8)),
								Ident("baz", Pos(14, 3, 1, 15)),
							),
						),
					),
				},
			),
		},
		"nil coalescing operator has the same precedence as or": {
			input: "foo ?? bar || baz ?? boo",
			want: Prog(
				Pos(0, 24, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 24, 1, 1),
						Bin(
							Pos(0, 24, 1, 1),
							Tok(lexer.QuestionQuestionToken, 18, 2, 1, 19),
							Bin(
								Pos(0, 17, 1, 1),
								Tok(lexer.OrOrToken, 11, 2, 1, 12),
								Bin(
									Pos(0, 10, 1, 1),
									Tok(lexer.QuestionQuestionToken, 4, 2, 1, 5),
									Ident("foo", Pos(0, 3, 1, 1)),
									Ident("bar", Pos(7, 3, 1, 8)),
								),
								Ident("baz", Pos(14, 3, 1, 15)),
							),
							Ident("boo", Pos(21, 3, 1, 22)),
						),
					),
				},
			),
		},
		"or is evaluated from left to right": {
			input: "foo || bar || baz",
			want: Prog(
				Pos(0, 17, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 17, 1, 1),
						Bin(
							Pos(0, 17, 1, 1),
							Tok(lexer.OrOrToken, 11, 2, 1, 12),
							Bin(
								Pos(0, 10, 1, 1),
								Tok(lexer.OrOrToken, 4, 2, 1, 5),
								Ident("foo", Pos(0, 3, 1, 1)),
								Ident("bar", Pos(7, 3, 1, 8)),
							),
							Ident("baz", Pos(14, 3, 1, 15)),
						),
					),
				},
			),
		},
		"and is evaluated from left to right": {
			input: "foo && bar && baz",
			want: Prog(
				Pos(0, 17, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 17, 1, 1),
						Bin(
							Pos(0, 17, 1, 1),
							Tok(lexer.AndAndToken, 11, 2, 1, 12),
							Bin(
								Pos(0, 10, 1, 1),
								Tok(lexer.AndAndToken, 4, 2, 1, 5),
								Ident("foo", Pos(0, 3, 1, 1)),
								Ident("bar", Pos(7, 3, 1, 8)),
							),
							Ident("baz", Pos(14, 3, 1, 15)),
						),
					),
				},
			),
		},
		"nil coalescing operator is evaluated from left to right": {
			input: "foo ?? bar ?? baz",
			want: Prog(
				Pos(0, 17, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 17, 1, 1),
						Bin(
							Pos(0, 17, 1, 1),
							Tok(lexer.QuestionQuestionToken, 11, 2, 1, 12),
							Bin(
								Pos(0, 10, 1, 1),
								Tok(lexer.QuestionQuestionToken, 4, 2, 1, 5),
								Ident("foo", Pos(0, 3, 1, 1)),
								Ident("bar", Pos(7, 3, 1, 8)),
							),
							Ident("baz", Pos(14, 3, 1, 15)),
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
			want: Prog(
				Pos(0, 36, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 36, 1, 1),
						Str(
							Pos(0, 36, 1, 1),
							StrCont("foo\nbar\rbaz\\car\t\b\"\v\f\x12\a", Pos(1, 34, 1, 2)),
						),
					),
				},
			),
		},
		"reports errors for invalid hex escapes": {
			input: `"foo \xgh bar"`,
			want: Prog(
				Pos(0, 14, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 14, 1, 1),
						Str(
							Pos(0, 14, 1, 1),
							StrCont("foo ", Pos(1, 4, 1, 2)),
							Invalid(VTok(lexer.ErrorToken, "invalid hex escape in string literal", 5, 4, 1, 6)),
							StrCont(" bar", Pos(9, 4, 1, 10)),
						),
					),
				},
			),
			err: ErrorList{
				&Error{Pos(5, 4, 1, 6), "invalid hex escape in string literal"},
			},
		},
		"reports errors for nonexistent escape sequences": {
			input: `"foo \q bar"`,
			want: Prog(
				Pos(0, 12, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 12, 1, 1),
						Str(
							Pos(0, 12, 1, 1),
							StrCont("foo ", Pos(1, 4, 1, 2)),
							Invalid(VTok(lexer.ErrorToken, "invalid escape sequence `\\q` in string literal", 5, 2, 1, 6)),
							StrCont(" bar", Pos(7, 4, 1, 8)),
						),
					),
				},
			),
			err: ErrorList{
				&Error{Pos(5, 2, 1, 6), "invalid escape sequence `\\q` in string literal"},
			},
		},
		"can contain interpolated expressions": {
			input: `"foo ${bar + 2} baz ${fudge}"`,
			want: Prog(
				Pos(0, 29, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 29, 1, 1),
						Str(
							Pos(0, 29, 1, 1),
							StrCont("foo ", Pos(1, 4, 1, 2)),
							StrInterp(
								Pos(5, 10, 1, 6),
								Bin(
									Pos(7, 7, 1, 8),
									Tok(lexer.PlusToken, 11, 1, 1, 12),
									Ident("bar", Pos(7, 3, 1, 8)),
									Int(lexer.DecIntToken, "2", 13, 1, 1, 14),
								),
							),
							StrCont(" baz ", Pos(15, 5, 1, 16)),
							StrInterp(
								Pos(20, 8, 1, 21),
								Ident("fudge", Pos(22, 5, 1, 23)),
							),
						),
					),
				},
			),
		},
		"can't contain string literals inside interpolation": {
			input: `"foo ${"bar" + 2} baza"`,
			want: Prog(
				Pos(0, 23, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 23, 1, 1),
						Str(
							Pos(0, 23, 1, 1),
							StrCont("foo ", Pos(1, 4, 1, 2)),
							StrInterp(
								Pos(5, 12, 1, 6),
								Bin(
									Pos(7, 9, 1, 8),
									Tok(lexer.PlusToken, 13, 1, 1, 14),
									Invalid(VTok(lexer.ErrorToken, "unexpected string literal in string interpolation, only raw strings delimited with `'` can be used in string interpolation", 7, 5, 1, 8)),
									Int(lexer.DecIntToken, "2", 15, 1, 1, 16),
								),
							),
							StrCont(" baza", Pos(17, 5, 1, 18)),
						),
					),
				},
			),
			err: ErrorList{
				&Error{Message: "unexpected string literal in string interpolation, only raw strings delimited with `'` can be used in string interpolation", Position: Pos(7, 5, 1, 8)},
			},
		},
		"can contain raw string literals inside interpolation": {
			input: `"foo ${'bar' + 2} baza"`,
			want: Prog(
				Pos(0, 23, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 23, 1, 1),
						Str(
							Pos(0, 23, 1, 1),
							StrCont("foo ", Pos(1, 4, 1, 2)),
							StrInterp(
								Pos(5, 12, 1, 6),
								Bin(
									Pos(7, 9, 1, 8),
									Tok(lexer.PlusToken, 13, 1, 1, 14),
									RawStr("bar", Pos(7, 5, 1, 8)),
									Int(lexer.DecIntToken, "2", 15, 1, 1, 16),
								),
							),
							StrCont(" baza", Pos(17, 5, 1, 18)),
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
			want: Prog(
				Pos(0, 36, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 36, 1, 1),
						RawStr(`foo\nbar\rbaz\\car\t\b\"\v\f\x12\a`, Pos(0, 36, 1, 1)),
					),
				},
			),
		},
		"can't contain interpolated expressions": {
			input: `'foo ${bar + 2} baz ${fudge}'`,
			want: Prog(
				Pos(0, 29, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 29, 1, 1),
						RawStr(`foo ${bar + 2} baz ${fudge}`, Pos(0, 29, 1, 1)),
					),
				},
			),
		},
		"can contain double quotes": {
			input: `'foo ${"bar" + 2} baza'`,
			want: Prog(
				Pos(0, 23, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 23, 1, 1),
						RawStr(`foo ${"bar" + 2} baza`, Pos(0, 23, 1, 1)),
					),
				},
			),
		},
		"doesn't allow escaping single quotes": {
			input: `'foo\'s house'`,
			want: Prog(
				Pos(0, 14, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 6, 1, 1),
						RawStr("foo\\", Pos(0, 6, 1, 1)),
					),
					ExprStmt(
						Pos(6, 1, 1, 7),
						Ident("s", Pos(6, 1, 1, 7)),
					),
					ExprStmt(
						Pos(8, 5, 1, 9),
						Ident("house", Pos(8, 5, 1, 9)),
					),
					ExprStmt(
						Pos(13, 1, 1, 14),
						Invalid(VTok(lexer.ErrorToken, "unterminated raw string literal, missing `'`", 13, 1, 1, 14)),
					),
				},
			),
			err: ErrorList{
				&Error{Message: "unterminated raw string literal, missing `'`", Position: Pos(13, 1, 1, 14)},
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
