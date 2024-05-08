package types

// Void is the type that is incompatible with
// any other type.
type Void struct{}

func (Void) IsSupertypeOf(other Type) bool {
	return false
}

func (Void) Inspect() string {
	return "void"
}
