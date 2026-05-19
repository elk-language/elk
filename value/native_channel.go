package value

import (
	"context"
	"fmt"
	"iter"
)

type NativeChannel[V ValueInterface] chan V

var _ Channel = make(NativeChannel[Float])

func MakeNativeChannel[V ValueInterface](size int) NativeChannel[V] {
	return make(NativeChannel[V], size)
}

func (ch NativeChannel[V]) ToWriteChannel() WriteChannel {
	return ch.ToNativeWriteChannel()
}

func (ch NativeChannel[V]) ToNativeWriteChannel() NativeWriteChannel[V] {
	return (chan V)(ch)
}

func (ch NativeChannel[V]) ToReadChannel() ReadChannel {
	return ch.ToNativeReadChannel()
}

func (ch NativeChannel[V]) ToNativeReadChannel() NativeReadChannel[V] {
	return (chan V)(ch)
}

func (ch NativeChannel[V]) Copy() Reference {
	return ch
}

func (c NativeChannel[V]) ToValue() Value {
	return Ref(c)
}

func (NativeChannel[V]) Class() *Class {
	return ChannelClass
}

func (NativeChannel[V]) DirectClass() *Class {
	return ChannelClass
}

func (NativeChannel[V]) SingletonClass() *Class {
	return nil
}

func (ch NativeChannel[V]) NativeChannelAny() any {
	return (chan V)(ch)
}

func (ch NativeChannel[V]) TransformToValue(v any) Value {
	panic("not a transformer channel")
}

func (ch NativeChannel[V]) TransformFromValue(v Value) any {
	panic("not a transformer channel")
}

func (NativeChannel[V]) IsTransformerChannel() bool {
	return false
}

func (ch NativeChannel[V]) Inspect() string {
	return fmt.Sprintf("Std::Channel{&: %p, length: %d, capacity: %d}", ch, ch.Length(), ch.Capacity())
}

func (ch NativeChannel[V]) Error() string {
	return ch.Inspect()
}

func (NativeChannel[V]) InstanceVariables() *InstanceVariables {
	return nil
}

func (ch NativeChannel[V]) Length() int {
	return len(ch)
}

func (ch NativeChannel[V]) Capacity() int {
	return cap(ch)
}

func (ch NativeChannel[V]) LeftCapacity() int {
	return ch.Capacity() - ch.Length()
}

func (ch NativeChannel[V]) Push(val Value) (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(NewError(ChannelClosedErrorClass, "cannot push values to a closed channel"))
		}
	}()

	v, ok := Downcast[V](val)
	if !ok {
		return NewInvalidValueInChannel(ch, val.Class()).ToValue()
	}
	ch <- v
	return Undefined
}

func (ch NativeChannel[V]) PushCtx(ctx context.Context, val Value) (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(NewError(ChannelClosedErrorClass, "cannot push values to a closed channel"))
		}
	}()

	v, ok := Downcast[V](val)
	if !ok {
		return NewInvalidValueInChannel(ch, val.Class()).ToValue()
	}
	select {
	case ch <- v:
		return Undefined
	case <-ctx.Done():
		return NewExecutionAbortedError().ToValue()
	}
}

func (ch NativeChannel[V]) Pop() (v Value, err Value) {
	result, ok := <-ch
	if !ok {
		return Undefined, NewError(ChannelClosedErrorClass, "cannot pop values from a closed channel").ToValue()
	}

	return result.ToValue(), Undefined
}

func (ch NativeChannel[V]) PopCtx(ctx context.Context) (v Value, err Value) {
	select {
	case result, ok := <-ch:
		if !ok {
			return Undefined, NewError(ChannelClosedErrorClass, "cannot pop values from a closed channel").ToValue()
		}
		return result.ToValue(), Undefined
	case <-ctx.Done():
		return Undefined, NewExecutionAbortedError().ToValue()
	}
}

func (ch NativeChannel[V]) NextValue() (Value, Value) {
	next, ok := <-ch
	if !ok {
		return Undefined, stopIterationSymbol.ToValue()
	}

	return next.ToValue(), Undefined
}

func (ch NativeChannel[V]) NextValueCtx(ctx context.Context) (Value, Value) {
	select {
	case next, ok := <-ch:
		if !ok {
			return Undefined, stopIterationSymbol.ToValue()
		}
		return next.ToValue(), Undefined
	case <-ctx.Done():
		return Undefined, NewExecutionAbortedError().ToValue()
	}
}

func (ch NativeChannel[V]) Iterate() iter.Seq2[Value, Value] {
	return func(yield func(Value, Value) bool) {
		for v := range ch {
			if !yield(v.ToValue(), Undefined) {
				return
			}
		}
	}
}

func (ch NativeChannel[V]) Iter() NativeIterator {
	return ch
}

func (ch NativeChannel[V]) Close() (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(NewError(ChannelClosedErrorClass, "cannot close a closed channel"))
		}
	}()

	close(ch)
	return Undefined
}
