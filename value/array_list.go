package value

// ::Std::ArrayList
//
// Represents a dynamically sized array with Elk value elements
// that can shrink and grow.
var ArrayListClass *Class

// ::Std::ArrayList::Iterator
//
// ArrayList iterator class.
var ArrayListIteratorClass *Class

type ArrayList interface {
	ArrayTuple
	Capacity() int
	LeftCapacity() int
	BoxOfVal(index Value) (Value, Value)
	SubscriptSet(index Value, v Value) Value
	SetAtVal(index int, v Value) Value
	Grow(int)
	AppendVal(elements ...Value) Value
	RemoveAt(i int)
	RemoveAtErr(index int) Value
	IterList() ArrayListIterator
	NewArrayList(capacity int) ArrayList
	CloneArrayList(capacity int) ArrayList
}

type ArrayListIterator interface {
	ArrayTupleIterator
}

func initArrayList() {
	ArrayListClass = NewClassWithOptions(ClassWithConstructor(ArrayListOfValueConstructor))
	ArrayListClass.IncludeMixin(ListMixin)
	StdModule.AddConstantString("ArrayList", Ref(ArrayListClass))
	RegisterNativeClass("Std::ArrayList", "value.ArrayListClass")

	ArrayListIteratorClass = NewClass()
	ArrayListClass.AddConstantString("Iterator", Ref(ArrayListIteratorClass))
	RegisterNativeClass("Std::ArrayList::Iterator", "value.ArrayListIteratorClass")
}
