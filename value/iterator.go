package value

// ::Std::Iterator
var IteratorInterface *Interface

// ::Std::Iterator::Base
var IteratorBaseMixin *Mixin

func initIterator() {
	IteratorInterface = NewInterface()
	StdModule.AddConstantString("Iterator", Ref(IteratorInterface))

	IteratorBaseMixin = NewMixin()
	IteratorInterface.AddConstantString("Base", Ref(IteratorBaseMixin))
}
