package value

import (
	"context"
	"fmt"
	"iter"
)

type ReadChannelOfValue struct {
	native <-chan Value
}

var _ ReadChannel = &ReadChannelOfValue{}

func NewReadChannelOfValue(size int) *ReadChannelOfValue {
	return &ReadChannelOfValue{
		native: make(chan Value, size),
	}
}

func (ch *ReadChannelOfValue) Copy() Reference {
	return NewReadChannelOfValue(ch.Length())
}

func (c *ReadChannelOfValue) ToValue() Value {
	return Ref(c)
}

func (*ReadChannelOfValue) Class() *Class {
	return ReadChannelClass
}

func (*ReadChannelOfValue) DirectClass() *Class {
	return ReadChannelClass
}

func (*ReadChannelOfValue) SingletonClass() *Class {
	return nil
}

func (ch *ReadChannelOfValue) Inspect() string {
	return fmt.Sprintf("Std::ReadChannel{&: %p, length: %d, capacity: %d}", ch, ch.Length(), ch.Capacity())
}

func (ch *ReadChannelOfValue) ToNativeChannel() <-chan Value {
	return ch.native
}

func (ch *ReadChannelOfValue) Error() string {
	return ch.Inspect()
}

func (*ReadChannelOfValue) InstanceVariables() *InstanceVariables {
	return nil
}

func (ch *ReadChannelOfValue) NativeChannelAny() any {
	return ch.native
}

func (*ReadChannelOfValue) TransformToValue(v any) Value {
	panic("not a transformer channel")
}

func (*ReadChannelOfValue) TransformFromValue(v Value) any {
	panic("not a transformer channel")
}

func (ch *ReadChannelOfValue) IsTransformerChannel() bool {
	return false
}

func (ch *ReadChannelOfValue) Length() int {
	return len(ch.native)
}

func (ch *ReadChannelOfValue) Capacity() int {
	return cap(ch.native)
}

func (ch *ReadChannelOfValue) LeftCapacity() int {
	return ch.Capacity() - ch.Length()
}

func (ch *ReadChannelOfValue) Pop() (v Value, err Value) {
	result, ok := <-ch.native
	if !ok {
		return Undefined, NewError(ChannelClosedErrorClass, "cannot pop values from a closed channel").ToValue()
	}
	return result, Undefined
}

func (ch *ReadChannelOfValue) PopCtx(ctx context.Context) (v Value, err Value) {
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

func (ch *ReadChannelOfValue) NextValue() (Value, Value) {
	next, ok := <-ch.native
	if !ok {
		return Undefined, stopIterationSymbol.ToValue()
	}

	return next, Undefined
}

func (ch *ReadChannelOfValue) NextValueCtx(ctx context.Context) (Value, Value) {
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

func (ch *ReadChannelOfValue) Iterate() iter.Seq2[Value, Value] {
	return func(yield func(Value, Value) bool) {
		for v := range ch.native {
			if !yield(v, Undefined) {
				return
			}
		}
	}
}

func (ch *ReadChannelOfValue) Iter() NativeIterator {
	return ch
}
