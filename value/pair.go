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
			return nil, NewIndexOutOfRangeError(key.Inspect(), pairLength)
		}
		return nil, NewCoerceError(IntClass, key.Class())
	}

	return p.Get(i)
}

// Get an element under the given index.
func (p *Pair) Get(index int) (Value, Value) {
	switch index {
	case 0, -2:
		return p.Key, nil
	case 1, -1:
		return p.Value, nil
	default:
		return nil, NewIndexOutOfRangeError(fmt.Sprint(index), pairLength)
	}
}

// Set an element under the given index.
func (p *Pair) SubscriptSet(key, val Value) Value {
	i, ok := ToGoInt(key)
	if !ok {
		if i == -1 {
			return NewIndexOutOfRangeError(key.Inspect(), pairLength)
		}
		return NewCoerceError(IntClass, key.Class())
	}

	return p.Set(i, val)
}

// Set an element under the given index.
func (p *Pair) Set(index int, val Value) Value {
	switch index {
	case 0, -2:
		p.Key = val
		return nil
	case 1, -1:
		p.Value = val
		return nil
	default:
		return NewIndexOutOfRangeError(fmt.Sprint(index), pairLength)
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

func (l *PairIterator) Copy() Value {
	return &PairIterator{
		Pair:  l.Pair,
		Index: l.Index,
	}
}

func (l *PairIterator) Inspect() string {
	return fmt.Sprintf("Std::Pair::Iterator{pair: %s, index: %d}", l.Pair.Inspect(), l.Index)
}

func (l *PairIterator) Error() string {
	return l.Inspect()
}

func (*PairIterator) InstanceVariables() SymbolMap {
	return nil
}

func (l *PairIterator) Next() (Value, Value) {
	if l.Index >= pairLength {
		return nil, stopIterationSymbol
	}

	next, err := l.Pair.Get(l.Index)
	if err != nil {
		return nil, err
	}

	l.Index++
	return next, nil
}

func initPair() {
	PairClass = NewClassWithOptions(
		ClassWithConstructor(PairConstructor),
	)
	PairClass.IncludeMixin(TupleMixin)
	StdModule.AddConstantString("Pair", PairClass)

	PairIteratorClass = NewClass()
	PairClass.AddConstantString("Iterator", PairIteratorClass)
}
