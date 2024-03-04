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
	ZERO_VALUE                           Type = iota // Zero value for Type
	ERROR                                            // ERROR Token with a message
	END_OF_FILE                                      // End Of File has been reached
	CHAR                                             // A character
	META_CHAR_ESCAPE                                 // A meta-character escape `\.`, `\+`
	QUOTED_TEXT                                      // Quoted text between \Q and \E
	DOT                                              // Dot `.`
	SINGLE_QUOTE                                     // Apostrophe `'`
	DASH                                             // Dash `-`
	COLON                                            // Colon `:`
	COMMA                                            // Comma `,`
	LANGLE                                           // Left angle bracket `<`
	RANGLE                                           // Right angle bracket `>`
	LPAREN                                           // Left parenthesis `(`
	RPAREN                                           // Right parenthesis `)`
	LBRACE                                           // Left brace `{`
	RBRACE                                           // Right brace `}`
	LBRACKET                                         // Left bracket `[`
	RBRACKET                                         // Right bracket `]`
	PIPE                                             // Pipe `|`
	STAR                                             // Star `*`
	PLUS                                             // Star `+`
	QUESTION                                         // Question mark `?`
	CARET                                            // Caret `^`
	DOLLAR                                           // Dollar `$`
	CARET_ESCAPE                                     // Caret escape `\c`
	UNICODE_ESCAPE                                   // Unicode escape `\u`
	HEX_ESCAPE                                       // Hex escape `\x`
	OCTAL_ESCAPE                                     // Octal escape `\o`
	SIMPLE_OCTAL_ESCAPE                              // Simple octal escape `\123`
	BELL_ESCAPE                                      // Bell escape `\a`
	FORM_FEED_ESCAPE                                 // Form feed escape `\f`
	TAB_ESCAPE                                       // Tab escape `\t`
	NEWLINE_ESCAPE                                   // Newline escape `\n`
	CARRIAGE_RETURN_ESCAPE                           // Carriage return escape `\r`
	ABSOLUTE_START_OF_STRING_ANCHOR                  // Beginning of text anchor `\A`
	ABSOLUTE_END_OF_STRING_ANCHOR                    // End of text anchor `\z`
	WORD_BOUNDARY_ANCHOR                             // Word boundary anchor `\b`
	NOT_WORD_BOUNDARY_ANCHOR                         // Not word boundary anchor `\B`
	UNICODE_CHAR_CLASS                               // Unicode char class `\p`
	NEGATED_UNICODE_CHAR_CLASS                       // Negated unicode char class `\P`
	WORD_CHAR_CLASS                                  // Word char class `\w`
	NOT_WORD_CHAR_CLASS                              // Not word char class `\W`
	DIGIT_CHAR_CLASS                                 // Digit char class `\d`
	NOT_DIGIT_CHAR_CLASS                             // Not digit char class `\D`
	WHITESPACE_CHAR_CLASS                            // Whitespace char class `\s`
	NOT_WHITESPACE_CHAR_CLASS                        // Not whitespace char class `\S`
	HORIZONTAL_WHITESPACE_CHAR_CLASS                 // Horizontal whitespace char class `\h`
	NOT_HORIZONTAL_WHITESPACE_CHAR_CLASS             // Not horizontal whitespace char class `\H`
	VERTICAL_WHITESPACE_CHAR_CLASS                   // Vertical whitespace char class `\v`
	NOT_VERTICAL_WHITESPACE_CHAR_CLASS               // Not vertical whitespace char class `\V`
)

var tokenNames = [...]string{
	ERROR:                                "ERROR",
	END_OF_FILE:                          "END_OF_FILE",
	META_CHAR_ESCAPE:                     "META_CHAR_ESCAPE",
	CHAR:                                 "CHAR",
	QUOTED_TEXT:                          "QUOTED_TEXT",
	DOT:                                  ".",
	SINGLE_QUOTE:                         "'",
	DASH:                                 "-",
	COLON:                                ":",
	COMMA:                                ",",
	LANGLE:                               "<",
	RANGLE:                               ">",
	LPAREN:                               "(",
	RPAREN:                               ")",
	LBRACE:                               "{",
	RBRACE:                               "}",
	LBRACKET:                             "[",
	RBRACKET:                             "]",
	PIPE:                                 "|",
	STAR:                                 "*",
	PLUS:                                 "+",
	QUESTION:                             "?",
	CARET:                                "^",
	DOLLAR:                               "$",
	CARET_ESCAPE:                         `\c`,
	UNICODE_ESCAPE:                       `\u`,
	HEX_ESCAPE:                           `\x`,
	OCTAL_ESCAPE:                         `\o`,
	SIMPLE_OCTAL_ESCAPE:                  `SIMPLE_OCTAL_ESCAPE`,
	BELL_ESCAPE:                          `\a`,
	FORM_FEED_ESCAPE:                     `\f`,
	TAB_ESCAPE:                           `\t`,
	NEWLINE_ESCAPE:                       `\n`,
	CARRIAGE_RETURN_ESCAPE:               `\r`,
	ABSOLUTE_START_OF_STRING_ANCHOR:      `\A`,
	ABSOLUTE_END_OF_STRING_ANCHOR:        `\z`,
	WORD_BOUNDARY_ANCHOR:                 `\b`,
	NOT_WORD_BOUNDARY_ANCHOR:             `\B`,
	UNICODE_CHAR_CLASS:                   `\p`,
	NEGATED_UNICODE_CHAR_CLASS:           `\P`,
	WORD_CHAR_CLASS:                      `\w`,
	NOT_WORD_CHAR_CLASS:                  `\W`,
	DIGIT_CHAR_CLASS:                     `\d`,
	NOT_DIGIT_CHAR_CLASS:                 `\D`,
	WHITESPACE_CHAR_CLASS:                `\s`,
	NOT_WHITESPACE_CHAR_CLASS:            `\S`,
	HORIZONTAL_WHITESPACE_CHAR_CLASS:     `\h`,
	NOT_HORIZONTAL_WHITESPACE_CHAR_CLASS: `\H`,
	VERTICAL_WHITESPACE_CHAR_CLASS:       `\v`,
	NOT_VERTICAL_WHITESPACE_CHAR_CLASS:   `\V`,
}
