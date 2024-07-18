package types

// Nothing represents no value. It is the bottom type.
// It is a subtype of all other types.
//
// It is similar to `never` but is used as a return type
// for calls to methods that do not exist to avoid cascading errors.
type Nothing struct{}

func (n Nothing) ToNonLiteral(env *GlobalEnvironment) Type {
	return n
}

func (Nothing) IsLiteral() bool {
	return false
}

func IsNothing(t Type) bool {
	_, ok := t.(Nothing)
	return ok
}

func (Nothing) inspect() string {
	return "nothing"
}
