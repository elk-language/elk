package value

import (
	"fmt"
	"sync"
)

var MutexClass *Class              // ::Std::Sync::Mutex
var MutexUnlockedErrorClass *Class // ::Std::Sync::Mutex::UnlockedError

// Wraps a Go mutex.
type Mutex struct {
	Native sync.Mutex
}

func NewMutex() *Mutex {
	return &Mutex{}
}

func MutexConstructor(class *Class) Value {
	return Ref(NewMutex())
}

func (m *Mutex) Copy() Reference {
	return NewMutex()
}

func (*Mutex) Class() *Class {
	return MutexClass
}

func (*Mutex) DirectClass() *Class {
	return MutexClass
}

func (*Mutex) SingletonClass() *Class {
	return nil
}

func (m *Mutex) Inspect() string {
	return fmt.Sprintf("Std::Sync::Mutex{&: %p}", m)
}

func (m *Mutex) Error() string {
	return m.Inspect()
}

func (*Mutex) InstanceVariables() SymbolMap {
	return nil
}

func (m *Mutex) Lock() {
	m.Native.Lock()
}

func (m *Mutex) Unlock() (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(NewError(MutexUnlockedErrorClass, "cannot unlock an unlocked mutex"))
		}
	}()

	m.Native.Unlock()
	return Undefined
}

func initMutex() {
	MutexClass = NewClassWithOptions(
		ClassWithConstructor(MutexConstructor),
	)
	SyncModule.AddConstantString("Mutex", Ref(MutexClass))

	MutexUnlockedErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	MutexClass.AddConstantString("UnlockedError", Ref(MutexUnlockedErrorClass))
}
