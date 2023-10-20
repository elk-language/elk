package value

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Numerical ID of a particular symbol.
type Symbol int

var SymbolClass *Class // ::Std::Symbol

func (s Symbol) Class() *Class {
	return StringClass
}

func (s Symbol) IsFrozen() bool {
	return true
}

func (s Symbol) SetFrozen() {}

func (s Symbol) Name() string {
	name, ok := SymbolTable.GetName(s)
	if !ok {
		panic(fmt.Sprintf("trying to get the name of a nonexistent symbol: %#v", s))
	}
	return name
}

func (s Symbol) InspectContent() string {
	var quotes bool
	var result strings.Builder
	firstLetter := true
	str := s.Name()

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
				result.WriteString(fmt.Sprintf(`\x%02x`, char))
			} else {
				result.WriteString(fmt.Sprintf(`\U%08X`, char))
			}
		}

		firstLetter = false
	}

	if quotes {
		return fmt.Sprintf(`"%s"`, result.String())
	}
	return result.String()
}

func (s Symbol) Inspect() string {
	return fmt.Sprintf(`:%s`, s.InspectContent())
}

func (s Symbol) InstanceVariables() SimpleSymbolMap {
	return nil
}

func initSymbol() {
	SymbolClass = NewClass(
		ClassWithImmutable(),
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstant("Symbol", SymbolClass)
}
