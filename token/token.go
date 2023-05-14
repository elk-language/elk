package token

import (
	"fmt"

	"github.com/elk-language/elk/position"
)

// Represents a single token produced by the lexer.
type Token struct {
	Type
	Value string // Literal value of the token
	*position.Position
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
	case HEX_INT:
		result = "0x" + t.Value
	case DUO_INT:
		result = "0d" + t.Value
	case OCT_INT:
		result = "0o" + t.Value
	case QUAT_INT:
		result = "0q" + t.Value
	case BIN_INT:
		result = "0b" + t.Value
	default:
		result = t.Value
	}

	if maxInspectLen < len(result) {
		return result[0:maxInspectLen] + "..."
	}

	return result
}

// Creates a new token.
func New(pos *position.Position, tokenType Type) *Token {
	return &Token{
		Position: pos,
		Type:     tokenType,
	}
}

// Creates a new token with the specified value.
func NewWithValue(pos *position.Position, tokenType Type, value string) *Token {
	return &Token{
		Position: pos,
		Type:     tokenType,
		Value:    value,
	}
}
