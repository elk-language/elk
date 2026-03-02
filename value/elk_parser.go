package value

var ElkParserClass *Class       // ::Std::Elk::Parser
var ElkParserResultClass *Class // ::Std::Elk::Parser::Result

func initElkParser() {
	ElkParserClass = NewClass()
	ElkModule.AddConstantString("Parser", Ref(ElkParserClass))
	RegisterNativeClass("Std::Elk::Parser", "value.ElkParserClass")

	ElkParserResultClass = NewClass()
	ElkParserClass.AddConstantString("Result", Ref(ElkParserResultClass))
	RegisterNativeClass("Std::Elk::Parser::Result", "value.ElkParserResultClass")
}
