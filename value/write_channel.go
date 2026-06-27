package value

import "context"

var WriteChannelClass *Class // ::Std::WriteChannel

type WriteChannel interface {
	ValueInterface
	AnyChannel
	Length() int
	Capacity() int
	LeftCapacity() int
	Push(val Value) (err Value)
	PushCtx(ctx context.Context, val Value) (err Value)
	Close() (err Value)
	IsTransformerChannel() bool
	TransformFromValue(v Value) any
	NativeChannelAny() any
}

func initWriteChannel() {
	WriteChannelClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	StdModule.AddConstantString("WriteChannel", Ref(WriteChannelClass))
	RegisterNativeClass("Std::WriteChannel", "value.WriteChannelClass")
}
