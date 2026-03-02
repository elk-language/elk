package value

var ImmutableBoxClass *Class // ::Std::ImmutableBox

type ImmutableBox interface {
	ValueInterface
	GetValue() Value
	Address() UInt
}

func initImmutableBox() {
	ImmutableBoxClass = NewClassWithOptions(ClassWithConstructor(ImmutableBoxOfValueConstructor))
	StdModule.AddConstantString("ImmutableBox", Ref(ImmutableBoxClass))
	RegisterNativeClass("Std::ImmutableBox", "value.ImmutableBoxClass")
}
