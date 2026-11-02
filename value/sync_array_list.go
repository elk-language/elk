package value

import "iter"

// ::Std::Sync::ArrayList
//
// Represents a thread-safe dynamically sized array with Elk value elements
// that can shrink and grow.
var SyncArrayListClass *Class

// ::Std::Sync::ArrayList::Iterator
//
// Sync ArrayList iterator class.
var SyncArrayListIteratorClass *Class

type SyncArrayList interface {
	Tuple
	Lock()
	Unlock()
	ArrayList() ArrayList
	Capacity() int
	LeftCapacity() int
	SubscriptSet(index Value, v Value) Value
	SubscriptSetInt(index int, v Value) Value
	SetAtVal(index int, v Value) Value
	Grow(int)
	AppendVal(elements ...Value) Value
	RemoveAt(i int)
	RemoveAtVal(index Value) Value
	// RemoveAtRange removes the elements arr[i..j] in place
	RemoveAtRange(i, j int)
	IterSyncList() SyncArrayListIterator
	NewSyncArrayList(capacity int) SyncArrayList
	CloneSyncArrayList(capacity int) SyncArrayList
	SliceSyncArrayList(from, to int) SyncArrayList

	// normalizes the index so that `-1` becomes `length-1`
	// does bounds checking and returns an elk error
	SubscriptInt(index int) (Value, Value)
	// normalizes the index so that `-1` becomes `length-1`
	// does bounds checking and returns an elk error
	Subscript(index Value) (Value, Value)
	// normalizes the index so that `-1` becomes `length-1`
	// does bounds checking and returns elk `nil` if the index is out of range
	TryAtVal(Value) (Value, Value)
	ConcatVal(Value) (Value, Value)
	ConcatValOk(Value) (Value, bool)
	RepeatVal(Value) (Value, Value)
	ToArrayList() ArrayList
	ToArrayTuple() ArrayTuple
}

type SyncArrayListIterator interface {
	NativeResettableIterator
	Elements() iter.Seq[Value]
}

func initSyncArrayList() {
	SyncArrayListClass = NewClassWithOptions(ClassWithConstructor(SyncArrayListOfValueConstructor))
	SyncArrayListClass.IncludeMixin(ListMixin)
	SyncModule.AddConstantString("ArrayList", Ref(SyncArrayListClass))
	RegisterNativeClass("Std::Sync::ArrayList", "value.SyncArrayListClass")

	SyncArrayListIteratorClass = NewClass()
	SyncArrayListClass.AddConstantString("Iterator", Ref(SyncArrayListIteratorClass))
	RegisterNativeClass("Std::Sync::ArrayList::Iterator", "value.SyncArrayListIteratorClass")
}
