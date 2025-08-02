package value

import (
	"fmt"
	"sync"
)

var RWMutexClass *Class              // ::Std::Sync::RWMutex
var RWMutexUnlockedErrorClass *Class // ::Std::Sync::RWMutex::UnlockedError

// Wraps a Go RWMutex.
type RWMutex struct {
	Native sync.RWMutex
}

func NewRWMutex() *RWMutex {
	return &RWMutex{}
}

func RWMutexConstructor(class *Class) Value {
	return Ref(NewRWMutex())
}

func (m *RWMutex) Copy() Reference {
	return NewRWMutex()
}

func (*RWMutex) Class() *Class {
	return RWMutexClass
}

func (*RWMutex) DirectClass() *Class {
	return RWMutexClass
}

func (*RWMutex) SingletonClass() *Class {
	return nil
}

func (m *RWMutex) Inspect() string {
	return fmt.Sprintf("Std::Sync::RWMutex{&: %p}", m)
}

func (m *RWMutex) Error() string {
	return m.Inspect()
}

func (*RWMutex) InstanceVariables() *InstanceVariables {
	return nil
}

func (m *RWMutex) Lock() {
	m.Native.Lock()
}

func (m *RWMutex) ReadLock() {
	m.Native.RLock()
}

func (m *RWMutex) Unlock() (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(NewError(RWMutexUnlockedErrorClass, "a rwmutex that is unlocked for writing cannot be unlocked for writing"))
		}
	}()

	m.Native.Unlock()
	return Undefined
}

func (m *RWMutex) ReadUnlock() (err Value) {
	defer func() {
		if r := recover(); r != nil {
			err = Ref(NewError(RWMutexUnlockedErrorClass, "a rwmutex that is unlocked for reading cannot be unlocked for reading"))
		}
	}()

	m.Native.RUnlock()
	return Undefined
}

func initRWMutex() {
	RWMutexClass = NewClassWithOptions(
		ClassWithConstructor(RWMutexConstructor),
	)
	SyncModule.AddConstantString("RWMutex", Ref(RWMutexClass))

	RWMutexUnlockedErrorClass = NewClassWithOptions(ClassWithParent(ErrorClass))
	RWMutexClass.AddConstantString("UnlockedError", Ref(RWMutexUnlockedErrorClass))
}
