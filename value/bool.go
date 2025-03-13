package value

var BoolClass *Class // ::Std::Bool

func initBool() {
	BoolClass = NewClassWithOptions(ClassWithParent(ValueClass))
	StdModule.AddConstantString("Bool", Ref(BoolClass))
}
