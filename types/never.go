package types

// Never represents no value.
// For example a function that never returns
// might use the type `never`.
type Never struct{}

func (n Never) ToNonLiteral(env *GlobalEnvironment) Type {
	return n
}

func IsNever(t Type) bool {
	_, ok := t.(Never)
	return ok
}

func (Never) inspect() string {
	return "never"
}
