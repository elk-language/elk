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
	"github.com/elk-language/elk/value"
	"github.com/fatih/color"
)

// Lexing mode which changes how characters are handled by the lexer.
type mode uint8

const (
	normalMode mode = iota // Initial mode

	afterMethodCallOperatorMode

	wordArrayListLiteralMode   // Triggered after entering the initial token `\w[` of a word array literal
	symbolArrayListLiteralMode // Triggered after entering the initial token `\s[` of a symbol array literal
	hexArrayListLiteralMode    // Triggered after entering the initial token `\x[` of a hex array literal
	binArrayListLiteralMode    // Triggered after entering the initial token `\b[` of a binary array literal

	wordHashSetLiteralMode   // Triggered after entering the initial token `^w[` of a word set literal
	symbolHashSetLiteralMode // Triggered after entering the initial token `^s[` of a symbol set literal
	hexHashSetLiteralMode    // Triggered after entering the initial token `^x[` of a hex set literal
	binHashSetLiteralMode    // Triggered after entering the initial token `^b[` of a binary set literal

	wordArrayTupleLiteralMode   // Triggered after entering the initial token `%w[` of a word arrayTuple literal
	symbolArrayTupleLiteralMode // Triggered after entering the initial token `%s[` of a symbol arrayTuple literal
	hexArrayTupleLiteralMode    // Triggered after entering the initial token `%x[` of a hex arrayTuple literal
	binArrayTupleLiteralMode    // Triggered after entering the initial token `%b[` of a binary arrayTuple literal

	stringLiteralMode           // Triggered after consuming the initial token `"` of a string literal
	invalidHexEscapeMode        // Triggered after encountering an invalid hex escape sequence in a string literal
	invalidUnicodeEscapeMode    // Triggered after encountering an invalid 4 character unicode escape sequence in a string literal
	invalidBigUnicodeEscapeMode // Triggered after encountering an invalid 8 character unicode escape sequence in a string literal
	invalidEscapeMode           // Triggered after encountering an invalid escape sequence in a string literal
	stringInterpolationMode     // Triggered after consuming the initial token `${` of string interpolation

	regexLiteralMode       // Triggered after consuming the initial token `%/` of a regex literal
	regexFlagMode          // Triggered during the lexing of regex literal flags
	regexInterpolationMode // Triggered after consuming the initial token `${` of regex interpolation
)

// Holds the current state of the lexing process.
type Lexer struct {
	// Path to the source file or some name.
	sourceName string
	// Elk source code.
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
	// Current lexing mode.
	modeStack []mode
}

// Implements the Colorizer interface
type Colorizer struct{}

func (c Colorizer) Colorize(source string) (string, error) {
	return Colorize(source), nil
}

// Lex the given string and construct a new one colouring every token.
func Colorize(source string) string {
	l := New(source)

	var result strings.Builder
	var previousEnd int
	for {
		tok := l.Next()
		if tok.Type == token.END_OF_FILE {
			missing := source[previousEnd:]
			if len(missing) > 0 {
				result.WriteString(missing)
			}
			break
		}
		span := tok.Span()
		between := source[previousEnd:span.StartPos.ByteOffset]
		result.WriteString(between)

		c := color.New(tok.AnsiStyling()...)
		lexeme := source[span.StartPos.ByteOffset : span.EndPos.ByteOffset+1]
		result.WriteString(c.Sprint(lexeme))
		previousEnd = span.EndPos.ByteOffset + 1
	}
	return result.String()
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

// Lex the given string and return a slice containing all the tokens.
func LexValue(source string) *value.ArrayList {
	l := New(source)

	tokens := value.NewArrayList(10)
	for {
		tok := l.Next()
		if tok.Type == token.END_OF_FILE {
			break
		}
		tokens.Append(value.Ref(tok))
	}
	return tokens
}

// Instantiates a new lexer for the given source code.
func New(source string) *Lexer {
	return NewWithName("<main>", source)
}

// Same as [New] but lets you specify the path to the source code file.
func NewWithName(sourceName string, source string) *Lexer {
	return &Lexer{
		sourceName:  sourceName,
		source:      source,
		line:        1,
		startLine:   1,
		column:      1,
		startColumn: 1,
		modeStack:   []mode{normalMode},
	}
}

func (*Lexer) Class() *value.Class {
	return value.ElkLexerClass
}

func (*Lexer) DirectClass() *value.Class {
	return value.ElkLexerClass
}

func (l *Lexer) Inspect() string {
	return fmt.Sprintf("Std::Elk::Lexer{&: %p, source_name: %s}", l, value.String(l.sourceName).Inspect())
}

func (l *Lexer) Error() string {
	return l.Inspect()
}

func (l *Lexer) SingletonClass() *value.Class {
	return nil
}

func (l *Lexer) InstanceVariables() value.SymbolMap {
	return nil
}

func (l *Lexer) Copy() value.Reference {
	return l
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

// Return the current mode
func (l *Lexer) mode() mode {
	return l.modeStack[len(l.modeStack)-1]
}

// Set a new mode.
func (l *Lexer) pushMode(m mode) {
	l.modeStack = append(l.modeStack, m)
}

// Get back to the previous mode.
func (l *Lexer) popMode() {
	l.modeStack = l.modeStack[:len(l.modeStack)-1]
}

// Attempts to scan and construct the next token.
func (l *Lexer) scanToken() *token.Token {
	switch l.mode() {
	case normalMode, stringInterpolationMode, regexInterpolationMode:
		return l.scanNormal(false)
	case afterMethodCallOperatorMode:
		l.popMode()
		return l.scanNormal(true)
	case stringLiteralMode:
		return l.scanStringLiteral()
	case regexLiteralMode:
		return l.scanRegexLiteral()
	case regexFlagMode:
		return l.scanRegexFlag()
	case wordArrayListLiteralMode:
		return l.scanWordCollectionLiteral(token.WORD_ARRAY_LIST_END)
	case symbolArrayListLiteralMode:
		return l.scanWordCollectionLiteral(token.SYMBOL_ARRAY_LIST_END)
	case wordHashSetLiteralMode:
		return l.scanWordCollectionLiteral(token.WORD_HASH_SET_END)
	case symbolHashSetLiteralMode:
		return l.scanWordCollectionLiteral(token.SYMBOL_HASH_SET_END)
	case wordArrayTupleLiteralMode:
		return l.scanWordCollectionLiteral(token.WORD_ARRAY_TUPLE_END)
	case symbolArrayTupleLiteralMode:
		return l.scanWordCollectionLiteral(token.SYMBOL_ARRAY_TUPLE_END)
	case hexArrayListLiteralMode:
		return l.scanIntCollectionLiteral(token.HEX_ARRAY_LIST_END, hexLiteralChars, "0x")
	case hexHashSetLiteralMode:
		return l.scanIntCollectionLiteral(token.HEX_HASH_SET_END, hexLiteralChars, "0x")
	case hexArrayTupleLiteralMode:
		return l.scanIntCollectionLiteral(token.HEX_ARRAY_TUPLE_END, hexLiteralChars, "0x")
	case binArrayListLiteralMode:
		return l.scanIntCollectionLiteral(token.BIN_ARRAY_LIST_END, binaryLiteralChars, "0b")
	case binHashSetLiteralMode:
		return l.scanIntCollectionLiteral(token.BIN_HASH_SET_END, binaryLiteralChars, "0b")
	case binArrayTupleLiteralMode:
		return l.scanIntCollectionLiteral(token.BIN_ARRAY_TUPLE_END, binaryLiteralChars, "0b")
	case invalidEscapeMode:
		return l.scanInvalidEscape()
	case invalidHexEscapeMode:
		return l.scanInvalidHexEscape()
	case invalidUnicodeEscapeMode:
		return l.scanInvalidUnicodeEscape()
	case invalidBigUnicodeEscapeMode:
		return l.scanInvalidBigUnicodeEscape()
	default:
		return l.lexError(fmt.Sprintf("unsupported lexing mode `%d`", l.mode()))
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

// Checks if the second next character matches the given one.
func (l *Lexer) acceptNextChar(char rune) bool {
	if !l.hasMoreTokens() {
		return false
	}

	if l.peekNextChar() == char {
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

// Assumes that `##[` has already been consumed.
// Builds the doc comment token.
func (l *Lexer) hashDocComment() *token.Token {
	nestCounter := 1
	docStrLines := make([]string, 0, 3)
	var lineBuffer strings.Builder

	leastIndented := math.MaxInt
	indent := 0
	nonIndentChars := false
	l.swallowNewLines()
charLoop:
	for {
		char, ok := l.advanceChar()
		if !ok {
			return l.lexError(fmt.Sprintf("unbalanced doc comments, expected %d more doc comment ending(s) `]##`", nestCounter))
		}

	charSwitch:
		switch char {
		case '#':
			nonIndentChars = true
			if l.matchChar('#') {
				if l.matchChar('[') {
					lineBuffer.WriteString("##[")
					nestCounter += 1
					break charSwitch
				}
				lineBuffer.WriteString("##")
				break charSwitch
			}
			lineBuffer.WriteString("#")
		case ']':
			nonIndentChars = true
			if l.matchChar('#') {
				if l.matchChar('#') {
					nestCounter -= 1
					if nestCounter == 0 {
						docStrLines = append(docStrLines, lineBuffer.String())
						lineBuffer.Reset()
						break charLoop
					}
					lineBuffer.WriteString("]##")
					break charSwitch
				}
				lineBuffer.WriteString("]#")
				break charSwitch
			}
			lineBuffer.WriteString("]")
		default:
			lineBuffer.WriteRune(char)
		}

		if !nonIndentChars && char == ' ' || char == '\t' {
			indent += 1
		} else if l.isNewLine(char) {
			l.incrementLine()
			docStrLines = append(docStrLines, lineBuffer.String())
			lineBuffer.Reset()
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
	var resultBuffer strings.Builder
	for _, line := range docStrLines {
		// add 1 because of the trailing newline
		if len(line) < leastIndented+1 {
			resultBuffer.WriteRune('\n')
			continue
		}

		resultBuffer.WriteString(line[leastIndented:])
	}

	return l.tokenWithValue(token.DOC_COMMENT, strings.TrimRight(resultBuffer.String(), "\t\n "))
}

// Assumes that `/**` has already been consumed.
// Builds the doc comment token.
func (l *Lexer) slashDocComment() *token.Token {
	nestCounter := 1
	docStrLines := make([]string, 0, 3)
	var lineBuffer strings.Builder

	leastIndented := math.MaxInt
	indent := 0
	nonIndentChars := false
	l.swallowNewLines()
charLoop:
	for {
		char, ok := l.advanceChar()
		if !ok {
			return l.lexError(fmt.Sprintf("unbalanced doc comments, expected %d more doc comment ending(s) `**/`", nestCounter))
		}

	charSwitch:
		switch char {
		case '/':
			nonIndentChars = true
			if l.matchChar('*') {
				if l.matchChar('*') {
					lineBuffer.WriteString("/**")
					nestCounter += 1
					break charSwitch
				}
				lineBuffer.WriteString("/*")
				break charSwitch
			}
			lineBuffer.WriteString("/")
		case '*':
			nonIndentChars = true
			if l.matchChar('*') {
				if l.matchChar('/') {
					nestCounter -= 1
					if nestCounter == 0 {
						docStrLines = append(docStrLines, lineBuffer.String())
						lineBuffer.Reset()
						break charLoop
					}
					lineBuffer.WriteString("**/")
					break charSwitch
				}
				lineBuffer.WriteString("**")
				break charSwitch
			}
			lineBuffer.WriteString("*")
		default:
			lineBuffer.WriteRune(char)
		}

		if !nonIndentChars && char == ' ' || char == '\t' {
			indent += 1
		} else if l.isNewLine(char) {
			l.incrementLine()
			docStrLines = append(docStrLines, lineBuffer.String())
			lineBuffer.Reset()
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
	var resultBuffer strings.Builder
	for _, line := range docStrLines {
		// add 1 because of the trailing newline
		if len(line) < leastIndented+1 {
			resultBuffer.WriteRune('\n')
			continue
		}

		resultBuffer.WriteString(line[leastIndented:])
	}

	return l.tokenWithValue(token.DOC_COMMENT, strings.TrimRight(resultBuffer.String(), "\t\n "))
}

// Assumes that "#" or "//" has already been consumed.
// Skips over a single line comment "#", "//" ...
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
func (l *Lexer) swallowHashBlockComments() *token.Token {
	nestCounter := 1
charLoop:
	for {
		char, ok := l.advanceChar()
		if !ok {
			return l.lexError(fmt.Sprintf("unbalanced block comments, expected %d more block comment ending(s) `]#`", nestCounter))
		}

		switch char {
		case '#':
			if l.matchChar('[') {
				nestCounter += 1
			}
		case ']':
			if l.matchChar('#') {
				nestCounter -= 1
				if nestCounter == 0 {
					break charLoop
				}
			}
		default:
			if l.isNewLine(char) {
				l.incrementLine()
			}
		}
	}

	l.skipToken()
	return nil
}

// Assumes that "/*" has already been consumed.
// Skips over a block comment "/*" ... "*/".
func (l *Lexer) swallowSlashBlockComments() *token.Token {
	nestCounter := 1
charLoop:
	for {
		char, ok := l.advanceChar()
		if !ok {
			return l.lexError(fmt.Sprintf("unbalanced block comments, expected %d more block comment ending(s) `*/`", nestCounter))
		}

		switch char {
		case '/':
			if l.matchChar('*') {
				nestCounter += 1
			}
		case '*':
			if l.matchChar('/') {
				nestCounter -= 1
				if nestCounter == 0 {
					break charLoop
				}
			}
		default:
			if l.isNewLine(char) {
				l.incrementLine()
			}
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

const unterminatedCharLiteralMessage = "unterminated character literal, missing backtick"

const charTerminator = '`'

// Assumes that the beginning ` has already been consumed.
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
		case charTerminator:
			lexemeBuff.WriteByte(charTerminator)
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
				if !l.swallowUntil(charTerminator) {
					return l.lexError(unterminatedCharLiteralMessage)
				}
				return l.lexError(invalidUnicodeEscapeError)
			}
			value, err := strconv.ParseUint(string(l.source[l.cursor-4:l.cursor]), 16, 16)
			if err != nil {
				if !l.swallowUntil(charTerminator) {
					return l.lexError(unterminatedCharLiteralMessage)
				}
				return l.lexError(invalidHexEscapeError)
			}
			lexemeBuff.WriteRune(rune(value))
		case 'U':
			if !l.matchCharsN(hexLiteralChars, 8) {
				if !l.swallowUntil(charTerminator) {
					return l.lexError(unterminatedCharLiteralMessage)
				}
				return l.lexError(invalidUnicodeEscapeError)
			}
			value, err := strconv.ParseUint(string(l.source[l.cursor-8:l.cursor]), 16, 32)
			if err != nil {
				if !l.swallowUntil(charTerminator) {
					return l.lexError(unterminatedCharLiteralMessage)
				}
				return l.lexError(invalidHexEscapeError)
			}
			lexemeBuff.WriteRune(rune(value))
		case 'x':
			if !l.matchCharsN(hexLiteralChars, 2) {
				if !l.swallowUntil(charTerminator) {
					return l.lexError(unterminatedCharLiteralMessage)
				}
				return l.lexError(invalidHexEscapeError)
			}
			value, err := strconv.ParseUint(string(l.source[l.cursor-2:l.cursor]), 16, 8)
			if err != nil {
				if !l.swallowUntil(charTerminator) {
					return l.lexError(unterminatedCharLiteralMessage)
				}
				return l.lexError(invalidHexEscapeError)
			}
			lexemeBuff.WriteByte(byte(value))
		case '\n':
			l.incrementLine()
			fallthrough
		default:
			l.matchChar(charTerminator)
			return l.lexError("invalid escape sequence in a character literal")
		}
	} else {
		ch, ok := l.advanceChar()
		if !ok {
			return l.lexError(unterminatedCharLiteralMessage)
		}
		lexemeBuff.WriteRune(ch)
	}
	if l.matchChar(charTerminator) {
		return l.tokenWithValue(token.CHAR_LITERAL, lexemeBuff.String())
	}

	if !l.swallowUntil(charTerminator) {
		return l.lexError(unterminatedCharLiteralMessage)
	}

	return l.lexError("invalid char literal with more than one character")
}

// Assumes that the beginning r` has already been consumed.
// Consumes a raw character literal.
func (l *Lexer) rawCharacter() *token.Token {
	var char string

	ch, ok := l.advanceChar()
	if !ok {
		return l.lexError(unterminatedCharLiteralMessage)
	}
	char = string(ch)
	if l.matchChar(charTerminator) {
		return l.tokenWithValue(token.RAW_CHAR_LITERAL, char)
	}

	if !l.swallowUntil(charTerminator) {
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
	nonDecimal := false
	digits := decimalLiteralChars
	var lexeme strings.Builder
	lexeme.WriteRune(startDigit)

	if startDigit == '0' {
		switch l.peekChar() {
		case 'x', 'X':
			// hexadecimal (base 16)
			l.advanceChar()
			lexeme.WriteRune('x')
			digits = hexLiteralChars
			nonDecimal = true
		case 'd', 'D':
			// duodecimal (base 12)
			l.advanceChar()
			lexeme.WriteRune('d')
			digits = duodecimalLiteralChars
			nonDecimal = true
		case 'o', 'O':
			// octal (base 8)
			l.advanceChar()
			lexeme.WriteRune('o')
			digits = octalLiteralChars
			nonDecimal = true
		case 'q', 'Q':
			// quaternary (base 4)
			l.advanceChar()
			lexeme.WriteRune('q')
			digits = quaternaryLiteralChars
			nonDecimal = true
		case 'b', 'B':
			// binary (base 2)
			l.advanceChar()
			lexeme.WriteRune('b')
			digits = binaryLiteralChars
			nonDecimal = true
		}
	}

	l.consumeDigits(digits, &lexeme)

	switch l.peekChar() {
	case 'i':
		l.advanceChar()
		switch ch, _ := l.advanceChar(); ch {
		case '6':
			if l.matchChar('4') {
				return l.tokenWithValue(token.INT64, lexeme.String())
			}
		case '3':
			if l.matchChar('2') {
				return l.tokenWithValue(token.INT32, lexeme.String())
			}
		case '1':
			if l.matchChar('6') {
				return l.tokenWithValue(token.INT16, lexeme.String())
			}
		case '8':
			return l.tokenWithValue(token.INT8, lexeme.String())
		}
		return l.lexError("invalid sized integer literal")
	case 'u':
		l.advanceChar()
		switch ch, _ := l.advanceChar(); ch {
		case '6':
			if l.matchChar('4') {
				return l.tokenWithValue(token.UINT64, lexeme.String())
			}
		case '3':
			if l.matchChar('2') {
				return l.tokenWithValue(token.UINT32, lexeme.String())
			}
		case '1':
			if l.matchChar('6') {
				return l.tokenWithValue(token.UINT16, lexeme.String())
			}
		case '8':
			return l.tokenWithValue(token.UINT8, lexeme.String())
		}
		return l.lexError("invalid sized integer literal")
	}
	if nonDecimal {
		return l.tokenWithValue(token.INT, lexeme.String())
	}

	tokenType := token.INT
	if l.acceptChar('.') && isDigit(l.peekNextChar()) {
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

	if l.matchChar('f') {
		switch ch, _ := l.advanceChar(); ch {
		case '6':
			if l.matchChar('4') {
				return l.tokenWithValue(token.FLOAT64, lexeme.String())
			}
		case '3':
			if l.matchChar('2') {
				return l.tokenWithValue(token.FLOAT32, lexeme.String())
			}
		}
		return l.lexError("invalid sized float literal")
	}
	if l.matchChar('b') {
		if l.matchChar('f') {
			return l.tokenWithValue(token.BIG_FLOAT, lexeme.String())
		}
		return l.lexError("invalid big numeric literal")
	}

	return l.tokenWithValue(tokenType, lexeme.String())
}

// Assumes that the initial letter has already been consumed.
// Consumes and constructs a public publicIdentifier token.
func (l *Lexer) publicIdentifier(init rune, afterMethodCallOperator bool) *token.Token {
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
		if afterMethodCallOperator {
			l.matchChar('=')
		}
		return l.tokenWithConsumedValue(token.PUBLIC_IDENTIFIER)
	}
}

// Assumes that the initial "_" has already been consumed.
// Consumes and constructs a private publicIdentifier token.
func (l *Lexer) privateIdentifier(afterMethodCallOperator bool) *token.Token {
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
		if afterMethodCallOperator {
			l.matchChar('=')
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

	lexeme := l.source[l.start+1 : l.cursor]
	if len(lexeme) == 0 {
		return l.lexError("empty instance variable name")
	}
	return l.tokenWithValue(token.INSTANCE_VARIABLE, string(lexeme))
}

// Assumes that the initial `$` has been consumed.
// Consumes and constructs a special identifier token.
func (l *Lexer) specialIdentifier() *token.Token {
	for isIdentifierChar(l.peekChar()) {
		l.advanceChar()
	}

	lexeme := l.source[l.start+1 : l.cursor]
	if len(lexeme) == 0 {
		return l.lexError("empty special identifier")
	}
	return l.tokenWithValue(token.SPECIAL_IDENTIFIER, string(lexeme))
}

const (
	unterminatedCollectionError = "unterminated %s literal, missing `%c`"
)

// Scans the content of word collection literals be it `\w[`, `\s[`, `%w[`, `%s[`, `^w[`, `^s[`
func (l *Lexer) scanWordCollectionLiteral(terminatorToken token.Type) *token.Token {
	var result strings.Builder
	var nonSpaceCharEncountered bool
	var endOfLiteral bool
	terminatorChar := ']'

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
		l.popMode()
		l.advanceChar()
		return l.token(terminatorToken)
	}
	return l.tokenWithValue(token.RAW_STRING, result.String())
}

// Scans the content of int collection literals be it `\x[`, `\b[`, `%x[`, `%b[`, `^x[`, `^b[`
func (l *Lexer) scanIntCollectionLiteral(terminatorToken token.Type, digitSet string, prefix string) *token.Token {
	var result strings.Builder
	var nonSpaceCharEncountered bool
	var endOfLiteral bool
	terminatorChar := ']'

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
			result.WriteString(prefix)
		}
		result.WriteRune(char)
	}

	if endOfLiteral && result.Len() == 0 {
		l.popMode()
		l.advanceChar()
		return l.token(terminatorToken)
	}
	return l.tokenWithValue(token.INT, result.String())
}

// Scan an invalid hex escape sequence in a string literal.
func (l *Lexer) scanInvalidHexEscape() *token.Token {
	l.popMode()
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

// Scan an invalid 8-character unicode escape sequence in a string/char literal.
func (l *Lexer) scanInvalidBigUnicodeEscape() *token.Token {
	l.popMode()
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

// Scan an invalid 4-character unicode escape sequence in a string/char literal.
func (l *Lexer) scanInvalidUnicodeEscape() *token.Token {
	l.popMode()
	// advance two chars since
	// we know that `\u` has to be present
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
	l.popMode()
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
		nextChar := l.peekNextChar()
		if char == '"' ||
			char == '$' && (nextChar == '{' || nextChar == '_' || unicode.IsLetter(nextChar)) ||
			char == '#' && (nextChar == '{' || nextChar == '_' || unicode.IsLetter(nextChar)) {
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
		case '$':
			lexemeBuff.WriteByte('$')
		case '#':
			lexemeBuff.WriteByte('#')
		case 'u':
			if !l.acceptCharsN(hexLiteralChars, 4) {
				l.pushMode(invalidUnicodeEscapeMode)
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
				l.pushMode(invalidBigUnicodeEscapeMode)
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
				l.pushMode(invalidHexEscapeMode)
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
			l.pushMode(invalidEscapeMode)
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
		nextChar := l.peekNextChar()
		if nextChar == '{' {
			l.advanceChar() // $
			l.advanceChar() // {
			l.pushMode(stringInterpolationMode)
			return l.token(token.STRING_INTERP_BEG)
		}
		if nextChar == '_' {
			// private identifier
			l.advanceChar() // $
			l.advanceChar() // _
			var tokenType token.Type
			if unicode.IsUpper(l.peekChar()) {
				tokenType = token.STRING_INTERP_CONSTANT
			} else {
				tokenType = token.STRING_INTERP_LOCAL
			}
			var buffer strings.Builder
			buffer.WriteRune('_')
			for isIdentifierChar(l.peekChar()) {
				ch, _ := l.advanceChar()
				buffer.WriteRune(ch)
			}
			return l.tokenWithValue(tokenType, buffer.String())
		} else if unicode.IsLetter(nextChar) {
			// public identifier
			l.advanceChar() // $
			var tokenType token.Type
			if unicode.IsUpper(nextChar) {
				tokenType = token.STRING_INTERP_CONSTANT
			} else {
				tokenType = token.STRING_INTERP_LOCAL
			}
			var buffer strings.Builder
			for isIdentifierChar(l.peekChar()) {
				ch, _ := l.advanceChar()
				buffer.WriteRune(ch)
			}
			return l.tokenWithValue(tokenType, buffer.String())
		}
	case '#':
		nextChar := l.peekNextChar()
		if nextChar == '{' {
			l.advanceChar() // #
			l.advanceChar() // {
			l.pushMode(stringInterpolationMode)
			return l.token(token.STRING_INSPECT_INTERP_BEG)
		}
		if nextChar == '_' {
			// private identifier
			l.advanceChar() // #
			l.advanceChar() // _
			var tokenType token.Type
			if unicode.IsUpper(l.peekChar()) {
				tokenType = token.STRING_INSPECT_INTERP_CONSTANT
			} else {
				tokenType = token.STRING_INSPECT_INTERP_LOCAL
			}
			var buffer strings.Builder
			buffer.WriteRune('_')
			for isIdentifierChar(l.peekChar()) {
				ch, _ := l.advanceChar()
				buffer.WriteRune(ch)
			}
			return l.tokenWithValue(tokenType, buffer.String())
		} else if unicode.IsLetter(nextChar) {
			// public identifier
			l.advanceChar() // #
			var tokenType token.Type
			if unicode.IsUpper(nextChar) {
				tokenType = token.STRING_INSPECT_INTERP_CONSTANT
			} else {
				tokenType = token.STRING_INSPECT_INTERP_LOCAL
			}
			var buffer strings.Builder
			for isIdentifierChar(l.peekChar()) {
				ch, _ := l.advanceChar()
				buffer.WriteRune(ch)
			}
			return l.tokenWithValue(tokenType, buffer.String())
		}
	case '"':
		l.popMode()
		l.advanceChar()
		return l.token(token.STRING_END)
	}

	return l.scanStringLiteralContent()
}

const (
	unterminatedRegexError = "unterminated regex literal, missing `/`"
)

// Scan characters when inside of a regex literal (after the initial `%/`)
// and when the next characters aren't `/`.
func (l *Lexer) scanRegexLiteralContent() *token.Token {
	var lexemeBuff strings.Builder
	for {
		char := l.peekChar()
		if char == '/' || char == '$' && l.peekNextChar() == '{' {
			return l.tokenWithValue(token.REGEX_CONTENT, lexemeBuff.String())
		}

		char, ok := l.advanceChar()
		if !ok {
			return l.lexError(unterminatedRegexError)
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
			return l.lexError(unterminatedRegexError)
		}
		switch char {
		case '/':
			lexemeBuff.WriteString(`\/`)
		case '\n':
			l.incrementLine()
			fallthrough
		default:
			lexemeBuff.WriteRune('\\')
			lexemeBuff.WriteRune(char)
		}
	}
}

// Scan flags after the ending `/` of a regex literal
func (l *Lexer) scanRegexFlag() *token.Token {
	if !unicode.IsLetter(l.peekNextChar()) {
		l.popMode()
		l.popMode()
	}

	char, ok := l.advanceChar()
	if !ok {
		return l.token(token.END_OF_FILE)
	}
	switch char {
	case 'i':
		return l.token(token.REGEX_FLAG_i)
	case 'm':
		return l.token(token.REGEX_FLAG_m)
	case 'U':
		return l.token(token.REGEX_FLAG_U)
	case 'a':
		return l.token(token.REGEX_FLAG_a)
	case 'x':
		return l.token(token.REGEX_FLAG_x)
	case 's':
		return l.token(token.REGEX_FLAG_s)
	default:
		return l.lexError("invalid regex flag")
	}
}

// Scan characters when inside of a regex literal (after the initial `%/`)
func (l *Lexer) scanRegexLiteral() *token.Token {
	char := l.peekChar()

	switch char {
	case '$':
		if l.peekNextChar() == '{' {
			l.advanceChar()
			l.advanceChar()
			l.pushMode(regexInterpolationMode)
			return l.token(token.REGEX_INTERP_BEG)
		}
	case '/':
		if unicode.IsLetter(l.peekNextChar()) {
			l.pushMode(regexFlagMode)
		} else {
			l.popMode()
		}
		l.advanceChar()
		return l.token(token.REGEX_END)
	}

	return l.scanRegexLiteralContent()
}

// Scan characters in normal mode.
func (l *Lexer) scanNormal(afterMethodCallOperator bool) *token.Token {
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
			switch l.mode() {
			case stringInterpolationMode:
				l.popMode()
				return l.token(token.STRING_INTERP_END)
			case regexInterpolationMode:
				l.popMode()
				return l.token(token.REGEX_INTERP_END)
			default:
				return l.token(token.RBRACE)
			}
		case ',':
			return l.token(token.COMMA)
		case '.':
			if l.matchChar(':') {
				return l.token(token.DOT_COLON)
			}
			if l.acceptChar('.') {
				if l.acceptNextChar('.') {
					l.advanceChars(2)
					return l.token(token.CLOSED_RANGE_OP)
				}
				if l.acceptNextChar('<') {
					l.advanceChars(2)
					return l.token(token.RIGHT_OPEN_RANGE_OP)
				}
				l.advanceChar()
				l.pushMode(afterMethodCallOperatorMode)
				return l.token(token.DOT_DOT)
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
				if l.matchChar('f') {
					switch ch, _ := l.advanceChar(); ch {
					case '6':
						if l.matchChar('4') {
							return l.tokenWithValue(token.FLOAT64, lexeme.String())
						}
					case '3':
						if l.matchChar('2') {
							return l.tokenWithValue(token.FLOAT32, lexeme.String())
						}
					}
					return l.lexError("invalid sized float literal")
				}
				if l.matchChar('b') {
					if l.matchChar('f') {
						return l.tokenWithValue(token.BIG_FLOAT, lexeme.String())
					}
					return l.lexError("invalid big numeric literal")
				}
				return l.tokenWithValue(token.FLOAT, lexeme.String())
			}

			l.pushMode(afterMethodCallOperatorMode)
			return l.token(token.DOT)
		case '-':
			if l.matchChar('=') {
				return l.token(token.MINUS_EQUAL)
			}
			if l.matchChar('@') {
				return l.token(token.MINUS_AT)
			}
			if l.matchChar('>') {
				return l.token(token.THIN_ARROW)
			}
			if l.matchChar('-') {
				return l.token(token.MINUS_MINUS)
			}
			return l.token(token.MINUS)
		case '+':
			if l.matchChar('=') {
				return l.token(token.PLUS_EQUAL)
			}
			if l.matchChar('@') {
				return l.token(token.PLUS_AT)
			}
			if l.matchChar('+') {
				return l.token(token.PLUS_PLUS)
			}
			return l.token(token.PLUS)
		case '^':
			if l.matchChar('=') {
				return l.token(token.XOR_EQUAL)
			}
			if l.matchChar('[') {
				return l.token(token.HASH_SET_LITERAL_BEG)
			}
			if l.acceptChar('w') {
				if l.acceptNextChar('[') {
					l.advanceChars(2)
					l.pushMode(wordHashSetLiteralMode)
					return l.token(token.WORD_HASH_SET_BEG)
				}
			}
			if l.acceptChar('s') {
				if l.acceptNextChar('[') {
					l.advanceChars(2)
					l.pushMode(symbolHashSetLiteralMode)
					return l.token(token.SYMBOL_HASH_SET_BEG)
				}
			}
			if l.acceptChar('x') {
				if l.acceptNextChar('[') {
					l.advanceChars(2)
					l.pushMode(hexHashSetLiteralMode)
					return l.token(token.HEX_HASH_SET_BEG)
				}
			}
			if l.acceptChar('b') {
				if l.acceptNextChar('[') {
					l.advanceChars(2)
					l.pushMode(binHashSetLiteralMode)
					return l.token(token.BIN_HASH_SET_BEG)
				}
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
			if l.matchChar('/') {
				l.swallowSingleLineComment()
			} else if l.matchChar('*') {
				if l.matchChar('*') {
					return l.slashDocComment()
				}
				if tok := l.swallowSlashBlockComments(); tok != nil {
					return tok
				}
			} else if l.matchChar('=') {
				return l.token(token.SLASH_EQUAL)
			} else {
				return l.token(token.SLASH)
			}
		case '=':
			if l.matchChar('=') {
				if l.matchChar('=') {
					return l.token(token.STRICT_EQUAL)
				}
				return l.token(token.EQUAL_EQUAL)
			}
			if l.matchChar('~') {
				return l.token(token.LAX_EQUAL)
			}
			if l.matchChar('>') {
				return l.token(token.THICK_ARROW)
			}
			return l.token(token.EQUAL_OP)
		case ':':
			if l.matchChar(':') {
				if l.matchChar('[') {
					return l.token(token.COLON_COLON_LBRACKET)
				}
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
			if l.acceptChar('.') {
				if l.acceptNextChar('.') {
					l.advanceChars(2)
					return l.token(token.LEFT_OPEN_RANGE_OP)
				}
				if l.acceptNextChar('<') {
					l.advanceChars(2)
					return l.token(token.OPEN_RANGE_OP)
				}
			}
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
			if l.matchChar('~') {
				return l.token(token.AND_TILDE)
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
				l.pushMode(afterMethodCallOperatorMode)
				if l.matchChar('.') {
					return l.token(token.QUESTION_DOT_DOT)
				}
				return l.token(token.QUESTION_DOT)
			}
			if l.matchChar('[') {
				return l.token(token.QUESTION_LBRACKET)
			}
			return l.token(token.QUESTION)
		case '!':
			if l.matchChar('=') {
				if l.matchChar('=') {
					return l.token(token.STRICT_NOT_EQUAL)
				}
				return l.token(token.NOT_EQUAL)
			}
			if l.matchChar('~') {
				return l.token(token.LAX_NOT_EQUAL)
			}
			return l.token(token.BANG)
		case '\\':
			if l.acceptChar('w') {
				if l.acceptNextChar('[') {
					l.advanceChars(2)
					l.pushMode(wordArrayListLiteralMode)
					return l.token(token.WORD_ARRAY_LIST_BEG)
				}
			}
			if l.acceptChar('s') {
				if l.acceptNextChar('[') {
					l.advanceChars(2)
					l.pushMode(symbolArrayListLiteralMode)
					return l.token(token.SYMBOL_ARRAY_LIST_BEG)
				}
			}
			if l.acceptChar('x') {
				if l.acceptNextChar('[') {
					l.advanceChars(2)
					l.pushMode(hexArrayListLiteralMode)
					return l.token(token.HEX_ARRAY_LIST_BEG)
				}
			}
			if l.acceptChar('b') {
				if l.acceptNextChar('[') {
					l.advanceChars(2)
					l.pushMode(binArrayListLiteralMode)
					return l.token(token.BIN_ARRAY_LIST_BEG)
				}
			}
			fallthrough
		case '%':
			if l.matchChar('/') {
				l.pushMode(regexLiteralMode)
				return l.token(token.REGEX_BEG)
			}
			if l.matchChar('[') {
				return l.token(token.TUPLE_LITERAL_BEG)
			}
			if l.matchChar('{') {
				return l.token(token.RECORD_LITERAL_BEG)
			}
			if l.matchChar('=') {
				return l.token(token.PERCENT_EQUAL)
			}
			if l.acceptChar('w') {
				if l.acceptNextChar('[') {
					l.advanceChars(2)
					l.pushMode(wordArrayTupleLiteralMode)
					return l.token(token.WORD_ARRAY_TUPLE_BEG)
				}
			}
			if l.acceptChar('s') {
				if l.acceptNextChar('[') {
					l.advanceChars(2)
					l.pushMode(symbolArrayTupleLiteralMode)
					return l.token(token.SYMBOL_ARRAY_TUPLE_BEG)
				}
			}
			if l.acceptChar('x') {
				if l.acceptNextChar('[') {
					l.advanceChars(2)
					l.pushMode(hexArrayTupleLiteralMode)
					return l.token(token.HEX_ARRAY_TUPLE_BEG)
				}
			}
			if l.acceptChar('b') {
				if l.acceptNextChar('[') {
					l.advanceChars(2)
					l.pushMode(binArrayTupleLiteralMode)
					return l.token(token.BIN_ARRAY_TUPLE_BEG)
				}
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
			if l.acceptChar('#') && l.acceptNextChar('[') {
				l.advanceChars(2)
				return l.hashDocComment()
			}

			if l.matchChar('[') {
				if tok := l.swallowHashBlockComments(); tok != nil {
					return tok
				}
			} else {
				l.swallowSingleLineComment()
			}
		case '`':
			return l.character()
		case 'r':
			if l.matchChar('`') {
				return l.rawCharacter()
			}
			return l.publicIdentifier('r', afterMethodCallOperator)
		case '\'':
			return l.rawString()
		case '"':
			if l.mode() == stringInterpolationMode {
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
			l.pushMode(stringLiteralMode)
			return l.token(token.STRING_BEG)
		case '_':
			return l.privateIdentifier(afterMethodCallOperator)
		case '@':
			return l.instanceVariable()
		case '$':
			return l.specialIdentifier()
		default:
			if isDigit(char) {
				return l.numberLiteral(char)
			} else if unicode.IsLetter(char) {
				return l.publicIdentifier(char, afterMethodCallOperator)
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
		position.NewLocation(
			l.sourceName,
			position.NewSpan(
				startPos,
				endPos,
			),
		),
		typ,
		value,
	)
	l.start = l.cursor
	l.startColumn = l.column
	l.startLine = l.line

	return token
}
