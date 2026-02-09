package value

var PairClass *Class // ::Std::Pair

// ::Std::Pair::Iterator
//
// Pair iterator class.
var PairIteratorClass *Class

type Pair interface {
	NativeIterable
	Key() Value
	SetKey(Value) Value
	Value() Value
	SetValue(Value) Value
}

type PairIterator interface {
	NativeIterator
}

func initPair() {
	PairClass = NewClassWithOptions(
		ClassWithConstructor(PairOfValueConstructor),
	)
	PairClass.IncludeMixin(TupleMixin)
	StdModule.AddConstantString("Pair", Ref(PairClass))
	RegisterNativeClass("Std::Pair", "value.PairClass")

	PairIteratorClass = NewClass()
	PairClass.AddConstantString("Iterator", Ref(PairIteratorClass))
	RegisterNativeClass("Std::Pair::Iterator", "value.PairIteratorClass")
}
