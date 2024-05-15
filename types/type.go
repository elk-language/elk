package types

type Type interface {
	IsSubtypeOf(Type, *GlobalEnvironment) bool
	ToNonLiteral(*GlobalEnvironment) Type
	Inspect() string
}

func Inspect(typ Type) string {
	if typ == nil {
		return "void"
	}

	return typ.Inspect()
}
