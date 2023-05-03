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
func Pos(startByte, byteLength, line, column int) *lexer.Position {
	return &lexer.Position{
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
func Unary(pos *lexer.Position, op *lexer.Token, right ast.ExpressionNode) *ast.UnaryExpressionNode {
	return &ast.UnaryExpressionNode{
		Position: pos,
		Op:       op,
		Right:    right,
	}
}

// Create a binary expression node.
func Bin(pos *lexer.Position, op *lexer.Token, left ast.ExpressionNode, right ast.ExpressionNode) *ast.BinaryExpressionNode {
	return &ast.BinaryExpressionNode{
		Position: pos,
		Left:     left,
		Op:       op,
		Right:    right,
	}
}

// Create a logical expression node.
func Logic(pos *lexer.Position, op *lexer.Token, left ast.ExpressionNode, right ast.ExpressionNode) *ast.LogicalExpressionNode {
	return &ast.LogicalExpressionNode{
		Position: pos,
		Left:     left,
		Op:       op,
		Right:    right,
	}
}

// Create a program node.
func Prog(pos *lexer.Position, body []ast.StatementNode) *ast.ProgramNode {
	return &ast.ProgramNode{
		Position: pos,
		Body:     body,
	}
}

// Create an expression statement node.
func ExprStmt(pos *lexer.Position, expr ast.ExpressionNode) *ast.ExpressionStatementNode {
	return &ast.ExpressionStatementNode{
		Position:   pos,
		Expression: expr,
	}
}

// Create an assignment expression node.
func Asgmt(pos *lexer.Position, op *lexer.Token, left ast.ExpressionNode, right ast.ExpressionNode) *ast.AssignmentExpressionNode {
	return &ast.AssignmentExpressionNode{
		Position: pos,
		Left:     left,
		Op:       op,
		Right:    right,
	}
}

// Create a raw string literal node.
func RawStr(value string, pos *lexer.Position) *ast.RawStringLiteralNode {
	return &ast.RawStringLiteralNode{
		Position: pos,
		Value:    value,
	}
}

// Create a string literal node.
func Str(pos *lexer.Position, content ...ast.StringLiteralContentNode) *ast.StringLiteralNode {
	return &ast.StringLiteralNode{
		Position: pos,
		Content:  content,
	}
}

// Create a string literal interpolation node.
func StrInterp(pos *lexer.Position, expr ast.ExpressionNode) *ast.StringInterpolationNode {
	return &ast.StringInterpolationNode{
		Position:   pos,
		Expression: expr,
	}
}

// Create a string literal content section node.
func StrCont(value string, pos *lexer.Position) *ast.StringLiteralContentSectionNode {
	return &ast.StringLiteralContentSectionNode{
		Position: pos,
		Value:    value,
	}
}

// Create a raw string literal node.
func Ident(value string, pos *lexer.Position) *ast.IdentifierNode {
	return &ast.IdentifierNode{
		Position: pos,
		Value:    value,
	}
}

// Create a raw string literal node.
func PrivIdent(value string, pos *lexer.Position) *ast.PrivateIdentifierNode {
	return &ast.PrivateIdentifierNode{
		Position: pos,
		Value:    value,
	}
}

// Create a raw string literal node.
func Const(value string, pos *lexer.Position) *ast.ConstantNode {
	return &ast.ConstantNode{
		Position: pos,
		Value:    value,
	}
}

// Create a raw string literal node.
func PrivConst(value string, pos *lexer.Position) *ast.PrivateConstantNode {
	return &ast.PrivateConstantNode{
		Position: pos,
		Value:    value,
	}
}

// Create a false literal node.
func False(pos *lexer.Position) *ast.FalseLiteralNode {
	return &ast.FalseLiteralNode{
		Position: pos,
	}
}

// Create a true literal node.
func True(pos *lexer.Position) *ast.TrueLiteralNode {
	return &ast.TrueLiteralNode{
		Position: pos,
	}
}

// Create a nil literal node.
func Nil(pos *lexer.Position) *ast.NilLiteralNode {
	return &ast.NilLiteralNode{
		Position: pos,
	}
}

// Create an empty statement node.
func EmptyStmt(pos *lexer.Position) *ast.EmptyStatementNode {
	return &ast.EmptyStatementNode{
		Position: pos,
	}
}

// Create an expression modifier node.
func Mod(pos *lexer.Position, mod *lexer.Token, left ast.ExpressionNode, right ast.ExpressionNode) *ast.ModifierNode {
	return &ast.ModifierNode{
		Position: pos,
		Left:     left,
		Modifier: mod,
		Right:    right,
	}
}

// Create an if...else expression modifier node.
func ModIfElse(pos *lexer.Position, then ast.ExpressionNode, cond ast.ExpressionNode, els ast.ExpressionNode) *ast.ModifierIfElseNode {
	return &ast.ModifierIfElseNode{
		Position:       pos,
		ThenExpression: then,
		Condition:      cond,
		ElseExpression: els,
	}
}

// Create an if expression node.
func IfExpr(pos *lexer.Position, cond ast.ExpressionNode, then []ast.StatementNode, els []ast.StatementNode) *ast.IfExpressionNode {
	return &ast.IfExpressionNode{
		Position:  pos,
		Condition: cond,
		ThenBody:  then,
		ElseBody:  els,
	}
}

// Create an if expression node.
func UnlessExpr(pos *lexer.Position, cond ast.ExpressionNode, then []ast.StatementNode, els []ast.StatementNode) *ast.UnlessExpressionNode {
	return &ast.UnlessExpressionNode{
		Position:  pos,
		Condition: cond,
		ThenBody:  then,
		ElseBody:  els,
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
		"has higher precedence than comparison operators": {
			input: "foo >= bar + baz",
			want: Prog(
				Pos(0, 16, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 16, 1, 1),
						Bin(
							Pos(0, 16, 1, 1),
							Tok(lexer.GreaterEqualToken, 4, 2, 1, 5),
							Ident("foo", Pos(0, 3, 1, 1)),
							Bin(
								Pos(7, 9, 1, 8),
								Tok(lexer.PlusToken, 11, 1, 1, 12),
								Ident("bar", Pos(7, 3, 1, 8)),
								Ident("baz", Pos(13, 3, 1, 14)),
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
						Pos(2, 4, 2, 1),
						Invalid(Tok(lexer.StarToken, 2, 1, 2, 1)),
					),
					ExprStmt(
						Pos(6, 3, 3, 1),
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
						Pos(2, 4, 2, 1),
						Invalid(Tok(lexer.SlashToken, 2, 1, 2, 1)),
					),
					ExprStmt(
						Pos(6, 3, 3, 1),
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
		"have higher precedence than multiplicative expressions": {
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
						Pos(2, 5, 2, 1),
						Invalid(Tok(lexer.StarStarToken, 2, 2, 2, 1)),
					),
					ExprStmt(
						Pos(7, 4, 3, 1),
						Invalid(Tok(lexer.StarStarToken, 7, 2, 3, 1)),
					),
				},
			),
			err: ErrorList{
				&Error{Pos(2, 2, 2, 1), "unexpected **, expected an expression"},
				&Error{Pos(7, 2, 3, 1), "unexpected **, expected an expression"},
			},
		},
		"has higher precedence than unary expressions": {
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
		"spaces can't separate statements": {
			input: "1 ** 2 \t 5 * 8",
			want: Prog(
				Pos(0, 14, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 6, 1, 1),
						Bin(
							Pos(0, 6, 1, 1),
							Tok(lexer.StarStarToken, 2, 2, 1, 3),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Int(lexer.DecIntToken, "2", 5, 1, 1, 6),
						),
					),
				},
			),
			err: ErrorList{
				&Error{Message: "unexpected DecInt, expected a statement separator `\\n`, `;` or end of file", Position: Pos(9, 1, 1, 10)},
			},
		},
		"can be empty with newlines": {
			input: "\n\n\n",
			want: Prog(
				Pos(0, 3, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 3, 1, 1)),
				},
			),
		},
		"can be empty with semicolons": {
			input: ";;;",
			want: Prog(
				Pos(0, 3, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					EmptyStmt(Pos(1, 1, 1, 2)),
					EmptyStmt(Pos(2, 1, 1, 3)),
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
						Pos(4, 6, 2, 1),
						Invalid(Tok(lexer.EqualToken, 4, 1, 2, 1)),
					),
					ExprStmt(
						Pos(10, 6, 3, 1),
						Invalid(Tok(lexer.EqualToken, 10, 1, 3, 1)),
					),
					ExprStmt(
						Pos(16, 3, 4, 1),
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
							Logic(
								Pos(4, 45, 1, 5),
								Tok(lexer.OrOrToken, 20, 2, 1, 21),
								Logic(
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
		"has lower precedence than equality": {
			input: "foo && bar == baz",
			want: Prog(
				Pos(0, 17, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 17, 1, 1),
						Logic(
							Pos(0, 17, 1, 1),
							Tok(lexer.AndAndToken, 4, 2, 1, 5),
							Ident("foo", Pos(0, 3, 1, 1)),
							Bin(
								Pos(7, 10, 1, 8),
								Tok(lexer.EqualEqualToken, 11, 2, 1, 12),
								Ident("bar", Pos(7, 3, 1, 8)),
								Ident("baz", Pos(14, 3, 1, 15)),
							),
						),
					),
				},
			),
		},
		"or has lower precedence than and": {
			input: "foo || bar && baz",
			want: Prog(
				Pos(0, 17, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 17, 1, 1),
						Logic(
							Pos(0, 17, 1, 1),
							Tok(lexer.OrOrToken, 4, 2, 1, 5),
							Ident("foo", Pos(0, 3, 1, 1)),
							Logic(
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
						Logic(
							Pos(0, 17, 1, 1),
							Tok(lexer.QuestionQuestionToken, 4, 2, 1, 5),
							Ident("foo", Pos(0, 3, 1, 1)),
							Logic(
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
						Logic(
							Pos(0, 24, 1, 1),
							Tok(lexer.QuestionQuestionToken, 18, 2, 1, 19),
							Logic(
								Pos(0, 17, 1, 1),
								Tok(lexer.OrOrToken, 11, 2, 1, 12),
								Logic(
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
						Logic(
							Pos(0, 17, 1, 1),
							Tok(lexer.OrOrToken, 11, 2, 1, 12),
							Logic(
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
						Logic(
							Pos(0, 17, 1, 1),
							Tok(lexer.AndAndToken, 11, 2, 1, 12),
							Logic(
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
						Logic(
							Pos(0, 17, 1, 1),
							Tok(lexer.QuestionQuestionToken, 11, 2, 1, 12),
							Logic(
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
				},
			),
			err: ErrorList{
				&Error{Message: "unexpected Identifier, expected a statement separator `\\n`, `;` or end of file", Position: Pos(6, 1, 1, 7)},
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

func TestEquality(t *testing.T) {
	tests := testTable{
		"is evaluated from left to right": {
			input: "bar == baz == 1",
			want: Prog(
				Pos(0, 15, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 15, 1, 1),
						Bin(
							Pos(0, 15, 1, 1),
							Tok(lexer.EqualEqualToken, 11, 2, 1, 12),
							Bin(
								Pos(0, 10, 1, 1),
								Tok(lexer.EqualEqualToken, 4, 2, 1, 5),
								Ident("bar", Pos(0, 3, 1, 1)),
								Ident("baz", Pos(7, 3, 1, 8)),
							),
							Int(lexer.DecIntToken, "1", 14, 1, 1, 15),
						),
					),
				},
			),
		},
		"can have endlines after the operator": {
			input: "bar ==\nbaz ==\n1",
			want: Prog(
				Pos(0, 15, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 15, 1, 1),
						Bin(
							Pos(0, 15, 1, 1),
							Tok(lexer.EqualEqualToken, 11, 2, 2, 5),
							Bin(
								Pos(0, 10, 1, 1),
								Tok(lexer.EqualEqualToken, 4, 2, 1, 5),
								Ident("bar", Pos(0, 3, 1, 1)),
								Ident("baz", Pos(7, 3, 2, 1)),
							),
							Int(lexer.DecIntToken, "1", 14, 1, 3, 1),
						),
					),
				},
			),
		},
		"can't have endlines before the operator": {
			input: "bar\n== baz\n== 1",
			want: Prog(
				Pos(0, 15, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 4, 1, 1),
						Ident("bar", Pos(0, 3, 1, 1)),
					),
					ExprStmt(
						Pos(4, 7, 2, 1),
						Invalid(Tok(lexer.EqualEqualToken, 4, 2, 2, 1)),
					),
					ExprStmt(
						Pos(11, 4, 3, 1),
						Invalid(Tok(lexer.EqualEqualToken, 11, 2, 3, 1)),
					),
				},
			),
			err: ErrorList{
				&Error{Message: "unexpected ==, expected an expression", Position: Pos(4, 2, 2, 1)},
				&Error{Message: "unexpected ==, expected an expression", Position: Pos(11, 2, 3, 1)},
			},
		},
		"has many versions": {
			input: "a == b != c === d !== e =:= f =!= g",
			want: Prog(
				Pos(0, 35, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 35, 1, 1),
						Bin(
							Pos(0, 35, 1, 1),
							Tok(lexer.RefNotEqualToken, 30, 3, 1, 31),
							Bin(
								Pos(0, 29, 1, 1),
								Tok(lexer.RefEqualToken, 24, 3, 1, 25),
								Bin(
									Pos(0, 23, 1, 1),
									Tok(lexer.StrictNotEqualToken, 18, 3, 1, 19),
									Bin(
										Pos(0, 17, 1, 1),
										Tok(lexer.StrictEqualToken, 12, 3, 1, 13),
										Bin(
											Pos(0, 11, 1, 1),
											Tok(lexer.NotEqualToken, 7, 2, 1, 8),
											Bin(
												Pos(0, 6, 1, 1),
												Tok(lexer.EqualEqualToken, 2, 2, 1, 3),
												Ident("a", Pos(0, 1, 1, 1)),
												Ident("b", Pos(5, 1, 1, 6)),
											),
											Ident("c", Pos(10, 1, 1, 11)),
										),
										Ident("d", Pos(16, 1, 1, 17)),
									),
									Ident("e", Pos(22, 1, 1, 23)),
								),
								Ident("f", Pos(28, 1, 1, 29)),
							),
							Ident("g", Pos(34, 1, 1, 35)),
						),
					),
				},
			),
		},
		"has higher precedence than logical operators": {
			input: "foo && bar == baz",
			want: Prog(
				Pos(0, 17, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 17, 1, 1),
						Logic(
							Pos(0, 17, 1, 1),
							Tok(lexer.AndAndToken, 4, 2, 1, 5),
							Ident("foo", Pos(0, 3, 1, 1)),
							Bin(
								Pos(7, 10, 1, 8),
								Tok(lexer.EqualEqualToken, 11, 2, 1, 12),
								Ident("bar", Pos(7, 3, 1, 8)),
								Ident("baz", Pos(14, 3, 1, 15)),
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
			want: Prog(
				Pos(0, 15, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 15, 1, 1),
						Bin(
							Pos(0, 15, 1, 1),
							Tok(lexer.GreaterToken, 10, 1, 1, 11),
							Bin(
								Pos(0, 9, 1, 1),
								Tok(lexer.GreaterToken, 4, 1, 1, 5),
								Ident("foo", Pos(0, 3, 1, 1)),
								Ident("bar", Pos(6, 3, 1, 7)),
							),
							Ident("baz", Pos(12, 3, 1, 13)),
						),
					),
				},
			),
		},
		"can have endlines after the operator": {
			input: "foo >\nbar >\nbaz",
			want: Prog(
				Pos(0, 15, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 15, 1, 1),
						Bin(
							Pos(0, 15, 1, 1),
							Tok(lexer.GreaterToken, 10, 1, 2, 5),
							Bin(
								Pos(0, 9, 1, 1),
								Tok(lexer.GreaterToken, 4, 1, 1, 5),
								Ident("foo", Pos(0, 3, 1, 1)),
								Ident("bar", Pos(6, 3, 2, 1)),
							),
							Ident("baz", Pos(12, 3, 3, 1)),
						),
					),
				},
			),
		},
		"can't have endlines before the operator": {
			input: "bar\n> baz\n> baz",
			want: Prog(
				Pos(0, 15, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 4, 1, 1),
						Ident("bar", Pos(0, 3, 1, 1)),
					),
					ExprStmt(
						Pos(4, 6, 2, 1),
						Invalid(Tok(lexer.GreaterToken, 4, 1, 2, 1)),
					),
					ExprStmt(
						Pos(10, 5, 3, 1),
						Invalid(Tok(lexer.GreaterToken, 10, 1, 3, 1)),
					),
				},
			),
			err: ErrorList{
				&Error{Message: "unexpected >, expected an expression", Position: Pos(4, 1, 2, 1)},
				&Error{Message: "unexpected >, expected an expression", Position: Pos(10, 1, 3, 1)},
			},
		},
		"has many versions": {
			input: "a < b <= c > d >= e <: f :> g <<: h :>> i <=> j",
			want: Prog(
				Pos(0, 47, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 47, 1, 1),
						Bin(
							Pos(0, 47, 1, 1),
							Tok(lexer.SpaceshipOpToken, 42, 3, 1, 43),
							Bin(
								Pos(0, 41, 1, 1),
								Tok(lexer.ReverseInstanceOfToken, 36, 3, 1, 37),
								Bin(
									Pos(0, 35, 1, 1),
									Tok(lexer.InstanceOfToken, 30, 3, 1, 31),
									Bin(
										Pos(0, 29, 1, 1),
										Tok(lexer.ReverseSubtypeToken, 25, 2, 1, 26),
										Bin(
											Pos(0, 24, 1, 1),
											Tok(lexer.SubtypeToken, 20, 2, 1, 21),
											Bin(
												Pos(0, 19, 1, 1),
												Tok(lexer.GreaterEqualToken, 15, 2, 1, 16),
												Bin(
													Pos(0, 14, 1, 1),
													Tok(lexer.GreaterToken, 11, 1, 1, 12),
													Bin(
														Pos(0, 10, 1, 1),
														Tok(lexer.LessEqualToken, 6, 2, 1, 7),
														Bin(
															Pos(0, 5, 1, 1),
															Tok(lexer.LessToken, 2, 1, 1, 3),
															Ident("a", Pos(0, 1, 1, 1)),
															Ident("b", Pos(4, 1, 1, 5)),
														),
														Ident("c", Pos(9, 1, 1, 10)),
													),
													Ident("d", Pos(13, 1, 1, 14)),
												),
												Ident("e", Pos(18, 1, 1, 19)),
											),
											Ident("f", Pos(23, 1, 1, 24)),
										),
										Ident("g", Pos(28, 1, 1, 29)),
									),
									Ident("h", Pos(34, 1, 1, 35)),
								),
								Ident("i", Pos(40, 1, 1, 41)),
							),
							Ident("j", Pos(46, 1, 1, 47)),
						),
					),
				},
			),
		},
		"has higher precedence than equality operators": {
			input: "foo == bar >= baz",
			want: Prog(
				Pos(0, 17, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 17, 1, 1),
						Bin(
							Pos(0, 17, 1, 1),
							Tok(lexer.EqualEqualToken, 4, 2, 1, 5),
							Ident("foo", Pos(0, 3, 1, 1)),
							Bin(
								Pos(7, 10, 1, 8),
								Tok(lexer.GreaterEqualToken, 11, 2, 1, 12),
								Ident("bar", Pos(7, 3, 1, 8)),
								Ident("baz", Pos(14, 3, 1, 15)),
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
			want: Prog(
				Pos(0, 16, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 16, 1, 1),
						Mod(
							Pos(0, 16, 1, 1),
							Tok(lexer.IfToken, 10, 2, 1, 11),
							Asgmt(
								Pos(0, 9, 1, 1),
								Tok(lexer.EqualToken, 4, 1, 1, 5),
								Ident("foo", Pos(0, 3, 1, 1)),
								Ident("bar", Pos(6, 3, 1, 7)),
							),
							Ident("baz", Pos(13, 3, 1, 14)),
						),
					),
				},
			),
		},
		"if can contain else": {
			input: "foo = bar if baz else car = red",
			want: Prog(
				Pos(0, 31, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 31, 1, 1),
						ModIfElse(
							Pos(0, 31, 1, 1),
							Asgmt(
								Pos(0, 9, 1, 1),
								Tok(lexer.EqualToken, 4, 1, 1, 5),
								Ident("foo", Pos(0, 3, 1, 1)),
								Ident("bar", Pos(6, 3, 1, 7)),
							),
							Ident("baz", Pos(13, 3, 1, 14)),
							Asgmt(
								Pos(22, 9, 1, 23),
								Tok(lexer.EqualToken, 26, 1, 1, 27),
								Ident("car", Pos(22, 3, 1, 23)),
								Ident("red", Pos(28, 3, 1, 29)),
							),
						),
					),
				},
			),
		},
		"has many versions": {
			input: "foo if bar\nfoo unless bar\nfoo while bar\nfoo until bar",
			want: Prog(
				Pos(0, 53, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 11, 1, 1),
						Mod(
							Pos(0, 10, 1, 1),
							Tok(lexer.IfToken, 4, 2, 1, 5),
							Ident("foo", Pos(0, 3, 1, 1)),
							Ident("bar", Pos(7, 3, 1, 8)),
						),
					),
					ExprStmt(
						Pos(11, 15, 2, 1),
						Mod(
							Pos(11, 14, 2, 1),
							Tok(lexer.UnlessToken, 15, 6, 2, 5),
							Ident("foo", Pos(11, 3, 2, 1)),
							Ident("bar", Pos(22, 3, 2, 12)),
						),
					),
					ExprStmt(
						Pos(26, 14, 3, 1),
						Mod(
							Pos(26, 13, 3, 1),
							Tok(lexer.WhileToken, 30, 5, 3, 5),
							Ident("foo", Pos(26, 3, 3, 1)),
							Ident("bar", Pos(36, 3, 3, 11)),
						),
					),
					ExprStmt(
						Pos(40, 13, 4, 1),
						Mod(
							Pos(40, 13, 4, 1),
							Tok(lexer.UntilToken, 44, 5, 4, 5),
							Ident("foo", Pos(40, 3, 4, 1)),
							Ident("bar", Pos(50, 3, 4, 11)),
						),
					),
				},
			),
		},
		"can't be nested": {
			input: "foo = bar if baz if false\n3",
			want: Prog(
				Pos(0, 27, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 16, 1, 1),
						Mod(
							Pos(0, 16, 1, 1),
							Tok(lexer.IfToken, 10, 2, 1, 11),
							Asgmt(
								Pos(0, 9, 1, 1),
								Tok(lexer.EqualToken, 4, 1, 1, 5),
								Ident("foo", Pos(0, 3, 1, 1)),
								Ident("bar", Pos(6, 3, 1, 7)),
							),
							Ident("baz", Pos(13, 3, 1, 14)),
						),
					),
					ExprStmt(
						Pos(26, 1, 2, 1),
						Int(lexer.DecIntToken, "3", 26, 1, 2, 1),
					),
				},
			),
			err: ErrorList{
				&Error{Message: "unexpected if, expected a statement separator `\\n`, `;` or end of file", Position: Pos(17, 2, 1, 18)},
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
			want: Prog(
				Pos(0, 31, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 30, 2, 1),
						IfExpr(
							Pos(1, 29, 2, 1),
							Bin(
								Pos(4, 7, 2, 4),
								Tok(lexer.GreaterToken, 8, 1, 2, 8),
								Ident("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(13, 9, 3, 2),
									Asgmt(
										Pos(13, 8, 3, 2),
										Tok(lexer.PlusEqualToken, 17, 2, 3, 6),
										Ident("foo", Pos(13, 3, 3, 2)),
										Int(lexer.DecIntToken, "2", 20, 1, 3, 9),
									),
								),
								ExprStmt(
									Pos(23, 4, 4, 2),
									Nil(Pos(23, 3, 4, 2)),
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
			want: Prog(
				Pos(0, 16, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 15, 2, 1),
						IfExpr(
							Pos(1, 14, 2, 1),
							Bin(
								Pos(4, 7, 2, 4),
								Tok(lexer.GreaterToken, 8, 1, 2, 8),
								Ident("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
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
			want: Prog(
				Pos(0, 39, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 34, 2, 1),
						Asgmt(
							Pos(1, 33, 2, 1),
							Tok(lexer.EqualToken, 5, 1, 2, 5),
							Ident("bar", Pos(1, 3, 2, 1)),
							IfExpr(
								Pos(8, 26, 3, 2),
								Bin(
									Pos(11, 7, 3, 5),
									Tok(lexer.GreaterToken, 15, 1, 3, 9),
									Ident("foo", Pos(11, 3, 3, 5)),
									Int(lexer.DecIntToken, "0", 17, 1, 3, 11),
								),
								Stmts{
									ExprStmt(
										Pos(21, 9, 4, 3),
										Asgmt(
											Pos(21, 8, 4, 3),
											Tok(lexer.PlusEqualToken, 25, 2, 4, 7),
											Ident("foo", Pos(21, 3, 4, 3)),
											Int(lexer.DecIntToken, "2", 28, 1, 4, 10),
										),
									),
								},
								nil,
							),
						),
					),
					ExprStmt(
						Pos(35, 4, 6, 1),
						Nil(Pos(35, 3, 6, 1)),
					),
				},
			),
		},
		"can be single line with then and without end": {
			input: `
if foo > 0 then foo += 2
nil
`,
			want: Prog(
				Pos(0, 30, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 25, 2, 1),
						IfExpr(
							Pos(1, 24, 2, 1),
							Bin(
								Pos(4, 7, 2, 4),
								Tok(lexer.GreaterToken, 8, 1, 2, 8),
								Ident("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(17, 8, 2, 17),
									Asgmt(
										Pos(17, 8, 2, 17),
										Tok(lexer.PlusEqualToken, 21, 2, 2, 21),
										Ident("foo", Pos(17, 3, 2, 17)),
										Int(lexer.DecIntToken, "2", 24, 1, 2, 24),
									),
								),
							},
							nil,
						),
					),
					ExprStmt(
						Pos(26, 4, 3, 1),
						Nil(Pos(26, 3, 3, 1)),
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
			want: Prog(
				Pos(0, 56, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 51, 2, 1),
						IfExpr(
							Pos(1, 50, 2, 1),
							Bin(
								Pos(4, 7, 2, 4),
								Tok(lexer.GreaterToken, 8, 1, 2, 8),
								Ident("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(13, 9, 3, 2),
									Asgmt(
										Pos(13, 8, 3, 2),
										Tok(lexer.PlusEqualToken, 17, 2, 3, 6),
										Ident("foo", Pos(13, 3, 3, 2)),
										Int(lexer.DecIntToken, "2", 20, 1, 3, 9),
									),
								),
								ExprStmt(
									Pos(23, 4, 4, 2),
									Nil(Pos(23, 3, 4, 2)),
								),
							},
							Stmts{
								ExprStmt(
									Pos(34, 9, 6, 3),
									Asgmt(
										Pos(34, 8, 6, 3),
										Tok(lexer.MinusEqualToken, 38, 2, 6, 7),
										Ident("foo", Pos(34, 3, 6, 3)),
										Int(lexer.DecIntToken, "2", 41, 1, 6, 10),
									),
								),
								ExprStmt(
									Pos(44, 4, 7, 2),
									Nil(Pos(44, 3, 7, 2)),
								),
							},
						),
					),
					ExprStmt(
						Pos(52, 4, 9, 1),
						Nil(Pos(52, 3, 9, 1)),
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
			want: Prog(
				Pos(0, 44, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 39, 2, 1),
						IfExpr(
							Pos(1, 38, 2, 1),
							Bin(
								Pos(4, 7, 2, 4),
								Tok(lexer.GreaterToken, 8, 1, 2, 8),
								Ident("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(17, 8, 2, 17),
									Asgmt(
										Pos(17, 8, 2, 17),
										Tok(lexer.PlusEqualToken, 21, 2, 2, 21),
										Ident("foo", Pos(17, 3, 2, 17)),
										Int(lexer.DecIntToken, "2", 24, 1, 2, 24),
									),
								),
							},
							Stmts{
								ExprStmt(
									Pos(31, 8, 3, 6),
									Asgmt(
										Pos(31, 8, 3, 6),
										Tok(lexer.MinusEqualToken, 35, 2, 3, 10),
										Ident("foo", Pos(31, 3, 3, 6)),
										Int(lexer.DecIntToken, "2", 38, 1, 3, 13),
									),
								),
							},
						),
					),
					ExprStmt(
						Pos(40, 4, 4, 1),
						Nil(Pos(40, 3, 4, 1)),
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
			want: Prog(
				Pos(0, 53, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 39, 2, 1),
						IfExpr(
							Pos(1, 38, 2, 1),
							Bin(
								Pos(4, 7, 2, 4),
								Tok(lexer.GreaterToken, 8, 1, 2, 8),
								Ident("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(17, 8, 2, 17),
									Asgmt(
										Pos(17, 8, 2, 17),
										Tok(lexer.PlusEqualToken, 21, 2, 2, 21),
										Ident("foo", Pos(17, 3, 2, 17)),
										Int(lexer.DecIntToken, "2", 24, 1, 2, 24),
									),
								),
							},
							Stmts{
								ExprStmt(
									Pos(31, 8, 3, 6),
									Asgmt(
										Pos(31, 8, 3, 6),
										Tok(lexer.MinusEqualToken, 35, 2, 3, 10),
										Ident("foo", Pos(31, 3, 3, 6)),
										Int(lexer.DecIntToken, "2", 38, 1, 3, 13),
									),
								),
							},
						),
					),
					ExprStmt(
						Pos(40, 9, 4, 1),
						Invalid(Tok(lexer.ElseToken, 40, 4, 4, 1)),
					),
					ExprStmt(
						Pos(49, 4, 5, 1),
						Nil(Pos(49, 3, 5, 1)),
					),
				},
			),
			err: ErrorList{
				&Error{Message: "unexpected else, expected an expression", Position: Pos(40, 4, 4, 1)},
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
			want: Prog(
				Pos(0, 104, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 99, 2, 1),
						IfExpr(
							Pos(1, 98, 2, 1),
							Bin(
								Pos(4, 7, 2, 4),
								Tok(lexer.GreaterToken, 8, 1, 2, 8),
								Ident("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(13, 9, 3, 2),
									Asgmt(
										Pos(13, 8, 3, 2),
										Tok(lexer.PlusEqualToken, 17, 2, 3, 6),
										Ident("foo", Pos(13, 3, 3, 2)),
										Int(lexer.DecIntToken, "2", 20, 1, 3, 9),
									),
								),
								ExprStmt(
									Pos(23, 4, 4, 2),
									Nil(Pos(23, 3, 4, 2)),
								),
							},
							Stmts{
								ExprStmt(
									Pos(27, 25, 5, 1),
									IfExpr(
										Pos(27, 25, 5, 1),
										Bin(
											Pos(33, 7, 5, 7),
											Tok(lexer.LessToken, 37, 1, 5, 11),
											Ident("foo", Pos(33, 3, 5, 7)),
											Int(lexer.DecIntToken, "5", 39, 1, 5, 13),
										),
										Stmts{
											ExprStmt(
												Pos(42, 10, 6, 2),
												Asgmt(
													Pos(42, 9, 6, 2),
													Tok(lexer.StarEqualToken, 46, 2, 6, 6),
													Ident("foo", Pos(42, 3, 6, 2)),
													Int(lexer.DecIntToken, "10", 49, 2, 6, 9),
												),
											),
										},
										Stmts{
											ExprStmt(
												Pos(52, 47, 7, 1),
												IfExpr(
													Pos(52, 47, 7, 1),
													Bin(
														Pos(58, 7, 7, 7),
														Tok(lexer.LessToken, 62, 1, 7, 11),
														Ident("foo", Pos(58, 3, 7, 7)),
														Int(lexer.DecIntToken, "0", 64, 1, 7, 13),
													),
													Stmts{
														ExprStmt(
															Pos(67, 9, 8, 2),
															Asgmt(
																Pos(67, 8, 8, 2),
																Tok(lexer.PercentEqualToken, 71, 2, 8, 6),
																Ident("foo", Pos(67, 3, 8, 2)),
																Int(lexer.DecIntToken, "3", 74, 1, 8, 9),
															),
														),
													},
													Stmts{
														ExprStmt(
															Pos(82, 9, 10, 2),
															Asgmt(
																Pos(82, 8, 10, 2),
																Tok(lexer.MinusEqualToken, 86, 2, 10, 6),
																Ident("foo", Pos(82, 3, 10, 2)),
																Int(lexer.DecIntToken, "2", 89, 1, 10, 9),
															),
														),
														ExprStmt(
															Pos(92, 4, 11, 2),
															Nil(Pos(92, 3, 11, 2)),
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
					ExprStmt(
						Pos(100, 4, 13, 1),
						Nil(Pos(100, 3, 13, 1)),
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
			want: Prog(
				Pos(0, 101, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 96, 2, 1),
						IfExpr(
							Pos(1, 95, 2, 1),
							Bin(
								Pos(4, 7, 2, 4),
								Tok(lexer.GreaterToken, 8, 1, 2, 8),
								Ident("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(17, 8, 2, 17),
									Asgmt(
										Pos(17, 8, 2, 17),
										Tok(lexer.PlusEqualToken, 21, 2, 2, 21),
										Ident("foo", Pos(17, 3, 2, 17)),
										Int(lexer.DecIntToken, "2", 24, 1, 2, 24),
									),
								),
							},
							Stmts{
								ExprStmt(
									Pos(26, 28, 3, 1),
									IfExpr(
										Pos(26, 28, 3, 1),
										Bin(
											Pos(32, 7, 3, 7),
											Tok(lexer.LessToken, 36, 1, 3, 11),
											Ident("foo", Pos(32, 3, 3, 7)),
											Int(lexer.DecIntToken, "5", 38, 1, 3, 13),
										),
										Stmts{
											ExprStmt(
												Pos(45, 9, 3, 20),
												Asgmt(
													Pos(45, 9, 3, 20),
													Tok(lexer.StarEqualToken, 49, 2, 3, 24),
													Ident("foo", Pos(45, 3, 3, 20)),
													Int(lexer.DecIntToken, "10", 52, 2, 3, 27),
												),
											),
										},
										Stmts{
											ExprStmt(
												Pos(55, 41, 4, 1),
												IfExpr(
													Pos(55, 41, 4, 1),
													Bin(
														Pos(61, 7, 4, 7),
														Tok(lexer.LessToken, 65, 1, 4, 11),
														Ident("foo", Pos(61, 3, 4, 7)),
														Int(lexer.DecIntToken, "0", 67, 1, 4, 13),
													),
													Stmts{
														ExprStmt(
															Pos(74, 8, 4, 20),
															Asgmt(
																Pos(74, 8, 4, 20),
																Tok(lexer.PercentEqualToken, 78, 2, 4, 24),
																Ident("foo", Pos(74, 3, 4, 20)),
																Int(lexer.DecIntToken, "3", 81, 1, 4, 27),
															),
														),
													},
													Stmts{
														ExprStmt(
															Pos(88, 8, 5, 6),
															Asgmt(
																Pos(88, 8, 5, 6),
																Tok(lexer.MinusEqualToken, 92, 2, 5, 10),
																Ident("foo", Pos(88, 3, 5, 6)),
																Int(lexer.DecIntToken, "2", 95, 1, 5, 13),
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
					ExprStmt(
						Pos(97, 4, 6, 1),
						Nil(Pos(97, 3, 6, 1)),
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
			want: Prog(
				Pos(0, 108, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 103, 2, 1),
						IfExpr(
							Pos(1, 102, 2, 1),
							Bin(
								Pos(4, 7, 2, 4),
								Tok(lexer.GreaterToken, 8, 1, 2, 8),
								Ident("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(13, 9, 3, 2),
									Asgmt(
										Pos(13, 8, 3, 2),
										Tok(lexer.PlusEqualToken, 17, 2, 3, 6),
										Ident("foo", Pos(13, 3, 3, 2)),
										Int(lexer.DecIntToken, "2", 20, 1, 3, 9),
									),
								),
								ExprStmt(
									Pos(23, 4, 4, 2),
									Nil(Pos(23, 3, 4, 2)),
								),
							},
							Stmts{
								ExprStmt(
									Pos(32, 71, 5, 6),
									IfExpr(
										Pos(32, 71, 5, 6),
										Bin(
											Pos(35, 7, 5, 9),
											Tok(lexer.LessToken, 39, 1, 5, 13),
											Ident("foo", Pos(35, 3, 5, 9)),
											Int(lexer.DecIntToken, "5", 41, 1, 5, 15),
										),
										Stmts{
											ExprStmt(
												Pos(44, 10, 6, 2),
												Asgmt(
													Pos(44, 9, 6, 2),
													Tok(lexer.StarEqualToken, 48, 2, 6, 6),
													Ident("foo", Pos(44, 3, 6, 2)),
													Int(lexer.DecIntToken, "10", 51, 2, 6, 9),
												),
											),
										},
										Stmts{
											ExprStmt(
												Pos(59, 44, 7, 6),
												IfExpr(
													Pos(59, 44, 7, 6),
													Bin(
														Pos(62, 7, 7, 9),
														Tok(lexer.LessToken, 66, 1, 7, 13),
														Ident("foo", Pos(62, 3, 7, 9)),
														Int(lexer.DecIntToken, "0", 68, 1, 7, 15),
													),
													Stmts{
														ExprStmt(
															Pos(71, 9, 8, 2),
															Asgmt(
																Pos(71, 8, 8, 2),
																Tok(lexer.PercentEqualToken, 75, 2, 8, 6),
																Ident("foo", Pos(71, 3, 8, 2)),
																Int(lexer.DecIntToken, "3", 78, 1, 8, 9),
															),
														),
													},
													Stmts{
														ExprStmt(
															Pos(86, 9, 10, 2),
															Asgmt(
																Pos(86, 8, 10, 2),
																Tok(lexer.MinusEqualToken, 90, 2, 10, 6),
																Ident("foo", Pos(86, 3, 10, 2)),
																Int(lexer.DecIntToken, "2", 93, 1, 10, 9),
															),
														),
														ExprStmt(
															Pos(96, 4, 11, 2),
															Nil(Pos(96, 3, 11, 2)),
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
					ExprStmt(
						Pos(104, 4, 13, 1),
						Nil(Pos(104, 3, 13, 1)),
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
			want: Prog(
				Pos(0, 35, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 34, 2, 1),
						UnlessExpr(
							Pos(1, 33, 2, 1),
							Bin(
								Pos(8, 7, 2, 8),
								Tok(lexer.GreaterToken, 12, 1, 2, 12),
								Ident("foo", Pos(8, 3, 2, 8)),
								Int(lexer.DecIntToken, "0", 14, 1, 2, 14),
							),
							Stmts{
								ExprStmt(
									Pos(17, 9, 3, 2),
									Asgmt(
										Pos(17, 8, 3, 2),
										Tok(lexer.PlusEqualToken, 21, 2, 3, 6),
										Ident("foo", Pos(17, 3, 3, 2)),
										Int(lexer.DecIntToken, "2", 24, 1, 3, 9),
									),
								),
								ExprStmt(
									Pos(27, 4, 4, 2),
									Nil(Pos(27, 3, 4, 2)),
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
			want: Prog(
				Pos(0, 20, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 19, 2, 1),
						UnlessExpr(
							Pos(1, 18, 2, 1),
							Bin(
								Pos(8, 7, 2, 8),
								Tok(lexer.GreaterToken, 12, 1, 2, 12),
								Ident("foo", Pos(8, 3, 2, 8)),
								Int(lexer.DecIntToken, "0", 14, 1, 2, 14),
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
			want: Prog(
				Pos(0, 43, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 38, 2, 1),
						Asgmt(
							Pos(1, 37, 2, 1),
							Tok(lexer.EqualToken, 5, 1, 2, 5),
							Ident("bar", Pos(1, 3, 2, 1)),
							UnlessExpr(
								Pos(8, 30, 3, 2),
								Bin(
									Pos(15, 7, 3, 9),
									Tok(lexer.GreaterToken, 19, 1, 3, 13),
									Ident("foo", Pos(15, 3, 3, 9)),
									Int(lexer.DecIntToken, "0", 21, 1, 3, 15),
								),
								Stmts{
									ExprStmt(
										Pos(25, 9, 4, 3),
										Asgmt(
											Pos(25, 8, 4, 3),
											Tok(lexer.PlusEqualToken, 29, 2, 4, 7),
											Ident("foo", Pos(25, 3, 4, 3)),
											Int(lexer.DecIntToken, "2", 32, 1, 4, 10),
										),
									),
								},
								nil,
							),
						),
					),
					ExprStmt(
						Pos(39, 4, 6, 1),
						Nil(Pos(39, 3, 6, 1)),
					),
				},
			),
		},
		"can be single line with then and without end": {
			input: `
unless foo > 0 then foo += 2
nil
`,
			want: Prog(
				Pos(0, 34, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 29, 2, 1),
						UnlessExpr(
							Pos(1, 28, 2, 1),
							Bin(
								Pos(8, 7, 2, 8),
								Tok(lexer.GreaterToken, 12, 1, 2, 12),
								Ident("foo", Pos(8, 3, 2, 8)),
								Int(lexer.DecIntToken, "0", 14, 1, 2, 14),
							),
							Stmts{
								ExprStmt(
									Pos(21, 8, 2, 21),
									Asgmt(
										Pos(21, 8, 2, 21),
										Tok(lexer.PlusEqualToken, 25, 2, 2, 25),
										Ident("foo", Pos(21, 3, 2, 21)),
										Int(lexer.DecIntToken, "2", 28, 1, 2, 28),
									),
								),
							},
							nil,
						),
					),
					ExprStmt(
						Pos(30, 4, 3, 1),
						Nil(Pos(30, 3, 3, 1)),
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
			want: Prog(
				Pos(0, 59, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 54, 2, 1),
						UnlessExpr(
							Pos(1, 53, 2, 1),
							Bin(
								Pos(8, 7, 2, 8),
								Tok(lexer.GreaterToken, 12, 1, 2, 12),
								Ident("foo", Pos(8, 3, 2, 8)),
								Int(lexer.DecIntToken, "0", 14, 1, 2, 14),
							),
							Stmts{
								ExprStmt(
									Pos(17, 9, 3, 2),
									Asgmt(
										Pos(17, 8, 3, 2),
										Tok(lexer.PlusEqualToken, 21, 2, 3, 6),
										Ident("foo", Pos(17, 3, 3, 2)),
										Int(lexer.DecIntToken, "2", 24, 1, 3, 9),
									),
								),
								ExprStmt(
									Pos(27, 4, 4, 2),
									Nil(Pos(27, 3, 4, 2)),
								),
							},
							Stmts{
								ExprStmt(
									Pos(37, 9, 6, 2),
									Asgmt(
										Pos(37, 8, 6, 2),
										Tok(lexer.MinusEqualToken, 41, 2, 6, 6),
										Ident("foo", Pos(37, 3, 6, 2)),
										Int(lexer.DecIntToken, "2", 44, 1, 6, 9),
									),
								),
								ExprStmt(
									Pos(47, 4, 7, 2),
									Nil(Pos(47, 3, 7, 2)),
								),
							},
						),
					),
					ExprStmt(
						Pos(55, 4, 9, 1),
						Nil(Pos(55, 3, 9, 1)),
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
			want: Prog(
				Pos(0, 48, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 43, 2, 1),
						UnlessExpr(
							Pos(1, 42, 2, 1),
							Bin(
								Pos(8, 7, 2, 8),
								Tok(lexer.GreaterToken, 12, 1, 2, 12),
								Ident("foo", Pos(8, 3, 2, 8)),
								Int(lexer.DecIntToken, "0", 14, 1, 2, 14),
							),
							Stmts{
								ExprStmt(
									Pos(21, 8, 2, 21),
									Asgmt(
										Pos(21, 8, 2, 21),
										Tok(lexer.PlusEqualToken, 25, 2, 2, 25),
										Ident("foo", Pos(21, 3, 2, 21)),
										Int(lexer.DecIntToken, "2", 28, 1, 2, 28),
									),
								),
							},
							Stmts{
								ExprStmt(
									Pos(35, 8, 3, 6),
									Asgmt(
										Pos(35, 8, 3, 6),
										Tok(lexer.MinusEqualToken, 39, 2, 3, 10),
										Ident("foo", Pos(35, 3, 3, 6)),
										Int(lexer.DecIntToken, "2", 42, 1, 3, 13),
									),
								),
							},
						),
					),
					ExprStmt(
						Pos(44, 4, 4, 1),
						Nil(Pos(44, 3, 4, 1)),
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
			want: Prog(
				Pos(0, 57, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 43, 2, 1),
						UnlessExpr(
							Pos(1, 42, 2, 1),
							Bin(
								Pos(8, 7, 2, 8),
								Tok(lexer.GreaterToken, 12, 1, 2, 12),
								Ident("foo", Pos(8, 3, 2, 8)),
								Int(lexer.DecIntToken, "0", 14, 1, 2, 14),
							),
							Stmts{
								ExprStmt(
									Pos(21, 8, 2, 21),
									Asgmt(
										Pos(21, 8, 2, 21),
										Tok(lexer.PlusEqualToken, 25, 2, 2, 25),
										Ident("foo", Pos(21, 3, 2, 21)),
										Int(lexer.DecIntToken, "2", 28, 1, 2, 28),
									),
								),
							},
							Stmts{
								ExprStmt(
									Pos(35, 8, 3, 6),
									Asgmt(
										Pos(35, 8, 3, 6),
										Tok(lexer.MinusEqualToken, 39, 2, 3, 10),
										Ident("foo", Pos(35, 3, 3, 6)),
										Int(lexer.DecIntToken, "2", 42, 1, 3, 13),
									),
								),
							},
						),
					),
					ExprStmt(
						Pos(44, 9, 4, 1),
						Invalid(Tok(lexer.ElseToken, 44, 4, 4, 1)),
					),
					ExprStmt(
						Pos(53, 4, 5, 1),
						Nil(Pos(53, 3, 5, 1)),
					),
				},
			),
			err: ErrorList{
				&Error{Message: "unexpected else, expected an expression", Position: Pos(44, 4, 4, 1)},
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
