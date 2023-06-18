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
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/elk-language/elk/position"
	"github.com/elk-language/elk/token"
)

// Lexing mode which changes how characters are handled by the lexer.
type mode uint8

const (
	normalMode mode = iota // Initial mode

	wordArrayLiteralMode   // Triggered after entering the initial token `%w[` of a word array literal
	symbolArrayLiteralMode // Triggered after entering the initial token `%s[` of a symbol array literal
	hexArrayLiteralMode    // Triggered after entering the initial token `%x[` of a hex array literal
	binArrayLiteralMode    // Triggered after entering the initial token `%b[` of a binary array literal

	wordSetLiteralMode   // Triggered after entering the initial token `%w{` of a word set literal
	symbolSetLiteralMode // Triggered after entering the initial token `%s{` of a symbol set literal
	hexSetLiteralMode    // Triggered after entering the initial token `%x{` of a hex set literal
	binSetLiteralMode    // Triggered after entering the initial token `%b{` of a binary set literal

	wordTupleLiteralMode   // Triggered after entering the initial token `%w(` of a word tuple literal
	symbolTupleLiteralMode // Triggered after entering the initial token `%s(` of a symbol tuple literal
	hexTupleLiteralMode    // Triggered after entering the initial token `%x(` of a hex tuple literal
	binTupleLiteralMode    // Triggered after entering the initial token `%b(` of a binary tuple literal

	stringLiteralMode           // Triggered after consuming the initial token `"` of a string literal
	invalidHexEscapeMode        // Triggered after encountering an invalid hex escape sequence in a string literal
	invalidUnicodeEscapeMode    // Triggered after encountering an invalid 4 character unicode escape sequence in a string literal
	invalidBigUnicodeEscapeMode // Triggered after encountering an invalid 8 character unicode escape sequence in a string literal
	invalidEscapeMode           // Triggered after encountering an invalid escape sequence in a string literal
	stringInterpolationMode     // Triggered after consuming the initial token `${` of string interpolation
)

// Holds the current state of the lexing process.
type Lexer struct {
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
func New(source []byte) *Lexer {
	return NewWithName("(eval)", source)
}

// Same as [New] but lets you specify the path to the source code file.
func NewWithName(sourceName string, source []byte) *Lexer {
	return &Lexer{
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

// Attempts to scan and construct the next token.
func (l *Lexer) scanToken() *token.Token {
	switch l.mode {
	case normalMode, stringInterpolationMode:
		return l.scanNormal()
	case stringLiteralMode:
		return l.scanStringLiteral()
	case wordArrayLiteralMode:
		return l.scanWordCollectionLiteral(']', token.WORD_LIST_END)
	case symbolArrayLiteralMode:
		return l.scanWordCollectionLiteral(']', token.SYMBOL_LIST_END)
	case wordSetLiteralMode:
		return l.scanWordCollectionLiteral('}', token.WORD_SET_END)
	case symbolSetLiteralMode:
		return l.scanWordCollectionLiteral('}', token.SYMBOL_SET_END)
	case wordTupleLiteralMode:
		return l.scanWordCollectionLiteral(')', token.WORD_TUPLE_END)
	case symbolTupleLiteralMode:
		return l.scanWordCollectionLiteral(')', token.SYMBOL_TUPLE_END)
	case hexArrayLiteralMode:
		return l.scanIntCollectionLiteral(']', token.HEX_LIST_END, hexLiteralChars, token.HEX_INT)
	case hexSetLiteralMode:
		return l.scanIntCollectionLiteral('}', token.HEX_SET_END, hexLiteralChars, token.HEX_INT)
	case hexTupleLiteralMode:
		return l.scanIntCollectionLiteral(')', token.HEX_TUPLE_END, hexLiteralChars, token.HEX_INT)
	case binArrayLiteralMode:
		return l.scanIntCollectionLiteral(']', token.BIN_LIST_END, binaryLiteralChars, token.BIN_INT)
	case binSetLiteralMode:
		return l.scanIntCollectionLiteral('}', token.BIN_SET_END, binaryLiteralChars, token.BIN_INT)
	case binTupleLiteralMode:
		return l.scanIntCollectionLiteral(')', token.BIN_TUPLE_END, binaryLiteralChars, token.BIN_INT)
	case invalidEscapeMode:
		return l.scanInvalidEscape()
	case invalidHexEscapeMode:
		return l.scanInvalidHexEscape()
	case invalidUnicodeEscapeMode:
		return l.scanInvalidUnicodeEscape()
	case invalidBigUnicodeEscapeMode:
		return l.scanInvalidBigUnicodeEscape()
	default:
		return l.lexError(fmt.Sprintf("unsupported lexing mode `%d`", l.mode))
	}
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

// Returns the next character and its length in bytes.
func (l *Lexer) nextChar() (rune, int) {
	return utf8.DecodeRune(l.source[l.cursor:])
}

// Returns the second next character and its length in bytes.
func (l *Lexer) nextNextChar() (rune, int) {
	if !l.hasMoreTokens() {
		return '\x00', 0
	}
	return utf8.DecodeRune(l.source[l.cursor+1:])
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

	for l.matchChar('\n') || (l.matchChar('\r') && l.matchChar('\n')) {
		l.incrementLine()
	}
}

// Assumes that `##[` has already been consumed.
// Builds the doc comment token.
func (l *Lexer) docComment() *token.Token {
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
			if l.matchCharN('#', 2) {
				nestCounter -= 1
				if nestCounter == 0 {
					break
				}
				docStrLines[docStrLine] += "]##"
			}
		}
		char, ok := l.advanceChar()
		if !ok {
			return l.lexError(fmt.Sprintf("unbalanced doc comments, expected %d more doc comment ending(s) `]##`", nestCounter))
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

	if leastIndented == math.MaxInt {
		leastIndented = indent
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
	result = strings.TrimLeft(result, "\n")
	result = strings.TrimRight(result, "\t\n ")

	return l.tokenWithValue(token.DOC_COMMENT, result)
}

// Assumes that "#" has already been consumed.
// Skips over a single line comment "#" ...
func (l *Lexer) swallowSingleLineComment() {
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

// Assumes that "#[" has already been consumed.
// Skips over a block comment "#[" ... "]#".
func (l *Lexer) swallowBlockComments() *token.Token {
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
			return l.lexError(fmt.Sprintf("unbalanced block comments, expected %d more block comment ending(s) `]#`", nestCounter))
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
func (l *Lexer) rawString() *token.Token {
	var result strings.Builder
	for {
		char, ok := l.advanceChar()
		if !ok {
			return l.lexError("unterminated raw string literal, missing `'`")
		}
		if char == '\'' {
			break
		}
		if char == '\n' {
			l.incrementLine()
		}
		result.WriteRune(char)
	}

	return l.tokenWithValue(token.RAW_STRING, result.String())
}

const unterminatedCharLiteralMessage = "unterminated character literal, missing quote"

// Assumes that the beginning c" has already been consumed.
// Consumes a character literal.
func (l *Lexer) character() *token.Token {
	var lexemeBuff strings.Builder

	if l.matchChar('\\') {
		next, ok := l.advanceChar()
		if !ok {
			return l.lexError(unterminatedCharLiteralMessage)
		}
		switch next {
		case '\\':
			lexemeBuff.WriteByte('\\')
		case 'n':
			lexemeBuff.WriteByte('\n')
		case 't':
			lexemeBuff.WriteByte('\t')
		case '"':
			lexemeBuff.WriteByte('"')
		case 'r':
			lexemeBuff.WriteByte('\r')
		case 'a':
			lexemeBuff.WriteByte('\a')
		case 'b':
			lexemeBuff.WriteByte('\b')
		case 'v':
			lexemeBuff.WriteByte('\v')
		case 'f':
			lexemeBuff.WriteByte('\f')
		case 'u':
			if !l.matchCharsN(hexLiteralChars, 4) {
				if !l.swallowUntil('"') {
					return l.lexError(unterminatedCharLiteralMessage)
				}
				return l.lexError(invalidUnicodeEscapeError)
			}
			value, err := strconv.ParseUint(string(l.source[l.cursor-4:l.cursor]), 16, 16)
			if err != nil {
				if !l.swallowUntil('"') {
					return l.lexError(unterminatedCharLiteralMessage)
				}
				return l.lexError(invalidHexEscapeError)
			}
			lexemeBuff.WriteRune(rune(value))
		case 'U':
			if !l.matchCharsN(hexLiteralChars, 8) {
				if !l.swallowUntil('"') {
					return l.lexError(unterminatedCharLiteralMessage)
				}
				return l.lexError(invalidUnicodeEscapeError)
			}
			value, err := strconv.ParseUint(string(l.source[l.cursor-8:l.cursor]), 16, 32)
			if err != nil {
				if !l.swallowUntil('"') {
					return l.lexError(unterminatedCharLiteralMessage)
				}
				return l.lexError(invalidHexEscapeError)
			}
			lexemeBuff.WriteRune(rune(value))
		case 'x':
			if !l.matchCharsN(hexLiteralChars, 2) {
				if !l.swallowUntil('"') {
					return l.lexError(unterminatedCharLiteralMessage)
				}
				return l.lexError(invalidHexEscapeError)
			}
			value, err := strconv.ParseUint(string(l.source[l.cursor-2:l.cursor]), 16, 8)
			if err != nil {
				if !l.swallowUntil('"') {
					return l.lexError(unterminatedCharLiteralMessage)
				}
				return l.lexError(invalidHexEscapeError)
			}
			lexemeBuff.WriteByte(byte(value))
		case '\n':
			l.incrementLine()
			fallthrough
		default:
			l.matchChar('"')
			return l.lexError("invalid escape sequence in a character literal")
		}
	} else {
		ch, ok := l.advanceChar()
		if !ok {
			return l.lexError(unterminatedCharLiteralMessage)
		}
		lexemeBuff.WriteRune(ch)
	}
	if l.matchChar('"') {
		return l.tokenWithValue(token.CHAR_LITERAL, lexemeBuff.String())
	}

	if !l.swallowUntil('"') {
		return l.lexError(unterminatedCharLiteralMessage)
	}

	return l.lexError("invalid char literal with more than one character")
}

// Assumes that the beginning c' has already been consumed.
// Consumes a character literal.
func (l *Lexer) rawCharacter() *token.Token {
	var char string

	ch, ok := l.advanceChar()
	if !ok {
		return l.lexError(unterminatedCharLiteralMessage)
	}
	char = string(ch)
	if l.matchChar('\'') {
		return l.tokenWithValue(token.RAW_CHAR_LITERAL, char)
	}

	if !l.swallowUntil('\'') {
		return l.lexError(unterminatedCharLiteralMessage)
	}

	return l.lexError("invalid raw char literal with more than one character")
}

const (
	hexLiteralChars        = "0123456789abcdefABCDEF"
	duodecimalLiteralChars = "0123456789abAB"
	decimalLiteralChars    = "0123456789"
	octalLiteralChars      = "01234567"
	quaternaryLiteralChars = "0123"
	binaryLiteralChars     = "01"
)

// Consumes digits from the given set
// and appends them to the given buffer.
// Underscores are ignored.
func (l *Lexer) consumeDigits(digitSet string, lexemeBuff *strings.Builder) {
	for {
		if l.peekChar() == '_' {
			l.advanceChar()
		}
		if !l.acceptChars(digitSet) {
			break
		}
		char, _ := l.advanceChar()
		lexemeBuff.WriteRune(char)
	}
}

// Assumes that the first digit has already been consumed.
// Consumes and constructs an Int or Float literal token.
func (l *Lexer) numberLiteral(startDigit rune) *token.Token {
	tokenType := token.DEC_INT
	digits := decimalLiteralChars
	var lexeme strings.Builder

	if startDigit == '0' {
		if l.matchChars("xX") {
			// hexadecimal (base 16)
			digits = hexLiteralChars
			tokenType = token.HEX_INT
		} else if l.matchChars("dD") {
			// duodecimal (base 12)
			digits = duodecimalLiteralChars
			tokenType = token.DUO_INT
		} else if l.matchChars("oO") {
			// octal (base 8)
			digits = octalLiteralChars
			tokenType = token.OCT_INT
		} else if l.matchChars("qQ") {
			// quaternary (base 4)
			digits = quaternaryLiteralChars
			tokenType = token.QUAT_INT
		} else if l.matchChars("bB") {
			// binary (base 2)
			digits = binaryLiteralChars
			tokenType = token.BIN_INT
		}
	}

	if tokenType != token.DEC_INT {
		l.consumeDigits(digits, &lexeme)
		return l.tokenWithValue(tokenType, lexeme.String())
	}
	lexeme.WriteRune(startDigit)
	l.consumeDigits(digits, &lexeme)

	if l.acceptChars(".") && isDigit(l.peekNextChar()) {
		l.advanceChar()
		lexeme.WriteByte('.')
		l.consumeDigits(digits, &lexeme)
		tokenType = token.FLOAT
	}
	if l.matchChars("eE") {
		lexeme.WriteByte('e')
		if ok, ch := l.matchCharsRune("+-"); ok {
			lexeme.WriteRune(ch)
		}

		l.consumeDigits(digits, &lexeme)
		tokenType = token.FLOAT
	}

	return l.tokenWithValue(tokenType, lexeme.String())
}

// Assumes that the initial letter has already been consumed.
// Consumes and constructs a public publicIdentifier token.
func (l *Lexer) publicIdentifier(init rune) *token.Token {
	if unicode.IsUpper(init) {
		// constant
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		return l.tokenWithConsumedValue(token.PUBLIC_CONSTANT)
	} else {
		// variable or method name
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		if lexType := token.Keywords[l.tokenValue()]; lexType.IsKeyword() {
			// Is a keyword
			return l.token(lexType)
		}
		return l.tokenWithConsumedValue(token.PUBLIC_IDENTIFIER)
	}
}

// Assumes that the initial "_" has already been consumed.
// Consumes and constructs a private publicIdentifier token.
func (l *Lexer) privateIdentifier() *token.Token {
	if unicode.IsUpper(l.peekChar()) {
		// constant
		l.advanceChar()
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		return l.tokenWithConsumedValue(token.PRIVATE_CONSTANT)
	} else if unicode.IsLower(l.peekChar()) {
		// variable or method name
		l.advanceChar()
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		return l.tokenWithConsumedValue(token.PRIVATE_IDENTIFIER)
	}

	return l.tokenWithConsumedValue(token.PRIVATE_IDENTIFIER)
}

// Assumes that the initial `@` has been consumed.
// Consumes and constructs an instance variable token.
func (l *Lexer) instanceVariable() *token.Token {
	for isIdentifierChar(l.peekChar()) {
		l.advanceChar()
	}

	return l.tokenWithValue(token.INSTANCE_VARIABLE, string(l.source[l.start+1:l.cursor]))
}

const (
	unterminatedCollectionError = "unterminated %s literal, missing `%c`"
)

// Scans the content of word collection literals be it `%w[`, `%s[`, `%w{`, `%s{`, `%w(`, `%s(`
func (l *Lexer) scanWordCollectionLiteral(terminatorChar rune, terminatorToken token.Type) *token.Token {
	var result strings.Builder
	var nonSpaceCharEncountered bool
	var endOfLiteral bool

	for {
		peek := l.peekChar()
		if peek == terminatorChar {
			endOfLiteral = true
			break
		}

		if unicode.IsSpace(peek) {
			if !nonSpaceCharEncountered {
				l.advanceChar()
				if peek == '\n' {
					l.incrementLine()
				}
				l.skipToken()
				continue
			}

			break
		}

		char, ok := l.advanceChar()
		if !ok {
			return l.lexError(fmt.Sprintf(unterminatedCollectionError, "word collection", terminatorChar))
		}

		if !nonSpaceCharEncountered {
			nonSpaceCharEncountered = true
		}
		result.WriteRune(char)
	}

	if endOfLiteral && result.Len() == 0 {
		l.mode = normalMode
		l.advanceChar()
		return l.token(terminatorToken)
	}
	return l.tokenWithValue(token.RAW_STRING, result.String())
}

// Scans the content of int collection literals be it `%x[`, `%b[`, `%x{`, `%b{`, `%x(`, `%b(`
func (l *Lexer) scanIntCollectionLiteral(terminatorChar rune, terminatorToken token.Type, digitSet string, elementToken token.Type) *token.Token {
	var result strings.Builder
	var nonSpaceCharEncountered bool
	var endOfLiteral bool

	for {
		peek := l.peekChar()
		if peek == terminatorChar {
			endOfLiteral = true
			break
		}

		if unicode.IsSpace(peek) {
			if !nonSpaceCharEncountered {
				l.advanceChar()
				if peek == '\n' {
					l.incrementLine()
				}
				l.skipToken()
				continue
			}

			break
		}

		char, ok := l.advanceChar()
		if !ok {
			return l.lexError(fmt.Sprintf(unterminatedCollectionError, "int collection", terminatorChar))
		}

		if char == '_' && l.peekChar() != '_' {
			continue
		}

		if !strings.ContainsRune(digitSet, char) {
			for {
				peek := l.peekChar()
				if unicode.IsSpace(peek) || peek == terminatorChar {
					break
				}
				_, ok := l.advanceChar()
				if !ok {
					break
				}
			}
			return l.lexError("invalid int literal")
		}

		if !nonSpaceCharEncountered {
			nonSpaceCharEncountered = true
		}
		result.WriteRune(char)
	}

	if endOfLiteral && result.Len() == 0 {
		l.mode = normalMode
		l.advanceChar()
		return l.token(terminatorToken)
	}
	return l.tokenWithValue(elementToken, result.String())
}

// Scan an invalid hex escape sequence in a string literal.
func (l *Lexer) scanInvalidHexEscape() *token.Token {
	l.mode = stringLiteralMode
	// advance two chars since
	// we know that `\x` has to be present
	l.advanceChars(2)
	// two more characters may be present but are not
	// guaranteed to be present, since the input may terminate at any point.
	if !l.advanceChars(2) {
		return l.lexError(unterminatedStringError)
	}

	return l.lexError(invalidHexEscapeError)
}

// Scan an invalid hex escape sequence in a string literal.
func (l *Lexer) scanInvalidBigUnicodeEscape() *token.Token {
	l.mode = stringLiteralMode
	// advance two chars since
	// we know that `\U` has to be present
	l.advanceChars(2)
	// 8 more characters may be present but are not
	// guaranteed to be present, since the input may terminate at any point.
	if !l.advanceChars(8) {
		return l.lexError(unterminatedStringError)
	}

	return l.lexError(invalidUnicodeEscapeError)
}

// Scan an invalid hex escape sequence in a string literal.
func (l *Lexer) scanInvalidUnicodeEscape() *token.Token {
	l.mode = stringLiteralMode
	// advance two chars since
	// we know that `\U` has to be present
	l.advanceChars(2)
	// 4 more characters may be present but are not
	// guaranteed to be present, since the input may terminate at any point.
	if !l.advanceChars(4) {
		return l.lexError(unterminatedStringError)
	}

	return l.lexError(invalidUnicodeEscapeError)
}

// Scan an invalid escape sequence in a string literal.
func (l *Lexer) scanInvalidEscape() *token.Token {
	l.mode = stringLiteralMode
	var lexemeBuff strings.Builder
	// advance two chars since
	// we know that `\` and a single character have to be present
	char, _ := l.advanceChar()
	lexemeBuff.WriteRune(char)

	char, _ = l.advanceChar()
	lexemeBuff.WriteRune(char)

	return l.lexError(fmt.Sprintf("invalid escape sequence `%s` in string literal", lexemeBuff.String()))
}

const (
	unterminatedStringError   = "unterminated string literal, missing `\"`"
	invalidHexEscapeError     = "invalid hex escape"
	invalidUnicodeEscapeError = "invalid unicode escape"
)

// Scan characters when inside of a string literal (after the initial `"`)
// and when the next characters aren't `"` or `}`.
func (l *Lexer) scanStringLiteralContent() *token.Token {
	var lexemeBuff strings.Builder
	for {
		char := l.peekChar()
		if char == '"' || char == '$' && l.peekNextChar() == '{' {
			return l.tokenWithValue(token.STRING_CONTENT, lexemeBuff.String())
		}

		char, ok := l.advanceChar()
		if !ok {
			return l.lexError(unterminatedStringError)
		}

		if char == '\n' {
			l.incrementLine()
		}

		if char != '\\' {
			lexemeBuff.WriteRune(char)
			continue
		}

		char, ok = l.advanceChar()
		if !ok {
			return l.lexError(unterminatedStringError)
		}
		switch char {
		case '\\':
			lexemeBuff.WriteByte('\\')
		case 'n':
			lexemeBuff.WriteByte('\n')
		case 't':
			lexemeBuff.WriteByte('\t')
		case '"':
			lexemeBuff.WriteByte('"')
		case 'r':
			lexemeBuff.WriteByte('\r')
		case 'a':
			lexemeBuff.WriteByte('\a')
		case 'b':
			lexemeBuff.WriteByte('\b')
		case 'v':
			lexemeBuff.WriteByte('\v')
		case 'f':
			lexemeBuff.WriteByte('\f')
		case 'u':
			if !l.acceptCharsN(hexLiteralChars, 4) {
				l.mode = invalidUnicodeEscapeMode
				l.backupChars(2)
				return l.tokenWithValue(token.STRING_CONTENT, lexemeBuff.String())
			}
			l.advanceChars(4)
			value, err := strconv.ParseUint(string(l.source[l.cursor-4:l.cursor]), 16, 16)
			if err != nil {
				return l.lexError(invalidUnicodeEscapeError)
			}
			lexemeBuff.WriteRune(rune(value))
		case 'U':
			if !l.acceptCharsN(hexLiteralChars, 8) {
				l.mode = invalidBigUnicodeEscapeMode
				l.backupChars(2)
				return l.tokenWithValue(token.STRING_CONTENT, lexemeBuff.String())
			}
			l.advanceChars(8)
			value, err := strconv.ParseUint(string(l.source[l.cursor-8:l.cursor]), 16, 32)
			if err != nil {
				return l.lexError(invalidUnicodeEscapeError)
			}
			lexemeBuff.WriteRune(rune(value))
		case 'x':
			if !l.acceptCharsN(hexLiteralChars, 2) {
				l.mode = invalidHexEscapeMode
				l.backupChars(2)
				return l.tokenWithValue(token.STRING_CONTENT, lexemeBuff.String())
			}
			l.advanceChars(2)
			value, err := strconv.ParseUint(string(l.source[l.cursor-2:l.cursor]), 16, 8)
			if err != nil {
				return l.lexError(invalidHexEscapeError)
			}
			lexemeBuff.WriteByte(byte(value))
		case '\n':
			l.incrementLine()
			fallthrough
		default:
			l.mode = invalidEscapeMode
			l.backupChars(2)
			return l.tokenWithValue(token.STRING_CONTENT, lexemeBuff.String())
		}
	}
}

// Scan characters when inside of a string literal (after the initial `"`)
func (l *Lexer) scanStringLiteral() *token.Token {
	char := l.peekChar()

	switch char {
	case '$':
		if l.peekNextChar() == '{' {
			l.advanceChar()
			l.advanceChar()
			l.mode = stringInterpolationMode
			return l.token(token.STRING_INTERP_BEG)
		}
	case '"':
		l.mode = normalMode
		l.advanceChar()
		return l.token(token.STRING_END)
	}

	return l.scanStringLiteralContent()
}

// Scan characters in normal mode.
func (l *Lexer) scanNormal() *token.Token {
	for {
		char, ok := l.advanceChar()
		if !ok {
			return l.token(token.END_OF_FILE)
		}

		switch char {
		case '[':
			return l.token(token.LBRACKET)
		case ']':
			return l.token(token.RBRACKET)
		case '(':
			return l.token(token.LPAREN)
		case ')':
			return l.token(token.RPAREN)
		case '{':
			return l.token(token.LBRACE)
		case '}':
			if l.mode == stringInterpolationMode {
				l.mode = stringLiteralMode
				return l.token(token.STRING_INTERP_END)
			}
			return l.token(token.RBRACE)
		case ',':
			return l.token(token.COMMA)
		case '.':
			if l.matchChar('.') {
				if l.matchChar('.') {
					return l.token(token.EXCLUSIVE_RANGE_OP)
				}
				return l.token(token.RANGE_OP)
			}
			if isDigit(l.peekChar()) {
				var lexeme strings.Builder
				lexeme.WriteByte('0')
				lexeme.WriteByte('.')
				l.consumeDigits(decimalLiteralChars, &lexeme)
				if l.matchChars("eE") {
					lexeme.WriteByte('e')
					if ok, ch := l.matchCharsRune("+-"); ok {
						lexeme.WriteRune(ch)
					}
					l.consumeDigits(decimalLiteralChars, &lexeme)
				}
				return l.tokenWithValue(token.FLOAT, lexeme.String())
			}
			return l.token(token.DOT)
		case '-':
			if l.matchChar('=') {
				return l.token(token.MINUS_EQUAL)
			}
			if l.matchChar('>') {
				return l.token(token.THIN_ARROW)
			}
			return l.token(token.MINUS)
		case '+':
			if l.matchChar('=') {
				return l.token(token.PLUS_EQUAL)
			}
			return l.token(token.PLUS)
		case '^':
			if l.matchChar('=') {
				return l.token(token.XOR_EQUAL)
			}
			return l.token(token.XOR)
		case '*':
			if l.matchChar('=') {
				return l.token(token.STAR_EQUAL)
			}
			if l.matchChar('*') {
				if l.matchChar('=') {
					return l.token(token.STAR_STAR_EQUAL)
				}
				return l.token(token.STAR_STAR)
			}
			return l.token(token.STAR)
		case '/':
			if l.matchChar('=') {
				return l.token(token.SLASH_EQUAL)
			}
			return l.token(token.SLASH)
		case '=':
			if l.matchChar('=') {
				if l.matchChar('=') {
					return l.token(token.STRICT_EQUAL)
				}
				return l.token(token.EQUAL_EQUAL)
			}
			if l.matchChar('~') {
				return l.token(token.MATCH_OP)
			}
			if l.matchChar('>') {
				return l.token(token.THICK_ARROW)
			}
			if l.peekChar() == ':' && l.peekNextChar() == '=' {
				l.advanceChar()
				l.advanceChar()
				return l.token(token.REF_EQUAL)
			}
			if l.peekChar() == '!' && l.peekNextChar() == '=' {
				l.advanceChar()
				l.advanceChar()
				return l.token(token.REF_NOT_EQUAL)
			}
			return l.token(token.EQUAL_OP)
		case ':':
			if l.matchChar(':') {
				return l.token(token.SCOPE_RES_OP)
			}
			if l.matchChar('=') {
				return l.token(token.COLON_EQUAL)
			}
			if l.matchChar('>') {
				if l.matchChar('>') {
					return l.token(token.REVERSE_INSTANCE_OF_OP)
				}
				return l.token(token.REVERSE_ISA_OP)
			}
			return l.token(token.COLON)
		case '~':
			if l.matchChar('=') {
				return l.token(token.TILDE_EQUAL)
			}
			if l.matchChar('>') {
				return l.token(token.WIGGLY_ARROW)
			}
			return l.token(token.TILDE)
		case ';':
			return l.token(token.SEMICOLON)
		case '>':
			if l.matchChar('=') {
				return l.token(token.GREATER_EQUAL)
			}
			if l.matchChar('>') {
				if l.matchChar('>') {
					if l.matchChar('=') {
						return l.token(token.RTRIPLE_BITSHIFT_EQUAL)
					}
					return l.token(token.RTRIPLE_BITSHIFT)
				}
				if l.matchChar('=') {
					return l.token(token.RBITSHIFT_EQUAL)
				}
				return l.token(token.RBITSHIFT)
			}
			return l.token(token.GREATER)
		case '<':
			if l.matchChar('=') {
				if l.matchChar('>') {
					return l.token(token.SPACESHIP_OP)
				}
				return l.token(token.LESS_EQUAL)
			}
			if l.matchChar(':') {
				return l.token(token.ISA_OP)
			}
			if l.matchChar('<') {
				if l.matchChar('<') {
					if l.matchChar('=') {
						return l.token(token.LTRIPLE_BITSHIFT_EQUAL)
					}
					return l.token(token.LTRIPLE_BITSHIFT)
				}
				if l.matchChar('=') {
					return l.token(token.LBITSHIFT_EQUAL)
				}
				if l.matchChar(':') {
					return l.token(token.INSTANCE_OF_OP)
				}
				return l.token(token.LBITSHIFT)
			}
			return l.token(token.LESS)
		case '&':
			if l.matchChar('&') {
				if l.matchChar('=') {
					return l.token(token.AND_AND_EQUAL)
				}
				return l.token(token.AND_AND)
			}
			if l.matchChar('!') {
				return l.token(token.AND_BANG)
			}
			if l.matchChar('=') {
				return l.token(token.AND_EQUAL)
			}
			return l.token(token.AND)
		case '|':
			if l.matchChar('|') {
				if l.matchChar('=') {
					return l.token(token.OR_OR_EQUAL)
				}
				return l.token(token.OR_OR)
			}
			if l.matchChar('>') {
				return l.token(token.PIPE_OP)
			}
			if l.matchChar('!') {
				return l.token(token.OR_BANG)
			}
			if l.matchChar('=') {
				return l.token(token.OR_EQUAL)
			}
			return l.token(token.OR)
		case '?':
			if l.matchChar('?') {
				if l.matchChar('=') {
					return l.token(token.QUESTION_QUESTION_EQUAL)
				}
				return l.token(token.QUESTION_QUESTION)
			}
			if l.matchChar('.') {
				return l.token(token.QUESTION_DOT)
			}
			return l.token(token.QUESTION)
		case '!':
			if l.matchChar('=') {
				if l.matchChar('=') {
					return l.token(token.STRICT_NOT_EQUAL)
				}
				return l.token(token.NOT_EQUAL)
			}
			return l.token(token.BANG)
		case '%':
			if l.matchChar('{') {
				return l.token(token.SET_LITERAL_BEG)
			}
			if l.matchChar('(') {
				return l.token(token.TUPLE_LITERAL_BEG)
			}
			if l.matchChar('=') {
				return l.token(token.PERCENT_EQUAL)
			}
			if l.matchChar('w') {
				if l.matchChar('[') {
					l.mode = wordArrayLiteralMode
					return l.token(token.WORD_LIST_BEG)
				}
				if l.matchChar('{') {
					l.mode = wordSetLiteralMode
					return l.token(token.WORD_SET_BEG)
				}
				if l.matchChar('(') {
					l.mode = wordTupleLiteralMode
					return l.token(token.WORD_TUPLE_BEG)
				}

				return l.lexError("invalid word collection literal delimiters `%w`")
			}
			if l.matchChar('s') {
				if l.matchChar('[') {
					l.mode = symbolArrayLiteralMode
					return l.token(token.SYMBOL_LIST_BEG)
				}
				if l.matchChar('{') {
					l.mode = symbolSetLiteralMode
					return l.token(token.SYMBOL_SET_BEG)
				}
				if l.matchChar('(') {
					l.mode = symbolTupleLiteralMode
					return l.token(token.SYMBOL_TUPLE_BEG)
				}

				return l.lexError("invalid symbol collection literal delimiters `%s`")
			}
			if l.matchChar('x') {
				if l.matchChar('[') {
					l.mode = hexArrayLiteralMode
					return l.token(token.HEX_LIST_BEG)
				}
				if l.matchChar('{') {
					l.mode = hexSetLiteralMode
					return l.token(token.HEX_SET_BEG)
				}
				if l.matchChar('(') {
					l.mode = hexTupleLiteralMode
					return l.token(token.HEX_TUPLE_BEG)
				}

				return l.lexError("invalid hex collection literal delimiters `%x`")
			}
			if l.matchChar('b') {
				if l.matchChar('[') {
					l.mode = binArrayLiteralMode
					return l.token(token.BIN_LIST_BEG)
				}
				if l.matchChar('{') {
					l.mode = binSetLiteralMode
					return l.token(token.BIN_SET_BEG)
				}
				if l.matchChar('(') {
					l.mode = binTupleLiteralMode
					return l.token(token.BIN_TUPLE_BEG)
				}

				return l.lexError("invalid binary collection literal delimiters `%b`")
			}
			return l.token(token.PERCENT)

		case '\n':
			l.foldNewLines()
			return l.token(token.NEWLINE)
		case '\r':
			if l.matchChar('\n') {
				l.foldNewLines()
				return l.token(token.NEWLINE)
			}
			fallthrough
		case ' ', '\t':
			l.skipByte()
		case '#':
			if l.matchChar('#') && l.matchChar('[') {
				return l.docComment()
			}

			if l.matchChar('[') {
				if tok := l.swallowBlockComments(); tok != nil {
					return tok
				}
			} else {
				l.swallowSingleLineComment()
			}
		case 'c':
			if l.matchChar('"') {
				return l.character()
			}
			if l.matchChar('\'') {
				return l.rawCharacter()
			}
			return l.publicIdentifier('c')
		case '\'':
			return l.rawString()
		case '"':
			if l.mode == stringInterpolationMode {
				for {
					_, ok := l.advanceChar()
					if !ok {
						return l.lexError(unterminatedStringError)
					}
					if l.matchChar('"') {
						break
					}
				}
				return l.lexError("unexpected string literal in string interpolation, only raw strings delimited with `'` can be used in string interpolation")
			}
			l.mode = stringLiteralMode
			return l.token(token.STRING_BEG)
		case '_':
			return l.privateIdentifier()
		case '@':
			return l.instanceVariable()
		default:
			if isDigit(char) {
				return l.numberLiteral(char)
			} else if unicode.IsLetter(char) {
				return l.publicIdentifier(char)
			}
			return l.lexError(fmt.Sprintf("unexpected character `%c`", char))
		}
	}
}

// Checks whether the given character is acceptable
// inside an publicIdentifier.
func isIdentifierChar(char rune) bool {
	return unicode.IsLetter(char) || unicode.IsNumber(char) || char == '_'
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

// Same as [token] but lets you specify the value of the token
// manually.
func (l *Lexer) tokenWithValue(typ token.Type, value string) *token.Token {
	token := token.NewWithValue(
		position.New(
			l.start,
			l.cursor-l.start,
			l.startLine,
			l.startColumn,
		),
		typ,
		value,
	)
	l.start = l.cursor
	l.startColumn = l.column
	l.startLine = l.line

	return token
}
