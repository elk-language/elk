package value

var DebugModule *Module // ::Std::Debug

func initDebug() {
	DebugModule = NewModule()
	StdModule.AddConstantString("Debug", Ref(DebugModule))
}
