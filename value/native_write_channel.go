package value

import (
	"context"
	"fmt"
)

type NativeWriteChannel[V ValueInterface] chan<- V

var _ WriteChannel = make(NativeWriteChannel[Float])

func MakeWriteNativeChannel[V ValueInterface](size int) NativeWriteChannel[V] {
	return make(NativeWriteChannel[V], size)
}

func (ch NativeWriteChannel[V]) Copy() Reference {
	return ch
}

func (c NativeWriteChannel[V]) ToValue() Value {
	return Ref(c)
}

func (NativeWriteChannel[V]) Class() *Class {
	return WriteChannelClass
}

func (NativeWriteChannel[V]) DirectClass() *Class {
	return WriteChannelClass
}

func (NativeWriteChannel[V]) SingletonClass() *Class {
	return nil
}

func (ch NativeWriteChannel[V]) NativeChannelAny() any {
	return (chan<- V)(ch)
}

func (ch NativeWriteChannel[V]) TransformFromValue(v Value) any {
	panic("not a transformer channel")
}

func (ch NativeWriteChannel[V]) TransformToValue(v any) Value {
	panic("not a transformer channel")
}

func (NativeWriteChannel[V]) IsTransformerChannel() bool {
	return false
}

func (ch NativeWriteChannel[V]) Inspect() string {
	return fmt.Sprintf("Std::WriteChannel{&: %p, length: %d, capacity: %d}", ch, ch.Length(), ch.Capacity())
}

func (ch NativeWriteChannel[V]) Error() string {
	return ch.Inspect()
}

func (NativeWriteChannel[V]) InstanceVariables() *InstanceVariables {
	return nil
}

func (ch NativeWriteChannel[V]) Length() int {
	return len(ch)
}

func (ch NativeWriteChannel[V]) Capacity() int {
	return cap(ch)
}

func (ch NativeWriteChannel[V]) LeftCapacity() int {
	return ch.Capacity() - ch.Length()
}

func (ch NativeWriteChannel[V]) Push(val Value) (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(ChannelClosedPushError)
		}
	}()

	v, ok := Downcast[V](val)
	if !ok {
		return NewInvalidValueInChannel(ch, val.Class()).ToValue()
	}
	ch <- v
	return Undefined
}

func (ch NativeWriteChannel[V]) PushCtx(ctx context.Context, val Value) (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(ChannelClosedPushError)
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
		return ExecutionAbortedError.ToValue()
	}
}

func (ch NativeWriteChannel[V]) Close() (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(ChannelClosedCloseError)
		}
	}()

	close(ch)
	return Undefined
}
