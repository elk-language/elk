package types

// Represents no value.
// It is used to mark constants of pure types, that do not have a runtime value.
type NoValue struct{}

func (n NoValue) ToNonLiteral(env *GlobalEnvironment) Type {
	return n
}

func (NoValue) IsLiteral() bool {
	return false
}

func IsNoValue(t Type) bool {
	_, ok := t.(NoValue)
	return ok
}

func (NoValue) inspect() string {
	return "NoValue"
}
