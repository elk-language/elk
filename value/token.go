package value

var TokenClass *Class // Std::Token

func initToken() {
	TokenClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	StdModule.AddConstantString("Token", Ref(TokenClass))
}
