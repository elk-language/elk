package value

var NilClass *Class // ::Std::Nil

type NilType struct{}

// Elk's Nil value
var Nil = NilType{}

func (NilType) Class() *Class {
	return NilClass
}

func (NilType) DirectClass() *Class {
	return NilClass
}

func (NilType) SingletonClass() *Class {
	return nil
}

func (NilType) Inspect() string {
	return "nil"
}

func (n NilType) Copy() Value {
	return n
}

func (NilType) InstanceVariables() SymbolMap {
	return nil
}

func initNil() {
	NilClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
	)
	StdModule.AddConstantString("Nil", NilClass)
}
