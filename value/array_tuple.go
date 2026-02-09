package value

import "iter"

// ::Std::ArrayTuple
//
// Represents an immutable array.
var ArrayTupleClass *Class

// ::Std::ArrayTuple::Iterator
//
// ArrayTuple iterator class.
var ArrayTupleIteratorClass *Class

type ArrayTuple interface {
	ValueInterface
	NativeIterable
	Elements() iter.Seq2[int, Value]
	Length() int
	ImmutableBoxOfVal(index Value) (Value, Value)
	AtVal(int) Value
	Subscript(index Value) (Value, Value)
	ConcatVal(Value) (Value, Value)
	RepeatVal(Value) (Value, Value)
	IterTuple() ArrayTupleIterator
}

type ArrayTupleIterator interface {
	NativeIterator
	ValueInterface
	Reset()
	Elements() iter.Seq[Value]
}

func initArrayTuple() {
	ArrayTupleClass = NewClassWithOptions(ClassWithConstructor(ArrayTupleOfValueConstructor))
	ArrayTupleClass.IncludeMixin(TupleMixin)
	StdModule.AddConstantString("ArrayTuple", Ref(ArrayTupleClass))
	RegisterNativeClass("Std::ArrayTuple", "value.ArrayTupleClass")

	ArrayTupleIteratorClass = NewClass()
	ArrayTupleClass.AddConstantString("Iterator", Ref(ArrayTupleIteratorClass))
	RegisterNativeClass("Std::ArrayTuple::Iterator", "value.ArrayTupleIteratorClass")
}
