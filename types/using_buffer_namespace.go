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

func (u *UsingBufferNamespace) Copy() *UsingBufferNamespace {
	return &UsingBufferNamespace{
		Module: u.Module,
	}
}

func (u *UsingBufferNamespace) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *UsingBufferNamespace {
	newNamespace := u.Copy()

	newNamespace.methods = MethodsDeepCopyEnv(u.methods, oldEnv, newEnv)
	newNamespace.instanceVariables = TypesDeepCopyEnv(u.instanceVariables, oldEnv, newEnv)
	newNamespace.constants = ConstantsDeepCopyEnv(u.constants, oldEnv, newEnv)
	newNamespace.subtypes = ConstantsDeepCopyEnv(u.subtypes, oldEnv, newEnv)

	if u.parent != nil {
		newNamespace.parent = DeepCopyEnv(u.parent, oldEnv, newEnv).(Namespace)
	}

	return newNamespace
}
