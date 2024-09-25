package types

type UsingBufferNamespace struct {
	Module
}

func NewUsingBufferNamespace() *UsingBufferNamespace {
	return &UsingBufferNamespace{
		Module: Module{
			NamespaceBase: MakeNamespaceBase("", "<using buffer namespace>"),
		},
	}
}
