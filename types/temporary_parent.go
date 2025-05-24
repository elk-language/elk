package types

// A temporary wrapper around a namespace parent
// used in the macro expansion phase before all namespaces
// are known.
type TemporaryParent struct {
	Namespace
	Child Namespace
}

func NewTemporaryParent(namespace, child Namespace) *TemporaryParent {
	return &TemporaryParent{
		Namespace: namespace,
		Child:     child,
	}
}
