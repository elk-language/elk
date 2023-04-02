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
	"math"
	"strings"
	"unicode"
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
func (l *lexer) advanceChar() (rune, bool) {
	if !l.HasMoreLexemes() {
		return 0, false
	}

	char, size := l.nextChar()

	l.cursor += size
	l.column += 1
	return char, true
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

// Consumes the next character if it's from the valid set.
func (l *lexer) acceptChars(validChars string) bool {
	if !l.HasMoreLexemes() {
		return false
	}

	if strings.ContainsRune(validChars, l.peekChar()) {
		l.advanceChar()
		return true
	}

	return false
}

// Consumes a series of characters from the given set.
func (l *lexer) acceptCharsRun(validChars string) bool {
	for {
		if strings.ContainsRune(validChars, l.peekChar()) {
			_, ok := l.advanceChar()
			if !ok {
				return false
			}
		} else {
			break
		}
	}
	return true
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
	if !l.HasMoreLexemes() {
		return '\x00'
	}
	char, _ := l.nextChar()
	return char
}

// Skips the current character.
func (l *lexer) skipChar() {
	l.start += 1
	l.startColumn += 1
}

// Skips the current accumulated lexeme.
func (l *lexer) skipLexeme() {
	l.start = l.cursor
}

// Swallow consecutive newlines and wrap them into a single lexeme.
func (l *lexer) foldNewLines() {
	l.incrementLine()

	for l.matchChar('\n') || (l.matchChar('\r') && l.matchChar('\n')) {
		l.incrementLine()
	}
}

// Assumes that "##[" has already been consumed.
// Builds the doc comment lexeme.
func (l *lexer) consumeDocComment() (string, error) {
	nestCounter := 1
	docStrLines := []string{""}
	docStrLine := 0
	leastIndented := math.MaxInt
	indent := 0
	nonIndentChars := false
	for {
		if l.matchChar('#') {
			nonIndentChars = true
			if l.matchChar('#') && l.matchChar('[') {
				docStrLines[docStrLine] += "##["
				nestCounter += 1
				continue
			}
		}
		if l.matchChar(']') {
			nonIndentChars = true
			if l.matchChar('#') && l.matchChar('#') {
				nestCounter -= 1
				if nestCounter == 0 {
					break
				}
				docStrLines[docStrLine] += "]##"
			}
		}
		char, ok := l.advanceChar()
		if !ok {
			return "", l.lexError(fmt.Sprintf("unbalanced doc comments, expected %d more doc comment ending(s) `%s`", nestCounter, color.New(color.Bold).Sprint("]##")))
		}
		docStrLines[docStrLine] += string(char)

		if !nonIndentChars && char == ' ' || char == '\t' {
			indent += 1
		} else if l.isNewLine() {
			l.incrementLine()
			docStrLines = append(docStrLines, "")
			docStrLine += 1
			if nonIndentChars && indent < leastIndented {
				leastIndented = indent
			}
			indent = 0
			nonIndentChars = false
		} else {
			nonIndentChars = true
		}
	}

	var result string
	for _, line := range docStrLines {
		if len(line) < leastIndented {
			result += line
			continue
		}

		result += line[leastIndented:]
	}
	result = strings.TrimPrefix(result, "\n")
	result = strings.TrimRight(result, "\t\n ")

	return result, nil
}

// Assumes that "#" has already been consumed.
// Skips over a single line comment "#" ...
func (l *lexer) swallowSingleLineComment() {
	for {
		if l.peekChar() == '\n' {
			break
		}
		if _, ok := l.advanceChar(); !ok {
			return
		}
	}
	l.skipLexeme()
}

// Assumes that "#[" has already been consumed.
// Skips over a block comment "#[" ... "]#".
func (l *lexer) swallowBlockComments() error {
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
		if _, ok := l.advanceChar(); !ok {
			return l.lexError(fmt.Sprintf("unbalanced block comments, expected %d more block comment ending(s) `%s`", nestCounter, color.New(color.Bold).Sprint("]#")))
		}
		if l.isNewLine() {
			l.incrementLine()
		}
	}
	l.skipLexeme()
	return nil
}

// Assumes that the beginning quote ' has already been consumed.
// Consumes a raw string delimited by single quotes.
func (l *lexer) consumeRawString() (string, error) {
	var result string
	for {
		char, ok := l.advanceChar()
		if !ok {
			return "", l.lexError("unterminated raw string, missing `'`")
		}
		if char == '\'' {
			break
		}
		if char == '\n' {
			l.line += 1
		}
		result += string(char)
	}

	return result, nil
}

// Assumes that the first digit has already been consumed.
// Consumes a number literal.
func (l *lexer) consumeNumber(startDigit rune) *Lexeme {
	nonDecimal := false
	digits := "0123456789_"
	if startDigit == '0' {
		if l.acceptChars("xX") {
			// hexadecimal (base 16)
			digits = "0123456789abcdefABCDEF_"
			nonDecimal = true
		} else if l.acceptChars("dD") {
			// duodecimal (base 12)
			digits = "0123456789ab"
			nonDecimal = true
		} else if l.acceptChars("oO") {
			// octal (base 8)
			digits = "01234567_"
			nonDecimal = true
		} else if l.acceptChars("qQ") {
			// quaternary (base 4)
			digits = "0123_"
			nonDecimal = true
		} else if l.acceptChars("bB") {
			// binary (base 2)
			digits = "01_"
			nonDecimal = true
		}
	}

	l.acceptCharsRun(digits)
	if nonDecimal {
		return l.buildLexeme(LexInt)
	}

	var isFloat bool
	if l.matchChar('.') {
		l.acceptCharsRun(digits)
		isFloat = true
	}
	if l.acceptChars("eE") {
		l.acceptChars("+-")
		l.acceptCharsRun(digits)
		isFloat = true
	}

	if isFloat {
		return l.buildLexeme(LexFloat)
	}
	return l.buildLexeme(LexInt)
}

// Assumes that the initial letter has already been consumed.
func (l *lexer) consumeIdentifier(init rune) *Lexeme {
	if unicode.IsUpper(init) {
		// constant
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		return l.buildLexeme(LexConstant)
	} else {
		// variable or method name
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		return l.buildLexeme(LexIdentifier)
	}
}

// Assumes that the initial "_" has already been consumed.
func (l *lexer) consumePrivateIdentifier() *Lexeme {
	if unicode.IsUpper(l.peekChar()) {
		// constant
		l.advanceChar()
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		return l.buildLexeme(LexPrivateConstant)
	} else if unicode.IsLower(l.peekChar()) {
		// variable or method name
		l.advanceChar()
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		return l.buildLexeme(LexPrivateIdentifier)
	}

	return l.buildLexeme(LexPrivateIdentifier)
}

// Attempts to scan and construct the next lexeme.
func (l *lexer) scanLexeme() (*Lexeme, error) {
	for {
		char, ok := l.advanceChar()
		if !ok {
			return nil, l.lexError("unexpected end of file")
		}

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
			fallthrough
		case '\t':
			fallthrough
		case ' ':
			l.skipChar()
		case '#':
			if l.matchChar('#') && l.matchChar('[') {
				str, err := l.consumeDocComment()
				if err != nil {
					return nil, err
				}
				l.start += 3
				l.cursor -= 3
				lexeme := l.buildLexemeWithValue(LexDocComment, str)
				l.start += 3
				l.cursor += 3
				return lexeme, nil
			}

			if l.matchChar('[') {
				err := l.swallowBlockComments()
				if err != nil {
					return nil, err
				}
			} else {
				l.swallowSingleLineComment()
			}
		case '\'':
			str, err := l.consumeRawString()
			if err != nil {
				return nil, err
			}
			return l.buildLexemeWithValue(LexRawString, str), nil
		case '_':
			return l.consumePrivateIdentifier(), nil
		default:
			if isDigit(char) {
				return l.consumeNumber(char), nil
			} else if unicode.IsLetter(char) {
				return l.consumeIdentifier(char), nil
			}
			return nil, l.unexpectedCharsError()
		}
	}
}

// Checks whether the given character is acceptable
// inside an identifier.
func isIdentifierChar(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsNumber(char) || char == '_'
}

// Checks whether the given character is a digit.
func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

// Checks whether the current char is a new line.
func (l *lexer) isNewLine() bool {
	char, _ := l.currentChar()

	return char == '\n' || (char == '\r' && l.matchChar('\n'))
}

// Increments the line number and resets the column number.
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

// Creates a new lexing error based on the current state
// of the lexer.
func (l *lexer) unexpectedCharsError() error {
	ellipsis := "(...)"
	maxErrLen := 35

	lexValue := l.lexemeValue()
	lexValue = lexValue[0:minInt(len(lexValue), maxErrLen)]
	i := l.start
	var trimmedLexValue []byte
	var byt byte

	for {
		if i == len(lexValue) || i == len(l.source) {
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
	lexValue = color.New(color.Bold, color.FgRed).Sprint(lexValue)
	return l.lexError(fmt.Sprintf("unexpected `%s`\n\n\t%s | %s%s\n\t%s", lexValue, lineStr, srcContext, lexValue, arrowStr))
}

// Creates a new lexing error.
func (l *lexer) lexError(message string) error {
	return fmt.Errorf("%s:%d:%d Lexing error, %s", l.sourceName, l.startLine, l.startColumn, message)
}

// Builds a lexeme based on the current state of the lexer and
// advances the cursors.
func (l *lexer) buildLexeme(typ LexemeType) *Lexeme {
	return l.buildLexemeWithValue(typ, l.lexemeValue())
}

// Same as [buildLexeme] but lets you specify the value of the lexeme
// manually.
func (l *lexer) buildLexemeWithValue(typ LexemeType, value string) *Lexeme {
	lexeme := &Lexeme{
		typ,
		value,
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
