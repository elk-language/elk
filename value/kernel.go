package value

var KernelModule *Module // ::Std::Kernel

func initKernel() {
	KernelModule = NewModule()
	StdModule.AddConstantString("Kernel", Ref(KernelModule))
	RegisterNativeModule("Std::Kernel", "value.KernelModule")
}
