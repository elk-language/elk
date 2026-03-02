package value

// ::Std::Iterator
var IteratorInterface *Interface

// ::Std::Iterator::Base
var IteratorBaseMixin *Mixin

func initIterator() {
	IteratorInterface = NewInterface()
	StdModule.AddConstantString("Iterator", Ref(IteratorInterface))
	RegisterNativeClass("Std::Iterator", "value.IteratorInterface")

	IteratorBaseMixin = NewMixin()
	IteratorInterface.AddConstantString("Base", Ref(IteratorBaseMixin))
	RegisterNativeMixin("Std::Iterator::Base", "value.IteratorBaseMixin")
}
