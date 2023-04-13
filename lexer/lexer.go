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
)

// Lexing mode which changes how characters are handled by the lexer.
type mode uint8

const (
	normalMode                mode = iota // Initial mode
	inWordArrayLiteralMode                // Triggered after entering the initial token `%w[` of a word array literal
	inSymbolArrayLiteralMode              // Triggered after entering the initial token `%s[` of a symbol array literal
	inWordSetLiteralMode                  // Triggered after entering the initial token `%w{` of a word set literal
	inSymbolSetLiteralMode                // Triggered after entering the initial token `%s{` of a symbol set literal
	inWordTupleLiteralMode                // Triggered after entering the initial token `%w(` of a word tuple literal
	inSymbolTupleLiteralMode              // Triggered after entering the initial token `%s(` of a symbol tuple literal
	inStringLiteralMode                   // Triggered after consuming the initial token `"` of a string literal
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
		l.advanceChar()
		return true
	}

	return false
}

// Same as [matchChars] but returns the consumed char.
func (l *lexer) matchCharsRune(validChars string) (bool, rune) {
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
func (l *lexer) matchChars(validChars string) bool {
	if !l.hasMoreTokens() {
		return false
	}

	if strings.ContainsRune(validChars, l.peekChar()) {
		l.advanceChar()
		return true
	}

	return false
}

// Checks if the next character is from the valid set.
func (l *lexer) acceptChars(validChars string) bool {
	if !l.hasMoreTokens() {
		return false
	}

	return strings.ContainsRune(validChars, l.peekChar())
}

// Checks if the second next character is from the valid set.
func (l *lexer) acceptCharsNext(validChars string) bool {
	if !l.hasMoreTokens() {
		return false
	}

	return strings.ContainsRune(validChars, l.peekNextChar())
}

// Consumes a series of characters from the given set.
func (l *lexer) matchCharsRun(validChars string) bool {
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

// Skips the current UTF-8 encoded character.
func (l *lexer) skipChar() {
	_, size := l.nextChar()

	l.start += size
	l.startColumn += 1
}

// Skips the current byte.
func (l *lexer) skipByte() {
	l.start += 1
	l.startColumn += 1
}

// Skips the current accumulated token.
func (l *lexer) skipToken() {
	l.start = l.cursor
	l.startColumn = l.column
	l.startLine = l.line
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
func (l *lexer) rawString() *Token {
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

	return l.tokenWithValue(RawStringToken, result.String())
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
func (l *lexer) consumeDigits(digitSet string, lexemeBuff *strings.Builder) {
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
func (l *lexer) numberLiteral(startDigit rune) *Token {
	tokenType := DecIntToken
	digits := decimalLiteralChars
	var lexeme strings.Builder

	if startDigit == '0' {
		if l.matchChars("xX") {
			// hexadecimal (base 16)
			digits = hexLiteralChars
			tokenType = HexIntToken
		} else if l.matchChars("dD") {
			// duodecimal (base 12)
			digits = duodecimalLiteralChars
			tokenType = DuoIntToken
		} else if l.matchChars("oO") {
			// octal (base 8)
			digits = octalLiteralChars
			tokenType = OctIntToken
		} else if l.matchChars("qQ") {
			// quaternary (base 4)
			digits = quaternaryLiteralChars
			tokenType = QuatIntToken
		} else if l.matchChars("bB") {
			// binary (base 2)
			digits = binaryLiteralChars
			tokenType = BinIntToken
		}
	}

	if tokenType != DecIntToken {
		l.consumeDigits(digits, &lexeme)
		return l.tokenWithValue(tokenType, lexeme.String())
	}
	lexeme.WriteRune(startDigit)
	l.consumeDigits(digits, &lexeme)

	if l.matchChar('.') {
		lexeme.WriteByte('.')
		l.consumeDigits(digits, &lexeme)
		tokenType = FloatToken
	}
	if l.matchChars("eE") {
		lexeme.WriteByte('e')
		if ok, ch := l.matchCharsRune("+-"); ok {
			lexeme.WriteRune(ch)
		}

		l.consumeDigits(digits, &lexeme)
		tokenType = FloatToken
	}

	return l.tokenWithValue(tokenType, lexeme.String())
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
		if l.matchChars("?!") {
			return l.tokenWithConsumedValue(IdentifierToken)
		}
		if lexType := keywords[l.tokenValue()]; lexType.IsKeyword() {
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
		l.matchChars("?!")
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
	case inWordArrayLiteralMode:
		return l.scanWordCollectionLiteral(']', WordArrayEndToken)
	case inSymbolArrayLiteralMode:
		return l.scanWordCollectionLiteral(']', SymbolArrayEndToken)
	case inWordSetLiteralMode:
		return l.scanWordCollectionLiteral('}', WordSetEndToken)
	case inSymbolSetLiteralMode:
		return l.scanWordCollectionLiteral('}', SymbolSetEndToken)
	case inWordTupleLiteralMode:
		return l.scanWordCollectionLiteral(')', WordTupleEndToken)
	case inSymbolTupleLiteralMode:
		return l.scanWordCollectionLiteral(')', SymbolTupleEndToken)
	default:
		return l.lexError("unsupported lexing mode")
	}
}

const (
	unterminatedWordCollectionError = "unterminated %s literal, missing `%c`"
)

// Scans the content of word collection literals be it `%w[`, `%s[`, `%w{`, `%s{`, `%w(`, `%s(`
func (l *lexer) scanWordCollectionLiteral(terminatorChar rune, terminatorToken TokenType) *Token {
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
			return l.lexError(fmt.Sprintf(unterminatedWordCollectionError, "word array", terminatorChar))
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
	return l.tokenWithValue(RawStringToken, result.String())
}

const (
	unterminatedStringError = "unterminated string literal, missing `\"`"
	invalidHexError         = "invalid hex escape"
)

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
		case 'x':
			if !l.acceptChars(hexLiteralChars) || !l.acceptCharsNext(hexLiteralChars) {
				return l.lexError(invalidHexError)
			}
			l.advanceChar()
			l.advanceChar()
			value, err := strconv.ParseUint(string(l.source[l.cursor-2:l.cursor]), 16, 8)
			if err != nil {
				return l.lexError(invalidHexError)
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
			return newEOF()
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
			return l.token(RBraceToken)
		case ',':
			return l.token(CommaToken)
		case '.':
			if l.matchChar('.') {
				if l.matchChar('.') {
					return l.token(ExclusiveRangeOpToken)
				}
				return l.token(RangeOpToken)
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
				return l.tokenWithValue(FloatToken, lexeme.String())
			}
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
				return l.token(MatchOpToken)
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
				return l.token(ScopeResOpToken)
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
			if ch := l.peekChar(); ch == '#' || unicode.IsSpace(ch) {
				return l.token(ColonToken)
			}

			return l.token(SymbolBegToken)
		case '~':
			if l.matchChar('=') {
				return l.token(TildeEqualToken)
			}
			if l.matchChar('>') {
				return l.token(WigglyArrowToken)
			}
			return l.token(TildeToken)
		case ';':
			return l.token(SemicolonToken)
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
				return l.token(PipeOpToken)
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
				if l.matchChar('[') {
					l.mode = inWordArrayLiteralMode
					return l.token(WordArrayBegToken)
				}
				if l.matchChar('{') {
					l.mode = inWordSetLiteralMode
					return l.token(WordSetBegToken)
				}
				if l.matchChar('(') {
					l.mode = inWordTupleLiteralMode
					return l.token(WordTupleBegToken)
				}

				return l.lexError("invalid word collection literal delimiters `%%w`")
			}
			if l.matchChar('s') {
				if l.matchChar('[') {
					l.mode = inSymbolArrayLiteralMode
					return l.token(SymbolArrayBegToken)
				}
				if l.matchChar('{') {
					l.mode = inSymbolSetLiteralMode
					return l.token(SymbolSetBegToken)
				}
				if l.matchChar('(') {
					l.mode = inSymbolTupleLiteralMode
					return l.token(SymbolTupleBegToken)
				}

				return l.lexError("invalid symbol collection literal `%s`")
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
			return l.token(EndLineToken)
		case '\r':
			if l.matchChar('\n') {
				l.foldNewLines()
				return l.token(EndLineToken)
			}
			fallthrough
		case '\t':
			fallthrough
		case ' ':
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
			return l.lexError(fmt.Sprintf("unexpected character `%c`", char))
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

// Creates a new lexing error token.
func (l *lexer) lexError(message string) *Token {
	return l.tokenWithValue(ErrorToken, message)
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
