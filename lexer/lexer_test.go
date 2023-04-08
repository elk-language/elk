package lexer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Represents a single lexer test case.
type testCase struct {
	input string
	want  []*Token
}

// Type of the lexer test table.
type testTable map[string]testCase

// Function which powers all lexer tests
// which inspects if the produced stream of tokens
// matches the expected one.
func tokenTest(tc testCase, t *testing.T) {
	lex := New([]byte(tc.input))
	var got []*Token
	for {
		tok := lex.Next()
		if tok.IsEOF() {
			break
		}
		got = append(got, tok)
	}
	diff := cmp.Diff(tc.want, got)
	if diff != "" {
		t.Fatalf(diff)
	}
}
