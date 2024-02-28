// Package token implements tokens produced
// by the regex lexer and used by the regex parser
package token

import (
	"fmt"

	"github.com/elk-language/elk/position"
)

// Represents a single token produced by the lexer.
type Token struct {
	Type
	Value string // Literal value of the token, will be empty for tokens with non-dynamic lexemes
	span  *position.Span
}

func (t *Token) Span() *position.Span {
	return t.span
}

func (t *Token) SetSpan(span *position.Span) {
	t.span = span
}

// When the Value field of the token is empty,
// the string will be fetched from a global map.
func (t *Token) StringValue() string {
	if t.Value == "" {
		return t.Type.String()
	}

	return t.Value
}

// Implements the fmt.Stringer interface.
func (t *Token) String() string {
	if len(t.Value) == 0 {
		return t.Type.String()
	}
	return fmt.Sprintf("`%s` (%s)", t.Value, t.Type.String())
}

// Creates a new token.
func New(span *position.Span, tokenType Type) *Token {
	return &Token{
		span: span,
		Type: tokenType,
	}
}

// Creates a new token with the specified value.
func NewWithValue(span *position.Span, tokenType Type, value string) *Token {
	return &Token{
		span:  span,
		Type:  tokenType,
		Value: value,
	}
}
