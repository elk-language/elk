package value

import "context"

var ChannelClass *Class            // ::Std::Channel
var ChannelClosedErrorClass *Class // ::Std::Channel::ClosedError

type Channel interface {
	ValueInterface
	NativeIterable
	NativeIterator
	Length() int
	Capacity() int
	LeftCapacity() int
	Push(val Value) (err Value)
	PushCtx(ctx context.Context, val Value) (err Value)
	Pop() (Value, bool)
	PopCtx(ctx context.Context) (val Value, ok bool, err Value)
	NextValueCtx(ctx context.Context) (val Value, err Value)
	Close() (err Value)
	IsTransformerChannel() bool
	TransformFromValue(v Value) any
	TransformToValue(v any) Value
	NativeChannelAny() any
}

func initChannel() {
	ChannelClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	StdModule.AddConstantString("Channel", Ref(ChannelClass))
	RegisterNativeClass("Std::Channel", "value.ChannelClass")

	ChannelClosedErrorClass = NewClassWithOptions(ClassWithSuperclass(ErrorClass))
	ChannelClass.AddConstantString("ClosedError", Ref(ChannelClosedErrorClass))
	RegisterNativeClass("Std::Channel::ClosedError", "value.ChannelClosedErrorClass")
}
