package value

var ElkModule *Module // Std::Elk

func initElk() {
	ElkModule = NewModule()
	StdModule.AddConstantString("Elk", Ref(ElkModule))
}
