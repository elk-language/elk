package types

type Type interface {
	IsSupertypeOf(Type) bool
	Inspect() string
}

func Inspect(typ Type) string {
	if typ == nil {
		return "void"
	}

	return typ.Inspect()
}
