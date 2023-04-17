package lexer

import (
	"fmt"
)

// Represents a single token produced by the lexer.
type Token struct {
	TokenType
	Value string // Literal value of the token
	Position
}

// Implements the fmt.Stringer interface.
func (t *Token) String() string {
	if len(t.Value) == 0 {
		return t.TokenType.String()
	}
	return fmt.Sprintf("`%s` (%s)", t.InspectValue(), t.TokenType.String())
}

const maxInspectLen = 20

// Returns a shortened version of the value
// which resembles source code.
func (t *Token) InspectValue() string {
	var result string

	switch t.TokenType {
	case InstanceVariableToken:
		result = "@" + t.Value
	case RawStringToken:
		result = "'" + t.Value + "'"
	case HexIntToken:
		result = "0x" + t.Value
	case DuoIntToken:
		result = "0d" + t.Value
	case OctIntToken:
		result = "0o" + t.Value
	case QuatIntToken:
		result = "0q" + t.Value
	case BinIntToken:
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
func NewToken(tokenType TokenType, startByte int, byteLength int, line int, column int) *Token {
	return &Token{
		TokenType: tokenType,
		Position: Position{
			StartByte:  startByte,
			ByteLength: byteLength,
			Line:       line,
			Column:     column,
		},
	}
}

// Creates a new token with the specified value.
func NewTokenWithValue(tokenType TokenType, value string, startByte int, byteLength int, line int, column int) *Token {
	return &Token{
		TokenType: tokenType,
		Value:     value,
		Position: Position{
			StartByte:  startByte,
			ByteLength: byteLength,
			Line:       line,
			Column:     column,
		},
	}
}
