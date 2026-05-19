package value

import (
	"context"
	"fmt"
)

var GLOBAL_ABORTER = NewAborter(context.WithCancel(context.Background()))

var AborterClass *Class                    // ::Std::Aborter
var AborterCannotBeClosedErrorClass *Class // ::Std::Aborter::CannotBeClosedError

type Aborter struct {
	ctx    context.Context
	cancel context.CancelFunc
}

var _ Reference = &Aborter{}

func NewClosedAborter() *Aborter {
	aborter := NewCancelAborter(GLOBAL_ABORTER)
	aborter.Close()
	return aborter
}

func NewTimeoutAborter(parent *Aborter, timeout TimeSpan) *Aborter {
	return NewAborter(context.WithTimeout(parent.ctx, timeout.Native()))
}

func NewDeadlineAborter(parent *Aborter, datetime *DateTime) *Aborter {
	return NewAborter(context.WithDeadline(parent.ctx, datetime.native))
}

func NewCancelAborter(parent *Aborter) *Aborter {
	return NewAborter(context.WithCancel(parent.ctx))
}

func NewAborter(ctx context.Context, cancel context.CancelFunc) *Aborter {
	return &Aborter{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (a *Aborter) Copy() Reference {
	return NewAborter(
		a.ctx,
		a.cancel,
	)
}

func (a *Aborter) ToValue() Value {
	return Ref(a)
}

func (*Aborter) Class() *Class {
	return AborterClass
}

func (*Aborter) DirectClass() *Class {
	return AborterClass
}

func (*Aborter) SingletonClass() *Class {
	return nil
}

func (a *Aborter) Inspect() string {
	return fmt.Sprintf("Std::Aborter{&: %p, is_closed: %t, is_closable: %t}", a, a.IsClosed(), a.IsClosable())
}

func (a *Aborter) Error() string {
	return a.Inspect()
}

func (*Aborter) InstanceVariables() *InstanceVariables {
	return nil
}

func (a *Aborter) Close() (err Value) {
	if a.cancel == nil {
		return NewError(AborterCannotBeClosedErrorClass, "tried to close an aborter that is not closable").ToValue()
	}

	a.cancel()
	return Undefined
}

func (a *Aborter) Closed() *NativeTransformerReadChannel[struct{}] {
	return NewNativeTransformerReadChannel(
		a.ctx.Done(),
		func(v struct{}) Value { return Nil },
	)
}

func (a *Aborter) IsClosed() bool {
	select {
	case <-a.ctx.Done():
		return true
	default:
		return false
	}
}

func (a *Aborter) IsClosable() bool {
	return a.cancel != nil
}

func (a *Aborter) Context() context.Context {
	return a.ctx
}

func (a *Aborter) CancelFunc() context.CancelFunc {
	return a.cancel
}

func initAborter() {
	AborterClass = NewClassWithOptions(ClassWithConstructor(UndefinedConstructor))
	StdModule.AddConstantString("Aborter", Ref(AborterClass))
	RegisterNativeClass("Std::Aborter", "value.AborterClass")

	AborterCannotBeClosedErrorClass = NewClassWithOptions(ClassWithSuperclass(ErrorClass))
	AborterClass.AddConstantString("CannotBeClosedError", Ref(AborterCannotBeClosedErrorClass))
	RegisterNativeClass("Std::Aborter::CannotBeClosedError", "value.AborterCannotBeClosedErrorClass")
}
