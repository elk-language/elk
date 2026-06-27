package value

var MacroModule *Module // Std::Macro

func initMacro() {
	MacroModule = NewModule()
	StdModule.AddConstantString("Macro", Ref(MacroModule))
}
