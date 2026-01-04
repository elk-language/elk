package value

var IntClass *Class         // ::Std::Int
var IntIteratorClass *Class // ::Std::Int::Iterator

// All simple Elk integer types (without BigInt)
type SimpleInt interface {
	SmallInt | Int64 | Int32 | Int16 | Int8 | UInt | UInt64 | UInt32 | UInt16 | UInt8
}

type SingedInt interface {
	SmallInt | Int64 | Int32 | Int16 | Int8
}

type UnsignedInt interface {
	UInt | UInt64 | UInt32 | UInt16 | UInt8
}

func initInt() {
	IntClass = NewClassWithOptions(ClassWithSuperclass(ValueClass))
	StdModule.AddConstantString("Int", Ref(IntClass))
	RegisterNativeClass("Std::Int", "value.IntClass")

	IntIteratorClass = NewClass()
	IntClass.AddConstantString("Iterator", Ref(IntIteratorClass))
	RegisterNativeClass("Std::Int::Iterator", "value.IntIteratorClass")

	IntClass.AddConstantString("Convertible", Ref(NewInterface()))
}
