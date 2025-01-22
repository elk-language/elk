package value

import (
	"fmt"
)

var ROMutexClass *Class // ::Std::Sync::ROMutex

// Wraps a RWMutex.
type ROMutex struct {
	RWMutex *RWMutex
}

func NewROMutex(rwmutex *RWMutex) *ROMutex {
	return &ROMutex{
		RWMutex: rwmutex,
	}
}

func ROMutexConstructor(class *Class) Value {
	return Ref(NewROMutex(nil))
}

func (m *ROMutex) Copy() Reference {
	return NewROMutex(m.RWMutex)
}

func (*ROMutex) Class() *Class {
	return ROMutexClass
}

func (*ROMutex) DirectClass() *Class {
	return ROMutexClass
}

func (*ROMutex) SingletonClass() *Class {
	return nil
}

func (m *ROMutex) Inspect() string {
	return fmt.Sprintf("Std::Sync::ROMutex{&: %p, rwmutex: %s}", m, m.RWMutex.Inspect())
}

func (m *ROMutex) Error() string {
	return m.Inspect()
}

func (*ROMutex) InstanceVariables() SymbolMap {
	return nil
}

func (m *ROMutex) Lock() {
	m.RWMutex.ReadLock()
}

func (m *ROMutex) Unlock() (err Value) {
	return m.RWMutex.ReadUnlock()
}

func initROMutex() {
	ROMutexClass = NewClassWithOptions(
		ClassWithConstructor(ROMutexConstructor),
	)
	SyncModule.AddConstantString("ROMutex", Ref(ROMutexClass))
}
