package value

var ElkLexerClass *Class // ::Std::Elk::Lexer

func initElkLexer() {
	ElkLexerClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	ElkModule.AddConstantString("Lexer", Ref(ElkLexerClass))
}
