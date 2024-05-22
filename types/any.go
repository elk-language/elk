package types

// All types are subtypes of any.
// Any is not a subtype of anything other than itself.
// It is the top type.
type Any struct{}

func (n Any) ToNonLiteral(env *GlobalEnvironment) Type {
	return n
}

func IsAny(t Type) bool {
	_, ok := t.(Any)
	return ok
}

func (Any) inspect() string {
	return "any"
}
