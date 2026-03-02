package value

var BoxClass *Class // ::Std::Box

type Box interface {
	ImmutableBox
	SetValue(Value)
	ToImmutableBoxInterface() ImmutableBox
}

func initBox() {
	BoxClass = NewClassWithOptions(ClassWithConstructor(BoxOfValueConstructor), ClassWithSuperclass(ImmutableBoxClass))
	StdModule.AddConstantString("Box", Ref(BoxClass))
	RegisterNativeClass("Std::Box", "value.BoxClass")
}
