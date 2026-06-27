package value

import (
	"context"
	"fmt"
	"iter"
)

// An adapter for a native Go read only channel to make it a valid Elk channel.
type NativeTransformerReadChannel[V any] struct {
	ch             <-chan V
	getTransformer func(V) Value
}

var _ ReadChannel = &NativeTransformerReadChannel[int]{}

func NewNativeTransformerReadChannel[V any](ch <-chan V, get func(V) Value) *NativeTransformerReadChannel[V] {
	return &NativeTransformerReadChannel[V]{
		ch:             ch,
		getTransformer: get,
	}
}

func (ch *NativeTransformerReadChannel[V]) NativeChannelAny() any {
	return ch.ch
}

func (ch *NativeTransformerReadChannel[V]) TransformToValue(v any) Value {
	return ch.getTransformer(v.(V))
}

func (ch *NativeTransformerReadChannel[V]) TransformFromValue(v Value) any {
	return v.ToInterface()
}

func (ch *NativeTransformerReadChannel[V]) IsTransformerChannel() bool {
	return true
}

func (ch *NativeTransformerReadChannel[V]) Copy() Reference {
	return ch
}

func (ch *NativeTransformerReadChannel[V]) ToValue() Value {
	return Ref(ch)
}

func (ch *NativeTransformerReadChannel[V]) Class() *Class {
	return ReadChannelClass
}

func (ch *NativeTransformerReadChannel[V]) DirectClass() *Class {
	return ReadChannelClass
}

func (ch *NativeTransformerReadChannel[V]) SingletonClass() *Class {
	return nil
}

func (ch *NativeTransformerReadChannel[V]) Inspect() string {
	return fmt.Sprintf("Std::ReadChannel{&: %p, length: %d, capacity: %d}", ch, ch.Length(), ch.Capacity())
}

func (ch *NativeTransformerReadChannel[V]) Error() string {
	return ch.Inspect()
}

func (ch *NativeTransformerReadChannel[V]) InstanceVariables() *InstanceVariables {
	return nil
}

func (ch *NativeTransformerReadChannel[V]) Length() int {
	return len(ch.ch)
}

func (ch *NativeTransformerReadChannel[V]) Capacity() int {
	return cap(ch.ch)
}

func (ch *NativeTransformerReadChannel[V]) LeftCapacity() int {
	return ch.Capacity() - ch.Length()
}

func (ch *NativeTransformerReadChannel[V]) Pop() (v Value, err Value) {
	result, ok := <-ch.ch
	if !ok {
		return Undefined, ChannelClosedPopError.ToValue()
	}
	return ch.getTransformer(result), Undefined
}

func (ch *NativeTransformerReadChannel[V]) PopCtx(ctx context.Context) (v Value, err Value) {
	select {
	case result, ok := <-ch.ch:
		if !ok {
			return Undefined, ChannelClosedPopError.ToValue()
		}
		return ch.getTransformer(result), Undefined
	case <-ctx.Done():
		return Undefined, ExecutionAbortedError.ToValue()
	}
}

func (ch *NativeTransformerReadChannel[V]) NextValue() (Value, Value) {
	next, ok := <-ch.ch
	if !ok {
		return Undefined, stopIterationSymbol.ToValue()
	}

	return ch.getTransformer(next), Undefined
}

func (ch *NativeTransformerReadChannel[V]) NextValueCtx(ctx context.Context) (Value, Value) {
	select {
	case next, ok := <-ch.ch:
		if !ok {
			return Undefined, stopIterationSymbol.ToValue()
		}
		return ch.getTransformer(next), Undefined
	case <-ctx.Done():
		return Undefined, ExecutionAbortedError.ToValue()
	}
}

func (ch *NativeTransformerReadChannel[V]) Iterate() iter.Seq2[Value, Value] {
	return func(yield func(Value, Value) bool) {
		for v := range ch.ch {
			if !yield(ch.getTransformer(v), Undefined) {
				return
			}
		}
	}
}

func (ch *NativeTransformerReadChannel[V]) Iter() NativeIterator {
	return ch
}
