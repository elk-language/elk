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
	"strings"
	"unicode/utf8"

	"github.com/fatih/color"
)

// Holds the current state of the lexing process.
type lexer struct {
	// Path to the source file or some name.
	sourceName string
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
	return NewWithName("(eval)", source)
}

// Same as [New] but lets you specify the path to the source code file.
func NewWithName(sourceName string, source []byte) *lexer {
	return &lexer{
		sourceName:  sourceName,
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

// Returns the current character and its length in bytes.
func (l *lexer) currentChar() (rune, int) {
	return utf8.DecodeRune(l.source[l.cursor-1:])
}

// Gets the next UTF-8 encoded character
// without incrementing the cursor.
func (l *lexer) peekChar() rune {
	char, _ := l.nextChar()
	return char
}

// Skips the current character.
func (l *lexer) skipChar() {
	l.start += 1
}

// Skips the current accumulated lexeme.
func (l *lexer) skipLexeme() {
	l.start = l.cursor
}

// Swallow consecutive newlines and wrap them into a single lexeme.
func (l *lexer) foldNewLines() {
	l.incrementLine()

	for {
		if l.matchChar('\n') || (l.matchChar('\r') && l.matchChar('\n')) {
			l.incrementLine()
			continue
		}

		if l.swallowCommentsIfPresent() {
			continue
		}

		break
	}
}

// Swallows a comment if it is present.
func (l *lexer) swallowCommentsIfPresent() bool {
	if !l.matchChar('#') {
		return false
	}

	l.swallowComments()
	return true
}

// Assumes that "#" has already been consumed.
// Skips over the entire comment, be it multiline "#[" ... "]#" or single line "#" ...
func (l *lexer) swallowComments() {
	if l.matchChar('[') {
		nestCounter := 1
		for {
			if l.matchChar('#') && l.matchChar('[') {
				nestCounter += 1
				continue
			}
			if l.matchChar(']') && l.matchChar('#') {
				nestCounter -= 1
				if nestCounter == 0 {
					break
				}
			}
			l.advanceChar()
			if l.isNewLine() {
				l.incrementLine()
			}
		}
		l.skipLexeme()
		return
	}

	for {
		l.advanceChar()
		if l.peekChar() == '\n' {
			break
		}
	}
	l.skipLexeme()
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
			return l.buildLexeme(LexSeparator), nil
		case '>':
			return l.buildLexeme(LexGreater), nil
		case '<':
			return l.buildLexeme(LexLess), nil
		case '\n':
			l.foldNewLines()
			return l.buildLexeme(LexSeparator), nil
		case '\r':
			if l.matchChar('\n') {
				l.foldNewLines()
				return l.buildLexeme(LexSeparator), nil
			}
			return nil, l.lexError()
		case ' ':
			l.skipChar()
		case '#':
			l.swallowComments()
		default:
			return nil, l.lexError()
		}
	}
}

// Checks whether the current char is a new line.
func (l *lexer) isNewLine() bool {
	char, _ := l.currentChar()

	return char == '\n' || (char == '\r' && l.matchChar('\n'))
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
	ellipsis := "(...)"
	maxErrLen := 35

	lexValue := l.lexemeValue()
	lexValue = lexValue[0:minInt(len(lexValue), maxErrLen)]
	i := l.start
	var trimmedLexValue []byte
	var byt byte

	for {
		if i == len(lexValue) {
			break
		}

		byt = l.source[i]
		if byt == '\n' {
			break
		}

		trimmedLexValue = append(trimmedLexValue, byt)
		i += 1
	}
	lexValue = string(trimmedLexValue)
	if len(lexValue) > maxErrLen {
		lexValue = lexValue + ellipsis
	}
	var srcContext []byte
	i = l.start
	for {
		if i == 0 {
			break
		}
		if i == l.start-maxErrLen {
			srcContext = append([]byte(ellipsis), srcContext...)
			break
		}

		i -= 1
		byt = l.source[i]
		if byt == '\n' {
			break
		}

		srcContext = append([]byte{byt}, srcContext...)
	}
	lineStr := fmt.Sprintf("%d", l.startLine)
	arrowStr := color.New(color.Italic, color.Faint).Sprintf("%s   %s^-- There", strings.Repeat(" ", len(lineStr)), strings.Repeat(" ", utf8.RuneCount((srcContext))))
	errFmtString := "%s:%s:%d Lexing error, unexpected %s\n\n\t%s | %s%s\n\t%s"
	lexValue = color.New(color.Bold, color.FgRed).Sprint(lexValue)
	return fmt.Errorf(errFmtString, l.sourceName, lineStr, l.startColumn, lexValue, lineStr, srcContext, lexValue, arrowStr)
}

// Builds a lexeme based on the current state of the lexer and
// advances the cursors.
func (l *lexer) buildLexeme(typ LexemeType) *Lexeme {
	lexeme := &Lexeme{
		typ,
		l.lexemeValue(),
		l.start,
		l.cursor - l.start,
		l.startLine,
		l.startColumn,
	}
	l.start = l.cursor
	l.startColumn = l.column
	l.startLine = l.line

	return lexeme
}
