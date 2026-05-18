package value

import (
	"context"
	"fmt"
	"iter"
)

// An adapter for a native Go channel to make it a valid Elk channel.
type NativeTransformerChannel[V any] struct {
	ch             chan V
	getTransformer func(V) Value
	setTransformer func(Value) V
}

var _ Channel = &NativeTransformerChannel[int]{}

func NewNativeTransformerChannel[V any](ch chan V, get func(V) Value, set func(Value) V) *NativeTransformerChannel[V] {
	return &NativeTransformerChannel[V]{
		ch:             ch,
		getTransformer: get,
		setTransformer: set,
	}
}

func (ch NativeTransformerChannel[V]) NativeChannelAny() any {
	return ch.ch
}

func (ch NativeTransformerChannel[V]) TransformToValue(v any) Value {
	return ch.getTransformer(v.(V))
}

func (ch NativeTransformerChannel[V]) TransformFromValue(v Value) any {
	return ch.setTransformer(v)
}

func (ch NativeTransformerChannel[V]) IsTransformerChannel() bool {
	return true
}

func (ch NativeTransformerChannel[V]) Copy() Reference {
	return ch
}

func (ch NativeTransformerChannel[V]) ToValue() Value {
	return Ref(ch)
}

func (ch NativeTransformerChannel[V]) Class() *Class {
	return ChannelClass
}

func (ch NativeTransformerChannel[V]) DirectClass() *Class {
	return ChannelClass
}

func (ch NativeTransformerChannel[V]) SingletonClass() *Class {
	return nil
}

func (ch NativeTransformerChannel[V]) Inspect() string {
	return fmt.Sprintf("Std::Channel{&: %p, length: %d, capacity: %d}", ch, ch.Length(), ch.Capacity())
}

func (ch NativeTransformerChannel[V]) Error() string {
	return ch.Inspect()
}

func (ch NativeTransformerChannel[V]) InstanceVariables() *InstanceVariables {
	return nil
}

func (ch NativeTransformerChannel[V]) Length() int {
	return len(ch.ch)
}

func (ch NativeTransformerChannel[V]) Capacity() int {
	return cap(ch.ch)
}

func (ch NativeTransformerChannel[V]) LeftCapacity() int {
	return ch.Capacity() - ch.Length()
}

func (ch NativeTransformerChannel[V]) Push(val Value) (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(NewError(ChannelClosedErrorClass, "cannot push values to a closed channel"))
		}
	}()

	ch.ch <- ch.setTransformer(val)
	return Undefined
}

func (ch NativeTransformerChannel[V]) PushCtx(ctx context.Context, val Value) (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(NewError(ChannelClosedErrorClass, "cannot push values to a closed channel"))
		}
	}()

	select {
	case ch.ch <- ch.setTransformer(val):
		return Undefined
	case <-ctx.Done():
		return NewExecutionAbortedError().ToValue()
	}
}

func (ch NativeTransformerChannel[V]) Pop() (Value, bool) {
	result, ok := <-ch.ch
	return ch.getTransformer(result), ok
}

func (ch NativeTransformerChannel[V]) PopCtx(ctx context.Context) (v Value, ok bool, err Value) {
	select {
	case result, ok := <-ch.ch:
		return ch.getTransformer(result), ok, Undefined
	case <-ctx.Done():
		return Undefined, false, NewExecutionAbortedError().ToValue()
	}
}

func (ch NativeTransformerChannel[V]) NextValue() (Value, Value) {
	next, ok := <-ch.ch
	if !ok {
		return Undefined, stopIterationSymbol.ToValue()
	}

	return ch.getTransformer(next), Undefined
}

func (ch NativeTransformerChannel[V]) NextValueCtx(ctx context.Context) (Value, Value) {
	select {
	case next, ok := <-ch.ch:
		if !ok {
			return Undefined, stopIterationSymbol.ToValue()
		}
		return ch.getTransformer(next), Undefined
	case <-ctx.Done():
		return Undefined, NewExecutionAbortedError().ToValue()
	}
}

func (ch NativeTransformerChannel[V]) Iterate() iter.Seq2[Value, Value] {
	return func(yield func(Value, Value) bool) {
		for v := range ch.ch {
			if !yield(ch.getTransformer(v), Undefined) {
				return
			}
		}
	}
}

func (ch NativeTransformerChannel[V]) Iter() NativeIterator {
	return ch
}

func (ch NativeTransformerChannel[V]) Close() (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(NewError(ChannelClosedErrorClass, "cannot close a closed channel"))
		}
	}()

	close(ch.ch)
	return Undefined
}
