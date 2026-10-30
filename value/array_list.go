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
	SubscriptSetInt(index int, v Value) Value
	SetAtVal(index int, v Value) Value
	Grow(int)
	AppendVal(elements ...Value) Value
	RemoveAt(i int)
	RemoveAtVal(index Value) Value
	// RemoveAtRange removes the elements arr[i..j] in place
	RemoveAtRange(i, j int)
	IterList() ArrayListIterator
	NewArrayList(capacity int) ArrayList
	CloneArrayList(capacity int) ArrayList
	SliceArrayList(from, to int) ArrayList
	ToArrayTuple() ArrayTuple
	ToSyncArrayList() SyncArrayList
}

type ArrayListIterator interface {
	ArrayTupleIterator
}

func ArrayListInterfaceRemoveAtRange(l ArrayList, rng Value) Value {
	var i, j int
	switch rng := rng.SafeAsReference().(type) {
	case *OpenRange:
		start := rng.Start.AsInt()
		end := rng.End.AsInt()
		i = start + 1
		j = end - 1
	case *LeftOpenRange:
		start := rng.Start.AsInt()
		end := rng.End.AsInt()
		i = start + 1
		j = end
	case *RightOpenRange:
		start := rng.Start.AsInt()
		end := rng.End.AsInt()
		i = start
		j = end - 1
	case *ClosedRange:
		start := rng.Start.AsInt()
		end := rng.End.AsInt()
		i = start
		j = end
	case *BeginlessOpenRange:
		end := rng.End.AsInt()
		i = 0
		j = end - 1
	case *BeginlessClosedRange:
		end := rng.End.AsInt()
		i = 0
		j = end
	case *EndlessOpenRange:
		start := rng.Start.AsInt()
		i = start + 1
		j = l.Length() - 1
	case *EndlessClosedRange:
		start := rng.Start.AsInt()
		i = start
		j = l.Length() - 1
	default:
		return NewCoerceError(RangeMixin, rng.Class()).ToValue()
	}

	l.RemoveAtRange(i, j)
	return Undefined
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
