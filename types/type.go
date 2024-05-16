package types

type Type interface {
	ToNonLiteral(*GlobalEnvironment) Type
	inspect() string
}

func Inspect(typ Type) string {
	if typ == nil {
		return "void"
	}

	return typ.inspect()
}
