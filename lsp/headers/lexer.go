package headers

import (
	"strings"
	"unicode/utf8"

	"github.com/elk-language/elk/lsp/headers/token"
)

// Holds the current state of the lexing process.
type Lexer struct {
	// headers source
	source []byte
	// Holds the index of the beginning byte
	// of the currently scanned token.
	start int
	// Holds the index of the current byte
	// the lexer is at.
	cursor         int
	newlineCounter int
}

// Instantiates a new lexer for the given regex.
func NewLexer(source []byte) *Lexer {
	return &Lexer{
		source: source,
	}
}

// Lex the given string and return a slice containing all the tokens.
func Lex(source []byte) []*token.Token {
	l := NewLexer(source)

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

func (l *Lexer) Cursor() int {
	return l.cursor
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
	return char, true
}

// Returns the next character and its length in bytes.
func (l *Lexer) nextChar() (rune, int) {
	return utf8.DecodeRune(l.source[l.cursor:])
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

// Gets the next UTF-8 encoded character
// without incrementing the cursor.
// Returns `'\x00', false` when at EOF.
func (l *Lexer) peekCharOk() (rune, bool) {
	if !l.hasMoreTokens() {
		return '\x00', false
	}
	char, _ := l.nextChar()
	return char, true
}

// Skips the current byte.
func (l *Lexer) skipByte() {
	l.start += 1
}

// Skips the current accumulated token.
func (l *Lexer) skipToken() {
	l.start = l.cursor
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
	token := token.NewWithValue(
		typ,
		value,
	)
	l.start = l.cursor

	return token
}

// Attempts to scan and construct the next token.
func (l *Lexer) scanToken() *token.Token {
	if l.newlineCounter == 2 {
		return l.token(token.END_OF_FILE)
	}

	for {
		char, ok := l.advanceChar()
		if !ok {
			l.newlineCounter = 0
			return l.token(token.END_OF_FILE)
		}

		switch char {
		case ':':
			l.newlineCounter = 0
			return l.token(token.COLON)
		case '\r':
			switch l.peekChar() {
			case '\n':
				l.advanceChar()
				l.newlineCounter++
				return l.token(token.NEWLINE)
			default:
				l.newlineCounter = 0
				return l.lexError("invalid char")
			}
		default:
			l.newlineCounter = 0
			return l.string()
		}
	}
}

func (l *Lexer) string() *token.Token {
	var result strings.Builder
	for {
		char, ok := l.peekCharOk()
		if !ok {
			return l.tokenWithValue(token.STRING, result.String())
		}
		switch char {
		case '\r':
			return l.tokenWithValue(token.STRING, result.String())
		}

		result.WriteRune(char)
		l.advanceChar()
	}
}
