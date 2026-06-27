package value

import "context"

var ReadChannelClass *Class // ::Std::ReadChannel

type ReadChannel interface {
	ValueInterface
	NativeIterable
	NativeIterator
	AnyChannel
	Length() int
	Capacity() int
	LeftCapacity() int
	Pop() (val Value, err Value)
	PopCtx(ctx context.Context) (val Value, err Value)
	NextValueCtx(ctx context.Context) (val Value, err Value)
	IsTransformerChannel() bool
	TransformToValue(v any) Value
	NativeChannelAny() any
}

func initReadChannel() {
	ReadChannelClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	StdModule.AddConstantString("ReadChannel", Ref(ReadChannelClass))
	RegisterNativeClass("Std::ReadChannel", "value.ReadChannelClass")
}
