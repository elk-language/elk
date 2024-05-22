package types

// Never represents no value. It is the bottom type.
// For example a function that never returns
// might use the type `never`.
// It is a subtype of all other types.
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
