package value

var RuntimeModule *Module // ::Std::Runtime

func initRuntime() {
	RuntimeModule = NewModule()
	StdModule.AddConstantString("Runtime", Ref(RuntimeModule))
}
