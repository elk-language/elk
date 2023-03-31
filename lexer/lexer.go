// Package lexer implements a lexical analyzer
// used by the Elk interpreter.
//
// Lexer expects a string containing Elk source code
// analyses it and returns a stream of lexemes/tokens.
//
// Lexemes are returned on demand.
package lexer

import (
	"fmt"
	"unicode/utf8"
)

// Holds the current state of the lexing process.
type lexer struct {
	// Elk source code.
	source []byte
	// Holds the index of the beginning byte
	// of the currently scanned lexeme.
	start int
	// Holds the index of the current byte
	// the lexer is at.
	cursor int
	// Column of the first character of the currently analysed lexeme.
	startColumn int
	// Column of the current character of the currently analysed lexeme.
	column int
	// First line number of the currently analysed lexeme.
	startLine int
	// Current line of source code being analysed.
	line int
}

// Instantiates a new lexer for the given source code.
func New(source []byte) *lexer {
	return &lexer{
		source:      source,
		line:        1,
		startLine:   1,
		column:      1,
		startColumn: 1,
	}
}

// Returns true if there is any code left to analyse.
func (l *lexer) HasMoreLexemes() bool {
	return l.cursor < len(l.source)
}

// Returns the next lexeme or an error if
// the input is malformed.
func (l *lexer) Next() (*Lexeme, error) {
	if !l.HasMoreLexemes() {
		return newEOF(), nil
	}

	return l.scanLexeme()
}

// Gets the next UTF-8 encoded character
// and increments the cursor.
func (l *lexer) advanceChar() rune {
	char, size := l.nextChar()

	l.cursor += size
	l.column += 1
	return char
}

// Checks if the given character matches
// the next UTF-8 encoded character in source code.
// If they match, the cursor gets incremented.
func (l *lexer) matchChar(char rune) bool {
	if !l.HasMoreLexemes() {
		return false
	}

	if l.peekChar() == char {
		l.cursor += 1
		return true
	}

	return false
}

// Returns the next character and its length in bytes.
func (l *lexer) nextChar() (rune, int) {
	return utf8.DecodeRune(l.source[l.cursor:])
}

// Gets the next UTF-8 encoded character
// without incrementing the cursor.
func (l *lexer) peekChar() rune {
	char, _ := l.nextChar()
	return char
}

// Attempts to scan and construct the next lexeme.
func (l *lexer) scanLexeme() (*Lexeme, error) {
	for {
		char := l.advanceChar()

		switch char {
		case '(':
			return l.buildLexeme(LexLParen), nil
		case ')':
			return l.buildLexeme(LexRParen), nil
		case '{':
			return l.buildLexeme(LexLBrace), nil
		case '}':
			return l.buildLexeme(LexLBrace), nil
		case ',':
			return l.buildLexeme(LexComma), nil
		case '.':
			return l.buildLexeme(LexDot), nil
		case '-':
			return l.buildLexeme(LexMinus), nil
		case '+':
			return l.buildLexeme(LexPlus), nil
		case ';':
			return l.buildLexeme(LexSemicolon), nil
		case '>':
			return l.buildLexeme(LexGreater), nil
		case '<':
			return l.buildLexeme(LexLess), nil
		case '\n':
			l.incrementLine()
			return l.buildLexeme(LexNewLine), nil
		case '\r':
			if l.matchChar('\n') {
				l.incrementLine()
				return l.buildLexeme(LexNewLine), nil
			}
			return nil, l.lexError()
		default:
			return nil, l.lexError()
		}
	}
}

// Increments the line number and resets the column number
func (l *lexer) incrementLine() {
	l.line += 1
	l.column = 1
}

// Returns the current lexeme value.
func (l *lexer) lexemeValue() string {
	return string(l.source[l.start:l.cursor])
}

// Returns the lesser integer.
func minInt(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

func (l *lexer) lexError() error {
	maxErrLen := 60
	lexValue := l.lexemeValue()
	origLexLen := len(lexValue)
	lexValue = lexValue[0:minInt(origLexLen, maxErrLen)]
	if origLexLen > 60 {
		lexValue = lexValue + "..."
	}

	return fmt.Errorf("%d:%d: lexing error, unexpected token: %s", l.startLine, l.startColumn, lexValue)
}

// Builds a lexeme based on the current state of the lexer and
// advances the cursors.
func (l *lexer) buildLexeme(typ LexemeType) *Lexeme {
	lexeme := &Lexeme{
		typ,
		l.lexemeValue(),
		l.start,
		l.cursor - l.start,
		l.line,
		l.startColumn,
	}
	l.start = l.cursor
	l.startColumn = l.column

	return lexeme
}
