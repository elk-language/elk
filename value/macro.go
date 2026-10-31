package value

var MacroModule *Module      // Std::Macro
var MacroContextClass *Class // Std::Macro::Context

func initMacro() {
	MacroModule = NewModule()
	StdModule.AddConstantString("Macro", Ref(MacroModule))
	RegisterNativeModule("Std::Macro", "value.MacroModule")

	MacroContextClass = NewClass()
	MacroModule.AddConstantString("Context", Ref(MacroContextClass))
	RegisterNativeClass("Std::Macro::Context", "value.MacroContextClass")
}
