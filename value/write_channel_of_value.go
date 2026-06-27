package value

import (
	"context"
	"fmt"
)

type WriteChannelOfValue struct {
	native chan<- Value
}

var _ WriteChannel = &WriteChannelOfValue{}

func NewWriteChannelOfValue(size int) *WriteChannelOfValue {
	return &WriteChannelOfValue{
		native: make(chan Value, size),
	}
}

func (ch *WriteChannelOfValue) Copy() Reference {
	return NewWriteChannelOfValue(ch.Length())
}

func (c *WriteChannelOfValue) ToValue() Value {
	return Ref(c)
}

func (*WriteChannelOfValue) Class() *Class {
	return WriteChannelClass
}

func (*WriteChannelOfValue) DirectClass() *Class {
	return WriteChannelClass
}

func (*WriteChannelOfValue) SingletonClass() *Class {
	return nil
}

func (ch *WriteChannelOfValue) Inspect() string {
	return fmt.Sprintf("Std::WriteChannel{&: %p, length: %d, capacity: %d}", ch, ch.Length(), ch.Capacity())
}

func (ch *WriteChannelOfValue) ToNativeChannel() chan<- Value {
	return ch.native
}

func (ch *WriteChannelOfValue) Error() string {
	return ch.Inspect()
}

func (*WriteChannelOfValue) InstanceVariables() *InstanceVariables {
	return nil
}

func (ch *WriteChannelOfValue) NativeChannelAny() any {
	return ch.native
}

func (*WriteChannelOfValue) TransformFromValue(v Value) any {
	panic("not a transformer channel")
}

func (*WriteChannelOfValue) TransformToValue(v any) Value {
	panic("not a transformer channel")
}

func (ch *WriteChannelOfValue) IsTransformerChannel() bool {
	return false
}

func (ch *WriteChannelOfValue) Length() int {
	return len(ch.native)
}

func (ch *WriteChannelOfValue) Capacity() int {
	return cap(ch.native)
}

func (ch *WriteChannelOfValue) LeftCapacity() int {
	return ch.Capacity() - ch.Length()
}

func (ch *WriteChannelOfValue) Push(val Value) (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(ChannelClosedPushError)
		}
	}()

	ch.native <- val
	return Undefined
}

func (ch *WriteChannelOfValue) PushCtx(ctx context.Context, val Value) (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(ChannelClosedPushError)
		}
	}()

	select {
	case ch.native <- val:
		return Undefined
	case <-ctx.Done():
		return ExecutionAbortedError.ToValue()
	}
}

func (ch *WriteChannelOfValue) Close() (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(ChannelClosedCloseError)
		}
	}()

	close(ch.native)
	return Undefined
}
