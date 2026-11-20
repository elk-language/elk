package value

import (
	"fmt"
	"strings"
)

type StringBuilder struct {
	buf strings.Builder
}

var _ Reference = &StringBuilder{}

var StringBuilderClass *Class // Std::StringBuilder

func stringBuilderConstructor(class *Class) Value {
	return (&StringBuilder{
		buf: strings.Builder{},
	}).ToValue()
}

func (b *StringBuilder) Class() *Class {
	return StringBuilderClass
}

func (b *StringBuilder) DirectClass() *Class {
	return StringBuilderClass
}

func (b *StringBuilder) SingletonClass() *Class {
	return nil
}

func (b *StringBuilder) Copy() *StringBuilder {
	newBuf := strings.Builder{}
	newBuf.WriteString(b.String())
	return &StringBuilder{
		buf: newBuf,
	}
}

func (b *StringBuilder) CopyReference() Reference {
	return b.Copy()
}

func (b *StringBuilder) String() string {
	return b.buf.String()
}

func (b *StringBuilder) ToString() String {
	return String(b.String())
}

func (b *StringBuilder) Inspect() string {
	return fmt.Sprintf(`Std::StringBuilder(%q)`, b.String())
}

func (b *StringBuilder) Error() string {
	return b.Inspect()
}

func (b *StringBuilder) InstanceVariables() *InstanceVariables {
	return nil
}

func (b *StringBuilder) ToValue() Value {
	return Ref(b)
}

func (b *StringBuilder) Write(other Value) SmallInt {
	if other.IsReference() {
		switch o := other.AsReference().(type) {
		case String:
			return b.WriteString(o)
		}
	}

	switch other.ValueFlag() {
	case CHAR_FLAG:
		return b.WriteChar(other.AsChar())
	case UINT8_FLAG:
		b.WriteByte(other.AsUInt8())
		return 1
	default:
		return 0
	}
}

func (b *StringBuilder) WriteChar(c Char) SmallInt {
	result, _ := b.buf.WriteRune(rune(c))
	return SmallInt(result)
}

func (b *StringBuilder) WriteByte(i UInt8) {
	_ = b.buf.WriteByte(byte(i))
}

func (b *StringBuilder) WriteString(s String) SmallInt {
	result, _ := b.buf.WriteString(s.String())
	return SmallInt(result)
}

func (b *StringBuilder) Capacity() int {
	return b.buf.Cap()
}

func (b *StringBuilder) Grow(n int) {
	b.buf.Grow(n)
}

func (b *StringBuilder) Length() int {
	return b.buf.Len()
}

func (b *StringBuilder) Reset() {
	b.buf.Reset()
}

func initStringBuilder() {
	StringBuilderClass = NewClassWithOptions(ClassWithConstructor(stringBuilderConstructor))
	StdModule.AddConstantString("StringBuilder", Ref(StringBuilderClass))
	RegisterNativeClass("Std::StringBuilder", "value.StringBuilderClass")
}
