package value

var BoolClass *Class // ::Std::Bool

func initBool() {
	BoolClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	StdModule.AddConstantString("Bool", Ref(BoolClass))
}
