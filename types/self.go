package types

type Self struct{}

func (s Self) ToNonLiteral(env *GlobalEnvironment) Type {
	return s
}

func (Self) IsLiteral() bool {
	return false
}

func IsSelf(t Type) bool {
	_, ok := t.(Self)
	return ok
}

func (Self) inspect() string {
	return "self"
}
