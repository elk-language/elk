package value

import (
	"fmt"
	"sync"
)

var OnceClass *Class // ::Std::Sync::Once

// Wraps a Go Once.
type Once struct {
	Native sync.Once
}

func NewOnce() *Once {
	return &Once{}
}

func OnceConstructor(class *Class) Value {
	return Ref(NewOnce())
}

func (o *Once) Copy() Reference {
	return NewOnce()
}

func (*Once) Class() *Class {
	return OnceClass
}

func (*Once) DirectClass() *Class {
	return OnceClass
}

func (*Once) SingletonClass() *Class {
	return nil
}

func (o *Once) Inspect() string {
	return fmt.Sprintf("Std::Sync::Once{&: %p}", o)
}

func (o *Once) Error() string {
	return o.Inspect()
}

func (*Once) InstanceVariables() SymbolMap {
	return nil
}

func initOnce() {
	OnceClass = NewClassWithOptions(
		ClassWithConstructor(OnceConstructor),
	)
	SyncModule.AddConstantString("Once", Ref(OnceClass))
}
