package value

var ThreadPoolClass *Class // ::Std::ThreadPool

func initThreadPool() {
	ThreadPoolClass = NewClass()
	StdModule.AddConstantString("ThreadPool", Ref(ThreadPoolClass))
}
