package value

import (
	"context"
	"fmt"
	"iter"
)

type ChannelOfValue struct {
	native chan Value
}

var _ Channel = &ChannelOfValue{}

func NewChannelOfValue(size int) *ChannelOfValue {
	return &ChannelOfValue{
		native: make(chan Value, size),
	}
}

func (ch *ChannelOfValue) ToReadChannel() ReadChannel {
	return ch.ToReadChannelOfValue()
}

func (ch *ChannelOfValue) ToReadChannelOfValue() *ReadChannelOfValue {
	return &ReadChannelOfValue{
		native: ch.native,
	}
}

func (ch *ChannelOfValue) ToWriteChannel() WriteChannel {
	return ch.ToWriteChannelOfValue()
}

func (ch *ChannelOfValue) ToWriteChannelOfValue() *WriteChannelOfValue {
	return &WriteChannelOfValue{
		native: ch.native,
	}
}

func (ch *ChannelOfValue) Copy() Reference {
	return NewChannelOfValue(ch.Length())
}

func (c *ChannelOfValue) ToValue() Value {
	return Ref(c)
}

func (*ChannelOfValue) Class() *Class {
	return ChannelClass
}

func (*ChannelOfValue) DirectClass() *Class {
	return ChannelClass
}

func (*ChannelOfValue) SingletonClass() *Class {
	return nil
}

func (ch *ChannelOfValue) Inspect() string {
	return fmt.Sprintf("Std::Channel{&: %p, length: %d, capacity: %d}", ch, ch.Length(), ch.Capacity())
}

func (ch *ChannelOfValue) ToNativeChannel() chan Value {
	return ch.native
}

func (ch *ChannelOfValue) Error() string {
	return ch.Inspect()
}

func (*ChannelOfValue) InstanceVariables() *InstanceVariables {
	return nil
}

func (ch *ChannelOfValue) NativeChannelAny() any {
	return ch.native
}

func (*ChannelOfValue) TransformToValue(v any) Value {
	panic("not a transformer channel")
}

func (*ChannelOfValue) TransformFromValue(v Value) any {
	panic("not a transformer channel")
}

func (ch *ChannelOfValue) IsTransformerChannel() bool {
	return false
}

func (ch *ChannelOfValue) Length() int {
	return len(ch.native)
}

func (ch *ChannelOfValue) Capacity() int {
	return cap(ch.native)
}

func (ch *ChannelOfValue) LeftCapacity() int {
	return ch.Capacity() - ch.Length()
}

func (ch *ChannelOfValue) Push(val Value) (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = NewError(ChannelClosedErrorClass, "cannot push values to a closed channel").ToValue()
		}
	}()

	ch.native <- val
	return Undefined
}

func (ch *ChannelOfValue) PushCtx(ctx context.Context, val Value) (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = NewError(ChannelClosedErrorClass, "cannot push values to a closed channel").ToValue()
		}
	}()

	select {
	case ch.native <- val:
		return Undefined
	case <-ctx.Done():
		return NewExecutionAbortedError().ToValue()
	}
}

func (ch *ChannelOfValue) Pop() (result Value, err Value) {
	result, ok := <-ch.native
	if !ok {
		return Undefined, NewError(ChannelClosedErrorClass, "cannot pop values from a closed channel").ToValue()
	}

	return result, Undefined
}

func (ch *ChannelOfValue) PopCtx(ctx context.Context) (result Value, err Value) {
	select {
	case result, ok := <-ch.native:
		if !ok {
			return Undefined, NewError(ChannelClosedErrorClass, "cannot pop values from a closed channel").ToValue()
		}
		return result, Undefined
	case <-ctx.Done():
		return Undefined, NewExecutionAbortedError().ToValue()
	}
}

func (ch *ChannelOfValue) NextValue() (Value, Value) {
	next, ok := <-ch.native
	if !ok {
		return Undefined, stopIterationSymbol.ToValue()
	}

	return next, Undefined
}

func (ch *ChannelOfValue) NextValueCtx(ctx context.Context) (Value, Value) {
	select {
	case next, ok := <-ch.native:
		if !ok {
			return Undefined, stopIterationSymbol.ToValue()
		}
		return next, Undefined
	case <-ctx.Done():
		return Undefined, NewExecutionAbortedError().ToValue()
	}
}

func (ch *ChannelOfValue) Iterate() iter.Seq2[Value, Value] {
	return func(yield func(Value, Value) bool) {
		for v := range ch.native {
			if !yield(v, Undefined) {
				return
			}
		}
	}
}

func (ch *ChannelOfValue) Iter() NativeIterator {
	return ch
}

func (ch *ChannelOfValue) Close() (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(NewError(ChannelClosedErrorClass, "cannot close a closed channel"))
		}
	}()

	close(ch.native)
	return Undefined
}
