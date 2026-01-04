package value

var ColorizerInterface *Interface // ::Std::Colorizer

func initColorizer() {
	ColorizerInterface = NewInterface()
	StdModule.AddConstantString("Colorizer", Ref(ColorizerInterface))
	RegisterNativeInterface("Std::Colorizer", "value.ColorizerInterface")
}
