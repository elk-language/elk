package value

var ElkTokenClass *Class // Std::ElkToken

func initToken() {
	ElkTokenClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	StdModule.AddConstantString("ElkToken", Ref(ElkTokenClass))
}
