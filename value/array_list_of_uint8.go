package value

import (
	"fmt"
)

// ::Std::NativeArrayList::UInt8Iterator
//
// ArrayListOfUInt8 iterator class.
var ArrayListOfUInt8IteratorClass *Class

type ArrayListOfUInt8 = nativeArrayList[UInt8]

func NewArrayListOfUInt8(capacity int) *ArrayListOfUInt8 {
	return newNativeArrayList[UInt8](capacity)
}

func NewArrayListOfUInt8WithLength(length int) *ArrayListOfUInt8 {
	return newNativeArrayListWithLength[UInt8](length)
}

func NewArrayListOfUInt8WithElements(capacity int, elements ...UInt8) *ArrayListOfUInt8 {
	return newNativeArrayListWithElements(capacity, elements...)
}

type ArrayListOfUInt8Iterator struct {
	nativeArrayListIterator[UInt8]
}

func NewArrayListOfUInt8Iterator(list *ArrayListOfUInt8) *ArrayListOfUInt8Iterator {
	return &ArrayListOfUInt8Iterator{
		nativeArrayListIterator: nativeArrayListIterator[UInt8]{
			ArrayList: list,
		},
	}
}

func NewArrayListOfUInt8IteratorWithIndex(list *ArrayListOfUInt8, index int) *ArrayListOfUInt8Iterator {
	return &ArrayListOfUInt8Iterator{
		nativeArrayListIterator: nativeArrayListIterator[UInt8]{
			ArrayList: list,
			Index:     index,
		},
	}
}

func (*ArrayListOfUInt8Iterator) Class() *Class {
	return ArrayListOfUInt8IteratorClass
}

func (*ArrayListOfUInt8Iterator) DirectClass() *Class {
	return ArrayListOfUInt8IteratorClass
}

func (l *ArrayListOfUInt8Iterator) Copy() Reference {
	return &ArrayListOfUInt8Iterator{
		nativeArrayListIterator: nativeArrayListIterator[UInt8]{
			ArrayList: l.ArrayList,
			Index:     l.Index,
		},
	}
}

func (i *ArrayListOfUInt8Iterator) ToValue() Value {
	return Ref(i)
}

func (l *ArrayListOfUInt8Iterator) Inspect() string {
	return fmt.Sprintf("Std::ArrayList::UInt8Iterator%s", l.inspect())
}

func (l *ArrayListOfUInt8Iterator) Error() string {
	return l.Inspect()
}

func initArrayListOfUInt8() {
	ArrayListOfUInt8IteratorClass = NewClass()
	ArrayListClass.AddConstantString("UInt8Iterator", Ref(ArrayListOfUInt8IteratorClass))
	RegisterNativeClass("Std::ArrayList::UInt8Iterator", "value.ArrayListOfUInt8IteratorClass")
}
