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

func (b *StringBuilder) Write(s) (int, error) {

}

func (b *StringBuilder) WriteChar(c Char) {

}

func (b *StringBuilder) WriteByte(b) {

}

func (b *StringBuilder) WriteString(s String) {

}

func initStringBuilder() {
	StringBuilderClass = NewClassWithOptions(ClassWithConstructor(stringBuilderConstructor))
	StdModule.AddConstantString("StringBuilder", Ref(StringBuilderClass))
	RegisterNativeClass("Std::StringBuilder", "value.StringBuilderClass")
}
