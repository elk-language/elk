package value

import (
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/cespare/xxhash/v2"
	"github.com/elk-language/elk/bitfield"
	"github.com/elk-language/elk/regex"
	"github.com/elk-language/elk/regex/flag"
	"github.com/google/go-cmp/cmp"
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

func CompileRegex(src string, flags bitfield.BitField8) (*Regex, error) {
	goSrc, errList := regex.Transpile(src, flags)
	if errList != nil {
		return nil, errList
	}

	goRe, err := regexp.Compile(goSrc)
	if err != nil {
		return nil, err
	}

	return NewRegex(*goRe, src, flags), nil
}

func MustCompileRegex(src string, flags bitfield.BitField8) *Regex {
	re, err := CompileRegex(src, flags)
	if err != nil {
		panic(err)
	}

	return re
}

func (*Regex) Class() *Class {
	return RegexClass
}

func (*Regex) DirectClass() *Class {
	return RegexClass
}

func (r *Regex) Copy() Value {
	return r
}

func (*Regex) SingletonClass() *Class {
	return nil
}

func (r *Regex) String() string {
	return r.Source
}

func (r *Regex) ToString() String {
	return String(r.String())
}

func (r *Regex) ToStringWithFlags() String {
	var buff strings.Builder
	buff.WriteString("(?")
	buff.WriteString(flag.ToStringWithDisabledFlags(r.Flags))
	buff.WriteRune(':')
	buff.WriteString(r.Source)
	buff.WriteRune(')')
	return String(buff.String())
}

func (r *Regex) Inspect() string {
	return fmt.Sprintf(`%%/%s/%s`, r.Source, flag.ToString(r.Flags))
}

func (r *Regex) InstanceVariables() SymbolMap {
	return nil
}

// Check whether r is equal to other
func (r *Regex) Equal(other Value) Value {
	switch o := other.(type) {
	case *Regex:
		return ToElkBool(r.Flags == o.Flags && r.Source == o.Source)
	default:
		return False
	}
}

// Check whether r is equal to other
func (r *Regex) StrictEqual(other Value) Value {
	return r.Equal(other)
}

// Check whether r is equal to other
func (r *Regex) LaxEqual(other Value) Value {
	return r.Equal(other)
}

func (r *Regex) Hash() UInt64 {
	d := xxhash.New()
	d.WriteString(string(r.Source))
	d.Write([]byte{r.Flags.Byte()})
	return UInt64(d.Sum64())
}

func (r *Regex) WriteSourceTo(w io.StringWriter) {
	w.WriteString("(?")
	w.WriteString(flag.ToString(r.Flags))
	w.WriteString(":")
	w.WriteString(r.Source)
	w.WriteString(")")
}

// Create a new regex concatenating r with other
func (r *Regex) Concat(other Value) (Value, *Error) {
	switch o := other.(type) {
	case *Regex:
		var buff strings.Builder
		r.WriteSourceTo(&buff)
		o.WriteSourceTo(&buff)
		re, err := CompileRegex(buff.String(), bitfield.BitField8{})
		if err != nil {
			return nil, NewError(RegexCompileErrorClass, err.Error())
		}
		return re, nil
	default:
		return nil, NewCoerceError(r.Class(), other.Class())
	}
}

// Repeat the content of this Regex n times and return a new Regex.
func (r *Regex) Repeat(other Value) (Value, *Error) {
	switch o := other.(type) {
	case SmallInt:
		var buff strings.Builder
		buff.WriteString("(?:")
		buff.WriteString(r.Source)
		buff.WriteRune(')')
		buff.WriteRune('{')
		buff.WriteString(strconv.Itoa(int(o)))
		buff.WriteRune('}')
		re, err := CompileRegex(buff.String(), r.Flags)
		if err != nil {
			return nil, NewError(RegexCompileErrorClass, err.Error())
		}
		return re, nil
	default:
		return nil, NewCoerceError(r.Class(), other.Class())
	}
}

// Check whether the regex matches the given string
func (r *Regex) Match(other Value) (Value, *Error) {
	switch o := other.(type) {
	case String:
		return ToElkBool(r.Re.MatchString(string(o))), nil
	default:
		return nil, NewCoerceError(StringClass, other.Class())
	}
}

func NewRegexComparer(opts *cmp.Options) cmp.Option {
	return cmp.Comparer(func(x, y *Regex) bool {
		if x == y {
			return true
		}

		return x.Flags == y.Flags && x.Source == y.Source
	})
}

func initRegex() {
	RegexClass = NewClassWithOptions(
		ClassWithSealed(),
		ClassWithNoInstanceVariables(),
	)
	StdModule.AddConstantString("Regex", RegexClass)
}
