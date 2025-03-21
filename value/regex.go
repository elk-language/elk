package value

import (
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

func (r *Regex) Copy() Reference {
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
	var buff strings.Builder
	buff.WriteString(`%/`)
	buff.WriteString(r.Source)
	buff.WriteRune('/')
	buff.WriteString(flag.ToString(r.Flags))
	return buff.String()
}

func (r *Regex) Error() string {
	return r.Inspect()
}

func (r *Regex) InstanceVariables() SymbolMap {
	return nil
}

// Check whether r is equal to other
func (r *Regex) Equal(other Value) Value {
	if !other.IsReference() {
		return False
	}

	switch o := other.AsReference().(type) {
	case *Regex:
		return ToElkBool(r.Flags == o.Flags && r.Source == o.Source)
	default:
		return False
	}
}

// Check whether r is equal to other
func (r *Regex) LaxEqualVal(other Value) Value {
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
func (r *Regex) ConcatVal(other Value) (Value, Value) {
	if !other.IsReference() {
		return Undefined, Ref(NewCoerceError(r.Class(), other.Class()))
	}

	switch o := other.AsReference().(type) {
	case *Regex:
		var buff strings.Builder
		r.WriteSourceTo(&buff)
		o.WriteSourceTo(&buff)
		re, err := CompileRegex(buff.String(), bitfield.BitField8{})
		if err != nil {
			return Undefined, Ref(NewError(RegexCompileErrorClass, err.Error()))
		}
		return Ref(re), Undefined
	default:
		return Undefined, Ref(NewCoerceError(r.Class(), other.Class()))
	}
}

// RepeatVal the content of this Regex n times and return a new Regex.
func (r *Regex) RepeatVal(other Value) (Value, Value) {
	if other.IsReference() {
		return Undefined, Ref(NewCoerceError(r.Class(), other.Class()))
	}

	switch other.ValueFlag() {
	case SMALL_INT_FLAG:
		var buff strings.Builder
		buff.WriteString("(?:")
		buff.WriteString(r.Source)
		buff.WriteRune(')')
		buff.WriteRune('{')
		buff.WriteString(strconv.Itoa(int(other.AsSmallInt())))
		buff.WriteRune('}')
		re, err := CompileRegex(buff.String(), r.Flags)
		if err != nil {
			return Undefined, Ref(NewError(RegexCompileErrorClass, err.Error()))
		}
		return Ref(re), Undefined
	default:
		return Undefined, Ref(NewCoerceError(r.Class(), other.Class()))
	}
}

// Check whether the regex matches the given string
func (r *Regex) Matches(other Value) (Value, Value) {
	if !other.IsReference() {
		return Undefined, Ref(NewCoerceError(r.Class(), other.Class()))
	}
	switch o := other.AsReference().(type) {
	case String:
		return ToElkBool(r.Re.MatchString(string(o))), Undefined
	default:
		return Undefined, Ref(NewCoerceError(r.Class(), other.Class()))
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
	RegexClass = NewClass()
	StdModule.AddConstantString("Regex", Ref(RegexClass))
}
