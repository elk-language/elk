package value

var LocalBoxClass *Class // ::Std::LocalBox

func initLocalBox() {
	LocalBoxClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	StdModule.AddConstantString("LocalBox", Ref(LocalBoxClass))
}
