package types

// A temporary wrapper around a namespace parent
// used in the macro expansion phase before all namespaces
// are known.
type TemporaryParent struct {
	Namespace
}

func NewTemporaryParent(namespace Namespace) *TemporaryParent {
	return &TemporaryParent{
		Namespace: namespace,
	}
}

func (t *TemporaryParent) DeepCopyEnv(oldEnv, newEnv *GlobalEnvironment) *TemporaryParent {
	newTemporaryParent := &TemporaryParent{
		Namespace: DeepCopyEnv(t.Namespace, oldEnv, newEnv).(Namespace),
	}

	return newTemporaryParent
}
