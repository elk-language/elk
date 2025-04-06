package value

var IntClass *Class         // ::Std::Int
var IntIteratorClass *Class // ::Std::Int::Iterator

// All simple Elk integer types (without BigInt)
type SimpleInt interface {
	SmallInt | Int64 | Int32 | Int16 | Int8 | UInt64 | UInt32 | UInt16 | UInt8
}

type SingedInt interface {
	SmallInt | Int64 | Int32 | Int16 | Int8
}

type UnsignedInt interface {
	UInt64 | UInt32 | UInt16 | UInt8
}

func initInt() {
	IntClass = NewClassWithOptions(ClassWithParent(ValueClass))
	StdModule.AddConstantString("Int", Ref(IntClass))

	IntIteratorClass = NewClass()
	IntClass.AddConstantString("Iterator", Ref(IntIteratorClass))

	IntClass.AddConstantString("Convertible", Ref(NewInterface()))
}
