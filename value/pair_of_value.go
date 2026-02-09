package value

import (
	"fmt"
	"iter"
)

type PairOfValue struct {
	key   Value
	value Value
}

var _ Pair = &PairOfValue{}

// Creates a new Pair.
func PairOfValueConstructor(class *Class) Value {
	return Ref(&PairOfValue{})
}

func NewPairOfValue(key, val Value) *PairOfValue {
	return &PairOfValue{
		key:   key,
		value: val,
	}
}

func MakePairOfValue(key, val Value) PairOfValue {
	return PairOfValue{
		key:   key,
		value: val,
	}
}

func (p *PairOfValue) Key() Value {
	return p.key
}

func (p *PairOfValue) SetKey(k Value) Value {
	p.key = k
	return Undefined
}

func (p *PairOfValue) Value() Value {
	return p.value
}

func (p *PairOfValue) SetValue(v Value) Value {
	p.value = v
	return Undefined
}

func (*PairOfValue) Class() *Class {
	return PairClass
}

func (*PairOfValue) DirectClass() *Class {
	return PairClass
}

func (*PairOfValue) SingletonClass() *Class {
	return nil
}

func (p *PairOfValue) Copy() Reference {
	return p
}

func (p *PairOfValue) ToValue() Value {
	return Ref(p)
}

func (p *PairOfValue) Inspect() string {
	return fmt.Sprintf("Std::Pair(%s, %s)", p.key.Inspect(), p.value.Inspect())
}

func (p *PairOfValue) Error() string {
	return p.Inspect()
}

func (*PairOfValue) InstanceVariables() *InstanceVariables {
	return nil
}

const pairLength = 2

func (*PairOfValue) Length() int {
	return pairLength
}

func (p *PairOfValue) Iterate() iter.Seq2[Value, Value] {
	return func(yield func(Value, Value) bool) {
		if !yield(p.key, Undefined) {
			return
		}
		if !yield(p.value, Undefined) {
			return
		}
	}
}

// Get an element under the given index.
func (p *PairOfValue) Subscript(key Value) (Value, Value) {
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
func (p *PairOfValue) Get(index int) (Value, Value) {
	switch index {
	case 0, -2:
		return p.key, Undefined
	case 1, -1:
		return p.value, Undefined
	default:
		return Undefined, Ref(NewIndexOutOfRangeError(fmt.Sprint(index), pairLength))
	}
}

// Set an element under the given index.
func (p *PairOfValue) SubscriptSet(key, val Value) Value {
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
func (p *PairOfValue) Set(index int, val Value) Value {
	switch index {
	case 0, -2:
		p.key = val
		return Undefined
	case 1, -1:
		p.value = val
		return Undefined
	default:
		return Ref(NewIndexOutOfRangeError(fmt.Sprint(index), pairLength))
	}
}

type PairOfValueIterator struct {
	Pair  *PairOfValue
	Index int
}

func NewPairOfValueIterator(pair *PairOfValue) *PairOfValueIterator {
	return &PairOfValueIterator{
		Pair: pair,
	}
}

func NewPairOfValueIteratorWithIndex(pair *PairOfValue, index int) *PairOfValueIterator {
	return &PairOfValueIterator{
		Pair:  pair,
		Index: index,
	}
}

func (*PairOfValueIterator) Class() *Class {
	return PairIteratorClass
}

func (*PairOfValueIterator) DirectClass() *Class {
	return PairIteratorClass
}

func (*PairOfValueIterator) SingletonClass() *Class {
	return nil
}

func (l *PairOfValueIterator) Copy() Reference {
	return &PairOfValueIterator{
		Pair:  l.Pair,
		Index: l.Index,
	}
}

func (i *PairOfValueIterator) ToValue() Value {
	return Ref(i)
}

func (l *PairOfValueIterator) Inspect() string {
	return fmt.Sprintf("Std::Pair::Iterator{&: %p, pair: %s, index: %d}", l, l.Pair.Inspect(), l.Index)
}

func (l *PairOfValueIterator) Error() string {
	return l.Inspect()
}

func (*PairOfValueIterator) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *PairOfValueIterator) NextValue() (Value, Value) {
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
