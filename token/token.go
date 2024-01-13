package token

import (
	"fmt"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/go-prompt"
	pstrings "github.com/elk-language/go-prompt/strings"
	"github.com/fatih/color"
)

// Represents a single token produced by the lexer.
type Token struct {
	Type
	Value string // Literal value of the token, will be empty for tokens with non-dynamic lexemes
	span  *position.Span
}

// Index of the first byte of the lexeme.
// Used by go-prompt.
func (t *Token) FirstByteIndex() pstrings.ByteNumber {
	return pstrings.ByteNumber(t.span.StartPos.ByteOffset)
}

// Index of the last byte of the lexeme.
// Used by go-prompt.
func (t *Token) LastByteIndex() pstrings.ByteNumber {
	return pstrings.ByteNumber(t.span.EndPos.ByteOffset)
}

// Text color for go-prompt.
func (t *Token) Color() prompt.Color {
	switch t.Type {
	case INSTANCE_VARIABLE:
		return prompt.Blue
	case PRIVATE_IDENTIFIER:
		return prompt.LightGray
	case PRIVATE_CONSTANT, PUBLIC_CONSTANT:
		return prompt.Turquoise
	case CHAR_LITERAL, RAW_CHAR_LITERAL:
		return prompt.Brown
	case STRING_BEG, STRING_CONTENT, STRING_END, RAW_STRING:
		return prompt.Yellow
	case STRING_INTERP_BEG, STRING_INTERP_END:
		return prompt.Red
	case ERROR:
		return prompt.Black
	}

	if t.IsIntLiteral() {
		return prompt.Blue
	}

	if t.IsFloatLiteral() {
		return prompt.Purple
	}

	if t.IsOperator() {
		return prompt.Fuchsia
	}

	if t.IsKeyword() {
		return prompt.DarkGreen
	}

	return prompt.DefaultColor
}

// Background color for go-prompt.
func (t *Token) BackgroundColor() prompt.Color {
	switch t.Type {
	case ERROR:
		return prompt.DarkRed
	default:
		return prompt.DefaultColor
	}
}

// Display attributes for go-prompt eg. bold, italic, underline.
func (t *Token) DisplayAttributes() []prompt.DisplayAttribute {
	switch t.Type {
	case PUBLIC_CONSTANT, PRIVATE_CONSTANT:
		return []prompt.DisplayAttribute{prompt.DisplayItalic}
	}

	return nil
}

// Returns the ANSI font styling for the github.com/fatih/color package.
func (t *Token) AnsiStyling() []color.Attribute {
	switch t.Type {
	case INSTANCE_VARIABLE:
		return []color.Attribute{color.FgBlue}
	case PRIVATE_IDENTIFIER:
		return []color.Attribute{color.FgHiBlack}
	case PRIVATE_CONSTANT, PUBLIC_CONSTANT:
		return []color.Attribute{color.FgHiCyan, color.Italic}
	case CHAR_LITERAL, RAW_CHAR_LITERAL:
		return []color.Attribute{color.FgYellow}
	case STRING_BEG, STRING_CONTENT, STRING_END, RAW_STRING:
		return []color.Attribute{color.FgHiYellow}
	case STRING_INTERP_BEG, STRING_INTERP_END:
		return []color.Attribute{color.FgHiRed}
	case ERROR:
		return []color.Attribute{color.FgBlack, color.BgRed}
	}

	if t.IsIntLiteral() {
		return []color.Attribute{color.FgHiBlue}
	}

	if t.IsFloatLiteral() {
		return []color.Attribute{color.FgMagenta}
	}

	if t.IsOperator() {
		return []color.Attribute{color.FgHiMagenta}
	}

	if t.IsKeyword() {
		return []color.Attribute{color.FgGreen}
	}

	return nil
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
	return fmt.Sprintf("`%s` (%s)", t.InspectValue(), t.Type.String())
}

const maxInspectLen = 20

// Returns a shortened version of the value
// which resembles source code.
func (t *Token) InspectValue() string {
	var result string

	switch t.Type {
	case INSTANCE_VARIABLE:
		result = "@" + t.Value
	case RAW_STRING:
		result = "'" + t.Value + "'"
	case CHAR_LITERAL:
		result = "`" + t.Value + "`"
	case RAW_CHAR_LITERAL:
		result = "r`" + t.Value + "`"
	default:
		result = t.Value
	}

	if maxInspectLen < len(result) {
		return result[0:maxInspectLen] + "..."
	}

	return result
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
