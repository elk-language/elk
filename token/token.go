package token

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/value"
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

func (*Token) Class() *value.Class {
	return value.ElkTokenClass
}

func (*Token) DirectClass() *value.Class {
	return value.ElkTokenClass
}

func (*Token) SingletonClass() *value.Class {
	return nil
}

func (t *Token) Copy() value.Reference {
	return t
}

func (t *Token) InstanceVariables() value.SymbolMap {
	return nil
}

func (t *Token) Inspect() string {
	return fmt.Sprintf(
		"Std::Token{&: %p, value: %s, span: %s}",
		t,
		value.String(t.Value).Inspect(),
		(*value.Span)(t.Span()).Inspect(),
	)
}

func (t *Token) Error() string {
	return t.Inspect()
}

func (t *Token) Equal(other *Token) bool {
	return t.Type == other.Type &&
		t.Value == other.Value &&
		t.span.Equal(other.span)
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
	case STRING_BEG, STRING_CONTENT, STRING_END, RAW_STRING, REGEX_CONTENT:
		return prompt.Yellow
	case STRING_INTERP_BEG, STRING_INTERP_END, REGEX_BEG, REGEX_END,
		REGEX_FLAG_i, REGEX_FLAG_m, REGEX_FLAG_s, REGEX_FLAG_U, REGEX_FLAG_a, REGEX_FLAG_x,
		STRING_INTERP_LOCAL, STRING_INTERP_CONSTANT, STRING_INSPECT_INTERP_BEG,
		STRING_INSPECT_INTERP_CONSTANT, STRING_INSPECT_INTERP_LOCAL:
		return prompt.Red
	case ERROR:
		return prompt.Black
	case HASH_SET_LITERAL_BEG, TUPLE_LITERAL_BEG, RECORD_LITERAL_BEG:
		return prompt.Fuchsia
	}

	if t.IsIntLiteral() {
		return prompt.Blue
	}

	if t.IsFloatLiteral() {
		return prompt.Purple
	}

	if t.IsOperator() || t.IsSpecialCollectionLiteral() {
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
	case STRING_BEG, STRING_CONTENT, STRING_END, RAW_STRING, REGEX_CONTENT:
		return []color.Attribute{color.FgHiYellow}
	case STRING_INTERP_BEG, STRING_INTERP_END, REGEX_BEG, REGEX_END,
		REGEX_FLAG_i, REGEX_FLAG_m, REGEX_FLAG_s, REGEX_FLAG_U, REGEX_FLAG_a, REGEX_FLAG_x,
		STRING_INTERP_LOCAL, STRING_INTERP_CONSTANT, STRING_INSPECT_INTERP_BEG,
		STRING_INSPECT_INTERP_CONSTANT, STRING_INSPECT_INTERP_LOCAL:
		return []color.Attribute{color.FgHiRed}
	case ERROR:
		return []color.Attribute{color.FgBlack, color.BgRed}
	case HASH_SET_LITERAL_BEG, TUPLE_LITERAL_BEG, RECORD_LITERAL_BEG:
		return []color.Attribute{color.FgHiMagenta}
	}

	if t.IsIntLiteral() {
		return []color.Attribute{color.FgHiBlue}
	}

	if t.IsFloatLiteral() {
		return []color.Attribute{color.FgMagenta}
	}

	if t.IsOperator() || t.IsSpecialCollectionLiteral() {
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
func (t *Token) FetchValue() string {
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
	var result strings.Builder

	switch t.Type {
	case INSTANCE_VARIABLE:
		result.WriteRune('@')
		result.WriteString(t.Value)
	case RAW_STRING:
		result.WriteRune('\'')
		result.WriteString(t.Value)
		result.WriteRune('\'')
	case CHAR_LITERAL:
		result.WriteRune('`')
		result.WriteString(t.Value)
		result.WriteRune('`')
	case RAW_CHAR_LITERAL:
		result.WriteString("r`")
		result.WriteString(t.Value)
		result.WriteRune('`')
	default:
		result.WriteString(t.Value)
	}

	if maxInspectLen < result.Len() {
		str := result.String()[0:maxInspectLen]
		result.Reset()
		result.WriteString(str)
		result.WriteString(`...`)
	}

	return result.String()
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
