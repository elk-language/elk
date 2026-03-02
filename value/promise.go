package value

var PromiseClass *Class // ::Std::Promise

func initPromise() {
	PromiseClass = NewClass()
	StdModule.AddConstantString("Promise", Ref(PromiseClass))
	RegisterNativeClass("Std::Promise", "value.PromiseClass")
}
