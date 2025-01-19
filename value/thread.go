package value

var ThreadClass *Class // ::Std::Thread

func initThread() {
	ThreadClass = NewClass()
	StdModule.AddConstantString("Thread", Ref(ThreadClass))
}
