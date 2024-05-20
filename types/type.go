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

func GetMethod(typ Type, name string, env *GlobalEnvironment) *Method {
	typ = typ.ToNonLiteral(env)

	switch t := typ.(type) {
	case *Class:
		return t.MethodString(name)
	case *Module:
		return t.MethodString(name)
	}

	return nil
}
