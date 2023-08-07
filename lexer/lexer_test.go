package lexer

import (
	"testing"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
	"github.com/google/go-cmp/cmp"
)

// Represents a single lexer test case.
type testCase struct {
	input string
	want  []*token.Token
}

// Type of the lexer test table.
type testTable map[string]testCase

// Create a new token in tests
var T = token.New

// Create a new token with value in tests
var V = token.NewWithValue

// Create a new position in tests
var P = position.New

// Create a new span in tests
var S = position.NewSpan

// Function which powers all lexer tests.
// Inspects if the produced stream of tokens
// matches the expected one.
func tokenTest(tc testCase, t *testing.T) {
	t.Helper()
	lex := New(tc.input)
	var got []*token.Token
	for {
		tok := lex.Next()
		if tok.IsEndOfFile() {
			break
		}
		got = append(got, tok)
	}
	diff := cmp.Diff(tc.want, got)
	if diff != "" {
		t.Fatal(diff)
	}
}
