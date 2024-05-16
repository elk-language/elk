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

func IsVoid(t Type) bool {
	_, ok := t.(Void)
	return ok
}

func (Void) inspect() string {
	return "void"
}
