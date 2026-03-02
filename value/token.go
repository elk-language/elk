package value

var ElkTokenClass *Class // Std::Elk::Token

func initElkToken() {
	ElkTokenClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ElkModule.AddConstantString("Token", Ref(ElkTokenClass))
	RegisterNativeClass("Std::Elk::Token", "value.ElkTokenClass")
}
