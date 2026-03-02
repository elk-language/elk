package value

var FalseClass *Class // ::Std::False

func initFalse() {
	FalseClass = NewClassWithOptions(
		ClassWithSuperclass(BoolClass),
	)
	StdModule.AddConstantString("False", Ref(FalseClass))
	RegisterNativeClass("Std::False", "value.FalseClass")
}
