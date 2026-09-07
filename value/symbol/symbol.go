package symbol

import (
	"encoding/binary"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/cespare/xxhash/v2"
)

// Numerical ID of a particular symbol.
type Symbol int

func (s Symbol) String() string {
	name, ok := SymbolTable.GetName(s)
	if !ok {
		panic(fmt.Sprintf("trying to get the name of a nonexistent symbol: %#v", s))
	}
	return name
}

func InspectSymbolContent(name string) string {
	var quotes bool
	var result strings.Builder
	firstLetter := true
	str := name

	for {
		if len(str) == 0 {
			break
		}
		char, bytes := utf8.DecodeRuneInString(str)
		str = str[bytes:]
		switch char {
		case '\\':
			result.WriteString(`\\`)
			quotes = true
		case '\n':
			result.WriteString(`\n`)
			quotes = true
		case '\t':
			result.WriteString(`\t`)
			quotes = true
		case '\r':
			result.WriteString(`\r`)
			quotes = true
		case '\a':
			result.WriteString(`\a`)
			quotes = true
		case '\b':
			result.WriteString(`\b`)
			quotes = true
		case '\v':
			result.WriteString(`\v`)
			quotes = true
		case '\f':
			result.WriteString(`\f`)
			quotes = true
		case '"':
			result.WriteString(`\"`)
			quotes = true
		case '_':
			result.WriteByte('_')
		default:
			if firstLetter && unicode.IsDigit(char) {
				quotes = true
			} else if !quotes && !unicode.IsDigit(char) && !unicode.IsLetter(char) {
				quotes = true
			}

			if unicode.IsGraphic(char) {
				result.WriteRune(char)
			} else if bytes == 1 {
				fmt.Fprintf(&result, `\x%02x`, char)
			} else {
				fmt.Fprintf(&result, `\U%08X`, char)
			}
		}

		firstLetter = false
	}

	if quotes {
		return fmt.Sprintf(`"%s"`, result.String())
	}
	return result.String()
}

func InspectSymbol(name string) string {
	return fmt.Sprintf(`:%s`, InspectSymbolContent(name))
}

func (s Symbol) InspectContent() string {
	return InspectSymbolContent(s.String())
}

func (s Symbol) Inspect() string {
	return InspectSymbol(s.String())
}

// Check whether s is equal to other
func (s Symbol) EqualSymbol(other Symbol) bool {
	return s == other
}

func (s Symbol) Hash() uint64 {
	d := xxhash.New()
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(s))
	d.Write(b)
	return d.Sum64()
}
