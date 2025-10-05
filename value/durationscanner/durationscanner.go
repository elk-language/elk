// Package durationscanner implements a tokenizer/lexer
// that analyses Elk duration strings.
package durationscanner

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Scans a duration string and produces
// tokens for each datetime component.
type Durationscanner struct {
	durString string
	cursor    int
	start     int
}

func New(durString string) *Durationscanner {
	return &Durationscanner{
		durString: durString,
	}
}

func (t *Durationscanner) Next() (Token, string) {
	if !t.hasMoreTokens() {
		return END_OF_FILE, ""
	}

	token, value := t.scan()
	t.start = t.cursor
	return token, value
}

func (t *Durationscanner) nextChar() (rune, int) {
	return utf8.DecodeRuneInString(t.durString[t.cursor:])
}

// Gets the next UTF-8 encoded character
// and increments the cursor.
func (t *Durationscanner) advanceChar() (rune, bool) {
	if !t.hasMoreTokens() {
		return 0, false
	}

	char, size := t.nextChar()

	t.cursor += size
	return char, true
}

func (t *Durationscanner) mustAdvanceChar() rune {
	ch, _ := t.advanceChar()
	return ch
}

// Returns true if there is any code left to analyse.
func (t *Durationscanner) hasMoreTokens() bool {
	return t.cursor < len(t.durString)
}

// Gets the next UTF-8 encoded character
// without incrementing the cursor.
func (t *Durationscanner) peekChar() rune {
	if !t.hasMoreTokens() {
		return '\x00'
	}
	char, _ := t.nextChar()
	return char
}

func (t *Durationscanner) value() string {
	return t.durString[t.start:t.cursor]
}

// Checks if the given character matches
// the next UTF-8 encoded character in source code.
// If they match, the cursor gets incremented.
func (t *Durationscanner) matchChar(char rune) bool {
	if !t.hasMoreTokens() {
		return false
	}

	if t.peekChar() == char {
		t.advanceChar()
		return true
	}

	return false
}

// Checks whether the given character is a digit.
func isDigit(char rune) bool {
	return char >= '0' && char <= '9'
}

// Checks if the next character is from the valid set.
func (d *Durationscanner) acceptChars(validChars string) bool {
	if !d.hasMoreTokens() {
		return false
	}

	return strings.ContainsRune(validChars, d.peekChar())
}

func (t *Durationscanner) scan() (Token, string) {
	for {
		char, ok := t.advanceChar()
		if !ok {
			return END_OF_FILE, ""
		}

		if unicode.IsSpace(char) {
			continue
		}
		if isDigit(char) || char == '.' {
			return t.numericComponent(char)
		}

		return ERROR, fmt.Sprintf("unexpected char '%c', expected a digit", char)
	}
}

// Consumes digits from the given set
// and appends them to the given buffer.
// Underscores are ignored.
func (d *Durationscanner) consumeDigits(digitSet string, lexemeBuff *strings.Builder) {
	for {
		if d.peekChar() == '_' {
			d.advanceChar()
		}
		if !d.acceptChars(digitSet) {
			break
		}
		char, _ := d.advanceChar()
		lexemeBuff.WriteRune(char)
	}
}

const (
	decimalLiteralChars = "0123456789"
)

func (d *Durationscanner) numericComponent(startDigit rune) (Token, string) {
	var lexeme strings.Builder

	if startDigit == '.' {
		lexeme.WriteString("0.")
		d.consumeDigits(decimalLiteralChars, &lexeme)
	} else {
		lexeme.WriteRune(startDigit)
		d.consumeDigits(decimalLiteralChars, &lexeme)

		if d.matchChar('.') {
			lexeme.WriteByte('.')
			d.consumeDigits(decimalLiteralChars, &lexeme)
		}
	}

	char, ok := d.advanceChar()
	if !ok {
		return ERROR, "unexpected EOF"
	}

	switch char {
	case 'Y':
		return YEARS, lexeme.String()
	case 'M':
		return MONTHS, lexeme.String()
	case 'D':
		return DAYS, lexeme.String()
	case 'h':
		return HOURS, lexeme.String()
	case 'm':
		if d.matchChar('s') {
			return MILLISECONDS, lexeme.String()
		}
		return MINUTES, lexeme.String()
	case 's':
		return SECONDS, lexeme.String()
	case 'u':
		if d.matchChar('s') {
			return MICROSECONDS, lexeme.String()
		}
		return ERROR, fmt.Sprintf("unexpected char '%c', expected 's'", d.mustAdvanceChar())
	case 'µ':
		if d.matchChar('s') {
			return MICROSECONDS, lexeme.String()
		}
		return ERROR, fmt.Sprintf("unexpected char '%c', expected 's'", d.mustAdvanceChar())
	case 'n':
		if d.matchChar('s') {
			return NANOSECONDS, lexeme.String()
		}
		return ERROR, fmt.Sprintf("unexpected char '%c', expected 's'", d.mustAdvanceChar())
	default:
		return ERROR, fmt.Sprintf("unexpected char '%c', expected a duration component like 'Y', 'h', 's'", char)
	}
}
