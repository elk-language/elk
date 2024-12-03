package value

var BoolClass *Class // ::Std::Bool

func initBool() {
	BoolClass = NewClass()
	StdModule.AddConstantString("Bool", ReferenceToValue(BoolClass))
}
