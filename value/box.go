package value

import (
	"fmt"
	"strings"

	"github.com/elk-language/elk/indent"
)

var BoxClass *Class // ::Std::Box

func initBox() {
	BoxClass = NewClassWithOptions(ClassWithConstructor(BoxConstructor))
	StdModule.AddConstantString("Box", Ref(BoxClass))
}

func BoxConstructor(class *Class) Value {
	return Ref(NewBox(Undefined))
}

// Box wraps another value, it's a pointer to another `Value`.
type Box Value

func NewBox(v Value) *Box {
	b := Box(v)
	return &b
}

// Retrieve the value stored in the box
func (b *Box) Get() Value {
	return Value(*b)
}

// Set the value in the box
func (b *Box) Set(v Value) {
	*b = Box(v)
}

func (*Box) Class() *Class {
	return BoxClass
}

func (*Box) DirectClass() *Class {
	return BoxClass
}

func (*Box) SingletonClass() *Class {
	return nil
}

func (b *Box) Copy() Reference {
	return b
}

func (*Box) InstanceVariables() SymbolMap {
	return nil
}

func (b *Box) Inspect() string {
	valInspect := b.Get().Inspect()
	if !strings.ContainsRune(valInspect, '\n') {
		return fmt.Sprintf("Std::Box{&: %p, %s}", b, valInspect)
	}

	var buff strings.Builder

	fmt.Fprintf(&buff, "Std::Box{\n  &: %p", b)

	buff.WriteString(",\n  ")
	indent.IndentStringFromSecondLine(&buff, b.Get().Inspect(), 1)

	buff.WriteString("\n}")

	return buff.String()
}

func (b *Box) Error() string {
	return b.Inspect()
}
