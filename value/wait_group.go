package value

import "sync"

var WaitGroupClass *Class // ::Std::WaitGroup

type WaitGroup struct {
	Native sync.WaitGroup
}

func (w *WaitGroup) Copy() Reference {
	return &WaitGroup{}
}

func (*WaitGroup) Class() *Class {
	return WaitGroupClass
}

func (*WaitGroup) DirectClass() *Class {
	return WaitGroupClass
}

func (*WaitGroup) SingletonClass() *Class {
	return nil
}

func (w *WaitGroup) Inspect() string {
	return "Std::WaitGroup{}"
}

func (w *WaitGroup) Error() string {
	return w.Inspect()
}

func (w *WaitGroup) InstanceVariables() SymbolMap {
	return nil
}

func (w *WaitGroup) Add(n int) {
	w.Native.Add(n)
}

func (w *WaitGroup) Done() {
	w.Native.Done()
}

func (w *WaitGroup) Remove(n int) {
	for range n {
		w.Native.Done()
	}
}

func (w *WaitGroup) Wait() {
	w.Native.Wait()
}

func initWaitGroup() {
	WaitGroupClass = NewClass()
	StdModule.AddConstantString("WaitGroup", Ref(WaitGroupClass))
}
