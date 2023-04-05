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
func (l *lexer) Next() *Token {
	if !l.hasMoreTokens() {
		return newEOF()
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
func (l *lexer) docComment() *Token {
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
					if indent < leastIndented {
						leastIndented = indent
					}
					break
				}
				docStrLines[docStrLine] += "]##"
			}
		}
		char, ok := l.advanceChar()
		if !ok {
			return l.tokenWithValue(ErrorToken, l.lexErrorWithHint(fmt.Sprintf("unbalanced doc comments, expected %d more doc comment ending(s) `%s`", nestCounter, warnFmt.Sprint("]##"))))
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

	return l.tokenWithValue(DocCommentToken, result)
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
func (l *lexer) swallowBlockComments() *Token {
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
			return l.tokenWithValue(ErrorToken, l.lexErrorWithHint(fmt.Sprintf("unbalanced block comments, expected %d more block comment ending(s) `%s`", nestCounter, warnFmt.Sprint("]#"))))
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
func (l *lexer) rawString() *Token {
	var result string
	for {
		char, ok := l.advanceChar()
		if !ok {
			return l.tokenWithValue(ErrorToken, l.lexErrorWithHint(fmt.Sprintf("unterminated raw string, missing `%s`", warnFmt.Sprint("'"))))
		}
		if char == '\'' {
			break
		}
		if char == '\n' {
			l.incrementLine()
		}
		result += string(char)
	}

	return l.tokenWithValue(RawStringToken, result)
}

// Assumes that the first digit has already been consumed.
// Consumes and constructs an Int or Float literal token.
func (l *lexer) numberLiteral(startDigit rune) *Token {
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
		return l.tokenWithConsumedValue(IntToken)
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
		return l.tokenWithConsumedValue(FloatToken)
	}
	return l.tokenWithConsumedValue(IntToken)
}

// Assumes that the initial letter has already been consumed.
// Consumes and constructs a public identifier token.
func (l *lexer) identifier(init rune) *Token {
	if unicode.IsUpper(init) {
		// constant
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		return l.tokenWithConsumedValue(ConstantToken)
	} else {
		// variable or method name
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		if lexType := keywords[l.tokenValue()]; lexType.isKeyword() {
			// Is a keyword
			return l.token(lexType)
		}
		return l.tokenWithConsumedValue(IdentifierToken)
	}
}

// Assumes that the initial "_" has already been consumed.
// Consumes and constructs a private identifier token.
func (l *lexer) privateIdentifier() *Token {
	if unicode.IsUpper(l.peekChar()) {
		// constant
		l.advanceChar()
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		return l.tokenWithConsumedValue(PrivateConstantToken)
	} else if unicode.IsLower(l.peekChar()) {
		// variable or method name
		l.advanceChar()
		for isIdentifierChar(l.peekChar()) {
			l.advanceChar()
		}
		return l.tokenWithConsumedValue(PrivateIdentifierToken)
	}

	return l.tokenWithConsumedValue(PrivateIdentifierToken)
}

// Assumes that the initial `@` has been consumed.
// Consumes and constructs an instance variable token.
func (l *lexer) instanceVariable() *Token {
	for isIdentifierChar(l.peekChar()) {
		l.advanceChar()
	}

	return l.tokenWithValue(InstanceVariableToken, string(l.source[l.start+1:l.cursor]))
}

// Attempts to scan and construct the next token.
func (l *lexer) scanToken() *Token {
	switch l.mode {
	case inStringLiteralMode:
		return l.scanStringLiteral()
	case inStringInterpolationMode:
		fallthrough
	case normalMode:
		return l.scanNormal()
	default:
		return l.tokenWithValue(ErrorToken, l.lexErrorWithHint("unsupported lexing mode"))
	}
}

const hexChars = "0123456789abcdefABCDEF"

// Scan characters when inside of a string literal (after the initial `"`)
// and when the next characters aren't `"` or `}`.
func (l *lexer) scanStringLiteralContent() *Token {
	var lexemeBuff strings.Builder
	for {
		char := l.peekChar()
		if char == '"' || char == '$' && l.peekNextChar() == '{' {
			return l.tokenWithValue(StringContentToken, lexemeBuff.String())
		}

		char, ok := l.advanceChar()
		if !ok {
			return l.tokenWithValue(ErrorToken, l.lexErrorWithHint("unterminated string literal"))
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
			return l.tokenWithValue(ErrorToken, l.lexErrorWithHint("unterminated string literal"))
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
		case 'x':
			if !l.acceptChars(hexChars) || !l.acceptChars(hexChars) {
				// TODO: hint should be based on the current cursor
				return l.tokenWithValue(ErrorToken, l.lexErrorWithHint("invalid hex escape"))
			}
			value, err := strconv.ParseUint(string(l.source[l.cursor-2:l.cursor]), 16, 8)
			if err != nil {
				return l.tokenWithValue(ErrorToken, l.lexErrorWithHint("invalid hex escape"))
			}
			byteValue := byte(value)
			lexemeBuff.WriteByte(byteValue)
		case '\n':
			l.incrementLine()
			fallthrough
		default:
			lexemeBuff.WriteByte('\\')
			lexemeBuff.WriteRune(char)
		}
	}
}

// Scan characters when inside of a string literal (after the initial `"`)
func (l *lexer) scanStringLiteral() *Token {
	char := l.peekChar()

	switch char {
	case '$':
		if l.peekNextChar() == '{' {
			l.advanceChar()
			l.advanceChar()
			l.mode = inStringInterpolationMode
			return l.token(StringInterpBegToken)
		}
	case '"':
		l.mode = normalMode
		l.advanceChar()
		return l.token(StringEndToken)
	}

	return l.scanStringLiteralContent()
}

// Scan characters in normal mode.
func (l *lexer) scanNormal() *Token {
	for {
		char, ok := l.advanceChar()
		if !ok {
			return l.tokenWithValue(ErrorToken, l.lexError("unexpected end of file"))
		}

		switch char {
		case '[':
			return l.token(LBracketToken)
		case ']':
			return l.token(RBracketToken)
		case '(':
			return l.token(LParenToken)
		case ')':
			return l.token(RParenToken)
		case '{':
			return l.token(LBraceToken)
		case '}':
			if l.mode == inStringInterpolationMode {
				l.mode = inStringLiteralMode
				return l.token(StringInterpEndToken)
			}
			return l.token(LBraceToken)
		case ',':
			return l.token(CommaToken)
		case '.':
			return l.token(DotToken)
		case '-':
			if l.matchChar('=') {
				return l.token(MinusEqualToken)
			}
			if l.matchChar('>') {
				return l.token(ThinArrowToken)
			}
			return l.token(MinusToken)
		case '+':
			if l.matchChar('=') {
				return l.token(PlusEqualToken)
			}
			return l.token(PlusToken)
		case '^':
			if l.matchChar('=') {
				return l.token(XorEqualToken)
			}
			return l.token(XorToken)
		case '*':
			if l.matchChar('=') {
				return l.token(StarEqualToken)
			}
			if l.matchChar('*') {
				if l.matchChar('=') {
					return l.token(PowerEqualToken)
				}
				return l.token(PowerToken)
			}
			return l.token(StarToken)
		case '=':
			if l.matchChar('=') {
				if l.matchChar('=') {
					return l.token(StrictEqualToken)
				}
				return l.token(EqualToken)
			}
			if l.matchChar('~') {
				return l.token(MatchOperatorToken)
			}
			if l.matchChar('>') {
				return l.token(ThickArrowToken)
			}
			if l.peekChar() == ':' && l.peekNextChar() == '=' {
				l.advanceChar()
				l.advanceChar()
				return l.token(RefEqualToken)
			}
			if l.peekChar() == '!' && l.peekNextChar() == '=' {
				l.advanceChar()
				l.advanceChar()
				return l.token(RefNotEqualToken)
			}
			return l.token(AssignToken)
		case ':':
			if l.matchChar(':') {
				return l.token(ScopeResOperatorToken)
			}
			if l.matchChar('=') {
				return l.token(ColonEqualToken)
			}
			if l.matchChar('>') {
				if l.matchChar('>') {
					return l.token(ReverseInstanceOfToken)
				}
				return l.token(ReverseSubtypeToken)
			}

			return l.token(ColonToken)
		case '~':
			if l.matchChar('=') {
				return l.token(TildeEqualToken)
			}
			if l.matchChar('>') {
				return l.token(WigglyArrowToken)
			}
			return l.token(TildeToken)
		case ';':
			return l.token(SeparatorToken)
		case '>':
			if l.matchChar('=') {
				return l.token(GreaterEqualToken)
			}
			if l.matchChar('>') {
				if l.matchChar('=') {
					return l.token(RBitShiftEqualToken)
				}
				return l.token(RBitShiftToken)
			}
			return l.token(GreaterToken)
		case '<':
			if l.matchChar('=') {
				return l.token(LessEqualToken)
			}
			if l.matchChar(':') {
				return l.token(SubtypeToken)
			}
			if l.matchChar('<') {
				if l.matchChar('=') {
					return l.token(LBitShiftEqualToken)
				}
				if l.matchChar(':') {
					return l.token(InstanceOfToken)
				}
				return l.token(LBitShiftToken)
			}
			return l.token(LessToken)
		case '&':
			if l.matchChar('&') {
				if l.matchChar('=') {
					return l.token(AndAndEqualToken)
				}
				return l.token(AndAndToken)
			}
			if l.matchChar('=') {
				return l.token(AndEqualToken)
			}
			return l.token(AndToken)
		case '|':
			if l.matchChar('|') {
				if l.matchChar('=') {
					return l.token(OrOrEqualToken)
				}
				return l.token(OrOrToken)
			}
			if l.matchChar('>') {
				return l.token(PipeOperatorToken)
			}
			if l.matchChar('=') {
				return l.token(OrEqualToken)
			}
			return l.token(OrToken)
		case '?':
			if l.matchChar('?') {
				if l.matchChar('=') {
					return l.token(NilCoalesceEqualToken)
				}
				return l.token(NilCoalesceToken)
			}
			return l.token(QuestionMarkToken)
		case '!':
			if l.matchChar('=') {
				if l.matchChar('=') {
					return l.token(StrictNotEqualToken)
				}
				return l.token(NotEqualToken)
			}
			return l.token(BangToken)
		case '%':
			if l.matchChar('=') {
				return l.token(PercentEqualToken)
			}
			if l.matchChar('w') {
				return l.token(PercentWToken)
			}
			if l.matchChar('s') {
				return l.token(PercentSToken)
			}
			if l.matchChar('i') {
				return l.token(PercentIToken)
			}
			if l.matchChar('f') {
				return l.token(PercentFToken)
			}
			if l.matchChar('{') {
				return l.token(SetLiteralBegToken)
			}
			if l.matchChar('(') {
				return l.token(TupleLiteralBegToken)
			}
			return l.token(PercentToken)

		case '\n':
			l.foldNewLines()
			return l.token(SeparatorToken)
		case '\r':
			if l.matchChar('\n') {
				l.foldNewLines()
				return l.token(SeparatorToken)
			}
			fallthrough
		case '\t':
			fallthrough
		case ' ':
			l.skipChar()
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
		case '\'':
			return l.rawString()
		case '"':
			l.mode = inStringLiteralMode
			return l.token(StringBegToken)
		case '_':
			return l.privateIdentifier()
		case '@':
			return l.instanceVariable()
		default:
			if isDigit(char) {
				return l.numberLiteral(char)
			} else if unicode.IsLetter(char) {
				return l.identifier(char)
			}
			return l.tokenWithValue(ErrorToken, l.unexpectedCharsError())
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

// Builds a lexing error message with a hint from source code
// based on the current state of the lexer.
func (l *lexer) lexErrorWithHint(message string) string {
	lexValue := l.tokenValue()
	lexValue = lexValue[0:minInt(len(lexValue), maxErrLen)]
	i := l.start
	var trimmedValueToken []byte
	var byt byte

	for {
		if i == len(lexValue) || i == len(l.source) {
			break
		}

		byt = l.source[i]
		if byt == '\n' {
			break
		}

		trimmedValueToken = append(trimmedValueToken, byt)
		i += 1
	}
	lexValue = string(trimmedValueToken)
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

// Creates a new lexing error message which shows unexpected characters.
func (l *lexer) unexpectedCharsError() string {
	return l.lexErrorWithHint("unexpected characters")
}

// Creates a new lexing error message.
func (l *lexer) lexError(message string) string {
	return fmt.Sprintf("%s:%d:%d Lexing error, %s", l.sourceName, l.startLine, l.startColumn, message)
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
