package value

var FSModule *Module // ::Std::FS

func initFS() {
	FSModule = NewModule()
	StdModule.AddConstantString("FS", Ref(FSModule))
}
