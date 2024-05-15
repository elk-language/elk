package types

// Void is the type that is incompatible with
// any other type.
type Void struct{}

func (Void) IsSubtypeOf(other Type, env *GlobalEnvironment) bool {
	return false
}

func (v Void) ToNonLiteral(env *GlobalEnvironment) Type {
	return v
}

func (Void) Inspect() string {
	return "void"
}
