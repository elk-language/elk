// Package lexer implements a regex lexer
package lexer

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/regex/token"
)

// Lexing mode which changes how characters are handled by the lexer.
type mode uint8

const (
	normalMode    mode = iota // Initial mode
	charClassMode             // triggered when in a character class (after `[`)
)

// Holds the current state of the lexing process.
type Lexer struct {
	// Regex source code.
	source string
	// Holds the index of the beginning byte
	// of the currently scanned token.
	start int
	// Holds the index of the current byte
	// the lexer is at.
	cursor int
	// Column of the first character of the currently analysed token.
	startColumn int
	// Column of the current character of the currently analysed token.
	column int
	// First line number of the currently analysed token.
	startLine int
	// Current line of source code being analysed.
	line int
	// whether the `x` flag is enabled
	// ignores all literal whitespace and allows for comments
	extendedSyntax bool
	// Current lexing mode.
	mode mode
}

// Instantiates a new lexer for the given regex.
func New(source string) *Lexer {
	return &Lexer{
		source:      source,
		line:        1,
		startLine:   1,
		column:      1,
		startColumn: 1,
		mode:        normalMode,
	}
}

// Lex the given string and return a slice containing all the tokens.
func Lex(source string) []*token.Token {
	l := New(source)

	var tokens []*token.Token
	for {
		tok := l.Next()
		if tok.Type == token.END_OF_FILE {
			break
		}
		tokens = append(tokens, tok)
	}
	return tokens
}

// Returns true if there is any code left to analyse.
func (l *Lexer) hasMoreTokens() bool {
	return l.cursor < len(l.source)
}

// Returns the next token or an error if
// the input is malformed.
func (l *Lexer) Next() *token.Token {
	if !l.hasMoreTokens() {
		return l.token(token.END_OF_FILE)
	}

	return l.scanToken()
}

// Gets the next UTF-8 encoded character
// and increments the cursor.
func (l *Lexer) advanceChar() (rune, bool) {
	if !l.hasMoreTokens() {
		return 0, false
	}

	char, size := l.nextChar()

	l.cursor += size
	l.column += 1
	return char, true
}

// Advance the next `n` characters
func (l *Lexer) advanceChars(n int) bool {
	for i := 0; i < n; i++ {
		_, ok := l.advanceChar()
		if !ok {
			return false
		}
	}

	return true
}

// Rewinds the cursor back to the previous char.
func (l *Lexer) backupChar() {
	l.cursor -= 1
	l.column -= 1
}

// Rewinds the cursor back n chars.
func (l *Lexer) backupChars(n int) {
	l.cursor -= n
	l.column -= n
}

// Swallows characters until the given char is seen.
func (l *Lexer) swallowUntil(char rune) bool {
	for {
		ch, ok := l.advanceChar()
		if !ok {
			return false
		}
		if ch == '\n' {
			l.incrementLine()
		}
		if ch == char {
			break
		}
	}

	return true
}

// Checks if the given character matches
// the next UTF-8 encoded character in source code.
// If they match, the cursor gets incremented.
func (l *Lexer) matchChar(char rune) bool {
	if !l.hasMoreTokens() {
		return false
	}

	if l.peekChar() == char {
		l.advanceChar()
		return true
	}

	return false
}

// Checks if the next `n` chars match the given char.
func (l *Lexer) matchCharN(char rune, n int) bool {
	for i := 0; i < n; i++ {
		if !l.matchChar(char) {
			return false
		}
	}

	return true
}

// Same as [matchChars] but returns the consumed char.
func (l *Lexer) matchCharsRune(validChars string) (bool, rune) {
	if !l.hasMoreTokens() {
		return false, 0
	}

	if strings.ContainsRune(validChars, l.peekChar()) {
		char, _ := l.advanceChar()
		return true, char
	}

	return false, 0
}

// Consumes the next character if it's from the valid set.
func (l *Lexer) matchChars(validChars string) bool {
	if !l.hasMoreTokens() {
		return false
	}

	if strings.ContainsRune(validChars, l.peekChar()) {
		l.advanceChar()
		return true
	}

	return false
}

// Consumes the next `n` characters if their from the valid set.
func (l *Lexer) matchCharsN(validChars string, n int) bool {
	for i := 0; i < n; i++ {
		if !l.matchChars(validChars) {
			return false
		}
	}

	return true
}

// Checks if the next `n` characters are from the valid set.
func (l *Lexer) acceptCharsN(validChars string, n int) bool {
	i := 0
	result := true
	for ; i < n; i++ {
		if !l.matchChars(validChars) {
			result = false
			break
		}
	}

	l.backupChars(i)

	return result
}

// Checks if the next character is from the valid set.
func (l *Lexer) acceptChars(validChars string) bool {
	if !l.hasMoreTokens() {
		return false
	}

	return strings.ContainsRune(validChars, l.peekChar())
}

// Checks if the next character matches the given one.
func (l *Lexer) acceptChar(char rune) bool {
	if !l.hasMoreTokens() {
		return false
	}

	if l.peekChar() == char {
		return true
	}

	return false
}

// Returns the next character and its length in bytes.
func (l *Lexer) nextChar() (rune, int) {
	return utf8.DecodeRuneInString(l.source[l.cursor:])
}

// Returns the second next character and its length in bytes.
func (l *Lexer) nextNextChar() (rune, int) {
	if !l.hasMoreTokens() {
		return '\x00', 0
	}
	return utf8.DecodeRuneInString(l.source[l.cursor+1:])
}

// Gets the next UTF-8 encoded character
// without incrementing the cursor.
func (l *Lexer) peekChar() rune {
	if !l.hasMoreTokens() {
		return '\x00'
	}
	char, _ := l.nextChar()
	return char
}

// Gets the second UTF-8 encoded character
// without incrementing the cursor.
func (l *Lexer) peekNextChar() rune {
	char, _ := l.nextNextChar()
	return char
}

// Skips the current byte.
func (l *Lexer) skipByte() {
	l.start += 1
	l.startColumn += 1
}

// Skips the current accumulated token.
func (l *Lexer) skipToken() {
	l.start = l.cursor
	l.startColumn = l.column
	l.startLine = l.line
}

// Swallow consecutive newlines and wrap them into a single token.
func (l *Lexer) foldNewLines() {
	l.incrementLine()
	l.swallowNewLines()
}

// Swallow consecutive newlines.
func (l *Lexer) swallowNewLines() {
	for l.matchChar('\n') || (l.matchChar('\r') && l.matchChar('\n')) {
		l.incrementLine()
	}
}

// Checks whether the given character is a digit.
func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

// Assumes that a character has already been consumed.
// Checks whether the current char is a new line.
func (l *Lexer) isNewLine(char rune) bool {
	return char == '\n' || (char == '\r' && l.matchChar('\n'))
}

// Increments the line number and resets the column number.
func (l *Lexer) incrementLine() {
	l.line += 1
	l.column = 1
}

// Returns the current token value.
func (l *Lexer) tokenValue() string {
	return string(l.source[l.start:l.cursor])
}

// Creates a new lexing error token.
func (l *Lexer) lexError(message string) *token.Token {
	return l.tokenWithValue(token.ERROR, message)
}

// Same as [tokenWithValue] but automatically adds
// the already consumed lexeme as the value of the new token.
func (l *Lexer) tokenWithConsumedValue(typ token.Type) *token.Token {
	return l.tokenWithValue(typ, l.tokenValue())
}

// Builds a token without a string value, based on the current state of the Lexer and
// advances the cursors.
func (l *Lexer) token(typ token.Type) *token.Token {
	return l.tokenWithValue(typ, "")
}

func (l *Lexer) tokenWithValue(typ token.Type, value string) *token.Token {
	startPos := position.New(
		l.start,
		l.startLine,
		l.startColumn,
	)
	endColumn := l.column - 1
	end := l.cursor - 1
	var endPos *position.Position

	if end == l.start {
		endPos = startPos
	} else {
		endPos = position.New(
			end,
			l.line,
			endColumn,
		)
	}
	token := token.NewWithValue(
		position.NewSpan(
			startPos,
			endPos,
		),
		typ,
		value,
	)
	l.start = l.cursor
	l.startColumn = l.column
	l.startLine = l.line

	return token
}

// Attempts to scan and construct the next token.
func (l *Lexer) scanToken() *token.Token {
	switch l.mode {
	case normalMode:
		return l.scanNormal()
	case charClassMode:
		return l.scanCharClass()
	default:
		return l.lexError(fmt.Sprintf("unsupported lexing mode `%d`", l.mode))
	}
}

// Scan characters in normal mode.
func (l *Lexer) scanNormal() *token.Token {
	for {
		char, ok := l.advanceChar()
		if !ok {
			return l.token(token.END_OF_FILE)
		}

		switch char {
		case '.':
			return l.token(token.DOT)
		case ',':
			return l.token(token.COMMA)
		case '|':
			return l.token(token.PIPE)
		case '{':
			return l.token(token.LBRACE)
		case '}':
			return l.token(token.RBRACE)
		case '[':
			l.mode = charClassMode
			return l.token(token.LBRACKET)
		case ']':
			return l.token(token.RBRACKET)
		case '\\':
			switch l.peekChar() {
			case 'x':
				l.advanceChar()
				return l.token(token.HEX_ESCAPE)
			case 'a':
				l.advanceChar()
				return l.token(token.BELL_ESCAPE)
			case 'f':
				l.advanceChar()
				return l.token(token.FORM_FEED_ESCAPE)
			case 't':
				l.advanceChar()
				return l.token(token.TAB_ESCAPE)
			case 'n':
				l.advanceChar()
				return l.token(token.NEWLINE_ESCAPE)
			case 'r':
				l.advanceChar()
				return l.token(token.CARRIAGE_RETURN_ESCAPE)
			case 'v':
				l.advanceChar()
				return l.token(token.VERTICAL_TAB_ESCAPE)
			case 'p':
				l.advanceChar()
				return l.token(token.UNICODE_CHAR_CLASS_ESCAPE)
			case 'P':
				l.advanceChar()
				return l.token(token.NEGATED_UNICODE_CHAR_CLASS_ESCAPE)
			case 'A':
				l.advanceChar()
				return l.token(token.BEGINNING_OF_TEXT_ESCAPE)
			case 'z':
				l.advanceChar()
				return l.token(token.END_OF_TEXT_ESCAPE)
			case 'b':
				l.advanceChar()
				return l.token(token.WORD_BOUNDARY_ESCAPE)
			case 'B':
				l.advanceChar()
				return l.token(token.NOT_WORD_BOUNDARY_ESCAPE)
			case 'w':
				l.advanceChar()
				return l.token(token.WORD_ESCAPE)
			case 'W':
				l.advanceChar()
				return l.token(token.NOT_WORD_ESCAPE)
			case 'd':
				l.advanceChar()
				return l.token(token.DIGIT_ESCAPE)
			case 'D':
				l.advanceChar()
				return l.token(token.NOT_DIGIT_ESCAPE)
			case 's':
				l.advanceChar()
				return l.token(token.WHITESPACE_ESCAPE)
			case 'S':
				l.advanceChar()
				return l.token(token.NOT_WHITESPACE_ESCAPE)
			}
			ch, ok := l.advanceChar()
			if !ok {
				return l.lexError("trailing backslash")
			}
			return l.lexError(fmt.Sprintf("invalid escape sequence: \\%c", ch))
		default:
			return l.tokenWithConsumedValue(token.CHAR)
		}
	}
}

// Scan characters in char class mode.
func (l *Lexer) scanCharClass() *token.Token {
	for {
		char, ok := l.advanceChar()
		if !ok {
			return l.token(token.END_OF_FILE)
		}

		switch char {
		case '[':
			for l.peekChar() != ']' {
				if _, ok := l.advanceChar(); !ok {
					return l.tokenWithConsumedValue(token.CHAR_LIST)
				}
			}

			charClassContent := l.source[l.start+1 : l.cursor]
			if charClassContent[0] == ':' && charClassContent[len(charClassContent)-1] == ':' {
				l.advanceChar()
				return l.tokenWithValue(token.POSIX_CLASS, charClassContent[1:len(charClassContent)-1])
			}
			return l.tokenWithConsumedValue(token.CHAR_LIST)
		case ']':
			l.mode = normalMode
			return l.token(token.RBRACKET)
		case '\\':

		default:
			for l.peekChar() != ']' {
				if _, ok := l.advanceChar(); !ok {
					break
				}
			}
			return l.tokenWithConsumedValue(token.CHAR_LIST)
		}
	}
}
