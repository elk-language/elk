// Package strings implements the rendering of Elk string and char
// literals back into Elk source code.
package strings

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Render a string as an Elk string literal, escaping it and wrapping it in double quotes.
func InspectString(s string) string {
	var buffer strings.Builder

	buffer.WriteString(`"`)
	leftStr := s
	for {
		char, size := utf8.DecodeRuneInString(leftStr)
		if size == 0 {
			// reached the end of the string
			break
		}
		if char == utf8.RuneError && size == 1 {
			// invalid UTF-8 character
			char = rune(leftStr[0])
		}
		switch char {
		case '\\':
			buffer.WriteString(`\\`)
		case '\n':
			buffer.WriteString(`\n`)
		case '\t':
			buffer.WriteString(`\t`)
		case '"':
			buffer.WriteString(`\"`)
		case '\r':
			buffer.WriteString(`\r`)
		case '\a':
			buffer.WriteString(`\a`)
		case '\b':
			buffer.WriteString(`\b`)
		case '\v':
			buffer.WriteString(`\v`)
		case '\f':
			buffer.WriteString(`\f`)
		case '$':
			buffer.WriteString(`\$`)
		case '#':
			buffer.WriteString(`\#`)
		default:
			if unicode.IsGraphic(char) {
				buffer.WriteRune(char)
			} else if char>>8 == 0 {
				fmt.Fprintf(&buffer, `\x%02x`, char)
			} else if char>>16 == 0 {
				fmt.Fprintf(&buffer, `\u%04x`, char)
			} else {
				fmt.Fprintf(&buffer, `\U%08X`, char)
			}
		}
		leftStr = leftStr[size:]
	}

	buffer.WriteString(`"`)
	return buffer.String()
}

// Render a rune as an Elk char literal, escaping it and wrapping it in backticks.
func InspectChar(c rune) string {
	var buff strings.Builder
	buff.WriteRune('`')
	switch c {
	case '\\':
		buff.WriteString(`\\`)
	case '\n':
		buff.WriteString(`\n`)
	case '\t':
		buff.WriteString(`\t`)
	case '`':
		buff.WriteString("\\`")
	case '\r':
		buff.WriteString(`\r`)
	case '\a':
		buff.WriteString(`\a`)
	case '\b':
		buff.WriteString(`\b`)
	case '\v':
		buff.WriteString(`\v`)
	case '\f':
		buff.WriteString(`\f`)
	default:
		if unicode.IsGraphic(c) {
			buff.WriteRune(c)
		} else if c>>8 == 0 {
			fmt.Fprintf(&buff, `\x%02x`, c)
		} else if c>>16 == 0 {
			fmt.Fprintf(&buff, `\u%04x`, c)
		} else {
			fmt.Fprintf(&buff, `\U%08X`, c)
		}
	}

	buff.WriteRune('`')
	return buff.String()
}
