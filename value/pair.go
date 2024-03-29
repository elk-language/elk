package value

import "fmt"

var PairClass *Class // ::Std::Pair

type Pair struct {
	Key   Value
	Value Value
}

// Creates a new Pair.
func PairConstructor(class *Class) Value {
	return &Pair{}
}

func NewPair(key, val Value) *Pair {
	return &Pair{
		Key:   key,
		Value: val,
	}
}

func (*Pair) Class() *Class {
	return PairClass
}

func (*Pair) DirectClass() *Class {
	return PairClass
}

func (*Pair) SingletonClass() *Class {
	return nil
}

func (p *Pair) Copy() Value {
	return p
}

func (p *Pair) Inspect() string {
	return fmt.Sprintf("Std::Pair{key: %s, value: %s}", p.Key.Inspect(), p.Value.Inspect())
}

func (*Pair) InstanceVariables() SymbolMap {
	return nil
}

func initPair() {
	PairClass = NewClassWithOptions(
		ClassWithNoInstanceVariables(),
		ClassWithSealed(),
		ClassWithConstructor(PairConstructor),
	)
	StdModule.AddConstantString("Pair", PairClass)
}
