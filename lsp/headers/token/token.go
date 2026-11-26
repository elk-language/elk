package token

import (
	"fmt"
	"unicode/utf8"
)

// Represents a single token produced by the lexer.
type Token struct {
	Type
	Value string // Literal value of the token, will be empty for tokens with non-dynamic lexemes
}

// When the Value field of the token is empty,
// the string will be fetched from a global map.
func (t *Token) StringValue() string {
	if t.Value == "" {
		return t.Type.String()
	}

	return t.Value
}

func (t *Token) Char() rune {
	char, _ := utf8.DecodeRuneInString(t.StringValue())
	return char
}

// Implements the fmt.Stringer interface.
func (t *Token) String() string {
	if len(t.Value) == 0 {
		return t.Type.String()
	}
	return fmt.Sprintf("`%s` (%s)", t.Value, t.Type.String())
}

// Creates a new token.
func New(tokenType Type) *Token {
	return &Token{
		Type: tokenType,
	}
}

// Creates a new token with the specified value.
func NewWithValue(tokenType Type, value string) *Token {
	return &Token{
		Type:  tokenType,
		Value: value,
	}
}
