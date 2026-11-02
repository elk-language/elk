package value

import "context"

var ChannelClass *Class            // ::Std::Channel
var ChannelClosedErrorClass *Class // ::Std::Channel::ClosedError

var ChannelClosedPopError *Object
var ChannelClosedPushError *Object
var ChannelClosedCloseError *Object

func NormaliseSentValueToChannel(channel AnyChannel, val Value) any {
	if !channel.IsTransformerChannel() {
		return val
	}

	return channel.TransformFromValueToAny(val)
}

type AnyChannel interface {
	ValueInterface
	Length() int
	Capacity() int
	LeftCapacity() int
	IsTransformerChannel() bool
	TransformFromValueToAny(v Value) any
	TransformAnyToValue(v any) Value
	NativeChannelAny() any
}

type Channel interface {
	ValueInterface
	NativeIterable
	NativeIterator
	AnyChannel
	Push(val Value) (err Value)
	PushCtx(ctx context.Context, val Value) (err Value)
	Pop() (val Value, err Value)
	PopCtx(ctx context.Context) (val Value, err Value)
	NextValueCtx(ctx context.Context) (val Value, err Value)
	Close() (err Value)
	ToReadChannel() ReadChannel
	ToWriteChannel() WriteChannel
}

func initChannel() {
	ChannelClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	StdModule.AddConstantString("Channel", Ref(ChannelClass))
	RegisterNativeClass("Std::Channel", "value.ChannelClass")

	ChannelClosedErrorClass = NewClassWithOptions(ClassWithSuperclass(ErrorClass))
	ChannelClass.AddConstantString("ClosedError", Ref(ChannelClosedErrorClass))
	RegisterNativeClass("Std::Channel::ClosedError", "value.ChannelClosedErrorClass")

	ChannelClosedPopError = NewError(ChannelClosedErrorClass, "cannot pop values from a closed channel")
	ChannelClosedPushError = NewError(ChannelClosedErrorClass, "cannot push values to a closed channel")
	ChannelClosedCloseError = NewError(ChannelClosedErrorClass, "cannot close a closed channel")
}
