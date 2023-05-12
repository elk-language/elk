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
		Token:    VTok(Pos(startByte, byteLength, line, column), tokenType, value),
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

// Create a nilable type node.
func Nilable(pos *lexer.Position, typ ast.TypeNode) *ast.NilableTypeNode {
	return &ast.NilableTypeNode{
		Position: pos,
		Type:     typ,
	}
}

// Create a binary type expression node.
func BinType(pos *lexer.Position, op *lexer.Token, left ast.TypeNode, right ast.TypeNode) *ast.BinaryTypeExpressionNode {
	return &ast.BinaryTypeExpressionNode{
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
func PubIdent(value string, pos *lexer.Position) *ast.PublicIdentifierNode {
	return &ast.PublicIdentifierNode{
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
func PubConst(value string, pos *lexer.Position) *ast.PublicConstantNode {
	return &ast.PublicConstantNode{
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

// Create a break expression node.
func Break(pos *lexer.Position) *ast.BreakExpressionNode {
	return &ast.BreakExpressionNode{
		Position: pos,
	}
}

// Create a return expression node.
func Return(pos *lexer.Position, val ast.ExpressionNode) *ast.ReturnExpressionNode {
	return &ast.ReturnExpressionNode{
		Position: pos,
		Value:    val,
	}
}

// Create a continue expression node.
func Continue(pos *lexer.Position, val ast.ExpressionNode) *ast.ContinueExpressionNode {
	return &ast.ContinueExpressionNode{
		Position: pos,
		Value:    val,
	}
}

// Create a throw expression node.
func Throw(pos *lexer.Position, val ast.ExpressionNode) *ast.ThrowExpressionNode {
	return &ast.ThrowExpressionNode{
		Position: pos,
		Value:    val,
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

// Create an if expression node.
func WhileExpr(pos *lexer.Position, cond ast.ExpressionNode, then []ast.StatementNode) *ast.WhileExpressionNode {
	return &ast.WhileExpressionNode{
		Position:  pos,
		Condition: cond,
		ThenBody:  then,
	}
}

// Create an if expression node.
func UntilExpr(pos *lexer.Position, cond ast.ExpressionNode, then []ast.StatementNode) *ast.UntilExpressionNode {
	return &ast.UntilExpressionNode{
		Position:  pos,
		Condition: cond,
		ThenBody:  then,
	}
}

// Create a loop expression node.
func LoopExpr(pos *lexer.Position, then []ast.StatementNode) *ast.LoopExpressionNode {
	return &ast.LoopExpressionNode{
		Position: pos,
		ThenBody: then,
	}
}

// Create a variable declaration node.
func Var(pos *lexer.Position, name *lexer.Token, typ ast.TypeNode, init ast.ExpressionNode) *ast.VariableDeclarationNode {
	return &ast.VariableDeclarationNode{
		Position:    pos,
		Name:        name,
		Type:        typ,
		Initialiser: init,
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
							Tok(Pos(6, 1, 1, 7), lexer.PlusToken),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(Pos(2, 1, 1, 3), lexer.PlusToken),
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
							Tok(Pos(6, 1, 2, 3), lexer.PlusToken),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(Pos(2, 1, 1, 3), lexer.PlusToken),
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
							Tok(Pos(2, 1, 2, 1), lexer.PlusToken),
							Int(lexer.DecIntToken, "2", 4, 1, 2, 3),
						),
					),
					ExprStmt(
						Pos(6, 3, 3, 1),
						Unary(
							Pos(6, 3, 3, 1),
							Tok(Pos(6, 1, 3, 1), lexer.PlusToken),
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
							Tok(Pos(4, 2, 1, 5), lexer.GreaterEqualToken),
							PubIdent("foo", Pos(0, 3, 1, 1)),
							Bin(
								Pos(7, 9, 1, 8),
								Tok(Pos(11, 1, 1, 12), lexer.PlusToken),
								PubIdent("bar", Pos(7, 3, 1, 8)),
								PubIdent("baz", Pos(13, 3, 1, 14)),
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
							Tok(Pos(6, 1, 1, 7), lexer.MinusToken),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(Pos(2, 1, 1, 3), lexer.MinusToken),
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
							Tok(Pos(6, 1, 2, 3), lexer.MinusToken),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(Pos(2, 1, 1, 3), lexer.MinusToken),
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
							Tok(Pos(2, 1, 2, 1), lexer.MinusToken),
							Int(lexer.DecIntToken, "2", 4, 1, 2, 3),
						),
					),
					ExprStmt(
						Pos(6, 3, 3, 1),
						Unary(
							Pos(6, 3, 3, 1),
							Tok(Pos(6, 1, 3, 1), lexer.MinusToken),
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
							Tok(Pos(6, 1, 1, 7), lexer.MinusToken),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(Pos(2, 1, 1, 3), lexer.PlusToken),
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
							Tok(Pos(6, 1, 1, 7), lexer.StarToken),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(Pos(2, 1, 1, 3), lexer.StarToken),
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
							Tok(Pos(6, 1, 2, 3), lexer.StarToken),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(Pos(2, 1, 1, 3), lexer.StarToken),
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
						Invalid(Tok(Pos(2, 1, 2, 1), lexer.StarToken)),
					),
					ExprStmt(
						Pos(6, 3, 3, 1),
						Invalid(Tok(Pos(6, 1, 3, 1), lexer.StarToken)),
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
							Tok(Pos(2, 1, 1, 3), lexer.PlusToken),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Bin(
								Pos(4, 5, 1, 5),
								Tok(Pos(6, 1, 1, 7), lexer.StarToken),
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
							Tok(Pos(6, 1, 1, 7), lexer.SlashToken),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(Pos(2, 1, 1, 3), lexer.SlashToken),
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
							Tok(Pos(6, 1, 2, 3), lexer.SlashToken),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(Pos(2, 1, 1, 3), lexer.SlashToken),
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
						Invalid(Tok(Pos(2, 1, 2, 1), lexer.SlashToken)),
					),
					ExprStmt(
						Pos(6, 3, 3, 1),
						Invalid(Tok(Pos(6, 1, 3, 1), lexer.SlashToken)),
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
							Tok(Pos(6, 1, 1, 7), lexer.SlashToken),
							Bin(
								Pos(0, 5, 1, 1),
								Tok(Pos(2, 1, 1, 3), lexer.StarToken),
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
							Tok(Pos(0, 1, 1, 1), lexer.PlusToken),
							Unary(
								Pos(1, 5, 1, 2),
								Tok(Pos(1, 1, 1, 2), lexer.PlusToken),
								Unary(
									Pos(2, 4, 1, 3),
									Tok(Pos(2, 1, 1, 3), lexer.PlusToken),
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
							Tok(Pos(0, 1, 1, 1), lexer.MinusToken),
							Unary(
								Pos(1, 5, 1, 2),
								Tok(Pos(1, 1, 1, 2), lexer.MinusToken),
								Unary(
									Pos(2, 4, 1, 3),
									Tok(Pos(2, 1, 1, 3), lexer.MinusToken),
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
							Tok(Pos(0, 1, 1, 1), lexer.BangToken),
							Unary(
								Pos(1, 5, 1, 2),
								Tok(Pos(1, 1, 1, 2), lexer.BangToken),
								Unary(
									Pos(2, 4, 1, 3),
									Tok(Pos(2, 1, 1, 3), lexer.BangToken),
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
							Tok(Pos(0, 1, 1, 1), lexer.TildeToken),
							Unary(
								Pos(1, 5, 1, 2),
								Tok(Pos(1, 1, 1, 2), lexer.TildeToken),
								Unary(
									Pos(2, 4, 1, 3),
									Tok(Pos(2, 1, 1, 3), lexer.TildeToken),
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
							Tok(Pos(0, 1, 1, 1), lexer.BangToken),
							Unary(
								Pos(1, 5, 1, 2),
								Tok(Pos(1, 1, 1, 2), lexer.PlusToken),
								Unary(
									Pos(2, 4, 1, 3),
									Tok(Pos(2, 1, 1, 3), lexer.TildeToken),
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
							Tok(Pos(10, 1, 1, 11), lexer.PlusToken),
							Bin(
								Pos(0, 9, 1, 1),
								Tok(Pos(6, 1, 1, 7), lexer.StarToken),
								Unary(
									Pos(0, 5, 1, 1),
									Tok(Pos(0, 1, 1, 1), lexer.BangToken),
									Unary(
										Pos(1, 4, 1, 2),
										Tok(Pos(1, 1, 1, 2), lexer.BangToken),
										Float("1.5", 2, 3, 1, 3),
									),
								),
								Int(lexer.DecIntToken, "2", 8, 1, 1, 9),
							),
							Unary(
								Pos(12, 3, 1, 13),
								Tok(Pos(12, 1, 1, 13), lexer.TildeToken),
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
							Tok(Pos(2, 2, 1, 3), lexer.StarStarToken),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Bin(
								Pos(5, 6, 1, 6),
								Tok(Pos(7, 2, 1, 8), lexer.StarStarToken),
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
							Tok(Pos(2, 2, 1, 3), lexer.StarStarToken),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Bin(
								Pos(5, 6, 2, 1),
								Tok(Pos(7, 2, 2, 3), lexer.StarStarToken),
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
						Invalid(Tok(Pos(2, 2, 2, 1), lexer.StarStarToken)),
					),
					ExprStmt(
						Pos(7, 4, 3, 1),
						Invalid(Tok(Pos(7, 2, 3, 1), lexer.StarStarToken)),
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
							Tok(Pos(0, 1, 1, 1), lexer.MinusToken),
							Bin(
								Pos(1, 6, 1, 2),
								Tok(Pos(3, 2, 1, 4), lexer.StarStarToken),
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
							Tok(Pos(2, 2, 1, 3), lexer.StarStarToken),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Int(lexer.DecIntToken, "2", 5, 1, 1, 6),
						),
					),
					ExprStmt(
						Pos(8, 5, 1, 9),
						Bin(
							Pos(8, 5, 1, 9),
							Tok(Pos(10, 1, 1, 11), lexer.StarToken),
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
							Tok(Pos(2, 2, 1, 3), lexer.StarStarToken),
							Int(lexer.DecIntToken, "1", 0, 1, 1, 1),
							Int(lexer.DecIntToken, "2", 5, 1, 1, 6),
						),
					),
					ExprStmt(
						Pos(7, 5, 2, 1),
						Bin(
							Pos(7, 5, 2, 1),
							Tok(Pos(9, 1, 2, 3), lexer.StarToken),
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
							Tok(Pos(2, 2, 1, 3), lexer.StarStarToken),
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
							Tok(Pos(2, 2, 1, 3), lexer.MinusEqualToken),
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
							Tok(Pos(6, 2, 1, 7), lexer.MinusEqualToken),
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
							Tok(Pos(6, 2, 1, 7), lexer.MinusEqualToken),
							PubConst("FooBa", Pos(0, 5, 1, 1)),
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
							Tok(Pos(6, 2, 1, 7), lexer.MinusEqualToken),
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
							Tok(Pos(4, 2, 1, 5), lexer.MinusEqualToken),
							PubIdent("foo", Pos(0, 3, 1, 1)),
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
							Tok(Pos(4, 2, 1, 5), lexer.MinusEqualToken),
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
							Tok(Pos(4, 1, 1, 5), lexer.EqualToken),
							PubIdent("foo", Pos(0, 3, 1, 1)),
							Asgmt(
								Pos(6, 13, 1, 7),
								Tok(Pos(10, 1, 1, 11), lexer.EqualToken),
								PubIdent("bar", Pos(6, 3, 1, 7)),
								Asgmt(
									Pos(12, 7, 1, 13),
									Tok(Pos(16, 1, 1, 17), lexer.EqualToken),
									PubIdent("baz", Pos(12, 3, 1, 13)),
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
							Tok(Pos(4, 1, 1, 5), lexer.EqualToken),
							PubIdent("foo", Pos(0, 3, 1, 1)),
							Asgmt(
								Pos(6, 13, 2, 1),
								Tok(Pos(10, 1, 2, 5), lexer.EqualToken),
								PubIdent("bar", Pos(6, 3, 2, 1)),
								Asgmt(
									Pos(12, 7, 3, 1),
									Tok(Pos(16, 1, 3, 5), lexer.EqualToken),
									PubIdent("baz", Pos(12, 3, 3, 1)),
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
						PubIdent("foo", Pos(0, 3, 1, 1)),
					),
					ExprStmt(
						Pos(4, 6, 2, 1),
						Invalid(Tok(Pos(4, 1, 2, 1), lexer.EqualToken)),
					),
					ExprStmt(
						Pos(10, 6, 3, 1),
						Invalid(Tok(Pos(10, 1, 3, 1), lexer.EqualToken)),
					),
					ExprStmt(
						Pos(16, 3, 4, 1),
						Invalid(Tok(Pos(16, 1, 4, 1), lexer.EqualToken)),
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
							Tok(Pos(2, 1, 1, 3), lexer.EqualToken),
							PubIdent("f", Pos(0, 1, 1, 1)),
							Logic(
								Pos(4, 45, 1, 5),
								Tok(Pos(20, 2, 1, 21), lexer.OrOrToken),
								Logic(
									Pos(4, 15, 1, 5),
									Tok(Pos(9, 2, 1, 10), lexer.AndAndToken),
									PubIdent("some", Pos(4, 4, 1, 5)),
									PubIdent("awesome", Pos(12, 7, 1, 13)),
								),
								Bin(
									Pos(23, 26, 1, 24),
									Tok(Pos(41, 2, 1, 42), lexer.EqualEqualToken),
									Bin(
										Pos(23, 17, 1, 24),
										Tok(Pos(37, 1, 1, 38), lexer.GreaterToken),
										Bin(
											Pos(23, 13, 1, 24),
											Tok(Pos(29, 1, 1, 30), lexer.PlusToken),
											PubIdent("thing", Pos(23, 5, 1, 24)),
											Bin(
												Pos(31, 5, 1, 32),
												Tok(Pos(33, 1, 1, 34), lexer.StarToken),
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
							Tok(Pos(2, 1, 1, 3), lexer.EqualToken),
							PubIdent("a", Pos(0, 1, 1, 1)),
							Asgmt(
								Pos(4, 82, 1, 5),
								Tok(Pos(6, 2, 1, 7), lexer.MinusEqualToken),
								PubIdent("b", Pos(4, 1, 1, 5)),
								Asgmt(
									Pos(9, 77, 1, 10),
									Tok(Pos(11, 2, 1, 12), lexer.PlusEqualToken),
									PubIdent("c", Pos(9, 1, 1, 10)),
									Asgmt(
										Pos(14, 72, 1, 15),
										Tok(Pos(16, 2, 1, 17), lexer.StarEqualToken),
										PubIdent("d", Pos(14, 1, 1, 15)),
										Asgmt(
											Pos(19, 67, 1, 20),
											Tok(Pos(21, 2, 1, 22), lexer.SlashEqualToken),
											PubIdent("e", Pos(19, 1, 1, 20)),
											Asgmt(
												Pos(24, 62, 1, 25),
												Tok(Pos(26, 3, 1, 27), lexer.StarStarEqualToken),
												PubIdent("f", Pos(24, 1, 1, 25)),
												Asgmt(
													Pos(30, 56, 1, 31),
													Tok(Pos(32, 2, 1, 33), lexer.TildeEqualToken),
													PubIdent("g", Pos(30, 1, 1, 31)),
													Asgmt(
														Pos(35, 51, 1, 36),
														Tok(Pos(37, 3, 1, 38), lexer.AndAndEqualToken),
														PubIdent("h", Pos(35, 1, 1, 36)),
														Asgmt(
															Pos(41, 45, 1, 42),
															Tok(Pos(43, 2, 1, 44), lexer.AndEqualToken),
															PubIdent("i", Pos(41, 1, 1, 42)),
															Asgmt(
																Pos(46, 40, 1, 47),
																Tok(Pos(48, 3, 1, 49), lexer.OrOrEqualToken),
																PubIdent("j", Pos(46, 1, 1, 47)),
																Asgmt(
																	Pos(52, 34, 1, 53),
																	Tok(Pos(54, 2, 1, 55), lexer.OrEqualToken),
																	PubIdent("k", Pos(52, 1, 1, 53)),
																	Asgmt(
																		Pos(57, 29, 1, 58),
																		Tok(Pos(59, 2, 1, 60), lexer.XorEqualToken),
																		PubIdent("l", Pos(57, 1, 1, 58)),
																		Asgmt(
																			Pos(62, 24, 1, 63),
																			Tok(Pos(64, 3, 1, 65), lexer.QuestionQuestionEqualToken),
																			PubIdent("m", Pos(62, 1, 1, 63)),
																			Asgmt(
																				Pos(68, 18, 1, 69),
																				Tok(Pos(70, 3, 1, 71), lexer.LBitShiftEqualToken),
																				PubIdent("n", Pos(68, 1, 1, 69)),
																				Asgmt(
																					Pos(74, 12, 1, 75),
																					Tok(Pos(76, 3, 1, 77), lexer.RBitShiftEqualToken),
																					PubIdent("o", Pos(74, 1, 1, 75)),
																					Asgmt(
																						Pos(80, 6, 1, 81),
																						Tok(Pos(82, 2, 1, 83), lexer.PercentEqualToken),
																						PubIdent("p", Pos(80, 1, 1, 81)),
																						PubIdent("q", Pos(85, 1, 1, 86)),
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
							Tok(Pos(4, 2, 1, 5), lexer.AndAndToken),
							PubIdent("foo", Pos(0, 3, 1, 1)),
							Bin(
								Pos(7, 10, 1, 8),
								Tok(Pos(11, 2, 1, 12), lexer.EqualEqualToken),
								PubIdent("bar", Pos(7, 3, 1, 8)),
								PubIdent("baz", Pos(14, 3, 1, 15)),
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
							Tok(Pos(4, 2, 1, 5), lexer.OrOrToken),
							PubIdent("foo", Pos(0, 3, 1, 1)),
							Logic(
								Pos(7, 10, 1, 8),
								Tok(Pos(11, 2, 1, 12), lexer.AndAndToken),
								PubIdent("bar", Pos(7, 3, 1, 8)),
								PubIdent("baz", Pos(14, 3, 1, 15)),
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
							Tok(Pos(4, 2, 1, 5), lexer.QuestionQuestionToken),
							PubIdent("foo", Pos(0, 3, 1, 1)),
							Logic(
								Pos(7, 10, 1, 8),
								Tok(Pos(11, 2, 1, 12), lexer.AndAndToken),
								PubIdent("bar", Pos(7, 3, 1, 8)),
								PubIdent("baz", Pos(14, 3, 1, 15)),
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
							Tok(Pos(18, 2, 1, 19), lexer.QuestionQuestionToken),
							Logic(
								Pos(0, 17, 1, 1),
								Tok(Pos(11, 2, 1, 12), lexer.OrOrToken),
								Logic(
									Pos(0, 10, 1, 1),
									Tok(Pos(4, 2, 1, 5), lexer.QuestionQuestionToken),
									PubIdent("foo", Pos(0, 3, 1, 1)),
									PubIdent("bar", Pos(7, 3, 1, 8)),
								),
								PubIdent("baz", Pos(14, 3, 1, 15)),
							),
							PubIdent("boo", Pos(21, 3, 1, 22)),
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
							Tok(Pos(11, 2, 1, 12), lexer.OrOrToken),
							Logic(
								Pos(0, 10, 1, 1),
								Tok(Pos(4, 2, 1, 5), lexer.OrOrToken),
								PubIdent("foo", Pos(0, 3, 1, 1)),
								PubIdent("bar", Pos(7, 3, 1, 8)),
							),
							PubIdent("baz", Pos(14, 3, 1, 15)),
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
							Tok(Pos(11, 2, 1, 12), lexer.AndAndToken),
							Logic(
								Pos(0, 10, 1, 1),
								Tok(Pos(4, 2, 1, 5), lexer.AndAndToken),
								PubIdent("foo", Pos(0, 3, 1, 1)),
								PubIdent("bar", Pos(7, 3, 1, 8)),
							),
							PubIdent("baz", Pos(14, 3, 1, 15)),
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
							Tok(Pos(11, 2, 1, 12), lexer.QuestionQuestionToken),
							Logic(
								Pos(0, 10, 1, 1),
								Tok(Pos(4, 2, 1, 5), lexer.QuestionQuestionToken),
								PubIdent("foo", Pos(0, 3, 1, 1)),
								PubIdent("bar", Pos(7, 3, 1, 8)),
							),
							PubIdent("baz", Pos(14, 3, 1, 15)),
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
							Invalid(VTok(Pos(5, 4, 1, 6), lexer.ErrorToken, "invalid hex escape in string literal")),
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
							Invalid(VTok(Pos(5, 2, 1, 6), lexer.ErrorToken, "invalid escape sequence `\\q` in string literal")),
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
									Tok(Pos(11, 1, 1, 12), lexer.PlusToken),
									PubIdent("bar", Pos(7, 3, 1, 8)),
									Int(lexer.DecIntToken, "2", 13, 1, 1, 14),
								),
							),
							StrCont(" baz ", Pos(15, 5, 1, 16)),
							StrInterp(
								Pos(20, 8, 1, 21),
								PubIdent("fudge", Pos(22, 5, 1, 23)),
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
									Tok(Pos(13, 1, 1, 14), lexer.PlusToken),
									Invalid(VTok(Pos(7, 5, 1, 8), lexer.ErrorToken, "unexpected string literal in string interpolation, only raw strings delimited with `'` can be used in string interpolation")),
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
									Tok(Pos(13, 1, 1, 14), lexer.PlusToken),
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
				&Error{Message: "unexpected PublicIdentifier, expected a statement separator `\\n`, `;` or end of file", Position: Pos(6, 1, 1, 7)},
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
							Tok(Pos(11, 2, 1, 12), lexer.EqualEqualToken),
							Bin(
								Pos(0, 10, 1, 1),
								Tok(Pos(4, 2, 1, 5), lexer.EqualEqualToken),
								PubIdent("bar", Pos(0, 3, 1, 1)),
								PubIdent("baz", Pos(7, 3, 1, 8)),
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
							Tok(Pos(11, 2, 2, 5), lexer.EqualEqualToken),
							Bin(
								Pos(0, 10, 1, 1),
								Tok(Pos(4, 2, 1, 5), lexer.EqualEqualToken),
								PubIdent("bar", Pos(0, 3, 1, 1)),
								PubIdent("baz", Pos(7, 3, 2, 1)),
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
						PubIdent("bar", Pos(0, 3, 1, 1)),
					),
					ExprStmt(
						Pos(4, 7, 2, 1),
						Invalid(Tok(Pos(4, 2, 2, 1), lexer.EqualEqualToken)),
					),
					ExprStmt(
						Pos(11, 4, 3, 1),
						Invalid(Tok(Pos(11, 2, 3, 1), lexer.EqualEqualToken)),
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
							Tok(Pos(30, 3, 1, 31), lexer.RefNotEqualToken),
							Bin(
								Pos(0, 29, 1, 1),
								Tok(Pos(24, 3, 1, 25), lexer.RefEqualToken),
								Bin(
									Pos(0, 23, 1, 1),
									Tok(Pos(18, 3, 1, 19), lexer.StrictNotEqualToken),
									Bin(
										Pos(0, 17, 1, 1),
										Tok(Pos(12, 3, 1, 13), lexer.StrictEqualToken),
										Bin(
											Pos(0, 11, 1, 1),
											Tok(Pos(7, 2, 1, 8), lexer.NotEqualToken),
											Bin(
												Pos(0, 6, 1, 1),
												Tok(Pos(2, 2, 1, 3), lexer.EqualEqualToken),
												PubIdent("a", Pos(0, 1, 1, 1)),
												PubIdent("b", Pos(5, 1, 1, 6)),
											),
											PubIdent("c", Pos(10, 1, 1, 11)),
										),
										PubIdent("d", Pos(16, 1, 1, 17)),
									),
									PubIdent("e", Pos(22, 1, 1, 23)),
								),
								PubIdent("f", Pos(28, 1, 1, 29)),
							),
							PubIdent("g", Pos(34, 1, 1, 35)),
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
							Tok(Pos(4, 2, 1, 5), lexer.AndAndToken),
							PubIdent("foo", Pos(0, 3, 1, 1)),
							Bin(
								Pos(7, 10, 1, 8),
								Tok(Pos(11, 2, 1, 12), lexer.EqualEqualToken),
								PubIdent("bar", Pos(7, 3, 1, 8)),
								PubIdent("baz", Pos(14, 3, 1, 15)),
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
							Tok(Pos(10, 1, 1, 11), lexer.GreaterToken),
							Bin(
								Pos(0, 9, 1, 1),
								Tok(Pos(4, 1, 1, 5), lexer.GreaterToken),
								PubIdent("foo", Pos(0, 3, 1, 1)),
								PubIdent("bar", Pos(6, 3, 1, 7)),
							),
							PubIdent("baz", Pos(12, 3, 1, 13)),
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
							Tok(Pos(10, 1, 2, 5), lexer.GreaterToken),
							Bin(
								Pos(0, 9, 1, 1),
								Tok(Pos(4, 1, 1, 5), lexer.GreaterToken),
								PubIdent("foo", Pos(0, 3, 1, 1)),
								PubIdent("bar", Pos(6, 3, 2, 1)),
							),
							PubIdent("baz", Pos(12, 3, 3, 1)),
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
						PubIdent("bar", Pos(0, 3, 1, 1)),
					),
					ExprStmt(
						Pos(4, 6, 2, 1),
						Invalid(Tok(Pos(4, 1, 2, 1), lexer.GreaterToken)),
					),
					ExprStmt(
						Pos(10, 5, 3, 1),
						Invalid(Tok(Pos(10, 1, 3, 1), lexer.GreaterToken)),
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
							Tok(Pos(42, 3, 1, 43), lexer.SpaceshipOpToken),
							Bin(
								Pos(0, 41, 1, 1),
								Tok(Pos(36, 3, 1, 37), lexer.ReverseInstanceOfToken),
								Bin(
									Pos(0, 35, 1, 1),
									Tok(Pos(30, 3, 1, 31), lexer.InstanceOfToken),
									Bin(
										Pos(0, 29, 1, 1),
										Tok(Pos(25, 2, 1, 26), lexer.ReverseSubtypeToken),
										Bin(
											Pos(0, 24, 1, 1),
											Tok(Pos(20, 2, 1, 21), lexer.SubtypeToken),
											Bin(
												Pos(0, 19, 1, 1),
												Tok(Pos(15, 2, 1, 16), lexer.GreaterEqualToken),
												Bin(
													Pos(0, 14, 1, 1),
													Tok(Pos(11, 1, 1, 12), lexer.GreaterToken),
													Bin(
														Pos(0, 10, 1, 1),
														Tok(Pos(6, 2, 1, 7), lexer.LessEqualToken),
														Bin(
															Pos(0, 5, 1, 1),
															Tok(Pos(2, 1, 1, 3), lexer.LessToken),
															PubIdent("a", Pos(0, 1, 1, 1)),
															PubIdent("b", Pos(4, 1, 1, 5)),
														),
														PubIdent("c", Pos(9, 1, 1, 10)),
													),
													PubIdent("d", Pos(13, 1, 1, 14)),
												),
												PubIdent("e", Pos(18, 1, 1, 19)),
											),
											PubIdent("f", Pos(23, 1, 1, 24)),
										),
										PubIdent("g", Pos(28, 1, 1, 29)),
									),
									PubIdent("h", Pos(34, 1, 1, 35)),
								),
								PubIdent("i", Pos(40, 1, 1, 41)),
							),
							PubIdent("j", Pos(46, 1, 1, 47)),
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
							Tok(Pos(4, 2, 1, 5), lexer.EqualEqualToken),
							PubIdent("foo", Pos(0, 3, 1, 1)),
							Bin(
								Pos(7, 10, 1, 8),
								Tok(Pos(11, 2, 1, 12), lexer.GreaterEqualToken),
								PubIdent("bar", Pos(7, 3, 1, 8)),
								PubIdent("baz", Pos(14, 3, 1, 15)),
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
							Tok(Pos(10, 2, 1, 11), lexer.IfToken),
							Asgmt(
								Pos(0, 9, 1, 1),
								Tok(Pos(4, 1, 1, 5), lexer.EqualToken),
								PubIdent("foo", Pos(0, 3, 1, 1)),
								PubIdent("bar", Pos(6, 3, 1, 7)),
							),
							PubIdent("baz", Pos(13, 3, 1, 14)),
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
								Tok(Pos(4, 1, 1, 5), lexer.EqualToken),
								PubIdent("foo", Pos(0, 3, 1, 1)),
								PubIdent("bar", Pos(6, 3, 1, 7)),
							),
							PubIdent("baz", Pos(13, 3, 1, 14)),
							Asgmt(
								Pos(22, 9, 1, 23),
								Tok(Pos(26, 1, 1, 27), lexer.EqualToken),
								PubIdent("car", Pos(22, 3, 1, 23)),
								PubIdent("red", Pos(28, 3, 1, 29)),
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
							Tok(Pos(4, 2, 1, 5), lexer.IfToken),
							PubIdent("foo", Pos(0, 3, 1, 1)),
							PubIdent("bar", Pos(7, 3, 1, 8)),
						),
					),
					ExprStmt(
						Pos(11, 15, 2, 1),
						Mod(
							Pos(11, 14, 2, 1),
							Tok(Pos(15, 6, 2, 5), lexer.UnlessToken),
							PubIdent("foo", Pos(11, 3, 2, 1)),
							PubIdent("bar", Pos(22, 3, 2, 12)),
						),
					),
					ExprStmt(
						Pos(26, 14, 3, 1),
						Mod(
							Pos(26, 13, 3, 1),
							Tok(Pos(30, 5, 3, 5), lexer.WhileToken),
							PubIdent("foo", Pos(26, 3, 3, 1)),
							PubIdent("bar", Pos(36, 3, 3, 11)),
						),
					),
					ExprStmt(
						Pos(40, 13, 4, 1),
						Mod(
							Pos(40, 13, 4, 1),
							Tok(Pos(44, 5, 4, 5), lexer.UntilToken),
							PubIdent("foo", Pos(40, 3, 4, 1)),
							PubIdent("bar", Pos(50, 3, 4, 11)),
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
							Tok(Pos(10, 2, 1, 11), lexer.IfToken),
							Asgmt(
								Pos(0, 9, 1, 1),
								Tok(Pos(4, 1, 1, 5), lexer.EqualToken),
								PubIdent("foo", Pos(0, 3, 1, 1)),
								PubIdent("bar", Pos(6, 3, 1, 7)),
							),
							PubIdent("baz", Pos(13, 3, 1, 14)),
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
								Tok(Pos(8, 1, 2, 8), lexer.GreaterToken),
								PubIdent("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(13, 9, 3, 2),
									Asgmt(
										Pos(13, 8, 3, 2),
										Tok(Pos(17, 2, 3, 6), lexer.PlusEqualToken),
										PubIdent("foo", Pos(13, 3, 3, 2)),
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
								Tok(Pos(8, 1, 2, 8), lexer.GreaterToken),
								PubIdent("foo", Pos(4, 3, 2, 4)),
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
							Tok(Pos(5, 1, 2, 5), lexer.EqualToken),
							PubIdent("bar", Pos(1, 3, 2, 1)),
							IfExpr(
								Pos(8, 26, 3, 2),
								Bin(
									Pos(11, 7, 3, 5),
									Tok(Pos(15, 1, 3, 9), lexer.GreaterToken),
									PubIdent("foo", Pos(11, 3, 3, 5)),
									Int(lexer.DecIntToken, "0", 17, 1, 3, 11),
								),
								Stmts{
									ExprStmt(
										Pos(21, 9, 4, 3),
										Asgmt(
											Pos(21, 8, 4, 3),
											Tok(Pos(25, 2, 4, 7), lexer.PlusEqualToken),
											PubIdent("foo", Pos(21, 3, 4, 3)),
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
								Tok(Pos(8, 1, 2, 8), lexer.GreaterToken),
								PubIdent("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(17, 8, 2, 17),
									Asgmt(
										Pos(17, 8, 2, 17),
										Tok(Pos(21, 2, 2, 21), lexer.PlusEqualToken),
										PubIdent("foo", Pos(17, 3, 2, 17)),
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
								Tok(Pos(8, 1, 2, 8), lexer.GreaterToken),
								PubIdent("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(13, 9, 3, 2),
									Asgmt(
										Pos(13, 8, 3, 2),
										Tok(Pos(17, 2, 3, 6), lexer.PlusEqualToken),
										PubIdent("foo", Pos(13, 3, 3, 2)),
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
										Tok(Pos(38, 2, 6, 7), lexer.MinusEqualToken),
										PubIdent("foo", Pos(34, 3, 6, 3)),
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
								Tok(Pos(8, 1, 2, 8), lexer.GreaterToken),
								PubIdent("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(17, 8, 2, 17),
									Asgmt(
										Pos(17, 8, 2, 17),
										Tok(Pos(21, 2, 2, 21), lexer.PlusEqualToken),
										PubIdent("foo", Pos(17, 3, 2, 17)),
										Int(lexer.DecIntToken, "2", 24, 1, 2, 24),
									),
								),
							},
							Stmts{
								ExprStmt(
									Pos(31, 8, 3, 6),
									Asgmt(
										Pos(31, 8, 3, 6),
										Tok(Pos(35, 2, 3, 10), lexer.MinusEqualToken),
										PubIdent("foo", Pos(31, 3, 3, 6)),
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
								Tok(Pos(8, 1, 2, 8), lexer.GreaterToken),
								PubIdent("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(17, 8, 2, 17),
									Asgmt(
										Pos(17, 8, 2, 17),
										Tok(Pos(21, 2, 2, 21), lexer.PlusEqualToken),
										PubIdent("foo", Pos(17, 3, 2, 17)),
										Int(lexer.DecIntToken, "2", 24, 1, 2, 24),
									),
								),
							},
							Stmts{
								ExprStmt(
									Pos(31, 8, 3, 6),
									Asgmt(
										Pos(31, 8, 3, 6),
										Tok(Pos(35, 2, 3, 10), lexer.MinusEqualToken),
										PubIdent("foo", Pos(31, 3, 3, 6)),
										Int(lexer.DecIntToken, "2", 38, 1, 3, 13),
									),
								),
							},
						),
					),
					ExprStmt(
						Pos(40, 9, 4, 1),
						Invalid(Tok(Pos(40, 4, 4, 1), lexer.ElseToken)),
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
								Tok(Pos(8, 1, 2, 8), lexer.GreaterToken),
								PubIdent("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(13, 9, 3, 2),
									Asgmt(
										Pos(13, 8, 3, 2),
										Tok(Pos(17, 2, 3, 6), lexer.PlusEqualToken),
										PubIdent("foo", Pos(13, 3, 3, 2)),
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
											Tok(Pos(37, 1, 5, 11), lexer.LessToken),
											PubIdent("foo", Pos(33, 3, 5, 7)),
											Int(lexer.DecIntToken, "5", 39, 1, 5, 13),
										),
										Stmts{
											ExprStmt(
												Pos(42, 10, 6, 2),
												Asgmt(
													Pos(42, 9, 6, 2),
													Tok(Pos(46, 2, 6, 6), lexer.StarEqualToken),
													PubIdent("foo", Pos(42, 3, 6, 2)),
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
														Tok(Pos(62, 1, 7, 11), lexer.LessToken),
														PubIdent("foo", Pos(58, 3, 7, 7)),
														Int(lexer.DecIntToken, "0", 64, 1, 7, 13),
													),
													Stmts{
														ExprStmt(
															Pos(67, 9, 8, 2),
															Asgmt(
																Pos(67, 8, 8, 2),
																Tok(Pos(71, 2, 8, 6), lexer.PercentEqualToken),
																PubIdent("foo", Pos(67, 3, 8, 2)),
																Int(lexer.DecIntToken, "3", 74, 1, 8, 9),
															),
														),
													},
													Stmts{
														ExprStmt(
															Pos(82, 9, 10, 2),
															Asgmt(
																Pos(82, 8, 10, 2),
																Tok(Pos(86, 2, 10, 6), lexer.MinusEqualToken),
																PubIdent("foo", Pos(82, 3, 10, 2)),
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
								Tok(Pos(8, 1, 2, 8), lexer.GreaterToken),
								PubIdent("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(17, 8, 2, 17),
									Asgmt(
										Pos(17, 8, 2, 17),
										Tok(Pos(21, 2, 2, 21), lexer.PlusEqualToken),
										PubIdent("foo", Pos(17, 3, 2, 17)),
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
											Tok(Pos(36, 1, 3, 11), lexer.LessToken),
											PubIdent("foo", Pos(32, 3, 3, 7)),
											Int(lexer.DecIntToken, "5", 38, 1, 3, 13),
										),
										Stmts{
											ExprStmt(
												Pos(45, 9, 3, 20),
												Asgmt(
													Pos(45, 9, 3, 20),
													Tok(Pos(49, 2, 3, 24), lexer.StarEqualToken),
													PubIdent("foo", Pos(45, 3, 3, 20)),
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
														Tok(Pos(65, 1, 4, 11), lexer.LessToken),
														PubIdent("foo", Pos(61, 3, 4, 7)),
														Int(lexer.DecIntToken, "0", 67, 1, 4, 13),
													),
													Stmts{
														ExprStmt(
															Pos(74, 8, 4, 20),
															Asgmt(
																Pos(74, 8, 4, 20),
																Tok(Pos(78, 2, 4, 24), lexer.PercentEqualToken),
																PubIdent("foo", Pos(74, 3, 4, 20)),
																Int(lexer.DecIntToken, "3", 81, 1, 4, 27),
															),
														),
													},
													Stmts{
														ExprStmt(
															Pos(88, 8, 5, 6),
															Asgmt(
																Pos(88, 8, 5, 6),
																Tok(Pos(92, 2, 5, 10), lexer.MinusEqualToken),
																PubIdent("foo", Pos(88, 3, 5, 6)),
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
								Tok(Pos(8, 1, 2, 8), lexer.GreaterToken),
								PubIdent("foo", Pos(4, 3, 2, 4)),
								Int(lexer.DecIntToken, "0", 10, 1, 2, 10),
							),
							Stmts{
								ExprStmt(
									Pos(13, 9, 3, 2),
									Asgmt(
										Pos(13, 8, 3, 2),
										Tok(Pos(17, 2, 3, 6), lexer.PlusEqualToken),
										PubIdent("foo", Pos(13, 3, 3, 2)),
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
											Tok(Pos(39, 1, 5, 13), lexer.LessToken),
											PubIdent("foo", Pos(35, 3, 5, 9)),
											Int(lexer.DecIntToken, "5", 41, 1, 5, 15),
										),
										Stmts{
											ExprStmt(
												Pos(44, 10, 6, 2),
												Asgmt(
													Pos(44, 9, 6, 2),
													Tok(Pos(48, 2, 6, 6), lexer.StarEqualToken),
													PubIdent("foo", Pos(44, 3, 6, 2)),
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
														Tok(Pos(66, 1, 7, 13), lexer.LessToken),
														PubIdent("foo", Pos(62, 3, 7, 9)),
														Int(lexer.DecIntToken, "0", 68, 1, 7, 15),
													),
													Stmts{
														ExprStmt(
															Pos(71, 9, 8, 2),
															Asgmt(
																Pos(71, 8, 8, 2),
																Tok(Pos(75, 2, 8, 6), lexer.PercentEqualToken),
																PubIdent("foo", Pos(71, 3, 8, 2)),
																Int(lexer.DecIntToken, "3", 78, 1, 8, 9),
															),
														),
													},
													Stmts{
														ExprStmt(
															Pos(86, 9, 10, 2),
															Asgmt(
																Pos(86, 8, 10, 2),
																Tok(Pos(90, 2, 10, 6), lexer.MinusEqualToken),
																PubIdent("foo", Pos(86, 3, 10, 2)),
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
								Tok(Pos(12, 1, 2, 12), lexer.GreaterToken),
								PubIdent("foo", Pos(8, 3, 2, 8)),
								Int(lexer.DecIntToken, "0", 14, 1, 2, 14),
							),
							Stmts{
								ExprStmt(
									Pos(17, 9, 3, 2),
									Asgmt(
										Pos(17, 8, 3, 2),
										Tok(Pos(21, 2, 3, 6), lexer.PlusEqualToken),
										PubIdent("foo", Pos(17, 3, 3, 2)),
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
								Tok(Pos(12, 1, 2, 12), lexer.GreaterToken),
								PubIdent("foo", Pos(8, 3, 2, 8)),
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
							Tok(Pos(5, 1, 2, 5), lexer.EqualToken),
							PubIdent("bar", Pos(1, 3, 2, 1)),
							UnlessExpr(
								Pos(8, 30, 3, 2),
								Bin(
									Pos(15, 7, 3, 9),
									Tok(Pos(19, 1, 3, 13), lexer.GreaterToken),
									PubIdent("foo", Pos(15, 3, 3, 9)),
									Int(lexer.DecIntToken, "0", 21, 1, 3, 15),
								),
								Stmts{
									ExprStmt(
										Pos(25, 9, 4, 3),
										Asgmt(
											Pos(25, 8, 4, 3),
											Tok(Pos(29, 2, 4, 7), lexer.PlusEqualToken),
											PubIdent("foo", Pos(25, 3, 4, 3)),
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
								Tok(Pos(12, 1, 2, 12), lexer.GreaterToken),
								PubIdent("foo", Pos(8, 3, 2, 8)),
								Int(lexer.DecIntToken, "0", 14, 1, 2, 14),
							),
							Stmts{
								ExprStmt(
									Pos(21, 8, 2, 21),
									Asgmt(
										Pos(21, 8, 2, 21),
										Tok(Pos(25, 2, 2, 25), lexer.PlusEqualToken),
										PubIdent("foo", Pos(21, 3, 2, 21)),
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
								Tok(Pos(12, 1, 2, 12), lexer.GreaterToken),
								PubIdent("foo", Pos(8, 3, 2, 8)),
								Int(lexer.DecIntToken, "0", 14, 1, 2, 14),
							),
							Stmts{
								ExprStmt(
									Pos(17, 9, 3, 2),
									Asgmt(
										Pos(17, 8, 3, 2),
										Tok(Pos(21, 2, 3, 6), lexer.PlusEqualToken),
										PubIdent("foo", Pos(17, 3, 3, 2)),
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
										Tok(Pos(41, 2, 6, 6), lexer.MinusEqualToken),
										PubIdent("foo", Pos(37, 3, 6, 2)),
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
								Tok(Pos(12, 1, 2, 12), lexer.GreaterToken),
								PubIdent("foo", Pos(8, 3, 2, 8)),
								Int(lexer.DecIntToken, "0", 14, 1, 2, 14),
							),
							Stmts{
								ExprStmt(
									Pos(21, 8, 2, 21),
									Asgmt(
										Pos(21, 8, 2, 21),
										Tok(Pos(25, 2, 2, 25), lexer.PlusEqualToken),
										PubIdent("foo", Pos(21, 3, 2, 21)),
										Int(lexer.DecIntToken, "2", 28, 1, 2, 28),
									),
								),
							},
							Stmts{
								ExprStmt(
									Pos(35, 8, 3, 6),
									Asgmt(
										Pos(35, 8, 3, 6),
										Tok(Pos(39, 2, 3, 10), lexer.MinusEqualToken),
										PubIdent("foo", Pos(35, 3, 3, 6)),
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
								Tok(Pos(12, 1, 2, 12), lexer.GreaterToken),
								PubIdent("foo", Pos(8, 3, 2, 8)),
								Int(lexer.DecIntToken, "0", 14, 1, 2, 14),
							),
							Stmts{
								ExprStmt(
									Pos(21, 8, 2, 21),
									Asgmt(
										Pos(21, 8, 2, 21),
										Tok(Pos(25, 2, 2, 25), lexer.PlusEqualToken),
										PubIdent("foo", Pos(21, 3, 2, 21)),
										Int(lexer.DecIntToken, "2", 28, 1, 2, 28),
									),
								),
							},
							Stmts{
								ExprStmt(
									Pos(35, 8, 3, 6),
									Asgmt(
										Pos(35, 8, 3, 6),
										Tok(Pos(39, 2, 3, 10), lexer.MinusEqualToken),
										PubIdent("foo", Pos(35, 3, 3, 6)),
										Int(lexer.DecIntToken, "2", 42, 1, 3, 13),
									),
								),
							},
						),
					),
					ExprStmt(
						Pos(44, 9, 4, 1),
						Invalid(Tok(Pos(44, 4, 4, 1), lexer.ElseToken)),
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

func TestWhile(t *testing.T) {
	tests := testTable{
		"can have a multiline body": {
			input: `
while foo > 0
	foo += 2
	nil
end
`,
			want: Prog(
				Pos(0, 34, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 33, 2, 1),
						WhileExpr(
							Pos(1, 32, 2, 1),
							Bin(
								Pos(7, 7, 2, 7),
								Tok(Pos(11, 1, 2, 11), lexer.GreaterToken),
								PubIdent("foo", Pos(7, 3, 2, 7)),
								Int(lexer.DecIntToken, "0", 13, 1, 2, 13),
							),
							Stmts{
								ExprStmt(
									Pos(16, 9, 3, 2),
									Asgmt(
										Pos(16, 8, 3, 2),
										Tok(Pos(20, 2, 3, 6), lexer.PlusEqualToken),
										PubIdent("foo", Pos(16, 3, 3, 2)),
										Int(lexer.DecIntToken, "2", 23, 1, 3, 9),
									),
								),
								ExprStmt(
									Pos(26, 4, 4, 2),
									Nil(Pos(26, 3, 4, 2)),
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
			want: Prog(
				Pos(0, 19, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 18, 2, 1),
						WhileExpr(
							Pos(1, 17, 2, 1),
							Bin(
								Pos(7, 7, 2, 7),
								Tok(Pos(11, 1, 2, 11), lexer.GreaterToken),
								PubIdent("foo", Pos(7, 3, 2, 7)),
								Int(lexer.DecIntToken, "0", 13, 1, 2, 13),
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
			want: Prog(
				Pos(0, 42, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 37, 2, 1),
						Asgmt(
							Pos(1, 36, 2, 1),
							Tok(Pos(5, 1, 2, 5), lexer.EqualToken),
							PubIdent("bar", Pos(1, 3, 2, 1)),
							WhileExpr(
								Pos(8, 29, 3, 2),
								Bin(
									Pos(14, 7, 3, 8),
									Tok(Pos(18, 1, 3, 12), lexer.GreaterToken),
									PubIdent("foo", Pos(14, 3, 3, 8)),
									Int(lexer.DecIntToken, "0", 20, 1, 3, 14),
								),
								Stmts{
									ExprStmt(
										Pos(24, 9, 4, 3),
										Asgmt(
											Pos(24, 8, 4, 3),
											Tok(Pos(28, 2, 4, 7), lexer.PlusEqualToken),
											PubIdent("foo", Pos(24, 3, 4, 3)),
											Int(lexer.DecIntToken, "2", 31, 1, 4, 10),
										),
									),
								},
							),
						),
					),
					ExprStmt(
						Pos(38, 4, 6, 1),
						Nil(Pos(38, 3, 6, 1)),
					),
				},
			),
		},
		"can be single line with then and without end": {
			input: `
while foo > 0 then foo += 2
nil
`,
			want: Prog(
				Pos(0, 33, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 28, 2, 1),
						WhileExpr(
							Pos(1, 27, 2, 1),
							Bin(
								Pos(7, 7, 2, 7),
								Tok(Pos(11, 1, 2, 11), lexer.GreaterToken),
								PubIdent("foo", Pos(7, 3, 2, 7)),
								Int(lexer.DecIntToken, "0", 13, 1, 2, 13),
							),
							Stmts{
								ExprStmt(
									Pos(20, 8, 2, 20),
									Asgmt(
										Pos(20, 8, 2, 20),
										Tok(Pos(24, 2, 2, 24), lexer.PlusEqualToken),
										PubIdent("foo", Pos(20, 3, 2, 20)),
										Int(lexer.DecIntToken, "2", 27, 1, 2, 27),
									),
								),
							},
						),
					),
					ExprStmt(
						Pos(29, 4, 3, 1),
						Nil(Pos(29, 3, 3, 1)),
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
			want: Prog(
				Pos(0, 58, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 53, 2, 1),
						WhileExpr(
							Pos(1, 52, 2, 1),
							Bin(
								Pos(7, 7, 2, 7),
								Tok(Pos(11, 1, 2, 11), lexer.GreaterToken),
								PubIdent("foo", Pos(7, 3, 2, 7)),
								Int(lexer.DecIntToken, "0", 13, 1, 2, 13),
							),
							Stmts{
								ExprStmt(
									Pos(16, 9, 3, 2),
									Asgmt(
										Pos(16, 8, 3, 2),
										Tok(Pos(20, 2, 3, 6), lexer.PlusEqualToken),
										PubIdent("foo", Pos(16, 3, 3, 2)),
										Int(lexer.DecIntToken, "2", 23, 1, 3, 9),
									),
								),
								ExprStmt(
									Pos(26, 4, 4, 2),
									Nil(Pos(26, 3, 4, 2)),
								),
								ExprStmt(
									Pos(30, 5, 5, 1),
									Invalid(Tok(Pos(30, 4, 5, 1), lexer.ElseToken)),
								),
								ExprStmt(
									Pos(36, 9, 6, 2),
									Asgmt(
										Pos(36, 8, 6, 2),
										Tok(Pos(40, 2, 6, 6), lexer.MinusEqualToken),
										PubIdent("foo", Pos(36, 3, 6, 2)),
										Int(lexer.DecIntToken, "2", 43, 1, 6, 9),
									),
								),
								ExprStmt(
									Pos(46, 4, 7, 2),
									Nil(Pos(46, 3, 7, 2)),
								),
							},
						),
					),
					ExprStmt(
						Pos(54, 4, 9, 1),
						Nil(Pos(54, 3, 9, 1)),
					),
				},
			),
			err: ErrorList{
				&Error{Message: "unexpected else, expected an expression", Position: Pos(30, 4, 5, 1)},
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
			want: Prog(
				Pos(0, 34, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 33, 2, 1),
						UntilExpr(
							Pos(1, 32, 2, 1),
							Bin(
								Pos(7, 7, 2, 7),
								Tok(Pos(11, 1, 2, 11), lexer.GreaterToken),
								PubIdent("foo", Pos(7, 3, 2, 7)),
								Int(lexer.DecIntToken, "0", 13, 1, 2, 13),
							),
							Stmts{
								ExprStmt(
									Pos(16, 9, 3, 2),
									Asgmt(
										Pos(16, 8, 3, 2),
										Tok(Pos(20, 2, 3, 6), lexer.PlusEqualToken),
										PubIdent("foo", Pos(16, 3, 3, 2)),
										Int(lexer.DecIntToken, "2", 23, 1, 3, 9),
									),
								),
								ExprStmt(
									Pos(26, 4, 4, 2),
									Nil(Pos(26, 3, 4, 2)),
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
			want: Prog(
				Pos(0, 19, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 18, 2, 1),
						UntilExpr(
							Pos(1, 17, 2, 1),
							Bin(
								Pos(7, 7, 2, 7),
								Tok(Pos(11, 1, 2, 11), lexer.GreaterToken),
								PubIdent("foo", Pos(7, 3, 2, 7)),
								Int(lexer.DecIntToken, "0", 13, 1, 2, 13),
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
			want: Prog(
				Pos(0, 42, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 37, 2, 1),
						Asgmt(
							Pos(1, 36, 2, 1),
							Tok(Pos(5, 1, 2, 5), lexer.EqualToken),
							PubIdent("bar", Pos(1, 3, 2, 1)),
							UntilExpr(
								Pos(8, 29, 3, 2),
								Bin(
									Pos(14, 7, 3, 8),
									Tok(Pos(18, 1, 3, 12), lexer.GreaterToken),
									PubIdent("foo", Pos(14, 3, 3, 8)),
									Int(lexer.DecIntToken, "0", 20, 1, 3, 14),
								),
								Stmts{
									ExprStmt(
										Pos(24, 9, 4, 3),
										Asgmt(
											Pos(24, 8, 4, 3),
											Tok(Pos(28, 2, 4, 7), lexer.PlusEqualToken),
											PubIdent("foo", Pos(24, 3, 4, 3)),
											Int(lexer.DecIntToken, "2", 31, 1, 4, 10),
										),
									),
								},
							),
						),
					),
					ExprStmt(
						Pos(38, 4, 6, 1),
						Nil(Pos(38, 3, 6, 1)),
					),
				},
			),
		},
		"can be single line with then and without end": {
			input: `
until foo > 0 then foo += 2
nil
`,
			want: Prog(
				Pos(0, 33, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 28, 2, 1),
						UntilExpr(
							Pos(1, 27, 2, 1),
							Bin(
								Pos(7, 7, 2, 7),
								Tok(Pos(11, 1, 2, 11), lexer.GreaterToken),
								PubIdent("foo", Pos(7, 3, 2, 7)),
								Int(lexer.DecIntToken, "0", 13, 1, 2, 13),
							),
							Stmts{
								ExprStmt(
									Pos(20, 8, 2, 20),
									Asgmt(
										Pos(20, 8, 2, 20),
										Tok(Pos(24, 2, 2, 24), lexer.PlusEqualToken),
										PubIdent("foo", Pos(20, 3, 2, 20)),
										Int(lexer.DecIntToken, "2", 27, 1, 2, 27),
									),
								),
							},
						),
					),
					ExprStmt(
						Pos(29, 4, 3, 1),
						Nil(Pos(29, 3, 3, 1)),
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
			want: Prog(
				Pos(0, 58, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 53, 2, 1),
						UntilExpr(
							Pos(1, 52, 2, 1),
							Bin(
								Pos(7, 7, 2, 7),
								Tok(Pos(11, 1, 2, 11), lexer.GreaterToken),
								PubIdent("foo", Pos(7, 3, 2, 7)),
								Int(lexer.DecIntToken, "0", 13, 1, 2, 13),
							),
							Stmts{
								ExprStmt(
									Pos(16, 9, 3, 2),
									Asgmt(
										Pos(16, 8, 3, 2),
										Tok(Pos(20, 2, 3, 6), lexer.PlusEqualToken),
										PubIdent("foo", Pos(16, 3, 3, 2)),
										Int(lexer.DecIntToken, "2", 23, 1, 3, 9),
									),
								),
								ExprStmt(
									Pos(26, 4, 4, 2),
									Nil(Pos(26, 3, 4, 2)),
								),
								ExprStmt(
									Pos(30, 5, 5, 1),
									Invalid(Tok(Pos(30, 4, 5, 1), lexer.ElseToken)),
								),
								ExprStmt(
									Pos(36, 9, 6, 2),
									Asgmt(
										Pos(36, 8, 6, 2),
										Tok(Pos(40, 2, 6, 6), lexer.MinusEqualToken),
										PubIdent("foo", Pos(36, 3, 6, 2)),
										Int(lexer.DecIntToken, "2", 43, 1, 6, 9),
									),
								),
								ExprStmt(
									Pos(46, 4, 7, 2),
									Nil(Pos(46, 3, 7, 2)),
								),
							},
						),
					),
					ExprStmt(
						Pos(54, 4, 9, 1),
						Nil(Pos(54, 3, 9, 1)),
					),
				},
			),
			err: ErrorList{
				&Error{Message: "unexpected else, expected an expression", Position: Pos(30, 4, 5, 1)},
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
			want: Prog(
				Pos(0, 25, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 24, 2, 1),
						LoopExpr(
							Pos(1, 23, 2, 1),
							Stmts{
								ExprStmt(
									Pos(7, 9, 3, 2),
									Asgmt(
										Pos(7, 8, 3, 2),
										Tok(Pos(11, 2, 3, 6), lexer.PlusEqualToken),
										PubIdent("foo", Pos(7, 3, 3, 2)),
										Int(lexer.DecIntToken, "2", 14, 1, 3, 9),
									),
								),
								ExprStmt(
									Pos(17, 4, 4, 2),
									Nil(Pos(17, 3, 4, 2)),
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
			want: Prog(
				Pos(0, 10, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 9, 2, 1),
						LoopExpr(
							Pos(1, 8, 2, 1),
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
			want: Prog(
				Pos(0, 33, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 28, 2, 1),
						Asgmt(
							Pos(1, 27, 2, 1),
							Tok(Pos(5, 1, 2, 5), lexer.EqualToken),
							PubIdent("bar", Pos(1, 3, 2, 1)),
							LoopExpr(
								Pos(8, 20, 3, 2),
								Stmts{
									ExprStmt(
										Pos(15, 9, 4, 3),
										Asgmt(
											Pos(15, 8, 4, 3),
											Tok(Pos(19, 2, 4, 7), lexer.PlusEqualToken),
											PubIdent("foo", Pos(15, 3, 4, 3)),
											Int(lexer.DecIntToken, "2", 22, 1, 4, 10),
										),
									),
								},
							),
						),
					),
					ExprStmt(
						Pos(29, 4, 6, 1),
						Nil(Pos(29, 3, 6, 1)),
					),
				},
			),
		},
		"can be single line without end": {
			input: `
loop foo += 2
nil
`,
			want: Prog(
				Pos(0, 19, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 14, 2, 1),
						LoopExpr(
							Pos(1, 13, 2, 1),
							Stmts{
								ExprStmt(
									Pos(6, 8, 2, 6),
									Asgmt(
										Pos(6, 8, 2, 6),
										Tok(Pos(10, 2, 2, 10), lexer.PlusEqualToken),
										PubIdent("foo", Pos(6, 3, 2, 6)),
										Int(lexer.DecIntToken, "2", 13, 1, 2, 13),
									),
								),
							},
						),
					),
					ExprStmt(
						Pos(15, 4, 3, 1),
						Nil(Pos(15, 3, 3, 1)),
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
			want: Prog(
				Pos(0, 49, 1, 1),
				Stmts{
					EmptyStmt(Pos(0, 1, 1, 1)),
					ExprStmt(
						Pos(1, 44, 2, 1),
						LoopExpr(
							Pos(1, 43, 2, 1),
							Stmts{
								ExprStmt(
									Pos(7, 9, 3, 2),
									Asgmt(
										Pos(7, 8, 3, 2),
										Tok(Pos(11, 2, 3, 6), lexer.PlusEqualToken),
										PubIdent("foo", Pos(7, 3, 3, 2)),
										Int(lexer.DecIntToken, "2", 14, 1, 3, 9),
									),
								),
								ExprStmt(
									Pos(17, 4, 4, 2),
									Nil(Pos(17, 3, 4, 2)),
								),
								ExprStmt(
									Pos(21, 5, 5, 1),
									Invalid(Tok(Pos(21, 4, 5, 1), lexer.ElseToken)),
								),
								ExprStmt(
									Pos(27, 9, 6, 2),
									Asgmt(
										Pos(27, 8, 6, 2),
										Tok(Pos(31, 2, 6, 6), lexer.MinusEqualToken),
										PubIdent("foo", Pos(27, 3, 6, 2)),
										Int(lexer.DecIntToken, "2", 34, 1, 6, 9),
									),
								),
								ExprStmt(
									Pos(37, 4, 7, 2),
									Nil(Pos(37, 3, 7, 2)),
								),
							},
						),
					),
					ExprStmt(
						Pos(45, 4, 9, 1),
						Nil(Pos(45, 3, 9, 1)),
					),
				},
			),
			err: ErrorList{
				&Error{Message: "unexpected else, expected an expression", Position: Pos(21, 4, 5, 1)},
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
			want: Prog(
				Pos(0, 5, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 5, 1, 1),
						Break(Pos(0, 5, 1, 1)),
					),
				},
			),
		},
		"can't have an argument": {
			input: `break 2`,
			want: Prog(
				Pos(0, 7, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 5, 1, 1),
						Break(Pos(0, 5, 1, 1)),
					),
				},
			),
			err: ErrorList{
				&Error{
					Message:  "unexpected DecInt, expected a statement separator `\\n`, `;` or end of file",
					Position: Pos(6, 1, 1, 7),
				},
			},
		},
		"is an expression": {
			input: `foo && break`,
			want: Prog(
				Pos(0, 12, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 12, 1, 1),
						Logic(
							Pos(0, 12, 1, 1),
							Tok(Pos(4, 2, 1, 5), lexer.AndAndToken),
							PubIdent("foo", Pos(0, 3, 1, 1)),
							Break(Pos(7, 5, 1, 8)),
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
			want: Prog(
				Pos(0, 6, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 6, 1, 1),
						Return(Pos(0, 6, 1, 1), nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "return\n1",
			want: Prog(
				Pos(0, 8, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 7, 1, 1),
						Return(Pos(0, 6, 1, 1), nil),
					),
					ExprStmt(
						Pos(7, 1, 2, 1),
						Int(lexer.DecIntToken, "1", 7, 1, 2, 1),
					),
				},
			),
		},
		"can have an argument": {
			input: `return 2`,
			want: Prog(
				Pos(0, 8, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 8, 1, 1),
						Return(
							Pos(0, 8, 1, 1),
							Int(lexer.DecIntToken, "2", 7, 1, 1, 8),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && return`,
			want: Prog(
				Pos(0, 13, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 13, 1, 1),
						Logic(
							Pos(0, 13, 1, 1),
							Tok(Pos(4, 2, 1, 5), lexer.AndAndToken),
							PubIdent("foo", Pos(0, 3, 1, 1)),
							Return(Pos(7, 6, 1, 8), nil),
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
			want: Prog(
				Pos(0, 8, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 8, 1, 1),
						Continue(Pos(0, 8, 1, 1), nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "continue\n1",
			want: Prog(
				Pos(0, 10, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 9, 1, 1),
						Continue(Pos(0, 8, 1, 1), nil),
					),
					ExprStmt(
						Pos(9, 1, 2, 1),
						Int(lexer.DecIntToken, "1", 9, 1, 2, 1),
					),
				},
			),
		},
		"can have an argument": {
			input: `continue 2`,
			want: Prog(
				Pos(0, 10, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 10, 1, 1),
						Continue(
							Pos(0, 10, 1, 1),
							Int(lexer.DecIntToken, "2", 9, 1, 1, 10),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && continue`,
			want: Prog(
				Pos(0, 15, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 15, 1, 1),
						Logic(
							Pos(0, 15, 1, 1),
							Tok(Pos(4, 2, 1, 5), lexer.AndAndToken),
							PubIdent("foo", Pos(0, 3, 1, 1)),
							Continue(Pos(7, 8, 1, 8), nil),
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
			want: Prog(
				Pos(0, 5, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 5, 1, 1),
						Throw(Pos(0, 5, 1, 1), nil),
					),
				},
			),
		},
		"can stand alone in the middle": {
			input: "throw\n1",
			want: Prog(
				Pos(0, 7, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 6, 1, 1),
						Throw(Pos(0, 5, 1, 1), nil),
					),
					ExprStmt(
						Pos(6, 1, 2, 1),
						Int(lexer.DecIntToken, "1", 6, 1, 2, 1),
					),
				},
			),
		},
		"can have an argument": {
			input: `throw 2`,
			want: Prog(
				Pos(0, 7, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 7, 1, 1),
						Throw(
							Pos(0, 7, 1, 1),
							Int(lexer.DecIntToken, "2", 6, 1, 1, 7),
						),
					),
				},
			),
		},
		"is an expression": {
			input: `foo && throw`,
			want: Prog(
				Pos(0, 12, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 12, 1, 1),
						Logic(
							Pos(0, 12, 1, 1),
							Tok(Pos(4, 2, 1, 5), lexer.AndAndToken),
							PubIdent("foo", Pos(0, 3, 1, 1)),
							Throw(Pos(7, 5, 1, 8), nil),
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
			want: Prog(
				Pos(0, 7, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 7, 1, 1),
						Var(
							Pos(0, 7, 1, 1),
							VTok(Pos(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							nil,
							nil,
						),
					),
				},
			),
		},
		"can have an initialiser without a type": {
			input: "var foo = 5",
			want: Prog(
				Pos(0, 11, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 11, 1, 1),
						Var(
							Pos(0, 11, 1, 1),
							VTok(Pos(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							nil,
							Int(lexer.DecIntToken, "5", 10, 1, 1, 11),
						),
					),
				},
			),
		},
		"can have an initialiser with a type": {
			input: "var foo: Int = 5",
			want: Prog(
				Pos(0, 16, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 16, 1, 1),
						Var(
							Pos(0, 16, 1, 1),
							VTok(Pos(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							PubConst("Int", Pos(9, 3, 1, 10)),
							Int(lexer.DecIntToken, "5", 15, 1, 1, 16),
						),
					),
				},
			),
		},
		"can have a type": {
			input: "var foo: Int",
			want: Prog(
				Pos(0, 12, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 12, 1, 1),
						Var(
							Pos(0, 12, 1, 1),
							VTok(Pos(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							PubConst("Int", Pos(9, 3, 1, 10)),
							nil,
						),
					),
				},
			),
		},
		"can have a nilable type": {
			input: "var foo: Int?",
			want: Prog(
				Pos(0, 13, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 13, 1, 1),
						Var(
							Pos(0, 13, 1, 1),
							VTok(Pos(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							Nilable(
								Pos(9, 4, 1, 10),
								PubConst("Int", Pos(9, 3, 1, 10)),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a union type": {
			input: "var foo: Int | String",
			want: Prog(
				Pos(0, 21, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 21, 1, 1),
						Var(
							Pos(0, 21, 1, 1),
							VTok(Pos(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							BinType(
								Pos(9, 12, 1, 10),
								Tok(Pos(13, 1, 1, 14), lexer.OrToken),
								PubConst("Int", Pos(9, 3, 1, 10)),
								PubConst("String", Pos(15, 6, 1, 16)),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nested union type": {
			input: "var foo: Int | String | Symbol",
			want: Prog(
				Pos(0, 30, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 30, 1, 1),
						Var(
							Pos(0, 30, 1, 1),
							VTok(Pos(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							BinType(
								Pos(9, 21, 1, 10),
								Tok(Pos(22, 1, 1, 23), lexer.OrToken),
								BinType(
									Pos(9, 12, 1, 10),
									Tok(Pos(13, 1, 1, 14), lexer.OrToken),
									PubConst("Int", Pos(9, 3, 1, 10)),
									PubConst("String", Pos(15, 6, 1, 16)),
								),
								PubConst("Symbol", Pos(24, 6, 1, 25)),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nilable union type": {
			input: "var foo: (Int | String)?",
			want: Prog(
				Pos(0, 24, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 24, 1, 1),
						Var(
							Pos(0, 24, 1, 1),
							VTok(Pos(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							Nilable(
								Pos(10, 14, 1, 11),
								BinType(
									Pos(10, 12, 1, 11),
									Tok(Pos(14, 1, 1, 15), lexer.OrToken),
									PubConst("Int", Pos(10, 3, 1, 11)),
									PubConst("String", Pos(16, 6, 1, 17)),
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
			want: Prog(
				Pos(0, 21, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 21, 1, 1),
						Var(
							Pos(0, 21, 1, 1),
							VTok(Pos(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							BinType(
								Pos(9, 12, 1, 10),
								Tok(Pos(13, 1, 1, 14), lexer.AndToken),
								PubConst("Int", Pos(9, 3, 1, 10)),
								PubConst("String", Pos(15, 6, 1, 16)),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nested intersection type": {
			input: "var foo: Int & String & Symbol",
			want: Prog(
				Pos(0, 30, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 30, 1, 1),
						Var(
							Pos(0, 30, 1, 1),
							VTok(Pos(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							BinType(
								Pos(9, 21, 1, 10),
								Tok(Pos(22, 1, 1, 23), lexer.AndToken),
								BinType(
									Pos(9, 12, 1, 10),
									Tok(Pos(13, 1, 1, 14), lexer.AndToken),
									PubConst("Int", Pos(9, 3, 1, 10)),
									PubConst("String", Pos(15, 6, 1, 16)),
								),
								PubConst("Symbol", Pos(24, 6, 1, 25)),
							),
							nil,
						),
					),
				},
			),
		},
		"can have a nilable intersection type": {
			input: "var foo: (Int & String)?",
			want: Prog(
				Pos(0, 24, 1, 1),
				Stmts{
					ExprStmt(
						Pos(0, 24, 1, 1),
						Var(
							Pos(0, 24, 1, 1),
							VTok(Pos(4, 3, 1, 5), lexer.PublicIdentifierToken, "foo"),
							Nilable(
								Pos(10, 14, 1, 11),
								BinType(
									Pos(10, 12, 1, 11),
									Tok(Pos(14, 1, 1, 15), lexer.AndToken),
									PubConst("Int", Pos(10, 3, 1, 11)),
									PubConst("String", Pos(16, 6, 1, 17)),
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

// func TestX(t *testing.T) {
// 	tests := testTable{}

// 	for name, tc := range tests {
// 		t.Run(name, func(t *testing.T) {
// 			parserTest(tc, t)
// 		})
// 	}
// }
