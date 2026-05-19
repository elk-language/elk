package value

import (
	"context"
	"fmt"
	"iter"
)

type NativeReadChannel[V ValueInterface] <-chan V

var _ ReadChannel = make(NativeReadChannel[Float])

func MakeNativeReadChannel[V ValueInterface](size int) NativeReadChannel[V] {
	return make(NativeReadChannel[V], size)
}

func (ch NativeReadChannel[V]) Copy() Reference {
	return ch
}

func (c NativeReadChannel[V]) ToValue() Value {
	return Ref(c)
}

func (NativeReadChannel[V]) Class() *Class {
	return ReadChannelClass
}

func (NativeReadChannel[V]) DirectClass() *Class {
	return ReadChannelClass
}

func (NativeReadChannel[V]) SingletonClass() *Class {
	return nil
}

func (ch NativeReadChannel[V]) NativeChannelAny() any {
	return (<-chan V)(ch)
}

func (ch NativeReadChannel[V]) TransformToValue(v any) Value {
	panic("not a transformer channel")
}

func (ch NativeReadChannel[V]) TransformFromValue(v Value) any {
	panic("not a transformer channel")
}

func (NativeReadChannel[V]) IsTransformerChannel() bool {
	return false
}

func (ch NativeReadChannel[V]) Inspect() string {
	return fmt.Sprintf("Std::ReadChannel{&: %p, length: %d, capacity: %d}", ch, ch.Length(), ch.Capacity())
}

func (ch NativeReadChannel[V]) Error() string {
	return ch.Inspect()
}

func (NativeReadChannel[V]) InstanceVariables() *InstanceVariables {
	return nil
}

func (ch NativeReadChannel[V]) Length() int {
	return len(ch)
}

func (ch NativeReadChannel[V]) Capacity() int {
	return cap(ch)
}

func (ch NativeReadChannel[V]) LeftCapacity() int {
	return ch.Capacity() - ch.Length()
}

func (ch NativeReadChannel[V]) Pop() (v Value, err Value) {
	result, ok := <-ch
	if !ok {
		return Undefined, ChannelClosedPopError.ToValue()
	}
	return result.ToValue(), Undefined
}

func (ch NativeReadChannel[V]) PopCtx(ctx context.Context) (v Value, err Value) {
	select {
	case result, ok := <-ch:
		if !ok {
			return Undefined, ChannelClosedPopError.ToValue()
		}
		return result.ToValue(), Undefined
	case <-ctx.Done():
		return Undefined, ExecutionAbortedError.ToValue()
	}
}

func (ch NativeReadChannel[V]) NextValue() (Value, Value) {
	next, ok := <-ch
	if !ok {
		return Undefined, stopIterationSymbol.ToValue()
	}

	return next.ToValue(), Undefined
}

func (ch NativeReadChannel[V]) NextValueCtx(ctx context.Context) (Value, Value) {
	select {
	case next, ok := <-ch:
		if !ok {
			return Undefined, stopIterationSymbol.ToValue()
		}
		return next.ToValue(), Undefined
	case <-ctx.Done():
		return Undefined, ExecutionAbortedError.ToValue()
	}
}

func (ch NativeReadChannel[V]) Iterate() iter.Seq2[Value, Value] {
	return func(yield func(Value, Value) bool) {
		for v := range ch {
			if !yield(v.ToValue(), Undefined) {
				return
			}
		}
	}
}

func (ch NativeReadChannel[V]) Iter() NativeIterator {
	return ch
}
