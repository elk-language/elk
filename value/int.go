package value

var IntClass *Class // ::Std::Int

// All simple Elk integer types (without BigInt)
type SimpleInt interface {
	SmallInt | Int64 | Int32 | Int16 | Int8 | UInt64 | UInt32 | UInt16 | UInt8
	Value
}

func initInt() {
	IntClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	StdModule.AddConstantString("Int", IntClass)
}
