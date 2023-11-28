// Package timescanner implements a tokenizer/lexer
// that analyses Elk time format strings.
package timescanner

import (
	"strings"
	"unicode/utf8"
)

// Scans a time format string and produces
// tokens from various formatting options.
type Timescanner struct {
	fmtString string
	cursor    int
	start     int
}

func New(fmtString string) *Timescanner {
	return &Timescanner{
		fmtString: fmtString,
	}
}

func (t *Timescanner) Next() (Token, string) {
	if !t.hasMoreTokens() {
		return END_OF_FILE, ""
	}

	token, value := t.scan()
	t.start = t.cursor
	return token, value
}

func (t *Timescanner) nextChar() (rune, int) {
	return utf8.DecodeRuneInString(t.fmtString[t.cursor:])
}

// Gets the next UTF-8 encoded character
// and increments the cursor.
func (t *Timescanner) advanceChar() (rune, bool) {
	if !t.hasMoreTokens() {
		return 0, false
	}

	char, size := t.nextChar()

	t.cursor += size
	return char, true
}

// Returns true if there is any code left to analyse.
func (t *Timescanner) hasMoreTokens() bool {
	return t.cursor < len(t.fmtString)
}

// Gets the next UTF-8 encoded character
// without incrementing the cursor.
func (t *Timescanner) peekChar() rune {
	if !t.hasMoreTokens() {
		return '\x00'
	}
	char, _ := t.nextChar()
	return char
}

func (t *Timescanner) value() string {
	return t.fmtString[t.start:t.cursor]
}

// Checks if the given character matches
// the next UTF-8 encoded character in source code.
// If they match, the cursor gets incremented.
func (t *Timescanner) matchChar(char rune) bool {
	if !t.hasMoreTokens() {
		return false
	}

	if t.peekChar() == char {
		t.advanceChar()
		return true
	}

	return false
}

func (t *Timescanner) scan() (Token, string) {
	for {
		char, ok := t.advanceChar()
		if !ok {
			return END_OF_FILE, ""
		}

		switch char {
		case '%':
			if t.matchChar('%') {
				return PERCENT, ""
			}
			if t.matchChar('n') {
				return NEWLINE, ""
			}
			if t.matchChar('t') {
				return TAB, ""
			}
			if t.matchChar(':') {
				if t.matchChar('z') {
					return TIMEZONE_OFFSET_COLON, ""
				}
				t.advanceChar()
				return INVALID_FORMAT_DIRECTIVE, t.value()
			}
			if t.matchChar('-') {
				if t.matchChar('Y') {
					return FULL_YEAR, ""
				}
				if t.matchChar('C') {
					return CENTURY, ""
				}
				if t.matchChar('y') {
					return YEAR_LAST_TWO, ""
				}
				if t.matchChar('m') {
					return MONTH, ""
				}
				if t.matchChar('d') {
					return DAY_OF_MONTH, ""
				}
				if t.matchChar('j') {
					return DAY_OF_YEAR, ""
				}
				if t.matchChar('H') {
					return HOUR_OF_DAY, ""
				}
				if t.matchChar('I') {
					return HOUR_OF_DAY12, ""
				}
				if t.matchChar('M') {
					return MINUTE_OF_HOUR, ""
				}
				if t.matchChar('S') {
					return SECOND_OF_MINUTE, ""
				}
				if t.matchChar('L') {
					return MILLISECOND_OF_SECOND, ""
				}
				if t.matchChar('G') {
					return FULL_YEAR_WEEK_BASED, ""
				}
				if t.matchChar('g') {
					return YEAR_LAST_TWO_WEEK_BASED, ""
				}
				if t.matchChar('V') {
					return WEEK_OF_WEEK_BASED_YEAR, ""
				}
				if t.matchChar('U') {
					return WEEK_OF_YEAR_ALT, ""
				}
				if t.matchChar('W') {
					return WEEK_OF_YEAR, ""
				}
				if t.matchChar('N') {
					return NANOSECOND_OF_SECOND, ""
				}
				if t.matchChar('3') && t.matchChar('N') {
					return MILLISECOND_OF_SECOND, ""
				}
				if t.matchChar('6') && t.matchChar('N') {
					return MICROSECOND_OF_SECOND, ""
				}
				if t.matchChar('9') && t.matchChar('N') {
					return NANOSECOND_OF_SECOND, ""
				}
				if t.matchChar('1') {
					if t.matchChar('2') && t.matchChar('N') {
						return PICOSECOND_OF_SECOND, ""
					}
					if t.matchChar('5') && t.matchChar('N') {
						return FEMTOSECOND_OF_SECOND, ""
					}
					if t.matchChar('8') && t.matchChar('N') {
						return ATTOSECOND_OF_SECOND, ""
					}
					t.advanceChar()
					return INVALID_FORMAT_DIRECTIVE, t.value()
				}
				if t.matchChar('2') {
					if t.matchChar('1') && t.matchChar('N') {
						return ZEPTOSECOND_OF_SECOND, ""
					}
					if t.matchChar('4') && t.matchChar('N') {
						return YOCTOSECOND_OF_SECOND, ""
					}
					t.advanceChar()
					return INVALID_FORMAT_DIRECTIVE, t.value()
				}
				t.advanceChar()
				return INVALID_FORMAT_DIRECTIVE, t.value()
			}
			if t.matchChar('_') {
				if t.matchChar('Y') {
					return FULL_YEAR_SPACE_PADDED, ""
				}
				if t.matchChar('C') {
					return CENTURY_SPACE_PADDED, ""
				}
				if t.matchChar('y') {
					return YEAR_LAST_TWO_SPACE_PADDED, ""
				}
				if t.matchChar('m') {
					return MONTH_SPACE_PADDED, ""
				}
				if t.matchChar('d') {
					return DAY_OF_MONTH_SPACE_PADDED, ""
				}
				if t.matchChar('j') {
					return DAY_OF_YEAR_SPACE_PADDED, ""
				}
				if t.matchChar('H') {
					return HOUR_OF_DAY_SPACE_PADDED, ""
				}
				if t.matchChar('I') {
					return HOUR_OF_DAY12_SPACE_PADDED, ""
				}
				if t.matchChar('M') {
					return MINUTE_OF_HOUR_SPACE_PADDED, ""
				}
				if t.matchChar('S') {
					return SECOND_OF_MINUTE_SPACE_PADDED, ""
				}
				if t.matchChar('L') {
					return MILLISECOND_OF_SECOND_SPACE_PADDED, ""
				}
				if t.matchChar('G') {
					return FULL_YEAR_WEEK_BASED_SPACE_PADDED, ""
				}
				if t.matchChar('g') {
					return YEAR_LAST_TWO_WEEK_BASED_SPACE_PADDED, ""
				}
				if t.matchChar('V') {
					return WEEK_OF_WEEK_BASED_YEAR_SPACE_PADDED, ""
				}
				if t.matchChar('U') {
					return WEEK_OF_YEAR_ALT_SPACE_PADDED, ""
				}
				if t.matchChar('W') {
					return WEEK_OF_YEAR_SPACE_PADDED, ""
				}
				if t.matchChar('N') {
					return NANOSECOND_OF_SECOND_SPACE_PADDED, ""
				}
				if t.matchChar('3') && t.matchChar('N') {
					return MILLISECOND_OF_SECOND_SPACE_PADDED, ""
				}
				if t.matchChar('6') && t.matchChar('N') {
					return MICROSECOND_OF_SECOND_SPACE_PADDED, ""
				}
				if t.matchChar('9') && t.matchChar('N') {
					return NANOSECOND_OF_SECOND_SPACE_PADDED, ""
				}
				if t.matchChar('1') {
					if t.matchChar('2') && t.matchChar('N') {
						return PICOSECOND_OF_SECOND_SPACE_PADDED, ""
					}
					if t.matchChar('5') && t.matchChar('N') {
						return FEMTOSECOND_OF_SECOND_SPACE_PADDED, ""
					}
					if t.matchChar('8') && t.matchChar('N') {
						return ATTOSECOND_OF_SECOND_SPACE_PADDED, ""
					}
					t.advanceChar()
					return INVALID_FORMAT_DIRECTIVE, t.value()
				}
				if t.matchChar('2') {
					if t.matchChar('1') && t.matchChar('N') {
						return ZEPTOSECOND_OF_SECOND_SPACE_PADDED, ""
					}
					if t.matchChar('4') && t.matchChar('N') {
						return YOCTOSECOND_OF_SECOND_SPACE_PADDED, ""
					}
					t.advanceChar()
					return INVALID_FORMAT_DIRECTIVE, t.value()
				}
				t.advanceChar()
				return INVALID_FORMAT_DIRECTIVE, t.value()
			}
			if t.matchChar('^') {
				if t.matchChar('B') {
					return MONTH_FULL_NAME_UPPERCASE, ""
				}
				if t.matchChar('b') {
					return MONTH_ABBREVIATED_NAME_UPPERCASE, ""
				}
				if t.matchChar('P') {
					return MERIDIEM_INDICATOR_UPPERCASE, ""
				}
				if t.matchChar('A') {
					return DAY_OF_WEEK_FULL_NAME_UPPERCASE, ""
				}
				if t.matchChar('a') {
					return DAY_OF_WEEK_ABBREVIATED_NAME_UPPERCASE, ""
				}
				if t.matchChar('c') {
					return DATE_AND_TIME_UPPERCASE, ""
				}
				if t.matchChar('+') {
					return DATE1_FORMAT_UPPERCASE, ""
				}
				t.advanceChar()
				return INVALID_FORMAT_DIRECTIVE, t.value()
			}
			if t.matchChar('Y') {
				return FULL_YEAR_ZERO_PADDED, ""
			}
			if t.matchChar('C') {
				return CENTURY_ZERO_PADDED, ""
			}
			if t.matchChar('y') {
				return YEAR_LAST_TWO_ZERO_PADDED, ""
			}
			if t.matchChar('m') {
				return MONTH_ZERO_PADDED, ""
			}
			if t.matchChar('B') {
				return MONTH_FULL_NAME, ""
			}
			if t.matchChar('b') || t.matchChar('h') {
				return MONTH_ABBREVIATED_NAME, ""
			}
			if t.matchChar('d') {
				return DAY_OF_MONTH_ZERO_PADDED, ""
			}
			if t.matchChar('e') {
				return DAY_OF_MONTH_SPACE_PADDED, ""
			}
			if t.matchChar('j') {
				return DAY_OF_YEAR_ZERO_PADDED, ""
			}
			if t.matchChar('H') {
				return HOUR_OF_DAY_ZERO_PADDED, ""
			}
			if t.matchChar('k') {
				return HOUR_OF_DAY_SPACE_PADDED, ""
			}
			if t.matchChar('I') {
				return HOUR_OF_DAY12_ZERO_PADDED, ""
			}
			if t.matchChar('l') {
				return HOUR_OF_DAY12_SPACE_PADDED, ""
			}
			if t.matchChar('p') {
				return MERIDIEM_INDICATOR_UPPERCASE, ""
			}
			if t.matchChar('P') {
				return MERIDIEM_INDICATOR_LOWERCASE, ""
			}
			if t.matchChar('M') {
				return MINUTE_OF_HOUR_ZERO_PADDED, ""
			}
			if t.matchChar('S') {
				return SECOND_OF_MINUTE_ZERO_PADDED, ""
			}
			if t.matchChar('s') {
				return UNIX_SECONDS, ""
			}
			if t.matchChar('Q') {
				return UNIX_MILLISECONDS, ""
			}
			if t.matchChar('L') {
				return MILLISECOND_OF_SECOND_ZERO_PADDED, ""
			}
			if t.matchChar('Z') {
				return TIMEZONE_NAME, ""
			}
			if t.matchChar('z') {
				return TIMEZONE_OFFSET, ""
			}
			if t.matchChar('A') {
				return DAY_OF_WEEK_FULL_NAME, ""
			}
			if t.matchChar('a') {
				return DAY_OF_WEEK_ABBREVIATED_NAME, ""
			}
			if t.matchChar('u') {
				return DAY_OF_WEEK_NUMBER, ""
			}
			if t.matchChar('w') {
				return DAY_OF_WEEK_NUMBER_ALT, ""
			}
			if t.matchChar('G') {
				return FULL_YEAR_WEEK_BASED_ZERO_PADDED, ""
			}
			if t.matchChar('g') {
				return YEAR_LAST_TWO_WEEK_BASED_ZERO_PADDED, ""
			}
			if t.matchChar('V') {
				return WEEK_OF_WEEK_BASED_YEAR_ZERO_PADDED, ""
			}
			if t.matchChar('U') {
				return WEEK_OF_YEAR_ALT_ZERO_PADDED, ""
			}
			if t.matchChar('W') {
				return WEEK_OF_YEAR_ZERO_PADDED, ""
			}
			if t.matchChar('c') {
				return DATE_AND_TIME, ""
			}
			if t.matchChar('D') {
				return DATE, ""
			}
			if t.matchChar('F') {
				return ISO8601_DATE, ""
			}
			if t.matchChar('F') {
				return ISO8601_DATE, ""
			}
			if t.matchChar('r') {
				return TIME12, ""
			}
			if t.matchChar('R') {
				return TIME24, ""
			}
			if t.matchChar('T') {
				return TIME24_SECONDS, ""
			}
			if t.matchChar('+') {
				return DATE1_FORMAT, ""
			}
			if t.matchChar('N') {
				return NANOSECOND_OF_SECOND_ZERO_PADDED, ""
			}
			if t.matchChar('3') {
				if t.matchChar('N') {
					return MILLISECOND_OF_SECOND_ZERO_PADDED, ""
				}
				if t.matchChar('s') {
					return UNIX_MILLISECONDS, ""
				}
				t.advanceChar()
				return INVALID_FORMAT_DIRECTIVE, t.value()
			}
			if t.matchChar('6') {
				if t.matchChar('N') {
					return MICROSECOND_OF_SECOND_ZERO_PADDED, ""
				}
				if t.matchChar('s') {
					return UNIX_MICROSECONDS, ""
				}
				t.advanceChar()
				return INVALID_FORMAT_DIRECTIVE, t.value()
			}
			if t.matchChar('9') {
				if t.matchChar('N') {
					return NANOSECOND_OF_SECOND_ZERO_PADDED, ""
				}
				if t.matchChar('s') {
					return UNIX_NANOSECONDS, ""
				}
				t.advanceChar()
				return INVALID_FORMAT_DIRECTIVE, t.value()
			}
			if t.matchChar('1') {
				if t.matchChar('2') {
					if t.matchChar('N') {
						return PICOSECOND_OF_SECOND_ZERO_PADDED, ""
					}
					if t.matchChar('s') {
						return UNIX_PICOSECONDS, ""
					}
					t.advanceChar()
					return INVALID_FORMAT_DIRECTIVE, t.value()
				}
				if t.matchChar('5') {
					if t.matchChar('N') {
						return FEMTOSECOND_OF_SECOND_ZERO_PADDED, ""
					}
					if t.matchChar('s') {
						return UNIX_FEMTOSECONDS, ""
					}
					t.advanceChar()
					return INVALID_FORMAT_DIRECTIVE, t.value()
				}
				if t.matchChar('8') {
					if t.matchChar('N') {
						return ATTOSECOND_OF_SECOND_ZERO_PADDED, ""
					}
					if t.matchChar('s') {
						return UNIX_ATTOSECONDS, ""
					}
					t.advanceChar()
					return INVALID_FORMAT_DIRECTIVE, t.value()
				}
				t.advanceChar()
				return INVALID_FORMAT_DIRECTIVE, t.value()
			}
			if t.matchChar('2') {
				if t.matchChar('1') {
					if t.matchChar('N') {
						return ZEPTOSECOND_OF_SECOND_ZERO_PADDED, ""
					}
					if t.matchChar('s') {
						return UNIX_ZEPTOSECONDS, ""
					}
					t.advanceChar()
					return INVALID_FORMAT_DIRECTIVE, t.value()
				}
				if t.matchChar('4') {
					if t.matchChar('N') {
						return YOCTOSECOND_OF_SECOND_ZERO_PADDED, ""
					}
					if t.matchChar('s') {
						return UNIX_YOCTOSECONDS, ""
					}
					t.advanceChar()
					return INVALID_FORMAT_DIRECTIVE, t.value()
				}

				t.advanceChar()
				return INVALID_FORMAT_DIRECTIVE, t.value()
			}

			t.advanceChar()
			return INVALID_FORMAT_DIRECTIVE, t.value()
		default:
			var buffer strings.Builder
			buffer.WriteRune(char)
			for t.peekChar() != '%' {
				ch, ok := t.advanceChar()
				if !ok {
					return TEXT, buffer.String()
				}
				buffer.WriteRune(ch)
			}

			return TEXT, buffer.String()
		}
	}
}
