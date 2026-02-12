package value

import (
	"fmt"
	"iter"
)

type NativePair[K ComparableValueInterface, V ValueInterface] struct {
	key   K
	value V
}

var _ Pair = &NativePair[String, String]{}

func NewNativePair[K ComparableValueInterface, V ValueInterface](key K, val V) *NativePair[K, V] {
	return &NativePair[K, V]{
		key:   key,
		value: val,
	}
}

func MakeNativePair[K ComparableValueInterface, V ValueInterface](key K, val V) NativePair[K, V] {
	return NativePair[K, V]{
		key:   key,
		value: val,
	}
}

func (p *NativePair[K, V]) ToPairOfValue() *PairOfValue {
	return NewPairOfValue(p.Key(), p.Value())
}

func (p *NativePair[K, V]) NativeKey() K {
	return p.key
}

func (p *NativePair[K, V]) Key() Value {
	return p.key.ToValue()
}

func (p *NativePair[K, V]) SetKey(k Value) Value {
	kVal, ok := Downcast[K](k)
	if !ok {
		return NewInvalidKeyInTypedPair(p, k.Class()).ToValue()
	}

	p.key = kVal
	return Undefined
}

func (p *NativePair[K, V]) SetNativeKey(k K) {
	p.key = k
}

func (p *NativePair[K, V]) NativeValue() V {
	return p.value
}

func (p *NativePair[K, V]) Value() Value {
	return p.value.ToValue()
}

func (p *NativePair[K, V]) SetValue(v Value) Value {
	vVal, ok := Downcast[V](v)
	if !ok {
		return NewInvalidValueInTypedPair(p, v.Class()).ToValue()
	}

	p.value = vVal
	return Undefined
}

func (p *NativePair[K, V]) SetNativeValue(v V) {
	p.value = v
}

func (*NativePair[K, V]) Class() *Class {
	return PairClass
}

func (*NativePair[K, V]) DirectClass() *Class {
	return PairClass
}

func (*NativePair[K, V]) SingletonClass() *Class {
	return nil
}

func (p *NativePair[K, V]) Copy() Reference {
	return p
}

func (p *NativePair[K, V]) ToValue() Value {
	return Ref(p)
}

func (p *NativePair[K, V]) Inspect() string {
	return fmt.Sprintf("Std::Pair(%s, %s)", p.key.Inspect(), p.value.Inspect())
}

func (p *NativePair[K, V]) Error() string {
	return p.Inspect()
}

func (*NativePair[K, V]) InstanceVariables() *InstanceVariables {
	return nil
}

func (*NativePair[K, V]) Length() int {
	return pairLength
}

func (p *NativePair[K, V]) Iterate() iter.Seq2[Value, Value] {
	return func(yield func(Value, Value) bool) {
		if !yield(p.Key(), Undefined) {
			return
		}
		if !yield(p.Value(), Undefined) {
			return
		}
	}
}

// Get an element under the given index.
func (p *NativePair[K, V]) Subscript(key Value) (Value, Value) {
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
func (p *NativePair[K, V]) Get(index int) (Value, Value) {
	switch index {
	case 0, -2:
		return p.Key(), Undefined
	case 1, -1:
		return p.Value(), Undefined
	default:
		return Undefined, Ref(NewIndexOutOfRangeError(fmt.Sprint(index), pairLength))
	}
}

// Set an element under the given index.
func (p *NativePair[K, V]) SubscriptSet(key, val Value) Value {
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
func (p *NativePair[K, V]) Set(index int, val Value) Value {
	switch index {
	case 0, -2:
		k, ok := Downcast[K](val)
		if !ok {
			return NewInvalidKeyInTypedPair(p, val.Class()).ToValue()
		}
		p.key = k
		return Undefined
	case 1, -1:
		v, ok := Downcast[V](val)
		if !ok {
			return NewInvalidValueInTypedPair(p, val.Class()).ToValue()
		}
		p.value = v
		return Undefined
	default:
		return Ref(NewIndexOutOfRangeError(fmt.Sprint(index), pairLength))
	}
}

type NativePairIterator[K ComparableValueInterface, V ValueInterface] struct {
	Pair  *NativePair[K, V]
	Index int
}

var _ PairIterator = &NativePairIterator[String, String]{}

func NewNativePairIterator[K ComparableValueInterface, V ValueInterface](pair *NativePair[K, V]) *NativePairIterator[K, V] {
	return &NativePairIterator[K, V]{
		Pair: pair,
	}
}

func NewNativePairIteratorWithIndex[K ComparableValueInterface, V ValueInterface](pair *NativePair[K, V], index int) *NativePairIterator[K, V] {
	return &NativePairIterator[K, V]{
		Pair:  pair,
		Index: index,
	}
}

func (*NativePairIterator[K, V]) Class() *Class {
	return PairIteratorClass
}

func (*NativePairIterator[K, V]) DirectClass() *Class {
	return PairIteratorClass
}

func (*NativePairIterator[K, V]) SingletonClass() *Class {
	return nil
}

func (l *NativePairIterator[K, V]) Copy() Reference {
	return &NativePairIterator[K, V]{
		Pair:  l.Pair,
		Index: l.Index,
	}
}

func (i *NativePairIterator[K, V]) ToValue() Value {
	return Ref(i)
}

func (l *NativePairIterator[K, V]) Inspect() string {
	return fmt.Sprintf("Std::Pair::Iterator{&: %p, pair: %s, index: %d}", l, l.Pair.Inspect(), l.Index)
}

func (l *NativePairIterator[K, V]) Error() string {
	return l.Inspect()
}

func (*NativePairIterator[K, V]) InstanceVariables() *InstanceVariables {
	return nil
}

func (l *NativePairIterator[K, V]) NextValue() (Value, Value) {
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
