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
		"int addition": {
			input: "1 + 2",
			want: &ast.ProgramNode{
				Position: P(0, 5, 1, 1),
				Body: []ast.StatementNode{
					&ast.ExpressionStatementNode{
						Position: P(0, 5, 1, 1),
						Expression: &ast.BinaryExpressionNode{
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
