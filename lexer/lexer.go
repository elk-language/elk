// Package lexer implements a lexical analyzer
// used by the Elk interpreter.
//
// Lexer expects a slice of bytes containing Elk source code
// analyses it and returns a stream of tokens/tokens.
//
// Tokens are returned on demand.
package lexer

import (
	"fmt"
	"math"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/fatih/color"
)

// Lexing mode which changes how characters are handled by the lexer.
type mode uint8

const (
	normalMode                mode = iota // Initial mode
	inStringLiteralMode                   // Triggered after consuming the initial token `"` of a string literal.
	inStringInterpolationMode             // Triggered after consuming the initial token `${` of string interpolation
)

// Holds the current state of the lexing process.
type lexer struct {
	// Path to the source file or some name.
	sourceName string
	// Elk source code.
	source []byte
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
	// Current lexing mode.
	mode mode
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
		mode:        normalMode,
	}
}

// Returns true if there is any code left to analyse.
func (l *lexer) hasMoreTokens() bool {
	return l.cursor < len(l.source)
}

// Returns the next token or an error if
// the input is malformed.
func (l *lexer) Next() (*Token, error) {
	if !l.hasMoreTokens() {
		return newEOF(), nil
	}

	return l.scanToken()
}

// Gets the next UTF-8 encoded character
// and increments the cursor.
func (l *lexer) advanceChar() (rune, bool) {
	if !l.hasMoreTokens() {
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
	if !l.hasMoreTokens() {
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
	if !l.hasMoreTokens() {
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
		if !strings.ContainsRune(validChars, l.peekChar()) {
			break
		}

		_, ok := l.advanceChar()
		if !ok {
			return false
		}
	}
	return true
}

// Returns the next character and its length in bytes.
func (l *lexer) nextChar() (rune, int) {
	return utf8.DecodeRune(l.source[l.cursor:])
}

// Returns the second next character and its length in bytes.
func (l *lexer) nextNextChar() (rune, int) {
	if !l.hasMoreTokens() {
		return '\x00', 0
	}
	return utf8.DecodeRune(l.source[l.cursor+1:])
}

// Gets the next UTF-8 encoded character
// without incrementing the cursor.
func (l *lexer) peekChar() rune {
	if !l.hasMoreTokens() {
		return '\x00'
	}
	char, _ := l.nextChar()
	return char
}

// Gets the second UTF-8 encoded character
// without incrementing the cursor.
func (l *lexer) peekNextChar() rune {
	char, _ := l.nextNextChar()
	return char
}

// Skips the current character.
func (l *lexer) skipChar() {
	l.start += 1
	l.startColumn += 1
}

// Skips the current accumulated token.
func (l *lexer) skipToken() {
	l.start = l.cursor
}

// Swallow consecutive newlines and wrap them into a single token.
func (l *lexer) foldNewLines() {
	l.incrementLine()

	for l.matchChar('\n') || (l.matchChar('\r') && l.matchChar('\n')) {
		l.incrementLine()
	}
}

// Assumes that `##[` has already been consumed.
// Builds the doc comment token.
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
			return "", l.lexErrorWithHint(fmt.Sprintf("unbalanced doc comments, expected %d more doc comment ending(s) `%s`", nestCounter, warnFmt.Sprint("]##")))
		}
		docStrLines[docStrLine] += string(char)

		if !nonIndentChars && char == ' ' || char == '\t' {
			indent += 1
		} else if l.isNewLine(char) {
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
		// add 1 because of the trailing newline
		if len(line) < leastIndented+1 {
			result += "\n"
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
	l.skipToken()
}

var (
	warnFmt  = color.New(color.Bold, color.FgHiYellow)
	errorFmt = color.New(color.Bold, color.FgHiRed)
	hintFmt  = color.New(color.Faint)
)

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
		char, ok := l.advanceChar()
		if !ok {
			return l.lexErrorWithHint(fmt.Sprintf("unbalanced block comments, expected %d more block comment ending(s) `%s`", nestCounter, warnFmt.Sprint("]#")))
		}
		if l.isNewLine(char) {
			l.incrementLine()
		}
	}
	l.skipToken()
	return nil
}

// Assumes that the beginning quote `'` has already been consumed.
// Consumes a raw string delimited by single quotes.
func (l *lexer) consumeRawString() (string, error) {
	var result string
	for {
		char, ok := l.advanceChar()
		if !ok {
			return "", l.lexErrorWithHint(fmt.Sprintf("unterminated raw string, missing `%s`", warnFmt.Sprint("'")))
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
func (l *lexer) consumeNumber(startDigit rune) *Token {
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
		return l.tokenWithConsumedValue(LexInt)
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
		return l.tokenWithConsumedValue(LexFloat)
	}
	return l.tokenWithConsumedValue(LexInt)
}

// Assumes that the initial letter has already been consumed.
func (l *lexer) consumeIdentifier(init rune) *Token {
	if unicode.IsUpper(init) {
		// constant
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		return l.tokenWithConsumedValue(LexConstant)
	} else {
		// variable or method name
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		if lexType := keywords[l.tokenValue()]; LexKeywordBeg < lexType && lexType < LexKeywordEnd {
			// Is a keyword
			return l.token(lexType)
		}
		return l.tokenWithConsumedValue(LexIdentifier)
	}
}

// Assumes that the initial "_" has already been consumed.
func (l *lexer) consumePrivateIdentifier() *Token {
	if unicode.IsUpper(l.peekChar()) {
		// constant
		l.advanceChar()
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		return l.tokenWithConsumedValue(LexPrivateConstant)
	} else if unicode.IsLower(l.peekChar()) {
		// variable or method name
		l.advanceChar()
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		return l.tokenWithConsumedValue(LexPrivateIdentifier)
	}

	return l.tokenWithConsumedValue(LexPrivateIdentifier)
}

// Attempts to scan and construct the next token.
func (l *lexer) scanToken() (*Token, error) {
	switch l.mode {
	case normalMode:
		return l.scanNormal()
	case inStringLiteralMode:
		return l.scanStringLiteral()
	case inStringInterpolationMode:
		return l.scanStringInterpolation()
	default:
		return nil, l.lexErrorWithHint("unsupported lexing mode")
	}
}

// Scan characters when inside of string interpolation.
func (l *lexer) scanStringInterpolation() (*Token, error) {
	return nil, l.lexErrorWithHint("not implemented yet")
}

// Scan characters when inside of a string literal.
func (l *lexer) scanStringLiteral() (*Token, error) {
	return nil, l.lexErrorWithHint("not implemented yet")
}

// Scan characters in normal mode.
func (l *lexer) scanNormal() (*Token, error) {
	for {
		char, ok := l.advanceChar()
		if !ok {
			return nil, l.lexError("unexpected end of file")
		}

		switch char {
		case '[':
			return l.token(LexLBracket), nil
		case ']':
			return l.token(LexRBracket), nil
		case '(':
			return l.token(LexLParen), nil
		case ')':
			return l.token(LexRParen), nil
		case '{':
			return l.token(LexLBrace), nil
		case '}':
			return l.token(LexLBrace), nil
		case ',':
			return l.token(LexComma), nil
		case '.':
			return l.token(LexDot), nil
		case '-':
			if l.matchChar('=') {
				return l.token(LexMinusEqual), nil
			}
			if l.matchChar('>') {
				return l.token(LexThinArrow), nil
			}
			return l.token(LexMinus), nil
		case '+':
			if l.matchChar('=') {
				return l.token(LexPlusEqual), nil
			}
			return l.token(LexPlus), nil
		case '*':
			if l.matchChar('=') {
				return l.token(LexStarEqual), nil
			}
			if l.matchChar('*') {
				if l.matchChar('=') {
					return l.token(LexPowerEqual), nil
				}
				return l.token(LexPower), nil
			}
			return l.token(LexStar), nil
		case '=':
			if l.matchChar('=') {
				if l.matchChar('=') {
					return l.token(LexStrictEqual), nil
				}
				return l.token(LexEqual), nil
			}
			if l.matchChar('~') {
				return l.token(LexMatchOperator), nil
			}
			if l.matchChar('>') {
				return l.token(LexThickArrow), nil
			}
			if l.peekChar() == ':' && l.peekNextChar() == '=' {
				l.advanceChar()
				l.advanceChar()
				return l.token(LexRefEqual), nil
			}
			if l.peekChar() == '!' && l.peekNextChar() == '=' {
				l.advanceChar()
				l.advanceChar()
				return l.token(LexRefNotEqual), nil
			}
			return l.token(LexAssign), nil
		case ':':
			if l.matchChar(':') {
				return l.token(LexScopeResOperator), nil
			}
			if l.matchChar('=') {
				return l.token(LexColonEqual), nil
			}
			if l.matchChar('>') {
				if l.matchChar('>') {
					return l.token(LexReverseInstanceOf), nil
				}
				return l.token(LexReverseSubtype), nil
			}

			return l.token(LexColon), nil
		case '~':
			if l.matchChar('=') {
				return l.token(LexTildeEqual), nil
			}
			if l.matchChar('>') {
				return l.token(LexWigglyArrow), nil
			}
			return l.token(LexTilde), nil
		case ';':
			return l.token(LexSeparator), nil
		case '>':
			if l.matchChar('=') {
				return l.token(LexGreaterEqual), nil
			}
			if l.matchChar('>') {
				if l.matchChar('=') {
					return l.token(LexRBitShiftEqual), nil
				}
				return l.token(LexRBitShift), nil
			}
			return l.token(LexGreater), nil
		case '<':
			if l.matchChar('=') {
				return l.token(LexLessEqual), nil
			}
			if l.matchChar(':') {
				return l.token(LexSubtype), nil
			}
			if l.matchChar('<') {
				if l.matchChar('=') {
					return l.token(LexLBitShiftEqual), nil
				}
				if l.matchChar(':') {
					return l.token(LexInstanceOf), nil
				}
				return l.token(LexLBitShift), nil
			}
			return l.token(LexLess), nil
		case '&':
			if l.matchChar('&') {
				if l.matchChar('=') {
					return l.token(LexAndAndEqual), nil
				}
				return l.token(LexAndAnd), nil
			}
			if l.matchChar('=') {
				return l.token(LexAndEqual), nil
			}
			return l.token(LexAnd), nil
		case '|':
			if l.matchChar('|') {
				if l.matchChar('=') {
					return l.token(LexOrOrEqual), nil
				}
				return l.token(LexOrOr), nil
			}
			if l.matchChar('>') {
				return l.token(LexPipeOperator), nil
			}
			if l.matchChar('=') {
				return l.token(LexOrEqual), nil
			}
			return l.token(LexOr), nil
		case '?':
			if l.matchChar('?') {
				if l.matchChar('=') {
					return l.token(LexNilCoalesceEqual), nil
				}
				return l.token(LexNilCoalesce), nil
			}
			return l.token(LexQuestionMark), nil
		case '!':
			if l.matchChar('=') {
				if l.matchChar('=') {
					return l.token(LexStrictNotEqual), nil
				}
				return l.token(LexNotEqual), nil
			}
			return l.token(LexBang), nil
		case '%':
			if l.matchChar('=') {
				return l.token(LexPercentEqual), nil
			}
			if l.matchChar('w') {
				return l.token(LexPercentW), nil
			}
			if l.matchChar('s') {
				return l.token(LexPercentS), nil
			}
			if l.matchChar('i') {
				return l.token(LexPercentI), nil
			}
			if l.matchChar('f') {
				return l.token(LexPercentF), nil
			}
			if l.matchChar('{') {
				return l.token(LexSetLiteralBeg), nil
			}
			if l.matchChar('(') {
				return l.token(LexTupleLiteralBeg), nil
			}
			return l.token(LexPercent), nil

		case '\n':
			l.foldNewLines()
			return l.token(LexSeparator), nil
		case '\r':
			if l.matchChar('\n') {
				l.foldNewLines()
				return l.token(LexSeparator), nil
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
				token := l.tokenWithValue(LexDocComment, str)
				l.start += 3
				l.cursor += 3
				return token, nil
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
			return l.tokenWithValue(LexRawString, str), nil
		case '"':
			l.mode = inStringLiteralMode
			return l.token(LexStringBeg), nil
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

// Assumes that a character has already been consumed.
// Checks whether the current char is a new line.
func (l *lexer) isNewLine(char rune) bool {
	return char == '\n' || (char == '\r' && l.matchChar('\n'))
}

// Increments the line number and resets the column number.
func (l *lexer) incrementLine() {
	l.line += 1
	l.column = 1
}

// Returns the current token value.
func (l *lexer) tokenValue() string {
	return string(l.source[l.start:l.cursor])
}

// Returns the lesser integer.
func minInt(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

const (
	ellipsis  = "(...)"
	maxErrLen = 35
)

// Builds a lexing error with a hint from source code
// based on the current state of the lexer.
func (l *lexer) lexErrorWithHint(message string) error {
	lexValue := l.tokenValue()
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
	var srcContextLen int
	i = l.start
	for {
		if i == 0 {
			srcContextLen = utf8.RuneCount(srcContext)
			break
		}
		if i == l.start-maxErrLen {
			srcContextLen = utf8.RuneCount(srcContext) + utf8.RuneCountInString(ellipsis)
			srcContext = append([]byte(hintFmt.Sprint(ellipsis)), srcContext...)
			break
		}

		i -= 1
		byt = l.source[i]
		if byt == '\n' {
			srcContextLen = utf8.RuneCount(srcContext)
			break
		}

		srcContext = append([]byte{byt}, srcContext...)
	}
	lineStr := fmt.Sprintf("%d", l.startLine)
	arrowStr := hintFmt.Sprintf("%s   %s^-- There", strings.Repeat(" ", len(lineStr)), strings.Repeat(" ", srcContextLen))
	lexValue = errorFmt.Sprint(lexValue)
	return l.lexError(fmt.Sprintf("%s\n\n\t%s | %s%s\n\t%s", message, lineStr, srcContext, lexValue, arrowStr))
}

// Creates a new lexing error which shows unexpected characters.
func (l *lexer) unexpectedCharsError() error {
	return l.lexErrorWithHint("unexpected characters")
}

// Creates a new lexing error.
func (l *lexer) lexError(message string) error {
	return fmt.Errorf("%s:%d:%d Lexing error, %s", l.sourceName, l.startLine, l.startColumn, message)
}

// Same as [tokenWithValue] but automatically adds
// the already consumed lexeme as the value of the new token.
func (l *lexer) tokenWithConsumedValue(typ TokenType) *Token {
	return l.tokenWithValue(typ, l.tokenValue())
}

// Builds a token without a string value, based on the current state of the lexer and
// advances the cursors.
func (l *lexer) token(typ TokenType) *Token {
	return l.tokenWithValue(typ, "")
}

// Same as [token] but lets you specify the value of the token
// manually.
func (l *lexer) tokenWithValue(typ TokenType, value string) *Token {
	token := &Token{
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

	return token
}
