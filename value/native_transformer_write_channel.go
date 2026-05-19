package value

import (
	"context"
	"fmt"
)

// An adapter for a native Go write only channel to make it a valid Elk channel.
type NativeTransformerWriteChannel[V any] struct {
	ch             chan<- V
	setTransformer func(Value) V
}

var _ WriteChannel = &NativeTransformerWriteChannel[int]{}

func NewNativeTransformerWriteChannel[V any](ch chan<- V, set func(Value) V) *NativeTransformerWriteChannel[V] {
	return &NativeTransformerWriteChannel[V]{
		ch:             ch,
		setTransformer: set,
	}
}

func (ch *NativeTransformerWriteChannel[V]) NativeChannelAny() any {
	return ch.ch
}

func (ch *NativeTransformerWriteChannel[V]) TransformFromValue(v Value) any {
	return ch.setTransformer(v)
}

func (ch *NativeTransformerWriteChannel[V]) TransformToValue(v any) Value {
	return v.(ValueInterface).ToValue()
}

func (ch *NativeTransformerWriteChannel[V]) IsTransformerChannel() bool {
	return true
}

func (ch *NativeTransformerWriteChannel[V]) Copy() Reference {
	return ch
}

func (ch *NativeTransformerWriteChannel[V]) ToValue() Value {
	return Ref(ch)
}

func (ch *NativeTransformerWriteChannel[V]) Class() *Class {
	return WriteChannelClass
}

func (ch *NativeTransformerWriteChannel[V]) DirectClass() *Class {
	return WriteChannelClass
}

func (ch *NativeTransformerWriteChannel[V]) SingletonClass() *Class {
	return nil
}

func (ch *NativeTransformerWriteChannel[V]) Inspect() string {
	return fmt.Sprintf("Std::WriteChannel{&: %p, length: %d, capacity: %d}", ch, ch.Length(), ch.Capacity())
}

func (ch *NativeTransformerWriteChannel[V]) Error() string {
	return ch.Inspect()
}

func (ch *NativeTransformerWriteChannel[V]) InstanceVariables() *InstanceVariables {
	return nil
}

func (ch *NativeTransformerWriteChannel[V]) Length() int {
	return len(ch.ch)
}

func (ch *NativeTransformerWriteChannel[V]) Capacity() int {
	return cap(ch.ch)
}

func (ch *NativeTransformerWriteChannel[V]) LeftCapacity() int {
	return ch.Capacity() - ch.Length()
}

func (ch *NativeTransformerWriteChannel[V]) Push(val Value) (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(NewError(ChannelClosedErrorClass, "cannot push values to a closed channel"))
		}
	}()

	ch.ch <- ch.setTransformer(val)
	return Undefined
}

func (ch *NativeTransformerWriteChannel[V]) PushCtx(ctx context.Context, val Value) (err Value) {
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

func (ch *NativeTransformerWriteChannel[V]) Close() (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(NewError(ChannelClosedErrorClass, "cannot close a closed channel"))
		}
	}()

	close(ch.ch)
	return Undefined
}
