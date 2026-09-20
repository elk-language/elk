package types

import "unsafe"

type UsingBufferNamespace struct {
	Module
}

func (u *UsingBufferNamespace) HashUint64() uint64 {
	return uint64(uintptr(unsafe.Pointer(u)))
}

func (u *UsingBufferNamespace) EqualAny(other any) bool {
	return u == other
}

func NewUsingBufferNamespace() *UsingBufferNamespace {
	return &UsingBufferNamespace{
		Module: Module{
			NamespaceBase: MakeNamespaceBase("", "<using buffer namespace>"),
		},
	}
}

func (u *UsingBufferNamespace) Copy() *UsingBufferNamespace {
	result := &UsingBufferNamespace{
		Module: u.Module,
	}
	result.id = ZERO_ID
	return result
}
