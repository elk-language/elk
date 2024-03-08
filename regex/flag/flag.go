// Package flag contains the definitions of Elk Regex flags
package flag

import (
	"github.com/elk-language/elk/bitfield"
)

const (
	CaseInsensitiveFlag bitfield.BitFlag8 = 1 << iota // i - case insensitive character matching
	MultilineFlag                                     // m - multi-line mode: ^ and $ match begin/end line in addition to begin/end text
	DotAllFlag                                        // s - let . match \n
	UngreedyFlag                                      // U - ungreedy: swap meaning of x* and x*?, x+ and x+?, etc
	ExtendedFlag                                      // x - ignore all whitespace and allow for comments with #
	ASCIIFlag                                         // a - ASCII mode, Perl char classes like \w, \d, \s only match ASCII characters
)

func IsSupportedByGo(flag bitfield.BitFlag8) bool {
	return flag <= UngreedyFlag
}

func ToChar(flag bitfield.BitFlag8) rune {
	return chars[flag]
}

var Flags = [...]bitfield.BitFlag8{
	CaseInsensitiveFlag,
	MultilineFlag,
	DotAllFlag,
	UngreedyFlag,
	ExtendedFlag,
	ASCIIFlag,
}

var chars = map[bitfield.BitFlag8]rune{
	CaseInsensitiveFlag: 'i',
	MultilineFlag:       'm',
	DotAllFlag:          's',
	UngreedyFlag:        'U',
	ExtendedFlag:        'x',
	ASCIIFlag:           'a',
}
