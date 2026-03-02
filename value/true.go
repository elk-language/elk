package value

var TrueClass *Class // ::Std::True

func initTrue() {
	TrueClass = NewClassWithOptions(
		ClassWithSuperclass(BoolClass),
	)
	StdModule.AddConstantString("True", Ref(TrueClass))
	RegisterNativeClass("Std::True", "value.TrueClass")
}
