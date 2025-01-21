package value

import (
	"fmt"
)

var ChannelClass *Class            // ::Std::Channel
var ChannelClosedErrorClass *Class // ::Std::Channel::ClosedError

type Channel struct {
	Native chan Value
}

func NewChannel(size int) *Channel {
	return &Channel{
		Native: make(chan Value, size),
	}
}

func ChannelConstructor(class *Class) Value {
	return Undefined
}

func (ch *Channel) Copy() Reference {
	return NewChannel(ch.Length())
}

func (*Channel) Class() *Class {
	return ChannelClass
}

func (*Channel) DirectClass() *Class {
	return ChannelClass
}

func (*Channel) SingletonClass() *Class {
	return nil
}

func (ch *Channel) Inspect() string {
	return fmt.Sprintf("Std::Channel{&: %p, length: %d, capacity: %d}", ch, ch.Length(), ch.Capacity())
}

func (ch *Channel) Error() string {
	return ch.Inspect()
}

func (*Channel) InstanceVariables() SymbolMap {
	return nil
}

func (ch *Channel) Length() int {
	return len(ch.Native)
}

func (ch *Channel) Capacity() int {
	return cap(ch.Native)
}

func (ch *Channel) LeftCapacity() int {
	return ch.Capacity() - ch.Length()
}

func (ch *Channel) Push(val Value) (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(NewError(ChannelClosedErrorClass, "cannot push values to a closed channel"))
		}
	}()

	ch.Native <- val
	return Undefined
}

func (ch *Channel) Pop() (Value, bool) {
	result, ok := <-ch.Native
	return result, ok
}

func (ch *Channel) Next() (Value, Value) {
	next, ok := <-ch.Native
	if !ok {
		return Undefined, stopIterationSymbol.ToValue()
	}

	return next, Undefined
}

func (ch *Channel) Close() (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(NewError(ChannelClosedErrorClass, "cannot close a closed channel"))
		}
	}()

	close(ch.Native)
	return Undefined
}

func initChannel() {
	ChannelClass = NewClassWithOptions(ClassWithConstructor(ChannelConstructor))
	StdModule.AddConstantString("Channel", Ref(ChannelClass))

	ChannelClosedErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	ChannelClass.AddConstantString("ClosedError", Ref(ChannelClosedErrorClass))
}
