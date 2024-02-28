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
	ZERO_VALUE                        Type = iota // Zero value for Type
	ERROR                                         // ERROR Token with a message
	END_OF_FILE                                   // End Of File has been reached
	CHAR_LIST                                     // A list of characters inside of []
	CHAR                                          // A character
	QUOTED_TEXT                                   // Quoted text between \Q and \E
	POSIX_CLASS                                   // POSIX
	LANGLE                                        // Left angle bracket `<`
	RANGLE                                        // Right angle bracket `>`
	DOT                                           // Dot `.`
	DASH                                          // Dash `-`
	COLON                                         // Colon `:`
	COMMA                                         // Comma `,`
	LPAREN                                        // Left parenthesis `(`
	RPAREN                                        // Right parenthesis `)`
	LBRACE                                        // Left brace `{`
	RBRACE                                        // Right brace `}`
	LBRACKET                                      // Left bracket `[`
	RBRACKET                                      // Right bracket `]`
	PIPE                                          // Pipe `|`
	STAR                                          // Star `*`
	STAR_QUESTION                                 // Star question mark `*?`
	PLUS                                          // Star `+`
	PLUS_QUESTION                                 // Star question mark `+?`
	QUESTION                                      // Question mark `?`
	QUESTION_QUESTION                             // Double question mark `??`
	CARET                                         // Caret `^`
	DOLLAR                                        // Dollar `$`
	HEX_ESCAPE                                    // Hex escape `\x`
	BELL_ESCAPE                                   // Bell escape `\a`
	FORM_FEED_ESCAPE                              // Form feed escape `\f`
	TAB_ESCAPE                                    // Tab escape `\t`
	NEWLINE_ESCAPE                                // Newline escape `\n`
	CARRIAGE_RETURN_ESCAPE                        // Carriage return escape `\r`
	VERTICAL_TAB_ESCAPE                           // Vertical tab escape `\v`
	UNICODE_CHAR_CLASS_ESCAPE                     // Unicode char class escape `\p`
	NEGATED_UNICODE_CHAR_CLASS_ESCAPE             // Negated unicode char class escape `\P`
	BEGINNING_OF_TEXT_ESCAPE                      // Beginning of text escape `\A`
	END_OF_TEXT_ESCAPE                            // End of text escape `\z`
	WORD_BOUNDARY_ESCAPE                          // Word boundary escape `\b`
	NOT_WORD_BOUNDARY_ESCAPE                      // Not word boundary escape `\B`
	WORD_ESCAPE                                   // Word escape `\w`
	NOT_WORD_ESCAPE                               // Not word escape `\W`
	DIGIT_ESCAPE                                  // Digit escape `\d`
	NOT_DIGIT_ESCAPE                              // Not digit escape `\D`
	WHITESPACE_ESCAPE                             // Whitespace escape `\s`
	NOT_WHITESPACE_ESCAPE                         // Not whitespace escape `\S`
)

var tokenNames = [...]string{
	ERROR:                             "ERROR",
	END_OF_FILE:                       "END_OF_FILE",
	CHAR_LIST:                         "CHAR_LIST",
	CHAR:                              "CHAR",
	QUOTED_TEXT:                       "QUOTED_TEXT",
	LANGLE:                            "<",
	RANGLE:                            ">",
	DOT:                               ".",
	DASH:                              "-",
	COLON:                             ":",
	COMMA:                             ",",
	LPAREN:                            "(",
	RPAREN:                            ")",
	LBRACE:                            "{",
	RBRACE:                            "}",
	LBRACKET:                          "[",
	RBRACKET:                          "]",
	PIPE:                              "|",
	STAR:                              "*",
	STAR_QUESTION:                     "*?",
	PLUS:                              "+",
	PLUS_QUESTION:                     "+?",
	QUESTION:                          "?",
	QUESTION_QUESTION:                 "??",
	CARET:                             "^",
	DOLLAR:                            "$",
	HEX_ESCAPE:                        `\x`,
	BELL_ESCAPE:                       `\a`,
	FORM_FEED_ESCAPE:                  `\f`,
	TAB_ESCAPE:                        `\t`,
	NEWLINE_ESCAPE:                    `\n`,
	CARRIAGE_RETURN_ESCAPE:            `\r`,
	VERTICAL_TAB_ESCAPE:               `\v`,
	UNICODE_CHAR_CLASS_ESCAPE:         `\p`,
	NEGATED_UNICODE_CHAR_CLASS_ESCAPE: `\P`,
	BEGINNING_OF_TEXT_ESCAPE:          `\A`,
	END_OF_TEXT_ESCAPE:                `\z`,
	WORD_BOUNDARY_ESCAPE:              `\b`,
	NOT_WORD_BOUNDARY_ESCAPE:          `\B`,
	WORD_ESCAPE:                       `\w`,
	NOT_WORD_ESCAPE:                   `\W`,
	DIGIT_ESCAPE:                      `\d`,
	NOT_DIGIT_ESCAPE:                  `\D`,
	WHITESPACE_ESCAPE:                 `\s`,
	NOT_WHITESPACE_ESCAPE:             `\S`,
}
