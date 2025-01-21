package value

import "fmt"

var PairClass *Class // ::Std::Pair

// ::Std::Pair::Iterator
//
// Pair iterator class.
var PairIteratorClass *Class

type Pair struct {
	Key   Value
	Value Value
}

// Creates a new Pair.
func PairConstructor(class *Class) Value {
	return Ref(&Pair{})
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

func (p *Pair) Copy() Reference {
	return p
}

func (p *Pair) Inspect() string {
	return fmt.Sprintf("Std::Pair(%s, %s)", p.Key.Inspect(), p.Value.Inspect())
}

func (p *Pair) Error() string {
	return p.Inspect()
}

func (*Pair) InstanceVariables() SymbolMap {
	return nil
}

const pairLength = 2

func (*Pair) Length() int {
	return pairLength
}

// Get an element under the given index.
func (p *Pair) Subscript(key Value) (Value, Value) {
	var i int

	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return Undefined, Ref(NewIndexOutOfRangeError(key.Inspect(), pairLength))
		}
		return Undefined, Ref(NewCoerceError(IntClass, key.Class()))
	}

	return p.Get(i)
}

// Get an element under the given index.
func (p *Pair) Get(index int) (Value, Value) {
	switch index {
	case 0, -2:
		return p.Key, Undefined
	case 1, -1:
		return p.Value, Undefined
	default:
		return Undefined, Ref(NewIndexOutOfRangeError(fmt.Sprint(index), pairLength))
	}
}

// Set an element under the given index.
func (p *Pair) SubscriptSet(key, val Value) Value {
	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return Ref(NewIndexOutOfRangeError(key.Inspect(), pairLength))
		}
		return Ref(NewCoerceError(IntClass, key.Class()))
	}

	return p.Set(i, val)
}

// Set an element under the given index.
func (p *Pair) Set(index int, val Value) Value {
	switch index {
	case 0, -2:
		p.Key = val
		return Undefined
	case 1, -1:
		p.Value = val
		return Undefined
	default:
		return Ref(NewIndexOutOfRangeError(fmt.Sprint(index), pairLength))
	}
}

type PairIterator struct {
	Pair  *Pair
	Index int
}

func NewPairIterator(pair *Pair) *PairIterator {
	return &PairIterator{
		Pair: pair,
	}
}

func NewPairIteratorWithIndex(pair *Pair, index int) *PairIterator {
	return &PairIterator{
		Pair:  pair,
		Index: index,
	}
}

func (*PairIterator) Class() *Class {
	return PairIteratorClass
}

func (*PairIterator) DirectClass() *Class {
	return PairIteratorClass
}

func (*PairIterator) SingletonClass() *Class {
	return nil
}

func (l *PairIterator) Copy() Reference {
	return &PairIterator{
		Pair:  l.Pair,
		Index: l.Index,
	}
}

func (l *PairIterator) Inspect() string {
	return fmt.Sprintf("Std::Pair::Iterator{&: %p, pair: %s, index: %d}", l, l.Pair.Inspect(), l.Index)
}

func (l *PairIterator) Error() string {
	return l.Inspect()
}

func (*PairIterator) InstanceVariables() SymbolMap {
	return nil
}

func (l *PairIterator) Next() (Value, Value) {
	if l.Index >= pairLength {
		return Undefined, stopIterationSymbol.ToValue()
	}

	next, err := l.Pair.Get(l.Index)
	if !err.IsUndefined() {
		return Undefined, err
	}

	l.Index++
	return next, Undefined
}

func initPair() {
	PairClass = NewClassWithOptions(
		ClassWithConstructor(PairConstructor),
	)
	PairClass.IncludeMixin(TupleMixin)
	StdModule.AddConstantString("Pair", Ref(PairClass))

	PairIteratorClass = NewClass()
	PairClass.AddConstantString("Iterator", Ref(PairIteratorClass))
}
