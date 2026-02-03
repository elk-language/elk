package value

import "iter"

// ::Std::ArrayTuple
//
// Represents an immutable array.
var ArrayTupleClass *Class

type ArrayTuple interface {
	ValueInterface
	Elements() iter.Seq2[int, Value]
	Length() int
}

func initArrayTuple() {
	ArrayTupleClass = NewClassWithOptions(ClassWithConstructor(ArrayTupleOfValueConstructor))
	ArrayTupleClass.IncludeMixin(TupleMixin)
	StdModule.AddConstantString("ArrayTuple", Ref(ArrayTupleClass))
	RegisterNativeClass("Std::ArrayTuple", "value.ArrayTupleClass")

	ArrayTupleOfValueIteratorClass = NewClass()
	ArrayTupleClass.AddConstantString("Iterator", Ref(ArrayTupleOfValueIteratorClass))
	RegisterNativeClass("Std::ArrayTuple::Iterator", "value.ArrayTupleIteratorClass")
}
