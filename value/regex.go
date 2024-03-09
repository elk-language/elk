package value

import (
	"fmt"
	"regexp"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/regex/flag"
)

var RegexClass *Class // ::Std::Regex

// Elk's compiled regex
type Regex struct {
	Re     regexp.Regexp
	Source string
	Flags  bitfield.BitField8
}

func NewRegex(re regexp.Regexp, src string, flags bitfield.BitField8) *Regex {
	return &Regex{
		Re:     re,
		Source: src,
		Flags:  flags,
	}
}

func (Regex) Class() *Class {
	return RegexClass
}

func (Regex) DirectClass() *Class {
	return RegexClass
}

func (r Regex) Copy() Value {
	return r
}

func (Regex) SingletonClass() *Class {
	return nil
}

func (r Regex) String() string {
	return r.Source
}

func (r Regex) ToString() String {
	return String(r.String())
}

func (r Regex) Inspect() string {
	return fmt.Sprintf(`%%/%s/%s`, r.Source, flag.ToString(r.Flags))
}

func (r Regex) InstanceVariables() SymbolMap {
	return nil
}

// Check whether r is equal to other
func (r Regex) Equal(other Value) Value {
	switch o := other.(type) {
	case Regex:
		return ToElkBool(r.Flags == o.Flags && r.Source == o.Source)
	default:
		return False
	}
}

// Check whether r is equal to other
func (r Regex) StrictEqual(other Value) Value {
	return r.Equal(other)
}

// Check whether r is equal to other
func (r Regex) LaxEqual(other Value) Value {
	return r.Equal(other)
}

func (r Regex) Hash() UInt64 {
	d := xxhash.New()
	d.WriteString(string(r.Source))
	d.Write([]byte{r.Flags.Byte()})
	return UInt64(d.Sum64())
}

func initRegex() {
	RegexClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Regex", RegexClass)
}
