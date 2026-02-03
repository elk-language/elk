package value

// ::Std::ArrayList
//
// Represents a dynamically sized array with Elk value elements
// that can shrink and grow.
var ArrayListClass *Class

type ArrayList interface {
	ArrayTuple
}

func initArrayList() {
	ArrayListClass = NewClassWithOptions(ClassWithConstructor(ArrayListOfValueConstructor))
	ArrayListClass.IncludeMixin(ListMixin)
	StdModule.AddConstantString("ArrayList", Ref(ArrayListClass))
	RegisterNativeClass("Std::ArrayList", "value.ArrayListClass")
}
