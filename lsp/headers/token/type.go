package token

// Represents the type of token
type Type uint8

// Name of the token.
func (t Type) String() string {
	if int(t) > len(tokenNames) {
		return "UNKNOWN"
	}

	return tokenNames[t]
}

const (
	ZERO_VALUE  Type = iota // Zero value for Type
	ERROR                   // ERROR Token with a message
	END_OF_FILE             // End Of File has been reached
	COLON                   // Colon `:`
	STRING                  // String
	NEWLINE                 // Two char newline `\r\n`
)

var tokenNames = [...]string{
	ERROR:       "ERROR",
	END_OF_FILE: "END_OF_FILE",
	COLON:       ":",
	STRING:      "STRING",
	NEWLINE:     `\r\n`,
}
